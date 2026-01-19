package tcp_handler

import (
	"context"
	"github.com/v587-zyf/gc/gcnet/tcp_session"
)

var defTcpHandler *TcpHandler

func Init(ctx context.Context, opts ...Option) (err error) {
	defTcpHandler = NewTcpHandler()
	if err = defTcpHandler.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *TcpHandler {
	return defTcpHandler
}

func GetCtx() context.Context {
	return defTcpHandler.ctx
}

func Register(msgID uint32, handler tcp_session.Recv) {
	defTcpHandler.Register(msgID, handler)
}

func GetHandler(msgID uint32) tcp_session.Recv {
	return defTcpHandler.GetHandler(msgID)
}

func HasHandler(msgID uint32) bool {
	return defTcpHandler.HasHandler(msgID)
}
