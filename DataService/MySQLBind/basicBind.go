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
	// Get MySQL client
	MySQLClient, err = sql.Open("mysql", conf.MySQLURI)
	if nil != err {
		log.Fatal(err)
	}

	err = MySQLClient.Ping()
	if nil != err {
		log.Fatal(err)
	}

	// Get SnowFlakeNode to generate ID
	SnowFlakeNode, err = snowflake.NewNode(0)
	if nil != err {
		log.Fatal(err)
	}

}

var (
	ErrAffectZeroCount           = errors.New("0 row affected")
	ErrUserNotFound              = errors.New("the user not found")
	ErrFriendNotFound            = errors.New("the friend not found")
	ErrNotFriendYet              = errors.New("you are not friends yet")
	ErrInBlackList               = errors.New("you are in the blacklist of target user")
	ErrFriendshipAlreadyInEffect = errors.New("your friendship already in effect")
	ErrNotTheFriendRequest       = errors.New("not have the friend request you can accept")
	ErrGroupChatNotFound         = errors.New("the group chat not found")
	ErrSubscriptionNotFound      = errors.New("the subscription not found")
)

// Private Function:
// Operate one row data by given sql string and values with open transaction.
func execSqlWithTransaction(sqlStr string, args ...interface{}) error {
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

// -----------------------------------------------------------------------

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

// Private Function:
// Scan one user's information from the 'row'.
func scanUserFromRow(row *sql.Row) (*TableUserBasic, error) {
	user := new(TableUserBasic)
	err := row.Scan(&(user.Id), &(user.Email), &(user.Name), &(user.Password), &(user.Mobile),
		&(user.Gender), &(user.Avatar), &(user.QrCode), &(user.IsDelete))
	if sql.ErrNoRows == err {
		return nil, ErrUserNotFound
	}
	if nil != err {
		return nil, err
	}
	return user, nil
}

const (
	InsertOneNewUserSQL = `INSERT INTO tb_user_basic (id, email, name, password, mobile, gender, avatar, qr_code, is_delete) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	DeleteOneUserRealSQL = `DELETE FROM tb_user_basic WHERE id = ?`
)

// Insert one row new data for saving information of new user.
// The 'id' will auto generate by 'SnowFlakeNode', the 'is_delete' will use default value 'false'.
func InsertOneNewUser(email, name, password, mobile string,
	gender int32, avatar, qrCode string, isDelete bool) (*TableUserBasic, error) {

	// generate an ID and insert the data
	id := SnowFlakeNode.Generate().Int64()
	err := execSqlWithTransaction(InsertOneNewUserSQL,
		id, email, name, password, mobile, gender, avatar, qrCode, isDelete)
	if nil != err {
		return nil, err
	}

	// return an user's basic information to follow the REST style
	user := &TableUserBasic{Id: id, Email: email, Name: name, Password: password,
		Mobile: mobile, Gender: gender, Avatar: avatar, QrCode: qrCode}

	return user, nil
}

// Delete one row data which find by 'id' in 'tb_user_basic' table really.
func DeleteOneUserByIdReal(id int64) error {
	return execSqlWithTransaction(DeleteOneUserRealSQL, id)
}

const (
	SelectUserBaseSQL = `SELECT id, email, name, password,mobile, gender, avatar, qr_code, is_delete FROM tb_user_basic`

	SelectOneUserByIdSQL            = SelectUserBaseSQL + ` WHERE id = ? AND is_delete = ?`
	SelectOneUserByEmailSQL         = SelectUserBaseSQL + ` WHERE email = ? AND is_delete = ?`
	SelectManyUserByNameSQL         = SelectUserBaseSQL + ` WHERE name = ? AND is_delete = ?`
	SelectOneUserPasswordByIdSQL    = `SELECT password FROM tb_user_basic WHERE id = ?`
	SelectOneUserPasswordByEmailSQL = `SELECT password FROM tb_user_basic WHERE email = ?`
)

// Select one row data from 'tb_user_basic' table by 'id' and 'is_delete'.
func SelectOneUserById(id int64, isDelete bool) (*TableUserBasic, error) {
	row := MySQLClient.QueryRow(SelectOneUserByIdSQL, id, isDelete)
	return scanUserFromRow(row)
}

// Select one row data from 'tb_user_basic' table by 'email' and 'is_delete'.
func SelectOneUserByEmail(email string, isDelete bool) (*TableUserBasic, error) {
	row := MySQLClient.QueryRow(SelectOneUserByEmailSQL, email, isDelete)
	return scanUserFromRow(row)

}

// Select many rows data from 'tb_user_basic' table by 'name' and 'is_delete'.
func SelectManyUserByName(name string, isDelete bool) ([]*TableUserBasic, error) {
	rows, err := MySQLClient.Query(SelectManyUserByNameSQL, name, isDelete)
	if nil != err {
		return nil, err
	}

	users := make([]*TableUserBasic, 0)
	for rows.Next() {
		user := new(TableUserBasic)
		err := rows.Scan(&(user.Id), &(user.Email), &(user.Name), &(user.Password), &(user.Mobile),
			&(user.Gender), &(user.Avatar), &(user.QrCode), &(user.IsDelete))
		if nil != err {
			continue
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
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

// Select all rows data from 'tb_user_basic' table.
func SelectAllUsers() ([]*TableUserBasic, error) {
	rows, err := MySQLClient.Query(SelectUserBaseSQL)
	if nil != err {
		return nil, err
	}
	result := make([]*TableUserBasic, 0)
	for rows.Next() {
		user := new(TableUserBasic)
		err := rows.Scan(&(user.Id), &(user.Email), &(user.Name), &(user.Password), &(user.Mobile),
			&(user.Gender), &(user.Avatar), &(user.QrCode), &(user.IsDelete))
		if nil != err {
			continue
		}
		result = append(result, user)
	}
	if len(result) == 0 {
		return nil, ErrUserNotFound
	}
	return result, nil
}

const (
	UpdateOneUserProfileSQL  = `UPDATE tb_user_basic SET name = ?, mobile = ?, gender = ? WHERE id = ?`
	UpdateOneUserAvatarSQL   = `UPDATE tb_user_basic SET avatar = ? WHERE id = ?`
	UpdateOneUserQrCodeSQL   = `UPDATE tb_user_basic SET qr_code = ? WHERE id = ?`
	UpdateOneUserPasswordSQL = `UPDATE tb_user_basic SET password = ? WHERE id = ?`
	UpdateOneUserIsDeleteSQL = `UPDATE tb_user_basic SET is_delete = ?  WHERE id = ?`
)

// Update the 'avatar' column of the one row data which find by 'id' in 'tb_user_basic' table.
func UpdateOneUserAvatarById(avatar string, id int64) error {
	return execSqlWithTransaction(UpdateOneUserAvatarSQL, avatar, id)

}

// Update the 'qr_code' column of the one row data which find by 'id' in 'tb_user_basic' table.
func UpdateOneUserQrCodeById(qrCode string, id int64) error {
	return execSqlWithTransaction(UpdateOneUserQrCodeSQL, qrCode, id)
}

// Update the 'password' column of the one row data which find by 'id' in 'tb_user_basic' table.
func UpdateOneUserPasswordById(password string, id int64) error {
	return execSqlWithTransaction(UpdateOneUserPasswordSQL, password, id)

}

// Update the 'is_delete' column of the one row data which find by 'id' in 'tb_user_basic' table.
func UpdateOneUserIsDeleteById(isDelete bool, id int64) error {
	return execSqlWithTransaction(UpdateOneUserIsDeleteSQL, isDelete, id)

}

// -----------------------------------------------------------------------

// friendship record information in 'tb_friendship' table
type TableFriendship struct {
	SelfId     int64  `json:"self_id"`
	FriendId   int64  `json:"friend_id"`
	FriendNote string `json:"friend_note"`
	IsAccept   bool   `json:"is_accept"`
	IsBlack    bool   `json:"is_black"`
	IsDelete   bool   `json:"is_delete"`
}

const (
	InsertOneNewFriendSQL = `INSERT INTO tb_friendship (self_id, friend_id, friend_note) VALUES (?, ?, ?) 
ON DUPLICATE KEY UPDATE friend_note = ?, is_accept = FALSE, is_black = FALSE, is_delete = FALSE`

	DeleteOneFriendshipRealSQL = `DELETE FROM tb_friendship WHERE self_id = ? AND friend_id = ?`
)

// Insert one new row data in 'tb_friendship' table.
// if the record was existed, update it.
func InsertOneNewFriend(selfId, friendId int64, friendNote string) error {
	return execSqlWithTransaction(InsertOneNewFriendSQL, selfId, friendId, friendNote, friendNote)
}

// Delete one row data which find by 'self_id' and 'friend_id' in 'tb_friendship' table really.
func DeleteOneFriendReal(selfId, friendId int64) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}
	ret, err := tx.Exec(DeleteOneFriendshipRealSQL, selfId, friendId)
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

const (
	SelectFriendshipBaseSQL = `SELECT self_id, friend_id, friend_note, is_accept, is_black, is_delete FROM
tb_friendship`

	SelectOneFriendshipSQL = SelectFriendshipBaseSQL + ` WHERE self_id= ? AND friend_id= ?`

	SelectFriendsIdByOptionsSQL = `SELECT friend_id FROM tb_friendship WHERE self_id= ? AND is_accept= ? AND 
is_black= ? AND is_delete= ?`
)

// Select one row data from 'tb_friendship' table by 'self_id' and 'friend_id'.
func SelectOneFriendship(selfId, friendId int64) (*TableFriendship, error) {
	row := MySQLClient.QueryRow(SelectOneFriendshipSQL, selfId, friendId)
	temp := new(TableFriendship)
	err := row.Scan(&(temp.SelfId), &(temp.FriendId), &(temp.FriendNote),
		&(temp.IsAccept), &(temp.IsBlack), &(temp.IsDelete))

	if sql.ErrNoRows == err {
		return nil, ErrFriendNotFound
	}
	if nil != err {
		return nil, err
	}
	return temp, nil
}

// Select the values of 'friend_id' column from 'tb_friendship' table by 'self_id', 'is_accept', 'is_black'
// and 'is_delete'
func SelectFriendsIdByOptions(selfId int64, isAccept, isBlack, isDelete bool) ([]int64, error) {
	rows, err := MySQLClient.Query(SelectFriendsIdByOptionsSQL, selfId, isAccept, isBlack, isDelete)
	if nil != err {
		return nil, err
	}
	result := make([]int64, 0)
	for rows.Next() {
		tempId := new(int64)
		err := rows.Scan(tempId)
		if nil != err {
			continue
		}
		result = append(result, *tempId)
	}
	if len(result) == 0 {
		return nil, ErrFriendNotFound
	}
	return result, nil
}

// Select all rows data from 'tb_friendship' table
func SelectAllFriendship() ([]*TableFriendship, error) {
	rows, err := MySQLClient.Query(SelectFriendshipBaseSQL)
	if nil != err {
		return nil, err
	}
	result := make([]*TableFriendship, 0)
	for rows.Next() {
		temp := new(TableFriendship)
		err := rows.Scan(&(temp.SelfId), &(temp.FriendId), &(temp.FriendNote),
			&(temp.IsAccept), &(temp.IsBlack), &(temp.IsDelete))
		if nil != err {
			continue
		}
		result = append(result, temp)
	}

	if len(result) == 0 {
		return nil, ErrFriendNotFound
	}
	return result, nil
}

const (
	UpdateOneFriendNoteSQL = `UPDATE tb_friendship SET friend_note= ? WHERE self_id= ? AND friend_id= ?`

	UpdateOneFriendIsAcceptSQL = `INSERT INTO tb_friendship (self_id, friend_id, is_accept) VALUES (?, ?, ?) 
ON DUPLICATE KEY UPDATE is_accept= ?`

	UpdateOneFriendIsBlackSQL = `INSERT INTO tb_friendship (self_id, friend_id, is_black) VALUES (?, ?, ?) 
ON DUPLICATE KEY UPDATE is_black= ?`

	UpdateOneFriendIsDeleteSQL = `UPDATE tb_friendship SET is_delete= ? WHERE self_id=? AND friend_id= ?`
)

// Update the 'friend_note' column of the one row data which find by 'self_id' and 'friend_id' in 'tb_friendship' table.
func UpdateOneFriendNote(selfId, friendId int64, friendNote string) error {
	return execSqlWithTransaction(UpdateOneFriendNoteSQL, friendNote, selfId, friendId)
}

// Update the 'is_accept' column of the one row data which find by 'self_id' and 'friend_id' in 'tb_friendship' table.
func UpdateOneFriendIsAccept(selfId, friendId int64, isAccept bool) error {
	return execSqlWithTransaction(UpdateOneFriendIsAcceptSQL, selfId, friendId, isAccept, isAccept)
}

// Update the 'is_black' column of the one row data which find by 'self_id' and 'friend_id' in 'tb_friendship' table.
func UpdateOneFriendIsBlack(selfId, friendId int64, isBlack bool) error {
	return execSqlWithTransaction(UpdateOneFriendIsBlackSQL, selfId, friendId, isBlack, isBlack)
}

// Update the 'is_delete' column of the one row data which find by 'self_id' and 'friend_id' in 'tb_friendship' table.
func UpdateOneFriendIsDelete(selfId, friendId int64, isDelete bool) error {
	return execSqlWithTransaction(UpdateOneFriendIsDeleteSQL, isDelete, selfId, friendId)
}

// -----------------------------------------------------------------------

// Group chat information in 'tb_group_chat'
type TableGroupChat struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	ManagerId int64  `json:"manager_id"`
	Avatar    string `json:"avatar"`
	QrCode    string `json:"qr_code"`
	IsDelete  bool   `json:"is_delete"`
}

const (
	InsertOneNewGroupChatSQL  = `INSERT INTO tb_group_chat(id, name, manager_id, avatar, qr_code) VALUES (?,?,?,?,?)`
	DeleteOneGroupChatRealSQL = `DELETE FROM tb_group_chat WHERE id = ?`
)

// Insert one new row data to save the information of an new group chat
// 'id' will generated by SnowFlakeNode, 'is_delete' use default value 'false'.
func InsertOneNewGroupChat(name, avatar, qrCode string, managerId int64) (*TableGroupChat, error) {

	// generate the ID and insert the data
	id := SnowFlakeNode.Generate().Int64()
	err := execSqlWithTransaction(InsertOneNewGroupChatSQL, id, name, managerId, avatar, qrCode)
	if nil != err {
		return nil, err
	}
	// return the inserted data for following the REST style
	groupChat := &TableGroupChat{Id: id, Name: name, ManagerId: managerId, Avatar: avatar, QrCode: qrCode}
	return groupChat, nil
}

// Delete one row data from 'tb_group_chat' by 'id' really
func DeleteOneGroupChatByIdReal(id int64) error {
	err := execSqlWithTransaction(DeleteOneGroupChatRealSQL, id)
	if nil != err {
		return err
	}
	return nil
}

const (
	SelectGroupChatBaseSQL         = `SELECT id, name, manager_id, avatar, qr_code, is_delete FROM tb_group_chat`
	SelectGroupChatByIdSQL         = SelectGroupChatBaseSQL + ` WHERE id= ? AND is_delete= ?`
	SelectGroupChatsByNameSQL      = SelectGroupChatBaseSQL + ` WHERE name= ? AND is_delete= ?`
	SelectGroupChatsByManagerIdSQl = SelectGroupChatBaseSQL + ` WHERE manager_id=? AND is_delete= ?`
)

// Private Function:
// Scan information of many group chat from the 'rows'.
// If the result is an empty slice, will return the 'ErrGroupChatNotFound'
func scanGroupChatFromRows(rows *sql.Rows) ([]*TableGroupChat, error) {
	groupChats := make([]*TableGroupChat, 0)
	for rows.Next() {
		temp := new(TableGroupChat)
		err := rows.Scan(&(temp.Id), &(temp.Name), &(temp.ManagerId),
			&(temp.Avatar), &(temp.QrCode), &(temp.IsDelete))
		if nil != err {
			continue
		}
		groupChats = append(groupChats, temp)
	}
	if len(groupChats) == 0 {
		return nil, ErrGroupChatNotFound
	}
	return groupChats, nil
}

// Select one row data from 'tb_group_chat' table by 'id'.
func SelectOneGroupChatById(id int64, isDelete bool) (*TableGroupChat, error) {
	row := MySQLClient.QueryRow(SelectGroupChatByIdSQL, id, isDelete)
	temp := new(TableGroupChat)
	err := row.Scan(&(temp.Id), &(temp.Name), &(temp.ManagerId),
		&(temp.Avatar), &(temp.QrCode), &(temp.IsDelete))

	if sql.ErrNoRows == err {
		return nil, ErrGroupChatNotFound
	}

	if nil != err {
		return nil, err
	}
	return temp, nil
}

// Select many rows data from 'tb_group_chat' table by 'name'.
func SelectManyGroupChatByName(name string, isDelete bool) ([]*TableGroupChat, error) {
	rows, err := MySQLClient.Query(SelectGroupChatsByNameSQL, name, isDelete)
	if nil != err {
		return nil, err
	}
	return scanGroupChatFromRows(rows)
}

// Select many rows data from 'tb_group_chat' table by 'manager_id'.
func SelectManyGroupChatByManagerId(managerId int64, isDelete bool) ([]*TableGroupChat, error) {
	rows, err := MySQLClient.Query(SelectGroupChatsByManagerIdSQl, managerId, isDelete)
	if nil != err {
		return nil, err
	}
	return scanGroupChatFromRows(rows)
}

// Select all rows data from 'tb_group_chat' table
func SelectAllGroupChat() ([]*TableGroupChat, error) {
	rows, err := MySQLClient.Query(SelectGroupChatBaseSQL)
	if nil != err {
		return nil, err
	}
	return scanGroupChatFromRows(rows)
}

const (
	UpdateOneGroupChatNameByIdSQL     = `UPDATE tb_group_chat SET name = ? WHERE id = ?`
	UpdateOneGroupChatManagerByIdSQL  = `UPDATE tb_group_chat SET manager_id = ? WHERE id = ?`
	UpdateOneGroupChatAvatarByIdSQL   = `UPDATE tb_group_chat SET avatar = ? WHERE id = ?`
	UpdateOneGroupChatQrCodeByIdSQL   = `UPDATE tb_group_chat SET qr_code = ? WHERE id = ?`
	UpdateOneGroupChatIsDeleteByIdSQL = `UPDATE tb_group_chat SET is_delete = ? WHERE id = ?`
)

// Update the 'name' column of the one row data which find by 'id' in 'tb_group_chat' table
func UpdateOneGroupChatNameById(id int64, newName string) error {
	return execSqlWithTransaction(UpdateOneGroupChatNameByIdSQL, newName, id)
}

// Update the 'manager_id' column of the one row data which find by 'id' in 'tb_group_chat' table
func UpdateOneGroupChatManagerById(id, newManagerId int64) error {
	return execSqlWithTransaction(UpdateOneGroupChatManagerByIdSQL, newManagerId, id)
}

// Update the 'avatar' column of the one row data which find by 'id' in 'tb_group_chat' table
func UpdateOneGroupChatAvatarById(id int64, newAvatar string) error {
	return execSqlWithTransaction(UpdateOneGroupChatAvatarByIdSQL, newAvatar, id)
}

// Update the 'qr_code' column of the one row data which find by 'id' in 'tb_group_chat' table
func UpdateOneGroupChatQrCodeById(id int64, newQrCode string) error {
	return execSqlWithTransaction(UpdateOneGroupChatQrCodeByIdSQL, newQrCode, id)
}

// Update the 'is_delete' column of the one row data which find by 'id' in 'tb_group_chat' table
func UpdateOneGroupChatIsDeleteById(id int64, isDelete bool) error {
	return execSqlWithTransaction(UpdateOneGroupChatIsDeleteByIdSQL, isDelete, id)
}

// -----------------------------------------------------------------------

// User and group chat information in 'tb_user_group_chat'
type TableUserGroupChat struct {
	GroupId  int64  `json:"group_id"`
	UserId   int64  `json:"user_id"`
	UserNote string `json:"user_note"`
	IsDelete bool   `json:"is_delete"`
}

const (
	InsertOneNewUserGroupChatSQL = `INSERT INTO tb_user_group_chat (group_id, user_id, user_note) VALUES (?,?,?) 
ON DUPLICATE KEY UPDATE user_note= ?, is_delete=FALSE `
	DeleteOneUserGroupChatRealSQL = `DELETE FROM tb_user_group_chat WHERE group_id= ? AND user_id= ?`
)

// Insert one new row data into 'tb_user_group_chat' table.
// If it is duplicated entry, will update the 'user_note' and 'is_delete'.
func InsertOneNewUserGroupChat(groupId, userId int64, userNote string) (*TableUserGroupChat, error) {
	err := execSqlWithTransaction(InsertOneNewUserGroupChatSQL, groupId, userId, userNote, userNote)
	if nil != err {
		return nil, err
	}
	temp := &TableUserGroupChat{GroupId: groupId, UserId: userId, UserNote: userNote}
	return temp, nil
}

// Delete one row data from 'tb_user_group_chat' table really
func DeleteOneUserGroupChatReal(groupId, userId int64) error {
	return execSqlWithTransaction(DeleteOneUserGroupChatRealSQL, groupId, userId)
}

const (
	SelectUserGroupChatBaseSQL = `SELECT group_id, user_id, user_note, is_delete FROM tb_user_group_chat`

	SelectOneUserGroupChatByGIDAndUIDSQL = SelectUserGroupChatBaseSQL + ` WHERE group_id= ? AND user_id= ?`
	SelectManyUserGroupChatByGroupIdSQL  = SelectUserGroupChatBaseSQL + ` WHERE group_id= ? AND is_delete= ?`
	SelectManyUserGroupChatByUserIdSQL   = SelectUserGroupChatBaseSQL + ` WHERE user_id= ? AND is_delete= ?`

	SelectUsersIdOfGroupChatSQL = `SELECT user_id FROM tb_user_group_chat WHERE group_id= ? AND is_delete= ?`

	SelectGroupChatsIdOfUserSQL = `SELECT group_id FROM tb_user_group_chat WHERE user_id= ? AND is_delete= ?`
)

// Select one row data from 'tb_user_group_chat' table by 'group_id' and 'user_id'
func SelectOneUserGroupChat(groupId, userId int64) (*TableUserGroupChat, error) {
	row := MySQLClient.QueryRow(SelectOneUserGroupChatByGIDAndUIDSQL, groupId, userId)
	temp := new(TableUserGroupChat)
	err := row.Scan(&(temp.GroupId), &(temp.UserId), &(temp.UserNote), &(temp.IsDelete))
	if nil != err {
		return nil, err
	}
	return temp, nil
}

// Scan information of many rows select from 'tb_user_group_chat' table
func scanUserGroupChatRows(rows *sql.Rows) []*TableUserGroupChat {
	result := make([]*TableUserGroupChat, 0)
	for rows.Next() {
		temp := new(TableUserGroupChat)
		err := rows.Scan(&(temp.GroupId), &(temp.UserId), &(temp.UserNote), &(temp.IsDelete))
		if nil != err {
			continue
		}
		result = append(result, temp)
	}
	return result
}

// Select all rows data from 'tb_user_group_chat' table.
func SelectAllUserGroupChat() ([]*TableUserGroupChat, error) {
	rows, err := MySQLClient.Query(SelectUserGroupChatBaseSQL)
	if nil != err {
		return nil, err
	}
	return scanUserGroupChatRows(rows), nil
}

// Select many rows data from 'tb_user_group_chat' table by 'group_id' and 'is_delete'.
func SelectManyUserGroupChatByGroupId(groupId int64, isDelete bool) ([]*TableUserGroupChat, error) {
	rows, err := MySQLClient.Query(SelectManyUserGroupChatByGroupIdSQL, groupId, isDelete)
	if nil != err {
		return nil, err
	}
	return scanUserGroupChatRows(rows), nil
}

// Select many rows data from 'tb_user_group_chat' table by 'user_id' and 'is_delete'.
func SelectManyUserGroupChatByUserId(userId int64, isDelete bool) ([]*TableUserGroupChat, error) {
	rows, err := MySQLClient.Query(SelectManyUserGroupChatByUserIdSQL, userId, isDelete)
	if nil != err {
		return nil, err
	}
	return scanUserGroupChatRows(rows), nil
}

// Select the values of 'user_id' column from 'tb_user_group_chat' table by 'group_id' and 'is_delete'
func SelectUsersIdOfGroupChat(groupId int64, isDelete bool) ([]int64, error) {
	rows, err := MySQLClient.Query(SelectUsersIdOfGroupChatSQL, groupId, isDelete)
	if nil != err {
		return nil, err
	}
	result := make([]int64, 0)
	for rows.Next() {
		temp := new(int64)
		err := rows.Scan(temp)
		if nil != err {
			continue
		}
		result = append(result, *temp)
	}
	return result, nil
}

// Select the values of 'group_id' column from 'tb_user_group_chat' table by 'user_id' and 'is_delete'.
func SelectGroupChatsIdOfUser(userId int64, isDelete bool) ([]int64, error) {
	rows, err := MySQLClient.Query(SelectGroupChatsIdOfUserSQL, userId, isDelete)
	if nil != err {
		return nil, err
	}
	result := make([]int64, 0)
	for rows.Next() {
		temp := new(int64)
		err := rows.Scan(temp)
		if nil != err {
			continue
		}
		result = append(result, *temp)
	}
	return result, nil
}

const (
	UpdateOneUserGroupChatNoteSQL     = `UPDATE tb_user_group_chat SET user_note= ? WHERE group_id= ? AND user_id= ?`
	UpdateOneUserGroupChatIsDeleteSQL = `UPDATE tb_user_group_chat SET is_delete= ? WHERE group_id= ? AND user_id= ?`
)

// Update the 'user_note' column of one row data which found by 'group_id' and 'user_id' in 'tb_user_group_chat' table.
func UpdateOneUserGroupChatNote(newNote string, groupId, userId int64) error {
	return execSqlWithTransaction(UpdateOneUserGroupChatNoteSQL, newNote, groupId, userId)
}

// Update the 'is_delete' column of one row data which found by 'group_id' and 'user_id' in 'tb_user_group_chat' table.
func UpdateOneUserGroupChatIsDelete(isDelete bool, groupId, userId int64) error {
	return execSqlWithTransaction(UpdateOneUserGroupChatIsDeleteSQL, isDelete, groupId, userId)
}

// -----------------------------------------------------------------------

// Subscription information in 'tb_subscription' table
type TableSubscription struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	ManagerId int64  `json:"manager_id"`
	Introduce string `json:"introduce"`
	Avatar    string `json:"avatar"`
	QrCode    string `json:"qr_code"`
	IsDelete  bool   `json:"is_delete"`
}

const (
	InsertOneNewSubscriptionSQL = `INSERT INTO tb_subscription(id, name, manager_id, introduce, avatar, qr_code) 
VALUES(?, ?, ?, ?, ?, ?)`

	DeleteOneSubscriptionRealSQL = `DELETE FROM tb_subscription WHERE id= ?`
)

// Insert one row data into 'tb_subscription' table to save the information of one new subscription.
// 'id' will generated by SnowFlakeNode, 'is_delete' use default value 'false'.
func InsertOneNewSubscription(name, introduce, avatar, qrCode string, managerId int64) (*TableSubscription, error) {
	id := SnowFlakeNode.Generate().Int64()
	err := execSqlWithTransaction(InsertOneNewSubscriptionSQL, id, name, managerId, introduce, avatar, qrCode)
	if nil != err {
		return nil, err
	}
	temp := &TableSubscription{Id: id, Name: name, ManagerId: managerId,
		Introduce: introduce, Avatar: avatar, QrCode: qrCode}
	return temp, nil
}

// Delete one row data from 'tb_subscription' table by 'id' really
func DeleteOneSubscriptionReal(id int64) error {
	return execSqlWithTransaction(DeleteOneSubscriptionRealSQL, id)
}

const (
	SelectSubscriptionBaseSQL = `SELECT id,name,manager_id,introduce,avatar,qr_code,is_delete FROM tb_subscription`

	SelectOneSubscriptionByIdSQL         = SelectSubscriptionBaseSQL + ` WHERE id= ? AND is_delete= ?`
	SelectOneSubscriptionByNameSQL       = SelectSubscriptionBaseSQL + ` WHERE name= ? AND is_delete= ?`
	SelectManySubscriptionByManagerIdSQL = SelectSubscriptionBaseSQL + ` WHERE manager_id= ? AND is_delete= ?`
)

// Select one row data from 'tb_subscription' table by 'id'.
func SelectOneSubscriptionById(id int64, isDelete bool) (*TableSubscription, error) {
	row := MySQLClient.QueryRow(SelectOneSubscriptionByIdSQL, id, isDelete)
	temp := new(TableSubscription)
	err := row.Scan(&(temp.Id), &(temp.Name), &(temp.ManagerId), &(temp.Introduce),
		&(temp.Avatar), &(temp.QrCode), &(temp.IsDelete))
	if sql.ErrNoRows == err {
		return nil, ErrSubscriptionNotFound
	}
	if nil != err {
		return nil, err
	}
	return temp, nil
}

// Select one row data from 'tb_subscription' table by 'name', the 'name' is unique key.
func SelectOneSubscriptionByName(name string, isDelete bool) (*TableSubscription, error) {
	row := MySQLClient.QueryRow(SelectOneSubscriptionByNameSQL, name, isDelete)
	temp := new(TableSubscription)
	err := row.Scan(&(temp.Id), &(temp.Name), &(temp.ManagerId), &(temp.Introduce),
		&(temp.Avatar), &(temp.QrCode), &(temp.IsDelete))
	if sql.ErrNoRows == err {
		return nil, ErrSubscriptionNotFound
	}
	if nil != err {
		return nil, err
	}
	return temp, nil
}

// Select many rows data from 'tb_subscription' table by 'manager_id'.
func SelectManySubscriptionByManagerId(managerId int64, isDelete bool) ([]*TableSubscription, error) {
	rows, err := MySQLClient.Query(SelectManySubscriptionByManagerIdSQL, managerId, isDelete)
	if nil != err {
		return nil, err
	}
	result := make([]*TableSubscription, 0)
	for rows.Next() {
		temp := new(TableSubscription)
		err := rows.Scan(&(temp.Id), &(temp.Name), &(temp.ManagerId), &(temp.Introduce),
			&(temp.Avatar), &(temp.QrCode), &(temp.IsDelete))
		if nil != err {
			continue
		}
		result = append(result, temp)
	}
	if len(result) == 0 {
		return nil, ErrSubscriptionNotFound
	}
	return result, nil
}

const (
	UpdateOneSubscriptionNameByIdSQL      = `UPDATE tb_subscription SET name= ? WHERE id=?`
	UpdateOneSubscriptionManagerByIdSQL   = `UPDATE tb_subscription SET manager_id= ? WHERE id=?`
	UpdateOneSubscriptionIntroduceByIdSQL = `UPDATE tb_subscription SET introduce= ? WHERE id=?`
	UpdateOneSubscriptionAvatarByIdSQL    = `UPDATE tb_subscription SET avatar= ? WHERE id= ?`
	UpdateOneSubscriptionQrCodeByIdSQL    = `UPDATE tb_subscription SET qr_code= ? WHERE id= ?`
	UpdateOneSubscriptionIsDeleteByIdSQL  = `UPDATE tb_subscription SET is_delete= ? WHERE id= ?`
)

// Update the 'name' column of the one row data which find by 'id' in 'tb_subscription' table.
func UpdateOneSubscriptionNameById(newName string, id int64) error {
	return execSqlWithTransaction(UpdateOneSubscriptionNameByIdSQL, newName, id)
}

// Update the 'manager_id' column of the one row data which find by 'id' in 'tb_subscription' table.
func UpdateOneSubscriptionManagerById(newManagerId, id int64) error {
	return execSqlWithTransaction(UpdateOneSubscriptionManagerByIdSQL, newManagerId, id)
}

// Update the 'introduce' column of the one row data which find by 'id' in 'tb_subscription' table.
func UpdateOneSubscriptionIntroduceById(newIntroduce string, id int64) error {
	return execSqlWithTransaction(UpdateOneSubscriptionIntroduceByIdSQL, newIntroduce, id)
}

// Update the 'avatar' column of the one row data which find by 'id' in 'tb_subscription' table.
func UpdateOneSubscriptionAvatarById(newAvatar string, id int64) error {
	return execSqlWithTransaction(UpdateOneSubscriptionAvatarByIdSQL, newAvatar, id)
}

// Update the 'qr_code' column of the one row data which find by 'id' in 'tb_subscription' table.
func UpdateOneSubscriptionQrCodeById(newQrCode string, id int64) error {
	return execSqlWithTransaction(UpdateOneSubscriptionQrCodeByIdSQL, newQrCode, id)
}

// Update the 'is_delete' column of the one row data which find by 'id' in 'tb_subscription' table.
func UpdateOneSubscriptionIsDeleteById(isDelete bool, id int64) error {
	return execSqlWithTransaction(UpdateOneSubscriptionIsDeleteByIdSQL, isDelete, id)
}

// -----------------------------------------------------------------------

// User and subscription information in 'tb_user_subscription' table.
type TableUserSubscription struct {
	SubsId   int64 `json:"subs_id"`
	UserId   int64 `json:"user_id"`
	IsDelete int64 `json:"is_delete"`
}

const (
	InsertOneNewUserSubscriptionSQL = `INSERT INTO tb_user_subscription(subs_id, user_id) VALUES(?, ?)`

	DeleteOneUserSubscriptionRealSQL = `DELETE FROM tb_user_subscription WHERE subs_id= ? AND user_id= ?`
)

// Insert one new data into 'tb_user_subscription' table.
func InsertOneNewUserSubscription(subsId, userId int64) (*TableUserSubscription, error) {
	err := execSqlWithTransaction(InsertOneNewUserSubscriptionSQL, subsId, userId)
	if nil != err {
		return nil, err
	}
	return &TableUserSubscription{SubsId: subsId, UserId: userId}, nil
}

// Delete one row data from 'tb_user_subscription' table really.
func DeleteOneUserSubscriptionReal(subsId, userId int64) error {
	return execSqlWithTransaction(DeleteOneUserSubscriptionRealSQL, subsId, userId)
}

const (
	SelectUserSubscriptionBaseSQL = `SELECT subs_id, user_id, is_delete FROM tb_user_subscription`

	SelectOneUserSubscriptionByIdsSQL = SelectUserSubscriptionBaseSQL + ` WHERE subs_id= ? AND user_id= ? 
AND is_delete= ?`

	SelectManyUserSubscriptionBySubsIdSQL = SelectUserSubscriptionBaseSQL + ` WHERE subs_id= ? AND is_delete= ?`
	SelectManyUserSubscriptionByUserIdSQL = SelectUserSubscriptionBaseSQL + ` WHERE user_id= ? AND is_delete= ?`

	SelectUsersIdOfSubscriptionSQL = `SELECT user_id FROM tb_user_subscription WHERE subs_id= ? AND is_delete= ?`

	SelectSubscriptionsIdOfUserSQL = `SELECT subs_id FROM tb_user_subscription WHERE user_id= ? AND is_delete= ?`
)

// Select one row data from 'tb_user_subscription' table by 'subs_id' and 'user_id'
func SelectOneUserSubscriptionByIds(subsId, userId int64, isDelete bool) (*TableUserSubscription, error) {
	row := MySQLClient.QueryRow(SelectOneUserSubscriptionByIdsSQL, subsId, userId, isDelete)
	temp := new(TableUserSubscription)
	err := row.Scan(&(temp.SubsId), &(temp.UserId), &(temp.IsDelete))
	if nil != err {
		return nil, err
	}
	return temp, nil
}

// Select many rows data from 'tb_user_subscription' table by 'subs_id'.
func SelectManyUserSubscriptionBySubsId(subsId int64, isDelete bool) ([]*TableUserSubscription, error) {
	rows, err := MySQLClient.Query(SelectManyUserSubscriptionBySubsIdSQL, subsId, isDelete)
	if nil != err {
		return nil, err
	}
	result := make([]*TableUserSubscription, 0)
	for rows.Next() {
		temp := new(TableUserSubscription)
		err := rows.Scan(&(temp.SubsId), &(temp.UserId), &(temp.IsDelete))
		if nil != err {
			continue
		}
		result = append(result, temp)
	}
	return result, nil
}

// Select many rows data from 'tb_user_subscription' table by 'user_id'.
func SelectManyUserSubscriptionByUserId(userId int64, isDelete bool) ([]*TableUserSubscription, error) {
	rows, err := MySQLClient.Query(SelectManyUserSubscriptionByUserIdSQL, userId, isDelete)
	if nil != err {
		return nil, err
	}
	result := make([]*TableUserSubscription, 0)
	for rows.Next() {
		temp := new(TableUserSubscription)
		err := rows.Scan(&(temp.SubsId), &(temp.UserId), &(temp.IsDelete))
		if nil != err {
			continue
		}
		result = append(result, temp)
	}
	return result, nil
}

// Select the values of the 'user_id' column in 'tb_user_subscription' table by 'subs_id'.
func SelectUsersIdOfSubscription(subsId int64, isDelete bool) ([]int64, error) {
	rows, err := MySQLClient.Query(SelectUsersIdOfSubscriptionSQL, subsId, isDelete)
	if nil != err {
		return nil, err
	}
	result := make([]int64, 0)
	for rows.Next() {
		userId := new(int64)
		err := rows.Scan(userId)
		if nil != err {
			continue
		}
		result = append(result, *userId)
	}
	return result, nil
}

// Select the values of the 'subs_id' column in 'tb_user_subscription' table by 'user_id'.
func SelectSubscriptionsIdOfUser(userId int64, isDelete bool) ([]int64, error) {
	rows, err := MySQLClient.Query(SelectSubscriptionsIdOfUserSQL, userId, isDelete)
	if nil != err {
		return nil, err
	}
	result := make([]int64, 0)
	for rows.Next() {
		subsId := new(int64)
		err := rows.Scan(subsId)
		if nil != err {
			continue
		}
		result = append(result, *subsId)
	}
	return result, nil
}

const (
	UpdateOneUserSubscriptionIsDeleteSQL = `UPDATE tb_user_subscription SET is_delete= ? WHERE subs_id= ? AND user_id= ?`
)

// Update the 'is_delete' of one row data in 'tb_user_subscription' table.
func UpdateOneUserSubscriptionIsDelete(isDelete bool, subsId, userId int64) error {
	return execSqlWithTransaction(UpdateOneUserSubscriptionIsDeleteSQL, isDelete, subsId, userId)
}

// -----------------------------------------------------------------------
