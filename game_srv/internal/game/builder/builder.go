package builder

import (
	"gameserver/internal/game/model"
	"msg"
)

func BuildApData(user_data *model.UserData) *msg.ApData {
	ap := &msg.ApData{
		BuyTimes:         user_data.BaseInfo.ApData.BuyTimes,
		RecoverStartTime: user_data.BaseInfo.ApData.RecoverStartTime,
	}
	return ap
}

func BuildGuideInfo(user_data *model.UserData) []*msg.GuideInfo {
	res := make([]*msg.GuideInfo, 0, len(user_data.BaseInfo.GuideData))
	for _, v := range user_data.BaseInfo.GuideData {
		res = append(res, &msg.GuideInfo{
			Id:    v.Id,
			Value: v.Value,
		})
	}
	return res
}

func BuildPopUps(user_data *model.UserData) []*msg.PopUpInfo {
	res := make([]*msg.PopUpInfo, 0, len(user_data.BaseInfo.PopUps))
	for _, v := range user_data.BaseInfo.PopUps {
		res = append(res, &msg.PopUpInfo{
			Id:      v.Id,
			PopType: v.PopUpType,
		})
	}
	return res
}

func BuildChargeInfo(user_data *model.UserData) *msg.RechargeInfo {
	res := &msg.RechargeInfo{
		RechargeIds: make([]uint32, 0, len(user_data.BaseInfo.Charge)),
		McInfo:      make([]*msg.MonthcardInfo, 0, len(user_data.BaseInfo.MonthCard)),
	}

	for _, v := range user_data.BaseInfo.Charge {
		if v.Value > 0 {
			res.RechargeIds = append(res.RechargeIds, uint32(v.Id))
		}
	}

	for _, v := range user_data.BaseInfo.MonthCard {
		res.McInfo = append(res.McInfo, &msg.MonthcardInfo{
			MonthCardId:       uint32(v.Id),
			EndTime:           uint32(v.EndTime),
			NextGetRewardTime: uint32(v.NextGetRewardTime),
		})
	}

	return res
}

func BuildFundInfo(user_data *model.UserData) []*msg.MainFundInfo {
	res := make([]*msg.MainFundInfo, 0, len(user_data.BaseInfo.MainFund))
	for _, v := range user_data.BaseInfo.MainFund {
		res = append(res, &msg.MainFundInfo{
			FundId:      uint32(v.Id),
			RewardMaxId: uint32(v.FreeId),
			BuyFlag:     uint32(v.BuyFlag),
		})
	}
	return res
}

func BuildAdData(user_data *model.UserData) []*msg.AdInfo {
	res := make([]*msg.AdInfo, 0, len(user_data.BaseInfo.Ad))
	for _, v := range user_data.BaseInfo.Ad {
		res = append(res, &msg.AdInfo{
			AdId:  v.AdId,
			Times: v.Times,
		})
	}
	return res
}

func BuildDailyApInfo(user_data *model.UserData) []*msg.DailyApInfo {
	res := make([]*msg.DailyApInfo, 0, len(user_data.BaseInfo.DailyApData))
	for _, v := range user_data.BaseInfo.DailyApData {
		res = append(res, &msg.DailyApInfo{
			Id:        v.Id,
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
			State:     v.State,
		})
	}
	return res
}

func BuildTalentData(user_data *model.UserData) *msg.TalentInfo {
	res := &msg.TalentInfo{
		NormalId: user_data.BaseInfo.TalentData.NormalId,
		KeyId:    user_data.BaseInfo.TalentData.KeyId,
		Attrs:    BuildAttrMap(user_data.BaseInfo.TalentData.Attrs),
		Parts:    user_data.BaseInfo.TalentData.Parts,
	}
	return res
}

func BuildAttrMap(data map[uint32]*model.Attr) []*msg.Attr {
	ret := make([]*msg.Attr, 0, len(data))
	for _, v := range data {
		ret = append(ret, BuildAttr(v))
	}
	return ret
}

func BuildAttr(data *model.Attr) *msg.Attr {
	return &msg.Attr{
		Id:        data.Id,
		Value:     data.InitValue + data.LevelValue + data.Add,
		CalcValue: data.FinalValue,
	}
}

