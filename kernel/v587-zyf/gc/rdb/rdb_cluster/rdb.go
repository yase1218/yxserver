package rdb_cluster

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
	"strings"
	"time"
)

type RedisCluster struct {
	options *RedisClusterOption
	client  *redis.ClusterClient

	ctx    context.Context
	cancel context.CancelFunc
}

func NewRedisCluster() *RedisCluster {
	rs := &RedisCluster{
		options: NewRedisClusterOption(),
	}

	return rs
}

func (r *RedisCluster) Init(ctx context.Context, opts ...any) (err error) {
	r.ctx, r.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(Option)(r.options)
		}
	}
	r.client = redis.NewClusterClient(&redis.ClusterOptions{
		Username:        r.options.username,
		Addrs:           r.options.addrs,  // 集群地址
		PoolSize:        3000,             // Redis连接池大小
		MaxRetries:      10,               // 最大重试次数
		MinIdleConns:    60,               // 空闲连接数量
		PoolTimeout:     10 * time.Second, // 空闲链接超时时间
		ConnMaxLifetime: 60 * time.Second, // 连接存活时长
		DialTimeout:     15 * time.Second, // 连接建立超时时间，默认5秒。
		ReadTimeout:     7 * time.Second,  // 读超时，默认3秒， -1表示取消读超时
		WriteTimeout:    7 * time.Second,  // 写超时，默认等于读超时
		Password:        r.options.pwd,    // redis 密码
	})

	if err = r.client.Ping(r.ctx).Err(); err != nil {
		log.Error("redis cluster ping err", zap.Error(err))
		return
	}

	return nil

}

func (r *RedisCluster) Get() any {
	return r.client
}

func (r *RedisCluster) GetCtx() context.Context {
	return r.ctx
}

func (r *RedisCluster) Locks(token string, keys string) bool {
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

func (r *RedisCluster) UnLocks(sid string, keys string) bool {
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
