package MongoBind

// ------------------------------------------------------------------------------------

// Update the 'friends' and 'blacklist' array to move a friend id into 'blacklist' from 'friends' .
func UpdateMoveFriendIntoBlacklistPlus(userId, friendId int64) error {
	err := updateIdArrayOfOneDocument(CollUserFriends, userId, friendId, "friends", "$pull")
	if nil != err {
		return err
	}
	err = updateIdArrayOfOneDocument(CollUserFriends, userId, friendId, "blacklist", "$addToSet")
	return err
}

// Update the 'friends' and 'blacklist' array to move a friend id into 'friends' from 'blacklist'
func UpdateMoveFriendOutFromBlacklistPlus(userId, friendId int64) error {
	err := updateIdArrayOfOneDocument(CollUserFriends, userId, friendId, "blacklist", "$pull")
	if nil != err {
		return err
	}
	err = updateIdArrayOfOneDocument(CollUserFriends, userId, friendId, "friends", "$addToSet")
	return err
}

// ------------------------------------------------------------------------------------

// Update the 'groups' array in document which in 'user_group_chats' collection, and the 'users' array in document
// which in 'group_chat_users' collection at the same time.
func UpdateMoveUserIntoGroupChat(userId, groupChatId int64) error {
	err := updateIdArrayOfOneDocument(CollGroupChatUsers, groupChatId, userId, "users", "$addToSet")
	if nil != err {
		return err
	}
	err = updateIdArrayOfOneDocument(CollUserGroupChats, userId, groupChatId, "groups", "$addToSet")
	return err
}

func UpdateMoveUserOutFromGroupChat(userId, groupChatId int64) error {
	err := updateIdArrayOfOneDocument(CollGroupChatUsers, groupChatId, userId, "users", "$pull")
	if nil != err {
		return err
	}
	err = updateIdArrayOfOneDocument(CollUserGroupChats, userId, groupChatId, "groups", "$pull")
	return err
}

// ------------------------------------------------------------------------------------

// Update the 'subscription' array in document which in 'user_subscriptions' collection, and the 'users' array in
// document which in 'subscription_users' collection at the same time.
func UpdateMakeUserFollowSubscription(userId, subsId int64) error {
	err := updateIdArrayOfOneDocument(CollSubscriptionUsers, subsId, userId, "users", "$addToSet")
	if nil != err {
		return err
	}
	err = updateIdArrayOfOneDocument(CollUserSubscriptions, userId, subsId, "subscriptions", "$addToSet")
	return err
}

func UpdateMakeUserUnFollowSubscription(userId, subsId int64) error {
	err := updateIdArrayOfOneDocument(CollSubscriptionUsers, subsId, userId, "users", "$pull")
	if nil != err {
		return err
	}
	err = updateIdArrayOfOneDocument(CollUserSubscriptions, userId, subsId, "subscriptions", "$pull")
	return err
}

// ------------------------------------------------------------------------------------
