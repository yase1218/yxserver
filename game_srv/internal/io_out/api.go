package io_out

import (
	"context"
	"errors"
	"kernel/kenum"
	"kernel/metric"
	"sync"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

var (
	instance *OutQueue
	once     sync.Once
)

func Init(outFunc OutFunc, panicFunc func(string), args ...int) {
	createInstance(outFunc, panicFunc, args...)
}

func createInstance(outFunc OutFunc, panicFunc func(string), args ...int) {
	once.Do(func() {
		if len(args) != 1 {
			log.Panic("out queue start param lenth err",
				zap.Int("need", 1), zap.Int("lenth", len(args)))
		}
		ctx, cancel := context.WithCancel(context.Background())
		instance = &OutQueue{
			ctx:       ctx,
			cancel:    cancel,
			stopChan:  make(chan struct{}),
			metrics:   &metric.ProcessorMetrics{},
			queue:     make(chan *OutMsg, 1024*args[0]),
			outFunc:   outFunc,
			panicFunc: panicFunc,
		}

		instance.state.Store(kenum.WorkState_Idle)
	})
}

func Start() error {
	if instance == nil {
		return errors.New("out queue not initialized, call Init first")
	}

	return instance.start()
}

func Stop() error {
	if instance == nil {
		return errors.New("out queue not initialized, call Init first")
	}
	return instance.stop()
}

func Push(m *OutMsg) error {
	if instance == nil {
		return errors.New("out queue not initialized, call Init first")
	}
	return instance.push(m)
}
