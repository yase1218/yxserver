package service

import (
	"kernel/tools"
	"msg"
	"sync"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"

	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
)

var (
	VerLock     sync.RWMutex
	ShopItemVer map[uint32]uint32
)

const (
	Obtained  = iota // 可领取
	Receive          // 已领取
	Forbidden        // 不可领取
)

func InitShop() {
	ShopItemVer = make(map[uint32]uint32)
	// TODO 去中心服获取
	// for i := 0; i < len(data); i++ {
	// 	ShopItemVer[data[i].ShopItemId] = data[i].Ver
	// }
}

func UpdateShopItemVer(shopItemId, ver uint32) {
	ShopItemVer[shopItemId] = ver
}

func DelShopItemVer(shopItemId uint32) {
	delete(ShopItemVer, shopItemId)
}

// GetShopItem 获取商城物品
//func GetShopItem(p *player.Player, tp uint32) (msg.ErrCode, []*model.ShopItem) {
//	if p.UserData.Shop == nil {
//		log.Error("shop nil", zap.Uint64("accountid", p.GetUserId()))
//		return msg.ErrCode_SYSTEM_ERROR, nil
//	}
//
//	shopItems := template.GetShopTemplate().GetShopItemByType(tp)
//	var ret []*model.ShopItem
//	var addItem []*model.ShopItem
//	resetDailyWeeklyItem(p, shopItems) // 刷新玩家商品
//	for i := 0; i < len(shopItems); i++ {
//		if item := getShopItem(p, shopItems[i].Id); item != nil {
//			if !item.Lock {
//				if refreshShopItem(p, item, shopItems[i], false) {
//					p.SaveShop()
//				}
//			} else {
//				lock := shopItemIsLock(p, shopItems[i])
//				if !lock && item.Lock { // 解锁
//					unlockShopItem(p, item, shopItems[i])
//				}
//			}
//			ret = append(ret, item)
//		} else {
//			if canAddShopItem(p, shopItems[i]) {
//				lock := shopItemIsLock(p, shopItems[i])
//				item := model.NewShopItem(shopItems[i].Id, lock, shopItems[i].ConfigVer)
//				refreshShopItem(p, item, shopItems[i], true)
//				addItem = append(addItem, item)
//				ret = append(ret, item)
//			}
//		}
//	}
//	if len(addItem) > 0 {
//		addShopItems(p, addItem)
//	}
//	return msg.ErrCode_SUCC, ret
//}

// GetShopItems 获取商城物品 重构
func GetShopItems(p *player.Player, tp uint32) (msg.ErrCode, []*model.ShopItem) {
	//if p.UserData.Shop == nil {
	//	log.Error("GetShopItems shop nil", zap.Uint64("accountid", p.GetUserId()))
	//	return msg.ErrCode_SYSTEM_ERROR, nil
	//}

	shopItems := template.GetShopTemplate().GetShopItemByType(tp)
	var ret []*model.ShopItem
	curTime := tools.GetCurTime()
	var refreshLockItems []*model.ShopItem
	for i := 0; i < len(shopItems); i++ {
		if shopItems[i].PreShopItem > 0 { // 有前置商品
			preItem := template.GetShopTemplate().GetShopItem(shopItems[i].PreShopItem) // 前置商品
			idMap := make(map[uint32]struct{})                                          // 前置商品id map，方便查找
			for j := len(shopItems[i].PreShopLink) - 1; j >= 0; j-- {
				pre := template.GetShopTemplate().GetShopItem(shopItems[i].PreShopLink[j])
				var nextRefreshTime uint32
				if item := getShopItem(p, pre.Id); item != nil {
					if item.NextRefreshTime > 0 {
						nextRefreshTime = item.NextRefreshTime
					}
				} else {
					nextRefreshTime = GetNextRefreshTime(pre)
				}
				if nextRefreshTime > 0 && curTime > nextRefreshTime && pre.LimitNum > 0 {
					idMap[pre.Id] = struct{}{}    // 加入过滤
					for k := j - 1; k >= 0; k-- { // 剩余商品加入过滤
						idMap[shopItems[i].PreShopLink[k]] = struct{}{} // 加入过滤
					}
					break
				}
			}

			if _, ok := idMap[preItem.Id]; !ok { // 前置商品没有被过滤
				if pr := getShopItem(p, preItem.Id); pr != nil && pr.BuyTimes >= preItem.LimitNum { // 前置商品符合条件
					item := getShopItem(p, shopItems[i].Id)
					if item == nil {
						item = genShopItem(p, shopItems[i], curTime)
						p.UserData.Shop.Items = append(p.UserData.Shop.Items, item)
					} else {
						if item.Lock { //最后刷新 防止前置商品未加载是否锁计算错误
							refreshLockItems = append(refreshLockItems, item)
							//item.Lock = shopItemIsLock(p, shopItems[i])
						}
						refreshShopItem(p, item, shopItems[i], false) // 更新刷新时间
					}
					ret = append(ret, item)
				}
			}
		} else { // 没有前置商品
			if item := getShopItem(p, shopItems[i].Id); item != nil { // 玩家商品列表里有
				if item.Lock {
					item.Lock = shopItemIsLock(p, shopItems[i])
				}
				refreshShopItem(p, item, shopItems[i], false)
				ret = append(ret, item)
			} else {
				item = genShopItem(p, shopItems[i], curTime)
				if item.NextRefreshTime > 0 { // 有刷新机制商品
					p.UserData.Shop.Items = append(p.UserData.Shop.Items, item)
				}
				ret = append(ret, item)
			}
		}
	}
	for i := 0; i < len(refreshLockItems); i++ {
		cfg := template.GetShopTemplate().GetShopItem(refreshLockItems[i].Id)
		refreshLockItems[i].Lock = shopItemIsLock(p, cfg)
	}
	p.SaveShop()
	return msg.ErrCode_SUCC, ret
}

