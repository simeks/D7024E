package main

import (
	"fmt"
	"net"
	"encoding/json"
)

type Msg struct {
	Type 	string
	Key		string
	Src 	string
	Dst 	string
}

type Transport struct {
	bindAddress string
}

func (transport *Transport) listen() {
	udpAddr, _ := net.ResolveUDPAddr("udp", transport.bindAddress)
	conn, _ := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		msg := Msg{}
		dec.Decode(&msg)
		// we got a message
		// ...
		fmt.Println("msg:",msg)
	}
} 

func (transport *Transport) send(msg *Msg) {
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)

	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()

	bytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println("Msg:",msg,"Sent:",string(bytes))
}
