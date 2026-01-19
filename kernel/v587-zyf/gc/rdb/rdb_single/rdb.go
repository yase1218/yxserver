package rdb_single

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
)

type RedisSingle struct {
	options *RedisSingleOption
	client  *redis.Client

	ctx    context.Context
	cancel context.CancelFunc
}

func NewRedisSingle() *RedisSingle {
	rs := &RedisSingle{
		options: NewRedisSingleOption(),
	}

	return rs
}

func (r *RedisSingle) Init(ctx context.Context, opts ...any) (err error) {
	r.ctx, r.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(Option)(r.options)
		}
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:         r.options.addr,
		Username:     r.options.username,
		Password:     r.options.pwd,
		DB:           0,
		PoolSize:     190,
		MinIdleConns: 20,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		PoolTimeout: 4 * time.Second, // 获取连接超时
	})
	if err = r.client.Ping(r.ctx).Err(); err != nil {
		return
	}

	return nil

}

func (r *RedisSingle) Get() any {
	return r.client
}

func (r *RedisSingle) GetCtx() context.Context {
	return r.ctx
}

func (r *RedisSingle) Locks(token string, keys string) bool {
	// 这个函数最多可以延迟 2s
	// 如果被锁住 会 隔 0.2s再请求一次 最多10次
	// 有错误就直接返回 无错误才会执行10次
	for i := 0; i < 10; i++ {
		ok, err := r.client.SetNX(r.ctx, keys, token, time.Second*15).Result()
		if err != nil {
			// 有错误直接返回
			log.Error("redis setnx err", zap.String("redisKey", keys), zap.String("redisValue", token))
			return false
		}

		// 加锁成功
		if ok {
			return true
		}

		// 未加锁成功 可能有其他程序正在使用 等待0.2s再次尝试
		time.Sleep(time.Millisecond * 200)
	}

	log.Error("wait long time, lock err", zap.String("redisKey", keys), zap.String("redisValue", token))
	// 多次尝试后依然被锁着 返回加锁失败
	return false
}

func (r *RedisSingle) UnLocks(sid string, keys string) bool {
	token, err := r.client.Get(r.ctx, keys).Result()
	if err != nil {
		log.Error("redis get err", zap.String("redisKey", keys))
		return false
	}

	tokenArray := strings.Split(token, "::")
	// 10s内其他服务无法删除这个锁
	if (tokenArray[0] != sid) && (time.Now().Unix() < utils.StrToInt64(tokenArray[1])+15) {
		return false
	}

	err = r.client.Del(r.ctx, keys).Err()
	if err != nil {
		log.Error("redis del err", zap.String("redisKey", keys), zap.Error(err))
		return false
	}

	return true
}

func (r *RedisSingle) Stop() {
	r.cancel()
	r.client.Close()
}
