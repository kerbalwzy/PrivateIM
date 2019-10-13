package MSGNode

import (
	"github.com/gorilla/websocket"
	"sync"
)

// structure for saving the int64 id
type Int64IdSet struct {
	data map[int64]struct{}
	wt   sync.RWMutex
}

func (obj *Int64IdSet) Add(id int64) {
	obj.wt.Lock()
	obj.data[id] = struct{}{}
	obj.wt.Unlock()
}

func (obj *Int64IdSet) Del(id int64) {
	obj.wt.Lock()
	delete(obj.data, id)
	obj.wt.Unlock()
}

func (obj *Int64IdSet) Exist(id int64) bool {
	obj.wt.RLock()
	_, ok := obj.data[id]
	obj.wt.RUnlock()
	return ok
}

func (obj *Int64IdSet) Keys() []int64 {
	obj.wt.RLock()
	temp := make([]int64, len(obj.data))
	index := 0
	for K := range obj.data {
		temp[index] = K
		index++
	}
	obj.wt.Unlock()
	return temp

}

// type for sort the group chat node by activeCount
type ActiveSorter []*GroupChatNode

// implement the sort.Interface on the ActiveSorter
func (obj ActiveSorter) Len() int {
	return len(obj)
}

func (obj ActiveSorter) Less(i, j int) bool {
	return obj[i].activeCount < obj[j].activeCount
}

func (obj ActiveSorter) Swap(i, j int) {
	obj[i], obj[j] = obj[j], obj[i]
}

// The connector for send and receive data with client
type Connector struct {
	conn        *websocket.Conn
	CloseSignal chan struct{}
	DataChan    chan []byte
}

func (obj *Connector) WriteMessage(messageType int, data []byte) error {
	return obj.conn.WriteMessage(messageType, data)
}

func (obj *Connector) ReadMessage() (messageType int, p []byte, err error) {
	return obj.conn.ReadMessage()
}

// Create a new connector for the connection.
func NewConnector(conn *websocket.Conn) *Connector {
	return &Connector{conn: conn, CloseSignal: make(chan struct{}), DataChan: make(chan []byte)}
}
