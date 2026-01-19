package ws_handler

import (
	"github.com/v587-zyf/gc/iface"
)

type WsHandlerOption struct {
	name    string
	startFn func(s iface.IWsSession)
	recvFn  func(s iface.IWsSession, data any)
	stopFn  func(s iface.IWsSession)
}

type Option func(opts *WsHandlerOption)

func NewWsHandlerOption() *WsHandlerOption {
	o := &WsHandlerOption{}

	return o
}

func WithName(name string) Option {
	return func(opts *WsHandlerOption) {
		opts.name = name
	}
}

func WithStartFn(fn func(s iface.IWsSession)) Option {
	return func(opts *WsHandlerOption) {
		opts.startFn = fn
	}
}

func WithRecvFn(fn func(s iface.IWsSession, data any)) Option {
	return func(opts *WsHandlerOption) {
		opts.recvFn = fn
	}
}

func WithStopFn(fn func(s iface.IWsSession)) Option {
	return func(opts *WsHandlerOption) {
		opts.stopFn = fn
	}
}
