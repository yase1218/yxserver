package iface

import (
	"context"
	"google.golang.org/grpc"
)

type IGrpc interface {
	Init(ctx context.Context, opts ...any) error
	Start()
}

type IGrpcStream interface {
	Init(ctx context.Context, opts ...any) error
	SetStream(stream *grpc.ClientStream)

	GetID() uint64
	GetStream() *grpc.ClientStream
	GetCtx() context.Context
}

type IGrpcMsg interface {
	Init(ctx context.Context, option ...any) (err error)
	SendToMsg(msg IProtoMessage)

	//Send2User(userID uint64, msgID int32, msg IProtoMessage)
	//SendErr2User(userID uint64, err error)
	//Broadcast(msgID int32, msg IProtoMessage)
}
