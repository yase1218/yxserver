package handle

import (
	"msg"

	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
)

// RequestLoadMissionsHandle 加载关卡
func RequestLoadMissionsHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponseLoadMissions{
		Result: msg.ErrCode_SUCC,
	}

	res.Data = append(res.Data, service.ToProtocolMissions(p.UserData.Mission.Missions)...)
	res.ChallengeData = append(res.ChallengeData, service.ToProtocolMissions(p.UserData.Mission.Challenges)...)

	for _, v := range p.UserData.Mission.ExtraIds {
		res.GetExtraRewardMissions = append(res.GetExtraRewardMissions, uint32(v))
	}
	p.SendResponse(packetId, res, res.Result)
	// 登录完成 此处认为玩家已经准备好登录, 可以接收服务端发放的其他消息
	//
	onLoginOver(p)
}

func onLoginOver(p *player.Player) {
	// TODO 社交相关
	// // 联盟新人申请, 发送消息给指挥官
	// // 检查玩家是否已有联盟
	// member, err := dao.GetMember(p.AccountInfo.AccountId)
	// if err == nil && member != nil {
	// 	// 获取目标联盟
	// 	alliance, err := dao.GetAlliance(member.AllianceID)
	// 	if err == nil && alliance != nil {
	// 		// 获取申请列表
	// 		applies, err := dao.GetApplicationList(member.AllianceID)
	// 		if err == nil && len(applies) > 0 {
	// 			service.ServMgr.GetAllianceService().BroadcastMsgToAllianceForOfficer(alliance.ID, &msg.NotifyAllianceApply{
	// 				Count: int32(len(applies)),
	// 			}, 0, msg.ErrCode_SUCC, 0, 0)
	// 		}
	// 	}
	// }
}

// RequestStartBattleHandle 开始战斗
func RequestStartBattleHandle(packetId uint32, args interface{}, p *player.Player) {
	// 	var (
	// 		req     = args.(*msg.RequestStartBattle)
	// 		stageId = req.GetMissionId()
	// 	)

	// 	err := service.StartBattle(p, stageId)
	// 	retMsg := &msg.ResponseStartBattle{
	// 		Result:     err,
	// 		MissionId:  req.MissionId,
	// 		StageEvent: make([]uint32, 0),
	// 	}

	// 	if err == msg.ErrCode_SUCC {
	// 		retMsg.Attrs = service.ToProtocolAttrs2(p.UserData.BaseInfo.Attrs)
	// 		retMsg.ShipId = p.UserData.BaseInfo.ShipId
	// 		retMsg.SupportId = p.UserData.BaseInfo.SupportId
	// 		retMsg.ComboSkills = p.UserData.BaseInfo.ComboSkill
	// 		if pet := service.GetPetById(p, p.UserData.BaseInfo.UsePet); pet != nil {
	// 			if len(pet.BaseAttr) > 0 {
	// 				retMsg.PetData = service.ToProtocolPet(pet)
	// 			}
	// 		}
	// 	}

	// 	if mission := service.GetMission(p, int(stageId), true); mission != nil {
	// 		retMsg.BeforeBattleStory = mission.BeforeStory
	// 		retMsg.AfterBattleStory = mission.AfterStory
	// 		service.UpdateStory(p, stageId, 1)
	// 	}

	// 	if p.UserData.Contract.TaskId != 0 && p.UserData.Contract.TaskType == msg.ConditionType_Condition_Contract_Rand {
	// 		contractFlag := false
	// 		switch msg.BattleType(template.GetMissionTemplate().GetMission(stageId).Type) {
	// 		case msg.BattleType_Battle_Main:
	// 			if stageId > p.UserData.StageInfo.MissionId {
	// 				contractFlag = true
	// 			}
	// 		case msg.BattleType_Battle_Challenge:
	// 			maxId := service.GetChallengeMaxMission(p)
	// 			if int(stageId) >= maxId {
	// 				contractFlag = true
	// 			}
	// 		}
	// 		if contractFlag {
	// 			contractCfg := template.GetContractTemplate().GetCfg(p.UserData.Contract.TaskId)
	// 			if contractCfg == nil {
	// 				log.Error("contract cfg nil", zap.Uint32("id", p.UserData.Contract.TaskId))
	// 				return
	// 			}

	// 			if _, ok := condition.GetCondition().Check(p, contractCfg.TaskType); !ok {
	// 				retMsg.StageEvent = append(retMsg.StageEvent, p.UserData.Contract.StageEventId)
	// 			}
	// 		}
	// 	}

	//	if p.UserData.BaseInfo.StageFirstEnter == nil {
	//		p.UserData.BaseInfo.StageFirstEnter = make([]int, 0)
	//	}
	//
	//	if !tools.ListIntContain(p.UserData.BaseInfo.StageFirstEnter, int(stageId)) {
	//		retMsg.IsFirstEnter = true
	//		p.UserData.BaseInfo.StageFirstEnter = append(p.UserData.BaseInfo.StageFirstEnter, int(stageId))
	//		dao.AccountDao.UpdateStageFirstEnter(p.GetAccountId(), p.UserData.BaseInfo.StageFirstEnter)
	//	}
	//
	// p.SendMsg(packetId, retMsg)
}

