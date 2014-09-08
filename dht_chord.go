// chord
package dht

import (
	"fmt"
	"math"
	"strconv"
)

type Node struct {
	nodeId      string
	nodeIp      string
	nodePort    string
	finger      [160]*Node
	successor   *Node
	predecessor *Node
}

var fingerTable [160]*Node

func makeDHTNode(id *string, ip string, port string) *Node {
	if id == nil {
		newid := generateNodeId()
		newNode := Node{newid, ip, port, fingerTable, nil, nil}
		newNode.successor = &newNode // add yourself as successor
		newNode.predecessor = &newNode
		return &newNode
	} else {
		newNode := Node{*id, ip, port, fingerTable, nil, nil}
		newNode.successor = &newNode
		newNode.predecessor = &newNode
		return &newNode
	}
}

func (n *Node) addToRing(node *Node) {
	responsibleNode := n.lookup(node.nodeId)
	node.successor = responsibleNode.successor
	node.predecessor = responsibleNode
	responsibleNode.successor = node
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

func (n *Node) updateFingerTables() {
	k := 1
	nodeid := []byte(n.nodeId)
	bits := n.getNumberOfBits()
	for k <= bits {
		s, _ := calcFinger(nodeid, k, bits)
		n.finger[k-1] = n.lookup(s)
		fmt.Println("Node " + n.nodeId + ", Finger " + strconv.Itoa(k) + ": " + n.finger[k-1].nodeId)
		k++
	}
	fmt.Println("")
}

func (n *Node) printRing() {
	fmt.Println(n.nodeId)
	for i := n.successor; i != n; i = i.successor {
		fmt.Println(i.nodeId)
	}
}

func (n *Node) testCalcFingers(k int, m int) {

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
