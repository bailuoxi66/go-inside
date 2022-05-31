package main

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"github.com/bailuoxi66/go-inside/demo-zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

// Handle 处理conn业务的钩子方法hook
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")

	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Printf("recv from client:msgID=", request.GetMsgID(),
		", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Printf("err:", err)
	}
}

/*
	基于demo-zinx框架来开发的,服务器端应用程序
*/
func main() {
	// 1.创建一个server句柄，使用demo-zinx api
	s := znet.NewServer("[zinx v0.5]")

	// 2.给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})

	// 3.启动server
	s.Server()
}
