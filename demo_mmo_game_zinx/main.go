package main

import "github.com/bailuoxi66/go-inside/demo-zinx/znet"

func main() {
	// 创建zinx server句柄
	s := znet.NewServer("MMO Game Zinx")

	// 连接创建和销毁的HOOK钩子函数
	// 注册一些路由业务
	// 启动服务
	s.Server()
}
