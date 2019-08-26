package ApiWS

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	//"../ApiRPC"
	"../DataLayer"
)

var UpGrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//authToken := r.URL.Query().Get("authToken")
		//if authToken == "" {
		//	return false
		//}
		//ok, _ := ApiRPC.CheckAuthToken(authToken)
		//return ok
		return true
	},
}
// ClientPool to save the client node
var ClientPool = NewNodePool()

// WaitSendChan to save the message which want send to offline user
// this chan has buffer, and the cap is 10000
var WaitSendChan = make(chan DataLayer.Message, 10000)

/*
Handle WebSocket upgrade request, create a connection node for the client and save it in ClientPool.
Before upgrading, check token by send it to UserCenter through gRPC call, if the result is false,
the upgrade request will be fail and close, else it will create a new connection node for the client,
and save it in ClientPool. After saved the new node, will start the "SendLoop" and "RecvLoop" of it.
Then try to find if there are messages which other users send to him when the client's user offline,
if existed, send them to the client at now.
*/
func BeginChat(w http.ResponseWriter, r *http.Request) {
	conn, err := UpGrader.Upgrade(w, r, nil)
	if nil != err {
		log.Println(err.Error())
		w.WriteHeader(400)
		_, _ = w.Write([]byte("upgrade connection fail"))
		return
	}
	id := r.URL.Query().Get("id")

	userId, err := strconv.ParseInt(id, 10, 64)
	if nil != err {
		log.Println(err.Error())
		w.WriteHeader(400)
		_, _ = w.Write([]byte("parse id to int64 error"))
		return
	}
	newNode := &Node{Id: userId, conn: conn, messageQueue: make(chan DataLayer.Message)}
	ClientPool.Add(newNode)

	go newNode.SendLoop()
	go newNode.RecvLoop()

}

// Dispatch the chat message from ordinary user
// It only support NormalMessage and GroupMessage as present.
// If the `TypeId` of message is not allow
// it will send a error information to user client.
func ChatDispatch(srcId int64, data []byte) {
	message := &DataLayer.ChatMessage{}
	err := json.Unmarshal(data, message)

	// send back the error message to sender
	if nil != err {
		log.Println(err.Error())
		SendErrorMessage(message.SrcId, 400, err)
	}
	// check the SrcId, prevent users from sending messages by counterfeiting others
	if srcId != message.SrcId {
		log.Printf("Disguise:user(%d) try to send message as user(%d)", srcId, message.SrcId)
		err := errors.New("SrcId required equal to your own userId")
		SendErrorMessage(srcId, 400, err)
	}
	message.SetCreateTime()
	switch message.TypeId {
	case DataLayer.NormalMessage:
		SendNormalMessage(message)
	default:
		err := errors.New("unsupported message type id")
		SendErrorMessage(srcId, 400, err)
	}
}

// Send error message to target user client,
// if the user is not online, the message will be discarded
func SendErrorMessage(dstId int64, code int, err error) {
	message := &DataLayer.ErrorMessage{
		BasicMessage: DataLayer.BasicMessage{
			TypeId:     DataLayer.DebugMessage,
			SrcId:      DataLayer.SystemId,
			DstId:      dstId,
			CreateTime: time.Now().Unix(),
		},
		Code: code, Error: err.Error()}
	if node, ok := ClientPool.Get(dstId); ok {
		log.Printf("ErrorMessage:send to user(%d) {code: %d, error: %s}", message.DstId, message.Code, message.Error)
		node.messageQueue <- message
	}
}

// Send normal message to target user client, if the user is not online,
// it will record the message as a WaitSendMessage into database
func SendNormalMessage(message DataLayer.Message) {
	dstId := message.GetDstId()
	log.Printf("NormalMessage:user(%d) send a message to user(%d)", message.GetSrcId(), dstId)
	node, ok := ClientPool.Get(dstId)
	if ok {
		node.messageQueue <- message
	} else {
		//
		WaitSendChan <- message
	}

}

// Record the message for offline user.
// Save the message as WaitSendMessage into database,
// when the user online again, send these message to the client,
// but if the type of message is DebugMessage or the DstId dose not existed,
// it would not be save.
func SaveWaitSendMessage() {
	for message := range WaitSendChan {
		if message.GetTypeId() == DataLayer.DebugMessage {
			continue
		}
		if !CheckDstIdExistence(message) {
			continue
		}

	}

}

// Check the user existence by DstId
func CheckDstIdExistence(message DataLayer.Message) bool {
	dstId := message.GetDstId()
	log.Printf("WaitSendMessage:SrdId(%d) check DisId(%d) return %t, ", message.GetSrcId(), dstId, true)
	return true
}

/*
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMTYyMjYyOTQ4Nzk0NTk3Mzc2LCJleHAiOjE1NjY1MjkwNjQsImlzcyI6InVzZXJDZW50ZXIifQ.nelQ8fHIgrUovgFOUguZCspdGQAYXmiM8_hcYKI2L8s
*/
/*
{
	"Type": 0,
	"src_id": 1,
	"dst_id": 2,
	"content_type": 1,
	"content": "I am a message from client 1"

}
*/
