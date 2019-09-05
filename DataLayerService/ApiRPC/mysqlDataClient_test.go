package ApiRPC

import "testing"

var (
	name     = "testName"
	email    = "testEmail@tt.com"
	mobile   = "13100000000"
	password = "this would be a hash value as password"
	gender   = 0

	userId int64
)

func TestSaveOneNewUser(t *testing.T) {
	user, err := SaveOneNewUser(name, email, mobile, password, gender)
	if nil != err {
		t.Error("NewOneUser Error: ", err)
	}
	if user.Id == 0 {
		t.Error("NewOneUser Error: the Id is zero")
	}
	if user.Name != name || user.Email != email || user.Mobile != mobile ||
		user.Password != password || user.Gender != int32(gender) {
		t.Error("NewOneUser Error: the raw data changed after insert")
	}
	if user.CreateTime == "" {
		t.Error("NewOneUser Error: the createTime is empty string")
	} else {
		t.Logf("NewOneUser Success: the user's id=%d, createTime=%s",
			user.Id, user.CreateTime)
	}

	userId = user.Id
}

func TestGetUserById(t *testing.T) {
	user, err := GetUserById(userId)
	if nil != err {
		t.Error("GetUserById Error: ", err)
	}
	if user.Id != userId {
		t.Error("GetUserById Error: the user's id is not equal the raw" +
			" query value")
	}
	if user.Name != name || user.Email != email || user.Mobile != mobile ||
		user.Password != password || user.Gender != int32(gender) {
		t.Error("GetUserById Error: the user's data is not equal the raw value")
	} else {
		t.Logf("GetUserById Success: the user's id=%d, createTime=%s",
			user.Id, user.CreateTime)
	}

}

func TestGetUserByEmail(t *testing.T) {
	user, err := GetUserByEmail(email)
	if nil != err {
		t.Error("GetUserByEmail Error: ", err)
	}
	if user.Id != userId {
		t.Error("GetUserByEmail Error: the user's id is not equal the raw" +
			" query value")
	}
	if user.Name != name || user.Email != email || user.Mobile != mobile ||
		user.Password != password || user.Gender != int32(gender) {
		t.Error("GetUserByEmail Error: the user's data is not equal the raw value")
	} else {
		t.Logf("GetUserByEmail Success: the user's id=%d, createTime=%s",
			user.Id, user.CreateTime)
	}

}

func TestGetUsersByName(t *testing.T) {
	users, err := GetUsersByName(name)
	if nil != err {
		t.Error("GetUsersByName Error: ", err)
	}
	if len(users.Data) == 0 {
		t.Error("GetUserByName Error: this result list is empty")
	} else {
		t.Logf("GetUsersByName Success: the zero index element: "+
			"user's id=%d, createTime=%s", users.Data[0].Id, users.Data[0].CreateTime)
	}

}

func TestPutUserBasicById(t *testing.T) {
	user, err := GetUserByEmail(email)
	if nil != err {
		t.Error("PutUserBasicById fail: ", err)
	}
	userId = user.Id
	newName, newMobile := "tNewName", "13211111111"
	newGender := 1
	newUser, err := PutUserBasicById(newName, newMobile, newGender, userId)
	if nil != err {
		t.Error("PutUserBasicById fail: ", err)
	}
	t.Logf("PutUserBasicById success: new name of the user is %s", newUser.Name)
}

func TestPutUserPasswordById(t *testing.T) {
	user, err := GetUserByEmail(email)
	if nil != err {
		t.Error("PutUserPasswordById fail: ", err)
	}
	newPassword := "this is a new password hash value 1"
	newUser, err := PutUserPasswordById(newPassword, user.Id)
	if nil != err {
		t.Error("PutUserPasswordById fail: ", err)
	} else {
		t.Logf("PutUserPasswordById success: new password of the user is:"+
			" %s", newUser.Password)
	}

}

func TestPutUserPasswordByEmail(t *testing.T) {
	newPassword := "this is a new password hash value 2"
	user, err := PutUserPasswordByEmail(newPassword, email)
	if nil != err {
		t.Error("PutUserPasswordByEmail fail: ", err)
	} else {
		t.Logf("PutUserPasswordByEmail success: new password of ther user is:"+
			" %s", user.Password)
	}

}

var (
	avatarPicName = "testAvatarPicName"
	qrCodePicName = "testQRCodePicName"
)

func TestPutUserAvatarById(t *testing.T) {
	user, err := GetUserByEmail(email)
	if nil != err {
		t.Error("PutUserAvatarById fail ", err)
	}
	data, err := PutUserAvatarById(avatarPicName, user.Id)
	if nil != err {
		t.Error("PutUserAvatarById fail: ", err)
	}
	if data.Avatar != avatarPicName {
		t.Error("PutUserAvatarById fail: the value after update not equal" +
			" the raw value")
	} else {
		t.Logf("PutUserAvatarById success: the name of the avatar picture "+
			"is: %s", data.Avatar)
	}

	userId = user.Id
}

func TestGetUserAvatarById(t *testing.T) {
	data, err := GetUserAvatarById(userId)
	if nil != err {
		t.Error("GetUserAvatarById fail: ", err)
	}
	if data.Avatar != avatarPicName {
		t.Error("GetUserAvatarById fail: the value queried is not equal to" +
			" raw value")
	} else {
		t.Logf("GetUserAvatarById success: the name of the avatar picture "+
			"is: %s", data.Avatar)
	}
}

func TestPutUserQRCodeById(t *testing.T) {
	data, err := PutUserQRCodeById(qrCodePicName, userId)
	if nil != err {
		t.Error("PutUserQRCodeById fail: ", err)
	}
	if data.QrCode != qrCodePicName {
		t.Error("PutUserQRCodeById fail: the value after update not equal the" +
			" raw value")
	} else {
		t.Logf("PutUserQRCodeById success: the name of the qrCode picture is:"+
			" %s", data.QrCode)
	}
}

