package tcp_session

import (
	"github.com/v587-zyf/gc/iface"
	"sync"
)

type Hooks struct {
	mu sync.Mutex

	onStartFns []Call
	onRecvFns  []Recv
	onStopFns  []Call
}

func NewHooks() *Hooks {
	return &Hooks{
		onStartFns: make([]Call, 0, 2),
		onRecvFns:  make([]Recv, 0, 2),
		onStopFns:  make([]Call, 0, 2),
	}
}

func (h *Hooks) OnStart(fns ...Call) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.onStartFns = append(h.onStartFns, fns...)
}

func (h *Hooks) OnRecv(fns ...Recv) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.onRecvFns = append(h.onRecvFns, fns...)
}

func (h *Hooks) OnStop(fns ...Call) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.onStopFns = append(h.onStopFns, fns...)
}

func (h *Hooks) OnMethod(method iface.ITcpSessionMethod) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.onStartFns = append(h.onStartFns, method.Start)
	h.onRecvFns = append(h.onRecvFns, method.Recv)
	h.onStopFns = append(h.onStopFns, method.Stop)
}

func (h *Hooks) OnHooks(hooks *Hooks) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.onStartFns = append(h.onStartFns, hooks.onStartFns...)
	h.onRecvFns = append(h.onRecvFns, hooks.onRecvFns...)
	h.onStopFns = append(h.onStopFns, hooks.onStopFns...)
}

func (h *Hooks) ExecuteStart(ss iface.ITcpSession) {
	for _, v := range h.onStartFns {
		v(ss)
	}
}

func (h *Hooks) ExecuteRecv(ss iface.ITcpSession, data []byte) {
	for _, v := range h.onRecvFns {
		v(ss, data)
	}
}

func (h *Hooks) ExecuteStop(ss iface.ITcpSession) {
	for _, v := range h.onStopFns {
		v(ss)
	}
}
