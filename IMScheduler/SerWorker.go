package main

import (
	"log"
	"net"
	"net/rpc"
)

/*
业务服务工作对象结构体
	SourceTag		对象在调度中心的资源标记
	Token			验证来自调度中心的Token，如果Token不正确表示任务不是来自我主动注册的任务中心
	ManagerTCPAdder	调度中心TCP地址

	listener		用来监听来自调度中心派发的任务的TCP监听对象
*/
type Worker struct {
	Tag             SourceTag
	Token           AuthToken
	ManagerTCPAdder string
	SelfTCPAdder    string

	listener net.Listener
}

func (wkr Worker) RegisterSelf() {
	client, err := rpc.Dial("tcp", wkr.ManagerTCPAdder)
	CheckErrFatal(err)
	err = client.Call("Manager.Register", wkr, &(wkr.Token))
	CheckErrFatal(err)
	log.Printf("Register self to Manager@(%s)\n", wkr.ManagerTCPAdder)
	log.Print(wkr.Token)
}
