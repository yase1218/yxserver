package service

import (
	"encoding/json"
	"fmt"
	"gameserver/internal/config"
	"gameserver/internal/game/player"
	"kernel/tools"
	"msg"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"github.com/v587-zyf/gc/utils"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

func RdbUserNameKey() string {
	return fmt.Sprintf("%v:UserNames", config.Conf.ServerId)
}

func RdbOfflineMailKey(uid uint64) string {
	return fmt.Sprintf("offline_mail:%v:%v", config.Conf.ServerId, uid)
}

func GetFirstPassRecordKeys(rankType template.RankType) string {
	switch rankType {
	case template.MainNormalRank:
		return fmt.Sprintf("first_pass_normal_player_record:%v", config.Conf.ServerId)
	case template.MainEliteRank:
		return fmt.Sprintf("first_pass_elite_player_record:%v", config.Conf.ServerId)
	}

	return ""
}

func GetRankKeys(rankCfg *template.JRank, isLike bool) string {
	if rankCfg == nil {
		return ""
	}

	name := rankCfg.RankName
	if isLike {
		name = fmt.Sprintf("%s_likes", name)
	}

	baseKey := fmt.Sprintf("%s:%v", name, config.Conf.ServerId)

	switch rankCfg.Id {
	case template.SessionRank:
		sessionStartTs := template.GetBattlePassTemplate().GetCurSeason(tools.GetCurTime()).BattleTime[0]
		timeStr := time.Unix(int64(sessionStartTs), 0).Format("2006_01_02")
		return fmt.Sprintf("%v:%v", baseKey, timeStr)
	case template.MainStarRank, template.MainEliteRank, template.MainNormalRank,
		template.WeekPassBlackBoosDamage:
		return baseKey
	case template.ResourcesPassCoinRank, template.ResourcesPassDeputyWeaponRank,
		template.ResourcesPassEquipRank, template.ResourcesPassPetRank:
		if !isLike {
			timeStr := getResetTimeString(rankCfg.ResetType)
			if timeStr == "" {
				return ""
			}
			return fmt.Sprintf("%s:%s", baseKey, timeStr)
		}
		return baseKey
	}

	return ""
}

func GetLastKey(cfg *template.JRank) string {
	if cfg == nil {
		return ""
	}

	var targetTime time.Time
	switch cfg.ResetType {
	case template.EveryDayZero:
		targetTime = utils.GetYesterdayZeroTime()
	case template.MondayZero:
		targetTime = utils.GetWeekdayZeroTime(time.Monday, true)
	default:
		log.Error("get last rank key error",
			zap.Any("resetType", cfg.ResetType),
			zap.Any("rankId", cfg.Id))
		return ""
	}

	return fmt.Sprintf("%v:%s:%s", config.Conf.ServerId, cfg.RankName, targetTime.Format("2006_01_02"))
}

func getResetTimeString(resetType template.RankResetType) string {
	switch resetType {
	case template.EveryDayZero:
		return time.Now().Format("2006_01_02")
	case template.MondayZero:
		return utils.GetWeekdayZeroTime(time.Monday, false).Format("2006_01_02")
	default:
		return ""
	}
}

func GetCommonRank(rankCfg *template.JRank, isLast bool) []*msg.CommonRankBaseData {
	key := getRankKey(rankCfg, isLast)
	if key == "" {
		return nil
	}

	rc := rdb_single.Get()

	res, err := rc.ZRevRangeByScoreWithScores(rdb_single.GetCtx(), key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  rankCfg.Max,
	}).Result()

	if err != nil {
		log.Error("get common rank from redis error",
			zap.Error(err),
			zap.String("key", key))
		return nil
	}

	if len(res) == 0 {
		return make([]*msg.CommonRankBaseData, 0)
	}

	rankList := make([]*msg.CommonRankBaseData, 0, len(res))

	for _, data := range res {
		memberStr, ok := data.Member.(string)
		if !ok {
			log.Warn("member is not string type",
				zap.Any("member", data.Member),
				zap.String("key", key))
			continue
		}

		var playerInfo *msg.CommonRankBaseData
		if err := json.Unmarshal([]byte(memberStr), &playerInfo); err != nil {
			log.Error("unmarshal player simple info error",
				zap.Error(err),
				zap.String("memberStr", memberStr),
				zap.String("key", key))
			continue
		}

		if playerInfo != nil {
			rankList = append(rankList, playerInfo)
		}
	}

	return rankList
}

