package tapping

import (
	"context"
	"errors"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"kernel/kenum"
	"sync"
)

var (
	instance *TapQueue
	once     sync.Once
)

func Init(tapFunc TapFunc, panicFunc func(string), args ...int) {
	createInstance(tapFunc, panicFunc, args...)
}

func createInstance(tapFunc TapFunc, panicFunc func(string), args ...int) {
	once.Do(func() {
		if len(args) != 1 {
			log.Panic("tap queue start param lenth err",
				zap.Int("need", 1), zap.Int("lenth", len(args)))
		}
		ctx, cancel := context.WithCancel(context.Background())
		instance = &TapQueue{
			ctx:       ctx,
			cancel:    cancel,
			stopChan:  make(chan struct{}),
			queue:     make(chan *TapData, 1024*args[0]),
			tapFunc:   tapFunc,
			panicFunc: panicFunc,
		}

		instance.state.Store(kenum.WorkState_Idle)
	})
}

func Start() error {
	if instance == nil {
		return errors.New("tap queue not initialized, call Init first")
	}

	return instance.start()
}

func Stop() error {
	if instance == nil {
		return errors.New("tap queue not initialized, call Init first")
	}
	return instance.stop()
}

func PushTapData(m *TapData) error {
	if instance == nil {
		return errors.New("tap queue not initialized, call Init first")
	}
	return instance.push(m)
}
