package main

import (
	"net/http"
	"os"
	"strings"
)

func main() {
	app := App{}

	localAddr := "127.0.0.1:13337"
	remote := ""
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-l" && (i+1) < len(os.Args) { // Port <port number>
			localAddr = os.Args[i+1]
			i++
		} else if os.Args[i] == "-j" && (i+1) < len(os.Args) { // Join <remote host>
			remote = os.Args[i+1]
			i++
		}
	}

	app.init(localAddr)

	if remote != "" { // Join existing ring
		go app.join(remote)
	}


	go func() {
		http.HandleFunc("/chord/", chordHandler)
		http.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
			postHandler(w, r, &app)
		})
		http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
			deleteHandler(w, r, &app)
		})
		http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
			getHandler(w, r, &app)
		})
		http.HandleFunc("/put/", func(w http.ResponseWriter, r *http.Request) {
			putHandler(w, r, &app)
		})
		http.ListenAndServe(":"+strings.Split(localAddr, ":")[1], nil)
	}()
	

	app.listen()

}
