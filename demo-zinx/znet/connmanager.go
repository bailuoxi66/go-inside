package znet

import (
	"errors"
	"fmt"
	"github.com/bailuoxi66/go-inside/demo-zinx/ziface"
	"sync"
)

// ConnManager 链接管理模块
type ConnManager struct {
	// 管理的链接集合
	connections map[uint32]ziface.IConnection
	// 保护链接集合的读写锁
	connLock sync.RWMutex
}

// NewConnManager 创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源 添加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn加入到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID=", conn.GetConnID(), " connection add to ConnManager successFully:conn num = ", connMgr.Len())
}

// Remove 删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源 添加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn加入到ConnManager中
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID=", conn.GetConnID(), " remove from ConnManager successFully:conn num = ", connMgr.Len())
}

// Get 根据ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源 添加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		// 找到了
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// Len 得到当前链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// ClearConn 清楚并种植所有的链接
func (connMgr *ConnManager) ClearConn() {
	// 保护共享资源 添加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除conn，并停止conn工作
	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connections succ! conn num = ", len(connMgr.connections))
}
