package main

import (
	"context"
	"github.com/Erickype/GoChatAppGRPC/proto"
	log "google.golang.org/grpc/grpclog"
	"os"
	"sync"
)

var _ log.LoggerV2

func init() {
	_ = log.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type Connection struct {
	stream proto.Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type Server struct {
	s          proto.UnimplementedBroadcastServer
	Connection []*Connection
}

func (s *Server) CreateStream(connect *proto.Connect, stream proto.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     connect.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, conn)
	return <-conn.error
}

func (s *Server) BroadCastMessage(_ context.Context, msg *proto.Message) (*proto.Close, error) {

	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		wait.Add(1)
		go func(msg *proto.Message, conn *Connection) {
			defer wait.Done()
			if conn.active {
				err := conn.stream.Send(msg)
				log.Infoln("Sending message to: ", conn.stream)

				if err != nil {
					log.Errorf("Error with stream: %s - Error: %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}
	go func() {
		wait.Wait()
		close(done)
	}()

	<-done

	return &proto.Close{}, nil

}

func main() {

}
