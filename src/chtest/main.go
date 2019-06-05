package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type configuration struct{
	DeviceId	string
}

var ip string

func main() {
	//log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	/*
		Please create a command line messaging program written in Go.
		When the program loads it should ask the user to input the IP address of the other device running the messaging program.
	*/

	fmt.Println("Please enter the IP(v4) address of the chtest you would like to connect to")

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var IPAddress = s.Text()
		ip := net.ParseIP(IPAddress)
		if ip != nil {
			fmt.Println("Connecting to",IPAddress )
			break
		}
		fmt.Println("The IP Address is not valid: ", IPAddress)
	}


	/*
		Each device will have a 2 byte unique ID to identify the user (this can be set in a config.json and read in when program starts) (eg. 0x0001)
	*/

	file, _ := os.Open("config.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Cannot read configuration file")
		panic(err)
	}

	fatal_if(err)

	fmt.Println("Device ID is: ", conf.DeviceId)


	/*

		Use a Websocket connection to send messages back and forth between the two devices.
		Once a connection has been established, neither program should ask the user to put in an IP address anymore. They will instead show the user a text input where messages can be typed. When enter is hit, it should will send the message the the second device where is should be shown.
	*/


	/*
		Each message should be sent with an 8 byte UUID, this should consist of a 6 byte timestamp (seconds from 01/01/1970) and the 2 byte UUID from the config
		These UUIDs need to be decoded on the other device to show when the message was sent and who sent it
		The second device should be hosted on an AWS EC2 (free tier), please send the IP address with the test
	*/


}


type fatalError error

// fatal doesn't exit because that would stop deferred functions being
// called.  Instead, its panic can be recovered at the top of the call
// stack.
func fatal(err fatalError) {
	panic(err)
}


func fatal_if(err error) {
	if err != nil {
		fatal(err)
	}
}