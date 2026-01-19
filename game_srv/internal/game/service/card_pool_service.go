package service

import (
	"fmt"
	"gameserver/internal/enum"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tda"
	"kernel/tools"
	"math/rand"
	"msg"
	"strconv"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

// lotteryTen 10连抽
func lotteryTen(p *player.Player, cardInfo *model.CardPool, cardConfig *template.JLotteryShip, times uint32) (msg.ErrCode, []*template.SimpleItem, []uint32) {
	var notifyItems []uint32
	notifyItems = append(notifyItems, cardConfig.CostItemId)
	//tdaUseItems := make([]*tda.Item, 0, 1)

	// 扣除道具
	if res := CostItem(p.GetUserId(), cardConfig.CostItemId, times, publicconst.CostTenLottery, false); res != msg.ErrCode_SUCC {
		return res, nil, nil
	}
	//tdaUseItems = append(tdaUseItems, &tda.Item{ItemId: strconv.Itoa(int(cardConfig.CostItemId)), ItemNum: times})

	//log.Debug("lottery ten times", zap.Uint32("cardCfgId", cardConfig.Id), zap.Uint32("BigQuaranteeTimes",
	//	cardConfig.BigQuaranteeTimes), zap.Uint32("TenRewardTimes", cardConfig.TenRewardTimes))

	var ret []*template.SimpleItem
	for i := 0; i < int(times); i++ {
		cardInfo.TenQuaranteeTimes += 1
		cardInfo.BigQuaranteeTimes += 1

		itemGroupId := cardConfig.GetLotteryPool(cardInfo.BigQuaranteeTimes)
		itemGroup := template.GetItemGroupTemplate().GetItemGroup(itemGroupId)
		if itemGroup == nil {
			log.Error("lottery not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("itemGroupId", itemGroupId))
			//log.Errorf("accountid:%v lotteryTen itemgroup:%v not exist", playerData.GetUserId(), itemGroupId)
			return msg.ErrCode_SYSTEM_ERROR, nil, nil
		}

		items := itemGroup.GenItem(p.Rand, 1)
		//log.Debug("---lotteryTen", zap.Uint32("itemGroupId", itemGroupId), zap.Any("items", items))

		// 触发保底重置次数
		//if cardConfig.IsBigQuaranteeItem(items) {
		//	cardInfo.TenQuaranteeTimes = 0
		//	cardInfo.BigQuaranteeTimes = 0
		//	if cardConfig.CardType == 2 {
		//		cardInfo.Status = 0
		//	}
		//} else if cardConfig.IsTenQuaranteeItem(items) {
		//	cardInfo.TenQuaranteeTimes = 0
		//}

		// 触发保底
		if cardInfo.BigQuaranteeTimes == cardConfig.BigQuaranteeTimes {
			bigGroup := template.GetItemGroupTemplate().GetItemGroup(cardConfig.BigQuaranteePool)
			items = bigGroup.GenItem(p.Rand, 1)
			cardInfo.BigQuaranteeTimes = 0
			cardInfo.TenQuaranteeTimes = 0
			if cardConfig.CardType == 2 {
				cardInfo.Status = 0
			}
			log.Debug("big", zap.Uint32("pool", cardConfig.BigQuaranteePool), zap.Any("item", items))
		} else if cardInfo.TenQuaranteeTimes == cardConfig.TenRewardTimes {
			tenGroup := template.GetItemGroupTemplate().GetItemGroup(cardConfig.TenQuarantePool)
			items = tenGroup.GenItem(p.Rand, 1)
			cardInfo.TenQuaranteeTimes = 0
			log.Debug("ten", zap.Uint32("pool", cardConfig.TenQuarantePool), zap.Any("item", items))
		}
		ret = append(ret, items...)
	}
	cardInfo.TotalTimes += times
	cardInfo.LotteryTotalTimes += times
	p.SaveCardPool()

	var finalItems []*template.SimpleItem
	tdaGainItems := make([]*tda.Item, 0, len(ret))
	for i := 0; i < len(ret); i++ {
		addItems := AddItem(p.GetUserId(), ret[i].ItemId, int32(ret[i].ItemNum), publicconst.AddTenLottery, false)

		//if len(addItems) > 0 {
		//	finalItems = append(finalItems, ToTemplateItem(addItems)...)
		//} else {
		//	finalItems = append(finalItems, &template.SimpleItem{ItemId: ret[i].ItemId, ItemNum: ret[i].ItemNum})
		//}

		templateItems := ToTemplateItem(addItems)
		finalItems = append(finalItems, templateItems...)

		realItemIds := GetSimpleItemIds(addItems)
		for m := 0; m < len(realItemIds); m++ {
			notifyItems = tools.ListUint32AddNoRepeat(notifyItems, realItemIds[m])
		}

		tdaGainItems = append(tdaGainItems, &tda.Item{ItemId: strconv.Itoa(int(ret[i].ItemId)), ItemNum: ret[i].ItemNum})
	}

	UpdateTask(p, true, publicconst.TASK_COND_LOTTERY, 10)
	processHistoryData(p, publicconst.TASK_COND_LOTTERY, 0, 10)

	// tda
	//tda.TdaLotteryResult(p.ChannelId, p.TdaCommonAttr, "十连抽", cardConfig.CardType, 10, tdaGainItems, tdaUseItems)

	return msg.ErrCode_SUCC, finalItems, notifyItems
}

// lotteryOne 单抽
func lotteryOne(p *player.Player, cardInfo *model.CardPool, cardConfig *template.JLotteryShip) (msg.ErrCode, []*template.SimpleItem, []uint32) {
	var notifyItems []uint32
	//tdaUseItems := make([]*tda.Item, 0, 1)

	if cardInfo.FreeTimes > 0 {
		cardInfo.FreeTimes -= 1
	} else {
		notifyItems = append(notifyItems, cardConfig.CostItemId)
		if res := CostItem(p.GetUserId(), cardConfig.CostItemId, 1, publicconst.CostOneLottery, false); res != msg.ErrCode_SUCC {
			return res, nil, nil
		}
		//tdaUseItems = append(tdaUseItems, &tda.Item{ItemId: strconv.Itoa(int(cardConfig.CostItemId)), ItemNum: 1})
	}

	// 扣除道具
	itemGroupId := cardConfig.GetLotteryPool(cardInfo.BigQuaranteeTimes + 1)
	log.Debug("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~", zap.Uint64("accountId",
		p.GetUserId()),
		zap.Uint32("times ", cardInfo.BigQuaranteeTimes+1),
		zap.Uint32("itemGroupId", itemGroupId),
	)
	itemGroup := template.GetItemGroupTemplate().GetItemGroup(itemGroupId)
	if itemGroup == nil {
		log.Error("lottery not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("itemGroupId", itemGroupId))
		//log.Errorf("accountid:%v lotteryTen itemgroup:%v not exist", playerData.GetUserId(), itemGroupId)
		return msg.ErrCode_SYSTEM_ERROR, nil, nil
	}

	var ret []*template.SimpleItem
	cardInfo.TenQuaranteeTimes += 1
	cardInfo.BigQuaranteeTimes += 1
	items := itemGroup.GenItem(p.Rand, 1)

	// 需要触发大保底
	if cardInfo.BigQuaranteeTimes == cardConfig.BigQuaranteeTimes {
		bigGroup := template.GetItemGroupTemplate().GetItemGroup(cardConfig.BigQuaranteePool)
		items = bigGroup.GenItem(p.Rand, 1)
		cardInfo.BigQuaranteeTimes = 0
		cardInfo.TenQuaranteeTimes = 0
		if cardConfig.CardType == 2 {
			cardInfo.Status = 0
		}
	} else if cardInfo.TenQuaranteeTimes == cardConfig.TenRewardTimes {
		tenGroup := template.GetItemGroupTemplate().GetItemGroup(cardConfig.TenQuarantePool)
		items = tenGroup.GenItem(p.Rand, 1)
		cardInfo.TenQuaranteeTimes = 0
	}

	cardInfo.TotalTimes += 1
	cardInfo.LotteryTotalTimes += 1
	p.SaveCardPool()

	ret = append(ret, items...)

	var finalItems []*template.SimpleItem
	tdaGainItems := make([]*tda.Item, 0, len(ret))
	for i := 0; i < len(ret); i++ {
		addItems := AddItem(p.GetUserId(),
			ret[i].ItemId, int32(ret[i].ItemNum), publicconst.AddOneLottery, false)

		templateItesm := ToTemplateItem(addItems)
		finalItems = append(finalItems, templateItesm...)
		realItemIds := GetSimpleItemIds(addItems)
		for m := 0; m < len(realItemIds); m++ {
			notifyItems = tools.ListUint32AddNoRepeat(notifyItems, realItemIds[m])
		}

		tdaGainItems = append(tdaGainItems, &tda.Item{ItemId: strconv.Itoa(int(ret[i].ItemId)), ItemNum: ret[i].ItemNum})
	}

	UpdateTask(p, true, publicconst.TASK_COND_LOTTERY, 1)
	processHistoryData(p, publicconst.TASK_COND_LOTTERY, 0, 1)

	// tda
	//tda.TdaLotteryResult(p.ChannelId, p.TdaCommonAttr, "单抽", cardConfig.CardType, 1, tdaGainItems, tdaUseItems)

	return msg.ErrCode_SUCC, finalItems, notifyItems
}

// firstLottery 首10抽
func firstLotteryTen(p *player.Player, cardInfo *model.CardPool) (msg.ErrCode, []*template.SimpleItem) {
	var notifyItems []uint32

	lotteryTimes := uint32(10)

	cardConfig := template.GetLotteryShipTemplate().GetCardPool(cardInfo.CardId)
	if cardConfig == nil {
		return msg.ErrCode_CARD_POOL_NOT_EXIST, nil
	}

	for i := uint32(0); i < lotteryTimes; i++ {
		cardInfo.TenQuaranteeTimes++
		cardInfo.BigQuaranteeTimes++

		if cardInfo.BigQuaranteeTimes == cardConfig.BigQuaranteeTimes {
			cardInfo.BigQuaranteeTimes = 0
		} else if cardInfo.TenQuaranteeTimes == cardConfig.TenRewardTimes {
			cardInfo.TenQuaranteeTimes = 0
		}
	}

	if cardInfo.TotalTimes >= cardConfig.DailyTimes ||
		(cardInfo.TotalTimes+lotteryTimes > cardConfig.DailyTimes) {
		return msg.ErrCode_DAILY_LOTTERY_TIMES_FULL, nil
	}

	num := uint32(GetItemNum(p.GetUserId(), cardConfig.CostItemId))
	if num < lotteryTimes {
		shopConfig := template.GetShopTemplate().GetShopItem(cardConfig.ShopId)
		if shopConfig == nil {
			return msg.ErrCode_NO_ENOUGH_ITEM, nil
		}
		if shopConfig.ChargeID > 0 || shopConfig.CostAd > 0 {
			return msg.ErrCode_BUY_SHOP_NEED_PAY, nil
		}
		if err, _, _, _ := BuyShopItem(p, cardConfig.ShopId, uint32(lotteryTimes-num)); err != msg.ErrCode_SUCC {
			return err, nil
		}
	}

	CostItem(p.GetUserId(), cardConfig.CostItemId, lotteryTimes, publicconst.CostTenLottery, true)

	cardInfo.FirstLotteryTen = 1
	cardInfo.TotalTimes += lotteryTimes
	cardInfo.LotteryTotalTimes += lotteryTimes

	// 前9次抽的卡池
	itemGroupNine := template.GetItemGroupTemplate().GetItemGroup(template.GetSystemItemTemplate().FirstLotteryItemGroup[0])
	if itemGroupNine == nil {
		log.Error("lottery cfg nil", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("cfgId", template.GetSystemItemTemplate().FirstLotteryItemGroup[0]))
		return msg.ErrCode_SYSTEM_ERROR, nil
	}
	// 第10次抽的卡池
	itemGroupTen := template.GetItemGroupTemplate().GetItemGroup(template.GetSystemItemTemplate().FirstLotteryItemGroup[1])
	if itemGroupTen == nil {
		log.Error("lottery cfg nil", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("cfgId", template.GetSystemItemTemplate().FirstLotteryItemGroup[1]))
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	var ret []*template.SimpleItem
	for i := 0; i < 9; i++ {
		ret = append(ret, itemGroupNine.GenItem(p.Rand, 1)...)
	}
	ret = append(ret, itemGroupTen.GenItem(p.Rand, 1)...)

	p.SaveCardPool()

	var finalItems []*template.SimpleItem
	for i := 0; i < len(ret); i++ {
		addItems := AddItem(p.GetUserId(),
			ret[i].ItemId, int32(ret[i].ItemNum), publicconst.AddFirstLottery, false)
		finalItems = append(finalItems, ret[i])

		realItemIds := GetSimpleItemIds(addItems)
		for m := 0; m < len(realItemIds); m++ {
			notifyItems = tools.ListUint32AddNoRepeat(notifyItems, realItemIds[m])
		}

		itemId := ret[i].ItemId
		//para += fmt.Sprintf("%v,%v|", itemId, ret[i].ItemNum)

		if template.GetTickerTemplate().CheckLotteryTicker(int(itemId)) {
			bannerMsg := &msg.NotifyBanner{}
			bannerMsg.BtType = msg.BannerType_Banner_Game
			bannerMsg.Content = fmt.Sprintf("%d", template.LotteryTickerID)
			bannerMsg.Params = append(bannerMsg.Params, p.UserData.Nick)
			bannerMsg.Params = append(bannerMsg.Params, fmt.Sprintf("%s_%d", template.LOTTERY_SHIP_FILE, cardInfo.CardId))
			bannerMsg.Params = append(bannerMsg.Params, fmt.Sprintf("%s_%d", template.ITEM_FILE, itemId))
			BoadCastMsg(bannerMsg)
		}
	}

	updateClientItemsChange(p.GetUserId(), notifyItems)

	UpdateTask(p, true, publicconst.TASK_COND_LOTTERY, lotteryTimes)
	processHistoryData(p, publicconst.TASK_COND_LOTTERY, 0, lotteryTimes)

	// 通知奖池变化
	notifyMsg := &msg.NotifyCardPoolInfoChange{}
	notifyMsg.Data = append(notifyMsg.Data, ToProtocolCardPool(cardInfo))
	p.SendNotify(notifyMsg)
	return msg.ErrCode_SUCC, ret
}

func firstLottery(p *player.Player, cardInfo *model.CardPool) (msg.ErrCode, []*template.SimpleItem) {
	var notifyItems []uint32

	lotteryTimes := uint32(1)

	cardConfig := template.GetLotteryShipTemplate().GetCardPool(cardInfo.CardId)
	if cardConfig == nil {
		return msg.ErrCode_CARD_POOL_NOT_EXIST, nil
	}

	for i := uint32(0); i < lotteryTimes; i++ {
		cardInfo.TenQuaranteeTimes++
		cardInfo.BigQuaranteeTimes++

		if cardInfo.BigQuaranteeTimes == cardConfig.BigQuaranteeTimes {
			cardInfo.BigQuaranteeTimes = 0
		} else if cardInfo.TenQuaranteeTimes == cardConfig.TenRewardTimes {
			cardInfo.TenQuaranteeTimes = 0
		}
	}

	if cardInfo.TotalTimes >= cardConfig.DailyTimes ||
		(cardInfo.TotalTimes+lotteryTimes > cardConfig.DailyTimes) {
		return msg.ErrCode_DAILY_LOTTERY_TIMES_FULL, nil
	}

	cardInfo.FreeTimes -= 1
	cardInfo.FirstLottery = 1
	cardInfo.TotalTimes += lotteryTimes
	cardInfo.LotteryTotalTimes += lotteryTimes

	itemGroup := template.GetItemGroupTemplate().GetItemGroup(template.GetSystemItemTemplate().FirstLotteryOnceItemGroup)
	if itemGroup == nil {
		log.Error("lottery cfg nil", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("cfgId", template.GetSystemItemTemplate().FirstLotteryItemGroup[0]))
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	playerShips := p.UserData.Ships.Ships
	playerShipIds := make(map[int32]bool)
	for _, ship := range playerShips {
		playerShipIds[int32(ship.Id)] = true
	}

	tempItems := make([]*template.SimpleItem, 0)
	for _, item := range itemGroup.Items {
		shipCfgId := template.GetItemTemplate().GetItem(item.ItemId).EffectArgs[0]
		if !playerShipIds[int32(shipCfgId)] && shipCfgId > 0 {
			tempItems = append(tempItems, &template.SimpleItem{
				ItemId:  item.ItemId,
				ItemNum: item.ItemNum,
			})
		}
	}

	var randomIndex int = -1
	var ret []*template.SimpleItem
	if len(tempItems) > 0 {
		randomIndex = rand.Intn(len(tempItems))
		ret = append(ret, tempItems[randomIndex])
	} else {
		randomIndex = rand.Intn(len(itemGroup.Items))
		ret = append(ret, &template.SimpleItem{
			ItemId:  itemGroup.Items[randomIndex].ItemId,
			ItemNum: itemGroup.Items[randomIndex].ItemNum,
		})
	}

	p.SaveCardPool()

	//para := fmt.Sprintf("type:%v,id:%v|", lotteryTimes, cardInfo.CardId)
	//var finalItems []*template.SimpleItem
	for i := 0; i < len(ret); i++ {
		addItems := AddItem(p.GetUserId(),
			ret[i].ItemId, int32(ret[i].ItemNum), publicconst.AddFirstLottery, false)
		//finalItems = append(finalItems, ret[i])

		realItemIds := GetSimpleItemIds(addItems)
		for m := 0; m < len(realItemIds); m++ {
			notifyItems = tools.ListUint32AddNoRepeat(notifyItems, realItemIds[m])
		}

		itemId := ret[i].ItemId
		//para += fmt.Sprintf("%v,%v|", itemId, ret[i].ItemNum)

		if template.GetTickerTemplate().CheckLotteryTicker(int(itemId)) {
			bannerMsg := &msg.NotifyBanner{}
			bannerMsg.BtType = msg.BannerType_Banner_Game
			bannerMsg.Content = fmt.Sprintf("%d", template.LotteryTickerID)
			bannerMsg.Params = append(bannerMsg.Params, p.UserData.Nick)
			bannerMsg.Params = append(bannerMsg.Params, fmt.Sprintf("%s_%d", template.LOTTERY_SHIP_FILE, cardInfo.CardId))
			bannerMsg.Params = append(bannerMsg.Params, fmt.Sprintf("%s_%d", template.ITEM_FILE, itemId))
			BoadCastMsg(bannerMsg)
		}
	}

	updateClientItemsChange(p.GetUserId(), notifyItems)

	UpdateTask(p, true, publicconst.TASK_COND_LOTTERY, lotteryTimes)
	processHistoryData(p, publicconst.TASK_COND_LOTTERY, 0, lotteryTimes)

	UpdateTask(p, true, publicconst.TASK_COND_SHIP_LOTTERY, lotteryTimes) // 累计进行XX次库鲁招募
	processHistoryData(p, publicconst.TASK_COND_SHIP_LOTTERY, 0, lotteryTimes)

	// 通知奖池变化
	notifyMsg := &msg.NotifyCardPoolInfoChange{}
	notifyMsg.Data = append(notifyMsg.Data, ToProtocolCardPool(cardInfo))
	p.SendNotify(notifyMsg)
	return msg.ErrCode_SUCC, ret
}

// Lottery 抽奖
func Lottery(p *player.Player, cardId, lotteryTimes uint32) (msg.ErrCode, []*template.SimpleItem) {
	if lotteryTimes == 0 {
		return msg.ErrCode_INVALID_DATA, nil
	}
	cardInfo := getCardPool(p, cardId)
	if cardInfo == nil || cardInfo.Status == 0 {
		return msg.ErrCode_CARD_POOL_NOT_EXIST, nil
	}

	//if cardInfo.CardId == 100 && lotteryTimes == 10 && cardInfo.FirstLotteryTen == 0 {
	//	return c.firstLotteryTen(playerData, cardInfo)
	//} else
	if cardInfo.CardId == 100 && lotteryTimes == 1 && cardInfo.FirstLottery == 0 {
		return firstLottery(p, cardInfo)
	}

	ResetCardPool(p, false)
	curTime := tools.GetCurTime()
	if cardInfo.StartTime > 0 && cardInfo.EndTime > 0 && curTime > cardInfo.EndTime {
		return msg.ErrCode_CARD_POOL_NOT_EXIST, nil
	}

	if cardInfo.LotteryEndTime > 0 && curTime > cardInfo.LotteryEndTime {
		return msg.ErrCode_LOTTERY_END, nil
	}

	cardConfig := template.GetLotteryShipTemplate().GetCardPool(cardId)
	if cardConfig == nil {
		return msg.ErrCode_CARD_POOL_NOT_EXIST, nil
	}

	// lotterType 0 单抽 1 10连抽 2 n次抽
	if cardInfo.TotalTimes >= cardConfig.DailyTimes ||
		(cardInfo.TotalTimes+lotteryTimes > cardConfig.DailyTimes) {
		return msg.ErrCode_DAILY_LOTTERY_TIMES_FULL, nil
	}

	if lotteryTimes == 1 {
		// 没有免费次数
		if cardInfo.FreeTimes == 0 {
			if !EnoughItem(p.GetUserId(), cardConfig.CostItemId, 1) {
				if cardConfig.ShopId == 0 {
					return msg.ErrCode_NO_ENOUGH_ITEM, nil
				}
				shopConfig := template.GetShopTemplate().GetShopItem(cardConfig.ShopId)
				if shopConfig == nil {
					return msg.ErrCode_NO_ENOUGH_ITEM, nil
				}
				if shopConfig.ChargeID > 0 || shopConfig.CostAd > 0 {
					return msg.ErrCode_BUY_SHOP_NEED_PAY, nil
				}
				if err, _, _, _ := BuyShopItem(p, cardConfig.ShopId, 1); err != msg.ErrCode_SUCC {
					return err, nil
				}
			}
		}
	} else {
		num := uint32(GetItemNum(p.GetUserId(), cardConfig.CostItemId))
		if num < lotteryTimes {
			shopConfig := template.GetShopTemplate().GetShopItem(cardConfig.ShopId)
			if shopConfig == nil {
				return msg.ErrCode_NO_ENOUGH_ITEM, nil
			}
			if shopConfig.ChargeID > 0 || shopConfig.CostAd > 0 {
				return msg.ErrCode_BUY_SHOP_NEED_PAY, nil
			}
			if err, _, _, _ := BuyShopItem(p, cardConfig.ShopId, uint32(lotteryTimes-num)); err != msg.ErrCode_SUCC {
				return err, nil
			}
		}
	}

	var notifyItems []uint32
	var changeItems []uint32
	var retItem []*template.SimpleItem
	var err msg.ErrCode

	if lotteryTimes > 1 {
		err, retItem, changeItems = lotteryTen(p, cardInfo, cardConfig, lotteryTimes)
	} else {
		err, retItem, changeItems = lotteryOne(p, cardInfo, cardConfig)
	}

	//str := ""
	//for _, item := range retItem {
	//	str += fmt.Sprintf("%d-%d,", item.ItemId, item.ItemNum)
	//}
	//fmt.Println(str)

	// tda
	if cardConfig.CardType == enum.Lottery_Type_Desert {
		tdaItemSlice := make([]*tda.Item, 0, len(retItem))
		for _, item := range retItem {
			tdaItemSlice = append(tdaItemSlice, &tda.Item{ItemId: strconv.Itoa(int(item.ItemId)), ItemNum: item.ItemNum})
		}
		tda.TdaEventDessertLottery(p.ChannelId, p.TdaCommonAttr, lotteryTimes, tdaItemSlice)
	}

	if err == msg.ErrCode_SUCC {
		notifyItems = append(notifyItems, changeItems...)

		updateClientItemsChange(p.GetUserId(), notifyItems)

		// 通知奖池变化
		notifyMsg := &msg.NotifyCardPoolInfoChange{}
		notifyMsg.Data = append(notifyMsg.Data, ToProtocolCardPool(cardInfo))
		p.SendNotify(notifyMsg)

		//para := fmt.Sprintf("type:%v,id:%v|", lotteryTimes, cardInfo.CardId)
		for i := 0; i < len(retItem); i++ {
			itemId := retItem[i].ItemId
			//para += fmt.Sprintf("%v,%v|", itemId, retItem[i].ItemNum)

			if template.GetTickerTemplate().CheckLotteryTicker(int(itemId)) {

				bannerMsg := &msg.NotifyBanner{}
				bannerMsg.BtType = msg.BannerType_Banner_Game
				bannerMsg.Content = fmt.Sprintf("%d", template.LotteryTickerID)
				bannerMsg.Params = append(bannerMsg.Params, p.UserData.Nick)
				bannerMsg.Params = append(bannerMsg.Params, fmt.Sprintf("%s_%d", template.LOTTERY_SHIP_FILE, cardId))
				bannerMsg.Params = append(bannerMsg.Params, fmt.Sprintf("%s_%d", template.ITEM_FILE, itemId))
				BoadCastMsg(bannerMsg)
			}
		}
		//ServMgr.GetCommonService().AddStaticsData(p, publicconst.Statics_Lottey, para)

		//report.ReportLottery(p.ChannelId, p.GetUserId(), config.Conf.ServerId, cardInfo.CardId, lotteryTimes)
		switch cardConfig.CardType {
		case enum.Lottery_Type_Normal:
			UpdateTask(p, true, publicconst.TASK_COND_SHIP_LOTTERY, lotteryTimes) // 累计进行XX次库鲁招募
			processHistoryData(p, publicconst.TASK_COND_SHIP_LOTTERY, 0, lotteryTimes)
		case enum.Lottery_Type_Disk:
			UpdateTask(p, true, publicconst.TASK_COND_DISK_BIND_BOX_LOTTERY, lotteryTimes) // 累计进行XX次库鲁招募
			processHistoryData(p, publicconst.TASK_COND_DISK_BIND_BOX_LOTTERY, 0, lotteryTimes)
		case enum.Lottery_Type_Pet:
			UpdateTask(p, true, publicconst.TASK_COND_SHIP_BEAST_EGG_LOTTERY, lotteryTimes) // 累计进行XX次库鲁招募
			processHistoryData(p, publicconst.TASK_COND_SHIP_BEAST_EGG_LOTTERY, 0, lotteryTimes)
		default:
		}

	}
	return err, retItem
}

// GetLotteryTimesReward 获取抽奖次数奖励
func GetLotteryTimesReward(p *player.Player, cardId uint32) (msg.ErrCode, []*model.SimpleItem) {
	cardInfo := getCardPool(p, cardId)
	if cardInfo == nil || cardInfo.Status == 0 {
		return msg.ErrCode_CARD_POOL_NOT_EXIST, nil
	}

	curTime := tools.GetCurTime()
	if cardInfo.StartTime > 0 && cardInfo.EndTime > 0 && curTime > cardInfo.EndTime {
		return msg.ErrCode_CARD_POOL_NOT_EXIST, nil
	}

	cardConfig := template.GetLotteryShipTemplate().GetCardPool(cardId)
	if cardConfig == nil {
		return msg.ErrCode_CARD_POOL_NOT_EXIST, nil
	}

	if cardConfig.ProgressTimes == 0 {
		return msg.ErrCode_INVALID_DATA, nil
	}

	times := cardInfo.LotteryTotalTimes / cardConfig.ProgressTimes
	if times == 0 {
		return msg.ErrCode_INVALID_DATA, nil
	}

	cardInfo.LotteryTotalTimes -= times * cardConfig.ProgressTimes
	p.SaveCardPool()

	var notifyItems []uint32
	var retItem []*model.SimpleItem

	for i := 0; i < len(cardConfig.LotteryTimesReward); i++ {
		addItems := AddItem(p.GetUserId(),
			cardConfig.LotteryTimesReward[i].ItemId,
			int32(cardConfig.LotteryTimesReward[i].ItemNum*times),
			publicconst.LotteryTimesAddItem, false)
		retItem = append(retItem, addItems...)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	// 通知奖池变化
	notifyMsg := &msg.NotifyCardPoolInfoChange{}
	notifyMsg.Data = append(notifyMsg.Data, ToProtocolCardPool(cardInfo))
	p.SendNotify(notifyMsg)
	return msg.ErrCode_SUCC, retItem
}

// ResetCardPool 重置数据
func ResetCardPool(p *player.Player, notifyClient bool) {
	if p.UserData.CardPool == nil {
		return
	}

	curTime := tools.GetCurTime()
	var resetData []*model.CardPool
	for i := 0; i < len(p.UserData.CardPool.CardPools); i++ {
		update := false
		if p.UserData.CardPool.CardPools[i].NextFreeTime > 0 && curTime >= p.UserData.CardPool.CardPools[i].NextFreeTime {
			p.UserData.CardPool.CardPools[i].FreeTimes = 1
			cardConfig := template.GetLotteryShipTemplate().GetCardPool(p.UserData.CardPool.CardPools[i].CardId)
			p.UserData.CardPool.CardPools[i].NextFreeTime = tools.GetDailyXRefreshTime(cardConfig.XDayFreeTime, template.GetSystemItemTemplate().RefreshHour)

			resetData = append(resetData, p.UserData.CardPool.CardPools[i])
			update = true
		}
		if p.UserData.CardPool.CardPools[i].NextResetTime > 0 && curTime >= p.UserData.CardPool.CardPools[i].NextResetTime {
			p.UserData.CardPool.CardPools[i].TotalTimes = 0
			p.UserData.CardPool.CardPools[i].NextResetTime = tools.GetDailyRefreshTime()

			resetData = append(resetData, p.UserData.CardPool.CardPools[i])
			update = true
		}
		if update {
			p.SaveCardPool()
		}
	}

	if notifyClient && len(resetData) > 0 {
		// 通知奖池变化
		notifyMsg := &msg.NotifyCardPoolInfoChange{}
		notifyMsg.Data = append(notifyMsg.Data, ToProtocolCardPools(resetData)...)
		p.SendNotify(notifyMsg)
	}
}

func GmResetDataCardPool(p *player.Player) {
	if p.UserData.CardPool == nil {
		return
	}

	var resetData []*model.CardPool
	for i := 0; i < len(p.UserData.CardPool.CardPools); i++ {
		update := false
		if p.UserData.CardPool.CardPools[i].NextFreeTime > 0 {
			p.UserData.CardPool.CardPools[i].FreeTimes = 1

			resetData = append(resetData, p.UserData.CardPool.CardPools[i])
			update = true
		}

		if update {
			p.SaveCardPool()
		}
	}

	if len(resetData) > 0 {
		// 通知奖池变化
		notifyMsg := &msg.NotifyCardPoolInfoChange{}
		notifyMsg.Data = append(notifyMsg.Data, ToProtocolCardPools(resetData)...)
		p.SendNotify(notifyMsg)
	}
}

// RefreshCardPoolActivity 刷新卡池活动
func RefreshCardPoolActivity(p *player.Player, cardId, start, end, lotteryEnd uint32) *model.CardPool {
	if p.UserData.CardPool == nil {
		log.Error("CardPool is nil", zap.Uint64("accountId", p.GetUserId()))
		return nil
	}

	if cardPool := getCardPool(p, cardId); cardPool != nil {
		log.Error("CardPool is exist", zap.Uint64("accountId", p.GetUserId()))
		return nil
	}

	cardConfig := template.GetLotteryShipTemplate().GetCardPool(cardId)
	var freeTimes uint32 = 0
	var nextResetTime uint32 = 0
	if cardConfig.XDayFreeTime > 0 {
		freeTimes = 1
	}

	if cardConfig.CardType != 2 {
		nextResetTime = tools.GetDailyRefreshTime()
	}

	newCardPool := model.NewCardPool(cardConfig.Id, freeTimes,
		tools.GetDailyXRefreshTime(cardConfig.XDayFreeTime, template.GetSystemItemTemplate().RefreshHour),
		nextResetTime, start, end, lotteryEnd)
	p.UserData.CardPool.CardPools = append(p.UserData.CardPool.CardPools, newCardPool)
	p.SaveCardPool()
	return newCardPool
}

// getCardPool 获得卡池
func getCardPool(p *player.Player, cardId uint32) *model.CardPool {
	for i := 0; i < len(p.UserData.CardPool.CardPools); i++ {
		if p.UserData.CardPool.CardPools[i].CardId == cardId {
			return p.UserData.CardPool.CardPools[i]
		}
	}
	return nil
}

func ToProtocolCardPools(data []*model.CardPool) []*msg.CardPoolInfo {
	var ret []*msg.CardPoolInfo
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolCardPool(data[i]))
	}
	return ret
}

func ToProtocolCardPool(data *model.CardPool) *msg.CardPoolInfo {
	return &msg.CardPoolInfo{
		CardId:              data.CardId,
		FreeLotteryTimes:    data.FreeTimes,
		NextFreeLotteryTime: data.NextFreeTime,
		UseLotteryTimes:     data.TotalTimes,
		StartTime:           data.StartTime,
		EndTime:             data.EndTime,
		TenQuaranteeTimes:   data.TenQuaranteeTimes,
		BigQuaranteeTimes:   data.BigQuaranteeTimes,
		Status:              data.Status,
		FirstLottery:        data.FirstLottery,
		LotteryEndTime:      data.LotteryEndTime,
		LotteryTotalTimes:   data.LotteryTotalTimes,
	}
}
