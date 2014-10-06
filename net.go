package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func chordHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Insert a new key/value pair</h1>"+
		"<form action=\"/post/\" method=\"POST\">"+
		"value:"+
		"<textarea name=\"insertvalue\"></textarea><br>"+
		"key:"+
		"<textarea name=\"insertkey\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>"+
		"<h1>Delete a key/value pair</h1>"+
		"<form action=\"/delete/\" method=\"POST\">"+
		"key:"+
		"<textarea name=\"deletekey\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>"+
		"<h1>Update value for key</h1>"+
		"<form action=\"/put/\" method=\"POST\">"+
		"value:"+
		"<textarea name=\"updatevalue\"></textarea><br>"+
		"key:"+
		"<textarea name=\"updatekey\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>"+
		"<h1>Get the value for key</h1>"+
		"<form action=\"/get/\" method=\"POST\">"+
		"key:"+
		"<textarea name=\"getkey\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>")
}

func postHandler(w http.ResponseWriter, r *http.Request, app *App) {
	value := r.FormValue("insertvalue")
	key := r.FormValue("insertkey")
	hashkey := sha1hash(key)

	responsibleNode := app.lookup(hashkey)

	if bytes.Compare(app.node.nodeId, responsibleNode.nodeId) == 0 {
		app.node.mutex.Lock()
		defer app.node.mutex.Unlock()
		app.node.keys[hashkey] = value

	} else {
		args := new(AddArgs)
		args.Key = hashkey
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
	fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
		"<p>Key/value pair inserted successfully!</p>")

}

func deleteHandler(w http.ResponseWriter, r *http.Request, app *App) {
	key := r.FormValue("deletekey")
	hashkey := sha1hash(key)

	responsibleNode := app.lookup(hashkey)

	if bytes.Compare(app.node.nodeId, responsibleNode.nodeId) == 0 {
		_, ok := app.node.keys[hashkey]
		if ok {
			app.node.mutex.Lock()
			defer app.node.mutex.Unlock()

			delete(app.node.keys, hashkey)
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Key deleted successfully!</p>")
		} else {
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Key was not found.</p>")
		}
	} else {
		args := new(AddArgs)
		args.Key = hashkey
		reply := new(AddReply)

		addr := responsibleNode.ip + ":" + responsibleNode.port
		err := app.nodeUDP.CallUDP("DeleteKey", addr, args, reply, time_out)

		if err != nil {
			fmt.Print("Call error - ")
			fmt.Println(err.Error())
			return
		}

		if reply != nil {
			if reply.WasDeleted == 1 {
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>Key deleted successfully!</p>")
			} else {
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>Key was not found.</p>")
			}
		}
	}
}

func getHandler(w http.ResponseWriter, r *http.Request, app *App) {
	key := r.FormValue("getkey")
	hashkey := sha1hash(key)

	responsibleNode := app.lookup(hashkey)

	if bytes.Compare(app.node.nodeId, responsibleNode.nodeId) == 0 {
		_, ok := app.node.keys[hashkey]
		if ok {
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Value: "+app.node.keys[hashkey]+"</p>")
		} else {
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Key was not found.</p>")
		}
	} else {
		args := new(AddArgs)
		args.Key = hashkey
		reply := new(AddReply)

		addr := responsibleNode.ip + ":" + responsibleNode.port
		err := app.nodeUDP.CallUDP("GetKey", addr, args, reply, time_out)

		if err != nil {
			fmt.Print("Call error - ")
			fmt.Println(err.Error())
			return
		}

		if reply != nil {
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Value: "+reply.Value+"</p>")
		} else {
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Key was not found.</p>")
		}
	}
}

func putHandler(w http.ResponseWriter, r *http.Request, app *App) {
	value := r.FormValue("updatevalue")
	key := r.FormValue("updatekey")
	hashkey := sha1hash(key)

	responsibleNode := app.lookup(hashkey)

	if bytes.Compare(app.node.nodeId, responsibleNode.nodeId) == 0 {
		_, ok := app.node.keys[hashkey]
		if ok {
			app.node.keys[hashkey] = value
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Value updated successfully!</p>")
		} else {
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Key was not found.</p>")
		}
	} else {
		args := new(AddArgs)
		args.Key = hashkey
		args.Value = value
		reply := new(AddReply)

		addr := responsibleNode.ip + ":" + responsibleNode.port
		err := app.nodeUDP.CallUDP("UpdateKey", addr, args, reply, time_out)

		if err != nil {
			fmt.Print("Call error - ")
			fmt.Println(err.Error())
			return
		}

		if reply != nil {
			if reply.WasUpdated == 1 {
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>Value updated successfully!</p>")
			} else {
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>Key was not found.</p>")
			}
		}
	}
}
