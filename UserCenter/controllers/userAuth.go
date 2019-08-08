package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"../models"
)

const (
	AuthTokenSalt = "uhfuwhfhw!23rp93242ihashf;3rbi;u2137974789y3kjnf&#^lknfa"
)

// Sign Up data struct , all field required.
// Verify data by the validators of gin binding.
type UserSignUp struct {
	Name            string `json:"name" binding:"required,max=10"`
	Email           string `json:"email"  binding:"required,email"`
	Password        string `json:"password"  binding:"required,min=8,max=12"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// SignUp(Register) HTTP API function
func SignUp(c *gin.Context) {
	var err error
	var item = UserSignUp{}

	if err = c.ShouldBindJSON(&item); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// check the email if registered
	user := models.UserBasic{Name: item.Name, Email: item.Email}
	err = user.MySQLGetByField("Email")
	if nil == err {
		c.JSON(400, gin.H{"error": "email is already sign up, please sign in"})
		return
	}

	// save user information to database
	user.SetPassword(item.Password)
	err = user.MySQLSignUp()
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	//ok, return user detail and AuthToken
	authToken, err := MakeAuthToken(user, AuthTokenSalt, int64(time.Hour)*24*30)
	data := gin.H{"user": user, "AuthToken": authToken}
	c.JSON(http.StatusOK, data)
}

type UserSignIn struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8,max=12"`
}

// SignIn(Login) HTTP API function
func SignIn(c *gin.Context) {
	var err error
	var item = UserSignIn{}

	if err = c.ShouldBindJSON(&item); nil != err {
		c.JSON(400, gin.H{"error": "verify fail"})
		return
	}

	// check email and password for user
	user := models.UserBasic{Email: item.Email}
	err = user.MySQLGetByField("Email")
	if nil != err {
		c.JSON(400, gin.H{"error": "verify fail"})
		return
	}
	if !user.CheckPassword(item.Password) {
		c.JSON(400, gin.H{"error": "verify fail"})
		return
	}

	//verify ok, return user detail and AuthToken
	authToken, err := MakeAuthToken(user, AuthTokenSalt, int64(time.Hour)*24*30)
	data := gin.H{"user": user, "AuthToken": authToken}
	c.JSON(http.StatusOK, data)
}
