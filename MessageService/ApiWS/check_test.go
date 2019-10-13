package ApiWS

import (
	"encoding/json"
	"testing"
	"time"

	"../MSGNode"
)

var (
	testReceiverId  int64 = 123
	testTextContent       = "<the is a test text message>"
)

func TestGetMessageTypeId(t *testing.T) {
	tempMessage := MSGNode.ChatMessage{
		BasicMessage: MSGNode.BasicMessage{
			TypeId:     MSGNode.UserChatMessageTypeId,
			SenderId:   MSGNode.SystemId,
			ReceiverId: testReceiverId,
		},
		ContentType: MSGNode.TextContent,
		Content:     testTextContent,
	}
	messageData, _ := json.Marshal(tempMessage)
	messageTypeId, err := GetRawJsonMessageTypeId(messageData)
	if nil != err {
		t.Fatal(err)
	}
	if messageTypeId != MSGNode.UserChatMessageTypeId {
		t.Fatal("the message type id is wrong")
	}
}

func TestBasicMessage_SetCreateTime(t *testing.T) {
	tempMessage := MSGNode.ChatMessage{
		BasicMessage: MSGNode.BasicMessage{
			TypeId:     MSGNode.UserChatMessageTypeId,
			SenderId:   MSGNode.SystemId,
			ReceiverId: testReceiverId,
		},
		ContentType: MSGNode.TextContent,
		Content:     testTextContent,
	}
	tempMessage.SetCreateTime()
	if tempMessage.CreateTime == 0 || tempMessage.CreateTime > time.Now().Unix() {
		t.Fatal("set message create time fail")
	}
}
