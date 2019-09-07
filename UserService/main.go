package main

import (
	"log"
)

func init() {
	log.SetPrefix("AuthCenter ")

}

func main() {
	go StartHttpServer()
	StartGRPCServer()
}
