package tools

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

// WaitExit 等待退出(cirt+c有效,一定加在main函数的最后一行)
func WaitExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("press cirt+c exit app")
	select {
	case sig := <-c:
		fmt.Println(sig, "exit app")
	}
}

// func GoSafe(fnName string, fn func()) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			buf := make([]byte, 4096)
// 			l := runtime.Stack(buf, false)
// 			err := fmt.Errorf("%v: %s", r, buf[:l])
// 			log.Error(fnName+" panic", zap.Error(err))
// 		}
// 	}()

// 	fn()
// }

func GoSafe(fnName string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			err := fmt.Errorf("%v: %s", r, buf[:l])
			log.Error(fnName+" panic", zap.Error(err))
		}
	}()

	fn()
}

func GoSafePost(fnName string, fn func(), post func(string)) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			err := fmt.Errorf("%v: %s", r, buf[:l])
			log.Error(fnName+" panic", zap.Error(err))
			if post != nil {
				post(fnName + "——" + fmt.Sprintf("%v", r))
			}
		}
	}()

	fn()
}

func GoSafeWithParam(fnName string, fn func(interface{}), post func(string), param interface{}) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			err := fmt.Errorf("%v: %s", r, buf[:l])
			log.Error(fnName+" panic", zap.Error(err))
			post(fnName + ":" + fmt.Sprintf("%v", r))
		}
	}()

	fn(param)
}
