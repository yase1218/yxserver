package event

import (
	"github.com/v587-zyf/gc/errcode"
	"reflect"
)

type ListenerFn func(...any)

const (
	MAX_LISTENER_CNT = 100
)

func NewEventEmitter(maxListeners int) *EventEmitter {
	return &EventEmitter{
		events: make(map[string][]ListenerFn),
		max:    maxListeners,
	}
}

func NewPool() *Pool {
	return &Pool{}
}

func GenListener(fn interface{}) (ListenerFn, error) {
	refValue := reflect.ValueOf(fn)

	if refValue.Kind() != reflect.Func {
		return nil, errcode.ERR_EVENT_PARAM_INVALID
	}

	lfn := func(params ...interface{}) {
		paramRefs := make([]reflect.Value, len(params))

		for k, v := range params {
			paramRefs[k] = reflect.ValueOf(v)
		}

		refValue.Call(paramRefs)
	}

	return lfn, nil
}
