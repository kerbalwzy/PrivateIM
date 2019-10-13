package ApiHTTP

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"

	conf "../Config"
)

// the source code of the gin.StyleLogger is add by myself.
// the source code explanation: https://github.com/gin-gonic/gin/pull/2096
var GlobalGinStyleLogger *gin.StyleLogger

func init() {
	GlobalGinStyleLogger = gin.NewGinStyleLogger(nil, nil)
}

func StartHttpServer() {
	r := gin.Default()
	r.Static("/static", conf.StaticFoldPath)

	auth := r.Group("/auth")
	auth.POST("/user", SignUp)
	auth.POST("/profile", SignIn)
	auth.GET("/password", GetResetPasswordEmail)

	info := r.Group("/info", JWTAuthMiddleware())
	info.GET("/profile", GetProfile)
	info.PUT("/profile", PutProfile)
	info.PUT("/avatar", PutAvatar)

	info.PUT("/password", PutPassword)
	info.POST("/password", ForgetPassword)

	info.POST("/qr_code", ParseQrCode)
	//
	relate := r.Group("/relation", JWTAuthMiddleware())
	relate.GET("/users", SearchUsers)

	relate.POST("/friend", AddFriend)
	relate.PUT("/friend", PutFriend)
	relate.GET("/friends", GetUsersFriendsInfo)
	relate.DELETE("/friend", DeleteFriend)



	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("nameValidator", NameValidator)
		_ = v.RegisterValidation("emailValidator", EmailValidator)
		_ = v.RegisterValidation("mobileValidator", MobileValidator)
		_ = v.RegisterValidation("passwordValidator", PasswordValidator)
		_ = v.RegisterValidation("genderValidator", GenderValidator)
		_ = v.RegisterValidation("relateActionValidator", RelateActionValidator)
	} else {
		log.Fatal("binding custom validators fail!!!")
	}

	err := r.Run(conf.UserCenterHttpServerAddress)
	if nil != err {
		log.Fatal(err)
	}
}
