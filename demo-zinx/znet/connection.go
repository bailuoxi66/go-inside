package znet

import (
	"errors"
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"io"
	"net"
)

type Connection struct {
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn
	// 链接ID
	ConnID uint32
	// 当前的链接状态
	isClosed bool
	// 告知当前链接已经退出、停止的channel
	ExitChan chan bool
	// 该链接处理的方法Router
	Router ziface.IRouter
}

func (c *Connection) StartReader() {
	fmt.Println("StartReader Goroutine is running...")
	defer fmt.Println("ConnID", c.ConnID, " StartReader is exit, remote addr is :", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中，最大512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err:", err)
		//	continue
		//}

		// 创建一个拆包解包对象
		dp := NewDataPack()

		// 读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error:", err)
			break
		}

		// 拆包，得到msgID和msgDatalen放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Printf("unpack error:", err)
			break
		}

		// 根据dataLen，再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err = io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Printf("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn:    c,
			message: msg,
		}

		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			// 从路由中，找到注册绑定的Conn对应的router调用
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

// Start 启动链接，让当前链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID:", c.ConnID)

	// 启动从当前链接的读数据的业务
	go c.StartReader()
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

// SendMsg 将数据发送给远程的客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将data进行封包 MsgDataLen|MsgID|Data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Printf("Pack error msg id:", msgId)
		return errors.New("Pack error msg")
	}

	// 将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Printf("Write msg id:", msgId, " error:", err)
		return errors.New("conn write error")
	}

	return nil
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}
