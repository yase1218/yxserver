package worker_pool

import (
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/iface"
)

func (p *WorkerPool) AssignWsTask(fn ws_session.Recv, ss iface.IWsSession, data any) error {
	return Assign(&WsTask{
		Func:    fn,
		Session: ss,
		Data:    data,
	})
}

type WsTask struct {
	Func    ws_session.Recv
	Session iface.IWsSession
	Data    any
}

func (t *WsTask) Do() {
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
