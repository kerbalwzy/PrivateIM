package ApiWS

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"

	"../ApiRPC"

	conf "../Config"
)

// The node of the user's clients for receiving and sending messages.
type UserNode struct {
	Id int64

	conns         [3]*Connector // the connector array. every node can have 3 connector max.
	connCount     int           // the count of the connector
	connsWatching bool          // mark weather the ConnsWatchingLoop goroutine is started
	wt            sync.Mutex    // the lock for operating the 'count' field

	Friends   Int64IdSet // the id of other users whom are the user's friend
	BlackList Int64IdSet // the id of other users whom are in the user' blacklist
}

// Add a connector for the node, the max count of connectors is 3.
func (obj *UserNode) AddConn(conn *Connector) {
	obj.wt.Lock()

	var oldestConn *Connector
	oldestConn, obj.conns[2], obj.conns[1], obj.conns[0] = obj.conns[2], obj.conns[1], obj.conns[0], conn
	if oldestConn != nil {
		close(oldestConn.CloseSignal)
		_ = oldestConn.conn.Close()
	}

	if !obj.connsWatching {
		obj.connsWatching = true
	}

	if obj.connCount < 3 {
		obj.connCount++
	}

	obj.wt.Unlock()

	log.Printf("[info] <Node.AddConn> add one new connetor for user(%d), the connector count= %d",
		obj.Id, obj.connCount)
}

