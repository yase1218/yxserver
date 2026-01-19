package api_rate_limiter

import (
	"context"
)

var defLimiter *APIRateLimiter

func Init(ctx context.Context, opts ...any) (err error) {
	defLimiter = NewAPIRateLimiter()
	if err = defLimiter.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func LimitCheck(api string) bool {
	return defLimiter.LimitCheck(api)
}

//func ExampleRateLimitCheck(c *gin.Context) {
//	ip := c.ClientIP()
//	if !GetLimiterInstance().LimitCheck(ip) {
//		render.RespJson(c, define.Code_RateLimitExceeded, nil)
//		c.Abort()
//		return
//	}
//	c.Next()
//}
