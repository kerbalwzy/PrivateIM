package DataLayer

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongoDBURI = "mongodb://localhost:27017"
)

var (
	Mongo           *mongo.Client
	WaitSendMsgColl *mongo.Collection
	UserFriendsColl *mongo.Collection
	UserBlackColl   *mongo.Collection
)

func init() {
	var err error
	Mongo, err = mongo.Connect(getTimeOutCtx(10), options.Client().ApplyURI(MongoDBURI))
	if nil != err {
		log.Fatal(err)
	}
	err = Mongo.Ping(getTimeOutCtx(3), readpref.Primary())
	if nil != err {
		log.Fatal(err)
	}
	WaitSendMsgColl = Mongo.Database("MsgCenter").Collection("WaitSendMsg")
	UserFriendsColl = Mongo.Database("MsgCenter").Collection("UserFriends")
	UserBlackColl = Mongo.Database("MsgCenter").Collection("UserBlackList")

}

func getTimeOutCtx(expire time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), expire*time.Second)
	return ctx
}

type TempWaitSendMsg struct {
	Id      int64    `bson:"_id"`
	Message [][]byte `bson:"message"`
}

//Save the message which sent failed because the target user is offline.
func MongoSaveWaitSendMessage(id int64, data []byte) error {
	_, err := WaitSendMsgColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"message": data}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: save WaitSendMessage fail for user(%d), error detail: %s", id, err.Error())
		return err
	} else {
		log.Printf("WaitSendMessage: save an message for user(%d)", id)
		return nil
	}
}

// Query the messages should be sent to current user.
func MongoQueryWaitSendMessage(id int64) ([][]byte, error) {
	temp := new(TempWaitSendMsg)
	err := WaitSendMsgColl.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": id}).Decode(temp)
	if nil != err {
		log.Printf("Error: query WaitSendMessage fail for user(%d), error detail: %s", id, err.Error())
		return nil, err
	}
	return temp.Message, nil
}

type TempFriends struct {
	Id      int64   `bson:"_id"`
	Friends []int64 `bson:"friends"`
}

// Add a friend's id for the current user.
func MongoAddFriendId(srcId, dstId int64) error {
	_, err := UserFriendsColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": srcId},
		bson.M{"$addToSet": bson.M{"friends": dstId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: add friends fail for user(%d), error detail: %s", srcId, err.Error())
		return err
	}
	return nil
}

// Query the id of friends of the current user
func MongoQueryFriendsId(id int64) ([]int64, error) {
	temp := new(TempFriends)
	err := UserFriendsColl.FindOne(getTimeOutCtx(3), bson.M{"_id": id}).Decode(temp)
	if nil != err {
		log.Printf("Error: query friends fail for user(%d), error detail: %s", id, err.Error())
		return nil, err
	}
	return temp.Friends, nil
}

// Remove a friend's id for the current user.
func MongoDelFriendId(srcId, dstId int64) error {
	_, err := UserFriendsColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": srcId},
		bson.M{"$pull": bson.M{"friends": dstId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: remove a friend fail for user(%d), error detail: %s", srcId, err.Error())
		return err
	}
	return nil
}

type TempBlackList struct {
	Id        int64   `bson:"_id"`
	BlackList []int64 `bson:"blackList"`
}

// Add a user'id into the blacklist of current user
func MongoBlackListAdd(srcId, dstId int64) error {
	_, err := UserBlackColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": srcId},
		bson.M{"$addToSet": bson.M{"blackList": dstId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: add friends fail for user(%d), error detail: %s", srcId, err.Error())
		return err
	}
	return nil
}

// Query the id of the friends who marked black by current user
func MongoQueryBlackList(id int64) ([]int64, error) {
	temp := new(TempBlackList)
	err := UserBlackColl.FindOne(getTimeOutCtx(3), bson.M{"_id": id}).Decode(temp)
	if nil != err {
		log.Printf("Error: query blacklist fail for user(%d), error detail: %s", id, err.Error())
		return nil, err
	}
	return temp.BlackList, nil
}

// Move a user's id out from the blacklist of current user
func MongoBlackListDel(srcId, dstId int64) error {
	_, err := UserBlackColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": srcId},
		bson.M{"$pull": bson.M{"blackList": dstId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: move a friend out from blacklist fail for user(%d), error detail: %s", srcId, err.Error())
		return err
	}
	return nil
}
