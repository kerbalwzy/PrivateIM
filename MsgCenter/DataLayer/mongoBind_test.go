package DataLayer

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
