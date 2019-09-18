package MySQLBind

import (
	"database/sql"
)

// Update 'name', 'mobile' and 'gender' columns of one row data which find by 'id' in 'tb_user_basic' table.
// It will checking the target user if existed and if nothing need be update before updating.
func UpdateOneUserProfileByIdPlus(name, mobile string, gender int32, id int64) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}

	// check the target user if existed
	row := tx.QueryRow(SelectOneUserByIdSQL, id, false)
	user, err := scanUserFromRow(row)
	if nil != err {
		_ = tx.Rollback()
		return err
	}

	// check if there nothing need to be update
	if user.Name == name && user.Mobile == mobile && user.Gender == gender {
		_ = tx.Rollback()
		return nil
	}

	// update the one row data
	_, err = tx.Exec(UpdateOneUserProfileSQL, name, mobile, gender, id)
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	return nil
}

const (
	SelectOneFriendIsBlackSQL  = `SELECT is_black FROM tb_friendship WHERE self_id= ? AND friend_id= ?`
	SelectOneFriendIsAcceptSQL = `SELECT is_accept FROM tb_friendship WHERE self_id= ? AND friend_id= ?`
)

// Insert one new row data into 'tb_friendship' table for the request that adding a new friend.
// It will checking the target user if existed and if there has an effect friendship record for them before insert
// the one new row data.
func InsertOneNewFriendPlus(selfId, friendId int64, friendNote string) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}

	// check the target user if existed
	row := tx.QueryRow(SelectOneUserByIdSQL, friendId, false)
	user, err := scanUserFromRow(row)
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	if friendNote == "" {
		friendNote = user.Name
	}

	// check the requester if in target user's blacklist
	isBlack := false
	row = tx.QueryRow(SelectOneFriendIsBlackSQL, friendId, selfId)
	err = row.Scan(&isBlack)
	if nil != err && sql.ErrNoRows != err {
		_ = tx.Rollback()
		return err
	}
	if isBlack {
		_ = tx.Rollback()
		return ErrInBlackList
	}

	// check the friendship is already effected.
	isAccept := false
	row = tx.QueryRow(SelectOneFriendIsAcceptSQL, friendId, selfId)
	err = row.Scan(&isAccept)
	if nil != err && sql.ErrNoRows != err {
		_ = tx.Rollback()
		return err
	}

	if isAccept {
		_ = tx.Rollback()
		return ErrFriendshipAlreadyInEffect
	}

	// insert the one row friendship record data finally
	_, err = tx.Exec(InsertOneNewFriendSQL, selfId, friendId, friendNote, friendNote)
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	return nil
}

// Insert one new row data into 'tb_friendship' table to accept the request that adding a new friend or not.
// It will checking the friend request record if existed and if there has an effect friendship record for then before
// insert the new one row data.And then, if the 'isAccept' is false, it will also insert one new data to record the
// requester is added into the recipient's blacklist
func UpdateAcceptOneNewFriendPlus(selfId, friendId int64, friendNote string, isAccept bool) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}
	// check the user who initiated the friend request if existed
	row := tx.QueryRow(SelectOneUserByIdSQL, friendId, false)
	user, err := scanUserFromRow(row)
	if nil != err {
		return err
	}
	if "" == friendNote {
		friendNote = user.Name
	}

	// check the friend request record if existed and the 'is_accept' value
	anotherIsAccept := false
	row = tx.QueryRow(SelectOneFriendIsAcceptSQL, friendId, selfId)
	err = row.Scan(&anotherIsAccept)
	if nil != err && sql.ErrNoRows != err {
		_ = tx.Rollback()
		return err
	}
	if sql.ErrNoRows == err {
		_ = tx.Rollback()
		return ErrNotTheFriendRequest
	}
	if anotherIsAccept {
		_ = tx.Rollback()
		return ErrFriendshipAlreadyInEffect
	}

	// accept or refuse the friendship request
	if isAccept {
		// add a friend relationship record for current user.
		_, err = tx.Exec(UpdateOneFriendIsAcceptSQL, selfId, friendId, isAccept, isAccept)
		if nil != err {
			_ = tx.Rollback()
			return err
		}
		// change the friend relationship record of requester, make the `is_accept` also be true
		_, err = tx.Exec(UpdateOneFriendIsAcceptSQL, friendId, selfId, isAccept, isAccept)
		if nil != err {
			_ = tx.Rollback()
			return err
		}

	} else {
		// when refuse the friend request, make the requester in blacklist of current user.
		_, err = tx.Exec(UpdateOneFriendIsBlackSQL, selfId, friendId, true, true)
		if nil != err {
			_ = tx.Rollback()
			return err
		}
	}

	// if the user give a note for the requester, execute it, if fail, don't need rollback
	if friendNote != "" {
		_, _ = tx.Exec(UpdateOneFriendNoteSQL, friendNote, selfId, friendId)
	}

	// commit all changes
	if err := tx.Commit(); nil != err {
		_ = tx.Rollback()
		return err
	}
	return nil
}

// Update the 'is_delete' column of the two row data which find by 'selfId' and 'friendId' in 'tb_friendship' table.
// Meaning if the user remove a friend from self's friend list, will also remove the user from the friend's friend list.
func UpdateDeleteOneFriendPlus(selfId, friendId int64) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}

	// update the record of current user.
	affect, err := tx.Exec(UpdateOneFriendIsDeleteSQL, true, selfId, friendId)
	if nil != err {
		_ = tx.Rollback()
		return err
	}

	affectCount, err := affect.RowsAffected()
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	if affectCount == 0 {
		_ = tx.Rollback()
		return ErrNotFriendYet
	}

	// update friend's record
	affect, err = tx.Exec(UpdateOneFriendIsDeleteSQL, true, friendId, selfId)
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	affectCount, err = affect.RowsAffected()
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	if affectCount == 0 {
		_ = tx.Rollback()
		return ErrNotFriendYet
	}

	err = tx.Commit()
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	return nil

}

