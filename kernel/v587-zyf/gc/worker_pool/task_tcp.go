package worker_pool

import (
	"github.com/v587-zyf/gc/gcnet/tcp_session"
	"github.com/v587-zyf/gc/iface"
)

func (p *WorkerPool) AssignTcpTask(fn tcp_session.Recv, ss iface.ITcpSession, data any) error {
	return Assign(&TcpTask{
		Func:    fn,
		Session: ss,
		Data:    data,
	})
}

type TcpTask struct {
	Func    tcp_session.Recv
	Session iface.ITcpSession
	Data    any
}

func (t *TcpTask) Do() {
	if t.Func == nil {
		//log.Warn("ws task func is nil", zap.Uint16("msgID", t.Data.(*iface.MessageFrame).MsgID), zap.String("msgName", pb.GetMsgName(t.Data.(*iface.MessageFrame).MsgID)))
		//log.Warn("ws task func is nil", zap.Uint16("msgID", t.Data.(*iface.MessageFrame).MsgID))
		if defaultWorkPoll.options.errHandler != nil {
			defaultWorkPoll.options.errHandler(t.Data)
		}
		return
	}
	t.Func(t.Session, t.Data)
}
