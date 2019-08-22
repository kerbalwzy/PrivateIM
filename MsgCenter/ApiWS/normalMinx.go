package ApiWS

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"

	"../DataLayer"
)

type NodePool struct {
	clients map[int64]*Node
	wt      sync.RWMutex
}

// add a new client note
func (obj *NodePool) Add(node *Node) {
	obj.wt.Lock()
	obj.clients[node.Id] = node
	obj.wt.Unlock()
}

// get a client note by id
func (obj *NodePool) Get(id int64) (*Node, bool) {
	obj.wt.RLock()
	node, ok := obj.clients[id]
	obj.wt.RUnlock()
	return node, ok

}

// delete a client note
func (obj *NodePool) Del(node *Node) {
	obj.wt.Lock()
	delete(obj.clients, node.Id)
	obj.wt.Unlock()

}

// client node
type Node struct {
	Id        int64 // userId
	conn      *websocket.Conn
	dataQueue chan []byte
}

// send message loop, send message to client
// If fail will remove the node from ClientPool, and record message into database as "waiting send message",
// When the client online again, then try send again.
func (obj *Node) SendLoop() {
	defer func() {
		ClientPool.Del(obj)
		obj.conn.Close()
	}()
	for {
		select {
		case data := <-obj.dataQueue:
			err := obj.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				DataLayer.RecordWaitSendMessage(obj.Id, data)
				return
			}
		}
	}
}

// recv message loop, recv message from client
// If recv error, will remove the node from ClientPool.
// Send the message to target user by Dispatch function,
// when dispatch error, it will also send back the error information to client
func (obj *Node) RecvLoop() {
	defer func() {
		ClientPool.Del(obj)
		obj.conn.Close()
	}()
	for {
		_, data, err := obj.conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}
		err = Dispatch(data)
		if nil != err {
			log.Println(err.Error())
			msg := ErrorMessage{Code: 400, Error: err.Error()}
			data, _ = msg.ToJson()
			obj.dataQueue <- data
		}
	}
}

// create a ClientPool to save nodes
func NewNodePool() *NodePool {
	return &NodePool{
		clients: make(map[int64]*Node),
		wt:      sync.RWMutex{},
	}
}
