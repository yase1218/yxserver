package worker_pool

import (
	"time"

	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/iface"
)

type DelaySendTask struct {
	Delay   time.Duration
	Func    ws_session.Recv
	Session iface.IWsSession
	Data    any
}

func (t *DelaySendTask) Do() {
	time.Sleep(t.Delay)
	t.Func(t.Session, t.Data)
}

func (p *WorkerPool) AssignDelaySendTask(delay time.Duration, fn ws_session.Recv, ss iface.IWsSession, data any) error {
	return Assign(&DelaySendTask{
		Delay:   delay,
		Func:    fn,
		Session: ss,
		Data:    data,
	})
}
