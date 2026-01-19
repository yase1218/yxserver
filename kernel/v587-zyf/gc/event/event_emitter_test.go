package event

import (
	"testing"
)

//func init() {
//	log.Init(context.Background())
//}

func TestEventEmitter(t *testing.T) {
	p := NewPool()

	e, _ := p.Get(123)

	fn1 := func(a int, b string, c int) {
		t.Logf("callback 1 %d %s %d", a, b, c)
	}

	fn2 := func(a int, b string, c int) {
		t.Logf("callback 2 %d %s %d", a, b, c)
	}

	wrapf1, _ := GenListener(fn1)
	wrapf2, _ := GenListener(fn2)

	e.On("abc", wrapf1)
	e.Once("aa", wrapf2)

	e.Emit("abc", 10, "abc", 123)
	e.Emit("aa", 10, "abc", 123)
	e.Emit("aa", 10, "abc", 123)
	e.Emit("aa", 10, "abc", 123)

	e.On("login", func(a ...any) {
		uid := a[0].(int32)
		t.Logf("login %d", uid)
	})
	e.Emit("login", int32(123))
}
