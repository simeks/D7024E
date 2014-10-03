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
	if args.Id != nil && predecessor != nil {
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
