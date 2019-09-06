package ApiRPC

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"

	"../MongoBind"

	conf "../Config"
	pb "../Protos"
)

type MongoData struct{}

func (obj *MongoData) SaveDelayedMessage(ctx context.Context, params *pb.DelayedMessage) (
	*pb.DelayedMessage, error) {

	return params, MongoBind.MongoSaveDelayedMessage(params.UserId, params.Message)
}

func (obj *MongoData) GetDelayedMessage(ctx context.Context, params *pb.DelayedMessage) (
	*pb.DelayedMessage, error) {

	data, err := MongoBind.MongoQueryDelayedMessage(params.UserId)
	params.MessageList = data
	return params, err
}

func (obj *MongoData) AddOneFriendId(ctx context.Context, params *pb.UserFriend) (*pb.UserFriend, error) {
	return params, MongoBind.MongoAddFriendId(params.SelfId, params.FriendId)
}

func (obj *MongoData) DelOneFriendId(ctx context.Context, params *pb.UserFriend) (*pb.UserFriend, error) {
	return params, MongoBind.MongoDelFriendId(params.SelfId, params.FriendId)
}

func (obj *MongoData) GetAllFriendId(ctx context.Context, params *pb.UserFriend) (*pb.UserFriend, error) {
	data, err := MongoBind.MongoQueryFriendsId(params.SelfId)
	params.FriendIdList = data
	return params, err
}

func (obj *MongoData) AddOneFriendToBlacklist(ctx context.Context, params *pb.UserBlacklist) (
	*pb.UserBlacklist, error) {
	return params, MongoBind.MongoBlackListAdd(params.SelfId, params.FriendId)
}

func (obj *MongoData) DelOneFriendFromBlacklist(ctx context.Context, params *pb.UserBlacklist) (
	*pb.UserBlacklist, error) {
	return params, MongoBind.MongoBlackListDel(params.SelfId, params.FriendId)
}

func (obj *MongoData) GetBlacklistOfUser(ctx context.Context, params *pb.UserBlacklist) (
	*pb.UserBlacklist, error) {

	data, err := MongoBind.MongoQueryBlackList(params.SelfId)
	params.FriendIdList = data
	return params, err
}

func (obj *MongoData) AddUserToGroupChat(ctx context.Context, params *pb.GroupChatInfo) (
	*pb.GroupChatInfo, error) {

	return params, MongoBind.MongoGroupChatAddUser(params.Id, params.UserId)
}

func (obj *MongoData) DelUserFromGroupChat(ctx context.Context, params *pb.GroupChatInfo) (
	*pb.GroupChatInfo, error) {

	return params, MongoBind.MongoGroupChatDelUser(params.Id, params.UserId)
}

func (obj *MongoData) GetUsersOfGroupChat(ctx context.Context, params *pb.GroupChatInfo) (
	*pb.GroupChatInfo, error) {
	data, err := MongoBind.MongoQueryGroupChatUsers(params.Id)
	params.UserIdList = data
	return params, err
}

func (obj *MongoData) AddUserToSubscription(ctx context.Context, params *pb.SubscriptionInfo) (
	*pb.SubscriptionInfo, error) {

	return params, MongoBind.MongoSubscriptionAddUser(params.Id, params.UserId)
}

func (obj *MongoData) DelUserFromSubscription(ctx context.Context, params *pb.SubscriptionInfo) (
	*pb.SubscriptionInfo, error) {

	return params, MongoBind.MongoSubscriptionDelUser(params.Id, params.UserId)
}

func (obj *MongoData) GetUsersOfSubscription(ctx context.Context, params *pb.SubscriptionInfo) (
	*pb.SubscriptionInfo, error) {

	data, err := MongoBind.MongoQuerySubscriptionUsers(params.Id)
	params.UserIdList = data
	return params, err
}

func StartMongoDataRPCServer() {
	// Start the gRPC server for MySQL data operation.

	// using CA TSL authentication
	caOption := getCAOption()

	// get an interceptor server option for Unary-Unary handler
	unaryOption := getUnaryInterceptorOption()

	server := grpc.NewServer(unaryOption, caOption)

	pb.RegisterMongoBindServiceServer(server, &MongoData{})

	log.Println(":::Start Mongo Data Layer gRPC Server")
	listener, err := net.Listen("tcp", conf.MongoDataRPCServerAddress)
	if nil != err {
		log.Fatalf("Start gRPC server error: %s", err.Error())
	}
	err = server.Serve(listener)
	if nil != err {
		log.Fatal(err)
	}

}
