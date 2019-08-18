package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"

	"../models"
)

type GetFriendParams struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
type GetFriendResult struct {
	Id     int64  `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Gender int    `json:"gender"`
	Note   string `json:"note"`
}

// Get friend HTTP API function
func GetFriend(c *gin.Context) {
	params := new(GetFriendParams)
	if err := c.ShouldBindJSON(params); nil != err {
		c.JSON(400, gin.H{"error": "invalid params " + err.Error()})
		return
	}
	if params.Id == 0 && params.Email == "" && params.Name == "" {
		c.JSON(400, gin.H{"error": "invalid params"})
		return
	}
	// search the other users by params
	selfId := c.MustGet(JWTDataKey).(int64)
	users, err := SearchOtherUsers(selfId, params)
	if nil != err {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	// search the current friends of user
	friendsIdAndNote, _ := GetFriendsIdAndNoteOfUser(selfId)
	// hidden the private information if the friend relationship not existed
	// create the result slice
	resultSlice := make([]*GetFriendResult, 0)
	for _, userP := range users {
		result := &GetFriendResult{Id: userP.Id, Email: userP.Email, Name: userP.Name}
		if note, ok := friendsIdAndNote[userP.Id]; ok {
			result.Gender = userP.Gender
			result.Note = note
		}
		resultSlice = append(resultSlice, result)
	}
	c.JSON(200, gin.H{"result": resultSlice})
}

// Search other users by params (GetFriendParams)
func SearchOtherUsers(selfId int64, params *GetFriendParams) ([]*models.UserBasic, error) {
	users, err := SearchUsers(params)
	if nil != err {
		return nil, err
	}

	// delete self from users, if the slice is empty after that, return error
	delIndex := -1
	for index, user := range users {
		if user.Id == selfId {
			delIndex = index
			break
		}
	}
	if delIndex != -1 {
		users = append(users[:delIndex], users[delIndex+1:]...)
	}

	if len(users) == 0 {
		return nil, errors.New("not found by params")
	}
	return users, nil
}

// Search users by params (GetFriendParams)
func SearchUsers(params *GetFriendParams) ([]*models.UserBasic, error) {
	users := make([]*models.UserBasic, 0)
	// if the Id not zero, use Id to query user first
	userP := new(models.UserBasic)
	if params.Id != 0 {
		userP.Id = params.Id
		_ = models.MySQLGetUserById(userP)
	}
	// if the email still empty string, mean not found by id
	if userP.Email == "" {
		userP.Email = params.Email
		_ = models.MySQLGetUserByEmail(userP)
	}
	// if the name still empty string, mean not found by id and email
	if userP.Name == "" {
		users, _ = models.MySQLGetUserByName(params.Name)
	}
	// if the id is not zero, mean found by id or email
	if userP.Id != 0 {
		users = append(users, userP)
	}
	if len(users) == 0 {
		return nil, errors.New("not found by params")
	}
	return users, nil
}

// Get the id and note of the friends, and return a map, key is id, value is note
func GetFriendsIdAndNoteOfUser(userId int64) (map[int64]string, error) {
	tempMap := make(map[int64]string)
	friends, err := models.MySQLGetUserAllFriends(userId)
	if nil != err {
		return tempMap, err
	}
	for _, friend := range friends {
		// Judging the validity of friendship
		if friend.IsAccept && !friend.IsRefuse && !friend.IsDelete {
			tempMap[friend.FriendId] = friend.FriendNote
		}
	}
	return tempMap, nil
}

type AddFriendParams struct {
	Id   int64  `json:"id" binding:"required"`
	Note string `json:"note" binding:"required,max=10"`
}

// Add friend HTTP API function
func AddFriend(c *gin.Context) {
	params := &AddFriendParams{}
	if err := c.ShouldBindJSON(params); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// check user if existed
	userP := &models.UserBasic{Id: params.Id}
	err := models.MySQLGetUserById(userP)
	if nil != err || userP.Email == "" {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	// if duplicate add or `is_refuse` is true, don't continue
	// `is_refuse` is true meaning the user don't accept the friend request
	// or some one between them move the other into black list
	selfId := c.MustGet(JWTDataKey).(int64)
	relateP := &models.UserRelate{SelfId: selfId, FriendId: userP.Id}
	if statusCode, err := CheckDuplicateAddAndBlackList(relateP); nil != err {
		c.JSON(statusCode, gin.H{"error": err.Error()})
	}

	// try save relation data into database
	relateP.FriendNote = params.Note


}

// check if duplicate add the user as friend or the user is in black list
func CheckDuplicateAddAndBlackList(relateP *models.UserRelate) (int, error) {
	if err := models.MySQLGetUserOneFriend(relateP); nil != err {
		return 500, errors.New("database error:" + err.Error())
	}
	if relateP.Id != 0 && relateP.IsAccept && !relateP.IsRefuse {
		return 400, errors.New("you are already friend, don't duplicate add")
	}
	if relateP.Id != 0 && relateP.IsRefuse {
		return 403, errors.New("please remove user from black list first")
	}
	return 200, nil
}

//Check the user relationship.
//If they are friends, return true, relation data id and friend note
//otherwise return false, 0 and empty string
func CheckFriendRelation(selfId, friendId int64) (bool, *models.UserRelate) {
	relateP := &models.UserRelate{SelfId: selfId, FriendId: friendId}
	_ = models.MySQLGetUserOneFriend(relateP)
	if relateP.Id == 0 || !relateP.IsAccept || relateP.IsDelete {
		return false, relateP
	}
	return true, relateP

}
