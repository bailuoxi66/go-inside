package ziface

// IServer 定义一个服务器接口
type IServer interface {

	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Server 运行服务器
	Server()

	// AddRouter 路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(router IRouter)
}
