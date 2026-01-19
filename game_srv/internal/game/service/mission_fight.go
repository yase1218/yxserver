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

/*
 * MainFightBeforeCheck
 *  @Description: 主线战斗战前检查
 *  @param player
 *  @param stageCfg
 *  @return msg.ErrCode
 */
func MainFightBeforeCheck(p *player.Player, stageCfg *template.JMission) msg.ErrCode {
	log.Debug("main fight before", zap.Uint64("uid", p.GetUserId()), zap.Int("fightStageId", stageCfg.Id))
	// 没有足够的体力
	if !EnoughItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_AP), uint32(stageCfg.PowerCost)) {
		return msg.ErrCode_MISSION_NO_ENOUGH_AP
	}

	// 解锁判定
	if err := CanLock(p, stageCfg.UnlockCond); err != msg.ErrCode_SUCC {
		return err
	}

	// 校验关卡
	if err := checkMission(stageCfg.Id, p); err != msg.ErrCode_SUCC {
		return err
	}

	// 扣除体力
	if stageCfg.PowerCost > 0 {
		if p.UserData.StageInfo.MissionId > 0 { // 1-1通关后才开始消耗体力
			if res := CostItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_AP), uint32(stageCfg.PowerCost), publicconst.MissionCostItem, true); res != msg.ErrCode_SUCC {
				return res
			}
		}
	}

	// 关卡解锁配件
	UpdateMissionTreasure(p, uint32(stageCfg.Id))
	PassMissionAddWeapon(p, stageCfg.Id)

	// 更新关卡数据
	p.UserData.BaseInfo.MissData = model.NewMissionData(stageCfg.Id, tools.GetCurTime())
	p.SaveBaseInfo()

	// 更新扑克数据
	UpdateMissionPoker(p, stageCfg.Id)

	// 更新任务
	UpdateTask(p, true, publicconst.TASK_COND_BATTLE_MISSIONTYPE, uint32(stageCfg.Type), 1)
	UpdateTask(p, true, publicconst.TASK_COND_ANY_BATTLE, 1)

	return msg.ErrCode_SUCC
}

/*
 * MainFightAfter
 *  @Description: 主线战斗战斗结束
 */
