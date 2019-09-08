package ApiRPC

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	conf "../Config"
	pb "../Protos"
)

var (
	mysqlDataClient pb.MySQLBindServiceClient
)

func init() {
	// Add CA TSL authentication data
	cert, err := tls.LoadX509KeyPair(conf.DataLayerSrvCAClientPem, conf.DataLayerSrvCAClientKey)
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
		ServerName:   "PrivateIM",
		RootCAs:      certPool,
	})

	conn, err := grpc.Dial(conf.MySQLDataRPCServerAddress, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatal(err.Error())
	}

	mysqlDataClient = pb.NewMySQLBindServiceClient(conn)
}

// Return the client of RPC call.
// Build this function because connection pools may need to be used in the future.
func getMySQLDataClient() pb.MySQLBindServiceClient {
	return mysqlDataClient
}

// Functions for operate the basic information of the user
func SaveOneNewUser(name, email, mobile, password string, gender int) (*pb.UserBasicInfo, error) {
	client := getMySQLDataClient()
	params := &pb.UserBasicInfo{
		Name: name, Email: email, Mobile: mobile,
		Password: password, Gender: int32(gender)}
	return client.NewOneUser(getTimeOutCtx(3), params)
}

func GetUserById(id int64) (*pb.UserBasicInfo, error) {
	client := getMySQLDataClient()
	params := &pb.QueryUserParams{Id: id, FilterField: pb.QueryUserField_ById}
	return client.GetUserById(getTimeOutCtx(3), params)
}

func GetUserByEmail(email string) (*pb.UserBasicInfo, error) {
	client := getMySQLDataClient()
	params := &pb.QueryUserParams{Email: email, FilterField: pb.QueryUserField_ByEmail}
	return client.GetUserByEmail(getTimeOutCtx(3), params)
}

// The detail data is saved in `pb.UserBasicInfoList.Data`
func GetUsersByName(name string) (*pb.UserBasicInfoList, error) {
	client := getMySQLDataClient()
	params := &pb.QueryUserParams{Name: name, FilterField: pb.QueryUserField_ByName}
	return client.GetUsersByName(getTimeOutCtx(3), params)
}

func PutUserBasicById(name, mobile string, gender int, id int64) (*pb.UserBasicInfo, error) {
	client := getMySQLDataClient()
	params := &pb.UpdateUserParams{Id: id, Name: name, Mobile: mobile,
		Gender: int32(gender), UpdateField: pb.UpdateUserField_NameMobileGender}
	return client.PutUserBasicById(getTimeOutCtx(3), params)
}

func PutUserPasswordById(password string, id int64) (*pb.UserBasicInfo, error) {
	client := getMySQLDataClient()
	params := &pb.UpdateUserParams{Password: password, Id: id,
		UpdateField: pb.UpdateUserField_Password}
	return client.PutUserPasswordById(getTimeOutCtx(3), params)
}

func PutUserPasswordByEmail(password, email string) (*pb.UserBasicInfo, error) {
	client := getMySQLDataClient()
	params := &pb.UpdateUserParams{Password: password, Email: email,
		UpdateField: pb.UpdateUserField_Password}
	return client.PutUserPasswordByEmail(getTimeOutCtx(3), params)
}

// Functions for operate avatar and qrCode file name of the user
func GetUserAvatarById(id int64) (*pb.UserAvatar, error) {
	client := getMySQLDataClient()
	params := &pb.UserAvatar{Id: id}
	return client.GetUserAvatarById(getTimeOutCtx(3), params)
}

func PutUserAvatarById(avatar string, id int64) (*pb.UserAvatar, error) {
	client := getMySQLDataClient()
	params := &pb.UserAvatar{Avatar: avatar, Id: id}
	return client.PutUserAvatarById(getTimeOutCtx(3), params)
}

func GetUserQRCodeById(id int64) (*pb.UserQRCode, error) {
	client := getMySQLDataClient()
	params := &pb.UserQRCode{Id: id}
	return client.GetUserQRCodeById(getTimeOutCtx(3), params)
}

func PutUserQRCodeById(qrCode string, id int64) (*pb.UserQRCode, error) {
	client := getMySQLDataClient()
	params := &pb.UserQRCode{QrCode: qrCode, Id: id}
	return client.PutUserQRCodeById(getTimeOutCtx(3), params)
}

// Functions for operate friendship record data between the user.
func AddOneNewFriend(selfId, friendId int64, note string) (*pb.Friendship, error) {
	client := getMySQLDataClient()
	params := &pb.Friendship{SelfId: selfId, FriendId: friendId, FriendNote:note}
	return client.AddOneNewFriend(getTimeOutCtx(3), params)
}

func PutOneFriendNote(selfId, friendId int64, note string) (*pb.Friendship, error) {
	client := getMySQLDataClient()
	params := &pb.Friendship{SelfId: selfId, FriendId: friendId, FriendNote: note}
	return client.PutOneFriendNote(getTimeOutCtx(3), params)
}

func AcceptOneNewFriend(selfId, friendId int64, note string, isAccept bool) (*pb.Friendship, error) {

	client := getMySQLDataClient()
	params := &pb.Friendship{
		SelfId:     selfId,
		FriendId:   friendId,
		FriendNote: note,
		IsAccept:   isAccept}
	return client.AcceptOneNewFriend(getTimeOutCtx(3), params)
}

func PutFriendBlacklist(selfId, friendId int64, isBlack bool) (*pb.Friendship, error) {
	client := getMySQLDataClient()
	params := &pb.Friendship{SelfId: selfId, FriendId: friendId, IsBlack: isBlack}
	return client.PutFriendBlacklist(getTimeOutCtx(3), params)
}

func DeleteOneFriend(selfId, friendId int64) (*pb.Friendship, error) {
	client := getMySQLDataClient()
	params := &pb.Friendship{SelfId: selfId, FriendId: friendId}
	return client.DeleteOneFriend(getTimeOutCtx(3), params)
}

// The detail data is saved in `pb.FriendshipList.Data`
func GetFriendshipInfo(selfId int64) (*pb.FriendshipList, error) {
	client := getMySQLDataClient()
	params := &pb.QueryFriendsParams{SelfId: selfId}
	return client.GetFriendshipInfo(getTimeOutCtx(3), params)
}

// The detail data is saved in `pb.FriendBasicInfoList.Data`
func GetFriendsBasicInfo(selfId int64) (*pb.FriendsBasicInfoList, error) {
	client := getMySQLDataClient()
	params := &pb.QueryFriendsParams{SelfId: selfId}
	return client.GetFriendsBasicInfo(getTimeOutCtx(3), params)
}