// 过滤掉需要刷新的商品（有前置商品）
func resetDailyWeeklyItem(p *player.Player, shopItems []*template.JShopItem) {
	var ids []uint32 // 玩家有这个商品且需要刷新的商品id
	curTime := tools.GetCurTime()
	for i := 0; i < len(shopItems); i++ {
		if item := getShopItem(p, shopItems[i].Id); item != nil {
			if item.NextRefreshTime > 0 && curTime >= item.NextRefreshTime { // 需要刷新
				arg0 := shopItems[i].RefreshArgs[0]
				if (arg0 == 1 || arg0 == 2 || arg0 == 3) && shopItems[i].PreShopItem > 0 { // 1日，2周，3月，4永久，5固定时间
					ids = append(ids, shopItems[i].Id)
				}
			}
		}
	}
	if len(ids) > 0 {
		var temp []*model.ShopItem
		for i := 0; i < len(p.UserData.Shop.Items); i++ {
			id := p.UserData.Shop.Items[i].Id
			if !tools.ListContain(ids, id) { // 过滤掉需要刷新的商品
				temp = append(temp, p.UserData.Shop.Items[i])
			}
		}
		p.UserData.Shop.Items = temp
		p.SaveShop()
	}
}

// GetShopItemByIds 获得指定商品
func GetShopItemByIds(p *player.Player, ids []uint32) (msg.ErrCode, []*model.ShopItem) {
	//tps := make(map[uint32]struct{})
	//var items []*model.ShopItem
	//for i := 0; i < len(ids); i++ {
	//	if shopItems := template.GetShopTemplate().GetShopItem(ids[i]); shopItems != nil {
	//		tps[shopItems.ShopType] = struct{}{}
	//	}
	//}
	//for tp, _ := range tps {
	//	if err, data := GetShopItem(p, tp); err == msg.ErrCode_SUCC {
	//		items = append(items, data...)
	//	}
	//}
	//
	//var ret []*model.ShopItem
	//for i := 0; i < len(items); i++ {
	//	if tools.ListContain(ids, items[i].Id) {
	//		ret = append(ret, items[i])
	//	}
	//}

	var ret []*model.ShopItem
	var dbSave bool
	for i := 0; i < len(ids); i++ {
		cfg := template.GetShopTemplate().GetShopItem(ids[i])
		if cfg == nil {
			log.Error("GetShopItemByIds id err", zap.Uint64("uid", p.GetUserId()), zap.Uint32("id", ids[i]))
			continue
		}
		if item := getShopItem(p, ids[i]); item != nil {
			if !item.Lock {
				if refreshShopItem(p, item, cfg, false) {
					dbSave = true
				}
			} else {
				if lock := shopItemIsLock(p, cfg); !lock { // 解锁
					unlockShopItem(p, item, cfg)
					dbSave = true
				}
			}
			ret = append(ret, item)
		} else {
			item = genShopItem(p, cfg, tools.GetCurTime())
			ret = append(ret, item)
		}
	}
	if dbSave {
		p.SaveShop()
	}
	return msg.ErrCode_SUCC, ret
}

