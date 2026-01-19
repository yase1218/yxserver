package service

import (
	"gameserver/internal/game/event"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tools"
	"msg"
	"time"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

// GetMissionExtraReward 获取关卡额外奖励
func GetMissionExtraReward(p *player.Player, missionId int) (msg.ErrCode, []*model.SimpleItem) {
	if missionId > p.UserData.StageInfo.MissionId {
		return msg.ErrCode_INVALID_DATA, nil
	}

	missionConfig := template.GetMissionTemplate().GetMission(int(missionId))
	if missionConfig == nil || len(missionConfig.ExtraRewardMap) == 0 {
		return msg.ErrCode_INVALID_DATA, nil
	}

	if tools.ListIntContain(p.UserData.Mission.ExtraIds, int(missionId)) {
		return msg.ErrCode_HAS_GET_MISSION_REWARD, nil
	}

	var notifyItems []uint32
	var temp []*model.SimpleItem
	for k, v := range missionConfig.ExtraRewardMap {
		addItems := AddItem(p.GetUserId(), uint32(k), int32(v), publicconst.MissionExtraAddItem, false)
		notifyItems = tools.ListUint32AddNoRepeats(notifyItems, GetSimpleItemIds(addItems))
		temp = append(temp, addItems...)
	}

	updateClientItemsChange(p.GetUserId(), notifyItems)

	p.UserData.Mission.ExtraIds = append(p.UserData.Mission.ExtraIds, int(missionId))
	p.SaveMission()
	return msg.ErrCode_SUCC, temp
}

// GetMissionBoxReward 获取关卡宝箱奖励
func GetMissionBoxReward(p *player.Player, missionId int) (msg.ErrCode, uint32, []*model.SimpleItem) {
	missionConfig := template.GetMissionTemplate().GetMission(int(missionId))
	if missionConfig == nil {
		return msg.ErrCode_MISSION_NOT_EXIST, 0, nil
	}

	var mission *model.Mission
	if missionConfig.Type == int(msg.BattleType_Battle_Challenge) {
		mission = findMission(p, missionId, false)
	} else {
		mission = findMission(p, missionId, true)
	}

	if mission == nil {
		return msg.ErrCode_MISSION_NOT_EXIST, 0, nil
	}

	var boxes []*template.JRuleBox
	for i := 0; i < len(missionConfig.RewardComplete); i++ {
		if tools.GetBit(mission.BoxState, uint32(i)) &&
			!tools.GetBit(mission.BoxRewardState, uint32(i)) {
			ruleBox := template.GetRuleBoxTemplate().GetRuleBox(missionConfig.RewardComplete[i])
			boxes = append(boxes, ruleBox)
			tools.SetBit(&mission.BoxRewardState, uint32(i))
		}
	}

	if len(boxes) == 0 {
		return msg.ErrCode_MISSION_BOX_REWARD_HAS_GET, 0, nil
	}

	var notifyItems []uint32
	var simpleItems []*model.SimpleItem
	simpleItemMap := make(map[uint32]uint32)
	for k := 0; k < len(boxes); k++ {
		ruleBox := boxes[k]
		for i := 0; i < len(ruleBox.Reward); i++ {
			addItems := AddItem(p.GetUserId(),
				ruleBox.Reward[i].ItemId,
				int32(ruleBox.Reward[i].ItemNum),
				publicconst.MissionBoxAddItem,
				false)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)

			if count, ok := simpleItemMap[ruleBox.Reward[i].ItemId]; ok {
				simpleItemMap[ruleBox.Reward[i].ItemId] = count + ruleBox.Reward[i].ItemNum
			} else {
				simpleItemMap[ruleBox.Reward[i].ItemId] = ruleBox.Reward[i].ItemNum
			}
		}
	}

	updateClientItemsChange(p.GetUserId(), notifyItems)
	p.SaveMission()

	for id, num := range simpleItemMap {
		simpleItems = append(simpleItems, model.NewSimpleItem(id, num))
	}

	UpdateTask(p, true, publicconst.TASK_COND_MISSION_BOX_REWARD, 1)

	return msg.ErrCode_SUCC, mission.BoxRewardState, simpleItems
}

