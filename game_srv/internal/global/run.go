package global

import (
	"fmt"
	"gameserver/internal/async"
	"gameserver/internal/tapping"
	"os/signal"
	"runtime"
	"syscall"

	"gameserver/internal/config"
	"gameserver/internal/content"
	"gameserver/internal/db"
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"gameserver/internal/io_out"
	"gameserver/internal/logic"
	"gameserver/internal/publicconst"

	"kernel/kenum"
	"kernel/nats"
	"kernel/tools"

	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"go.uber.org/zap"
)

func Run() {
	defer log.Sync()
	defer func() {
		log.Info("game_server closing...")
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			err := fmt.Errorf("%v: %s", r, buf[:l])
			log.Error("run panic", zap.Error(err))
			service.PostPanic("global.Run" + ":" + fmt.Sprintf("%v", r))
		}
	}()

	start()

	go tools.GoSafePost("signal down", func() {
		<-signalChan
		stop()
		close(exitChan)
	}, service.PostPanic)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}

func start() {
	start_register()
	start_discover()

	if err := content.Start(); err != nil {
		panic(fmt.Sprintf("content start err:%s", err))
	}
	if err := async.Start(); err != nil {
		panic(fmt.Sprintf("async start err:%s", err))
	}
	if err := tapping.Start(); err != nil {
		panic(fmt.Sprintf("tapping start err:%s", err))
	}
	if err := io_out.Start(); err != nil {
		panic(fmt.Sprintf("io start err:%s", err))
	}
	if err := logic.Start(); err != nil {
		panic(fmt.Sprintf("logic start err:%s", err))
	}
	sub_nats() // sub gate_way and fight and sdk
	log.Info("game_server started")
}

func start_register() {
	// addr := config.Conf.Ws.LinkAddr
	// if addr == "" {
	// 	addr = config.Conf.Ws.Addr
	// }

	info := fmt.Sprintf("%v", config.Conf.ServerId)
	if err := game_register.Register(int64(publicconst.MAX_SERVER_TTL), info, service.PostPanic); err != nil {
		panic(fmt.Sprintf("gs register err:%s", err))
	}
	log.Info("regist to discover", zap.String("info", info))
}

func start_discover() {
	fight_discover.Start(kenum.SER_FIGHT_PREFIX, service.PostPanic)
	log.Info("start_discover succ")
}

func sub_nats() {
	// sub gate
	sub_gate := fmt.Sprintf("gate.game.%v", config.Conf.ServerId)
	if err := nats.Subscribe(sub_gate, logic.DispatchGateMsg); err != nil {
		panic(fmt.Sprintf("nats subscribe gate subscribe:%s, err:%s", sub_gate, err))
	}
	log.Info("nats subscribe", zap.String("sub name", sub_gate))

	// sub admin boradcast
	admin_boradcast_game := "admin.game"
	if err := nats.Subscribe(admin_boradcast_game, logic.DispatchAdminMsg); err != nil {
		panic(fmt.Sprintf("nats subscribe admin subscribe:%s, err:%s", admin_boradcast_game, err))
	}
	log.Info("nats subscribe", zap.String("sub name", admin_boradcast_game))

	// sub admin
	admin_game := fmt.Sprintf("admin.game.%v", config.Conf.ServerId)
	if err := nats.Subscribe(admin_game, logic.DispatchAdminMsg); err != nil {
		panic(fmt.Sprintf("nats subscribe admin subscribe:%s, err:%s", admin_game, err))
	}
	log.Info("nats subscribe", zap.String("sub name", admin_game))

	// sub sdk
	sub_sdk := fmt.Sprintf("sdk.game.%v", config.Conf.ServerId)
	if err := nats.Subscribe(sub_sdk, logic.DispatchSDKMsg); err != nil {
		panic(fmt.Sprintf("nats subscribe sdk subscribe:%s, err:%s", sub_sdk, err))
	}
	log.Info("nats subscribe", zap.String("sub name", sub_sdk))

	// sub fight
	sub_fight := fmt.Sprintf("fight.game.%v", config.Conf.ServerId)
	if err := nats.Subscribe(sub_fight, logic.DispatchFightMsg); err != nil {
		panic(fmt.Sprintf("nats subscribe fight subscribe:%s, err:%s", sub_fight, err))
	}
	log.Info("nats subscribe", zap.String("sub name", sub_fight))
}

func stop() {
	log.Info("Server Stop")

	if err := content.Stop(); err != nil {
		log.Error("content stop error", zap.Error(err))
	}
	if err := async.Stop(); err != nil {
		log.Error("async stop error", zap.Error(err))
	}
	if err := tapping.Stop(); err != nil {
		log.Error("tapping stop error", zap.Error(err))
	}
	if err := logic.Stop(); err != nil {
		log.Error("logic stop error", zap.Error(err))
	}
	if err := io_out.Stop(); err != nil {
		log.Error("io stop error", zap.Error(err))
	}

	nats.Stop()
	player.Stop()
	db.Stop()
	db.StopGlobal()
	rdb_single.Stop()

	log.Info("Server Stopped")
}
