package main

import "./ApiRPC"

func main() {
	go ApiRPC.StartMySQLDataRPCServer()
	ApiRPC.StartMongoDataRPCServer()
}
