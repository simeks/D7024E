// chord
package dht

import (
	"fmt"
)

var ring []Node

type Node struct {
	nodeId   string
	nodeIp   string
	nodePort string
}

func makeDHTNode(id *string, ip string, port string) Node {
	if id == nil {
		newid := generateNodeId()
		return Node{newid, ip, port}
	} else {
		return Node{*id, ip, port}
	}
}

func (n *Node) addToRing(node Node) {
	ring = append(ring, node)
}

func (n *Node) printRing() {
	for i := range ring {
		fmt.Println(ring[i].nodeId + "\n") // print all nodes in the ring
	}
}

func (n *Node) testCalcFingers(i int, j int) { // i = ?? j = number of bits

}

func (n *Node) lookup(hashKey string) Node {
	return Node{"test1", "test2", "test3"} // returnera rätt node här
}
