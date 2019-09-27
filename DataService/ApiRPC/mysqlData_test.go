/*!!!!!! WARNING NOTE MESSAGE !!!!!!

You should start the rpc server of mysql data before execute the testing file.
	start function: ApiRPC.StartMySQLDataRPCServer()

And all the test functions are also example for using client to call the rpc server.
*/

package ApiRPC

import (
	mysqlPb "../Protos/mysqlProto"
	"testing"
)

// client for mysql data rpc server.
var testClient mysqlPb.MySQLBindServiceClient

func init() {
	testClient = GetMySQLDateClient()
}

var testAvatar = "<the test avatar pic name>"

// -----------------------------------------------------------------

var (
	testUserId1, testUserId2         int64
	testUserName                           = "testName"
	testUserEmail1, testUserEmail2         = "testUser1@test.com", "testUser2@test.com"
	testUserPassword                       = "<password: a hash value>"
	testUserGender                   int32 = 1
	testUserMobile                         = "13100000000"
	testUserQrCode1, testUserQrCode2       = "<the user qrCode pic name1>", "<the user qrCode pic name2>"
)

func TestMySQLData_PostSaveOneNewUser(t *testing.T) {
	// insert one new user with "isDelete= false"
	param := &mysqlPb.UserBasic{
		Name:     testUserName,
		Email:    testUserEmail1,
		Password: testUserPassword,
		Gender:   testUserGender,
		Mobile:   testUserMobile,
		Avatar:   testAvatar,
		QrCode:   testUserQrCode1,
		IsDelete: false,
	}

	user, err := testClient.PostSaveOneNewUser(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	testUserId1 = user.Id

	// insert one new user with "isDelete= true"
	param2 := &mysqlPb.UserBasic{
		Name:     testUserName,
		Email:    testUserEmail2,
		Password: testUserPassword,
		Gender:   testUserGender,
		Mobile:   testUserMobile,
		Avatar:   testAvatar,
		QrCode:   testUserQrCode2,
		IsDelete: true,
	}

	user2, err := testClient.PostSaveOneNewUser(getTimeOutCtx(3), param2)
	if nil != err {
		t.Fatal(err)
	}
	testUserId2 = user2.Id
}

func TestMySQLData_GetOneUserById(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testUserId1, IsDelete: false}
	user1, err := testClient.GetOneUserById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nuser1(%v)", user1)

	param = &mysqlPb.IdAndIsDelete{Id: testUserId2, IsDelete: true}
	user2, err := testClient.GetOneUserById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nuser2(%v)", user2)

}

