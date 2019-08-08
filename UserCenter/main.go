package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"./controllers"
	_ "./models"
)

var ConfigMap = make(map[string]string)

func init() {
	log.SetPrefix("AuthCenter ")

	//path, _ := filepath.Abs(os.Args[0])
	//BaseDir := filepath.Dir(path)
	//ConfigMap["SQLite3_ADDR"] = filepath.Join(BaseDir, "userDb.SQLite3")

	ConfigMap["SQLite3_ADDR"] = "/Users/wzy/GitProrgram/PrivateIM/UserCenter/userDb.SQLite3"
}

func main() {
	router := gin.Default()

	router.POST("/user", controllers.SignUp)
	router.POST("/profile", controllers.SignIn)
	router.Run(":8080")
}
