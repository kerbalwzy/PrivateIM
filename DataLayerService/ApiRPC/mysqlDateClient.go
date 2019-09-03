package ApiRPC

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	conf "../Config"
	pb "../Protos"
)

var (
	theClient pb.MySQLBindServiceClient
)

func init() {
	conn, err := grpc.Dial(conf.MySQLDataRPCServerAddress, grpc.WithInsecure())
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

func SaveOneNewUser(name, email, mobile, password string,
	gender int) (*pb.UserBasicInfo, error) {
	client := getClient()
	params := &pb.UserBasicInfo{
		Name: name, Email: email, Mobile: mobile,
		Password: password, Gender: int32(gender)}
	return client.NewOneUser(getTimeOutCtx(3), params)
}

func GetUserById(id int64) (*pb.UserBasicInfo, error) {
	client := getClient()
	params := &pb.QueryUserParams{Id: id, FilterField: pb.QueryUserField_ById}
	return client.GetUserById(getTimeOutCtx(3), params)
}

func GetUserByEmail(email string) (*pb.UserBasicInfo, error) {
	client := getClient()
	params := &pb.QueryUserParams{Email: email, FilterField: pb.QueryUserField_ByEmail}
	return client.GetUserByEmail(getTimeOutCtx(3), params)
}

func GetUsersByName(name string) (*pb.UserBasicInfoList, error) {
	client := getClient()
	params := &pb.QueryUserParams{Name: name, FilterField: pb.QueryUserField_ByName}
	return client.GetUsersByName(getTimeOutCtx(3), params)
}

func PutUserBasicById(name, mobile string, gender int, id int64) (*pb.UserBasicInfo, error) {
	client := getClient()
	params := &pb.UpdateUserParams{Id: id, Name: name, Mobile: mobile,
		Gender: int32(gender), UpdateField: pb.UpdateUserField_NameMobileGender}
	return client.PutUserBasicById(getTimeOutCtx(3), params)
}

func PutUserPasswordById(password string, id int64) (*pb.UserBasicInfo, error) {
	client := getClient()
	params := &pb.UpdateUserParams{Password: password, Id: id,
		UpdateField: pb.UpdateUserField_Password}
	return client.PutUserPasswordById(getTimeOutCtx(3), params)
}

func PutUserPasswordByEmail(password, email string) (*pb.UserBasicInfo, error) {
	client := getClient()
	params := &pb.UpdateUserParams{Password: password, Email: email,
		UpdateField: pb.UpdateUserField_Password}
	return client.PutUserPasswordByEmail(getTimeOutCtx(3), params)
}

func GetUserAvatarById(id int64) (*pb.UserAvatar, error) {
	client := getClient()
	params := &pb.UserAvatar{Id: id}
	return client.GetUserAvatarById(getTimeOutCtx(3), params)
}

func PutUserAvatarById(avatar string, id int64) (*pb.UserAvatar, error) {
	client := getClient()
	params := &pb.UserAvatar{Avatar: avatar, Id: id}
	return client.PutUserAvatarById(getTimeOutCtx(3), params)
}

func GetUserQRCodeById(id int64) (*pb.UserQRCode, error) {
	client := getClient()
	params := &pb.UserQRCode{Id: id}
	return client.GetUserQRCodeById(getTimeOutCtx(3), params)
}

func PutUserQRCodeById(qrCode string, id int64) (*pb.UserQRCode, error) {
	client := getClient()
	params := &pb.UserQRCode{QrCode: qrCode, Id: id}
	return client.PutUserQRCodeById(getTimeOutCtx(3), params)
}
