package ziface

/*
	封包、拆包、模块
	直接面向TCP链接中的数据流，用于处理TCP粘包问题
*/
type IDataPack interface {
	// GetHeadLen 获取包的头的长度方法
	GetHeadLen() uint32
	// Pack 封包方法
	Pack(msg IMessage) ([]byte, error)
	// Unpack 拆包方法
	Unpack([]byte) (IMessage, error)
}