func BuildLuckSaleAck(user_data *model.UserData) *msg.LuckSaleAck {
	pbTask := make([]*msg.Task, 0, len(user_data.LuckSale.Task))
	for _, v := range user_data.LuckSale.Task {
		pbTask = append(pbTask, BuildLuckSaleTask(v))
	}

	ack := &msg.LuckSaleAck{
		Jackpot: int32(user_data.LuckSale.Jackpot),
		Tasks:   pbTask,
	}
	if user_data.LuckSale.Jackpot != -1 {
		ack.Times = user_data.LuckSale.Data[user_data.LuckSale.Jackpot].Times
		ack.Ids = user_data.LuckSale.Data[user_data.LuckSale.Jackpot].Ids
	}

	return ack
}

func BuildLuckSaleTask(v *model.LuckSaleTaskUnit) *msg.Task {
	return &msg.Task{
		TaskId:       v.TaskId,
		TaskValue:    v.Value,
		State:        msg.TaskState(v.State),
		CompleteTime: v.UpdateTime,
	}
}

func BuildContractAck(user_data *model.UserData, ack *msg.ContractAck, rand_time uint32) {
	ack.TaskIds = user_data.Contract.TaskIds
	ack.TaskId = user_data.Contract.TaskId
	ack.FinishCount = user_data.Contract.FinishCount
	ack.Reward = user_data.Contract.Reward
	ack.SignNum = user_data.Contract.SignNum
	ack.RandNum = user_data.Contract.RandNum
	ack.RandCd = int64(rand_time)
	ack.ResetTime = int64(user_data.Contract.ResetDate)
	ack.RandDiamondNum = user_data.Contract.RandDiamondNum
	ack.Point = user_data.Contract.Point
}

func BuildContractCancelAck(user_data *model.UserData, ack *msg.ContractCancelAck) {
	ack.TaskIds = user_data.Contract.TaskIds
	ack.TaskId = user_data.Contract.TaskId
	ack.FinishCount = user_data.Contract.FinishCount
	ack.Reward = user_data.Contract.Reward
}

func BuildContractRandAck(user_data *model.UserData, ack *msg.ContractRandAck, rand_time uint32) {
	ack.TaskIds = user_data.Contract.TaskIds
	ack.RandNum = user_data.Contract.RandNum
	ack.RandCd = int64(rand_time)
	ack.RandDiamondNum = user_data.Contract.RandDiamondNum
}

func BuildContractRewardAck(user_data *model.UserData, ack *msg.ContractRewardAck) {
	ack.TaskIds = user_data.Contract.TaskIds
	ack.TaskId = user_data.Contract.TaskId
	ack.FinishCount = user_data.Contract.FinishCount
	ack.Reward = true
	ack.Point = user_data.Contract.Point
}

func BuildContractResetNtf(user_data *model.UserData) *msg.ContractResetNtf {
	return &msg.ContractResetNtf{
		TaskIds:        user_data.Contract.TaskIds,
		TaskId:         user_data.Contract.TaskId,
		FinishCount:    user_data.Contract.FinishCount,
		Reward:         user_data.Contract.Reward,
		SignNum:        user_data.Contract.SignNum,
		RandNum:        user_data.Contract.RandNum,
		ResetTime:      int64(user_data.Contract.ResetDate),
		RandDiamondNum: user_data.Contract.RandDiamondNum,
		Point:          user_data.Contract.Point,
	}
}

func BuildFunctionPreviewAck(user_data *model.UserData, ack *msg.FunctionPreviewAck) {
	ack.Data = make(map[uint32]uint32, len(user_data.FunctionPreview.Data))
	for k, v := range user_data.FunctionPreview.Data {
		ack.Data[k] = uint32(v)
	}
}

func BuildAtlas(data *model.Atlas) []*msg.AtlasAck_AtlasUnit {
	pbSlice := make([]*msg.AtlasAck_AtlasUnit, 0, len(data.Data))
	for _, v := range data.Data {
		pbSlice = append(pbSlice, BuildAtlasUnit(v))
	}
	return pbSlice
}
func BuildAtlasUnit(data *model.AtlasUnit) *msg.AtlasAck_AtlasUnit {
	return &msg.AtlasAck_AtlasUnit{
		Id:       uint32(data.Id),
		IsReward: data.Reward,
	}
}
