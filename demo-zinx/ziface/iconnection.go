package ziface

import "net"

// IConnection 定义链接模块的抽象层
type IConnection interface {

	// Start 启动链接，让当前链接准备开始工作
	Start()
	// Stop 停止链接，让当前链接结束当前的工作
	Stop()
	// GetTCPConnection 获取当前链接的绑定 socket conn
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前链接模块的链接ID
	GetConnID() uint32
	// Send 获取远程客户端的TCP状态 IP Port
	Send(data []byte) error
}

// HandleFunc 定义一个处理链接业务的方法
type HandleFunc func(net.TCPConn, []byte, int) error