func GetChallengeMaxMission(p *player.Player) int {
	var maxId = 0
	for i := 0; i < len(p.UserData.Mission.Challenges); i++ {
		if p.UserData.Mission.Challenges[i].IsPass && p.UserData.Mission.Challenges[i].MissionId > maxId {
			maxId = p.UserData.Mission.Challenges[i].MissionId
		}
	}
	return maxId
}

// addMission 添加关卡
func addMission(p *player.Player, mission *model.Mission, isMain bool) {
	if isMain {
		p.UserData.Mission.Missions = append(p.UserData.Mission.Missions, mission)
	} else {
		p.UserData.Mission.Challenges = append(p.UserData.Mission.Challenges, mission)
	}
	p.SaveMission()
}

// UpdateStory 更新剧情
func UpdateStory(p *player.Player, missionId int, isBefore int) msg.ErrCode {
	missionConfig := template.GetMissionTemplate().GetMission(int(missionId))
	if missionConfig == nil {
		return msg.ErrCode_MISSION_NOT_EXIST
	}

	mission := findMission(p, missionId, true)
	if mission == nil {
		mission = model.NewMission(int(missionId), 10000, 0, false)
		addMission(p, mission, true)
	}

	if isBefore == 1 {
		if mission.BeforeStory == 1 {
			return msg.ErrCode_INVALID_DATA
		}
		mission.BeforeStory = 1
	} else {
		if mission.AfterStory == 1 {
			return msg.ErrCode_INVALID_DATA
		}
		mission.AfterStory = 1
	}
	p.SaveMission()
	return msg.ErrCode_SUCC
}
func GetMission(playerData *player.Player, missionId int, isMain bool) *model.Mission {
	var missions []*model.Mission
	if isMain {
		missions = playerData.UserData.Mission.Missions
	} else {
		missions = playerData.UserData.Mission.Challenges
	}
	for i := 0; i < len(missions); i++ {
		if missions[i].MissionId == missionId {
			return missions[i]
		}
	}
	return nil
}

// challengeStartBattle 挑战开始战斗
func challengeStartBattle(p *player.Player, missionConfig *template.JMission) msg.ErrCode {
	// 更新当前关卡数据
	p.UserData.BaseInfo.MissData = model.NewMissionData(missionConfig.Id, tools.GetCurTime())
	p.SaveBaseInfo()

	// 挑战、扫荡X类型关卡X次
	UpdateTask(p, true,
		publicconst.TASK_COND_BATTLE_MISSIONTYPE, uint32(missionConfig.Type), 1)
	UpdateTask(p, true, publicconst.TASK_COND_ANY_BATTLE, 1)

	UpdateMissionPoker(p, missionConfig.Id)

	return msg.ErrCode_SUCC
}

// checkPassMissionCond 校验通关条件
func checkPassMissionCond(condition int, args []int, result *msg.BattleResult, minTime uint32) bool {
	switch condition {
	case template.KillMonsterId:
		for _, arg := range args {
			exists := false
			for _, d := range result.Data {
				if d.GetBossIds() == uint32(arg) {
					exists = true
					break
				}
			}
			if !exists {
				return false
			}
		}
	case template.Interactive:
		for _, arg := range args {
			if !tools.ListContain(result.InteractiveIds, uint32(arg)) {
				return false
			}
		}
	case template.LiveTime:
		if minTime < uint32(args[0]) ||
			result.BattleTime/1000 < uint32(args[0]) {
			return false
		}
	case template.KillElite:
		for _, arg := range args {
			if !tools.ListContain(result.EliteMonsterIds, uint32(arg)) {
				return false
			}
		}
	}
	return true
}

