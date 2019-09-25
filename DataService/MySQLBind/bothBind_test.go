package MySQLBind

import (
	"testing"
)

// Testing for operating 'tb_user_basic'
var (
	userId1, userId2 int64
	tempEmail1             = "test_1_email@test.com"
	tempEmail2             = "test_2_email@test.com"
	tempName               = "testName"
	tempPassword           = "<test password 1, should be hash value>"
	tempMobile             = "13100000000"
	tempGender       int32 = 1
	tempAvatar             = "<temp avatar pic name>"
	tempQrCode1            = "<temp qr_code pic name 1>"
	tempQrCode2            = "<temp qr_code pic name 2>"
)

func TestInsertOneNewUser(t *testing.T) {
	user, err := InsertOneNewUser(tempEmail1, tempName, tempPassword, tempMobile,
		tempGender, tempAvatar, tempQrCode1, false)
	if nil != err {
		t.Fatal(err)
	}
	if user.Id == 0 {
		t.Fatal("last insert id = 0")
	}
	userId1 = user.Id

	user2, err := InsertOneNewUser(tempEmail2, tempName, tempPassword, tempMobile,
		tempGender, tempAvatar, tempQrCode2, false)
	if nil != err {
		t.Fatal(err)
	}
	if user2.Id == 0 {
		t.Fatal("last insert id = 0")
	}
	userId2 = user2.Id
}

func TestSelectOneUserById(t *testing.T) {
	user, err := SelectOneUserById(userId1, false)
	if nil != err {
		t.Fatal(err)
	}
	if user.Email != tempEmail1 {
		t.Fatal("query an wrong user with wrong email:", user.Email)
	}

	user2, err := SelectOneUserById(userId2, false)
	if nil != err {
		t.Fatal(err)
	}
	if user2.Email != tempEmail2 {
		t.Fatal("query an wrong user with wrong email:", user2.Email)
	}

}

func TestSelectOneUserByEmail(t *testing.T) {
	user, err := SelectOneUserByEmail(tempEmail1, false)
	if nil != err {
		t.Fatal(err)
	}
	if user.Id != userId1 {
		t.Fatal("query a wrong user with wrong id:", user.Id)
	}

	user2, err := SelectOneUserByEmail(tempEmail2, false)
	if nil != err {
		t.Fatal(err)
	}
	if user2.Id != userId2 {
		t.Fatal("query a wrong user with wrong id:", user2.Id)
	}
}

func TestSelectManyUserByName(t *testing.T) {
	users, err := SelectManyUserByName(tempName, false)
	if nil != err {
		t.Fatal(err)
	}
	isErr := 0
	for _, user := range users {
		t.Logf("user: %v\n", user)
		if user.Id == userId1 {
			isErr += 1
		}
		if user.Id == userId2 {
			isErr += 1
		}
	}

	if isErr < 2 {
		t.Fatal("query users by name error, should get 2 but only get:", isErr)
	}

}

func TestSelectOneUserPasswordById(t *testing.T) {
	password, err := SelectOneUserPasswordById(userId1)
	if nil != err {
		t.Fatal(err)
	}
	if password != tempPassword {
		t.Fatal("query password by id fail, should be:", tempPassword, "but get:", password)
	}
}

func TestSelectOneUserPasswordByEmail(t *testing.T) {
	password, err := SelectOneUserPasswordByEmail(tempEmail1)
	if nil != err {
		t.Fatal(err)
	}
	if password != tempPassword {
		t.Fatal("query password by id fail, should be:", tempPassword, "but get:", password)
	}
}

func TestSelectAllUsers(t *testing.T) {
	users, err := SelectAllUsers()
	if nil != err {
		t.Fatal(err)
	}
	if len(users) < 2 {
		t.Fatalf("here should have 2 user at least, but not.")
	}
	for index, user := range users {
		t.Logf("SelectAllUsers: %d >> %v", index, user)
	}
}

func TestUpdateOneUserAvatarById(t *testing.T) {
	newAvatar := "<new avatar pic name>"
	err := UpdateOneUserAvatarById(newAvatar, userId1)
	if nil != err {
		t.Fatal(err)
	}
	user, err := SelectOneUserById(userId1, false)
	if nil != err {
		t.Fatal(err)
	}
	if user.Avatar != newAvatar {
		t.Fatal("update avatar by id fail, new avatar should be :", newAvatar, "but get:", user.Avatar)
	}
}

func TestUpdateOneUserQrCodeById(t *testing.T) {
	newQrcode := "<new qr_code pic name>"
	err := UpdateOneUserQrCodeById(newQrcode, userId1)
	if nil != err {
		t.Fatal(err)
	}
	user, err := SelectOneUserById(userId1, false)
	if nil != err {
		t.Fatal(err)
	}
	if user.QrCode != newQrcode {
		t.Fatal("update qr_code by id fail, should be:", newQrcode, "but get:", user.QrCode)
	}
}

