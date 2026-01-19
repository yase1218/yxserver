package event

import (
	"fmt"
	"runtime"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

var (
	EventMgr eventMgr
)

type eventMgr struct {
	PoolSize  uint32
	EventSize uint32
	EventChan []chan IEvent
}

func (e *eventMgr) Init(poolSize, eventSize uint32) {
	e.PoolSize = poolSize
	e.EventSize = eventSize
}

func (e *eventMgr) Run() {
	//log.Info("eventMgr Run ")
	if e.PoolSize == 0 {
		e.PoolSize = 1
	}

	// for i := 0; i < int(e.PoolSize); i++ {
	// 	e.EventChan = append(e.EventChan, make(chan IEvent, e.EventSize))
	// 	go tools.GoSafe("do work", func() {
	// 		eventWorker(i, e.EventChan[i])
	// 	}, service.PostPanic)
	// }
}

func (e *eventMgr) Close() {
	for i := 0; i < len(e.EventChan); i++ {
		close(e.EventChan[i])
	}
}

func (e *eventMgr) PublishEvent(event IEvent) {
	// pos := event.RouteID() % int64(e.PoolSize)
	// e.EventChan[pos] <- event
	event.CallBack()(event)
}

func eventWorker(worker int, c chan IEvent) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			err := fmt.Errorf("%v: %s", r, buf[:l])
			log.Error("panic", zap.Error(err))

			// go tools.GoSafe("do work", func() {
			// 	eventWorker(worker, c)
			// })
		}
	}()
	//log.Debug("EventWorker worker:%v start", worker)
	//defer func() {
	//	if err := recover(); err != nil {
	//		log.Error("EventWorker err", zap.Int("worker", worker), zap.Any("err", err))
	//		//log.Errorf("EventWorker worker:%v err:%v", worker, err)
	//		go eventWorker(worker, c)
	//	}
	//	//log.Infof("EventWorker worker %v exit", worker)
	//}()
	for event := range c {
		if f := event.CallBack(); f != nil {
			f(event)
		}
	}
}
