package MongoBind

import "testing"

var (
	testUserId1      int64 = 1234
	testUserId2      int64 = 2345
	testGroupChatId  int64 = 3456
	testSubsId       int64 = 4567
	testTodayDateNum       = getTodayNum()
)
// ------------------------------------------------------------------------------------

func TestUpdateDelayMessage(t *testing.T) {
	testDelayMessage := []byte("<test delay message>")
	err := UpdateDelayMessage(testUserId1, testDelayMessage)
	if nil != err {
		t.Fatal(err)
	}
}

func TestFindAndDeleteDelayMessage(t *testing.T) {
	data, err := FindAndDeleteDelayMessage(testUserId1)
	if nil != err {
		t.Fatal(err)
	}
	if len(data.Message) < 1 {
		t.Fatal("should have 1 delay message at least, but not")
	}
	for index, item := range data.Message {
		t.Logf("delayMessage: %d >> %s", index, item)
	}
}

// ------------------------------------------------------------------------------------

func TestUpdateUserChatHistoryByJoinId(t *testing.T) {
	tempJoinUserID := GetJoinUserId(testUserId1, testUserId2)
	err := UpdateUserChatHistoryByJoinId(tempJoinUserID, []byte("<test message for user chat history>"))
	if nil != err {
		t.Fatal(err)
	}
}

func TestFindUserAllChatHistoryByJoinId(t *testing.T) {
	tempJoinUserID := GetJoinUserId(testUserId1, testUserId2)
	data, err := FindUserAllChatHistoryByJoinId(tempJoinUserID)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("UserChatHistory: %s", data.Id)
	for _, item := range data.History {
		t.Logf("\tMessageDate: %d ", item.Date)
		for index, msg := range item.Messages {
			t.Logf("\t\tMessage: %d >> %s", index, msg)
		}
	}
}

func TestFindUserChatHistoryByJoinIdAndDateRange(t *testing.T) {
	data, err := FindUserChatHistoryByJoinIdAndDateRange(
		GetJoinUserId(testUserId1, testUserId2),
		testTodayDateNum-1, testTodayDateNum+1,
	)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("UserChatHistory: %s", data.Id)
	for _, item := range data.History {
		t.Logf("\tMessageDate: %d ", item.Date)
		for index, msg := range item.Messages {
			t.Logf("\t\tMessage: %d >> %s", index, msg)
		}
	}
	// test with not existed id
	_, err = FindUserChatHistoryByJoinIdAndDateRange("wrongId", testTodayDateNum-1, testTodayDateNum+1)
	if nil == err {
		t.Fatal("should have en error")
	}
	t.Logf("WantError: %s", err.Error())

	// test with wrong date
	_, err = FindUserChatHistoryByJoinIdAndDateRange(GetJoinUserId(testUserId1, testUserId2),
		testTodayDateNum+2, testTodayDateNum+3)
}

func TestFindUserChatHistoryByJoinIdAndDate(t *testing.T) {
	data, err := FindUserChatHistoryByJoinIdAndDate(GetJoinUserId(testUserId1, testUserId2), testTodayDateNum)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("UserChatHistory: %s", data.Id)
	for _, item := range data.History {
		t.Logf("\tMessageDate: %d ", item.Date)
		for index, msg := range item.Messages {
			t.Logf("\t\tMessage: %d >> %s", index, msg)
		}
	}

	// test with the date have not record
	_, err = FindUserChatHistoryByJoinIdAndDate(GetJoinUserId(testUserId1, testUserId2), testTodayDateNum+2)
	if nil == err {
		t.Fatal("should have en error")
	}
	t.Logf("WantError: %s", err.Error())
}

// ------------------------------------------------------------------------------------

func TestUpdateGroupChatHistoryById(t *testing.T) {
	tempMessage := []byte("<test group chat message>")
	err := UpdateGroupChatHistoryById(testGroupChatId, tempMessage)
	if nil != err {
		t.Fatal(err)
	}
}

func TestFindAllGroupChatHistoryById(t *testing.T) {
	data, err := FindAllGroupChatHistoryById(testGroupChatId)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("GroupChatHistory: %d", data.GroupId)
	for _, item := range data.History {
		t.Logf("HistoryDate: %d ", item.Date)
		for index, msg := range item.Messages {
			t.Logf("\t\tHistoryMsg: %d >> %s", index, msg)
		}
	}
}

func TestFindGroupChatHistoryByIdAndDateRange(t *testing.T) {
	data, err := FindGroupChatHistoryByIdAndDateRange(testGroupChatId, testTodayDateNum-1, testTodayDateNum+1)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("GroupChatHistory: %d", data.GroupId)
	for _, item := range data.History {
		t.Logf("HistoryDate: %d ", item.Date)
		for index, msg := range item.Messages {
			t.Logf("\t\tHistoryMsg: %d >> %s", index, msg)
		}
	}

	// test with not existed group id
	_, err = FindGroupChatHistoryByIdAndDateRange(0, testTodayDateNum-1, testTodayDateNum+1)
	if err == nil {
		t.Fatal("should have an error, but not")
	}

	// test with the date have not history message record
	_, err = FindGroupChatHistoryByIdAndDateRange(testGroupChatId, testTodayDateNum+2, testTodayDateNum+3)
	if err == nil {
		t.Fatal("should have an error, but not")
	}
}