func TestMySQLData_GetOneUserByEmail(t *testing.T) {
	param := &mysqlPb.EmailAndIsDelete{Email: testUserEmail1, IsDelete: false}
	user, err := testClient.GetOneUserByEmail(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nuser(%v)", user)

}

func TestMySQLData_GetUserListByName(t *testing.T) {
	param := &mysqlPb.NameAndIsDelete{Name: testUserName}
	userList, err := testClient.GetUserListByName(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, user := range userList.Data {
		t.Logf("\nuser(%v)", user)
	}
}

func TestMySQLData_GetOneUserPasswordById(t *testing.T) {
	param := &mysqlPb.Id{Value: testUserId1}
	password, err := testClient.GetOneUserPasswordById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nuser id= %d, password= %s", param.Value, password.Value)
}

func TestMySQLData_GetOneUserPasswordByEmail(t *testing.T) {
	param := &mysqlPb.Email{Value: testUserEmail1}
	password, err := testClient.GetOneUserPasswordByEmail(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nuser email= %s, password= %s", param.Value, password.Value)
}

func TestMySQLData_GetAllUserList(t *testing.T) {
	param := &mysqlPb.EmptyParam{}
	userList, err := testClient.GetAllUserList(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, user := range userList.Data {
		t.Logf("\nuser(%v)", user)
	}
}

func TestMySQLData_PutUserAvatarById(t *testing.T) {
	param := &mysqlPb.IdAndAvatar{Id: testUserId1, Avatar: "<new avatar pic name>"}
	_, err := testClient.PutUserAvatarById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutUserQrCodeById(t *testing.T) {
	param := &mysqlPb.IdAndQrCode{Id: testUserId1, QrCode: "<new qrCode pic name>"}
	_, err := testClient.PutUserQrCodeById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutUserPasswordById(t *testing.T) {
	param := &mysqlPb.IdAndPassword{Id: testUserId1, Password: "<new password: a hash value>"}
	_, err := testClient.PutUserPasswordById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutUserIsDeleteById(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testUserId2, IsDelete: false}
	_, err := testClient.PutUserIsDeleteById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutUserProfileByIdPlus(t *testing.T) {
	param := &mysqlPb.UserProfilePlus{Name: "newName", Mobile: "13199999999", Gender: 2, Id: testUserId2}
	_, err := testClient.PutUserProfileByIdPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------

var (
	testSelfId, testFriendId         int64 = 1, 2
	testFriendNote1, testFriendNote2       = "note1", "note2"
)

func TestMySQLData_PostSaveOneNewFriendship(t *testing.T) {
	param := &mysqlPb.FriendshipBasic{SelfId: testSelfId, FriendId: testFriendId, FriendNote: testFriendNote1}
	friendship, err := testClient.PostSaveOneNewFriendship(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nfriendship(%v)", friendship)
}

func TestMySQLData_PostSaveOneNewFriendPlus(t *testing.T) {
	param := &mysqlPb.FriendshipBasic{SelfId: testUserId1, FriendId: testUserId2, FriendNote: testFriendNote2}
	friendship, err := testClient.PostSaveOneNewFriendPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nfriendship(%v)", friendship)
}

func TestMySQLData_GetOneFriendship(t *testing.T) {
	param := &mysqlPb.FriendshipBasic{SelfId: testSelfId, FriendId: testFriendId}
	friendship, err := testClient.GetOneFriendship(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nfriendship(%v)", friendship)
}

func TestMySQLData_GetAllFriendshipList(t *testing.T) {
	param := &mysqlPb.EmptyParam{}
	friendshipList, err := testClient.GetAllFriendshipList(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, friendship := range friendshipList.Data {
		t.Logf("\nfriendship(%v)", friendship)
	}
}

func TestMySQLData_PutOneFriendIsAccept(t *testing.T) {
	param := &mysqlPb.FriendshipBasic{SelfId: testSelfId, FriendId: testFriendId, IsAccept: true}
	_, err := testClient.PutOneFriendIsAccept(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}

}

func TestMySQLData_PutOneFriendIsBlack(t *testing.T) {
	param := &mysqlPb.FriendshipBasic{SelfId: testSelfId, FriendId: testFriendId, IsBlack: true}
	_, err := testClient.PutOneFriendIsBlack(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutOneFriendIsDelete(t *testing.T) {
	param := &mysqlPb.FriendshipBasic{SelfId: testSelfId, FriendId: testFriendId, IsDelete: true}
	_, err := testClient.PutOneFriendIsDelete(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutOneFriendNote(t *testing.T) {
	param := &mysqlPb.FriendshipBasic{SelfId: testUserId1, FriendId: testUserId2, FriendNote: "newNote"}
	_, err := testClient.PutOneFriendNote(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}

}

func TestMySQLData_GetFriendsIdListByOptions(t *testing.T) {
	param := &mysqlPb.FriendshipBasic{SelfId: testSelfId, IsAccept: true, IsBlack: true, IsDelete: true}

	// after the previous functions, should found a friendship here
	idList, err := testClient.GetFriendsIdListByOptions(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nidList= %v", idList.Data)
}

func TestMySQLData_PutAcceptOneNewFriendPlus(t *testing.T) {
	param := &mysqlPb.FriendshipBasic{SelfId: testUserId2, FriendId: testUserId1, FriendNote: "cpNote", IsAccept: true}
	_, err := testClient.PutAcceptOneNewFriendPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}

}

func TestMySQLData_GetEffectiveFriendsIdListByIdPlus(t *testing.T) {
	param := &mysqlPb.Id{Value: testUserId1}
	idList, err := testClient.GetEffectiveFriendsIdListByIdPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nidList= %v", idList.Data)
}

func TestMySQLData_GetBlacklistFriendsIdListByIdPlus(t *testing.T) {
	param := &mysqlPb.Id{Value: testUserId1}
	idList, err := testClient.GetBlacklistFriendsIdListByIdPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nidList= %v", idList.Data)
}

func TestMySQLData_GetAllFriendsInfoPlus(t *testing.T) {
	param := &mysqlPb.Id{Value: testUserId1}
	infoList, err := testClient.GetAllFriendsInfoPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, info := range infoList.Data {
		t.Logf("\ninfo= (%v)", info)
	}
}

func TestMySQLData_GetEffectiveFriendsInfoPlus(t *testing.T) {
	param := &mysqlPb.Id{Value: testUserId1}
	infoList, err := testClient.GetEffectiveFriendsInfoPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, info := range infoList.Data {
		t.Logf("\ninfo= (%v)", info)
	}
}

func TestMySQLData_GetBlacklistFriendsInfoPlus(t *testing.T) {
	param := &mysqlPb.Id{Value: testUserId1}
	infoList, err := testClient.GetBlacklistFriendsInfoPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, info := range infoList.Data {
		t.Logf("\ninfo= (%v)", info)
	}
}

func TestMySQLData_PutDeleteOneFriendPlus(t *testing.T) {
	param := &mysqlPb.FriendshipBasic{SelfId: testUserId2, FriendId: testUserId1}
	_, err := testClient.PutDeleteOneFriendPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------

var (
	testGroupChatId1, testGroupChatId2 int64
	testGroupName                      = "groupName"
	testGroupQrCode1, testGroupQrCode2 = "<the group chat qr code pic name 1>", "<the group chat qr code pic name 2>"
)

func TestMySQLData_PostSaveOneNewGroupChat(t *testing.T) {
	param := &mysqlPb.GroupChatBasic{Name: testGroupName, Avatar: testAvatar, QrCode: testGroupQrCode1, ManagerId: testUserId1}
	groupChat, err := testClient.PostSaveOneNewGroupChat(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	testGroupChatId1 = groupChat.Id
}

func TestMySQLData_PostSaveOneNewGroupChatPlus(t *testing.T) {
	param := &mysqlPb.GroupChatBasic{Name: testGroupName, Avatar: testAvatar, QrCode: testGroupQrCode2, ManagerId: testUserId2}
	groupChat, err := testClient.PostSaveOneNewGroupChatPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	testGroupChatId2 = groupChat.Id
}

func TestMySQLData_GetOneGroupChatById(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testGroupChatId1}
	groupChat, err := testClient.GetOneGroupChatById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\ngroupChat= (%v)", groupChat)
}

func TestMySQLData_GetGroupChatListByName(t *testing.T) {
	param := &mysqlPb.NameAndIsDelete{Name: testGroupName}
	groupChatList, err := testClient.GetGroupChatListByName(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, groupChat := range groupChatList.Data {
		t.Logf("\ngroupChat= (%v)", groupChat)
	}
}

func TestMySQLData_GetGroupChatListByManagerId(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testUserId2}
	groupChatList, err := testClient.GetGroupChatListByManagerId(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, groupChat := range groupChatList.Data {
		t.Logf("\ngroupChat= (%v)", groupChat)
	}
}

func TestMySQLData_GetAllGroupChatList(t *testing.T) {
	param := &mysqlPb.EmptyParam{}
	groupChatList, err := testClient.GetAllGroupChatList(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, groupChat := range groupChatList.Data {
		t.Logf("\ngroupChat= (%v)", groupChat)
	}

}

func TestMySQLData_PutOneGroupChatNameById(t *testing.T) {
	param := &mysqlPb.IdAndName{Id: testGroupChatId1, Name: "newName"}
	_, err := testClient.PutOneGroupChatNameById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}

}

func TestMySQLData_PutOneGroupChatManagerById(t *testing.T) {
	param := &mysqlPb.GroupAndManagerId{GroupId: testGroupChatId1, ManagerId: testUserId2}
	_, err := testClient.PutOneGroupChatManagerById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutOneGroupChatAvatarById(t *testing.T) {
	param := &mysqlPb.IdAndAvatar{Id: testGroupChatId1, Avatar: "<new avatar pic name>"}
	_, err := testClient.PutOneGroupChatAvatarById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutOneGroupChatQrCodeById(t *testing.T) {
	param := &mysqlPb.IdAndQrCode{Id: testGroupChatId1, QrCode: "<new qr code pic name>"}
	_, err := testClient.PutOneGroupChatQrCodeById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutOneGroupChatIsDeleteById(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testGroupChatId1, IsDelete: true}
	_, err := testClient.PutOneGroupChatIsDeleteById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------

func TestMySQLData_PostSaveOneNewUserGroupChat(t *testing.T) {
	param := &mysqlPb.UserGroupChatRelate{GroupId: testGroupChatId2, UserId: testUserId1, UserNote: "note3"}
	_, err := testClient.PostSaveOneNewUserGroupChat(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_GetOneUserGroupChat(t *testing.T) {
	param := &mysqlPb.UserAndGroupId{UserId: testUserId1, GroupId: testGroupChatId2}
	userGroupChat, err := testClient.GetOneUserGroupChat(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nuserGroupChat= (%v)", userGroupChat)
}

func TestMySQLData_GetAllUserGroupChatList(t *testing.T) {
	param := &mysqlPb.EmptyParam{}
	dataList, err := testClient.GetAllUserGroupChatList(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, data := range dataList.Data {
		t.Logf("\nuserGroupChat= (%v)", data)
	}
}

func TestMySQLData_GetUserGroupChatListByGroupId(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testGroupChatId2}
	dataList, err := testClient.GetUserGroupChatListByGroupId(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\ndata count= %d", len(dataList.Data))
}

func TestMySQLData_GetUserGroupChatListByUserId(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testUserId2}
	dataList, err := testClient.GetUserGroupChatListByUserId(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\ndata count= %d", len(dataList.Data))
}

func TestMySQLData_GetUserIdListOfGroupChat(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testGroupChatId2}
	dataList, err := testClient.GetUserIdListOfGroupChat(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\ndata count= %d", len(dataList.Data))
}

func TestMySQLData_GetGroupChatIdListOfUser(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testUserId2}
	dataList, err := testClient.GetGroupChatIdListOfUser(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\ndata count= %d", len(dataList.Data))
}

func TestMySQLData_GetGroupChatUsersInfoPlus(t *testing.T) {
	param := &mysqlPb.Id{Value: testGroupChatId2}
	infoList, err := testClient.GetGroupChatUsersInfoPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, info := range infoList.Data {
		t.Logf("\ninfo= (%v)", info)
	}
}

func TestMySQLData_GetUserGroupChatsInfoPlus(t *testing.T) {
	param := &mysqlPb.Id{Value: testUserId2}
	infoList, err := testClient.GetUserGroupChatsInfoPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, info := range infoList.Data {
		t.Logf("\ninfo= (%v)", info)
	}
}

func TestMySQLData_PutOneUserGroupChatNote(t *testing.T) {
	param := &mysqlPb.UserGroupChatRelate{GroupId: testGroupChatId2, UserId: testUserId2, UserNote: "newNote2"}
	_, err := testClient.PutOneUserGroupChatNote(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutOneUserGroupChatIsDelete(t *testing.T) {
	param := &mysqlPb.UserGroupChatRelate{GroupId: testGroupChatId2, UserId: testUserId2, IsDelete: true}
	_, err := testClient.PutOneUserGroupChatIsDelete(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}

}

// -----------------------------------------------------------------
var (
	testSubsId1, testSubsId2         int64
	testSubsName1, testSubsName2     = "subsName1", "subsName2"
	testSubsIntro                    = "<the test subscription intro>"
	testSubsQrCode1, testSubsQrCode2 = "<the subs qr code pic name 1>", "<the subs qr code pic name2>"
)

func TestMySQLData_PostSaveOneNewSubscription(t *testing.T) {
	param := &mysqlPb.SubscriptionBasic{Name: testSubsName1, Intro: testSubsIntro, Avatar: testAvatar,
		QrCode: testSubsQrCode1, ManagerId: testUserId1}
	subscription, err := testClient.PostSaveOneNewSubscription(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	testSubsId1 = subscription.Id
}

func TestMySQLData_PostSaveOneNewSubscriptionPlus(t *testing.T) {
	param := &mysqlPb.SubscriptionBasic{Name: testSubsName2, Intro: testSubsIntro, Avatar: testAvatar,
		QrCode: testSubsQrCode2, ManagerId: testUserId2}
	subscription, err := testClient.PostSaveOneNewSubscriptionPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	testSubsId2 = subscription.Id
}

func TestMySQLData_GetOneSubscriptionById(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testSubsId1}
	subs, err := testClient.GetOneSubscriptionById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nsubscription= (%v)", subs)
}

func TestMySQLData_GetOneSubscriptionByName(t *testing.T) {
	param := &mysqlPb.NameAndIsDelete{Name: testSubsName1}
	subs, err := testClient.GetOneSubscriptionByName(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nsubscription= (%v)", subs)
}

func TestMySQLData_GetSubscriptionListByManagerId(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testUserId2}
	subsList, err := testClient.GetSubscriptionListByManagerId(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\n data conut= %d", len(subsList.Data))
}

func TestMySQLData_PutOneSubscriptionNameById(t *testing.T) {
	param := &mysqlPb.IdAndName{Id: testSubsId1, Name: "newName"}
	_, err := testClient.PutOneSubscriptionNameById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutOneSubscriptionManagerById(t *testing.T) {
	param := &mysqlPb.SubsAndManagerId{SubsId: testSubsId1, ManagerId: testUserId2}
	_, err := testClient.PutOneSubscriptionManagerById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutOneSubscriptionIntroById(t *testing.T) {
	param := &mysqlPb.IdAndIntro{Id: testSubsId1, Intro: "<new intro>"}
	_, err := testClient.PutOneSubscriptionIntroById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutOneSubscriptionAvatarById(t *testing.T) {
	param := &mysqlPb.IdAndAvatar{Id: testSubsId1, Avatar: "<new avatar pic name>"}
	_, err := testClient.PutOneSubscriptionAvatarById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutOneSubscriptionQrCodeById(t *testing.T) {
	param := &mysqlPb.IdAndQrCode{Id: testSubsId1, QrCode: "<new qr code pic name>"}
	_, err := testClient.PutOneSubscriptionQrCodeById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_PutOneSubscriptionIsDeleteById(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testSubsId1, IsDelete: true}
	_, err := testClient.PutOneSubscriptionIsDeleteById(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------

func TestMySQLData_PostSaveOneNewUserSubscription(t *testing.T) {
	param := &mysqlPb.UserSubscriptionRelate{SubsId: testSubsId2, UserId: testUserId1}
	_, err := testClient.PostSaveOneNewUserSubscription(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_GetOneUserSubscription(t *testing.T) {
	param := &mysqlPb.UserAndSubsId{SubsId: testSubsId2, UserId: testUserId1}
	relate, err := testClient.GetOneUserSubscription(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\nuseSubscription= (%v)", relate)
}

func TestMySQLData_GetUserSubscriptionListBySubsId(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testSubsId2}
	relateList, err := testClient.GetUserSubscriptionListBySubsId(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\ndata count= %d", len(relateList.Data))
}

func TestMySQLData_GetUserSubscriptionListByUserId(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testUserId1}
	relateList, err := testClient.GetUserSubscriptionListByUserId(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\ndata count= %d", len(relateList.Data))
}

func TestMySQLData_GetUserIdListOfSubscription(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testSubsId2}
	idList, err := testClient.GetUserIdListOfSubscription(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\ndata count= %d", len(idList.Data))
}

func TestMySQLData_GetSubscriptionIdListOfUser(t *testing.T) {
	param := &mysqlPb.IdAndIsDelete{Id: testUserId1}
	idList, err := testClient.GetSubscriptionIdListOfUser(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("\ndata count= %d", len(idList.Data))
}

func TestMySQLData_GetSubscriptionUsersInfoPlus(t *testing.T) {
	param := &mysqlPb.Id{Value: testSubsId2}
	infoList, err := testClient.GetSubscriptionUsersInfoPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, info := range infoList.Data {
		t.Logf("\ninfo= (%v)", info)
	}
}

func TestMySQLData_GetUserSubscriptionsInfoPlus(t *testing.T) {
	param := &mysqlPb.Id{Value: testUserId2}
	infoList, err := testClient.GetUserSubscriptionsInfoPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	for _, info := range infoList.Data {
		t.Logf("\ninfo= (%v)", info)
	}
}

// -----------------------------------------------------------------

// Clean testing data
func TestMySQLData_DeleteOneUserReal(t *testing.T) {
	param := &mysqlPb.Id{Value: testUserId1}
	_, err := testClient.DeleteOneUserReal(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}

	param2 := &mysqlPb.Id{Value: testUserId2}
	_, err = testClient.DeleteOneUserReal(getTimeOutCtx(3), param2)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_DeleteOneFriendshipReal(t *testing.T) {
	param := &mysqlPb.FriendshipBasic{SelfId: testSelfId, FriendId: testFriendId}
	_, err := testClient.DeleteOneFriendshipReal(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	param2 := &mysqlPb.FriendshipBasic{SelfId: testUserId1, FriendId: testUserId2}
	_, err = testClient.DeleteOneFriendshipReal(getTimeOutCtx(3), param2)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_DeleteOneGroupChatReal(t *testing.T) {
	param := &mysqlPb.Id{Value: testGroupChatId1}
	_, err := testClient.DeleteOneGroupChatReal(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}

	param2 := &mysqlPb.Id{Value: testGroupChatId2}
	_, err = testClient.DeleteOneGroupChatReal(getTimeOutCtx(3), param2)
	if nil != err {
		t.Fatal(err)
	}

}

func TestMySQLData_DeleteOneUserGroupChatReal(t *testing.T) {
	param := &mysqlPb.UserAndGroupId{GroupId: testGroupChatId2, UserId: testUserId1}
	_, err := testClient.DeleteOneUserGroupChatReal(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}

	param.UserId = testUserId2
	_, err = testClient.DeleteOneUserGroupChatReal(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_DeleteOneSubscriptionReal(t *testing.T) {
	param := &mysqlPb.Id{Value: testSubsId1}
	_, err := testClient.DeleteOneSubscriptionReal(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	param.Value = testSubsId2
	_, err = testClient.DeleteOneSubscriptionReal(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMySQLData_DeleteOneUserSubscriptionReal(t *testing.T) {
	param := &mysqlPb.UserAndSubsId{SubsId: testSubsId2, UserId: testUserId2}
	_, err := testClient.DeleteOneUserSubscriptionReal(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}

	param.UserId = testUserId1

	_, err = testClient.DeleteOneUserSubscriptionReal(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}
