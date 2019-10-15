package ApiHTTP

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
	"net/http"
	"time"

	"../ApiRPC"
	"../RpcClientPbs/mysqlPb"
	"../utils"

	conf "../Config"
)

const (
	JWTGetUserId = "user_id"
)

// Just used to validate a basic type value, not struct instance.
// The main use for validate the query string params
var (
	simpleFieldValidate *validator.Validate
)

func init() {
	config := &validator.Config{TagName: "validate"}
	simpleFieldValidate = validator.New(config)

}

// Sign Up data struct , all field required.
// Verify data by the validators of gin binding.
type UserSignUp struct {
	Name            string `json:"name" binding:"nameValidator"`
	Email           string `json:"email"  binding:"emailValidator"`
	Password        string `json:"password"  binding:"passwordValidator"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// SignUp(Register) HTTP API function.
// Saving the user's basic information which the user input and set the default value for some value like gender,
// and create the qr code image for this user.
func SignUp(c *gin.Context) {
	var err error
	var tempUserSignUp = UserSignUp{}

	if err = c.ShouldBindJSON(&tempUserSignUp); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// initial some data
	passwordHash := utils.GetPasswordHash(tempUserSignUp.Password, conf.PasswordHashSalt)
	qrCodePicData, _ := utils.CreatQRCodeBytes("user_email=" + tempUserSignUp.Email)
	qrCodePicHashName := utils.BytesDataHash(qrCodePicData)

	// save user information to db
	userBasic, err := ApiRPC.SaveOneNewUser(
		tempUserSignUp.Name,
		tempUserSignUp.Email,
		"",
		passwordHash,
		conf.DefaultAvatarPicName,
		qrCodePicHashName,
		0)

	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// save the qr code pic
	_ = SaveQRCodeLocal(qrCodePicData, qrCodePicHashName)

	HidePasswordAndCompleteAvatarAndQrCodeURL(userBasic)
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

	HidePasswordAndCompleteAvatarAndQrCodeURL(userBasic)
	//verify ok, return user detail and AuthToken
	c.JSON(http.StatusOK, detailAndToken(userBasic))
}

// Create auth token by user, and return data
func detailAndToken(user *mysqlPb.UserBasic) gin.H {
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

//Hide password and complete the URL of avatar and qr code
func HidePasswordAndCompleteAvatarAndQrCodeURL(userBasic *mysqlPb.UserBasic) {
	userBasic.Password = ""
	userBasic.QrCode = conf.PhotosUrlPrefix + userBasic.QrCode + conf.PhotoSuffix
	userBasic.Avatar = conf.PhotosUrlPrefix + userBasic.Avatar + conf.PhotoSuffix
}

// GetResetPasswordEmail HTTP API function
// Send a email which content with the `Auth-Token` to the user's email box for reset password.
func GetResetPasswordEmail(c *gin.Context) {
	tempEmail := c.Request.URL.Query().Get("email")
	if tempEmail == "" {
		c.JSON(400, gin.H{"error": "query string param: email is required"})
		return
	}
	err := simpleFieldValidate.Field(tempEmail, "email,lte=100")
	if nil != err {
		c.JSON(400, gin.H{"error": "param: email validate failed"})
		return
	}
	data, err := ApiRPC.GetUserByEmail(tempEmail)
	if nil != err {
		c.JSON(400, gin.H{"error": "email not registered"})
		return
	}
	if ok := utils.CheckEmailSentIn3Minute(tempEmail); ok {
		c.JSON(200, gin.H{"message": "the email was sent, please check your email box"})
		return
	}
	// send the reset password authentication link email
	claims := utils.CustomJWTClaims{
		Id: data.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + conf.ResetPasswordTokenAliveTIme), // expire time
			Issuer:    conf.AuthTokenIssuer,                                        //signal issuer
		},
	}
	authToken, _ := utils.CreateJWTToken(claims, []byte(conf.AuthTokenSalt))
	go func() {
		err := utils.SendPasswordResetEmail(tempEmail, authToken)
		if err != nil {
			errInfo := fmt.Sprintf("send password reset email to user(%s) fail: %s", tempEmail, err.Error())
			GlobalGinStyleLogger.Fprintln(c, 500, errInfo)
		}
	}()

	c.Status(200)

}