const (
	SelectAllFriendsInfoPlusSQL = `SELECT id, friend_note, email, name, mobile, gender, avatar, is_accept, is_black, 
friends.is_delete FROM tb_user_basic AS basic, (SELECT friend_id, friend_note, is_accept,is_black, is_delete FROM
tb_friendship WHERE self_id= ? AND is_accept=TRUE) AS friends WHERE basic.id = friends.friend_id`

	SelectEffectiveFriendsInfoPlusSQL = `SELECT id, friend_note, email, name, mobile, gender, avatar, is_accept, 
is_black, friends.is_delete FROM tb_user_basic AS basic, (SELECT friend_id, friend_note, is_accept,is_black, is_delete 
FROM tb_friendship WHERE self_id= ? AND is_accept = TRUE AND is_black = FALSE AND is_delete = FALSE) AS friends WHERE 
basic.id = friends.friend_id`

	SelectBlacklistFriendsInfoPlusSQL = `SELECT id, friend_note, email, name, mobile, gender, avatar, is_accept,
is_black, friends.is_delete FROM tb_user_basic AS basic, (SELECT friend_id, friend_note, is_accept,is_black, is_delete 
FROM tb_friendship WHERE self_id= ? AND is_black=TRUE AND is_delete= FALSE) AS friends WHERE basic.id=friends.friend_id`

	SelectEffectiveFriendsIdPlusSQL = `SELECT friend_id FROM tb_friendship WHERE self_id = ? AND is_accept = TRUE 
AND is_black = FALSE  AND is_delete = FALSE`

	SelectBlacklistFriendsIdPlusSQL = `SELECT friend_id FROM tb_friendship WHERE self_id = ? AND is_black = TRUE 
AND is_delete = FALSE`
)

// Friend's information from 'tb_user_basic' and 'tb_friendship' table joined.
// Because protocol buffer 3 only have int32, so 'Gender' also use int32 here.
type JoinTableFriendInfo struct {
	Id       int64  `json:"id"`
	Note     string `json:"note"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	Gender   int32  `json:"gender"`
	Avatar   string `json:"avatar"`
	IsAccept bool   `json:"is_accept"`
	IsBlack  bool   `json:"is_black"`
	IsDelete bool   `json:"is_delete"`
}

// Private Function:
// Get the information of user's friends whom selected by the sql string
func selectFriendsInfo(sqlStr string, selfId int64) ([]*JoinTableFriendInfo, error) {
	rows, err := MySQLClient.Query(sqlStr, selfId)
	if nil != err {
		return nil, err
	}

	friendsInfo := make([]*JoinTableFriendInfo, 0)
	for rows.Next() {
		temp := new(JoinTableFriendInfo)
		err = rows.Scan(&(temp.Id), &(temp.Note), &(temp.Email), &(temp.Name), &(temp.Mobile),
			&(temp.Gender), &(temp.Avatar), &(temp.IsAccept), &(temp.IsBlack), &(temp.IsDelete))
		if nil != err {
			continue
		}
		friendsInfo = append(friendsInfo, temp)
	}
	return friendsInfo, nil

}

// Get the information of user's all friends, which just require the 'is_accept' is 'true' in 'tb_friendship' table.
// Meaning it would include the friends who in user's blacklist or the friends user have deleted.
func SelectAllFriendsInfoPlus(selfId int64) ([]*JoinTableFriendInfo, error) {
	return selectFriendsInfo(SelectAllFriendsInfoPlusSQL, selfId)
}

// Get the information of user's effective friends, which 'is_accept' is 'true', 'is_black' and 'is_delete' is 'false'
// in 'tb_friendship' table, meaning it would not include the friends who in user's blacklist or the friends have
// deleted.
func SelectEffectiveFriendsInfoPlus(selfId int64) ([]*JoinTableFriendInfo, error) {
	return selectFriendsInfo(SelectEffectiveFriendsInfoPlusSQL, selfId)
}

// Get the information of user's friends whom are in user's blacklist
func SelectBlacklistFriendsInfoPlus(selfId int64) ([]*JoinTableFriendInfo, error) {
	return selectFriendsInfo(SelectBlacklistFriendsInfoPlusSQL, selfId)
}

// Private Function:
// Get the id of user's friends whom selected by the sql string
func selectFiendsId(sqlStr string, selfId int64) ([]int64, error) {
	rows, err := MySQLClient.Query(sqlStr, selfId)
	if nil != err {
		return nil, err
	}
	friendsId := make([]int64, 0)
	for rows.Next() {
		temp := new(int64)
		err = rows.Scan(temp)
		if nil != err {
			continue
		}
		friendsId = append(friendsId, *temp)
	}
	return friendsId, nil
}

// Get value of 'friend_id' column belong to the one row data which selected from 'tb_friendship' table by 'self_id'.
// Require the 'is_accept' is 'true' and 'is_black' is 'false' at the same time, meaning it would not include the
// friends who in user's blacklist.
func SelectEffectiveFriendsIdPlus(selfId int64) ([]int64, error) {
	return selectFiendsId(SelectEffectiveFriendsIdPlusSQL, selfId)
}

// Get the id of user's friends whom are in user's blacklist
func SelectBlacklistFriendsIdPlus(selfId int64) ([]int64, error) {
	return selectFiendsId(SelectBlacklistFriendsIdPlusSQL, selfId)
}
