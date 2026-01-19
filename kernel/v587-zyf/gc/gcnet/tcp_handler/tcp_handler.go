package tcp_handler

import (
	"context"
	"github.com/v587-zyf/gc/gcnet/tcp_session"
	"github.com/v587-zyf/gc/iface"
)

type TcpHandlerUnit struct {
	msgID   uint32
	handler tcp_session.Recv
}
type TcpHandler struct {
	options *TcpHandlerOption

	ctx    context.Context
	cancel context.CancelFunc

	handlers map[uint32]*TcpHandlerUnit
}

func NewTcpHandler() *TcpHandler {
	return &TcpHandler{
		options:  NewTcpHandlerOption(),
		handlers: make(map[uint32]*TcpHandlerUnit),
	}
}

func (h *TcpHandler) Init(ctx context.Context, option ...Option) (err error) {
	h.ctx, h.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt(h.options)
	}

	return nil
}

func (h *TcpHandler) Name() string {
	name := ""
	if h.options.name != "" {
		name = h.options.name
	}
	return name
}

func (h *TcpHandler) Register(msgID uint32, handler tcp_session.Recv) {
	h.handlers[msgID] = &TcpHandlerUnit{
		msgID:   msgID,
		handler: handler,
	}
}

func (h *TcpHandler) GetHandler(msgID uint32) tcp_session.Recv {
	if h, ok := h.handlers[msgID]; ok {
		return h.handler
	}
	return nil
}

func (h *TcpHandler) HasHandler(msgID uint32) bool {
	_, ok := h.handlers[msgID]
	return ok
}

func (h *TcpHandler) Start(ss iface.ITcpSession) {
	if h.options.startFn != nil {
		h.options.startFn(ss)
	}
}

func (h *TcpHandler) Recv(ss iface.ITcpSession, data any) {
	if h.options.recvFn != nil {
		h.options.recvFn(ss, data)
	}
}

func (h *TcpHandler) Stop(ss iface.ITcpSession) {
	if h.options.stopFn != nil {
		h.options.stopFn(ss)
	}
}