func TestUpdateOneUserPasswordById(t *testing.T) {
	newPassword := "<new password , should be hash value>"
	err := UpdateOneUserPasswordById(newPassword, userId1)
	if nil != err {
		t.Fatal(err)
	}
	password, err := SelectOneUserPasswordById(userId1)
	if nil != err {
		t.Fatal(err)
	}
	if password != newPassword {
		t.Fatal("update password by id fail, new password should be:", newPassword, "but get:", password)
	}
}

func TestUpdateOneUserIsDeleteById(t *testing.T) {
	err := UpdateOneUserIsDeleteById(true, userId1)
	if nil != err {
		t.Fatal(err)
	}
	user, err := SelectOneUserById(userId1, true)
	if nil != err {
		t.Fatal(err)
	}
	if user.IsDelete != true {
		t.Fatal("update isDelete by id fail, should be true, but get false")
	}
	_ = UpdateOneUserIsDeleteById(false, userId1)

}

// Plus function test
func TestUpdateOneUserProfileByIdPlus(t *testing.T) {
	newName := "NewName"
	newMobile := "13199999999"
	var newGender int32 = 2
	err := UpdateOneUserProfileByIdPlus(newName, newMobile, newGender, userId1)
	if nil != err {
		t.Fatal(err)
	}
	user, err := SelectOneUserById(userId1, false)
	if nil != err {
		t.Fatal(err)
	}
	if user.Name != newName || user.Mobile != newMobile || user.Gender != newGender {
		t.Fatal("update user profile fail, some data changed wrong")
	}
}

// ------------------------------------------------------------------

// Testing for operating 'tb_friendship'
var (
	selfId, friendId         int64
	friendNote1, friendNote2 = "note1", "note2"
)

func TestInsertOneNewFriend(t *testing.T) {
	selfId = userId1
	friendId = userId2
	err := InsertOneNewFriend(selfId, friendId, friendNote1)
	if nil != err {
		t.Fatal(err)
	}

	data, err := SelectOneFriendship(selfId, friendId)
	if nil != err {
		t.Fatal(err)
	}
	if data.FriendNote != friendNote1 || data.SelfId != selfId || data.FriendId != friendId {
		t.Fatalf("insert one new friend fail, the data was wrong, get data:\n\t%v", data)
	}
}

func TestUpdateOneFriendNote(t *testing.T) {
	newNote := "newNote"
	err := UpdateOneFriendNote(selfId, friendId, newNote)
	if nil != err && ErrAffectZeroCount != err {
		t.Fatal(err)
	}
	data, _ := SelectOneFriendship(selfId, friendId)
	if data.FriendNote != newNote {
		t.Fatal("the new note should be:", newNote, " but get:", data.FriendNote)
	}
}

func TestUpdateOneFriendIsAccept(t *testing.T) {
	err := UpdateOneFriendIsAccept(selfId, friendId, true)
	if nil != err {
		t.Fatal(err)
	}
	data, _ := SelectOneFriendship(selfId, friendId)
	if data.IsAccept != true {
		t.Fatal("the is_accept should be true, but is false")
	}
	_ = UpdateOneFriendIsAccept(selfId, friendId, false)
}

func TestUpdateOneFriendIsBlack(t *testing.T) {
	err := UpdateOneFriendIsBlack(selfId, friendId, true)
	if nil != err {
		t.Fatal(err)
	}
	data, _ := SelectOneFriendship(selfId, friendId)
	if data.IsBlack != true {
		t.Fatal("the is_black should be true, but is false")
	}

	_ = UpdateOneFriendIsBlack(selfId, friendId, false)
}

func TestUpdateOneFriendIsDelete(t *testing.T) {
	err := UpdateOneFriendIsDelete(selfId, friendId, true)
	if nil != err {
		t.Fatal(err)
	}
	data, _ := SelectOneFriendship(selfId, friendId)
	if data.IsDelete != true {
		t.Fatalf("the is_delete should be ture, but is false")
	}
	_ = UpdateOneFriendIsDelete(selfId, friendId, false)
}

func TestSelectOneFriendship(t *testing.T) {
	data, err := SelectOneFriendship(selfId, friendId)
	if nil != err {
		t.Fatal(err)
	}
	if data == nil {
		t.Fatal("the data should not be nil")
	}
	t.Logf("the friendship: %v", data)
}

