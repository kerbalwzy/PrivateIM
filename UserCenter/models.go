package MySQLBind

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

const (
	PasswordSalt = "fasdfasf87tr3h87sf23386t123!@e23BLfishf"
)

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