// BuyShopItem 购买商品
func BuyShopItem(p *player.Player, id, num uint32) (msg.ErrCode, []*model.SimpleItem, []*model.ShopItem, uint32) {
	if p.UserData.Shop == nil {
		log.Error("shop nil", zap.Uint64("accountid", p.GetUserId()))
		return msg.ErrCode_SYSTEM_ERROR, nil, nil, 0
	}

	if num == 0 {
		return msg.ErrCode_INVALID_DATA, nil, nil, 0
	}

	if num >= 5000 {
		return msg.ErrCode_INVALID_DATA, nil, nil, 0
	}

	itemConfig := template.GetShopTemplate().GetShopItem(id)
	if itemConfig == nil {
		return msg.ErrCode_SHOP_ITEM_NOT_EXIST, nil, nil, 0
	}

	if itemConfig.ChargeID > 0 || itemConfig.CostAd > 0 {
		return msg.ErrCode_BUY_SHOP_NEED_PAY, nil, nil, 0
	}

	item := getShopItem(p, id)
	if item == nil {
		return msg.ErrCode_SHOP_ITEM_NOT_EXIST, nil, nil, 0
	}

	if item.Lock {
		return msg.ErrCode_SHOP_ITEM_LOCK, nil, nil, 0
	}

	// 到时间则刷新
	if refreshShopItem(p, item, itemConfig, false) {
		p.SaveShop()
	}

	if itemConfig.LimitNum > 0 && item.BuyTimes >= itemConfig.LimitNum {
		return msg.ErrCode_BUY_SHOP_ITEM_OVER_LIMIT, nil, nil, 0
	}

	costItem := itemConfig.GetCostItem(item.ItemId, item.ItemNum)
	for i := 0; i < len(costItem); i++ {
		needNum := costItem[i].ItemNum * num //* itemConfig.Discount * num / 100
		if !EnoughItem(p.GetUserId(), itemConfig.CostItem[i].ItemId, needNum) {
			return msg.ErrCode_NO_ENOUGH_ITEM, nil, nil, 0
		}
	}

	isCostDiamond := false
	var notifyItems []uint32
	for i := 0; i < len(costItem); i++ {
		needNum := costItem[i].ItemNum * num //* itemConfig.Discount * num / 100
		CostItem(p.GetUserId(),
			costItem[i].ItemId,
			needNum,
			publicconst.ShopCostItem,
			false)
		notifyItems = append(notifyItems, itemConfig.CostItem[i].ItemId)
		if itemConfig.CostItem[i].ItemId == uint32(publicconst.ITEM_CODE_DIAMOND) {
			isCostDiamond = true
		}
	}
	if isCostDiamond {
		UpdateTask(p, true, publicconst.TASK_COND_RESUME_DIAMOND_BUY_SHOP, num)
		processHistoryData(p, publicconst.TASK_COND_RESUME_DIAMOND_BUY_SHOP, 0, num)
	}

	var addItems []*model.SimpleItem
	var rewardItems []*model.SimpleItem
	if itemConfig.IsFix == 1 {
		for i := 0; i < len(itemConfig.GetItem); i++ {
			temp := AddItem(p.GetUserId(),
				itemConfig.GetItem[i].ItemId,
				int32(itemConfig.GetItem[i].ItemNum*num),
				publicconst.ShopAddItem,
				false)
			addItems = append(addItems, temp...)

			rewardItems = append(rewardItems, &model.SimpleItem{
				Id:  itemConfig.GetItem[i].ItemId,
				Num: itemConfig.GetItem[i].ItemNum * num,
			})
		}
	} else {
		addItems = AddItem(p.GetUserId(),
			item.ItemId,
			int32(item.ItemNum*num),
			publicconst.ShopAddItem,
			false)

		rewardItems = append(rewardItems, &model.SimpleItem{
			Id:  item.ItemId,
			Num: item.ItemNum * num,
		})
	}

	// TODO 公会红包
	// if itemConfig.RedPacket > 0 {
	// 	ServMgr.GetAllianceService().AddRedPacket(p, itemConfig.RedPacket)
	// }

	notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)

	// 通知客户端
	updateClientItemsChange(p.GetUserId(), notifyItems)

	item.BuyTimes += num
	item.UpdateTime = tools.GetCurTime()
	p.SaveShop()

	// 购买满了如果有下一级就加下一级商品
	var addItem []*model.ShopItem
	if itemConfig.LimitNum > 0 && item.BuyTimes >= itemConfig.LimitNum {
		for i := 0; i < len(itemConfig.NextShopItem); i++ {
			if nextConfig := template.GetShopTemplate().GetShopItem(itemConfig.NextShopItem[i]); nextConfig != nil {
				item := model.NewShopItem(itemConfig.NextShopItem[i], false, nextConfig.ConfigVer)
				refreshShopItem(p, item, nextConfig, true)
				addItem = append(addItem, item)
			}
		}
		if len(addItem) > 0 {
			addShopItems(p, addItem)
		}
	}

	// // 统计记录
	// para := fmt.Sprintf("id:%v,num:%v|", id, num)
	// for i := 0; i < len(addItems); i++ {
	// 	para += fmt.Sprintf("%v,%v|", addItems[i].Id, addItems[i].Num)
	// }
	// ServMgr.GetCommonService().AddStaticsData(p, publicconst.Statics_Buy_Shop_Item_Id, para)

	// // tda
	// costItemId, costItemNum := uint32(0), uint32(0)
	// if len(costItem) > 0 {
	// 	costItemId = costItem[0].ItemId
	// 	costItemNum = costItem[0].ItemNum
	// }
	// tda.TdaShopBuy(p.ChannelId, p.TdaCommonAttr, itemConfig.ShopType, itemConfig.Id, num, costItemId, costItemNum)

	return msg.ErrCode_SUCC, rewardItems, addItem, item.BuyTimes
}

