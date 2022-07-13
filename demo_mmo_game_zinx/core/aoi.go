package core

import "fmt"

// 定义一些AOI的边界值
const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

// AOIManager AOI区域管理模块
type AOIManager struct {
	// 区域的左边界坐标
	MinX int
	// 区域的右边界坐标
	MaxX int
	// X方向格子的数量
	CntsX int
	// 区域的上边界坐标
	MinY int
	// 区域的下边界坐标
	MaxY int
	// Y方向格子的数量
	CntsY int
	// 当前区域中有哪些格子map key；格子ID value：格子对象
	grids map[int]*Grid
}

// NewAOIManager 初始化一个AOI区域管理模块
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	//给AOI初始化区域的格子所有的格子进行编号 和 初始化
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//计算格子id，根据x，y编号
			gid := y*cntsX + x

			//初始化gid格子
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
		}
	}
	return aoiMgr
}

//得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

//得到每个格子在Y轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

//打印格子信息
func (m *AOIManager) String() string {
	//打印AOIManager信息
	s := fmt.Sprintf("AOIManager:\nMinX:%d, MaxX:%d, cntsX:%d, MinY:%d, MaxY:%d, cntsY:%d\ngrids:",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)

	//打印全部格子信息
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

// GetSurroundGridsByGid 根据GID得到周边九宫格格子集合
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	// 判断gID是否在AOIManager中
	if _, ok := m.grids[gID]; !ok {
		return
	}

	// 将当前gid本身加入到九宫格切片中
	grids = append(grids, m.grids[gID])

	// 判断当前gID左边或者右边是否还有格子？
	idx := gID % m.CntsX
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1])
	}

	// 将x轴当前的格子都取出，进行遍历
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	for _, v := range gidsX {
		// 得到当前格子id的y轴编号
		idy := v / m.CntsX
		// 是否上面还有格子
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])
		}
		// 是否下面还有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}
	}

	return grids
}

// GetGidByPos 通过x、y横纵坐标得到当前GID格子编号
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridLength()
	idy := (int(y) - m.MinY) / m.gridWidth()

	return idy*m.CntsX + idx
}

// GetPidsByPos x、y横纵坐标得到周边九宫格全部的PlayerIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	//得到当前玩家GID格子id
	gID := m.GetGidByPos(x, y)
	//通过gID得到九宫格信息
	grids := m.GetSurroundGridsByGid(gID)
	//将九宫格信息放置到playerIDs
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GID)
	}

	return playerIDs
}

//添加一个playerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

//移除一个格子中的指定playerID
func (m *AOIManager) RemovePidToGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

//通过gID获取全部的PlayerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

//通过坐标将player添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]

	grid.Add(pID)
}

//通过坐标将player从一个格子中移除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]

	grid.Remove(pID)
}
