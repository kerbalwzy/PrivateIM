package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"./ApiRPC"
)

func StartGRPCServer() {
	listener, err := net.Listen("tcp", ":23333")
	if nil != err {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	ApiRPC.RegisterUserCenterServer(server, &ApiRPC.MygRPCServer{})
	log.Println("------------------start Golang gRPC server")
	err = server.Serve(listener)
	if nil != err {
		log.Fatal(err)
	}
}