func UpdateCommonRankInfo(p *player.Player, args interface{}, rankType template.RankType) {
	rankCfg := template.GetRankTemplate().GetRank(rankType)
	if rankCfg == nil {
		log.Error("update player rank error: config is nil",
			zap.Int("rankType", int(rankType)),
			zap.Any("args", args))
		return
	}

	score, rankInfo, err := parseRankArgs(args, rankType)
	if err != nil {
		log.Error("parse rank args error",
			zap.Error(err),
			zap.Int("rankType", int(rankType)))
		return
	}

	code, playerInfo := GetPlayerSimpleInfo(p.GetUserId())
	if code != msg.ErrCode_SUCC {
		log.Error("get player simple info failed",
			zap.Int32("code", int32(code)),
			zap.Int("rankType", int(rankType)))
		return
	}

	simpleInfo := ToPlayerSimpleInfo(playerInfo)
	rankInfo.PlayerInfo = append(rankInfo.PlayerInfo, simpleInfo)

	data, err := json.Marshal(rankInfo)
	if err != nil {
		log.Error("marshal player simple info error",
			zap.Error(err),
			zap.Int("rankType", int(rankType)))
		return
	}

	if err := updatePlayerRankData(rankCfg, p.GetUserId(), score, data); err != nil {
		log.Error("update player rank data error",
			zap.Error(err),
			zap.Int("rankType", int(rankType)))
	}
}

func AddRankLikesNum(p *player.Player, rankType template.RankType, targetId uint64) int64 {
	rankCfg := template.GetRankTemplate().GetRank(rankType)
	if rankCfg == nil {
		log.Error("get rank cfg error", zap.Any("rankType", rankType))
		return -1
	}

	key := GetRankKeys(rankCfg, true)
	if key == "" {
		log.Error("rank key is empty", zap.Any("rankType", rankType))
		return -1
	}

	rc := rdb_single.Get()

	accountStr := strconv.FormatUint(uint64(targetId), 10)
	newLikesNum, err := rc.HIncrBy(rdb_single.GetCtx(), key, accountStr, 1).Result()
	if err != nil {
		log.Error("increment rank likes num error",
			zap.Uint64("targetId", targetId),
			zap.String("key", key),
			zap.Error(err))
		return -1
	}

	log.Debug("add rank likes num success",
		zap.Uint64("targetId", targetId),
		zap.String("key", key),
		zap.Int64("newLikesNum", newLikesNum))

	return newLikesNum
}

func getRankKey(rankCfg *template.JRank, isLast bool) string {
	if isLast {
		return GetLastKey(rankCfg)
	}
	return GetRankKeys(rankCfg, false)
}

func parseRankArgs(args interface{}, rankType template.RankType) (float64, *msg.CommonRankBaseData, error) {
	rankInfo := &msg.CommonRankBaseData{}
	var score float64

	switch val := args.(type) {
	case uint32:
		score = float64(val)
		rankInfo.Extra = []float64{score}
	case []float64:
		if len(val) >= 2 {
			score = val[0]*1000000 + val[1]
		} else if len(val) == 1 {
			score = val[0]
		}
		rankInfo.Extra = val
	default:
		return 0, nil, fmt.Errorf("unsupported args type: %T for rankType: %d", args, rankType)
	}

	return score, rankInfo, nil
}

func updatePlayerRankData(rankCfg *template.JRank, playerId uint64, score float64, data []byte) error {
	key := GetRankKeys(rankCfg, false)
	if key == "" {
		return fmt.Errorf("rank key is empty for rankId: %d", rankCfg.Id)
	}

	rc := rdb_single.Get()

	ctx := rdb_single.GetCtx()

	oldMembers, err := rc.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()
	if err != nil {
		return fmt.Errorf("get old rank data error: %w", err)
	}

	playerIdStr := strconv.FormatUint(playerId, 10)
	for _, member := range oldMembers {
		var oldRankInfo msg.CommonRankBaseData
		if err := json.Unmarshal([]byte(member), &oldRankInfo); err == nil {
			if len(oldRankInfo.PlayerInfo) > 0 &&
				strconv.FormatUint(uint64(oldRankInfo.PlayerInfo[0].GetAccountId()), 10) == playerIdStr {
				if _, err := rc.ZRem(ctx, key, member).Result(); err != nil {
					log.Error("remove old rank data error",
						zap.Error(err),
						zap.String("key", key))
				}
				break
			}
		}
	}

	_, err = rc.ZAdd(ctx, key, redis.Z{Score: score, Member: data}).Result()
	if err != nil {
		return fmt.Errorf("add new rank data error: %w", err)
	}

	if _, err := rc.ZRemRangeByRank(ctx, key, rankCfg.Max, -1).Result(); err != nil {
		log.Error("remove overflow rank data error",
			zap.Error(err),
			zap.String("key", key))
	}

	log.Debug("update player rank success",
		zap.String("key", key),
		zap.Uint64("playerId", playerId),
		zap.Float64("score", score))

	return nil
}

