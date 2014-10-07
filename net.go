package main

import (
	"encoding/json"
	"fmt"
)


type NotifyMsg struct {
	NodeId []byte
	Ip     string
	Port   string
}

// insertKey, updateKey
type KeyValueMsg struct {
	Key, Value string
}

// getKey, deleteKey
type KeyMsg struct {
	Key string
}

type ValueMsg struct {
	Value string
}

type DeleteValueReply struct {
	Deleted bool
}
type UpdateValueReply struct {
	Updated bool
}


type JoinRequest struct {
	Ip string
	Port string
	Id []byte
}

type JoinReply struct {
	Id []byte
}

// findSuccessor, findPredecessor
type FindNodeReq struct {
	Id []byte
}

// findSuccessor, findPredecessor, getSuccessor, getPredecessor
type FindNodeReply struct {
	Id []byte
	Ip string
	Port string
	Found bool
}

type Net struct {
	app *App
}


func (n *Net) notify(msg Msg) {
	m := NotifyMsg{}
	json.Unmarshal(msg.Data, &m)

	node := ExternalNode{m.NodeId, m.Ip, m.Port}
	n.app.node.notify(&node)
}

func (n *Net) insertKey(msg Msg) {
	m := KeyValueMsg{}
	json.Unmarshal(msg.Data, &m)

	n.app.keyValue[m.Key] = m.Value
	
}

func (n *Net) getKey(req Request) {
	r := KeyMsg{}
	json.Unmarshal(req.Data, &r)

	reply := ValueMsg{}

	val, ok := n.app.keyValue[r.Key]
	if ok {
		reply.Value = val
	}

	bytes, _ := json.Marshal(reply)
	n.app.transport.sendReply(req.SN, bytes)
}

func (n *Net) deleteKey(req Request) {
	r := KeyMsg{}
	json.Unmarshal(req.Data, &r)
	
	reply := DeleteValueReply{}

	_, ok := n.app.keyValue[r.Key]
	if ok {
		delete(n.app.keyValue, r.Key)
		reply.Deleted = true
	} else {
		reply.Deleted = false
	}	

	bytes, _ := json.Marshal(reply)
	n.app.transport.sendReply(req.SN, bytes)
}

func (n *Net) updateKey(req Request) {
	r := KeyValueMsg{}
	json.Unmarshal(req.Data, &r)
	
	reply := UpdateValueReply{}	

	_, ok := n.app.keyValue[r.Key]
	if ok {
		n.app.keyValue[r.Key] = r.Value
		reply.Updated = true
	} else {
		reply.Updated = false
	}

	bytes, _ := json.Marshal(reply)
	n.app.transport.sendReply(req.SN, bytes)
}



func (n *Net) ping(req Request) {
	n.app.transport.sendReply(req.SN, []byte{})
}

func (n *Net) join(req Request) {
	r := JoinRequest{}
	json.Unmarshal(req.Data, &r)

	fmt.Println("Received a Join message from ", r.Ip+r.Port, "\n")

	reply := JoinReply{}
	reply.Id = n.app.node.nodeId
	bytes, _ := json.Marshal(reply)
	n.app.transport.sendReply(req.SN, bytes)
}

func (n *Net) findSuccessor(req Request) {
	r := FindNodeReq{}
	json.Unmarshal(req.Data, &r)

	successor := n.app.findSuccessor(r.Id)

	reply := FindNodeReply{}
	if successor != nil {
		reply.Id = successor.nodeId
		reply.Ip = successor.ip
		reply.Port = successor.port
		reply.Found = true
	} else {
		reply.Found = false
	}

	bytes, _ := json.Marshal(reply)
	n.app.transport.sendReply(req.SN, bytes)
}

func (n *Net) findPredecessor(req Request) {
	r := FindNodeReq{}
	json.Unmarshal(req.Data, &r)

	predecessor := n.app.findPredecessor(r.Id)

	reply := FindNodeReply{}
	if predecessor != nil {
		reply.Id = predecessor.nodeId
		reply.Ip = predecessor.ip
		reply.Port = predecessor.port
		reply.Found = true
	} else {
		reply.Found = false
	}

	bytes, _ := json.Marshal(reply)
	n.app.transport.sendReply(req.SN, bytes)
}

func (n *Net) getSuccessor(req Request) {
	successor := n.app.node.finger[0].node

	reply := FindNodeReply{}
	if successor != nil {
		reply.Id = successor.nodeId
		reply.Ip = successor.ip
		reply.Port = successor.port
		reply.Found = true
	} else {
		reply.Found = false
	}

	bytes, _ := json.Marshal(reply)
	n.app.transport.sendReply(req.SN, bytes)
}

func (n *Net) getPredecessor(req Request) {
	predecessor := n.app.node.predecessor

	reply := FindNodeReply{}
	if predecessor != nil {
		reply.Id = predecessor.nodeId
		reply.Ip = predecessor.ip
		reply.Port = predecessor.port
		reply.Found = true
	} else {
		reply.Found = false
	}

	bytes, _ := json.Marshal(reply)
	n.app.transport.sendReply(req.SN, bytes)
}

