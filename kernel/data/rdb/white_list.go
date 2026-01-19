package rdb

import (
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"go.uber.org/zap"

	"kernel/kenum"
)

func UpdateWhiteList(list []string) (err error) {
	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
	)

	if err = rc.Del(rCtx, kenum.Redis_Key_Whitelist).Err(); err != nil {
		log.Error("del whitelist err", zap.Error(err))
		return err
	}

	if err = rc.SAdd(rCtx, kenum.Redis_Key_Whitelist, list).Err(); err != nil {
		log.Error("sadd whitelist err", zap.Error(err))
		return
	}

	return
}

func IsInWhiteList(accountId string) (bool, error) {
	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
	)

	return rc.SIsMember(rCtx, kenum.Redis_Key_Whitelist, accountId).Result()
}
