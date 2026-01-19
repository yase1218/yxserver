package event

import (
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/log"
	"reflect"
	"runtime/debug"
	"sync"

	"go.uber.org/zap"
)

type EventEmitter struct {
	events map[string][]ListenerFn

	mu sync.RWMutex

	max int
}

func (e *EventEmitter) On(eventName string, fns ...ListenerFn) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	v, ok := e.events[eventName]
	if !ok {
		v = make([]ListenerFn, 0, 10)
	}

	if len(v)+len(fns) > e.max {
		return errcode.ERR_EVENT_LISTENER_LIMIT
	}

	e.events[eventName] = append(v, fns...)

	return nil
}

func (e *EventEmitter) Off(eventName string, fn ListenerFn) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	fns, ok := e.events[eventName]
	if !ok {
		return errcode.ERR_EVENT_LISTENER_EMPTY
	}

	var idx int = -1
	fnp := reflect.ValueOf(fn).Pointer()
	for k, v := range fns {
		vp := reflect.ValueOf(v).Pointer()
		if fnp == vp {
			idx = k
			break
		}
	}
	if idx < 0 {
		return errcode.ERR_EVENT_LISTENER_NOT_FIND
	}

	newListeners := make([]ListenerFn, 0, len(fns)-1)

	newListeners = append(newListeners, fns[0:idx]...)
	newListeners = append(newListeners, fns[idx+1:]...)

	e.events[eventName] = newListeners
	return nil
}

func (e *EventEmitter) Once(eventName string, fns ...ListenerFn) error {
	v, ok := e.events[eventName]
	if !ok {
		v = make([]ListenerFn, 0, 10)
	}

	if len(v)+len(fns) > e.max {
		return errcode.ERR_EVENT_LISTENER_LIMIT
	}

	wrapFns := make([]ListenerFn, 0, len(fns))
	for _, fn := range fns {
		var wrapFn ListenerFn
		wrapFn = func(params ...interface{}) {
			fn(params...)
			e.Off(eventName, wrapFn)
		}
		wrapFns = append(wrapFns, wrapFn)
	}

	e.On(eventName, wrapFns...)

	return nil
}

func (e *EventEmitter) Emit(eventName string, params ...interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			log.Error("readMsgLoop panic", zap.Any("r", r), zap.String("stack", string(debug.Stack())))
		}
	}()

	e.mu.RLock()
	v, ok := e.events[eventName]
	if !ok {
		e.mu.RUnlock()
		return errcode.ERR_EVENT_LISTENER_EMPTY
	}
	listeners := make([]ListenerFn, len(v))
	copy(listeners, v)
	e.mu.RUnlock()

	for _, v := range listeners {
		v(params...)
	}

	return nil
}
