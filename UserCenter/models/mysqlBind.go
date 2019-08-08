package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/snowflake"
	_ "github.com/go-sql-driver/mysql"
)

const (
	UserDbMySQLURI = "root:mysql@tcp(10.211.55.4:3306)/IMUserCenter?charset=utf8&parseTime=true"

	UserSignUpSql = "INSERT INTO tb_user_basic (id, name, email, password)VALUES (?, ?, ?, ?);"

	GetUserByFieldSql = "SELECT id, name, mobile, email, gender, create_time, update_time, password FROM tb_user_basic WHERE %s=?;"
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

// Save user with id, name, email,password to database.
// If successful, get full information of user from database and update to user.
func (user *UserBasic) MySQLSignUp() error {
	// start a Transaction
	tx, err := MySQLClient.Begin()
	if nil != err {
		return err
	}

	// try to insert user data into database
	id := SnowFlakeNode.Generate()
	ret, err := tx.Exec(UserSignUpSql, id, user.Name, user.Email, user.password)
	if nil != err {
		tx.Rollback()
		return err
	}
	aff, err := ret.RowsAffected()
	if 0 == aff || nil != err {
		tx.Rollback()
		return err
	}

	// try to get full information of user from database, and update to user.
	tmpSql := fmt.Sprintf(GetUserByFieldSql, "id")
	err = tx.QueryRow(tmpSql, id).Scan(&(user.Id), &(user.Name), &(user.Mobile), &(user.Email), &(user.Gender),
		&(user.CreateTime), &(user.UpdateTime), &(user.password))
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

// Get an user basic by id or email or mobile
func (user *UserBasic) MySQLGetByField(field string) error {
	value, err := GetReflectValueByField(*user, field)
	if nil != err {
		return err
	}
	tempSQL := fmt.Sprintf(GetUserByFieldSql, field)
	row := MySQLClient.QueryRow(tempSQL, value)
	err = row.Scan(&(user.Id), &(user.Name), &(user.Mobile), &(user.Email), &(user.Gender),
		&(user.CreateTime), &(user.UpdateTime), &(user.password))
	if nil != err {
		return err
	}
	if user.Id == 0 {
		return fmt.Errorf("get user by <%s> fail", field)
	}
	return nil
}