func TestSelectFriendsIdByOptions(t *testing.T) {
	ids, err := SelectFriendsIdByOptions(selfId, false, false, false)
	if nil != err {
		t.Fatal(err)
	}

	isErr := true
	for _, id := range ids {
		if id == friendId {
			isErr = false
		}
	}
	if isErr {
		t.Fatal("the friendId should in ids, but not")
	}
	t.Logf("ids: %v", ids)
}

func TestSelectAllFriendship(t *testing.T) {
	data, err := SelectAllFriendship()
	if nil != err {
		t.Fatal(err)
	}
	if len(data) < 1 {
		t.Fatal("there should have 1 friendship record at least, but not")
	}
	for index, temp := range data {
		t.Logf("AllFrienship: %d >> %v", index, temp)
	}
}

// Plus function test
func TestInsertOneNewFriendPlus(t *testing.T) {
	// try to add a not existed user as friend
	err := InsertOneNewFriendPlus(selfId, 0, "")
	if nil == err {
		t.Fatal("should have an error, but not")
	}
	t.Logf("WantError: %s", err.Error())

	// try add the friend who add me into his blacklist
	_ = UpdateOneFriendIsBlack(friendId, selfId, true)
	err = InsertOneNewFriendPlus(selfId, friendId, "")
	if nil == err {
		t.Fatal("should have an error, but not")
	}
	t.Logf("WantError: %s", err.Error())

	// try add the friend who already have the effect friendship between the two user.
	_ = UpdateOneFriendIsBlack(friendId, selfId, false)
	_ = UpdateOneFriendIsAccept(friendId, selfId, true)
	err = InsertOneNewFriendPlus(selfId, friendId, "")
	if nil == err {
		t.Fatal("should have an error, but not")
	}
	t.Logf("WantError: %s", err.Error())

	// normal insert one new friend
	_ = DeleteOneFriendReal(friendId, selfId)
	err = InsertOneNewFriendPlus(selfId, friendId, friendNote1)
	if nil != err {
		t.Fatal(err)
	}
	data, _ := SelectOneFriendship(selfId, friendId)
	if data == nil {
		t.Fatal("the data should not be nil")
	}
	t.Logf("friendship: %v", data)

}

func TestUpdateAcceptOneNewFriendPlus(t *testing.T) {
	// try to accept one not existed user as friend
	tempSelfId := friendId
	tempFriendId := selfId
	err := UpdateAcceptOneNewFriendPlus(tempSelfId, 0, "", true)
	if nil == err {
		t.Fatal("should have an error, but not")
	}
	t.Logf("WantError: %s", err.Error())

	// try to accept one already effect friendship
	_ = UpdateOneFriendIsAccept(selfId, friendId, true)
	err = UpdateAcceptOneNewFriendPlus(tempSelfId, tempFriendId, "", true)
	if nil == err {
		t.Fatal("should have an error, but not")
	}
	t.Logf("WantError: %s", err.Error())
	_ = UpdateOneFriendIsAccept(selfId, friendId, false)

	// normal refuse one friend request
	err = UpdateAcceptOneNewFriendPlus(tempSelfId, tempFriendId, "", false)
	if nil != err {
		t.Fatal(err)
	}
	data, _ := SelectOneFriendship(selfId, friendId)
	if data.IsAccept {
		t.Fatal("the is_accept of request should be false")
	}
	data, _ = SelectOneFriendship(tempSelfId, tempFriendId)
	if data.IsAccept {
		t.Fatal("the is_accept of recipient should be false")
	}
	if data.IsBlack != true {
		t.Fatal("the is_black of recipient should be true")
	}
	_ = DeleteOneFriendReal(tempSelfId, tempFriendId)

	// normal accept one friend request
	err = UpdateAcceptOneNewFriendPlus(tempSelfId, tempFriendId, "", true)
	if nil != err {
		t.Fatal(err)
	}
	data, _ = SelectOneFriendship(selfId, friendId)
	if data.IsAccept != true {
		t.Fatal("the is_accept of friendship of requester should be true")
	}
	data, _ = SelectOneFriendship(tempSelfId, tempFriendId)
	if data.IsAccept != true {
		t.Fatal("the is_accept of friendship of recipient should be true")
	}

}

func TestUpdateDeleteOneFriendPlus(t *testing.T) {
	err := UpdateDeleteOneFriendPlus(selfId, friendId)
	if nil != err {
		t.Fatal(err)
	}
	data, _ := SelectOneFriendship(selfId, friendId)
	if data.IsDelete != true {
		t.Fatal("the is_delete should be true")
	}
	data, _ = SelectOneFriendship(friendId, selfId)
	if data.IsDelete != true {
		t.Fatal("the is_delete should be true")
	}
	_ = UpdateOneFriendIsDelete(selfId, friendId, false)
	_ = UpdateOneFriendIsDelete(friendId, selfId, false)
}

