package event

// BrodcastMsgEvent 广播消息
type BrodcastMsgEvent struct {
	IEvent
	BroadMsg interface{}
	f        EventFunc
}

func NewBrodcastMsgEvent(msg interface{}, cb EventFunc) *BrodcastMsgEvent {
	return &BrodcastMsgEvent{
		BroadMsg: msg,
		f:        cb,
	}
}

func (b *BrodcastMsgEvent) RouteID() int64 {
	return 0
}

func (b *BrodcastMsgEvent) CallBack() EventFunc {
	return b.f
}
