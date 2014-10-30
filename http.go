package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func chordHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Insert a new key/value pair</h1>"+
		"<form action=\"/post/\" method=\"POST\">"+
		"Value:"+
		"<textarea name=\"insertvalue\"></textarea><br>"+
		"Key:"+
		"<textarea name=\"insertkey\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>"+
		"<h1>Delete a key/value pair</h1>"+
		"<form action=\"/delete/\" method=\"POST\">"+
		"Key:"+
		"<textarea name=\"deletekey\"></textarea><br>"+
		"Ecryption key:"+
		"<textarea name=\"deleteencryptionkey\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>"+
		"<h1>Update value for key</h1>"+
		"<form action=\"/put/\" method=\"POST\">"+
		"Value:"+
		"<textarea name=\"updatevalue\"></textarea><br>"+
		"Key:"+
		"<textarea name=\"updatekey\"></textarea><br>"+
		"Encryption key:"+
		"<textarea name=\"updateencryptionkey\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>"+
		"<h1>Get the value for key</h1>"+
		"<form action=\"/get/\" method=\"POST\">"+
		"Key:"+
		"<textarea name=\"getkey\"></textarea><br>"+
		"Ecryption key:"+
		"<textarea name=\"getencryptionkey\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>")
}

func postHandler(w http.ResponseWriter, r *http.Request, app *App) {
	value := r.FormValue("insertvalue")
	key := r.FormValue("insertkey")
	hashkey := sha1hash(key)

	// Encrypt the value
	secret, _ := GenerateAesSecret()
	value = EncryptAes(secret, value)

	responsibleNode := app.lookup(hashkey)

	if bytes.Compare(app.node.nodeId, responsibleNode.nodeId) == 0 {
		app.node.mutex.Lock()
		defer app.node.mutex.Unlock()

		_, ok := app.keyValue[hashkey]

		if ok { // if the key already exist
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>That key is already taken. Please choose a new key.</p>")
		} else { // insert the new key/value-pair
			app.keyValue[hashkey] = value
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Key/value pair inserted successfully!</p>"+
				"<p>Your encryption key is: "+secret+".")
		}

	} else {
		req := KeyValueMsg{}
		req.Key = hashkey
		req.Value = value

		bytes, _ := json.Marshal(req)

		r := app.transport.sendRequest(responsibleNode.addr, "keyValueExists", bytes)

		if r == nil {
			fmt.Println("Call error (keyValueExists)")
			return
		}

		if r != nil {
			reply := KeyValueExistsReply{}
			json.Unmarshal(r.Data, &reply)

			if reply.Exists { // key exists, dont insert it
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>not ins</p>")
			} else {
				app.transport.sendMsg(responsibleNode.addr, "insertKey", bytes)
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>ins</p>")
			}
		}
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request, app *App) {
	key := r.FormValue("deletekey")
	hashkey := sha1hash(key)

	responsibleNode := app.lookup(hashkey)

	if bytes.Compare(app.node.nodeId, responsibleNode.nodeId) == 0 {
		_, ok := app.keyValue[hashkey]
		if ok {
			app.node.mutex.Lock()
			defer app.node.mutex.Unlock()

			delete(app.keyValue, hashkey)
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Key deleted successfully!</p>")
		} else {
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Key was not found.</p>")
		}
	} else {
		req := KeyMsg{}
		req.Key = hashkey

		bytes, _ := json.Marshal(req)
		r := app.transport.sendRequest(responsibleNode.addr, "deleteKey", bytes)

		if r == nil {
			fmt.Println("Call error (deleteKey)")
			return
		}

		if r != nil {
			reply := DeleteValueReply{}
			json.Unmarshal(r.Data, &reply)

			if reply.Deleted {
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
	//decryptionkey := r.FormValue("getdecryptionkey")
	hashkey := sha1hash(key)

	responsibleNode := app.lookup(hashkey)

	if bytes.Compare(app.node.nodeId, responsibleNode.nodeId) == 0 {
		_, ok := app.keyValue[hashkey]
		if ok {
			// decrypta app.keyValue[hashkey]
			// ...
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Value: "+app.keyValue[hashkey]+"</p>")
		} else {
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Key was not found.</p>")
		}
	} else {
		req := KeyMsg{}
		req.Key = hashkey

		bytes, _ := json.Marshal(req)
		r := app.transport.sendRequest(responsibleNode.addr, "getKey", bytes)

		if r == nil {
			fmt.Println("Call error (getKey)")
			return
		}

		if r != nil {
			reply := ValueMsg{}
			json.Unmarshal(r.Data, &reply)

			if reply.Value != "" {
				// decrypta reply.Value
				// ...
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>Value: "+reply.Value+"</p>")
				return
			}
		}
		fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
			"<p>Key was not found.</p>")

	}
}

func putHandler(w http.ResponseWriter, r *http.Request, app *App) {
	value := r.FormValue("updatevalue")
	key := r.FormValue("updatekey")
	//encryptionkey := r.FormValue("updateencryptionkey")
	hashkey := sha1hash(key)

	responsibleNode := app.lookup(hashkey)

	if bytes.Compare(app.node.nodeId, responsibleNode.nodeId) == 0 {
		_, ok := app.keyValue[hashkey]
		if ok {
			app.keyValue[hashkey] = value
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Value updated successfully!</p>")
		} else {
			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>Key was not found.</p>")
		}
	} else {
		req := KeyValueMsg{}
		req.Key = hashkey
		req.Value = value

		bytes, _ := json.Marshal(req)
		r := app.transport.sendRequest(responsibleNode.addr, "updateKey", bytes)

		if r == nil {
			fmt.Println("Call error (updateKey)")
			return
		}

		if r != nil {
			reply := UpdateValueReply{}
			json.Unmarshal(r.Data, &reply)

			if reply.Updated {
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>Value updated successfully!</p>")
			} else {
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>Key was not found.</p>")
			}
		}
	}
}