// RequestEndBattleHandle 结束战斗
func RequestEndBattleHandle(packetId uint32, args interface{}, p *player.Player) {
	// var (
	// 	req     = args.(*msg.RequestEndBattle)
	// 	stageId int
	// )

	// if p.UserData.BaseInfo.MissData != nil {
	// 	stageId = p.UserData.BaseInfo.MissData.MissionId
	// }
	// missionConfig := template.GetMissionTemplate().GetMission(uint32(stageId))
	// if missionConfig == nil {
	// 	retMsg := &msg.ResponseEndBattle{
	// 		Result: msg.ErrCode_MISSION_NOT_EXIST,
	// 	}
	// 	p.SendMsg(packetId, retMsg)
	// 	return
	// }

	// switch msg.BattleType(missionConfig.Type) {
	// case msg.BattleType_Battle_Main: // 主线
	// 	normalEndBattle(packetId, req, p)
	// case msg.BattleType_Battle_Challenge: // 挑战玩法
	// 	challengeEndBattle(packetId, req, p)
	// default:
	// 	playMethodEndBattle(packetId, req, p)
	// }

	// // tda main battle done
	// //tda.TdaMainBattleDoneServer(playerData.ChannelId, playerData.TdaCommonAttr, req.GetSerial(), req.GetResult().BattleTime)

	// p.FightType = msg.BattleType_Battle_None
	// p.LastFightTime = time.Time{}
}

// RequestGetMissionRewardHandle 领取关卡奖励
func RequestGetMissionRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestGetMissionReward)
	err, rewardState, getItems := service.GetMissionBoxReward(p, int(req.MissionId))
	res := &msg.ResponseGetMissionReward{
		Result: err,
	}
	if err == msg.ErrCode_SUCC {
		res.MissionId = req.MissionId
		res.BoxRewardState = rewardState
		res.GetItems = service.ToProtocolSimpleItems(getItems)
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestExitBattleHandle 退出战斗
func RequestExitBattleHandle(packetId uint32, args interface{}, playerData *player.Player) {
	// err := service.ServMgr.GetMissionService().ExitBattle(playerData)
	// retMsg := &msg.ResponseExitBattle{
	// 	Result: err,
	// }
	// tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
}

// RequestUpdateBattleStoryHandle 剧情
func RequestUpdateBattleStoryHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestUpdateBattleStory)
	err := service.UpdateStory(p, int(req.MissionId), int(req.IsBeforeStory))
	res := &msg.ResponseUpdateBattleStory{
		Result:        err,
		MissionId:     req.MissionId,
		IsBeforeStory: req.IsBeforeStory,
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestSetBattleSpeedHandle 设置战斗速度
func RequestSetBattleSpeedHandle(packetId uint32, args interface{}, playerData *player.Player) {
	// req := args.(*msg.RequestSetBattleSpeed)
	// err := service.ServMgr.GetMissionService().SetBattleSpeed(playerData, req.Speed)
	// retMsg := &msg.ResponseSetBattleSpeed{
	// 	Result: err,
	// }
	// tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
}

// RequestMissionExtraRewardHandle 请求关卡额外奖励
func RequestMissionExtraRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestMissionExtraReward)
	err, items := service.GetMissionExtraReward(p, int(req.MissionId))
	res := &msg.ResponseMissionExtraReward{
		Result:    err,
		MissionId: req.MissionId,
		GetItems:  service.ToProtocolSimpleItems(items),
	}

	p.SendResponse(packetId, res, res.Result)
}

// func normalEndBattle(packetId uint32, req *msg.RequestEndBattle, p *player.Player) {
// 	var (
// 		err          msg.ErrCode
// 		missionId    int
// 		completeTime uint32
// 		boxState     uint32
// 		isPass       uint32
// 		rewardItems  []*model.SimpleItem
// 	)