// BuyShopItemWithCost 购买商品
func BuyShopItemWithCost(p *player.Player, id, num, costId uint32) (msg.ErrCode, []*model.SimpleItem, []*model.ShopItem, uint32) {
	if p.UserData.Shop == nil {
		log.Error("shop nil", zap.Uint64("accountid", p.GetUserId()))
		return msg.ErrCode_SYSTEM_ERROR, nil, nil, 0
	}

	if num == 0 {
		return msg.ErrCode_INVALID_DATA, nil, nil, 0
	}

	if num >= 5000 {
		return msg.ErrCode_INVALID_DATA, nil, nil, 0
	}

	itemConfig := template.GetShopTemplate().GetShopItem(id)
	if itemConfig == nil {
		return msg.ErrCode_SHOP_ITEM_NOT_EXIST, nil, nil, 0
	}

	if itemConfig.ChargeID > 0 || itemConfig.CostAd > 0 {
		return msg.ErrCode_BUY_SHOP_NEED_PAY, nil, nil, 0
	}

	item := getShopItem(p, id)
	if item == nil {
		item = genShopItem(p, itemConfig, tools.GetCurTime())
		p.UserData.Shop.Items = append(p.UserData.Shop.Items, item)
	} else {
		item.UpdateTime = tools.GetCurTime()
	}
	if item.Lock {
		return msg.ErrCode_SHOP_ITEM_LOCK, nil, nil, 0
	}

	if itemConfig.LimitNum > 0 && (num > itemConfig.LimitNum || item.BuyTimes+num > itemConfig.LimitNum) {
		return msg.ErrCode_BUY_SHOP_ITEM_OVER_LIMIT, nil, nil, 0
	}

	if len(itemConfig.CostItem) == 1 && costId == 0 {
		costId = itemConfig.CostItem[0].ItemId
	}

	// 消耗数量检查
	var costNum int32 = 0
	for i := 0; i < len(itemConfig.CostItem); i++ {
		if itemConfig.CostItem[i].ItemId == costId {
			costNum = int32(itemConfig.CostItem[i].ItemNum) * int32(num)
			break
		}
	}
	//if costNum == -1 { // 参数有误
	//	log.Error("BuyShopItemWithCost costItemId err", zap.Uint32("itemID", id), zap.Uint32("costItemId", costId))
	//	return msg.ErrCode_NO_ENOUGH_ITEM, nil, nil, 0
	//}

	var notifyItems []uint32
	if costNum > 0 {
		if !EnoughItem(p.GetUserId(), costId, uint32(costNum)) {
			return msg.ErrCode_NO_ENOUGH_ITEM, nil, nil, 0
		}
		CostItem(p.GetUserId(), costId, uint32(costNum), publicconst.ShopCostItem, false)
		notifyItems = append(notifyItems, costId)
	}

	isCostDiamond := false
	if costId == uint32(publicconst.ITEM_CODE_DIAMOND) {
		isCostDiamond = true
	}
	if isCostDiamond {
		UpdateTask(p, true, publicconst.TASK_COND_RESUME_DIAMOND_BUY_SHOP, num)
		processHistoryData(p, publicconst.TASK_COND_RESUME_DIAMOND_BUY_SHOP, 0, num)
	}

	var rewardItems []*model.SimpleItem
	addItems := AddItem(p.GetUserId(), item.ItemId, int32(num*item.ItemNum), publicconst.ShopAddItem, false)
	rewardItems = append(rewardItems, &model.SimpleItem{Id: item.ItemId, Num: num * item.ItemNum})
	notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	// 通知客户端
	updateClientItemsChange(p.GetUserId(), notifyItems)

	item.BuyTimes += num
	item.UpdateTime = tools.GetCurTime()
	p.SaveShop()

	// 购买满了如果有下一级就加下一级商品
	var addItem []*model.ShopItem
	if itemConfig.LimitNum > 0 && item.BuyTimes >= itemConfig.LimitNum {
		for i := 0; i < len(itemConfig.NextShopItem); i++ {
			if nextConfig := template.GetShopTemplate().GetShopItem(itemConfig.NextShopItem[i]); nextConfig != nil {
				item := model.NewShopItem(itemConfig.NextShopItem[i], false, nextConfig.ConfigVer)
				item.NextRefreshTime = GetNextRefreshTime(itemConfig)
				addItem = append(addItem, item)
			}
		}
	}
	return msg.ErrCode_SUCC, rewardItems, addItem, item.BuyTimes
}

