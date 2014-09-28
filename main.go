package main

import (
	"fmt"
	"time"
	"os"
)

type App struct {
	transport 	Transport
	node 		*Node

} 

func (this *App) init(bindAddr, bindPort string) {
	this.node = makeDHTNode(nil, bindAddr, bindPort)
	this.transport = Transport{bindAddr+":"+bindPort}

	// call stabilize and fixFingers periodically
	go func() {
		c := time.Tick(3 * time.Second)
		for now := range c {
			fmt.Println(now)
			this.node.stabilize()
			this.node.fixFingers()
		}
	}()
}

// Tries to join the node at the specified address.
func (this *App) join(addr string) {
	msg := Msg{"JOIN", "<Key>", this.transport.bindAddress, addr}
	this.transport.send(&msg)
}

func (this *App) listen() {
	this.transport.listen()
}

func main() {
	app := App{}

	if len(os.Args) > 1 {
		if os.Args[1] == "server" {
			app.init("127.0.0.1", "13337")
			app.listen()
		} else if os.Args[1] == "client" {
			app.init("127.0.0.1", "13338")
			app.join("127.0.0.1:13337")
		}

	}

}
