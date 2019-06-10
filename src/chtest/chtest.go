package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Configuration struct{
	DeviceID	string
}

type Message struct {
	UUID 		string
	TimeStamp 	int64
	Line 		string
}

var outgoingMessages = make(chan Message)

var deviceID string

var ip string



func main() {
	var err error

	// prevent any date / time etc appearing on cmd line
	log.SetFlags(0)

	/*
	Please create a command line messaging program written in Go.
	When the program loads it should ask the user to input the IP
	address of the other device running the messaging program.
	*/

	fmt.Println("Please enter the IP(v4) address of the server you would like to connect to")

	ip, err = getIPAddress()
	if err != nil {
		fmt.Println("error ", err.Error())
	}

	/*
	Each device will have a 2 byte unique ID to identify the user
	(this can be set in a config.json and read in when program starts)
	(eg. 0x0001)
	*/

	deviceID, err = readDeviceID()

	if err != nil {
		fmt.Println("Device ID err: ", err)
 	} else {
		fmt.Println("Device ID is: ", deviceID)
	}

	/*
	Use a Websocket connection to send messages back and forth between the two devices.
	Once a connection has been established, neither program should ask the user to put in an IP address anymore.
	*/

	var addr = join(ip, ":8080")

	// NB following code is based heavily on gorilla websocket examples - client
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/chtest"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	log.Printf("You are now connected to %s, type your message and hit return to send", u.String())

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()


	/*
	They will instead show the user a text input where messages can be typed.
	When enter is hit, it should will send the message the the second device where is should be shown.
	*/

	go userInputHandler()

	// set up handlers for channels
	for {
		select {
		case message := <-outgoingMessages:

			/*
				Each message should be sent with an 8 byte UUID,
				this should consist of a 6 byte timestamp (seconds from 01/01/1970) and the 2 byte UUID from the config

				These UUIDs need to be decoded on the other device to show when the message was sent and who sent it
				The second device should be hosted on an AWS EC2 (free tier), please send the IP address with the test
			*/

			var UUIDpart1 = strconv.FormatInt(message.TimeStamp, 16)

			// Could use AppendInt?
			var UUID = join(UUIDpart1, deviceID[2:])

			messageJSON, err := json.Marshal(struct {
				Message   	string `json:"message"`
				UUID    	string `json:"uuid"`
			}{
				message.Line,
				UUID,
			})

			err = c.WriteMessage(websocket.TextMessage, []byte(messageJSON))
			if err != nil {
				log.Println("write:", err)
				return
			}

		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			// wait a second
			case <-time.After(time.Second):
			}
			return
		}
	}
}


func getIPAddress() (IPAddress string, err error) {
	var _ip = ""

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var IPAddress = s.Text()
		ip := net.ParseIP(IPAddress)

		if ip == nil {
			fmt.Printf("Could not parse IPAddress %s ", IPAddress)
			fmt.Println("")
			continue
		}

		return IPAddress, nil
	}
	return _ip, nil
}



func userInputHandler()  {
	s := bufio.NewScanner(os.Stdin)
	for {
		for s.Scan() {
			var userInput = s.Text()
			/*
				this should consist of a 6 byte timestamp (seconds from 01/01/1970) and the 2 byte UUID from the config
			*/
			now := time.Now()
			secs := now.Unix()
			outgoingMessages <- Message{
				UUID: deviceID,
				TimeStamp: secs,
				Line: userInput,
			}
		}
	}
}


func readDeviceID() (DeviceID string, err error){

	file, _ := os.Open("config.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	error := decoder.Decode(&configuration)

	return configuration.DeviceID, error
}


func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

//type fatalError error

// fatal doesn't exit because that would stop deferred functions being
// called.  Instead, its panic can be recovered at the top of the call
// stack.
//func fatal(err fatalError) {
//	panic(err)
//}


//func fatal_if(err error) {
//	if err != nil {
//		fatal(err)
//	}
//}