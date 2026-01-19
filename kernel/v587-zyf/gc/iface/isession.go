package iface

import (
	"context"
	"net"
	"time"
)

type ITcpSession interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Remove(key string)

	GetID() uint64
	SetID(id uint64)

	Start()
	Close() error

	GetConn() net.Conn
	GetCtx() context.Context

	SendMsg(fn func(args ...any) ([]byte, error), args ...any) error

	DoSomething(fn func(args ...any) bool) bool
	CheckSomething(fn func(args ...any) bool) bool

	IsHeartbeatTimeout(now time.Time) bool
}

type IWsSession interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Remove(key string)

	GetID() uint64
	SetID(id uint64)

	Start()
	Close() error

	GetConn() IConn
	GetCtx() context.Context
	GetRemoteAddr() string

	SendMsg(fn func(args ...any) ([]byte, error), args ...any) error

	DoSomething(fn func(args ...any) bool) bool
	CheckSomething(fn func(args ...any) bool) bool

	IsHeartbeatTimeout(now time.Time) bool
}

type ITcpSessionMethod interface {
	Name() string
	Start(ss ITcpSession)
	Recv(conn ITcpSession, data any)
	Stop(ss ITcpSession)
}
type IWsSessionMethod interface {
	Name() string
	Start(ss IWsSession)
	Recv(conn IWsSession, data any)
	Stop(ss IWsSession)
}
