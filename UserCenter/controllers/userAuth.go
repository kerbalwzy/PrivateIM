package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"../models"
	"../utils"
)

// Sign Up data struct , all field required.
// Verify data by the validators of gin binding.
type UserSignUp struct {
	Name            string `json:"name" binding:"nameValidator"`
	Email           string `json:"email"  binding:"emailValidator"`
	Password        string `json:"password"  binding:"passwordValidator"`
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
	userP := &models.UserBasic{Name: item.Name, Email: item.Email}
	err = models.MySQLGetUserByEmail(userP)
	if nil == err || userP.Id != 0 {
		c.JSON(400, gin.H{"error": "email is already sign up, please sign in"})
		return
	}

	// save user information to database
	userP.SetPassword(item.Password)
	err = models.MySQLUserSignUp(userP)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	//ok, return user detail and auth token
	c.JSON(201, detailAndToken(userP))
}

type UserSignIn struct {
	Email    string `json:"email"  binding:"emailValidator"`
	Password string `json:"password"  binding:"passwordValidator"`
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
	userP := &models.UserBasic{Email: item.Email}
	err = models.MySQLGetUserByEmail(userP)
	if nil != err {
		c.JSON(400, gin.H{"error": "verify fail"})
		return
	}
	if !userP.CheckPassword(item.Password) {
		c.JSON(400, gin.H{"error": "verify fail"})
		return
	}

	//verify ok, return user detail and AuthToken
	c.JSON(http.StatusOK, detailAndToken(userP))
}

// Create auth token by user, and return data
func detailAndToken(user *models.UserBasic) gin.H {
	claims := utils.CustomJWTClaims{
		Id: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + AuthTokenAliveTime), // expire time
			Issuer:    AuthTokenIssuer,                               //signal issuer
		},
	}
	authToken, _ := utils.CreateJWTToken(claims, []byte(AuthTokenSalt))
	return gin.H{"user": user, "auth_token": authToken}
}
