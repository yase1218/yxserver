package tcp_server

import "context"

var defTcpSer *TcpServer

func Init(ctx context.Context, opts ...Option) (err error) {
	defTcpSer = NewTcpServer()
	if err = defTcpSer.Init(ctx, opts...); err != nil {
		return err
	}
	return nil
}

func Get() *TcpServer {
	return defTcpSer
}

func GetCtx() context.Context {
	return defTcpSer.ctx
}

func Start() {
	defTcpSer.Start()
}

func Stop() {
	defTcpSer.Stop()
}
