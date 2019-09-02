package ApiRPC

import "testing"

var (
	name     = "testName"
	email    = "testEmail@tt.com"
	mobile   = "13100000000"
	password = "this would be a hash value as password"
	gender   = 0

	userId int64
)

func TestNewOneUser(t *testing.T) {
	user, err := NewOneUser(name, email, mobile, password, gender)
	if nil != err {
		t.Error("NewOneUser Error: ", err)
	}
	if user.Id == 0 {
		t.Error("NewOneUser Error: the Id is zero")
	}
	if user.Name != name || user.Email != email || user.Mobile != mobile ||
		user.Password != password || user.Gender != int32(gender) {
		t.Error("NewOneUser Error: the raw data changed after insert")
	}
	if user.CreateTime == "" {
		t.Error("NewOneUser Error: the createTime is empty string")
	}
	t.Logf("NewOneUser Success: the user's id=%d, createTime=%s",
		user.Id, user.CreateTime)

	userId = user.Id
}

func TestQueryUserById(t *testing.T) {
	user, err := QueryUserById(userId)
	if nil != err {
		t.Error("QueryUserById Error: ", err)
	}
	if user.Id != userId {
		t.Error("QueryUserById Error: the user's id is not equal the raw query value")
	}
	if user.Name != name || user.Email != email || user.Mobile != mobile ||
		user.Password != password || user.Gender != int32(gender) {
		t.Error("QueryUserById Error: the user's data is not equal the raw value")
	}
	t.Logf("QueryUserById Success: the user's id=%d, createTime=%s",
		user.Id, user.CreateTime)

}

func TestQueryUserByEmail(t *testing.T) {
	user, err := QueryUserByEmail(email)
	if nil != err {
		t.Error("QueryUserByEmail Error: ", err)
	}
	if user.Id != userId {
		t.Error("QueryUserByEmail Error: the user's id is not equal the raw query value")
	}
	if user.Name != name || user.Email != email || user.Mobile != mobile ||
		user.Password != password || user.Gender != int32(gender) {
		t.Error("QueryUserByEmail Error: the user's data is not equal the raw value")
	}
	t.Logf("QueryUserByEmail Success: the user's id=%d, createTime=%s",
		user.Id, user.CreateTime)
}

func TestQueryUsersByName(t *testing.T) {
	users, err := QueryUsersByName(name)
	if nil != err {
		t.Error("QueryUsersByName Error: ", err)
	}
	if len(users.Users) == 0 {
		t.Error("QueryUserByName Error: this result list is empty")
	}
	t.Logf("QueryUsersByName Success: the zero index element: user's id=%d, createTime=%s", users.Users[0].Id, users.Users[0].CreateTime)
}
