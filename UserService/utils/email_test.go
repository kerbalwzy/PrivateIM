package utils

import "testing"

func TestSendPasswordResetEmail(t *testing.T) {
	email := "634640761@qq.com"
	token := "this is a test reset password auth token, hahahahahahahahahahahaha"
	err := SendPasswordResetEmail(email, token)
	if nil != err {
		t.Error("SendPasswordResetEmail fail: ", err)
	} else {
		t.Logf("SendPasswordResetEmail success: reciver email is %s", email)
	}
}
