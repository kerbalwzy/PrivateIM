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

// Check auth token by call UserCenter gRPC server
func CheckAuthToken(token string) (int64, error) {
	authToken := &AuthToken{Token: token, ClientInfo: MyClientInfo}
	ctx := context.Background()
	context.WithValue(ctx, "client_profile", MyClientInfo)
	result, err := AuthClient.CheckAuthToken(ctx, authToken)
	if nil != err {
		log.Println(err)
		return 0, err
	}
	return result.UserId, nil
}
