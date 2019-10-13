package ApiHTTP

//
import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"

	//"encoding/json"
	//"errors"
	//"fmt"
	//"github.com/dgrijalva/jwt-go"
	//"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
	//"gopkg.in/go-playground/validator.v8"
	"io/ioutil"
	//"time"

	"../ApiRPC"
	"../utils"

	conf "../Config"
)

//
// GetProfile HTTP API function
func GetProfile(c *gin.Context) {
	userId := c.MustGet(JWTGetUserId).(int64)
	userBasic, err := ApiRPC.GetUserById(userId)
	if nil != err {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	HidePasswordAndCompleteAvatarAndQrCodeURL(userBasic)
	c.JSON(200, userBasic)
}

type TempProfile struct {
	Name   string `json:"name" binding:"nameValidator"`
	Mobile string `json:"mobile" binding:"mobileValidator"`
	Gender int    `json:"gender" binding:"genderValidator"`
}

var (
	ErrNoJsonBody         = errors.New("not have any JsonBodyParams")
	ERrPutProfileNoName   = errors.New("`name` not exited in JsonBodyParams")
	ERrPutProfileNoMobile = errors.New("`mobile` not exited in JsonBodyParams")
	ERrPutProfileNoGender = errors.New("`gender` not exited in JsonBodyParams")
)
// Parse the JsonBodyParams to TempProfile
func parseTempProfile(c *gin.Context) (*TempProfile, error) {
	// Parse the JsonBodyParams to map
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	if n == 0 {
		return nil, ErrNoJsonBody
	}

	tempDict := make(map[string]interface{})
	_ = json.Unmarshal(buf[0:n], &tempDict)
	// Check the integrity of parameters
	name, ok := tempDict["name"]
	if !ok {
		return nil, ERrPutProfileNoName
	}
	mobile, ok := tempDict["mobile"]
	if !ok {
		return nil, ERrPutProfileNoMobile
	}
	gender, ok := tempDict["gender"]
	if !ok {
		return nil, ERrPutProfileNoGender
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
	err = binding.Validator.ValidateStruct(tempProfile)
	if nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Update user info
	userId := c.MustGet(JWTGetUserId).(int64)
	userBasic, err := ApiRPC.PutUserProfileById(
		userId,
		tempProfile.Name,
		tempProfile.Mobile,
		tempProfile.Gender)

	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, userBasic)
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
	oldPasswordHash, err := ApiRPC.GetUserPasswordById(userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if oldPasswordHash.Value != utils.GetPasswordHash(temp.OldPassword, conf.PasswordHashSalt) {
		c.JSON(400, gin.H{"error": "old password error"})
		return
	}
	err = ApiRPC.PutUserPasswordById(utils.GetPasswordHash(temp.Password, conf.PasswordHashSalt), userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
	}
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
	err := ApiRPC.PutUserPasswordById(passwordHash, userId)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)

}

//
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

// save QRCode file to local
func SaveQRCodeLocal(data []byte, hashName string) error {

	savePath := conf.PhotoSaveFoldPath + hashName + conf.PhotoSuffix
	err := ioutil.WriteFile(savePath, data, 0644)
	if nil != err {
		log.Printf("@@@@@@@%s", err.Error())
		return err
	}
	return nil
}

// todo upload the data to cloud with a hashName
func UploadDataToCloud(data []byte, hashName string) error {
	return nil
}
