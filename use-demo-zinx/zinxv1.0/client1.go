package main

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/znet"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client1 start...")

	time.Sleep(1 * time.Second)
	// 1. 直接链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	for {
		// 发送封包的message消息 MsgID：0
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte("Zinx client1 Test Message")))
		if err != nil {
			fmt.Printf("Pack error:", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Printf("write error:", err)
			return
		}

		// 服务器应该给我们回复一个message数据， MsgID：1 Data：ping...ping...ping...
		// 先读取流中head部分，得到ID和dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Printf("read head error:", err)
			break
		}

		// 将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error:", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error:", err)
				return
			}

			fmt.Println("---->Recv Server Msg:ID=", msg.Id, " len=", msg.DataLen, " data=", string(msg.Data))
		}

		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
