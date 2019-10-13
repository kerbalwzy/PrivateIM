package MSGNode

import "time"

const SystemId = 1024

// The type code of message
const (
	UserChatMessageTypeId     = iota // Users chat with each other one to one
	GroupChatMessageTypeId           // User group chat
	SubscriptionMessageTypeId        // From system notification or subscription
	ErrorMessageTypeId               // Tell the client what error happened
)

// The type code of message content
const (
	TextContent  = iota // the text message
	ImageContent        // the picture message
	VideoContent        // the video message
	VoiceContent        // the voice message
)

// Saving the message which want sent to offline user.
// This chan has buffer, and the cap is 10000. Every element is an array,
// which saved the target user's id and bytes data of the message.
var DelayMessageChat = make(chan [2]interface{}, 10000)

type Message interface {
	SetCreateTime()
	GetDeliveryTime() int64
	GetReceiverId() int64
}

// Basic message struct
type BasicMessage struct {
	TypeId       int   `json:"type_id"`                 // the type number of message
	SenderId     int64 `json:"sender_id"`               // who send this message, the sender id
	ReceiverId   int64 `json:"receiver_id"`             // who will recv this message, the receiver id
	CreateTime   int64 `json:"create_time,omitempty"`   // set by the message center, timestamp, unit:sec.
	DeliveryTime int64 `json:"delivery_time,omitempty"` // the time for message want be sent, use for timing message
}

// Set the create time for the message
func (obj *BasicMessage) SetCreateTime() {
	obj.CreateTime = time.Now().Unix()
}

// Get the value of 'DeliveryTime' field
func (obj *BasicMessage) GetDeliveryTime() int64 {
	return obj.DeliveryTime
}

// Get the value of 'ReceiverId' field
func (obj *BasicMessage) GetReceiverId() int64 {
	return obj.ReceiverId
}

// ChatMessage, used to send NormalMessage and GroupMessage mainly
// the Content-Type can be of {0:text, 1:picture，2:video, 3:voice}
type ChatMessage struct {
	BasicMessage
	ContentType int    `json:"content_type"`           // how to show the message in client
	Content     string `json:"content,omitempty"`      // text content
	PreviewPic  string `json:"preview_pic,omitempty"`  // preview picture url
	ResourceUrl string `json:"resource_url,omitempty"` // resource URL
	Description string `json:"description,omitempty"`  // simple description
	Additional  string `json:"additional,omitempty"`   // other additional information
}

// SubscriptionMessage, used to send ChannelsMessage mainly
type SubscriptionMessage struct {
	BasicMessage
	Title       string `json:"title"`                  // the title
	Abstract    string `json:"abstract"`               // the brief introduction of this message
	PreviewPic  string `json:"preview_pic,omitempty"`  // the preview picture url
	ResourceUrl string `json:"resource_url,omitempty"` // resource URL
}

// ErrorMessage, used to send DebugMessage mainly
type ErrorMessage struct {
	BasicMessage
	Code       int    `json:"code"`                  // the code of error type
	Error      string `json:"error"`                 // the detail error information
	RawMessage []byte `json:"raw_message,omitempty"` // the row message which the user want to send.
}
