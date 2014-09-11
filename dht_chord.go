// chord
package dht

import (
	"fmt"
	//"math"
	"math/big"
	//"strconv"
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

	node.successor = n.findSuccessor(node.nodeId)
	node.predecessor = node.successor.predecessor
	node.successor.predecessor = node
	node.predecessor.successor = node

	//init finger table
	//n.initFingerTable(node)

	//update others
	//node.updateOthers()

	//move keys in (predecessor, n] from successor
	//...

}

func (n *Node) initFingerTable(node *Node) {
	node.finger[0].node = n.findSuccessor(node.finger[0].start)

	id1 := node.nodeId
	var id2 []byte
	var keyId []byte

	for i := 1; i < 3; i++ {
		id2 = node.finger[i-1].node.nodeId
		keyId = node.finger[i].start
		if between(id1, id2, keyId) {
			node.finger[i].node = node.finger[i-1].node
		} else {
			node.finger[i].node = n.findSuccessor(node.finger[i].start)
		}
	}
}

func (n *Node) updateOthers() {
	for i := 0; i < 3; i++ {

		x := new(big.Int)
		two := big.NewInt(2)
		sum := new(big.Int)
		sum.Exp(two, big.NewInt(int64(i)), nil)
		x.SetString(string(n.nodeId), 16)
		x.Sub(x, sum)
		result := x.Bytes()

		p := n.findPredecessor(result)
		p.updateFingerTable(n, i)
	}
}

func (n *Node) updateFingerTable(node *Node, i int) {
	id1 := n.nodeId
	id2 := n.finger[i].node.nodeId
	keyId := node.nodeId

	if between(id1, id2, keyId) {
		n.finger[i].node = node
		p := n.predecessor
		p.updateFingerTable(node, i)
	}
}

func (n *Node) findSuccessor(id []byte) *Node {
	predecessor := n.findPredecessor(id)
	return predecessor.successor
}

func (n *Node) findPredecessor(id []byte) *Node {
	//id1 := n.nodeId
	//id2 := n.successor.nodeId
	//predecessor := n

	//for i := n; between(id1, id2, id) == false; i = i.closestPrecedingFinger(id) {
	//	id1 = i.closestPrecedingFinger(id).nodeId
	//	id2 = i.closestPrecedingFinger(id).successor.nodeId
	//	fmt.Println("id1: " + string(id1) + " id2: " + string(id2))
	//	predecessor = i.closestPrecedingFinger(id)
	//}
	//return predecessor
	return n.lookup(id)
}

func (n *Node) closestPrecedingFinger(id []byte) *Node {
	id1 := n.nodeId
	id2 := id

	for i := 3; i > 0; i-- {
		keyId := n.finger[i-1].node.nodeId
		if between(id1, id2, keyId) {
			return n.finger[i-1].node
		}
	}
	return n
}

func (n *Node) printRing() {
	//fmt.Println("Node: " + string(n.nodeId) + " Successor: " + string(n.successor.nodeId) + " Predecessor: " + string(n.predecessor.nodeId))
	fmt.Println("Node "+":", n.nodeId)
	fmt.Println("Successor "+":", n.successor.nodeId)
	fmt.Println("Predecessors"+":", n.predecessor.nodeId)
	fmt.Println("")
	for i := n.successor; i != n; i = i.successor {
		//fmt.Println("Node: " + string(i.nodeId) + " Successor: " + string(i.successor.nodeId) + " Predecessor: " + string(i.predecessor.nodeId))
		fmt.Println("Node "+":", i.nodeId)
		fmt.Println("Successors"+":", i.successor.nodeId)
		fmt.Println("Predecessor "+":", i.predecessor.nodeId)
		fmt.Println("")
	}
}

func (n *Node) testCalcFingers(k int, m int) {

}

//only for test
//func (n *Node) updateFingerTables() {
//	k := 1
//	fmt.Println("Node ", n.nodeId)
//	for k <= 3 {
//		_, s := calcFinger(n.nodeId, k, 3)
//		n.finger[k-1].node = n.lookup(s)

//		// printa bara ut 3 första fingrarna
//		if k <= 3 {
//			fmt.Println("Finger "+strconv.Itoa(k)+": ", n.finger[k-1].node.nodeId)
//		}
//		k++
//	}
//	fmt.Println("")
//}

func (n *Node) lookup(key []byte) *Node {
	id1 := n.nodeId
	id2 := n.successor.nodeId

	if between(id1, id2, key) {
		return n
	} else {
		return n.successor.lookup(key)
	}
}
