package grpc_msg

import (
	"context"
	"github.com/v587-zyf/gc/iface"
)

var defGrpcMsg *GrpcMsg

func InitGrpcMsg(ctx context.Context, opts ...any) (err error) {
	defGrpcMsg = NewGrpcMsg()
	if err = defGrpcMsg.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func SendToMsg(msg iface.IProtoMessage) {
	defGrpcMsg.SendToMsg(msg)
}

func GetMsg() <-chan iface.IProtoMessage {
	return defGrpcMsg.GetMsg()
}

func Get() *GrpcMsg { return defGrpcMsg }

//func Send2User(userID uint64, msgID int32, msg iface.IProtoMessage) {
//	defGrpcMsg.Send2User(userID, msgID, msg)
//}
//
//func SendErr2User(userID uint64, err error) {
//	defGrpcMsg.SendErr2User(userID, err)
//}
//
//func Broadcast(msgID int32, msg iface.IProtoMessage) {
//	defGrpcMsg.Broadcast(msgID, msg)
//}
