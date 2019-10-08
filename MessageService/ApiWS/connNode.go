package ApiWS

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"

	"../ApiRPC"
)

type NodePool struct {
	clients map[int64]*Node
	wt      sync.RWMutex
}

// Get a client node by id from ClientPool
func (obj *NodePool) Get(id int64) (*Node, bool) {
	obj.wt.RLock()
	node, ok := obj.clients[id]
	obj.wt.RUnlock()
	return node, ok

}

// Add a new client node into ClientPool, if the user is had a node, replace and close the old one
func (obj *NodePool) Add(node *Node) {
	obj.wt.Lock()
	obj.clients[node.Id] = node
	obj.wt.Unlock()

	// todo test code used in separate development, need remove later
	fmt.Printf("client node list: \n")
	for k, v := range obj.clients {
		fmt.Printf("\t%d, %p\n", k, v)
	}
}

// Delete a client node from ClientPool, and close the connection of the node.
func (obj *NodePool) Del(node *Node) {
	obj.wt.Lock()
	delete(obj.clients, node.Id)
	obj.wt.Unlock()

	// todo test code used in separate development, need remove later
	fmt.Printf("client node list: \n")
	for k, v := range obj.clients {
		fmt.Printf("\t%d, %p\n", k, v)
	}
}

// The connection node of the client for receiving and sending messages.
type Node struct {
	Id int64

	conns     [3]*Connector // the connector array
	connCount int           // the count of the connector
	wt        sync.Mutex    // the lock for operating the 'count' field

	Friends   sync.Map
	BlackList sync.Map
}

// Create a new node instance for the connection
func NewNode(userId int64, conn *websocket.Conn) *Node {
	connector := &Connector{conn: conn, CloseSignal: make(chan struct{}), DataChan: make(chan []byte)}
	conns := [3]*Connector{connector}
	node := &Node{
		Id: userId,

		conns:     conns,
		connCount: 0,
		wt:        sync.Mutex{},

		Friends:   sync.Map{},
		BlackList: sync.Map{}}

	// load the user's friends and blacklist
	friends, blacklist, err := ApiRPC.GetUserFriendIdList(userId)
	if nil == err {
		for _, id := range friends {
			node.Friends.Store(id, struct{}{})
		}

		for _, id := range blacklist {
			node.BlackList.Store(id, struct{}{})
		}
	} else {
		log.Printf("[error] <NewNode> load friends and blacklist for user(%d) fail, detail: %s", userId, err)
	}

	return node
}

// Add a connector for the node, the max count of connectors is 3.
func (obj *Node) AddConn(conn *Connector) {
	obj.wt.Lock()
	var oldestConn *Connector
	oldestConn, obj.conns[2], obj.conns[1], obj.conns[0] = obj.conns[2], obj.conns[1], obj.conns[0], conn
	if oldestConn != nil {
		close(oldestConn.CloseSignal)
	}
	if obj.connCount < 3 {
		obj.connCount++
	}
	obj.wt.Unlock()

}

// Watching the connectors of the node, when a connector is closed, remove it from the node and reduce the value
// of count of the connectors whom are belong to the node. When the node have not connectors, don't save the node
// in ClientPool anymore.
func (obj *Node) ConnsWatchLoop() {
	for {
		obj.wt.Lock()
		for index, conn := range obj.conns {
			if nil != conn {
				if _, ok := <-conn.CloseSignal; !ok {
					obj.connCount--
					switch index {
					case 0:
						obj.conns[0], obj.conns[1], obj.conns[2] = obj.conns[1], obj.conns[2], nil
					case 1:
						obj.conns[1], obj.conns[2] = obj.conns[2], nil
					case 3:
						obj.conns[2] = nil
					}
				}
			}
		}
		// when the node have not connectors, don't save the node in ClientPool anymore.
		if obj.connCount == 0 {
			ClientPool.Del(obj)
		}
		obj.wt.Unlock()

		time.Sleep(1 * time.Second)
	}
}

// add the message data to the node's every connector data channel.If have not any connector here, the message would
// be saved as delay message
func (obj *Node) AddMessageData(data []byte) {
	if obj.connCount == 0 {
		delayMessageChat <- [2]interface{}{obj.Id, data}
		return
	}
	for _, conn := range obj.conns {
		if nil != conn {
			conn.DataChan <- data
		}
	}

}

// The connector for send and receive data with client really
type Connector struct {
	conn        *websocket.Conn
	CloseSignal chan struct{}
	DataChan    chan []byte
}

// Watching the connector close signal, keep the connect safe.
func (obj *Connector) CloseWatchLoop() {
	for {
		select {
		case <-obj.CloseSignal:
			err := obj.conn.Close()
			log.Printf("[error] <Transfer.CloseWatchLoop> : %s", err)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

// Keep send the data to client by connector, when have a new message for it.
// When the connector was closed, stop the goroutine.
func SendDateLoop(userId int64, connector *Connector) {
	for {
		select {
		case <-connector.CloseSignal:
			return
		case data := <-connector.DataChan:
			err := connector.conn.WriteMessage(websocket.TextMessage, data)
			if nil != err {
				delayMessageChat <- [2]interface{}{userId, data}
				close(connector.CloseSignal)
			}
			log.Printf("[info] <SendDateLoop> send data to user(%d)", userId)
		}
	}
}

// Keep trying to receive message from the client by connector.
// When the connector was closed, stop the goroutine.
func RecvDataLoop(userId int64, connector *Connector) {
	for {
		select {
		case <-connector.CloseSignal:
			return
		default:
			_, data, err := connector.conn.ReadMessage()
			if err != nil {
				log.Printf("[error] <RecvLoop> recevie data fail from user(%d), detail: %s", userId, err.Error())
				close(connector.CloseSignal)
			}
			MessageDispatch(userId, data)
		}
	}
}
