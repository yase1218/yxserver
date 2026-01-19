package grpc_server_stream_mgr

import (
	"google.golang.org/grpc"
	"sync"
)

type GrpcServerStreamMgr struct {
	options *GrpcOption

	mu      sync.Mutex
	streams map[int32]map[uint64]grpc.ServerStream // serverType:id:stream
}

func NewGrpcClientStream() *GrpcServerStreamMgr {
	return &GrpcServerStreamMgr{
		options: NewGrpcOption(),
		streams: make(map[int32]map[uint64]grpc.ServerStream),
	}
}

var _ GrpcServerStreamMgr

func (g *GrpcServerStreamMgr) Init(option ...any) (err error) {
	for _, opt := range option {
		opt.(Option)(g.options)
	}

	return nil
}

func (g *GrpcServerStreamMgr) Add(st int32, id uint64, stream grpc.ServerStream) {
	g.mu.Lock()
	defer g.mu.Unlock()

	streamMap, ok := g.streams[st]
	if !ok {
		g.streams[st] = make(map[uint64]grpc.ServerStream)
		streamMap = g.streams[st]
	}
	streamMap[id] = stream
}

func (g *GrpcServerStreamMgr) Del(st int32, id uint64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.streams[st]; ok {
		delete(g.streams[st], id)
	}
}

func (g *GrpcServerStreamMgr) GetStreamByType(st int32) map[uint64]grpc.ServerStream {
	return g.streams[st]
}

func (g *GrpcServerStreamMgr) RandStreamByType(st int32) grpc.ServerStream {
	m, ok := g.streams[st]
	if !ok {
		return nil
	}

	for _, stream := range m {
		return stream
	}
	return nil
}