func GetAllRankLikes(rankType template.RankType) map[uint64]uint32 {
	rankCfg := template.GetRankTemplate().GetRank(rankType)
	if rankCfg == nil {
		log.Error("get rank cfg error", zap.Any("rankType", rankType))
		return nil
	}

	key := GetRankKeys(rankCfg, true)
	if key == "" {
		return nil
	}

	rc := rdb_single.Get()

	result, err := rc.HGetAll(rdb_single.GetCtx(), key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Error("get all rank likes error",
				zap.String("key", key),
				zap.Error(err))
		}
		return nil
	}

	likesMap := make(map[uint64]uint32, len(result))
	for keyStr, value := range result {
		targetId, err := strconv.ParseUint(keyStr, 10, 32)
		if err != nil {
			log.Warn("parse target id error, skip",
				zap.String("key", keyStr),
				zap.Error(err))
			continue
		}

		likesNum, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Warn("parse likes num error, skip",
				zap.String("value", value),
				zap.Error(err))
			continue
		}

		likesMap[uint64(targetId)] = uint32(likesNum)
	}

	return likesMap
}

func GetRankMaxInfo(cfg *template.JRank) *msg.CommonRankBaseData {
	var (
		keys = GetRankKeys(cfg, false)
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
	)

	res, err := rc.ZRevRangeByScoreWithScores(rCtx, keys, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  1,
	}).Result()

	if err != nil {
		if err == redis.Nil {
			log.Info("redis rank data is empty", zap.String("key", keys))
		} else {
			log.Error("redis query error", zap.String("key", keys), zap.Error(err))
		}
		return nil
	}

	if len(res) == 0 {
		log.Info("no rank data found", zap.String("key", keys))
		return nil
	}

	memberStr, ok := res[0].Member.(string)
	if !ok {
		log.Error("invalid member type in rank data",
			zap.String("key", keys),
			zap.Any("member", res[0].Member))
		return nil
	}

	info := &msg.CommonRankBaseData{}
	if err := json.Unmarshal([]byte(memberStr), info); err != nil {
		log.Error("unmarshal rank data error",
			zap.Error(err),
			zap.String("memberStr", memberStr),
			zap.String("key", keys))
		return nil
	}

	return info
}

func RecordFirstPassPlayerInfo(p *player.Player, missionId int, rankType template.RankType) {
	var (
		keys = GetFirstPassRecordKeys(rankType)
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
	)

	code, playerInfo := GetPlayerSimpleInfo(p.GetUserId())
	if code != msg.ErrCode_SUCC {
		log.Error("get player simple info failed",
			zap.Int32("code", int32(code)),
		)
		return
	}

	simpleInfo := ToPlayerSimpleInfo(playerInfo)
	infoStr, err := json.Marshal(simpleInfo)
	if err != nil {
		log.Error("marshal player simple data error",
			zap.Error(err),
			zap.Any("simpleInfo", simpleInfo))
		return
	}

	if _, err := rc.HSet(rCtx, keys, missionId, infoStr).Result(); err != nil {
		log.Error("recore player pass info error", zap.Int("missionId", missionId), zap.Any("info", infoStr), zap.Error(err))
		return
	}

	NotifyRedPointToAllPlayer(rankType)
}

func GetFirstPassRecord(rankType template.RankType) map[string]string {
	var (
		keys = GetFirstPassRecordKeys(rankType)
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
	)

	if keys == "" {
		return nil
	}

	res, err := rc.HGetAll(rCtx, keys).Result()
	if err != nil {
		log.Error("get first pass data hgetall error",
			zap.String("key", keys),
			zap.Error(err),
		)
		return nil
	}

	return res
}

func GetMaxFirstPassRecord() []uint32 {
	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
	)

	getMaxPassID := func(rankType template.RankType) uint32 {
		key := GetFirstPassRecordKeys(rankType)
		res, err := rc.HGetAll(rCtx, key).Result()
		if err != nil {
			log.Error("hgetall max first pass record error",
				zap.Error(err),
				zap.Any("rankType", rankType))
			return 0
		}

		var maxPassID uint32
		for passIDStr := range res {
			passID, err := strconv.ParseUint(passIDStr, 10, 32)
			if err != nil {
				continue
			}

			if uint32(passID) > maxPassID {
				maxPassID = uint32(passID)
			}
		}
		return maxPassID
	}

	normalMaxPassID := getMaxPassID(template.MainNormalRank)
	eliteMaxPassID := getMaxPassID(template.MainEliteRank)

	return []uint32{normalMaxPassID, eliteMaxPassID}
}
