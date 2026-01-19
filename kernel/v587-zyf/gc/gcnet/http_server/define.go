package http_server

import "github.com/gofiber/fiber/v2"

type Ctx struct {
	*fiber.Ctx
}

type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}
type ResponseHandlerFn func(*Ctx) (any, error)

type OriginHandlerFn func(*Ctx) error

type OriginHandlerDelayFn func(*Ctx) (func() error, error) // 带延迟回调的handler 发送返回值后执行回调
