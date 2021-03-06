package znet

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/utils"
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
	// 当前server的消息管理模块，用来绑定MsgID和对应的处理业务API的关系
	MsgHandle ziface.IMsgHandler
	// 该server的链接管理器
	ConMgr ziface.IConnManager
	// 该server创建链接之后自动调用Hook函数
	OnConnStart func(conn ziface.IConnection)
	// 该server销毁链接之前自动调用Hook函数
	OnConnStop func(conn ziface.IConnection)
}

// Start 启动服务器
func (s *Server) Start() {

	fmt.Printf("[Zinx] Server Name:%s, listenner at IP:%s, Port:%d is starting...\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version:%s, MaxConn:%d, MaxPacketSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

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
			// 开启消息队列及worker工作池
			s.MsgHandle.StartWorkerPool()
			// 如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}

			// 设置最大链接个数的判断，如果超过最大链接，则关闭此新的链接
			if s.ConMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端相应一个超出最大链接的错误包
				fmt.Println("Too Many Connections MaxConn=", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			// 将处理新链接的业务方法和conn 进行绑定 得到我们的链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandle)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {

	// 将一些服务器的资源、状态或者一些已经开辟的链接信息 进行停止或者回收
	fmt.Println("[STOP]Zinx server name:", s.Name)
	s.ConMgr.ClearConn()
}

// Server 运行服务器
func (s *Server) Server() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务
	// 阻塞
	select {}
}

// AddRouter 路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgID, router)
	fmt.Println("AddRouter Success!...")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConMgr
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
		ConMgr:    NewConnManager(),
	}
	return s
}

// SetOnConnStart 注册OnConnStart 钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 注册OnConnStop 钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用OnConnStart 钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("--------------> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用OnConnStop 钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("--------------> Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
