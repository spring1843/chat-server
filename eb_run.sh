#!/bin/bash
#Script to run project under AWS beanstalk
go get github.com/emicklei/go-restful/
go get github.com/gorilla/websocket/
PORT=80 go run $GOPATH/src/github.com/spring1843/chat-server/main.go -config $GOPATH/src/github.com/spring1843/chat-server/config.json