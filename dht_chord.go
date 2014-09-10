// chord
package dht

import (
	"fmt"
	"math"
	"strconv"
)

type Finger struct {
	node  *Node
	start string
}

type Node struct {
	nodeId      string
	nodeIp      string
	nodePort    string
	finger      [160]Finger
	successor   *Node
	predecessor *Node
}

var fingerTable [160]Finger

func makeDHTNode(id *string, ip string, port string) *Node {
	var newId string

	if id == nil {
		newId = generateNodeId()
	} else {
		newId = *id
	}

	for i := 0; i < 160; i++ {
		start, _ := calcFinger([]byte(newId), i+1, 160)
		fingerTable[i].start = start
	}
	newNode := Node{newId, ip, port, fingerTable, nil, nil}
	newNode.successor = &newNode
	newNode.predecessor = &newNode

	// fill the finger table
	for i := 0; i < 160; i++ {
		newNode.finger[i].node = &newNode
	}
	return &newNode
}

func (n *Node) addToRing(node *Node) {

	node.successor = n.findSuccessor(node.nodeId)
	node.predecessor = node.successor.predecessor
	node.successor.predecessor = node

	//init finger table
	//n.initFingerTable(node)

	//update others
	//node.updateOthers()

	//move keys in (predecessor, n] from successor
	//...

}

//func (n *Node) initFingerTable(node *Node) {
//	 node.finger[0] = node.successor

//	for i := 1; i < 160; i++ {
//		if node.finger[i].start is in [n, node.finger[i-1]) {
//			node.finger[i] = node.finger[i-1]
//		} else {
//			node.finger[i] = n.findSuccessor(finger[i].start)
//		}
//	}
//}

//func (n *Node) updateOthers() {
//	for i := 0; i < 160; i++ {
//		 find last node p whose ith finger might be n
//		 p := findPredecessor(n - 2^(i-1))
//		 p.updateFingerTable(n, i)

//	}
//}

//func (n *Node) updateFingerTable(node *Node, i int) {
//	if node is in [n, n.finger[i-1]) {
//		finger[i-1] = node
//		p := n.predecessor
//		p.updateFingerTable(node, i)
//	}
//}

func (n *Node) findSuccessor(id string) *Node {
	predecessor := n.findPredecessor(id)
	return predecessor.successor
}

func (n *Node) findPredecessor(id string) *Node {
	//id1 := []byte(n.nodeId)
	//id2 := []byte(n.successor.nodeId)
	//keyId := []byte(id)
	//predecessor := n

	//for i := n; between(id1, id2, keyId) == false; i = i.closestPrecedingFinger(id) {
	//	id1 = []byte(i.closestPrecedingFinger(id).nodeId)
	//	id2 = []byte(i.closestPrecedingFinger(id).successor.nodeId)
	//	predecessor = i.closestPrecedingFinger(id)
	//}
	//return predecessor
	return n.lookup(id)
}

func (n *Node) closestPrecedingFinger(id string) *Node {
	id1 := []byte(n.nodeId)
	id2 := []byte(id)

	for i := 160; i > 0; i-- {
		keyId := []byte(n.finger[i-1].node.nodeId)
		if between(id1, id2, keyId) {
			return n.finger[i].node
		}
	}
	return n
}

func (n *Node) getNumberOfBits() int {
	i := 0
	m := 0
	length := n.ringLength()
	for length > int(math.Pow(float64(2), float64(i))) {
		i++
		m++
	}
	return m
}

func (n *Node) ringLength() int {
	length := 1
	for i := n.successor; i != n; i = i.successor {
		length++
	}
	return length
}

// only for test
func (n *Node) setSuccessor(node *Node) {
	n.successor = node
}

func (n *Node) printRing() {
	fmt.Println("Node: " + n.nodeId + " Successor: " + n.successor.nodeId + " Predecessor: " + n.predecessor.nodeId)
	for i := n.successor; i != n; i = i.successor {
		fmt.Println("Node: " + i.nodeId + " Successor: " + i.successor.nodeId + " Predecessor: " + i.predecessor.nodeId)
	}
}

func (n *Node) testCalcFingers(k int, m int) {

}

// only for test
func (n *Node) updateFingerTables() {
	k := 1
	nodeid := []byte(n.nodeId)
	//bits := n.getNumberOfBits()
	for k <= 160 {
		s, _ := calcFinger(nodeid, k, 160)
		n.finger[k-1].node = n.lookup(s)

		// printa bara ut 4 första fingrarna
		if k <= 4 {
			fmt.Println("s: " + s)
			fmt.Println("Node " + n.nodeId + ", Finger " + strconv.Itoa(k) + ": " + n.finger[k-1].node.nodeId)
		}
		k++
	}
	fmt.Println("")
}

func (n *Node) lookup(key string) *Node {
	id1 := []byte(n.nodeId)
	id2 := []byte(n.successor.nodeId)
	keyId := []byte(key)

	if between(id1, id2, keyId) {
		return n
	} else {
		return n.successor.lookup(key)
	}
}
