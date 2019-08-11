package controllers

import (
	"../models"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"mime/multipart"
)

const (
	DefaultAvatarUrl        = "this is the default avatar url"
	StaticResourceUrlPrefix = "this is the static resource url prefix"
	MaxAvatarUploadSize     = 100 * 2 << 10
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
	c.JSON(200, gin.H{"avatar": StaticResourceUrlPrefix + *avatar})
}

// PutAvatar HTTP API function
func PutAvatar(c *gin.Context) {
	//userId := c.MustGet("user_id")
	file, err := c.FormFile("new_avatar")
	if nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if file.Size > MaxAvatarUploadSize || file.Size == 0 {
		c.JSON(400, gin.H{"error": "the upload image size need gt=0kb and lte=100kb"})
		return
	}
	hashName, err := saveAvatarFile(file)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userId := c.MustGet("user_id")
	err = models.MySQLPutUserAvatar(userId.(int64), hashName)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"avatar": StaticResourceUrlPrefix + hashName})
}

// try to save avatar file data with a hashName
func saveAvatarFile(file *multipart.FileHeader) (hashName string, err error) {
	// read data from file
	buff := new(bytes.Buffer)
	fp, err := file.Open()
	defer fp.Close()
	if nil != err {
		return
	}
	n, err := io.Copy(buff, fp)
	if nil != err {
		return
	}
	if n == 0 {
		err = errors.New("read file data error")
		return
	}
	// get hashName of data
	h := md5.New()
	h.Write(buff.Bytes())
	hashName = hex.EncodeToString(h.Sum(nil))
	err = UploadDataToStaticServer(buff.Bytes(), hashName)
	return
}

// upload the data to static file server with a hashName
func UploadDataToStaticServer(data []byte, hashName string) error {
	// todo
	return nil
}
