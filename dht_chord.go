package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

const num_bits int = 160
const num_successors int = 6

type Finger struct {
	node  *ExternalNode
	start []byte
}

type Node struct {
	nodeId        []byte
	addr          string
	finger        [num_bits]Finger
	predecessor   *ExternalNode
	successorList [num_successors]*ExternalNode
	mutex         sync.Mutex
}

type ExternalNode struct {
	nodeId []byte
	addr   string
}

func makeDHTNode(id *string, addr string) *Node {
	if id == nil {
		idStr := generateNodeId()
		id = &idStr
	}

	idBytes := stringToId(*id)

	externalNode := new(ExternalNode)
	externalNode.nodeId = idBytes
	externalNode.addr = addr

	newNode := new(Node)
	newNode.nodeId = idBytes
	newNode.addr = addr
	newNode.predecessor = externalNode

	for i := 0; i < num_bits; i++ {
		_, start := calcFinger(idBytes, i+1, num_bits)
		newNode.finger[i].start = start
	}
	newNode.finger[0].node = externalNode

	return newNode
}

// return closest finger preceding id
func (this *Node) closestPrecedingFinger(id []byte) *ExternalNode {
	for i := num_bits - 1; i >= 0; i-- {
		if this.finger[i].node != nil && between3(this.nodeId, id, this.finger[i].node.nodeId) {
			return this.finger[i].node
		}
	}

	extNode := new(ExternalNode)
	extNode.nodeId = this.nodeId
	extNode.addr = this.addr
	return extNode
}

func (n *ExternalNode) findSuccessor(t *Transport, id []byte) *ExternalNode {
	req := FindNodeReq{}
	req.Id = id

	// call FindSuccessor on np, which is already in the ring
	bytes, _ := json.Marshal(req)
	r := t.sendRequest(n.addr, "findSuccessor", bytes)
	if r == nil {
		fmt.Println("Call error (findSuccessor)")
		return nil
	}

	if r != nil {
		reply := FindNodeReply{}
		json.Unmarshal(r.Data, &reply)

		if reply.Found {
			extNode := new(ExternalNode)
			extNode.nodeId = reply.Id
			extNode.addr = reply.Addr

			return extNode
		}
	}
	return nil
}

func (n *ExternalNode) findPredecessor(t *Transport, id []byte) *ExternalNode {
	req := FindNodeReq{}
	req.Id = id

	bytes, _ := json.Marshal(req)
	r := t.sendRequest(n.addr, "findPredecessor", bytes)
	if r == nil {
		fmt.Println("Call error (findPredecessor)")
		return nil
	}

	if r != nil {
		reply := FindNodeReply{}
		json.Unmarshal(r.Data, &reply)

		if reply.Found {
			extNode := new(ExternalNode)
			extNode.nodeId = reply.Id
			extNode.addr = reply.Addr

			return extNode
		}
	}
	return nil
}

func (n *ExternalNode) getSuccessor(t *Transport) *ExternalNode {
	r := t.sendRequest(n.addr, "getSuccessor", []byte{})
	if r == nil {
		fmt.Println("Call error (getSuccessor)")
		return nil
	}

	if r != nil {
		reply := FindNodeReply{}
		json.Unmarshal(r.Data, &reply)

		if reply.Found {
			extNode := new(ExternalNode)
			extNode.nodeId = reply.Id
			extNode.addr = reply.Addr

			return extNode
		}
	}
	return nil
}

func (n *ExternalNode) getSuccessorList(t *Transport) *[num_successors]*ExternalNode {
	r := t.sendRequest(n.addr, "getSuccessorList", []byte{})
	if r == nil {
		fmt.Println("Call error (getSuccessorList)")
		return nil
	}

	if r != nil {
		reply := SuccessorListReply{}
		json.Unmarshal(r.Data, &reply)

		succlist := new([num_successors]*ExternalNode)

		for i := 0; i < num_successors; i++ {
			extNode := new(ExternalNode)
			extNode.nodeId = reply.Id[i]
			extNode.addr = reply.Addr[i]
			succlist[i] = extNode
		}
		return succlist
	}
	return nil
}

func (n *ExternalNode) getPredecessor(t *Transport) *ExternalNode {
	r := t.sendRequest(n.addr, "getPredecessor", []byte{})
	if r == nil {
		fmt.Println("Call error (getPredecessor)")
		return nil
	}

	if r != nil {
		reply := FindNodeReply{}
		json.Unmarshal(r.Data, &reply)

		if reply.Found {
			extNode := new(ExternalNode)
			extNode.nodeId = reply.Id
			extNode.addr = reply.Addr

			return extNode
		}
	}
	return nil
}

func (n *ExternalNode) notify(t *Transport, en *ExternalNode) {
	msg := NotifyMsg{}
	msg.NodeId = en.nodeId
	msg.Addr = en.addr

	bytes, _ := json.Marshal(msg)
	t.sendMsg(n.addr, "notify", bytes)
}

func (n *ExternalNode) closestPrecedingFinger(t *Transport, id []byte) *ExternalNode {
	req := FindNodeReq{}
	req.Id = id

	bytes, _ := json.Marshal(req)
	r := t.sendRequest(n.addr, "closestPrecedingFinger", bytes)
	if r == nil {
		fmt.Println("Call error (closestPrecedingFinger)")
		return nil
	}

	if r != nil {
		reply := FindNodeReply{}
		json.Unmarshal(r.Data, &reply)

		if reply.Found {
			extNode := new(ExternalNode)
			extNode.nodeId = reply.Id
			extNode.addr = reply.Addr

			return extNode
		}
	}
	return nil
}

func (n *ExternalNode) transferData(t *Transport, kv *map[string]string) {
	msg := TransferDataMsg{}
	msg.KeyValue = *kv

	bytes, _ := json.Marshal(msg)
	t.sendMsg(n.addr, "transferData", bytes)
}
