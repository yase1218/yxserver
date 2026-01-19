package ws_handler

import (
	"context"
	"github.com/v587-zyf/gc/gcnet/ws_session"
)

var defWsHandler *WsHandler

func Init(ctx context.Context, opts ...any) (err error) {
	defWsHandler = NewWsHandler()
	if err = defWsHandler.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *WsHandler {
	return defWsHandler
}

func GetCtx() context.Context {
	return defWsHandler.ctx
}

func Register(msgID uint32, handler ws_session.Recv) {
	defWsHandler.Register(msgID, handler)
}

func GetHandler(msgID uint32) ws_session.Recv {
	return defWsHandler.GetHandler(msgID)
}

func HasHandler(msgID uint32) bool {
	return defWsHandler.HasHandler(msgID)
}