// BuyShopItemCallBack 购买商品回调
func BuyShopItemCallBack(p *player.Player, itemConfig *template.JShopItem) {
	if itemConfig == nil {
		return
	}
	item := getShopItem(p, itemConfig.Id)
	if item == nil {
		return
	}

	if item.Lock {
		log.Error("BuyShopItemCallBack lock", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("id", itemConfig.Id))
		return
	}

	if itemConfig.LimitNum > 0 && item.BuyTimes >= itemConfig.LimitNum {
		log.Error("BuyShopItemCallBack lock LimitNum", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("id", itemConfig.Id))
		return
	}

	var notifyItems []uint32
	var addItems []*model.SimpleItem
	var rewardItems []*model.SimpleItem
	if itemConfig.IsFix == 1 {
		for i := 0; i < len(itemConfig.GetItem); i++ {
			temp := AddItem(p.GetUserId(),
				itemConfig.GetItem[i].ItemId,
				int32(itemConfig.GetItem[i].ItemNum*item.Times),
				publicconst.ShopAddItem,
				false)
			addItems = append(addItems, temp...)

			rewardItems = append(rewardItems, &model.SimpleItem{
				Id:  itemConfig.GetItem[i].ItemId,
				Num: itemConfig.GetItem[i].ItemNum * item.Times,
			})
		}
	} else {
		addItems = AddItem(p.GetUserId(),
			item.ItemId,
			int32(item.ItemNum*item.Times),
			publicconst.ShopAddItem,
			false)

		rewardItems = append(rewardItems, &model.SimpleItem{
			Id:  item.ItemId,
			Num: item.ItemNum * item.Times,
		})
	}

	notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	//event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(p, notifyItems, ListenNotifyClientItemEventEvent))

	updateClientItemsChange(p.GetUserId(), notifyItems)
	item.BuyTimes += 1
	if item.Times == 2 {
		item.Times = 1
	}
	item.UpdateTime = tools.GetCurTime()
	p.SaveShop()

	// 购买满了如果有下一级就加下一级商品
	var addItem []*model.ShopItem
	if itemConfig.LimitNum > 0 && item.BuyTimes >= itemConfig.LimitNum {
		for i := 0; i < len(itemConfig.NextShopItem); i++ {
			if nextConfig := template.GetShopTemplate().GetShopItem(itemConfig.NextShopItem[i]); nextConfig != nil {
				item := model.NewShopItem(itemConfig.NextShopItem[i], false, nextConfig.ConfigVer)
				refreshShopItem(p, item, nextConfig, true)
				addItem = append(addItem, item)
			}
		}
		if len(addItem) > 0 {
			addShopItems(p, addItem)
		}
	}
	p.SendNotify(&msg.NotifyShopBuyItem{
		ShopItemId:   itemConfig.Id,
		BuyNum:       1,
		BuyTimes:     item.BuyTimes,
		Times:        item.Times,
		GetItems:     ToProtocolSimpleItems(rewardItems),
		NextShopItem: ToProtocolShopItems(addItem),
	})
}

// GetShopRedPoint 获取商城红点
func GetShopRedPoint(p *player.Player) *msg.RedPointInfo {
	ret := &msg.RedPointInfo{RdType: msg.RedPointType_Shop_Red_Point}
	if p.UserData.Shop == nil {
		log.Error("shop nil", zap.Uint64("accountid", p.GetUserId()))
		return ret
	}
	templateItem := template.GetShopTemplate().GetFreeShopItem()
	for i := 0; i < len(templateItem); i++ {
		item := getShopItem(p, templateItem[i].Id)
		if item != nil && !item.Lock {
			// 需要刷新重置的则刷新重置
			if refreshShopItem(p, item, templateItem[i], false) {
				p.SaveShop()
			}
			if item.BuyTimes == 0 {
				ret.RdData = append(ret.RdData, templateItem[i].Id)
			}
		}
	}
	return ret
}

// RefreshShop 刷新商品
func RefreshShop(p *player.Player, shopItemIds []uint32) (msg.ErrCode, []*model.ShopItem) {
	var ret []*model.ShopItem
	var dbSave bool
	for i := 0; i < len(shopItemIds); i++ {
		if item := getShopItem(p, shopItemIds[i]); item != nil {
			if item.Lock {
				continue
			}

			shopItemCfg := template.GetShopTemplate().GetShopItem(shopItemIds[i])
			if refreshShopItem(p, item, shopItemCfg, false) {
				ret = append(ret, item)
				dbSave = true
			}
		}
	}
	if dbSave {
		p.SaveShop()
	}

	return msg.ErrCode_SUCC, ret
}

func canAddShopItem(p *player.Player, itemConfig *template.JShopItem) bool {
	if itemConfig.PreShopItem > 0 {
		item := getShopItem(p, itemConfig.PreShopItem)
		preConfig := template.GetShopTemplate().GetShopItem(itemConfig.PreShopItem)
		if item == nil || item.BuyTimes < preConfig.LimitNum {
			return false
		}
	}
	return true
}

