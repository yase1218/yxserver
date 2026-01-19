package logic

import (
	"context"
	"errors"
	"fmt"
	"gameserver/internal/async"
	"gameserver/internal/content"
	"gameserver/internal/game/handle"
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"kernel/kenum"
	"kernel/metric"
	"kernel/tools"
	"msg"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 消息管道类型
type MsgChannel = int

const (
	MsgChannel_Gate MsgChannel = iota
	MsgChannel_Sdk
	MsgChannel_Fight
	MsgChannel_Out

	MsgChannel_Max
)

func channelTypeToString(channelType MsgChannel) string {
	switch channelType {
	case MsgChannel_Gate:
		return "gate"
	case MsgChannel_Sdk:
		return "sdk"
	case MsgChannel_Fight:
		return "fight"
	case MsgChannel_Out:
		return "out"
	default:
		return "unknown"
	}
}

type OutMsg struct {
	Subject string
	Msg     proto.Message
}

type GameLogic struct {
	ctx       context.Context
	cancel    context.CancelFunc
	stopChan  chan struct{}
	wg        sync.WaitGroup
	state     atomic.Uint32              // 状态 WorkState
	metrics   []*metric.ProcessorMetrics // 指标收集
	gateMsg   chan *msg.GateToGame
	sdkMsg    chan *msg.RequestCommonInterMsg
	fightMsg  chan *msg.FightToGame
	contentCb chan *content.ContentCb
	asyncCb   chan *async.AsyncReadUserCb
	tickCount uint64

	lastMsg *msg.GateToGame
}

func (g *GameLogic) metrcisInfo() map[string]interface{} {
	infos := make(map[string]interface{})
	for i := 0; i < MsgChannel_Max; i++ {
		infos[channelTypeToString(i)] = g.metrics[i].Info()
	}
	return infos
}

type OrderInfo struct {
	Uid      uint64
	OrderNum string
	// ...
}

// 开服加载构造两个map 运行时同步更新两个map 更新操作异步入库 订单信息永久保留
type OrderManager struct {
	byOrderNum map[string]*OrderInfo
	byUid      map[uint64]*OrderInfo
}

func (g *GameLogic) start() error {
	if !g.state.CompareAndSwap(kenum.WorkState_Idle, kenum.WorkState_Running) {
		return errors.New("game logic can't start, current state : " + kenum.StateToString(g.state.Load()))
	}

	log.Info("game logic start running")
	g.wg.Add(1)
	go tools.GoSafePost("game logic run", func() {
		g.run()
	}, service.PostPanic)

	g.wg.Add(1)
	go tools.GoSafePost("game logic monitor", func() {
		g.monitor()
	}, service.PostPanic)

	log.Info("game logic start success")
	return nil
}

func (g *GameLogic) monitor() {
	defer g.wg.Done()
	defer func() {
		if g.cancel != nil {
			g.cancel()
		}
	}()

	for {
		select {
		case <-g.stopChan:
			log.Info("game logic received stop signal, draining messages")
			g.drainMessages()
			return
		case <-g.ctx.Done():
			log.Info("game logic context canceled, exiting")
			return
		case <-time.After(time.Second * 3):
			log.Info("game logic queue size ", zap.Int("gate queue size", len(g.gateMsg)), zap.Int("db queue size", len(g.asyncCb)))
		}
	}
}

func (g *GameLogic) run() {
	defer g.wg.Done()
	defer func() {
		if g.cancel != nil {
			g.cancel()
		}
	}()

	log.Info("game logic main loop started")

	secondTicker := time.NewTicker(time.Second * 1)
	defer secondTicker.Stop()

	battlePassTimer := time.NewTimer(time.Duration(tools.GetWeeklyRefreshTime(0)-tools.GetCurTime()) * time.Second) // 海边派对周环境刷新
	defer battlePassTimer.Stop()

	// for {
	// 	// 第一优先级：退出信号和 gate 消息
	// 	select {
	// 	case <-g.stopChan:
	// 		log.Info("game logic received stop signal, draining messages")
	// 		g.drainMessages()
	// 		return
	// 	case <-g.ctx.Done():
	// 		log.Info("game logic context canceled, exiting")
	// 		return
	// 	case m, ok := <-g.gateMsg:
	// 		if !ok {
	// 			return
	// 		}
	// 		g.process_gate_msg(m)

	// 		continue // 处理完 gate 消息后继续检查是否有更多高优先级消息
	// 	// case cb, ok := <-g.asyncCb:
	// 	// 	if !ok {
	// 	// 		log.Error("game logic async message channel closed")
	// 	// 		return
	// 	// 	}
	// 	// 	if cb != nil && cb.Cb != nil {
	// 	// 		cb.Cb(cb.Err, cb.Account, cb.UserData)
	// 	// 	}
	// 	// 	continue
	// 	default:
	// 		// 没有高优先级消息，进入普通消息处理
	// 	}

	// 	// 第二优先级 db
	// 	select {
	// 	case <-g.stopChan:
	// 		log.Info("game logic received stop signal, draining messages")
	// 		g.drainMessages()
	// 		return
	// 	case <-g.ctx.Done():
	// 		log.Info("game logic context canceled, exiting")
	// 		return
	// 	case cb, ok := <-g.asyncCb:
	// 		if !ok {
	// 			log.Error("game logic async message channel closed")
	// 			return
	// 		}
	// 		if cb != nil && cb.Cb != nil {
	// 			cb.Cb(cb.Err, cb.Account, cb.UserData)
	// 		}
	// 		continue
	// 	default:
	// 	}

	// 	// 第三优先级：其他消息
	// 	select {
	// 	case <-g.stopChan:
	// 		log.Info("game logic received stop signal, draining messages")
	// 		g.drainMessages()
	// 		return
	// 	case <-g.ctx.Done():
	// 		log.Info("game logic context canceled, exiting")
	// 		return
	// 	// case m, ok := <-g.gateMsg:
	// 	// 	if !ok {
	// 	// 		return
	// 	// 	}
	// 	// 	g.process_gate_msg(m)
	// 	case <-secondTicker.C:
	// 		g.on_second_ticker(time.Now())
	// 	case cb, ok := <-g.contentCb:
	// 		if !ok {
	// 			log.Error("game logic content message channel closed")
	// 			return
	// 		}
	// 		if cb != nil && cb.Cb != nil {
	// 			cb.Cb(cb.Ec)
	// 		}
	// 	// case cb, ok := <-g.asyncCb:
	// 	// 	if !ok {
	// 	// 		log.Error("game logic async message channel closed")
	// 	// 		return
	// 	// 	}
	// 	// 	if cb != nil && cb.Cb != nil {
	// 	// 		cb.Cb(cb.Err, cb.Account, cb.UserData)
	// 	// 	}
	// 	case m, ok := <-g.sdkMsg:
	// 		if !ok {
	// 			log.Error("game logic sdk message channel closed")
	// 			return
	// 		}
	// 		g.process_sdk_msg(m)
	// 	case m, ok := <-g.fightMsg:
	// 		if !ok {
	// 			log.Error("game logic fight message channel closed")
	// 			return
	// 		}
	// 		g.process_fight_msg(m)
	// 	case <-battlePassTimer.C:
	// 		g.refresh_battle_pass_week_fields(battlePassTimer)
	// 	}
	// }
	idleDelay := time.NewTimer(5 * time.Millisecond)
	defer idleDelay.Stop()
	for {
		idleDelay.Reset(5 * time.Millisecond)
		select {
		case <-g.stopChan:
			// 收到停止信号 处理剩余消息后退出
			log.Info("game logic received stop signal, draining messages")
			g.drainMessages()
			return
		case <-g.ctx.Done():
			// ctx取消 立即退出
			log.Info("game logic context canceled, exiting")
			return
		// case m, ok := <-g.gateMsg:
		// 	g.lastMsg = m
		// 	if !ok {
		// 		log.Error("game logic gate message channel closed")
		// 		return
		// 	}
		// 	g.process_gate_msg(m)
		case cb, ok := <-g.asyncCb:
			if !ok {
				log.Error("game logic async message channel closed")
				return
			}
			if cb != nil {
				if cb.Cb != nil {
					cb.Cb(cb.Err, cb.Account, cb.UserData)
				}
			}
		default:
			select {
			case <-g.stopChan:
				// 收到停止信号 处理剩余消息后退出
				log.Info("game logic received stop signal, draining messages")
				g.drainMessages()
				return
			case <-g.ctx.Done():
				// ctx取消 立即退出
				log.Info("game logic context canceled, exiting")
				return
			case <-secondTicker.C:
				// 秒级定时器
				g.on_second_ticker(time.Now())
			case cb, ok := <-g.contentCb:
				if !ok {
					log.Error("game logic content message channel closed")
					return
				}
				if cb != nil {
					if cb.Cb != nil {
						cb.Cb(cb.Ec)
					}
				}
			case m, ok := <-g.gateMsg:
				g.lastMsg = m
				if !ok {
					log.Error("game logic gate message channel closed")
					return
				}
				g.process_gate_msg(m)
			// case cb, ok := <-g.asyncCb:
			// 	if !ok {
			// 		log.Error("game logic async message channel closed")
			// 		return
			// 	}
			// 	if cb != nil {
			// 		if cb.Cb != nil {
			// 			cb.Cb(cb.Err, cb.Account, cb.UserData)
			// 		}
			// 	}
			case m, ok := <-g.sdkMsg:
				if !ok {
					log.Error("game logic sdk message channel closed")
					return
				}
				g.process_sdk_msg(m)
			case m, ok := <-g.fightMsg:
				if !ok {
					log.Error("game logic fight message channel closed")
					return
				}
				g.process_fight_msg(m)
			case <-battlePassTimer.C:
				// 海边派对黄历周环境刷新
				g.refresh_battle_pass_week_fields(battlePassTimer)
			case <-idleDelay.C: // 5毫秒结束重来
			}
		}
	}
}

func (g *GameLogic) stop() error {
	timeOut := time.Second * 10
	if !g.state.CompareAndSwap(kenum.WorkState_Running, kenum.WorkState_Stopping) {
		return errors.New("game logic can't stop, current state : " + kenum.StateToString(g.state.Load()))
	}

	log.Info("game logic stopping")

	// 发送停止信号
	close(g.stopChan)

	// 等待goroutine退出
	stopped := make(chan struct{})
	go tools.GoSafePost("game logic wait stop", func() {
		g.wg.Wait()
		close(stopped)
	}, service.PostPanic)

	select {
	case <-stopped:
		g.state.Store(kenum.WorkState_Stopped)
		log.Info("game logic stopped", zap.Any("metrics", g.metrcisInfo()))
		return nil
	case <-time.After(timeOut):
		if g.cancel != nil {
			g.cancel()
		}
		log.Warn("game logic stop timeout, forcing shutdown", zap.Duration("timeout", timeOut))
		return errors.New("stop timeout after " + timeOut.String())
	}
}

func (g *GameLogic) drainMessages() {
	log.Info("game logic processing remaining messages during shutdown")
	const maxAttemps = 10
	const internal = 100 * time.Millisecond

	for i := 0; i < maxAttemps; i++ {
		drained := false
		for {
			select {
			case m := <-g.gateMsg:
				g.process_gate_msg(m)
				drained = true
			case m := <-g.sdkMsg:
				g.process_sdk_msg(m)
				drained = true
			case m := <-g.fightMsg:
				g.process_fight_msg(m)
				drained = true
			default:
				if !drained {
					// 本轮没有处理消息，认为已处理完毕
					log.Info("game logic drained all messages", zap.Int("attemp", i+1),
						zap.Int("cur gate msg lenth", len(g.gateMsg)),
						zap.Int("cur sdk msg lenth", len(g.sdkMsg)),
						zap.Int("cur fight msg lenth", len(g.fightMsg)))
					return
				}
				// 下一轮尝试
				goto nextRound
			}
		nextRound:
			time.Sleep(internal)
		}
	}
	log.Warn("game logic may have remaining messages after max attempts",
		zap.Int("cur gate msg lenth", len(g.gateMsg)),
		zap.Int("cur sdk msg lenth", len(g.sdkMsg)),
		zap.Int("cur fight msg lenth", len(g.fightMsg)))
}

func (g *GameLogic) pushGate(m *msg.GateToGame) error {
	g.metrics[MsgChannel_Gate].PushAttempts.Add(1)

	state := g.state.Load()
	if state != kenum.WorkState_Running {
		g.metrics[MsgChannel_Gate].PushRejects.Add(1)
		return errors.New("game logic cannot push gate msg, state is : " + kenum.StateToString(state))
	}

	select {
	case g.gateMsg <- m:
		g.metrics[MsgChannel_Gate].PushSuccess.Add(1)
		traceFunc := func(mid, gid uint32, sid uint64, account string) {
			msgName := ""
			switch mid {
			case msg.MsgID_RequestLoginId:
				msgName = "RequestLogin"
			case msg.MsgID_RandomPlayerNameReqId:
				msgName = "RandomPlayerNameReq"
			case msg.MsgID_InitPlayerNameAndShipReqId:
				msgName = "InitPlayerNameAndShipReq"
			default:
				return
			}
			log.Info("trace msg push to queue now",
				zap.Uint32("gate id", gid),
				zap.Uint32("msg", mid),
				zap.String("msgName", msgName),
				zap.Uint64("session id", sid),
				zap.String("account", account),
			)
		}

		traceFunc(m.MsgId, m.GateId, m.Session, m.AccountId)
		return nil
	case <-g.stopChan:
		g.metrics[MsgChannel_Gate].PushRejects.Add(1)
		return errors.New("game logic is stopping, gate msg reject")
	case <-g.ctx.Done():
		g.metrics[MsgChannel_Gate].PushRejects.Add(1)
		return errors.New("ctx cancelled, gate msg reject")
	default:
		g.metrics[MsgChannel_Gate].PushRejects.Add(1)
		return errors.New("game logic gate msg queue is full, gate msg reject")
	}
}

func (g *GameLogic) pushContencCb(cb *content.ContentCb) error {
	state := g.state.Load()
	if state != kenum.WorkState_Running {
		return errors.New("game logic cannot push content cb, state is : " + kenum.StateToString(state))
	}

	select {
	case g.contentCb <- cb:
		return nil
	case <-g.stopChan:
		return errors.New("game logic is stopping, content cb reject")
	case <-g.ctx.Done():
		return errors.New("ctx cancelled, content cb reject")
	default:
		return errors.New("game logic content cb queue is full, content cb reject")
	}
}

func (g *GameLogic) pushAsyncCb(cb *async.AsyncReadUserCb) error {
	state := g.state.Load()
	if state != kenum.WorkState_Running {
		return errors.New("game logic cannot push async cb, state is : " + kenum.StateToString(state))
	}

	select {
	case g.asyncCb <- cb:
		return nil
	case <-g.stopChan:
		return errors.New("game logic is stopping, async cb reject")
	case <-g.ctx.Done():
		return errors.New("ctx cancelled, async cb reject")
	default:
		return errors.New("game logic async cb queue is full, async cb reject")
	}
}

func (g *GameLogic) pushSDK(m *msg.RequestCommonInterMsg) error {
	g.metrics[MsgChannel_Sdk].PushAttempts.Add(1)

	state := g.state.Load()
	if state != kenum.WorkState_Running {
		g.metrics[MsgChannel_Sdk].PushRejects.Add(1)
		return errors.New("game logic cannot push sdk msg, state is : " + kenum.StateToString(state))
	}

	select {
	case g.sdkMsg <- m:
		g.metrics[MsgChannel_Sdk].PushSuccess.Add(1)
		return nil
	case <-g.stopChan:
		g.metrics[MsgChannel_Sdk].PushRejects.Add(1)
		return errors.New("game logic is stopping, sdk msg reject")
	case <-g.ctx.Done():
		g.metrics[MsgChannel_Sdk].PushRejects.Add(1)
		return errors.New("ctx cancelled, sdk msg reject")
	default:
		g.metrics[MsgChannel_Sdk].PushRejects.Add(1)
		return errors.New("game logic sdk msg queue is full, sdk msg reject")
	}
}

func (g *GameLogic) pushFight(m *msg.FightToGame) error {
	g.metrics[MsgChannel_Fight].PushAttempts.Add(1)

	state := g.state.Load()
	if state != kenum.WorkState_Running {
		g.metrics[MsgChannel_Fight].PushRejects.Add(1)
		return errors.New("game logic cannot push fight msg, state is : " + kenum.StateToString(state))
	}

	select {
	case g.fightMsg <- m:
		g.metrics[MsgChannel_Fight].PushSuccess.Add(1)
		return nil
	case <-g.stopChan:
		g.metrics[MsgChannel_Fight].PushRejects.Add(1)
		return errors.New("game logic is stopping, fight msg reject")
	case <-g.ctx.Done():
		g.metrics[MsgChannel_Fight].PushRejects.Add(1)
		return errors.New("ctx cancelled, fight msg reject")
	default:
		g.metrics[MsgChannel_Fight].PushRejects.Add(1)
		return errors.New("game logic fight msg queue is full, fight msg reject")
	}
}

func (g *GameLogic) pushOut(m *OutMsg) error {
	g.metrics[MsgChannel_Out].PushAttempts.Add(1)

	state := g.state.Load()
	if state != kenum.WorkState_Running {
		g.metrics[MsgChannel_Out].PushRejects.Add(1)
		return errors.New("game logic cannot push out msg, state is : " + kenum.StateToString(state))
	}

	select {
	case <-g.stopChan:
		g.metrics[MsgChannel_Out].PushRejects.Add(1)
		return errors.New("game logic is stopping, out msg reject")
	case <-g.ctx.Done():
		g.metrics[MsgChannel_Out].PushRejects.Add(1)
		return errors.New("ctx cancelled, out msg reject")
	default:
		g.metrics[MsgChannel_Out].PushRejects.Add(1)
		return errors.New("game logic out msg queue is full, out msg reject")
	}

}
func (g *GameLogic) process_gate_msg(m *msg.GateToGame) {
	traceFunc := func(mid uint32, sid uint64, account, content string) {
		msgName := ""
		switch mid {
		case msg.MsgID_RequestLoginId:
			msgName = "RequestLogin"
		case msg.MsgID_RandomPlayerNameReqId:
			msgName = "RandomPlayerNameReq"
		case msg.MsgID_InitPlayerNameAndShipReqId:
			msgName = "InitPlayerNameAndShipReq async"
		default:
			return
		}
		log.Info("trace msg process "+content,
			zap.Uint32("msg", mid),
			zap.String("msgName", msgName),
			zap.Uint64("session id", sid),
			zap.String("account", account),
		)
	}
	traceFunc(m.MsgId, m.Session, m.AccountId, "start")
	g.metrics[MsgChannel_Gate].ProcessCount.Add(1)
	defer func() {
		if r := recover(); r != nil {
			g.metrics[MsgChannel_Gate].ProcessErrors.Add(1)
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			err := fmt.Errorf("%v: %s", r, buf[:l])
			log.Error("game logic gate msg handler panic",
				zap.Error(err),
				zap.Any("message", m))
		}
	}()

	accountId := m.AccountId
	pid := m.PacketId
	msgId := m.MsgId
	data := m.Data
	sessionId := m.Session
	p := player.FindByAccount(accountId)

	now := time.Now()

	defer func() {
		cost := time.Since(now).Milliseconds()
		if cost > 50 {
			log.Debug("process_gate_msg too slow",
				zap.String("account id", accountId),
				zap.Uint32("msg id", msgId),
				zap.Int64("cost millisecond", cost),
			)
		}
	}()

	var (
		er       error
		user_msg proto.Message
	)
	if msgId >= 9000 || msgId == uint32(msg.MsgId_ID_RequestLogout) {
		_, _, user_msg, er = msg.PbProcessor.Unmarshal(data)
	} else {
		_, _, user_msg, er = msg.PbProcessor.UnmarshalUnlen(data)
	}
	if er != nil {
		log.Error("process_gate_msg Unmarshl failed",
			zap.String("account id", accountId),
			zap.Uint32("package id", pid),
			zap.Uint32("msg id", msgId),
			zap.Error(er))
		traceFunc(m.MsgId, m.Session, m.AccountId, "unmarshal failed")
		return
	}

	if p == nil { // 账号未登录 第一条必须是登录消息
		if msgId != uint32(msg.MsgId_ID_RequestLogin) {
			log.Error("process_gate_msg player not found, maybe logout",
				zap.String("account id", accountId),
				zap.Uint32("package id", pid),
				zap.Uint32("msg id", msgId),
			)
			traceFunc(m.MsgId, m.Session, m.AccountId, "player not login")
			return
		} else {
			//  login
			req := user_msg.(*msg.RequestLogin)
			login(accountId, pid, m.Ip, m.GateId, sessionId, req, 0)
			traceFunc(m.MsgId, m.Session, m.AccountId, "succ")
			return
		}
	} else { // 账号已登录
		if msgId == uint32(msg.MsgId_ID_RequestLogin) { // 重登录 踢人
			log.Info("player already login",
				zap.String("account id", accountId),
			)

			service.KickOutPlayer(p)

			//  login
			req := user_msg.(*msg.RequestLogin)
			login(accountId, pid, m.Ip, m.GateId, sessionId, req, p.GetFsId())
			traceFunc(m.MsgId, m.Session, m.AccountId, "succ")
			return
		}
	}
	handle.HandleGateMsg(msgId, pid, user_msg, p)
	traceFunc(m.MsgId, m.Session, m.AccountId, "succ")
	//log.Debug("Gate =============>>>> Game", zap.Uint32("msg id:", msg_id))
}

func login(accountId string, packetId uint32, ip string, gateId uint32, sessionId uint64, req *msg.RequestLogin, fsId uint32) {
	service.OnLogin(accountId, gateId, fsId, packetId, sessionId, ip, req)
	// now := time.Now()
	// p := player.CreatePlayer(now)
	// p.BindGateId(gateId)
	// p.BindSessionId(sessionId)

	// if fsId != 0 {
	// 	p.SetFsId(fsId)
	// }

	// {
	// 	res := &msg.ResponseLogin{
	// 		Result: msg.ErrCode_SUCC,
	// 		Info: &msg.AccountInfo{
	// 			Uid: 1,
	// 		},
	// 	}

	// 	p.SendResponse(packetId, res, res.Result)
	// 	return p
	// }

	// 登录加载
	// if err := service.OnLogin(p, accountId, packetId, ip, req); err != nil {
	// 	log.Error("login failed",
	// 		zap.String("account_id", accountId), zap.Error(err))
	// 	player.DelPlayer(p)
	// 	return nil
	// } else {
	// 	log.Info("login succ", service.ZapUser(p))
	// 	service.ResLogin(packetId, now, p)
	// 	service.AfterLogin(p)
	// 	player.AddPlayer(p)
	// }
	//return p
}

func (g *GameLogic) process_sdk_msg(m *msg.RequestCommonInterMsg) {
	g.metrics[MsgChannel_Sdk].ProcessCount.Add(1)
	defer func() {
		if r := recover(); r != nil {
			g.metrics[MsgChannel_Sdk].ProcessErrors.Add(1)
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			err := fmt.Errorf("%v: %s", r, buf[:l])
			log.Error("game logic sdk msg handler panic",
				zap.Error(err),
				zap.Any("message", m))
		}
	}()
	handle.RouteInterMsg(m)
}

func (g *GameLogic) process_fight_msg(m *msg.FightToGame) {
	g.metrics[MsgChannel_Fight].ProcessCount.Add(1)
	defer func() {
		if r := recover(); r != nil {
			g.metrics[MsgChannel_Fight].ProcessErrors.Add(1)
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			err := fmt.Errorf("%v: %s", r, buf[:l])
			log.Error("game logic fight msg handler panic",
				zap.Error(err),
				zap.Any("message", m))
		}
	}()
	handle.HandleFightMsg(m.GetMsgId(), m)
}

func (g *GameLogic) on_second_ticker(now time.Time) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			err := fmt.Errorf("%v: %s", r, buf[:l])
			log.Error("game logic second ticker panic",
				zap.Error(err))
		}
	}()

	//if !config.Conf.IsDebug() {
	if g.tickCount%10 == 0 {
		log.Info("~~~~~~~~~~~~~~~~~~~~~~~~game logic metrics",
			zap.Any("metrics", g.metrcisInfo()),
		)
		log.Info("~~~~~~~~~~~~~~~~~~~~~~~~current online player", zap.Int("num", int(player.OnlineNum())))

		log.Info("~~~~~~~~~~~~~~~~~~~~~~~~login logout cost", zap.Any("cost", service.CostRecordData))
	}
	//}
	service.OnSecondTicker(now)
	g.tickCount++
}

func (g *GameLogic) refresh_battle_pass_week_fields(timer *time.Timer) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			err := fmt.Errorf("%v: %s", r, buf[:l])
			log.Error("game logic battle pass timer panic",
				zap.Error(err))
		}
	}()
	service.BattlePassRefreshTime()
	timer.Reset(time.Duration(tools.GetWeeklyRefreshTime(0)-tools.GetCurTime()) * time.Second)
}
