package znet

import (
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"net"
)

type Connection struct {
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn
	// 链接ID
	ConnID int32
	// 当前的链接状态
	isClosed bool
	// 当前链接所绑定的处理业务方法API
	handleAPI ziface.HandleFunc
	// 告知当前链接已经退出、停止的channel
	ExitChan chan bool
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID int32, callbackAPI ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callbackAPI,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}

	return c
}
