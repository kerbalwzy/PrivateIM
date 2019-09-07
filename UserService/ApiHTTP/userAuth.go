package ApiHTTP

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"../ApiRPC"
	"../utils"

	conf "../Config"
	pb "../Protos"
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
	var tempUserSignUp = UserSignUp{}

	if err = c.ShouldBindJSON(&tempUserSignUp); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// check the email if registered
	userBasic, err := ApiRPC.GetUserByEmail(tempUserSignUp.Email)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else if userBasic.Id != 0 {
		c.JSON(400, gin.H{"error": "email is already sign up, please sign in"})
		return
	}

	// save user information to database
	passwordHash := utils.GetPasswordHash(tempUserSignUp.Password, conf.PasswordHashSalt)
	userBasic, err = ApiRPC.SaveOneNewUser(tempUserSignUp.Name, tempUserSignUp.Email,
		"", passwordHash, 0)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	//ok, return user detail and auth token
	c.JSON(201, detailAndToken(userBasic))
}

type UserSignIn struct {
	Email    string `json:"email"  binding:"emailValidator"`
	Password string `json:"password"  binding:"passwordValidator"`
}

// SignIn(Login) HTTP API function
func SignIn(c *gin.Context) {
	var err error
	var temp = UserSignIn{}

	if err = c.ShouldBindJSON(&temp); nil != err {
		c.JSON(400, gin.H{"error": "verify fail"})
		return
	}

	// check email and password for user

	userBasic, err := ApiRPC.GetUserByEmail(temp.Email)
	if nil != err {
		c.JSON(400, gin.H{"error": "verify fail, email or password error"})
		return
	}
	if userBasic.Password != utils.GetPasswordHash(temp.Password, conf.PasswordHashSalt) {
		c.JSON(400, gin.H{"error": "verify fail, email or password error"})
		return
	}

	//verify ok, return user detail and AuthToken
	c.JSON(http.StatusOK, detailAndToken(userBasic))
}

// Create auth token by user, and return data
func detailAndToken(user *pb.UserBasicInfo) gin.H {
	claims := utils.CustomJWTClaims{
		Id: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + conf.AuthTokenAliveTime), // expire time
			Issuer:    conf.AuthTokenIssuer,                               //signal issuer
		},
	}
	authToken, _ := utils.CreateJWTToken(claims, []byte(conf.AuthTokenSalt))
	return gin.H{"user": user, "auth_token": authToken}
}
