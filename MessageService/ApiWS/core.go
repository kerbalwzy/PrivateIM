package ApiWS

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"../ApiRPC"
	"../MSGNode"
)

// WebSocket upgrade worker
var UpGrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }}

// Handle WebSocket upgrade request.
// Before upgrading, check token by send it to UserCenter through gRPC call. Then trying to get the user's connect node
// from UserNodesPool. if have not, new one and add into the UserNodesPool. Create a new connector for this connection and
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

	// try to get the user's connect node from UserNodesPool. if have not, new one and add into the UserNodesPool.
	theNode, ok := MSGNode.GlobalUsers.Get(userId)
	if !ok {
		friends, blacklist, err := ApiRPC.GetUserFriendIdList(userId)
		if nil != err {
			log.Printf("[error] <BeginChat> get friends and blacklist for user(%d) fail, detail: %s", userId, err)
		}
		theNode = MSGNode.NewUserNode(userId, friends, blacklist)
		go theNode.ConnsWatchingLoop()

		MSGNode.GlobalUsers.Add(theNode)

	}

	// create a new connector for this connection and add the new connector into theNode's conns
	tempConnector := MSGNode.NewConnector(conn)
	theNode.AddConn(tempConnector)

	log.Printf("[info] <BeginChat> new WebSockt connection for user(%d)-address(%s), connector(%p)",
		userId, r.RemoteAddr, tempConnector)

	// start the goroutines for this connector
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
	case MSGNode.UserChatMessageTypeId:
		SendUserChatMessage(senderId, data)
	case MSGNode.GroupChatMessageTypeId:
		SendGroupChatMessage(senderId, data)
	case MSGNode.SubscriptionMessageTypeId:
		SendSubscriptionMessage(senderId, data)
	default:
		SendErrorMessage(senderId, 400, ErrUnSupportMsgTypeId, nil)
	}
}

// Send error message to target user client, if the user is not online, the message will be save as delay message.
func SendErrorMessage(receiverId int64, code int, err error, rawMessage []byte) {
	tempMessage := MSGNode.ErrorMessage{
		BasicMessage: MSGNode.BasicMessage{
			TypeId:     MSGNode.ErrorMessageTypeId,
			SenderId:   MSGNode.SystemId,
			ReceiverId: receiverId,
			CreateTime: time.Now().Unix(),
		},
		Code:       code,
		Error:      err.Error(),
		RawMessage: rawMessage,
	}
	message, _ := json.Marshal(tempMessage)

	if node, ok := MSGNode.GlobalUsers.Get(receiverId); ok {
		log.Printf("[info] <SendErrorMessage> send to receiver(%d) {code: %d, error: %s}",
			receiverId, code, err.Error())
		node.AddMessageData(message)
	} else {
		MSGNode.DelayMessageChat <- [2]interface{}{receiverId, message}
	}
}

// Send normal message to target user client, if the user is not online, it will record the message as a
// WaitSendMessage into database
func SendUserChatMessage(senderId int64, message []byte) {
	tempMessage := new(MSGNode.ChatMessage)
	err := json.Unmarshal(message, tempMessage)
	if nil != err {
		SendErrorMessage(senderId, 400, err, message)
		return
	}

	// check the chat message data legality
	code, err := checkUserChatMessageData(senderId, tempMessage)
	if nil != err {
		SendErrorMessage(senderId, code, err, message)
		return
	}

	recipientNode, ok := MSGNode.GlobalUsers.Get(tempMessage.ReceiverId)
	// check whether the recipient should receive the message
	code, err = checkWhetherReceiverShouldReceive(ok, recipientNode, senderId, tempMessage.ReceiverId)
	if nil != err {
		SendErrorMessage(senderId, code, err, message)
		return
	}

	// check whether the message is an timing message
	yes, err := checkWeatherTimingMessage(tempMessage)
	if nil != err {
		SendErrorMessage(senderId, code, err, message)
		return
	}

	// send the message or save as timing message or delay message
	tempMessage.SetCreateTime()
	message, _ = json.Marshal(tempMessage)

	// todo: save the timing message
	if yes {
		SendErrorMessage(senderId, 500, errors.New("not support at present"), message)
		return
	}

	if ok {
		recipientNode.AddMessageData(message)
	} else {
		// when the receiver not online, save the message as delay message.
		MSGNode.DelayMessageChat <- [2]interface{}{tempMessage.ReceiverId, message}
	}
	log.Printf("[info] <SendUsertempMessage> sender(%d) send a message to receiver(%d)",
		senderId, tempMessage.ReceiverId)

	// save the chat history and tell the sender the message was send ok
	SaveUserChatHistory(senderId, tempMessage.ReceiverId, message)
	SendErrorMessage(senderId, 200, ErrSendMessageOk, message)
}

