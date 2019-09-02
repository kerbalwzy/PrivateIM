package ApiWS

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"../ApiRPC"
	"../DataLayer"
)

// WebSocket upgrade worker
var UpGrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }}

// Saving all node for every client connection.
var ClientPool = NewNodePool()

// Saving the message which want sent to offline user.
// This chan has buffer, and the cap is 10000. Every element is an array,
// which saved the target user's id and bytes data of the message.
var WaitSendChan = make(chan [2]interface{}, 10000)

// Handle WebSocket upgrade request, create a connection node for the client
// and save it into ClientPool. Before upgrading, check token by send it to
// UserCenter through gRPC call, if the result is false, the upgrade request
// will be fail and close, else it will create a new connection node for the
// client, and save it in ClientPool. After saved the new node, will start
// the "SendLoop" and "RecvLoop" of it. Then try to query the messages which
// other users sent to him when the client's user is offline, and send them.
func BeginChat(w http.ResponseWriter, r *http.Request) {
	// check the token by gRPC call , try to get the user's id.
	token := r.URL.Query().Get("authToken")
	userId, err := ApiRPC.CheckAuthToken(token)
	if nil != err {
		log.Printf("Error: check auth token fail for client(%s)", r.RemoteAddr)
		w.WriteHeader(400)
		_, _ = w.Write([]byte("authToken authentication fail: " + err.Error()))
		return
	}
	conn, err := UpGrader.Upgrade(w, r, nil)
	if nil != err {
		log.Println(err.Error())
		w.WriteHeader(400)
		_, _ = w.Write([]byte("upgrade connection fail: " + err.Error()))
		return
	}

	// new a node for the client connection and add it into ClientPool, and then
	// turn on the send and recv loop of the node.
	theNode := NewNode(userId, conn)
	ClientPool.Add(theNode)
	go theNode.SendLoop()
	go theNode.RecvLoop()

	// query and send `WaitSendMessage` for this user and send them.
	if messages, err := DataLayer.MongoQueryWaitSendMessage(userId); nil == err {
		for _, message := range messages {
			theNode.messageChan <- message
		}
	}

}

// Dispatch the chat message from ordinary user
// It only support NormalMessage and GroupMessage as present.
// If the `TypeId` of message is not allow
// it will send a error information to user client.
func ChatMessageDispatch(srcId int64, data []byte) {
	message := ChatMessage{}
	err := json.Unmarshal(data, message)

	// send back the error message to sender
	if nil != err {
		log.Println(err.Error())
		SendErrorMessage(message.SrcId, 400, err)
	}
	// check the SrcId, prevent users from sending messages by counterfeiting others
	if srcId != message.SrcId {
		log.Printf("Disguise:user(%d) try to send message as user(%d)",
			srcId, message.SrcId)
		err := errors.New("SrcId required equal to your own userId")
		SendErrorMessage(srcId, 400, err)
	}
	message.SetCreateTime()
	switch message.TypeId {
	case NormalMessage:
		SendNormalMessage(message)
	default:
		err := errors.New("unsupported message type id")
		SendErrorMessage(srcId, 400, err)
	}
}

// Send error message to target user client,
// if the user is not online, the message will be discarded
func SendErrorMessage(dstId int64, code int, err error) {
	message := ErrorMessage{
		BasicMessage: BasicMessage{
			TypeId:     DebugMessage,
			SrcId:      SystemId,
			DstId:      dstId,
			CreateTime: time.Now().Unix(),
		},
		Code: code, Error: err.Error()}
	if node, ok := ClientPool.Get(dstId); ok {
		log.Printf("ErrorMessage:send to user(%d) {code: %d, error: %s}",
			message.DstId, message.Code, message.Error)
		node.messageChan <- message.ToJson()
	}
}

// Send normal message to target user client, if the user is not online,
// it will record the message as a WaitSendMessage into database
func SendNormalMessage(message ChatMessage) {
	dstId := message.GetDstId()
	log.Printf("NormalMessage:user(%d) send a message to user(%d)",
		message.GetSrcId(), dstId)
	node, ok := ClientPool.Get(dstId)
	if ok {
		node.messageChan <- message.ToJson()
	} else {
		//
		WaitSendChan <- [2]interface{}{message.GetDstId(), message}
	}

}

// Check the user existence by DstId
func CheckDstIdExistence(message Message) bool {
	dstId := message.GetDstId()
	log.Printf("WaitSendMessage:SrdId(%d) check DisId(%d) return %t, ",
		message.GetSrcId(), dstId, true)
	return true
}

// Record the message for offline user.
// Save the message as WaitSendMessage into database,
// when the user online again, send these message to the client,
// but if the type of message is DebugMessage or the DstId dose not existed,
// it would not be save.
func SaveWaitSendMessage() {
	for message := range WaitSendChan {
		dstId := message[0].(int64)
		messageData := message[1].([]byte)

	}

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
