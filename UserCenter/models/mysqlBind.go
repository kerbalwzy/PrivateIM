package models

import (
	"database/sql"
	"errors"
	"github.com/bwmarrin/snowflake"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	MySQLClient   = new(sql.DB)
	SnowFlakeNode = new(snowflake.Node)
)

func init() {
	var err error
	MySQLClient, err = sql.Open("mysql", UserDbMySQLURI)
	if nil != err {
		log.Fatal(err)
	}
	SnowFlakeNode, err = snowflake.NewNode(0)
	if nil != err {
		log.Fatal(err)
	}

}

// user basic information sql strings
const (
	UserDbMySQLURI = "root:mysql@tcp(10.211.55.4:3306)/IMUserCenter?charset=utf8&parseTime=true"

	UserNewOne = "INSERT INTO tb_user_basic (id, name, email, password)VALUES (?, ?, ?, ?);"

	UserGetProfileBasic = "SELECT id, name, mobile, email, gender, create_time, password FROM tb_user_basic "

	UserGetProfileById = UserGetProfileBasic + "WHERE id = ?"

	UserGetProfileByEmail = UserGetProfileBasic + "WHERE email = ?"

	UserGetProfileByName = UserGetProfileBasic + "WHERE name = ?"

	UserUpdateProfile = "UPDATE tb_user_basic SET name=?, mobile=?, gender=? WHERE id = ?"
)

// scan user from the row
func ScanUserFromRow(rowP *sql.Row, userP *UserBasic) error {
	err := rowP.Scan(&(userP.Id), &(userP.Name), &(userP.Mobile), &(userP.Email), &(userP.Gender),
		&(userP.CreateTime), &(userP.password))
	if nil != err {
		return err
	}
	if userP.Id == 0 {
		return errors.New("get user profile fail")
	}
	return nil
}

// Get user all information in `tb_user_basic` table by ID
func MySQLGetUserById(userP *UserBasic) error {
	rowP := MySQLClient.QueryRow(UserGetProfileById, userP.Id)
	err := ScanUserFromRow(rowP, userP)
	if nil != err {
		return err
	}
	return nil
}

// Get user all information in `tb_user_basic` table by email
func MySQLGetUserByEmail(userP *UserBasic) error {
	rowP := MySQLClient.QueryRow(UserGetProfileByEmail, userP.Email)
	err := ScanUserFromRow(rowP, userP)
	if nil != err {
		return err
	}
	return nil
}

// Get users all information in `tb_user_basic` table by name
func MySQLGetUserByName(name string) ([]*UserBasic, error) {
	rows, err := MySQLClient.Query(UserGetProfileByName, name)
	if nil != err {
		return nil, err
	}
	users := make([]*UserBasic, 0)
	for rows.Next() {
		user := &UserBasic{}
		err := rows.Scan(&(user.Id), &(user.Name), &(user.Mobile), &(user.Email), &(user.Gender),
			&(user.CreateTime), &(user.password))
		if nil != err {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

// Save user with id, name, email,password to database.
// If successful, get full information of user from database and update to user.
func MySQLUserSignUp(user *UserBasic) error {
	// start a Transaction
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}

	// try to insert user data into database
	id := SnowFlakeNode.Generate()
	_, err = tx.Exec(UserNewOne, id, user.Name, user.Email, user.password)
	if nil != err {
		tx.Rollback()
		return err
	}

	// try to get full information of user from database, and update to user.
	err = tx.QueryRow(UserGetProfileById, id).Scan(&(user.Id), &(user.Name), &(user.Mobile), &(user.Email), &(user.Gender),
		&(user.CreateTime), &(user.password))
	if nil != err {
		tx.Rollback()
		return err
	}

	// commit Transaction
	err = tx.Commit()
	if nil != err {
		tx.Rollback()
	}
	return nil
}

// Update name,mobile and gender of user basic by id
func MySQLUpdateProfile(name, mobile string, gender int, userId int64) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		tx.Rollback()
		return err
	}
	// update user profile
	ret, err := tx.Exec(UserUpdateProfile, name, mobile, gender, userId)
	if nil != err {
		tx.Rollback()
		return err
	}
	aff, err := ret.RowsAffected()
	if aff == 0 {
		tx.Rollback()
		return errors.New("not thing need to update")
	}
	if nil != err {
		tx.Rollback()
		return err
	}
	// commit Transaction
	err = tx.Commit()
	if nil != err {
		tx.Rollback()
	}
	return nil
}

