package main

import (
	//"encoding/hex"
	//"fmt"
	"math/big"
	//"math/rand"
	//"strconv"
	"sync"
)

const num_bits int = 160

type Finger struct {
	node  *ExternalNode
	start []byte
}

type Node struct {
	nodeId      []byte
	ip          string
	port        string
	finger      [num_bits]Finger
	predecessor *ExternalNode
	mutex       sync.Mutex
}

type ExternalNode struct {
	nodeId []byte
	ip     string
	port   string
}

func makeDHTNode(id *string, ip string, port string) *Node {
	if id == nil {
		idStr := generateNodeId()
		id = &idStr
	}

	x := big.Int{}
	x.SetString(*id, 16)
	idBytes := x.Bytes()

	externalNode := new(ExternalNode)
	externalNode.nodeId = idBytes
	externalNode.ip = ip
	externalNode.port = port

	newNode := new(Node)
	newNode.nodeId = idBytes
	newNode.ip = ip
	newNode.port = port
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
	extNode.ip = this.ip
	extNode.port = this.port
	return extNode
}

// np thinks it might be our predecessor.
func (this *Node) notify(np *ExternalNode) {
	if this.predecessor == nil || between3(this.predecessor.nodeId, this.nodeId, np.nodeId) {
		this.predecessor = np
	}
}
