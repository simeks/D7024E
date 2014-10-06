package main

import (
	"fmt"
)

func (s *AddService) Join(args *AddArgs, reply *AddReply) {
	fmt.Println("Received a Join message from ", args.Ip+args.Port, "\n")
	reply.Id = s.app.node.nodeId
	reply.Ip = s.app.node.ip
	reply.Port = s.app.node.port
}

func (s *AddService) FindSuccessor(args *AddArgs, reply *AddReply) {
	successor := s.app.findSuccessor(args.Id)

	if successor != nil {
		reply.Id = successor.nodeId
		reply.Ip = successor.ip
		reply.Port = successor.port
	}
}

func (s *AddService) FindPredecessor(args *AddArgs, reply *AddReply) {
	predecessor := s.app.findPredecessor(args.Id)
	if predecessor != nil {
		reply.Id = predecessor.nodeId
		reply.Ip = predecessor.ip
		reply.Port = predecessor.port
	}
}

func (s *AddService) GetSuccessor(args *AddArgs, reply *AddReply) {
	reply.Id = s.app.node.finger[0].node.nodeId
	reply.Ip = s.app.node.finger[0].node.ip
	reply.Port = s.app.node.finger[0].node.port
}

func (s *AddService) GetPredecessor(args *AddArgs, reply *AddReply) {
	if s.app.node.predecessor != nil {
		reply.Id = s.app.node.predecessor.nodeId
		reply.Ip = s.app.node.predecessor.ip
		reply.Port = s.app.node.predecessor.port
	}
}

func (s *AddService) Notify(args *AddArgs, reply *AddReply) {
	extNode := new(ExternalNode)
	extNode.nodeId = args.Id
	extNode.ip = args.Ip
	extNode.port = args.Port
	s.app.node.notify(extNode)
}

func (s *AddService) InsertKey(args *AddArgs, reply *AddReply) {
	s.app.node.mutex.Lock()
	defer s.app.node.mutex.Unlock()

	s.app.node.keys[args.Key] = args.Value
}

func (s *AddService) DeleteKey(args *AddArgs, reply *AddReply) {
	_, ok := s.app.node.keys[args.Key]
	if ok {
		s.app.node.mutex.Lock()
		defer s.app.node.mutex.Unlock()

		delete(s.app.node.keys, args.Key)
		reply.WasDeleted = 1
	} else {
		reply.WasDeleted = 0
	}
}

func (s *AddService) GetKey(args *AddArgs, reply *AddReply) {
	_, ok := s.app.node.keys[args.Key]
	if ok {
		reply.Value = s.app.node.keys[args.Key]
	}
}

func (s *AddService) UpdateKey(args *AddArgs, reply *AddReply) {
	_, ok := s.app.node.keys[args.Key]
	if ok {
		s.app.node.keys[args.Key] = args.Value
		reply.WasUpdated = 1
	} else {
		reply.WasUpdated = 0
	}
}

func (s *AddService) Ping(args *AddArgs, reply *AddReply) {

}
