package MongoBind

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	conf "../Config"
)

var (
	CollDelayMessage,           // the message for the user when him was offline. should be sent and delete when the user online.
	CollUserChatHistory,        // the history of the message which between two user.
	CollGroupChatHistory,       // the history of the message which belong the group chat.
	CollSubscriptionMsgHistory, // the history of the message which belong the subscription.

	CollUserFriends,       // the id of user's friends, include the blacklist.
	CollUserGroupChats,    // the id of the group chat which the user was joined.
	CollUserSubscriptions, // the id of the subscription which the user was followed

	CollGroupUsers, // the id of users whom are belong to the group chat.
	CollSubscriptionUsers *mongo.Collection // the id of users whom are followed the subscription.
)

func init() {
	conn, err := mongo.Connect(getTimeOutCtx(3), options.Client().ApplyURI(conf.MongoDBURI))
	if nil != err {
		log.Fatal(err)
	}
	err = conn.Ping(getTimeOutCtx(3), readpref.Primary())
	if nil != err {
		log.Fatal(err)
	}
	db := conn.Database(conf.MongoDBName)
	CollDelayMessage = db.Collection(conf.CollDelayMessageName)
	CollUserChatHistory = db.Collection(conf.CollChatHistoryName)
	CollGroupChatHistory = db.Collection(conf.CollGroupChatHistoryName)
	CollSubscriptionMsgHistory = db.Collection(conf.CollSubscriptionMsgHistoryName)

	CollUserFriends = db.Collection(conf.CollUserFriendsName)
	CollUserGroupChats = db.Collection(conf.CollUserGroupChatsName)
	CollUserSubscriptions = db.Collection(conf.CollUserSubscriptions)

	CollGroupUsers = db.Collection(conf.CollGroupChatUsersName)
	CollSubscriptionUsers = db.Collection(conf.CollSubscriptionUsersName)

}

// Update options:
var (
	// when the document want be update not existed, insert new one.
	upsertTrueOption = options.Update().SetUpsert(true)
)

// Custom errors:
var (
	ErrFoundCount             = errors.New("the found count should be 0 or 1")
	ErrMessageHistoryNotFound = errors.New("the message history not found")
)

// Return a context instance with deadline
func getTimeOutCtx(expire time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), expire*time.Second)
	return ctx
}

// Return a number of today
func getTodayNum() int32 {
	year, month, day := time.Now().Date()
	dateStr := fmt.Sprintf("%d%02d%d", year, month, day)
	dateNum, _ := strconv.Atoi(dateStr)
	return int32(dateNum)
}

// Return a string joined by two user's id
func GetJoinUserId(userId1, userId2 int64) string {
	if userId2 > userId1 {
		return fmt.Sprintf("%d_%d", userId1, userId2)
	} else {
		return fmt.Sprintf("%d_%d", userId2, userId1)
	}
}

// information for history messages(message recorded for every date).
// Because there only have int32 and int64 in protocol buffers3, so the 'Date' defined with int32
type historyMessage struct {
	Date     int32    `json:"date"`
	Messages [][]byte `json:"messages"`
}

// Update the history message record for today.
// It will checking if  the document, which found by id in the target collection, had saved history message record
// for today. If the result is true, will append it, else add new one sub document for the message of today.
func updateHistoryMessageForToday(coll *mongo.Collection, id interface{}, message []byte) error {
	dateNum := getTodayNum()

	// checking if had saved chat history for today
	count, err := coll.CountDocuments(getTimeOutCtx(2),
		bson.M{"_id": id, "history": bson.M{"$elemMatch": bson.M{"date": bson.M{"$eq": dateNum}}}})
	if nil != err {
		return err
	}
	switch count {
	case 1:
		// had record for today
		_, err = coll.UpdateOne(getTimeOutCtx(3),
			bson.M{"_id": id, "history": bson.M{"$elemMatch": bson.M{"date": bson.M{"$eq": dateNum}}}},
			bson.M{"$push": bson.M{"history.$.messages": message}})
	case 0:
		// not had record for today
		_, err = coll.UpdateOne(getTimeOutCtx(3),
			bson.M{"_id": id},
			bson.M{"$push": bson.M{"history": bson.M{"date": dateNum, "messages": bson.A{message}}}},
			upsertTrueOption)
	default:
		return ErrFoundCount
	}
	return err
}

// Find all history message in the document which found by id in target collection.
func findAllHistoryMessageById(coll *mongo.Collection, id interface{}) *mongo.SingleResult {
	return coll.FindOne(getTimeOutCtx(3), bson.M{"_id": id})
}

