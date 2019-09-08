package ApiHTTP

import (
	"../ApiRPC"
	"../utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"io/ioutil"
	"time"

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

// GetProfile HTTP API function
func GetProfile(c *gin.Context) {
	userId := c.MustGet(JWTGetUserId).(int64)
	userBasic, err := ApiRPC.GetUserById(userId)
	if nil != err {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	userBasic.Password = ""
	c.JSON(200, userBasic)
}

type TempProfile struct {
	Name   string `json:"name" binding:"nameValidator"`
	Mobile string `json:"mobile" binding:"mobileValidator"`
	Gender int    `json:"gender" binding:"genderValidator"`
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

// PutProfile HTTP API function
func PutProfile(c *gin.Context) {
	// Validate the params
	tempProfile, err := parseTempProfile(c)
	if nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	errs := binding.Validator.ValidateStruct(tempProfile)
	if nil != errs {
		c.JSON(400, gin.H{"errors": errs.Error()})
		return
	}
	// Update user info
	userId := c.MustGet(JWTGetUserId).(int64)
	userBasic, err := ApiRPC.PutUserBasicById(
		tempProfile.Name, tempProfile.Mobile, tempProfile.Gender, userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userBasic.Password = ""
	c.JSON(200, userBasic)
}

// GetAvatar HTTP API function
func GetAvatar(c *gin.Context) {
	userId := c.MustGet(JWTGetUserId).(int64)
	data, err := ApiRPC.GetUserAvatarById(userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if data.Avatar == "" {
		c.JSON(200, gin.H{"avatar_url": conf.DefaultAvatarUrl})
		return
	}
	c.JSON(200, gin.H{"avatar_url": conf.PhotosUrlPrefix + data.Avatar + conf.PhotoSuffix})
}

// PutAvatar HTTP API function
func PutAvatar(c *gin.Context) {
	// get file data and hash value as name
	file, err := c.FormFile("new_avatar")
	if nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if file.Size > conf.AvatarPicUploadMaxSize || file.Size == 0 {
		c.JSON(400, gin.H{"error": "the upload image size need gt=0kb and lte=100kb"})
		return
	}
	hashName, data, err := utils.GinFormFileHash(file)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := UploadAvatarLocal(data, hashName); nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// save the information into database
	userId := c.MustGet(JWTGetUserId).(int64)
	userAvatar, err := ApiRPC.PutUserAvatarById(hashName, userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"avatar_url": conf.PhotosUrlPrefix + userAvatar.Avatar + conf.PhotoSuffix})
}

type TempPassword struct {
	OldPassword     string `json:"old_password" binding:"passwordValidator"`
	Password        string `json:"password" binding:"passwordValidator,nefield=OldPassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// PutPassword HTTP API function
func PutPassword(c *gin.Context) {
	temp := new(TempPassword)
	if err := c.ShouldBindJSON(temp); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userId := c.MustGet(JWTGetUserId).(int64)
	userBasic, err := ApiRPC.GetUserById(userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if userBasic.Password != utils.GetPasswordHash(temp.OldPassword, conf.PasswordHashSalt) {
		c.JSON(400, gin.H{"error": "old password error"})
		return
	}
	_, err = ApiRPC.PutUserPasswordById(utils.GetPasswordHash(temp.Password, conf.PasswordHashSalt), userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.Status(200)
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
	go utils.SendPasswordResetEmail(tempEmail, authToken)
	c.Status(200)
}

type TempForgetPassword struct {
	Password        string `json:"password" binding:"passwordValidator"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// ForgetPassword HTTP API function
// the auth-token is from the `reset password email`
func ForgetPassword(c *gin.Context) {
	temp := new(TempForgetPassword)
	if err := c.ShouldBindJSON(temp); nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userId := c.MustGet(JWTGetUserId).(int64)
	passwordHash := utils.GetPasswordHash(temp.Password, conf.PasswordHashSalt)
	_, err := ApiRPC.PutUserPasswordById(passwordHash, userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)

}

// GetQRCode HTT API function
func GetQrCode(c *gin.Context) {
	// try to get QrCode hash name from database. if existed, return.
	userId := c.MustGet(JWTGetUserId).(int64)
	userQrCode, err := ApiRPC.GetUserQRCodeById(userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if userQrCode.QrCode != "" {
		c.JSON(200, gin.H{"qr_code": conf.PhotosUrlPrefix + userQrCode.QrCode + conf.PhotoSuffix})
		return
	}

	// if the qr code hash name is not existed, create an new and save
	queryStrParam := fmt.Sprintf("user=%d", userId)
	content := QRCodeContent(queryStrParam)
	data, _ := utils.CreatQRCodeBytes(content)
	hashName := utils.BytesDataHash(data)
	err = SaveQRCodeLocal(data, hashName)
	if nil != err {
		c.JSON(500, gin.H{"error": "create QRCode fail"})
		return
	}
	userQrCode, err = ApiRPC.PutUserQRCodeById(hashName, userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"qr_code_url": conf.PhotosUrlPrefix + userQrCode.QrCode + conf.PhotoSuffix})

}

// ParseQrCode HTTP API function
func ParseQrCode(c *gin.Context) {
	// get file data
	file, err := c.FormFile("qr_code")
	if nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if file.Size == 0 {
		c.JSON(400, gin.H{"error": "the upload image size need gt=0kb and lte=2MB"})
		return
	}
	_, data, err := utils.GinFormFileHash(file)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	content, err := utils.ParseQRCodeBytes(data)
	if nil != err {
		c.JSON(400, gin.H{"error": "QRCode parse fail"})
		return
	}
	c.JSON(200, gin.H{"qr_content": content})

}

// Save avatar file to local dir
func UploadAvatarLocal(data []byte, hashName string) error {
	path := conf.PhotoSaveFoldPath + hashName + conf.PhotoSuffix
	if err := utils.UploadFileToLocal(data, path); nil != err {
		return err
	}
	return nil
}

// make the content for create a QRCode, infect the content is a query string param.
func QRCodeContent(content string) string {
	return conf.QRCodeBaseUrl + content
}

// save QRCode file to local
func SaveQRCodeLocal(data []byte, hashName string) error {
	savePath := conf.PhotoSaveFoldPath + hashName + conf.PhotoSuffix
	err := ioutil.WriteFile(savePath, data, 0644)
	if nil != err {
		return err
	}
	return nil
}

// todo upload the data to cloud with a hashName
func UploadDataToCloud(data []byte, hashName string) error {
	return nil
}
