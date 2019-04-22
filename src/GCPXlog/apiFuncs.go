/*
生成日志处理对象
*/
package GCPXlog

import (
	"log"
)

// 定义Logger结构体对象
type Logger struct {
	writerSlice []LogWriter // to save writers
}

// 给logger对象添加writer
func (logger *Logger) RegisterWriter(writer LogWriter) {
	logger.writerSlice = append(logger.writerSlice, writer)
}

func (logger *Logger) Error(fmtMsg string, args ...interface{}) {
	log.Printf(fmtMsg, args...)
}

func (logger *Logger) Info(fmtMsg string, args ...interface{}) {
	log.Printf(fmtMsg, args...)
}

// 定义一个写入器接口，因为未来可能会有多个写入器共同工作
type LogWriter interface {
	Write(data interface{}) error
}
