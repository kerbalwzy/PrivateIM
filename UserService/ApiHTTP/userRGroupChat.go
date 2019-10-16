package ApiHTTP

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"

	"../ApiRPC"
	"../ElasticClient"
	"../utils"

	conf "../Config"
)

var (
	ErrCreatedSameNameGroupChat = errors.New("you created the same group chat already")
	ErrFindManagerFail          = errors.New("find the manager fail")
)

type NewGroupChatParam struct {
	Name string `json:"name" binding:"nameValidator"`
}

// Create one new group chat with
func NewOneGroupChat(c *gin.Context) {
	param := new(NewGroupChatParam)
	if err := c.ShouldBindJSON(param); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	managerId := c.MustGet(JWTGetUserId).(int64)
	tempGroupChat, _ := ApiRPC.GetOneGroupChatByNameAndManger(param.Name, managerId)
	if tempGroupChat != nil {
		c.JSON(400, gin.H{"error": ErrCreatedSameNameGroupChat.Error()})
		return
	}
	manager, err := ApiRPC.GetUserById(managerId)
	if nil != err {
		c.JSON(400, gin.H{"error": ErrFindManagerFail})
		return
	}

	qrCodePicData, _ := utils.CreatQRCodeBytes(fmt.Sprintf("group_chat_name=%s&manager=%d", param.Name, managerId))
	qrCodePicHashName := utils.BytesDataHash(qrCodePicData)

	tempGroupChat, err = ApiRPC.SaveOneNewGroupChat(
		param.Name,
		conf.DefaultAvatarPicName,
		qrCodePicHashName,
		managerId,
	)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// save the qr code pic
	_ = SaveQRCodeLocal(qrCodePicData, qrCodePicHashName)

	manager.Avatar = conf.PhotosUrlPrefix + manager.Avatar + conf.PhotoSuffix
	tempGroupChat.Avatar = conf.PhotosUrlPrefix + tempGroupChat.Avatar + conf.PhotoSuffix

	// update the index of ElasticSearch
	_ = elasticClient.GroupChatIndexDocSave(
		tempGroupChat.Id,
		tempGroupChat.Name,
		tempGroupChat.Avatar,
		manager.Name,
		manager.Avatar)

	c.JSON(201, gin.H{"group_chat": tempGroupChat})
}

// todo: update the 'is_delete' to false of the group chat
func DismissOneGroupChat(c *gin.Context) {
	c.JSON(200, "waiting implement")
}

// todo:
func PutOneGroupChatName(c *gin.Context) {
	c.JSON(200, "waiting implement")
}

// todo:
func PutOneGroupChatAvatar(c *gin.Context) {
	c.JSON(200, "waiting implement")
}

// todo: change the manager of the group chat to other user.
func PutOneGroupChatManager(c *gin.Context) {
	c.JSON(200, "waiting implement")
}

// todo:
func GetOneGroupChatInfo(c *gin.Context) {
	c.JSON(200, "waiting implement")
}

type SearchParam struct {
	KW   string `json:"kw" binding:"required"`
	Page int    `json:"page"`
	Size int    `json:"size"`
}

// Search group chat information from ElasticSearch
func SearchGroupChats(c *gin.Context) {
	param := new(SearchParam)
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
	data, err := elasticClient.GroupChatIndexDocSearch(param.KW, param.Page, param.Size)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

// todo
func JoinOneGroupChat(c *gin.Context) {
	c.JSON(200, "waiting implement")
}

// todo
func QuitOneGroupChat(c *gin.Context) {
	c.JSON(200, "waiting implement")
}

// todo
func PutSelfNoteInGroupChat(c *gin.Context) {
	c.JSON(200, "waiting implement")
}

// The all the group chat information which the user joined
func GetGroupChatsUserJoined(c *gin.Context) {
	userId := c.MustGet(JWTGetUserId).(int64)
	infoList, err := ApiRPC.GetGroupChatsInfoTheUserJoined(userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"group_chat_list": infoList.Data})
}

// Get the information of users whom joined the group chat
func GetUsersInfoOfGroupChat(c *gin.Context) {
	groupChatId := c.Query("group_id")
	if groupChatId == "" {
		c.JSON(400, gin.H{"error": "query param error"})
		return
	}
	groupId, err := strconv.ParseInt(groupChatId, 10, 64)
	if nil != err {
		c.JSON(400, gin.H{"error": "query param error"})
		return
	}
	userId := c.MustGet(JWTGetUserId).(int64)

	tempUserGroupChat, _ := ApiRPC.GetOneUserGroupChat(userId, groupId)
	if nil == tempUserGroupChat {
		c.JSON(400, gin.H{"error": "you are not the member of the group chat"})
		return
	}

	infoList, err := ApiRPC.GetUsersInfoOfTheGroupChat(groupId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	for _, info := range infoList.Data {
		info.UserAvatar = conf.PhotosUrlPrefix + info.UserAvatar + conf.PhotoSuffix
	}
	c.JSON(200, gin.H{"users": infoList.Data})
}
