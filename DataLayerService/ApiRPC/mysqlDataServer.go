package ApiRPC

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"

	conf "../Config"
	"../MySQLBind"
	pb "../Protos"
)

var (
	ParamsErr      = errors.New("bad request because of the wrong params value")
	CtxCanceledErr = errors.New("the client canceled or connection time out")
)

type MySQLData struct{}

func (obj *MySQLData) NewOneUser(ctx context.Context,
	params *pb.UserBasicInfo) (*pb.UserBasicInfo, error) {

	err := checkCtxCanceled(ctx)
	if nil != err {
		return nil, err
	}
	data, err := MySQLBind.InsertOneUser(params.Name, params.Email,
		params.Mobile, params.Password, int(params.Gender))
	if nil != err {
		return nil, err
	}
	user := initUserBasic(data)
	return user, nil
}

func (obj *MySQLData) GetUserById(ctx context.Context,
	params *pb.QueryUserParams) (*pb.UserBasicInfo, error) {

	err := checkCtxCanceled(ctx)
	if nil != err {
		return nil, err
	}
	if params.FilterField != pb.QueryUserField_ById {
		return nil, ParamsErr
	}
	data, err := MySQLBind.SelectUserById(params.Id)
	if nil != err {
		return nil, err
	}
	user := initUserBasic(data)
	return user, nil

}

func (obj *MySQLData) GetUserByEmail(ctx context.Context,
	params *pb.QueryUserParams) (*pb.UserBasicInfo, error) {

	err := checkCtxCanceled(ctx)
	if nil != err {
		return nil, err
	}
	if params.FilterField != pb.QueryUserField_ByEmail {
		return nil, ParamsErr
	}
	data, err := MySQLBind.SelectUserByEmail(params.Email)
	if nil != err {
		return nil, err
	}
	user := initUserBasic(data)
	return user, nil

}

func (obj *MySQLData) GetUsersByName(ctx context.Context,
	params *pb.QueryUserParams) (*pb.UserBasicInfoList, error) {
	err := checkCtxCanceled(ctx)
	if nil != err {
		return nil, err
	}
	if params.FilterField != pb.QueryUserField_ByName {
		return nil, ParamsErr
	}
	data, err := MySQLBind.SelectUsersByName(params.Name)
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

// Updating the name, mobile, gender of the target user, which found by id.
func (obj *MySQLData) PutUserBasicById(ctx context.Context,
	params *pb.UpdateUserParams) (*pb.UserBasicInfo, error) {

	err := checkCtxCanceled(ctx)
	if nil != err {
		return nil, err
	}

	if params.UpdateField != pb.UpdateUserField_NameMobileGender {
		return nil, ParamsErr
	}

	data, err := MySQLBind.UpdateUserBasicById(params.Name, params.Mobile,
		int(params.Gender), params.Id)
	if nil != err {
		return nil, err
	}

	user := initUserBasic(data)
	return user, nil
}

func (obj *MySQLData) PutUserPasswordById(ctx context.Context,
	params *pb.UpdateUserParams) (*pb.UserBasicInfo, error) {

	err := checkCtxCanceled(ctx)
	if nil != err {
		return nil, err
	}
	if params.UpdateField != pb.UpdateUserField_Password {
		return nil, ParamsErr
	}

	data, err := MySQLBind.UpdateUserPasswordById(params.Password, params.Id)
	if nil != err {
		return nil, err
	}
	user := initUserBasic(data)
	return user, nil

}

func (obj *MySQLData) PutUserPasswordByEmail(ctx context.Context,
	params *pb.UpdateUserParams) (*pb.UserBasicInfo, error) {

	err := checkCtxCanceled(ctx)
	if nil != err {
		return nil, err
	}
	if params.UpdateField != pb.UpdateUserField_Password {
		return nil, ParamsErr
	}

	data, err := MySQLBind.UpdateUserPasswordByEmail(params.Password, params.Email)
	if nil != err {
		return nil, err
	}
	user := initUserBasic(data)
	return user, nil
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
// the degree of coupling. Because the CreateTime saved as CCT time zone,
// it not need to translate the time zone.
func initUserBasic(data *MySQLBind.TempUserBasic) *pb.UserBasicInfo {
	return &pb.UserBasicInfo{
		Id:         data.Id,
		Name:       data.Name,
		Email:      data.Email,
		Mobile:     data.Mobile,
		Password:   data.Password,
		Gender:     int32(data.Gender),
		CreateTime: data.CreateTime.Format(conf.TimeDisplayFormat)}
}

// Start the gRPC server for MySQL data operation.
// Using CA TSL authentication
func StartMySQLgRPCServer() {
	listener, err := net.Listen("tcp", conf.MySQLDataRPCServerAddress)
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
