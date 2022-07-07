package core

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
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
