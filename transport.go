package main

import (
	"fmt"
)

func (s *AddService) Join(args *AddArgs, reply *AddReply) {
	fmt.Println("Received a Join message from ", args.Ip+args.Port)
	reply.Id = s.app.node.nodeId
	reply.Ip = s.app.node.ip
	reply.Port = s.app.node.port
}

func (s *AddService) FindSuccessor(args *AddArgs, reply *AddReply) {
	successor := s.app.findSuccessor(args.Id)
	reply.Id = successor.nodeId
	reply.Ip = successor.ip
	reply.Port = successor.port
}

func (s *AddService) FindPredecessor(args *AddArgs, reply *AddReply) {
	s.app.findPredecessor(args.Id)
}

func (s *AddService) GetSuccessor(args *AddArgs, reply *AddReply) {
	reply.Id = s.app.node.finger[0].node.nodeId
	reply.Ip = s.app.node.finger[0].node.ip
	reply.Port = s.app.node.finger[0].node.port
}
