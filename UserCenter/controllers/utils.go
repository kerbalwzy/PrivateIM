package controllers

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//Generator a AuthToken by data(interface), salt(string),expire(unit:second)
//if the expire <= 0, the real expire will be three thousand years, ha.
func MakeAuthToken(data interface{}, salt string, expiry int64) (token string, err error) {
	// make payload base64 string
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err = encoder.Encode(data)
	if nil != err {
		return "", err
	}
	gobData := buff.Bytes()
	bs64DataStr := base64.StdEncoding.EncodeToString(gobData)

	// make expiry datetimeStamp base64 string
	var liveTimeStamp int64
	now := time.Now()
	if expiry <= 0 {
		// 如果有效时间<=0, 则将有效截止时间设置为到3000年后到此时
		liveTimeStamp = NYearLaterTimeStamp(now, 3000)
	} else {
		liveTimeStamp = time.Now().Unix() + expiry
	}
	liveTimeStampStr := fmt.Sprintf("%d", liveTimeStamp)
	bs64TimeStamp := base64.StdEncoding.EncodeToString([]byte(liveTimeStampStr))

	// make the signal base64 string
	h := md5.Sum([]byte(bs64DataStr + bs64TimeStamp + salt))
	signature := fmt.Sprintf("%X", h)
	token = fmt.Sprintf("%s.%s.%s", bs64TimeStamp, bs64DataStr, signature)

	return token, err
}

//Parse a JWTToken string, and check boolean, if true, return the payload data
func CheckAuthToken(token, salt string, data interface{}) (ok bool, err error) {
	//var bs64TimeStamp, bs64DataStr, signature string
	ret := strings.Split(token, ".")
	if len(ret) < 3 {
		err = errors.New("Invalid token string ")
		return
	}
	bs64TimeStamp := ret[0]
	bs64DataStr := ret[1]
	signature := ret[2]
	// check the salt
	h := md5.Sum([]byte(bs64DataStr + bs64TimeStamp + salt))
	if signature != fmt.Sprintf("%X", h) {
		err = errors.New("Invalid token string ")
		return
	}

	// check the expiry time
	timeStampStr, err := base64.StdEncoding.DecodeString(bs64TimeStamp)
	if nil != err {
		return
	}
	timeStamp, err := strconv.Atoi(string(timeStampStr))
	if nil != err || int64(timeStamp) < time.Now().Unix() {
		return
	}

	// save the data to structure instance
	gobData, err := base64.StdEncoding.DecodeString(bs64DataStr)
	if nil != err {
		return
	}
	buffP := bytes.NewBuffer(gobData)
	decoder := gob.NewDecoder(buffP)
	err = decoder.Decode(data)
	if nil != err {
		return
	}
	return true, nil

}

//Generate today's timestamp after N years
func NYearLaterTimeStamp(now time.Time, year int) int64 {
	t := time.Date(now.Year()+year,
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
		time.UTC)
	return t.Unix()
}