func TestSelectAllFriendsInfoPlus(t *testing.T) {
	friends, err := SelectAllFriendsInfoPlus(selfId)
	t.Logf("SelfId = %d", selfId)
	if nil != err {
		t.Fatal(err)
	}
	if len(friends) < 1 {
		t.Fatal("select 1 friend at least, but not")
	}
	for index, friend := range friends {
		t.Logf("all friend:%d >>> %v", index, friend)
	}
}

func TestSelectEffectiveFriendsInfoPlus(t *testing.T) {
	friends, err := SelectEffectiveFriendsInfoPlus(selfId)
	if nil != err {
		t.Fatal(err)
	}
	friendCount1 := len(friends)
	if friendCount1 < 1 {
		t.Fatal("get 1 friend information at least, but not")
	}
	for index, friend := range friends {
		t.Logf("effect friend: %d >> %v", index, friend)
	}

}

func TestSelectBlacklistFriendsInfoPlus(t *testing.T) {
	_ = UpdateOneFriendIsBlack(selfId, friendId, true)
	friends, err := SelectBlacklistFriendsInfoPlus(selfId)
	if nil != err {
		t.Fatal(err)
	}
	friendCount1 := len(friends)
	if friendCount1 < 1 {
		t.Fatal("should have 1 friend in blacklist at least, but not")
	}
	for index, friend := range friends {
		t.Logf("blacklist: %d >> %v", index, friend)
	}

	// move the friend out from the blacklist
	_ = UpdateOneFriendIsBlack(selfId, friendId, false)
}

func TestSelectEffectiveFriendsIdPlus(t *testing.T) {
	ids, err := SelectEffectiveFriendsIdPlus(selfId)
	if nil != err {
		t.Fatal(err)
	}
	if len(ids) < 1 {
		t.Fatal("the count of id should right than 1, but not")
	}
}

func TestSelectBlacklistFriendsId(t *testing.T) {
	_ = UpdateOneFriendIsBlack(selfId, friendId, true)
	ids, err := SelectBlacklistFriendsIdPlus(selfId)
	if nil != err {
		t.Fatal(err)
	}
	if len(ids) < 1 {
		t.Fatal("the count of id should right than 1, but not")
	}
	_ = UpdateOneFriendIsBlack(selfId, friendId, false)
}

// ------------------------------------------------------------------

// Testing for operating 'tb_group_chat'
var (
	groupChatName = "测试群聊1"
	managerId     int64
	groupAvatar   = "<the group chat avatar pic name>"
	groupQrCode   = "<the group qrCode pic name>"

	testGroupChatId int64
)

func TestInsertOneNewGroupChat(t *testing.T) {
	managerId = userId1
	groupChat, err := InsertOneNewGroupChat(groupChatName, groupAvatar, groupQrCode, managerId)
	if nil != err {
		t.Fatal(err)
	}
	testGroupChatId = groupChat.Id
}

func TestSelectOneGroupChatById(t *testing.T) {
	groupChat, err := SelectOneGroupChatById(testGroupChatId, false)
	if nil != err {
		t.Fatal(err)
	}
	if groupChat.ManagerId != managerId || groupChat.Name != groupChatName || groupChat.QrCode != groupQrCode ||
		groupChat.Avatar != groupAvatar {
		t.Fatal("select one group chat wrong, the get values not equal the set values")
	}
}

func TestSelectManyGroupChatByName(t *testing.T) {
	groupChats, err := SelectManyGroupChatByName(groupChatName, false)
	if nil != err {
		t.Fatal(err)
	}
	if len(groupChats) < 1 {
		t.Fatalf("there are 1 group chat named: %s at least, but not", groupChatName)
	}
	for index, groupChat := range groupChats {
		t.Logf("group chat: %d >> %v", index, groupChat)
	}

}

func TestSelectManyGroupChatByManagerId(t *testing.T) {
	groupChats, err := SelectManyGroupChatByManagerId(managerId, false)
	if nil != err {
		t.Fatal(err)
	}
	if len(groupChats) < 1 {
		t.Fatalf("the manger(%d) have 1 group chat at least, but not", managerId)
	}
	for index, groupChat := range groupChats {
		t.Logf("group chat: %d >> %v", index, groupChat)
	}
}

func TestSelectAllGroupChat(t *testing.T) {
	groupChats, err := SelectAllGroupChat()
	if nil != err {
		t.Fatal(err)
	}
	if len(groupChats) < 1 {
		t.Fatalf("there should have 1 group chat at least")
	}

}

