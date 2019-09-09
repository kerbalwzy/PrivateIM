package ApiRPC

import (
	"google.golang.org/grpc"
	"log"

	conf "../Config"
	pb "../Protos"
)

var userAuthClient pb.UserAuthClient

func init() {
	// don't use TSL
	conn, err := grpc.Dial(conf.UserAuthRPCServerAddress, grpc.WithInsecure())
	if nil != err {
		log.Fatal(err)
	}
	userAuthClient = pb.NewUserAuthClient(conn)
}

// The `pb.UserAuthClient` is a interface type, so there is passed by reference
func getClient() pb.UserAuthClient {
	return userAuthClient
}

func CheckAuthToken(token string) (*pb.TokenCheckResult, error) {
	client := getClient()
	params := &pb.AuthToken{Data: token}
	return client.CheckAuthToken(getTimeOutCtx(3), params)
}
