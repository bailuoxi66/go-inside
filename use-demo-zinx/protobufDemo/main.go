package main

import (
	"fmt"
	"github.com/bailuoxi66/go-inside/use-demo-zinx/protobufDemo/pb"
	"github.com/golang/protobuf/proto"
)

func main() {
	// 定义一个Person结构对象
	person := &pb.Person{
		Name:   "luoyu",
		Age:    17,
		Emails: []string{"bailuoxihaha@163.com", "bailuoxi@gmail.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "11111111",
				Type:   pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "22222222",
				Type:   pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "33333333",
				Type:   pb.PhoneType_WORK,
			},
		},
	}

	// 将person对象进行序列化，得到二进制信息
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal err:", err)
	}

	// 解码
	newdata := &pb.Person{}
	err = proto.Unmarshal(data, newdata)
	if err != nil {
		fmt.Println("unmarshal err:", err)
	}
	fmt.Println("源数据：", person)
	fmt.Println("目标数据：", newdata)
}
