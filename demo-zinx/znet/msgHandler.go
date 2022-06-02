package znet

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/utils"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"strconv"
)

// MsgHandle 消息处理模块的实现
type MsgHandle struct {
	// 存放每个msgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池的worker数量
	WorkerPollSize uint32
}

// NewMsgHandle 初始化创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPollSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.MaxPackageSize),
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

// StartWorkerPool 启动一个Worker工作池（开启工作池的动作只能发生一次，一个zinx框架只能有一个worker工作池）
func (mh *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize分别开启Worker，每个worker用一个go来承载
	for i := 0; i < int(mh.WorkerPollSize); i++ {
		// 一个worker被启动
		// 1. 当前worker对应的channel消息队列 开辟空间 第0个worker 就用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2. 启动当前的worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID：", workerID, " is started...")

	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue, 由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1. 将消息平均分配给不同的worker, 根据客户端建立ConnID来进行分配
	workID := request.GetConnection().GetConnID() % mh.WorkerPollSize
	fmt.Println("Add ConnID:", request.GetConnection().GetConnID(),
		" request MsgID:", request.GetMsgID(),
		"to WorkerID:", workID)

	// 2. 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workID] <- request
}
