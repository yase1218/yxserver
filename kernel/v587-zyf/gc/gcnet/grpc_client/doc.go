package grpc_client

import (
	"context"
	"google.golang.org/grpc"
)

var defGrpcClient *GrpcClient

func InitGrpcClient(ctx context.Context, opts ...any) (err error) {
	defGrpcClient = NewGrpcClient()
	if err = defGrpcClient.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *GrpcClient {
	return defGrpcClient
}

func GetClient() *grpc.ClientConn {
	return defGrpcClient.GetClient()
}

func GetCtx() context.Context {
	return defGrpcClient.GetCtx()
}

func Start() {
	defGrpcClient.Start()
}

func Stop() {
	defGrpcClient.Stop()
}
