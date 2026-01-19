package grpc_msg

import (
	"context"
	"github.com/v587-zyf/gc/iface"
)

const (
	DEF_SIZE = 1024 * 1024 * 5
)

type GrpcMsg struct {
	options *GrpcOption

	ctx    context.Context
	cancel context.CancelFunc

	msgCh chan iface.IProtoMessage
}

func NewGrpcMsg() *GrpcMsg {
	s := &GrpcMsg{
		options: NewGrpcOption(),
	}

	return s
}

func (g *GrpcMsg) Init(ctx context.Context, option ...any) (err error) {
	g.ctx, g.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt.(Option)(g.options)
	}

	if g.options.size != 0 {
		g.msgCh = make(chan iface.IProtoMessage, g.options.size)
	} else {
		g.msgCh = make(chan iface.IProtoMessage, DEF_SIZE)
	}

	return nil
}

func (g *GrpcMsg) SendToMsg(msg iface.IProtoMessage) {
	g.msgCh <- msg
}

//func (g *GrpcMsg) Send2User(userID uint64, msgID int32, msg iface.IProtoMessage) {
//	msgBytes, err := msg.Marshal()
//	if err != nil {
//		log.Error("enterNtf marshal err", zap.Error(err))
//		return
//	}
//
//	content := &server.Send2User{MsgID: msgID, Content: msgBytes}
//	reqBytes, err := handler.GetClientWsHandler().Marshal(server.MsgID_Send2UserId, 0, userID, content)
//	if err != nil {
//		panic(err)
//		return
//	}
//
//	msgData := &server.MessageData{Sender: enums.SERVER_GAME, Receiver: enums.SERVER_GATE, Content: reqBytes}
//	g.SendToMsg(msgData)
//}
//
//func (g *GrpcMsg) SendErr2User(userID uint64, err error) {
//	errNtf := new(pb.ErrNtf)
//	var errCode errcode.ErrCode
//	if errors.As(err, &errCode) {
//		errNtf.ErrNo = errCode.Int32()
//		errNtf.ErrMsg = errCode.Error()
//	} else {
//		errNtf.ErrNo = errcode.ERR_STANDARD_ERR.Int32()
//		errNtf.ErrMsg = err.Error()
//	}
//	g.Send2User(userID, pb.MsgID_Err_NtfId, errNtf)
//}
//
//func (g *GrpcMsg) Broadcast(msgID int32, msg iface.IProtoMessage) {
//	msgBytes, err := msg.Marshal()
//	if err != nil {
//		log.Error("enterNtf marshal err", zap.Error(err))
//		return
//	}
//
//	content := &server.Broadcast{MsgID: msgID, Content: msgBytes}
//	reqBytes, err := handler.GetClientWsHandler().Marshal(server.MsgID_BroadcastId, 0, 0, content)
//	if err != nil {
//		panic(err)
//		return
//	}
//
//	msgData := &server.MessageData{Sender: enums.SERVER_GAME, Receiver: enums.SERVER_GATE, Content: reqBytes}
//	g.SendToMsg(msgData)
//}

func (g *GrpcMsg) GetMsg() <-chan iface.IProtoMessage {
	return g.msgCh
}
