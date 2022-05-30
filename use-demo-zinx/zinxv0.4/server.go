package main

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"github.com/bailuoxi66/go-inside/demo-zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping\n"))
	if err != nil {
		fmt.Println("Call Router PreHandle err:", err)
	}
}

// Handle 处理conn业务的钩子方法hook
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping\n"))
	if err != nil {
		fmt.Println("Call Router Handle err:", err)
	}
}

// PostHandle 在处理conn业务之后的钩子方法hook
func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
	if err != nil {
		fmt.Println("Call Router PostHandle err:", err)
	}
}

/*
	基于demo-zinx框架来开发的,服务器端应用程序
*/
func main() {
	// 1.创建一个server句柄，使用demo-zinx api
	s := znet.NewServer("[zinx v0.3]")

	// 2.给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})

	// 3.启动server
	s.Server()
}
