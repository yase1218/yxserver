package module

import (
	"context"
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"kernel/tools"
	"sync"
)

type ModuleMgr struct {
	modules sync.Map // name:iface.IModule
}

func (mm *ModuleMgr) Add(m iface.IModule) {
	if _, ok := mm.modules.Load(m.Name()); ok {
		log.Warn("module already exists", zap.String("name", m.Name()))
		return
	}

	mm.modules.Store(m.Name(), m)
}

func (mm *ModuleMgr) Get(name string) iface.IModule {
	if module, ok := mm.modules.Load(name); ok {
		return module.(iface.IModule)
	}

	log.Warn("module not exists", zap.String("name", name))
	return nil
}

func (mm *ModuleMgr) Del(name string) {
	if _, ok := mm.modules.Load(name); ok {
		mm.modules.Delete(name)
		log.Info("module deleted", zap.String("name", name))
	} else {
		log.Warn("module not found for deletion", zap.String("name", name))
	}
}

func (mm *ModuleMgr) Length() int {
	length := 0
	mm.modules.Range(func(key, value any) bool {
		length++
		return true
	})

	return length
}

func (mm *ModuleMgr) Init(ctx context.Context, opts ...iface.Option) (err error) {
	moduleLen := mm.Length()
	if moduleLen <= 0 {
		return fmt.Errorf("no module")
	}

	var wg sync.WaitGroup
	wg.Add(moduleLen)

	mm.modules.Range(func(key, value any) bool {
		go tools.GoSafe("module init", func() {
			defer wg.Done()
			m := value.(iface.IModule)
			err := m.Init(ctx, opts...)
			if err != nil {
				log.Error("module init failed", zap.String("name", m.Name()), zap.Error(err))
			}
		})
		return true
	})

	wg.Wait()

	return nil
}

func (mm *ModuleMgr) Start() (err error) {
	moduleLen := mm.Length()
	if moduleLen <= 0 {
		return fmt.Errorf("no module")
	}

	var wg sync.WaitGroup
	wg.Add(moduleLen)

	mm.modules.Range(func(key, value any) bool {
		go tools.GoSafe("module init", func() {
			defer wg.Done()
			m := value.(iface.IModule)
			err := m.Start()
			if err != nil {
				log.Error("module init failed", zap.String("name", m.Name()), zap.Error(err))
			}
		})

		return true
	})

	wg.Wait()

	return nil
}

func (mm *ModuleMgr) Run() {
	moduleLen := mm.Length()
	if moduleLen <= 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(moduleLen)

	mm.modules.Range(func(key, value any) bool {
		go tools.GoSafe("module init", func() {
			defer wg.Done()
			m := value.(iface.IModule)
			m.Run()
		})

		return true
	})

	wg.Wait()
}

func (mm *ModuleMgr) Stop() {
	moduleLen := mm.Length()
	if moduleLen <= 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(moduleLen)

	mm.modules.Range(func(key, value any) bool {
		go tools.GoSafe("module init", func() {
			defer wg.Done()
			m := value.(iface.IModule)
			m.Stop()
		})

		return true
	})

	wg.Wait()
}
