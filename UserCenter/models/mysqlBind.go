package models

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/snowflake"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	UserSignUpSql     = "INSERT INTO tb_user_basic (id, name, email, password)VALUES (?, ?, ?, ?);"
	GetUserByFieldSql = "SELECT(id, name, mobile, email, gender, create_time, update_time)FROM tb_user_basic WHERE ?=?;"
)

var (
	UserDbSQLite3Path = "/Users/wzy/GitProrgram/PrivateIM/UserCenter/userDb.SQLite3"
	SQLite3           = new(sql.DB)
	SnowFlakeNode     = new(snowflake.Node)
)

func init() {
	//path, _ := filepath.Abs(os.Args[0])
	//BaseDir := filepath.Dir(path)
	//UserDbSQLite3Path = filepath.Join(BaseDir, "userDb.SQLite3")
	var err error
	SQLite3, err = sql.Open("sqlite3", UserDbSQLite3Path)
	if nil != err {
		log.Fatal(err)
	}
	SnowFlakeNode, err = snowflake.NewNode(0)
	if nil != err {
		log.Fatal(err)
	}

}

// If success, saving the id, name, email on the user instance
func (u *UserBasic) Save() error {
	UserInsertStmt, err := SQLite3.Prepare(UserSignUpSql)
	if nil != err {
		return err
	}
	id := SnowFlakeNode.Generate()
	ret, err := UserInsertStmt.Exec(id, u.Name, u.Email, u.password)
	if nil != err {
		return err
	}
	aff, err := ret.RowsAffected()
	if 0 == aff || nil != err {
		return err
	}
	u.Id = id.Int64()
	return nil
}

// Query an user by name and password, For SignIn API.
// And the password is a hash value, not a plaintext
func (u *UserBasic) QueryUserByX(name, password string) error {

	return nil
}

// Get an user by id or email or mobile
func (u *UserBasic) GetUserByField(field string) error {
	switch field {
	case "id":
		raw, err := SQLite3.Query(GetUserByFieldSql, field, u.Id)
		if nil != err {
			return err
		}
		for raw.Next() {
			err = raw.Scan(u.Id, u.Name, u.Mobile, u.Email, u.Gender, u.CreateTime, u.UpdateTIme)
			if nil != err {
				return err
			}
		}

	case "email":
		raw, err := SQLite3.Query(GetUserByFieldSql, field, u.Email)
		if nil != err {
			log.Printf("@@@@@@@@@@")
			return err
		}
		for raw.Next() {
			log.Println("xxxxxxxxxx")
			err = raw.Scan(u.Id, u.Name, u.Mobile, u.Email, u.Gender, u.CreateTime, u.UpdateTIme)
			if nil != err {
				return err
			}
		}

	case "mobile":
		raw, err := SQLite3.Query(GetUserByFieldSql, field, u.Mobile)
		if nil != err {
			return err
		}
		for raw.Next() {
			err = raw.Scan(u.Id, u.Name, u.Mobile, u.Email, u.Gender, u.CreateTime, u.UpdateTIme)
			if nil != err {
				return err
			}
		}

	default:
		return fmt.Errorf("can only chose in id, email and mobile")
	}

	//}
	return nil
}
