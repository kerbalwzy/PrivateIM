package ApiRPC

import (
	conf "../Config"
	"../MySQLBind"
	pb "../Protos"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
)

var (
	ParamsErr      = errors.New("bad request because of the wrong params value")
	CtxCanceledErr = errors.New("the client canceled or connection time out")
)

type MySQLData struct{}

// Functions for operate the basic information of the user
func (obj *MySQLData) NewOneUser(ctx context.Context,
	params *pb.UserBasicInfo) (*pb.UserBasicInfo, error) {

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
		userList.Data = append(userList.Data, user)
	}
	return userList, nil
}

// Updating the name, mobile, gender of the target user, which found by id.
func (obj *MySQLData) PutUserBasicById(ctx context.Context,
	params *pb.UpdateUserParams) (*pb.UserBasicInfo, error) {

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

// Translate the user information from TempUserBasic to UserBasicInfo
// Using the different struct between MySQLBind and gRPCpb for reduce
// the degree of coupling. Because the CreateTime saved as CCT time zone
// in MySQL, it not need to translate the time zone.
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

// Functions for operate avatar and qrCode file name of the user
func (obj *MySQLData) GetUserAvatarById(ctx context.Context,
	params *pb.UserAvatar) (*pb.UserAvatar, error) {

	var err error
	params.Avatar, err = MySQLBind.SelectUserAvatarById(params.Id)
	if nil != err {
		return nil, err
	}

	return params, nil
}

func (obj *MySQLData) PutUserAvatarById(ctx context.Context,
	params *pb.UserAvatar) (*pb.UserAvatar, error) {

	err := MySQLBind.UpdateUserAvatarById(params.Id, params.Avatar)
	if nil != err {
		return nil, err
	}
	return params, nil
}

func (obj *MySQLData) GetUserQRCodeById(ctx context.Context,
	params *pb.UserQRCode) (*pb.UserQRCode, error) {

	var err error
	params.QrCode, err = MySQLBind.SelectUserQRCodeById(params.Id)
	if nil != err {
		return nil, err
	}
	return params, nil
}

func (obj *MySQLData) PutUserQRCodeById(ctx context.Context,
	params *pb.UserQRCode) (*pb.UserQRCode, error) {

	err := MySQLBind.UpdateUserQRCode(params.Id, params.QrCode)
	if nil != err {
		return nil, err
	}
	return params, nil
}

// Functions for operate friendship record data between the user.
func (obj *MySQLData) AddOneNewFriend(ctx context.Context,
	params *pb.Friendship) (*pb.Friendship, error) {

	err := MySQLBind.InsertOneNewFriend(params.SelfId, params.FriendId,
		params.FriendNote)
	if nil != err {
		return nil, err
	}
	return params, nil

}

func (obj *MySQLData) PutOneFriendNote(ctx context.Context,
	params *pb.Friendship) (*pb.Friendship, error) {

	err := MySQLBind.UpdateOneFriendNote(params.SelfId, params.FriendId,
		params.FriendNote)
	if nil != err {
		return nil, err
	}
	return params, nil

}

func (obj *MySQLData) AcceptOneNewFriend(ctx context.Context,
	params *pb.Friendship) (*pb.Friendship, error) {

	err := MySQLBind.UpdateAcceptNewFriend(params.SelfId, params.FriendId,
		params.FriendNote, params.IsAccept)
	if nil != err {
		return nil, err
	}
	return params, nil
}

func (obj *MySQLData) PutFriendBlacklist(ctx context.Context,
	params *pb.Friendship) (*pb.Friendship, error) {

	err := MySQLBind.UpdateFriendBlacklist(params.SelfId, params.FriendId,
		params.IsBlack)
	if nil != err {
		return nil, err
	}
	return params, nil
}

func (obj *MySQLData) DeleteOneFriend(ctx context.Context,
	params *pb.Friendship) (*pb.Friendship, error) {

	err := MySQLBind.DeleteOneFriend(params.SelfId, params.FriendId)
	if nil != err {
		return nil, err
	}
	return params, nil
}

func (obj *MySQLData) GetFriendshipInfo(ctx context.Context,
	params *pb.QueryFriendsParams) (*pb.FriendshipList, error) {

	dataList, err := MySQLBind.SelectFriendsRelates(params.SelfId)
	if nil != err {
		return nil, err
	}

	friendshipList := new(pb.FriendshipList)
	for _, data := range dataList {
		temp := &pb.Friendship{
			SelfId:     data.SelfId,
			FriendId:   data.FriendId,
			FriendNote: data.FriendNote,
			IsAccept:   data.IsAccept,
			IsBlack:    data.IsBlack,
			IsDelete:   data.IsDelete}

		friendshipList.Data = append(friendshipList.Data, temp)
	}
	return friendshipList, nil
}

func (obj *MySQLData) GetFriendsBasicInfo(ctx context.Context,
	params *pb.QueryFriendsParams) (*pb.FriendsBasicInfoList, error) {

	dataList, err := MySQLBind.SelectFriendsInfo(params.SelfId)
	if nil != err {
		return nil, err
	}

	result := new(pb.FriendsBasicInfoList)
	for _, data := range dataList {
		temp := &pb.FriendsBasicInfo{
			FriendId: data.FriendId,
			Name:     data.Name,
			Email:    data.Email,
			Mobile:   data.Mobile,
			Gender:   int32(data.Gender),
			Note:     data.Note,
			IsBlack:  data.IsBlack}
		result.Data = append(result.Data, temp)
	}
	return result, nil
}

// Check the client if canceled the calling or connection is time out
func checkCtxCanceled(ctx context.Context) error {
	if ctx.Err() == context.Canceled {
		return CtxCanceledErr
	}
	return nil
}

// Start the gRPC server for MySQL data operation.
func StartMySQLgRPCServer() {
	// using CA TSL authentication
	cert, err := tls.LoadX509KeyPair(conf.DataLayerSrvCAServerPem, conf.DataLayerSrvCAServerKey)
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(conf.DataLayerSrvCAPem)
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	caOption := grpc.Creds(c)

	// new an interceptor, similar to middleware
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		// check the call if canceled by client or time out before every handler
		err = checkCtxCanceled(ctx)
		if err != nil {
			return
		}
		// continue the handel
		return handler(ctx, req)
	}
	// add the interceptor for every Unary-Unary handler
	unaryOption := grpc.UnaryInterceptor(interceptor)

	server := grpc.NewServer(unaryOption, caOption)

	pb.RegisterMySQLBindServiceServer(server, &MySQLData{})

	log.Println(":::Start MySQL Data Layer gRPC Server")
	listener, err := net.Listen("tcp", conf.MySQLDataRPCServerAddress)
	if nil != err {
		log.Fatalf("Start gRPC server error: %s", err.Error())
	}
	err = server.Serve(listener)
	if nil != err {
		log.Fatal(err)
	}
}
