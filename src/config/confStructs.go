package config

var Setting Config

// 创建配置对象结构体
type Config struct {
	LogConfig    LogConfig    `json:"logConfig"`
	ServerConfig ServerConfig `json:"serverConfig"`
}

type LogConfig struct {
	Global  Global  `json:"global"`
	Console Console `json:"console"`
	File    File    `json:"file"`
	Mongodb Mongodb `json:"mongodb"`
}

type Global struct {
	OutPut       []string `json:"outPut"`
	GlobalFormat string   `json:"globalFormat"`
	GlobalLevel  string   `json:"globalLevel"`
}

type Console struct {
	basicFormat
}

type File struct {
	basicFormat
}

type Mongodb struct {
	basicFormat
}

type basicFormat struct {
	Format string `json:"format"`
	Level  string `json:"level"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}
