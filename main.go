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

func main() {

}
