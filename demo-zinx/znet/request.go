package znet

import "github.com/bailuoxi66/go-inside/demo-zinx/ziface"

type Request struct {
	// 已经和客户端建立好的链接
	conn ziface.IConnection
	// 客户端请求的数据
	message ziface.IMessage
}

// GetConnection 得到当前链接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 得到请求的数据消息
func (r *Request) GetData() []byte {
	return r.message.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.message.GetMsgId()
}
