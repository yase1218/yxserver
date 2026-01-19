package io_out

import (
	"context"
	"errors"
	"kernel/kenum"
	"kernel/metric"
	"kernel/tools"
	"sync"
	"sync/atomic"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type OutFunc func(string, proto.Message) error

type OutMsg struct {
	Subject string
	Msg     proto.Message
}
type OutQueue struct {
	ctx      context.Context
	cancel   context.CancelFunc
	stopChan chan struct{}
	wg       sync.WaitGroup
	state    atomic.Uint32            // 状态 WorkState
	metrics  *metric.ProcessorMetrics // 指标收集

	queue     chan *OutMsg
	outFunc   OutFunc
	panicFunc func(string)
}

func (o *OutQueue) start() error {
	if o.outFunc == nil {
		return errors.New("out queue outFunc nil")
	}
	if !o.state.CompareAndSwap(kenum.WorkState_Idle, kenum.WorkState_Running) {
		return errors.New("out queue can't start, current state : " + kenum.StateToString(o.state.Load()))
	}

	log.Info("out queue start running")
	o.wg.Add(1)
	go tools.GoSafePost("out queue run", func() {
		o.run()
	}, o.panicFunc)

	o.wg.Add(1)
	go tools.GoSafePost("out queue monitor", func() {
		o.monitor()
	}, o.panicFunc)

	log.Info("out queue start success")
	return nil
}

func (o *OutQueue) monitor() {
	defer o.wg.Done()
	defer func() {
		if o.cancel != nil {
			o.cancel()
		}
	}()

	for {
		select {
		case <-o.stopChan:
			// 收到停止信号 处理剩余消息后退出
			log.Info("out queue received stop signal, draining messages")
			o.drainMessages()
			return
		case <-o.ctx.Done():
			// ctx取消 立即退出
			log.Info("out queue context canceled, exiting")
			return
		case <-time.After(time.Second * 5):
			log.Info("out queue size", zap.Int("size", len(o.queue)))
		}
	}
}

func (o *OutQueue) run() {
	defer o.wg.Done()
	defer func() {
		if o.cancel != nil {
			o.cancel()
		}
	}()

	for {
		select {
		case <-o.stopChan:
			// 收到停止信号 处理剩余消息后退出
			log.Info("out queue received stop signal, draining messages")
			o.drainMessages()
			return
		case <-o.ctx.Done():
			// ctx取消 立即退出
			log.Info("out queue context canceled, exiting")
			return
		case m, ok := <-o.queue:
			if !ok {
				log.Error("out queue msg queue channel closed")
				return
			}
			o.outMsg(m)
		}
	}
}

func (o *OutQueue) stop() error {
	timeOut := time.Second * 10
	if !o.state.CompareAndSwap(kenum.WorkState_Running, kenum.WorkState_Stopping) {
		return errors.New("out queue can't stop, current state : " + kenum.StateToString(o.state.Load()))
	}

	log.Info("out queue stopping")

	// 发送停止信号
	close(o.stopChan)

	// 等待goroutine退出
	stopped := make(chan struct{})
	go tools.GoSafePost("out queue wait stop", func() {
		o.wg.Wait()
		close(stopped)
	}, o.panicFunc)

	select {
	case <-stopped:
		o.state.Store(kenum.WorkState_Stopped)
		log.Info("out queue stopped", zap.Any("metrics", o.metrics.Info()))
		return nil
	case <-time.After(timeOut):
		if o.cancel != nil {
			o.cancel()
		}
		log.Warn("out queue stop timeout, forcing shutdown", zap.Duration("timeout", timeOut))
		return errors.New("out queue stop timeout after " + timeOut.String())
	}
}

func (o *OutQueue) drainMessages() {
	log.Info("out queue processing remaining messages during shutdown")
	const maxAttemps = 10
	const internal = 100 * time.Millisecond

	for i := 0; i < maxAttemps; i++ {
		drained := false
		for {
			select {
			case m := <-o.queue:
				o.outMsg(m)
				drained = true
			default:
				if !drained {
					// 本轮没有处理消息，认为已处理完毕
					log.Info("out queue drained all messages", zap.Int("attemp", i+1),
						zap.Int("cur out queue lenth", len(o.queue)))
					return
				}
				// 下一轮尝试
				goto nextRound
			}
		nextRound:
			time.Sleep(internal)
		}
	}

	log.Warn("out queue may have remaining messages after max attempts",
		zap.Int("cur out queue lenth", len(o.queue)))
}

func (o *OutQueue) outMsg(m *OutMsg) {
	if m == nil {
		o.metrics.ProcessErrors.Add(1)
		return
	}
	if o.outFunc == nil {
		o.metrics.ProcessErrors.Add(1)
		return
	}

	o.metrics.ProcessCount.Add(1)
	defer func() {
		if r := recover(); r != nil {
			o.metrics.ProcessErrors.Add(1)
			log.Error("out queue out msg handler panic",
				zap.Any("recover", r),
				zap.Any("message", m))
		}
	}()

	if err := o.outFunc(m.Subject, m.Msg); err != nil {
		log.Error("out msg error", zap.Error(err))
	}
}

func (o *OutQueue) push(m *OutMsg) error {
	o.metrics.PushAttempts.Add(1)

	state := o.state.Load()
	if state != kenum.WorkState_Running {
		o.metrics.PushRejects.Add(1)
		return errors.New("out queue cannot push gate msg, state is : " + kenum.StateToString(state))
	}

	select {
	case o.queue <- m:
		o.metrics.PushSuccess.Add(1)
		return nil
	case <-o.stopChan:
		o.metrics.PushRejects.Add(1)
		return errors.New("out queue is stopping, gate msg reject")
	case <-o.ctx.Done():
		o.metrics.PushRejects.Add(1)
		return errors.New("ctx cancelled, gate msg reject")
	default:
		o.metrics.PushRejects.Add(1)
		return errors.New("out queue gate msg queue is full, gate msg reject")
	}
}