// Find many history message  by date range in the document. The document was found by id in target collection.
// The 'startDate' and 'endDate' are the number value of the date. Because there only have 'int32' and 'int64'
// in protocol buffers3, so the type of 'startDate' and 'endDate' are 'int32'.
func findManyHistoryMessageByIdAndDateRange(coll *mongo.Collection, id interface{}, startDate, endDate int32) (
	*mongo.Cursor, error) {

	cur, err := coll.Aggregate(getTimeOutCtx(3),
		bson.A{
			bson.M{"$match": bson.M{"_id": id}},
			bson.M{"$project": bson.M{
				"_id": 1,
				"history": bson.M{
					"$filter": bson.M{
						"input": "$history",
						"as":    "message",
						"cond": bson.M{"$and": bson.A{
							bson.M{"$gte": bson.A{"$$message.date", startDate}},
							bson.M{"$lte": bson.A{"$$message.date", endDate}},},
						},
					},
				},
			}},
		})
	if nil != err {
		return nil, err
	}
	// check if have found document
	if !cur.Next(getTimeOutCtx(3)) {
		return nil, ErrMessageHistoryNotFound
	}
	return cur, nil
}

// Find many history message the a special date in the document. The document was found by id in target collection.
// The 'date' is the number value of the date. Because there only have 'int32' and 'int64' in protocol buffers3,
// so the type of 'date' is 'int32'.
func findManyHistoryMessageByIdAndDate(coll *mongo.Collection, id interface{}, date int32) (*mongo.Cursor, error) {
	cur, err := coll.Aggregate(getTimeOutCtx(3),
		bson.A{
			bson.M{"$match": bson.M{"_id": id}},
			bson.M{"$project": bson.M{
				"_id": 1,
				"history": bson.M{
					"$filter": bson.M{
						"input": "$history",
						"as":    "message",
						"cond":  bson.M{"$eq": bson.A{"$$message.date", date}},
					},
				},
			}}})
	if nil != err {
		return nil, err
	}
	// check if have found data
	if !cur.Next(getTimeOutCtx(3)) {
		return nil, ErrMessageHistoryNotFound
	}

	return cur, err
}

// ------------------------------------------------------------------------------------

/* delay_message document eg.:
{
	"_id": <user_id>,
	"messages": [
		<the message bytes data: eg. []byte("test message string")>,
		...
	]
}
*/
// information for delay message
type DocDelayMessage struct {
	Id      int64    `bson:"_id"`
	Message [][]byte `bson:"messages"`
}

// Update the delay message of the user. If the document not existed, will insert new one with the data.
func UpdateDelayMessage(userId int64, message []byte) error {
	_, err := CollDelayMessage.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": userId},
		bson.M{"$push": bson.M{"messages": message}},
		upsertTrueOption)
	return err
}

// Find and delete the delay message of the user from collection.
func FindAndDeleteDelayMessage(userId int64) (*DocDelayMessage, error) {
	temp := new(DocDelayMessage)
	err := CollDelayMessage.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": userId}).Decode(temp)
	if nil != err {
		return nil, err
	}
	return temp, nil
}

// ------------------------------------------------------------------------------------

/* user_chat_history document eg.:
{
	"_id" : <join_user_id eg.: 1_2>,
	"history": [
		{
			"date": <date number: eg. 20190303>,
			"messages": [
				<the message bytes data eg.: []byte("test message string")>,
				....
			]
		},
	...
	]
}
*/
// information for user chat history
type DocUserChatHistory struct {
	Id      string           `bson:"_id"`
	History []historyMessage `bson:"history"`
}

// Update the chat history between two users.
func UpdateUserChatHistoryByJoinId(joinUserId string, message []byte) error {
	return updateHistoryMessageForToday(CollUserChatHistory, joinUserId, message)
}

// Find all chat history between two users by their id.
func FindUserAllChatHistoryByJoinId(joinUserId string) (*DocUserChatHistory, error) {
	temp := new(DocUserChatHistory)
	err := findAllHistoryMessageById(CollUserChatHistory, joinUserId).Decode(temp)
	return temp, err
}

// Find many chat history between two users by their id and the date range
func FindUserChatHistoryByJoinIdAndDateRange(joinUserId string, startDate, endDate int32) (*DocUserChatHistory, error) {
	cur, err := findManyHistoryMessageByIdAndDateRange(CollUserChatHistory, joinUserId, startDate, endDate)
	if nil != err {
		return nil, err
	}

	// decode data
	temp := new(DocUserChatHistory)
	err = cur.Decode(temp)
	if nil != err {
		return nil, err
	}
	if len(temp.History) == 0 {
		return nil, ErrMessageHistoryNotFound
	}

	return temp, nil
}

