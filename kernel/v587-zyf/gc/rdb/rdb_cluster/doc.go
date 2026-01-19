package rdb_cluster

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var defRedis *RedisCluster

func InitCluster(ctx context.Context, opts ...any) (err error) {
	defRedis = NewRedisCluster()
	if err = defRedis.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *redis.ClusterClient {
	return defRedis.Get().(*redis.ClusterClient)
}

func GetCtx() context.Context {
	return defRedis.GetCtx()
}

// Locks 分布式锁 锁住需要锁的key 返回true加锁成功 false 加锁失败
func GetLocker(token string, keys string) *Locker {
	uk := fmt.Sprint("{lock}", keys)
	if !defRedis.Locks(fmt.Sprint(token, "::", time.Now().Unix()), uk) {
		return nil
	}

	// 加锁成功的 记录一个guid
	return &Locker{
		keys: uk,
		rdb:  defRedis,
		sid:  token,
	}
}

func UnLocks(sid string, keys string) bool {
	return defRedis.UnLocks(sid, keys)
}
