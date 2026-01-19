package db_model

import (
	"context"
	"github.com/v587-zyf/gc/iface"
	"gopkg.in/gorp.v1"
)

var defDb *DBModelMap

func Init(ctx context.Context, opts ...any) (err error) {
	defDb = NewDBModelMap()
	if err = defDb.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Register(dbKey string, model iface.IDBModel, initer func(dbMap *gorp.DbMap)) {
	defDb.Register(dbKey, model, initer)
}
func Start() error {
	return defDb.Start()
}
