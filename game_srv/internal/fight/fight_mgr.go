package fight

import (
	"fmt"
	"kernel/loadbalancer"
	"strconv"
	"strings"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

type FightServer struct {
	Id   uint32
	Info string
}

type FightManager struct {
	lb *loadbalancer.LoadBalancer
}

var FightMgrIns *FightManager

func init() {
	FightMgrIns = &FightManager{
		lb: loadbalancer.NewLoadBalancer(),
	}
}

func UpsertServer(k, v string) {
	id, err := parse_fight_id(k)
	if err != nil {
		log.Error("parse fight id err", zap.Error(err))
		return
	}
	addr, load, err := parse_server_info(v)
	if err != nil {
		log.Error("parse server info err", zap.Error(err))
		return
	}
	if FightMgrIns.lb.Upsert(id, addr, load) {
		log.Info("discover new fight",
			zap.Uint32("id", id),
			zap.String("addr", addr),
			zap.Uint32("load", load),
		)
	} else {
		log.Info("discover update fight",
			zap.Uint32("id", id),
			zap.String("addr", addr),
			zap.Uint32("load", load),
		)
	}
}

func RemoveServer(k string) {
	id, err := parse_fight_id(k)
	if err != nil {
		log.Error("parse fight id err", zap.Error(err))
		return
	}
	info := FightMgrIns.lb.FindServer(id)
	if info == nil {
		log.Error("server nil when remove", zap.Uint32("id", id))
		return
	}
	FightMgrIns.lb.RemoveServer(id)
	log.Info("discover remove fight",
		zap.Uint32("id", id),
		zap.String("addr", info.Addr),
		zap.Uint32("load", info.Load),
	)
}

func SelectFight() *loadbalancer.ServeInfo {
	return FightMgrIns.lb.SelectServer()
}

func parse_fight_id(k string) (uint32, error) {
	ss := strings.Split(k, "/")
	if len(ss) != 4 {
		return 0, fmt.Errorf("parse fight id err %s", k)
	}
	id, err := strconv.ParseUint(ss[3], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parse fight id err %s", k)
	}
	return uint32(id), nil
}

func parse_server_info(v string) (string, uint32, error) {
	strs := strings.Split(v, "|")
	if len(strs) != 2 {
		return "", 0, fmt.Errorf("parse server info err %s", v)
	}

	load, err := strconv.ParseUint(strs[1], 10, 32)
	if err != nil {
		return "", 0, fmt.Errorf("parse server's load failed %s", v)
	}

	return strs[0], uint32(load), nil
}
