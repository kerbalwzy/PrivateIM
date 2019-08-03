package main

import (
	"net"
	"net/rpc"
)

/*
调度中心：
核心功能一:
	功能描述:	接收来自其业务程序的注册请求，分析业务资源类型，保存业务功能信息，并开始为注册的业务提供代理和调度功能
	通信描述:	RPC	TCP Gob

核心功能二：
	功能描述:	接收来自客户端的业务请求，将业务转发给相应的业务服务程序，接收服务返回的响应，然后再返回给客户端
	通信描述:	HTTP RestAp JSON
*/

func main() {
	//startManagerRpcSer()
	//startAnWorker()
}
func startAnWorker() {
	var worker = Worker{ManagerTCPAdder: "127.0.0.1:1234",
		SelfTCPAdder: "127.0.0.1:2222"}
	worker.RegisterSelf()
}

func startManagerRpcSer() {
	manager := new(Manager)
	err := rpc.Register(manager)
	CheckErrFatal(err)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	CheckErrFatal(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	CheckErrFatal(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)

	}
}