func TestUpdateOneGroupChatNameById(t *testing.T) {
	newName := "newName"
	err := UpdateOneGroupChatNameById(testGroupChatId, newName)
	if nil != err {
		t.Fatal(err)
	}
	groupChat, _ := SelectOneGroupChatById(testGroupChatId, false)
	if groupChat.Name != newName {
		t.Fatalf("the new name should be: %s, but is: %s", newName, groupChat.Name)
	}
}

func TestUpdateOneGroupChatManagerById(t *testing.T) {
	err := UpdateOneGroupChatManagerById(testGroupChatId, userId2)
	if nil != err {
		t.Fatal(err)
	}
	groupChat, _ := SelectOneGroupChatById(testGroupChatId, false)
	if groupChat.ManagerId != userId2 {
		t.Fatalf("the new manager_id should be: %d, but is: %d", userId2, groupChat.ManagerId)
	}
}

func TestUpdateOneGroupChatAvatarById(t *testing.T) {
	newAvatar := "<new avatar pic name>"
	err := UpdateOneGroupChatAvatarById(testGroupChatId, newAvatar)
	if nil != err {
		t.Fatal(err)
	}
	groupChat, _ := SelectOneGroupChatById(testGroupChatId, false)
	if groupChat.Avatar != newAvatar {
		t.Fatalf("the new avatar should be: %s, but is: %s", newAvatar, groupChat.Avatar)
	}
}

func TestUpdateOneGroupChatQrCodeById(t *testing.T) {
	newQrCode := "<new qrCode pic name>"
	err := UpdateOneGroupChatQrCodeById(testGroupChatId, newQrCode)
	if nil != err {
		t.Fatal(err)
	}
	groupChat, _ := SelectOneGroupChatById(testGroupChatId, false)
	if groupChat.QrCode != newQrCode {
		t.Fatalf("the new qrCode should be: %s, but is: %s", newQrCode, groupChat.QrCode)
	}

}

func TestUpdateOneGroupChatIsDeleteById(t *testing.T) {
	err := UpdateOneGroupChatIsDeleteById(testGroupChatId, true)
	if nil != err {
		t.Fatal(err)
	}
	groupChat, _ := SelectOneGroupChatById(testGroupChatId, true)
	if !groupChat.IsDelete {
		t.Fatalf("the is_delete of group chat(%d) should be true, but is false", testGroupChatId)
	}

}

// Plus function test
func TestInsertOneNewGroupChatPlus(t *testing.T) {
	_ = DeleteOneGroupChatByIdReal(testGroupChatId)
	// insert with not existed user as manager
	_, err := InsertOneNewGroupChatPlus(groupChatName, groupAvatar, groupQrCode, 0)
	if nil == err {
		t.Fatal("should have an error, but not")
	}
	t.Logf("WantError: %s", err.Error())

	// insert one group chat normally
	managerId = userId1
	groupChat, err := InsertOneNewGroupChatPlus(groupChatName, groupAvatar, groupQrCode, managerId)
	if nil != err {
		t.Fatal(err)
	}
	testGroupChatId = groupChat.Id

	userGroupChat, err := SelectOneUserGroupChat(testGroupChatId, managerId)
	if nil != err {
		t.Fatal(err)
	}
	if userGroupChat.GroupId != testGroupChatId || userGroupChat.UserId != managerId {
		t.Fatal("insert one new row into tb_user_group_chat with wrong data")
	}

}

// ------------------------------------------------------------------

// Testing for operating 'tb_user_group_chat
func TestInsertOneNewUserGroupChat(t *testing.T) {
	_, err := InsertOneNewUserGroupChat(testGroupChatId, userId2, "user2")
	if nil != err {
		t.Fatal(err)
	}
}

func TestSelectOneUserGroupChat(t *testing.T) {
	userGroupChat, err := SelectOneUserGroupChat(testGroupChatId, userId2)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("userGroupChat: %v", userGroupChat)
}

func TestSelectAllUserGroupChat(t *testing.T) {
	data, err := SelectAllUserGroupChat()
	if nil != err {
		t.Fatal(err)
	}
	// after 'TestInsertOneNewGroupChatPlus' and 'TestInsertOneNewUserGroupChat', should have 2 rows data at least
	if len(data) < 2 {
		t.Fatal("there should have 1 userGroupChat at least, but not")
	}

}

func TestSelectManyUserGroupChatByGroupId(t *testing.T) {
	data, err := SelectManyUserGroupChatByGroupId(testGroupChatId, false)
	if nil != err {
		t.Fatal(err)
	}
	// after 'TestInsertOneNewGroupChatPlus' and 'TestInsertOneNewUserGroupChat', should have 2 rows data at least
	if len(data) < 2 {
		t.Fatal("should have 2 rows data at least, but not")
	}
	for index, userGroupChat := range data {
		t.Logf("selectByGroupId\tuserGroupChat: %d >> %v", index, userGroupChat)
	}
}

