package ws_server

import (
	"context"
)

var defWsSer *WsServer

func Init(ctx context.Context, opts ...Option) (err error) {
	defWsSer = NewWsServer()
	if err = defWsSer.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *WsServer {
	return defWsSer
}

func GetCtx() context.Context {
	return defWsSer.ctx
}

func Start() {
	defWsSer.Start()
}

func Stop() {
	defWsSer.Stop()
}
