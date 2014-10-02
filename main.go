package main

import (
	"fmt"
	"github.com/liamzebedee/go-qrp"
	"os"
	"time"
)

type App struct {
	transport Transport
	node      *Node
	nodeUDP   *qrp.Node
}

type AddService struct {
	app *App
}

type AddArgs struct {
	Id       []byte
	Ip, Port string
}
type AddReply struct {
	Id       []byte
	Ip, Port string
}

func (s *AddService) Join(args *AddArgs, reply *AddReply) {
	fmt.Println("Join")
	fmt.Println("args Id: ", args.Id)
	fmt.Println("args Ip: ", args.Ip)
	fmt.Println("args Port: ", args.Port)
	reply.Id = s.app.node.nodeId
	reply.Ip = s.app.node.ip
	reply.Port = s.app.node.port
}

func (this *App) init(bindAddr, bindPort string) {
	this.node = makeDHTNode(nil, bindAddr, bindPort)
	this.transport = Transport{bindAddr + ":" + bindPort}

	node, err := qrp.CreateNodeUDP("udp", bindAddr+":"+bindPort, 512)
	if err != nil {
		fmt.Print("ERROR: Can't create node -", err.Error())
		return
	}
	this.nodeUDP = node
	fmt.Println("Node created")

	join := new(AddService)
	join.app = this
	node.Register(join)
	fmt.Println("Join service registered")

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

//Tries to join the node at the specified address.
func (this *App) join(addr string) {
	//msg := Msg{"JOIN", "<Key>", this.transport.bindAddress, addr}
	//this.transport.send(&msg)

	args := new(AddArgs)
	args.Id = this.node.nodeId
	args.Ip = this.node.ip
	args.Port = this.node.port

	reply := new(AddReply)

	fmt.Println("Calling Join on server")
	err := this.nodeUDP.CallUDP("Join", addr, args, reply, 3)

	if err != nil {
		fmt.Print("Call error - ")
		fmt.Println(err.Error())
		return
	}

	if reply != nil {
		fmt.Println("reply Id: ", reply.Id)
		fmt.Println("reply Ip: ", reply.Ip)
		fmt.Println("reply Port: ", reply.Port)
	}
}

//func (this *App) listen() {
//	this.transport.listen()
//}

func main() {
	app := App{}

	if len(os.Args) > 1 {
		if os.Args[1] == "server" {
			app.init("127.0.0.1", "13337")
			//app.listen()

			for {
				err := app.nodeUDP.ListenAndServe()
				if err != nil {
					fmt.Println("Error serving -", err.Error())
					return
				}
			}

		} else if os.Args[1] == "client" {
			app.init("127.0.0.1", "13338")

			go func() {
				err := app.nodeUDP.ListenAndServe()
				if err != nil {
					fmt.Println("Error serving -", err.Error())
					return
				}
			}()
			// run Join on the server
			app.join("127.0.0.1:13337")
		}
	}
}
