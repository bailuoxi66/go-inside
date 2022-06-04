package znet

import (
	"errors"
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/utils"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"io"
	"net"
)

type Connection struct {
	// 当前Conn隶属于那个Server
	TcpServer ziface.IServer
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn
	// 链接ID
	ConnID uint32
	// 当前的链接状态
	isClosed bool
	// 告知当前链接已经退出、停止的channel(由Reader告知Writer退出)
	ExitChan chan bool
	// 无缓冲管道，用于读写 goroutine之间的消息通信
	msgChan chan []byte
	// 消息的管理MsgID 和对应业务的处理业务API关系
	MsgHandler ziface.IMsgHandler
}

func (c *Connection) StartReader() {
	fmt.Println("[StartReader Goroutine is running...]")
	defer fmt.Println("ConnID", c.ConnID, "[StartReader is exit], remote addr is :", c.Conn.RemoteAddr().String())
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中，找到注册绑定的Conn对应的router调用
			// 根据绑定好的MsgID 找到对应处理api业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// StartWriter 写消息Goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[StartWriter Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!!!]")

	// 不断阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

// Start 启动链接，让当前链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID:", c.ConnID)

	// 启动从当前链接的读数据的业务
	go c.StartReader()
	// 启动从写前链接的读数据的业务
	go c.StartWriter()
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

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前链接从ConnMgr中拆除掉
	c.TcpServer.GetConnMgr().Remove(c)
	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

	// 将数据发送给chan
	c.msgChan <- binaryMsg
	return nil
}

// NewConnection 初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}

	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}