// user more information sql string
const (
	UserGetAvatar = "SELECT avatar FROM tb_user_more WHERE user_id = ?"

	UserInsertOrUpdateAvatar = "INSERT INTO tb_user_more (user_id, avatar) VALUES (?, ?)  ON DUPLICATE KEY UPDATE avatar=?;"

	UserAvatarHashNameCount = "SELECT COUNT(user_id) FROM tb_user_more WHERE avatar=?"

	UserGetQRCode = "SELECT qr_code FROM tb_user_more WHERE user_id = ?"

	UserInsertOrUpdateQRCode = "INSERT INTO tb_user_more (user_id, qr_code) VALUES (?, ?)  ON DUPLICATE KEY UPDATE qr_code=?;"
)

// Get user avatar name by user id
func MySQLGetUserAvatar(userId int64, avatarP *string) error {
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
func MySQLPutUserAvatar(userId int64, hashName string) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(UserInsertOrUpdateAvatar, userId, hashName, hashName)
	if nil != err {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if nil != err {
		tx.Rollback()
		return err
	}
	return nil
}

// Get the count of avatar hash name in tb_user_more table
func MySQLAvatarHashNameCount(hashName string) int {
	row := MySQLClient.QueryRow(UserAvatarHashNameCount, hashName)
	count := new(int)
	err := row.Scan(count)
	if nil != err {
		return 0
	}
	return *count
}

// Get user QRCode name by user id
func MySQLGetUserQRCode(userId int64, hashNameP *string) error {
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
func MySQLPutUserQRCode(userId int64, hashName string) error {
	tx, err := MySQLClient.Begin()
	if nil != err {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(UserInsertOrUpdateQRCode, userId, hashName, hashName)
	if nil != err {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if nil != err {
		tx.Rollback()
		return err
	}
	return nil
}

// user relationship information sql strings
const (
	UserGetFriendBasic = "SELECT id, dst_id, src_id, note, is_accept, is_refuse, is_delete FROM tb_friend_relation "

	UserGetOneFriend = UserGetFriendBasic + "WHERE src_id = ? AND dst_id = ?"

	UserGetAllFriends = UserGetFriendBasic + "WHERE src_id = ?"

	UserAddOneFriend = `INSERT INTO tb_friend_relation(id, src_id, dst_id, note,is_delete) VALUES(?,?,?,?,?) 
ON DUPLICATE KEY UPDATE note = ?,is_delete=false`

	UserAcceptOneFriend = ""

	UserBlackOneFriend = ""

	UserDeleteOneFriend = ""
)

// Get one friend relation information of user
func MySQLGetUserOneFriend(relateP *UserRelate) error {
	rowP := MySQLClient.QueryRow(UserGetOneFriend, relateP.SelfId, relateP.FriendId)
	err := rowP.Scan(&(relateP.Id), &(relateP.SelfId), &(relateP.FriendId),
		&(relateP.FriendNote), &(relateP.IsAccept), &(relateP.IsRefuse), &(relateP.IsDelete))
	if nil != err {
		return err
	}
	return nil
}

// Get all friends relation information of uer
func MySQLGetUserAllFriends(userId int64) ([]*UserRelate, error) {
	rows, err := MySQLClient.Query(UserGetAllFriends, userId)
	if nil != err {
		return nil, err
	}
	friends := make([]*UserRelate, 0)
	for rows.Next() {
		relateP := new(UserRelate)
		err := rows.Scan(&(relateP.Id), &(relateP.SelfId), &(relateP.FriendId),
			&(relateP.FriendNote), &(relateP.IsAccept), &(relateP.IsRefuse), &(relateP.IsDelete))
		if nil != err {
			continue
		}
		friends = append(friends, relateP)
	}
	return friends, nil
}

// Add one friend relation information of user
func MySQLAddOneFriend(relateP *UserRelate) error {
	// open a Transaction
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}

	// try save or update a relationship row data
	relateP.Id = SnowFlakeNode.Generate().Int64()
	_, err = tx.Exec(UserAddOneFriend, relateP.Id, relateP.SelfId, relateP.FriendId, relateP.FriendNote, relateP.FriendNote)

	return nil
}
