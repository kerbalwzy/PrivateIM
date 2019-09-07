package ApiHTTP

import (
	"errors"
	"github.com/gin-gonic/gin"

	"../DataLayer"
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
	selfId := c.MustGet(JWTGetUserId).(int64)
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
func SearchOtherUsers(selfId int64, params *GetFriendParams) ([]*DataLayer.UserBasic, error) {
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
func SearchUsers(params *GetFriendParams) ([]*DataLayer.UserBasic, error) {
	users := make([]*DataLayer.UserBasic, 0)
	// if the Id not zero, use Id to query user first
	userP := new(DataLayer.UserBasic)
	if params.Id != 0 {
		userP.Id = params.Id
		_ = DataLayer.MySQLGetUserById(userP)
	}
	// if the email still empty string, mean not found by id
	if userP.Email == "" {
		userP.Email = params.Email
		_ = DataLayer.MySQLGetUserByEmail(userP)
	}
	// if the name still empty string, mean not found by id and email
	if userP.Name == "" {
		users, _ = DataLayer.MySQLGetUserByName(params.Name)
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
	friendsRelates, err := DataLayer.MySQLGetUserFriendsRelates(userId)
	if nil != err {
		return tempMap, err
	}
	for _, relate := range friendsRelates {
		// Judging the validity of relationship
		if relate.IsAccept && !relate.IsBlack && !relate.IsDelete {
			tempMap[relate.FriendId] = relate.FriendNote
		}
	}
	return tempMap, nil
}

type AddFriendParams struct {
	FriendId int64  `json:"dst_id" binding:"required"`
	Note     string `json:"note" binding:"nameValidator"`
}

var (
	ErrAddSelfAsFriend = errors.New("can't build a friendship with yourself")
)

// Add friend HTTP API function
func AddFriend(c *gin.Context) {
	params := &AddFriendParams{}
	if err := c.ShouldBindJSON(params); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// check the params
	selfId := c.MustGet(JWTGetUserId).(int64)
	statusCode, err := CheckAndAddFriend(selfId, params.FriendId, params.Note)
	if nil != err {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.JSON(statusCode, gin.H{"message": "initiate and add friends successfully, wait for the target user to agree"})

	// todo: Let the communication center notify the target user
	//NotifyTargetUser(params.FriendId)

}

// Check if duplicate add the user as friend or self is in black list of target user
func CheckAndAddFriend(selfId, friendId int64, note string) (int, error) {
	if selfId == friendId {
		return 400, ErrAddSelfAsFriend
	}
	err := DataLayer.MySQLAddOneFriend(selfId, friendId, note)
	if err == DataLayer.ErrTargetUserNotExisted || err == DataLayer.ErrFriendshipAlreadyInEffect {
		return 400, err
	}

	if err == DataLayer.ErrInBlackList {
		return 403, err
	}

	return 200, nil
}

type PutFriendParams struct {
	Action   int    `json:"action" binding:"relateActionValidator"`
	FriendId int64  `json:"dst_id" binding:"required"`
	Note     string `json:"note" binding:"max=10"`
	IsAccept bool   `json:"is_accept"`
	IsBlack  bool   `json:"is_black"`
}

// Put friend HTTP API function
func PutFriend(c *gin.Context) {
	params := new(PutFriendParams)
	if err := c.ShouldBindJSON(params); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	selfId := c.MustGet(JWTGetUserId).(int64)
	if params.FriendId == selfId {
		c.JSON(400, gin.H{"action": params.Action,
			"error": "can't modify a friendship with yourself"})
		return
	}

	statusCode := 200
	message := ""
	// modify friend note
	if params.Action == 1 {
		statusCode, message = ModifyFriendNote(selfId, params)
	}

	// handle friend request
	if params.Action == 2 {
		statusCode, message = CheckAndAcceptFriend(selfId, params)
	}

	// move friend to blacklist in or out
	if params.Action == 3 {
		statusCode, message = ManageFriendShipBlacklist(selfId, params.FriendId, params.IsBlack)
	}

	if 200 != statusCode {
		c.JSON(statusCode, gin.H{"action": params.Action, "error": message})
	} else {
		c.JSON(statusCode, gin.H{"action": params.Action, "message": message})
	}

}

// Modify note on my friends
func ModifyFriendNote(selfId int64, params *PutFriendParams) (int, string) {
	if params.Note == "" {
		return 400, "note for friend not allow be an empty string"
	}
	if err := DataLayer.MySQLModifyNoteOfFriend(selfId, params.FriendId, params.Note); nil != err {
		if err == DataLayer.ErrNoFriendship {
			return 400, err.Error()
		} else {
			return 500, err.Error()
		}
	}

	return 200, "modify note for friend successfully"
}

// Handle a friend relationship request
func CheckAndAcceptFriend(selfId int64, params *PutFriendParams) (int, string) {
	// check if the friend request existed
	err := DataLayer.MySQLAcceptOneFriend(selfId, params.FriendId, params.Note, params.IsAccept)
	if err == DataLayer.ErrFriendRequestNotExisted || err == DataLayer.ErrFriendshipAlreadyInEffect {
		return 400, err.Error()
	}

	if nil != err {
		return 500, err.Error()
	}

	if params.IsAccept {
		return 200, "you are friend now, chat happy"
	} else {
		return 200, "you refused and move the user to you blacklist"
	}

}

// Manage the friend blacklist
func ManageFriendShipBlacklist(selfId, friendId int64, isBlack bool) (int, string) {
	err := DataLayer.MySQLManageFriendBlacklist(selfId, friendId, isBlack)
	if DataLayer.ErrFriendBlacklistNoChange == err {
		return 400, err.Error()
	}

	if nil != err {
		return 500, err.Error()
	}
	if isBlack {
		return 200, "move friend into blacklist successfully"
	} else {
		return 200, "move friend out from blacklist successfully"
	}
}

type DeleteFriendParams struct {
	FriendId int64 `json:"dst_id" binding:"required"`
}

// Delete friend HTTP API function
func DeleteFriend(c *gin.Context) {
	params := new(DeleteFriendParams)
	if err := c.ShouldBindJSON(params); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	selfId := c.MustGet(JWTGetUserId).(int64)
	err := DataLayer.MySQLDeleteOneFriend(selfId, params.FriendId)
	if DataLayer.ErrNoFriendship == err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "the record has been deleted and will not be notified to your friend."})
}

// Get All Friend HTTP API function
func AllFriends(c *gin.Context) {
	selfId := c.MustGet(JWTGetUserId).(int64)
	data, err := DataLayer.MySQLGetUserFriendsInfo(selfId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"friends": data})
}