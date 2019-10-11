package ApiWS

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
	obj.wt.Unlock()
	return ok
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

// The connector for send and receive data with client really
type Connector struct {
	conn        *websocket.Conn
	CloseSignal chan struct{}
	DataChan    chan []byte
}

// Create a new connector for the connection.
func NewConnector(conn *websocket.Conn) *Connector {
	return &Connector{conn: conn, CloseSignal: make(chan struct{}), DataChan: make(chan []byte)}
}
