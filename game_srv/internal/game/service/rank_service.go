package service

import (
	"encoding/json"
	"fmt"
	"gameserver/internal/common"
	"gameserver/internal/config"
	common2 "gameserver/internal/game/common"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
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

const (
	Main_Rank                 = "main_rank"      // 主线
	Challenge_Rank            = "challenge_rank" // 挑战
	Coin_Rank                 = "coin_rank"      // 金币
	Equip_Rank                = "equip_rank"     // 装备
	Weapon_Rank               = "weapon_rank"    // 武器
	Pet_Rank                  = "pet_rank"
	Stage_Star_Rank           = "mission_star_rank"
	LastMainStarRankResetTime = "last_main_star_rank_reset_time"
	rankCacheTime             = 1 * time.Minute
	mainStarRankClearTime     = 3
	StarRankCrycle            = 3
)

var rankTypeToResourcesPassType map[template.RankType]uint32

func init() {
	rankTypeToResourcesPassType = make(map[template.RankType]uint32)
	rankTypeToResourcesPassType[template.ResourcesPassCoinRank] = MoneyPass
	rankTypeToResourcesPassType[template.ResourcesPassDeputyWeaponRank] = SidearmPass
	rankTypeToResourcesPassType[template.ResourcesPassEquipRank] = EquipPass
	rankTypeToResourcesPassType[template.ResourcesPassPetRank] = PetPass
}

type SimpleMissionReward struct {
	*model.MissionReward
	Nick      string
	Head      uint32
	HeadFrame uint32
	State     uint32
}

func GetRankData(p *player.Player, rankType template.RankType) ([]*msg.CommonRankBaseData, []float64, uint64) {
	rankCfg := template.GetRankTemplate().GetRank(rankType)
	if rankCfg == nil {
		log.Error("get rank failed rankcfg is nil", zap.Any("rankType", rankType))
		return nil, nil, 0
	}

	selfInfo := make([]float64, 0)
	var nextRewardTime int64 = 0

	switch rankCfg.Id {
	case template.SessionRank, template.SessionFriendRank:
		nextRewardTime = int64(template.GetBattlePassTemplate().GetCurSeason(tools.GetCurTime()).WeekTime[1])
		selfInfo = append(selfInfo, float64(p.UserData.PeakFight.Cup))
	case template.MainEliteRank:
		selfInfo = appendMissionRankData(selfInfo, p.UserData.Mission.Challenges)

	case template.MainNormalRank:
		selfInfo = appendMissionRankData(selfInfo, p.UserData.Mission.Missions)

	case template.MainStarRank:
		selfInfo = append(selfInfo, float64(p.UserData.StageInfo.StageStar))
		info := common2.GetServerInfo()
		nextRewardTime = nextCycleTimeDays(info.OpenTime, time.Now().Unix(), StarRankCrycle)

	case template.ResourcesPassCoinRank, template.ResourcesPassDeputyWeaponRank,
		template.ResourcesPassEquipRank, template.ResourcesPassPetRank:
		selfInfo = appendResourcesPassData(selfInfo, p.UserData.ResourcesPass.PassList, rankTypeToResourcesPassType[rankCfg.Id])
		nextRewardTime = utils.GetNextWeekdayTimestamp(time.Monday)

	case template.WeekPassBlackBoosDamage:
		info := p.UserData.PlayMethod.Data
		var max int = 0
		for _, v := range info {
			if v.MaxDamage > max {
				max = v.MaxDamage
			}
		}
		selfInfo = append(selfInfo, float64(max))
	default:
		log.Error("unhandled rank type",
			zap.String("rank name", rankCfg.RankName),
			zap.Any("rank id", rankCfg.Id))
	}

	rankList := GetCommonRank(rankCfg, false)
	likesMap := GetAllRankLikes(rankType)

	for _, v := range rankList {
		for i := range v.PlayerInfo {
			if likesNum, ok := likesMap[v.PlayerInfo[i].AccountId]; ok {
				v.Likes = likesNum
			}
		}

	}

	return rankList, selfInfo, uint64(nextRewardTime)
}

func HandleAddLikesForRank(p *player.Player, rankType template.RankType, targetId uint64) (msg.ErrCode, int64) {
	likesInfo := p.UserData.Likes
	if v, ok := likesInfo.LikesMap[rankType]; ok && v {
		return msg.ErrCode_Rank_Is_Already_Like, -1
	}

	if targetId <= 0 {
		return msg.ErrCode_Rank_Target_Is_Not_Exist, -1
	}

	num := AddRankLikesNum(p, rankType, targetId)
	if num < 0 {
		return msg.ErrCode_Rank_Like_Failed, -1
	}

	likesInfo.LikesMap[rankType] = true
	p.UserData.Likes = likesInfo
	p.SaveLikes()
	return msg.ErrCode_SUCC, num
}

func HandleRefreshRankData() {
	cfg := template.GetRankTemplate().Data
	rankDataCache := make(map[template.RankType][]*msg.CommonRankBaseData)

	now := time.Now()

	for _, v := range cfg {
		if v.RankName == "" {
			continue
		}

		var shouldRefresh bool
		switch v.ResetType {
		case template.EveryDayZero:
			shouldRefresh = true
		case template.MondayZero:
			shouldRefresh = now.Weekday() == time.Monday && now.Hour() == 0
		case template.NeverReset:
			if v.Id == template.MainStarRank {
				key := fmt.Sprintf("%s:%v", LastMainStarRankResetTime, config.Conf.ServerId)
				rc := rdb_single.Get()
				timeStr, err := rc.Get(rdb_single.GetCtx(), key).Result()
				if err != nil {
					if err == redis.Nil {
						timeStr = "0"
					} else {
						log.Error("get main star rank last reset time error", zap.Error(err))
						continue
					}

				}

				var lastResetTime int64
				if timeStr == "" {
					lastResetTime = 0
				} else {
					lastResetTime, err = strconv.ParseInt(timeStr, 10, 64)
					if err != nil {
						log.Error("reset resources pass time trans error",
							zap.String("timeStr", timeStr),
							zap.Error(err))
						lastResetTime = 0
					}
				}

				currentTime := time.Now().Unix()
				days := utils.GetDaysBetweenTimestamps(lastResetTime, currentTime)
				if days >= mainStarRankClearTime {
					shouldRefresh = true
					_, err := rc.Set(rdb_single.GetCtx(), key, currentTime, 0).Result()
					if err != nil {
						log.Error("update main star rank redis key error",
							zap.String("key", key),
							zap.Error(err))
					}
				}
			}
		}

		if shouldRefresh {
			data := GetCommonRank(v, true)
			if len(data) > 0 {
				rankDataCache[v.Id] = data
			}
		}
	}

	if len(rankDataCache) == 0 {
		return
	}

	rankingTemplate := template.GetRankingTemplate()
	mailTemplate := template.GetMailTemplate()

	for rankType, rankList := range rankDataCache {
		for rankIndex, rankData := range rankList {
			rankNum := rankIndex + 1

			rankRewardCfg := rankingTemplate.GetRankReward(rankType, uint32(rankNum))
			if rankRewardCfg == nil {
				log.Error("rank reward cfg is nil",
					zap.Any("rankType", rankType),
					zap.Int("rank", rankNum))
				continue
			}

			mailCfg := mailTemplate.GetMail(rankRewardCfg.Id)
			if mailCfg == nil {
				log.Error("mail cfg is nil",
					zap.Any("mailId", rankRewardCfg.Id),
					zap.Any("rankType", rankType),
					zap.Int("rank", rankNum))
				continue
			}

			var items []*model.SimpleItem
			for _, item := range mailCfg.Reward {
				items = append(items, &model.SimpleItem{
					Id:  item.ItemId,
					Num: item.ItemNum,
				})
			}

			endTime := now.AddDate(0, 0, 60)

			for _, playerInfo := range rankData.PlayerInfo {
				accountIDStr := strconv.FormatUint(playerInfo.AccountId, 10)
				p := player.FindByAccount(accountIDStr)

				mail := model.NewMail(
					common.GenSnowFlake(),
					fmt.Sprintf("%v", mailCfg.Title),
					fmt.Sprintf("%v", mailCfg.Content),
					items,
					uint32(endTime.Unix()),
				)
				mail.MailType = 1

				if p == nil {
					AddOfflineMail(playerInfo.AccountId, mail)
				} else {
					AddSystemMail(p, mail)
				}

				log.Info("rank reward mail sent",
					zap.Uint64("accountId", playerInfo.AccountId),
					zap.Any("rankType", rankType),
					zap.Int("rank", rankNum))
			}
		}
	}

	log.Info("rank refresh completed",
		zap.Int("totalRankTypes", len(rankDataCache)),
		zap.Time("processTime", now))
}

func ToProtocolRankMissionReward(data *SimpleMissionReward) *msg.RankMissionRewardInfo {
	ret := &msg.RankMissionRewardInfo{}
	ret.MissionId = uint32(data.MissionId)
	ret.AccountId = uint32(data.Uid)
	ret.PassTime = data.PassTime
	ret.State = data.State
	if data.State == 0 {
		if _, data := GetPlayerSimpleInfo(data.Uid); data != nil {
			ret.Name = data.Name
			ret.Head = data.Head
			ret.HeadFrame = data.HeadFrame
		}
	}
	return ret
}

func HandleRefreshLikesInfo(p *player.Player) {
	likesInfo := p.UserData.Likes
	for v := range likesInfo.LikesMap {
		likesInfo.LikesMap[v] = false
	}

	p.UserData.Likes = likesInfo
	p.SaveLikes()
}

func GetFirstPassMaxReward(p *player.Player, rankType template.RankType) (msg.ErrCode, []*msg.SimpleItem) {
	if rankType != template.MainNormalRank && rankType != template.MainEliteRank {
		return msg.ErrCode_CONDITION_NOT_MET, nil
	}

	cfg := template.GetRankTemplate().GetRank(rankType)
	if cfg == nil {
		return msg.ErrCode_CONFIG_NIL, nil

	}

	var recordInfo map[uint32]bool
	maxInfo := GetMaxFirstPassRecord()
	var maxPassId uint32 = 0
	if rankType == template.MainNormalRank {
		recordInfo = p.UserData.Ranks.NormalPassRewardInfo
		maxPassId = maxInfo[0]
	} else {
		recordInfo = p.UserData.Ranks.ElitePassRewardInfo
		maxPassId = maxInfo[1]
	}

	rankRewardMap := template.GetRankRewardTemplate().GetMaxRankRewardsList(int(maxPassId))

	userId := p.GetUserId()
	itemList := make([]*msg.SimpleItem, 0)
	for passId, rewardList := range rankRewardMap {
		if _, ok := recordInfo[passId]; !ok {
			for i := range rewardList {
				AddItem(userId, rewardList[i].ItemId, int32(rewardList[i].ItemNum), publicconst.MaxFirstPassReward, true)
				recordInfo[passId] = true

				var item = &msg.SimpleItem{
					ItemId:  rewardList[i].ItemId,
					ItemNum: rewardList[i].ItemNum,
					Src:     rewardList[i].Src,
				}

				itemList = append(itemList, item)
			}
		}
	}

	return msg.ErrCode_SUCC, itemList
}

func GetFirstPassRecordData(p *player.Player, rankType template.RankType) (msg.ErrCode, map[uint32]*msg.PlayerSimpleInfo, []uint32) {
	if rankType != template.MainNormalRank && rankType != template.MainEliteRank {
		return msg.ErrCode_CONDITION_NOT_MET, nil, nil
	}

	data := GetFirstPassRecord(rankType)
	if data == nil {
		return msg.ErrCode_INVALID_DATA, nil, nil
	}

	resMap := make(map[uint32]*msg.PlayerSimpleInfo, 0)
	for passIdStr, playerInfoStr := range data {
		passId, err := strconv.ParseUint(passIdStr, 10, 64)
		var info = &msg.PlayerSimpleInfo{}
		err1 := json.Unmarshal([]byte(playerInfoStr), &info)

		if err != nil || err1 != nil {
			continue
		}

		resMap[uint32(passId)] = info
	}

	var selfRecordInfo map[uint32]bool
	if rankType == template.MainNormalRank {
		selfRecordInfo = p.UserData.Ranks.NormalPassRewardInfo
	} else {
		selfRecordInfo = p.UserData.Ranks.ElitePassRewardInfo
	}

	selfRecordList := make([]uint32, 0)
	for k, _ := range selfRecordInfo {
		selfRecordList = append(selfRecordList, k)
	}

	return msg.ErrCode_SUCC, resMap, selfRecordList
}

func GetRankRedPoint(p *player.Player) *msg.RedPointInfo {
	normalInfo := p.UserData.Ranks.NormalPassRewardInfo
	eliteInfo := p.UserData.Ranks.ElitePassRewardInfo
	for k, _ := range eliteInfo {
		if _, ok := normalInfo[k]; !ok {
			normalInfo[k] = true
		}
	}

	maxList := GetMaxFirstPassRecord()
	if len(maxList) <= 0 {
		return nil
	}

	isShow := false
	recordIdx := make([]int, 2)
	for i := range maxList {
		flag := false
		for passId, _ := range normalInfo {
			if maxList[i] == passId {
				flag = true
				break
			}
		}

		if !flag {
			isShow = true
			recordIdx = append(recordIdx, i+1)
		}
	}

	if isShow {
		ret := &msg.RedPointInfo{RdType: msg.RedPointType_MissionReward_Point}
		for i := range recordIdx {
			ret.RdData = append(ret.RdData, uint32(recordIdx[i]))
		}
		return ret
	}

	return nil
}

func NotifyRedPointToAllPlayer(rankType template.RankType) {
	for _, p := range player.AllPlayers() {
		ret := GetRankRedPoint(p)
		retMsg := &msg.ResponseRedPoint{
			Result: msg.ErrCode_SUCC,
		}

		retMsg.Data = append(retMsg.Data, ret)
		p.SendNotify(retMsg)
	}
}

func appendMissionRankData(selfInfo []float64, missionData []*model.Mission) []float64 {
	var maxID int
	var completeTime int

	for _, v := range missionData {
		if v.MissionId > maxID {
			maxID = v.MissionId
			completeTime = int(v.CompleteTime)
		}
	}

	return append(selfInfo, float64(maxID), float64(completeTime))
}

func appendResourcesPassData(selfInfo []float64, passList []*model.ResourcesPass, passType uint32) []float64 {
	for _, v := range passList {
		if v.PassType == passType {
			selfInfo = append(selfInfo, float64(v.Total))
			break
		}
	}
	return selfInfo
}

func nextCycleTime(startTimestamp, currentTimestamp, cycleSeconds int64) int64 {
	elapsedTime := currentTimestamp - startTimestamp
	completedCycles := elapsedTime / cycleSeconds

	if elapsedTime%cycleSeconds == 0 && currentTimestamp != startTimestamp {
		completedCycles++
	}

	nextCycleStart := startTimestamp + (completedCycles+1)*cycleSeconds

	nextCycleStart = adjustToZeroHour(nextCycleStart, cycleSeconds)

	return nextCycleStart
}

func nextCycleTimeDays(startTimestamp, currentTimestamp int64, cycleDays int) int64 {
	cycleSeconds := int64(cycleDays) * 24 * 60 * 60
	return nextCycleTime(startTimestamp, currentTimestamp, cycleSeconds)
}

func adjustToZeroHour(timestamp, cycleSeconds int64) int64 {
	secondsInDay := int64(24 * 60 * 60)
	days := timestamp / secondsInDay

	return days * secondsInDay
}
