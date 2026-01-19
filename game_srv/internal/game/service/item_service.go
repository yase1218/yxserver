package service

import (
	"gameserver/internal/game/event"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tda"
	"kernel/tools"
	"msg"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

var (
	useItemMap map[publicconst.UseItemType]UseItemFunc
)

type UseItemFunc func(*player.Player, uint32, uint32, publicconst.ItemSource, bool, []*model.SimpleItem, bool) (msg.ErrCode, []*model.SimpleItem)

func init() {
	useItemMap = make(map[publicconst.UseItemType]UseItemFunc)
	useItemMap[publicconst.USE_ITEM_COMPOSE] = useItemCommon
	useItemMap[publicconst.USE_ITEM_DECOMPOSE] = useItemCommon
	useItemMap[publicconst.USE_ITEM_COST] = useItemCommon
	useItemMap[publicconst.USE_ITEM_SELECT] = selectItem
	useItemMap[publicconst.USE_ITEM_FROM_ITEM_GROUP] = getItemFromItemGroup
	useItemMap[publicconst.USE_ITEM_GET_SHIP] = useItemGetShip
	useItemMap[publicconst.USE_ITEM_ADD_PLAYMETHOD_TIMES] = useItemAddPlayMethodTimes
	useItemMap[publicconst.USE_ITEM_ADD_Equip] = useItemAddEquip
	useItemMap[publicconst.USE_ITEM_ADD_PET] = useItemAddPet
	useItemMap[publicconst.USE_ITEM_ADV_CARD_PACKAGE] = useItemAdvCardPackage
	useItemMap[publicconst.USE_ITEM_ADV_CARD] = useItemAdvCard
	useItemMap[publicconst.USE_ITEM_ADD_GEM] = useItemAddGem
}

func _useItem(p *player.Player, itemId, num uint32, source publicconst.ItemSource, itemTemplate *template.Item, inPacket bool, items []*model.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	if f := useItemMap[publicconst.UseItemType(itemTemplate.EffectType)]; f != nil {
		return f(p, itemId, num, source, inPacket, items, ntf)
	}
	return msg.ErrCode_ITEM_NOT_USE, nil
}

func useItemAddGem(p *player.Player, itemId, num uint32, source publicconst.ItemSource, inPacket bool, items []*model.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	itemTemplate := template.GetItemTemplate().GetItem(itemId)
	if itemTemplate == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	if len(itemTemplate.EffectArgs) == 0 {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	// 消耗道具
	if inPacket {
		if res := CostItem(p.GetUserId(), itemId, num, source, ntf); res != msg.ErrCode_SUCC {
			return res, nil
		}
	}

	gemId := itemTemplate.EffectArgs[0]
	gem := template.GetGemTemplate().GetGem(gemId)
	if gem == nil {
		return msg.ErrCode_EQUIP_NOT_EXIST, nil
	}

	var retItems []*model.SimpleItem
	retItems = append(retItems, &model.SimpleItem{
		Id:  itemId,
		Num: num,
	})

	changeGem := AddGem(p, gemId, num, 0, 0, true)
	if changeGem == nil {
		log.Error("useItemAddGem err",
			zap.Uint64("accountId", p.GetUserId()),
			zap.Uint32("gemId", gemId))
		return msg.ErrCode_SUCC, nil
	}

	NotifyGemChange(p, changeGem)
	return msg.ErrCode_SUCC, retItems
}

func useItemAdvCard(p *player.Player, itemId, num uint32, source publicconst.ItemSource, inPacket bool, items []*model.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	itemTemplate := template.GetItemTemplate().GetItem(itemId)
	if itemTemplate == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	if len(itemTemplate.EffectArgs) == 0 {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	// 消耗道具
	if inPacket {
		if res := CostItem(p.GetUserId(), itemId, num, source, ntf); res != msg.ErrCode_SUCC {
			return res, nil
		}
	}

	cardId := itemTemplate.EffectArgs[0]
	advCard := template.GetAdvCardTemplate().GetCardById(int(cardId))
	if advCard == nil {
		return msg.ErrCode_EQUIP_NOT_EXIST, nil
	}

	var retItems []*model.SimpleItem
	retItems = append(retItems, &model.SimpleItem{
		Id:  itemId,
		Num: num,
	})

	return msg.ErrCode_SUCC, retItems
}

func useItemAdvCardPackage(p *player.Player, itemId, num uint32, source publicconst.ItemSource, inPacket bool, items []*model.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	itemTemplate := template.GetItemTemplate().GetItem(itemId)
	if itemTemplate == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	if len(itemTemplate.EffectArgs) == 0 {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	// 消耗道具
	var notifyItems []uint32
	if inPacket {
		if res := CostItem(p.GetUserId(), itemId, num, source, ntf); res != msg.ErrCode_SUCC {
			return res, nil
		}
		notifyItems = append(notifyItems, itemId)
	}

	// 获取卡包ID
	bagId := int(itemTemplate.EffectArgs[0])
	// 获取卡包配置
	bagTemplate := template.GetAdvCardBagTemplate().GetBagById(bagId)
	if bagTemplate == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	var retItems []*model.SimpleItem

	// 如果有指定卡牌，直接使用指定卡牌
	if len(bagTemplate.CradID) > 0 {
		randomCards := make(map[uint32]uint32)
		for _, cardId := range bagTemplate.CradID {
			randomCards[uint32(cardId)]++
		}

		// 将随机得到的卡牌添加到玩家背包
		for cardId, cardNum := range randomCards {
			// 添加卡牌到玩家背包

			AddItem(p.GetUserId(), cardId, int32(cardNum), source, ntf)
			notifyItems = append(notifyItems, cardId)

			// 添加到返回结果中
			retItems = append(retItems, &model.SimpleItem{
				Id:  cardId,
				Num: cardNum,
			})
		}
	} else {
		// 使用RandomWeightItem实现随机卡牌
		// 根据Num确定随机次数
		randomTimes := 1
		if len(bagTemplate.Num) > 0 {
			randomTimes = bagTemplate.Num[0]
		}

		// 先根据品质权重随机选择品质
		for i := 0; i < randomTimes; i++ {
			// 创建品质的权重项
			var qualityWeightItems []*template.WeightItem
			for qualityIndex, weight := range bagTemplate.QualityWeight {
				qualityWeightItems = append(qualityWeightItems, &template.WeightItem{
					ItemId:  uint32(qualityIndex),
					ItemNum: 1,
					Weight:  uint32(weight),
				})
			}

			// 随机选择品质
			qualityItems := template.RandomWeightItem(p.Rand, qualityWeightItems, 1)
			if len(qualityItems) == 0 {
				continue
			}

			quality := qualityItems[0].ItemId

			// 根据BaseId查找对应品质的卡牌
			var cardWeightItems []*template.WeightItem
			for _, baseId := range bagTemplate.BaseId {
				cards := template.GetAdvCardTemplate().GetCardsByBaseId(baseId)
				for _, card := range cards {
					// 检查卡牌对应的道具品质是否匹配
					cardItem := template.GetItemTemplate().GetItem(uint32(card.Id))
					if cardItem != nil && cardItem.Quality == quality {
						cardWeightItems = append(cardWeightItems, &template.WeightItem{
							ItemId:  uint32(card.Id),
							ItemNum: 1,
							Weight:  uint32(card.Weight),
						})
					}
				}
			}

			// 如果找到符合条件的卡牌，根据权重随机选择
			if len(cardWeightItems) > 0 {
				randomCardItems := template.RandomWeightItem(p.Rand, cardWeightItems, 1)

				for _, item := range randomCardItems {
					// 添加卡牌到玩家背包
					AddItem(p.GetUserId(), item.ItemId, int32(item.ItemNum), source, ntf)
					notifyItems = append(notifyItems, item.ItemId)

					// 添加到返回结果中
					retItems = append(retItems, &model.SimpleItem{
						Id:  item.ItemId,
						Num: item.ItemNum,
					})
				}
			} else {
				// 没有指定卡牌，随机获取一张
				log.Error("[bug开卡包] 没有指定卡牌，随机获取一张",
					zap.Uint64("accountid", p.GetUserId()),
				)
			}
		}
		if randomTimes == 1 {
			log.Error("[bug开卡包] 就是填了一张卡牌?",
				zap.Uint64("accountid", p.GetUserId()),
			)
		}
		if len(retItems) != randomTimes {
			log.Error("[bug开卡包] 随机出来的卡牌数量和随机次数不匹配",
				zap.Uint64("accountid", p.GetUserId()),
			)
		}
	}

	// 通知客户端道具变化
	if len(notifyItems) > 0 && ntf {
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	return msg.ErrCode_SUCC, retItems
}

func useItemAddPet(p *player.Player, itemId, num uint32, source publicconst.ItemSource, inPacket bool, items []*model.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	itemTemplate := template.GetItemTemplate().GetItem(itemId)
	log.Debug("useItemAddPet", zap.Uint64("uid", p.GetUserId()), zap.Uint32("itemId", itemId), zap.Uint32("num", num))
	if itemTemplate == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	// if len(itemTemplate.EffectArgs) == 0 {
	// 	return msg.ErrCode_SYSTEM_ERROR, nil
	// }

	// 消耗道具
	var notifyItems []uint32
	if inPacket {
		if res := CostItem(p.GetUserId(), itemId, num, source, ntf); res != msg.ErrCode_SUCC {
			return res, nil
		}
		notifyItems = append(notifyItems, itemId)
	}

	petId := itemTemplate.Id
	petConfig := template.GetPetTemplate().GetPet(petId)
	if petConfig == nil {
		return msg.ErrCode_EQUIP_NOT_EXIST, nil
	}

	var retItems []*model.SimpleItem
	retItems = append(retItems, &model.SimpleItem{
		Id:  itemId,
		Num: num,
	})

	AddPetByItem(p, petId, ntf)
	UpdateTask(p, true, publicconst.TASK_COND_GET_PET, petConfig.Rarity, num) // 激活XX个XX品质库鲁兽
	return msg.ErrCode_SUCC, retItems
}

func useItemAddEquip(p *player.Player, itemId, num uint32, source publicconst.ItemSource, inPacket bool, items []*model.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	itemTemplate := template.GetItemTemplate().GetItem(itemId)
	if itemTemplate == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	if len(itemTemplate.EffectArgs) == 0 {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	// 消耗道具
	if inPacket {
		if res := CostItem(p.GetUserId(), itemId, num, source, ntf); res != msg.ErrCode_SUCC {
			return res, nil
		}
	}

	equipId := itemTemplate.EffectArgs[0]
	equip := template.GetEquipTemplate().GetEquip(equipId)
	if equip == nil {
		return msg.ErrCode_EQUIP_NOT_EXIST, nil
	}

	var retItems []*model.SimpleItem
	retItems = append(retItems, &model.SimpleItem{
		Id:  itemId,
		Num: num,
	})

	changeEquip := AddEquip(p, equipId, num, true)
	if changeEquip == nil {
		log.Error("useItemAddEquip err", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("equipId", equipId))
		return msg.ErrCode_SUCC, nil
	}

	NotifyEquipChange(p, changeEquip)
	return msg.ErrCode_SUCC, retItems
}

// useItemCommon 通用使用道具
func useItemCommon(p *player.Player, itemId, num uint32, source publicconst.ItemSource, inPacket bool, items []*model.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	accountId := p.GetUserId()
	itemTemplate := template.GetItemTemplate().GetItem(itemId)
	if itemTemplate == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	// 消耗道具
	if inPacket {
		if res := CostItem(accountId, itemId, num, source, ntf); res != msg.ErrCode_SUCC {
			return res, nil
		}
	}

	var reward []*model.SimpleItem
	rewardMap := make(map[uint32]uint32)
	// 生成道具
	for i := 0; i < len(itemTemplate.EffectArgs); i += 2 {
		itemId := itemTemplate.EffectArgs[i]
		num := itemTemplate.EffectArgs[i+1] * num
		if _, ok := rewardMap[itemId]; ok {
			rewardMap[itemId] += num
		} else {
			rewardMap[itemId] = num
		}
	}

	for id, num := range rewardMap {
		reward = append(reward, &model.SimpleItem{
			Id:  id,
			Num: num,
		})
	}

	return msg.ErrCode_SUCC, reward
}

func selectItem(p *player.Player, itemId, num uint32, source publicconst.ItemSource, inPacket bool, selectItems []*model.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	itemTemplate := template.GetItemTemplate().GetItem(itemId)
	if itemTemplate == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	if len(itemTemplate.EffectArgs) == 0 {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	var total uint32 = 0
	for i := 0; i < len(selectItems); i++ {
		total += selectItems[i].Num
	}

	if num != total {
		return msg.ErrCode_INVALID_DATA, nil
	}

	itemGroup := template.GetItemGroupTemplate().GetItemGroup(itemTemplate.EffectArgs[0])
	if itemGroup == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	var rewardItems []*model.SimpleItem
	for i := 0; i < len(selectItems); i++ {
		if item := itemGroup.GetItem(selectItems[i].Id); item != nil {
			rewardItems = append(rewardItems, &model.SimpleItem{
				Id:  selectItems[i].Id,
				Num: selectItems[i].Num * item.ItemNum,
			})
		} else {
			return msg.ErrCode_INVALID_DATA, nil
		}
	}

	var notifyItems []uint32
	if res := CostItem(p.GetUserId(), itemId, num, source, ntf); res != msg.ErrCode_SUCC {
		return res, nil
	}
	notifyItems = append(notifyItems, itemId)

	var reward []*model.SimpleItem
	rewardMap := make(map[uint32]uint32)
	// 生成道具
	for i := 0; i < len(rewardItems); i++ {
		addItems := AddItem(p.GetUserId(), rewardItems[i].Id,
			int32(rewardItems[i].Num), source, ntf)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		if _, ok := rewardMap[rewardItems[i].Id]; ok {
			rewardMap[rewardItems[i].Id] += rewardItems[i].Num
		} else {
			rewardMap[rewardItems[i].Id] = rewardItems[i].Num
		}
	}

	for id, num := range rewardMap {
		reward = append(reward, &model.SimpleItem{
			Id:  id,
			Num: num,
		})
	}
	if ntf {
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	return msg.ErrCode_SUCC, reward
}

func getItemFromItemGroup(p *player.Player, itemId, num uint32, source publicconst.ItemSource, inPacket bool, selectItems []*model.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	accountId := p.GetUserId()
	itemTemplate := template.GetItemTemplate().GetItem(itemId)
	if itemTemplate == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	// 消耗道具
	var notifyItems []uint32
	if inPacket {
		if res := CostItem(accountId, itemId, num, source, ntf); res != msg.ErrCode_SUCC {
			return res, nil
		}
		notifyItems = append(notifyItems, itemId)
	}

	var items []*template.SimpleItem
	for i := 0; i < len(itemTemplate.EffectArgs); i += 1 {
		if group := template.GetItemGroupTemplate().GetItemGroup(itemTemplate.EffectArgs[i]); group != nil {
			items = append(items, group.GenItem(p.Rand, num)...)
		}
	}

	var reward []*model.SimpleItem
	rewardMap := make(map[uint32]*model.SimpleItem)
	for i := 0; i < len(items); i++ {
		addItems := AddItem(accountId, items[i].ItemId, int32(items[i].ItemNum), publicconst.UseAddItem, ntf)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)

		for m := 0; m < len(addItems); m++ {
			if data, ok := rewardMap[addItems[m].Id]; ok {
				data.Num += addItems[m].Num
			} else {
				rewardMap[addItems[m].Id] = &model.SimpleItem{
					Id:  addItems[m].Id,
					Num: addItems[m].Num,
					Src: addItems[m].Src,
				}
			}
		}
	}

	for id, data := range rewardMap {
		reward = append(reward, &model.SimpleItem{
			Id:  id,
			Num: data.Num,
			Src: data.Src,
		})
	}
	//	event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(playerData, notifyItems, ListenNotifyClientItemEventEvent))

	if ntf {
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	return msg.ErrCode_SUCC, reward
}

func useItemGetShip(p *player.Player, itemId, num uint32, source publicconst.ItemSource, inPacket bool, items []*model.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	itemTemplate := template.GetItemTemplate().GetItem(itemId)
	if itemTemplate == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	if len(itemTemplate.EffectArgs) == 0 {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	// 消耗道具
	var notifyItems []uint32
	if inPacket {
		if res := CostItem(p.GetUserId(), itemId, num, source, ntf); res != msg.ErrCode_SUCC {
			return res, nil
		}
		notifyItems = append(notifyItems, itemId)
	}

	var retItems []*model.SimpleItem
	shipId := itemTemplate.EffectArgs[0]

	bAddRaw := false
	if ship := getShip(p, shipId); ship == nil {
		AddShip(p, shipId, source, true)
		num -= 1
		bAddRaw = true
	}
	if num > 0 {
		shipConfig := template.GetShipTemplate().GetShip(shipId)
		for i := 0; i < len(shipConfig.DecomposeItem); i++ {
			temp := AddItem(p.GetUserId(),
				shipConfig.DecomposeItem[i].ItemId,
				int32(shipConfig.DecomposeItem[i].ItemNum*num),
				source,
				ntf)
			retItems = append(retItems, temp...)
			notifyItems = append(notifyItems, shipConfig.DecomposeItem[i].ItemId)
		}
		UpdateTask(p, true, publicconst.TASK_COND_ACTIVE_SHIP, shipConfig.Rarity, 1) // 激活XX个XX品质库鲁
	}

	// 重复转换来的
	for i := 0; i < len(retItems); i++ {
		retItems[i].Src = 1
	}

	if bAddRaw {
		retItems = append(retItems, &model.SimpleItem{
			Id:  itemId,
			Num: 1,
		})
	}
	event.EventMgr.PublishEvent(event.NewShipChangeEvent(p, shipId, 1,
		1, publicconst.Ship_Star_Level, ListenShipChangeEvent))
	return msg.ErrCode_SUCC, retItems
}

func useItemAddPlayMethodTimes(p *player.Player, itemId, num uint32, source publicconst.ItemSource, inPacket bool, items []*model.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	itemTemplate := template.GetItemTemplate().GetItem(itemId)
	if itemTemplate == nil {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	if len(itemTemplate.EffectArgs) == 0 {
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	var retItems []*model.SimpleItem
	// 消耗道具
	if inPacket {
		if res := CostItem(p.GetUserId(), itemId, num, source, ntf); res != msg.ErrCode_SUCC {
			return res, nil
		}
		retItems = append(retItems, &model.SimpleItem{Id: itemId, Num: num})
	}

	btType := itemTemplate.EffectArgs[0]
	AddPlayMethodTimes(p, int(btType), int(num))

	return msg.ErrCode_SUCC, retItems
}

func EnoughItem(uid uint64, itemId, num uint32) bool {
	playerData := player.FindByUserId(uid)
	if playerData == nil {
		log.Error("EnoughItem err", zap.Uint64("uid", uid))
		return false
	}

	item := findItem(playerData, itemId)
	if item == nil {
		return false
	}

	if num == 0 {
		return true
	}

	if item.Num < uint64(num) {
		return false
	}
	return true
}

func findItem(p *player.Player, itemId uint32) *model.Item {
	for i := 0; i < len(p.UserData.Items.Items); i++ {
		if p.UserData.Items.Items[i].Id == itemId {
			return p.UserData.Items.Items[i]
		}
	}
	return nil
}

func GetItemNum(uid uint64, itemId uint32) uint64 {
	p := player.FindByUserId(uid)
	if p == nil {
		log.Error("GetItemNum err", zap.Uint64("accountId", uid))
		return 0
	}

	item := findItem(p, itemId)
	if item == nil {
		return 0
	}

	return item.Num
}

func AddPlayerItem(p *player.Player, itemId uint32, num int32, source publicconst.ItemSource, notifyClient bool) []*model.SimpleItem {
	itemConfig := template.GetItemTemplate().GetItem(itemId)
	if itemConfig == nil {
		log.Error("AddItem err", zap.String("accountId", p.GetOpenId()), zap.Uint32("itemId", itemId), zap.Int32("num", num), zap.Int64("source", int64(source)))
		return nil
	}

	// 立即使用
	if itemConfig.Use == 1 {
		if err, items := _useItem(p, itemId, uint32(num), source, itemConfig, false, nil, notifyClient); err == msg.ErrCode_SUCC {
			event.EventMgr.PublishEvent(event.NewAddItemEvent(p, itemId, uint32(num), 0, "", source, notifyClient, ListenAddItemEvent))
			if notifyClient && len(items) > 0 {
				var notifyItems []uint32
				notifyItems = append(notifyItems, GetSimpleItemIds(items)...)
				if notifyClient {
					updateClientItemsChange(p.GetUserId(), notifyItems)
				}
			}

			return items
		}
		return nil
	}

	// 封装一层添加数量,方便进行一些处理
	addNum := uint32(num)

	// switch itemId {
	// case uint32(publicconst.ITEM_ALLIANCE_EXP):
	// 	// 检查玩家是否已有联盟
	// 	member, err := dao.GetMember(p.UserData.AccountId)
	// 	if err == nil && member != nil {
	// 		ServMgr.GetAllianceService().AddAllianceExp(member.AllianceID, num)
	// 	}
	// 	return nil
	// case enum.ItemID_Alliance_Medal:
	// 	dbAllianceMember, err := dao.GetMember(p.UserData.AccountId)
	// 	if err == nil && dbAllianceMember != nil {
	// 		// alliance active rank
	// 		//if err = rdb.AddAllianceRankByType(dbAllianceMember.AllianceID,
	// 		//	msg.AllianceRankType_Alliance_Rank_Single_Active, accountId, float64(num)); err != nil {
	// 		//	log.Error("add alliance rank single active err", zap.Error(err),
	// 		//		zap.Int64("accountId", accountId), zap.Uint32("allianceId", dbAllianceMember.AllianceID))
	// 		//}

	// 		if err := ServMgr.GetRankService().AddAllianceRank(dbAllianceMember.AllianceID,
	// 			msg.AllianceRankType_Alliance_Rank_Single_Active, accountId, uint64(num)); err != nil {
	// 			log.Error("add alliance rank err", zap.Error(err),
	// 				zap.Uint64("accountId", p.GetUserId()),
	// 				zap.Uint32("allianceId", dbAllianceMember.AllianceID))
	// 		}

	// 		// alliance weekly task active
	// 		dbAlliance, err := dao.GetAlliance(dbAllianceMember.AllianceID)
	// 		if err == nil {
	// 			templateAlliance := template.GetGuildTemplate().GetAlliance(dbAlliance.Level)
	// 			if templateAlliance != nil {
	// 				itemMax := uint32(templateAlliance.ItemMax)
	// 				if p.UserData.Task.AllianceWeeklyActiveValue < itemMax {
	// 					if p.UserData.Task.AllianceWeeklyActiveValue+uint32(num) > itemMax {
	// 						addNum = itemMax - p.UserData.Task.AllianceWeeklyActiveValue
	// 						p.UserData.Task.AllianceWeeklyActiveValue = itemMax
	// 					} else {
	// 						p.UserData.Task.AllianceWeeklyActiveValue += uint32(num)
	// 					}
	// 					dao.TaskDao.UpdateAllianceWeeklyActiveValue(accountId, p.UserData.Task.AllianceWeeklyActiveValue)
	// 					NotifyActiveChange(p)
	// 				} else {
	// 					return nil
	// 				}

	// 				// tda
	// 				tda.TdaGuildActiveNess(p.ChannelId, p.TdaCommonAttr, dbAlliance.Level, dbAlliance.ID, dbAlliance.Name, dbAlliance.MemberCount, uint32(num))
	// 			} else {
	// 				log.Error("table alliance nil", zap.Uint32("lv", dbAlliance.Level))
	// 			}
	// 		}
	// 	}
	// }

	var ret []*model.SimpleItem
	item := findItem(p, itemId)
	if item != nil {
		item.Num += uint64(addNum)
		switch itemId {
		case uint32(publicconst.ITEM_CODE_AP):
			if item.Num >= uint64(template.GetSystemItemTemplate().ApMax) {
				addNum = addNum - uint32(item.Num-uint64(template.GetSystemItemTemplate().ApMax)) // 实际增加数量
				item.Num = uint64(template.GetSystemItemTemplate().ApMax)
			}
		case uint32(publicconst.ITEM_CODE_EXP):
			//if item.Num >= int64(template.GetLevelTemplate().GetMaxExp()) {
			//	item.Num = int64(template.GetLevelTemplate().GetMaxExp())
			//}
		}

		p.SaveItems()

	} else {
		item = model.NewItem(itemId, 0, addNum)
		p.UserData.Items.Items = append(p.UserData.Items.Items, item)
		p.SaveItems()

	}

	ret = append(ret, &model.SimpleItem{
		Id:  itemId,
		Num: addNum,
	})

	if source != publicconst.InitAddItem {
		event.EventMgr.PublishEvent(event.NewAddItemEvent(p, itemId, addNum, int64(item.Num), "", source, notifyClient, ListenAddItemEvent))
	} else {
		if itemConfig.BigType == 20 {
			// TODO 老宠物移除暂留
			// PutPetEgg(p, itemId)
		}
	}
	UpdateTask(p, true, publicconst.TASK_COND_ADD_DISK, itemConfig.BigType, addNum) // 获得XX个磁盘
	//processHistoryData(p, publicconst.TASK_COND_ADD_DISK, 0, addNum)
	UpdateTask(p, true, publicconst.TASK_COND_ADD_RARITY_DISK, itemConfig.BigType, itemConfig.Quality, addNum) // 获得X个XX品质磁盘

	tda.TdaItemChange(p.ChannelId, p.TdaCommonAttr, "0", itemId, uint32(num), uint32(source), uint32(item.Num))

	return ret
}

func AddItem(uid uint64, itemId uint32, num int32, source publicconst.ItemSource, notifyClient bool) []*model.SimpleItem {
	if num <= 0 {
		log.Error("AddItem err", zap.Uint64("accountId", uid), zap.Uint32("itemId", itemId), zap.Int32("num", num), zap.Int64("source", int64(source)))
		return nil
	}
	p := player.FindByUserId(uid)
	if p == nil {
		log.Error("AddItem err", zap.Uint64("accountId", uid))
		return nil
	}

	return AddPlayerItem(p, itemId, num, source, notifyClient)
}

func CostItem(uid uint64, itemId, num uint32, source publicconst.ItemSource, notifyClient bool) (res msg.ErrCode) {
	res = msg.ErrCode_SUCC
	p := player.FindByUserId(uid)
	if p == nil {
		log.Error("CostItem err", zap.Uint64("uid", uid))
		res = msg.ErrCode_PLAYER_NOT_EXIST
		return
	}

	item := findItem(p, itemId)
	if item == nil {
		log.Error("CostItem err", zap.Uint64("uid", uid), zap.Uint32("itemId", itemId), zap.Uint32("num", num), zap.Int64("source", int64(source)))
		res = msg.ErrCode_ITEM_NOT_EXIST
		return
	}

	if item.Num < uint64(num) {
		log.Error("CostItem err", zap.Uint64("uid", uid), zap.Uint32("itemId", itemId), zap.Uint32("num", num), zap.Uint64("item.Num", item.Num))
		res = msg.ErrCode_NO_ENOUGH_ITEM
		return
	}
	item.Num -= uint64(num)
	p.SaveItems()

	// 体力
	// if itemId == uint32(publicconst.ITEM_CODE_AP) {
	// 	recoveryAp(p)
	// }
	event.EventMgr.PublishEvent(event.NewCostItemEvent(p, itemId, num, int64(item.Num), "", source, notifyClient, ListenCostItemEvent))

	tda.TdaItemChange(p.ChannelId, p.TdaCommonAttr, "1", itemId, num, uint32(source), uint32(item.Num))
	return
}

// UseItem 使用道具
func UseItem(p *player.Player, itemId, num uint32, source publicconst.ItemSource, selectItems []*msg.SimpleItem, ntf bool) (msg.ErrCode, []*model.SimpleItem) {
	item := findItem(p, itemId)
	if item == nil {
		return msg.ErrCode_ITEM_NOT_EXIST, nil
	}

	if num == 0 {
		return msg.ErrCode_INVALID_DATA, nil
	}

	var items []*model.SimpleItem
	for i := 0; i < len(selectItems); i++ {
		items = append(items, &model.SimpleItem{
			Id:  selectItems[i].ItemId,
			Num: selectItems[i].ItemNum,
		})
	}

	itemTemplate := template.GetItemTemplate().GetItem(itemId)
	if itemTemplate == nil {
		return msg.ErrCode_ITEM_NOT_EXIST, nil
	}

	if itemTemplate.Lv > 0 && p.UserData.Level < itemTemplate.Lv { // 库鲁等级不够
		return msg.ErrCode_COND_LEVEL_LOW, nil
	}

	if !EnoughItem(p.GetUserId(), itemId, num) {
		return msg.ErrCode_NO_ENOUGH_ITEM, nil
	}

	return _useItem(p, itemId, num, source, itemTemplate, true, items, ntf)
}

func updateClientItemsChange(uid uint64, itemIds []uint32) {
	if len(itemIds) == 0 {
		return
	}
	if p := player.FindByUserId(uid); p != nil {
		res := &msg.NotifyUpdateItem{}
		itemIds = tools.ListRemoveRepeat(itemIds) // 去重

		var sendItems []uint32
		for k := 0; k < len(itemIds); k++ {
			if tools.ListContain(sendItems, itemIds[k]) {
				continue
			}
			sendItems = append(sendItems, itemIds[k])
			if item := findItem(p, itemIds[k]); item != nil {
				res.Items = append(res.Items, ToProtocolItem(item))
			}
		}
		if len(sendItems) > 0 {
			p.SendNotify(res)
		}
	}
}

func GetSimpleItemIds(data []*model.SimpleItem) []uint32 {
	var ret []uint32
	for i := 0; i < len(data); i++ {
		ret = append(ret, data[i].Id)
	}
	return ret
}

func ToProtocolItem(item *model.Item) *msg.Item {
	return &msg.Item{
		ItemId:    item.Id,
		ItemNum:   int64(item.Num),
		LimitDate: item.LimitDate,
	}
}
func ToProtocolSimpleItem(item *model.SimpleItem) *msg.SimpleItem {
	return &msg.SimpleItem{
		ItemId:  item.Id,
		ItemNum: item.Num,
	}
}

func ToTemplateItem(items []*model.SimpleItem) []*template.SimpleItem {
	var ret []*template.SimpleItem
	for i := 0; i < len(items); i++ {
		ret = append(ret, &template.SimpleItem{
			ItemId:  items[i].Id,
			ItemNum: items[i].Num,
			Src:     items[i].Src,
		})
	}
	return ret
}

func ToProtocolItems(items []*model.Item) []*msg.Item {
	var ret []*msg.Item
	for i := 0; i < len(items); i++ {
		ret = append(ret, ToProtocolItem(items[i]))
	}
	return ret
}

func TemplateSimpleItemToProtocolSImpleItems(items []*template.SimpleItem) []*msg.SimpleItem {
	var ret []*msg.SimpleItem
	for i := 0; i < len(items); i++ {
		ret = append(ret, &msg.SimpleItem{
			ItemId:  items[i].ItemId,
			ItemNum: items[i].ItemNum,
		})
	}
	return ret
}

func TemplateItemToProtocolItems(items []*template.SimpleItem) []*msg.SimpleItem {

	var ret []*msg.SimpleItem
	for i := 0; i < len(items); i++ {
		ret = append(ret, &msg.SimpleItem{
			ItemId:  items[i].ItemId,
			ItemNum: items[i].ItemNum,
			Src:     items[i].Src,
		})
	}
	return ret
}

func ToProtocolSimpleItems(items []*model.SimpleItem) []*msg.SimpleItem {
	var ret []*msg.SimpleItem
	mapItem := make(map[uint32]*model.SimpleItem)
	for i := 0; i < len(items); i++ {
		if _, ok := mapItem[items[i].Id]; ok {
			mapItem[items[i].Id].Num += items[i].Num
		} else {
			mapItem[items[i].Id] = &model.SimpleItem{
				Id:  items[i].Id,
				Num: items[i].Num,
			}
		}
	}

	for _, data := range mapItem {
		ret = append(ret, ToProtocolSimpleItem(data))
	}
	return ret
}

func ProtocolSimpleItemsToItems(items []*model.SimpleItem) []*msg.Item {
	var ret []*msg.Item
	for i := 0; i < len(items); i++ {
		ret = append(ret, &msg.Item{
			ItemId:  items[i].Id,
			ItemNum: int64(items[i].Num),
		})
	}
	return ret
}

func mergeSimpleItem(data []*template.SimpleItem) []*template.SimpleItem {
	var ret []*template.SimpleItem

	itemMap := make(map[uint32]uint32)
	for i := 0; i < len(data); i++ {
		if num, ok := itemMap[data[i].ItemId]; ok {
			itemMap[data[i].ItemId] = num + data[i].ItemNum
		} else {
			itemMap[data[i].ItemId] = data[i].ItemNum
		}
	}

	for id, num := range itemMap {
		ret = append(ret, &template.SimpleItem{
			ItemId:  id,
			ItemNum: num,
		})
	}
	return ret
}

func NotifyEquipChange(p *player.Player, equip *model.Equip) {
	p.SendNotify(&msg.NotifyEquipChange{
		Data: ToProtocolEquips([]*model.Equip{equip}),
	})
}

// UpdateClientItemChange 通知客户端道具变化
func updateClientItemChange(p *player.Player, itemId uint32) {
	if item := findItem(p, itemId); item != nil {
		res := &msg.NotifyUpdateItem{}
		res.Items = append(res.Items, ToProtocolItem(item))
		p.SendNotify(res)
	}
}

func ToProtocolGems(data []*model.GemBagSlot) []*msg.Gem {
	var ret []*msg.Gem
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolGem(data[i]))
	}
	return ret
}

func ToProtocolGemmap(data map[uint64]*model.GemBagSlot) []*msg.Gem {
	var ret []*msg.Gem
	for _, v := range data {
		ret = append(ret, ToProtocolGem(v))
	}
	return ret
}

func ToProtocolGem(data *model.GemBagSlot) *msg.Gem {
	return &msg.Gem{
		Uuid: data.Uuid,
		Num:  data.Num,
		Lock: data.Lock,
	}
}

func ToProtocolGemPos(pos [][]uint64) []*msg.GemPos {
	var ret []*msg.GemPos
	for i := 0; i < len(pos); i++ {
		ret = append(ret, &msg.GemPos{
			Uuid: pos[i],
		})
	}
	return ret
}

func ToProtocolEquip(data *model.Equip) *msg.Equip {
	return &msg.Equip{
		Id:  data.Id,
		Num: data.Num,
	}
}

func ToProtocolSuit(data *model.AccountEquip) []*msg.SuitRewardInfo {
	var ret []*msg.SuitRewardInfo
	for i := 0; i < len(data.SuitReward); i++ {
		temp := &msg.SuitRewardInfo{}
		temp.SuitId = data.SuitReward[i].SuitId
		for k := 0; k < len(data.SuitReward[i].PosData); k++ {
			temp.Reward = append(temp.Reward, &msg.SuitPosReward{
				Pos:     data.SuitReward[i].PosData[k].Pos,
				EquipId: data.SuitReward[i].PosData[k].EquipId,
			})
		}
		ret = append(ret, temp)
	}
	return ret
}

func ToProtocolEquips(data []*model.Equip) []*msg.Equip {
	var ret []*msg.Equip
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolEquip(data[i]))
	}
	return ret
}

func ToProtocolEquipMap(data map[uint32]*model.Equip) []*msg.Equip {
	var ret []*msg.Equip
	for _, v := range data {
		ret = append(ret, ToProtocolEquip(v))
	}
	return ret
}

func ToProtocolEquipPos(data *model.EquipPos) *msg.UseEquip {
	if data.EquipId == 0 {
		return nil
	}
	jEquip := template.GetEquipTemplate().GetEquip(data.EquipId)
	if jEquip == nil {
		return nil
	}
	ret := &msg.UseEquip{
		Pos:      data.Pos,
		PosLevel: data.Level,
		EquipId:  data.EquipId,
	}

	if ret.PosLevel > jEquip.Data.LevelMax {
		ret.PosLevel = jEquip.Data.LevelMax
	}

	var ids []uint32
	for i := 0; i < len(jEquip.InitAttr); i++ {
		ids = append(ids, jEquip.InitAttr[i].Id)
	}

	for id, data := range data.Attr {
		if tools.ListContain(ids, id) {
			ret.Attrs = append(ret.Attrs, &msg.Attr{
				Id:        data.Id,
				Value:     data.InitValue + data.LevelValue + data.Add,
				CalcValue: data.FinalValue,
			})
		}
	}

	return ret
}

func ToProtocolEquipPosList(data []*model.EquipPos) []*msg.UseEquip {
	var ret []*msg.UseEquip
	for i := 0; i < len(data); i++ {
		if temp := ToProtocolEquipPos(data[i]); temp != nil {
			ret = append(ret, ToProtocolEquipPos(data[i]))
		}
	}
	return ret
}

func ToProtocolEquipSuit(data *model.EquipSuit) *msg.EquipSuit {
	if data == nil {
		return nil
	}
	return &msg.EquipSuit{
		SuitId:   data.SuitId,
		EquipIds: data.EquipIds,
		SkillId:  data.SkillId,
	}
}

func ToProtocolEquipSuits(data []*model.EquipSuit) []*msg.EquipSuit {
	var ret []*msg.EquipSuit
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolEquipSuit(data[i]))
	}
	return ret
}

// ComposeItem 道具合成
func ComposeItem(p *player.Player, itemId, composeItemID uint32, source publicconst.ItemSource) (msg.ErrCode, []*model.SimpleItem) {
	item := findItem(p, itemId)
	if item == nil {
		return msg.ErrCode_ITEM_NOT_EXIST, nil
	}

	itemTemplate := template.GetComposeItemTemplate().GetComposeItem(itemId)
	if itemTemplate == nil {
		return msg.ErrCode_ITEM_NOT_EXIST, nil
	}

	if itemTemplate.BigType == 1 { // 合成唯一
		composeItemID = itemTemplate.Items[0].ItemId
	} else { // 可选合成
		if !itemTemplate.IsInItems(composeItemID) { // 选择合成的道具不在配置表里
			log.Error("composeItemID is invalid", zap.Uint64("uid", p.GetUserId()), zap.Uint32("itemId", itemId), zap.Uint32("composeItemID", composeItemID))
			return msg.ErrCode_ITEM_NOT_USE, nil
		}
	}

	if item.Num < uint64(itemTemplate.ComposeItem.Number) {
		return msg.ErrCode_NO_ENOUGH_ITEM, nil
	}

	composeItem := itemTemplate.GetItem(composeItemID) // 合成的道具
	composeItemNum := composeItem.ItemNum              // 合成道具个数
	costItemNum := itemTemplate.ComposeItem.Number     // 消耗碎片个数
	// 消耗道具
	if res := CostItem(p.GetUserId(), itemId, costItemNum, source, true); res != msg.ErrCode_SUCC {
		return res, nil
	}
	// 添加道具
	reward := AddItem(p.GetUserId(), composeItemID, int32(composeItemNum), source, false)
	return msg.ErrCode_SUCC, reward
}
