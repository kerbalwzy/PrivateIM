package utils

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

var TempToken string

func TestCreateJWTToken(t *testing.T) {
	var err error
	claims := CustomJWTClaims{
		Id: 1234124312342,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + 1000),
			Issuer:    "test"}}

	TempToken, err = CreateJWTToken(claims, []byte("test salt"))
	if nil != err {
		t.Error(err)
	} else {
		t.Logf("Token: %s", TempToken)
	}

}

func TestParseJWTToken(t *testing.T) {
	claims, err := ParseJWTToken(TempToken, []byte("test salt"))
	if nil != err {
		t.Error(err)
	}
	if claims == nil {
		t.Fail()
	} else {
		t.Logf("claims Id: %d, Expire: %d", claims.Id, claims.ExpiresAt)
	}

}


func TestRefreshJWTToken(t *testing.T) {
	var err error
	TempToken, err = RefreshJWTToken(TempToken, []byte("test salt"), time.Hour)
	if nil != err {
		t.Error(err)
	} else {
		t.Logf("Token: %s", TempToken)
	}

}
