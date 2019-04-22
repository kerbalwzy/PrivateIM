/*
问答游戏-网页版-服务器程序-程序主入口
author	kerbalwzy@gmail.com
*/
package main

import (
	"config"
	"fmt"
)

func init() {
	Setting := config.LoadConfigFile("./config.json")
	fmt.Printf("%s\n", Setting.ServerConfig.Port)
	fmt.Printf("%s\n", Setting.LogConfig.Global.OutPut[1])
	fmt.Printf("%s\n", Setting.LogConfig.File)

}

func main() {

}

func TempPrint(args ...interface{}) {
	fmtStr := ""
	for range args {
		fmtStr += "%s\n---\n"
	}
	fmt.Printf(fmtStr, args...)

}