func TestFindGroupChatHistoryByIdAndDate(t *testing.T) {
	data, err := FindGroupChatHistoryByIdAndDate(testGroupChatId, testTodayDateNum)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("GroupChatHistory: %d", data.GroupId)
	for _, item := range data.History {
		t.Logf("HistoryDate: %d ", item.Date)
		for index, msg := range item.Messages {
			t.Logf("\t\tHistoryMsg: %d >> %s", index, msg)
		}
	}
	// test with not existed group id
	_, err = FindGroupChatHistoryByIdAndDate(0, testTodayDateNum)
	if err == nil {
		t.Fatal("should have an error, but not")
	}

	// test with the date have not history message record
	_, err = FindGroupChatHistoryByIdAndDate(testGroupChatId, testTodayDateNum+2)
}

// ------------------------------------------------------------------------------------

func TestUpdateSubscriptionHistoryById(t *testing.T) {
	tempMessage := []byte("<test subscription message>")
	err := UpdateSubscriptionHistoryById(testSubsId, tempMessage)
	if nil != err {
		t.Fatal(err)
	}
}

func TestFindAllSubscriptionHistoryById(t *testing.T) {
	data, err := FindAllSubscriptionHistoryById(testSubsId)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("SubscriptionHistory: %d", data.SubsId)
	for _, item := range data.History {
		t.Logf("HistoryDate: %d ", item.Date)
		for index, msg := range item.Messages {
			t.Logf("\t\tHistoryMsg: %d >> %s", index, msg)
		}
	}
}

func TestFindSubscriptionHistoryByIdAndDateRange(t *testing.T) {
	data, err := FindSubscriptionHistoryByIdAndDateRange(testSubsId, testTodayDateNum-1, testTodayDateNum+1)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("SubscriptionHistory: %d", data.SubsId)
	for _, item := range data.History {
		t.Logf("HistoryDate: %d ", item.Date)
		for index, msg := range item.Messages {
			t.Logf("\t\tHistoryMsg: %d >> %s", index, msg)
		}
	}

	// test with not existed subscription id
	_, err = FindSubscriptionHistoryByIdAndDateRange(0, testTodayDateNum-1, testTodayDateNum+1)
	if err == nil {
		t.Fatal("should have an error, but not")
	}

	// test with the date have not history message record
	_, err = FindSubscriptionHistoryByIdAndDateRange(testSubsId, testTodayDateNum+2, testTodayDateNum+3)
	if err == nil {
		t.Fatal("should have an error, but not")
	}
}

func TestFindSubscriptionHistoryByIdAndDate(t *testing.T) {
	data, err := FindSubscriptionHistoryByIdAndDate(testSubsId, testTodayDateNum)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("SubscriptionHistory: %d", data.SubsId)
	for _, item := range data.History {
		t.Logf("HistoryDate: %d ", item.Date)
		for index, msg := range item.Messages {
			t.Logf("\t\tHistoryMsg: %d >> %s", index, msg)
		}
	}

	// test with not existed subscription id
	_, err = FindSubscriptionHistoryByIdAndDate(0, testTodayDateNum)
	if err == nil {
		t.Fatal("should have an error, but not")
	}

	// test with the date have not history message record
	_, err = FindSubscriptionHistoryByIdAndDate(testSubsId, testTodayDateNum+2)
	if err == nil {
		t.Fatal("should have an error, but not")
	}
}

// ------------------------------------------------------------------------------------

var (
	testUserId   int64 = 12345
	testFriendId int64 = 23456
)

func TestUpdateUserFriendsToAddFriend(t *testing.T) {
	err := UpdateUserFriendsToAddFriend(testUserId, testFriendId)
	if nil != err {
		t.Fatal(err)
	}

	_ = UpdateUserFriendsToAddFriend(testUserId2, testFriendId)
}

func TestUpdateUserFriendsToDelFriend(t *testing.T) {
	err := UpdateUserFriendsToDelFriend(testUserId, testFriendId)
	if nil != err {
		t.Fatal(err)
	}
}

func TestUpdateUserBlacklistToAddUser(t *testing.T) {
	err := UpdateUserBlacklistToAddUser(testUserId, testFriendId)
	if nil != err {
		t.Fatal(err)
	}

	_ = UpdateUserBlacklistToAddUser(testUserId2, testUserId1)
}

func TestUpdateUserBlacklistToDelUser(t *testing.T) {
	err := UpdateUserBlacklistToDelUser(testUserId, testFriendId)
	if nil != err {
		t.Fatal(err)
	}
}

