package event

import "sync"

type Pool struct {
	sync.Map
}

func (p *Pool) Get(userID int32) (*EventEmitter, bool) {
	v, ok := p.Load(userID)
	if ok {
		return v.(*EventEmitter), false
	}

	v, loaded := p.LoadOrStore(userID, NewEventEmitter(MAX_LISTENER_CNT))

	return v.(*EventEmitter), !loaded
}
