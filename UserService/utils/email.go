package utils

import (
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/gomail.v2"
	"time"

	conf "../Config"
)

const ResetPasswordEmailTemplate = `
<div style="border:solid 2px gray;
            padding:10px;
            outline:gray dashed 10px;
            width:618px;">
  <h2>
    Reset Password By Email, Sent from PrivateIM. 
  </h2>
  <h4>
    Please click or copy the follow link redirect to the page for reset password, it will expire after 5 minute.
  </h4>
  <a href="%s">%s</a>
  <h4>
	In addition, after 5 minutes, you can ask for a new email to be sent.
  </h4>
</div>
`

var rdb *redis.Client
var emailClient *gomail.Dialer

func init() {
	rdb = GetRedisClient(conf.RedisAddr, "", 0)
	emailClient = gomail.NewDialer(conf.EmailServerHost, conf.EmailServerPort,
		conf.EmailAuthUserName, conf.EmailAuthPassword)
}

// Check the record of sent an email for reset password to the user in 10 minute.
// If sent already, return true, else return false
func CheckEmailSentIn3Minute(email string) bool {
	_, err := rdb.Get(fmt.Sprintf("ResetPassword_%s", email)).Result()
	if nil != err {
		// err happened, thinking not found the value, meaning not sent in 3 minute
		return false
	} else {
		return true
	}
}

// Save a tag value in redis, to mark have sent the reset password email for the user.
// The expire time of the data is 3 minute.
func SetEmailSentTag(email string) {
	_ = rdb.SetNX(fmt.Sprintf("ResetPassword_%s", email),
		1, conf.RestPasswordEmailSentTagAliveTime*time.Second)

}

// Send an email for reset password to the user, the core content is a authentication
// string, that can mark user was sign in, use to update the password
func SendPasswordResetEmail(email, authToken string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", conf.EmailAuthUserName)
	m.SetHeader("To", email)
	m.SetHeader("Subject", conf.RestPasswordEmailSubject)

	coreLink := conf.RestPasswordPageBaseLink + authToken
	contentBody := fmt.Sprintf(ResetPasswordEmailTemplate, coreLink, coreLink)
	m.SetBody("text/html", contentBody)

	err := emailClient.DialAndSend(m)
	if nil != err {
		return err
	}
	SetEmailSentTag(email)
	return nil
}
