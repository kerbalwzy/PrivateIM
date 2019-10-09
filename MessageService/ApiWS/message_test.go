package ApiWS

import (
	"encoding/json"
	"testing"
	"time"
)

var (
	testReceiverId  int64 = 123
	testTextContent       = "<the is a test text message>"
)

func TestGetMessageTypeId(t *testing.T) {
	tempMessage := ChatMessage{
		BasicMessage: BasicMessage{
			TypeId:     UserChatMessageTypeId,
			SenderId:   SystemId,
			ReceiverId: testReceiverId,
		},
		ContentType: TextContent,
		Content:     testTextContent,
	}
	messageData, _ := json.Marshal(tempMessage)
	messageTypeId, err := GetRawJsonMessageTypeId(messageData)
	if nil != err {
		t.Fatal(err)
	}
	if messageTypeId != UserChatMessageTypeId {
		t.Fatal("the message type id is wrong")
	}
}

func TestBasicMessage_SetCreateTime(t *testing.T) {
	tempMessage := ChatMessage{
		BasicMessage: BasicMessage{
			TypeId:     UserChatMessageTypeId,
			SenderId:   SystemId,
			ReceiverId: testReceiverId,
		},
		ContentType: TextContent,
		Content:     testTextContent,
	}
	tempMessage.SetCreateTime()
	if tempMessage.CreateTime == 0 || tempMessage.CreateTime > time.Now().Unix() {
		t.Fatal("set message create time fail")
	}
}
