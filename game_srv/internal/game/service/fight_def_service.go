package service

import (
	"encoding/json"
	"errors"
	"gameserver/internal/fight"
	"gameserver/internal/game/common"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"kernel/errCode"
	"kernel/protocol"
	"kernel/tools"
	"msg"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/rdb/rdb_single"

	"github.com/v587-zyf/gc/errcode"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

var (
	user_fight_map map[uint64]uint32
)

func init() {
	user_fight_map = make(map[uint64]uint32)
}

func AddFight(uid uint64, fightId uint32) {
	user_fight_map[uid] = fightId
}
func LoadFightId(uid uint64) uint32 {
	if v, ok := user_fight_map[uid]; ok {
		return v
	}
	return 0
}
func DelFight(uid uint64) {
	delete(user_fight_map, uid)
}
func RangeFight(fn func(key, value any) bool) {
	for k, v := range user_fight_map {
		if !fn(k, v) {
			break
		}
	}
}

func ClearFight(sendFight bool) {
	for k := range user_fight_map {
		p := player.FindByUserId(k)
		if p == nil {
			continue
		}

		if sendFight {
			p.IsSendEndFight = true
		}

		LeaveFight(p)
		DelFight(k)
	}
}

// 公共离开战斗 向战斗服发送离开消息
func LeaveFight(p *player.Player) {
	bc := GetBattleCacheByUid(p.GetUserId())
	if bc != nil {
		if bc.FsId != 0 {
			req := &msg.FsLeaveFightReq{
				FightId: bc.BattleId,
				EndType: p.EndFightType,
			}

			SendToFight(p, 0, req)
		}

		p.IsSendEndFight = false
		p.UserData.Fight.Clear()
		//p.SetFsId(0)
		p.SaveFight()
	}
}

func PauseFight(p *player.Player) {
	if p.GetBattleId() > 0 {
		req := &msg.FightSwitchReq{
			Type: msg.FightSwitchType_Fight_Switch_Pause,
		}

		SendToFight(p, 0, req)
	}
}

// 公共战前检查
func CheckEnterFightBefore(dbFightInfo *model.Fight, stageId int) bool {
	if dbFightInfo.FightStageId > 0 && dbFightInfo.FightStageId == stageId {
		return false
	}
	return true
}

func CreateFight(p *player.Player, stageId uint32, packetId uint32, extra []byte) msg.ErrCode {
	log.Info("create fight", zap.Uint64("uid", p.GetUserId()), zap.Uint32("fightStageId", stageId))

	if _, ok := msg.OutMap_name[int32(stageId)]; !ok {
		stageCfg := template.GetMissionTemplate().GetMission(int(stageId))
		if stageCfg == nil {
			log.Error("stage cfg nil", zap.Uint32("stageId", stageId))
			return msg.ErrCode_CONFIG_NIL
		}

		curWeek := time.Now().Weekday()

		// different fight enter before check
		switch msg.BattleType(stageCfg.Type) {
		case msg.BattleType_Battle_Main:
			pbErr := MainFightBeforeCheck(p, stageCfg)
			if pbErr != msg.ErrCode_SUCC {
				log.Error("main fight before check err", zap.Uint64("uid", p.GetUserId()), zap.String("pbErr", pbErr.String()))
				return pbErr
			}
		case msg.BattleType_Battle_Challenge:
			pbErr := ChallengeFightBeforeCheck(p, stageCfg)
			if pbErr != msg.ErrCode_SUCC {
				log.Error("challenge fight before check err", zap.Uint64("uid", p.GetUserId()), zap.String("pbErr", pbErr.String()))
				return pbErr
			}
		case msg.BattleType_Battle_Union:
			pbErr := PlayerMethodFightBeforeCheck(p, stageCfg)
			if pbErr != msg.ErrCode_SUCC {
				log.Error("play method fight not check ok", zap.Uint64("uid", p.GetUserId()),
					zap.Uint32("stageId", stageId))
				return pbErr
			}

		case msg.BattleType_Battle_Peak:
			if ec := PeakFightCheck(p, stageId); ec != msg.ErrCode_SUCC {
				log.Error("peak fight not check ok", zap.Uint64("uid", p.GetUserId()),
					zap.Uint32("stageId", stageId))
				return ec
			}

			battlePassCfg := template.GetBattlePassTemplate().GetCurSeason(tools.GetCurTime())
			if battlePassCfg == nil {
				log.Error("battle pass cfg nil", zap.Uint32("curTime", tools.GetCurTime()))
				return msg.ErrCode_CONFIG_NIL
			}

			fightExtra := &protocol.PeakFightExtra{
				StageBuff: []uint32{uint32(battlePassCfg.StageBuff)},
			}
			weekFields, _ := template.GetBattlePassTemplate().GetWeeklyBattleField(tools.GetCurTime())
			for i := 0; i < len(weekFields); i++ {
				fightExtra.StageBuff = append(fightExtra.StageBuff, uint32(weekFields[i]))
			}

			extraBytes, err := fightExtra.Marshal()
			if err != nil {
				log.Error("zombie extra marshal err", zap.Error(err), zap.Uint64("uid", p.GetUserId()), zap.Any("extra", extra))
				return msg.ErrCode_SYSTEM_ERROR
			}
			extra = extraBytes
		case msg.BattleType_Battle_Zombie:
			// todo 竞技场战前检查 重构完再写
			//if ec := ServMgr.GetArenaService().CheckBattle(p); ec != msg.ErrCode_SUCC {
			//	return errCode.ERR_FIGHT_BEFORE_CHECK
			//}
			//
			//fightExtra := ServMgr.GetArenaService().GetFightParam(p)
			//if fightExtra == nil {
			//	log.Error("zombie GetFightParam err", zap.Uint64("uid", p.GetUserId()), zap.Any("extra", extra))
			//	return errCode.ERR_FIGHT_BEFORE_CHECK
			//}
			//extraBytes, err := fightExtra.Marshal()
			//if err != nil {
			//	log.Error("zombie extra marshal err", zap.Error(err), zap.Uint64("uid", p.GetUserId()), zap.Any("extra", extra))
			//	err = errcode.ERR_JSON_MARSHAL_ERR
			//	return err
			//}
			//extra = extraBytes
		case msg.BattleType_Battle_EquipStage:

		case msg.BattleType_Battle_Desert: // 黑频Boss
			if curWeek != time.Monday && curWeek != time.Wednesday && curWeek != time.Friday {
				return msg.ErrCode_CONDITION_NOT_MET
			}

			info := common.GetServerInfo()
			weeks := tools.GetWeekCount(info.OpenTime)
			stageIds := template.GetPlayMethodStageTemplate().GetTargetStageId(uint32(msg.BATTLETYPE_BATTLE_DESERT), uint32(weeks))[0]
			if stageId != stageIds {
				return msg.ErrCode_CONFIG_NIL
			}

		case msg.BattleType_Battle_Contract:
			if _, ok := p.UserData.WeekPass.ContractInfo[stageId]; ok {
				return msg.ErrCode_CONDITION_NOT_MET
			}

			if curWeek != time.Tuesday && curWeek != time.Thursday {
				return msg.ErrCode_CONDITION_NOT_MET
			}

			info := common.GetServerInfo()
			weeks := tools.GetWeekCount(info.OpenTime)
			stageIds := template.GetPlayMethodStageTemplate().GetTargetStageId(uint32(msg.BattleType_Battle_Contract), uint32(weeks))
			isExist := false
			for _, v := range stageIds {
				if v == stageId {
					isExist = true
					break
				}
			}

			if !isExist {
				return msg.ErrCode_CONDITION_NOT_MET
			}

		case msg.BattleType_Battle_Secret:
			if curWeek != time.Saturday && curWeek != time.Sunday {
				return msg.ErrCode_CONDITION_NOT_MET
			}

			SecretMaxTimes := template.GetPlayMethodTemplate().GetPlayMethod(int(msg.BattleType_Battle_Secret)).Limit
			if SecretMaxTimes <= 0 || p.UserData.WeekPass.SecretCount >= int32(SecretMaxTimes) {
				return msg.ErrCode_ARENA_PK_CNT_NOT_ENOUGH
			}

			info := common.GetServerInfo()
			weeks := tools.GetWeekCount(info.OpenTime)
			stageIds := template.GetPlayMethodStageTemplate().GetTargetStageId(uint32(msg.BattleType_Battle_Secret), uint32(weeks))[0]
			if stageId != stageIds {
				return msg.ErrCode_CONFIG_NIL
			}
		}

		p.FightType = msg.BattleType(stageCfg.Type)
		p.LastFightTime = time.Now()
	} else {
		switch msg.OutMap(stageId) {
		case msg.OutMap_Out_Map_Union:
			// todo 联盟战外地图 联盟重构完再写
			//member, err := dao.GetMember(p.AccountInfo.AccountId)
			//if err != nil || member == nil {
			//	log.Error("not join alliance", zap.Uint64("uid", p.GetUserId()))
			//	return errCode.ERR_FIGHT_BEFORE_CHECK
			//}
			//
			//extraInfo := &protocol.AlliancePublicExtra{
			//	AllianceId: member.AllianceID,
			//}
			//extraBytes, err := extraInfo.Marshal()
			//if err != nil {
			//	log.Error("alliance public extra marshal err", zap.Error(err), zap.Uint64("uid", p.GetUserId()), zap.Any("extra", extra))
			//	err = errcode.ERR_JSON_MARSHAL_ERR
			//	return err
			//}
			//extra = extraBytes
		}
	}

	//p.SaveFight()

	createFightReq := &msg.FsCreateFightReq{
		StageId: stageId,
		Extra:   extra,
	}
	fs := fight.SelectFight()
	if fs == nil {
		log.Error("create fight but no fight server can use", zap.Uint32("stageId", stageId), ZapUser(p))
		return msg.ErrCode_SYSTEM_ERROR
	}
	p.SetFsId(fs.ID)
	SetBattleCacheFsId(p.GetUserId(), fs.ID)
	SendToFight(p, packetId, createFightReq)
	return msg.ErrCode_SUCC
}

// 开始战斗
func StartFight(p *player.Player, req *msg.StartFightReq, packetId uint32, extra []byte) error {
	bc := GetBattleCacheByUid(p.GetUserId())
	if bc.BattleId == 0 {
		return errCode.ERR_USER_NOT_IN_FIGHT
	}

	stageCfg := template.GetMissionTemplate().GetMission(int(bc.StageId))
	if stageCfg == nil {
		log.Error("stage cfg nil", zap.Uint32("stageId", bc.StageId))
		return errcode.ERR_CONFIG_NIL
	}

	startFight := &msg.FsStartFight{FightId: bc.FsId}
	switch msg.BattleType(stageCfg.Type) {
	case msg.BattleType_Battle_Peak:
		return EnterPeakFight(p)
	case msg.BattleType_Battle_Zombie:
		UpdateMissionPoker(p, stageCfg.Id)
		SendToFight(p, packetId, startFight)
	case msg.BattleType_Battle_Desert:
		fallthrough
	case msg.BattleType_Battle_Union:
		fallthrough
	case msg.BattleType_Battle_EquipStage:
		fallthrough
	case msg.BattleType_Battle_Challenge:
		fallthrough
	case msg.BattleType_Battle_Contract:
		fallthrough
	case msg.BattleType_Battle_Secret:
		fallthrough
	case msg.BattleType_Battle_Main:
		SendToFight(p, packetId, startFight)
	default:
		log.Error("battle type not support", zap.Uint32("stageId", bc.StageId),
			zap.Int32("battleType", int32(p.FightType)))
	}

	//grpc.GetGameGrpc().SendStartFight(p.UserData.GetAccountId(), p.UserData.GetFightId())

	return nil
}

// 战斗结束
func EndFight(p *player.Player, req *msg.EndFightReq, packetId uint32, extra []byte) error {
	if p.GetBattleId() == 0 {
		return errCode.ERR_USER_NOT_IN_FIGHT
	}
	if p.UserData.Fight.FightStageId == 0 {
		log.Error("user fight stageId zero", zap.Uint64("userId", p.GetUserId()),
			zap.Int("fightStageId", p.UserData.Fight.FightStageId), zap.Uint32("fight", p.UserData.Fight.FightId))
		return errCode.ERR_USER_NOT_IN_FIGHT
	}

	if req.GetMapType() == msg.MapType_Map_Type_Out {
		_, ok := msg.OutMap_name[int32(p.UserData.Fight.FightStageId)]
		if !ok {
			log.Error("end fight. not out_map fight", zap.Uint64("userId", p.GetUserId()),
				zap.Int("fightStageId", p.UserData.Fight.FightStageId), zap.Uint32("fight", p.UserData.Fight.FightId))
			return nil
		}
	}

	//var (
	//	err error
	//)

	_, inOutMap := msg.OutMap_name[int32(p.UserData.Fight.FightStageId)]
	if !inOutMap {
		stageCfg := template.GetMissionTemplate().GetMission(p.UserData.Fight.FightStageId)
		if stageCfg == nil {
			log.Error("stage cfg nil", zap.Int("stageId", p.UserData.Fight.FightStageId))
			return errcode.ERR_CONFIG_NIL
		}

		switch msg.BattleType(stageCfg.Type) {
		case msg.BattleType_Battle_Peak:
			// todo 重构pvp后看怎么处理玩家数据
			PFUserExit(p, req.EndType)
		case msg.BattleType_Battle_Zombie:
			// todo 重构竞技场后看怎么处理玩家数据
			//ServMgr.GetArenaService().FightResult(p.UserData, false, 0, 0)
		case msg.BattleType_Battle_Main:
		case msg.BattleType_Battle_Challenge:
		case msg.BattleType_Battle_Desert:
		case msg.BattleType_Battle_Union:
		case msg.BattleType_Battle_EquipStage:
		case msg.BattleType_Battle_Contract:
		case msg.BattleType_Battle_Secret:
		default:
			log.Error("battle type not support", zap.Int("stageId", p.UserData.Fight.FightStageId))
			return errcode.ERR_PARAM
		}

		//if err != nil {
		//	return err
		//}
	}
	p.EndFightType = req.GetEndType()

	LeaveFight(p)
	//ServMgr.GetFightService().LeaveFight(p.UserData)

	return nil
}

// 广播
func Broadcast(p *player.Player, req *msg.FightBroadcastReq, packetId uint32, extra []byte) error {
	if p.GetBattleId() == 0 {
		return nil
	}

	stageCfg := template.GetMissionTemplate().GetMission(p.UserData.Fight.FightStageId)
	if stageCfg == nil {
		log.Error("stage cfg nil", zap.Int("stageId", p.UserData.Fight.FightStageId))
		return errcode.ERR_CONFIG_NIL
	}

	switch msg.BattleType(stageCfg.Type) {
	case msg.BattleType_Battle_Peak:
		// todo pvp 广播 重构完再写
		PFBroadcast(p, req)
	}

	return nil
}

// 战斗一共有哪些怪
func FightMonster(p *player.Player, req *msg.FightMonsterReq, packetId uint32, extra []byte) error {
	ActivateAtlasByType(p, template.HandBook_Type_Monster, req.GetMonsterIds())
	return nil
}

// 获取海边派对周环境
func GetBattlePassWeekFields() ([]int32, uint32) {
	return template.GetBattlePassTemplate().GetWeeklyBattleField(tools.GetCurTime())
}

// 获取海边派对周环境
func InitBattlePassWeekField() {
	curTime := tools.GetCurTime()
	cfg := template.GetBattlePassTemplate().GetCurSeason(curTime)
	nextWeekRefreshTime := tools.GetWeeklyRefreshTime(0)
	if cfg != nil {
		rc := rdb_single.Get()
		rCtx := rdb_single.GetCtx()
		key := cfg.GetWeeklyBattleFieldRedisKey()
		nextTimeStr, err := rc.HGet(rCtx, key, template.BattlePassNextWeekRefreshTime).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				refreshWeeklyBattleField(cfg, nextWeekRefreshTime)
				return
			} else {
				log.Error("JBattlePass init get nextWeekRefreshTime err", zap.Uint32("id", cfg.Id), zap.Error(err))
				return
			}
		}
		nextTime, err := strconv.ParseUint(nextTimeStr, 10, 32)
		if err != nil {
			log.Error("JBattlePass init ParseUint err", zap.Uint32("id", cfg.Id), zap.String("nextTimeStr", nextTimeStr), zap.Error(err))
			return
		}
		if curTime > uint32(nextTime) {
			refreshWeeklyBattleField(cfg, nextWeekRefreshTime)
			return
		}
		log.Debug("JBattlePass init nextTime", zap.Uint32("cfgId", cfg.Id), zap.Any("nextTime", nextTime))
		cfg.SetNextWeekRefreshTime(uint32(nextTime))
		weekFieldBytes, err := rc.HGet(rCtx, key, template.BattlePassWeekField).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				refreshWeeklyBattleField(cfg, nextWeekRefreshTime)
				return
			} else {
				log.Error("JBattlePass init get weekField err", zap.Uint32("id", cfg.Id), zap.Error(err))
				return
			}
		}
		var weeklyField []int32
		log.Debug("JBattlePass init weekFieldBytes", zap.Uint32("cfgId", cfg.Id), zap.String("weekFieldBytes", weekFieldBytes))
		if err = json.Unmarshal([]byte(weekFieldBytes), &weeklyField); err != nil {
			log.Error("JBattlePass init Unmarshal weeklyField err", zap.Uint32("id", cfg.Id), zap.String("weekFieldBytes", weekFieldBytes), zap.Error(err))
		}
		log.Debug("JBattlePass init weeklyField", zap.Uint32("cfgId", cfg.Id), zap.Any("weeklyField", weeklyField))
		cfg.SetWeeklyField(weeklyField)
		if len(weeklyField) == 0 {
			cfg.SetWeeklyField(cfg.RandWeekField())
		}
	}
}

