package main

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"log"

	"github.com/gin-gonic/gin"

	"./controllers"
	_ "./models"
)

func init() {
	log.SetPrefix("AuthCenter ")

}

const StaticFoldPath = "/Users/wzy/GitProrgram/PrivateIM/UserCenter/static/"

func main() {
	r := gin.Default()
	r.Static("/static", StaticFoldPath)

	auth := r.Group("/auth")
	auth.POST("/user", controllers.SignUp)
	auth.POST("/profile", controllers.SignIn)

	info := r.Group("/info", controllers.JWTAuthMiddleware())
	info.GET("/profile", controllers.GetProfile)
	info.PUT("/profile", controllers.PutProfile)
	info.GET("/avatar", controllers.GetAvatar)
	info.PUT("/avatar", controllers.PutAvatar)
	info.GET("/qrcode", controllers.GetQrCode)
	info.POST("/qrcode", controllers.ParseQrCode)

	relate := r.Group("/relation", controllers.JWTAuthMiddleware())
	relate.GET("/friend", controllers.GetFriend)
	relate.POST("/friend", controllers.AddFriend)
	relate.PUT("/friend", controllers.PutFriend)
	relate.DELETE("/friend", controllers.DeleteFriend)
	relate.GET("/friends", controllers.AllFriends)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("nameValidator", controllers.NameValidator)
		_ = v.RegisterValidation("emailValidator", controllers.EmailValidator)
		_ = v.RegisterValidation("mobileValidator", controllers.MobileValidator)
		_ = v.RegisterValidation("passwordValidator", controllers.PasswordValidator)
		_ = v.RegisterValidation("genderValidator", controllers.GenderValidator)
		_ = v.RegisterValidation("relateActionValidator", controllers.RelateActionValidator)
	} else {
		log.Fatal("binding custom validators fail!!!")
	}

	err := r.Run(":8080")
	if nil != err {
		log.Fatal(err)
	}
}
