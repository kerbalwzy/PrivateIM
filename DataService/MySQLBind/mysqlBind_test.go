package MySQLBind

import (
	"testing"
)

var (
	userId1, userId2 int64
	tempEmail1             = "test_1_email@test.com"
	tempEmail2             = "test_2_email@test.com"
	tempName               = "testName"
	tempPassword           = "test password 1, should be hash value"
	tempMobile             = "13100000000"
	tempGender       int32 = 1
	tempAvatar             = "temp avatar pic name"
	tempQrCode1            = "temp qr_code pic name 1"
	tempQrCode2            = "temp qr_code pic name 2"
)

func TestInsertOneNewUser(t *testing.T) {
	user, err := InsertOneNewUser(tempEmail1, tempName, tempPassword, tempMobile,
		tempGender, tempAvatar, tempQrCode1)
	if nil != err {
		t.Fatal(err)
	}
	if user.Id == 0 {
		t.Fatal("last insert id = 0")
	}
	userId1 = user.Id

	user2, err := InsertOneNewUser(tempEmail2, tempName, tempPassword, tempMobile,
		tempGender, tempAvatar, tempQrCode2)
	if nil != err {
		t.Fatal(err)
	}
	if user2.Id == 0 {
		t.Fatal("last insert id = 0")
	}
	userId2 = user2.Id
}

func TestSelectOneUserById(t *testing.T) {
	user, err := SelectOneUserById(userId1)
	if nil != err {
		t.Fatal(err)
	}
	if user.Email != tempEmail1 {
		t.Fatal("query an wrong user with wrong email:", user.Email)
	}

	user2, err := SelectOneUserById(userId2)
	if nil != err {
		t.Fatal(err)
	}
	if user2.Email != tempEmail2 {
		t.Fatal("query an wrong user with wrong email:", user2.Email)
	}

}

func TestSelectOneUserByEmail(t *testing.T) {
	user, err := SelectOneUserByEmail(tempEmail1)
	if nil != err {
		t.Fatal(err)
	}
	if user.Id != userId1 {
		t.Fatal("query a wrong user with wrong id:", user.Id)
	}

	user2, err := SelectOneUserByEmail(tempEmail2)
	if nil != err {
		t.Fatal(err)
	}
	if user2.Id != userId2 {
		t.Fatal("query a wrong user with wrong id:", user2.Id)
	}
}

