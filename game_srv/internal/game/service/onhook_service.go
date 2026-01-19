package service

import (
	"gameserver/internal/config"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/report"
	"kernel/tools"
	"msg"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

// GetOnHookData 获得挂机数据
func GetOnHookData(p *player.Player) msg.ErrCode {
	if p.UserData.BaseInfo.HookData != nil {
		settleOnHook(p, p.UserData.StageInfo.MissionId)
	}

	return msg.ErrCode_SUCC
}

// GetOnHookReward 获得挂机奖励
func GetOnHookReward(p *player.Player) (msg.ErrCode, uint32, []*model.SimpleItem) {
	if p.UserData.BaseInfo.HookData == nil {
		return msg.ErrCode_NO_HOOK_REWARD, 0, nil
	}
	settleOnHook(p, p.UserData.StageInfo.MissionId)
	var addItems []*model.SimpleItem
	for i := 0; i < len(p.UserData.BaseInfo.HookData.Items); i++ {
		temp := uint32(p.UserData.BaseInfo.HookData.Items[i].Num)
		if temp > 0 {
			addItems = append(addItems, &model.SimpleItem{
				Id:  p.UserData.BaseInfo.HookData.Items[i].Id,
				Num: temp,
			})
			p.UserData.BaseInfo.HookData.Items[i].Num -= float64(temp)
		}
	}
	if len(addItems) == 0 {
		return msg.ErrCode_NO_HOOK_REWARD, 0, nil
	}

	var notifyItems []uint32
	var resItems = make([]*model.SimpleItem, 0)
	for i := 0; i < len(addItems); i++ {
		items := AddItem(p.GetUserId(), addItems[i].Id, int32(addItems[i].Num), publicconst.OnHookAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(items)...)
		resItems = append(resItems, items...)
	}

	// 通知客户端
	//event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(player, notifyItems, ListenNotifyClientItemEventEvent))

	updateClientItemsChange(p.GetUserId(), notifyItems)

	// 防止主线关卡挂机数据没有
	addOnHookData(p, p.UserData.StageInfo.MissionId)

	temp := p.UserData.BaseInfo.HookData.StartTime
	if p.UserData.BaseInfo.HookData.TotalTime >= template.GetSystemItemTemplate().MaxOnHookTime {
		temp = tools.GetCurTime()
	}

	p.UserData.BaseInfo.HookData.StartTime = tools.GetCurTime()
	p.UserData.BaseInfo.HookData.TotalTime = 0
	p.SaveBaseInfo()

	UpdateTask(p, true, publicconst.TASK_COND_GET_ON_HOOK_REWARD, 1)
	processHistoryData(p, publicconst.TASK_COND_GET_ON_HOOK_REWARD, 0, 1)

	report.ReportHookReward(p.ChannelId, p.GetUserId(), config.Conf.ServerId, 0, ProtocolSimpleItemsToItems(resItems))
	return msg.ErrCode_SUCC, tools.GetCurTime() - temp, resItems
}

// GetQuickOnHookInfo 获得快速挂机信息
func GetQuickOnHookInfo(p *player.Player) (msg.ErrCode, uint32) {
	if p.UserData.BaseInfo.QuickOnHookData == nil {
		p.UserData.BaseInfo.QuickOnHookData = &model.QuickOnHookData{}
		p.SaveBaseInfo()
	}

	nextBuyTime := p.UserData.BaseInfo.QuickOnHookData.NextBuyTime
	curTime := tools.GetCurTime()
	if nextBuyTime > 0 && curTime >= nextBuyTime {
		p.UserData.BaseInfo.QuickOnHookData.BuyTimes = 0
		p.UserData.BaseInfo.QuickOnHookData.NextBuyTime = tools.GetDailyRefreshTime()
		p.SaveBaseInfo()
	}

	return msg.ErrCode_SUCC, p.UserData.BaseInfo.QuickOnHookData.BuyTimes
}

// GetQuickOnHookReward 获得快速挂机奖励
func GetQuickOnHookReward(p *player.Player) (msg.ErrCode, uint32, []*model.Item) {
	nextBuyTime := p.UserData.BaseInfo.QuickOnHookData.NextBuyTime
	curTime := tools.GetCurTime()
	if nextBuyTime > 0 && curTime >= nextBuyTime {
		p.UserData.BaseInfo.QuickOnHookData.BuyTimes = 0
		p.UserData.BaseInfo.QuickOnHookData.NextBuyTime = tools.GetDailyRefreshTime()
		p.SaveBaseInfo()
	}

	buyTimes := p.UserData.BaseInfo.QuickOnHookData.BuyTimes
	if buyTimes >= uint32(len(template.GetSystemItemTemplate().QuickOnHookBuy)) {
		return msg.ErrCode_QUICK_ON_HOOK_TIMES_FULL, 0, nil
	}

	costNum := template.GetSystemItemTemplate().QuickOnHookBuy[buyTimes]

	var notifyClientItem []uint32
	if costNum > 0 {
		if !EnoughItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_AP), uint32(costNum)) {
			return msg.ErrCode_NO_ENOUGH_ITEM, 0, nil
		}
		if res := CostItem(p.GetUserId(),
			uint32(publicconst.ITEM_CODE_AP),
			costNum, publicconst.QuickOnHookCostItem, false); res != msg.ErrCode_SUCC {
			return res, 0, nil
		}
		notifyClientItem = append(notifyClientItem, uint32(publicconst.ITEM_CODE_AP))
	}

	addItems := calcQuickOnHookReward(p, false)
	for i := 0; i < len(addItems); i++ {
		add_items := AddItem(p.GetUserId(), addItems[i].Id, int32(addItems[i].Num), publicconst.QuickOnHookAddItem, false)
		notifyClientItem = append(notifyClientItem, GetSimpleItemIds(add_items)...)
	}

	updateClientItemsChange(p.GetUserId(), notifyClientItem)

	p.UserData.BaseInfo.QuickOnHookData.BuyTimes += 1
	p.UserData.BaseInfo.QuickOnHookData.NextBuyTime = tools.GetDailyRefreshTime()
	p.SaveBaseInfo()

	UpdateTask(p, true, publicconst.TASK_COND_QUICK_ON_HOOK, 1)
	UpdateTask(p, true, publicconst.TASK_COND_GET_ON_HOOK_REWARD, 1)

	// TODO
	//ServMgr.GetCommonService().AddStaticsData(p, publicconst.Statics_Buy_OnHook, "")
	//report.ReportHookReward(p.ChannelId, p.GetUserId(), config.Conf.ServerId, 1, ToProtocolItems(addItems))

	return msg.ErrCode_SUCC, p.UserData.BaseInfo.QuickOnHookData.BuyTimes, addItems
}

