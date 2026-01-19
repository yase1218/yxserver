package elastic

import (
	"context"
	"github.com/olivere/elastic/v7"
)

var defElastic *Elastic

func Init(ctx context.Context, opts ...any) (err error) {
	defElastic = NewElastic()
	if err = defElastic.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *elastic.Client {
	return defElastic.Get()
}
func GetCtx() context.Context {
	return defElastic.ctx
}
