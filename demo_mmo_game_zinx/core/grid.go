package core

import (
	"fmt"
	"sync"
)

// Grid 一个AOI地图中的格子类型
type Grid struct {
	//格子ID
	GID int
	//格子的左边边界坐标
	MinX int
	//格子的右边边界坐标
	MaxX int
	//格子的左边边界坐标
	MinY int
	//格子的右边边界坐标
	MaxY int
	//当前格子内玩家或者物体成员的ID集合
	playerIDs map[int]bool
	//保护当前集合的锁
	pIDLock sync.RWMutex
}

// NewGrid 初始化当前的格子的方法
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// Add 格子中添加一个玩家
func (g *Grid) Add(playerID int) {
	fmt.Println(playerID)
	fmt.Println(g)
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
	fmt.Println("add gID:", g.GID)
	fmt.Println("add playerID:", playerID)
	fmt.Println("add Grid:", g)
}

// Remove 格子中移除一个玩家
func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// GetPlayerIDs 得到当前格子中所有玩家ID
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}

// 调试使用-打印出格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id:%d MinX:%d MaxX:%d MinY:%d MaxY:%d playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
