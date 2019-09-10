package ApiRPC

import "testing"

var (
	testUserId1        int64 = 100
	testUserId2        int64 = 200
	testDelayedMessage       = []byte("this is a test message")
)

func TestSaveDelayedMessage(t *testing.T) {
	data, err := SaveDelayedMessage(testUserId1, testDelayedMessage)
	if nil != err {
		t.Error("SaveDelayedMessage fail: ", err)
	} else {
		t.Logf("SaveDelayedMessage success: save for user = %d", data.UserId)
	}
}
//
//func TestGetDelayedMessage(t *testing.T) {
//	data, err := GetDelayedMessage(testUserId1)
//	if nil != err {
//		t.Error("GetDelayedMessage fail: ", err)
//	} else {
//		t.Logf("GetDelayedMessage success: the message is = %s", data.MessageList)
//	}
//}
//
//func TestAddOneFriendId(t *testing.T) {
//	data, err := AddOneFriendId(testUserId1, testUserId2)
//	if nil != err {
//		t.Error("AddOneFriendId fail: ", err)
//	} else {
//		t.Logf("AddOneFriendId success: the selfId = %d, friendId = %d", data.SelfId, data.FriendId)
//	}
//}
//
//func TestGetAllFriendId(t *testing.T) {
//	data, err := GetAllFriendId(testUserId1)
//	if nil != err {
//		t.Error("GetAllFriendId fail: ", err)
//	} else {
//		t.Logf("GetAllFriendId success: data = %v", data)
//	}
//}
//
//func TestDelOneFriendId(t *testing.T) {
//	data, err := DelOneFriendId(testUserId1, testUserId2)
//	if nil != err {
//		t.Error("DelOneFriendId fail: ", err)
//	} else {
//		t.Logf("DelOneFriendId success: the selfid = %d, freindId = %d", data.SelfId, data.FriendId)
//	}
//}
//
//func TestAddOneFriendToBlacklist(t *testing.T) {
//	data, err := AddOneFriendToBlacklist(testUserId1, testUserId2)
//	if nil != err {
//		t.Error("AddOneFriendToBlacklist fail: ", err)
//	} else {
//		t.Logf("AddOneFriendToBlacklist success: selfId = %d, friendId = %d", data.SelfId, data.FriendId)
//	}
//}
//
//func TestGetBlacklistOfUser(t *testing.T) {
//	data, err := GetBlacklistOfUser(testUserId1)
//	if nil != err {
//		t.Error("GetBlacklistOfUser fail: ", err)
//	} else {
//		t.Logf("GetBlacklistOfUser success: the data = %v", data)
//	}
//}
//
//func TestDelOneFriendFromBlacklist(t *testing.T) {
//	data, err := DelOneFriendFromBlacklist(testUserId1, testUserId2)
//	if nil != err {
//		t.Error("DelOneFriendFromBlacklist fail:", err)
//	} else {
//		t.Logf("DelOneFriendFromBlacklist success: the data = %v", data)
//	}
//
//}

var (
	testGroupId int64 = 111
	testSubsId  int64 = 222
)

//func TestAddUserToGroupChat(t *testing.T) {
//	data, err := AddUserToGroupChat(testGroupId, testUserId1)
//	if nil != err {
//		t.Error("AddUserToGroupChat fail: ", err)
//	} else {
//		t.Logf("AddUserToGroupChat success: the data = %v", data)
//	}
//}
//
//func TestGetUsersOfGroupChat(t *testing.T) {
//	data, err := GetUsersOfGroupChat(testGroupId)
//	if nil != err {
//		t.Error("GetUsersOfCroupChat fail: ", err)
//	} else {
//		t.Logf("GetUsersOfGroupChat success: the data = %v", data)
//	}
//}
//
//func TestDelUserFromGroupChat(t *testing.T) {
//	data, err := DelUserFromGroupChat(testGroupId, testUserId1)
//	if nil != err {
//		t.Error("DelUserFromGroupChat fail: ", err)
//	} else {
//		t.Logf("DelUserFromGourpChat success: the data = %v", data)
//	}
//}

//func TestAddUserToSubscription(t *testing.T) {
//	data, err := AddUserToSubscription(testSubsId, testUserId1)
//	if nil != err {
//		t.Error("AddUserToSubscription fail: ", err)
//	} else {
//		t.Logf("AddUserToSubscription success: the data = %v", data)
//	}
//}
//
//func TestGetUsersOfSubscription(t *testing.T) {
//	data, err := GetUsersOfSubscription(testSubsId)
//	if nil != err {
//		t.Error("GetUsersOfSubscription fail: ", err)
//	} else {
//		t.Logf("GetUsersOfSubscription success: the data = %v", data)
//	}
//}
//
//func TestDelUserFromSubscription(t *testing.T) {
//	data, err := DelUserFromSubscription(testSubsId, testUserId1)
//	if nil != err {
//		t.Error("DelUserFromSubscription fail: ", err)
//	} else {
//		t.Logf("DelUserFromSubscription success: the data = %v", data)
//	}
//}
