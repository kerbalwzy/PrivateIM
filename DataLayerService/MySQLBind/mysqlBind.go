package MySQLBind

import (
	"database/sql"
	"errors"
	"log"
	"time"

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

// User's basic information sql strings
const (
	UserNewOne = "INSERT INTO tb_user_basic (id, name, email, mobile, password, " +
		"gender) VALUES (?, ?, ?, ?, ?, ?);"

	UserGetProfileBasic = "SELECT id, name, mobile, email, gender, create_time, " +
		"password FROM tb_user_basic "

	UserGetProfileById = UserGetProfileBasic + "WHERE id = ?"

	UserGetProfileByEmail = UserGetProfileBasic + "WHERE email = ?"

	UserGetProfileByName = UserGetProfileBasic + "WHERE name = ?"

	UserUpdateProfile = "UPDATE tb_user_basic SET name=?, mobile=?, gender=? " +
		"WHERE id = ?"
)

var UserNotExitedError = errors.New("the user not existed")

// user basic information in `tb_user_basic` table
type TempUserBasic struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name" `
	Mobile     string    `json:"mobile"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Gender     int       `json:"gender"`
	CreateTime time.Time `json:"create_time" time_format:"2006-01-02 15:04:05"`
}

// Save user with id, name, email,password to database.
// If successful, get full information of user from database and update to user.
func InsertOneUser(name, email, mobile, password string, gender int) (
	*TempUserBasic, error) {
	// start a Transaction
	tx, err := MySQLClient.Begin()
	if nil != err {
		return nil, err
	}

	// try to insert user data into database
	id := SnowFlakeNode.Generate()
	_, err = tx.Exec(UserNewOne, id, name, email, mobile, password, gender)
	if nil != err {
		_ = tx.Rollback()
		return nil, err
	}
	user := new(TempUserBasic)
	// try to get full information of user from database, and update to user.
	err = tx.QueryRow(UserGetProfileById, id).Scan(&(user.Id), &(user.Name),
		&(user.Mobile), &(user.Email), &(user.Gender), &(user.CreateTime),
		&(user.Password))
	if nil != err {
		_ = tx.Rollback()
		return nil, err
	}

	// commit Transaction
	err = tx.Commit()
	if nil != err {
		_ = tx.Rollback()
		return nil, err
	}
	return user, nil
}

// Scan user information from the row
func ScanUserFromRow(row *sql.Row) (*TempUserBasic, error) {
	user := new(TempUserBasic)
	err := row.Scan(&(user.Id), &(user.Name), &(user.Mobile), &(user.Email),
		&(user.Gender), &(user.CreateTime), &(user.Password))
	if nil != err {
		return nil, err
	}
	if user.Id == 0 {
		return nil, UserNotExitedError
	}
	return user, nil
}

// Get user all information in `tb_user_basic` table by id
func QueryUserById(id int64) (*TempUserBasic, error) {
	row := MySQLClient.QueryRow(UserGetProfileById, id)
	return ScanUserFromRow(row)
}

// Get user all information in `tb_user_basic` table by email
func QueryUserByEmail(email string) (*TempUserBasic, error) {
	row := MySQLClient.QueryRow(UserGetProfileByEmail, email)
	return ScanUserFromRow(row)

}

// Get users all information in `tb_user_basic` table by name
func QueryUsersByName(name string) ([]*TempUserBasic, error) {
	rows, err := MySQLClient.Query(UserGetProfileByName, name)
	if nil != err {
		return nil, err
	}
	users := make([]*TempUserBasic, 0)
	for rows.Next() {
		user := new(TempUserBasic)
		err := rows.Scan(&(user.Id), &(user.Name), &(user.Mobile), &(user.Email), &(user.Gender),
			&(user.CreateTime), &(user.Password))
		if nil != err {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

// Update name,mobile and gender of user basic by id
func UpdateProfile(name, mobile string, gender int, userId int64) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}
	// update user profile
	ret, err := tx.Exec(UserUpdateProfile, name, mobile, gender, userId)
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	aff, err := ret.RowsAffected()
	if aff == 0 {
		_ = tx.Rollback()
		return errors.New("not thing need to update")
	}
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	// commit Transaction
	err = tx.Commit()
	if nil != err {
		_ = tx.Rollback()
	}
	return nil
}

// User's more information sql string
const (
	UserGetAvatar = "SELECT avatar FROM tb_user_more WHERE user_id = ?"

	UserInsertOrUpdateAvatar = "INSERT INTO tb_user_more (user_id, avatar) VALUES (?, ?)  ON DUPLICATE KEY UPDATE avatar=?;"

	UserAvatarHashNameCount = "SELECT COUNT(user_id) FROM tb_user_more WHERE avatar=?"

	UserGetQRCode = "SELECT qr_code FROM tb_user_more WHERE user_id = ?"

	UserInsertOrUpdateQRCode = "INSERT INTO tb_user_more (user_id, qr_code) VALUES (?, ?)  ON DUPLICATE KEY UPDATE qr_code=?;"
)

// user more information in `tb_user_more` table
type UserMore struct {
	UserId int64  `json:"user_id"`
	Avatar string `json:"avatar"`
	QrCode string `json:"qr_code"`
}

// Get user avatar name by user id
func GetUserAvatar(userId int64, avatarP *string) error {
	row := MySQLClient.QueryRow(UserGetAvatar, userId)
	err := row.Scan(avatarP)

	// if not found, it dose not need to abort en error, but return.
	if err == sql.ErrNoRows {
		return nil
	}
	if nil != err {
		return err
	}
	return nil
}

// Insert or Update avatar hash name into database
func PutUserAvatar(userId int64, hashName string) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}
	_, err = tx.Exec(UserInsertOrUpdateAvatar, userId, hashName, hashName)
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

// Get the count of avatar hash name in tb_user_more table
func AvatarHashNameCount(hashName string) int {
	row := MySQLClient.QueryRow(UserAvatarHashNameCount, hashName)
	count := new(int)
	err := row.Scan(count)
	if nil != err {
		return 0
	}
	return *count
}

// Get user QRCode name by user id
func GetUserQRCode(userId int64, hashNameP *string) error {
	row := MySQLClient.QueryRow(UserGetQRCode, userId)
	err := row.Scan(hashNameP)

	// if not found, it dose not need to abort en error, but return.
	if err == sql.ErrNoRows {
		return nil
	}
	if nil != err {
		return err
	}
	return nil
}

// Insert or Update QRCode hash name into database
func PutUserQRCode(userId int64, hashName string) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}
	_, err = tx.Exec(UserInsertOrUpdateQRCode, userId, hashName, hashName)
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

// User's relationship information sql strings
const (
	UserGetFriendsRelate = `SELECT id, src_id, dst_id, note, is_accept, is_black, is_delete FROM tb_friend_relation 
WHERE src_id = ?`

	UserCheckTargetFriendExisted = "SELECT id FROM tb_user_basic WHERE id = ?"

	UserCheckBlackList = `SELECT is_black FROM tb_friend_relation WHERE src_id = ? AND dst_id = ?`

	UserCheckFriendshipAlreadyInEffect = "SELECT is_accept FROM tb_friend_relation WHERE src_id = ? AND dst_id = ?"

	UserAddOneFriend = `INSERT INTO tb_friend_relation(id, src_id, dst_id, note) VALUES(?,?,?,?) 
ON DUPLICATE KEY UPDATE note = ?,is_accept = FALSE, is_black=FALSE, is_delete = FALSE`

	UserCheckFriendRequest = `SELECT id from tb_friend_relation WHERE src_id =? AND dst_id = ?`

	UserAcceptOneFriend = `INSERT INTO tb_friend_relation(id, src_id, dst_id, is_accept) VALUES(?,?,?,?) 
ON DUPLICATE KEY UPDATE is_accept = TRUE, is_black = FALSE, is_delete = FALSE `

	UserCheckBlacklist = `SELECT id, is_black FROM tb_friend_relation WHERE src_id = ? AND dst_id = ?`

	UserBlackOneFriend = `INSERT INTO tb_friend_relation(id, src_id, dst_id, is_black) VALUES(?,?,?,?) 
ON DUPLICATE KEY UPDATE is_black = ? `

	UserNoteOneFriend = `UPDATE tb_friend_relation SET note = ? WHERE src_id = ? AND dst_id = ?`

	UserDeleteOneFriend = `UPDATE tb_friend_relation SET is_accept = FALSE, is_black = FALSE, is_delete=TRUE 
WHERE src_id=? AND dst_id = ?`

	UserGetFriendsInfo = `SELECT id, name, email, mobile, gender, note, is_black FROM tb_user_basic as basic, 
(SELECT dst_id, note, is_black from tb_friend_relation where src_id= ? and is_delete = FALSE and is_accept=TRUE) 
as friends where friends.dst_id = basic.id`
)

var (
	ErrNoFriendship              = errors.New("you are not friends yet")
	ErrTargetUserNotExisted      = errors.New("the target user you want dose not existed")
	ErrInBlackList               = errors.New("you are in the black list of target user")
	ErrFriendshipAlreadyInEffect = errors.New("your friendship already in effect")
	ErrFriendRequestNotExisted   = errors.New("there have not a friend request you can accept")
	ErrFriendBlacklistNoChange   = errors.New("the status of friend blacklist is not change")
)

// Add one friend relation information of user
func AddOneFriend(selfId, friendId int64, note string) error {
	// open a Transaction
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}
	// check the target user if existed
	row := tx.QueryRow(UserCheckTargetFriendExisted, friendId)
	err = row.Scan(&friendId)
	if nil != err {
		_ = tx.Rollback()
		return ErrTargetUserNotExisted
	}
	// check the self if existed in target user's black list
	isBlack := new(bool)
	row = tx.QueryRow(UserCheckBlackList, friendId, selfId)
	_ = row.Scan(isBlack)
	if *isBlack {
		_ = tx.Rollback()
		return ErrInBlackList
	}
	// check the friendship is already in effect
	isAccept := new(bool)
	row = tx.QueryRow(UserCheckFriendshipAlreadyInEffect, selfId, friendId)
	_ = row.Scan(isAccept)
	if *isAccept {
		_ = tx.Rollback()
		return ErrFriendshipAlreadyInEffect
	}
	// every thing ok, add a friendship record
	relateId := SnowFlakeNode.Generate()
	_, err = tx.Exec(UserAddOneFriend, relateId, selfId, friendId, note, note)
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

// Update one friend note
func ModifyNoteOfFriend(selfId, friendId int64, note string) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}
	affect, err := tx.Exec(UserNoteOneFriend, note, selfId, friendId)
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	// if affect 0 row, means the friendship not existed
	if count, _ := affect.RowsAffected(); count == 0 {
		_ = tx.Rollback()
		return ErrNoFriendship
	}
	err = tx.Commit()
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	return nil

}

// Handle a friend request, chose accept or not
func AcceptOneFriend(selfId, friendId int64, note string, isAccept bool) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}

	// check the friend request if existed
	row := tx.QueryRow(UserCheckFriendRequest, friendId, selfId)
	friendRecordId := new(int64)
	if err := row.Scan(friendRecordId); nil != err {
		_ = tx.Rollback()
		return ErrFriendRequestNotExisted
	}

	selfRecordId := SnowFlakeNode.Generate()

	// check the friendship if already in effect
	isEffect := new(bool)
	row = tx.QueryRow(UserCheckFriendshipAlreadyInEffect, selfId, friendId)
	_ = row.Scan(isEffect)
	if *isEffect {
		_ = tx.Rollback()
		return ErrFriendshipAlreadyInEffect
	}

	// accept or refuse the friendship request
	if isAccept {
		// add a friend relationship record for self
		_, err = tx.Exec(UserAcceptOneFriend, selfRecordId, selfId, friendId, isAccept)
		if nil != err {
			_ = tx.Rollback()
			return err
		}
		// change the friend relationship record of requester, make the `is_accept` also be true
		_, err = tx.Exec(UserAcceptOneFriend, friendRecordId, friendId, selfId, isAccept)
		if nil != err {
			_ = tx.Rollback()
			return err
		}

	} else {
		// refuse the friend request, also need add one record for self, make the requester in blacklist
		_, err = tx.Exec(UserBlackOneFriend, selfRecordId, selfId, friendId, !isAccept, !isAccept)
		if nil != err {
			_ = tx.Rollback()
			return err
		}
	}
	if note != "" {
		// change the note for friend, if fail not need rollback
		_, _ = tx.Exec(UserNoteOneFriend, note, selfId, friendId)
	}
	if err := tx.Commit(); nil != err {
		_ = tx.Rollback()
		return err
	}
	return nil
}

// Move friend to blacklist in or out
func ManageFriendBlacklist(selfId, friendId int64, isBlack bool) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}
	// check the friendship data if recorded by self
	relateId := new(int64)
	blackRecord := new(bool)
	row := tx.QueryRow(UserCheckBlacklist, selfId, friendId)
	_ = row.Scan(relateId, blackRecord)

	// if the friend blacklist status if not change, don't continue
	if *relateId != 0 && *blackRecord == isBlack {
		_ = tx.Rollback()
		return ErrFriendBlacklistNoChange
	}

	if *relateId == 0 {
		*relateId = SnowFlakeNode.Generate().Int64()
	}

	// move friend to blacklist in or out
	_, err = tx.Exec(UserBlackOneFriend, relateId, selfId, friendId, isBlack, isBlack)
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if nil != err {
		return err
	}
	return nil
}

// Delete friend relationship record
func DeleteOneFriend(selfId, friendId int64) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}

	// update self record
	affect, err := tx.Exec(UserDeleteOneFriend, selfId, friendId)
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
		return ErrNoFriendship
	}

	// update friend record
	affect, err = tx.Exec(UserDeleteOneFriend, friendId, selfId)
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
		return ErrNoFriendship
	}

	err = tx.Commit()
	if nil != err {
		_ = tx.Rollback()
		return err
	}
	return nil

}

// Get all friends relation information of uer
func GetUserFriendsRelates(userId int64) ([]*UserRelate, error) {
	rows, err := MySQLClient.Query(UserGetFriendsRelate, userId)
	if nil != err {
		return nil, err
	}
	friends := make([]*UserRelate, 0)
	for rows.Next() {
		relateP := new(UserRelate)
		err := rows.Scan(&(relateP.Id), &(relateP.SelfId), &(relateP.FriendId),
			&(relateP.FriendNote), &(relateP.IsAccept), &(relateP.IsBlack), &(relateP.IsDelete))
		if nil != err {
			continue
		}
		friends = append(friends, relateP)
	}
	return friends, nil
}

// Get the friends basic and relate information of user
func GetUserFriendsInfo(selfId int64) ([]*FriendInformation, error) {
	rows, err := MySQLClient.Query(UserGetFriendsInfo, selfId)
	if nil != err {
		return nil, err
	}

	friendsInfo := make([]*FriendInformation, 0)
	for rows.Next() {
		temp := new(FriendInformation)
		_ = rows.Scan(&(temp.FriendId), &(temp.Name), &(temp.Email), &(temp.Mobile),
			&(temp.Gender), &(temp.Note), &(temp.IsBlack))

		friendsInfo = append(friendsInfo, temp)
	}

	return friendsInfo, nil
}
