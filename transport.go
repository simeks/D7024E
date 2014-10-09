package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"time"
	"sync"
	"strings"
)

const time_out time.Duration = 3

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
	channel chan Reply
}

type RequestContext struct {
	req *Request
	replyChan chan []byte
}

type Transport struct {
	bindAddress string
	lastReqSN int

	sentRequests map[int]*PendingRequest

	lock sync.Mutex
}

func (t *Transport) init(bindAddr string) {
	t.bindAddress = bindAddr
	t.lastReqSN = 0
	t.sentRequests = make(map[int]*PendingRequest)
}

func (t *Transport) listen(msgChan chan *Msg, reqChan chan *RequestContext) {
	port := strings.Split(t.bindAddress, ":")[1]

	udpAddr, _ := net.ResolveUDPAddr("udp", ":"+port)
	conn, _ := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	for {
		packet := Packet{}
		dec.Decode(&packet)

		if packet.Msg != nil {
			msgChan <- packet.Msg
		} else if packet.Request != nil {
			rc := RequestContext{packet.Request, make(chan []byte)}
			go t.waitForReply(packet.Request.SN, packet.Src, rc.replyChan)

			reqChan <- &rc
		} else if packet.Reply != nil {
			t.processReply(packet.Reply)
		}

	}
}

func (t *Transport) waitForReply(sn int, src string, c chan []byte) {
	data := <- c

	packet := Packet{t.bindAddress, nil, nil, &Reply{sn, data}}
	t.sendPacket(src, packet)
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
	pr := PendingRequest{make(chan Reply)}

	t.sentRequests[t.lastReqSN] = &pr
	t.lastReqSN++

	t.lock.Unlock()
	t.sendPacket(dst, packet)

	timeout := time.NewTimer(time.Second * time_out)
	select {
		case <- timeout.C:
			return nil
		case reply := <- pr.channel:
			return &reply
	}
}
