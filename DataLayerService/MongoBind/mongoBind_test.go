package MongoBind

import (
	"log"
	"testing"
)

var (
	srcId   int64 = 1111111
	dstId   int64 = 2222222
	testMsg       = []byte("this is a test message")
)

func TestMongoSaveWaitSendMessage(t *testing.T) {
	err := MongoSaveWaitSendMessage(srcId, testMsg)
	if nil != err {
		t.Error(err)
	}
}

func TestMongoQueryWaitSendMessage(t *testing.T) {
	messages, err := MongoQueryWaitSendMessage(srcId)
	if nil != err {
		t.Error(err)
	}
	for _, msg := range messages {
		log.Printf("%s\n", msg)
	}
}

func TestMongoAddFriendId(t *testing.T) {
	err := MongoAddFriendId(srcId, dstId)
	if nil != err {
		t.Error(err)
	}
}

func TestMongoQueryFriendsId(t *testing.T) {
	friends, err := MongoQueryFriendsId(srcId)
	if nil != err {
		t.Error(err)
	}
	for _, id := range friends {
		log.Printf("friend id :%d\n", id)
	}
}

func TestMongoDelFriendId(t *testing.T) {
	err := MongoDelFriendId(srcId, dstId)
	if nil != err {
		t.Error(err)
	}
}

func TestMongoBlackListAdd(t *testing.T) {
	err := MongoBlackListAdd(srcId, dstId)
	if nil != err {
		t.Error(err)
	}
}

func TestMongoQueryBlackList(t *testing.T) {
	blackList, err := MongoQueryBlackList(srcId)
	if nil != err {
		t.Error(err)
	}
	for _, id := range blackList {
		log.Printf("%d\n", id)
	}
}

func TestMongoBlackListDel(t *testing.T) {
	err := MongoBlackListDel(srcId, dstId)
	if nil != err {
		t.Error(err)
	}
}

var (
	groupId1 int64 = 12345678
	groupId2 int64 = 98765432
)

func TestMongoGroupChatAddUser(t *testing.T) {
	err := MongoGroupChatAddUser(groupId1, srcId)
	if nil != err {
		t.Error(err)
	}
	err = MongoGroupChatAddUser(groupId1, dstId)
	if nil != err {
		t.Error(err)
	}
	err = MongoGroupChatAddUser(groupId2, srcId)
	if nil != err {
		t.Error(err)
	}
	err = MongoGroupChatAddUser(groupId2, dstId)
	if nil != err {
		t.Error(err)
	}
}

func TestMongoQueryGroupChatUsers(t *testing.T) {
	users, err := MongoQueryGroupChatUsers(groupId1)
	if nil != err {
		t.Error(err)
	}
	log.Printf("user count= %d", len(users))
}

func TestMongoQueryGroupChatAll(t *testing.T) {
	data, err := MongoQueryGroupChatAll()
	if nil != err {
		t.Error(err)
	}
	log.Printf("data count= %d", len(data))
	for _, detail := range data {
		log.Printf("group(%d) has users: %v", detail.Id, detail.UsersId)
	}
}

func TestMongoGroupChatDelUser(t *testing.T) {
	err := MongoGroupChatDelUser(groupId1, srcId)
	if nil != err {
		t.Error(err)
	}
}

func TestMongoQueryGroupChatUsers2(t *testing.T) {
	TestMongoQueryGroupChatUsers(t)
}

var (
	subsId1 int64 = 56789098765
	subsId2 int64 = 5678987656789
)

func TestMongoSubscriptionAddUser(t *testing.T) {
	err := MongoSubscriptionAddUser(subsId1, srcId)
	if nil != err {
		t.Error(err)
	}
	err = MongoSubscriptionAddUser(subsId1, dstId)
	if nil != err {
		t.Error(err)
	}
	err = MongoSubscriptionAddUser(subsId2, srcId)
	if nil != err {
		t.Error(err)
	}
	err = MongoSubscriptionAddUser(subsId2, dstId)
	if nil != err {
		t.Error(err)
	}

}

func TestMongoQuerySubscriptionUsers(t *testing.T) {
	users, err := MongoQuerySubscriptionUsers(subsId1)
	if nil != err {
		t.Error(err)
	}
	log.Printf("user conut of subscription(%d) is %d", subsId1, len(users))
}

func TestMongoQuerySubscriptionAll(t *testing.T) {
	data, err := MongoQuerySubscriptionAll()
	if nil != err {
		t.Error(err)
	}
	for _, detail := range data {
		log.Printf("subscription(%d) has users(%v)", detail.Id, detail.UsersId)
	}
}

func TestMongoSubscriptionDelUser(t *testing.T) {
	err := MongoSubscriptionDelUser(subsId1, srcId)
	if nil != err {
		t.Error(err)
	}
}

func TestMongoQuerySubscriptionAll2(t *testing.T) {
	TestMongoQuerySubscriptionAll(t)
}
