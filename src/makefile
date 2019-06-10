#! /usr/bin/make -f

# GOROOT may need to be set if Go's installed somewhere unusual, e.g.
#     GOROOT=/home/ralph/go/go1.5.1.linux-amd64
# PATH needs to contain $GOROOT/bin.
# GOPATH needs to be set, e.g. export GOPATH=/home/ec2-user/chtest
#
# To build from source for local architecture do:
#     make var check all
#
# To build from source for new executable suitable for deployment, e.g.
# another server's architecture, do:
#     make var check deploy

SHELL = /bin/bash

all:
	go install chtest   # Native executable in ../bin/chtest.
	go install server   # Native executable in ../bin/server.

var:
	[[ -d $$GOPATH ]]
	[[ -f $$GOPATH/src/chtest/main.go ]]

check:
	go vet chtest
	go test chtest

bench:
	go test -bench ^ -benchmem chtest

fmt:
	go fmt chtest

deploy:
	GOOS=linux GOARCH=amd64 go install -installsuffix deploy chtest

clean:
	rm -f ../bin/chtest
	rm -f ../bin/*/chtest
	rm -f ../bin/server
	rm -f ../bin/*/server
