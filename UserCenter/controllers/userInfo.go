package controllers

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"../models"
	"../utils"
)

const (
	AuthTokenSalt      = "this is a auth token salt"
	AuthTokenAliveTime = 3600 * 24 //unit:second
	AuthTokenIssuer    = "userCenter"

	PhotoSaveFoldPath   = "/Users/wzy/GitProrgram/PrivateIM/UserCenter/static/photos/"
	PhotoSuffix         = ".png"
	DefaultAvatarUrl    = "/static/photos/defaultAvatar.png"
	PhotosUrlPrefix     = "/static/photos/" // if you use oss , should change this value
	MaxAvatarUploadSize = 100 * 2 << 10
)

// GetProfile HTTP API function
func GetProfile(c *gin.Context) {
	user := models.UserBasic{Id: c.MustGet("user_id").(int64)}
	err := models.MySQLGetUserByField("Id", &user)
	if nil != err {
		c.JSON(404, gin.H{"error": "get user information fail"})
		return
	}
	c.JSON(200, user)
}

type TempProfile struct {
	Name   string `json:"name" binding:"nameValidator"`
	Mobile string `json:"mobile" binding:"mobileValidator"`
	Gender int    `json:"gender" binding:"genderValidator"`
}

// PutProfile HTTP API function
func PutProfile(c *gin.Context) {
	// Validate the params
	tempProfileP, err := parseTempProfile(c)
	if nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	errs := binding.Validator.ValidateStruct(tempProfileP)
	if nil != errs {
		c.JSON(400, gin.H{"errors": errs.Error()})
		return
	}
	// Update user info
	userId := c.MustGet("user_id")
	err = models.MySQLUpdateProfile(tempProfileP.Name, tempProfileP.Mobile, tempProfileP.Gender, userId.(int64))
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tempProfileP)
}

// Parse the JsonBodyParams to TempProfile
func parseTempProfile(c *gin.Context) (*TempProfile, error) {
	// Parse the JsonBodyParams to map
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	if n == 0 {
		return nil, errors.New("not have any JsonBodyParams")
	}

	tempDict := make(map[string]interface{})
	_ = json.Unmarshal(buf[0:n], &tempDict)
	// Check the integrity of parameters
	name, ok := tempDict["name"]
	if !ok {
		return nil, errors.New("`name` not exited in JsonBodyParams")
	}
	mobile, ok := tempDict["mobile"]
	if !ok {
		return nil, errors.New("`mobile` not exited in JsonBodyParams")
	}
	gender, ok := tempDict["gender"]
	if !ok {
		return nil, errors.New("`gender` not exited in JsonBodyParams")
	}

	tempProfileP := &TempProfile{
		Name:   name.(string),
		Mobile: mobile.(string),
		Gender: int(gender.(float64))}
	return tempProfileP, nil
}

// GetAvatar HTTP API function
func GetAvatar(c *gin.Context) {
	userId := c.MustGet("user_id")
	avatar := new(string)
	err := models.MySQLGetUserAvatar(userId.(int64), avatar)
	if nil != err {
		c.JSON(500, gin.H{"error": "query avatar fail"})
		return
	}
	if *avatar == "" {
		c.JSON(200, gin.H{"avatar": DefaultAvatarUrl})
		return
	}
	c.JSON(200, gin.H{"avatar": PhotosUrlPrefix + *avatar + PhotoSuffix})
}

// PutAvatar HTTP API function
func PutAvatar(c *gin.Context) {
	// get file data and hash value as name
	file, err := c.FormFile("new_avatar")
	if nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if file.Size > MaxAvatarUploadSize || file.Size == 0 {
		c.JSON(400, gin.H{"error": "the upload image size need gt=0kb and lte=100kb"})
		return
	}
	hashName, data, err := utils.GinFormFileHash(file)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// check if the hash name is existed more then one in the table
	// if true, it meanings that has the same file already upload.
	// not need to save the file again.
	count := models.MySQLAvatarHashNameCount(hashName)
	if count == 0 {
		// save the file data to local or static server
		if err := UploadAvatarLocal(data, hashName); nil != err {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	// save the information into database
	userId := c.MustGet("user_id")
	err = models.MySQLPutUserAvatar(userId.(int64), hashName)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"avatar": PhotosUrlPrefix + hashName + PhotoSuffix})
}

// save avatar file to local
func UploadAvatarLocal(data []byte, hashName string) error {
	prefix := PhotoSaveFoldPath
	suffix := PhotoSuffix
	path := prefix + hashName + suffix
	if err := utils.UploadFileToLocal(data, path); nil != err {
		return err
	}
	return nil
}

// todo upload the data to cloud with a hashName
func UploadDataToCloud(data []byte, hashName string) error {
	return nil
}
