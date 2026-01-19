package rdb_single

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var defRedis *RedisSingle

func InitSingle(ctx context.Context, opts ...any) (err error) {
	defRedis = NewRedisSingle()
	if err = defRedis.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *redis.Client {
	return defRedis.Get().(*redis.Client)
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

func Stop() {
	defRedis.Stop()
}
