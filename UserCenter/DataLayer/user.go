package DataLayer

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"time"
)

const (
	PasswordSalt = "fasdfasf87tr3h87sf23386t123!@e23BLfishf"
)

// user basic information in `tb_user_basic` table
type UserBasic struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name" `
	Mobile     string    `json:"mobile"`
	Email      string    `json:"email"`
	Gender     int       `json:"gender"`
	CreateTime time.Time `json:"create_time" time_format:"2006-01-02 15:04:05"`

	password string
}

// get a password hash vale with salts
func GetPasswordHash(password string, salts ...string) string {
	var tempStr string
	tempStr += password
	for _, salt := range salts {
		tempStr += salt
	}

	mac := hmac.New(md5.New, nil)
	mac.Write([]byte(tempStr))
	return hex.EncodeToString(mac.Sum(nil))
}

// set password for user by MD5 hash value
func (user *UserBasic) SetPassword(password string) string {
	user.password = GetPasswordHash(password, user.Email, PasswordSalt)
	return user.password
}

// check password for user
func (user *UserBasic) CheckPassword(password string) bool {
	return user.password == GetPasswordHash(password, user.Email, PasswordSalt)
}

// user more information in `tb_user_more` table
type UserMore struct {
	UserId int64  `json:"user_id"`
	Avatar string `json:"avatar"`
	QrCode string `json:"qr_code"`
}

// user relationship information in `tb_friend_relation` table
type UserRelate struct {
	Id         int64  `json:"id"`
	SelfId     int64  `json:"self_id"`
	FriendId   int64  `json:"friend_id"`
	FriendNote string `json:"friend_note"`
	IsAccept   bool   `json:"is_accept"`
	IsBlack    bool   `json:"is_black"`
	IsDelete   bool   `json:"is_delete"`
}

// user basic and relate information from `tb_user_basic` and `tb_friend_relation` table
type FriendInformation struct {
	FriendId int64  `json:"friend_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Gender   int    `json:"gender"`
	Note     string `json:"note"`
	IsBlack  bool   `json:"is_black"`
}
