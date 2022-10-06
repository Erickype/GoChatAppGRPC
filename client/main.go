package main

import (
	"github.com/Erickype/GoChatAppGRPC/proto"
	"sync"
)

var _ proto.BroadcastClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func main() {

}
