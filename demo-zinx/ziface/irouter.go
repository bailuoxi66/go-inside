package ziface

/*
	路由抽象接口
	路由里的数据都是IRequest
*/
type IRouter interface {
	// PreHandle 在处理conn业务之前的钩子方法hook
	PreHandle(request IRequest)
	// Handle 处理conn业务的钩子方法hook
	Handle(request IRequest)
	// PostHandle 在处理conn业务之后的钩子方法hook
	PostHandle(request IRequest)
}
