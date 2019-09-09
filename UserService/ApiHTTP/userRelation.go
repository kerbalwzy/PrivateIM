package ApiHTTP

import (
	"errors"
	"github.com/gin-gonic/gin"

	"../ApiRPC"

	pb "../Protos"
)

type GetFriendParams struct {
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
	Avatar string `json:"avatar"`
}

// Get friend HTTP API function.
// Search the user by email or name. Only when the friendship is effect, the value of Mobile,
// Gender and Note can be show for searcher.
func GetFriend(c *gin.Context) {
	params := new(GetFriendParams)
	if err := c.ShouldBindJSON(params); nil != err {
		c.JSON(400, gin.H{"error": "invalid params " + err.Error()})
		return
	}
	if params.Email == "" && params.Name == "" {
		c.JSON(400, gin.H{"error": "invalid params"})
		return
	}
	// search users by param, if the length of result is zero, return now.
	users, err := GetUsersByEmailOrName(params)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(users) == 0 {
		c.JSON(404, gin.H{"error": "not found by params"})
		return
	}

	// search the friendship record data of the user
	selfId := c.MustGet(JWTGetUserId).(int64)
	friendsIdNote, err := GetFriendsIdAndNoteOfUser(selfId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// organizational results
	result := make([]*GetFriendResult, 1)
	for _, user := range users {
		ret := &GetFriendResult{Id: user.Id, Name: user.Name, Email: user.Email}
		if note, ok := friendsIdNote[user.Id]; ok {
			ret.Mobile = user.Mobile
			ret.Gender = int(user.Gender)
			ret.Note = note
		}
		data, _ := ApiRPC.GetUserAvatarById(user.Id)
		if data != nil {
			ret.Avatar = data.Avatar
		}

		result = append(result, ret)
	}

	c.JSON(200, gin.H{"result": result})
}

// Search users by email or name, prefer to use email
func GetUsersByEmailOrName(params *GetFriendParams) ([]*pb.UserBasicInfo, error) {
	users := make([]*pb.UserBasicInfo, 1)
	if params.Email != "" {
		ret, err := ApiRPC.GetUserByEmail(params.Email)
		if nil != err {
			return nil, err
		}
		users = append(users, ret)
		return users, nil
	} else {
		ret, err := ApiRPC.GetUsersByName(params.Name)
		if nil != err {
			return nil, err
		}
		users = append(users, ret.Data...)
		return users, nil
	}
}

// Get the id and note of the friends, and return a map, key is id, value is note
func GetFriendsIdAndNoteOfUser(userId int64) (map[int64]string, error) {
	tempMap := make(map[int64]string)
	ret, err := ApiRPC.GetFriendshipInfo(userId)
	if nil != err {
		return tempMap, err
	}
	for _, relate := range ret.Data {
		// Judging the validity of relationship, only the effect friendship can be save
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

	// todo: Let the MessageService know the change
	//NotifyTargetUser(params.FriendId)

}

// Check if duplicate add the user as friend or self is in black list of target user
func CheckAndAddFriend(selfId, friendId int64, note string) (int, error) {
	if selfId == friendId {
		return 400, ErrAddSelfAsFriend
	}
	_, err := ApiRPC.AddOneNewFriend(selfId, friendId, note)
	if nil != err {
		return 500, err
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

	// do different thing by action
	switch params.Action {
	case 1:
		statusCode, message = ModifyFriendNote(selfId, params)
	case 2:
		statusCode, message = CheckAndAcceptFriend(selfId, params)
	case 3:
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
	if _, err := ApiRPC.PutOneFriendNote(selfId, params.FriendId, params.Note); nil != err {
		return 500, err.Error()
	}

	return 200, "modify note for friend successfully"
}

// Handle a friend relationship request
func CheckAndAcceptFriend(selfId int64, params *PutFriendParams) (int, string) {
	// check if the friend request existed
	_, err := ApiRPC.AcceptOneNewFriend(selfId, params.FriendId, params.Note, params.IsAccept)

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
	_, err := ApiRPC.PutFriendBlacklist(selfId, friendId, isBlack)

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
	_, err := ApiRPC.DeleteOneFriend(selfId, params.FriendId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "the record has been deleted and will not be notified to your friend."})
}

// Get All Friend HTTP API function
func AllFriends(c *gin.Context) {
	selfId := c.MustGet(JWTGetUserId).(int64)
	ret, err := ApiRPC.GetFriendsBasicInfo(selfId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"friends": ret.Data})
}
