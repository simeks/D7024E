// chord
package dht

import (
	"fmt"
)

type Node struct {
	nodeId      string
	nodeIp      string
	nodePort    string
	finger      [3]*Node
	successor   *Node
	predecessor *Node
}

var fingerTable [3]*Node

func makeDHTNode(id *string, ip string, port string) *Node {
	if id == nil {
		newid := generateNodeId()
		newNode := Node{newid, ip, port, fingerTable, nil, nil}
		newNode.successor = &newNode // add yourself as successor
		return &newNode
	} else {
		newNode := Node{*id, ip, port, fingerTable, nil, nil}
		newNode.successor = &newNode
		return &newNode
	}
}

func (n *Node) addToRing(node *Node) {
	responsibleNode := n.lookup(node.nodeId)
	// fmt.Println(responsibleNode.nodeId) // lookup returnar fel då man använder 160 bitar
	node.successor = responsibleNode.successor
	node.predecessor = responsibleNode
	responsibleNode.successor = node
}

func (n *Node) updateFingerTables() {
	k := 1
	nodeid := []byte(n.nodeId)
	for k <= 3 {
		s, _ := calcFinger(nodeid, k, 3) // calculates every finger for node n, 3 bits
		n.finger[k-1] = n.lookup(s)
		k++
	}
}

func (n *Node) printRing() {
	fmt.Println(n.nodeId + "\n")
	for i := n.successor; i != n; i = i.successor {
		fmt.Println(i.nodeId + "\n")
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
