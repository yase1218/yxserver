package grpc_server_stream_mgr

type GrpcOption struct{}

type Option func(opts *GrpcOption)

func NewGrpcOption() *GrpcOption {
	o := &GrpcOption{}

	return o
}