func MainFightAfter(p *player.Player, ntf *msg.FsFightResultNtf) {
	log.Info("main fight end", zap.Uint64("accountId", p.GetUserId()), zap.Any("ntf", ntf))

	defer LeaveFight(p)

	stageCfg := template.GetMissionTemplate().GetMission(int(ntf.GetStageId()))
	if stageCfg == nil {
		log.Error("stageCfg nil", zap.Uint32("stageId", ntf.GetStageId()))
		return
	}

	var (
		updateMission       bool
		refreshCompleteTime bool
		notifyItems         []uint32
		rewardItems         []*model.SimpleItem

		killCnt   = uint32(ntf.GetKillCnt())
		fightTime = uint32(ntf.GetFightTime())

		pbNtf = &msg.MainFightResultNtf{
			Win:        ntf.GetVictory(),
			StageId:    ntf.GetStageId(),
			KillCnt:    killCnt,
			FightTime:  fightTime,
			BoxState:   0,
			MaxStageId: 0,
			FailType:   ntf.GetReason(),
			RewardItem: nil,
		}
	)

	// get stage db data
	mission := findMission(p, int(ntf.GetStageId()), true)
	if mission == nil {
		completeTime := fightTime
		if fightTime == 0 {
			completeTime = 1000
		}
		mission = model.NewMission(int(ntf.GetStageId()), completeTime, killCnt, ntf.GetVictory())
		addMission(p, mission, true)

		if ntf.GetVictory() {
			if fightTime == 0 {
				mission.CompleteTime = 1000
			}
			refreshCompleteTime = true
		}
	}

	// update kill monster num
	if killCnt > mission.KillMonsterNum {
		mission.KillMonsterNum = killCnt
		updateMission = true
	}

	// update complete time
	if fightTime < mission.CompleteTime {
		mission.CompleteTime = fightTime
		updateMission = true
		refreshCompleteTime = true
	}

	// update pass mission
	if !mission.IsPass && ntf.GetVictory() {
		mission.IsPass = true
		updateMission = true
	}

	// update box state
	isUpdateMissionData, addStar := updateBoxState(stageCfg, mission, ntf)
	if !updateMission && isUpdateMissionData {
		updateMission = true
	}

	if addStar > 0 {
		p.UserData.StageInfo.StageStar += addStar
		p.SaveStageInfo()
		//  TODO 排行
		//ServMgr.GetRankService().updateStageStarRank(p)
		UpdateCommonRankInfo(p, p.UserData.StageInfo.StageStar, template.MainStarRank)
	}

	if ntf.GetVictory() {
		// 当前主线通关需要更新主线结算挂机奖励 主线关卡变化
		if p.UserData.StageInfo.MissionId == 0 ||
			template.GetMissionTemplate().IsNextMission(int(p.UserData.StageInfo.MissionId), int(p.UserData.BaseInfo.MissData.MissionId)) {
			// 结算挂机奖励
			if p.UserData.StageInfo.MissionId == 0 {
				p.UserData.BaseInfo.HookData.StartTime = (tools.GetCurTime() - 10*60)
				p.SaveBaseInfo()
			}
			settleOnHook(p, p.UserData.BaseInfo.MissData.MissionId)

			// 合约任务
			if p.UserData.StageInfo.MissionId < p.UserData.BaseInfo.MissData.MissionId &&
				p.UserData.Contract.TaskId != 0 &&
				p.UserData.Contract.TaskType == msg.ConditionType_Condition_Contract_Rand {
				FinishContract(p, 1)
			}

			p.UserData.StageInfo.MissionId = p.UserData.BaseInfo.MissData.MissionId
			p.SaveStageInfo()
			// 添加关卡流失统计
			// dao.UserStaticDao.AddLossMission(model.NewLossMission(tools.GetStaticTime(tools.GetCurTime()),
			// 	uint32(p.GetUserId()), p.AccountInfo.MissionId))

			// tda update mission
			// tdaData := &tda.CommonUser{
			// 	Max_battle_id: strconv.Itoa(int(p.AccountInfo.MissionId)),
			// }
			// tda.TdaUpdateCommonUser(p.TdaCommonAttr.AccountId, p.TdaCommonAttr.DistinctId, tdaData)
		}

		// 通关奖励
		for k, v := range stageCfg.NormalReward {
			num := v * GetMonthCardMainScale(p) / 100
			addItems := AddItem(p.GetUserId(), uint32(k), int32(num), publicconst.PassMissionAddItem, false)
			rewardItems = append(rewardItems, addItems...)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		}

		// TODO 排行榜 更新排名
		if refreshCompleteTime {
			args := make([]float64, 0)
			args = append(args, float64(mission.MissionId))
			args = append(args, float64(mission.CompleteTime))
			UpdateCommonRankInfo(p, args, template.MainNormalRank)
			// 	ServMgr.GetRankService().updateMissionRankMission(p, mission)
			// 	ServMgr.GetRankService().updateSpecialMissionRank(p, mission)
		}

		// 触发调查问卷
		trigerQuestion(p, p.UserData.StageInfo.MissionId)

		// 记录首通信息
		RecordFirstPassPlayerInfo(p, mission.MissionId, template.MainNormalRank)

		// 通关事件
		event.EventMgr.PublishEvent(event.NewPassMissionEvent(p, p.UserData.StageInfo.MissionId, ListenPassMissionEvent))
		UpdateTask(p, true, publicconst.TASK_COND_PASS_TYPE_MISSION, uint32(stageCfg.Type), 1) // 完成XX关卡类型XX次
	} else {
		// 失败奖励
		scale := fightTime / 60
		scale = tools.LimitUint32(scale, uint32(stageCfg.RewardQuitLimit))
		if scale > 0 {
			for k, v := range stageCfg.ExitReward {
				finalScale := float64(scale) * float64(GetMonthCardMainScale(p)) / float64(100)
				num := uint32(v * finalScale)
				if num > 0 {
					addItems := AddItem(p.GetUserId(), uint32(k), int32(num), publicconst.PassMissionAddItem, false)
					rewardItems = append(rewardItems, addItems...)
					notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
				}
			}
		}
	}

	if len(notifyItems) > 0 {
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	if updateMission {
		p.SaveMission()
	}

	// write statics log
	//para := fmt.Sprintf("missionId:%v,start:%v end:%v pass:%v|", p.UserData.BaseInfo.MissData.MissionId,
	//	p.UserData.BaseInfo.MissData.StartTime, tools.GetCurTime(), ntf.GetVictory())
	//for i := 0; i < len(rewardItems); i++ {
	//	para += fmt.Sprintf("%v,%v|", rewardItems[i].Id, rewardItems[i].Num)
	//}

	//ServMgr.GetCommonService().AddStaticsData(p, publicconst.Statics_Main_Mission_Id, para)

	pbNtf.BoxState = mission.BoxState
	pbNtf.MaxStageId = uint32(p.UserData.StageInfo.MissionId)
	pbNtf.RewardItem = ToProtocolSimpleItems(rewardItems)

	p.ResetMissData()

	fightEndUpdateTask(p, ntf)

	p.UserData.BaseInfo.BattleData = nil

	p.SaveBaseInfo()
	p.SendNotify(pbNtf)
}

func ContractFightAfter(p *player.Player, ntf *msg.FsFightResultNtf) {
	log.Info("main fight end", zap.Uint64("accountId", p.GetUserId()), zap.Any("ntf", ntf))

	defer LeaveFight(p)

	stageCfg := template.GetMissionTemplate().GetMission(int(ntf.GetStageId()))
	if stageCfg == nil {
		log.Error("stageCfg nil", zap.Uint32("stageId", ntf.GetStageId()))
		return
	}

	var (
		notifyItems []uint32
		rewardItems []*model.SimpleItem

		killCnt   = uint32(ntf.GetKillCnt())
		fightTime = uint32(ntf.GetFightTime())

		pbNtf = &msg.MainFightResultNtf{
			Win:        ntf.GetVictory(),
			StageId:    ntf.GetStageId(),
			KillCnt:    killCnt,
			FightTime:  fightTime,
			BoxState:   0,
			MaxStageId: 0,
			FailType:   ntf.GetReason(),
			RewardItem: nil,
		}
	)

	if ntf.GetVictory() {
		if _, ok := p.UserData.WeekPass.ContractInfo[pbNtf.StageId]; !ok {
			// 通关奖励
			for k, v := range stageCfg.NormalReward {
				num := v * GetMonthCardMainScale(p) / 100
				addItems := AddItem(p.GetUserId(), uint32(k), int32(num), publicconst.PassMissionAddItem, false)
				rewardItems = append(rewardItems, addItems...)
				notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
			}

			p.UserData.WeekPass.ContractInfo[pbNtf.StageId] = true
			p.SaveWeekPass()
		}
	}

	if len(notifyItems) > 0 {
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	pbNtf.RewardItem = ToProtocolSimpleItems(rewardItems)

	p.SendNotify(pbNtf)
}

func SecretFightAfter(p *player.Player, ntf *msg.FsFightResultNtf) {
	log.Info("main fight end", zap.Uint64("accountId", p.GetUserId()), zap.Any("ntf", ntf))

	defer LeaveFight(p)

	stageCfg := template.GetMissionTemplate().GetMission(int(ntf.GetStageId()))
	if stageCfg == nil {
		log.Error("stageCfg nil", zap.Uint32("stageId", ntf.GetStageId()))
		return
	}

	var (
		notifyItems []uint32
		killCnt     = uint32(ntf.GetKillCnt())
		fightTime   = uint32(ntf.GetFightTime())

		pbNtf = &msg.MainFightResultNtf{
			Win:        ntf.GetVictory(),
			StageId:    ntf.GetStageId(),
			KillCnt:    killCnt,
			FightTime:  fightTime,
			BoxState:   0,
			MaxStageId: 0,
			FailType:   ntf.GetReason(),
			RewardItem: nil,
		}
	)

	isUpdate, _, rewards := updateSecretBoxState(stageCfg, &p.UserData.WeekPass.SecretBoxState, ntf)

	if ntf.GetVictory() {
		p.UserData.WeekPass.SecretCount++
		for _, v := range rewards {
			addItems := AddItem(p.GetUserId(), uint32(v.Id), int32(v.Num), publicconst.PassMissionAddItem, false)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		}
	}

	if len(notifyItems) > 0 {
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	if ntf.GetVictory() && isUpdate {
		p.SaveWeekPass()
	}

	pbNtf.BoxState = p.UserData.WeekPass.SecretBoxState
	pbNtf.RewardItem = ToProtocolSimpleItems(rewards)

	p.SendNotify(pbNtf)
}

/*
 * ChallengeFightBeforeCheck
 *  @Description: 挑战关卡战前检查
 *  @param player
 *  @param stageCfg
 *  @return msg.ErrCode
 */
func ChallengeFightBeforeCheck(p *player.Player, stageCfg *template.JMission) msg.ErrCode {
	log.Debug("challenge fight before", zap.Uint64("uid", p.GetUserId()), zap.Int("fightStageId", stageCfg.Id))
	// 没有足够的体力
	if !EnoughItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_AP), uint32(stageCfg.PowerCost)) {
		return msg.ErrCode_MISSION_NO_ENOUGH_AP
	}

	// 解锁判定
	if err := CanLock(p, stageCfg.UnlockCond); err != msg.ErrCode_SUCC {
		return err
	}

	// 扣除体力
	if stageCfg.PowerCost > 0 {
		if res := CostItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_AP), uint32(stageCfg.PowerCost), publicconst.MissionCostItem, true); res != msg.ErrCode_SUCC {
			return res
		}
	}

	// 更新关卡数据
	p.UserData.BaseInfo.MissData = model.NewMissionData(stageCfg.Id, tools.GetCurTime())
	p.SaveBaseInfo()

	// 更新扑克数据
	UpdateMissionPoker(p, stageCfg.Id)

	// 更新任务
	UpdateTask(p, true, publicconst.TASK_COND_BATTLE_MISSIONTYPE, uint32(stageCfg.Type), 1)
	UpdateTask(p, true, publicconst.TASK_COND_ANY_BATTLE, 1)

	return msg.ErrCode_SUCC
}

