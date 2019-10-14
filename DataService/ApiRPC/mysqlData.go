package ApiRPC

import (
	mysqlBind "../MySQLBind"
	mysqlPb "../Protos/mysqlProto"
	"context"
)

type MySQLData struct{}

// Translate the information of user to 'pb.UserBasic' from 'MySQLBind.TableUserBasic'.
func makePBUserBasic(user *mysqlBind.TableUserBasic) *mysqlPb.UserBasic {
	return &mysqlPb.UserBasic{
		Id:       user.Id,
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
		Mobile:   user.Mobile,
		Gender:   user.Gender,
		Avatar:   user.Avatar,
		QrCode:   user.QrCode,
		IsDelete: user.IsDelete}
}

// Translate the information of friendship to 'pb.FriendshipBasic' from 'MySQLBind.TableFriendship'.
func makePBFriendshipBasic(friendship *mysqlBind.TableFriendship) *mysqlPb.FriendshipBasic {
	return &mysqlPb.FriendshipBasic{
		SelfId:     friendship.SelfId,
		FriendId:   friendship.FriendId,
		FriendNote: friendship.FriendNote,
		IsAccept:   friendship.IsAccept,
		IsBlack:    friendship.IsBlack,
		IsDelete:   friendship.IsDelete}
}

// Translate the information of user' friend to 'pb.FriendInfoPlus' from 'MySQLBind.JoinTableFriendInfo'.
func makePBFriendInfoPlus(friendInfo *mysqlBind.JoinTableFriendInfo) *mysqlPb.FriendsInfoPlus {
	return &mysqlPb.FriendsInfoPlus{
		Id:       friendInfo.Id,
		Name:     friendInfo.Name,
		Email:    friendInfo.Email,
		Mobile:   friendInfo.Mobile,
		Gender:   friendInfo.Gender,
		Note:     friendInfo.Note,
		Avatar:   friendInfo.Avatar,
		IsAccept: friendInfo.IsAccept,
		IsBlack:  friendInfo.IsBlack,
		IsDelete: friendInfo.IsDelete}
}

// Translate the information of group chat to 'pb.GroupChatBasic' from 'MySQLBind.TableGroupChat'.
func makePBGroupChatBasic(groupChat *mysqlBind.TableGroupChat) *mysqlPb.GroupChatBasic {
	return &mysqlPb.GroupChatBasic{
		Id:        groupChat.Id,
		Name:      groupChat.Name,
		ManagerId: groupChat.ManagerId,
		Avatar:    groupChat.Avatar,
		QrCode:    groupChat.QrCode,
		IsDelete:  groupChat.IsDelete}
}

// Translate the information of user and group chat to 'pb.UserGroupChatRelate' from 'MySQLBind.TableUserGroupChat'.
func makePBUserGroupChatRelate(userGroupChat *mysqlBind.TableUserGroupChat) *mysqlPb.UserGroupChatRelate {
	return &mysqlPb.UserGroupChatRelate{
		GroupId:  userGroupChat.GroupId,
		UserId:   userGroupChat.UserId,
		UserNote: userGroupChat.UserNote,
		IsDelete: userGroupChat.IsDelete}
}

// Translate the information of user whom joined the group chat to 'pb.UserInfoInGroupChat' from
// 'MySQLBind.JoinTableGroupChatUsersInfo'
func makePBUserInfoInGroupChat(info *mysqlBind.JoinTableGroupChatUsersInfo) *mysqlPb.UserInfoInGroupChat {
	return &mysqlPb.UserInfoInGroupChat{
		GroupId:    info.GroupId,
		UserId:     info.UserId,
		UserNote:   info.UserNote,
		UserName:   info.UserName,
		UserGender: info.UserGender,
		UserEmail:  info.UserEmail,
		UserAvatar: info.UserAvatar}
}

