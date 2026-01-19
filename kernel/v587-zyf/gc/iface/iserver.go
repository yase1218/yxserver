package iface

import "context"

type IServer interface {
	Init(ctx context.Context, opts ...any) error
	Start()
	Stop()
}
