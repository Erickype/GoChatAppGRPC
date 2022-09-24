package main

import (
	"context"
	"github.com/Erickype/GoChatAppGRPC/proto"
	log "google.golang.org/grpc/grpclog"
	"os"
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

func (s *Server) BroadCastMessage(_ context.Context, _ *proto.Message) (*proto.Close, error) {

	return &proto.Close{}, nil

}

func main() {

}