var (
	ErrNotJoinedTheGroupChat = errors.New("you are not the member of the group chat")
	ErrGroupChatFindFail     = errors.New("find the target group chat fail, may not existed")
)

// Send the message to a group chat.
// In fact, is send the message to every other user whom joined the group chat.
func SendGroupChatMessage(senderId int64, message []byte) {
	tempMessage := new(MSGNode.ChatMessage)
	err := json.Unmarshal(message, tempMessage)
	if nil != err {
		SendErrorMessage(senderId, 400, err, message)
		return
	}

	// check the chat message data legality
	code, err := checkUserChatMessageData(senderId, tempMessage)
	if nil != err {
		SendErrorMessage(senderId, code, err, message)
		return
	}

	// get the group chat node
	groupChatNode, ok := MSGNode.GlobalGroupChats.Get(tempMessage.ReceiverId)

	if !ok {
		users, err := ApiRPC.GetGroupChatUsers(tempMessage.ReceiverId)
		if nil != err {
			log.Printf("[error] <SendGroupChatMessage> get group chat(%d) info fail: %s",
				tempMessage.ReceiverId, err)
			SendErrorMessage(senderId, 500, ErrGroupChatFindFail, message)
			return
		}
		groupChatNode = MSGNode.NewGroupChatNode(tempMessage.ReceiverId, users)

		MSGNode.GlobalGroupChats.Add(groupChatNode)
	}

	// and check weather the user has join the group chat
	if !groupChatNode.Users.Exist(senderId) {
		SendErrorMessage(senderId, 400, ErrNotJoinedTheGroupChat, message)
		return
	}

	// send the message to every members of the group chat
	tempMessage.SetCreateTime()
	message, _ = json.Marshal(tempMessage)

	for _, memberId := range groupChatNode.Users.Keys() {
		if userNode, ok := MSGNode.GlobalUsers.Get(memberId); ok {
			userNode.AddMessageData(message)
		} else {
			MSGNode.DelayMessageChat <- [2]interface{}{memberId, message}
		}
	}

	// add the activity count and save the group chat history
	groupChatNode.AddActiveCount()
	SaveGroupChatHistory(tempMessage.ReceiverId, message)
	log.Printf("[info] <SendGroupChatMessage> sender(%d) send a message to group chat(%d)", senderId, tempMessage.ReceiverId)

}

var (
	ErrSubscriptionFindFail = errors.New("find the target subscription fail, maybe not existed")
	ErrNotSubsManager       = errors.New("you are not the subscription's manager")
)

// The manager send a message to every fans of the subscription.
func SendSubscriptionMessage(senderId int64, message []byte) {
	tempMessage := new(MSGNode.SubscriptionMessage)
	err := json.Unmarshal(message, tempMessage)
	if nil != err {
		SendErrorMessage(senderId, 400, err, message)
		return
	}

	// check the subscription message legality
	code, err := checkSubscriptionMessageData(senderId, tempMessage)
	if nil != err {
		SendErrorMessage(senderId, code, err, message)
		return
	}

	// try to get the subscription node, and check weather the sender is the manager of the subscription
	subsNode, ok := MSGNode.GlobalSubscriptions.Get(tempMessage.ReceiverId)
	if !ok {
		manager, fans, err := ApiRPC.GetSubscriptionInfo(tempMessage.ReceiverId)
		if nil != err {
			SendErrorMessage(senderId, 500, ErrSubscriptionFindFail, message)
			return
		}
		if manager != senderId {
			SendErrorMessage(senderId, 403, ErrNotSubsManager, message)
		}
		subsNode = MSGNode.NewSubsNode(tempMessage.ReceiverId, manager, fans)
		MSGNode.GlobalSubscriptions.Add(subsNode)
	}

	if subsNode.ManagerId != senderId {
		SendErrorMessage(senderId, 403, ErrNotSubsManager, message)
		return
	}

	// send the message to every fans of the subscription
	tempMessage.SetCreateTime()
	message, _ = json.Marshal(tempMessage)

	for _, memberId := range subsNode.Fans.Keys() {
		if userNode, ok := MSGNode.GlobalUsers.Get(memberId); ok {
			userNode.AddMessageData(message)
		} else {
			MSGNode.DelayMessageChat <- [2]interface{}{memberId, message}
		}
	}

	// save the subscription message history
	SaveSubscriptionHistory(tempMessage.ReceiverId, message)
	log.Printf("[info] <SendSubscriptionMessage> subscription(%d) send a message to fans", tempMessage.ReceiverId)

}

