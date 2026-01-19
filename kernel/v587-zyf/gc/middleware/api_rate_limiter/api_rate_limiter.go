package api_rate_limiter

import (
	"context"
	"golang.org/x/time/rate"
	"sync"
)

type APIRateLimiter struct {
	options *LimiterOption

	limiterMap sync.Map

	ctx    context.Context
	cancel context.CancelFunc
}

func NewAPIRateLimiter() *APIRateLimiter {
	return &APIRateLimiter{
		options: NewLimiterOption(),
	}
}

func (arl *APIRateLimiter) Init(ctx context.Context, opts ...any) (err error) {
	arl.ctx, arl.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(Option)(arl.options)
		}
	}

	return
}

func (arl *APIRateLimiter) GetLimiter(api string) *rate.Limiter {
	newLimiter := rate.NewLimiter(rate.Limit(arl.options.rateLimit), arl.options.burst)
	limiter, _ := arl.limiterMap.LoadOrStore(api, newLimiter)
	return limiter.(*rate.Limiter)
}

func (arl *APIRateLimiter) LimitCheck(api string) bool {
	limiter := arl.GetLimiter(api)
	return limiter.Allow()
}
