package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"time"
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
	Src string

	Msg *Msg
	Request *Request
	Reply *Reply
}

type PendingRequest struct {
	sn int
	channel chan Reply
}
type RecvRequest struct {
	src string
}

type Transport struct {
	bindAddress string
	lastReqSN int

	sentRequests map[int]*PendingRequest
	recvRequests map[int]*RecvRequest
}

func (transport *Transport) init(bindAddr string) {
	transport.bindAddress = bindAddr
	transport.lastReqSN = 0
	transport.sentRequests = make(map[int]*PendingRequest)
	transport.recvRequests = make(map[int]*RecvRequest)
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
			recvReq := RecvRequest{packet.Src}
			transport.recvRequests[packet.Request.SN] = &recvReq

			reqChan <- *packet.Request
		} else if packet.Reply != nil {
			transport.processReply(packet.Reply)
		}

	}
} 

func (transport *Transport) processReply(reply *Reply) {
	pr := transport.sentRequests[reply.SN]
	pr.channel <- *reply

	delete(transport.sentRequests, reply.SN)
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
	packet := Packet{transport.bindAddress, &Msg{msgId, data}, nil, nil}
	transport.sendPacket(dst, packet)
}

// Blocks until either a reply is retrieved or the request times out
func (transport *Transport) sendRequest(dst, msgId string, data []byte) *Reply {
	packet := Packet{transport.bindAddress, nil, &Request{transport.lastReqSN, msgId, data}, nil}
	transport.sendPacket(dst, packet)

	pr := PendingRequest{transport.lastReqSN, make(chan Reply)}
	transport.sentRequests[transport.lastReqSN] = &pr
	transport.lastReqSN++

	timeout := time.NewTimer(time.Second * 3)
	select {
		case <- timeout.C:
			return nil
		case reply := <- pr.channel:
			return &reply
	}
}

func (transport *Transport) sendReply(sn int, data []byte) {
	rr := transport.recvRequests[sn]
	delete(transport.recvRequests, sn)

	packet := Packet{transport.bindAddress, nil, nil, &Reply{sn, data}}
	transport.sendPacket(rr.src, packet)
}
