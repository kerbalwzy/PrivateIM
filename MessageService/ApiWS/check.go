package ApiWS

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"../ApiRPC"
	"../MSGNode"
)

var (
	ErrNoMessageTypeId    = errors.New("not have 'type_id' field in the json string message")
	ErrUnSupportMsgTypeId = errors.New("the value of 'type_id' is not support")
)

// Get the value of 'type_id' from the json string message
func GetRawJsonMessageTypeId(message []byte) (int, error) {
	temp := make(map[string]interface{})
	err := json.Unmarshal(message, &temp)
	if nil != err {
		return -1, err
	}
	if typeId, ok := temp["type_id"]; ok {
		return int(reflect.ValueOf(typeId).Float()), nil
	}
	return -1, ErrNoMessageTypeId
}

var (
	ErrSendMessageOk        = errors.New("the message send ok")
	ErrUserDisguise         = errors.New("the sender id is not identical, don't disguise other to send message")
	ErrUnSupportContentType = errors.New("the message content type is not support")
	ErrTextContentEmpty     = errors.New("the text content is empty")
	ErrPreviewPicEmpty      = errors.New("the media preview picture url is empty")
	ErrResourceURLEmpty     = errors.New("the media resource url is empty")
	ErrReceiverRefuseRecv   = errors.New("the receiver refuse receive the message")
	ErrFriendshipNotExisted = errors.New("your are not friend still")
	ErrReceiverNotExisted   = errors.New("the receiver is not existed")
)

// Check the chat message data legality.
// Requiring the sender id is identical, and the content type is supported, and the message real content not be empty
func checkUserChatMessageData(senderId int64, chatMessage *MSGNode.ChatMessage) (int, error) {
	if chatMessage.SenderId != senderId {
		return 400, ErrUserDisguise
	}

	switch chatMessage.ContentType {
	case MSGNode.TextContent:
		if chatMessage.Content == "" {
			return 400, ErrTextContentEmpty
		}
	case MSGNode.ImageContent, MSGNode.VoiceContent:
		if chatMessage.PreviewPic == "" {
			return 400, ErrPreviewPicEmpty
		}
		if chatMessage.ResourceUrl == "" {
			return 400, ErrResourceURLEmpty
		}
	case MSGNode.VideoContent:
		if chatMessage.ResourceUrl == "" {
			return 400, ErrResourceURLEmpty
		}

	default:
		return 400, ErrUnSupportContentType
	}

	return 200, nil
}

// Check whether the receiver should receive the message.
// Requiring the sender have an effective friendship with the receiver.
func checkWhetherReceiverShouldReceive(ok bool, receiver *MSGNode.UserNode, senderId, receiverId int64) (int, error) {
	switch ok {
	case true:
		// when the receiver is online
		if receiver.BlackList.Exist(senderId) {
			return 403, ErrReceiverRefuseRecv
		}
		if !receiver.Friends.Exist(senderId) {
			return 403, ErrFriendshipNotExisted
		}
	case false:
		// when the receiver is offline
		friends, err := ApiRPC.GetUserFriends(receiverId)
		if nil != err {
			return 404, ErrReceiverNotExisted
		}

		if blacklist, _ := ApiRPC.GetUserBlacklist(receiverId); nil != blacklist {
			for _, id := range blacklist {
				if id == senderId {
					return 403, ErrReceiverRefuseRecv
				}
			}
		}
		isNotFriend := true
		for _, id := range friends {
			if id == senderId {
				isNotFriend = false
				break
			}
		}
		if isNotFriend {
			return 403, ErrFriendshipNotExisted
		}
	}
	return 200, nil
}

var (
	ErrTitleEmpty    = errors.New("the title can not be empty")
	ErrAbstractEmpty = errors.New("the abstract can not be empty")
)

// Check the subscription message data legality
// Requiring the Title and Abstract not be empty
func checkSubscriptionMessageData(senderId int64, subsMessage *MSGNode.SubscriptionMessage) (int, error) {
	if subsMessage.SenderId != senderId {
		return 400, ErrUserDisguise
	}
	if subsMessage.Title == "" {
		return 400, ErrTitleEmpty
	}
	if subsMessage.Abstract == "" {
		return 400, ErrAbstractEmpty
	}
	return 200, nil
}

var (
	ErrDeliveryTime = errors.New("invalid delivery time")
)

// Check the value of 'DeliveryTime' field in the message. If the value is zero, meaning it is not a timing message.
// When it is a timing message, requiring the delivery time is after now at least 2 minute.
func checkWeatherTimingMessage(message MSGNode.Message) (bool, error) {
	deliveryTime := message.GetDeliveryTime()
	if 0 == deliveryTime {
		return false, nil
	}

	// leave 20 seconds free
	if deliveryTime < time.Now().Add(100 * time.Second).Unix() {
		return true, ErrDeliveryTime
	}
	return true, nil
}
