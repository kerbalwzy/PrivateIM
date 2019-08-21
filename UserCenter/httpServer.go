package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"

	"./ApiHTTP"
	_ "./models"
)

const StaticFoldPath = "/Users/wzy/GitProrgram/PrivateIM/UserCenter/static/"

func StartHttpServer() {
	r := gin.Default()
	r.Static("/static", StaticFoldPath)

	auth := r.Group("/auth")
	auth.POST("/user", ApiHTTP.SignUp)
	auth.POST("/profile", ApiHTTP.SignIn)

	info := r.Group("/info", ApiHTTP.JWTAuthMiddleware())
	info.GET("/profile", ApiHTTP.GetProfile)
	info.PUT("/profile", ApiHTTP.PutProfile)
	info.GET("/avatar", ApiHTTP.GetAvatar)
	info.PUT("/avatar", ApiHTTP.PutAvatar)
	info.GET("/qrcode", ApiHTTP.GetQrCode)
	info.POST("/qrcode", ApiHTTP.ParseQrCode)

	relate := r.Group("/relation", ApiHTTP.JWTAuthMiddleware())
	relate.GET("/friend", ApiHTTP.GetFriend)
	relate.POST("/friend", ApiHTTP.AddFriend)
	relate.PUT("/friend", ApiHTTP.PutFriend)
	relate.DELETE("/friend", ApiHTTP.DeleteFriend)
	relate.GET("/friends", ApiHTTP.AllFriends)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("nameValidator", ApiHTTP.NameValidator)
		_ = v.RegisterValidation("emailValidator", ApiHTTP.EmailValidator)
		_ = v.RegisterValidation("mobileValidator", ApiHTTP.MobileValidator)
		_ = v.RegisterValidation("passwordValidator", ApiHTTP.PasswordValidator)
		_ = v.RegisterValidation("genderValidator", ApiHTTP.GenderValidator)
		_ = v.RegisterValidation("relateActionValidator", ApiHTTP.RelateActionValidator)
	} else {
		log.Fatal("binding custom validators fail!!!")
	}

	err := r.Run(":8080")
	if nil != err {
		log.Fatal(err)
	}
}