func TestSelectManyUserGroupChatByUserId(t *testing.T) {
	data, err := SelectManyUserGroupChatByUserId(userId2, false)
	if nil != err {
		t.Fatal(err)
	}
	if len(data) < 1 {
		t.Fatal("should have 1 row data at least, but not ")
	}
}

func TestSelectUsersIdOfGroupChat(t *testing.T) {
	ids, err := SelectUsersIdOfGroupChat(testGroupChatId, false)
	if nil != err {
		t.Fatal(err)
	}
	if len(ids) < 2 {
		t.Fatal("should have 2 id at least, but not")
	}
	for index, id := range ids {
		t.Logf("usersIdOfGroupChat: %d >> %d", index, id)
	}
}

func TestSelectGroupChatsIdOfUser(t *testing.T) {
	ids, err := SelectGroupChatsIdOfUser(userId2, false)
	if nil != err {
		t.Fatal(err)
	}
	if len(ids) < 1 {
		t.Fatal("should have 1 id at least, but not")
	}
	for index, id := range ids {
		t.Logf("groupChatsIdOfUser: %d >> %d", index, id)
	}
}

func TestUpdateOneUserGroupChatNote(t *testing.T) {
	newNote := "newNote"
	err := UpdateOneUserGroupChatNote(newNote, testGroupChatId, managerId)
	if nil != err {
		t.Fatal(err)
	}
	userGroupChat, _ := SelectOneUserGroupChat(testGroupChatId, managerId)
	if userGroupChat.UserNote != newNote {
		t.Fatal("Update userGroupChat note fail")
	}
}

func TestUpdateOneUserGroupChatIsDelete(t *testing.T) {
	err := UpdateOneUserGroupChatIsDelete(true, testGroupChatId, userId2)
	if nil != err {
		t.Fatal(err)
	}
	userGroupChat, _ := SelectOneUserGroupChat(testGroupChatId, userId2)
	if !userGroupChat.IsDelete {
		t.Fatal("'is_delete should be true, but false")
	}
}

// Plus function test
func TestSelectGroupChatUsersInfoPlus(t *testing.T) {
	data, err := SelectGroupChatUsersInfoPlus(testGroupChatId)
	if nil != err {
		t.Fatal(err)
	}
	// after 'TestUpdateOneUserGroupChatIsDelete', only should have 1 row data at least
	if len(data) < 1 {
		t.Fatal("should have 1 row data at least, but not")
	}
	for index, info := range data {
		t.Logf("groupChatUsersInfo: %d >> %v", index, info)
	}

}

func TestSelectUserGroupChatsInfoPlus(t *testing.T) {
	data, err := SelectUserGroupChatsInfoPlus(managerId)
	if nil != err {
		t.Fatal(err)
	}
	if len(data) < 1 {
		t.Fatal("should have 1 row data at least, but not")
	}
	for index, info := range data {
		t.Logf("userGroupChatsInfo: %d >> %v", index, info)
	}

}

// ------------------------------------------------------------------

// Testing for operating 'tb_subscription' table
var (
	testSubsId1, testSubsManagerId int64
	testSubsName                   = "subsName1"
	testIntroduce                  = "<test subscription introduce>"
	testSubsAvatar                 = "<test subscription avatar pic name>"
	testSubsQrCode                 = "<test subscription qr code pic name>"
)

func TestInsertOneNewSubscription(t *testing.T) {
	testSubsManagerId = userId1
	subs, err := InsertOneNewSubscription(testSubsName, testIntroduce, testSubsAvatar,
		testSubsQrCode, testSubsManagerId)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("insertOneNewSubscription: %v", subs)
	testSubsId1 = subs.Id
}

func TestSelectOneSubscriptionById(t *testing.T) {
	_, err := SelectOneSubscriptionById(testSubsId1, false)
	if nil != err {
		t.Fatal()
	}
}

func TestSelectOneSubscriptionByName(t *testing.T) {
	_, err := SelectOneSubscriptionByName(testSubsName, false)
	if nil != err {
		t.Fatal(err)
	}
}

func TestSelectManySubscriptionByManagerId(t *testing.T) {
	data, err := SelectManySubscriptionByManagerId(testSubsManagerId, false)
	if nil != err {
		t.Fatal(err)
	}
	if len(data) < 1 {
		t.Fatal("should have 1 subscription at least, but not")
	}
}

