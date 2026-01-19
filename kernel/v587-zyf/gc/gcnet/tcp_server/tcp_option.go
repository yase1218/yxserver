package tcp_server

import (
	"github.com/v587-zyf/gc/iface"
)

type TcpOption struct {
	listenAddr string

	method iface.ITcpSessionMethod
}

type Option func(opts *TcpOption)

func NewTcpOption() *TcpOption {
	o := &TcpOption{}

	return o
}

func WithListenAddr(addr string) Option {
	return func(opts *TcpOption) {
		opts.listenAddr = addr
	}
}

func WithMethod(m iface.ITcpSessionMethod) Option {
	return func(opts *TcpOption) {
		opts.method = m
	}
}
