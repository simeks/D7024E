package main

import (
	"fmt"
	"os"
)



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
