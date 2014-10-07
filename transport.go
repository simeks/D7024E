package main

import (
	"fmt"
	"net"
	"encoding/gob"
)

type Msg struct {
	Id string
	Data []byte
}

type Request struct {
	SN int
	Id string
	Data []byte
}

type Reply struct {
	SN int
	Data []byte
}


type Packet struct {
	Msg *Msg
	Request *Request
	Reply *Reply
}

type Transport struct {
	bindAddress string
	lastReqSN int
}

func (transport *Transport) listen(msgChan chan Msg, reqChan chan Request) {
	udpAddr, _ := net.ResolveUDPAddr("udp", transport.bindAddress)
	conn, _ := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	for {
		packet := Packet{}
		dec.Decode(&packet)

		if packet.Msg != nil {
			msgChan <- *packet.Msg

		} else if packet.Request != nil {

		} else if packet.Reply != nil {

		}

	}
} 

func (transport *Transport) sendPacket(dst string, packet Packet) {
	udpAddr, err := net.ResolveUDPAddr("udp", dst)

	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()

	enc := gob.NewEncoder(conn)

	err = enc.Encode(packet)
	if err != nil {
		fmt.Println(err)
		return
	}	
}

func (transport *Transport) sendMsg(dst, msgId string, data []byte) {
	packet := Packet{&Msg{msgId, data}, nil, nil}
	transport.sendPacket(dst, packet)
}

func (transport *Transport) sendRequest(dst, msgId string, data []byte) Reply {
	packet := Packet{nil, &Request{transport.lastReqId, msgId, data}, nil}
	transport.sendPacket(dst, packet)
	transport.lastReqId++
}

func (transport *Transport) sendReply(dst, msgId string, data []byte) {
	packet := Packet{nil, nil, &Reply{0, msgId, data}}
	transport.sendPacket(dst, packet)
}
