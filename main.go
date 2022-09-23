package main

import (
	"github.com/Erickype/GoChatAppGRPC/proto"
	log "google.golang.org/grpc/grpclog"
	"os"
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
	s          proto.UnimplementedBroadcastServer
	Connection []*Connection
}

func (s *Server) CreateStream(pconn *proto.Connect, stream proto.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, conn)
	return <-conn.error
}

func main() {

}