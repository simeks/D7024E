package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"
)

func chordHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Insert string</h1>"+
		"<form action=\"/inserted/\" method=\"POST\">"+
		"<textarea name=\"body\"></textarea><br>"+
		"<input type=\"submit\" value=\"Insert\">"+
		"</form>")
}

func insertedHandler(w http.ResponseWriter, r *http.Request, app *App) {
	value := r.FormValue("body")
	key := sha1hash(value)

	responsibleNode := app.lookup(key)

	fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
		"<p>Value: "+value+"</p>"+
		"<p>Key: "+key+"</p>"+
		"<p>Responsible node: "+hex.EncodeToString(responsibleNode.nodeId)+"</p>")

	if bytes.Compare(app.node.nodeId, responsibleNode.nodeId) == 0 {
		app.node.keys[key] = value
	} else {
		args := new(AddArgs)
		args.Key = key
		args.Value = value
		reply := new(AddReply)

		addr := responsibleNode.ip + ":" + responsibleNode.port
		err := app.nodeUDP.CallUDP("InsertKey", addr, args, reply, time_out)

		if err != nil {
			fmt.Print("Call error - ")
			fmt.Println(err.Error())
			return
		}
	}

}
