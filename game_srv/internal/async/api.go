package async

import (
	"context"
	"errors"
	"kernel/kenum"
	"sync"
)

var (
	instance *AsncQueue
	once     sync.Once
)

func Init(panicFn func(string), pushBackFn PushReadUserFn) {
	createInstance(panicFn, pushBackFn)
}

func createInstance(panicFn func(string), pushBackFn PushReadUserFn) {
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		instance = &AsncQueue{
			ctx:        ctx,
			cancel:     cancel,
			stopChan:   make(chan struct{}),
			userQueue:  make(chan *AyncReadUser, 1024*10),
			panicFunc:  panicFn,
			pushBackFn: pushBackFn,
		}

		instance.state.Store(kenum.WorkState_Idle)
	})
}

func Start() error {
	if instance == nil {
		return errors.New("AsncQueue not initialized, call Init first")
	}
	return instance.start()
}

func Stop() error {
	if instance == nil {
		return errors.New("AsncQueue not initialized, call Init first")
	}
	return instance.stop()
}

func Push(r *AyncReadUser) error {
	if instance == nil {
		return errors.New("AsncQueue not initialized, call Init first")
	}

	if r == nil {
		return errors.New("AyncQueue Push nil")
	}

	if r.Cb == nil {
		return errors.New("AyncQueue Push cb nil")
	}

	return instance.push(r)
}
