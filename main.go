package main

import (
	"encoding/hex"
	"fmt"
	"github.com/liamzebedee/go-qrp"
	"os"
	"time"
)

type App struct {
	node    *Node
	nodeUDP *qrp.Node
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

func (this *App) init(bindAddr, bindPort string) {
	this.node = makeDHTNode(nil, bindAddr, bindPort)

	node, err := qrp.CreateNodeUDP("udp", bindAddr+":"+bindPort, 512)
	if err != nil {
		fmt.Print("ERROR: Can't create node -", err.Error())
		return
	}
	this.nodeUDP = node
	fmt.Println("\n" + bindAddr + ":" + bindPort + ": Node " + hex.EncodeToString(this.node.nodeId) + " created\n")

	join := new(AddService)
	join.app = this
	node.Register(join)
	fmt.Println("Join service registered")

	findsuccessor := new(AddService)
	findsuccessor.app = this
	node.Register(findsuccessor)
	fmt.Println("FindSuccessor service registered")

	findpredecessor := new(AddService)
	findpredecessor.app = this
	node.Register(findpredecessor)
	fmt.Println("FindPredecessor service registered")

	getsuccessor := new(AddService)
	getsuccessor.app = this
	node.Register(getsuccessor)
	fmt.Println("GetSuccessor service registered")

	getpredecessor := new(AddService)
	getpredecessor.app = this
	node.Register(getpredecessor)
	fmt.Println("GetPredecessor service registered")

	notify := new(AddService)
	notify.app = this
	node.Register(notify)
	fmt.Println("Notify service registered")

	fmt.Println("")

	// call stabilize and fixFingers periodically
	go func() {
		c := time.Tick(3 * time.Second)
		for now := range c {
			fmt.Println(now)
			this.stabilize()
			this.fixFingers()

			fmt.Println("Successor: ", hex.EncodeToString(this.node.finger[0].node.nodeId))
			fmt.Println("")
		}
	}()
}

//Tries to join the node at the specified address.
func (this *App) join(addr string) {

	args := new(AddArgs)
	args.Id = this.node.nodeId
	args.Ip = this.node.ip
	args.Port = this.node.port

	reply := new(AddReply)

	// get a node that is already in the ring
	fmt.Println("Calling Join on ", addr)
	err := this.nodeUDP.CallUDP("Join", addr, args, reply, 3)

	if err != nil {
		fmt.Print("Call error - ")
		fmt.Println(err.Error())
		return
	}

	if reply != nil {
		extNode := new(ExternalNode)
		extNode.nodeId = reply.Id
		extNode.ip = reply.Ip
		extNode.port = reply.Port

		// extNode is already in the ring
		this.addToRing(extNode)
	}
}

func main() {
	app := App{}

	if len(os.Args) > 1 {
		if os.Args[1] == "server" {
			app.init("127.0.0.1", "13337")

			for {
				err := app.nodeUDP.ListenAndServe()
				if err != nil {
					fmt.Println("Error serving -", err.Error())
					return
				}
			}

		} else if os.Args[1] == "client" {
			app.init("127.0.0.1", "13338")

			// run Join on the server
			go app.join("127.0.0.1:13337")

			for {
				err := app.nodeUDP.ListenAndServe()
				if err != nil {
					fmt.Println("Error serving -", err.Error())
					return
				}
			}
		}
	}
}
