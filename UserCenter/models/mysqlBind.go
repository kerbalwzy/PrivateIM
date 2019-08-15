package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	_ "github.com/go-sql-driver/mysql"
	"log"

	"../utils"
)

const (
	UserDbMySQLURI = "root:mysql@tcp(10.211.55.4:3306)/IMUserCenter?charset=utf8&parseTime=true"

	UserSignUpSql = "INSERT INTO tb_user_basic (id, name, email, password)VALUES (?, ?, ?, ?);"

	UserGetByFieldSql = "SELECT id, name, mobile, email, gender, create_time, password FROM tb_user_basic WHERE %s = ?;"

	UserPutProfileSql = "UPDATE tb_user_basic SET name=?, mobile=?, gender=? WHERE id = ?"

	UserGetAvatarSql = "SELECT avatar FROM tb_user_more WHERE user_id = ?"

	UserInsertOrUpdateAvatar = "INSERT INTO tb_user_more (user_id) VALUES (?)  ON DUPLICATE KEY UPDATE avatar=?;"

	UserAvatarHashNameCount = "SELECT COUNT(user_id) FROM tb_user_more WHERE avatar=?"

	UserGetQRCodeSql = "SELECT qr_code FROM tb_user_more WHERE user_id = ?"

	UserInsertOrUpdateQRCode = "INSERT INTO tb_user_more (user_id) VALUES (?)  ON DUPLICATE KEY UPDATE qr_code=?;"
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

// Get an user basic by id or email or mobile
func MySQLGetUserByField(field string, user *UserBasic) error {
	value, err := utils.GetReflectValueByField(*user, field)
	if nil != err {
		return err
	}
	tempSQL := fmt.Sprintf(UserGetByFieldSql, field)
	row := MySQLClient.QueryRow(tempSQL, value)
	err = row.Scan(&(user.Id), &(user.Name), &(user.Mobile), &(user.Email), &(user.Gender),
		&(user.CreateTime), &(user.password))
	if nil != err {
		return err
	}
	if user.Id == 0 {
		return fmt.Errorf("get user by <%s> fail", field)
	}
	return nil
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
	_, err = tx.Exec(UserSignUpSql, id, user.Name, user.Email, user.password)
	if nil != err {
		tx.Rollback()
		return err
	}

	// try to get full information of user from database, and update to user.
	tmpSql := fmt.Sprintf(UserGetByFieldSql, "id")
	err = tx.QueryRow(tmpSql, id).Scan(&(user.Id), &(user.Name), &(user.Mobile), &(user.Email), &(user.Gender),
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
	ret, err := tx.Exec(UserPutProfileSql, name, mobile, gender, userId)
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

// Get user avatar name by user id
func MySQLGetUserAvatar(userId int64, avatar *string) error {
	row := MySQLClient.QueryRow(UserGetAvatarSql, userId)
	err := row.Scan(avatar)

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
	_, err = tx.Exec(UserInsertOrUpdateAvatar, userId, hashName)
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
func MySQLGetUserQRCode(userId int64, qrCode *string) error {
	row := MySQLClient.QueryRow(UserGetQRCodeSql, userId)
	err := row.Scan(qrCode)

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
	_, err = tx.Exec(UserInsertOrUpdateQRCode, userId, hashName)
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
