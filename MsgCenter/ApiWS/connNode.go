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

// Delete a client node from ClientPool
func (obj *NodePool) Del(node *Node) {
	obj.wt.Lock()
	delete(obj.clients, node.Id)
	obj.wt.Unlock()

}

var nodePool = &NodePool{
	clients: make(map[int64]*Node),
	wt:      sync.RWMutex{},
}

// Create a ClientPool to save nodes, (Singleton Design Pattern)
func NewNodePool() *NodePool {
	return nodePool
}

// client node
type Node struct {
	Id           int64 // userId
	conn         *websocket.Conn
	messageQueue chan DataLayer.Message
}

// Send message loop, send message to client.
// If fail will remove the node from ClientPool,
// and move the message into `WaitSendChan`.
func (obj *Node) SendLoop() {
	defer func() {
		ClientPool.Del(obj)
		_ = obj.conn.Close()
	}()
	for {
		select {
		case message := <-obj.messageQueue:
			err := obj.conn.WriteMessage(websocket.TextMessage, message.ToJson())
			if err != nil {
				log.Println(err.Error())
				WaitSendChan <- message
				return
			}
		}
	}
}

// Recv message loop, recv message from client
// If recv error, will remove the node from ClientPool.
// Send the message to target user by Dispatch function,
// when dispatch error, it will also send back the error information to client
func (obj *Node) RecvLoop() {
	defer func() {
		ClientPool.Del(obj)
		_ = obj.conn.Close()
	}()
	for {
		_, data, err := obj.conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}
		// dispatch the chat message from ordinary user
		ChatDispatch(obj.Id, data)
	}
}

func NewNode(userId int64, conn *websocket.Conn) (*Node, error) {
	return nil, nil
}