/*
 * ChallengeFightAfter
 *  @Description: 挑战关卡结算
 *  @param player
 *  @param ntf
 */
func ChallengeFightAfter(p *player.Player, ntf *msg.FsFightResultNtf) {
	log.Info("challenge fight end", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("stageId", ntf.StageId), zap.Any("ntf", ntf))
	defer LeaveFight(p)

	stageCfg := template.GetMissionTemplate().GetMission(int(ntf.GetStageId()))
	if stageCfg == nil {
		log.Error("stageCfg nil", zap.Uint32("stageId", ntf.GetStageId()))
		return
	}

	var (
		updateMission       bool
		refreshCompleteTime bool
		notifyItems         []uint32
		rewardItems         []*model.SimpleItem

		killCnt   = uint32(ntf.GetKillCnt())
		fightTime = uint32(ntf.GetFightTime())

		pbNtf = &msg.MainFightResultNtf{
			Win:        ntf.GetVictory(),
			StageId:    ntf.GetStageId(),
			KillCnt:    killCnt,
			FightTime:  fightTime,
			BoxState:   0,
			MaxStageId: 0,
			FailType:   ntf.GetReason(),
			RewardItem: nil,
		}
	)

	// get stage db data
	mission := findMission(p, int(ntf.GetStageId()), false)
	if mission == nil {
		completeTime := fightTime
		if fightTime == 0 {
			completeTime = 1000
		}
		mission = model.NewMission(int(ntf.GetStageId()), completeTime, killCnt, ntf.GetVictory())
		addMission(p, mission, false)

		if ntf.GetVictory() {
			if fightTime == 0 {
				mission.CompleteTime = 1000
			}
			refreshCompleteTime = true
		}
	}

	// update kill monster num
	if killCnt > mission.KillMonsterNum {
		mission.KillMonsterNum = killCnt
		updateMission = true
	}

	// update complete time
	if fightTime < mission.CompleteTime {
		mission.CompleteTime = fightTime
		updateMission = true
		refreshCompleteTime = true
	}

	// update pass mission
	if !mission.IsPass && ntf.GetVictory() {
		mission.IsPass = true
		updateMission = true
	}

	// update box state
	isUpdateMissionData, addStar := updateBoxState(stageCfg, mission, ntf)
	if !updateMission && isUpdateMissionData {
		updateMission = true
	}

	if addStar > 0 {
		p.UserData.StageInfo.StageStar += addStar
		p.SaveStageInfo()

		// TODO 排行榜
		// ServMgr.GetRankService().updateStageStarRank(p)
	}

	if ntf.GetVictory() {
		// 合约任务
		if p.UserData.StageInfo.MissionId < p.UserData.BaseInfo.MissData.MissionId &&
			p.UserData.Contract.TaskId != 0 &&
			p.UserData.Contract.TaskType == msg.ConditionType_Condition_Contract_Rand {
			FinishContract(p, 1)
		}

		if p.UserData.StageInfo.StageFirstPass == nil {
			p.UserData.StageInfo.StageFirstPass = make([]int, 0)
		}
		if !tools.ListIntContain(p.UserData.StageInfo.StageFirstPass, mission.MissionId) {
			p.UserData.StageInfo.StageFirstPass = append(p.UserData.StageInfo.StageFirstPass, mission.MissionId)
			p.SaveStageInfo()
		}

		// 通关奖励
		for k, v := range stageCfg.NormalReward {
			num := v * GetMonthCardMainScale(p) / 100
			addItems := AddItem(p.GetUserId(), uint32(k), int32(num), publicconst.PassMissionAddItem, false)
			rewardItems = append(rewardItems, addItems...)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		}

		// 更新排名 TODO 排行榜
		if refreshCompleteTime {
			args := make([]float64, 0)
			args = append(args, float64(mission.MissionId), float64(mission.CompleteTime))
			UpdateCommonRankInfo(p, args, template.MainEliteRank)
			// ServMgr.GetRankService().updateMissionRankMission(p, mission)
			// ServMgr.GetRankService().updateSpecialMissionRank(p, mission)
		}

		// 首通记录
		RecordFirstPassPlayerInfo(p, mission.MissionId, template.MainEliteRank)

		// 触发调查问卷
		trigerQuestion(p, p.UserData.StageInfo.MissionId)

		// 任务
		UpdateTask(p, true, publicconst.TASK_COND_FIRST_PASS_CHALLENGE_MISSION, 1)
		processHistoryData(p, publicconst.TASK_COND_FIRST_PASS_CHALLENGE_MISSION, 0, 1)
		UpdateTask(p, true, publicconst.TASK_COND_PASS_TYPE_MISSION, uint32(stageCfg.Type), 1) // 完成XX关卡类型XX次

		// 通关事件
		event.EventMgr.PublishEvent(event.NewPassMissionEvent(p, p.UserData.BaseInfo.MissData.MissionId, ListenPassMissionEvent))
	} else {
		// 失败奖励
		scale := fightTime / 60
		scale = tools.LimitUint32(scale, uint32(stageCfg.RewardQuitLimit))
		if scale > 0 {
			for k, v := range stageCfg.ExitReward {
				finalScale := float64(scale) * float64(GetMonthCardMainScale(p)) / float64(100)
				num := uint32(v * finalScale)
				if num > 0 {
					addItems := AddItem(p.GetUserId(), uint32(k), int32(num), publicconst.PassMissionAddItem, false)
					rewardItems = append(rewardItems, addItems...)
					notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
				}
			}
		}
	}

	if len(notifyItems) > 0 {
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	if updateMission {
		p.SaveMission()
	}

	// write statics log
	// para := fmt.Sprintf("missionId:%v,start:%v end:%v pass:%v|", p.UserData.BaseInfo.MissData.MissionId,
	// 	p.UserData.BaseInfo.MissData.StartTime, tools.GetCurTime(), ntf.GetVictory())
	// for i := 0; i < len(rewardItems); i++ {
	// 	para += fmt.Sprintf("%v,%v|", rewardItems[i].Id, rewardItems[i].Num)
	// }
	// ServMgr.GetCommonService().AddStaticsData(p, publicconst.Statics_Challege_Mission_Id, para)

	pbNtf.BoxState = mission.BoxState
	pbNtf.MaxStageId = uint32(tools.FindMaxElemForArray(p.UserData.StageInfo.StageFirstPass))
	pbNtf.RewardItem = ToProtocolSimpleItems(rewardItems)

	p.ResetMissData()

	fightEndUpdateTask(p, ntf)

	p.UserData.BaseInfo.BattleData = nil
	p.SaveBaseInfo()

	p.SendNotify(pbNtf)
}

/*
 * updateBoxState
 *  @Description: 优化后的更新宝箱状态
 *  @param stageCfg	关卡配置
 *  @param mission	关卡数据
 *  @param ntf		战斗结果通知
 *  @return bool	是否更新关卡数据
 *  @return uint32	加星数
 */
func updateBoxState(stageCfg *template.JMission, mission *model.Mission, ntf *msg.FsFightResultNtf) (bool, uint32) {
	var (
		update  = false
		addStar = uint32(0)
	)

	for k, v := range stageCfg.RewardComplete {
		if !tools.GetBit(mission.BoxState, uint32(k)) {
			ruleBoxCfg := template.GetRuleBoxTemplate().GetRuleBox(v)
			if ruleBoxCfg == nil {
				log.Error("ruleBoxCfg is nil", zap.Int("ruleBoxId", v))
				break
			}

			checkPass := true
			for _, vv := range ruleBoxCfg.Conds {
				if !boxCondition(vv, ntf) {
					checkPass = false
					break
				}
			}
			if !checkPass {
				break
			}

			addStar++
			update = true

			tools.SetBit(&mission.BoxState, uint32(k))
		}
	}
	return update, addStar
}

/*
 * boxCondition
 *  @Description: 优化后的宝箱校验
 *  @param para	宝箱条件
 *  @param ntf	战斗结果通知
 *  @return bool
 */
func boxCondition(para *template.BoxCond, ntf *msg.FsFightResultNtf) bool {
	switch para.Cond {
	case template.BattleTime:
		// 战斗时长
		if len(para.Args) >= 1 {
			if ntf.GetFightTime()/1000 >= int64(para.Args[0]) {
				return true
			}
		}
	case template.KillBossId:
		// 杀死指定的怪物
		for i := 0; i < len(para.Args); i += 1 {
			bossId := uint32(para.Args[i])
			if _, ok := ntf.KillMonsters[bossId]; !ok {
				return false
			}
		}
		return true

	case template.PASSMISSION:
		// 通关关卡
		if ntf.GetVictory() {
			return true
		}
	case template.ShipHp:
		// 剩余血量
		if len(para.Args) >= 1 && ntf.GetVictory() {
			return ntf.GetHpRate() >= uint32(para.Args[0])
		}
	}
	return false
}

/*
 * mainFightEnd
 *  @Description: 主线战斗结束 更新任务相关
 *  @param playerData
 *  @param ntf
 */
func fightEndUpdateTask(p *player.Player, ntf *msg.FsFightResultNtf) {
	if ntf.KillCnt > 0 {
		UpdateTask(p, true, publicconst.TASK_COND_KILL_MONSTER_NUM, uint32(ntf.KillCnt))
		processHistoryData(p, publicconst.TASK_COND_KILL_MONSTER_NUM, 0, uint32(ntf.KillCnt))
	}

	if ntf.GetChip() > 0 {
		UpdateTask(p, true, publicconst.TASK_COND_GET_CHIP_NUM, ntf.GetChip())
	}

	if len(ntf.GetPokerTypes()) > 0 {
		args := make([]uint32, 0, len(ntf.GetPokerTypes())*2)
		for k, v := range ntf.GetPokerTypes() {
			args = append(args, k, v)
		}
		UpdateTask(p, true, publicconst.TASK_COND_POKER_NUM, args...)
	}

	if len(ntf.GetKillMonsters()) > 0 {
		UpdateTask(p, true, publicconst.TASK_COND_KILL_ELITE_MONSTER_BOSS_NUM, uint32(len(ntf.GetKillMonsters())))
		processHistoryData(p, publicconst.TASK_COND_KILL_ELITE_MONSTER_BOSS_NUM, 0, uint32(len(ntf.GetKillMonsters())))
	}
}

func updateSecretBoxState(stageCfg *template.JMission, boxState *uint32, ntf *msg.FsFightResultNtf) (bool, uint32, []*model.SimpleItem) {
	var (
		update  = false
		addStar = uint32(0)
		rewards = make([]*model.SimpleItem, 0)
	)

	for k, v := range stageCfg.RewardComplete {
		if !tools.GetBit(*boxState, uint32(k)) {
			ruleBoxCfg := template.GetRuleBoxTemplate().GetRuleBox(v)
			if ruleBoxCfg == nil {
				log.Error("ruleBoxCfg is nil", zap.Int("ruleBoxId", v))
				break
			}

			checkPass := true
			for _, vv := range ruleBoxCfg.Conds {
				if !boxCondition(vv, ntf) {
					checkPass = false
					break
				}
			}
			if !checkPass {
				continue
			}

			addStar++
			update = true

			for i := range ruleBoxCfg.Reward {
				item := ruleBoxCfg.Reward[i]
				rewards = append(rewards, &model.SimpleItem{
					Id:  item.ItemId,
					Num: item.ItemNum,
					Src: item.Src,
				})
			}

			tools.SetBitFlag(boxState, uint32(k), uint32(publicconst.CanGet))
		}
	}
	return update, addStar, rewards
}
