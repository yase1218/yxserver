package mongo

import (
	"context"
	"github.com/qiniu/qmgo"
)

type Mongo struct {
	options *MongoOption
	client  *qmgo.Client
	db      *qmgo.Database

	ctx    context.Context
	cancel context.CancelFunc
}

func NewMongo() *Mongo {
	m := &Mongo{
		options: NewMongoOption(),
	}

	return m
}

func (m *Mongo) Init(ctx context.Context, opts ...any) (err error) {
	m.ctx, m.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(Option)(m.options)
		}
	}

	m.client, err = qmgo.NewClient(context.Background(), &qmgo.Config{Uri: m.options.uri})
	if err != nil {
		return err
	}

	if m.options.db != "" {
		m.db = m.client.Database(m.options.db)
	}

	if err = m.client.Ping(10); err != nil {
		return err
	}

	return nil
}
func (m *Mongo) Get() *qmgo.Client {
	return m.client
}
func (m *Mongo) GetDB() *qmgo.Database {
	return m.db
}
func (m *Mongo) GetCtx() context.Context {
	return m.ctx
}
func (m *Mongo) DB(dbName string) *qmgo.Database {
	return m.client.Database(dbName)
}

func (m *Mongo) Collection(colName string) *qmgo.Collection {
	return m.db.Collection(colName)
}
