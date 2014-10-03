package main

import (
	//"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
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

	return newNode
}

// this joins the network;
// node is an arbitrary node in the network
//func (this *Node) addToRing(np *Node) {
//	this.mutex.Lock()
//	this.predecessor = nil
//	this.finger[0].node = np.findSuccessor(this.nodeId)
//	this.mutex.Unlock()
//}

func (this *App) addToRing(np *ExternalNode) {
	this.node.mutex.Lock()
	this.node.predecessor = nil
	this.node.mutex.Unlock()

	args := new(AddArgs)
	args.Id = this.node.nodeId
	args.Ip = this.node.ip
	args.Port = this.node.port

	reply := new(AddReply)

	addr := np.ip + ":" + np.port

	// call FindSuccessor on np, which is already in the ring
	err := this.nodeUDP.CallUDP("FindSuccessor", addr, args, reply, 3)

	if err != nil {
		fmt.Print("Call error - ")
		fmt.Println(err.Error())
		return
	}

	if reply != nil {
		extNode := new(ExternalNode)
		extNode.nodeId = reply.Id
		extNode.ip = reply.Ip
		extNode.port = reply.Port

		this.node.mutex.Lock()
		this.node.finger[0].node = extNode
		this.node.mutex.Unlock()
	}
}

// ask node n to find id's successor
//func (this *Node) findSuccessor(id []byte) *Node {
//	np := this.findPredecessor(id)
//	return np.finger[0].node
//}

func (this *App) findSuccessor(id []byte) *ExternalNode {
	np := this.findPredecessor(id)

	args := new(AddArgs)
	reply := new(AddReply)

	addr := np.ip + ":" + np.port

	// call GetSuccessor on np
	err := this.nodeUDP.CallUDP("GetSuccessor", addr, args, reply, 3)

	if err != nil {
		fmt.Print("Call error - ")
		fmt.Println(err.Error())
		return nil
	}

	if reply != nil {
		// now we have the successor
		successor := new(ExternalNode)
		successor.nodeId = reply.Id
		successor.ip = reply.Ip
		successor.port = reply.Port
		return successor
	}
	return nil
}

// ask node n to find id's predecessor
//func (this *Node) findPredecessor(id []byte) *Node {
//	np := this

//	for between(np.nodeId, np.finger[0].node.nodeId, id) == false {
//		np = np.closestPrecedingFinger(id)
//	}
//	return np
//}

func (this *App) findPredecessor(id []byte) *ExternalNode {
	if between(this.node.nodeId, this.node.finger[0].node.nodeId, id) {
		extNode := new(ExternalNode)
		extNode.nodeId = this.node.nodeId
		extNode.ip = this.node.ip
		extNode.port = this.node.port
		return extNode
	} else {
		np := this.node.closestPrecedingFinger(id)

		args := new(AddArgs)
		args.Id = id

		reply := new(AddReply)

		addr := np.ip + ":" + np.port

		// call FindPredecessor on np
		err := this.nodeUDP.CallUDP("FindPredecessor", addr, args, reply, 3)

		if err != nil {
			fmt.Print("Call error - ")
			fmt.Println(err.Error())
			return nil
		}
		return nil
	}
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

//func (this *Node) lookup(key string) *Node {
//	id := big.Int{}
//	id.SetString(key, 16)
//	idBytes := id.Bytes()

//	return this.findSuccessor(idBytes)
//}

func (this *App) lookup(key string) *ExternalNode {
	id := big.Int{}
	id.SetString(key, 16)
	idBytes := id.Bytes()

	return this.findSuccessor(idBytes)
}

// periodically verify n’s immediate successor,
// and tell the successor about n.
//func (this *Node) stabilize() {
//	this.mutex.Lock()
//	x := this.finger[0].node.predecessor
//	if x != nil && between3(this.nodeId, this.finger[0].node.nodeId, x.nodeId) {
//		this.finger[0].node = x
//	}
//	this.finger[0].node.notify(this)
//	this.mutex.Unlock()
//}

func (this *App) stabilize() {

	args := new(AddArgs)
	reply := new(AddReply)

	addr := this.node.finger[0].node.ip + ":" + this.node.finger[0].node.port

	// call GetPredecessor on this.node's successor
	err := this.nodeUDP.CallUDP("GetPredecessor", addr, args, reply, 3)

	if err != nil {
		fmt.Print("Call error - ")
		fmt.Println(err.Error())
		return
	}

	if reply != nil {
		// now we have the predecessor
		predecessor := new(ExternalNode)
		predecessor.nodeId = reply.Id
		predecessor.ip = reply.Ip
		predecessor.port = reply.Port

		if predecessor != nil && between3(this.node.nodeId, this.node.finger[0].node.nodeId, predecessor.nodeId) {
			this.node.mutex.Lock()
			this.node.finger[0].node = predecessor
			this.node.mutex.Unlock()
		}

		addr := this.node.finger[0].node.ip + ":" + this.node.finger[0].node.port

		args.Id = this.node.nodeId
		args.Ip = this.node.ip
		args.Port = this.node.port

		err := this.nodeUDP.CallUDP("Notify", addr, args, reply, 3)
		if err != nil {
			fmt.Print("Call error - ")
			fmt.Println(err.Error())
			return
		}
	}
}

// np thinks it might be our predecessor.
func (this *Node) notify(np *ExternalNode) {
	if this.predecessor == nil || between3(this.predecessor.nodeId, this.nodeId, np.nodeId) {
		this.predecessor = np
	}
}

func (this *App) fixFingers() {
	i := rand.Intn(num_bits)
	successor := this.findSuccessor(this.node.finger[i].start)
	this.node.mutex.Lock()
	this.node.finger[i].node = successor
	this.node.mutex.Unlock()
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
