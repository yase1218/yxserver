package http_server

import (
	"encoding/json"
	"errors"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/log"
)

func SendErrCode(c *fiber.Ctx, errCode errcode.ErrCode) error {
	resp := Response{
		Code: errCode.Int(),
		Msg:  errCode.Error(),
		Data: nil,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	// return c.SendString(string(out))

	c.Response().SetStatusCode(200)
	c.Response().SetBodyRaw(out)
	return nil
}

func SendError(c *fiber.Ctx, err error) error {
	resp := Response{
		Code: errcode.ERR_STANDARD_ERR.Int(),
		Msg:  err.Error(),
		Data: nil,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return SendErrCode(c, errcode.ERR_JSON_MARSHAL_ERR)
	}

	c.Response().SetStatusCode(200)
	c.Response().SetBodyRaw(out)
	return nil
}

func SendResponse(c *fiber.Ctx, data any) error {
	resp := Response{
		Code:    errcode.ERR_SUCCEED.Int(),
		Msg:     "ok",
		Success: true,
		Data:    data,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	c.Response().SetStatusCode(200)
	c.Response().SetBodyRaw(out)

	return nil
}

func NewResponseHandlerFn(fn ResponseHandlerFn) fiber.Handler {
	return func(c *fiber.Ctx) (retErr error) {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 1<<10)
				runtime.Stack(buf, true)
				if err, ok := r.(error); ok {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.String("err", err.Error()), zap.ByteString("core", buf))
				} else if err, ok := r.(string); ok {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.String("err", err), zap.ByteString("core", buf))
				} else {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.Reflect("err", err), zap.ByteString("core", buf))
				}
				retErr = SendErrCode(c, errcode.ERR_SERVER_INTERNAL)

				return
			}
		}()

		ctx := &Ctx{Ctx: c}
		resp, err := fn(ctx)
		if err != nil {
			var errCode errcode.ErrCode
			if errors.As(err, &errCode) && !errors.Is(errCode, errcode.ERR_SUCCEED) {
				return SendErrCode(c, errCode)
			} else {
				return SendError(c, err)
			}
		}
		return SendResponse(c, resp)
	}
}

func NewOriginHandlerFn(fn func(*Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) (retErr error) {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 1<<10)
				runtime.Stack(buf, true)
				if err, ok := r.(error); ok {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.String("err", err.Error()), zap.ByteString("core", buf))
				} else if err, ok := r.(string); ok {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.String("err", err), zap.ByteString("core", buf))
				} else {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.Reflect("err", err), zap.ByteString("core", buf))
				}

				retErr = errcode.ERR_SERVER_INTERNAL
				return
			}
		}()

		ctx := &Ctx{Ctx: c}
		retErr = fn(ctx)

		return retErr
	}
}

func NewOriginHandlerDelayFn(fn func(*Ctx) (func() error, error)) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) (retErr error) {
		var cb func() error
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 1<<10)
				runtime.Stack(buf, true)
				if err, ok := r.(error); ok {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.String("err", err.Error()), zap.ByteString("core", buf))
				} else if err, ok := r.(string); ok {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.String("err", err), zap.ByteString("core", buf))
				} else {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.Reflect("err", err), zap.ByteString("core", buf))
				}

				retErr = errcode.ERR_SERVER_INTERNAL
				return
			}
		}()

		ctx := &Ctx{Ctx: c}
		cb, retErr = fn(ctx)
		if retErr != nil {
			return
		}
		if cb != nil {
			retErr = cb()
		}

		return retErr
	}
}
