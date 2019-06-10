# Command line messaging project

Example of a command line messaging project running over websockets, with configuration.

Basis for websocket implementation taken from Gorilla WebSocket [Echo Example][1] project.

## Configuration

Copy the `config.json.example` file to `src/chtest/config.json`

ie,

```bash
$ cd <project_path>
$ cp config.json.example src/chtest/config.json
```

Edit config.json with value for client - ie change the value of deviceID:

```json
{
  "deviceID": "0x0001"
}
```

## Requirements

Get the Golang libraries

```bash
$ export GOPATH=`pwd`
$ cd src
$ go get chtest

```



## Client

Start the client by typing:

```bash
$ go run chtest.go
```

Connect to the running server on AWS Free Tier by entering  54.229.136.220 when asked for the IP address.

## Server

Server can be run locally, to see unpacking of date and deviceID by typing  
  
```bash
$ cd <project_path>
$ export GOPATH=`pwd`
$ cd <project_path>/src/server
$ go run server.go
```

## Build

An executable can be built for the chtest client: 

```bash
$ cd <project_path>
$ export GOPATH=`pwd`
$ go build chtest
```
This can be run, after the config file is copied to the bin directory:

```bash
$ cp config.json.example bin/config.json
$ ./bin/chtest
```



[1]: https://github.com/gorilla/websocket/tree/master/examples/echo
