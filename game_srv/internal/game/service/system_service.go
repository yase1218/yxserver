package service

import (
	"gameserver/internal/config"
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

// GetAdventureGuideReward 获取奇遇引导奖励
func GetAdventureGuideReward(p *player.Player, id uint32) (msg.ErrCode, []*model.SimpleItem) {
	jAdventure := template.GetAdventureTemplate().GetAdventure(id)
	if jAdventure == nil {
		return msg.ErrCode_SUCC, nil
	}
	if p.UserData.BaseInfo.Adventures == nil {
		p.UserData.BaseInfo.Adventures = make(map[uint32]*model.AdventureInfo)
	}
	data, ok := p.UserData.BaseInfo.Adventures[id]
	if ok {
		if data.State == 1 {
			return msg.ErrCode_SUCC, nil
		}
	} else {
		data = model.NewAdventureInfo(id, 1)
		p.UserData.BaseInfo.Adventures[id] = data
	}

	var notifyItems []uint32
	var temp []*model.SimpleItem
	for i := 0; i < len(jAdventure.RewardItem); i++ {
		addItems := AddItem(p.GetUserId(),
			jAdventure.RewardItem[i].ItemId,
			int32(jAdventure.RewardItem[i].ItemNum),
			publicconst.AdventureAddItem, false)
		notifyItems = tools.ListUint32AddNoRepeats(notifyItems, GetSimpleItemIds(addItems))
		temp = append(temp, addItems...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)
	p.SaveBaseInfo()
	return msg.ErrCode_SUCC, temp
}

// ConstructionLottery 交互物抽奖
func ConstructionLottery(p *player.Player, id uint32) msg.ErrCode {
	//config := template.GetConstructionTemplate().GetConstruction(id)
	//if config == nil {
	//	return msg.ErrCode_INVALID_DATA
	//}
	//
	//if config.Tp != 4 || len(config.Paras) != 4 {
	//	return msg.ErrCode_INVALID_DATA
	//}
	//
	//if config.Tp == 4 && config.Paras[0] != 36000 {
	//	return msg.ErrCode_INVALID_DATA
	//}
	//
	//costNum := config.Paras[3]
	//if !ServMgr.GetItemService().EnoughItem(playerData.GetUserId(), uint32(publicconst.ITEM_CODE_DIAMOND), uint32(costNum)) {
	//	return msg.ErrCode_NO_ENOUGH_ITEM
	//}
	//
	//var notifyItems []uint32
	//ServMgr.GetItemService().CostItem(playerData.GetUserId(),
	//	uint32(publicconst.ITEM_CODE_DIAMOND),
	//	uint32(costNum),
	//	publicconst.ConstructionLotteryCostItem,
	//	false)
	//notifyItems = append(notifyItems, uint32(publicconst.ITEM_CODE_DIAMOND))
	//// 通知客户端
	//ServMgr.GetItemService().updateClientItemsChange(playerData.AccountInfo.AccountId, notifyItems)
	return msg.ErrCode_SUCC
}

// UpgradeTalent 升级天赋
func UpgradeTalent(p *player.Player, tp uint32) msg.ErrCode {
	if p.UserData.BaseInfo.TalentData == nil {
		p.UserData.BaseInfo.TalentData = model.NewTalentInfo(0, 0)
	}
	var target *template.JTalent
	//var cur *template.JTalent
	if tp == 0 {
		if p.UserData.BaseInfo.TalentData.NormalId == 0 {
			target = template.GetTalentTemplate().InitNoramlTalent
		} else {
			target = template.GetTalentTemplate().GetTalent(p.UserData.BaseInfo.TalentData.NormalId).Next
		}
	} else {
		if p.UserData.BaseInfo.TalentData.KeyId == 0 {
			target = template.GetTalentTemplate().InitKeyTalent
		} else {
			target = template.GetTalentTemplate().GetTalent(p.UserData.BaseInfo.TalentData.KeyId).Next
		}
	}

	if target == nil {
		return msg.ErrCode_TALENT_LEVEL_FULL
	}

	for i := 0; i < len(target.CostItem); i++ {
		if !EnoughItem(p.GetUserId(), target.CostItem[i].ItemId, target.CostItem[i].ItemNum) {
			return msg.ErrCode_NO_ENOUGH_ITEM
		}
	}

	var items []uint32
	for i := 0; i < len(target.CostItem); i++ {
		CostItem(p.GetUserId(), target.CostItem[i].ItemId,
			target.CostItem[i].ItemNum,
			publicconst.TalentUpgradeCostItem, false)
		items = append(items, target.CostItem[i].ItemId)
	}
	updateClientItemsChange(p.GetUserId(), items)
	if tp == 0 {
		p.UserData.BaseInfo.TalentData.NormalId = target.Id
	} else {
		p.UserData.BaseInfo.TalentData.KeyId = target.Id
	}

	InitAttr(p.UserData.BaseInfo.TalentData.Attrs, target.Attr)
	// 计算天赋
	for id, data := range p.UserData.BaseInfo.TalentData.Attrs {
		if idConfig := template.GetAttrListTemplate().GetAttr(id); idConfig != nil {
			data.CalcFinalValue()
		}
	}
	p.UserData.BaseInfo.TalentData.Parts = append(p.UserData.BaseInfo.TalentData.Parts, target.Parts...)

	// 通知天赋变化
	notifyMsg := &msg.NotifyTalentChange{
		Data: ToProtocolTalentData(p.UserData.BaseInfo.TalentData),
	}
	p.SendNotify(notifyMsg)
	p.SaveBaseInfo()
	return msg.ErrCode_SUCC
}

// GetDailyAp 领取每日体力
func GetDailyAp(p *player.Player, id uint32) (msg.ErrCode, uint32) {
	pos := -1
	for i := 0; i < len(p.UserData.BaseInfo.DailyApData); i++ {
		if p.UserData.BaseInfo.DailyApData[i].Id == id {
			pos = i
			break
		}
	}

	if pos == -1 {
		return msg.ErrCode_INVALID_DATA, 0
	}
	curTime := tools.GetCurTime()
	if curTime < p.UserData.BaseInfo.DailyApData[pos].StartTime || curTime > p.UserData.BaseInfo.DailyApData[pos].EndTime {
		return msg.ErrCode_NOT_ARRIVE_GET_DAILY_AP_TIME, 0
	}
	if p.UserData.BaseInfo.DailyApData[pos].State == 1 {
		return msg.ErrCode_HAS_GET_DAILY_AP, 0
	}

	config := template.GetDailyApTemplate().GetDailAp(id)
	curNum := GetItemNum(p.GetUserId(), uint32(publicconst.ITEM_CODE_AP))
	if curNum+uint64(config.Reward[0].ItemNum) > uint64(template.GetSystemItemTemplate().ApMax) {
		return msg.ErrCode_AP_OVER_LIMIT, 0
	}

	p.UserData.BaseInfo.DailyApData[pos].State = 1
	p.SaveBaseInfo()

	var notifyItems []uint32
	for i := 0; i < len(config.Reward); i++ {
		addItems := AddItem(p.GetUserId(),
			config.Reward[i].ItemId,
			int32(config.Reward[i].ItemNum), publicconst.DailApAddItem, false)
		notifyItems = tools.ListUint32AddNoRepeats(notifyItems, GetSimpleItemIds(addItems))
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)
	return msg.ErrCode_SUCC, p.UserData.BaseInfo.DailyApData[pos].State
}

// GetMonthCardDailyReward 领取月卡每日奖励
func GetMonthCardDailyReward(p *player.Player) (msg.ErrCode, uint32) {
	curTime := tools.GetCurTime()
	if curTime < p.UserData.BaseInfo.MonthCardDailyRewardTime {
		return msg.ErrCode_HAS_GET_MONTHCARD_DAILY_REWARD, 0
	}

	reward := template.GetSystemItemTemplate().MonthCardDailyItem
	var notifyItems []uint32
	for i := 0; i < len(reward); i++ {
		addItems := AddItem(p.GetUserId(),
			reward[i].ItemId,
			int32(reward[i].ItemNum), publicconst.MonthCardDailyAddItem, false)
		notifyItems = tools.ListUint32AddNoRepeats(notifyItems, GetSimpleItemIds(addItems))
	}

	updateClientItemsChange(p.GetUserId(), notifyItems)
	p.UserData.BaseInfo.MonthCardDailyRewardTime = tools.GetDailyRefreshTime()
	p.SaveBaseInfo()
	return msg.ErrCode_SUCC, p.UserData.BaseInfo.MonthCardDailyRewardTime
}

func getMainFundReward(p *player.Player, id int) (map[int]int, *model.MainFundInfo) {
	ret := make(map[int]int)
	lst := template.GetFundTemplate().GetFundById(id, int(p.UserData.StageInfo.MissionId))

	var info *model.MainFundInfo
	for i := 0; i < len(p.UserData.BaseInfo.MainFund); i++ {
		if p.UserData.BaseInfo.MainFund[i].Id == id {
			info = p.UserData.BaseInfo.MainFund[i]
			break
		}
	}
	if info == nil {
		info = model.NewMainFundInfo(id)
		p.UserData.BaseInfo.MainFund = append(p.UserData.BaseInfo.MainFund, info)
	}

	for _, v := range lst {
		// 免费的还没领取
		if info.FreeId < v.StageId {
			for _, vv := range v.FreeRewardItems {
				ret[int(vv.ItemId)] += int(vv.ItemNum)
			}
			info.FreeId = v.StageId
		}

		if info.BuyFlag == 1 && info.PayId < v.StageId {
			for _, vv := range v.PayRewardItems {
				ret[int(vv.ItemId)] += int(vv.ItemNum)
			}
			info.PayId = v.StageId
		}
	}
	return ret, info
}

// GetMainFund 获取主线基金
func GetMainFund(p *player.Player, id int) (msg.ErrCode, []*template.SimpleItem, *model.MainFundInfo) {
	var notifyItems []uint32
	var temp []*template.SimpleItem
	itemMap, info := getMainFundReward(p, id)
	if len(itemMap) == 0 {
		return msg.ErrCode_HAS_GET_MAINFUND_REWARD, nil, nil
	}

	for id, num := range itemMap {
		addItems := AddItem(p.GetUserId(), uint32(id), int32(num), publicconst.MainFundAddItem, false)
		notifyItems = tools.ListUint32AddNoRepeats(notifyItems, GetSimpleItemIds(addItems))
		temp = append(temp, ToTemplateItem(addItems)...)
	}

	updateClientItemsChange(p.GetUserId(), notifyItems)
	p.SaveBaseInfo()
	return msg.ErrCode_SUCC, mergeSimpleItem(temp), info
}

// GetMonthCard 获取月卡奖励
func GetMonthCard(p *player.Player) (msg.ErrCode, []*template.SimpleItem) {
	curTime := int(tools.GetCurTime())
	var notifyItems []uint32
	var temp []*template.SimpleItem
	for i := 0; i < len(p.UserData.BaseInfo.MonthCard); i++ {
		if curTime > p.UserData.BaseInfo.MonthCard[i].EndTime {
			continue
		}

		// 可以领取奖励
		if curTime >= p.UserData.BaseInfo.MonthCard[i].NextGetRewardTime {
			config := template.GetMonthCardTemplate().GetMonthCard(p.UserData.BaseInfo.MonthCard[i].Id)
			for m := 0; m < len(config.DailyRewardSlice); m++ {
				addItems := AddItem(p.GetUserId(),
					config.DailyRewardSlice[m].ItemId,
					int32(config.DailyRewardSlice[m].ItemNum), publicconst.MonthCardRewardAddItem, false)
				notifyItems = tools.ListUint32AddNoRepeats(notifyItems, GetSimpleItemIds(addItems))
				temp = append(temp, ToTemplateItem(addItems)...)
			}
			p.UserData.BaseInfo.MonthCard[i].NextGetRewardTime = int(tools.GetDailyRefreshTime())
		}
	}

	if len(temp) == 0 {
		return msg.ErrCode_HAS_GET_MONTHCARD_REWARD, nil
	}

	var finalItems []*template.SimpleItem
	itemMap := make(map[uint32]uint32)
	for i := 0; i < len(temp); i++ {
		if _, ok := itemMap[temp[i].ItemId]; ok {
			itemMap[temp[i].ItemId] += temp[i].ItemNum
		} else {
			itemMap[temp[i].ItemId] = temp[i].ItemNum
		}
	}
	for id, num := range itemMap {
		finalItems = append(finalItems, &template.SimpleItem{
			ItemId:  id,
			ItemNum: num,
		})
	}

	updateClientItemsChange(p.GetUserId(), notifyItems)
	p.SaveBaseInfo()

	return msg.ErrCode_SUCC, finalItems
}

// GetQuestionReward 获得问卷调查奖励
func GetQuestionReward(p *player.Player, questionId string) (msg.ErrCode, []*template.SimpleItem) {
	questionInfo := template.GetSystemItemTemplate().GetQuestionInfo(questionId)
	if questionInfo == nil {
		return msg.ErrCode_INVALID_DATA, nil
	}

	// if !CanUnlockOneCond(playerData, questionInfo.Cond) {
	// 	return msg.ErrCode_FUNCTION_NOT_OPEN, nil
	// }

	if tools.ListStrContain(p.UserData.BaseInfo.RewardQuestionIds, questionId) {
		return msg.ErrCode_INVALID_DATA, nil
	}

	p.UserData.BaseInfo.RewardQuestionIds = append(p.UserData.BaseInfo.RewardQuestionIds, questionId)
	p.SaveBaseInfo()

	var notifyItems []uint32
	var finalItems []*template.SimpleItem
	for i := 0; i < len(questionInfo.Reward); i++ {
		addItems := AddItem(p.GetUserId(), questionInfo.Reward[i].ItemId,
			int32(questionInfo.Reward[i].ItemNum), publicconst.QuestionRewardAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		finalItems = append(finalItems, ToTemplateItem(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)
	return msg.ErrCode_SUCC, finalItems
}

func SetPopUp(p *player.Player, id, popType uint32) msg.ErrCode {
	for i := 0; i < len(p.UserData.BaseInfo.PopUps); i++ {
		if p.UserData.BaseInfo.PopUps[i].Id == id && p.UserData.BaseInfo.PopUps[i].PopUpType == popType {
			return msg.ErrCode_SYSTEM_ERROR
		}
	}

	info := &model.PopUpInfo{
		Id:        id,
		PopUpType: popType,
	}
	p.UserData.BaseInfo.PopUps = append(p.UserData.BaseInfo.PopUps, info)
	p.SaveBaseInfo()
	return msg.ErrCode_SUCC
}

// SetVideo 设置引导信息
func SetVideo(p *player.Player) msg.ErrCode {
	if p.UserData.BaseInfo.VideoFlag == 1 {
		return msg.ErrCode_SUCC
	}
	p.UserData.BaseInfo.VideoFlag = 1
	p.SaveBaseInfo()
	return msg.ErrCode_SUCC
}

// SaveBattleData 保存战斗数据
func SaveBattleData(p *player.Player, data *msg.BattleData) msg.ErrCode {
	p.UserData.BaseInfo.BattleData = data
	p.SaveBaseInfo()
	return msg.ErrCode_SUCC
}

// SetGuideInfo 设置引导信息
func SetGuideInfo(p *player.Player, id, value uint32) msg.ErrCode {
	pos := -1
	for i := 0; i < len(p.UserData.BaseInfo.GuideData); i++ {
		if p.UserData.BaseInfo.GuideData[i].Id == id {
			pos = i
			break
		}
	}

	if pos >= 0 {
		p.UserData.BaseInfo.GuideData[pos].Value = value
	} else {
		p.UserData.BaseInfo.GuideData = append(p.UserData.BaseInfo.GuideData, &model.GuideInfo{
			Id:    id,
			Value: value,
		})
	}
	p.SaveBaseInfo()

	// // 统计记录
	// ServMgr.GetCommonService().AddStaticsData(p, publicconst.Guide_Id, fmt.Sprintf("%v.%v", id, value))

	// guideStatics := model.NewGuideStep(tools.GetStaticTime(tools.GetCurTime()), uint32(p.GetUserId()), p.AccountInfo.ChannelId, id, value)
	// dao.UserStaticDao.AddGuidStep(guideStatics)
	return msg.ErrCode_SUCC
}

// InterUseCdkResponse 使用cdk
func InterUseCdkResponse(req *msg.InterResponseUseCdk, p *player.Player) {
	retMsg := &msg.ResponseUseCdk{}
	retMsg.Result = req.Result
	if req.Result == msg.ErrCode_SUCC {
		var notifyItems []uint32
		rewardItem := make(map[uint32]uint32)
		for i := 0; i < len(req.CdkRet); i++ {
			retMsg.CdkRet = append(retMsg.CdkRet, &msg.CdkResult{
				Cdk:    req.CdkRet[i].Cdk,
				Result: req.CdkRet[i].Result,
			})
			if req.CdkRet[i].Result == msg.ErrCode_SUCC {
				// items := template.GetKeyValueInt2Int(req.CdkRet[i].Items)
				// for k := 0; k < len(items); k++ {
				// 	itemId := uint32(items[k][0])
				// 	itemNum := uint32(items[k][1])
				// 	addItems := AddItem(p.GetUserId(),
				// 		itemId, int32(itemNum), publicconst.CdkAddItem, false)
				// 	notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
				// 	if _, ok := rewardItem[itemId]; ok {
				// 		rewardItem[itemId] += itemNum
				// 	} else {
				// 		rewardItem[itemId] = itemNum
				// 	}
				// }
				//dao.UserStaticDao.AddCdkRecord(model.NewCdkRecord(req.CdkRet[i].Cdk, p.GetUserId(), p.AccountInfo.Nick))
			}
		}

		for itemId, num := range rewardItem {
			retMsg.Data = append(retMsg.Data, &msg.SimpleItem{
				ItemId:  itemId,
				ItemNum: num,
			})
		}

		if len(notifyItems) > 0 {
			//event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(playerData, notifyItems, ListenNotifyClientItemEventEvent))
			updateClientItemsChange(p.GetUserId(), notifyItems)
		}
	}
	p.SendNotify(retMsg)
}

func InterEndAd(arg interface{}, p *player.Player) {
	AdCallBack(p)
}

func UseCdk(p *player.Player, cdk []string, packetId uint32) msg.ErrCode {
	if len(cdk) == 0 || len(cdk) > 10 {
		return msg.ErrCode_INVALID_DATA
	}
	for i := 0; i < len(cdk); i++ {
		if len(cdk[i]) == 0 {
			return msg.ErrCode_INVALID_DATA
		}
	}

	log.Debug("UseCdk", zap.Any("cdk", cdk))
	// 发送到cdk服务器
	interMsg := &msg.InterRequestUseCdk{}
	interMsg.Cdk = cdk
	interMsg.AccountId = int64(p.GetUserId())
	interMsg.PacketId = packetId
	PublisInterMsg(config.Conf.CdkServerId, uint32(msg.InterMsgId_ID_InterRequestUseCdk), interMsg)
	return msg.ErrCode_SUCC
}

// BuyAp 购买体力
func BuyAp(p *player.Player) (msg.ErrCode, uint32, []*model.SimpleItem) {
	curTime := tools.GetCurTime()
	accountId := p.GetUserId()

	// 次数到达上限
	buyTimes := p.UserData.BaseInfo.ApData.BuyTimes
	if buyTimes == uint32(len(template.GetSystemItemTemplate().BuyApCostDiamond)) {
		return msg.ErrCode_OVER_BUY_AP_TIMES, 0, nil
	}

	// 体力到达上限
	if EnoughItem(accountId, uint32(publicconst.ITEM_CODE_AP), template.GetSystemItemTemplate().ApMax) {
		return msg.ErrCode_AP_OVER_LIMIT, 0, nil
	}

	// 没有足够的道具
	itemNum := template.GetSystemItemTemplate().BuyApCostDiamond[buyTimes]
	if res := CostItem(accountId, uint32(publicconst.ITEM_CODE_DIAMOND), itemNum, publicconst.BuyApCostItem, false); res != msg.ErrCode_SUCC {
		return res, 0, nil
	}
	addItems := AddItem(accountId, uint32(publicconst.ITEM_CODE_AP), int32(template.GetSystemItemTemplate().BuyGetAp), publicconst.BuyApAddItem, false)

	// 通知客户端
	var itemIds []uint32
	itemIds = append(itemIds, uint32(publicconst.ITEM_CODE_DIAMOND))
	itemIds = append(itemIds, GetSimpleItemIds(addItems)...)

	// 通知客户端
	updateClientItemsChange(p.GetUserId(), itemIds)

	// 需要重置
	if curTime >= p.UserData.BaseInfo.ApData.NextBuyTime {
		p.UserData.BaseInfo.ApData.NextBuyTime = tools.GetDailyRefreshTime()
		p.UserData.BaseInfo.ApData.BuyTimes = 0
	}

	p.UserData.BaseInfo.ApData.BuyTimes += 1
	p.UserData.BaseInfo.ApData.NextBuyTime = tools.GetDailyRefreshTime()
	p.SaveBaseInfo()
	UpdateTask(p, true, publicconst.TASK_COND_BUY_AP, 1)

	// // 统计记录
	// para := fmt.Sprintf("%v", p.AccountInfo.ApData.BuyTimes)
	// ServMgr.GetCommonService().AddStaticsData(p, publicconst.Statics_Buy_Ap_Id, para)

	return msg.ErrCode_SUCC, p.UserData.BaseInfo.ApData.BuyTimes, addItems
}

// recoveryAp 恢复体力
func recoveryAp(p *player.Player, now time.Time) {
	levelConfig := template.GetLevelTemplate().GetLevelData(p.UserData.Level)
	if levelConfig == nil {
		return
	}

	if len(p.UserData.Items.Items) == 0 {
		return
	}

	uid := p.GetUserId()
	ap_data := p.UserData.BaseInfo.ApData
	// 体力已满
	if EnoughItem(uid, uint32(publicconst.ITEM_CODE_AP), levelConfig.ApMax) {
		if ap_data.RecoverStartTime != 0 {
			ap_data.RecoverStartTime = 0
			//dao.AccountDao.UpdateApRecoveryStartTime(accountId, ap_data.RecoverStartTime)
		}
		return
	}
	curTime := uint32(now.Unix())
	if ap_data.RecoverStartTime == 0 {
		ap_data.RecoverStartTime = curTime
		//dao.AccountDao.UpdateApRecoveryStartTime(accountId, ap_data.RecoverStartTime)
		NotifyApRecovery(p)
	} else {
		interval := (curTime - ap_data.RecoverStartTime) / template.GetSystemItemTemplate().ApRecoveryInterval
		// 结算体力
		if interval > 0 {
			ap_data.RecoverStartTime += interval * template.GetSystemItemTemplate().ApRecoveryInterval
			//dao.AccountDao.UpdateApRecoveryStartTime(accountId, ap_data.RecoverStartTime)

			addValue := int32(interval * template.GetSystemItemTemplate().ApRecoveryValue)
			curValue := int32(GetItemNum(uid, uint32(publicconst.ITEM_CODE_AP)))
			if curValue+addValue > int32(levelConfig.ApMax) {
				addValue = int32(levelConfig.ApMax) - curValue
			}

			if addValue > 0 {
				AddItem(uid, uint32(publicconst.ITEM_CODE_AP), addValue, publicconst.RecoveryItem, true)
			}

			NotifyApRecovery(p)
		}
	}
}

func ResetDailyAp(p *player.Player, notifyClient bool) {
	if len(p.UserData.BaseInfo.DailyApData) > 0 {
		curTime := tools.GetCurTime()
		if curTime <= p.UserData.BaseInfo.DailyApData[0].EndTime {
			return
		}
		p.UserData.BaseInfo.DailyApData = p.UserData.BaseInfo.DailyApData[0:0]
	}
	lst := template.GetDailyApTemplate().GetAllDailAp()
	dailyStart := tools.GetDateStart(time.Now())

	refreshHour := template.GetSystemItemTemplate().RefreshHour
	for i := 0; i < len(lst); i++ {
		start := uint32(dailyStart.Add(time.Minute * time.Duration(lst[i].StartMin)).Unix())
		end := uint32(dailyStart.AddDate(0, 0, 1).Add(time.Hour * time.Duration(refreshHour)).Unix())
		p.UserData.BaseInfo.DailyApData = append(p.UserData.BaseInfo.DailyApData,
			model.NewDailApInfo(lst[i].Id, start, end))
	}
	p.SaveBaseInfo()
	if notifyClient {
		retMsg := &msg.NotifyDailyApInfo{}
		retMsg.Data = ToProtocolDailApInfo(p.UserData.BaseInfo.DailyApData)
		p.SendNotify(retMsg)
	}
}

func ToProtocolDailApInfo(data []*model.DailyApInfo) []*msg.DailyApInfo {
	var ret []*msg.DailyApInfo
	for i := 0; i < len(data); i++ {
		ret = append(ret, &msg.DailyApInfo{
			Id:        data[i].Id,
			StartTime: data[i].StartTime,
			EndTime:   data[i].EndTime,
			State:     data[i].State,
		})
	}
	return ret
}

func ToProtocolTalentData(data *model.TalentInfo) *msg.TalentInfo {
	if data == nil {
		return nil
	}
	var ret = &msg.TalentInfo{}
	ret.NormalId = data.NormalId
	ret.KeyId = data.KeyId
	ret.Parts = data.Parts

	var attrs []*model.Attr
	for _, d := range data.Attrs {
		attrs = append(attrs, d)
	}
	ret.Attrs = ToProtocolAttrs(attrs)
	return ret
}
