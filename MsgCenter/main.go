package main

import (
	"log"
	"ruozhuo"
)

func init() {
	log.SetPrefix("MsgCenter ")
}

func main() {
	worker := &ruozhuo.Worker{
		Tag:          "msg",
		Address:      "127.0.0.1:10000",
		ManagerAdder: "127.0.0.1:7070",
	}
	ruozhuo.StartWorker(worker)
}