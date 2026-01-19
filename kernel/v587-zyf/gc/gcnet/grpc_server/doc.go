package grpc_server

import (
	"context"
	"google.golang.org/grpc"
)

var defGrpcServer *GrpcServer

func InitGrpcClient(ctx context.Context, opts ...any) (err error) {
	defGrpcServer = NewGrpcServer()
	if err = defGrpcServer.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *GrpcServer {
	return defGrpcServer
}

func GetCtx() context.Context {
	return defGrpcServer.GetCtx()
}

func GetServer() *grpc.Server {
	return defGrpcServer.GetServer()
}

func Start() {
	defGrpcServer.Start()
}
