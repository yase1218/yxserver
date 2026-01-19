package ws_conn

import (
	"runtime/debug"
	"sync"

	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

type Recv func(conn iface.IWsSession, data []byte)

type Call func(ss iface.IWsSession)

type Hooks struct {
	mu sync.RWMutex

	startMethods map[string]Call
	recvMethods  map[string]Recv
	stopMethods  map[string]Call
}

func NewHooks() *Hooks {
	return &Hooks{
		startMethods: make(map[string]Call),
		recvMethods:  make(map[string]Recv),
		stopMethods:  make(map[string]Call),
	}
}

func (h *Hooks) OnStart(key string, fn Call) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if key == "" || fn == nil {
		return
	}
	h.startMethods[key] = fn
}

func (h *Hooks) OnRecv(key string, fn Recv) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if key == "" || fn == nil {
		return
	}
	h.recvMethods[key] = fn
}

func (h *Hooks) OnStop(key string, fn Call) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if key == "" || fn == nil {
		return
	}
	h.stopMethods[key] = fn
}

func (h *Hooks) OnMethod(method iface.IWsSessionMethod) {
	if method == nil {
		return
	}

	h.OnStart(method.Name(), method.Start)
	h.OnRecv(method.Name(), method.Recv)
	h.OnStop(method.Name(), method.Stop)
}

func (h *Hooks) OnHooks(hooks *Hooks) {
	if hooks == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for key, fn := range hooks.startMethods {
		if key != "" && fn != nil {
			h.startMethods[key] = fn
		}
	}

	for key, fn := range hooks.recvMethods {
		if key != "" && fn != nil {
			h.recvMethods[key] = fn
		}
	}

	for key, fn := range hooks.stopMethods {
		if key != "" && fn != nil {
			h.stopMethods[key] = fn
		}
	}
}

func (h *Hooks) ExecuteStart(ss iface.IWsSession) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, v := range h.startMethods {
		go func(fn Call) {
			defer func() {
				if r := recover(); r != nil {
					log.Error("ExecuteStart panic", zap.Any("r", r), zap.String("stack", string(debug.Stack())))
				}
			}()
			fn(ss)
		}(v)
	}
}

func (h *Hooks) ExecuteRecv(ss iface.IWsSession, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, v := range h.recvMethods {
		go func(fn Recv) {
			defer func() {
				if r := recover(); r != nil {
					log.Error("ExecuteRecv panic", zap.Any("r", r), zap.String("stack", string(debug.Stack())))
				}
			}()
			fn(ss, data)
		}(v)
	}
}

func (h *Hooks) ExecuteStop(ss iface.IWsSession) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, v := range h.stopMethods {
		go func(fn Call) {
			defer func() {
				if r := recover(); r != nil {
					log.Error("ExecuteStop panic", zap.Any("r", r), zap.String("stack", string(debug.Stack())))
				}
			}()
			fn(ss)
		}(v)
	}
}
