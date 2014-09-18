package dht

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
)

var num_bits int = 160

type Finger struct {
	node  *Node
	start []byte
}

type Node struct {
	nodeId      []byte
	ip          string
	port        string
	finger      [160]Finger
	predecessor *Node
}

func makeDHTNode(id *string, ip string, port string) *Node {
	if id == nil {
		idStr := generateNodeId()
		id = &idStr
	}

	x := big.Int{}
	x.SetString(*id, 16)
	idBytes := x.Bytes()

	newNode := new(Node)
	newNode.nodeId = idBytes
	newNode.ip = ip
	newNode.port = port
	newNode.predecessor = newNode

	for i := 0; i < num_bits; i++ {
		_, start := calcFinger(idBytes, i+1, num_bits)
		newNode.finger[i].start = start
	}
	newNode.finger[0].node = newNode

	return newNode
}

// this joins the network;
// node is an arbitrary node in the network
func (this *Node) addToRing(np *Node) {
	this.predecessor = nil
	this.finger[0].node = np.findSuccessor(this.nodeId)
}

// ask node n to find id's successor
func (this *Node) findSuccessor(id []byte) *Node {
	np := this.findPredecessor(id)
	return np.finger[0].node
}

// ask node n to find id's predecessor
func (this *Node) findPredecessor(id []byte) *Node {
	np := this

	for ; between(np.nodeId, np.finger[0].node.nodeId, id) == false; {
		np = np.closestPrecedingFinger(id)
	}
	return np
}

// return closest finger preceding id
func (this *Node) closestPrecedingFinger(id []byte) *Node {
	for i := num_bits-1; i >= 0; i-- {
		if this.finger[i].node != nil && between3(this.nodeId, id, this.finger[i].node.nodeId) {
			return this.finger[i].node
		}
	}
	return this
}

func (this *Node) lookup(key string) *Node {
	return new(Node)
}

// periodically verify n’s immediate successor,
// and tell the successor about n. 
func (this *Node) stabilize() {
	x := this.finger[0].node.predecessor
	if x != nil && between3(this.nodeId, this.finger[0].node.nodeId, x.nodeId) {
		this.finger[0].node = x
	}
	this.finger[0].node.notify(this)
}

// np thinks it might be our predecessor. 
func (this *Node) notify(np *Node) {
	if this.predecessor == nil || between3(this.predecessor.nodeId, this.nodeId, np.nodeId) {
		this.predecessor = np
	}
}

func (this *Node) fixFingers() {
	i := rand.Intn(num_bits)
	this.finger[i].node = this.findSuccessor(this.finger[i].start)
}

func (this *Node) printRing() {
	this.printNode()

	for i := this.finger[0].node; i != this; i = i.finger[0].node {
		i.printNode()
	}
}

func (this *Node) printNode() {
	fmt.Println("Node "+":", hex.EncodeToString(this.nodeId))
	if this.finger[0].node != nil {
		fmt.Println("Successor "+":", hex.EncodeToString(this.finger[0].node.nodeId))
	}
	if this.predecessor != nil {
		fmt.Println("Predecessors"+":", hex.EncodeToString(this.predecessor.nodeId))
	}
	fmt.Println("Fingers")
	for i := 0; i < num_bits; i++ {
		if this.finger[i].node != nil {
			fmt.Println("Start:",hex.EncodeToString(this.finger[i].start),"Id:", hex.EncodeToString(this.finger[i].node.nodeId))
		}
	}
	fmt.Println("")	
}

func (this *Node) testCalcFingers(k int, m int) {
}