// 海边派对刷新周环境
func refreshWeeklyBattleField(j *template.JBattlePass, nextWeekRefreshTime uint32) {
	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()
	key := j.GetWeeklyBattleFieldRedisKey()
	_, err := rc.HSet(rCtx, key, template.BattlePassNextWeekRefreshTime, nextWeekRefreshTime).Result()
	if err != nil {
		log.Error("refreshWeeklyBattleField set nextWeekRefreshTime err", zap.Uint32("id", j.Id), zap.Uint32("nextWeekRefreshTime", nextWeekRefreshTime), zap.Error(err))
		return
	}
	weekField := j.RandWeekField()
	bytes, err := json.Marshal(weekField)
	if err != nil {
		log.Error("refreshWeeklyBattleField err", zap.Uint32("id", j.Id), zap.Any("weekField", weekField), zap.Error(err))
		return
	}
	_, err = rc.HSet(rCtx, key, template.BattlePassWeekField, string(bytes)).Result()
	if err != nil {
		log.Error("refreshWeeklyBattleField set weekField err", zap.Uint32("id", j.Id), zap.Any("weekField", weekField), zap.Error(err))
		return
	}
	j.SetNextWeekRefreshTime(nextWeekRefreshTime)
	j.SetWeeklyField(weekField)
}

func BattlePassRefreshTime() {
	curTime := tools.GetCurTime()
	nextWeekRefreshTime := tools.GetWeeklyRefreshTime(0)
	log.Debug("BattlePassTemplate RefreshTime", zap.Uint32("curTime", curTime), zap.Uint32("nextWeekRefreshTime", nextWeekRefreshTime))
	cfg := template.GetBattlePassTemplate().GetCurSeason(curTime)
	if cfg != nil {
		refreshWeeklyBattleField(cfg, nextWeekRefreshTime)
	}
}
