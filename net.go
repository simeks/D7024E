package main

import (
	"encoding/json"
	"fmt"
)

type NotifyMsg struct {
	NodeId []byte
	Addr   string
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

type KeyValueExistsReply struct {
	Exists bool
}

type UpdateValueReply struct {
	Updated bool
}

type JoinRequest struct {
	Addr string
	Id   []byte
}

type JoinReply struct {
	Id []byte
}

// findSuccessor, findPredecessor, closestPrecedingFinger
type FindNodeReq struct {
	Id []byte
}

// findSuccessor, findPredecessor, getSuccessor, getPredecessor
type FindNodeReply struct {
	Id    []byte
	Addr  string
	Found bool
}

type SuccessorListReply struct {
	Id   [num_successors][]byte
	Addr [num_successors]string
}

type TransferDataMsg struct {
	KeyValue map[string]string
}

type Net struct {
	app *App
}

func (n *Net) notify(msg *Msg) {
	m := NotifyMsg{}
	json.Unmarshal(msg.Data, &m)

	node := ExternalNode{m.NodeId, m.Addr}
	n.app.notify(&node)
}

func (n *Net) insertKey(msg *Msg) {
	m := KeyValueMsg{}
	json.Unmarshal(msg.Data, &m)

	n.app.node.mutex.Lock()
	n.app.keyValue[m.Key] = m.Value
	n.app.node.mutex.Unlock()

}

func (n *Net) getKey(rc *RequestContext) {
	r := KeyMsg{}
	json.Unmarshal(rc.req.Data, &r)

	reply := ValueMsg{}

	val, ok := n.app.keyValue[r.Key]
	if ok {
		reply.Value = val
	} else {
		reply.Value = ""
	}

	bytes, _ := json.Marshal(reply)
	rc.replyChan <- bytes
}

func (n *Net) deleteKey(rc *RequestContext) {
	r := KeyMsg{}
	json.Unmarshal(rc.req.Data, &r)

	reply := DeleteValueReply{}

	_, ok := n.app.keyValue[r.Key]
	if ok {
		delete(n.app.keyValue, r.Key)
		reply.Deleted = true
	} else {
		reply.Deleted = false
	}

	bytes, _ := json.Marshal(reply)
	rc.replyChan <- bytes
}

func (n *Net) keyValueExists(rc *RequestContext) {
	r := KeyValueMsg{}
	json.Unmarshal(rc.req.Data, &r)

	reply := KeyValueExistsReply{}

	_, ok := n.app.keyValue[r.Key]
	if ok {
		reply.Exists = true
	} else {
		reply.Exists = false
	}

	bytes, _ := json.Marshal(reply)
	rc.replyChan <- bytes
}

func (n *Net) updateKey(rc *RequestContext) {
	r := KeyValueMsg{}
	json.Unmarshal(rc.req.Data, &r)

	reply := UpdateValueReply{}

	_, ok := n.app.keyValue[r.Key]
	if ok {
		n.app.node.mutex.Lock()
		n.app.keyValue[r.Key] = r.Value
		n.app.node.mutex.Unlock()
		reply.Updated = true
	} else {
		reply.Updated = false
	}

	bytes, _ := json.Marshal(reply)
	rc.replyChan <- bytes
}

func (n *Net) ping(rc *RequestContext) {
	rc.replyChan <- []byte{}
}

func (n *Net) join(rc *RequestContext) {
	r := JoinRequest{}
	json.Unmarshal(rc.req.Data, &r)

	fmt.Println("Received a Join message from ", r.Addr, "\n")

	reply := JoinReply{}
	reply.Id = n.app.node.nodeId
	bytes, _ := json.Marshal(reply)
	rc.replyChan <- bytes
}

func (n *Net) findSuccessor(rc *RequestContext) {
	r := FindNodeReq{}
	json.Unmarshal(rc.req.Data, &r)

	successor := n.app.findSuccessor(r.Id)

	reply := FindNodeReply{}
	if successor != nil {
		reply.Id = successor.nodeId
		reply.Addr = successor.addr
		reply.Found = true
	} else {
		reply.Found = false
	}

	bytes, _ := json.Marshal(reply)
	rc.replyChan <- bytes
}

func (n *Net) findPredecessor(rc *RequestContext) {
	r := FindNodeReq{}
	json.Unmarshal(rc.req.Data, &r)

	predecessor := n.app.findPredecessor(r.Id)

	reply := FindNodeReply{}
	if predecessor != nil {
		reply.Id = predecessor.nodeId
		reply.Addr = predecessor.addr
		reply.Found = true
	} else {
		reply.Found = false
	}

	bytes, _ := json.Marshal(reply)
	rc.replyChan <- bytes
}

func (n *Net) getSuccessor(rc *RequestContext) {
	successor := n.app.node.finger[0].node

	reply := FindNodeReply{}
	if successor != nil {
		reply.Id = successor.nodeId
		reply.Addr = successor.addr
		reply.Found = true
	} else {
		reply.Found = false
	}

	bytes, _ := json.Marshal(reply)
	rc.replyChan <- bytes
}

func (n *Net) getSuccessorList(rc *RequestContext) {
	succList := n.app.node.successorList

	reply := SuccessorListReply{}
	for i := 0; i < num_successors; i++ {
		reply.Id[i] = succList[i].nodeId
		reply.Addr[i] = succList[i].addr
	}

	bytes, _ := json.Marshal(reply)
	rc.replyChan <- bytes
}

func (n *Net) getPredecessor(rc *RequestContext) {
	predecessor := n.app.node.predecessor

	reply := FindNodeReply{}
	if predecessor != nil {
		reply.Id = predecessor.nodeId
		reply.Addr = predecessor.addr
		reply.Found = true
	} else {
		reply.Found = false
	}

	bytes, _ := json.Marshal(reply)
	rc.replyChan <- bytes
}

func (n *Net) closestPrecedingFinger(rc *RequestContext) {
	r := FindNodeReq{}
	json.Unmarshal(rc.req.Data, &r)

	node := n.app.node.closestPrecedingFinger(r.Id)

	reply := FindNodeReply{}

	if node != nil {
		reply.Id = node.nodeId
		reply.Addr = node.addr
		reply.Found = true
	} else {
		reply.Found = false
	}

	bytes, _ := json.Marshal(reply)
	rc.replyChan <- bytes
}

func (n *Net) transferData(msg *Msg) {
	m := TransferDataMsg{}
	json.Unmarshal(msg.Data, &m)

	n.app.transferData(&m.KeyValue)
}
