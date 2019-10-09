package ApiWS

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"../ApiRPC"
)

// WebSocket upgrade worker
var UpGrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }}

// Saving all node for every client connection.
var ClientPool = &NodePool{
	clients: make(map[int64]*Node),
	wt:      sync.RWMutex{},
}

// Saving the message which want sent to offline user.
// This chan has buffer, and the cap is 10000. Every element is an array,
// which saved the target user's id and bytes data of the message.
var delayMessageChat = make(chan [2]interface{}, 10000)

// Handle WebSocket upgrade request.
// Before upgrading, check token by send it to UserCenter through gRPC call. Then trying to get the user's connect node
// from ClientPool. if have not, new one and add into the ClientPool. Create a new connector for this connection and
// add the new connector into theNode's conns. Then start the send and receive data loop goroutine for this connector,
// and start the connector status watching loop goroutine.
func BeginChat(w http.ResponseWriter, r *http.Request) {
	// check the token by gRPC call , try to get the user's id.
	token := r.URL.Query().Get("authToken")
	userId, err := ApiRPC.CheckAuthToken(token)
	if nil != err {
		log.Printf("[error] <BeginChat> check auth token fail for user(%d)-address(%s), detail: %s",
			userId, r.RemoteAddr, err)
		w.WriteHeader(400)
		_, _ = w.Write([]byte("authToken authentication fail: " + err.Error()))
		return
	}

	// upgrade the connection with WebSocket protocol
	conn, err := UpGrader.Upgrade(w, r, nil)
	if nil != err {
		log.Printf("[error] <BeginChat> upgrade the request fail for user(%d)-address(%s), detail: %s",
			userId, r.RemoteAddr, err)
		w.WriteHeader(400)
		_, _ = w.Write([]byte("upgrade connection fail: " + err.Error()))
		return
	}

	// try to get the user's connect node from ClientPool. if have not, new one and add into the ClientPool.
	theNode, ok := ClientPool.Get(userId)
	if !ok {
		theNode = NewNode(userId)
		go theNode.ConnsWatchingLoop()

		ClientPool.Add(theNode)

	}

	// create a new connector for this connection and add the new connector into theNode's conns
	tempConnector := NewConnector(conn)
	theNode.AddConn(tempConnector)

	// start the goroutines for this connector
	go tempConnector.CloseWatchingLoop()
	go SendDateLoop(userId, tempConnector, r.RemoteAddr)
	go RecvDataLoop(userId, tempConnector, r.RemoteAddr)

	// query and send `DelayMessage` for this user and send them.
	if !ok {
		if messages, err := ApiRPC.GetDelayMessages(userId); nil == err {
			for _, message := range messages {
				theNode.AddMessageData(message)
			}
		}
	}
	log.Printf("[info] <BeginChat> new WebSockt connection for user(%d)-address(%s), connector(%p)",
		userId, r.RemoteAddr, tempConnector)
}

// Dispatch the chat message from ordinary user.
// If the `TypeId` of message is got error or not supported, it will send a error information to user client that sent
// the message.
func MessageDispatch(senderId int64, data []byte) {
	messageTypeId, err := GetRawJsonMessageTypeId(data)
	if nil != err {
		SendErrorMessage(senderId, 400, err, nil)
		return
	}
	switch messageTypeId {
	case UserChatMessageTypeId:
		SendUserChatMessage(senderId, data)
	case GroupChatMessageTypeId:
		SendGroupChatMessage(senderId, data)
	case SubscriptionMessageTypeId:
		SendSubscriptionMessage(senderId, data)
	default:
		SendErrorMessage(senderId, 400, ErrUnSupportMsgTypeId, nil)
	}
}

// Send error message to target user client, if the user is not online, the message will be save as delay message.
func SendErrorMessage(receiverId int64, code int, err error, rawMessage []byte) {
	tempMessage := ErrorMessage{
		BasicMessage: BasicMessage{
			TypeId:     ErrorMessageTypeId,
			SenderId:   SystemId,
			ReceiverId: receiverId,
			CreateTime: time.Now().Unix(),
		},
		Code:       code,
		Error:      err.Error(),
		RawMessage: rawMessage,
	}
	message, _ := json.Marshal(tempMessage)

	if node, ok := ClientPool.Get(receiverId); ok {
		log.Printf("[error] <SendErrorMessage> send to receiver(%d) {code: %d, error: %s}",
			receiverId, code, err.Error())
		node.AddMessageData(message)
	} else {
		delayMessageChat <- [2]interface{}{receiverId, message}
	}
}

