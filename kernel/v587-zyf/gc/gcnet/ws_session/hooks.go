package ws_session

import (
	"sync"

	"github.com/v587-zyf/gc/iface"
)

type Hooks struct {
	mu sync.Mutex

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
	h.lock()

	if fn != nil {
		h.startMethods[key] = fn
	}
}

func (h *Hooks) OnRecv(key string, fn Recv) {
	h.lock()

	if fn != nil {
		h.recvMethods[key] = fn
	}
}

func (h *Hooks) OnStop(key string, fn Call) {
	h.lock()

	if fn != nil {
		h.stopMethods[key] = fn
	}
}

func (h *Hooks) OnMethod(method iface.IWsSessionMethod) {
	h.lock()

	if method != nil {
		h.OnStart(method.Name(), method.Start)
		h.OnRecv(method.Name(), method.Recv)
		h.OnStop(method.Name(), method.Stop)
	}
}

func (h *Hooks) OnHooks(hooks *Hooks) {
	h.lock()

	for key, fn := range hooks.startMethods {
		h.startMethods[key] = fn
	}

	for key, fn := range hooks.recvMethods {
		h.recvMethods[key] = fn
	}

	for key, fn := range hooks.stopMethods {
		h.stopMethods[key] = fn
	}
}

func (h *Hooks) ExecuteStart(ss iface.IWsSession) {
	h.lock()

	for _, v := range h.startMethods {
		v(ss)
	}
}

func (h *Hooks) ExecuteRecv(ss iface.IWsSession, data []byte) {
	h.lock()

	for _, v := range h.recvMethods {
		v(ss, data)
	}
}

func (h *Hooks) ExecuteStop(ss iface.IWsSession) {
	h.lock()

	for _, v := range h.stopMethods {
		v(ss)
	}
}

func (h *Hooks) lock() {
	// h.mu.Lock()
	// defer h.mu.Unlock()
}
