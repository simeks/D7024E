package dht

import (
	"fmt"
	"encoding/hex"
	"math/big"
)

var num_bits int = 3

type Finger struct {
	node  *Node
	start []byte
}

type Node struct {
	nodeId      []byte
	ip      	string
	port    	string
	finger      [3]Finger
	successor   *Node
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
	newNode.predecessor = nil
	newNode.successor = newNode

	newNode.addToRing(nil)

	return newNode
}

// node joins the network;
// this is an arbitrary node in the network
func (this *Node) addToRing(node *Node) {
	if node != nil {
		node.initFingerTables(node)
		this.updateOthers()
		// move keys in (predecessor, n] from successor
	} else { // this is the only node in the network
		for i := 0; i < num_bits; i++ {
			this.finger[i].node = this
		}
		this.predecessor = this
	}

}

// initialize finger table of local node;
// node is an arbitrary node already in the network
func (this *Node) initFingerTables(node *Node) {
	this.finger[0].node = node.findSuccessor(this.finger[0].start)
	this.predecessor = this.finger[0].node.predecessor
	this.finger[0].node.predecessor = this

	for i := 0; i < num_bits-1; i++ {
		if between(this.nodeId, this.finger[i].node.nodeId, this.finger[i+1].start) {
			this.finger[i+1].node = this.finger[i].node
		} else {
			this.finger[i+1].node = node.findSuccessor(this.finger[i+1].start)
		}
	}
}

// update all nodes whose finger
// tables should refer to node
func (this *Node) updateOthers() {
	for i := 0; i < num_bits; i++ {
		// find last node p whose i:th finger might be n

		pow := big.Int{}
		pow.Exp(big.NewInt(2), big.NewInt(int64(i)), nil)

		thisId := big.Int{}
		thisId.SetBytes(this.nodeId)

		id := big.Int{}
		id.Sub(&thisId, &pow)

		p := this.findPredecessor(id.Bytes())
		p.updateFingerTable(this, i)
	}
}

// If s is the i:th finger of n, update n's finger table with s
func (this *Node) updateFingerTable(s *Node, i int) {
	if between(this.nodeId, this.finger[i].node.nodeId, s.nodeId) {
		this.finger[i].node = s
		p := this.predecessor
		p.updateFingerTable(s, i)
	}
}

func (this *Node) findSuccessor(id []byte) *Node {
	n2 := this.findPredecessor(id)
	return n2.successor
}

func (this *Node) findPredecessor(id []byte) *Node {
	np := this

	for ; between(np.nodeId, np.successor.nodeId, id) == false; {
		np = np.closestPrecedingFinger(id)
	}
	return np
}

func (this *Node) closestPrecedingFinger(id []byte) *Node {
	for i := num_bits-1; i >= 0; i-- {
		if between(this.nodeId, id, this.finger[i].node.nodeId) {
			return this.finger[i].node
		}
	}
	return this
}

func (this *Node) lookup(key string) *Node {
	return new(Node)
}
func (this *Node) printRing() {
	fmt.Println("Node "+":", hex.EncodeToString(this.nodeId))
	if this.successor != nil {
		fmt.Println("Successor "+":", hex.EncodeToString(this.successor.nodeId))
	}
	if this.predecessor != nil {
		fmt.Println("Predecessors"+":", hex.EncodeToString(this.predecessor.nodeId))
	}
	fmt.Println("")

	for i := this.successor; i != this; i = i.successor {
		fmt.Println("Node "+":", hex.EncodeToString(i.nodeId))
		if this.successor != nil {
			fmt.Println("Successor "+":", hex.EncodeToString(i.successor.nodeId))
		}
		if this.predecessor != nil {
			fmt.Println("Predecessors"+":", hex.EncodeToString(i.predecessor.nodeId))
		}
		fmt.Println("")
	}
}

func (this *Node) testCalcFingers(k int, m int) {

}
