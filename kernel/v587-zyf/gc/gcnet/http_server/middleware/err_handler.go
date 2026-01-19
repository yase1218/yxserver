package middleware

import (
	"errors"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/gcnet/http_server"
)

func NewErrHandler() http_server.OriginHandlerFn {
	return func(c *http_server.Ctx) (err error) {
		if err = c.Next(); err != nil {
			var errCode errcode.ErrCode
			if errors.As(err, &errCode) {
				if err = http_server.SendErrCode(c.Ctx, errCode); err != nil {
					return err
				}
			}
		}
		return
	}
}
