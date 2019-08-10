package models

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"time"
)

const (
	PasswordSalt = "fasdfasf87tr3h87sf23386t123!@e23BLfishf"
)

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


type UserMore struct {
	UserId int64  `json:"user_id"`
	Avatar string `json:"avatar"`
	QrCode string `json:"qr_code"`
}
