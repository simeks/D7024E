package main

import (
	//"encoding/hex"
	//"fmt"
	"math/big"
	//"math/rand"
	//"strconv"
	"sync"
)

const num_bits int = 160

type Finger struct {
	node  *ExternalNode
	start []byte
}

type Node struct {
	nodeId      []byte
	ip          string
	port        string
	finger      [num_bits]Finger
	predecessor *ExternalNode
	mutex       sync.Mutex
	keys        map[string]string
}

type ExternalNode struct {
	nodeId []byte
	ip     string
	port   string
}

func makeDHTNode(id *string, ip string, port string) *Node {
	if id == nil {
		idStr := generateNodeId()
		id = &idStr
	}

	x := big.Int{}
	x.SetString(*id, 16)
	idBytes := x.Bytes()

	externalNode := new(ExternalNode)
	externalNode.nodeId = idBytes
	externalNode.ip = ip
	externalNode.port = port

	newNode := new(Node)
	newNode.nodeId = idBytes
	newNode.ip = ip
	newNode.port = port
	newNode.predecessor = externalNode

	for i := 0; i < num_bits; i++ {
		_, start := calcFinger(idBytes, i+1, num_bits)
		newNode.finger[i].start = start
	}
	newNode.finger[0].node = externalNode

	newNode.keys = make(map[string]string)

	return newNode
}

// return closest finger preceding id
func (this *Node) closestPrecedingFinger(id []byte) *ExternalNode {
	for i := num_bits - 1; i >= 0; i-- {
		if this.finger[i].node != nil && between3(this.nodeId, id, this.finger[i].node.nodeId) {
			return this.finger[i].node
		}
	}

	extNode := new(ExternalNode)
	extNode.nodeId = this.nodeId
	extNode.ip = this.ip
	extNode.port = this.port
	return extNode
}

// np thinks it might be our predecessor.
func (this *Node) notify(np *ExternalNode) {
	if this.predecessor == nil || between3(this.predecessor.nodeId, this.nodeId, np.nodeId) {
		this.predecessor = np
	}
}

//func (this *Node) printRing() {
//	this.printNode()

//	for i := this.finger[0].node; i != this; i = i.finger[0].node {
//		i.printNode()
//	}
//}

//func (this *Node) printNode() {
//	fmt.Println("Node "+":", hex.EncodeToString(this.nodeId))
//	if this.finger[0].node != nil {
//		fmt.Println("Successor "+":", hex.EncodeToString(this.finger[0].node.nodeId))
//	}
//	if this.predecessor != nil {
//		fmt.Println("Predecessors"+":", hex.EncodeToString(this.predecessor.nodeId))
//	}
//	//fmt.Println("Fingers")
//	//for i := 0; i < num_bits; i++ {
//	//	if this.finger[i].node != nil {
//	//		fmt.Println(strconv.Itoa(i)+" "+"Start:", hex.EncodeToString(this.finger[i].start), "Id:", hex.EncodeToString(this.finger[i].node.nodeId))
//	//	}
//	//}
//	fmt.Println("")
//}

//func (this *Node) testCalcFingers(k int, m int) {
//	fmt.Println("calulcating result = (n+2^(k-1)) mod (2^m)")

//	// convert the n to a bigint
//	nBigInt := big.Int{}
//	nBigInt.SetBytes(this.nodeId)
//	//fmt.Printf("n            %s\n", nBigInt.String())
//	fmt.Printf("n            %s\n", hex.EncodeToString(this.nodeId))

//	fmt.Printf("k            %d\n", k)

//	fmt.Printf("m            %d\n", m)

//	// get the right addend, i.e. 2^(k-1)
//	two := big.NewInt(2)
//	addend := big.Int{}
//	addend.Exp(two, big.NewInt(int64(k-1)), nil)

//	fmt.Printf("2^(k-1)      %s\n", addend.String())

//	// calculate sum
//	sum := big.Int{}
//	sum.Add(&nBigInt, &addend)

//	fmt.Printf("(n+2^(k-1))  %s\n", sum.String())

//	// calculate 2^m
//	ceil := big.Int{}
//	ceil.Exp(two, big.NewInt(int64(m)), nil)

//	fmt.Printf("2^m          %s\n", ceil.String())

//	// apply the mod
//	result := big.Int{}
//	result.Mod(&sum, &ceil)

//	fmt.Printf("finger       %s\n", result.String())

//	resultBytes := result.Bytes()
//	resultHex := fmt.Sprintf("%x", resultBytes)

//	fmt.Printf("finger (hex) %s\n", resultHex)

//	fmt.Println("successor   ", hex.EncodeToString(this.findSuccessor(resultBytes).nodeId))

//	dist := distance(this.nodeId, resultBytes, num_bits)

//	fmt.Println("distance     " + dist.String())
//}
