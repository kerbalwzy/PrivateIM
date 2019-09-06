package ApiRPC

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"

	conf "../Config"
	pb "../Protos"
)

var (
	mongoDataClient pb.MongoBindServiceClient
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

	conn, err := grpc.Dial(conf.MongoDataRPCServerAddress, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatal(err.Error())
	}

	mongoDataClient = pb.NewMongoBindServiceClient(conn)

}

// Return the client of RPC call.
// Build this function because connection pools may need to be used in the future.
func getMongoDataClient() pb.MongoBindServiceClient {
	return mongoDataClient
}

func SaveDelayedMessage(id int64, message []byte) (*pb.DelayedMessage, error) {
	client := getMongoDataClient()
	params := &pb.DelayedMessage{UserId: id, Message: message}
	return client.SaveDelayedMessage(getTimeOutCtx(3), params)
}

func GetDelayedMessage(id int64) (*pb.DelayedMessage, error) {
	client := getMongoDataClient()
	params := &pb.DelayedMessage{UserId: id}
	return client.GetDelayedMessage(getTimeOutCtx(3), params)
}

func AddOneFriendId(selfId, friendId int64) (*pb.UserFriend, error) {
	client := getMongoDataClient()
	params := &pb.UserFriend{SelfId: selfId, FriendId: friendId}
	return client.AddOneFriendId(getTimeOutCtx(3), params)
}

func DelOneFriendId(selfId, friendId int64) (*pb.UserFriend, error) {
	client := getMongoDataClient()
	params := &pb.UserFriend{SelfId: selfId, FriendId: friendId}
	return client.DelOneFriendId(getTimeOutCtx(3), params)
}

func GetAllFriendId(id int64) (*pb.UserFriend, error) {
	client := getMongoDataClient()
	params := &pb.UserFriend{SelfId: id}
	return client.GetAllFriendId(getTimeOutCtx(3), params)
}

func AddOneFriendToBlacklist(selfId, friendId int64) (*pb.UserBlacklist, error) {
	client := getMongoDataClient()
	params := &pb.UserBlacklist{SelfId: selfId, FriendId: friendId}
	return client.AddOneFriendToBlacklist(getTimeOutCtx(3), params)
}

func DelOneFriendFromBlacklist(selfId, friendId int64) (*pb.UserBlacklist, error) {
	client := getMongoDataClient()
	params := &pb.UserBlacklist{SelfId: selfId, FriendId: friendId}
	return client.DelOneFriendFromBlacklist(getTimeOutCtx(3), params)
}

func GetBlacklistOfUser(id int64) (*pb.UserBlacklist, error) {
	client := getMongoDataClient()
	params := &pb.UserBlacklist{SelfId: id}
	return client.GetBlacklistOfUser(getTimeOutCtx(3), params)
}

func AddUserToGroupChat(groupId, userId int64) (*pb.GroupChatInfo, error) {
	client := getMongoDataClient()
	params := &pb.GroupChatInfo{Id: groupId, UserId: userId}
	return client.AddUserToGroupChat(getTimeOutCtx(3), params)

}

func DelUserFromGroupChat(groupId, userId int64) (*pb.GroupChatInfo, error) {
	client := getMongoDataClient()
	params := &pb.GroupChatInfo{Id: groupId, UserId: userId}
	return client.DelUserFromGroupChat(getTimeOutCtx(3), params)
}

func GetUsersOfGroupChat(groupId int64) (*pb.GroupChatInfo, error) {
	client := getMongoDataClient()
	params := &pb.GroupChatInfo{Id: groupId}
	return client.GetUsersOfGroupChat(getTimeOutCtx(3), params)
}

func AddUserToSubscription(subsId, userId int64) (*pb.SubscriptionInfo, error) {
	client := getMongoDataClient()
	params := &pb.SubscriptionInfo{Id: subsId, UserId: userId}
	return client.AddUserToSubscription(getTimeOutCtx(3), params)
}

func DelUserFromSubscription(subsId, userId int64) (*pb.SubscriptionInfo, error) {
	client := getMongoDataClient()
	params := &pb.SubscriptionInfo{Id: subsId, UserId: userId}
	return client.DelUserFromSubscription(getTimeOutCtx(3), params)
}

func GetUsersOfSubscription(subsId int64) (*pb.SubscriptionInfo, error) {
	client := getMongoDataClient()
	params := &pb.SubscriptionInfo{Id: subsId}
	return client.GetUsersOfSubscription(getTimeOutCtx(3), params)
}