// calcQuickOnHookReward 计算快速挂机奖励
func calcQuickOnHookReward(p *player.Player, saveDb bool) []*model.Item {
	var addItems []*model.Item
	// 倍数
	scale := float64(template.GetSystemItemTemplate().QuickOnHookTime*1.0) / float64(template.GetSystemItemTemplate().OnHookInterval)
	scale = scale * float64(GetMonthCardOnHookScale(p)) / float64(100)

	if missionConfig := template.GetMissionTemplate().GetMission(int(p.UserData.StageInfo.MissionId)); missionConfig != nil {
		for k, v := range missionConfig.OnHookReward {
			flag := false
			for i := 0; i < len(p.UserData.BaseInfo.QuickOnHookData.Items); i++ {
				if p.UserData.BaseInfo.QuickOnHookData.Items[i].Id == uint32(k) {
					p.UserData.BaseInfo.QuickOnHookData.Items[i].Num += v * scale
					flag = true
					break
				}
			}
			if !flag {
				p.UserData.BaseInfo.QuickOnHookData.Items = append(p.UserData.BaseInfo.QuickOnHookData.Items, &model.FloatItem{Id: uint32(k), Num: v * scale})
			}
		}
	}

	for _, v := range p.UserData.BaseInfo.QuickOnHookData.Items {
		if v.Num > 0 {

		}
	}
	for i := 0; i < len(p.UserData.BaseInfo.QuickOnHookData.Items); i++ {
		num := uint32(p.UserData.BaseInfo.QuickOnHookData.Items[i].Num)
		if num > 0 {
			addItems = append(addItems, &model.Item{
				Id:  p.UserData.BaseInfo.QuickOnHookData.Items[i].Id,
				Num: uint64(num)})
			p.UserData.BaseInfo.QuickOnHookData.Items[i].Num -= float64(num)
		}
	}

	if saveDb {
		p.SaveBaseInfo()
	}
	return addItems
}

