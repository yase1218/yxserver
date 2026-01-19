package loadbalancer

import (
	"container/heap"
	"sync"
	"time"
)

/*
最小堆维护服务信息 堆顶负载最小
上报端负载阈值设置 减少更新频率
*/

// ServeInfo 一个服务
type ServeInfo struct {
	ID       uint32
	Load     uint32
	Addr     string
	UpdateAt time.Time

	index int // 在堆中的索引，由堆内部维护
}

type MinHeap []*ServeInfo

func (h MinHeap) Len() int { return len(h) }
func (h MinHeap) Less(i, j int) bool {
	if h[i].Load == h[j].Load {
		return h[i].UpdateAt.Before(h[j].UpdateAt)
	}
	return h[i].Load < h[j].Load
}
func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *MinHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*ServeInfo)
	item.index = n
	*h = append(*h, item)
}
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	item.index = -1 // 标记为移除
	*h = old[0 : n-1]
	return item
}

// LoadBalancer 负载均衡器
type LoadBalancer struct {
	mu     sync.RWMutex
	severs map[uint32]*ServeInfo
	heap   MinHeap
}

func NewLoadBalancer() *LoadBalancer {
	lb := &LoadBalancer{
		severs: make(map[uint32]*ServeInfo),
		heap:   make(MinHeap, 0),
	}
	heap.Init(&lb.heap) // 初始化堆
	return lb
}

// Upsert 更新或添加服务信息
func (lb *LoadBalancer) Upsert(id uint32, addr string, load uint32) bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	cur := time.Now()
	if gw, exists := lb.severs[id]; exists {
		oldLoad := gw.Load
		gw.Load = load
		gw.Addr = addr // 地址也可能更新
		gw.UpdateAt = cur
		if gw.Load != oldLoad {
			heap.Fix(&lb.heap, gw.index)
		}
		return false
	} else {
		newGw := &ServeInfo{
			ID:       id,
			Load:     load,
			Addr:     addr,
			UpdateAt: cur,
		}
		heap.Push(&lb.heap, newGw)
		lb.severs[id] = newGw
		return true
	}
}

// SelectServer 获取负载最小的服务
func (lb *LoadBalancer) SelectServer() *ServeInfo {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if lb.heap.Len() == 0 {
		return nil
	}
	// 堆顶元素即为负载最小的服务
	minGw := lb.heap[0]
	// 返回一个副本，防止调用方修改内部数据
	return &ServeInfo{
		ID:   minGw.ID,
		Load: minGw.Load,
		Addr: minGw.Addr,
	}
}

// RemoveServer 移除指定服务
func (lb *LoadBalancer) RemoveServer(id uint32) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if gw, exists := lb.severs[id]; exists {
		heap.Remove(&lb.heap, gw.index)
		delete(lb.severs, id)
	}
}

func (lb *LoadBalancer) FindServer(id uint32) *ServeInfo {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return lb.severs[id]
}

func (lb *LoadBalancer) Len() int {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return len(lb.severs)
}
