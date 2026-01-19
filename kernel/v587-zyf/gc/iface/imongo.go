package iface

import (
	"context"
	"github.com/qiniu/qmgo"
)

type IMongo interface {
	Init(ctx context.Context, opts ...any) (err error)
	Get() *qmgo.Client
	GetDB() *qmgo.Database
	GetCtx() context.Context
	DB(dbName string) *qmgo.Database
	Collection(colName string) *qmgo.Collection
}
