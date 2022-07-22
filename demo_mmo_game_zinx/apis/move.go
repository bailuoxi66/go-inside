package apis

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"github.com/bailuoxi66/go-inside/demo-zinx/znet"
	"github.com/bailuoxi66/go-inside/demo_mmo_game_zinx/pb"
	"github.com/golang/protobuf/proto"
)

// MoveApi 玩家移动
type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	// 解析客户端传递进来的proto协议
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move: Position Unmarshal error:", err)
		return
	}

	// 得到当前发送位置的是那个玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error, ", err)
		return
	}

	fmt.Println("Player pid=%d, move(%f, %f, %f, %f)\n", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	// 给其他玩家进行当前玩家的位置信息广播
}
