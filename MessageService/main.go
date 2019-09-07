package main

import (
	"log"
	"net/http"

	"./ApiWS"
)

func init() {
	log.SetPrefix("MsgCenter ")
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {
	http.HandleFunc("/", ApiWS.BeginChat)
	log.Println("start message center websocket server...")
	err := http.ListenAndServe(":8000", nil)
	if nil != err {
		log.Fatal(err)
	}
}
