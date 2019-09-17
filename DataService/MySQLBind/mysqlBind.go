package MySQLBind

import (
	"database/sql"
	"errors"
	"log"

	"github.com/bwmarrin/snowflake"
	_ "github.com/go-sql-driver/mysql"

	conf "../Config"
)

var (
	MySQLClient   = new(sql.DB)
	SnowFlakeNode = new(snowflake.Node)
)

func init() {
	var err error
	MySQLClient, err = sql.Open("mysql", conf.UserDbMySQLURI)
	if nil != err {
		log.Fatal(err)
	}
	SnowFlakeNode, err = snowflake.NewNode(0)
	if nil != err {
		log.Fatal(err)
	}
}

var (
	ErrAffectZeroCount = errors.New("0 row affected")
)

// User basic information in `tb_user_basic` table.
// Because protocol buffer 3 only have int32, so 'Gender' also use int32 here.
type TableUserBasic struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Mobile   string `json:"mobile"`
	Gender   int32  `json:"gender"`
	Avatar   string `json:"avatar"`
	QrCode   string `json:"qr_code"`
	IsDelete bool   `json:"is_delete"`
}

// Insert or update one row data by given sql string and new values.
func insertOrUpdateOneRowData(sqlStr string, args ...interface{}) error {
	// start a transaction
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}
	// update the row data with new value
	ret, err := tx.Exec(sqlStr, args...)
	if nil != err {
		_ = tx.Rollback()
		return err
	}

	affectCount, err := ret.RowsAffected()
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	if affectCount == 0 {
		_ = tx.Rollback()
		return ErrAffectZeroCount
	}
	// commit the transaction finally
	err = tx.Commit()
	if nil != err {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// Scan one user's information from the 'row'.
func scanUserFromRow(row *sql.Row) (*TableUserBasic, error) {
	user := new(TableUserBasic)
	err := row.Scan(&(user.Id), &(user.Email), &(user.Name), &(user.Mobile),
		&(user.Gender), &(user.Avatar), &(user.QrCode), &(user.IsDelete))
	if nil != err {
		return nil, err
	}
	return user, nil
}

const InsertOneNewUserSQL = `INSERT INTO tb_user_basic (id, email, name, password, mobile, gender, avatar, qr_code) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

// Insert one row new data for saving information of new user.
// The 'id' will auto generate by 'SnowFlakeNode', the 'is_delete' will use default value 'false', and others will
// get from params of the function.
func InsertOneNewUser(email, name, password, mobile string,
	gender int32, avatar, qrCode string) (*TableUserBasic, error) {

	// try to insert user data into database
	id := SnowFlakeNode.Generate().Int64()
	err := insertOrUpdateOneRowData(InsertOneNewUserSQL, id, email, name, password, mobile, gender, avatar, qrCode)
	if nil != err {
		return nil, err
	}

	// return an user's basic information to follow the REST style
	user := &TableUserBasic{Id: id, Email: email, Name: name, Password: password,
		Mobile: mobile, Gender: gender, Avatar: avatar, QrCode: qrCode}

	return user, nil
}

const (
	SelectUserBaseSQL = `SELECT id, email, name, mobile, gender, avatar, qr_code, is_delete FROM tb_user_basic`

	SelectOneUserByIdSQL            = SelectUserBaseSQL + ` WHERE id = ?`
	SelectOneUserByEmailSQL         = SelectUserBaseSQL + ` WHERE email = ?`
	SelectManyUserByNameSQL         = SelectUserBaseSQL + ` WHERE name = ?`
	SelectOneUserPasswordByIdSQL    = `SELECT password FROM tb_user_basic WHERE id = ?`
	SelectOneUserPasswordByEmailSQL = `SELECT password FROM tb_user_basic WHERE email = ?`
)

// Select one row data from 'tb_user_basic' table by 'id' that given.
func SelectOneUserById(id int64) (*TableUserBasic, error) {
	row := MySQLClient.QueryRow(SelectOneUserByIdSQL, id)
	return scanUserFromRow(row)
}

// Select one row data from 'tb_user_basic' table by 'email' that given.
func SelectOneUserByEmail(email string) (*TableUserBasic, error) {
	row := MySQLClient.QueryRow(SelectOneUserByEmailSQL, email)
	return scanUserFromRow(row)

}

// Select many row data from 'tb_user_basic' table by 'name' that given.
func SelectManyUserByName(name string) ([]*TableUserBasic, error) {
	rows, err := MySQLClient.Query(SelectManyUserByNameSQL, name)
	if nil != err {
		return nil, err
	}

	users := make([]*TableUserBasic, 0)
	for rows.Next() {
		user := new(TableUserBasic)
		err := rows.Scan(&(user.Id), &(user.Email), &(user.Name), &(user.Mobile),
			&(user.Gender), &(user.Avatar), &(user.QrCode), &(user.IsDelete))
		if nil != err {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

// Get the value of 'password' column belong to one row data which selected from 'tb_user_basic' table by 'id'.
func SelectOneUserPasswordById(id int64) (string, error) {
	row := MySQLClient.QueryRow(SelectOneUserPasswordByIdSQL, id)
	var password string
	err := row.Scan(&password)
	return password, err
}

// Get the value of 'password' column belong to one row data which selected from 'tb_user_basic' table by 'email'
func SelectOneUserPasswordByEmail(email string) (string, error) {
	row := MySQLClient.QueryRow(SelectOneUserPasswordByEmailSQL, email)
	var password string
	err := row.Scan(&password)
	return password, err
}

const (
	UpdateOneUserProfileSQL  = `UPDATE tb_user_basic SET name = ?, mobile = ?, gender = ? WHERE id = ?`
	UpdateOneUserAvatarSQL   = `UPDATE tb_user_basic SET avatar = ? WHERE id = ?`
	UpdateOneUserQrCodeSQL   = `UPDATE tb_user_basic SET qr_code = ? WHERE id = ?`
	UpdateOneUserPasswordSQL = `UPDATE tb_user_basic SET password = ? WHERE id = ?`
	UpdateOneUserIsDeleteSQL = `UPDATE tb_user_basic SET is_delete = ?  WHERE id = ?`
)

// Update 'name', 'mobile' and 'gender' columns of one row data which find by 'id' in 'tb_user_basic' table.
func UpdateOneUserProfileById(name, mobile string, gender int32, id int64) error {
	return insertOrUpdateOneRowData(UpdateOneUserProfileSQL, name, mobile, gender, id)
}

// Update the 'avatar' column of the one row data which find by 'id' in 'tb_user_basic' table.
func UpdateOneUserAvatarById(avatar string, id int64) error {
	return insertOrUpdateOneRowData(UpdateOneUserAvatarSQL, avatar, id)

}

// Update the 'qr_code' column of the one row data which find by 'id' in 'tb_user_basic' table.
func UpdateOneUserQrCodeById(qrCode string, id int64) error {
	return insertOrUpdateOneRowData(UpdateOneUserQrCodeSQL, qrCode, id)
}

// Update the 'password' column of the one row data which find by 'id' in 'tb_user_basic' table.
func UpdateOneUserPasswordById(password string, id int64) error {
	return insertOrUpdateOneRowData(UpdateOneUserPasswordSQL, password, id)

}

// Update the 'is_delete' column of the one row data which find by 'id' in 'tb_user_basic' table.
func UpdateOneUserIsDeleteById(isDelete bool, id int64) error {
	return insertOrUpdateOneRowData(UpdateOneUserIsDeleteSQL, isDelete, id)

}

const DeleteOneUserSQL = `DELETE FROM tb_user_basic WHERE id = ?`

// Delete one row data which find by 'id' in 'tb_user_basic' table really.
func DeleteOneUserById(id int64) error {
	return insertOrUpdateOneRowData(DeleteOneUserSQL, id)
}

var (
	ErrNotFriendYet              = errors.New("you are not friends yet")
	ErrInBlackList               = errors.New("you are in the blacklist of target user")
	ErrFriendshipAlreadyInEffect = errors.New("your friendship already in effect")
	ErrNotTheFriendRequest       = errors.New("not have the friend request you can accept")
)

// Friend basic information in `tb_user_basic` and `tb_friendship` table.
// Because protocol buffer 3 only have int32, so 'Gender' also use int32 here.
type TableFriendInfo struct {
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

const (
	CheckOneFriendIsBlackSQL  = `SELECT is_black FROM tb_friendship WHERE self_id = ? AND friend_id = ?`
	CheckOneFriendIsAcceptSQL = `SELECT is_accept FROM tb_friendship WHERE self_id = ? AND friend_id = ?`

	InsertOneNewFriendSQL = `INSERT INTO tb_friendship (self_id, friend_id, friend_note) VALUES (?, ?, ?) 
ON DUPLICATE KEY UPDATE friend_note = ?, is_accept = FALSE, is_black = FALSE, is_delete = FALSE`

	UpdateAcceptOneNewFriendSQL = `INSERT INTO tb_friendship (self_id, friend_id, is_accept) VALUES (?, ?, ?) 
ON DUPLICATE KEY UPDATE is_accept = TRUE, is_black = FALSE, is_delete = FALSE`

	UpdateOneFriendIsBlackSQL = `INSERT INTO tb_friendship (self_id, friend_id, is_black) VALUES (?, ?, ?) 
ON DUPLICATE KEY UPDATE is_black = ?`

	UpdateOneFriendNoteSQL = `UPDATE tb_friendship SET friend_note = ? WHERE self_id = ? and friend_id = ?`

	UpdateOneFriendIsDeleteSQL = `UPDATE tb_friendship SET is_accept = FALSE, is_black = FALSE, is_delete = TRUE
WHERE self_id=? AND friend_id = ?`
)

// Insert one new row data in 'tb_friendship' table.
// First, check the target user haven add the requester in blacklist. If yes, refuse the friend request by return
// an error. Then check the friend request record is already existed and the 'is_accept' is 'true', if all yes,
// tell the requester they are already have an effect friendship record by return an error
func InsertOneNewFriend(selfId, friendId int64, friendNote string) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}

	// check blacklist
	isBlack := false
	row := tx.QueryRow(CheckOneFriendIsBlackSQL, friendId, selfId)
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
	row = tx.QueryRow(CheckOneFriendIsAcceptSQL, selfId, friendId)
	err = row.Scan(&isAccept)
	if nil != err && sql.ErrNoRows != err {
		_ = tx.Rollback()
		return err
	}
	if isAccept {
		_ = tx.Rollback()
		return ErrFriendshipAlreadyInEffect
	}

	// insert the data finally
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

// Handle a friend request, the value of 'isAccept' can be true or false.
// First, check the friend request record data if existed and the 'is_accept' value, if the data not existed or the
// value of 'is_accept' is 'true', return the error.
// When the value of 'isAccept' is 'false', it will add the requester in blacklist of the current user by the way.
func UpdateAcceptOneNewFriend(selfId, friendId int64, friendNote string, isAccept bool) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}

	// check the friend request if existed and the 'is_accept value'
	anotherIsAccept := false
	row := tx.QueryRow(CheckOneFriendIsAcceptSQL, friendId, selfId)
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
		_, err = tx.Exec(UpdateAcceptOneNewFriendSQL, selfId, friendId, isAccept)
		if nil != err {
			_ = tx.Rollback()
			return err
		}
		// change the friend relationship record of requester, make the `is_accept` also be true
		_, err = tx.Exec(UpdateAcceptOneNewFriendSQL, friendId, selfId, isAccept)
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

// Update the 'friend_note' column of the one row data which find by 'selfId' and 'friendId' in 'tb_friendship' table.
func UpdateOneFriendNote(selfId, friendId int64, friendNote string) error {
	return insertOrUpdateOneRowData(UpdateOneFriendNoteSQL, friendNote, selfId, friendId)
}

// Update the 'is_black' column of the one row data which find by 'selfId' and 'friendId' in 'tb_friendship' table.
func UpdateOneFriendIsBlack(selfId, friendId int64, isBlack bool) error {
	return insertOrUpdateOneRowData(UpdateOneFriendIsBlackSQL, selfId, friendId, isBlack, isBlack)
}

// Update the 'is_delete' column of the one row data which find by 'selfId' and 'friendId' in 'tb_friendship' table.
// When the 'isDelete' is 'true', it will change the friend's data by the way. It meaning that will remove the current
// user from friend's friendship list at the same time.
func UpdateOneFriendIsDelete(selfId, friendId int64) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}

	// update the record of current user.
	affect, err := tx.Exec(UpdateOneFriendIsDeleteSQL, selfId, friendId)
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
	affect, err = tx.Exec(UpdateOneFriendIsDeleteSQL, friendId, selfId)
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
	SelectAllFriendsInfoSQL = `SELECT id, friend_note, email, name, mobile, gender, avatar, is_accept, is_black, 
friends.is_delete FROM tb_user_basic AS basic, (SELECT friend_id, friend_note, is_accept,is_black, is_delete FROM
tb_friendship WHERE self_id= ? AND is_accept=TRUE AND is_delete=FALSE) AS friends WHERE basic.id = friends.friend_id`

	SelectEffectFriendsInfoSQL = `SELECT id, friend_note, email, name, mobile, gender, avatar, is_accept, is_black,
friends.is_delete FROM tb_user_basic AS basic, (SELECT friend_id, friend_note, is_accept,is_black, is_delete FROM 
tb_friendship WHERE self_id= ? AND is_accept = TRUE AND is_black = FALSE AND is_delete = FALSE) AS friends WHERE 
basic.id = friends.friend_id`

	SelectBlacklistFriendsInfoSQL = `SELECT id, friend_note, email, name, mobile, gender, avatar, is_accept,is_black,
friends.is_delete FROM tb_user_basic AS basic, (SELECT friend_id, friend_note, is_accept,is_black, is_delete FROM 
tb_friendship WHERE self_id= ? AND is_black=TRUE AND is_delete = FALSE) AS friends WHERE basic.id = friends.friend_id`

	SelectEffectFriendsIdSQL = `SELECT friend_id FROM tb_friendship WHERE self_id = ? AND is_accept = TRUE 
AND is_black = FALSE  AND is_delete = FALSE`

	SelectBlacklistFriendsIdSQL = `SELECT friend_id FROM tb_friendship WHERE self_id = ? AND is_black = TRUE 
AND is_delete = FALSE`
)

// Get the information of user's friends whom selected by the sql string
func selectFriendsInfo(sqlStr string, selfId int64) ([]*TableFriendInfo, error) {
	rows, err := MySQLClient.Query(sqlStr, selfId)
	if nil != err {
		return nil, err
	}

	friendsInfo := make([]*TableFriendInfo, 0)
	for rows.Next() {
		temp := new(TableFriendInfo)
		err = rows.Scan(&(temp.Id), &(temp.Note), &(temp.Email), &(temp.Name), &(temp.Mobile),
			&(temp.Gender), &(temp.Avatar), &(temp.IsAccept), &(temp.IsBlack), &(temp.IsDelete))
		if nil != err {
			continue
		}
		friendsInfo = append(friendsInfo, temp)
	}
	return friendsInfo, nil

}

// Get the information of user's all friends, which 'is_accept' is 'true' and 'is_delete' is 'false' in
// 'tb_friendship' table, meaning it would include the friends who in user's blacklist.
func SelectAllFriendsInfo(selfId int64) ([]*TableFriendInfo, error) {
	return selectFriendsInfo(SelectAllFriendsInfoSQL, selfId)
}

// Get the information of user's effective friends, which 'is_accept' is 'true', 'is_black' and 'is_delete' is 'false'
// in 'tb_friendship' table, meaning it would not include the friends who in user's blacklist.
func SelectEffectFriendsInfo(selfId int64) ([]*TableFriendInfo, error) {
	return selectFriendsInfo(SelectEffectFriendsInfoSQL, selfId)
}

// Get the information of user's friends whom are in user's blacklist
func SelectBlacklistFriendsInfo(selfId int64) ([]*TableFriendInfo, error) {
	return selectFriendsInfo(SelectBlacklistFriendsInfoSQL, selfId)
}

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
func SelectEffectFriendsId(selfId int64) ([]int64, error) {
	return selectFiendsId(SelectEffectFriendsIdSQL, selfId)
}

// Get the id of user's friends whom are in user's blacklist
func SelectBlacklistFriendsId(selfId int64) ([]int64, error) {
	return selectFiendsId(SelectBlacklistFriendsIdSQL, selfId)
}

const DeleteOneFriendshipRecordSQL = `DELETE FROM tb_friendship WHERE self_id = ? AND friend_id = ?`

// Delete one row data which find by 'self_id' and 'friend_id' in 'tb_friendship' table really.
func DeleteOneFriendshipRecord(selfId, friendId int64) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}
	ret, err := tx.Exec(DeleteOneFriendshipRecordSQL, selfId, friendId)
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	affectCount, err := ret.RowsAffected()
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	if affectCount == 0 {
		_ = tx.Rollback()
		return ErrAffectZeroCount
	}

	err = tx.Commit()
	if nil != err {
		_ = tx.Rollback()
		return err
	}

	return nil
}
