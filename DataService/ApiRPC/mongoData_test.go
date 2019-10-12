package ApiRPC

import (
	"../MongoBind"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"testing"
	"time"

	mongoPb "../Protos/mongoProto"
)

// -----------------------------------------------------------
var (
	testUId1, testUId2  int64 = 123, 234
	testJoinId                = "123_234"
	testGId, testSubsId int64 = 111, 222
	testDelayMessage          = []byte("<this is a delay message>")
	testHistoryMessage        = []byte("<this is a message have sent>")
)

var tClient mongoPb.MongoBindServiceClient

func init() {
	tClient = GetMongoDateClient()
}

func TestMongoData_PutSaveDelayMessage(t *testing.T) {
	param := &mongoPb.IdAndMessage{Id: testUId1, Message: testDelayMessage}
	_, err := tClient.PutSaveDelayMessage(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_GetDelayMessage(t *testing.T) {
	param := &mongoPb.Id{Value: testUId1}
	data, err := tClient.GetDelayMessage(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("The delay message for user(%d):", data.Id)
	for index, message := range data.MessageList {
		t.Logf("\t(%d): %s", index, message)
	}
}

func TestMongoData_PutSaveUserChatHistory(t *testing.T) {
	param := &mongoPb.JoinIdAndMessage{JoinId: testJoinId, Message: testHistoryMessage}
	_, err := tClient.PutSaveUserChatHistory(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func getTodayDateNum() int32 {
	year, month, day := time.Now().Date()
	dateStr := fmt.Sprintf("%d%02d%d", year, month, day)
	dateNum, _ := strconv.ParseInt(dateStr, 10, 32)
	return int32(dateNum)
}

func TestMongoData_GetAllUserChatHistory(t *testing.T) {
	param := &mongoPb.JoinId{Value: testJoinId}
	data, err := tClient.GetAllUserChatHistory(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("User chat history for join user(%s)", data.JoinId)
	for _, dateAndMessage := range data.Data {
		t.Logf("\thistory(%d)", dateAndMessage.Date)
		for index, message := range dateAndMessage.MessageList {
			t.Logf("\t(%d): %s", index, message)
		}
	}
}

func TestMongoData_GetUserChatHistoryByDate(t *testing.T) {
	param := &mongoPb.JoinIdAndDate{JoinId: testJoinId, Date: 20190230}
	_, err := tClient.GetUserChatHistoryByDate(getTimeOutCtx(3), param)
	if nil == err {
		t.Fatal("should have an error")
	}
	t.Logf("WantError: %s", err.Error())

	// normal testing
	param = &mongoPb.JoinIdAndDate{JoinId: testJoinId, Date: getTodayDateNum()}
	data, err := tClient.GetUserChatHistoryByDate(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("User chat history for join user(%s)", data.JoinId)
	for _, dateAndMessage := range data.Data {
		t.Logf("\thistory(%d)", dateAndMessage.Date)
		for index, message := range dateAndMessage.MessageList {
			t.Logf("\t(%d): %s", index, message)
		}
	}
}

func TestMongoData_GetUserChatHistoryByDateRange(t *testing.T) {
	param := &mongoPb.JoinIdAndDateRange{JoinId: testJoinId, StartDate: getTodayDateNum() - 1, EndDate: getTodayDateNum() + 1}
	data, err := tClient.GetUserChatHistoryByDateRange(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("User chat history for join user(%s)", data.JoinId)
	for _, dateAndMessage := range data.Data {
		t.Logf("\thistory(%d)", dateAndMessage.Date)
		for index, message := range dateAndMessage.MessageList {
			t.Logf("\t(%d): %s", index, message)
		}
	}
}

func TestMongoData_PutSaveGroupChatHistory(t *testing.T) {
	param := &mongoPb.IdAndMessage{Id: testGId, Message: testHistoryMessage}
	_, err := tClient.PutSaveGroupChatHistory(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_GetAllGroupChatHistory(t *testing.T) {
	param := &mongoPb.Id{Value: testGId}
	data, err := tClient.GetAllGroupChatHistory(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("User chat history for group chat (%d)", data.Id)
	for _, dateAndMessage := range data.Data {
		t.Logf("\thistory(%d)", dateAndMessage.Date)
		for index, message := range dateAndMessage.MessageList {
			t.Logf("\t(%d): %s", index, message)
		}
	}
}

func TestMongoData_GetGroupChatHistoryByDate(t *testing.T) {
	param := &mongoPb.IdAndDate{Id: testGId, Date: getTodayDateNum()}
	data, err := tClient.GetGroupChatHistoryByDate(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("User chat history for group chat (%d)", data.Id)
	for _, dateAndMessage := range data.Data {
		t.Logf("\thistory(%d)", dateAndMessage.Date)
		for index, message := range dateAndMessage.MessageList {
			t.Logf("\t(%d): %s", index, message)
		}
	}
}

func TestMongoData_GetGroupChatHistoryByDateRange(t *testing.T) {
	param := &mongoPb.IdAndDateRange{Id: testGId, StartDate: getTodayDateNum() - 1, EndDate: getTodayDateNum() + 1}
	data, err := tClient.GetGroupChatHistoryByDateRange(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("User chat history for group chat (%d)", data.Id)
	for _, dateAndMessage := range data.Data {
		t.Logf("\thistory(%d)", dateAndMessage.Date)
		for index, message := range dateAndMessage.MessageList {
			t.Logf("\t(%d): %s", index, message)
		}
	}
}

func TestMongoData_PutSaveSubscriptionHistory(t *testing.T) {
	param := &mongoPb.IdAndMessage{Id: testSubsId, Message: testHistoryMessage}
	_, err := tClient.PutSaveSubscriptionHistory(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_GetAllSubscriptionHistory(t *testing.T) {
	param := &mongoPb.Id{Value: testSubsId}
	data, err := tClient.GetAllSubscriptionHistory(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("User chat history for subscription(%d)", data.Id)
	for _, dateAndMessage := range data.Data {
		t.Logf("\thistory(%d)", dateAndMessage.Date)
		for index, message := range dateAndMessage.MessageList {
			t.Logf("\t(%d): %s", index, message)
		}
	}
}

func TestMongoData_GetSubscriptionHistoryByDate(t *testing.T) {
	param := &mongoPb.IdAndDate{Id: testSubsId, Date: getTodayDateNum()}
	data, err := tClient.GetSubscriptionHistoryByDate(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("User chat history for subscription(%d)", data.Id)
	for _, dateAndMessage := range data.Data {
		t.Logf("\thistory(%d)", dateAndMessage.Date)
		for index, message := range dateAndMessage.MessageList {
			t.Logf("\t(%d): %s", index, message)
		}
	}
}

func TestMongoData_GetSubscriptionHistoryByDateRange(t *testing.T) {
	param := &mongoPb.IdAndDateRange{Id: testSubsId, StartDate: getTodayDateNum() - 1, EndDate: getTodayDateNum() + 1}
	data, err := tClient.GetSubscriptionHistoryByDateRange(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("User chat history for subscription(%d)", data.Id)
	for _, dateAndMessage := range data.Data {
		t.Logf("\thistory(%d)", dateAndMessage.Date)
		for index, message := range dateAndMessage.MessageList {
			t.Logf("\t(%d): %s", index, message)
		}
	}
}

// -----------------------------------------------------------
var (
	testUId3, testUId4 int64 = 444, 555
)

func TestMongoData_PutUserFriendsAdd(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId1, OtherId: testUId3}
	_, err := tClient.PutUserFriendsAdd(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_PutUserBlacklistAdd(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId1, OtherId: testUId4}
	_, err := tClient.PutUserBlacklistAdd(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_GetUserFriendsAndBlacklist(t *testing.T) {
	param := &mongoPb.Id{Value: testUId1}
	data, err := tClient.GetUserFriendsAndBlacklist(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("get user(%d)'s friends and blacklist:", data.Id)
	t.Logf("\t friends: %v", data.Friends)
	t.Logf("\t blacklist: %v", data.Blacklist)
}

func TestMongoData_PutUserFriendsDel(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId1, OtherId: testUId3}
	_, err := tClient.PutUserFriendsDel(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_PutUserBlacklistDel(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId1, OtherId: testUId4}
	_, err := tClient.PutUserBlacklistDel(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_PutUserGroupChatsAdd(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId1, OtherId: testGId}
	_, err := tClient.PutUserGroupChatsAdd(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_GetUserGroupChats(t *testing.T) {
	param := &mongoPb.Id{Value: testUId1}
	data, err := tClient.GetUserGroupChats(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("user(%d)'s group chats: %v", data.Id, data.Groups)
}

func TestMongoData_PutUserGroupChatsDel(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId1, OtherId: testGId}
	_, err := tClient.PutUserGroupChatsDel(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_PutUserSubscriptionsAdd(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId1, OtherId: testSubsId}
	_, err := tClient.PutUserSubscriptionsAdd(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_GetUserSubscriptions(t *testing.T) {
	param := &mongoPb.Id{Value: testUId1}
	data, err := tClient.GetUserSubscriptions(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("user(%d)'s subscriptions: %v", data.Id, data.Subscriptions)

}

func TestMongoData_PutUserSubscriptionsDel(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId1, OtherId: testSubsId}
	_, err := tClient.PutUserSubscriptionsDel(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------

func TestMongoData_PutGroupChatUsersAdd(t *testing.T) {
	param := &mongoPb.XAndManagerAndUserId{Id: testGId, ManagerId: testUId1, UserId: testUId3}
	_, err := tClient.PutGroupChatUsersAdd(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_GetGroupChatUsers(t *testing.T) {
	param := &mongoPb.Id{Value: testGId}
	data, err := tClient.GetGroupChatUsers(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("the group chat(%d)'s manager id= %d", data.Id, data.ManagerId)
	t.Logf("the group chat(%d)'s users: %v", data.Id, data.Users)
}

func TestMongoData_PutGroupChatUsersDel(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testGId, OtherId: testUId3}
	_, err := tClient.PutGroupChatUsersDel(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_PutSubscriptionUsersAdd(t *testing.T) {
	param := &mongoPb.XAndManagerAndUserId{Id: testSubsId, ManagerId: testUId1, UserId: testUId3}
	_, err := tClient.PutSubscriptionUsersAdd(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_GetSubscriptionUsers(t *testing.T) {
	param := &mongoPb.Id{Value: testSubsId}
	data, err := tClient.GetSubscriptionUsers(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("the subscription(%d)'s manager id= %d", data.Id, data.ManagerId)
	t.Logf("the subscription(%d)'s users: %v", data.Id, data.Users)
}

func TestMongoData_PutSubscriptionUsersDel(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testSubsId, OtherId: testUId3}
	_, err := tClient.PutSubscriptionUsersDel(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------

// -----------------------------------------------------------
func TestMongoData_PutMoveFriendIntoBlacklistPlus(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId2, OtherId: testUId4}
	_, _ = tClient.PutUserFriendsAdd(getTimeOutCtx(3), param)
	_, err := tClient.PutMoveFriendIntoBlacklistPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_PutMoveFriendOutFromBlacklistPlus(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId2, OtherId: testUId4}
	_, err := tClient.PutMoveFriendOutFromBlacklistPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_PutUserJoinGroupChatPlus(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId4, OtherId: testGId}
	_, err := tClient.PutUserJoinGroupChatPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_PutUserQuitGroupChatPlus(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId4, OtherId: testGId}
	_, err := tClient.PutUserQuitGroupChatPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_PutUserFollowSubscriptionPlus(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId4, OtherId: testSubsId}
	_, err := tClient.PutUserFollowSubscriptionPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestMongoData_PutUserUnFollowSubscriptionPlus(t *testing.T) {
	param := &mongoPb.DoubleId{MainId: testUId4, OtherId: testSubsId}
	_, err := tClient.PutUserUnFollowSubscriptionPlus(getTimeOutCtx(3), param)
	if nil != err {
		t.Fatal(err)
	}
}

func TestCleanTheTestData(t *testing.T) {
	_ = MongoBind.CollDelayMessage.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testUId1})
	_ = MongoBind.CollUserChatHistory.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testJoinId})
	_ = MongoBind.CollGroupChatHistory.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testGId})
	_ = MongoBind.CollSubscriptionMsgHistory.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testSubsId})

	_ = MongoBind.CollUserFriends.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testUId1})
	_ = MongoBind.CollUserGroupChats.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testUId1})
	_ = MongoBind.CollUserSubscriptions.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testUId1})

	_ = MongoBind.CollUserFriends.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testUId2})
	_ = MongoBind.CollUserGroupChats.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testUId4})
	_ = MongoBind.CollUserSubscriptions.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testUId4})

	_ = MongoBind.CollGroupChatUsers.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testGId})
	_ = MongoBind.CollSubscriptionUsers.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": testSubsId})
}
