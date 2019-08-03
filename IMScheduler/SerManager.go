package main

import (
	"log"
	"net"
)

/*
业务注册中心
功能描述:	接收业务程序的请求，将业务注册到调度中心Map
通信描述:	RPC TCP Gob
*/

// 服务管理对象结构体
type Manager struct {
	WorkerMap map[SourceTag][]Worker
	listener  net.Listener
	sender    net.Dialer
}

// 接收业务服务的注册请求，将业务
func (m *Manager) Register(wkr Worker, token *AuthToken) error {
	log.Printf("Get register request from Worker@(%s)", wkr.SelfTCPAdder)
	//*token = MakeToken()
	return nil
}
