package main

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"log"

	"github.com/gin-gonic/gin"

	"./controllers"
	_ "./models"
	"./utils"
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
	auth.POST("profile", controllers.SignIn)

	info := r.Group("/info", controllers.JWTAuthMiddleware())
	info.GET("/profile", controllers.GetProfile)
	info.PUT("/profile", controllers.PutProfile)
	info.GET("/avatar", controllers.GetAvatar)
	info.PUT("/avatar", controllers.PutAvatar)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("nameValidator", utils.NameValidator)
		_ = v.RegisterValidation("emailValidator", utils.EmailValidator)
		_ = v.RegisterValidation("mobileValidator", utils.MobileValidator)
		_ = v.RegisterValidation("passwordValidator", utils.PasswordValidator)
		_ = v.RegisterValidation("genderValidator", utils.GenderValidator)
	} else {
		log.Fatal("binding custom validators fail!!!")
	}

	err := r.Run(":8080")
	if nil != err {
		log.Fatal(err)
	}
}
