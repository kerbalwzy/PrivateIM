package ApiWS

import (
	"log"
	"net/http"

	conf "../Config"
)

func StartMessageWebSocketServer() {
	go SaveDelayMessageLoop()
	go GlobalGroupChats.CleanGroupChatLoop()
	go GlobalSubscriptions.CleanByLifeTimeLoop()

	// start the message transfer WebSocket server
	http.HandleFunc("/", BeginChat)

	log.Printf("[info] start MessageService with address: %s", conf.MessageServerAddress)
	err := http.ListenAndServe(conf.MessageServerAddress, nil)
	if nil != err {
		log.Fatalf("[error] start MessageService fail: %s", err.Error())
	}
}
