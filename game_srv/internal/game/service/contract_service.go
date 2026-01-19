package service

import (
	"kernel/tda"
	"kernel/tools"
	"msg"
	"strconv"

	"github.com/v587-zyf/gc/utils"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"

	"gameserver/internal/enum"
	"gameserver/internal/game/builder"
	"gameserver/internal/game/condition"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
)

func RandCrontractTaskIds() []uint32 {
	var (
		res     = make([]uint32, 0, 3)
		taskMap = make(map[uint32]struct{}, 3)
		n       = 0
	)

	for {
		g := utils.RandWeightByMap(template.GetSystemItemTemplate().ContractRandWeight)
		cfg := template.GetContractTemplate().RandByGroups(uint32(g))
		if _, ok := taskMap[cfg.Id]; ok {
			continue
		}
		taskMap[cfg.Id] = struct{}{}

		n++
		res = append(res, cfg.Id)

		if n == 3 {
			break
		}
	}

	return res
}

func Contract(p *player.Player, ack *msg.ContractAck) msg.ErrCode {
	if p.UserData.Contract == nil {
		log.Error("contract nil", zap.Uint64("uid", p.GetUserId()))
		return msg.ErrCode_SYSTEM_ERROR
	}
	ResetContractAll(p, false, false)
	builder.BuildContractAck(p.UserData, ack, p.ContractRandTime)

	return msg.ErrCode_SUCC
}

func ContractSign(p *player.Player, req *msg.ContractSignReq, ack *msg.ContractSignAck) msg.ErrCode {
	contractCfg := template.GetContractTemplate().GetCfg(req.GetTaskId())
	if contractCfg == nil {
		log.Error("contract cfg nil", zap.Uint32("id", req.GetTaskId()))
		return msg.ErrCode_CONFIG_NIL
	}

	if p.UserData.Contract.TaskId != 0 {
		return msg.ErrCode_CONTRACT_NOT_FINISH
	}
	//if playerData.Contract.SignNum >= enum.Contract_Max_Sign_Num {
	if p.UserData.Contract.SignNum >= template.GetSystemItemTemplate().ContractSignNum {
		return msg.ErrCode_CONTRACT_SIGN_NUM_FULL
	}
	if !tools.ListContain(p.UserData.Contract.TaskIds, req.GetTaskId()) {
		return msg.ErrCode_CONTRACT_NOT_EXIST
	}

	p.UserData.Contract.TaskId = req.GetTaskId()
	p.UserData.Contract.SignNum++

	switch msg.ConditionType(contractCfg.TaskType[0]) {
	case msg.ConditionType_Condition_Contract_Rand:
		p.UserData.Contract.StageEventId = contractCfg.TaskType[1]
	case msg.ConditionType_Condition_Contract_Kill_Monster:
	}
	p.UserData.Contract.TaskType = msg.ConditionType(contractCfg.TaskType[0])

	ack.SignNum = p.UserData.Contract.SignNum
	p.SaveContract()

	// // tda contract accept
	// tda.TdaContractAccept(p.ChannelId, p.TdaCommonAttr, strconv.Itoa(int(contractCfg.Id)), strconv.Itoa(int(p.AccountInfo.MissionId)))

	return msg.ErrCode_SUCC
}

func ContractCancel(p *player.Player, ack *msg.ContractCancelAck) msg.ErrCode {
	if p.UserData.Contract.TaskId == 0 {
		return msg.ErrCode_CONTRACT_NOT_SIGN
	}

	ResetContract(p)

	builder.BuildContractCancelAck(p.UserData, ack)

	return msg.ErrCode_SUCC
}

func ContractRand(p *player.Player, ack *msg.ContractRandAck) msg.ErrCode {
	if p.UserData.Contract.TaskId != 0 {
		return msg.ErrCode_ILLEGAL_OPERATIONS
	}
	//if playerData.Contract.RandNum >= enum.Contract_Max_Rand_Num {
	randMax := false
	if p.UserData.Contract.RandNum >= template.GetSystemItemTemplate().ContractRandNum {
		randMax = true
		if p.UserData.Contract.RandDiamondNum >= template.GetSystemItemTemplate().ContractRandDiamondsTimes {
			return msg.ErrCode_CONTRACT_RAND_NUM_FILL
		} else {
			if !EnoughItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_DIAMOND),
				template.GetSystemItemTemplate().ContractRandDiamonds) {
				return msg.ErrCode_NO_ENOUGH_ITEM
			}
		}
	}

	curTime := tools.GetCurTime()
	if p.ContractRandTime > 0 && (curTime-p.ContractRandTime) < uint32(enum.Contract_Rand_Cd) {
		return msg.ErrCode_IN_CD
	}

	if randMax {
		CostItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_DIAMOND),
			template.GetSystemItemTemplate().ContractRandDiamonds, publicconst.ContractRand, true)

		p.UserData.Contract.RandDiamondNum++
	} else {
		p.UserData.Contract.RandNum++
	}

	p.UserData.Contract.TaskIds = RandCrontractTaskIds()
	p.ContractRandTime = curTime

	builder.BuildContractRandAck(p.UserData, ack, p.ContractRandTime)

	return msg.ErrCode_SUCC
}

