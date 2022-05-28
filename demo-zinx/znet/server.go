package znet

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"net"
)

// Server iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
}

// Start 启动服务器
func (s *Server) Start() {

	fmt.Printf("[Start] Server Listener at IP:%s, Port:%d  is starting\n", s.IP, s.Port)

	// 为了不让Start阻塞在Accept
	go func() {
		// 1. 获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve tcp addr err:", err)
			return
		}
		// 2. 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("ListenTCP err:", err)
		}
		fmt.Println("Start Zinx server succ, ", s.Name, "Listen...")
		// 3. 阻塞的等待客户端链接，处理客户端链接业务（读写）
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}

			// 已经与客户端建立链接，做一些业务，做一个最基本的最大512字节长度的回显业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("Read buf err:", err)
						continue
					}

					// 回显功能
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("Write back buf err:", err)
						continue
					}
				}
			}()
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {

	// TODO 将一些服务器的资源、状态或者一些已经开辟的链接信息 进行停止或者回收
}

// Server 运行服务器
func (s *Server) Server() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务
	// 阻塞
	select {}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
