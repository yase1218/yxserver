package protobuf

import (
	"context"
	"fmt"
	"github.com/v587-zyf/gc/examples/protobuf/out/server"
	"github.com/v587-zyf/gc/gcnet/grpc_client"
)

func SayHello(name string) (string, error) {
	var (
		err  error
		addr = ":50051"
		ctx  = context.Background()
	)
	if err = grpc_client.InitGrpcClient(ctx, grpc_client.WithListenAddr(addr)); err != nil {
		fmt.Println("grpc client init failed:", err)
		return "", err
	}
	defer grpc_client.Stop()

	client := server.NewHelloServiceClient(grpc_client.GetClient())
	resp, err := client.SayHello(ctx, &server.HelloRequest{Name: name, Enum: server.MyEnum_ENUM_VALUE_1})
	if err != nil {
		return "", err
	}

	return resp.Message, nil
}
