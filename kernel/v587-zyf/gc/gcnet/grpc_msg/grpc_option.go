package grpc_msg

type GrpcOption struct {
	size int
}

type Option func(opts *GrpcOption)

func NewGrpcOption() *GrpcOption {
	o := &GrpcOption{}

	return o
}

func WithSize(size int) Option {
	return func(opts *GrpcOption) {
		opts.size = size
	}
}
