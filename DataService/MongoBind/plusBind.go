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
