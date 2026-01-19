package http_server

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"kernel/tools"
	"net"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/v587-zyf/gc/log"
)

type HttpServer struct {
	options *HttpOption

	ln net.Listener

	app *fiber.App

	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup
}

func NewHttpServer() *HttpServer {
	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  90 * time.Second,
		//DisableKeepalive:      true,
		DisableStartupMessage: true,
		//Prefork:               true,
		//ErrorHandler:          config.ErrorHandler,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// nothing to do
			return nil
		},
	})
	// app.Server().KeepHijackedConns = true

	s := &HttpServer{
		options: NewHttpOption(),
		app:     app,
	}

	return s
}

func (s *HttpServer) Init(ctx context.Context, opts ...any) (err error) {
	s.ctx, s.cancel = context.WithCancel(ctx)

	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(Option)(s.options)
		}
	}

	//if s.options.allowOrigins != "" {
	//	s.app.Use(cors.New(cors.Config{
	//		AllowOrigins:     s.options.allowOrigins,              // 只允许来自这些特定源的请求
	//		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",       // 允许的 HTTP 方法
	//		AllowHeaders:     "Authorization,Content-Type,Accept", // 只允许这些特定的头部
	//		AllowCredentials: true,                                // 允许发送 cookies 和其他凭据
	//	}))
	//}

	if s.options.isHttps {
		err = s.InitHttps()
	} else {
		err = s.InitHttp()
	}

	return nil
}

func (s *HttpServer) InitHttp() (err error) {
	s.ln, err = net.Listen("tcp", s.options.listenAddr)
	if err != nil {
		log.Error("net listen err", zap.Error(err))
		return
	}

	return nil
}
func (s *HttpServer) InitHttps() error {
	// 加载证书和私钥
	cert, err := tls.LoadX509KeyPair(s.options.pem, s.options.key)
	if err != nil {
		log.Error("load https cert err", zap.Error(err))
		return err
	}

	// 配置TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// 创建一个TCP监听器
	//s.ln, err = tls.Listen("tcp", ":443", tlsConfig)
	s.ln, err = tls.Listen("tcp", s.options.listenAddr, tlsConfig)
	if err != nil {
		log.Error("net listen err", zap.Error(err))
		return err
	}

	return nil
}

func (s *HttpServer) GetApp() *fiber.App {
	return s.app
}

func (s *HttpServer) Start() {
	s.wg.Add(1)

	go tools.GoSafe("http_server start listener", func() {
		err := s.app.Listener(s.ln)
		if err != nil {
			log.Error("httpserver stopped", zap.Error(err))
		} else {
			log.Info("httpserver stopped")
		}

		s.wg.Done()
	})

	return
}

func (s *HttpServer) Stop() {}

func (s *HttpServer) Wait() error {
	s.wg.Wait()

	return nil
}

func (s *HttpServer) Post(path string, fn ResponseHandlerFn) {
	s.app.Post(path, NewResponseHandlerFn(fn))
}

func (s *HttpServer) Get(path string, fn ResponseHandlerFn) {
	s.app.Get(path, NewResponseHandlerFn(fn))
}

func (s *HttpServer) PostOrigin(path string, fn OriginHandlerFn) {
	s.app.Post(path, NewOriginHandlerFn(fn))
}

func (s *HttpServer) GetOrigin(path string, fn OriginHandlerFn) {
	s.app.Get(path, NewOriginHandlerFn(fn))
}

func (s *HttpServer) PostOriginDelay(path string, fn OriginHandlerDelayFn) {
	s.app.Post(path, NewOriginHandlerDelayFn(fn))
}

func (s *HttpServer) GetOriginDelay(path string, fn OriginHandlerDelayFn) {
	s.app.Get(path, NewOriginHandlerDelayFn(fn))
}

func (s *HttpServer) Use(fn OriginHandlerFn) {
	s.app.Use(NewOriginHandlerFn(fn))
}

// Use registers a middleware route that will match requests
// with the provided prefix (which is optional and defaults to "/").
//
//	app.Use(func(c *fiber.Ctx) error {
//	     return c.Next()
//	})
//	app.Use("/api", func(c *fiber.Ctx) error {
//	     return c.Next()
//	})
//	app.Use("/api", handler, func(c *fiber.Ctx) error {
//	     return c.Next()
//	})
//
// This method will match all HTTP verbs: GET, POST, PUT, HEAD etc...
func (s *HttpServer) UseOrigin(args ...any) {
	s.app.Use(args...)
}
