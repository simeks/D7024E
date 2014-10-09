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

	this.changePredecessor(nil)

	req := FindNodeReq{}
	req.Id = this.node.nodeId
	this.changeSuccessor(np.findSuccessor(&this.transport, this.node.nodeId))

	if this.node.finger[0].node == nil {
		fmt.Println("Could not join ring: Successor not found.")
	}
}

func (this *App) findSuccessor(id []byte) *ExternalNode {
	np := this.findPredecessor(id)

	if np != nil {
		// call GetSuccessor on np
		s := np.getSuccessor(&this.transport)

		if s != nil {
			return s
		}
	}
	extNode := new(ExternalNode)
	extNode.nodeId = this.node.nodeId
	extNode.addr = this.node.addr
	return extNode
}

func (this *App) findPredecessor(id []byte) *ExternalNode {
	n := &ExternalNode{this.node.nodeId, this.node.addr}
	succ := this.node.finger[0].node

	for between(n.nodeId, succ.nodeId, id) == false {
		n = n.closestPrecedingFinger(&this.transport, id)
		succ = n.getSuccessor(&this.transport)
	}

	return n
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
		p := successor.getPredecessor(&this.transport)

		if p != nil && between3(this.node.nodeId, successor.nodeId, p.nodeId) {
			this.changeSuccessor(p)
		}

		successor.notify(&this.transport, &ExternalNode{this.node.nodeId, this.node.addr})
	}
}

func (this *App) fixFingers() {
	this.node.mutex.Lock()
	defer this.node.mutex.Unlock()

	i := rand.Intn(num_bits)
	successor := this.findSuccessor(this.node.finger[i].start)

	if i == 0 {
		this.changeSuccessor(successor)
	} else {
		this.node.finger[i].node = successor
	}
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
				case "closestPrecedingFinger":
					go this.net.closestPrecedingFinger(req)
					break
				}				

		}
	}
}

func (this *App) pingFinger(i int) {
	finger := this.node.finger[i].node
	if finger != nil {
		r := this.transport.sendRequest(finger.addr, "ping", []byte{})

		if r == nil {
			// Finger[i] has timed out
			fmt.Println("finger[" + strconv.Itoa(i) + "] has timed out")
			if i == 0 {
				// We always need a successor to be set.
				this.changeSuccessor(&ExternalNode{this.node.nodeId, this.node.addr})
			} else {
				this.node.finger[i].node = nil
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
			this.changePredecessor(nil)
		}
	}

	for i := 0; i < num_bits; i++ {
		go this.pingFinger(i)
	}

}


// np thinks it might be our predecessor.
func (this *App) notify(np *ExternalNode) {
	if this.node.predecessor == nil || between3(this.node.predecessor.nodeId, this.node.nodeId, np.nodeId) {
		this.changePredecessor(np)
	}
}

func (this *App) updateKeyValue() {
	/*
	* 
	* (predecessor, this] - My values
	*/

	//for k, v := range this.keyValue {

	//}
}

func (this *App) changeSuccessor(n *ExternalNode) {
	/*
	- Notify successor about 
	*/

	this.node.finger[0].node = n
}

func (this *App) changePredecessor(n *ExternalNode) {
	// If the predecessor 

	this.node.predecessor = n
}