// Find one chat history between two users by their id and a special date
func FindUserChatHistoryByJoinIdAndDate(joinUserId string, date int32) (*DocUserChatHistory, error) {
	cur, err := findManyHistoryMessageByIdAndDate(CollUserChatHistory, joinUserId, date)
	if nil != err {
		return nil, err
	}

	// decode data
	temp := new(DocUserChatHistory)
	err = cur.Decode(temp)
	if nil != err {
		return nil, err
	}
	if len(temp.History) == 0 {
		return nil, ErrMessageHistoryNotFound
	}

	return temp, nil

}

// ------------------------------------------------------------------------------------

/* group_chat_history document eg.:
{
	"_id": <the group chat id>,
	"history": [
		{
			"date": <date number: eg. 20190303>,
			"messages": [
				<the message bytes data eg.: []byte("test message string")>,
				....
			]
		},
	...
	]
}
*/
// information for group chat history.
type DocGroupChatHistory struct {
	GroupId int64            `bson:"_id"`
	History []historyMessage `bson:"history"`
}

// Update the message history record for the group chat.
func UpdateGroupChatHistoryById(groupId int64, message []byte) error {
	return updateHistoryMessageForToday(CollGroupChatHistory, groupId, message)
}

// Find the all chat history of the group chat by group id.
func FindAllGroupChatHistoryById(groupId int64) (*DocGroupChatHistory, error) {
	temp := new(DocGroupChatHistory)
	err := findAllHistoryMessageById(CollGroupChatHistory, groupId).Decode(temp)
	return temp, err
}

// Find many group chat history of the group chat by group id and date range.
func FindGroupChatHistoryByIdAndDateRange(groupId int64, startDate, endDate int32) (*DocGroupChatHistory, error) {
	cur, err := findManyHistoryMessageByIdAndDateRange(CollGroupChatHistory, groupId, startDate, endDate)
	if nil != err {
		return nil, err
	}

	// decode data
	temp := new(DocGroupChatHistory)
	err = cur.Decode(temp)
	if nil != err {
		return nil, err
	}
	if len(temp.History) == 0 {
		return nil, ErrMessageHistoryNotFound
	}
	return temp, nil
}

// Find many group chat history of the group chat by group id and a special date.
func FindGroupChatHistoryByIdAndDate(groupId int64, date int32) (*DocGroupChatHistory, error) {
	cur, err := findManyHistoryMessageByIdAndDate(CollGroupChatHistory, groupId, date)
	if nil != err {
		return nil, err
	}

	// decode data
	temp := new(DocGroupChatHistory)
	err = cur.Decode(temp)
	if nil != err {
		return nil, err
	}
	if len(temp.History) == 0 {
		return nil, ErrMessageHistoryNotFound
	}

	return temp, nil

}

// ------------------------------------------------------------------------------------

/* subscription_msg_history document eg.:
{
	"_id": <the subscription id>,
	"history": [
		{
			"date": <date number: eg. 20190303>,
			"messages": [
				<the message bytes data eg.: []byte("test message string")>,
				....
			]
		},
	...
	]
}
*/
// information for subscription message history.
type DocSubscriptionHistory struct {
	SubsId  int64            `bson:"_id"`
	History []historyMessage `bson:"history"`
}

// Update the message history record for the subscription
func UpdateSubscriptionHistoryById(subsId int64, message []byte) error {
	return updateHistoryMessageForToday(CollSubscriptionMsgHistory, subsId, message)
}

// Find the all subscription message history of the subscription by the id.
func FindAllSubscriptionHistoryById(subsId int64) (*DocSubscriptionHistory, error) {
	temp := new(DocSubscriptionHistory)
	err := findAllHistoryMessageById(CollSubscriptionMsgHistory, subsId).Decode(temp)
	return temp, err
}

// Find many subscription history of the subscription by id and date range.
func FindSubscriptionHistoryByIdAndDateRange(subsId int64, startDate, endDate int32) (*DocSubscriptionHistory, error) {
	cur, err := findManyHistoryMessageByIdAndDateRange(CollSubscriptionMsgHistory, subsId, startDate, endDate)
	if nil != err {
		return nil, err
	}

	// decode data
	temp := new(DocSubscriptionHistory)
	err = cur.Decode(temp)
	if nil != err {
		return nil, err
	}
	if len(temp.History) == 0 {
		return nil, ErrMessageHistoryNotFound
	}
	return temp, nil
}

// Find many subscription history of the subscription by id and a special date.
func FindSubscriptionHistoryByIdAndDate(subsId int64, date int32) (*DocSubscriptionHistory, error) {
	cur, err := findManyHistoryMessageByIdAndDate(CollSubscriptionMsgHistory, subsId, date)
	if nil != err {
		return nil, err
	}

	// decode data
	temp := new(DocSubscriptionHistory)
	err = cur.Decode(temp)
	if nil != err {
		return nil, err
	}
	if len(temp.History) == 0 {
		return nil, ErrMessageHistoryNotFound
	}

	return temp, nil
}

