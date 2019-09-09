package main

import (
	"log"

	"./ApiHTTP"
	"./ApiRPC"
)

func init() {
	log.SetPrefix("AuthCenter ")

}

func main() {
	go ApiHTTP.StartHttpServer()
	ApiRPC.StartUserAuthRPCServer()
}