// isPassMission是否通关
func isPassMission(p *player.Player, result *msg.BattleResult, missionId int) bool {
	if p.CheckBattle == 0 {
		return true
	}

	missionConfig := template.GetMissionTemplate().GetMission(missionId)
	if missionConfig == nil {
		return false
	}

	delta := float64(tools.GetCurTime()-p.UserData.BaseInfo.MissData.StartTime) * float64(p.UserData.BaseInfo.MissData.Speed)
	minTime := p.UserData.BaseInfo.MissData.Total + uint32(delta/10)

	for k, v := range missionConfig.PassConditionMap {
		if !checkPassMissionCond(k, v, result, minTime) {
			return false
		}
	}
	return true
}

// ChallengeEndBattle 挑战关结束战斗
func ChallengeEndBattle(p *player.Player, result *msg.BattleResult) (msg.ErrCode, int, uint32, uint32, uint32, []*model.SimpleItem) {
	curTime := tools.GetCurTime()
	if p.UserData.BaseInfo.MissData == nil || p.UserData.BaseInfo.MissData.MissionId == 0 {
		return msg.ErrCode_MISSION_NO_START, 0, 0, 0, 0, nil
	}

	if curTime-p.UserData.BaseInfo.MissData.StartTime < uint32(publicconst.MISSION_BATTLE_MIN_SECONDS) {
		return msg.ErrCode_INVALID_DATA, 0, 0, 0, 0, nil
	}

	missionConfig := template.GetMissionTemplate().GetMission(int(p.UserData.BaseInfo.MissData.MissionId))

	log.Info("challenge mission end battle", zap.Uint64("accountId", p.GetUserId()),
		zap.Int("missionId", p.UserData.BaseInfo.MissData.MissionId), zap.Reflect("result", result))

	passValue := 0
	isPass := isPassMission(p, result, p.UserData.BaseInfo.MissData.MissionId)
	if isPass {
		passValue = 1
	}
	// 战斗结算
	mission := findMission(p, p.UserData.BaseInfo.MissData.MissionId, false)
	completeTime := result.BattleTime
	// 新关卡
	updateMission := false
	refreshCompleteTime := false
	if mission == nil {
		mission = model.NewMission(int(p.UserData.BaseInfo.MissData.MissionId), completeTime, result.KillMonsterNum, isPass)
		addMission(p, mission, false)

		if isPass {
			if completeTime == 0 {
				mission.CompleteTime = 1000
			}
			refreshCompleteTime = true
		}
	} else {
		if result.KillMonsterNum > mission.KillMonsterNum {
			mission.KillMonsterNum = result.KillMonsterNum
			updateMission = true
		}

		updateFlag := false
		if mission.CompleteTime == 0 {
			if completeTime == 0 {
				mission.CompleteTime = 1000
			} else {
				mission.CompleteTime = completeTime
			}
			updateFlag = true
		} else {
			if completeTime < mission.CompleteTime {
				updateFlag = true
				mission.CompleteTime = completeTime
			}
		}
		if updateFlag {
			updateMission = true
			refreshCompleteTime = true
		}

		if !mission.IsPass && isPass {
			mission.IsPass = isPass
			updateMission = true
		}
	}
	// 宝箱奖励条件
	temp, addStar := updateBattleBoxState(p.GetUserId(), mission, completeTime, result)
	if !updateMission && temp {
		updateMission = temp
	}

	// 满足通关条件
	if isPass {
		// 通关事件
		event.EventMgr.PublishEvent(event.NewPassMissionEvent(p, missionConfig.Id, ListenPassMissionEvent))

		// 刷新排行榜 TODO
		if refreshCompleteTime {
			// ServMgr.GetRankService().updateMissionRankMission(p, mission)
			// ServMgr.GetRankService().updateSpecialMissionRank(p, mission)
		}

		if p.UserData.StageInfo.MissionId < p.UserData.BaseInfo.MissData.MissionId &&
			p.UserData.Contract.TaskId != 0 &&
			p.UserData.Contract.TaskType == msg.ConditionType_Condition_Contract_Rand {
			FinishContract(p, 1)
		}

		if p.UserData.StageInfo.StageFirstPass == nil {
			p.UserData.StageInfo.StageFirstPass = make([]int, 0)
		}
	}

	var (
		notifyItems []uint32
		rewardItems []*model.SimpleItem
	)
	if !tools.ListIntContain(p.UserData.StageInfo.StageFirstPass, mission.MissionId) && isPass {
		for k, v := range missionConfig.NormalReward {
			n := float64(v*GetMonthCardMainScale(p)) / float64(100)
			addItems := AddItem(p.GetUserId(), uint32(k), int32(n), publicconst.PassMissionAddItem, false)
			rewardItems = append(rewardItems, addItems...)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		}

		UpdateTask(p, true, publicconst.TASK_COND_FIRST_PASS_CHALLENGE_MISSION, 1)
		processHistoryData(p, publicconst.TASK_COND_FIRST_PASS_CHALLENGE_MISSION, 0, 1)

		updateClientItemsChange(p.GetUserId(), notifyItems)

		p.UserData.StageInfo.StageFirstPass = append(p.UserData.StageInfo.StageFirstPass, mission.MissionId)
		p.SaveStageInfo()
	}

	if updateMission {
		p.SaveMission()
	}

	p.ResetMissData()

	endMissionTask(p, result)

	p.UserData.BaseInfo.BattleData = nil
	p.SaveBaseInfo()

	if addStar > 0 {
		p.UserData.StageInfo.StageStar += addStar
		p.SaveStageInfo()
		// TODO 排行榜
		//ServMgr.GetRankService().updateStageStarRank(p)
	}

	return msg.ErrCode_SUCC, mission.MissionId, completeTime, mission.BoxState, uint32(passValue), rewardItems
}

