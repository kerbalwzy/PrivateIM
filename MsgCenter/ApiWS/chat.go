package ApiWS

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	//"../ApiRPC"
	"../DataLayer"
)

const (
	NormalMessage = 0
	GroupMessage  = 1
)

type Message interface {
	ToJson() ([]byte, error)
}

type TextMessage struct {
	Type        int    `json:"type"`                   // message type: NormalMessage or GroupMessage
	SrcId       int64  `json:"src_id"`                 // who send this message: the client user id
	DstId       int64  `json:"dst_id"`                 // send to who: user id or group id
	ContentType int    `json:"content_type"`           // how to show the message in client
	Content     string `json:"content,omitempty"`      // message content
	PreviewPic  string `json:"preview_pic,omitempty"`  // preview picture url
	Url         string `json:"url,omitempty"`          // resource URL
	Description string `json:"description,omitempty"`  // simple description
	Others      string `json:"others,omitempty"`       // other additional information
	ProduceTime int64  `json:"produce_time,omitempty"` // the produce timestamp of message (Unit:sec.)
}

func (obj *TextMessage) ToJson() ([]byte, error) {
	data, err := json.Marshal(obj)
	return data, err
}

type ErrorMessage struct {
	Code        int    `json:"code"`
	Error       string `json:"error,omitempty"`
	ProduceTime int64  `json:"produce_time,omitempty"` // the produce timestamp of message (Unit:sec.)
}

func (obj *ErrorMessage) ToJson() ([]byte, error) {
	data, err := json.Marshal(obj)
	return data, err
}

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

var ClientPool = NewNodePool()

// Handle websocket upgrade request, create a node for this connection and save in ClientPool
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
	newNode := &Node{Id: userId, conn: conn, dataQueue: make(chan []byte)}
	ClientPool.Add(newNode)

	go newNode.SendLoop()
	SendErrorMessage(userId, []byte("hello"))
	go newNode.RecvLoop()

}

// Send the message data to target user client.
// It only support NormalMessage and GroupMessage as present. If the `Type` of message,
// it will send a error information to user client.
func Dispatch(data []byte) error {
	msg := TextMessage{}
	err := json.Unmarshal(data, &msg)
	if nil != err {
		log.Println(err.Error())
		return err
	}
	switch msg.Type {
	case NormalMessage:
		SendNormalMessage(msg.DstId, data)
	default:
		tempErrMsg := ErrorMessage{Code: 400, Error: "no support Type of message", ProduceTime: time.Now().Unix()}
		data, _ := tempErrMsg.ToJson()
		SendErrorMessage(msg.DstId, data)
	}
	return nil
}

// Send error message to target user client, if the user is not online, will not send.
func SendErrorMessage(dstId int64, data []byte) {
	if node, ok := ClientPool.Get(dstId); ok {
		node.dataQueue <- data
	}
}

// Send normal message to target user client, if the user is not online,
// it will record the message as a WaitSendMessage into database
func SendNormalMessage(dstId int64, data []byte) {
	log.Printf("send a normal message to user: %d", dstId)
	node, ok := ClientPool.Get(dstId)
	if ok {
		node.dataQueue <- data
	} else {
		tempMsg := new(TextMessage)
		err := json.Unmarshal(data, tempMsg)
		if nil != err {
			log.Println(err.Error())
			return
		}
		RecordWaitSendMessage(dstId, tempMsg)
	}

}

// Record the message for offline user.
// Save the message as WaitSendMessage into database, when the user online again, send
// these message to the client
func RecordWaitSendMessage(dstId int64, msg Message) {
	// set ProduceTime
	v := reflect.ValueOf(msg)
	p := v.Elem().FieldByName("ProduceTime")
	p.SetInt(time.Now().Unix())

	// record the waiting send message to database
	data, _ := msg.ToJson()
	log.Printf("record a WaitSendMessage for user: %d", dstId)
	DataLayer.RecordWaitSendMessage(dstId, data)
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
