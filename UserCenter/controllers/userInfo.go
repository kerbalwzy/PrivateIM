package controllers

import (
	"../models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// GetProfile HTTP API function
func GetProfile(c *gin.Context) {
	user := models.UserBasic{Id: c.MustGet("user_id").(int64)}
	err := user.MySQLGetByField("Id")
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
	userP := &models.UserBasic{Id: userId.(int64)}
	err = userP.MySQLUpdateProfile(tempProfileP.Name, tempProfileP.Mobile, tempProfileP.Gender)
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
