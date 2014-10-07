package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"time"
	"sync"
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

	lock sync.Mutex
}

func (t *Transport) init(bindAddr string) {
	t.bindAddress = bindAddr
	t.lastReqSN = 0
	t.sentRequests = make(map[int]*PendingRequest)
	t.recvRequests = make(map[int]*RecvRequest)
}

func (t *Transport) listen(msgChan chan Msg, reqChan chan Request) {
	udpAddr, _ := net.ResolveUDPAddr("udp", t.bindAddress)
	conn, _ := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	for {
		packet := Packet{}
		dec.Decode(&packet)

		if packet.Msg != nil {
			msgChan <- *packet.Msg
		} else if packet.Request != nil {
			t.lock.Lock()
			recvReq := RecvRequest{packet.Src}
			t.recvRequests[packet.Request.SN] = &recvReq
			t.lock.Unlock()

			reqChan <- *packet.Request
		} else if packet.Reply != nil {
			t.processReply(packet.Reply)
		}

	}
} 

func (t *Transport) processReply(reply *Reply) {
	t.lock.Lock()
	defer t.lock.Unlock()

	pr := t.sentRequests[reply.SN]
	pr.channel <- *reply

	delete(t.sentRequests, reply.SN)
}

func (t *Transport) sendPacket(dst string, packet Packet) {
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

func (t *Transport) sendMsg(dst, msgId string, data []byte) {
	packet := Packet{t.bindAddress, &Msg{msgId, data}, nil, nil}
	t.sendPacket(dst, packet)
}

// Blocks until either a reply is retrieved or the request times out
func (t *Transport) sendRequest(dst, msgId string, data []byte) *Reply {
	t.lock.Lock()

	packet := Packet{t.bindAddress, nil, &Request{t.lastReqSN, msgId, data}, nil}
	pr := PendingRequest{t.lastReqSN, make(chan Reply)}
	t.sentRequests[t.lastReqSN] = &pr
	t.lastReqSN++

	t.lock.Unlock()

	t.sendPacket(dst, packet)

	timeout := time.NewTimer(time.Second * 3)
	select {
		case <- timeout.C:
			return nil
		case reply := <- pr.channel:
			return &reply
	}
}

func (t *Transport) sendReply(sn int, data []byte) {
	t.lock.Lock()
	rr := t.recvRequests[sn]
	delete(t.recvRequests, sn)
	t.lock.Unlock()

	packet := Packet{t.bindAddress, nil, nil, &Reply{sn, data}}
	t.sendPacket(rr.src, packet)
}
