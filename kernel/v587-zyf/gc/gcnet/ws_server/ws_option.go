package ws_server

import (
	"github.com/v587-zyf/gc/iface"
	"net/http"
)

type HandlerFunc struct {
	path    string
	fn      func(http.ResponseWriter, *http.Request)
	methods string
}

type WsOption struct {
	addr string
	pem  string
	key  string

	https bool

	handler      http.Handler
	handlerFuncs []HandlerFunc
	method       iface.IWsSessionMethod
}

type Option func(opts *WsOption)

func NewWsOption() *WsOption {
	o := &WsOption{
		handlerFuncs: make([]HandlerFunc, 0),
	}

	return o
}

func WithAddr(addr string) Option {
	return func(opts *WsOption) {
		opts.addr = addr
	}
}

func WithPem(pem string) Option {
	return func(opts *WsOption) {
		opts.pem = pem
	}
}

func WithKey(key string) Option {
	return func(opts *WsOption) {
		opts.key = key
	}
}

func WithHttps(https bool) Option {
	return func(opts *WsOption) {
		opts.https = https
	}
}

func WithWsFunc(handler http.Handler) Option {
	return func(opts *WsOption) {
		opts.handler = handler
	}
}

func WithMethod(m iface.IWsSessionMethod) Option {
	return func(opts *WsOption) {
		opts.method = m
	}
}

func WithHandlerFunc(path string, fn func(http.ResponseWriter, *http.Request), methods string) Option {
	return func(opts *WsOption) {
		opts.handlerFuncs = append(opts.handlerFuncs, HandlerFunc{path, fn, methods})
	}
}
