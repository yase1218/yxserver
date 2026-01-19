package grpc_server

import (
	"context"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

type GrpcServer struct {
	options *GrpcOption

	ctx    context.Context
	cancel context.CancelFunc

	listener net.Listener
	server   *grpc.Server
}

func NewGrpcServer() *GrpcServer {
	s := &GrpcServer{
		options: NewGrpcOption(),
	}

	return s
}

func (s *GrpcServer) Init(ctx context.Context, option ...any) (err error) {
	s.ctx, s.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt.(Option)(s.options)
	}

	s.listener, err = net.Listen("tcp", s.options.listenAddr)
	if err != nil {
		log.Error("net listen err", zap.Error(err))
		return
	}

	keepalivePolicy := keepalive.EnforcementPolicy{
		MinTime:             30 * time.Second,
		PermitWithoutStream: true,
	}
	keepaliveOptions := keepalive.ServerParameters{
		Time:    30 * time.Second,
		Timeout: 20 * time.Second,
	}
	s.server = grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(keepalivePolicy),
		grpc.KeepaliveParams(keepaliveOptions),
	)

	return nil
}

func (s *GrpcServer) GetServer() *grpc.Server {
	return s.server
}
func (s *GrpcServer) GetCtx() context.Context {
	return s.ctx
}

func (s *GrpcServer) Start() {
	err := s.server.Serve(s.listener)
	if err != nil {
		log.Error("grpc server start err", zap.Error(err))
		return
	}
}

func (s *GrpcServer) Stop() {
	s.listener.Close()
	s.server.Stop()
}
