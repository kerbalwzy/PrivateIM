package main

import (
	"log"
	"rouzhuo"
)

func init()  {
	log.SetFlags(log.LstdFlags|log.Llongfile)
	log.SetPrefix("IMServer ")
}

func main() {
	manager := &rouzhuo.Manager{ApiAddress:"127.0.0.1:8080"}
	manager.StartApiServer()
}