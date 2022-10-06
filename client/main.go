package main

import (
	"context"
	"fmt"
	"github.com/Erickype/GoChatAppGRPC/proto"
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
		return fmt.Errorf("error connecting: %v", err.Error())
	}

	return streamError
}

func main() {

}
