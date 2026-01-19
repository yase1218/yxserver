package module

import (
	"context"
	"github.com/v587-zyf/gc/iface"
)

type DefModule struct {
	ctx context.Context
}

func (m *DefModule) Name() string {
	return ""
}

func (m *DefModule) Init(ctx context.Context, opts ...iface.Option) error {
	m.ctx = ctx
	return nil
}

func (m *DefModule) Start() error {
	return nil
}

func (m *DefModule) Run() {}

func (m *DefModule) Stop() {}

func (m *DefModule) GetCtx() context.Context {
	return m.ctx
}
