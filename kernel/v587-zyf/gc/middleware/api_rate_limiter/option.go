package api_rate_limiter

type LimiterOption struct {
	rateLimit float64 // 每秒允许的请求数
	burst     int     // 突发容量
}

type Option func(o *LimiterOption)

func NewLimiterOption() *LimiterOption {
	return &LimiterOption{}
}

func WithRateLimit(rateLimit float64) Option {
	return func(o *LimiterOption) {
		o.rateLimit = rateLimit
	}
}
func WithBurst(burst int) Option {
	return func(o *LimiterOption) {
		o.burst = burst
	}
}
