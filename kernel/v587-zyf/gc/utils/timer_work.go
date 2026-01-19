package utils

import (
	"sync"
	"time"
)

// _Work 一个循环的计时工作
type _Work struct {
	Timer *time.Timer
	Func  func()
	Time  time.Duration

	//关闭通道
	CloseChan chan struct{}
}

// Stop 停止工作
func (t *_Work) Stop() {
	t.CloseChan <- struct{}{}
}

// Reset 重新计时
func (t *_Work) Reset(ti time.Duration) {
	if t.Timer == nil {
		return
	}

	t.Time = ti
	t.Timer.Reset(t.Time)
}

// Work 开始工作
func (t *_Work) Work() {
	defer func() {
		t.Timer.Stop()
	}()

	t.Timer = time.NewTimer(t.Time)
	for {
		select {
		case <-t.Timer.C:
			t.Func()
			t.Timer.Reset(t.Time)
			break
		case <-t.CloseChan:
			return
		}
	}
}

// TimeWork 计时工作
type TimeWork struct {
	mWorks   map[int32]*_Work
	mWorksCS sync.Mutex
}

// NewTimeWork 新建一个计时工作
func NewTimeWork() *TimeWork {
	return &TimeWork{
		mWorks: make(map[int32]*_Work),
	}
}

// Start 开启工作
func (t *TimeWork) Start(id int32, fun func(), ti time.Duration) {
	t.mWorksCS.Lock()
	defer t.mWorksCS.Unlock()

	work, ok := t.mWorks[id]
	if ok {
		work.Stop()
		delete(t.mWorks, id)
	}

	work = &_Work{
		Func:      fun,
		Time:      ti,
		CloseChan: make(chan struct{}, 1),
	}
	go work.Work()

	t.mWorks[id] = work
}

// Stop 关闭某个id的计时器任务
func (t *TimeWork) Stop(id int32) {
	t.mWorksCS.Lock()
	defer t.mWorksCS.Unlock()

	work, ok := t.mWorks[id]
	if !ok {
		// 未找到
		return
	}

	work.Stop()
	delete(t.mWorks, id)
}

// Close 关闭所有
func (t *TimeWork) Close() {
	t.mWorksCS.Lock()
	defer t.mWorksCS.Unlock()

	for _, work := range t.mWorks {
		work.Stop()
	}

	t.mWorks = make(map[int32]*_Work)
}

// Reset 重置时间间隔
func (t *TimeWork) Reset(id int32, ti time.Duration) {
	t.mWorksCS.Lock()
	defer t.mWorksCS.Unlock()

	work, ok := t.mWorks[id]
	if !ok {
		// 未找到
		return
	}

	work.Reset(ti)
}
