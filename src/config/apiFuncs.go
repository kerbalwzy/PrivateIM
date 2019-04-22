/*
生成配置处理对象
*/
package config

import (
	"bytes"
	"encoding/json"
	"os"
	"regexp"

	"GCPXlog"
)

// 读取JSON格式的配置文件，并将配置信息保存到config对象
func LoadConfigFile(filePath string) Config {
	var logger GCPXlog.Logger
	var config Config
	// open config file
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("配置文件<%s>打开失败 INFO:%s", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	// get file information
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		logger.Error("配置文件<%s>内容为空！")
		os.Exit(1)
	}

	// get config content from the file
	buffer := make([]byte, fileInfo.Size())
	_, err = file.Read(buffer)

	// remove Notes from the config byte slice
	buffer, err = StripNotes(buffer)
	if err != nil {
		logger.Error("去除配置文件<%s>中的注释信息失败 INFO:%s", filePath, err)
		os.Exit(1)
	}
	buffer = []byte(os.ExpandEnv(string(buffer)))

	// pares json date and give it to config object
	err = json.Unmarshal(buffer, &config) //解析json格式数据
	if err != nil {
		logger.Error("配置文件<%s>解析失败 INFO:%s", filePath, err)
		os.Exit(1)
	}

	logger.Info("程序初始化 - - 加载配置文件成功")
	return config
}

// 去除掉配置文件中可能存在的以"#"号开头的注释内容
func StripNotes(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0) // 去除掉Windows系统中使用时可能会出现的"\r"
	lines := bytes.Split(data, []byte("\n"))
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}
	return bytes.Join(filtered, []byte("\n")), nil
}
