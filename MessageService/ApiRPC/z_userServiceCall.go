package ApiRPC

import (
	"../RpcClientPbs/userAuthPb"
	"strconv"
)

func CheckAuthToken(token string) (int64, error) {
	// todo test code used in separate development, need remove later
	tempId, _ := strconv.Atoi(token)
	return int64(tempId), nil

	// code to actually use
	param := &userAuthPb.Token{Value: token}
	data, err := GetUserAuthClient().CheckAuthToken(getTimeOutCtx(3), param)
	if nil != err {
		return -1, err
	}
	return data.Value, err
}
