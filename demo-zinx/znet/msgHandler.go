package znet

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"strconv"
)

// MsgHandle 消息处理模块的实现
type MsgHandle struct {
	// 存放每个msgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// NewMsgHandle 初始化创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1. 从request中找到MsgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgID(), " is Not Found! Need Register")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1. 判断当前 msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgID:" + strconv.Itoa(int(msgID)))
	}
	// 2. 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID:", msgID, " Success!")
}
