package utils

import (
	"testing"
	"time"
)

type TestUserBasic struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name" `
	Mobile     string    `json:"mobile"`
	Email      string    `json:"email"`
	Gender     int       `json:"gender"`
	CreateTime time.Time `json:"create_time" time_format:"2006-01-02 15:04:05"`
	UpdateTIme time.Time `json:"update_time" time_format:"2006-01-02 15:04:05"`

	password string
}

func TestGetReflectValueByField(t *testing.T) {
	user := TestUserBasic{Id: 100, Name: "wang", Email: "test@test.com", Mobile: "13122222222"}
	if id, _ := GetReflectValueByField(user, "Id"); id != user.Id {
		t.Logf("id = %v", id)
		t.Fail()
	}
	if name, _ := GetReflectValueByField(user, "Name"); name != user.Name {
		t.Logf("name = %v", name)
		t.Fail()
	}
	if email, _ := GetReflectValueByField(user, "Email"); email != user.Email {
		t.Logf("email = %v", email)
		t.Fail()
	}
	if mobile, _ := GetReflectValueByField(user, "Mobile"); mobile != user.Mobile {
		t.Logf("mobile = %v", mobile)
		t.Fail()
	}
}
