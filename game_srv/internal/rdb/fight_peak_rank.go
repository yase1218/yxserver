package rdb

import (
	"errors"
	"gameserver/internal/enum"
	"kernel/tools"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"github.com/v587-zyf/gc/utils"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

const (
	PeakFightRankScoreHighPartOffset = 1e12
	PeakFightRankScoreLowPartOffset  = 1e6
)

func MarshalPFRScore(score float64) float64 {
	highPart := int64(score) << 32
	lowPart := int64(^uint32(time.Now().Unix()))

	score = float64(highPart | lowPart)

	return score
}

func UnmarshalPFRScore(score float64) float64 {
	return float64(uint64(score) >> 32)
}

func AddPeakFightRank(uid uint64, data float64) error {
	log.Info("AddPeakFightRank", zap.Uint64("uid", uid), zap.Float64("data", data))
	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()

		accountStr = strconv.Itoa(int(uid))

		//timestamp  = time.Now().UnixMilli()
		//highPart = int64(data) * PeakFightRankScoreHighPartOffset
		//lowPart  = timestamp / PeakFightRankScoreLowPartOffset
		//score    = highPart + lowPart
		score = MarshalPFRScore(data)

		season  = template.GetBattlePassTemplate().GetCurSeason(tools.GetCurTime()).Season
		rankKey = FormatPeakFightRank(season)
	)

	isExists, err := rc.Exists(rCtx, rankKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Error("peak fight rank", zap.Error(err))
		return err
	}

	member := redis.Z{
		Member: accountStr,
		Score:  float64(score),
	}
	if err := rc.ZAdd(rCtx, rankKey, member).Err(); err != nil {
		log.Error("peak fight rank", zap.Error(err))
		return err
	}

	if isExists == 0 {
		if err := rc.Expire(rCtx, rankKey, enum.Redis_PeakFight_Rank_Expire).Err(); err != nil {
			log.Error("peak fight rank expire", zap.Error(err))
			return err
		}
	}
	return nil
}

func GetPeakFightRank() ([]redis.Z, error) {
	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()

		season = template.GetBattlePassTemplate().GetCurSeason(tools.GetCurTime()).Season
	)

	// 为了减轻服务器压力,暂时只查前100名数据
	vals, err := rc.ZRevRangeWithScores(rCtx, FormatPeakFightRank(season), 0, 99).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	for i, val := range vals {
		//vals[i].Score = val.Score / PeakFightRankScoreHighPartOffset
		vals[i].Score = UnmarshalPFRScore(val.Score)
	}

	return vals, nil
}

func GetUserPeakFightRankScore(uid uint64) (float64, error) {
	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()

		accountStr = strconv.FormatUint(uid, 10)
		season     = template.GetBattlePassTemplate().GetCurSeason(tools.GetCurTime()).Season
	)

	score, err := rc.ZScore(rCtx, FormatPeakFightRank(season), accountStr).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}

	//score = score / PeakFightRankScoreHighPartOffset
	score = UnmarshalPFRScore(score)

	return score, nil
}

func GetUserPeakFightRanking(uid uint64) (int64, error) {
	var (
		rc     = rdb_single.Get()
		rCtx   = rdb_single.GetCtx()
		season = template.GetBattlePassTemplate().GetCurSeason(tools.GetCurTime()).Season
	)
	return rc.ZRevRank(rCtx, FormatPeakFightRank(season), strconv.FormatUint(uid, 10)).Result()
}

// accountId:ranking
func GetPeakFightRankMembers(season uint32) (map[int64]uint32, error) {
	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()

		rankKey = FormatPeakFightRank(season)
		members = make(map[int64]uint32)
	)

	vals, err := rc.ZRevRangeWithScores(rCtx, rankKey, 0, 99).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Error("get rank err", zap.Error(err), zap.String("key", rankKey))
		return nil, err
	}
	for ranking, val := range vals {
		accountId := utils.StrToInt64(val.Member.(string))
		members[accountId] = uint32(ranking + 1)
		log.Info("peak fight rank", zap.String("rankKey", rankKey), zap.Int64("accountId", accountId), zap.Uint32("ranking", members[accountId]))
	}
	return members, nil
}