// Translate the information of group chat which the user is joined to 'pb.GroupChatInfoOfUser' from
// 'MySQLBind.JoinTableUserGroupChatsInfo'.
func makePBGroupChatInfoOfUser(info *mysqlBind.JoinTableUserGroupChatsInfo) *mysqlPb.GroupChatInfoOfUser {
	return &mysqlPb.GroupChatInfoOfUser{
		UserId:      info.UserId,
		GroupId:     info.GroupId,
		GroupName:   info.GroupName,
		GroupAvatar: info.GroupAvatar,
		GroupQrCode: info.GroupQrCode}
}

// Translate the information of subscription to 'pb.SubscriptionBasic' from 'MySQLBind.TableSubscription'.
func makePBSubscriptionBasic(subscription *mysqlBind.TableSubscription) *mysqlPb.SubscriptionBasic {
	return &mysqlPb.SubscriptionBasic{
		Id:        subscription.Id,
		Name:      subscription.Name,
		ManagerId: subscription.ManagerId,
		Intro:     subscription.Intro,
		Avatar:    subscription.Avatar,
		QrCode:    subscription.QrCode,
		IsDelete:  subscription.IsDelete}
}

// Translate the information of user and subscription to 'pb.UserSubscriptionRelate' from 
// 'MySQLBind.TableUserSubscription'
func makePBUserSubscriptionRelate(relate *mysqlBind.TableUserSubscription) *mysqlPb.UserSubscriptionRelate {
	return &mysqlPb.UserSubscriptionRelate{
		SubsId:   relate.SubsId,
		UserId:   relate.UserId,
		IsDelete: relate.IsDelete}
}

// Translate the information of user whom followed the subscription to 'pb.UserInfoOfSubscription' from
// 'MySQLBind.JoinTableSubscriptionUsersInfo'.
func makePBUserInfoOfSubscription(info *mysqlBind.JoinTableSubscriptionUsersInfo) *mysqlPb.UserInfoOfSubscription {
	return &mysqlPb.UserInfoOfSubscription{
		SubsId:     info.SubsId,
		UserId:     info.UserId,
		UserEmail:  info.Email,
		UserName:   info.Name,
		UserGender: info.Gender}
}

// Translate the information of subscription which the user was followed to 'pb.SubscriptionInfoOfUser' from
// 'MySQLBind.JoinTableUserSubscriptionsInfo'.
func makePBSubscriptionInfoOfUser(info *mysqlBind.JoinTableUserSubscriptionsInfo) *mysqlPb.SubscriptionInfoOfUser {
	return &mysqlPb.SubscriptionInfoOfUser{
		UserId:     info.UserId,
		SubsId:     info.SubsId,
		SubsName:   info.Name,
		SubsIntro:  info.Intro,
		SubsAvatar: info.Avatar,
		SubsQrCode: info.QrCode}
}

func (obj *MySQLData) PostSaveOneNewUser(ctx context.Context, param *mysqlPb.UserBasic) (*mysqlPb.UserBasic, error) {
	user, err := mysqlBind.InsertOneNewUser(param.Email, param.Name, param.Password, param.Mobile,
		param.Gender, param.Avatar, param.QrCode, param.IsDelete)
	if nil != err {
		return nil, err
	}
	param.Id = user.Id
	return param, nil
}

