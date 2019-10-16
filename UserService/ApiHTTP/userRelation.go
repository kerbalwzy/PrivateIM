package ApiHTTP

import (
	"errors"
	"github.com/gin-gonic/gin"

	"../ApiRPC"
	"../ElasticClient"
)

type SearchUsersParam struct {
	KW   string `json:"kw" binding:"required"`
	Page int    `json:"page"`
	Size int    `json:"size"`
}

// Search user HTTP API function.
// search data from ElasticSearch
func SearchUsers(c *gin.Context) {
	param := new(SearchUsersParam)
	if err := c.ShouldBindJSON(param); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if param.Page == 0 {
		param.Page = 1
	}
	if param.Size == 0 {
		param.Size = 10
	}
	data, err := elasticClient.UserIndexDocSearch(param.KW, param.Page, param.Size)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
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

}

// Check if duplicate add the user as friend or self is in black list of target user
func CheckAndAddFriend(selfId, friendId int64, note string) (int, error) {
	if selfId == friendId {
		return 400, ErrAddSelfAsFriend
	}
	err := ApiRPC.AddOneNewFriend(selfId, friendId, note)
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

//// Modify note on my friends
func ModifyFriendNote(selfId int64, params *PutFriendParams) (int, string) {
	if params.Note == "" {
		return 400, "note for friend not allow be an empty string"
	}
	if err := ApiRPC.PutOneFriendNote(selfId, params.FriendId, params.Note); nil != err {
		return 500, err.Error()
	}

	return 200, "modify note for friend successfully"
}

// Handle a friend relationship request
func CheckAndAcceptFriend(selfId int64, params *PutFriendParams) (int, string) {
	// check if the friend request existed
	err := ApiRPC.AcceptOneNewFriend(selfId, params.FriendId, params.Note, params.IsAccept)
	if nil != err {
		return 500, err.Error()
	}
	if params.IsAccept {
		ApiRPC.MSGUserNodeAddFriend(selfId, params.FriendId)
		ApiRPC.MSGUserNodeAddFriend(params.FriendId, selfId)

		return 200, "you are friend now, chat happy"
	} else {
		ApiRPC.MSGUserNodeAddBlacklist(selfId, params.FriendId)
		return 200, "you refused and move the user to you blacklist"
	}

}

// Manage the friend blacklist
func ManageFriendShipBlacklist(selfId, friendId int64, isBlack bool) (int, string) {
	err := ApiRPC.PutOneFriendIsBlack(selfId, friendId, isBlack)

	if nil != err {
		return 500, err.Error()
	}
	if isBlack {
		ApiRPC.MSGUserNodeMoveFriendIntoBlacklist(selfId, friendId)
		return 200, "move friend into blacklist successfully"
	} else {
		ApiRPC.MSGUserNodeMoveFriendOutFromBlacklist(selfId, friendId)
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
	err := ApiRPC.DeleteOneFriend(selfId, params.FriendId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ApiRPC.MSGUserNodeDelFriend(selfId, params.FriendId)
	c.JSON(200, gin.H{"message": "the record has been deleted and will not be notified to your friend."})
}

// Get All the user's friends information HTTP API function
func GetFriendsInfo(c *gin.Context) {
	selfId := c.MustGet(JWTGetUserId).(int64)
	ret, err := ApiRPC.GetUserFriendsInfo(selfId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"friends": ret.Data})
}
