package iface

import "context"

type IModule interface {
	Name() string
	Init(ctx context.Context, opts ...Option) error
	Start() error
	Run()
	Stop()
}

type IOptions interface {
	ToDo()
}

type Option func(opts IOptions)

type IModuleMgr interface {
	Add(m IModule)
	Get(name string) IModule
	Del(name string)

	Init(ctx context.Context, opts ...Option) error
	Start() error
	Run()
	Stop()
}