// Keep send the data to client by connector, when have a new message for it.
// When the connector was closed, stop the goroutine.
func SendDateLoop(userId int64, connector *MSGNode.Connector, clientAddr string) {
	defer func() {
		recover()
	}()

	userConnectInfo := fmt.Sprintf("user(%d)-address(%s)-connector(%p)", userId, clientAddr, connector)
	log.Printf("[info] <SendDateLoop> start a send data goroutine for %s", userConnectInfo)

	for {
		select {
		case <-connector.CloseSignal:
			log.Printf("[info] <SendDateLoop> stop a send data goroutine for %s", userConnectInfo)
			return
		case data := <-connector.DataChan:
			err := connector.WriteMessage(websocket.TextMessage, data)
			if nil != err {
				MSGNode.DelayMessageChat <- [2]interface{}{userId, data}
				close(connector.CloseSignal)
				return
			}
			//log.Printf("[info] <SendDateLoop> send data to %s", userConnectInfo)
		}
	}
}

// Keep trying to receive message from the client by connector.
// When the connector was closed, stop the goroutine.
func RecvDataLoop(userId int64, connector *MSGNode.Connector, clientAddr string) {
	defer func() {
		recover()
	}()

	userConnectInfo := fmt.Sprintf("user(%d)-address(%s)-connector(%p)", userId, clientAddr, connector)
	log.Printf("[info] <RecvDataLoop> start a recv data goroutine for %s", userConnectInfo)
	for {
		select {
		case <-connector.CloseSignal:
			log.Printf("[info] <RecvDataLoop> stop a recv data goroutine for %s", userConnectInfo)
			return
		default:
			_, data, err := connector.ReadMessage()
			if err != nil {
				log.Printf("[error] <RecvDataLoop> recevie data fail from %s, detail: %s", userConnectInfo, err.Error())
				close(connector.CloseSignal)
				return
			}
			//log.Printf("[info] <SendDateLoop> recv data from  %s", userConnectInfo)
			MessageDispatch(userId, data)
		}
	}
}

// Record the message for offline user.
// Save the message as delay message into database, when the user online again, send these message to the client.
// but if some thing error happened when save the delay message, it would only output the error log, would not block
// the program.
func SaveDelayMessageLoop() {
	for message := range MSGNode.DelayMessageChat {
		receiverId := message[0].(int64)
		messageData := message[1].([]byte)
		err := ApiRPC.SaveDelayMessage(receiverId, messageData)
		if nil != err {
			log.Printf("[error] <SaveDelayMessageLoop> save data for recipient(%d) fail: %s",
				receiverId, err.Error())
		}
		log.Printf("[info] <SaveDelayMessageLoop> save data for recipient(%d) success", receiverId)
	}
}

// Record the normal user chat history
func SaveUserChatHistory(senderId, receiverId int64, message []byte) {
	if senderId > receiverId {
		senderId, receiverId = receiverId, senderId
	}
	joinId := fmt.Sprintf("%d_%d", senderId, receiverId)
	err := ApiRPC.SaveUserChatHistory(joinId, message)
	if nil != err {
		log.Printf("[error] <SaveUserChatHistory> save chat history for (%s) fail: %s", joinId, err.Error())
	}
}

// Record the group chat history
func SaveGroupChatHistory(groupId int64, message []byte) {
	err := ApiRPC.SaveGroupChatHistory(groupId, message)
	if nil != err {
		log.Printf("[error] <SaveGroupChatHistory> save group chat history for (%d) fail: %s", groupId, err.Error())
	}
}

// Record the subscription message history
func SaveSubscriptionHistory(subsId int64, message []byte) {
	err := ApiRPC.SaveSubscriptionMessageHistory(subsId, message)
	if nil != err {
		log.Printf("[error] <SaveSubscriptionHistory> save subscription history for (%d) fail: %s", subsId, err.Error())

	}
}
