package protobuf

import (
	"context"
	"fmt"
	"github.com/v587-zyf/gc/examples/protobuf/out/server"
	"github.com/v587-zyf/gc/gcnet/grpc_server"
	"log"
)

type Server struct {
	server.UnimplementedHelloServiceServer
}

// UnimplementedHelloServiceServer must be embedded to have forward compatible implementations.
var _ server.HelloServiceServer = (*Server)(nil)

func (s *Server) SayHello(ctx context.Context, req *server.HelloRequest) (*server.HelloResponse, error) {
	return &server.HelloResponse{Message: "Hello " + req.GetName(), Enum: req.GetEnum()}, nil
}
func (s *Server) mustEmbedUnimplementedHelloServiceServer() {}

func StartServer() (err error) {
	var (
		addr       = ":50051"
		ser        = new(Server)
		ctx        = context.Background()
		GrpcServer = grpc_server.NewGrpcServer()
	)

	if err = GrpcServer.Init(ctx, grpc_server.WithListenAddr(addr)); err != nil {
		fmt.Println("grpc init err:", err)
		return
	}
	server.RegisterHelloServiceServer(GrpcServer.GetServer(), ser)

	go GrpcServer.Start()

	log.Println("Starting gRPC server on ", addr)
	return nil
}
