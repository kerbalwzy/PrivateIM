package ApiRPC

import (
	"../RpcClientPbs/mysqlPb"
)

func GetUserByEmail(email string) (*mysqlPb.UserBasic, error) {
	param := &mysqlPb.EmailAndIsDelete{Email: email}
	return GetMySQLDataClient().GetOneUserByEmail(getTimeOutCtx(3), param)
}

func SaveOneNewUser(name, email, mobile, password, avatar, qrCode string, gender int) (*mysqlPb.UserBasic, error) {
	param := &mysqlPb.UserBasic{
		Name:     name,
		Email:    email,
		Mobile:   mobile,
		Password: password,
		Gender:   int32(gender),
		Avatar:   avatar,
		QrCode:   qrCode,
	}
	return GetMySQLDataClient().PostSaveOneNewUser(getTimeOutCtx(3), param)

}

func GetUserById(id int64) (*mysqlPb.UserBasic, error) {
	param := &mysqlPb.IdAndIsDelete{Id: id}
	return GetMySQLDataClient().GetOneUserById(getTimeOutCtx(3), param)
}

func GetUsersByName(name string) (*mysqlPb.UserBasicList, error) {
	param := &mysqlPb.NameAndIsDelete{Name: name}
	return GetMySQLDataClient().GetUserListByName(getTimeOutCtx(3), param)
}

func PutUserProfileById(userId int64, name, mobile string, gender int) (*mysqlPb.UserProfilePlus, error) {
	param := &mysqlPb.UserProfilePlus{Id: userId, Name: name, Mobile: mobile, Gender: int32(gender)}
	return GetMySQLDataClient().PutUserProfileByIdPlus(getTimeOutCtx(3), param)
}

func PutUserAvatarById(avatar string, id int64) (*mysqlPb.IdAndAvatar, error) {
	param := &mysqlPb.IdAndAvatar{Id: id, Avatar: avatar}
	return GetMySQLDataClient().PutUserAvatarById(getTimeOutCtx(3), param)
}

func GetUserPasswordById(id int64) (*mysqlPb.Password, error) {
	param := &mysqlPb.Id{Value: id}
	return GetMySQLDataClient().GetOneUserPasswordById(getTimeOutCtx(3), param)
}

func PutUserPasswordById(password string, id int64) error {
	param := &mysqlPb.IdAndPassword{Id: id, Password: password}
	_, err := GetMySQLDataClient().PutUserPasswordById(getTimeOutCtx(3), param)
	return err
}

func AddOneNewFriend(selfId, friendId int64, friendNote string) error {
	param := &mysqlPb.FriendshipBasic{SelfId: selfId, FriendId: friendId, FriendNote: friendNote}
	_, err := GetMySQLDataClient().PostSaveOneNewFriendPlus(getTimeOutCtx(3), param)
	return err
}

func AcceptOneNewFriend(selfId, friendId int64, friendNote string, isAccept bool) error {
	param := &mysqlPb.FriendshipBasic{SelfId: selfId, FriendId: friendId, FriendNote: friendNote, IsAccept: isAccept}
	_, err := GetMySQLDataClient().PutAcceptOneNewFriendPlus(getTimeOutCtx(3), param)
	return err

}

func PutOneFriendNote(selfId, friendId int64, friendNote string) error {
	param := &mysqlPb.FriendshipBasic{SelfId: selfId, FriendId: friendId, FriendNote: friendNote}
	_, err := GetMySQLDataClient().PutOneFriendNote(getTimeOutCtx(3), param)
	return err
}

func PutOneFriendIsBlack(selfId, friendId int64, isBlack bool) error {
	param := &mysqlPb.FriendshipBasic{SelfId: selfId, FriendId: friendId, IsBlack: isBlack}
	_, err := GetMySQLDataClient().PutOneFriendIsBlack(getTimeOutCtx(3), param)
	return err
}

func GetUserFriendsInfo(id int64) (*mysqlPb.FriendsInfoListPlus, error) {
	param := &mysqlPb.Id{Value: id}
	return GetMySQLDataClient().GetAllFriendsInfoPlus(getTimeOutCtx(3), param)
}

func DeleteOneFriend(selfId, friendId int64) error {
	param := &mysqlPb.FriendshipBasic{SelfId: selfId, FriendId: friendId}
	_, err := GetMySQLDataClient().PutDeleteOneFriendPlus(getTimeOutCtx(3), param)
	return err
}

func SaveOneNewGroupChat(name, avatar, qrCode string, managerId int64) (*mysqlPb.GroupChatBasic, error) {
	param := &mysqlPb.GroupChatBasic{
		Name:      name,
		ManagerId: managerId,
		Avatar:    avatar,
		QrCode:    qrCode,
	}
	return GetMySQLDataClient().PostSaveOneNewGroupChat(getTimeOutCtx(3), param)
}

func GetOneGroupChatByNameAndManger(name string, managerId int64) (*mysqlPb.GroupChatBasic, error) {
	param := &mysqlPb.IdAndName{Id: managerId, Name: name}
	return GetMySQLDataClient().GetOneGroupChatByNameAndManager(getTimeOutCtx(3), param)
}

func GetGroupChatsInfoTheUserJoined(userId int64) (*mysqlPb.GroupChatInfoListOfUserPlus, error) {
	param := &mysqlPb.Id{Value: userId}
	return GetMySQLDataClient().GetUserGroupChatsInfoPlus(getTimeOutCtx(3), param)
}

func GetOneUserGroupChat(userId, groupId int64) (*mysqlPb.UserGroupChatRelate, error) {
	param := &mysqlPb.UserAndGroupId{UserId: userId, GroupId: groupId}
	return GetMySQLDataClient().GetOneUserGroupChat(getTimeOutCtx(3), param)
}

func GetUsersInfoOfTheGroupChat(groupId int64) (*mysqlPb.UserInfoInGroupChatListPlus, error) {
	param := &mysqlPb.Id{Value: groupId}
	return GetMySQLDataClient().GetGroupChatUsersInfoPlus(getTimeOutCtx(3), param)
}
