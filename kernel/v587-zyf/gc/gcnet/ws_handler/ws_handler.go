package ws_handler

import (
	"context"

	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/iface"
)

type WsHandlerUnit struct {
	msgID   uint32
	handler ws_session.Recv
}
type WsHandler struct {
	options *WsHandlerOption

	ctx    context.Context
	cancel context.CancelFunc

	handlers map[uint32]*WsHandlerUnit
}

func NewWsHandler() *WsHandler {
	return &WsHandler{
		options:  NewWsHandlerOption(),
		handlers: make(map[uint32]*WsHandlerUnit),
	}
}

func (h *WsHandler) Init(ctx context.Context, option ...any) (err error) {
	h.ctx, h.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt.(Option)(h.options)
	}

	return nil
}

func (h *WsHandler) Name() string {
	name := ""
	if h.options.name != "" {
		name = h.options.name
	}
	return name
}

func (h *WsHandler) Register(msgID uint32, handler ws_session.Recv) {
	h.handlers[msgID] = &WsHandlerUnit{
		msgID:   msgID,
		handler: handler,
	}
}

func (h *WsHandler) GetHandler(msgID uint32) ws_session.Recv {
	if h, ok := h.handlers[msgID]; ok {
		return h.handler
	}
	return nil
}

func (h *WsHandler) HasHandler(msgID uint32) bool {
	_, ok := h.handlers[msgID]
	return ok
}

func (h *WsHandler) Start(ss iface.IWsSession) {
	if h.options.startFn != nil {
		h.options.startFn(ss)
	}
}

func (h *WsHandler) Recv(ss iface.IWsSession, data any) {
	if h.options.recvFn != nil {
		h.options.recvFn(ss, data)
	}
}

func (h *WsHandler) Stop(ss iface.IWsSession) {
	if h.options.stopFn != nil {
		h.options.stopFn(ss)
	}
}
