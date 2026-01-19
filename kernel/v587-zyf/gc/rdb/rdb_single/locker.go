package rdb_single

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

// Locker 锁
type Locker struct {
	keys string
	rdb  *RedisSingle
	sid  string
	guid string
}

func (l *Locker) Get() *redis.Client {
	return l.rdb.Get().(*redis.Client)
}

func (l *Locker) GetCtx() context.Context {
	return l.rdb.GetCtx()
}

func (l *Locker) Unlock() {
	if !l.rdb.UnLocks(l.sid, l.keys) {
		log.Error("unlock locker", zap.String("key", l.keys))
	}
}
