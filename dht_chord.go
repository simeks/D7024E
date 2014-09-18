// chord
package dht

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	//"math/rand"
	"strconv"
)

type Finger struct {
	node  *Node
	start []byte
}

type Node struct {
	nodeId      []byte
	nodeIp      string
	nodePort    string
	finger      [3]Finger
	successor   *Node
	predecessor *Node
}

var fingerTable [3]Finger

func makeDHTNode(id *string, ip string, port string) *Node {
	var newId string

	if id == nil {
		newId = generateNodeId()
	} else {
		newId = *id
	}

	x := new(big.Int)
	x.SetString(newId, 16)
	result := x.Bytes()

	for i := 0; i < 3; i++ {
		_, start := calcFinger(result, i+1, 3)
		fingerTable[i].start = start
	}
	newNode := Node{result, ip, port, fingerTable, nil, nil}
	newNode.successor = &newNode
	newNode.predecessor = &newNode

	// fill the finger table
	for i := 0; i < 3; i++ {
		newNode.finger[i].node = &newNode
	}
	return &newNode
}

func (n *Node) addToRing(node *Node) {

	//init finger table
	n.initFingerTable(node)

	//update others
	node.updateOthers()

	//move keys in (predecessor, n] from successor
	//...

}

func (n *Node) initFingerTable(node *Node) {
	node.finger[0].node = n.findSuccessor(node.finger[0].start)

	node.successor = node.finger[0].node
	node.predecessor = node.successor.predecessor
	node.successor.predecessor = node
	node.predecessor.successor = node

	fmt.Println("")
	fmt.Println("------------------------------------------------- Node: ", node.nodeId)

	//fmt.Println("node.start: ", node.finger[0].start)
	//fmt.Println("finger[0]: ", node.finger[0].node.nodeId)
	//fmt.Println("")

	id1 := node.nodeId
	var id2 []byte
	var keyId []byte

	for i := 1; i < 3; i++ {
		id2 = node.finger[i-1].node.nodeId
		keyId = node.finger[i].start
		//if between(id1, id2, keyId) { // if keyId is in [node, finger[i-1].node)
		if between2(id1, id2, keyId) {
			node.finger[i].node = node.finger[i-1].node
		} else {
			node.finger[i].node = n.findSuccessor(node.finger[i].start)
		}
		//fmt.Println("node.start: ", node.finger[i].start)
		//fmt.Println("finger[i]: ", node.finger[i].node.nodeId)
		//fmt.Println("")
	}
}

func (n *Node) updateOthers() {
	for i := 0; i < 3; i++ {

		x := new(big.Int)
		two := big.NewInt(2)
		sum := new(big.Int)
		sum.Exp(two, big.NewInt(int64(i)), nil)
		x.SetString(fmt.Sprintf("%x", n.nodeId), 16)
		x.Sub(x, sum)

		result := x.Bytes()
		p := n.findSuccessor(result)

		if bytes.Compare(p.nodeId, result) != 0 {
			p = p.predecessor
		}

		if x.Int64() >= 0 {
			fmt.Println("n - 2^i: ", x)
			fmt.Println("byte array result: ", result)
			fmt.Println("update this node: ", p.nodeId)

			p.updateFingerTable(n, i)
		}
	}
}

func (n *Node) updateFingerTable(s *Node, i int) {
	id1 := n.nodeId
	id2 := n.finger[i].node.nodeId
	keyId := s.nodeId

	//if between(id1, id2, keyId) {
	if strictlyBetween(id1, id2, keyId) {
		n.finger[i].node = s
		fmt.Println("")
		fmt.Println("")
		fmt.Print("Node ", n.nodeId)
		fmt.Print(", finger["+strconv.Itoa(i)+"] is now ", s.nodeId)
		fmt.Println("")
		fmt.Print("Node ", s.nodeId)
		fmt.Print(" has now updated node ", n.nodeId)
		fmt.Println("")
		fmt.Println("")

		p := n.predecessor
		p.updateFingerTable(s, i)
	}
}

func (n *Node) findSuccessor(id []byte) *Node {
	predecessor := n.findPredecessor(id)
	return predecessor.successor
}

func (n *Node) findPredecessor(id []byte) *Node {
	np := n

	for between2(np.nodeId, np.successor.nodeId, id) == false {
		np = np.closestPrecedingFinger(id)
	}
	return np
}

func (n *Node) closestPrecedingFinger(id []byte) *Node {
	id1 := n.nodeId

	for i := 3; i > 0; i-- {
		keyId := n.finger[i-1].node.nodeId
		if strictlyBetween(id1, id, keyId) { // if keyId is in (n, id)
			return n.finger[i-1].node
		}
	}
	return n
}

// periodically verify n's immediate successor,
// and tell the successor about n
//func (n *Node) stabilize() {
//	x := n.successor.predecessor

//	if strictlyBetween(n.nodeId, n.successor.nodeId, x.nodeId) {
//		n.successor = x
//		n.successor.notify(n)
//	}
//}

// n2 thinks it might be our predecessor
//func (n *Node) notify(n2 *Node) {
//	if n.predecessor == nil || strictlyBetween(n.predecessor.nodeId, n.nodeId, n2.nodeId) {
//		n.predecessor = n2
//	}
//}

// periodically refresh finger table entries
//func (n *Node) fixFingers() {
//	i := rand.Int63n(3)
//	n.finger[i].node = n.findSuccessor(n.finger[i].start)
//}

func (n *Node) printRing() {
	fmt.Println("Node "+":", hex.EncodeToString(n.nodeId))
	fmt.Println("Successor "+":", hex.EncodeToString(n.successor.nodeId))
	fmt.Println("Predecessors"+":", hex.EncodeToString(n.predecessor.nodeId))
	fmt.Println("")

	for i := n.successor; i != n; i = i.successor {
		fmt.Println("Node "+":", hex.EncodeToString(i.nodeId))
		fmt.Println("Successor "+":", hex.EncodeToString(i.successor.nodeId))
		fmt.Println("Predecessors"+":", hex.EncodeToString(i.predecessor.nodeId))
		fmt.Println("")
	}
}

func (n *Node) testCalcFingers(k int, m int) {

}

//func (n *Node) lookup(key []byte) *Node {
//	id1 := n.nodeId
//	id2 := n.successor.nodeId

//	if between2(id1, id2, key) {
//		return n
//	} else {
//		return n.successor.lookup(key)
//	}
//}
