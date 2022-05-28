package znet

import "github.com/bailuoxi66/go-inside/demo-zinx/ziface"

// BaseRouter 实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类方法重写就好了
type BaseRouter struct{}

/*
	这里之所以BaseRouter的方法都为空
	是因为有的Router不希望有PreHandler和PostHandler这两个业务
	所以Router全部继承BaseRouter的好处是，不需要实现PreHandler和PostHandler
*/
// PreHandle 在处理conn业务之前的钩子方法hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// Handle 处理conn业务的钩子方法hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// PostHandle 在处理conn业务之后的钩子方法hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