// shopItemIsLock 商品是否解锁
func shopItemIsLock(p *player.Player, itemConfig *template.JShopItem) bool {
	if itemConfig.UnlockLv > 0 {
		if p.UserData.Level < itemConfig.UnlockLv {
			return true
		}
	}
	if itemConfig.UnlockMission > 0 {
		if !IsPassMission(p, itemConfig.UnlockMission, true) {
			return true
		}
	}

	if itemConfig.PreShopItem > 0 {
		item := getShopItem(p, itemConfig.PreShopItem)
		preConfig := template.GetShopTemplate().GetShopItem(itemConfig.PreShopItem)
		if item == nil || item.BuyTimes < preConfig.LimitNum {
			return true
		}
	}

	// check guild lv TODO 公会联盟
	// if itemConfig.UnlockGuildLv > 0 {
	// 	member, err := dao.GetMember(p.AccountInfo.AccountId)
	// 	if err != nil || member == nil {
	// 		return false
	// 	}

	// 	alliance, err := dao.GetAlliance(member.AllianceID)
	// 	if err != nil || alliance == nil {
	// 		return false
	// 	}

	// 	if int(alliance.Level) > itemConfig.UnlockGuildLv {
	// 		return true
	// 	}
	// }
	return false
}

func GmDayRefreshShopItems(p *player.Player) {
	if p.UserData.Shop == nil {
		log.Error("shop nil", zap.Uint64("accountid", p.GetUserId()))
		return
	}

	for _, item := range p.UserData.Shop.Items {
		shopItemCfg := template.GetShopTemplate().GetShopItem(item.Id)
		if shopItemCfg == nil {
			continue
		}

		if shopItemCfg.ShopType != 1 {
			continue
		}

		item.NextRefreshTime = tools.GetDailyRefreshTime()
		item.BuyTimes = 0
	}
	p.SaveShop()
}

// refreshShopItem 刷新
func refreshShopItem(p *player.Player, item *model.ShopItem, itemConfig *template.JShopItem, refresh bool) bool {
	curTime := tools.GetCurTime()
	curVer := getShopItemVer(item.Id)
	if !refresh {
		if item.NextRefreshTime > 0 && curTime >= item.NextRefreshTime {
			refresh = true
		}

		if !refresh && item.Ver < curVer {
			refresh = true
		}
	}

	if !refresh {
		return false
	}

	if itemConfig.IsFix != 1 {
		getItem := itemConfig.RefreshGetItem(p.Rand)
		if getItem == nil {
			log.Error("refreshShopItem getItem is null", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("shopItemId", item.Id))
		} else {
			item.ItemId = getItem.ItemId
			item.ItemNum = getItem.ItemNum
		}
	} else {
		item.ItemId = itemConfig.GetItem[0].ItemId
		item.ItemNum = itemConfig.GetItem[0].ItemNum
	}
	item.Ver = curVer
	item.UpdateTime = curTime
	switch itemConfig.RefreshArgs[0] {
	case 1:
		item.NextRefreshTime = tools.GetDailyRefreshTime()
		item.BuyTimes = 0
	case 2:
		item.NextRefreshTime = tools.GetWeeklyRefreshTime(template.GetSystemItemTemplate().RefreshHour)
		item.BuyTimes = 0
	case 3:
		item.NextRefreshTime = tools.GetMonthRefreshTime()
		item.BuyTimes = 0
	case 4:
	case 5:
		var minTime uint32
		for i := 1; i < len(itemConfig.RefreshArgs); i++ {
			temp := tools.GetHourRefreshTime(itemConfig.RefreshArgs[i])
			if i == 1 {
				minTime = temp
			} else {
				if temp < minTime {
					minTime = temp
				}
			}
		}
		item.NextRefreshTime = minTime
		item.BuyTimes = 0
	}
	return true
}

func unlockShopItem(p *player.Player, item *model.ShopItem, itemConfig *template.JShopItem) {
	item.Lock = false
	switch itemConfig.RefreshArgs[0] {
	case 1:
		item.NextRefreshTime = tools.GetDailyRefreshTime()
		item.BuyTimes = 0
	case 2:
		item.NextRefreshTime = tools.GetWeeklyRefreshTime(template.GetSystemItemTemplate().RefreshHour)
		item.BuyTimes = 0
	case 3:
		item.NextRefreshTime = tools.GetMonthRefreshTime()
		item.BuyTimes = 0
	case 4:
	case 5:
		var minTime uint32
		for i := 1; i < len(itemConfig.RefreshArgs); i++ {
			temp := tools.GetHourRefreshTime(itemConfig.RefreshArgs[i])
			if i == 1 {
				minTime = temp
			} else {
				if temp < minTime {
					minTime = temp
				}
			}
		}
		item.NextRefreshTime = minTime
		item.BuyTimes = 0
	}
	//p.SaveShop()
}

func getShopItemVer(shopItemId uint32) uint32 {
	if data, ok := ShopItemVer[shopItemId]; ok {
		return data
	}
	return 0
}

