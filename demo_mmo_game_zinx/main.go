package main

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"github.com/bailuoxi66/go-inside/demo-zinx/znet"
	"github.com/bailuoxi66/go-inside/demo_mmo_game_zinx/apis"
	"github.com/bailuoxi66/go-inside/demo_mmo_game_zinx/core"
)

// OnConnectionAdd 当前客户端建立之后的Hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个Player对象
	player := core.NewPlayer(conn)

	fmt.Println("11:", player)
	// 给客户端发送MsgID:1的消息，同步当前Player的ID给客户端
	player.SyncPid()
	// 给客户端发送MsgID:200的消息，同步当前Player的初始位置给客户端
	player.BroadCastStartPosition()

	// 将当前新上线的玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	s := core.WorldMgrObj.GetAllPlayers()
	fmt.Println("333:", s)

	// 同步玩家上线的位置消息
	player.SyncSurrounding()

	// 将该链接绑定一个Pid，玩家ID的属性
	conn.SetProperty("pid", player.Pid)

	fmt.Println("=========> player pid = ", player.Pid, " is arrived<========")
}

// 给当前链接端口之前触发的Hook钩子函数
func OnConnectionLost(conn ziface.IConnection) {
	// 通过链接属性得到当前链接所绑定pid
	pid, _ := conn.GetProperty("pid")
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 触发玩家下线业务
	player.Offline()

	fmt.Println("=========>Player pid=", pid, " offline...<===========")
}

func main() {
	// 创建zinx server句柄
	s := znet.NewServer("MMO Game Zinx")

	// 连接创建和销毁的HOOK钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	// 注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	// 启动服务
	s.Server()
}
