package ApiRPC

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"

	"../MySQLBind"
	pb "../Protos"
)

const (
	MySQLDataServerAddress = "0.0.0.0:23331"
)

var (
	ParamsErr      = errors.New("bad request because of the wrong params value")
	CtxCanceledErr = errors.New("the client canceled or connection time out")
)

type MySQLData struct{}

func (obj *MySQLData) NewOneUser(ctx context.Context,
	params *pb.UserBasicInfo) (*pb.UserBasicInfo, error) {
	if err := checkCtxCanceled(ctx); nil != err {
		return nil, err
	}
	data, err := MySQLBind.InsertOneUser(params.Name, params.Email,
		params.Password, int(params.Gender))
	if nil != err {
		return nil, err
	}
	tempCreateTime := data.CreateTime.Local().String()
	user := &pb.UserBasicInfo{Id: data.Id, Email: data.Email, Mobile: data.Mobile,
		Password: data.Password, Gender: int32(data.Gender), CreateTime: tempCreateTime}
	return user, nil
}

func (obj *MySQLData) GetUserById(ctx context.Context,
	params *pb.QueryUserParams) (*pb.UserBasicInfo, error) {
	if err := checkCtxCanceled(ctx); nil != err {
		return nil, err
	}
	if params.FilterField != pb.QueryField_ById {
		return nil, ParamsErr
	}
	data, err := MySQLBind.QueryUserById(params.Id)
	if nil != err {
		return nil, err
	}
	user := initUserBasic(data)
	return user, nil

}

func (obj *MySQLData) GetUserByEmail(ctx context.Context,
	params *pb.QueryUserParams) (*pb.UserBasicInfo, error) {
	if err := checkCtxCanceled(ctx); nil != err {
		return nil, err
	}
	if params.FilterField != pb.QueryField_ByEmail {
		return nil, ParamsErr
	}
	data, err := MySQLBind.QueryUserByEmail(params.Email)
	if nil != err {
		return nil, err
	}
	user := initUserBasic(data)
	return user, nil

}

func (obj *MySQLData) GetUsersByName(ctx context.Context,
	params *pb.QueryUserParams) (*pb.UserBasicInfoList, error) {
	if err := checkCtxCanceled(ctx); nil != err {
		return nil, err
	}
	if params.FilterField != pb.QueryField_ByName {
		return nil, ParamsErr
	}
	data, err := MySQLBind.QueryUsersByName(params.Name)
	if nil != err {
		return nil, err
	}
	userList := new(pb.UserBasicInfoList)
	for _, elem := range data {
		user := initUserBasic(elem)
		userList.Users = append(userList.Users, user)
	}
	return userList, nil
}

// Check the client if canceled the calling or connection is time out
func checkCtxCanceled(ctx context.Context) error {
	if ctx.Err() == context.Canceled {
		return CtxCanceledErr
	}
	return nil
}

// Translate the user information from TempUserBasic to UserBasicInfo
// Using the different struct between MySQLBind and gRPCpb for reduce
// the degree of coupling
func initUserBasic(data *MySQLBind.TempUserBasic) *pb.UserBasicInfo {
	return &pb.UserBasicInfo{
		Id:         data.Id,
		Name:       data.Name,
		Email:      data.Email,
		Mobile:     data.Mobile,
		Password:   data.Password,
		Gender:     int32(data.Gender),
		CreateTime: data.CreateTime.String()}
}

func StartMySQLgRPCServer() {
	listener, err := net.Listen("tcp", MySQLDataServerAddress)
	if nil != err {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	pb.RegisterMySQLBindServiceServer(server, &MySQLData{})
	log.Println(":::Start MySQL Data Layer gRPC Server")
	err = server.Serve(listener)
	if nil != err {
		log.Fatal(err)
	}

}
