package async

import (
	"context"
	"errors"
	"gameserver/internal/config"
	"gameserver/internal/db"
	"gameserver/internal/game/model"
	"kernel/kenum"
	"kernel/tools"
	"sync"
	"sync/atomic"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type ReadUserCallBack func(error, string, *model.UserData)
type PushReadUserFn func(*AsyncReadUserCb) error

type AyncReadUser struct {
	AccountId    string
	PacketId     uint32
	Ip           string
	GateId       uint32
	FsId         uint32
	SessionId    uint64
	Extra        string
	SdkChannelNo string
	Os           int
	Cb           ReadUserCallBack
}

type AsyncReadUserCb struct {
	Err      error
	Account  string
	UserData *model.UserData
	Cb       ReadUserCallBack
}

type AsncQueue struct {
	ctx      context.Context
	cancel   context.CancelFunc
	stopChan chan struct{}
	wg       sync.WaitGroup
	state    atomic.Uint32 // 状态 WorkState

	userQueue chan *AyncReadUser
	panicFunc func(string)

	pushBackFn PushReadUserFn
}

func (a *AsncQueue) start() error {
	if !a.state.CompareAndSwap(kenum.WorkState_Idle, kenum.WorkState_Running) {
		return errors.New("AyncQueue can't start, current state : " + kenum.StateToString(a.state.Load()))
	}

	log.Info("AyncQueuee start running")
	a.wg.Add(1)
	go tools.GoSafePost("AyncQueue run", func() {
		a.run()
	}, a.panicFunc)

	a.wg.Add(1)
	go tools.GoSafePost("AyncQueue run", func() {
		a.run()
	}, a.panicFunc)

	a.wg.Add(1)
	go tools.GoSafePost("AyncQueue monitor", func() {
		a.monitor()
	}, a.panicFunc)

	log.Info("AyncQueue start success")
	return nil
}

func (a *AsncQueue) monitor() {
	defer a.wg.Done()
	defer func() {
		if a.cancel != nil {
			a.cancel()
		}
	}()

	for {
		select {
		case <-a.stopChan:
			// 收到停止信号 处理剩余消息后退出
			log.Info("AyncQueue received stop signal, draining messages")
			return
		case <-a.ctx.Done():
			// ctx取消 立即退出
			log.Info("AyncQueue context canceled, exiting")
			return
		case <-time.After(time.Second * 3):
			log.Info("async load queue size", zap.Int("size", len(a.userQueue)))
		}
	}
}

func (a *AsncQueue) run() {
	defer a.wg.Done()
	defer func() {
		if a.cancel != nil {
			a.cancel()
		}
	}()

	for {
		select {
		case <-a.stopChan:
			// 收到停止信号 处理剩余消息后退出
			log.Info("AyncQueue received stop signal, draining messages")
			return
		case <-a.ctx.Done():
			// ctx取消 立即退出
			log.Info("AyncQueue context canceled, exiting")
			return
		case u, ok := <-a.userQueue:
			if !ok {
				log.Error("AyncQueue msg queue channel closed")
				return
			}
			log.Info("trace login async load ll", zap.String("accountId", u.AccountId))
			a.process(u)
		}
	}
}

func (a *AsncQueue) stop() error {
	timeOut := time.Second * 10
	if !a.state.CompareAndSwap(kenum.WorkState_Running, kenum.WorkState_Stopping) {
		return errors.New("AyncQueue can't stop, current state : " + kenum.StateToString(a.state.Load()))
	}

	log.Info("AyncQueue stopping")

	// 发送停止信号
	close(a.stopChan)

	// 等待goroutine退出
	stopped := make(chan struct{})
	go tools.GoSafePost("AyncQueue wait stop", func() {
		a.wg.Wait()
		close(stopped)
	}, a.panicFunc)

	select {
	case <-stopped:
		a.state.Store(kenum.WorkState_Stopped)
		log.Info("AyncQueue stopped")
		return nil
	case <-time.After(timeOut):
		if a.cancel != nil {
			a.cancel()
		}
		log.Warn("AyncQueue stop timeout, forcing shutdown", zap.Duration("timeout", timeOut))
		return errors.New("AyncQueue stop timeout after " + timeOut.String())
	}
}

func (a *AsncQueue) process(r *AyncReadUser) {
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	res := &AsyncReadUserCb{
		Account: r.AccountId,
		Cb:      r.Cb,
	}

	now := time.Now()
	user_data := &model.UserData{}
	db_res := db.LocalMongoReader.FindOne(
		ctx,
		config.Conf.LocalMongo.DB,
		db.CollectionName_User,
		bson.M{"accountid": r.AccountId},
	)

	res.Err = db_res.Decode(user_data)
	if res.Err == nil {
		res.UserData = user_data
	}

	cost := time.Since(now).Milliseconds()
	log.Info("trace login async load user cost", zap.Int64("cost", cost), zap.String("accountId", r.AccountId))
	a.pushBackFn(res)
}

func (a *AsncQueue) push(r *AyncReadUser) error {
	state := a.state.Load()
	if state != kenum.WorkState_Running {
		return errors.New("AyncQueue cannot push AyncQueue, state is : " + kenum.StateToString(state))
	}

	log.Info("trace login async load push start", zap.String("accountId", r.AccountId))
	select {
	case a.userQueue <- r:
		log.Info("trace login async load push end", zap.String("accountId", r.AccountId))
		return nil
	case <-a.stopChan:
		return errors.New("AyncQueue is stopping, AyncQueue reject")
	case <-a.ctx.Done():
		return errors.New("ctx cancelled, AyncQueue reject")
	default:
		return errors.New("AyncQueue AyncQueue queue is full, AyncQueue reject")
	}
}
