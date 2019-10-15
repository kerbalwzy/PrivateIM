package elasticClient

import (
	"testing"
)

var (
	testId     int64 = 1183956383159029784
	testName         = "wang王@123"
	testEmail        = "test@demo.com"
	testAvatar       = "<test avatar url>"
	testGender       = 1
)

func TestUserIndexDocSave(t *testing.T) {

	err := UserIndexDocSave(testId, testName, testEmail, testAvatar, testGender)
	if nil != err {
		t.Fatal(err)
	}
}

func TestUserIndexDocSearch(t *testing.T) {
	data, err := UserIndexDocSearch("wang", 1, 10)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("%s", data)
}

func TestUserIndexDocUpdate(t *testing.T) {
	err := UserIndexDocUpdate(testId, "gender", 2)
	if nil != err {
		t.Fatal(err)
	}
}

var (
	testGId   int64 = 1183956383159029799
	testGName       = "studyES"
)

func TestGroupChatIndexDocSave(t *testing.T) {
	err := GroupChatIndexDocSave(testGId, testGName, testAvatar, testName, testAvatar)
	if nil != err {
		t.Fatal(err)
	}
}

func TestGroupChatIndexDocSearch(t *testing.T) {
	data, err := GroupChatIndexDocSearch("wang", 2, 1)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("%s", data)
}

func TestGroupChatIndexDocUpdate(t *testing.T) {
	err := GroupChatIndexDocUpdate(testGId, "manager_name", "wang")
	if nil != err {
		t.Fatal(err)
	}
}

var (
	testSId    int64 = 1183956383159029777
	testSName        = "hello word"
	testSIntro       = "<this is hello word subscription for you>"
)

func TestSubscriptionIndexSave(t *testing.T) {
	err := SubscriptionIndexSave(testSId, testSName, testSIntro, testAvatar, testName, testAvatar)
	if nil != err {
		t.Fatal(err)
	}
}

func TestSubscriptionIndexSearch(t *testing.T) {
	data, err := SubscriptionIndexSearch("word", 1, 1)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("%s", data)
}

func TestSubscriptionIndexUpdate(t *testing.T) {
	err := SubscriptionIndexUpdate(testSId, "avatar", "<new avatar url>")
	if nil != err {
		t.Fatal(err)
	}
}

func TestIndexDocDelete(t *testing.T) {
	var err error
	err = IndexDocDelete(UserIndexName, testId)
	if nil != err {
		t.Fatal(err)
	}
	err = IndexDocDelete(GroupChatIndexName, testGId)
	if nil != err {
		t.Fatal(err)
	}
	err = IndexDocDelete(SubscriptionIndexName, testSId)
	if nil != err {
		t.Fatal(err)
	}
}
