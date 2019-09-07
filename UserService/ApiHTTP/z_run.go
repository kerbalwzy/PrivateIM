package ApiHTTP

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"

	conf "../Config"
)

func StartHttpServer() {
	r := gin.Default()
	r.Static("/static", conf.StaticFoldPath)

	auth := r.Group("/auth")
	auth.POST("/user", SignUp)
	auth.POST("/profile", SignIn)

	info := r.Group("/info", JWTAuthMiddleware())
	info.GET("/profile", GetProfile)
	info.PUT("/profile", PutProfile)
	info.PUT("/password", PutPassword)
	info.GET("/password", GetResetPasswordEmail)
	info.POST("/password", )
	info.GET("/avatar", GetAvatar)
	info.PUT("/avatar", PutAvatar)
	info.GET("/qrcode", GetQrCode)
	info.POST("/qrcode", ParseQrCode)

	relate := r.Group("/relation", JWTAuthMiddleware())
	relate.GET("/friend", GetFriend)
	relate.POST("/friend", AddFriend)
	relate.PUT("/friend", PutFriend)
	relate.DELETE("/friend", DeleteFriend)
	relate.GET("/friends", AllFriends)

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

	err := r.Run(":8080")
	if nil != err {
		log.Fatal(err)
	}
}
