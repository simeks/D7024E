package main

import (
	"net/http"
	"os"
)

func main() {
	app := App{}

	if len(os.Args) > 1 {
		if os.Args[1] == "server" {
			app.init("127.0.0.1", "13337")

			go func() {
				http.HandleFunc("/chord/", chordHandler)
				http.HandleFunc("/inserted/", func(w http.ResponseWriter, r *http.Request) {
					insertedHandler(w, r, &app)
				})
				http.ListenAndServe(":13337", nil)
			}()

			app.listen()

		} else if os.Args[1] == "client" {
			app.init("127.0.0.1", os.Args[2])

			go app.join("127.0.0.1:13337")

			go func() {
				http.HandleFunc("/chord/", chordHandler)
				http.HandleFunc("/inserted/", func(w http.ResponseWriter, r *http.Request) {
					insertedHandler(w, r, &app)
				})
				http.ListenAndServe(":"+os.Args[2], nil)
			}()

			app.listen()

		}
	}
}