func ContractReward(p *player.Player, ack *msg.ContractRewardAck) msg.ErrCode {
	if p.UserData.Contract.TaskId == 0 {
		return msg.ErrCode_CONTRACT_NOT_SIGN
	}

	contractCfg := template.GetContractTemplate().GetCfg(p.UserData.Contract.TaskId)
	if contractCfg == nil {
		log.Error("contract cfg nil", zap.Uint32("id", p.UserData.Contract.TaskId))
		return msg.ErrCode_CONFIG_NIL
	}

	if _, ok := condition.GetCondition().Check(p, contractCfg.TaskType); !ok {
		return msg.ErrCode_CONDITION_NOT_MET
	}

	if p.UserData.Contract.Reward {
		return msg.ErrCode_REPEAT_REWARD
	}

	tdaItems := make([]*tda.Item, 0, len(contractCfg.Reward))
	if contractCfg.Reward != nil {
		var notifyItems []uint32
		for _, item := range contractCfg.Reward {
			addItems := AddItem(p.GetUserId(), item.ItemId, int32(item.ItemNum), publicconst.ContractReward, false)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
			tdaItems = append(tdaItems, &tda.Item{ItemId: strconv.Itoa(int(item.ItemId)), ItemNum: item.ItemNum})
		}
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	ResetContract(p)

	builder.BuildContractRewardAck(p.UserData, ack)
	ack.RewardItem = TemplateItemToProtocolItems(contractCfg.Reward)

	p.UserData.Contract.Point += contractCfg.Point

	// TODO 排行榜
	// updateContractRank(p.GetUserId(), p.UserData.Contract.Point)

	// // tda contract reward
	// tda.TdaContractReward(p.ChannelId, p.TdaCommonAttr, strconv.Itoa(int(contractCfg.Id)), tdaItems)

	return msg.ErrCode_SUCC
}

func ResetContractAll(p *player.Player, sendNtf bool, gm bool) {
	resetTime := tools.GetDailyRefreshTime()
	if !gm {
		if p.UserData.Contract.ResetDate == resetTime {
			return
		}
	}

	ResetContract(p)

	p.UserData.Contract.SignNum = 0
	p.UserData.Contract.RandNum = 0
	p.UserData.Contract.ResetDate = resetTime
	p.UserData.Contract.RandDiamondNum = 0

	p.SaveContract()
	if sendNtf {
		p.SendNotify(builder.BuildContractResetNtf(p.UserData))
	}
}

func FinishContract(p *player.Player, finishNum uint32) {
	if p.UserData.Contract.TaskId == 0 {
		return
	}

	contractCfg := template.GetContractTemplate().GetCfg(p.UserData.Contract.TaskId)
	if contractCfg == nil {
		log.Error("contract cfg nil", zap.Uint32("id", p.UserData.Contract.TaskId))
		return
	}

	if _, ok := condition.GetCondition().Check(p, contractCfg.TaskType); ok {
		return
	}

	p.UserData.Contract.FinishCount += finishNum
	switch msg.ConditionType(contractCfg.TaskType[0]) {
	case msg.ConditionType_Condition_Contract_Rand:
		if p.UserData.Contract.FinishCount >= contractCfg.TaskType[2] {
			p.UserData.Contract.FinishCount = contractCfg.TaskType[2]
		}
	case msg.ConditionType_Condition_Contract_Kill_Monster:
		if p.UserData.Contract.FinishCount >= contractCfg.TaskType[1] {
			p.UserData.Contract.FinishCount = contractCfg.TaskType[1]
		}
	}

	p.SendNotify(&msg.ContractFinishNtf{
		FinishCount: p.UserData.Contract.FinishCount,
		TaskId:      p.UserData.Contract.TaskId,
	})
}

func ResetContract(p *player.Player) {
	p.UserData.Contract.TaskIds = RandCrontractTaskIds()
	p.UserData.Contract.TaskId = 0
	p.UserData.Contract.TaskType = 0
	p.UserData.Contract.StageEventId = 0
	p.UserData.Contract.FinishCount = 0
	p.UserData.Contract.Reward = false
}
