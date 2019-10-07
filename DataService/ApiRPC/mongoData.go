package ApiRPC

import (
	mongoBind "../MongoBind"
	mongoPb "../Protos/mongoProto"

	"context"
)

type MongoData struct{}

// Translate the data of delay message to 'pb.UserChatHistory' from 'mongoBind.DocUserChatHistory'.
func makePBUserChatHistory(data *mongoBind.DocUserChatHistory) *mongoPb.UserChatHistory {
	history := new(mongoPb.UserChatHistory)

	// copy data
	history.JoinId = data.Id
	for _, item := range data.History {
		temp := &mongoPb.DateAndMessage{Date: item.Date, MessageList: item.Messages}
		history.Data = append(history.Data, temp)
	}
	return history
}

// Translate the data of group chat history to 'pb.GroupChatHistory' from 'mongoBind.DocGroupChatHistory'
func makePBGroupChatHistory(data *mongoBind.DocGroupChatHistory) *mongoPb.GroupChatHistory {
	history := new(mongoPb.GroupChatHistory)

	// copy data
	history.Id = data.GroupId
	for _, item := range data.History {
		temp := &mongoPb.DateAndMessage{Date: item.Date, MessageList: item.Messages}
		history.Data = append(history.Data, temp)
	}
	return history
}

// Translate the data of message sent by subscription to 'pb.SubscriptionHistory' from 'mongoBind.DocSubscriptionHistory'
func makePBSubscriptionHistory(data *mongoBind.DocSubscriptionHistory) *mongoPb.SubscriptionHistory {
	history := new(mongoPb.SubscriptionHistory)

	// copy data
	history.Id = data.SubsId
	for _, item := range data.History {
		temp := &mongoPb.DateAndMessage{Date: item.Date, MessageList: item.Messages}
		history.Data = append(history.Data, temp)
	}
	return history
}

func (obj *MongoData) PutSaveDelayMessage(ctx context.Context, param *mongoPb.IdAndMessage) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateDelayMessage(param.Id, param.Message)
}

func (obj *MongoData) GetDelayMessage(ctx context.Context, param *mongoPb.Id) (*mongoPb.DelayMessage, error) {
	data, err := mongoBind.FindAndDeleteDelayMessage(param.Value)
	if nil != err {
		return nil, err
	}
	return &mongoPb.DelayMessage{Id: data.Id, MessageList: data.Message}, nil
}

func (obj *MongoData) PutSaveUserChatHistory(ctx context.Context, param *mongoPb.JoinIdAndMessage) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateUserChatHistoryByJoinId(param.JoinId, param.Message)
}

func (obj *MongoData) GetAllUserChatHistory(ctx context.Context, param *mongoPb.JoinId) (*mongoPb.UserChatHistory, error) {
	data, err := mongoBind.FindUserAllChatHistoryByJoinId(param.Value)
	if nil != err {
		return nil, err
	}
	return makePBUserChatHistory(data), nil
}

func (obj *MongoData) GetUserChatHistoryByDate(ctx context.Context, param *mongoPb.JoinIdAndDate) (*mongoPb.UserChatHistory, error) {
	data, err := mongoBind.FindUserChatHistoryByJoinIdAndDate(param.JoinId, param.Date)
	if nil != err {
		return nil, err
	}
	return makePBUserChatHistory(data), nil
}

func (obj *MongoData) GetUserChatHistoryByDateRange(ctx context.Context, param *mongoPb.JoinIdAndDateRange) (*mongoPb.UserChatHistory, error) {
	data, err := mongoBind.FindUserChatHistoryByJoinIdAndDateRange(param.JoinId, param.StartDate, param.EndDate)
	if nil != err {
		return nil, err
	}
	return makePBUserChatHistory(data), nil
}

func (obj *MongoData) PutSaveGroupChatHistory(ctx context.Context, param *mongoPb.IdAndMessage) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateGroupChatHistoryById(param.Id, param.Message)
}

func (obj *MongoData) GetAllGroupChatHistory(ctx context.Context, param *mongoPb.Id) (*mongoPb.GroupChatHistory, error) {
	data, err := mongoBind.FindAllGroupChatHistoryById(param.Value)
	if nil != err {
		return nil, err
	}
	return makePBGroupChatHistory(data), nil
}

func (obj *MongoData) GetGroupChatHistoryByDate(ctx context.Context, param *mongoPb.IdAndDate) (*mongoPb.GroupChatHistory, error) {
	data, err := mongoBind.FindGroupChatHistoryByIdAndDate(param.Id, param.Date)
	if nil != err {
		return nil, err
	}
	return makePBGroupChatHistory(data), nil
}

func (obj *MongoData) GetGroupChatHistoryByDateRange(ctx context.Context, param *mongoPb.IdAndDateRange) (*mongoPb.GroupChatHistory, error) {
	data, err := mongoBind.FindGroupChatHistoryByIdAndDateRange(param.Id, param.StartDate, param.EndDate)
	if nil != err {
		return nil, err
	}
	return makePBGroupChatHistory(data), nil
}

func (obj *MongoData) PutSaveSubscriptionHistory(ctx context.Context, param *mongoPb.IdAndMessage) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateSubscriptionHistoryById(param.Id, param.Message)
}

func (obj *MongoData) GetAllSubscriptionHistory(ctx context.Context, param *mongoPb.Id) (*mongoPb.SubscriptionHistory, error) {
	data, err := mongoBind.FindAllSubscriptionHistoryById(param.Value)
	if nil != err {
		return nil, err
	}
	return makePBSubscriptionHistory(data), nil
}

