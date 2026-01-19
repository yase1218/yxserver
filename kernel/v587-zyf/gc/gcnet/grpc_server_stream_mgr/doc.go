package grpc_server_stream_mgr

import (
	"google.golang.org/grpc"
	"time"
)

const (
	CHAN_SIZE = 1024 * 1024 * 5

	NIL_SLEEP_TIME    = 5 * time.Second
	NO_MSG_SLEEP_TIME = 50 * time.Millisecond
)

var defGrpcStreamClientMgr *GrpcServerStreamMgr

func InitGrpcClientStream(opts ...any) (err error) {
	defGrpcStreamClientMgr = NewGrpcClientStream()
	if err = defGrpcStreamClientMgr.Init(opts...); err != nil {
		return err
	}

	return nil
}
func Get() *GrpcServerStreamMgr {
	return defGrpcStreamClientMgr
}
func Add(st int32, id uint64, stream grpc.ServerStream) {
	defGrpcStreamClientMgr.Add(st, id, stream)
}

func Del(st int32, id uint64) {
	defGrpcStreamClientMgr.Del(st, id)
}

func GetStreamByType(st int32) map[uint64]grpc.ServerStream {
	return defGrpcStreamClientMgr.GetStreamByType(st)
}

func RandStreamByType(st int32) grpc.ServerStream {
	return defGrpcStreamClientMgr.RandStreamByType(st)
}
