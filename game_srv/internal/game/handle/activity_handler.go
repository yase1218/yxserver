package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

// RequestActConfigHandle 请求活动配置
func RequestActConfigHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponseActConfig{
		Result: msg.ErrCode_SUCC,
	}
	res.Data = service.ToProtocolActConfigs(p.UserData.AccountActivity.Activities)
	res.GetPreRewardTps = p.UserData.AccountActivity.PreRewardTps
	p.SendResponse(packetId, res, res.Result)
}

// RequestLoadActDataHandle 请求活动数据
func RequestLoadActDataHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestLoadActData)
	err, data := service.GetActivityData(p, req.ActId)
	res := &msg.ResponseLoadActData{
		Result: err,
	}
	if err == msg.ErrCode_SUCC {
		res.ActId = req.ActId
		res.Data = service.ToProtocolActDatas(data)
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestGetActRewardHandle 请求活动奖励
func RequestGetActRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestGetActReward)
	err, data := service.GetActivityReward(p, req.ActId, req.SubActId, 1)
	res := &msg.ResponseGetActReward{
		Result: err,
	}
	if err == msg.ErrCode_SUCC {
		res.ActId = req.ActId
		res.GetItems = service.ToProtocolSimpleItems(data)
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestBuyPassGradeHandle 购买战令等级
func RequestBuyPassGradeHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestBuyPassGrade)
	err := service.BuyPassGrade(p, req.ActId, req.Grade)
	res := &msg.ResponseBuyPassGrade{
		Result: err,
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestSuppSignHandle 请求补签
func RequestSuppSignHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestSuppSign)
	err, data := service.SuppSign(p, req.ActId)
	res := &msg.ResponseSuppSign{
		Result:   err,
		GetItems: service.ToProtocolSimpleItems(data),
		ActId:    req.ActId,
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestEnterDesertActHandle 请求进入沙漠活动
// func RequestEnterDesertActHandle(packetId uint32, args interface{}, p *player.Player) {
// 	err, rank, data := service.EnterDesert(p)
// 	res := &msg.ResponseEnterDesertAct{
// 		Result: err,
// 	}
// 	if err == msg.ErrCode_SUCC {
// 		res.Rank = rank
// 		res.Data = service.ToProtocolActDatas(data)
// 		res.SettleTime = tools.GetWeeklyRefreshTime(0)
// 	}
// 	p.SendResponse(packetId, res, res.Result)
// }

// RequestActPreviewRewardHandle 请求活动预告奖励
func RequestActPreviewRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestActPreviewReward)
	err, reward := service.GetActPreviewReward(p, req.ActType)
	res := &msg.ResponseActPreviewReward{
		Result: err,
	}
	if err == msg.ErrCode_SUCC {
		res.GetItems = service.ToProtocolSimpleItems(reward)
		res.ActType = req.ActType
	}
	p.SendResponse(packetId, res, res.Result)
}
