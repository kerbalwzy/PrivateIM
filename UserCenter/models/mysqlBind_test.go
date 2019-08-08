package models

import (
	"testing"
)

var testUser = UserBasic{
	Name:   "wang",
	Email:  "test@test333.com",
	Mobile: "13122222222"}

func TestUserBasic_Save(t *testing.T) {
	testUser.SetPassword("nihaoceshia")
	if err := testUser.MySQLSignUp(); nil != err {
		t.Error(err)
	}

}

func TestUserBasic_GetByField(t *testing.T) {
	if err := testUser.MySQLGetByField("Name"); nil != err {
		t.Error(err)
	}
	if testUser.Id == 0 {
		t.Error("don`t get value from databases")
	}
}
