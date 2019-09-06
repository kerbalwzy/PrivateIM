package MongoBind

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	conf "../Config"
)

var (
	MongoClient *mongo.Client
	MsgCenterDB *mongo.Database

	DelayedMessageColl *mongo.Collection
	UserFriendsColl    *mongo.Collection
	UserBlackColl      *mongo.Collection

	GroupChatsColl    *mongo.Collection
	SubscriptionsColl *mongo.Collection
)

func init() {
	var err error
	MongoClient, err = mongo.Connect(getTimeOutCtx(10),
		options.Client().ApplyURI(conf.MsgDbMongoURI))
	if nil != err {
		log.Fatal(err)
	}
	err = MongoClient.Ping(getTimeOutCtx(3), readpref.Primary())
	if nil != err {
		log.Fatal(err)
	}

	MsgCenterDB = MongoClient.Database("MsgCenter")

	DelayedMessageColl = MsgCenterDB.Collection("DelayedMessage")
	UserFriendsColl = MsgCenterDB.Collection("UserFriends")
	UserBlackColl = MsgCenterDB.Collection("UserBlackList")

	GroupChatsColl = MsgCenterDB.Collection("GroupChats")
	SubscriptionsColl = MsgCenterDB.Collection("Subscriptions")

}

// Return a context instance with deadline
func getTimeOutCtx(expire time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), expire*time.Second)
	return ctx
}

// Id is for user's id
type TempDelayedMessage struct {
	Id      int64    `bson:"_id"`
	Message [][]byte `bson:"message"`
}

//Save the message which sent failed because the target user is offline.
func MongoSaveDelayedMessage(id int64, data []byte) error {
	_, err := DelayedMessageColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"message": data}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: save DelayedMessage fail for user(%d), error detail: %s", id, err.Error())
		return err
	} else {
		log.Printf("DelayedMessage: save an message for user(%d)", id)
		return nil
	}
}

// Query the messages should be sent to current user.
func MongoQueryDelayedMessage(id int64) ([][]byte, error) {
	temp := new(TempDelayedMessage)
	err := DelayedMessageColl.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": id}).Decode(temp)
	if nil != err {
		log.Printf("Error: query DelayedMessage fail for user(%d), error detail: %s", id, err.Error())
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
	BlackList []int64 `bson:"blacklist"`
}

// Add a user'id into the blacklist of current user
func MongoBlackListAdd(srcId, dstId int64) error {
	_, err := UserBlackColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": srcId},
		bson.M{"$addToSet": bson.M{"blacklist": dstId}},
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
		bson.M{"$pull": bson.M{"blacklist": dstId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: move a friend out from blacklist fail for user(%d), error detail: %s", srcId, err.Error())
		return err
	}
	return nil
}

type TempGroupChat struct {
	Id      int64   `bson:"_id"`
	UsersId []int64 `bson:"users_id"`
}

// Add a user's id into the group chat
func MongoGroupChatAddUser(groupId, userId int64) error {
	_, err := GroupChatsColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": groupId},
		bson.M{"$addToSet": bson.M{"users_id": userId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: add user fail for groupChat(%d), error detail: %s", groupId, err.Error())
		return err
	}
	return nil
}

// Query the user's id of the group
func MongoQueryGroupChatUsers(groupId int64) ([]int64, error) {
	temp := new(TempGroupChat)
	err := GroupChatsColl.FindOne(getTimeOutCtx(3), bson.M{"_id": groupId}).Decode(temp)
	if nil != err {
		log.Printf("Error: query users fail for groupChat(%d), error detail: %s", groupId, err.Error())
		return nil, err
	}
	return temp.UsersId, nil
}

// Query the all groups information
func MongoQueryGroupChatAll() ([]TempGroupChat, error) {
	ctx := getTimeOutCtx(30)
	curs, err := GroupChatsColl.Find(ctx, bson.D{})
	if nil != err {
		log.Printf("Error: query all group chat information fail")
		return nil, err
	}
	defer curs.Close(ctx)
	data := make([]TempGroupChat, 0)
	for curs.Next(ctx) {
		temp := new(TempGroupChat)
		err := curs.Decode(temp)
		if nil != err {
			log.Printf("Error: query all group chat error, detail: %s", err.Error())
			continue
		}
		data = append(data, *temp)
	}
	return data, nil
}

// Move a user's id out from a group chat
func MongoGroupChatDelUser(groupId, userId int64) error {
	_, err := GroupChatsColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": groupId},
		bson.M{"$pull": bson.M{"users_id": userId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: remove user fail for groupChat(%d), error detail: %s", groupId, err.Error())
		return err
	}
	return nil
}

type TempSubscription struct {
	Id      int64   `bson:"_id"`
	UsersId []int64 `bson:"users_id"`
}

// Add a user's into the subscription
func MongoSubscriptionAddUser(subsId, userId int64) error {
	_, err := SubscriptionsColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": subsId},
		bson.M{"$addToSet": bson.M{"users_id": userId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: add user fail for subscription(%d), error detail: %s", subsId, err.Error())
		return err
	}
	return nil
}

// Query the user's id of the subscription
func MongoQuerySubscriptionUsers(subsId int64) ([]int64, error) {
	temp := new(TempSubscription)
	err := SubscriptionsColl.FindOne(getTimeOutCtx(3), bson.M{"_id": subsId}).Decode(temp)
	if nil != err {
		log.Printf("Error: query user fail for subscription(%d), error detail: %s", subsId, err.Error())
		return nil, err
	}
	return temp.UsersId, nil
}

// Query the all subscription information
func MongoQuerySubscriptionAll() ([]TempSubscription, error) {
	ctx := getTimeOutCtx(30)
	curs, err := SubscriptionsColl.Find(ctx, bson.D{})
	if nil != err {
		log.Printf("Error: query all subscription fail")
		return nil, err
	}
	defer curs.Close(ctx)
	data := make([]TempSubscription, 0)
	for curs.Next(ctx) {
		temp := new(TempSubscription)
		err := curs.Decode(temp)
		if nil != err {
			log.Printf("Error: query all subscription error, detail: %s", err.Error())
			continue
		}
		data = append(data, *temp)

	}
	return data, nil
}

// Move a user's id out from a subscription
func MongoSubscriptionDelUser(subsId, userId int64) error {
	_, err := SubscriptionsColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": subsId},
		bson.M{"$pull": bson.M{"users_id": userId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: remove user fail for subscription(%d), error detail: %s", subsId, err.Error())
	}
	return nil
}
