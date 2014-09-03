// chord
package dht

import (
	"fmt"
)

type Node struct {
	nodeId    string
	nodeIp    string
	nodePort  string
	successor *Node
}

func makeDHTNode(id *string, ip string, port string) *Node {
	if id == nil {
		newid := generateNodeId()
		return &Node{newid, ip, port, nil}
	} else {
		return &Node{*id, ip, port, nil}
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

func (n *Node) testCalcFingers(i int, j int) { // i = ?? j = number of bits

}

func (n *Node) lookup(hashKey string) Node {
	return Node{"test1", "test2", "test3", nil} // returnera rätt node här
}