func TestUpdateOneSubscriptionNameById(t *testing.T) {
	newName := "subsNameN"
	err := UpdateOneSubscriptionNameById(newName, testSubsId1)
	if nil != err {
		t.Fatal(err)
	}
	data, err := SelectOneSubscriptionByName(newName, false)
	if nil != err {
		t.Fatal(err)
	}
	if data.Id != testSubsId1 {
		t.Fatal("the id is wrong")
	}
}

func TestUpdateOneSubscriptionManagerById(t *testing.T) {
	err := UpdateOneSubscriptionManagerById(userId2, testSubsId1)
	if nil != err {
		t.Fatal(err)
	}
	data, _ := SelectOneSubscriptionById(testSubsId1, false)
	if data.ManagerId != userId2 {
		t.Fatal("update manager id fail")
	}

}

func TestUpdateOneSubscriptionIntroById(t *testing.T) {
	newIntro := "<new introduce test>"
	err := UpdateOneSubscriptionIntroById(newIntro, testSubsId1)
	if nil != err {
		t.Fatal(err)
	}
	data, _ := SelectOneSubscriptionById(testSubsId1, false)
	if data.Intro != newIntro {
		t.Fatal("update introduce fail")
	}
}

func TestUpdateOneSubscriptionAvatarById(t *testing.T) {
	newAvatar := "<new avatar pic name>"
	err := UpdateOneSubscriptionAvatarById(newAvatar, testSubsId1)
	if nil != err {
		t.Fatal(err)
	}
	data, _ := SelectOneSubscriptionById(testSubsId1, false)
	if data.Avatar != newAvatar {
		t.Fatal("update the avatar fail")
	}

}

func TestUpdateOneSubscriptionQrCodeById(t *testing.T) {
	newQrCode := "<new qr code pic name>"
	err := UpdateOneSubscriptionQrCodeById(newQrCode, testSubsId1)
	if nil != err {
		t.Fatal(err)
	}
	data, _ := SelectOneSubscriptionById(testSubsId1, false)
	if data.QrCode != newQrCode {
		t.Fatal("update the qr_code fail")
	}
}

func TestUpdateOneSubscriptionIsDeleteById(t *testing.T) {
	err := UpdateOneSubscriptionIsDeleteById(true, testSubsId1)
	if nil != err {
		t.Fatal(err)
	}
	_, err = SelectOneSubscriptionById(testSubsId1, false)
	if err == nil {
		t.Fatal("should have an error")
	}
	data, _ := SelectOneSubscriptionById(testSubsId1, true)
	if data == nil {
		t.Fatal("update the is_delete fail")
	}

	// recover the data
	_ = UpdateOneSubscriptionIsDeleteById(false, testSubsId1)
}

// Plus function test
func TestInsertOneNewSubscriptionPlus(t *testing.T) {

	// insert with not existed manager
	_, err := InsertOneNewSubscriptionPlus(testSubsName, testIntroduce, testSubsAvatar,
		testSubsQrCode, 0)
	if nil == err {
		t.Fatal(err)
	}
	t.Logf("WantError: %s", err.Error())

	// insert normally
	data, err := InsertOneNewSubscriptionPlus(testSubsName, testIntroduce, testSubsAvatar, testSubsQrCode, userId2)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("newSubscription: %v", data)
	userSubs, err := SelectOneUserSubscription(data.Id, userId2)
	if nil != err {
		t.Fatal(err)
	}
	if data.Id != userSubs.SubsId || userId2 != userSubs.UserId {
		t.Fatal("the id should be same")
	}

	t.Logf("newUserSubscription: %v", userSubs)

	// clean the test data
	_ = DeleteOneSubscriptionReal(data.Id)

	err = DeleteOneUserSubscriptionReal(userSubs.SubsId, userSubs.UserId)
	if nil != err {
		t.Fatal(err)
	}
}

// ------------------------------------------------------------------
func TestInsertOneNewUserSubscription(t *testing.T) {
	_, err := InsertOneNewUserSubscription(testSubsId1, userId1)
	if nil != err {
		t.Fatal(err)
	}
	_, err = InsertOneNewUserSubscription(testSubsId1, userId1)
	if nil == err {
		t.Fatal("should have an error")
	}
	t.Logf("WantError: %s", err.Error())

	_, _ = InsertOneNewUserSubscription(testSubsId1, userId2)
}

func TestSelectOneUserSubscription(t *testing.T) {
	data, err := SelectOneUserSubscription(testSubsId1, userId1)
	if nil != err {
		t.Fatal(err)
	}
	if data.SubsId != testSubsId1 || data.UserId != userId1 {
		t.Fatal("ids is wrong")
	}
}

