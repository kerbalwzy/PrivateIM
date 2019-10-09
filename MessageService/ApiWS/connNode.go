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

	conns         [3]*Connector // the connector array. every node can have 3 connector max.
	connCount     int           // the count of the connector
	wt            sync.Mutex    // the lock for operating the 'count' field
	connsWatching bool          //

	Friends   sync.Map
	BlackList sync.Map
}

// Create a new node instance for the user's connection
func NewNode(userId int64) *Node {
	node := &Node{
		Id: userId,

		conns:     [3]*Connector{},
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

// Create a new connector for the connection.
func NewConnector(conn *websocket.Conn) *Connector {
	return &Connector{conn: conn, CloseSignal: make(chan struct{}), DataChan: make(chan []byte)}
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
	if !obj.connsWatching {
		obj.connsWatching = true
	}
	obj.wt.Unlock()
	log.Printf("[info] <Node.AddConn> add one new connetor for user(%d), the connector count= %d",
		obj.Id, obj.connCount)
}

// Watching the connectors of the node, when a connector is closed, remove it from the node and reduce the value
// of count of the connectors whom are belong to the node. When the node have not connectors, don't save the node
// in ClientPool anymore.
func (obj *Node) ConnsWatchingLoop() {
	log.Printf("[info] <Node.ConnsWatchingLoop> start a node's conns watching goroutine")
	for {
		obj.wt.Lock()
		for index, conn := range obj.conns {
			if nil != conn {
				select {
				case <-time.After(1 * time.Second):
					continue
				case <-conn.CloseSignal:
					obj.connCount--
					switch index {
					case 0:
						obj.conns[0], obj.conns[1], obj.conns[2] = obj.conns[1], obj.conns[2], nil
					case 1:
						obj.conns[1], obj.conns[2] = obj.conns[2], nil
					case 3:
						obj.conns[2] = nil
					}
					log.Printf("[info] <Node.ConnsWatchingLoop> reduce one connector of user(%d),"+
						" the connector count= %d", obj.Id, obj.connCount)
				}
			}
		}
		obj.wt.Unlock()

		time.Sleep(1 * time.Second)

		// when the node have not connectors and the loop is not the first start.
		// Don't save the node in ClientPool anymore, and stop this watching goroutine.
		if obj.connCount <= 0 && obj.connsWatching {
			ClientPool.Del(obj)
			log.Printf("[info] <Node.ConnsWatchingLoop> stop a node's conns watching goroutine")
			return
		} else {
			// todo test code used in separate development, need remove later
			log.Printf("[info] <Node.ConnsWatchingLoop> -- user(%d) connectors: %v\n", obj.Id, obj.conns)
		}
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

// Watching the connector close signal, keep the connect safe. When receive the close signal, would close the connect
// and stop this watching.
func (obj *Connector) CloseWatchingLoop() {
	log.Printf("[info] <Connetor.CloseWatchingLoop> start a connector close watching goroutine")
	for {
		select {
		case <-obj.CloseSignal:
			err := obj.conn.Close()
			if nil != err {
				log.Printf("[error] <Connetor.CloseWatchingLoop> : %s", err)
			}
			log.Printf("[info] <Connetor.CloseWatchingLoop> stop a connector close watching goroutine")
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

// Keep send the data to client by connector, when have a new message for it.
// When the connector was closed, stop the goroutine.
func SendDateLoop(userId int64, connector *Connector, clientAddr string) {
	defer func() {
		_ = recover()
	}()

	userConnectInfo := fmt.Sprintf("user(%d)-address(%s)-connector(%p)", userId, clientAddr, connector)
	log.Printf("[info] <SendDateLoop> start a send data goroutine for %s", userConnectInfo)

	for {
		select {
		case <-connector.CloseSignal:
			log.Printf("[info] <SendDateLoop> stop a send data goroutine for %s", userConnectInfo)
			return
		case data := <-connector.DataChan:
			err := connector.conn.WriteMessage(websocket.TextMessage, data)
			if nil != err {
				delayMessageChat <- [2]interface{}{userId, data}
				close(connector.CloseSignal)
				return
			}
			log.Printf("[info] <SendDateLoop> send data to %s", userConnectInfo)
		}
	}
}

// Keep trying to receive message from the client by connector.
// When the connector was closed, stop the goroutine.
func RecvDataLoop(userId int64, connector *Connector, clientAddr string) {
	defer func() {
		_ = recover()
	}()

	userConnectInfo := fmt.Sprintf("user(%d)-address(%s)-connector(%p)", userId, clientAddr, connector)
	log.Printf("[info] <RecvDataLoop> start a recv data goroutine for %s", userConnectInfo)
	for {
		select {
		case <-connector.CloseSignal:
			log.Printf("[info] <RecvDataLoop> stop a recv data goroutine for %s", userConnectInfo)
			return
		default:
			_, data, err := connector.conn.ReadMessage()
			if err != nil {
				log.Printf("[error] <RecvDataLoop> recevie data fail from %s, detail: %s", userConnectInfo, err.Error())

				close(connector.CloseSignal)
				return
			}
			log.Printf("[info] <SendDateLoop> recv data from  %s", userConnectInfo)
			MessageDispatch(userId, data)
		}
	}
}
