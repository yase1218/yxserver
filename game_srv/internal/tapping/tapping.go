package tapping

import (
	"context"
	"errors"
	"gameserver/internal/config"
	"gameserver/internal/db"
	"kernel/dao"
	"kernel/kenum"
	"kernel/tda"
	"kernel/tools"
	"sync"
	"sync/atomic"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

const (
	TapSwitch       = iota
	TapSwitchShushu // 只有数数打点
	TapSwitchLocal  // 只有本地打点
	TapSwitchAll    // 全都打点
)

var (
	TapEventId uint64
)

type TapFunc func(data *TapData) error

type TapData struct {
	AccountId  string
	DistinctId string
	EventName  string
	ChannelId  uint32
	Switch     int32 //1只有sdk打点，2只有本地打点，3sdk和本地都打点
	Data       any
}

type TapQueue struct {
	ctx      context.Context
	cancel   context.CancelFunc
	stopChan chan struct{}
	wg       sync.WaitGroup
	state    atomic.Uint32 // 状态 WorkState

	queue     chan *TapData
	tapFunc   TapFunc
	panicFunc func(string)
}

func (o *TapQueue) start() error {
	if o.tapFunc == nil {
		return errors.New("tap queue tapFunc nil")
	}
	if !o.state.CompareAndSwap(kenum.WorkState_Idle, kenum.WorkState_Running) {
		return errors.New("tap queue can't start, current state : " + kenum.StateToString(o.state.Load()))
	}

	log.Info("tap queue start running")
	o.wg.Add(1)
	go tools.GoSafePost("tap queue run", func() {
		o.run()
	}, o.panicFunc)

	log.Info("tap queue start success")
	return nil
}

func (o *TapQueue) run() {
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
			log.Info("tap queue received stop signal, draining messages")
			o.drainMessages()
			return
		case <-o.ctx.Done():
			// ctx取消 立即退出
			log.Info("tap queue context canceled, exiting")
			return
		case m, ok := <-o.queue:
			if !ok {
				log.Error("tap queue msg queue channel closed")
				return
			}
			o.processTap(m)
		}
	}
}

func (o *TapQueue) stop() error {
	timeOut := time.Second * 10
	if !o.state.CompareAndSwap(kenum.WorkState_Running, kenum.WorkState_Stopping) {
		return errors.New("tap queue can't stop, current state : " + kenum.StateToString(o.state.Load()))
	}

	log.Info("tap queue stopping")

	// 发送停止信号
	close(o.stopChan)

	// 等待goroutine退出
	stopped := make(chan struct{})
	go tools.GoSafePost("tap queue wait stop", func() {
		o.wg.Wait()
		close(stopped)
	}, o.panicFunc)

	select {
	case <-stopped:
		o.state.Store(kenum.WorkState_Stopped)
		return nil
	case <-time.After(timeOut):
		if o.cancel != nil {
			o.cancel()
		}
		log.Warn("tap queue stop timeout, forcing shutdown", zap.Duration("timeout", timeOut))
		return errors.New("stop timeout after " + timeOut.String())
	}
}

func (o *TapQueue) drainMessages() {
	log.Info("tap queue processing remaining messages during shutdown")
	const maxAttemps = 10
	const internal = 100 * time.Millisecond

	for i := 0; i < maxAttemps; i++ {
		drained := false
		for {
			select {
			case m := <-o.queue:
				o.processTap(m)
				drained = true
			default:
				if !drained {
					// 本轮没有处理消息，认为已处理完毕
					log.Info("tap queue drained all messages", zap.Int("attemp", i+1),
						zap.Int("cur tap queue lenth", len(o.queue)))
					return
				}
				// 下一轮尝试
				goto nextRound
			}
		nextRound:
			time.Sleep(internal)
		}
	}

	log.Warn("tap queue may have remaining messages after max attempts",
		zap.Int("cur tap queue lenth", len(o.queue)))
}

func (o *TapQueue) processTap(m *TapData) {
	if m == nil {
		return
	}
	if o.tapFunc == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("tap queue tap msg handler panic",
				zap.Any("recover", r),
				zap.Any("message", m))
		}
	}()

	if err := o.tapFunc(m); err != nil {
		log.Error("tap msg error", zap.Error(err))
	}
}

func (o *TapQueue) push(m *TapData) error {
	state := o.state.Load()
	if state != kenum.WorkState_Running {
		return errors.New("tap queue cannot push, state is : " + kenum.StateToString(state))
	}

	select {
	case o.queue <- m:
		return nil
	case <-o.stopChan:
		return errors.New("tap queue is stopping, reject")
	case <-o.ctx.Done():
		return errors.New("ctx cancelled, reject")
	default:
		return errors.New("tap queue is full, tap queue reject")
	}
}

func ProcessTapData(data *TapData) error {
	log.Debug("ProcessTapData", zap.Any("data", data))
	if data.Switch == 1 || data.Switch == 3 && config.Conf.Tap.Switch == 1 || config.Conf.Tap.Switch == 3 { // sdk打点
		log.Debug("ProcessTapData sdk", zap.Any("data", data))
		if dataMap, err := tda.FlattenStructToMap(data); err != nil {
			log.Error("ProcessTapData tap data FlattenStructToMap", zap.Any("data", data), zap.Error(err))
		} else {
			if err = tda.GetTa().Track(data.AccountId, data.DistinctId, data.EventName, dataMap); err != nil {
				return err
			}
		}
	}
	if data.Switch == 2 || data.Switch == 3 && config.Conf.Tap.Switch == 2 || config.Conf.Tap.Switch == 3 { // 本地打点
		log.Debug("ProcessTapData local", zap.Any("data", data))
		// mongo落地
		if err := CreateEvent(data); err != nil {
			log.Error("ProcessTapData create event error", zap.Any("data", data), zap.Error(err))
		}
	}
	return nil
}

func CreateEvent(data *TapData) error { // 异步落地
	TapEventId++
	dbName := config.Conf.GlobalMongo.DB
	collection := data.EventName + "_event"
	//data.Data["accountid"] = data.AccountId
	//data.Data["distinctid"] = data.DistinctId
	//data.Data["channelid"] = data.ChannelId
	//data.Data["create_at"] = time.Now()
	op := &dao.WriteOperation{
		Database:   dbName,
		Collection: collection,
		Type:       dao.Insert,
		Uuid:       TapEventId,
		Document:   data.Data,
		Tms:        time.Now().UnixMilli(),
		Cover:      true,
	}
	err := db.GlobalMongoWriter.AsyncWrite(op)
	return err
}
