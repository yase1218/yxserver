package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"kernel/tools"
	"msg"

	template2 "github.com/zy/game_data/template"
)

// RequestOnHookDataHandle 请求挂机数据
func RequestOnHookDataHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestOnHookData)
	res := &msg.ResponseOnHookData{}
	res.Result = service.GetOnHookData(p)
	if p.UserData.BaseInfo.HookData != nil {
		res.OnHookTime = p.UserData.BaseInfo.HookData.TotalTime + tools.GetCurTime() - p.UserData.BaseInfo.HookData.StartTime
		if res.OnHookTime >= template2.GetSystemItemTemplate().MaxOnHookTime {
			res.OnHookTime = template2.GetSystemItemTemplate().MaxOnHookTime
		}
		for i := 0; i < len(p.UserData.BaseInfo.HookData.Items); i++ {
			item := p.UserData.BaseInfo.HookData.Items[i]
			tempNum := int64(item.Num)
			if tempNum > 0 {
				res.Data = append(res.Data, &msg.Item{
					ItemId:  item.Id,
					ItemNum: int64(item.Num),
				})
			}
		}
	}
	res.ClientData = req.ClientData
	p.SendResponse(packetId, res, res.Result)
}

// RequestOnHookRewardHandle 请求挂机奖励
func RequestOnHookRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponseOnHookReward{}
	err, startTime, items := service.GetOnHookReward(p)
	res.Result = err
	if err == msg.ErrCode_SUCC {
		res.StartTime = startTime
		res.Data = service.ProtocolSimpleItemsToItems(items)
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestQuickOnHookInfoHandle 请求快速挂机
func RequestQuickOnHookInfoHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponseQuickOnHookInfo{}
	res.Result, res.Times = service.GetQuickOnHookInfo(p)
	p.SendResponse(packetId, res, res.Result)
}

// RequestQuickOnHookRewardHandle 请求快速挂机奖励
func RequestQuickOnHookRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponseQuickOnHookReward{}
	err, times, items := service.GetQuickOnHookReward(p)
	res.Result = err
	if err == msg.ErrCode_SUCC {
		res.Times = times
		res.Data = service.ToProtocolItems(items)
	}

	p.SendResponse(packetId, res, res.Result)
}
