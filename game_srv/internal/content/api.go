package content

import (
	"context"
	"errors"
	"kernel/kenum"
	"sync"
)

var (
	instance *ContentService
	once     sync.Once
	id       uint64
)

func Init(panicFn func(string), pushBackFn PushContentFn) {
	createInstance(panicFn, pushBackFn)
}

func createInstance(panicFn func(string), pushBackFn PushContentFn) {
	once.Do(func() {
		ctx, cacel := context.WithCancel(context.Background())
		instance = &ContentService{
			ctx:        ctx,
			cancel:     cacel,
			stopChan:   make(chan struct{}),
			queue:      make(chan *ContentData, 1024*10),
			panicFunc:  panicFn,
			pushBackFn: pushBackFn,
		}

		instance.state.Store(kenum.WorkState_Idle)
	})
}
func Start() error {
	if instance == nil {
		return errors.New("ContentService not initialized, call Init first")
	}

	return instance.start()
}

func Stop() error {
	if instance == nil {
		return errors.New("ContentService not initialized, call Init first")
	}
	return instance.stop()
}

func PushContent(data *ContentData) error {
	id++
	data.Id = id
	if instance == nil {
		return errors.New("ContentService not initialized, call Init first")
	}
	if data == nil {
		return errors.New("ContentService push nil")
	}

	if data.Cb == nil {
		return errors.New("ContentService cb nil")
	}

	return instance.push(data)
}
