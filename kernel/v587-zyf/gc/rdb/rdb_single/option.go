package rdb_single

type RedisSingleOption struct {
	addr     string
	username string
	pwd      string
}

type Option func(o *RedisSingleOption)

func NewRedisSingleOption() *RedisSingleOption {
	return &RedisSingleOption{}
}

func WithAddr(addr string) Option {
	return func(o *RedisSingleOption) {
		o.addr = addr
	}
}

func WithUsername(un string) Option {
	return func(o *RedisSingleOption) {
		o.username = un
	}
}

func WithPwd(pwd string) Option {
	return func(o *RedisSingleOption) {
		o.pwd = pwd
	}
}
