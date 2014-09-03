// chord
package dht

import ()

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

}

func (n *Node) printRing() {

}

func (n *Node) testCalcFingers(i int, j int) { // i = ?? j = number of bits

}

func (n *Node) lookup(hashKey string) Node {
	return Node{"test1", "test2", "test3"} // returnera rätt node här
}