// Send normal message to target user client, if the user is not online, it will record the message as a
// WaitSendMessage into database
func SendUserChatMessage(senderId int64, message []byte) {
	chatMessage := new(ChatMessage)
	err := json.Unmarshal(message, chatMessage)
	if nil != err {
		SendErrorMessage(senderId, 400, err, message)
		return
	}
	// check the chat message data legality
	code, err := checkUserChatMessageData(senderId, chatMessage)
	if nil != err {
		SendErrorMessage(senderId, code, err, message)
		return
	}

	recipientNode, ok := ClientPool.Get(chatMessage.ReceiverId)
	// check whether the recipient should receive the message
	code, err = checkWhetherRecipientShouldReceive(ok, recipientNode, senderId, chatMessage.ReceiverId)
	if nil != err {
		SendErrorMessage(senderId, code, err, message)
		return
	}

	// todo: check whether the message is an timing message
	yes, err := checkAndSaveTimingMessage(chatMessage)
	if nil != err {
		SendErrorMessage(senderId, code, err, message)
		return
	}
	if yes {
		return
	}

	// send the message or save as delay message
	if ok {
		recipientNode.AddMessageData(message)
	} else {
		// when the receiver not online, save the message as delay message.
		delayMessageChat <- [2]interface{}{chatMessage.ReceiverId, message}
	}

	// know the sender the message was send ok
	SendErrorMessage(senderId, 200, ErrSendMessageOk, message)
}

var (
	ErrSendMessageOk        = errors.New("the message send ok")
	ErrUserDisguise         = errors.New("the sender id is not identical, don't disguise other to send message")
	ErrUnSupportContentType = errors.New("the message content type is not support")
	ErrTextContentEmpty     = errors.New("the text content is empty")
	ErrPreviewPicEmpty      = errors.New("the media preview picture url is empty")
	ErrResourceURLEmpty     = errors.New("the media resource url is empty")
	ErrRecipientRefuseRecv  = errors.New("the recipient refuse receive the message")
	ErrFriendshipNotExisted = errors.New("your are not friend still")
	ErrRecipientNotExisted  = errors.New("the recipient is not existed")
)

// Check the chat message data legality.
// Requiring the sender id is identical, and the content type is supported, and the message real content not be empty
func checkUserChatMessageData(senderId int64, chatMessage *ChatMessage) (int, error) {
	if chatMessage.SenderId != senderId {
		return 400, ErrUserDisguise
	}

	switch chatMessage.ContentType {
	case TextContent:
		if chatMessage.Content == "" {
			return 400, ErrTextContentEmpty
		}
	case ImageContent, VoiceContent, VideoContent:
		if chatMessage.PreviewPic == "" {
			return 400, ErrPreviewPicEmpty
		}

		if chatMessage.ResourceUrl == "" {
			return 400, ErrResourceURLEmpty
		}
	default:
		return 400, ErrUnSupportContentType
	}

	return 200, nil
}

// Check whether the recipient should receive the message.
// Requiring the sender have an effective friendship with the receiver.
func checkWhetherRecipientShouldReceive(ok bool, receiverNode *Node, senderId, recipientId int64) (int, error) {
	switch ok {
	case true:
		// when the recipient is online
		if _, isBlack := receiverNode.BlackList.Load(senderId); isBlack {
			return 403, ErrRecipientRefuseRecv
		}
		if _, effective := receiverNode.Friends.Load(senderId); !effective {
			return 403, ErrFriendshipNotExisted
		}
	case false:
		// when the recipient is offline
		friends, blacklist, err := ApiRPC.GetUserFriendIdList(recipientId)
		if nil != err {
			return 404, ErrRecipientNotExisted
		}
		for _, id := range blacklist {
			if id == senderId {
				return 403, ErrRecipientRefuseRecv
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

func SendSubscriptionMessage(senderId int64, message []byte) {
	panic("implement the function")
}

func SendGroupChatMessage(senderId int64, message []byte) {
	panic("implement the function")
}

// Record the message for offline user.
// Save the message as delay message into database, when the user online again, send these message to the client.
// but if some thing error happened when save the delay message, it would only output the error log, would not block
// the program.
func SaveDelayMessageLoop() {
	for message := range delayMessageChat {
		receiverId := message[0].(int64)
		messageData := message[1].([]byte)
		err := ApiRPC.SaveDelayMessage(receiverId, messageData)
		if nil != err {
			log.Printf("[error] <SaveDelayMessageLoop> save data for recipient(%d) fail: %s",
				receiverId, err.Error())
		}
		log.Printf("[info] <SaveDelayMessageLoop> save data for recipient(%d) success: %s",
			receiverId, messageData)
	}
}

var (
	ErrDeliveryTime = errors.New("invalid delivery time")
)

// Check the value of 'DeliveryTime' field in the message. If the value is zero, meaning it is not a timing message.
// When it is a timing message, requiring the delivery time is after now at least 2 minute.
func checkAndSaveTimingMessage(message Message) (bool, error) {
	// don't support timing message at present
	return false, nil

	deliveryTime := message.GetDeliveryTime()
	if 0 == deliveryTime {
		return false, nil
	}

	// leave 20 seconds free
	if deliveryTime < time.Now().Add(100 * time.Second).Unix() {
		return true, ErrDeliveryTime
	}

	// todo: save timing message
	return true, nil
}