func ToProtocolShopItem(item *model.ShopItem) *msg.ShopItem {
	ret := &msg.ShopItem{
		ShopItemId:  item.Id,
		BuyTimes:    item.BuyTimes,
		RefreshTime: item.NextRefreshTime,
		Times:       item.Times,
	}
	if jShop := template.GetShopTemplate().GetShopItem(item.Id); jShop != nil {
		if jShop.IsFix == 1 {
			for i := 0; i < len(jShop.GetItem); i++ {
				ret.GetItems = append(ret.GetItems, &msg.SimpleItem{
					ItemId:  jShop.GetItem[i].ItemId,
					ItemNum: jShop.GetItem[i].ItemNum,
				})
			}
		} else {
			ret.GetItems = append(ret.GetItems, &msg.SimpleItem{
				ItemId:  item.ItemId,
				ItemNum: item.ItemNum,
			})
		}

	}

	if item.Lock {
		ret.Lock = 1
	}
	return ret
}

func ToProtocolShopItems(items []*model.ShopItem) []*msg.ShopItem {
	var ret []*msg.ShopItem
	for i := 0; i < len(items); i++ {
		ret = append(ret, ToProtocolShopItem(items[i]))
	}
	return ret
}

func RefreshShopLogin(p *player.Player) {
	//if p.UserData.Shop == nil {
	//	log.Error("shop nil", zap.Uint64("accountid", p.GetUserId()))
	//	return
	//}
	//// 删除配置低的商品
	//var shopItems []*model.ShopItem
	//var deleteIds []uint32
	//for _, v := range p.UserData.Shop.Items {
	//	if config := template.GetShopTemplate().GetShopItem(v.Id); config != nil {
	//		if config.ConfigVer > v.ConfigVer {
	//			deleteIds = append(deleteIds, v.Id)
	//		} else {
	//			shopItems = append(shopItems, v)
	//		}
	//	} else {
	//		deleteIds = append(deleteIds, v.Id)
	//	}
	//}
	//
	//p.UserData.Shop.Items = shopItems
	//if len(deleteIds) > 0 {
	//	p.SaveShop()
	//}
	//
	//// 默认加载类型1的商品
	//GetShopItem(p, 1) // TODO注释
	//GetShopItem(p, 2)
	//GetShopItem(p, 5)
	//curTime := tools.GetCurTime()
	//for i := len(p.UserData.Shop.Items) - 1; i >= 0; i-- {
	//	if config := template.GetShopTemplate().GetShopItem(p.UserData.Shop.Items[i].Id); config != nil && p.UserData.Shop.Items[i].ConfigVer == config.ConfigVer {
	//		if p.UserData.Shop.Items[i].NextRefreshTime > 0 { // 有刷新时间
	//			if p.UserData.Shop.Items[i].NextRefreshTime > curTime { // 没到刷新时间 保留
	//				if _, ok := p.UserData.ShopItemMap[p.UserData.Shop.Items[i].Id]; !ok {
	//					p.UserData.ShopItemMap[p.UserData.Shop.Items[i].Id] = p.UserData.Shop.Items[i]
	//				}
	//			}
	//		} else { // 没有刷新时间直接保留
	//			p.UserData.ShopItemMap[p.UserData.Shop.Items[i].Id] = p.UserData.Shop.Items[i]
	//		}
	//	}
	//}
}

// func loadShop(p *player.Player) msg.ErrCode {
// 	err, temp := dao.ShopDao.LoadShop(p.GetUserId())
// 	if err != msg.ErrCode_SUCC {
// 		return err
// 	}

// 	// 删除配置低的商品
// 	var shopItems []*model.ShopItem
// 	var deleteIds []uint32
// 	for i := 0; i < len(temp.Items); i++ {
// 		if config := template.GetShopTemplate().GetShopItem(temp.Items[i].Id); config != nil {
// 			if config.ConfigVer > temp.Items[i].ConfigVer {
// 				deleteIds = append(deleteIds, temp.Items[i].Id)
// 			} else {
// 				shopItems = append(shopItems, temp.Items[i])
// 			}
// 		} else {
// 			deleteIds = append(deleteIds, temp.Items[i].Id)
// 		}
// 	}
// 	if len(deleteIds) > 0 {
// 		dao.ShopDao.DelShopItems(p.GetUserId(), deleteIds)
// 	}

// 	temp.Items = shopItems
// 	p.UserData.Shop = temp

// 	// 默认加载类型1的商品
// 	GetShopItem(p, 1)
// 	GetShopItem(p, 2)
// 	GetShopItem(p, 5)
// 	return msg.ErrCode_SUCC
// }

// getShopItem 获取商品
func getShopItem(p *player.Player, id uint32) *model.ShopItem {
	for i := 0; i < len(p.UserData.Shop.Items); i++ {
		if p.UserData.Shop.Items[i].Id == id {
			return p.UserData.Shop.Items[i]
		}
	}
	return nil
}

// addShopItem 添加商品
func addShopItem(p *player.Player, item *model.ShopItem) {
	p.UserData.Shop.Items = append(p.UserData.Shop.Items, item)
	p.SaveShop()
}

