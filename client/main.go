package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/Erickype/GoChatAppGRPC/proto"
	"sync"
	"time"
)

var client proto.BroadcastClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *proto.User) error {
	var streamError error

	stream, err := client.CreateStream(context.Background(), &proto.Connect{
		User:   user,
		Active: true,
	})
	if err != nil {
		return fmt.Errorf("error connecting: %v", err.Error())
	}

	wait.Add(1)

	go func(str proto.Broadcast_CreateStreamClient) {
		defer wait.Done()
		for {
			msg, err := str.Recv()
			if err != nil {
				streamError = fmt.Errorf("error reading message: %v", err.Error())
				break
			}
			fmt.Printf("%v: %s\n", msg.Id, msg.Content)
		}
	}(stream)

	return streamError
}

func main() {
	timestamp := time.Now()
	done := make(chan int)
	name := flag.String("N", "Anon", "The name of the user")
	flag.Parse()
	id := sha256.Sum256([]byte(timestamp.String() + *name))
}
