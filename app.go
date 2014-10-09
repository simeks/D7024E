package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"time"
	"encoding/json"
	"strconv"
)

type App struct {
	node    *Node
	transport Transport
	net Net

	keyValue map[string]string
}

func (this *App) init(bindAddr string) {
	this.keyValue = make(map[string]string)

	this.node = makeDHTNode(nil, bindAddr)
	this.transport.init(bindAddr)
	this.net = Net{this}

	// call stabilize and fixFingers periodically
	go func() {
		c := time.Tick(3 * time.Second)
		for {
			select {
			case <-c:
				this.stabilize()
				this.fixFingers()

				fmt.Println("Node: ", hex.EncodeToString(this.node.nodeId))
				fmt.Println("Successor: ", hex.EncodeToString(this.node.finger[0].node.nodeId))
				if this.node.predecessor != nil {
					fmt.Println("Predecessor: ", hex.EncodeToString(this.node.predecessor.nodeId))
				}

				fmt.Println("Keys: ")
				for x := range this.keyValue {
					fmt.Println(x)
				}
				fmt.Println("")

				//// check what the 80th finger is
				//if this.node.finger[79].node != nil {
				//	fmt.Println("80th finger: ", hex.EncodeToString(this.node.finger[79].node.nodeId))
				//}

				//// check what the 120th finger is
				//if this.node.finger[119].node != nil {
				//	fmt.Println("120th finger: ", hex.EncodeToString(this.node.finger[119].node.nodeId))
				//}

				//// check what the 160th finger is
				//if this.node.finger[159].node != nil {
				//	fmt.Println("160th finger: ", hex.EncodeToString(this.node.finger[159].node.nodeId))
				//}

				//fmt.Println("")
			}
		}
	}()

	// Ping nodes periodically
	go func() {
		c := time.Tick(5 * time.Second)
		for {
			select {
			case <-c:
				this.sendPing()
			}
		}
	}()
}


//Tries to join the node at the specified address.
func (this *App) join(addr string) {
	req := JoinRequest{}
	req.Id = this.node.nodeId
	req.Addr = this.node.addr

	// get a node that is already in the ring
	fmt.Println("Calling Join on ", addr)

	bytes, _ := json.Marshal(req)
	r := this.transport.sendRequest(addr, "join", bytes)
	if r == nil {
		fmt.Println("Call error (join)")
		return
	}

	if r != nil {
		reply := JoinReply{}
		json.Unmarshal(r.Data, &reply)

		extNode := ExternalNode{}
		extNode.nodeId = reply.Id
		extNode.addr = addr

		// extNode is already in the ring
		this.addToRing(&extNode)
	}
	
}

func (this *App) addToRing(np *ExternalNode) {
	this.node.mutex.Lock()
	defer this.node.mutex.Unlock()

	this.node.predecessor = nil

	req := FindNodeReq{}
	req.Id = this.node.nodeId

	// call FindSuccessor on np, which is already in the ring
	bytes, _ := json.Marshal(req)
	r := this.transport.sendRequest(np.addr, "findSuccessor", bytes)
	if r == nil {
		fmt.Println("Call error (findSuccessor)")
		return
	}

	if r != nil {
		reply := FindNodeReply{}
		json.Unmarshal(r.Data, &reply)

		if reply.Found {
			extNode := new(ExternalNode)
			extNode.nodeId = reply.Id
			extNode.addr = reply.Addr

			this.node.finger[0].node = extNode
		} else {
			fmt.Println("Could not join ring: Successor not found.")
		}
	}
}

func (this *App) findSuccessor(id []byte) *ExternalNode {
	np := this.findPredecessor(id)

	if np != nil {
		// call GetSuccessor on np
		r := this.transport.sendRequest(np.addr, "getSuccessor", []byte{})
		if r == nil {
			fmt.Println("Call error (getSuccessor)")
			return nil
		}

		if r != nil {
			reply := FindNodeReply{}
			json.Unmarshal(r.Data, &reply)

			if reply.Found {
				// now we have the successor
				successor := new(ExternalNode)
				successor.nodeId = reply.Id
				successor.addr = reply.Addr
				return successor
			}
		}
	}
	extNode := new(ExternalNode)
	extNode.nodeId = this.node.nodeId
	extNode.addr = this.node.addr
	return extNode
}

