package mongo

import (
	"context"
	"github.com/qiniu/qmgo"
)

var defMongo *Mongo

func Init(ctx context.Context, opts ...any) (err error) {
	defMongo = NewMongo()
	if err = defMongo.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *Mongo {
	return defMongo
}

func DB(dbName string) *qmgo.Database {
	return defMongo.client.Database(dbName)
}

func Collection(colName string) *qmgo.Collection {
	return defMongo.db.Collection(colName)
}
