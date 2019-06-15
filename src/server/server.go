// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Updated by Stephen Wolff, June 2019 for CHTest

// +build ignore

package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func chtest(w http.ResponseWriter, r *http.Request) {

	// upgrade request to websocket
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {

		// Replace _ with mt here, if returning writing messages / acks?
		_, messageBytes, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		/*
			Each message should be sent with an 8 byte UUID,
			this should consist of a 6 byte timestamp (seconds from 01/01/1970) and the 2 byte UUID from the config

			These UUIDs need to be decoded on the other device to show when the message was sent and who sent it
			The second device should be hosted on an AWS EC2 (free tier), please send the IP address with the test
		*/

		decoder := json.NewDecoder(strings.NewReader(string(messageBytes)))
		var message struct {
			Message  string `json:"message"`
			UUID     int64 `json:"uuid"`
		}
		err = decoder.Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// make byte array and put the UUID bytes into it
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(message.UUID))

		// grab the device ID
		d := make([]byte, 2)
		d[0] = b[6]
		d[1] = b[7]

		// zero the device ID bytes to avoid sending us into the future
		b[6] = 0
		b[7] = 0

		// get bytes into more usable types
		unixTimestamp := int64(binary.LittleEndian.Uint64(b))
		deviceIDStr := hex.EncodeToString(d)
		tm := time.Unix(unixTimestamp, 0)

		// log when message is received, from whom, and when
		log.Printf("recv: %s", message.Message)
		log.Printf("recv at: %s", tm)
		log.Printf("recv from: 0x%s", deviceIDStr)
	}
}


func home(w http.ResponseWriter, r *http.Request) {
	chTemplate.Execute(w, "ws://"+r.Host+"/chtest")
}


func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/chtest", chtest)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// basic homepage template (from example)
// could be extended to send message in same format as command line client
// Also, golang application could be extended to provide chat facilities, ie broadcast, rooms etc
var chTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };

	window.addEventListener('load', function() {
	    console.log('Connect to websocket on backend');
 		if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>
"Send" to send a message to the server and all listeners
You can change the message and send multiple times.
<p>
<form>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
