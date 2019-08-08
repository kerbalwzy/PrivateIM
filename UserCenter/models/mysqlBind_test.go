package models

import (
	"testing"
)

func TestUserBasic_Save(t *testing.T) {
	testUser := UserBasic{
		Name:  "wang",
		Email: "test@test.com"}
	if err := testUser.Save(); nil != err {
		t.Error(err)
	}

}
