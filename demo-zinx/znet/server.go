package znet

import (
	"errors"
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

// CallBackToClient 定义当前客户端链接所绑定的handle api（目前这个handle是写死的，以后优化应该由用户自定义handle方法）
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	// 回显的业务
	fmt.Println("[Conn Handle] CallBackToClient")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err:", err)
		return errors.New("CallBackToClient error")
	}

	return nil
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

		var cid uint32
		cid = 0
		// 3. 阻塞的等待客户端链接，处理客户端链接业务（读写）
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}
			// 将处理新链接的业务方法和conn 进行绑定 得到我们的链接模块
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
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
