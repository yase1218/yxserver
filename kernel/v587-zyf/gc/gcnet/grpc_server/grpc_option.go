package grpc_server

type GrpcOption struct {
	listenAddr string
}

type Option func(opts *GrpcOption)

func NewGrpcOption() *GrpcOption {
	o := &GrpcOption{}

	return o
}

func WithListenAddr(addr string) Option {
	return func(opts *GrpcOption) {
		opts.listenAddr = addr
	}
}
