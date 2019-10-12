package main

import (
	"log"

	"./ApiRPC"
	"./ApiWS"
)

func init() {
	log.SetPrefix("MsgCenter ")
	log.SetFlags(log.LstdFlags)
}

func main() {
	go ApiRPC.StartMsgTransferRPCServer()
	ApiWS.StartMessageWebSocketServer()
}
