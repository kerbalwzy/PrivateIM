package ApiRPC

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

const (
	UserCenterGPRCServerAddress = "127.0.0.1:23333"
)

var AuthClient UserCenterClient
var MyClientInfo = &ClientInfo{
	Host:    "127.0.0.1",
	Name:    "MsgCenter",
	AuthKey: "this is auth key for MsgCenter gRPC client"}

func init() {
	conn, err := grpc.Dial(UserCenterGPRCServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}

	AuthClient = NewUserCenterClient(conn)
}

// check auth token by call UserCenter gRPC server
func CheckAuthToken(token string) (bool, error) {
	authToken := &AuthToken{Token: token, ClientInfo: MyClientInfo}
	checkResult, err := AuthClient.CheckAuthToken(context.Background(), authToken)
	if nil != err {
		log.Println(err)
		return false, err
	}
	return checkResult.Ok, nil
}