// settleOnHook 结算挂机
func settleOnHook(p *player.Player, newMissionId int) {
	if p.UserData.BaseInfo.HookData == nil && newMissionId > 0 {
		p.UserData.BaseInfo.HookData = &model.OnHookData{
			StartTime: tools.GetCurTime(),
			IsNotice:  true,
		}
		addOnHookData(p, newMissionId)
		p.SaveBaseInfo()
		return
	}

	maxOnHookTime := template.GetSystemItemTemplate().MaxOnHookTime
	if p.UserData.BaseInfo.HookData == nil || p.UserData.BaseInfo.HookData.TotalTime >= maxOnHookTime {
		return
	}

	curTime := tools.GetCurTime()
	deltaTime := curTime - p.UserData.BaseInfo.HookData.StartTime
	if p.UserData.BaseInfo.HookData.IsNotice {
		deltaTime = template.GetSystemItemTemplate().OnHookInterval
	}

	if p.UserData.BaseInfo.HookData.TotalTime+deltaTime > maxOnHookTime {
		deltaTime = maxOnHookTime - p.UserData.BaseInfo.HookData.TotalTime
	}

	if deltaTime > maxOnHookTime {
		log.Error("settleOnHook err", zap.Uint64("accountId", p.GetUserId()),
			zap.Uint32("daltaTime", deltaTime), zap.Uint32("total", p.UserData.BaseInfo.HookData.TotalTime),
			zap.Uint32("startTime", p.UserData.BaseInfo.HookData.StartTime), zap.Uint32("curTime", curTime))
	}

	if deltaTime <= 0 {
		return
	}

	times := float64(deltaTime * 1.0 / template.GetSystemItemTemplate().OnHookInterval)
	times = times * float64(GetMonthCardOnHookScale(p)) / 100

	notifyClient := false
	if p.UserData.StageInfo.MissionId != newMissionId {
		settleCurMission(p, times)
		addOnHookData(p, newMissionId)
		p.UserData.BaseInfo.HookData.StartTime = curTime
		if p.UserData.BaseInfo.HookData.IsNotice {
			p.UserData.BaseInfo.HookData.IsNotice = false
			notifyClient = true
		}
		p.UserData.BaseInfo.HookData.TotalTime += deltaTime
		p.UserData.BaseInfo.HookData.TotalTime = tools.LimitUint32(p.UserData.BaseInfo.HookData.TotalTime, maxOnHookTime)
		p.SaveBaseInfo()
	} else {
		if p.UserData.StageInfo.MissionId == 0 {
			times = 0
		}

		if times <= 0 {
			return
		}

		addTime := template.GetSystemItemTemplate().OnHookInterval * uint32(times)
		if p.UserData.BaseInfo.HookData.IsNotice {
			p.UserData.BaseInfo.HookData.IsNotice = false
			notifyClient = true
		} else {
			p.UserData.BaseInfo.HookData.StartTime += addTime
		}
		p.UserData.BaseInfo.HookData.TotalTime += addTime

		if deltaTime > maxOnHookTime {
			log.Error("settleOnHook err", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("totalTime", p.UserData.BaseInfo.HookData.TotalTime))
			//log.Errorf("settleOnHook accountid:%v, TotalTime:%v", player.GetUserId(), player.AccountInfo.HookData.TotalTime)
		}

		p.UserData.BaseInfo.HookData.TotalTime = tools.LimitUint32(p.UserData.BaseInfo.HookData.TotalTime, maxOnHookTime)
		settleCurMission(p, times)
		p.SaveBaseInfo()
	}

	if notifyClient {
		p.SendNotify(&msg.NotifyOnHookData{
			OnHookTime: p.UserData.BaseInfo.HookData.TotalTime,
		})
	}
	return
}

// addOnHookData 添加挂机数据
func addOnHookData(player *player.Player, missionId int) {
	missionConfig := template.GetMissionTemplate().GetMission(int(missionId))
	if missionConfig == nil {
		log.Error("addOnHookData err", zap.Uint64("accountId", player.GetUserId()),
			zap.Int("missionId", missionId))
		return
	}

	if player.UserData.BaseInfo.HookData == nil {
		return
	}

	for k := range missionConfig.OnHookReward {
		flag := false
		for i := 0; i < len(player.UserData.BaseInfo.HookData.Items); i++ {
			if player.UserData.BaseInfo.HookData.Items[i].Id == uint32(k) {
				flag = true
				break
			}
		}
		if !flag {
			player.UserData.BaseInfo.HookData.Items = append(player.UserData.BaseInfo.HookData.Items, &model.FloatItem{
				Id: uint32(k),
			})
		}
	}
}

// settleCurMission 结算当前关卡
func settleCurMission(player *player.Player, scale float64) {
	tempId := player.UserData.StageInfo.MissionId
	if tempId == 0 {
		tempId = 10001
	}
	missionConfig := template.GetMissionTemplate().GetMission(int(tempId))
	if missionConfig == nil {
		return
	}

	for k, v := range missionConfig.OnHookReward {
		// 特殊处理下第一关赠送时间无奖励问题
		if tempId == 10001 && len(missionConfig.OnHookReward) != len(player.UserData.BaseInfo.HookData.Items) {
			player.UserData.BaseInfo.HookData.Items = append(player.UserData.BaseInfo.HookData.Items, &model.FloatItem{
				Id: uint32(k),
			})
		}

		for i := 0; i < len(player.UserData.BaseInfo.HookData.Items); i++ {
			if player.UserData.BaseInfo.HookData.Items[i].Id == uint32(k) {
				player.UserData.BaseInfo.HookData.Items[i].Num += v * scale
				break
			}
		}
	}
}