// addShopItems 添加商品
func addShopItems(p *player.Player, items []*model.ShopItem) {
	p.UserData.Shop.Items = append(p.UserData.Shop.Items, items...)
	p.SaveShop()
}

// updateShopItem 更新商品
func updateShopItem(p *player.Player, item *model.ShopItem) {
	p.SaveShop()
}

func GetFirstChargePackageReward(p *player.Player, configId uint32, day uint32) (msg.ErrCode, []*msg.SimpleItem) {
	var itemsList []*msg.SimpleItem
	var dayIdx = day - 1
	if day <= 0 || configId <= 0 {
		return msg.ErrCode_CONDITION_NOT_MET, itemsList
	}

	cfg := template.GetOneChargeTemplate().GetTargetCfg(configId)
	if cfg == nil {
		log.Error("one charge config is not exists", zap.Uint32("cfgId", configId), zap.Uint32("day", day))
		return msg.ErrCode_CONDITION_NOT_MET, itemsList
	}

	itemsList = cfg.ParsedItemRewards[dayIdx]
	if len(itemsList) <= 0 || itemsList == nil {
		log.Error("one charge config is not exists", zap.Uint32("cfgId", configId), zap.Uint32("day", day))
		return msg.ErrCode_CONDITION_NOT_MET, itemsList
	}

	chargePackageInfo := p.UserData.BaseInfo.FirstChargePackage
	isObtained := false
	recordIdx := -1
	var loginCnt int32 = 0
	for i := 0; i < len(chargePackageInfo); i++ {
		if chargePackageInfo[i].Id == configId {
			for j := 0; j < len(chargePackageInfo[i].State); j++ {
				if j == (int(dayIdx)) && chargePackageInfo[i].State[j] == Obtained {
					isObtained = true
					recordIdx = i
					break
				}
			}
			loginCnt = chargePackageInfo[i].LoginCount
			break
		}
	}

	if loginCnt < int32(day) {
		return msg.ErrCode_CONDITION_NOT_MET, itemsList
	}

	if !isObtained || recordIdx < 0 {
		return msg.ErrCode_CONDITION_NOT_MET, itemsList
	}

	for _, v := range itemsList {
		AddPlayerItem(p, v.ItemId, int32(v.ItemNum), publicconst.ItemSource(publicconst.Charge_First), true)
	}

	p.UserData.BaseInfo.FirstChargePackage[recordIdx].State[dayIdx] = Receive
	p.SaveBaseInfo()
	return msg.ErrCode_SUCC, itemsList
}

func HandleUpdateChargePackageState(p *player.Player) {
	chargeInfo := p.UserData.BaseInfo.FirstChargePackage
	isUpdate := false
	for i := 0; i < len(chargeInfo); i++ {
		chargeInfo[i].LoginCount++
		loginCount := chargeInfo[i].LoginCount
		loginIdx := loginCount - 1
		if loginCount > 0 && loginCount <= int32(len(chargeInfo[i].State)) && (chargeInfo[i].State[loginIdx] == Forbidden) {
			chargeInfo[i].State[loginIdx] = Obtained
			isUpdate = true
		}
	}

	if isUpdate {
		p.SaveBaseInfo()
	}
}

func genShopItem(p *player.Player, shopItem *template.JShopItem, curTime uint32) *model.ShopItem {
	lock := shopItemIsLock(p, shopItem)
	item := model.NewShopItem(shopItem.Id, lock, shopItem.ConfigVer)
	if shopItem.IsFix != 1 {
		getItem := shopItem.RefreshGetItem(p.Rand)
		if getItem == nil {
			log.Error("genShopItem getItem is null", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("shopItemId", item.Id))
		} else {
			item.ItemId = getItem.ItemId
			item.ItemNum = getItem.ItemNum
		}
	} else {
		item.ItemId = shopItem.GetItem[0].ItemId
		item.ItemNum = shopItem.GetItem[0].ItemNum
	}
	item.UpdateTime = curTime
	item.NextRefreshTime = GetNextRefreshTime(shopItem)
	return item
}

func GetNextRefreshTime(s *template.JShopItem) uint32 {
	var refreshTime uint32
	switch s.RefreshArgs[0] {
	case 1:
		refreshTime = tools.GetDailyRefreshTime()
	case 2:
		refreshTime = tools.GetWeeklyRefreshTime(template.GetSystemItemTemplate().RefreshHour)
	case 3:
		refreshTime = tools.GetMonthRefreshTime()
	case 4:
	case 5:
		var minTime uint32
		for i := 1; i < len(s.RefreshArgs); i++ {
			temp := tools.GetHourRefreshTime(s.RefreshArgs[i])
			if i == 1 {
				minTime = temp
			} else {
				if temp < minTime {
					minTime = temp
				}
			}
		}
		refreshTime = minTime
	}
	return refreshTime
}