// 	if p.UserData.StageInfo.MissionId < p.UserData.BaseInfo.MissData.MissionId &&
// 		p.UserData.Contract.TaskId != 0 &&
// 		p.UserData.Contract.TaskType == msg.ConditionType_Condition_Contract_Kill_Monster {
// 		service.ServMgr.GetContractService().Finish(p, req.GetResult().GetKillMonsterNum())
// 	}

// 	if req.Result.ExitBattle == 1 {
// 		err, missionId, completeTime, boxState, isPass, rewardItems = service.ServMgr.GetMissionService().ExitBattleSettle(p, req.Result)
// 	} else {
// 		err, missionId, completeTime, boxState, isPass, rewardItems = service.ServMgr.GetMissionService().EndBattle(p, req.Result)
// 	}
// 	p.MissionAdRewardItems = rewardItems
// 	retMsg := &msg.ResponseEndBattle{
// 		Result:         err,
// 		MissionId:      uint32(missionId),
// 		KillMonsterNum: req.Result.KillMonsterNum,
// 		BattleTime:     completeTime,
// 		BoxState:       boxState,
// 		MainMissionId:  uint32(p.AccountInfo.MissionId),
// 		IsPass:         isPass,
// 		ExitBattle:     req.Result.ExitBattle,
// 		GetItems:       service.ToProtocolSimpleItems(rewardItems),
// 		MissionPass:    service.ServMgr.GetMissionService().GetMissionPassValue(p, missionId, true),
// 	}
// 	tools.SendMsg(p.PlayerAgent, retMsg, packetId, retMsg.Result)
// }

// func challengeEndBattle(packetId uint32, req *msg.RequestEndBattle, playerData *player.Player) {
// 	var (
// 		err          msg.ErrCode
// 		missionId    int
// 		completeTime uint32
// 		boxState     uint32
// 		isPass       uint32
// 		rewardItems  []*model.SimpleItem
// 	)

// 	if req.Result.ExitBattle == 1 {
// 		err, missionId, completeTime, boxState, isPass, rewardItems = service.ServMgr.GetMissionService().ExitBattleSettle(playerData, req.Result)
// 	} else {
// 		err, missionId, completeTime, boxState, isPass, rewardItems = service.ServMgr.GetMissionService().ChallengeEndBattle(playerData, req.Result)
// 	}
// 	playerData.MissionAdRewardItems = rewardItems
// 	retMsg := &msg.ResponseEndBattle{
// 		Result:         err,
// 		MissionId:      uint32(missionId),
// 		KillMonsterNum: req.Result.KillMonsterNum,
// 		BattleTime:     completeTime,
// 		BoxState:       boxState,
// 		MainMissionId:  uint32(playerData.AccountInfo.MissionId),
// 		IsPass:         isPass,
// 		ExitBattle:     req.Result.ExitBattle,
// 		GetItems:       service.ToProtocolSimpleItems(rewardItems),
// 		MissionPass:    service.ServMgr.GetMissionService().GetMissionPassValue(playerData, missionId, false),
// 	}
// 	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
// }

// // PlayMethodEndBattle 玩法结束战斗
// func playMethodEndBattle(packetId uint32, req *msg.RequestEndBattle, playerData *player.Player) {
// 	retMsg := &msg.ResponsePlayMethodEndBattle{
// 		Result: msg.ErrCode_SUCC,
// 	}
// 	if playerData.AccountInfo.MissData != nil {
// 		retMsg.MissionId = uint32(playerData.AccountInfo.MissData.MissionId)
// 	}

// 	//if err := service.ServMgr.GetCommonService().FunctionOpen(playerData, publicconst.PlayMethod); err != msg.ErrCode_SUCC {
// 	//	retMsg.Result = err
// 	//	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
// 	//	return
// 	//}

// 	log.Info("play method end battle", zap.Int64("accountId", playerData.GetAccountId()),
// 		zap.Uint32("missionId", uint32(playerData.AccountInfo.MissData.MissionId)), zap.Reflect("req", req))

// 	err, isPass, items := service.ServMgr.GetPlayMethodService().EndBattle(playerData, req.Result)
// 	retMsg.Result = err
// 	if isPass {
// 		retMsg.IsPass = 1
// 	}
// 	if len(items) > 0 {
// 		retMsg.GetItems = append(retMsg.GetItems, service.ToProtocolSimpleItems(items)...)
// 	}
// 	retMsg.KillMonsterNum = req.Result.KillMonsterNum
// 	retMsg.BattleTime = req.Result.BattleTime
// 	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
// }

// // PlayMethodEndBattle 玩法结束战斗
// func peakEndBattle(packetId uint32, req *msg.RequestEndBattle, playerData *player.Player) {

// }