func TestSelectManyUserSubscriptionBySubsId(t *testing.T) {
	data, err := SelectManyUserSubscriptionBySubsId(testSubsId1, false)
	if nil != err {
		t.Fatal(err)
	}
	if len(data) < 2 {
		t.Fatal("should have 2 rows data at least")
	}
	for index, temp := range data {
		t.Logf("userSubs: %d >> %v", index, temp)
	}
}

func TestSelectManyUserSubscriptionByUserId(t *testing.T) {
	data, err := SelectManyUserSubscriptionByUserId(userId1, false)
	if nil != err {
		t.Fatal(err)
	}
	if len(data) < 1 {
		t.Fatal("should have 1 row data at least")
	}
	for index, temp := range data {
		t.Logf("userSubs: %d >> %v", index, temp)
	}
}

func TestSelectUsersIdOfSubscription(t *testing.T) {
	ids, err := SelectUsersIdOfSubscription(testSubsId1, false)
	if nil != err {
		t.Fatal(err)
	}
	if len(ids) < 2 {
		t.Fatal("should have tow rows data at least")
	}
}

func TestSelectSubscriptionsIdOfUser(t *testing.T) {
	ids, err := SelectSubscriptionsIdOfUser(userId1, false)
	if nil != err {
		t.Fatal(err)
	}
	if len(ids) < 1 {
		t.Fatal("should have 1 row data at least")
	}
}

func TestUpdateOneUserSubscriptionIsDelete(t *testing.T) {
	err := UpdateOneUserSubscriptionIsDelete(true, testSubsId1, userId1)
	if nil != err {
		t.Fatal(err)
	}
	_, err = SelectOneUserSubscription(testSubsId1, userId1)
	if nil != err {
		t.Fatal(err)
	}
	_ = UpdateOneUserSubscriptionIsDelete(false, testSubsId1, userId1)
}

// Plus function test

func TestSelectSubscriptionsOfUserPlus(t *testing.T) {
	subss, err := SelectSubscriptionsOfUserPlus(userId1)
	if nil != err {
		t.Fatal(err)
	}
	if len(subss) < 1 {
		t.Fatal("should have 1 row data at least")
	}

	for index, item := range subss {
		t.Logf("subsOfUser: %d >> %v", index, item)
	}
}

func TestSelectUsersOfSubscriptionPlus(t *testing.T) {
	users, err := SelectUsersOfSubscriptionPlus(testSubsId1)
	if nil != err {
		t.Fatal(err)
	}
	if len(users) < 2 {
		t.Fatal("should have 2 rows data at least")
	}

	for index, item := range users {
		t.Logf("userOfSubs: %d >> %v", index, item)
	}
}

// ------------------------------------------------------------------

//Clean the test data
func TestDeleteOneUserByIdReal(t *testing.T) {
	// this is delete one row data real
	err := DeleteOneUserByIdReal(userId1)
	if nil != err {
		t.Fatal(err)
	}
	user1, _ := SelectOneUserById(userId1, false)
	if nil != user1 {
		t.Fatal("delete one user by id fail, user1 id=:", userId1, "should be delete but not.")
	}

	err = DeleteOneUserByIdReal(userId2)
	if nil != err {
		t.Fatal(err)
	}
	user2, _ := SelectOneUserById(userId2, false)
	if nil != user2 {
		t.Fatal("delete one user by id fail, user2 id=:", userId2, "should be delete but not.")
	}

}

func TestDeleteOneFriendReal(t *testing.T) {
	err := DeleteOneFriendReal(selfId, friendId)
	if nil != err {
		t.Fatal(err)
	}
	err = DeleteOneFriendReal(friendId, selfId)
	if nil != err {
		t.Fatal(err)
	}
}

func TestDeleteOneGroupChatByIdReal(t *testing.T) {
	err := DeleteOneGroupChatByIdReal(testGroupChatId)
	if nil != err {
		t.Fatal(err)
	}
}

func TestDeleteOneUserGroupChatReal(t *testing.T) {
	err := DeleteOneUserGroupChatReal(testGroupChatId, managerId)
	if nil != err {
		t.Fatal(err)
	}

	err = DeleteOneUserGroupChatReal(testGroupChatId, userId2)
	if nil != err {
		t.Fatal(err)
	}

}

func TestDeleteOneSubscriptionReal(t *testing.T) {
	err := DeleteOneSubscriptionReal(testSubsId1)
	if nil != err {
		t.Fatal(err)
	}
}

func TestDeleteOneUserSubscriptionReal(t *testing.T) {
	err := DeleteOneUserSubscriptionReal(testSubsId1, userId1)
	if nil != err {
		t.Fatal(err)
	}
	err = DeleteOneUserSubscriptionReal(testSubsId1, userId2)
	if nil != err {
		t.Fatal(err)
	}
}
