package main

import (
	"context"
	"github.com/Erickype/GoChatAppGRPC/proto"
	log "google.golang.org/grpc/grpclog"
	"sync"
)

var client proto.BroadcastClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *proto.User) error {
	var streamError error

	_, err := client.CreateStream(context.Background(), &proto.Connect{
		User:   user,
		Active: true,
	})
	if err != nil {
		log.Fatalf("Error connecting: %v", err.Error())
	}

	return streamError
}

func main() {

}
