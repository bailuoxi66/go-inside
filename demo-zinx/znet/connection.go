package znet

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"net"
)

type Connection struct {
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn
	// 链接ID
	ConnID uint32
	// 当前的链接状态
	isClosed bool
	// 当前链接所绑定的处理业务方法API
	handleAPI ziface.HandleFunc
	// 告知当前链接已经退出、停止的channel
	ExitChan chan bool
}

func (c *Connection) StartReader() {
	fmt.Println("StartReader Goroutine is running...")
	defer fmt.Println("ConnID", c.ConnID, " StartReader is exit, remote addr is :", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err:", err)
			continue
		}

		// 调用当前所绑定的HandleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID:", c.ConnID, " handle is error", err)
			break
		}
	}
}

// Start 启动链接，让当前链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID:", c.ConnID)

	// 启动从当前链接的读数据的业务
	// TODO 启动从写前链接的读数据的业务
}

// Stop 停止链接，让当前链接结束当前的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID:", c.ConnID)

	// 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true
	// 关闭socket链接
	c.Conn.Close()

	// 回收资源
	close(c.ExitChan)
}

// GetTCPConnection 获取当前链接的绑定 socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// GetConnID 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// Send 获取远程客户端的TCP状态 IP Port
func (c *Connection) Send(data []byte) error {
	return nil
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callbackAPI ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callbackAPI,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}

	return c
}