func (obj *MongoData) GetSubscriptionHistoryByDate(ctx context.Context, param *mongoPb.IdAndDate) (*mongoPb.SubscriptionHistory, error) {
	data, err := mongoBind.FindSubscriptionHistoryByIdAndDate(param.Id, param.Date)
	if nil != err {
		return nil, err
	}
	return makePBSubscriptionHistory(data), nil
}

func (obj *MongoData) GetSubscriptionHistoryByDateRange(ctx context.Context, param *mongoPb.IdAndDateRange) (*mongoPb.SubscriptionHistory, error) {
	data, err := mongoBind.FindSubscriptionHistoryByIdAndDateRange(param.Id, param.StartDate, param.EndDate)
	if nil != err {
		return nil, err
	}
	return makePBSubscriptionHistory(data), nil
}

func (obj *MongoData) PutUserFriendsAdd(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateUserFriendsToAddFriend(param.MainId, param.OtherId)
}

func (obj *MongoData) PutUserFriendsDel(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateUserFriendsToDelFriend(param.MainId, param.OtherId)
}

func (obj *MongoData) PutUserBlacklistAdd(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateUserBlacklistToAddUser(param.MainId, param.OtherId)
}

func (obj *MongoData) PutUserBlacklistDel(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateUserBlacklistToDelUser(param.MainId, param.OtherId)
}

func (obj *MongoData) GetUserFriendsAndBlacklist(ctx context.Context, param *mongoPb.Id) (*mongoPb.UserFriendsAndBlacklist, error) {
	data, err := mongoBind.FindUserFriendsAndBlacklistById(param.Value)
	if nil != err {
		return nil, err
	}
	return &mongoPb.UserFriendsAndBlacklist{Id: data.UserId, Friends: data.Friends, Blacklist: data.Blacklist}, nil
}

func (obj *MongoData) PutUserGroupChatsAdd(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateUserGroupChatsToAddOne(param.MainId, param.OtherId)
}

func (obj *MongoData) PutUserGroupChatsDel(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateUserGroupChatsToDelOne(param.MainId, param.OtherId)
}

func (obj *MongoData) GetUserGroupChats(ctx context.Context, param *mongoPb.Id) (*mongoPb.UserGroupChats, error) {
	data, err := mongoBind.FindUserGroupChatsById(param.Value)
	if nil != err {
		return nil, err
	}
	return &mongoPb.UserGroupChats{Id: data.UserId, Groups: data.Groups}, nil
}

func (obj *MongoData) PutUserSubscriptionsAdd(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateUserSubscriptionsToAddOne(param.MainId, param.OtherId)
}

func (obj *MongoData) PutUserSubscriptionsDel(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateUserSubscriptionsToDelOne(param.MainId, param.OtherId)
}

func (obj *MongoData) GetUserSubscriptions(ctx context.Context, param *mongoPb.Id) (*mongoPb.UserSubscriptions, error) {
	data, err := mongoBind.FindUserSubscriptionsById(param.Value)
	if nil != err {
		return nil, err
	}
	return &mongoPb.UserSubscriptions{Id: data.UserId, Subscriptions: data.Subscriptions}, nil
}

func (obj *MongoData) PutGroupChatUsersAdd(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateGroupChatUserToAddOne(param.MainId, param.OtherId)

}

func (obj *MongoData) PutGroupChatUsersDel(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateGroupChatUsersToDelOne(param.MainId, param.OtherId)
}

func (obj *MongoData) GetGroupChatUsers(ctx context.Context, param *mongoPb.Id) (*mongoPb.GroupChatUsers, error) {
	data, err := mongoBind.FindGroupChatUsersById(param.Value)
	if nil != err {
		return nil, err
	}
	return &mongoPb.GroupChatUsers{Id: data.GroupId, Users: data.Users}, nil
}

func (obj *MongoData) PutSubscriptionUsersAdd(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateSubscriptionUsersToAddOne(param.MainId, param.OtherId)
}

func (obj *MongoData) PutSubscriptionUsersDel(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateSubscriptionUsersToDelOne(param.MainId, param.OtherId)
}

func (obj *MongoData) GetSubscriptionUsers(ctx context.Context, param *mongoPb.Id) (*mongoPb.SubscriptionUsers, error) {
	data, err := mongoBind.FindSubscriptionUsersById(param.Value)
	if nil != err {
		return nil, err
	}
	return &mongoPb.SubscriptionUsers{Id: data.SubsId, Users: data.Users}, nil
}

func (obj *MongoData) PutMoveFriendIntoBlacklistPlus(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateMoveFriendIntoBlacklistPlus(param.MainId, param.OtherId)
}

func (obj *MongoData) PutMoveFriendOutFromBlacklistPlus(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateMoveFriendOutFromBlacklistPlus(param.MainId, param.OtherId)
}

func (obj *MongoData) PutUserJoinGroupChatPlus(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateMoveUserIntoGroupChat(param.MainId, param.OtherId)
}

func (obj *MongoData) PutUserQuitGroupChatPlus(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateMoveUserOutFromGroupChat(param.MainId, param.OtherId)
}

func (obj *MongoData) PutUserFollowSubscriptionPlus(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateMakeUserFollowSubscription(param.MainId, param.OtherId)
}

func (obj *MongoData) PutUserUnFollowSubscriptionPlus(ctx context.Context, param *mongoPb.DoubleId) (*mongoPb.EmptyResult, error) {
	return &mongoPb.EmptyResult{}, mongoBind.UpdateMakeUserUnFollowSubscription(param.MainId, param.OtherId)
}
