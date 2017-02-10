#!/bin/bash
#Script to run project under AWS beanstalk
PORT=80 go run $GOPATH/src/github.com/spring1843/chat-server/src/main.go -config $GOPATH/src/github.com/spring1843/chat-server/src/config.json