package http_server

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

var defHttpServer *HttpServer

func Init(ctx context.Context, opts ...any) (err error) {
	defHttpServer = NewHttpServer()
	if err = defHttpServer.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func GetHttpServer() *HttpServer {
	return defHttpServer
}

func InitHttp() (err error) {
	return defHttpServer.InitHttp()
}
func InitHttps() error {
	return defHttpServer.InitHttps()
}

func GetApp() *fiber.App {
	return defHttpServer.GetApp()
}

func Start() {
	defHttpServer.Start()
}

func Stop() {
	defHttpServer.Stop()
}

func Wait() error {
	return defHttpServer.Wait()
}

func Post(path string, fn ResponseHandlerFn) {
	defHttpServer.Post(path, fn)
}

func Get(path string, fn ResponseHandlerFn) {
	defHttpServer.Get(path, fn)
}

func PostOrigin(path string, fn OriginHandlerFn) {
	defHttpServer.PostOrigin(path, fn)
}

func GetOrigin(path string, fn OriginHandlerFn) {
	defHttpServer.GetOrigin(path, fn)
}

func PostOriginDelay(path string, fn OriginHandlerDelayFn) {
	defHttpServer.PostOriginDelay(path, fn)
}

func GetOriginDelay(path string, fn OriginHandlerDelayFn) {
	defHttpServer.GetOriginDelay(path, fn)
}

func Use(fn OriginHandlerFn) {
	defHttpServer.Use(fn)
}

func UseOrigin(args ...any) {
	defHttpServer.UseOrigin(args...)
}
