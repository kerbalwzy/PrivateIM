package ApiRPC

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)
import "../utils"
import conf "../Config"

var testUserId int64 = 10000

func TestCheckAuthToken(t *testing.T) {

	claims := utils.CustomJWTClaims{
		Id: testUserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + conf.AuthTokenAliveTime),
			Issuer:    conf.AuthTokenIssuer,
		}}
	authToken, _ := utils.CreateJWTToken(claims, []byte(conf.AuthTokenSalt))

	data, err := CheckAuthToken(authToken)
	if nil != err {
		t.Error("CheckAuthToken fail: ", err)
	} else {
		t.Logf("CheckAuthToken success: raw user's id = %d, parse from token get id = %d", testUserId, data.UserId)
	}
}
