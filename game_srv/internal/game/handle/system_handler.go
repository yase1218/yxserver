package handle

import (
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"gameserver/internal/publicconst"
	"msg"
)

// RequestBuyApHandle 购买体力
func RequestBuyApHandle(packetId uint32, args interface{}, p *player.Player) {
	retMsg := &msg.ResponseBuyAp{}
	var items []*model.SimpleItem
	retMsg.Result, retMsg.BuyTimes, items = service.BuyAp(p)
	if retMsg.Result == msg.ErrCode_SUCC {
		retMsg.GetItems = service.ToProtocolSimpleItems(items)
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestUseCdkHandle 使用cdk
func RequestUseCdkHandle(packetId uint32, args interface{}, p *player.Player) {
	reqMsg := args.(*msg.RequestUseCdk)
	retMsg := &msg.ResponseUseCdk{}
	err := service.UseCdk(p, reqMsg.Cdk, packetId)
	retMsg.Result = err
	if err != msg.ErrCode_SUCC {
		p.SendResponse(packetId, retMsg, retMsg.Result)
	}
}

// RequestRedPointHandle 红点
func RequestRedPointHandle(packetId uint32, args interface{}, p *player.Player) {
	retMsg := &msg.ResponseRedPoint{
		Result: msg.ErrCode_SUCC,
	}
	if redPoint := service.GetShopRedPoint(p); redPoint != nil {
		retMsg.Data = append(retMsg.Data, redPoint)
	}
	if redPoint := service.GetActivityRedPoint(p); redPoint != nil {
		retMsg.Data = append(retMsg.Data, redPoint)
	}
	// TODO 排行榜
	// if redPoint := service.GetMissionRewardRedPoint(playerData); redPoint != nil {
	// 	retMsg.Data = append(retMsg.Data, redPoint)
	// }
	if redPoin := service.GetRankRedPoint(p); redPoin != nil {
		retMsg.Data = append(retMsg.Data, redPoin)
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestSetGuideInfoHandle 设置新手引导
func RequestSetGuideInfoHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestSetGuideInfo)
	retMsg := &msg.ResponseSetGuideInfo{
		Result: msg.ErrCode_SUCC,
		Id:     req.Id,
		Value:  req.Value,
	}
	retMsg.Result = service.SetGuideInfo(p, req.Id, req.Value)
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestSetVideoHandle 设置视频
func RequestSetVideoHandle(packetId uint32, args interface{}, p *player.Player) {
	retMsg := &msg.ResponseSetVideo{
		Result: msg.ErrCode_SUCC,
	}
	retMsg.Result = service.SetVideo(p)
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestClickUpdateNickHandle 点击设置昵称
func RequestClickUpdateNickHandle(packetId uint32, args interface{}, p *player.Player) {
	retMsg := &msg.ResponseClickUpdateNick{
		Result: msg.ErrCode_SUCC,
	}

	service.UpdateTask(p, true, publicconst.TASK_COND_CLICK_UPDATE_NICK, 1)
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestSetPopUpHandle 设置弹窗
func RequestSetPopUpHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestSetPopUp)
	retMsg := &msg.ResponseSetPopUp{
		Result: msg.ErrCode_SUCC,
	}

	service.SetPopUp(p, req.Id, req.PopType)
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestStaticsActionHandle 请求统计信息
func RequestStaticsActionHandle(packetId uint32, args interface{}, p *player.Player) {
	retMsg := &msg.ResponseStaticsAction{}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestReconnectInfoHandle 请求重连数据
func RequestReconnectInfoHandle(packetId uint32, args interface{}, p *player.Player) {
	retMsg := &msg.ResponseReconnectInfo{
		Result: msg.ErrCode_SUCC,
	}
	retMsg.Data = p.UserData.BaseInfo.BattleData
	retMsg.Attrs = service.ToProtocolAttrs2(p.UserData.BaseInfo.Attrs)
	p.SendResponse(packetId, retMsg, retMsg.Result)

	if p.UserData.BaseInfo.MissData != nil {
		service.UpdateMissionPoker(p, p.UserData.BaseInfo.MissData.MissionId)
	}
}

// RequestUploadBattleDataHandle 请求上传战斗数据
func RequestUploadBattleDataHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestUploadBattleData)
	retMsg := &msg.ResponseUploadBattleData{
		Result: msg.ErrCode_SUCC,
	}
	if p.UserData.BaseInfo.MissData != nil {
		retMsg.Result = service.SaveBattleData(p, req.Data)
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestGetQuestionRewardHandle 请求问卷调查奖励
func RequestGetQuestionRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestGetQuestionReward)
	retMsg := &msg.ResponseGetQuestionReward{
		Result: msg.ErrCode_SUCC,
	}
	err, items := service.GetQuestionReward(p, req.QuestionId)
	retMsg.Result = err
	if err == msg.ErrCode_SUCC {
		retMsg.QuestionId = req.QuestionId
		retMsg.GetItems = service.TemplateItemToProtocolItems(items)
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

func RequestStartAdHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestStartAd)
	retMsg := &msg.ResponseStartAd{
		Result: service.StartAd(p, req.AdId, req.Para),
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestGetMonthCardRewardHandle 领取月卡奖励
func RequestGetMonthCardRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	retMsg := &msg.ResponseGetMonthCardReward{}
	err, items := service.GetMonthCard(p)
	retMsg.Result = err
	if err == msg.ErrCode_SUCC {
		retMsg.GetItems = service.TemplateItemToProtocolItems(items)
		for _, v := range p.UserData.BaseInfo.MonthCard {
			retMsg.Data = append(retMsg.Data, &msg.MonthcardInfo{
				MonthCardId:       uint32(v.Id),
				EndTime:           uint32(v.EndTime),
				NextGetRewardTime: uint32(v.NextGetRewardTime),
			})
		}
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestGetMainFundRewardHandle 获取主线基金
func RequestGetMainFundRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestGetMainFundReward)
	retMsg := &msg.ResponseGetMainFundReward{}
	err, items, info := service.GetMainFund(p, int(req.FundId))
	retMsg.Result = err
	if err == msg.ErrCode_SUCC {
		retMsg.GetItems = service.TemplateItemToProtocolItems(items)
		retMsg.FundInfo = &msg.MainFundInfo{}
		retMsg.FundInfo.FundId = uint32(info.Id)
		retMsg.FundInfo.RewardMaxId = uint32(info.FreeId)
		retMsg.FundInfo.BuyFlag = uint32(info.BuyFlag)
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestMonthCardDailyRewardHandle 领取月卡每日奖励
func RequestMonthCardDailyRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	retMsg := &msg.ResponseMonthCardDailyReward{}
	err, nextTime := service.GetMonthCardDailyReward(p)
	retMsg.Result = err
	if err == msg.ErrCode_SUCC {
		retMsg.NextRewardTime = nextTime
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestGetApHandle 领取体力奖励
func RequestGetApHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestGetAp)
	retMsg := &msg.ResponseGetAp{}
	err, state := service.GetDailyAp(p, req.Id)
	retMsg.Result = err
	retMsg.Id = req.Id
	if err == msg.ErrCode_SUCC {
		retMsg.State = state
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestTalentUpgradeHandle 请求升级天赋
func RequestTalentUpgradeHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestTalentUpgrade)
	retMsg := &msg.ResponseTalentUpgrade{
		Id: req.Id,
	}
	retMsg.Result = service.UpgradeTalent(p, req.TalentType)
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestConstructionLotteryHandle 请求升级天赋
func RequestConstructionLotteryHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestConstructionLottery)
	retMsg := &msg.ResponseConstructionLottery{
		Result: service.ConstructionLottery(p, req.Id),
		Id:     req.Id,
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestGetAdventureGuideRewardHandle 请求设置奇遇
func RequestGetAdventureGuideRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestAdventureGuideReward)
	err, items := service.GetAdventureGuideReward(p, req.Id)
	retMsg := &msg.ResponseAdventureGuideReward{
		Result: err,
	}
	if err == msg.ErrCode_SUCC {
		retMsg.GetItems = service.ToProtocolSimpleItems(items)
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}
