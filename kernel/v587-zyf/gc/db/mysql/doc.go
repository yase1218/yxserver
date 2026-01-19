package mysql

import (
	"context"
	"gorm.io/gorm"
)

var defDb *Mysql

func Init(ctx context.Context, opts ...any) (err error) {
	defDb = NewMysql()
	if err = defDb.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func GetDB() *gorm.DB {
	return defDb.db
}
func GetCtx() context.Context {
	return defDb.ctx
}

func AddAutoCreateDb(db any) {
	defDb.AddAutoCreateDb(db)
}
func AutoCreateDbs() (err error) {
	return defDb.AutoCreateDbs()
}
