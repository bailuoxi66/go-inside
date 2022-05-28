package main

import "github.com/bailuoxi66/go-inside/demo-zinx/znet"

/**
基于demo-zinx框架来开发的,服务器端应用程序
*/
func main() {
	// 1.创建一个server句柄，使用demo-zinx api
	s := znet.NewServer("[zinx v0.1]")

	// 2.启动server
	s.Server()
}
