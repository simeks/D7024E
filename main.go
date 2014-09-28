package main

import (
	"fmt"
	"time"
)

func main() {

	node := makeDHTNode(nil, "localhost", "1111")

	//node.addToRing(...)

	// call stabilize and fixFingers periodically
	go func() {
		c := time.Tick(3 * time.Second)
		for now := range c {
			fmt.Println(now)
			node.stabilize()
			node.fixFingers()
		}
	}()

	// used to see the output of fmt.Println() from inside the goroutines
	var input string
	fmt.Scanln(&input)
}
