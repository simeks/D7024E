package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	//"strconv"
	"sync"
)

const num_bits int = 160

type Finger struct {
	node  *Node
	start []byte
}

type Node struct {
	nodeId      []byte
	ip          string
	port        string
	finger      [num_bits]Finger
	predecessor *Node
	mutex       sync.Mutex
}

func makeDHTNode(id *string, ip string, port string) *Node {
	if id == nil {
		idStr := generateNodeId()
		id = &idStr
	}

	x := big.Int{}
	x.SetString(*id, 16)
	idBytes := x.Bytes()

	newNode := new(Node)
	newNode.nodeId = idBytes
	newNode.ip = ip
	newNode.port = port
	newNode.predecessor = newNode

	for i := 0; i < num_bits; i++ {
		_, start := calcFinger(idBytes, i+1, num_bits)
		newNode.finger[i].start = start
	}
	newNode.finger[0].node = newNode

	return newNode
}

// this joins the network;
// node is an arbitrary node in the network
func (this *Node) addToRing(np *Node) {
	this.mutex.Lock()
	this.predecessor = nil
	this.finger[0].node = np.findSuccessor(this.nodeId)
	this.mutex.Unlock()
}

// ask node n to find id's successor
func (this *Node) findSuccessor(id []byte) *Node {
	np := this.findPredecessor(id)
	return np.finger[0].node
}

// ask node n to find id's predecessor
func (this *Node) findPredecessor(id []byte) *Node {
	np := this

	for between(np.nodeId, np.finger[0].node.nodeId, id) == false {
		np = np.closestPrecedingFinger(id)
	}
	return np
}

// return closest finger preceding id
func (this *Node) closestPrecedingFinger(id []byte) *Node {
	for i := num_bits - 1; i >= 0; i-- {
		if this.finger[i].node != nil && between3(this.nodeId, id, this.finger[i].node.nodeId) {
			return this.finger[i].node
		}
	}
	return this
}

func (this *Node) lookup(key string) *Node {
	id := big.Int{}
	id.SetString(key, 16)
	idBytes := id.Bytes()

	return this.findSuccessor(idBytes)
}

// periodically verify n’s immediate successor,
// and tell the successor about n.
func (this *Node) stabilize() {
	this.mutex.Lock()
	x := this.finger[0].node.predecessor
	if x != nil && between3(this.nodeId, this.finger[0].node.nodeId, x.nodeId) {
		this.finger[0].node = x
	}
	this.finger[0].node.notify(this)
	this.mutex.Unlock()
}

// np thinks it might be our predecessor.
func (this *Node) notify(np *Node) {
	if this.predecessor == nil || between3(this.predecessor.nodeId, this.nodeId, np.nodeId) {
		this.predecessor = np
	}
}

func (this *Node) fixFingers() {
	this.mutex.Lock()
	i := rand.Intn(num_bits)
	this.finger[i].node = this.findSuccessor(this.finger[i].start)
	this.mutex.Unlock()
}

func (this *Node) printRing() {
	this.printNode()

	for i := this.finger[0].node; i != this; i = i.finger[0].node {
		i.printNode()
	}
}

func (this *Node) printNode() {
	fmt.Println("Node "+":", hex.EncodeToString(this.nodeId))
	if this.finger[0].node != nil {
		fmt.Println("Successor "+":", hex.EncodeToString(this.finger[0].node.nodeId))
	}
	if this.predecessor != nil {
		fmt.Println("Predecessors"+":", hex.EncodeToString(this.predecessor.nodeId))
	}
	//fmt.Println("Fingers")
	//for i := 0; i < num_bits; i++ {
	//	if this.finger[i].node != nil {
	//		fmt.Println(strconv.Itoa(i)+" "+"Start:", hex.EncodeToString(this.finger[i].start), "Id:", hex.EncodeToString(this.finger[i].node.nodeId))
	//	}
	//}
	fmt.Println("")
}

/*
 * Example of expected output.
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            0
 * m            160
 * 2^(k-1)      1
 * (n+2^(k-1))  682874255151879437996522856919401519827635625587
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996522856919401519827635625587
 * finger (hex) 779d240121ed6d5e8bd0cb6529b08e5c617b5e73
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0
 */
func (this *Node) testCalcFingers(k int, m int) {
	fmt.Println("calulcating result = (n+2^(k-1)) mod (2^m)")

	// convert the n to a bigint
	nBigInt := big.Int{}
	nBigInt.SetBytes(this.nodeId)
	//fmt.Printf("n            %s\n", nBigInt.String())
	fmt.Printf("n            %s\n", hex.EncodeToString(this.nodeId))

	fmt.Printf("k            %d\n", k)

	fmt.Printf("m            %d\n", m)

	// get the right addend, i.e. 2^(k-1)
	two := big.NewInt(2)
	addend := big.Int{}
	addend.Exp(two, big.NewInt(int64(k-1)), nil)

	fmt.Printf("2^(k-1)      %s\n", addend.String())

	// calculate sum
	sum := big.Int{}
	sum.Add(&nBigInt, &addend)

	fmt.Printf("(n+2^(k-1))  %s\n", sum.String())

	// calculate 2^m
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(m)), nil)

	fmt.Printf("2^m          %s\n", ceil.String())

	// apply the mod
	result := big.Int{}
	result.Mod(&sum, &ceil)

	fmt.Printf("finger       %s\n", result.String())

	resultBytes := result.Bytes()
	resultHex := fmt.Sprintf("%x", resultBytes)

	fmt.Printf("finger (hex) %s\n", resultHex)

	fmt.Println("successor   ", hex.EncodeToString(this.findSuccessor(resultBytes).nodeId))

	dist := distance(this.nodeId, resultBytes, num_bits)

	fmt.Println("distance     " + dist.String())
}
