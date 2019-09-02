package ApiWS

import (
	"encoding/json"
	"time"
)

const SystemId = 1024

// The type code of message
const (
	NormalMessage   = 0 // Users chat with each other one to one
	GroupMessage    = 1 // User group chat
	ChannelsMessage = 2 // From system notification or user subscription
	DebugMessage    = 3 // Tell the client what error happened
)

type Message interface {
	ToJson() []byte
	SetCreateTime()
	GetTypeId() int
	GetSrcId() int64
	GetDstId() int64
}

// Basic message struct
type BasicMessage struct {
	TypeId     int   `json:"type_id"`               // the type number of message
	SrcId      int64 `json:"src_id"`                // who send this message, the sender id
	DstId      int64 `json:"dst_id"`                // who will recv this message, the receiver id
	CreateTime int64 `json:"create_time,omitempty"` // Added by the message center, timestamp, unit:sec.
}

// Encode the message to json string
func (obj *BasicMessage) ToJson() []byte {
	data, _ := json.Marshal(obj)
	return data
}

// Set the create time for the message
func (obj *BasicMessage) SetCreateTime() {
	obj.CreateTime = time.Now().Unix()
}

// Get the TypeId of the message
func (obj *BasicMessage) GetTypeId() int {
	return obj.TypeId
}

// Get the SrcId of the message
func (obj *BasicMessage) GetSrcId() int64 {
	return obj.SrcId
}

// Get the DstId of the message
func (obj *BasicMessage) GetDstId() int64 {
	return obj.DstId
}

// ChatMessage, used to send NormalMessage and GroupMessage mainly
// the ContentType can be of {0:text, 1:picture，2:video, 3:voice}
type ChatMessage struct {
	BasicMessage
	ContentType int    `json:"content_type"`          // how to show the message in client
	Content     string `json:"content,omitempty"`     // text content
	PreviewPic  string `json:"preview_pic,omitempty"` // preview picture url
	Url         string `json:"url,omitempty"`         // resource URL
	Description string `json:"description,omitempty"` // simple description
	Additional  string `json:"additional,omitempty"`  // other additional information
}

// SubscriptionMessage, used to send ChannelsMessage mainly
type SubscriptionMessage struct {
	BasicMessage
	Title      string `json:"title"`                 // the title
	Abstract   string `json:"abstract,omitempty"`    // the brief introduction of this message
	PreviewPic string `json:"preview_pic,omitempty"` // the preview picture url
	Url        string `json:"url,omitempty"`         // the link of complete content
}

// ErrorMessage, used to send DebugMessage mainly
type ErrorMessage struct {
	BasicMessage
	Code  int    `json:"code"`            // the code of error type
	Error string `json:"error,omitempty"` // the detail error information
}
