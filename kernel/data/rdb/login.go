package rdb

import (
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"github.com/v587-zyf/gc/utils"
	"kernel/kenum"
	"time"
)

type LoginInfo []any

func (t *LoginInfo) String(index int32) (out string) {
	if index < 0 || index >= int32(len(*t)) {
		return ""
	}

	data, ok := (*t)[index].(string)
	if !ok {
		return ""
	}

	return data
}

func (t *LoginInfo) Uint64(index int32) uint64 {
	if index < 0 || index >= int32(len(*t)) {
		return 0
	}
	data, ok := (*t)[index].(string)
	if !ok {
		return 0
	}
	return utils.StrToUInt64(data)
}

func (t *LoginInfo) Int64(index int32) int64 {
	if index < 0 || index >= int32(len(*t)) {
		return 0
	}
	data, ok := (*t)[index].(string)
	if !ok {
		return 0
	}

	return utils.StrToInt64(data)
}

func (t *LoginInfo) UInt32(index int32) uint32 {
	if index < 0 || index >= int32(len(*t)) {
		return 0
	}
	data, ok := (*t)[index].(string)
	if !ok {
		return 0
	}

	return utils.StrToUInt32(data)
}

func SetUserLoginInfo(key string, value map[string]any) error {
	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
	)

	if err := rc.HMSet(rCtx, key, value).Err(); err != nil {
		return err
	}

	return rc.Expire(rCtx, key, Second(TD_OneDaySecond)).Err()
}

func GetUserLoginInfo(loginKey string) ([]any, error) {
	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
	)

	return rc.HMGet(rCtx, loginKey,
		kenum.Redis_Login_Token,
		kenum.Redis_Login_Gate,
		kenum.Redis_Login_UID,
		kenum.Redis_Login_Tda_Comm_Attr).Result()
}

func SetUserLoginToken(token string) error {
	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
	)

	return rc.Set(rCtx, kenum.Redis_Key_Token, token, time.Duration(TD_OneDaySecond)).Err()
}