// EndBattle 结束战斗
func EndBattle(p *player.Player, result *msg.BattleResult) (msg.ErrCode, uint32, uint32, uint32, uint32, []*model.SimpleItem) {
	curTime := tools.GetCurTime()
	if p.UserData.BaseInfo.MissData == nil || p.UserData.BaseInfo.MissData.MissionId == 0 {
		return msg.ErrCode_MISSION_NO_START, 0, 0, 0, 0, nil
	}

	if curTime-p.UserData.BaseInfo.MissData.StartTime < uint32(publicconst.MISSION_BATTLE_MIN_SECONDS) {
		return msg.ErrCode_INVALID_DATA, 0, 0, 0, 0, nil
	}

	missionConfig := template.GetMissionTemplate().GetMission(int(p.UserData.BaseInfo.MissData.MissionId))
	if msg.BattleType(missionConfig.Type) != msg.BattleType_Battle_Main {
		return msg.ErrCode_INVALID_DATA, 0, 0, 0, 0, nil
	}
	log.Info("mission end battle", zap.Uint64("accountId", p.GetUserId()),
		zap.Int("missionId", p.UserData.BaseInfo.MissData.MissionId), zap.Reflect("result", result))

	passValue := 0
	isPass := isPassMission(p, result, p.UserData.BaseInfo.MissData.MissionId)
	if isPass {
		passValue = 1
	}
	// 战斗结算
	mission := findMission(p, p.UserData.BaseInfo.MissData.MissionId, true)
	//completeTime := curTime - playerData.MissData.StartTime
	completeTime := result.BattleTime

	// 新关卡
	updateMission := false
	refreshCompleteTime := false
	if mission == nil {
		mission = model.NewMission(int(p.UserData.BaseInfo.MissData.MissionId), completeTime, result.KillMonsterNum, isPass)
		addMission(p, mission, true)

		if isPass {
			if completeTime == 0 {
				mission.CompleteTime = 1000
			}
			refreshCompleteTime = true
		}
	} else {
		if result.KillMonsterNum > mission.KillMonsterNum {
			mission.KillMonsterNum = result.KillMonsterNum
			updateMission = true
		}

		updateFlag := false
		if mission.CompleteTime == 0 {
			if completeTime == 0 {
				mission.CompleteTime = 1000
			} else {
				mission.CompleteTime = completeTime
			}
			updateFlag = true
		} else {
			if completeTime < mission.CompleteTime {
				updateFlag = true
				mission.CompleteTime = completeTime
			}
		}
		if updateFlag {
			updateMission = true
			refreshCompleteTime = true
		}

		if !mission.IsPass && isPass {
			mission.IsPass = isPass
			updateMission = true
		}
	}
	// 宝箱奖励条件
	temp, addStar := updateBattleBoxState(p.GetUserId(), mission, completeTime, result)
	if !updateMission && temp {
		updateMission = temp
	}

	// 满足通关条件
	if isPass {
		// 当前主线通关需要更新主线结算挂机奖励 主线关卡变化
		if p.UserData.StageInfo.MissionId == 0 ||
			template.GetMissionTemplate().IsNextMission(int(p.UserData.StageInfo.MissionId), int(p.UserData.BaseInfo.MissData.MissionId)) {
			// 结算挂机奖励
			settleOnHook(p, p.UserData.BaseInfo.MissData.MissionId)

			if p.UserData.StageInfo.MissionId < p.UserData.BaseInfo.MissData.MissionId &&
				p.UserData.Contract.TaskId != 0 &&
				p.UserData.Contract.TaskType == msg.ConditionType_Condition_Contract_Rand {
				FinishContract(p, 1)
			}

			p.UserData.StageInfo.MissionId = p.UserData.BaseInfo.MissData.MissionId
			p.SaveStageInfo()
			// // 添加关卡流失统计
			// dao.UserStaticDao.AddLossMission(model.NewLossMission(tools.GetStaticTime(tools.GetCurTime()),
			// 	uint32(p.GetUserId()), p.AccountInfo.MissionId))

			// // tda update mission
			// tdaData := &tda.CommonUser{
			// 	Max_battle_id: strconv.Itoa(int(p.AccountInfo.MissionId)),
			// }
			// tda.TdaUpdateCommonUser(p.TdaCommonAttr.AccountId, p.TdaCommonAttr.DistinctId, tdaData)
		}

		if refreshCompleteTime {
			// ServMgr.GetRankService().updateMissionRankMission(p, mission)
			// ServMgr.GetRankService().updateSpecialMissionRank(p, mission)
		}

		// 触发调查问卷
		trigerQuestion(p, p.UserData.StageInfo.MissionId)

		// 通关事件
		event.EventMgr.PublishEvent(event.NewPassMissionEvent(p, p.UserData.BaseInfo.MissData.MissionId, ListenPassMissionEvent))
	}

	var notifyItems []uint32
	var rewardItems []*model.SimpleItem
	// 发放通关奖励
	if isPass {
		for k, v := range missionConfig.NormalReward {
			num := float64(v*GetMonthCardMainScale(p)) / float64(100)
			addItems := AddItem(p.GetUserId(), uint32(k), int32(num), publicconst.PassMissionAddItem, false)
			rewardItems = append(rewardItems, addItems...)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		}
	} else {
		scale := result.BattleTime / 60
		scale = tools.LimitUint32(scale, uint32(missionConfig.RewardQuitLimit))
		if scale > 0 {
			for k, v := range missionConfig.ExitReward {
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

	// 通知客户端
	if len(notifyItems) > 0 {
		//event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(playerData, notifyItems, ListenNotifyClientItemEventEvent))
		updateClientItemsChange(p.GetUserId(), notifyItems)

	}

	// 更新关卡
	if updateMission {
		p.SaveMission()
	}

	// 主线关卡记录
	// para := fmt.Sprintf("missionId:%v,start:%v end:%v pass:%v|", p.AccountInfo.MissData.MissionId, p.AccountInfo.MissData.StartTime, curTime, isPass)
	// for i := 0; i < len(rewardItems); i++ {
	// 	para += fmt.Sprintf("%v,%v|", rewardItems[i].Id, rewardItems[i].Num)
	// }
	// ServMgr.GetCommonService().AddStaticsData(p, publicconst.Statics_Main_Mission_Id, para)

	curMissionId := p.UserData.BaseInfo.MissData.MissionId
	p.ResetMissData()
	p.SaveBaseInfo()

	endMissionTask(p, result)

	p.UserData.BaseInfo.BattleData = nil
	p.SaveBaseInfo()

	if addStar > 0 {
		p.UserData.StageInfo.StageStar += addStar
		p.SaveStageInfo()
		// TODO 排行
		// ServMgr.GetRankService().updateStageStarRank(p)
	}

	return msg.ErrCode_SUCC, uint32(curMissionId), completeTime, mission.BoxState, uint32(passValue), rewardItems
}

// StartBattle 开始战斗
func StartBattle(p *player.Player, missionId int) msg.ErrCode {
	accountId := p.GetUserId()
	mission := template.GetMissionTemplate().GetMission(int(missionId))
	if mission == nil {
		return msg.ErrCode_MISSION_NOT_EXIST
	}

	btType := msg.BattleType(mission.Type)
	p.FightType = btType
	p.LastFightTime = time.Now()

	switch btType {
	case msg.BattleType_Battle_Challenge:
		return challengeStartBattle(p, mission)
	case msg.BattleType_Battle_Coin,
		msg.BattleType_Battle_Equip,
		msg.BattleType_Battle_Weapon,
		msg.BattleType_Battle_Desert,
		msg.BattleType_Battle_Union,
		msg.BattleType_Battle_Zombie:
		return PlayMethodStartBattle(p, missionId)
	default:
		// do something
	}

	// 没有足够的体力
	if !EnoughItem(accountId, uint32(publicconst.ITEM_CODE_AP), uint32(mission.PowerCost)) {
		return msg.ErrCode_MISSION_NO_ENOUGH_AP
	}

	// 解锁判定
	if err := CanLock(p, mission.UnlockCond); err != msg.ErrCode_SUCC {
		return err
	}

	// 校验关卡
	if err := checkMission(missionId, p); err != msg.ErrCode_SUCC {
		return err
	}

	// 扣除体力
	if mission.PowerCost > 0 {
		if res := CostItem(accountId, uint32(publicconst.ITEM_CODE_AP), uint32(mission.PowerCost), publicconst.MissionCostItem, true); res != msg.ErrCode_SUCC {
			return res
		}
	}

	// 更新当前关卡数据
	p.UserData.BaseInfo.MissData = model.NewMissionData(missionId, tools.GetCurTime())
	p.SaveBaseInfo()

	UpdateMissionPoker(p, missionId)

	// 挑战、扫荡X类型关卡X次
	UpdateTask(p, true, publicconst.TASK_COND_BATTLE_MISSIONTYPE, uint32(mission.Type), 1)
	UpdateTask(p, true, publicconst.TASK_COND_ANY_BATTLE, 1)

	return msg.ErrCode_SUCC
}

func trigerQuestion(p *player.Player, missionId int) {
	//if tools.ListContain(playerData.AccountInfo.QuestionIds, missionId) {
	//	return
	//}
	//
	//info := template.GetSystemItemTemplate().GetQuestionInfo(missionId)
	//if info == nil {
	//	return
	//}
	//
	//mailConfig := template.GetMailTemplate().GetMail(info.MailId)
	//if mailConfig == nil {
	//	return
	//}
	//
	//var items []*model.SimpleItem
	//for i := 0; i < len(info.Reward); i++ {
	//	items = append(items, &model.SimpleItem{Id: info.Reward[i].ItemId, Num: info.Reward[i].ItemNum})
	//}
	//
	//endTime := time.Now().AddDate(0, 0, 90)
	//mail := model.NewMail(playerData.GenId(), mailConfig.StrTitle, mailConfig.StrContent, items, uint32(endTime.Unix()))
	//ServMgr.GetMailService().AddMail(playerData, mail)
	//event.EventMgr.PublishEvent(event.NewMailEvent(playerData, mail.MailId, false, ListenMailEvent))
	//
	//playerData.AccountInfo.QuestionIds = append(playerData.AccountInfo.QuestionIds, missionId)
	//dao.AccountDao.UpdateQuestionIds(playerData.GetUserId(), playerData.AccountInfo.QuestionIds)
}

// checkMission 关卡校验
func checkMission(missionId int, playerData *player.Player) msg.ErrCode {
	mission := template.GetMissionTemplate().GetMission(int(missionId))
	if mission == nil {
		return msg.ErrCode_MISSION_NOT_EXIST
	}

	// 没有打过的关卡 需要是当前主线下一关
	if findMission(playerData, missionId, true) == nil {
		if playerData.UserData.StageInfo.MissionId == 0 {
			if int(missionId) != template.GetMissionTemplate().FirstMission.Id {
				return msg.ErrCode_MISSION_LAST_NOT_PASS
			}
		} else {
			curMission := template.GetMissionTemplate().GetMission(int(playerData.UserData.StageInfo.MissionId))
			if curMission != nil && curMission.NextId != int(missionId) {
				return msg.ErrCode_MISSION_LAST_NOT_PASS
			}
		}
	}
	return msg.ErrCode_SUCC
}

// findMission 查找关卡
func findMission(p *player.Player, missionId int, isMain bool) *model.Mission {
	if p.UserData.Mission == nil {
		log.Error("mission data nil", zap.Uint64("accountId", p.GetUserId()))
		return nil
	}
	var missions []*model.Mission
	if isMain {
		missions = p.UserData.Mission.Missions
	} else {
		missions = p.UserData.Mission.Challenges
	}
	for i := 0; i < len(missions); i++ {
		if missions[i].MissionId == missionId {
			return missions[i]
		}
	}
	return nil
}

// passMissionNum 通关主线数量
func passMissionNum(p *player.Player, isMain bool) uint32 {
	var missions []*model.Mission
	if isMain {
		missions = p.UserData.Mission.Missions
	} else {
		missions = p.UserData.Mission.Challenges
	}
	var ret uint32 = 0
	for i := 0; i < len(missions); i++ {
		if missions[i].IsPass {
			ret += 1
		}
	}
	return ret
}

// IsPassMission 是否通关
func IsPassMission(p *player.Player, missionid int, isMain bool) bool {
	if p.UserData.Mission == nil {
		log.Error("mission data nil", zap.Uint64("accountId", p.GetUserId()))
		return false
	}
	var missions []*model.Mission
	if isMain {
		missions = p.UserData.Mission.Missions
	} else {
		missions = p.UserData.Mission.Challenges
	}

	for i := 0; i < len(missions); i++ {
		if missions[i].MissionId == missionid && missions[i].IsPass {
			return true
		}
	}
	return false
}

// getBattleBoxState 获得战斗宝箱状态
func updateBattleBoxState(accountId uint64, mission *model.Mission, battleTime uint32, result *msg.BattleResult) (bool, uint32) {
	if mission == nil {
		return false, 0
	}
	missionConfig := template.GetMissionTemplate().GetMission(mission.MissionId)
	if missionConfig == nil {
		return false, 0
	}

	update := false
	addStar := uint32(0)
	for i := 0; i < len(missionConfig.RewardComplete); i++ {
		if !tools.GetBit(mission.BoxState, uint32(i)) {
			if ruleBox := template.GetRuleBoxTemplate().GetRuleBox(missionConfig.RewardComplete[i]); ruleBox != nil {
				pass := true
				for k := 0; k < len(ruleBox.Conds); k++ {
					if !checkBattBoxCond(mission, ruleBox.Conds[k], battleTime, result) {
						pass = false
						break
					}
				}
				if !pass {
					break
				}
				addStar++
				//oldState := mission.BoxState
				tools.SetBit(&mission.BoxState, uint32(i))
				//if mission.BoxState == 2 ||
				//	mission.BoxState == 4 ||
				//	mission.BoxState == 6 {
				//	log.Debug("mission boxState err", zap.Int64("accountID", accountId),
				//		zap.Int("i", i), zap.Any("boxRuleIds", missionConfig.RewardComplete),
				//		zap.Any("cond", ruleBox.PassCondition), zap.Uint32("missionId", mission.MissionId),
				//		zap.Uint32("oldStage", oldState), zap.Uint32("box", missionConfig.RewardComplete[i]))
				//}
				update = true
			}
		}
	}
	return update, addStar
}

// checkBattBoxCond 检测宝箱条件
func checkBattBoxCond(mission *model.Mission, para *template.BoxCond, battleTime uint32, result *msg.BattleResult) bool {
	switch para.Cond {
	case template.BattleTime: // 战斗时长
		if len(para.Args) >= 1 {
			if battleTime/1000 >= uint32(para.Args[0]) {
				return true
			}
		}
	case template.KillBossId: // 杀死指定的怪物
		for i := 0; i < len(para.Args); i += 1 {
			bossId := uint32(para.Args[i])
			exist := false
			for j := 0; j < len(result.Data); j++ {
				if result.Data[j].BossIds == bossId {
					exist = true
					break
				}
			}
			if !exist {
				return false
			}
		}
		return true
	case template.PASSMISSION:
		if mission.IsPass && result.GetExitBattle() != 1 {
			return true
		}
	case template.ShipHp:
		if len(para.Args) >= 1 && result.GetExitBattle() != 1 && mission.IsPass {
			return result.ShipHp >= uint32(para.Args[0])
		}
	}
	return false
}

func endMissionTask(p *player.Player, result *msg.BattleResult) {
	if result.KillMonsterNum > 0 {
		UpdateTask(p, true,
			publicconst.TASK_COND_KILL_MONSTER_NUM, result.KillMonsterNum)
		processHistoryData(p, publicconst.TASK_COND_KILL_MONSTER_NUM, 0, result.KillMonsterNum)
	}

	if result.ChipNum > 0 {
		UpdateTask(p, true,
			publicconst.TASK_COND_GET_CHIP_NUM, result.ChipNum)
	}

	if len(result.PokerData) > 0 {
		var args []uint32
		for i := 0; i < len(result.PokerData); i++ {
			args = append(args, result.PokerData[i].PokerId)
			args = append(args, result.PokerData[i].Num)
		}
		UpdateTask(p, true,
			publicconst.TASK_COND_POKER_NUM, args...)
	}

	eliteAndBoosNum := len(result.EliteMonsterIds) + len(result.Data)
	if eliteAndBoosNum > 0 {
		UpdateTask(p, true,
			publicconst.TASK_COND_KILL_ELITE_MONSTER_BOSS_NUM, uint32(eliteAndBoosNum))
		processHistoryData(p, publicconst.TASK_COND_KILL_ELITE_MONSTER_BOSS_NUM, 0, uint32(eliteAndBoosNum))
	}

	updateDesertLamp(p, result)
}

// updateMagicLamp 更新沙漠神灯
func updateDesertLamp(p *player.Player, result *msg.BattleResult) {
	// 本功能挪到了desert_service
	addNum := uint32(0)
	for _, monster := range result.Monsters {
		addNum += monster.Num
	}
	if addNum > 0 {
		AddKillTimes(p, addNum)
	}
}

// ToProtocolMission 协议关卡数据
func ToProtocolMission(mission *model.Mission) *msg.Mission {
	isPass := 0
	if mission.IsPass {
		isPass = 1
	}
	return &msg.Mission{
		MissionId:      uint32(mission.MissionId),
		BoxState:       mission.BoxState,
		BoxRewardState: mission.BoxRewardState,
		IsPass:         uint32(isPass),
	}
}

// ToProtocolMissions 协议关卡数据
func ToProtocolMissions(missions []*model.Mission) []*msg.Mission {
	var ret []*msg.Mission
	for i := 0; i < len(missions); i++ {
		ret = append(ret, ToProtocolMission(missions[i]))
	}
	return ret
}