func TestSelectManyUserByName(t *testing.T) {
	users, err := SelectManyUserByName(tempName)
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

func TestUpdateOneUserProfileById(t *testing.T) {
	newName := "NewName"
	newMobile := "13199999999"
	var newGender int32 = 2
	err := UpdateOneUserProfileById(newName, newMobile, newGender, userId1)
	if nil != err {
		t.Fatal(err)
	}
	user, err := SelectOneUserById(userId1)
	if nil != err {
		t.Fatal(err)
	}
	if user.Name != newName || user.Mobile != newMobile || user.Gender != newGender {
		t.Fatal("update user profile fail, some data changed wrong")
	}
}

func TestUpdateOneUserAvatarById(t *testing.T) {
	newAvatar := "new avatar pic name"
	err := UpdateOneUserAvatarById(newAvatar, userId1)
	if nil != err {
		t.Fatal(err)
	}
	user, err := SelectOneUserById(userId1)
	if nil != err {
		t.Fatal(err)
	}
	if user.Avatar != newAvatar {
		t.Fatal("update avatar by id fail, new avatar should be :", newAvatar, "but get:", user.Avatar)
	}
}

func TestUpdateOneUserQrCodeById(t *testing.T) {
	newQrcode := "new qr_code pic name"
	err := UpdateOneUserQrCodeById(newQrcode, userId1)
	if nil != err {
		t.Fatal(err)
	}
	user, err := SelectOneUserById(userId1)
	if nil != err {
		t.Fatal(err)
	}
	if user.QrCode != newQrcode {
		t.Fatal("update qr_code by id fail, should be:", newQrcode, "but get:", user.QrCode)
	}
}

func TestUpdateOneUserPasswordById(t *testing.T) {
	newPassword := "new password , should be hash value"
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
	user, err := SelectOneUserById(userId1)
	if nil != err {
		t.Fatal(err)
	}
	if user.IsDelete != true {
		t.Fatal("update isDelete by id fail, should be true, but get false")
	}

}

var (
	selfId, friendId         int64
	friendNote1, friendNote2 = "note1", "note2"
)

const selectFriendshipSQL = `SELECT self_id, friend_id, friend_note, is_accept, is_black, is_delete FROM tb_friendship 
WHERE self_id = ? AND friend_id = ?`

type tempFriendship struct {
	SelfId     int64  `json:"self_id"`
	FriendId   int64  `json:"friend_id"`
	FriendNote string `json:"friend_note"`
	IsAccept   bool   `json:"is_accept"`
	IsBlack    bool   `json:"is_black"`
	IsDelete   bool   `json:"is_delete"`
}

func selectFriendship(selfId, friendId int64) (*tempFriendship, error) {
	row := MySQLClient.QueryRow(selectFriendshipSQL, selfId, friendId)
	temp := new(tempFriendship)
	err := row.Scan(&(temp.SelfId), &(temp.FriendId), &(temp.FriendNote),
		&(temp.IsAccept), &(temp.IsBlack), &(temp.IsDelete))
	if nil != err {
		return nil, err
	}
	return temp, nil
}

func TestInsertOneNewFriend(t *testing.T) {
	selfId = userId1
	friendId = userId2
	err := InsertOneNewFriend(selfId, friendId, friendNote1)
	if nil != err {
		t.Fatal(err)
	}

	data, err := selectFriendship(selfId, friendId)
	if nil != err {
		t.Fatal(err)
	}
	if data.FriendNote != friendNote1 || data.SelfId != selfId || data.FriendId != friendId {
		t.Fatalf("insert one new friend fail, the data was wrong, get data:\n\t%v", data)
	}
}

func TestUpdateAcceptOneNewFriend(t *testing.T) {
	// accept the friend request
	err := UpdateAcceptOneNewFriend(friendId, selfId, friendNote2, true)
	if nil != err {
		t.Fatal(err)
	}
	data1, err := selectFriendship(selfId, friendId)
	if nil != err {
		t.Fatal(err)
	}
	if data1.IsAccept != true {
		t.Fatal("after accepted , requester's is_accept should be true, but still false")
	}
	data2, err := selectFriendship(friendId, selfId)
	if nil != err {
		t.Fatal()
	}
	if data2.IsAccept != true {
		t.Fatal("after accepted, the self's is_accept should be true, but still false")
	}
	if data2.FriendNote != friendNote2 {
		t.Fatal("save friend note error, should be:", friendNote2, "but get:", data2.FriendNote)
	}

	err = UpdateAcceptOneNewFriend(friendId, selfId, friendNote2, true)
	if err != ErrFriendshipAlreadyInEffect {
		t.Fatal("should get err:", ErrFriendshipAlreadyInEffect, "\nbut get:", err)
	}
	err = UpdateAcceptOneNewFriend(friendId, 10010, friendNote2, true)
	if err != ErrNotTheFriendRequest {
		t.Fatal("should get err:", ErrNotTheFriendRequest, "\nbut get:", err)
	}

}

func TestUpdateOneFriendNote(t *testing.T) {
	newNote := "newNote"
	err := UpdateOneFriendNote(selfId, friendId, newNote)
	if nil != err && ErrAffectZeroCount != err {
		t.Fatal(err)
	}
	data, _ := selectFriendship(selfId, friendId)
	if data.FriendNote != newNote {
		t.Fatal("update friend note fail, new note should be:", newNote, " but get:", data.FriendNote)
	}
}

func TestUpdateOneFriendIsBlack(t *testing.T) {
	err := UpdateOneFriendIsBlack(selfId, friendId, true)
	if nil != err && ErrAffectZeroCount != err {
		t.Fatal(err)
	}
	data, _ := selectFriendship(selfId, friendId)
	if data.IsBlack != true {
		t.Fatal("update friend is black fail, should be true, but get false")
	}

	err = UpdateOneFriendIsBlack(selfId, friendId, false)
	if nil != err {
		t.Fatal(err)
	}
	data, _ = selectFriendship(selfId, friendId)
	if data.IsBlack != false {
		t.Fatal("update friend is black fail, should be false, but get true")
	}
}

func TestUpdateOneFriendIsDelete(t *testing.T) {
	err := UpdateOneFriendIsDelete(selfId, friendId)
	if nil != err {
		t.Fatal(err)
	}
	data1, _ := selectFriendship(selfId, friendId)
	data2, _ := selectFriendship(friendId, selfId)

	if data1.IsDelete != true || data1.IsAccept != false || data1.IsBlack != false {
		t.Fatalf("change ship control columns value failed, \n data1:%v", data1)
	}

	if data2.IsDelete != true || data2.IsAccept != false || data2.IsBlack != false {
		t.Fatalf("change ship control columns value failed, \n data2:%v", data2)
	}
}

func TestUpdateAcceptOneNewFriend2(t *testing.T) {
	// refuse the friend request
	_ = InsertOneNewFriend(selfId, friendId, friendNote1)
	err := UpdateAcceptOneNewFriend(friendId, selfId, friendNote2, false)
	if nil != err {
		t.Fatal(err)
	}
	data, _ := selectFriendship(friendId, selfId)
	if data.IsBlack != true {
		t.Fatal("when refuse one friend request should add the requester into blacklist, but there not")
	}
	err = InsertOneNewFriend(selfId, friendId, friendNote1)
	if err != ErrInBlackList {
		t.Fatal("should get err:", ErrInBlackList, " but get:", err)
	}

	// re accept the friend request
	_ = UpdateOneFriendIsBlack(friendId, selfId, false)
	_ = InsertOneNewFriend(selfId, friendId, friendNote1)
	_ = UpdateAcceptOneNewFriend(friendId, selfId, friendNote2, true)
}

func TestSelectAllFriendsInfo(t *testing.T) {
	friends, err := SelectAllFriendsInfo(selfId)
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

func TestSelectEffectFriendsInfo(t *testing.T) {
	friends, err := SelectEffectFriendsInfo(selfId)
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

	// move the friend to blacklist
	_ = UpdateOneFriendIsBlack(selfId, friendId, true)
	friends, err = SelectEffectFriendsInfo(selfId)
	if nil != err {
		t.Fatal(err)
	}
	friendCount2 := len(friends)
	if friendCount2+1 != friendCount1 {
		t.Logf("should reduce 1 effect friend, but not")
	}
	t.Logf("freinds after move friend(%d) into blacklist:", friendId)
	for index, friend := range friends {
		t.Logf("effect friend: %d >> %v", index, friend)
	}

}

func TestSelectBlacklistFriendsInfo(t *testing.T) {
	_ = UpdateOneFriendIsBlack(selfId, friendId, true)
	friends, err := SelectBlacklistFriendsInfo(selfId)
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
	friends, err = SelectBlacklistFriendsInfo(selfId)
	if nil != err {
		t.Fatal(err)
	}
	friendCount2 := len(friends)
	if friendCount2+1 != friendCount1 {
		t.Fatal("should reduce 1 friend from blacklist, but not")
	}
	t.Logf("after move the friend(%d) out from blacklist:", friendId)
	for index, friend := range friends {
		t.Logf("blacklist: %d >> %v", index, friend)
	}
}

func TestSelectEffectFriendsId(t *testing.T) {
	ids, err := SelectEffectFriendsId(selfId)
	if nil != err {
		t.Fatal(err)
	}
	if len(ids) < 1 {
		t.Fatal("the count of id should right than 1, but not")
	}
}

func TestSelectBlacklistFriendsId(t *testing.T) {
	_ = UpdateOneFriendIsBlack(selfId, friendId, true)
	ids, err := SelectBlacklistFriendsId(selfId)
	if nil != err {
		t.Fatal(err)
	}
	if len(ids) < 1 {
		t.Fatal("the count of id should right than 1, but not")
	}
	_ = UpdateOneFriendIsBlack(selfId, friendId, false)
}

// clean the test data
func TestDeleteOneUserById(t *testing.T) {
	// this is delete one row data real
	err := DeleteOneUserById(userId1)
	if nil != err {
		t.Fatal(err)
	}
	user1, _ := SelectOneUserById(userId1)
	if nil != user1 {
		t.Fatal("delete one user by id fail, user1 id=:", userId1, "should be delete but not.")
	}

	err = DeleteOneUserById(userId2)
	if nil != err {
		t.Fatal(err)
	}
	user2, _ := SelectOneUserById(userId2)
	if nil != user2 {
		t.Fatal("delete one user by id fail, user2 id=:", userId2, "should be delete but not.")
	}

}

func TestDeleteOneFriendshipRecord(t *testing.T) {
	err := DeleteOneFriendshipRecord(selfId, friendId)
	if nil != err {
		t.Fatal(err)
	}
	err = DeleteOneFriendshipRecord(friendId, selfId)
	if nil != err {
		t.Fatal(err)
	}
}