// ------------------------------------------------------------------------------------

// Update an id array in a document which found by id in target collection.
// The document would be found by 'queryId' in 'coll'.
func updateIdArrayOfOneDocument(coll *mongo.Collection, queryId, updateId int64, arrayName, operate string, ) error {
	_, err := coll.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": queryId},
		bson.M{operate: bson.M{arrayName: updateId}},
		upsertTrueOption)
	return err
}

/* user_friends document eg.:
{
	"_id": <the user id>,
    "friends": [
        <another user id>,
        <another user id>,
		...
    ],
    "blacklist": [
    	<another user id>,
        <another user id>,
		...
    ]
}
*/
// information for user's friends.
type DocUserFriends struct {
	UserId    int64   `bson:"_id"`
	Friends   []int64 `bson:"friends"`
	Blacklist []int64 `bson:"blacklist"`
}

// Update the 'friends' array to add one new friend's id.
func UpdateUserFriendsToAddFriend(userId, friendId int64) error {
	return updateIdArrayOfOneDocument(CollUserFriends, userId, friendId, "friends", "$addToSet")
}

// Update the 'friends' array to delete one friend's id.
func UpdateUserFriendsToDelFriend(userId, friendId int64) error {
	return updateIdArrayOfOneDocument(CollUserFriends, userId, friendId, "friends", "$pull")
}

// Update the 'blacklist' array to add one user id.
func UpdateUserBlacklistToAddUser(userId, anotherId int64) error {
	return updateIdArrayOfOneDocument(CollUserFriends, userId, anotherId, "blacklist", "$addToSet")
}

// Update the 'blacklist' array to delete one user id.
func UpdateUserBlacklistToDelUser(userId, anotherId int64) error {
	return updateIdArrayOfOneDocument(CollUserFriends, userId, anotherId, "blacklist", "$pull")
}

func FindUserFriendsAndBlacklistById(userId int64) (*DocUserFriends, error) {
	temp := new(DocUserFriends)
	err := CollUserFriends.FindOne(getTimeOutCtx(3), bson.M{"_id": userId}).Decode(temp)
	return temp, err
}

// ------------------------------------------------------------------------------------

/* user_group_chats document eg.:
{
	"_id": <the user id>,
	"groups": [
		<the group chat id>,
		<the group chat id>,
		...
	]
}
*/
// information for the group chat which the user had joined.
type DocUserGroupChats struct {
	UserId int64   `bson:"_id"`
	Groups []int64 `bson:"groups"`
}

// Update the 'groups' array to add one group chat id.
func UpdateUserGroupChatsToAddOne(userId, groupChatId int64) error {
	return updateIdArrayOfOneDocument(CollUserGroupChats, userId, groupChatId, "groups", "$addToSet")
}

// Update the 'groups' array to delete one group chat id.
func UpdateUserGroupChatsToDelOne(userId, groupChatId int64) error {
	return updateIdArrayOfOneDocument(CollUserGroupChats, userId, groupChatId, "groups", "$pull")
}

func FindUserGroupChatsById(userId int64) (*DocUserGroupChats, error) {
	temp := new(DocUserGroupChats)
	err := CollUserGroupChats.FindOne(getTimeOutCtx(3), bson.M{"_id": userId}).Decode(temp)
	return temp, err
}

// ------------------------------------------------------------------------------------

/* user_subscriptions document eg.:
{
	"_id": <the user id>,
	"subscriptions": [
		<the subscription id>,
		<the subscription id>,
		...
	]
}
*/
// information for the subscription which the user had followed.
type DocUserSubscriptions struct {
	UserId        int64   `bson:"_id"`
	Subscriptions []int64 `bson:"subscriptions"`
}

// Update the 'subscription' array to add one subscription id.
func UpdateUserSubscriptionsToAddOne(userId, subsId int64) error {
	return updateIdArrayOfOneDocument(CollUserSubscriptions, userId, subsId, "subscriptions", "$addToSet")
}

// Update the 'subscription' array to delete one subscription id.
func UpdateUserSubscriptionsToDelOne(userId, subsId int64) error {
	return updateIdArrayOfOneDocument(CollUserSubscriptions, userId, subsId, "subscriptions", "$pull")
}

func FindUserSubscriptionsById(userId int64) (*DocUserSubscriptions, error) {
	temp := new(DocUserSubscriptions)
	err := CollUserSubscriptions.FindOne(getTimeOutCtx(3), bson.M{"_id": userId}).Decode(temp)
	return temp, err
}

// ------------------------------------------------------------------------------------
