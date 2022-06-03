package ziface

// IConnManager 链接管理模块抽象层
type IConnManager interface {
	// Add 添加链接
	Add(conn IConnection)
	// Remove 删除链接
	Remove(conn IConnection)
	// Get 根据ConnID获取链接
	Get(connID uint32) (IConnection, error)
	// Len 得到当前链接总数
	Len() int
	// ClearConn 清楚并种植所有的链接
	ClearConn()
}
