package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

const (
	Nil = redis.Nil
)

func DelByKeys(rc *redis.Client, rCtx context.Context, keys ...string) {
	if _, err := rc.Del(rCtx, keys...).Result(); err != nil {
		log.Error("del rank err", zap.Error(err), zap.Strings("keys", keys))
	}
}