func (this *App) findPredecessor(id []byte) *ExternalNode {

	if between(this.node.nodeId, this.node.finger[0].node.nodeId, id) {
		extNode := new(ExternalNode)
		extNode.nodeId = this.node.nodeId
		extNode.addr = this.node.addr
		return extNode
	} else {
		np := this.node.closestPrecedingFinger(id)
		req := FindNodeReq{}
		req.Id = id

		// call FindPredecessor on np
		bytes, _ := json.Marshal(req)
		r := this.transport.sendRequest(np.addr, "findPredecessor", bytes)

		if r == nil {
			fmt.Println("Call error (findPredecessor)")
			return nil
		}


		return nil
	}
}

func (this *App) lookup(key string) *ExternalNode {
	id := big.Int{}
	id.SetString(key, 16)
	idBytes := id.Bytes()

	return this.findSuccessor(idBytes)
}

func (this *App) stabilize() {

	if this.node.finger[0].node != nil {
		this.node.mutex.Lock()
		defer this.node.mutex.Unlock()

		successor := this.node.finger[0].node

		// call GetPredecessor on this.node's successor
		r := this.transport.sendRequest(successor.addr, "getPredecessor", []byte{})

		if r == nil {
			fmt.Println("Call error (getPredecessor)")
			return
		}
		if r != nil {
			reply := FindNodeReply{}
			json.Unmarshal(r.Data, &reply)

			// now we have the predecessor
			predecessor := new(ExternalNode)
			predecessor.nodeId = reply.Id
			predecessor.addr = reply.Addr
			if predecessor != nil && reply.Found && between3(this.node.nodeId, successor.nodeId, predecessor.nodeId) {
				this.node.finger[0].node = predecessor
			}

			msg := NotifyMsg{}
			msg.NodeId = this.node.nodeId
			msg.Addr = this.node.addr

			bytes, _ := json.Marshal(msg)
			this.transport.sendMsg(successor.addr, "notify", bytes)
		}
	}
}

func (this *App) fixFingers() {
	this.node.mutex.Lock()
	defer this.node.mutex.Unlock()

	i := rand.Intn(num_bits)
	successor := this.findSuccessor(this.node.finger[i].start)
	this.node.finger[i].node = successor
}

func (this *App) listen() {

	msgChan := make(chan *Msg)
	reqChan := make(chan *RequestContext)
	go this.transport.listen(msgChan, reqChan)

	for {
		select {
			case msg := <- msgChan:
				switch msg.Id {
				case "notify":
					go this.net.notify(msg)
					break
				case "insertKey":
					go this.net.insertKey(msg)
					break
				}
			case req := <- reqChan:
				switch req.req.Id {
				case "join":
					go this.net.join(req)
					break
				case "findSuccessor":
					go this.net.findSuccessor(req)
					break
				case "findPredecessor":
					go this.net.findPredecessor(req)
					break
				case "getSuccessor":
					go this.net.getSuccessor(req)
					break
				case "getPredecessor":
					go this.net.getPredecessor(req)
					break
				case "getKey":
					go this.net.getKey(req)
					break
				case "deleteKey":
					go this.net.deleteKey(req)
					break
				case "updateKey":
					go this.net.updateKey(req)
					break
				case "ping":
					go this.net.ping(req)
					break
				}				

		}
	}
}


func (this *App) sendPing() {
	if this.node.predecessor != nil {
		r := this.transport.sendRequest(this.node.predecessor.addr, "ping", []byte{})
		
		if r == nil {
			// Predecessor has timed out
			fmt.Println("Predecessor has timed out")
			this.node.predecessor = nil
		}
	}

	////////////////
	// We always need a successor to be set.
	/////////////

	for i := 0; i < num_bits; i++ {
		finger := this.node.finger[i].node
		if finger != nil {
			r := this.transport.sendRequest(finger.addr, "ping", []byte{})

			if r == nil {
				// Finger[i] has timed out
				fmt.Println("finger[" + strconv.Itoa(i) + "] has timed out")
				this.node.finger[i].node = nil
			}
		}
	}
}
