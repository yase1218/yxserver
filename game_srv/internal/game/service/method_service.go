package service

import (
	"gameserver/internal/game/event"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tools"
	"msg"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

func checkPreMission(p *player.Player, playMethod *model.PlayMethodData, missionConfig *template.JMission) msg.ErrCode {
	if len(playMethod.MissionData) == 0 {
		if missionConfig.PreId != 0 {
			return msg.ErrCode_FUNCTION_LOCK
		}
	} else {
		if missionConfig.PreId != playMethod.MissionData[0].MissionId {
			return msg.ErrCode_FUNCTION_LOCK
		}
	}
	return msg.ErrCode_SUCC
}

func checkDesertBattle(p *player.Player, playMethod *model.PlayMethodData, missionConfig *template.JMission) msg.ErrCode {
	act := getActivityByType(p, publicconst.Desert)
	if act == nil || !isValidActivity(act) {
		return msg.ErrCode_ACTIVITY_NOT_EXIST
	}

	if missionConfig.NextId == 0 &&
		len(playMethod.MissionData) > 0 &&
		playMethod.MissionData[0].MissionId == missionConfig.Id {
		return msg.ErrCode_SUCC
	}
	if err := checkPreMission(p, playMethod, missionConfig); err != msg.ErrCode_SUCC {
		return err
	}

	return msg.ErrCode_SUCC
}

func checkUnionBattle(p *player.Player, playMethod *model.PlayMethodData, missionConfig *template.JMission) msg.ErrCode {
	if missionConfig.NextId == 0 &&
		len(playMethod.MissionData) > 0 &&
		playMethod.MissionData[0].MissionId == missionConfig.Id {
		return msg.ErrCode_SUCC
	}
	if err := checkPreMission(p, playMethod, missionConfig); err != msg.ErrCode_SUCC {
		return err
	}
	return msg.ErrCode_SUCC
}

func PlayerMethodFightBeforeCheck(p *player.Player, stageCfg *template.JMission) msg.ErrCode {
	//// 解锁判定
	//if err := ServMgr.GetCommonService().CanLock(player, stageCfg.UnlockCond); err != msg.ErrCode_SUCC {
	//	return err
	//}
	//
	//// 玩法校验
	//playCfg := template.GetPlayMethodTemplate().GetPlayMethod(stageCfg.Type)
	//if playCfg == nil {
	//	log.Error("desert play cfg nil", zap.Int("type", stageCfg.Type))
	//	return msg.ErrCode_CONFIG_NIL
	//}
	//
	//// 解锁判定
	//if err := ServMgr.GetCommonService().CanLock(player, stageCfg.UnlockCond); err != msg.ErrCode_SUCC {
	//	return err
	//}
	//
	playMethod := getPlayMethod(p, stageCfg.Type)
	if playMethod == nil {
		log.Error("playMethod nil", zap.Uint64("accountId", p.GetUserId()), zap.Int("type", stageCfg.Type))
		return msg.ErrCode_PLAYMENTOD_NOT_OPEN
	}
	//if playMethod.TotalTimes == 0 {
	//	return msg.ErrCode_NO_ENOUGH_ITEM
	//}

	switch stageCfg.Type {
	case msg.BATTLETYPE_BATTLE_DESERT:
		if ret := checkDesertBattle(p, playMethod, stageCfg); ret != msg.ErrCode_SUCC {
			// return ret
			return msg.ErrCode_SUCC
		}
	case msg.BATTLETYPE_BATTLE_UNION:
		if ret := checkUnionBattle(p, playMethod, stageCfg); ret != msg.ErrCode_SUCC {
			return ret
		}

		UpdateTask(p, true, publicconst.TASK_COND_ALLIANCE_BOSS_FIGHT, 1)
	}

	AddPlayMethodTimes(p, stageCfg.Type, -1)

	p.UserData.BaseInfo.MissData = model.NewMissionData(stageCfg.Id, tools.GetCurTime())
	p.SaveBaseInfo()

	UpdateTask(p, true, publicconst.TASK_COND_BATTLE_MISSIONTYPE, uint32(stageCfg.Type), 1)
	UpdateTask(p, true, publicconst.TASK_COND_ANY_BATTLE, 1)
	UpdateMissionPoker(p, stageCfg.Id)

	return msg.ErrCode_SUCC
}

func PlayerMethodFightAfter(p *player.Player, ntf *msg.FsFightResultNtf) {
	log.Info("player method fight end", zap.Uint64("userId", p.GetUserId()), zap.Any("ntf", ntf))

	defer LeaveFight(p)

	stageCfg := template.GetMissionTemplate().GetMission(int(ntf.GetStageId()))
	if stageCfg == nil {
		log.Error("stageCfg nil", zap.Uint32("stageId", ntf.GetStageId()))
		return
	}

	playMethod := getPlayMethod(p, stageCfg.Type)
	if playMethod == nil {
		log.Error("playMethod nil", zap.Int("type", stageCfg.Type))
		return
	}

	var (
		killCnt   = uint32(ntf.GetKillCnt())
		fightTime = uint32(ntf.GetFightTime())

		mission = getMission(p, stageCfg.Type, stageCfg.Id)
		pbNtf   = &msg.PlayMethodFightResultNtf{
			Win:        ntf.GetVictory(),
			StageId:    ntf.GetStageId(),
			KillCnt:    killCnt,
			FightTime:  fightTime,
			RewardItem: nil,
			BattleType: msg.BattleType(stageCfg.Type),
			BossHurt:   ntf.GetBossHurt(),
		}
	)

	if mission == nil {
		mission = model.NewMission(stageCfg.Id, fightTime, killCnt, ntf.GetVictory())
		playMethod.MissionData = playMethod.MissionData[:0]
		playMethod.MissionData = append(playMethod.MissionData, mission)
	}

	if !mission.IsPass && ntf.GetVictory() {
		mission.IsPass = true
		mission.UpdateTime = tools.GetCurTime()
	}

	var (
		rewardItems []*model.SimpleItem
		notifyItems []uint32
	)
	switch stageCfg.Type {
	case msg.BATTLETYPE_BATTLE_DESERT:
		hurt := ntf.GetBossHurt()
		info := p.UserData.PlayMethod.Data
		var curMax int = 0
		for _, v := range info {
			if v.BtType == msg.BATTLETYPE_BATTLE_DESERT {
				if v.MaxDamage < int(hurt) {
					v.MaxDamage = int(hurt)
					curMax = v.MaxDamage
					break
				}
			}
		}
		p.UserData.PlayMethod.Data = info
		UpdateCommonRankInfo(p, uint32(curMax), template.WeekPassBlackBoosDamage)
		p.SavePlayMethod()
		fallthrough
	case msg.BATTLETYPE_BATTLE_UNION:
		// 联盟 || 沙漠boss奖励按血量百分比算
		totalHp := stageCfg.GetBossInitHp()
		totalHurt := uint64(ntf.GetBossHurt()) * 100
		percentage := int(totalHurt / totalHp)

		progressRewardSlice := stageCfg.GetProgressReward(percentage)
		progressRewardMap := make(map[int]int, len(progressRewardSlice))
		for _, v := range progressRewardSlice {
			progressRewardMap[int(v.ItemId)] += int(v.ItemNum)
		}

		for k, v := range progressRewardMap {
			addItems := AddItem(p.GetUserId(), uint32(k), int32(v), publicconst.PlayMethodAddItem, false)
			rewardItems = append(rewardItems, addItems...)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		}
	}

	// todo 联盟boss伤害排行
	//if stageCfg.Type == msg.BATTLETYPE_BATTLE_UNION {
	//	member, err := dao.GetMember(player.GetAccountId())
	//	if err == nil && member != nil {
	//		if err = ServMgr.GetRankService().AddAllianceRank(member.AllianceID,
	//			msg.AllianceRankType_Alliance_Rank_Boss_Single, player.GetAccountId(), ntf.GetBossHurt()); err != nil {
	//			log.Error("add alliance rank err", zap.Error(err), zap.Int64("accountId", player.GetAccountId()), zap.Uint32("allianceId", member.AllianceID))
	//		}
	//
	//		if err = ServMgr.GetRankService().AddAllianceRank(member.AllianceID,
	//			msg.AllianceRankType_Alliance_Rank_Boss, player.GetAccountId(), ntf.GetBossHurt()); err != nil {
	//			log.Error("add alliance rank err", zap.Error(err), zap.Int64("accountId", player.GetAccountId()), zap.Uint32("allianceId", member.AllianceID))
	//		}
	//
	//		if player.AccountInfo.AllianceBossHurtMax < int64(ntf.GetBossHurt()) {
	//			player.AccountInfo.AllianceBossHurtMax = int64(ntf.GetBossHurt())
	//			dao.AccountDao.UpdateAllianceBossHurtMax(player.GetAccountId(), int64(ntf.GetBossHurt()))
	//		}
	//	}
	//}

	if ntf.GetVictory() {
		notifyMsg := &msg.NotifyPlayMethondDataChange{}
		notifyMsg.Data = append(notifyMsg.Data, ToProtocolPlayMethodData(playMethod, p.UserData.Fight.Faction))
		p.SendNotify(notifyMsg)

		event.EventMgr.PublishEvent(event.NewPassMissionEvent(p, stageCfg.Id, ListenPassMissionEvent)) // 通关事件
	}

	fightEndUpdateTask(p, ntf) // 结束关卡任务相关
	p.UserData.BaseInfo.MissData.MissionId = 0
	p.UserData.BaseInfo.MissData.StartTime = 0
	p.UserData.BaseInfo.BattleData = nil
	p.SaveBaseInfo()

	pbNtf.RewardItem = ToProtocolSimpleItems(rewardItems)
	p.SendNotify(pbNtf)
}

func IsOpen(p *player.Player, bt msg.BattleType) msg.ErrCode {
	return msg.ErrCode_SUCC
}

// UpdateWeapon
func UpdateWeapon(p *player.Player, btType int, weaponIds []uint32) msg.ErrCode {
	if err := IsOpen(p, msg.BattleType(btType)); err != msg.ErrCode_SUCC {
		return err
	}

	idLen := len(weaponIds)
	if idLen > 5 {
		return msg.ErrCode_INVALID_DATA
	}

	for i := 0; i < len(weaponIds); i++ {
		if weaponIds[i] > 0 {
			if weapon := getWeapon(p, weaponIds[i]); weapon == nil {
				return msg.ErrCode_WEAPON_NOT_EXIST
			}
		}
	}

	playMethod := getPlayMethod(p, btType)
	if playMethod == nil {
		return msg.ErrCode_FUNCTION_NOT_OPEN
	}

	if tools.ListUint32Equal(playMethod.WeaponIds, weaponIds) {
		return msg.ErrCode_SUCC
	}

	playMethod.WeaponIds = weaponIds
	p.SavePlayMethod()

	return msg.ErrCode_SUCC
}

func GetFightWeapons(p *player.Player) []uint32 {
	res := make([]uint32, 0)
	playMethod := getPlayMethod(p, int(p.FightType))
	if playMethod == nil {
		return res
	}
	return playMethod.WeaponIds
}

// PlayMethodStartBattle 开始战斗
func PlayMethodStartBattle(p *player.Player, missionId int) msg.ErrCode {
	mission := template.GetMissionTemplate().GetMission(missionId)
	if mission == nil {
		return msg.ErrCode_MISSION_NOT_EXIST
	}

	btType := msg.BattleType(mission.Type)

	// 解锁判定
	if err := CanLock(p, mission.UnlockCond); err != msg.ErrCode_SUCC {
		return err
	}

	playMethodConfig := template.GetPlayMethodTemplate().GetPlayMethod(int(btType))
	if playMethodConfig == nil {
		return msg.ErrCode_INVALID_DATA
	}

	playMethod := getPlayMethod(p, mission.Type)
	if playMethod == nil {
		return msg.ErrCode_PLAYMENTOD_NOT_OPEN
	}

	if btType == msg.BattleType_Battle_Challenge {
		return msg.ErrCode_INVALID_DATA
	}

	if btType == msg.BattleType_Battle_Coin ||
		btType == msg.BattleType_Battle_Equip ||
		btType == msg.BattleType_Battle_Weapon {
		if err := checkPreMission(p, playMethod, mission); err != msg.ErrCode_SUCC {
			return err
		}
	} else if btType == msg.BattleType_Battle_Desert {
		if ret := checkDesertBattle(p, playMethod, mission); ret != msg.ErrCode_SUCC {
			return ret
		}
	} else if btType == msg.BattleType_Battle_Union {
		if ret := checkUnionBattle(p, playMethod, mission); ret != msg.ErrCode_SUCC {
			return ret
		}
	}

	if playMethod.TotalTimes == 0 {
		return msg.ErrCode_NO_ENOUGH_ITEM
	}

	AddPlayMethodTimes(p, mission.Type, -1)

	// 更新当前关卡数据
	p.UserData.BaseInfo.MissData = model.NewMissionData(missionId, tools.GetCurTime())
	p.SaveBaseInfo()

	UpdateTask(p, true, publicconst.TASK_COND_BATTLE_MISSIONTYPE, uint32(mission.Type), 1)

	if btType == msg.BattleType_Battle_Union {
		UpdateTask(p, true,
			publicconst.TASK_COND_ALLIANCE_BOSS_FIGHT, 1)
	}

	UpdateTask(p, true, publicconst.TASK_COND_ANY_BATTLE, 1)
	UpdateMissionPoker(p, missionId)
	return msg.ErrCode_SUCC
}

func CheckNewMethod(p *player.Player) {
	lstPlayMethod := template.GetPlayMethodTemplate().GetAllPlayMethod()
	// 比对是否新增了玩法
	var newPlayMethod []*template.JPlayMethod
	for i := 0; i < len(lstPlayMethod); i++ {
		exist := false
		for k := 0; k < len(p.UserData.PlayMethod.Data); k++ {
			if p.UserData.PlayMethod.Data[k].BtType == lstPlayMethod[i].Type {
				exist = true
				break
			}
		}
		if !exist {
			newPlayMethod = append(newPlayMethod, lstPlayMethod[i])
		}
	}
	var newPlayMethodModel []*model.PlayMethodData
	for i := 0; i < len(newPlayMethod); i++ {
		newPlayMethodModel = append(newPlayMethodModel,
			model.NewPlayMethodInfo(newPlayMethod[i].Type, newPlayMethod[i].Limit))
	}
	if len(newPlayMethodModel) > 0 {
		p.UserData.PlayMethod.Data = append(p.UserData.PlayMethod.Data, newPlayMethodModel...)
		p.SavePlayMethod()
	}

	// var isFirstLogin uint32 = 0
	// if tools.GetCurTime() >= p.UserData.PlayMethod.NextRefreshTime {
	// 	isFirstLogin = 1
	// }
	RefreshPlayMethod(p, false, false)
}

func GmRefreshPlayMethod(p *player.Player) {
	RefreshPlayMethod(p, true, true)
}

func RefreshPlayMethod(p *player.Player, notifyClient, force bool) {
	if p.UserData.PlayMethod == nil {
		return
	}
	curTime := tools.GetCurTime()
	if !force && curTime < p.UserData.PlayMethod.NextRefreshTime {
		return
	}

	for i := 0; i < len(p.UserData.PlayMethod.Data); i++ {
		// 刷新道具
		if config := template.GetPlayMethodTemplate().GetPlayMethod(p.UserData.PlayMethod.Data[i].BtType); config != nil {
			if config.Limit > 0 && p.UserData.PlayMethod.Data[i].TotalTimes < config.Limit {
				times := GetMonthCardSweepTimes(p)
				p.UserData.PlayMethod.Data[i].TotalTimes = config.Limit + times
			}
		}
	}

	p.UserData.PlayMethod.NextRefreshTime = tools.GetDailyRefreshTime()
	p.SavePlayMethod()

	if notifyClient {
		p.SendNotify(&msg.NotifyPlayMethondDataChange{
			Data: ToProtocolPlayMethodDatas(p.UserData.PlayMethod.Data, p.UserData.Fight.Faction),
		})
	}
}

func getMission(p *player.Player, bt, missionId int) *model.Mission {
	var missions []*model.Mission
	for i := 0; i < len(p.UserData.PlayMethod.Data); i++ {
		if p.UserData.PlayMethod.Data[i].BtType == bt {
			missions = p.UserData.PlayMethod.Data[i].MissionData
			break
		}
	}
	for i := 0; i < len(missions); i++ {
		if missions[i].MissionId == missionId {
			return missions[i]
		}
	}
	return nil
}

func ToProtocolPlayMethodData(data *model.PlayMethodData, faction msg.FactionType) *msg.PlayMethondData {
	ret := &msg.PlayMethondData{
		BtType:     msg.BattleType(data.BtType),
		TotalTimes: uint32(data.TotalTimes),
		Type:       faction,
	}
	for i := 0; i < len(data.WeaponIds); i++ {
		ret.WeaponIds = append(ret.WeaponIds, data.WeaponIds[i])
	}

	for i := 0; i < len(data.MissionData); i++ {
		ret.Missions = append(ret.Missions, ToProtocolMission(data.MissionData[i]))
	}
	return ret
}

func ToProtocolPlayMethodDatas(data []*model.PlayMethodData, factionType msg.FactionType) []*msg.PlayMethondData {
	var ret []*msg.PlayMethondData
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolPlayMethodData(data[i], factionType))
	}
	return ret
}
func getPlayMethod(p *player.Player, btType int) *model.PlayMethodData {
	for i := 0; i < len(p.UserData.PlayMethod.Data); i++ {
		if p.UserData.PlayMethod.Data[i].BtType == btType {
			return p.UserData.PlayMethod.Data[i]
		}
	}
	return nil
}

// AddPlayMethodTimes 添加玩法次数
func AddPlayMethodTimes(p *player.Player, byType int, times int) {
	data := getPlayMethod(p, byType)
	if data == nil {
		log.Error("AddPlayMethodTimes err", zap.Uint64("accountId", p.GetUserId()), zap.Int("byType", byType))
		return
	}

	temp := data.TotalTimes + times
	data.TotalTimes = temp
	p.SavePlayMethod()

	p.SendNotify(&msg.NotifyPlayMethondDataChange{
		Data: []*msg.PlayMethondData{ToProtocolPlayMethodData(data, p.UserData.Fight.Faction)},
	})
}
