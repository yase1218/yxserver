package service

import (
	"gameserver/internal/config"
	"gameserver/internal/io_out"
	"msg"

	"google.golang.org/protobuf/proto"
)

func PublisInterMsg(subject string, msgId uint32, message proto.Message) {
	interMsg := &msg.RequestCommonInterMsg{
		ServerId: int64(config.Conf.ServerId),
		MsgId:    msgId,
	}

	data, _ := proto.Marshal(message)
	interMsg.Data = data
	//nats.Publish(subject, interMsg)
	io_out.Push(&io_out.OutMsg{
		Subject: subject,
		Msg:     interMsg,
	})
}
