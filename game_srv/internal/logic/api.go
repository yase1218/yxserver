package logic

import (
	"context"
	"errors"
	"gameserver/internal/async"
	"gameserver/internal/content"
	"gameserver/internal/game/handle"
	"kernel/kenum"
	"kernel/metric"
	"kernel/tools"
	"msg"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	instance *GameLogic
	once     sync.Once
)

func Init(args ...int) {
	createInstance(args...)
}

func createInstance(args ...int) {
	once.Do(func() {
		if len(args) != MsgChannel_Max {
			log.Panic("game logic start param lenth err",
				zap.Int("need", MsgChannel_Max), zap.Int("lenth", len(args)))
		}
		ctx, cacel := context.WithCancel(context.Background())
		instance = &GameLogic{
			ctx:       ctx,
			cancel:    cacel,
			stopChan:  make(chan struct{}),
			metrics:   make([]*metric.ProcessorMetrics, MsgChannel_Max),
			gateMsg:   make(chan *msg.GateToGame, 1024*args[MsgChannel_Gate]),
			sdkMsg:    make(chan *msg.RequestCommonInterMsg, 1024*args[MsgChannel_Sdk]),
			fightMsg:  make(chan *msg.FightToGame, 1024*args[MsgChannel_Fight]),
			contentCb: make(chan *content.ContentCb, 1024*10),
			asyncCb:   make(chan *async.AsyncReadUserCb, 1024*10),
		}

		for i := 0; i < MsgChannel_Max; i++ {
			instance.metrics[i] = &metric.ProcessorMetrics{}
		}

		instance.state.Store(kenum.WorkState_Idle)
	})
}

func Start() error {
	if instance == nil {
		return errors.New("game logic not initialized, call Init first")
	}

	return instance.start()
}

func Stop() error {
	if instance == nil {
		return errors.New("game logic not initialized, call Init first")
	}
	return instance.stop()
}

func PushGate(m *msg.GateToGame) error {
	if instance == nil {
		return errors.New("game logic not initialized, call Init first")
	}
	return instance.pushGate(m)
}

func PushSDK(m *msg.RequestCommonInterMsg) error {
	if instance == nil {
		return errors.New("game logic not initialized, call Init first")
	}
	return instance.pushSDK(m)
}

func PushFight(m *msg.FightToGame) error {
	if instance == nil {
		return errors.New("game logic not initialized, call Init first")
	}
	return instance.pushFight(m)
}

func PushContencCb(cb *content.ContentCb) error {
	if instance == nil {
		return errors.New("game logic not initialized, call Init first")
	}
	return instance.pushContencCb(cb)
}

func PushAsyncCb(cb *async.AsyncReadUserCb) error {
	if instance == nil {
		return errors.New("game logic not initialized, call Init first")
	}
	return instance.pushAsyncCb(cb)
}

func DispatchGateMsg(req *nats.Msg) {
	gate_msg := &msg.GateToGame{}
	if err := proto.Unmarshal(req.Data, gate_msg); err != nil {
		log.Error("DisptchGateMsg", zap.String("subject", req.Subject))
	} else {
		traceFunc := func(mid uint32, sid uint64, account, content string) {
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
			log.Info("trace msg recv from gate "+content,
				zap.Uint32("msg", mid),
				zap.String("msgName", msgName),
				zap.Uint64("session id", sid),
				zap.String("account", account),
			)
		}
		traceFunc(gate_msg.MsgId, gate_msg.Session, gate_msg.AccountId, "push to queue")
		if err := PushGate(gate_msg); err != nil {
			log.Error("PushGate failed", zap.Uint32("msg id", gate_msg.MsgId), zap.Error(err))
			traceFunc(gate_msg.MsgId, gate_msg.Session, gate_msg.AccountId, "push to queue failed")
		} else {
			traceFunc(gate_msg.MsgId, gate_msg.Session, gate_msg.AccountId, "push to queue success")
		}
	}
}

func DispatchSDKMsg(req *nats.Msg) {
	sdk_msg := &msg.RequestCommonInterMsg{}
	if err := proto.Unmarshal(req.Data, sdk_msg); err != nil {
		log.Error("DispatchSDKMsg", zap.String("subject", req.Subject))
	} else {
		PushSDK(sdk_msg)
	}
}

func DispatchFightMsg(req *nats.Msg) {
	fight_msg := &msg.FightToGame{}
	if err := proto.Unmarshal(req.Data, fight_msg); err != nil {
		log.Error("DisptchFightMsg", zap.String("subject", req.Subject))
	} else {
		PushFight(fight_msg)
	}
}

func DispatchAdminMsg(req *nats.Msg) {
	log.Info("DispatchAdminMsg", zap.String("subject", req.Subject), zap.Any("req", req))
	admin_msg := &msg.RequestCommonInterMsg{}
	if err := proto.Unmarshal(req.Data, admin_msg); err != nil {
		log.Error("DisptchGateMsg", zap.String("subject", req.Subject))
	} else {
		if admin_msg.MsgId >= 100000 {
			go tools.GoSafe("admin async trigger", func() {
				resMsg := handle.RouteAdminMsg(admin_msg)
				if resMsg == nil {
					log.Error("RouteAdminMsg returned nil", zap.Any("admin_msg", admin_msg))
					return
				}

				data, err := proto.Marshal(resMsg)
				if err != nil {
					log.Error("Failed to marshal response", zap.Error(err))
					return
				}

				if err := req.Respond(data); err != nil {
					log.Error("Failed to send response", zap.Error(err))
				}
				log.Info("respond admin msg", zap.String("subject", req.Subject), zap.Any("msg", admin_msg))
			})
		} else {
			PushSDK(admin_msg)
		}
		log.Info("handle admin msg", zap.String("subject", req.Subject), zap.Any("msg", admin_msg))
	}
}
