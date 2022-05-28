package ziface

/*
	IRequest接口：
	实际上是把客户端请求的链接信息和请求数据 包装到了一个Request中
*/
type IRequest interface {

	// GetConnection 得到当前链接
	GetConnection() IConnection
	// GetData 得到请求的消息数据
	GetData() []byte
}