func (obj *MySQLData) DeleteOneUserReal(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.Id, error) {
	err := mysqlBind.DeleteOneUserByIdReal(param.Value)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) GetOneUserById(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.UserBasic, error) {
	user, err := mysqlBind.SelectOneUserById(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return makePBUserBasic(user), nil
}

func (obj *MySQLData) GetOneUserByEmail(ctx context.Context, param *mysqlPb.EmailAndIsDelete) (*mysqlPb.UserBasic, error) {
	user, err := mysqlBind.SelectOneUserByEmail(param.Email, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return makePBUserBasic(user), nil
}

func (obj *MySQLData) GetUserListByName(ctx context.Context, param *mysqlPb.NameAndIsDelete) (*mysqlPb.UserBasicList, error) {
	users, err := mysqlBind.SelectManyUserByName(param.Name, param.IsDelete)
	if nil != err {
		return nil, err
	}

	data := make([]*mysqlPb.UserBasic, 0)
	for _, user := range users {
		data = append(data, makePBUserBasic(user))
	}
	return &mysqlPb.UserBasicList{Data: data}, nil
}

func (obj *MySQLData) GetOneUserPasswordById(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.Password, error) {
	password, err := mysqlBind.SelectOneUserPasswordById(param.Value)
	if nil != err {
		return nil, err
	}
	return &mysqlPb.Password{Value: password}, nil
}

func (obj *MySQLData) GetOneUserPasswordByEmail(ctx context.Context, param *mysqlPb.Email) (*mysqlPb.Password, error) {
	password, err := mysqlBind.SelectOneUserPasswordByEmail(param.Value)
	if nil != err {
		return nil, err
	}
	return &mysqlPb.Password{Value: password}, nil
}

func (obj *MySQLData) GetAllUserList(ctx context.Context, param *mysqlPb.EmptyParam) (*mysqlPb.UserBasicList, error) {
	users, err := mysqlBind.SelectAllUsers()
	if nil != err {
		return nil, err
	}

	data := make([]*mysqlPb.UserBasic, 0)
	for _, user := range users {
		data = append(data, makePBUserBasic(user))
	}
	return &mysqlPb.UserBasicList{Data: data}, nil
}

func (obj *MySQLData) PutUserAvatarById(ctx context.Context, param *mysqlPb.IdAndAvatar) (*mysqlPb.IdAndAvatar, error) {
	err := mysqlBind.UpdateOneUserAvatarById(param.Avatar, param.Id)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutUserQrCodeById(ctx context.Context, param *mysqlPb.IdAndQrCode) (*mysqlPb.IdAndQrCode, error) {
	err := mysqlBind.UpdateOneUserQrCodeById(param.QrCode, param.Id)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutUserPasswordById(ctx context.Context, param *mysqlPb.IdAndPassword) (*mysqlPb.IdAndPassword, error) {
	err := mysqlBind.UpdateOneUserPasswordById(param.Password, param.Id)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutUserIsDeleteById(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.IdAndIsDelete, error) {
	err := mysqlBind.UpdateOneUserIsDeleteById(param.IsDelete, param.Id)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutUserProfileByIdPlus(ctx context.Context, param *mysqlPb.UserProfilePlus) (*mysqlPb.UserProfilePlus, error) {
	err := mysqlBind.UpdateOneUserProfileByIdPlus(param.Name, param.Mobile, param.Gender, param.Id)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PostSaveOneNewFriendship(ctx context.Context, param *mysqlPb.FriendshipBasic) (*mysqlPb.FriendshipBasic, error) {
	err := mysqlBind.InsertOneNewFriend(param.SelfId, param.FriendId, param.FriendNote)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PostSaveOneNewFriendPlus(ctx context.Context, param *mysqlPb.FriendshipBasic) (*mysqlPb.FriendshipBasic, error) {
	err := mysqlBind.InsertOneNewFriendPlus(param.SelfId, param.FriendId, param.FriendNote)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) DeleteOneFriendshipReal(ctx context.Context, param *mysqlPb.FriendshipBasic) (*mysqlPb.FriendshipBasic, error) {
	err := mysqlBind.DeleteOneFriendReal(param.SelfId, param.FriendId)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) GetOneFriendship(ctx context.Context, param *mysqlPb.FriendshipBasic) (*mysqlPb.FriendshipBasic, error) {
	friendship, err := mysqlBind.SelectOneFriendship(param.SelfId, param.FriendId)
	if nil != err {
		return nil, err
	}
	return makePBFriendshipBasic(friendship), nil
}

func (obj *MySQLData) GetFriendsIdListByOptions(ctx context.Context, param *mysqlPb.FriendshipBasic) (*mysqlPb.IdList, error) {
	data, err := mysqlBind.SelectFriendsIdByOptions(param.SelfId, param.IsAccept, param.IsBlack, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return &mysqlPb.IdList{Data: data}, nil
}

func (obj *MySQLData) GetAllFriendshipList(ctx context.Context, param *mysqlPb.EmptyParam) (*mysqlPb.FriendshipBasicList, error) {
	friendships, err := mysqlBind.SelectAllFriendship()
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.FriendshipBasic, 0)
	for _, friendship := range friendships {
		data = append(data, makePBFriendshipBasic(friendship))
	}
	return &mysqlPb.FriendshipBasicList{Data: data}, nil
}

func (obj *MySQLData) GetEffectiveFriendsIdListByIdPlus(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.IdList, error) {
	data, err := mysqlBind.SelectEffectiveFriendsIdPlus(param.Value)
	if nil != err {
		return nil, err
	}
	return &mysqlPb.IdList{Data: data}, nil
}

func (obj *MySQLData) GetBlacklistFriendsIdListByIdPlus(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.IdList, error) {
	data, err := mysqlBind.SelectBlacklistFriendsIdPlus(param.Value)
	if nil != err {
		return nil, err
	}
	return &mysqlPb.IdList{Data: data}, nil
}

func (obj *MySQLData) GetAllFriendsInfoPlus(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.FriendsInfoListPlus, error) {
	friendInfoList, err := mysqlBind.SelectAllFriendsInfoPlus(param.Value)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.FriendsInfoPlus, 0)
	for _, friendInfo := range friendInfoList {
		data = append(data, makePBFriendInfoPlus(friendInfo))
	}
	return &mysqlPb.FriendsInfoListPlus{Data: data}, nil

}

func (obj *MySQLData) GetEffectiveFriendsInfoPlus(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.FriendsInfoListPlus, error) {
	friendInfoList, err := mysqlBind.SelectEffectiveFriendsInfoPlus(param.Value)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.FriendsInfoPlus, 0)
	for _, friendInfo := range friendInfoList {
		data = append(data, makePBFriendInfoPlus(friendInfo))
	}
	return &mysqlPb.FriendsInfoListPlus{Data: data}, nil
}

func (obj *MySQLData) GetBlacklistFriendsInfoPlus(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.FriendsInfoListPlus, error) {
	friendInfoList, err := mysqlBind.SelectBlacklistFriendsInfoPlus(param.Value)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.FriendsInfoPlus, 0)
	for _, friendInfo := range friendInfoList {
		data = append(data, makePBFriendInfoPlus(friendInfo))
	}
	return &mysqlPb.FriendsInfoListPlus{Data: data}, nil
}

func (obj *MySQLData) PutOneFriendNote(ctx context.Context, param *mysqlPb.FriendshipBasic) (*mysqlPb.FriendshipBasic, error) {
	err := mysqlBind.UpdateOneFriendNote(param.SelfId, param.FriendId, param.FriendNote)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneFriendIsAccept(ctx context.Context, param *mysqlPb.FriendshipBasic) (*mysqlPb.FriendshipBasic, error) {
	err := mysqlBind.UpdateOneFriendIsAccept(param.SelfId, param.FriendId, param.IsAccept)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneFriendIsBlack(ctx context.Context, param *mysqlPb.FriendshipBasic) (*mysqlPb.FriendshipBasic, error) {
	err := mysqlBind.UpdateOneFriendIsBlack(param.SelfId, param.FriendId, param.IsBlack)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneFriendIsDelete(ctx context.Context, param *mysqlPb.FriendshipBasic) (*mysqlPb.FriendshipBasic, error) {
	err := mysqlBind.UpdateOneFriendIsDelete(param.SelfId, param.FriendId, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutAcceptOneNewFriendPlus(ctx context.Context, param *mysqlPb.FriendshipBasic) (*mysqlPb.FriendshipBasic, error) {
	err := mysqlBind.UpdateAcceptOneNewFriendPlus(param.SelfId, param.FriendId, param.FriendNote, param.IsAccept)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutDeleteOneFriendPlus(ctx context.Context, param *mysqlPb.FriendshipBasic) (*mysqlPb.FriendshipBasic, error) {
	err := mysqlBind.UpdateDeleteOneFriendPlus(param.SelfId, param.FriendId)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PostSaveOneNewGroupChat(ctx context.Context, param *mysqlPb.GroupChatBasic) (*mysqlPb.GroupChatBasic, error) {
	groupChat, err := mysqlBind.InsertOneNewGroupChat(param.Name, param.Avatar, param.QrCode, param.ManagerId)
	if nil != err {
		return nil, err
	}
	param.Id = groupChat.Id
	return param, nil

}

func (obj *MySQLData) PostSaveOneNewGroupChatPlus(ctx context.Context, param *mysqlPb.GroupChatBasic) (*mysqlPb.GroupChatBasic, error) {
	groupChat, err := mysqlBind.InsertOneNewGroupChatPlus(param.Name, param.Avatar, param.QrCode, param.ManagerId)
	if nil != err {
		return nil, err
	}
	param.Id = groupChat.Id
	return param, nil
}

func (obj *MySQLData) DeleteOneGroupChatReal(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.Id, error) {
	err := mysqlBind.DeleteOneGroupChatByIdReal(param.Value)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) GetOneGroupChatById(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.GroupChatBasic, error) {
	groupChat, err := mysqlBind.SelectOneGroupChatById(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return makePBGroupChatBasic(groupChat), nil
}

func (obj *MySQLData) GetOneGroupChatByNameAndManager(ctx context.Context, param *mysqlPb.IdAndName) (*mysqlPb.GroupChatBasic, error) {
	groupChat, err := mysqlBind.SelectOneGroupChatByNameAndManager(param.Name, param.Id)
	if nil != err {
		return nil, err
	}
	return makePBGroupChatBasic(groupChat), nil
}

func (obj *MySQLData) GetGroupChatListByName(ctx context.Context, param *mysqlPb.NameAndIsDelete) (*mysqlPb.GroupChatList, error) {
	groupChatList, err := mysqlBind.SelectManyGroupChatByName(param.Name, param.IsDelete)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.GroupChatBasic, 0)
	for _, groupChat := range groupChatList {
		data = append(data, makePBGroupChatBasic(groupChat))
	}
	return &mysqlPb.GroupChatList{Data: data}, nil
}

func (obj *MySQLData) GetGroupChatListByManagerId(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.GroupChatList, error) {
	groupChatList, err := mysqlBind.SelectManyGroupChatByManagerId(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.GroupChatBasic, 0)
	for _, groupChat := range groupChatList {
		data = append(data, makePBGroupChatBasic(groupChat))
	}
	return &mysqlPb.GroupChatList{Data: data}, nil
}

func (obj *MySQLData) GetAllGroupChatList(ctx context.Context, param *mysqlPb.EmptyParam) (*mysqlPb.GroupChatList, error) {
	groupChatList, err := mysqlBind.SelectAllGroupChat()
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.GroupChatBasic, 0)
	for _, groupChat := range groupChatList {
		data = append(data, makePBGroupChatBasic(groupChat))
	}
	return &mysqlPb.GroupChatList{Data: data}, nil
}

func (obj *MySQLData) PutOneGroupChatNameById(ctx context.Context, param *mysqlPb.IdAndName) (*mysqlPb.IdAndName, error) {
	err := mysqlBind.UpdateOneGroupChatNameById(param.Id, param.Name)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneGroupChatManagerById(ctx context.Context, param *mysqlPb.GroupAndManagerId) (*mysqlPb.GroupAndManagerId, error) {
	err := mysqlBind.UpdateOneGroupChatManagerById(param.GroupId, param.ManagerId)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneGroupChatAvatarById(ctx context.Context, param *mysqlPb.IdAndAvatar) (*mysqlPb.IdAndAvatar, error) {
	err := mysqlBind.UpdateOneGroupChatAvatarById(param.Id, param.Avatar)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneGroupChatQrCodeById(ctx context.Context, param *mysqlPb.IdAndQrCode) (*mysqlPb.IdAndQrCode, error) {
	err := mysqlBind.UpdateOneGroupChatQrCodeById(param.Id, param.QrCode)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneGroupChatIsDeleteById(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.IdAndIsDelete, error) {
	err := mysqlBind.UpdateOneGroupChatIsDeleteById(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PostSaveOneNewUserGroupChat(ctx context.Context, param *mysqlPb.UserGroupChatRelate) (*mysqlPb.UserGroupChatRelate, error) {
	userGroupChat, err := mysqlBind.InsertOneNewUserGroupChat(param.GroupId, param.UserId, param.UserNote)
	if nil != err {
		return nil, err
	}
	return makePBUserGroupChatRelate(userGroupChat), nil
}

func (obj *MySQLData) DeleteOneUserGroupChatReal(ctx context.Context, param *mysqlPb.UserAndGroupId) (*mysqlPb.UserAndGroupId, error) {
	err := mysqlBind.DeleteOneUserGroupChatReal(param.GroupId, param.UserId)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) GetOneUserGroupChat(ctx context.Context, param *mysqlPb.UserAndGroupId) (*mysqlPb.UserGroupChatRelate, error) {
	userGroupChat, err := mysqlBind.SelectOneUserGroupChat(param.GroupId, param.UserId)
	if nil != err {
		return nil, err
	}
	return makePBUserGroupChatRelate(userGroupChat), nil
}

func (obj *MySQLData) GetAllUserGroupChatList(ctx context.Context, param *mysqlPb.EmptyParam) (*mysqlPb.UserGroupChatRelateList, error) {
	userGroupChatList, err := mysqlBind.SelectAllUserGroupChat()
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.UserGroupChatRelate, 0)
	for _, userGroupChat := range userGroupChatList {
		data = append(data, makePBUserGroupChatRelate(userGroupChat))
	}
	return &mysqlPb.UserGroupChatRelateList{Data: data}, nil
}

func (obj *MySQLData) GetUserGroupChatListByGroupId(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.UserGroupChatRelateList, error) {
	userGroupChatList, err := mysqlBind.SelectManyUserGroupChatByGroupId(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.UserGroupChatRelate, 0)
	for _, userGroupChat := range userGroupChatList {
		data = append(data, makePBUserGroupChatRelate(userGroupChat))
	}
	return &mysqlPb.UserGroupChatRelateList{Data: data}, nil
}

func (obj *MySQLData) GetUserGroupChatListByUserId(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.UserGroupChatRelateList, error) {
	userGroupChatList, err := mysqlBind.SelectManyUserGroupChatByUserId(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.UserGroupChatRelate, 0)
	for _, userGroupChat := range userGroupChatList {
		data = append(data, makePBUserGroupChatRelate(userGroupChat))
	}
	return &mysqlPb.UserGroupChatRelateList{Data: data}, nil
}

func (obj *MySQLData) GetUserIdListOfGroupChat(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.IdList, error) {
	data, err := mysqlBind.SelectUsersIdOfGroupChat(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return &mysqlPb.IdList{Data: data}, nil
}

func (obj *MySQLData) GetGroupChatIdListOfUser(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.IdList, error) {
	data, err := mysqlBind.SelectGroupChatsIdOfUser(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return &mysqlPb.IdList{Data: data}, nil
}

func (obj *MySQLData) GetGroupChatUsersInfoPlus(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.UserInfoInGroupChatListPlus, error) {
	infoList, err := mysqlBind.SelectGroupChatUsersInfoPlus(param.Value)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.UserInfoInGroupChat, 0)
	for _, info := range infoList {
		data = append(data, makePBUserInfoInGroupChat(info))
	}
	return &mysqlPb.UserInfoInGroupChatListPlus{Data: data}, nil
}

func (obj *MySQLData) GetUserGroupChatsInfoPlus(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.GroupChatInfoListOfUserPlus, error) {
	infoList, err := mysqlBind.SelectUserGroupChatsInfoPlus(param.Value)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.GroupChatInfoOfUser, 0)
	for _, info := range infoList {
		data = append(data, makePBGroupChatInfoOfUser(info))
	}
	return &mysqlPb.GroupChatInfoListOfUserPlus{Data: data}, nil
}

func (obj *MySQLData) PutOneUserGroupChatNote(ctx context.Context, param *mysqlPb.UserGroupChatRelate) (*mysqlPb.UserGroupChatRelate, error) {
	err := mysqlBind.UpdateOneUserGroupChatNote(param.UserNote, param.GroupId, param.UserId)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneUserGroupChatIsDelete(ctx context.Context, param *mysqlPb.UserGroupChatRelate) (*mysqlPb.UserGroupChatRelate, error) {
	err := mysqlBind.UpdateOneUserGroupChatIsDelete(param.IsDelete, param.GroupId, param.UserId)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PostSaveOneNewSubscription(ctx context.Context, param *mysqlPb.SubscriptionBasic) (*mysqlPb.SubscriptionBasic, error) {
	subscription, err := mysqlBind.InsertOneNewSubscription(param.Name, param.Intro, param.Avatar, param.QrCode, param.ManagerId)
	if nil != err {
		return nil, err
	}
	param.Id = subscription.Id
	return param, nil
}

func (obj *MySQLData) PostSaveOneNewSubscriptionPlus(ctx context.Context, param *mysqlPb.SubscriptionBasic) (*mysqlPb.SubscriptionBasic, error) {
	subscription, err := mysqlBind.InsertOneNewSubscriptionPlus(param.Name, param.Intro, param.Avatar, param.QrCode, param.ManagerId)
	if nil != err {
		return nil, err
	}
	param.Id = subscription.Id
	return param, nil
}

func (obj *MySQLData) DeleteOneSubscriptionReal(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.Id, error) {
	err := mysqlBind.DeleteOneSubscriptionReal(param.Value)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) GetOneSubscriptionById(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.SubscriptionBasic, error) {
	subscription, err := mysqlBind.SelectOneSubscriptionById(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return makePBSubscriptionBasic(subscription), nil
}

func (obj *MySQLData) GetOneSubscriptionByName(ctx context.Context, param *mysqlPb.NameAndIsDelete) (*mysqlPb.SubscriptionBasic, error) {
	subscription, err := mysqlBind.SelectOneSubscriptionByName(param.Name, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return makePBSubscriptionBasic(subscription), nil
}

func (obj *MySQLData) GetSubscriptionListByManagerId(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.SubscriptionBasicList, error) {
	subscriptionList, err := mysqlBind.SelectManySubscriptionByManagerId(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.SubscriptionBasic, 0)
	for _, subscription := range subscriptionList {
		data = append(data, makePBSubscriptionBasic(subscription))
	}
	return &mysqlPb.SubscriptionBasicList{Data: data}, nil
}

func (obj *MySQLData) PutOneSubscriptionNameById(ctx context.Context, param *mysqlPb.IdAndName) (*mysqlPb.IdAndName, error) {
	err := mysqlBind.UpdateOneSubscriptionNameById(param.Name, param.Id)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneSubscriptionManagerById(ctx context.Context, param *mysqlPb.SubsAndManagerId) (*mysqlPb.SubsAndManagerId, error) {
	err := mysqlBind.UpdateOneSubscriptionManagerById(param.ManagerId, param.SubsId)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneSubscriptionIntroById(ctx context.Context, param *mysqlPb.IdAndIntro) (*mysqlPb.IdAndIntro, error) {
	err := mysqlBind.UpdateOneSubscriptionIntroById(param.Intro, param.Id)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneSubscriptionAvatarById(ctx context.Context, param *mysqlPb.IdAndAvatar) (*mysqlPb.IdAndAvatar, error) {
	err := mysqlBind.UpdateOneSubscriptionAvatarById(param.Avatar, param.Id)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneSubscriptionQrCodeById(ctx context.Context, param *mysqlPb.IdAndQrCode) (*mysqlPb.IdAndQrCode, error) {
	err := mysqlBind.UpdateOneSubscriptionQrCodeById(param.QrCode, param.Id)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PutOneSubscriptionIsDeleteById(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.IdAndIsDelete, error) {
	err := mysqlBind.UpdateOneSubscriptionIsDeleteById(param.IsDelete, param.Id)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) PostSaveOneNewUserSubscription(ctx context.Context, param *mysqlPb.UserSubscriptionRelate) (*mysqlPb.UserSubscriptionRelate, error) {
	_, err := mysqlBind.InsertOneNewUserSubscription(param.SubsId, param.UserId)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) DeleteOneUserSubscriptionReal(ctx context.Context, param *mysqlPb.UserAndSubsId) (*mysqlPb.UserAndSubsId, error) {
	err := mysqlBind.DeleteOneUserSubscriptionReal(param.SubsId, param.UserId)
	if nil != err {
		return nil, err
	}
	return param, nil
}

func (obj *MySQLData) GetOneUserSubscription(ctx context.Context, param *mysqlPb.UserAndSubsId) (*mysqlPb.UserSubscriptionRelate, error) {
	relate, err := mysqlBind.SelectOneUserSubscription(param.SubsId, param.UserId)
	if nil != err {
		return nil, err
	}
	return makePBUserSubscriptionRelate(relate), nil
}

func (obj *MySQLData) GetUserSubscriptionListBySubsId(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.UserSubscriptionRelateList, error) {
	relateList, err := mysqlBind.SelectManyUserSubscriptionBySubsId(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.UserSubscriptionRelate, 0)
	for _, relate := range relateList {
		data = append(data, makePBUserSubscriptionRelate(relate))
	}
	return &mysqlPb.UserSubscriptionRelateList{Data: data}, nil
}

func (obj *MySQLData) GetUserSubscriptionListByUserId(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.UserSubscriptionRelateList, error) {
	relateList, err := mysqlBind.SelectManyUserSubscriptionByUserId(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.UserSubscriptionRelate, 0)
	for _, relate := range relateList {
		data = append(data, makePBUserSubscriptionRelate(relate))
	}
	return &mysqlPb.UserSubscriptionRelateList{Data: data}, nil
}

func (obj *MySQLData) GetUserIdListOfSubscription(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.IdList, error) {
	data, err := mysqlBind.SelectUsersIdOfSubscription(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return &mysqlPb.IdList{Data: data}, nil
}

func (obj *MySQLData) GetSubscriptionIdListOfUser(ctx context.Context, param *mysqlPb.IdAndIsDelete) (*mysqlPb.IdList, error) {
	data, err := mysqlBind.SelectSubscriptionsIdOfUser(param.Id, param.IsDelete)
	if nil != err {
		return nil, err
	}
	return &mysqlPb.IdList{Data: data}, nil
}

func (obj *MySQLData) GetSubscriptionUsersInfoPlus(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.UserInfoOfSubscriptionList, error) {
	infoList, err := mysqlBind.SelectUsersOfSubscriptionPlus(param.Value)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.UserInfoOfSubscription, 0)
	for _, info := range infoList {
		data = append(data, makePBUserInfoOfSubscription(info))
	}
	return &mysqlPb.UserInfoOfSubscriptionList{Data: data}, nil
}

func (obj *MySQLData) GetUserSubscriptionsInfoPlus(ctx context.Context, param *mysqlPb.Id) (*mysqlPb.SubscriptionInfoOfUserList, error) {
	infoList, err := mysqlBind.SelectSubscriptionsOfUserPlus(param.Value)
	if nil != err {
		return nil, err
	}
	data := make([]*mysqlPb.SubscriptionInfoOfUser, 0)
	for _, info := range infoList {
		data = append(data, makePBSubscriptionInfoOfUser(info))
	}
	return &mysqlPb.SubscriptionInfoOfUserList{Data: data}, nil
}

func (obj *MySQLData) PutOneUserSubscriptionIsDelete(ctx context.Context, param *mysqlPb.UserSubscriptionRelate) (*mysqlPb.UserSubscriptionRelate, error) {
	err := mysqlBind.UpdateOneUserSubscriptionIsDelete(param.IsDelete, param.SubsId, param.UserId)
	if nil != err {
		return nil, err
	}
	return param, nil
}
