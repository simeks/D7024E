package main

import (
	//"encoding/hex"
	"fmt"
	"math/big"
	//"math/rand"
	"time"
	"encoding/json"
)

// 0 = no time out
const time_out int = 3

type App struct {
	node    *Node
	transport Transport

}

func (this *App) init(bindAddr, bindPort string) {
	this.node = makeDHTNode(nil, bindAddr, bindPort)
	this.transport.bindAddress = bindAddr+":"+bindPort

	// call stabilize and fixFingers periodically
	go func() {
		c := time.Tick(3 * time.Second)
		for {
			select {
			case <-c:
				this.stabilize()
				this.fixFingers()

				//fmt.Println("Node: ", hex.EncodeToString(this.node.nodeId))
				//fmt.Println("Successor: ", hex.EncodeToString(this.node.finger[0].node.nodeId))
				if this.node.predecessor != nil {
				//	fmt.Println("Predecessor: ", hex.EncodeToString(this.node.predecessor.nodeId))
				}

				//fmt.Println("Keys: ")
				for x := range this.node.keys {
					fmt.Println(x)
				}
				//fmt.Println("")

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

type MsgA struct {
	Stuff string
}

//Tries to join the node at the specified address.
func (this *App) join(addr string) {
	msg := MsgA{}
	msg.Stuff = "asd"
	bytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	this.transport.sendMsg(addr, "asd", bytes)
	this.transport.sendMsg(addr, "asd2", bytes)
	this.transport.sendMsg(addr, "asd3", bytes)

	/*
	args := new(AddArgs)
	args.Id = this.node.nodeId
	args.Ip = this.node.ip
	args.Port = this.node.port

	reply := new(AddReply)

	// get a node that is already in the ring
	fmt.Println("Calling Join on ", addr, "\n")
	err := this.nodeUDP.CallUDP("Join", addr, args, reply, time_out)

	if err != nil {
		fmt.Print("Call error - ")
		fmt.Println(err.Error())
		return
	}

	if reply != nil {
		extNode := new(ExternalNode)
		extNode.nodeId = reply.Id
		extNode.ip = reply.Ip
		extNode.port = reply.Port

		// extNode is already in the ring
		this.addToRing(extNode)
	}
	*/
}

func (this *App) addToRing(np *ExternalNode) {
	/*
	this.node.mutex.Lock()
	defer this.node.mutex.Unlock()

	this.node.predecessor = nil

	args := new(AddArgs)
	args.Id = this.node.nodeId
	args.Ip = this.node.ip
	args.Port = this.node.port

	reply := new(AddReply)

	addr := np.ip + ":" + np.port

	// call FindSuccessor on np, which is already in the ring
	err := this.nodeUDP.CallUDP("FindSuccessor", addr, args, reply, time_out)

	if err != nil {
		fmt.Print("Call error - ")
		fmt.Println(err.Error())
		return
	}

	if reply != nil {
		extNode := new(ExternalNode)
		extNode.nodeId = reply.Id
		extNode.ip = reply.Ip
		extNode.port = reply.Port

		this.node.finger[0].node = extNode
	}
	*/
}

func (this *App) findSuccessor(id []byte) *ExternalNode {
	/*
	np := this.findPredecessor(id)

	if np != nil {
		args := new(AddArgs)
		reply := new(AddReply)

		addr := np.ip + ":" + np.port

		// call GetSuccessor on np
		err := this.nodeUDP.CallUDP("GetSuccessor", addr, args, reply, time_out)

		if err != nil {
			fmt.Print("Call error - ")
			fmt.Println(err.Error())
			return nil
		}

		if reply != nil {
			// now we have the successor
			successor := new(ExternalNode)
			successor.nodeId = reply.Id
			successor.ip = reply.Ip
			successor.port = reply.Port
			return successor
		}
	}
	extNode := new(ExternalNode)
	extNode.nodeId = this.node.nodeId
	extNode.ip = this.node.ip
	extNode.port = this.node.port
	return extNode
	*/
	return nil
}

func (this *App) findPredecessor(id []byte) *ExternalNode {
/*
	if between(this.node.nodeId, this.node.finger[0].node.nodeId, id) {
		extNode := new(ExternalNode)
		extNode.nodeId = this.node.nodeId
		extNode.ip = this.node.ip
		extNode.port = this.node.port
		return extNode
	} else {
		np := this.node.closestPrecedingFinger(id)
		args := new(AddArgs)
		args.Id = id

		reply := new(AddReply)

		addr := np.ip + ":" + np.port

		// call FindPredecessor on np
		err := this.nodeUDP.CallUDP("FindPredecessor", addr, args, reply, time_out)

		if err != nil {
			fmt.Print("Call error - ")
			fmt.Println(err.Error())
			return nil
		}
		return nil
	}
	*/
	return nil
}

func (this *App) lookup(key string) *ExternalNode {
	id := big.Int{}
	id.SetString(key, 16)
	idBytes := id.Bytes()

	return this.findSuccessor(idBytes)
}

func (this *App) stabilize() {
/*
	if this.node.finger[0].node != nil {
		this.node.mutex.Lock()
		defer this.node.mutex.Unlock()
		args := new(AddArgs)
		reply := new(AddReply)

		addr := this.node.finger[0].node.ip + ":" + this.node.finger[0].node.port
		// call GetPredecessor on this.node's successor
		//err := this.nodeUDP.CallUDP("GetPredecessor", addr, args, reply, time_out)

		if err != nil {
			fmt.Print("Call error - ")
			fmt.Println(err.Error())
			return
		}
		if reply != nil {
			// now we have the predecessor
			predecessor := new(ExternalNode)
			predecessor.nodeId = reply.Id
			predecessor.ip = reply.Ip
			predecessor.port = reply.Port
			if predecessor != nil && between3(this.node.nodeId, this.node.finger[0].node.nodeId, predecessor.nodeId) {
				this.node.finger[0].node = predecessor
			}
			addr := this.node.finger[0].node.ip + ":" + this.node.finger[0].node.port

			//args.Id = this.node.nodeId
			//args.Ip = this.node.ip
			//args.Port = this.node.port
			//err := this.nodeUDP.CallUDP("Notify", addr, args, reply, time_out)
			//if err != nil {
			//	fmt.Print("Call error - ")
			//	fmt.Println(err.Error())
			//	return
			//}
		}
	}
	*/
}

func (this *App) fixFingers() {
	/*this.node.mutex.Lock()
	defer this.node.mutex.Unlock()

	i := rand.Intn(num_bits)
	successor := this.findSuccessor(this.node.finger[i].start)
	this.node.finger[i].node = successor
	*/
}

func (this *App) listen() {

	msgChan := make(chan Msg)
	reqChan := make(chan Request)
	go this.transport.listen(msgChan, reqChan)

	for {
		select {
			case msg := <- msgChan:
				fmt.Println("Msg:",msg.Id)

				msga := MsgA{}
				err := json.Unmarshal(msg.Data, &msga)
				if err != nil {
					fmt.Println("Error:",err)
				}
				fmt.Println("msg",msga.Stuff)

			case req := <- reqChan:
				fmt.Println("Req:",req.Id)

		}
	}
}

func (this *App) sendPing() {
	/*
	args := new(AddArgs)
	reply := new(AddReply)
	args.Id = this.node.nodeId

	if this.node.predecessor != nil {

		//err := this.nodeUDP.CallUDP("Ping", this.node.predecessor.ip+":"+this.node.predecessor.port, args, reply, 3)
		
		//if err != nil {
			// Predecessor has timed out
		//	fmt.Println("Predecessor has timed out")
		//	this.node.predecessor = nil
		//}
	}

	for i := 0; i < num_bits; i++ {
		finger := this.node.finger[i].node
		if finger != nil {

			//err := this.nodeUDP.CallUDP("Ping", finger.ip+":"+finger.port, args, reply, 3)

			//if err != nil {
				// Finger[i] has timed out
			//	fmt.Println("finger[" + strconv.Itoa(i) + "] has timed out")
			//	this.node.finger[i].node = nil
			//}
		}
	}
	*/
}
