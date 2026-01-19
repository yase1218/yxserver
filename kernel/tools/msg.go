package tools

// import (
// 	"leaf/gate"
// 	"msg"

// 	"google.golang.org/protobuf/proto"
// )

// type Sender func(uint32, interface{})

// var RedirectSend Sender

// func RedirectSendFn(fn Sender) {
// 	RedirectSend = fn
// }

// func SendNotifyMsg(agent gate.Agent, sendMsg interface{}) {
// 	agent.WriteMsg(0, sendMsg)
// }

// func SendMsg(agent gate.Agent, sendMsg interface{}, packetId uint32, err msg.ErrCode) {
// 	if RedirectSend != nil {
// 		RedirectSend(packetId, sendMsg) // 暂不发送统一err 客户端没用到
// 		return
// 	}
// 	msgId := msg.Processor.GetMsgId(sendMsg.(proto.Message))
// 	agent.WriteMsg(packetId, sendMsg)
// 	if err != msg.ErrCode_SUCC && err != msg.ErrCode_ERR_NONE {
// 		errMsg := &msg.NotifyErrMsg{
// 			Id:     msg.MsgId(msgId
// 			Result: err,
// 		}
// 		agent.WriteMsg(packetId, errMsg)
// 	}
// }
