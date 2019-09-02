package ApiRPC

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "../Protos"
)

const (
	MySQLDataServerAddr = "localhost:23331"
)

var (
	theClient pb.MySQLBindServiceClient
)

func init() {
	conn, err := grpc.Dial(MySQLDataServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}

	theClient = pb.NewMySQLBindServiceClient(conn)
}

// Return the client of RPC call.
// Build this function because connection pools may need to be used in the future.
func getClient() pb.MySQLBindServiceClient {
	return theClient
}

// Return a context instance with deadline
func getTimeOutCtx(expire time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), expire*time.Second)
	return ctx
}

func QueryUserById(id int64) (*pb.UserBasicInfo, error) {
	client := getClient()
	params := &pb.QueryUserParams{Id: id, FilterField: pb.QueryField_ById}
	return client.GetUserById(getTimeOutCtx(3), params)
}

func QueryUserByEmail(email string) (*pb.UserBasicInfo, error) {
	client := getClient()
	params := &pb.QueryUserParams{Email: email, FilterField: pb.QueryField_ByEmail}
	return client.GetUserByEmail(getTimeOutCtx(3), params)
}

func QueryUsersByName(name string) (*pb.UserBasicInfoList, error) {
	client := getClient()
	params := &pb.QueryUserParams{Name: name, FilterField: pb.QueryField_ByName}
	return client.GetUsersByName(getTimeOutCtx(3), params)
}

func NewOneUser(name, email, mobile, password string,
	gender int) (*pb.UserBasicInfo, error) {
	client := getClient()
	params := &pb.UserBasicInfo{
		Name: name, Email: email, Mobile: mobile,
		Password: password, Gender: int32(gender)}
	return client.NewOneUser(getTimeOutCtx(3), params)
}
