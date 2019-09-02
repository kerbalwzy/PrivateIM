package ApiWS

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"

	"../DataLayer"
)

type NodePool struct {
	clients map[int64]*Node
	wt      sync.RWMutex
}

// Add a new client node into ClientPool
func (obj *NodePool) Add(node *Node) {
	obj.wt.Lock()
	obj.clients[node.Id] = node
	obj.wt.Unlock()
}

// Get a client node by id from ClientPool
func (obj *NodePool) Get(id int64) (*Node, bool) {
	obj.wt.RLock()
	node, ok := obj.clients[id]
	obj.wt.RUnlock()
	return node, ok

}

// Delete a client node from ClientPool, and close the connection of the node.
func (obj *NodePool) Del(node *Node) {
	obj.wt.Lock()
	delete(obj.clients, node.Id)
	obj.wt.Unlock()
	_ = node.conn.Close()

}

var nodePool = &NodePool{
	clients: make(map[int64]*Node),
	wt:      sync.RWMutex{},
}

// Create a ClientPool to save nodes, (Singleton Design Pattern)
func NewNodePool() *NodePool {
	return nodePool
}

// The connection node of the client for receiving and sending messages.
// The "Id" of the node saved the user's id which using the client, the "conn" is
// used to send or receive data, the "messageQueue" saved the message those need
// send to the client, the "Friends" saved the id of user's friends, the "BlackList"
// saved the id of user who is marked black by the user.
type Node struct {
	Id          int64 // userId
	conn        *websocket.Conn
	messageChan chan []byte
	Friends     sync.Map
	BlackList   sync.Map
}

// Sending the message to the client which the node marked.
// Continuously try to get messages from the "messageChan" of the node and send them
// to the client immediately. If send fail, moving the message into "WaitSendChan",
// closing the connection and removing from node pool.
func (obj *Node) SendLoop() {
	defer func() {
		ClientPool.Del(obj)
	}()
	for message := range obj.messageChan {
		log.Printf("SendMessage: send data to client(user_id=%d)", obj.Id)
		err := obj.conn.WriteMessage(websocket.TextMessage, message)
		if nil != err {
			log.Printf("Error: send data fail to client(user_id=%d), error detail: %s", obj.Id, err.Error())
			WaitSendChan <- [2]interface{}{obj.Id, message}
			return
		}
	}
}

// Receiving message from the client which the node marked.
// Continuously try to receive data from the "conn" of the node, if having error
// happened when receiving, will close the connection and remove from node pool.
// The data will be handed over to the "Chat Dispatch" function for subsequent processing
func (obj *Node) RecvLoop() {
	defer func() {
		ClientPool.Del(obj)
	}()
	for {
		_, data, err := obj.conn.ReadMessage()
		log.Printf("RecvMessage: receive data from client(user_id=%d)", obj.Id)
		if err != nil {
			log.Printf("Error: recevie data fail from client(user_id=%d), error detail: %s", obj.Id, err.Error())
			return
		}
		ChatMessageDispatch(obj.Id, data)
	}
}

//Create a new node instance for the connection
func NewNode(userId int64, conn *websocket.Conn) *Node {
	node := &Node{
		Id: userId, conn: conn,
		messageChan: make(chan []byte),
		Friends:     sync.Map{},
		BlackList:   sync.Map{}}
	if friends, err := DataLayer.MongoQueryFriendsId(userId); nil == err {
		for _, id := range friends {
			node.Friends.Store(id, struct{}{})
		}
	}
	if blackList, err := DataLayer.MongoQueryBlackList(userId); nil == err {
		for _, id := range blackList {
			node.BlackList.Store(id, struct{}{})
		}
	}
	return node
}
