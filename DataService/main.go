package main

import (
	"./ApiRPC"

	"log"
)

func init() {
	log.SetFlags(log.LstdFlags)
	log.SetPrefix("<DataService> ")
}

func main() {
	go ApiRPC.StartMySQLDataRPCServer()
	ApiRPC.StartMongoDataRPCServer()
}