// add the message data to the node's every connector data channel.If have not any connector here, the message would
// be saved as delay message
func (obj *UserNode) AddMessageData(data []byte) {
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

// Watching the connectors of the node, when a connector is closed, remove it from the node and reduce the value
// of count of the connectors whom are belong to the node. When the node have not connectors, don't save the node
// in UserNodesPool anymore.
func (obj *UserNode) ConnsWatchingLoop() {
	log.Printf("[info] <Node.ConnsWatchingLoop> start a node's conns watching goroutine")
	for {
		obj.wt.Lock()

		for index, conn := range obj.conns {
			if nil != conn {
				// create a timeout monitor for this check
				timeOut := time.NewTimer(time.Second * 1)
				select {
				case <-timeOut.C:
					break // only break this select work
				case <-conn.CloseSignal:
					_ = conn.conn.Close()

					// actively close this timeout monitoring
					timeOut.Stop()

					obj.connCount--
					switch index {
					case 0:
						obj.conns[0], obj.conns[1], obj.conns[2] = obj.conns[1], obj.conns[2], nil
					case 1:
						obj.conns[1], obj.conns[2] = obj.conns[2], nil
					case 2:
						obj.conns[2] = nil
					}

					log.Printf("[info] <Node.ConnsWatchingLoop> reduce one connector of user(%d),"+
						" the connector count= %d", obj.Id, obj.connCount)
				}
			}
		}

		// when the node have not connectors and the loop is not the first start.
		// Don't save the node in UserNodesPool anymore, and stop this watching goroutine.
		if obj.connCount <= 0 && obj.connsWatching {
			GlobalUsers.Del(obj)
			log.Printf("[info] <Node.ConnsWatchingLoop> stop a node's conns watching goroutine")

			obj.wt.Unlock()
			return
		}

		obj.wt.Unlock()
		// todo test code used in separate development, need remove later
		//log.Printf("[info] <Node.ConnsWatchingLoop> -- user(%d) connectors: %v\n", obj.Id, obj.conns)
		time.Sleep(1 * time.Second)
	}
}

// Create a new node instance for the user's connection
func NewUserNode(userId int64) *UserNode {
	node := new(UserNode)
	node.Id = userId
	node.Friends = Int64IdSet{data: map[int64]struct{}{}, wt: sync.RWMutex{}}
	node.BlackList = Int64IdSet{data: map[int64]struct{}{}, wt: sync.RWMutex{}}

	// load the user's friends and blacklist
	friends, blacklist, err := ApiRPC.GetUserFriendIdList(userId)
	if nil == err {
		for _, id := range friends {
			node.Friends.data[id] = struct{}{} // don't need lock here
		}

		for _, id := range blacklist {
			node.BlackList.data[id] = struct{}{} // don't need lock here
		}
	} else {
		log.Printf("[error] <NewNode> load friends and blacklist for user(%d) fail, detail: %s", userId, err)
	}
	log.Printf("[info] <NewUserNode> new a user client node, id= %d", userId)
	return node
}

type UserNodePool struct {
	clients map[int64]*UserNode
	wt      sync.RWMutex
}

// Get a client node by id from UserNodesPool
func (obj *UserNodePool) Get(id int64) (*UserNode, bool) {
	obj.wt.RLock()
	node, ok := obj.clients[id]
	obj.wt.RUnlock()
	return node, ok

}

// Add a new client node into UserNodesPool, if the user is had a node, replace and close the old one
func (obj *UserNodePool) Add(node *UserNode) {
	obj.wt.Lock()
	obj.clients[node.Id] = node
	obj.wt.Unlock()

	// todo test code used in separate development, need remove later
	fmt.Printf("client node list: \n")
	for k, v := range obj.clients {
		fmt.Printf("\t%d, %p\n", k, v)
	}
}

// Delete a client node from UserNodesPool, and close the connection of the node.
func (obj *UserNodePool) Del(node *UserNode) {
	obj.wt.Lock()
	delete(obj.clients, node.Id)
	obj.wt.Unlock()

	// todo test code used in separate development, need remove later
	fmt.Printf("client node list: \n")
	for k, v := range obj.clients {
		fmt.Printf("\t%d, %p\n", k, v)
	}
}

// Delete all the client node.
// It would close every connectors and stop the data translate goroutines for every nodes before reset the clients map.
func (obj *UserNodePool) CleanAll() {
	obj.wt.Lock()
	// close every node's connectors, stop the data translate goroutines
	for _, node := range obj.clients {
		for _, conn := range node.conns {
			close(conn.CloseSignal)
		}
	}

	// reset the clients map
	obj.clients = make(map[int64]*UserNode)
	obj.wt.Unlock()

}

// ---------------------------------------------------------------------------------

// The group chat node, saving some information of the group chat
type GroupChatNode struct {
	Id    int64      // group chat id
	Users Int64IdSet // the id of users whom joined the group chat

	activeCount int       // This value is increased every time a message is sent to this group chat
	initTime    time.Time // The group chat node initial timestamp, unit is second
	wt          sync.Mutex
}

// Increased the activeCount of the group chat node, keep the concurrent security.
func (obj *GroupChatNode) AddActiveCount() {
	obj.wt.Lock()
	obj.activeCount++
	obj.wt.Unlock()
}

// Reset the value of activeCount of the group chat node
func (obj *GroupChatNode) ResetActiveCount() {
	obj.wt.Lock()
	obj.activeCount = 0
	obj.wt.Unlock()
}

var (
	ErrGroupChatFindFail     = errors.New("find the target group chat fail, may not existed")
	ErrNotJoinedTheGroupChat = errors.New("you are not the member of the group chat")
)

// Initial a new group chat node.
func NewGroupChatNode(id int64) (*GroupChatNode, error) {
	tempGroupChat := new(GroupChatNode)
	tempGroupChat.Id = id
	tempGroupChat.initTime = time.Now()
	tempGroupChat.Users = Int64IdSet{data: map[int64]struct{}{}, wt: sync.RWMutex{}}

	userIdSlice, err := ApiRPC.GetGroupChatUsers(id)

	// load the id of users whom haven joined the group chat
	if nil == err {
		for _, id := range userIdSlice {
			tempGroupChat.Users.data[id] = struct{}{} // don't need lock here
		}
	} else {
		log.Printf("[error] <NewGroupChatNode> load users for group chat(%d) fail, detail: %s", id, err)
		return nil, ErrGroupChatFindFail
	}
	log.Printf("[info] <NewGroupChatNode> new a group chat node, the id= %d", id)
	return tempGroupChat, nil
}

// The group chat node pool. Save and manage the group chat nodes
type GroupChatNodePool struct {
	groups map[int64]*GroupChatNode
	wt     sync.RWMutex
}

// Get a group chat node from the GroupChatPool
func (obj *GroupChatNodePool) Get(id int64) (*GroupChatNode, bool) {
	obj.wt.RLock()
	groupChat, ok := obj.groups[id]
	obj.wt.RUnlock()
	return groupChat, ok
}

// Add a group chat node into the GroupChatPool
func (obj *GroupChatNodePool) Add(groupChat *GroupChatNode) {
	obj.wt.Lock()
	obj.groups[groupChat.Id] = groupChat
	obj.wt.Unlock()

}

// Delete a group chat node from the GroupChatPool
func (obj *GroupChatNodePool) Del(id int64) {
	obj.wt.Lock()
	delete(obj.groups, id)
	obj.wt.Unlock()
}

// Reset the groups map
func (obj *GroupChatNodePool) CleanAll() {
	obj.wt.Lock()
	obj.groups = make(map[int64]*GroupChatNode)
	obj.wt.Unlock()
}

// Clear up nodes whose lifetime exceeds the limit or low activity.
// Working on at NN:00:00 every day by config.
func (obj *GroupChatNodePool) CleanGroupChatLoop() {
	log.Printf("[info] <GroupChatNodePool.CleanGroupChatLoop> start the group chat pool clear up goroutine")
	for {
		// get next cleaning up execute time
		todayDateStr := time.Now().Format("2006-01-02")
		todayZeroH, _ := time.ParseInLocation("2006-01-02", todayDateStr, time.Local)
		tomorrowZeroH := todayZeroH.AddDate(0, 0, 1)
		nextCleanTime := tomorrowZeroH.Add(conf.GroupChatNodeCleanTime * time.Hour)

		// waiting for cleaning up, would blocking here
		select {
		case <-time.After(nextCleanTime.Sub(time.Now())):
			obj.CleanByLifeTime()
			obj.CleanByActiveCount()
		}
	}
}

// Clear up nodes whose lifetime exceeds the limit
func (obj *GroupChatNodePool) CleanByLifeTime() {
	targets := make([]int64, 0)
	timeNow := time.Now()
	count := 0

	obj.wt.Lock()
	defer obj.wt.Unlock()
	// find all the group chat nodes whose life time exceeds the limit
	for id, groupChatNode := range obj.groups {
		lifeTime := timeNow.Sub(groupChatNode.initTime)
		if lifeTime.Seconds() > conf.GroupChatNodeLifeTime {
			targets = append(targets, id)
		}
	}

	// clear up the nodes
	for _, id := range targets {
		delete(obj.groups, id)
		count++
	}
	log.Printf("[info] <GroupChatNodePool.CleanByLifeTime> clear up the group chat node, count= %d", count)

}

// Clear up nodes whose activity percentage is under the limit
func (obj *GroupChatNodePool) CleanByActiveCount() {
	obj.wt.Lock()
	defer obj.wt.Unlock()

	nodeCount := len(obj.groups)
	if nodeCount <= conf.GroupChatNodeLowActivityCleanStartLimit {
		log.Printf("[info] <GroupChatNodePool.CleanByActiveCount> don't need clean")
		return
	}

	tempActiveSorter := make(ActiveSorter, 0, nodeCount)
	// add element into the sorter, and sort the data
	for _, groupChatNode := range obj.groups {
		tempActiveSorter = append(tempActiveSorter, groupChatNode)
	}
	sort.Sort(tempActiveSorter)

	// get the count of the nodes that needs to be cleaned up, and clean those nodes
	cleanCount := float64(nodeCount) * conf.GroupChatNodeLowActivityCleanPercentage / 100
	cleanCount = math.Ceil(cleanCount)
	for _, node := range tempActiveSorter[:int(cleanCount)+1] {
		delete(obj.groups, node.Id)
	}

	// reset the active count of every group chat node
	for _, node := range obj.groups {
		node.ResetActiveCount()
	}

	log.Printf("[info] <GroupChatNodePool.CleanByActiveCount> clear up the group chat node, count= %d",
		int(cleanCount))

}
