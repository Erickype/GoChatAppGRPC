package main

import (
	"context"
	"github.com/Erickype/GoChatAppGRPC/proto"
	"google.golang.org/grpc"
	log "google.golang.org/grpc/grpclog"
	"net"
	"os"
	"sync"
)

var grpcLog log.LoggerV2

func init() {
	grpcLog = log.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type Connection struct {
	stream proto.Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type Server struct {
	proto.UnimplementedBroadcastServer
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
				grpcLog.Infoln("Sending message to: ", conn.stream)

				if err != nil {
					grpcLog.Errorf("Error with stream: %s - Error: %v", conn.stream, err)
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
	var connections []*Connection

	server := &Server{
		Connection: connections,
	}

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		grpcLog.Fatalf("Error creating TCP conn: %v", err.Error())
	}
	grpcLog.Info("Starting server at port: 8080")

	proto.RegisterBroadcastServer(grpcServer, server)
	err = grpcServer.Serve(listener)
	if err != nil {
		grpcLog.Fatalf("Cannot initialize the server: %v", err.Error())
	}

}
