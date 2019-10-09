package main

import (
	"log"

	"./ApiWS"
)

func init() {
	log.SetPrefix("MsgCenter ")
	log.SetFlags(log.LstdFlags)
}

func main() {
	ApiWS.StartMessageWebSocketServer()
}
