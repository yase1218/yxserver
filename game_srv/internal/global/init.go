package global

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"kernel/data/db_global"
	"kernel/discover"
	"kernel/kenum"
	"kernel/nats"
	"kernel/rbi"
	"kernel/tda"
	"kernel/tools"

	"gameserver/internal/async"
	common2 "gameserver/internal/common"
	"gameserver/internal/config"
	"gameserver/internal/content"
	"gameserver/internal/db"
	"gameserver/internal/fight"
	"gameserver/internal/game/common"
	"gameserver/internal/game/service"
	"gameserver/internal/game/uid"
	"gameserver/internal/io_out"
	"gameserver/internal/logic"
	"gameserver/internal/tapping"

	"github.com/v587-zyf/gc/db/mongo"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"github.com/v587-zyf/gc/utils"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

var (
	game_register  *discover.ServiceRegister
	admin_register *discover.ServiceRegister
	fight_discover *discover.ServerDiscover
)

func Init(serverId uint32, serverName, localDb string) {
	rand.NewSource(time.Now().UnixNano())

	ctx := context.Background()
	init_log(ctx, serverId)
	init_conf(ctx, serverId, serverName, localDb)
	init_tda(ctx)
	init_db(ctx)
	init_redis(ctx)
	init_snow()
	init_nats(ctx)
	init_cfg_template(ctx)
	init_prof()
	init_rbi()
	init_register(ctx)
	init_discover(ctx)

	// 暂时保留 需要改为和中心服通信 start
	init_db_global(ctx)
	init_server_data()
	// 暂时保留 需要改为和中心服通信 end

	// 暂时保留grpc
	init_grpc(ctx)

	init_service(ctx)

	init_content(ctx)
	init_async(ctx)
	init_tapping(ctx)
	init_out(ctx)
	// last init
	init_logic(ctx)
}

func init_log(ctx context.Context, serverId uint32) {
	// if err := log.Init(ctx, log.WithSerName("game"), log.WithSkipCaller(2)); err != nil {
	// 	panic(fmt.Sprintf("log init err:%s", err.Error()))
	// }
	log.InitDevelopment("game", "game:"+strconv.Itoa(int(serverId)))
}

func init_conf(ctx context.Context, serverId uint32, serverName, localDb string) {
	config.Conf = new(config.Config)
	if err := utils.Load(config.Conf, ""); err != nil {
		panic(fmt.Sprintf("load config file failed:%s", err.Error()))
	}
	if serverId > 0 {
		config.Conf.ServerId = serverId
	}
	if serverName != "" {
		config.Conf.ServerName = serverName
	}
	if localDb != "" {
		config.Conf.LocalMongo.DB = localDb
	}
	log.Info("config load succ", zap.String("version", "1:50:50"), zap.Any("conf", config.Conf))
}

func init_tda(ctx context.Context) {
	if err := tda.Init(); err != nil {
		panic(fmt.Sprintf("tda init err:%s", err.Error()))
	}
	log.Info("tda init succ")
}

func init_db(ctx context.Context) {
	db.Init(service.PostPanic)
	db.InitGlobal(service.PostPanic)
	log.Info("db init succ")
}

func init_redis(ctx context.Context) {
	if err := rdb_single.InitSingle(
		ctx,
		rdb_single.WithAddr(config.Conf.Redis.Addr),
		rdb_single.WithPwd(config.Conf.Redis.Pass),
	); err != nil {
		panic(fmt.Sprintf("redis init err:%s", err.Error()))
	}
	log.Info("redis init succ")
}

func init_nats(ctx context.Context) {
	if err := nats.InitService(config.Conf.Nats.Addr); err != nil {
		panic(fmt.Sprintf("nats init addr:%s, err:%s", config.Conf.Nats.Addr, err.Error()))
	}
	log.Info("nats init succ")
}

func init_cfg_template(ctx context.Context) {
	template.InitTemplate(config.Conf.CsvPath)
	log.Info("template init succ")
}

func init_prof() {
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪，block
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪，mutex

	go tools.GoSafePost("prof init", func() {
		http.ListenAndServe(":6060", nil)
	}, service.PostPanic)
	log.Info("prof init succ")
}

func init_rbi() {
	rbi.SetOpen(config.Conf.Rbi.Open)
	if err := rbi.Init(config.Conf.Rbi.Url, config.Conf.Rbi.Port); err != nil {
		log.Error("rbi init failed", zap.Error(err))
	}
	log.Info("rbi init succ")
}

func init_register(ctx context.Context) {
	serverKey := fmt.Sprintf("%v%v", kenum.SER_GAME_PREFIX, config.Conf.ServerId)
	game_register = discover.NewServiceRegister(config.Conf.Etcd.Addrs, serverKey)
	if game_register == nil {
		panic("game_register create failed")
	}

	log.Info("init game_register")
}

func init_discover(ctx context.Context) {
	fight_discover = discover.NewServerDiscover()
	if err := fight_discover.Init(
		ctx,
		discover.WithEndpoints(config.Conf.Etcd.Addrs),
		discover.WithUpdateFn(fight.UpsertServer),
		discover.WithRemoveFn(fight.RemoveServer),
	); err != nil {
		panic(fmt.Sprintf("game_discover init err:%s", err.Error()))
	}
	log.Info("init fight_discover")
}

func init_db_global(ctx context.Context) {
	if err := mongo.Init(ctx, mongo.WithUri(config.Conf.GlobalMongo.Url), mongo.WithDb(config.Conf.GlobalMongo.DB)); err != nil {
		panic(fmt.Sprintf("mongo init err:%s", err.Error()))
	}
	db_global.SetDBByName(config.Conf.GlobalMongo.DB)
}

func init_server_data() {
	common.LoadServerInfo()
	log.Info("server data init succ")
}

func init_snow() {
	common2.InitSnowFlake(int64(config.Conf.ServerId))
	uid.Init()
	log.Info("snowflake init succ")
}

func init_grpc(ctx context.Context) {
}

func init_service(ctx context.Context) {
	service.Init(ctx)
}

func init_content(ctx context.Context) {
	content.Init(service.PostPanic, logic.PushContencCb)
	log.Info("content init succ")
}

func init_async(ctx context.Context) {
	async.Init(service.PostPanic, logic.PushAsyncCb)
	log.Info("async init succ")
}

func init_out(ctx context.Context) {
	io_out.Init(nats.Publish, service.PostPanic, 10)
	log.Info("out queue init succ")
}

func init_logic(ctx context.Context) {
	logic.Init(100, 1, 1, 1)
	log.Info("logic init succ")
}

func init_tapping(ctx context.Context) {
	tapping.Init(tapping.ProcessTapData, service.PostPanic, 10)
	log.Info("tapping init succ")
}
