package controllers

import (
	"database/sql"
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
			result.Mobile = userP.Mobile
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
	Note string `json:"note" binding:"nameValidator"`
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
		return
	}

	// try save relation data into database
	relateP.FriendNote = params.Note
	err = models.MySQLAddOneFriend(relateP)
	if nil != err {
		c.JSON(500, gin.H{"error": "save relation error: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"email": userP.Email, "name": userP.Name,
		"note": relateP.FriendNote, "is_accept": relateP.IsAccept,
		"is_refuse": relateP.IsRefuse, "is_delete": relateP.IsDelete})

	// todo: Let the communication center notify the target user
	NotifyTargetUser(relateP)

}

// check if duplicate add the user as friend or the user is in black list
func CheckDuplicateAddAndBlackList(relateP *models.UserRelate) (int, error) {
	if err := models.MySQLGetUserOneFriend(relateP); nil != err && sql.ErrNoRows != err {
		return 500, errors.New("database error:" + err.Error())
	}
	if relateP.Id != 0 && relateP.IsAccept && !relateP.IsRefuse && !relateP.IsDelete {
		return 400, errors.New("you are already friends")
	}
	if relateP.Id != 0 && relateP.IsRefuse {
		return 403, errors.New("there's a blacklist relationship between you")
	}
	return 200, nil
}

type PutFriendParams struct {
	Action   int    `json:"action binding:relateActionValidator"`
	SelfId   int64  `json:"src_id" binding:"required"`
	FriendId int64  `json:"dst_id" binding:"required"`
	Note     string `json:"note binding:max=10"`
	IsAccept bool   `json:"is_accept"`
	IsRefuse bool   `json:"is_refuse"`
}

// Put friend HTTP API function
func PutFriend(c *gin.Context) {
	params := new(PutFriendParams)
	if err := c.ShouldBindJSON(params); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// modify friend note
	if params.Action == 1 {

	}

	// handle friend request
	if params.Action == 2 {
		if status, message := HandleFriendRequest(params); 200 != status {
			c.JSON(status, gin.H{"error": message})
			return
		} else {
			c.JSON(status, gin.H{"message": message})
			return
		}

	}

}

func ModifyFriendNote(params *PutFriendParams) (int, string) {
	return 200, ""
}

func HandleFriendRequest(params *PutFriendParams) (int, string) {
	return 200, "you are friend now, chat happy"
}