func TestGetUserQRCodeById(t *testing.T) {
	data, err := GetUserQRCodeById(userId)
	if nil != err {
		t.Error("GetUserQRCodeById fail: ", err)
	}
	if data.QrCode != qrCodePicName {
		t.Error("GetUserQRCodeById fail: the value queried is not equal to raw value")
	} else {
		t.Logf("GetUserQRCodeById success: the name of the qrCode picture is:"+
			" %s", data.QrCode)
	}

}

var (
	name1   = "demoName"
	email1  = "demoEmail@tt.com"
	mobile1 = "13155555555"
)

func TestAddOneNewFriend(t *testing.T) {
	user1, _ := GetUserByEmail(email)
	user2, err := SaveOneNewUser(name1, email1, mobile1, password, gender)
	if nil != err {
		user2, _ = GetUserByEmail(email1)
	}
	selfId := user1.Id
	friendId := user2.Id
	_, err = AddOneNewFriend(selfId, friendId)
	if nil != err {
		t.Error("AddOneNewFriend fail: ", err)
	} else {
		t.Logf("AddOneNewFriend success: selfId = %d, friendId = %d", selfId, friendId)
	}
}

func TestPutOneFriendNote(t *testing.T) {
	user1, _ := GetUserByEmail(email)
	user2, _ := GetUserByEmail(email1)
	selfId, friendId := user1.Id, user2.Id
	newNote := "testNote"
	_, err := PutOneFriendNote(selfId, friendId, newNote)
	if nil != err {
		t.Error("PutOneFriendNote fail: ", err)
	} else {
		t.Logf("PutOneFriendNote success: selfId = %d, friendId = %d, "+
			"note = %s", selfId, friendId, newNote)
	}

}

func TestAcceptOneNewFriend(t *testing.T) {
	user1, _ := GetUserByEmail(email)
	user2, _ := GetUserByEmail(email1)
	selfId, friendId := user2.Id, user1.Id
	isAccept := true
	_, err := AcceptOneNewFriend(selfId, friendId, "", isAccept)
	if nil != err {
		t.Error("AcceptOneNewFriend fail: ", err)
	} else {
		t.Logf("AcceptOneNewFriend success: selfId = %d, friendId = %d, "+
			"isAccept = %t", selfId, friendId, isAccept)
	}
}

func TestAcceptOneNewFriend2(t *testing.T) {
	user1, _ := GetUserByEmail(email)
	user2, _ := GetUserByEmail(email1)
	selfId, friendId := user2.Id, user1.Id
	isAccept := false
	_, err := AcceptOneNewFriend(selfId, friendId, "", isAccept)
	if nil != err {
		t.Error("AcceptOneNewFriend fail: ", err)
	} else {
		t.Logf("AcceptOneNewFriend success: selfId = %d, friendId = %d, "+
			"isAccept = %t", selfId, friendId, isAccept)
	}
}

func TestPutFriendBlacklist(t *testing.T) {
	user1, _ := GetUserByEmail(email)
	user2, _ := GetUserByEmail(email1)
	selfId, friendId := user1.Id, user2.Id
	isBlack := true
	_, err := PutFriendBlacklist(selfId, friendId, isBlack)
	if nil != err {
		t.Error("PutFriendBlacklist fail: ", err)
	} else {
		t.Logf("PutFriendBlacklist success: selfId = %d, friendId = %d,"+
			" isBlack = %t", selfId, friendId, isBlack)
	}

}

func TestPutFriendBlacklist2(t *testing.T) {
	user1, _ := GetUserByEmail(email)
	user2, _ := GetUserByEmail(email1)
	selfId, friendId := user1.Id, user2.Id
	isBlack := false
	_, err := PutFriendBlacklist(selfId, friendId, isBlack)
	if nil != err {
		t.Error("PutFriendBlacklist fail: ", err)
	} else {
		t.Logf("PutFriendBlacklist success: selfId = %d, friendId = %d,"+
			" isBlack = %t", selfId, friendId, isBlack)
	}
}

func TestDeleteOneFriend(t *testing.T) {
	user1, _ := GetUserByEmail(email)
	user2, _ := GetUserByEmail(email1)
	selfId, friendId := user1.Id, user2.Id
	_, err := DeleteOneFriend(selfId, friendId)
	if nil != err {
		t.Error("DeleteOneFriend fail: ", err)
	} else {
		t.Logf("DeleteOneFriend success: selfId = %d, friendId = %d", selfId, friendId)
	}
}

func TestGetFriendshipInfo(t *testing.T) {
	user, _ := GetUserByEmail(email)
	selfId := user.Id
	ret, err := GetFriendshipInfo(selfId)
	if nil != err {
		t.Error("GetFriendshipInfo fail: ", err)
	}
	if len(ret.Data) == 0 {
		t.Error("GetFriendshipInfo fail: there should have friendship data " +
			"after TestAcceptOneNewFriend")
	} else {
		t.Logf("GetFriendshipInfo success: \n\t %v", ret.Data)
	}

}

func TestGetFriendsBasicInfo(t *testing.T) {
	user, _ := GetUserByEmail(email)
	selfId := user.Id
	ret, err := GetFriendsBasicInfo(selfId)
	if nil != err {
		t.Error("GetFriendsBasicInfo fail: ", err)
	}
	if len(ret.Data) == 0 {
		t.Error("GetFriendsBasicInfo fail: there should have friends data " +
			"after TestAcceptOneNewFriend")
	} else {
		t.Logf("GetFriendsBasicInfo success: \n\t%v", ret.Data)
	}

}
