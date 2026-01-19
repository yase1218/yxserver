package rdb_cluster

type RedisClusterOption struct {
	addrs    []string
	username string
	pwd      string
}

type Option func(o *RedisClusterOption)

func NewRedisClusterOption() *RedisClusterOption {
	return &RedisClusterOption{}
}

func WithAddr(addrs []string) Option {
	return func(o *RedisClusterOption) {
		o.addrs = addrs
	}
}

func WithUsername(un string) Option {
	return func(o *RedisClusterOption) {
		o.username = un
	}
}

func WithPwd(pwd string) Option {
	return func(o *RedisClusterOption) {
		o.pwd = pwd
	}
}
