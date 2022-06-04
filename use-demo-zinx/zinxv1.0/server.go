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
	fmt.Println("recv from client:msgID=", request.GetMsgID(),
		", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Printf("err:", err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

// Handle 处理conn业务的钩子方法hook
func (hzr *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")

	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client:msgID=", request.GetMsgID(),
		", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello Welcome Zinx..."))
	if err != nil {
		fmt.Printf("err:", err)
	}
}

// DoConnectionBegin 创建链接之后执行钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("======> DoConnectionBegin is Called...")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}

	// 给当前链接设置一些属性
	fmt.Println("Set conn Name, home..")
	conn.SetProperty("Name", "落渔-lane")
	conn.SetProperty("Github", "https://github.com/bailuoxi66/go-inside")
}

// DoConnectionLost 链接端口之后执行钩子函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("======> DoConnectionLost is Called...")
	fmt.Println("conn ID: ", conn.GetConnID(), " is Lost...")

	// 获取链接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
	if github, err := conn.GetProperty("Github"); err == nil {
		fmt.Println("Github = ", github)
	}
}

/*
	基于demo-zinx框架来开发的,服务器端应用程序
*/
func main() {
	// 1.创建一个server句柄，使用demo-zinx api
	s := znet.NewServer("[zinx v0.8]")

	// 2.注册链接Hook的钩子函数
	// 因为函数在go语言里面也是一种数据类型
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 3.给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	// 3.启动server
	s.Server()
}
