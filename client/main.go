package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/Erickype/GoChatAppGRPC/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	log "google.golang.org/grpc/grpclog"
	"os"
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

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err.Error())
	}

	client = proto.NewBroadcastClient(conn)
	user := &proto.User{
		Id:   hex.EncodeToString(id[:]),
		Name: *name,
	}

	err = connect(user)
	if err != nil {
		return
	}

	go func() {
		defer wait.Done()
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			msg := &proto.Message{
				Id:        user.Id,
				Content:   scanner.Text(),
				Timestamp: timestamp.String(),
			}

			_, err := client.BroadCastMessage(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error sending message: %v", err.Error())
				break
			}
		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
}