func TestFindUserFriendsAndBlacklistById(t *testing.T) {
	data, err := FindUserFriendsAndBlacklistById(testUserId2)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("UserFriendsAndBlacklist: for user id= %d", data.UserId)
	for index, item := range data.Friends {
		t.Logf("friends: %d >> %d", index, item)
	}

	for index, item := range data.Blacklist {
		t.Logf("blacklist: %d >> %d", index, item)
	}
}

// Plus functions test
func TestUpdateMoveFriendIntoBlacklistPlus(t *testing.T) {
	err := UpdateMoveFriendIntoBlacklistPlus(testUserId, testFriendId)
	if nil != err {
		t.Fatal(err)
	}
}

func TestUpdateMoveFriendOutFromBlacklistPlus(t *testing.T) {
	err := UpdateMoveFriendOutFromBlacklistPlus(testUserId, testFriendId)
	if nil != err {
		t.Fatal(err)
	}
}

// ------------------------------------------------------------------------------------

func TestUpdateUserGroupChatsToAddOne(t *testing.T) {
	err := UpdateUserGroupChatsToAddOne(testUserId1, testGroupChatId)
	if nil != err {
		t.Fatal(err)
	}
}

func TestFindUserGroupChatsById(t *testing.T) {
	data, err := FindUserGroupChatsById(testUserId1)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("user(%d)'s group chat id", data.UserId)
	for index, id := range data.Groups {
		t.Logf("%d >> %d", index, id)
	}
}

func TestUpdateUserGroupChatsToDelOne(t *testing.T) {
	err := UpdateUserGroupChatsToDelOne(testUserId1, testGroupChatId)
	if nil != err {
		t.Fatal(err)
	}
}

// ------------------------------------------------------------------------------------

func TestUpdateUserSubscriptionsToAddOne(t *testing.T) {
	err := UpdateUserSubscriptionsToAddOne(testUserId, testSubsId)
	if nil != err {
	}
}

func TestFindUserSubscriptionsById(t *testing.T) {
	data, err := FindUserSubscriptionsById(testUserId)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("user(%d)'s subscriptions id", data.UserId)
	for index, id := range data.Subscriptions {
		t.Logf("%d >> %d", index, id)
	}
}

func TestUpdateUserSubscriptionsToDelOne(t *testing.T) {
	err := UpdateUserSubscriptionsToDelOne(testUserId, testSubsId)
	if nil != err {
		t.Fatal(err)
	}
}

// ------------------------------------------------------------------------------------

func TestUpdateGroupChatUserToAddOne(t *testing.T) {
	err := UpdateGroupChatUserToAddOne(testGroupChatId, testUserId2, testUserId)
	if nil != err {
		t.Fatal(err)
	}
}

func TestFindGroupChatUsersById(t *testing.T) {
	data, err := FindGroupChatUsersById(testGroupChatId)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("group chat(%d)'s manager id= %d", data.GroupId, data.ManagerId)
	t.Logf("group chat(%d)'s user id", data.GroupId)
	for index, id := range data.Users {
		t.Logf("%d >> %d", index, id)
	}
}

func TestUpdateGroupChatUsersToDelOne(t *testing.T) {
	err := UpdateGroupChatUsersToDelOne(testGroupChatId, testUserId)
	if nil != err {
		t.Fatal(err)
	}
}

// ------------------------------------------------------------------------------------

func TestUpdateSubscriptionUsersToAddOne(t *testing.T) {
	err := UpdateSubscriptionUsersToAddOne(testSubsId, testUserId2, testUserId)
	if nil != err {
		t.Fatal(err)
	}
}

func TestFindSubscriptionUsersById(t *testing.T) {
	data, err := FindSubscriptionUsersById(testSubsId)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("subscription(%d)'s manager id= %d", data.SubsId, data.ManagerId)
	t.Logf("subscription(%d)'s user id", data.SubsId)
	for index, id := range data.Users {
		t.Logf("%d >> %d", index, id)
	}
}

func TestUpdateSubscriptionUsersToDelOne(t *testing.T) {
	err := UpdateSubscriptionUsersToDelOne(testSubsId, testUserId)
	if nil != err {
		t.Fatal(err)
	}
}

// ------------------------------------------------------------------------------------

// Plus functions test
func TestUpdateMoveUserIntoGroupChat(t *testing.T) {
	err := UpdateMoveUserIntoGroupChat(testUserId1, testGroupChatId)
	if nil != err {
		t.Fatal(err)
	}
}

func TestUpdateMoveUserOutFromGroupChat(t *testing.T) {
	err := UpdateMoveUserOutFromGroupChat(testUserId1, testGroupChatId)
	if nil != err {
		t.Fatal(err)
	}
}

func TestUpdateMakeUserFollowSubscription(t *testing.T) {
	err := UpdateMakeUserFollowSubscription(testUserId1, testSubsId)
	if nil != err {
		t.Fatal(err)
	}
}

func TestUpdateMakeUserUnFollowSubscription(t *testing.T) {
	err := UpdateMakeUserUnFollowSubscription(testUserId1, testSubsId)
	if nil != err {
		t.Fatal(err)
	}
}
