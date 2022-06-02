package utils

import (
	"encoding/json"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"io/ioutil"
)

/*
	存储一切有关Zinx框架的全局参数，供其他模块使用
	一些参数是可以通过zinx.json由用户进行配置的
*/
type GlobalObj struct {
	/*
		server
	*/
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         //当前服务器的名称
	/*
		zinx
	*/
	Version        string // 当前Zinx的版本号
	MaxConn        int    //当前服务器主机允许的最大链接数
	MaxPackageSize uint32 //当前zinx框架数据包的最大值
}

// GlobalObject 定义一个全局的对外Globalobj
var GlobalObject *GlobalObj

// Reload 从zinx.json去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	// 将json文件解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法，初始化当前的GlobalObject
func init() {
	GlobalObject := &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.6",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	GlobalObject.Reload()
}
