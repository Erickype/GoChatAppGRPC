package main

import (
	log "google.golang.org/grpc/grpclog"
	"os"
)

var grpcLog log.LoggerV2

func init() {
	grpcLog = log.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

func main() {

}
