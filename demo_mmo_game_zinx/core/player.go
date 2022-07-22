package core

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"github.com/bailuoxi66/go-inside/demo_mmo_game_zinx/pb"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
)

// Player 玩家对象
type Player struct {
	Pid  int32              // 玩家ID
	Conn ziface.IConnection // 当前玩家的链接（用于和客户端的链接）
	X    float32            // 平面X坐标
	Y    float32            // 高度
	Z    float32            // 平面Y坐标（注意 不是Y）
	V    float32            // 旋转0-360 角度
}

var PidGen int32 = 1  // 用来生成玩家ID的生成器
var IdLock sync.Mutex // 保护PidGen的锁

// NewPlayer 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	// 创建一个玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}

	return p
}

//SendMsg 提供一个发送个客户端消息的方法,主要是将pb的protobuf数据序列化之后，再调用zinx的SendMsg方法
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	// 将proto Message结构体序列化之后 转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return
	}

	// 将二进制文件 通过zinx框架的SendMsg将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg error")
		return
	}

	return
}

// SyncPid 告知客户端玩家Pid，同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	// 组件MsgID：1 的proto数据
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	// 将消息发送给客户端
	p.SendMsg(1, proto_msg)
}

// BroadCastStartPosition 告知客户端玩家Pid，同步已经生成的玩家ID给客户端
func (p *Player) BroadCastStartPosition() {
	// 组件MsgID：200 的proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 将消息发送给客户端
	p.SendMsg(200, proto_msg)
}

// Talk 玩家广播世界聊天消息
func (p *Player) Talk(content string) {
	// 1. 组件MsgID 200 proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, // tp-1 代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	// 2. 得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	// 3. 向所有的玩家（包括自己）发送MsgID 200的消息
	for _, player := range players {
		// player分别给对应的客户端发送Msg
		player.SendMsg(200, proto_msg)
	}
}

// SyncSurrounding 同步玩家上线的位置消息
func (p *Player) SyncSurrounding() {
	// 1. 获取当前玩家周围有哪些（九宫格）
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	fmt.Println("444:", pids)
	fmt.Println("4444:", WorldMgrObj.AoiMgr.GetGidByPos(p.X, p.Z))

	gId := WorldMgrObj.AoiMgr.GetGidByPos(p.X, p.Z)

	pp := WorldMgrObj.AoiMgr.GetPidsByGid(gId)
	fmt.Println("pp:", pp)

	players := make([]*Player, 0, len(pids))

	for _, pid := range pids {
		if P := WorldMgrObj.GetPlayerByPid(int32(pid)); P != nil {
			players = append(players, P)
		}
	}
	fmt.Println("555:", players)

	// 2. 将当前玩家的位置信息通过MsgID:200 发送给周围的玩家（让其他玩家看到自己）
	// 2.1 组件MsgID：200 proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //Tp2 代表广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 2.2 全部周围的玩家都向格子的客户端发送200消息，proto_msg
	for _, player := range players {
		fmt.Println("..............")
		fmt.Println(player)
		player.SendMsg(200, proto_msg)
	}

	// 3 将周围的全部玩家的位置信息发送给当前的玩家MsgID：202 客户端（让自己看到其他玩家）
	// 3.1 组建MsgID：202 proto数据
	// 3.1.1 制作pb.Player slice
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		// 制作一个message Player
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}

		players_proto_msg = append(players_proto_msg, p)
	}

	// 3.1.2 封装SyncPlayer protobuf数据
	SyncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}

	// 3.2 将组件好的数据发送给当前玩家的客户端
	p.SendMsg(202, SyncPlayers_proto_msg)
}
