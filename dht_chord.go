// chord
package dht

import (
	"fmt"
)

type Node struct {
	nodeId    string
	nodeIp    string
	nodePort  string
	finger    []*Node
	successor *Node
}

func makeDHTNode(id *string, ip string, port string) *Node {
	if id == nil {
		newid := generateNodeId()
		newNode := Node{newid, ip, port, nil, nil}
		newNode.successor = &newNode // add yourself as successor
		return &newNode
	} else {
		newNode := Node{*id, ip, port, nil, nil}
		newNode.successor = &newNode
		return &newNode
	}
}

func (n *Node) addToRing(node *Node) {
	responsibleNode := n.lookup(node.nodeId)
	node.successor = responsibleNode.successor
	responsibleNode.successor = node
}

func (n *Node) printRing() {
	fmt.Println(n.nodeId + "\n")
	for i := n.successor; i != nil && i != n; i = i.successor {
		fmt.Println(i.nodeId + "\n")
	}
}

func (n *Node) testCalcFingers(k int, m int) { // k = finger index, m = number of bits
	// calcFinger(n []byte, k int, m int)
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
