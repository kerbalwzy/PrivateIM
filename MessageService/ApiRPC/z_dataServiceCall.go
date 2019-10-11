package ApiRPC

import (
	"../RpcClientPbs/mongoPb"
	"log"
)

func SaveDelayMessage(userId int64, message []byte) error {
	// todo test code used in separate development, need remove later
	return nil

	// code to actually use
	param := &mongoPb.IdAndMessage{Id: userId, Message: message}
	_, err := GetMongoDateClient().PutSaveDelayMessage(getTimeOutCtx(3), param)
	return err
}

func GetDelayMessages(userId int64) ([][]byte, error) {
	// todo test code used in separate development, need remove later
	testDelayMessages := [][]byte{
		[]byte("<test delay message 1>"),
		[]byte("<test delay message 2>"),
		[]byte("<test delay message 3>"),
		[]byte("<test delay message 4>"),
	}
	return testDelayMessages, nil

	// code to actually use
	param := &mongoPb.Id{Value: userId}
	data, err := GetMongoDateClient().GetDelayMessage(getTimeOutCtx(3), param)
	if nil != err {
		return nil, err
	}
	return data.MessageList, nil
}

func GetUserFriendIdList(userId int64) ([]int64, []int64, error) {
	// todo test code used in separate development, need remove later
	return []int64{userId - 1, userId + 1, userId + 2}, []int64{userId + 10}, nil

	// code to actually use
	param := &mongoPb.Id{Value: userId}
	data, err := GetMongoDateClient().GetUserFriendsAndBlacklist(getTimeOutCtx(3), param)
	if nil != err {
		log.Printf(
			"[error] <GetUserFriendIdList> load friends and blacklist for user(%d) fail, detail: %s",
			userId, err.Error())
		return nil, nil, err
	}
	return data.Friends, data.Blacklist, nil
}

func GetUserGroupChats(userId int64) ([]int64, error) {
	// todo test code used in separate development, need remove later
	return []int64{111, 222, 333}, nil

	// code to actually use
	param := &mongoPb.Id{Value: userId}
	data, err := GetMongoDateClient().GetUserGroupChats(getTimeOutCtx(3), param)
	if nil != err {
		return nil, err
	}
	return data.Groups, nil
}

func GetGroupChatUsers(groupId int64) ([]int64, error) {
	// todo test code used in separate development, need remove later
	return []int64{0, 1, 2, 3}, nil

	// code to actually use
	param := &mongoPb.Id{Value: groupId}
	data, err := GetMongoDateClient().GetGroupChatUsers(getTimeOutCtx(3), param)
	if nil != err {
		return nil, err
	}
	return data.Users, nil

}
