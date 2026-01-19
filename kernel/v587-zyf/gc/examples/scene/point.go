package scene

import (
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"sync"
)

type Point struct {
	id     int
	minX   int
	maxX   int
	minY   int
	maxY   int
	isSafe bool

	objs  map[iface.TileType]struct{}
	oLock sync.RWMutex
}

func NewGrid(id int, minX, maxX, minY, maxY int) *Point {
	p := &Point{
		id:   id,
		minX: minX,
		maxX: maxX,
		minY: minY,
		maxY: maxY,
		objs: make(map[iface.TileType]struct{}),
	}

	return p
}

func (p *Point) IsSafe() bool   { return p.isSafe }
func (p *Point) SetSafe(i bool) { p.isSafe = i }
func (p *Point) ID() int        { return p.id }
func (p *Point) AddObj(t iface.TileType) {
	p.oLock.Lock()
	defer p.oLock.Unlock()

	p.objs[t] = struct{}{}
}
func (p *Point) RemoveObj(t iface.TileType) {
	p.oLock.Lock()
	defer p.oLock.Unlock()

	delete(p.objs, t)
}
func (p *Point) GetObjs() map[iface.TileType]struct{} {
	p.oLock.RLock()
	defer p.oLock.RUnlock()

	return p.objs
}
func (p *Point) GetObjsSlice() []iface.TileType {
	p.oLock.RLock()
	defer p.oLock.RUnlock()

	slice := make([]iface.TileType, len(p.objs))

	i := 0
	for tileType := range p.objs {
		slice[i] = tileType
	}

	return slice
}

func (p *Point) String() string {
	return fmt.Sprintf("id:%d minX:%d maxX:%d minY:%d maxY:%d",
		p.id, p.minX, p.maxX, p.minY, p.maxY)
}
