package tcp_handler

import (
	"github.com/v587-zyf/gc/iface"
)

type TcpHandlerOption struct {
	name    string
	startFn func(s iface.ITcpSession)
	recvFn  func(s iface.ITcpSession, data any)
	stopFn  func(s iface.ITcpSession)
}

type Option func(opts *TcpHandlerOption)

func NewTcpHandlerOption() *TcpHandlerOption {
	o := &TcpHandlerOption{}

	return o
}

func WithName(name string) Option {
	return func(opts *TcpHandlerOption) {
		opts.name = name
	}
}

func WithStartFn(fn func(s iface.ITcpSession)) Option {
	return func(opts *TcpHandlerOption) {
		opts.startFn = fn
	}
}

func WithRecvFn(fn func(s iface.ITcpSession, data any)) Option {
	return func(opts *TcpHandlerOption) {
		opts.recvFn = fn
	}
}

func WithStopFn(fn func(s iface.ITcpSession)) Option {
	return func(opts *TcpHandlerOption) {
		opts.stopFn = fn
	}
}
