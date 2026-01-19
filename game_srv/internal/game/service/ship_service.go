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
	template2 "github.com/zy/game_data/template"
	"go.uber.org/zap"
)

// GmUpgradeStarShip 通过Gm命令升星机甲
func GmUpgradeStarShip(p *player.Player, shipId, addLv uint32) msg.ErrCode {

	ship := getShip(p, shipId)
	if ship == nil {
		return msg.ErrCode_SHIP_NOT_EXIST
	}

	starLevel := template2.GetRoleStarTemplate().GetShipStarLevel(shipId, ship.StarLevel)
	if starLevel == nil {
		return msg.ErrCode_INVALID_DATA
	}

	nextStarLevel := template2.GetRoleStarTemplate().GetShipStarLevel(shipId, ship.StarLevel+addLv)
	if nextStarLevel == nil {
		return msg.ErrCode_SHIP_STAR_FULL
	}

	oldLevel := ship.StarLevel
	ship.StarLevel += addLv

	// 发布机甲变化事件
	event.EventMgr.PublishEvent(event.NewShipChangeEvent(p, shipId, oldLevel,
		ship.StarLevel, publicconst.Ship_Star_Level, ListenShipChangeEvent))

	p.SaveShips()

	// 添加升星属性
	addAllShipStarAttr(p, ship, int32(oldLevel))

	AddComboSkill(p, nextStarLevel.Data.PickSkillcombo)

	msg2 := &msg.NotifyShipsChange{}
	msg2.Data = append(msg2.Data, ToProtocolShip(ship, p))
	p.SendNotify(msg2)

	return msg.ErrCode_SUCC
}

// UpgradeStarShip 升星机甲
func UpgradeStarShip(p *player.Player, shipId uint32, upgradeMax bool) (msg.ErrCode, uint32, uint32) {
	ship := getShip(p, shipId)
	if ship == nil {
		return msg.ErrCode_SHIP_NOT_EXIST, 0, 0
	}

	shipConfig := template2.GetShipTemplate().GetShip(shipId)
	if shipConfig == nil {
		log.Error("ship cfg nil", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("shipId", shipId))
		return msg.ErrCode_CONFIG_NIL, 0, 0
	}

	starLevel := template2.GetRoleStarTemplate().GetShipStarLevel(shipId, ship.StarLevel)
	if starLevel == nil {
		return msg.ErrCode_INVALID_DATA, 0, 0
	}

	if template2.GetRoleStarTemplate().GetShipStarLevel(shipId, ship.StarLevel+1) == nil {
		return msg.ErrCode_SHIP_STAR_FULL, 0, 0
	}

	// for i := 0; i < len(nextStarLevel.CostItems); i++ {
	// 	if !ServMgr.GetItemService().EnoughItem(playerData.AccountInfo.AccountId,
	// 		nextStarLevel.CostItems[i].ItemId, nextStarLevel.CostItems[i].ItemNum) {
	// 		return msg.ErrCode_NO_ENOUGH_ITEM, 0
	// 	}
	// }

	reward := make(map[uint32]uint32)
	cost := make(map[uint32]uint32)
	targetLv := ship.StarLevel
	oldLv := ship.StarLevel

	loopCnt := 1
	if upgradeMax {
		loopCnt = 100
	}

	for i := 0; i < loopCnt; i++ {
		nextStarLevel := template2.GetRoleStarTemplate().GetShipStarLevel(shipId, targetLv+1)
		if nextStarLevel == nil {
			break
		}

		nextCost := make(map[uint32]uint32)
		for itemId, itemNum := range cost {
			nextCost[itemId] = itemNum
		}

		for i := 0; i < len(nextStarLevel.CostItems); i++ {
			itemId := nextStarLevel.CostItems[i].ItemId
			itemNum := nextStarLevel.CostItems[i].ItemNum
			nextCost[itemId] += itemNum
		}

		itemLack := false
		for itemId, itemNum := range nextCost {
			if !EnoughItem(p.GetUserId(), itemId, itemNum) {
				itemLack = true
				break
			}
		}

		if itemLack {
			break
		}

		for _, v := range nextStarLevel.Reward {
			reward[v.ItemId] += v.ItemNum
		}
		cost = nextCost
		targetLv += 1
	}

	if targetLv == ship.StarLevel {
		return msg.ErrCode_NO_ENOUGH_ITEM, 0, 0
	}

	// 扣除道具
	var notifyClientItems []uint32
	//tdaItems := make([]*tda.Item, 0, len(cost))
	for itemId, itemNum := range cost {
		CostItem(p.GetUserId(),
			itemId,
			itemNum,
			publicconst.ShipUpgradeStarCostItem,
			false)

		notifyClientItems = append(notifyClientItems, itemId)
		//tdaItems = append(tdaItems, &tda.Item{ItemId: strconv.Itoa(int(itemId)), ItemNum: itemNum})
	}

	// 获得升星奖励
	for itemId, itemNum := range reward {

		addItems := AddItem(p.GetUserId(),
			itemId,
			int32(itemNum),
			publicconst.ShipUpgradeStarAddItem,
			false)
		notifyClientItems = append(notifyClientItems, GetSimpleItemIds(addItems)...)
	}

	// 通知客户端
	//	event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(playerData, notifyClientItems, ListenNotifyClientItemEventEvent))
	// 通知客户端
	updateClientItemsChange(p.GetUserId(), notifyClientItems)

	for curLv := ship.StarLevel; curLv < targetLv; curLv++ {

		nextLv := curLv + 1
		nextStarLevel := template2.GetRoleStarTemplate().GetShipStarLevel(shipId, nextLv)
		if nextStarLevel == nil {
			continue
		}

		old_starlvl := curLv

		ship.StarLevel = nextLv

		// 发布机甲变化事件
		event.EventMgr.PublishEvent(event.NewShipChangeEvent(p, shipId, old_starlvl,
			ship.StarLevel, publicconst.Ship_Star_Level, ListenShipChangeEvent))

		// 添加升星属性
		addAllShipStarAttr(p, ship, int32(curLv))

		AddComboSkill(p, nextStarLevel.Data.PickSkillcombo)
	}

	p.SaveShips()

	msg2 := &msg.NotifyShipsChange{}
	msg2.Data = append(msg2.Data, ToProtocolShip(ship, p))
	p.SendNotify(msg2)

	// tda
	//tda.TdaKuluRankUp(p.ChannelId, p.TdaCommonAttr, shipId, shipConfig.Rarity, ship.StarLevel, tdaItems)

	return msg.ErrCode_SUCC, oldLv, targetLv
}

// GetShipPreview 获得机甲预览
func GetShipPreview(p *player.Player, shipId uint32) (msg.ErrCode, *model.Ship) {
	shipConfig := template2.GetShipTemplate().GetShip(shipId)
	if shipConfig == nil {
		return msg.ErrCode_SHIP_NOT_EXIST, nil
	}
	ship := model.NewShip(shipId)

	for i := 0; i < len(shipConfig.InitAttr); i++ {
		attr := model.NewAttr(shipConfig.InitAttr[i].Id, shipConfig.InitAttr[i].Value)
		ship.Attrs[attr.Id] = attr
	}

	// 等级属性
	addShipLevelAttr(ship, template2.GetShipLevelTempalte().GetLevelAttr(p.UserData.Level))
	calcShipAttr(ship)

	return msg.ErrCode_SUCC, ship
}

// ExchangeShip 请求机甲兑换
func ExchangeShip(p *player.Player, shipId uint32) msg.ErrCode {
	ship := getShip(p, shipId)
	if ship != nil {
		return msg.ErrCode_SHIP_EXIST
	}

	shipConfig := template2.GetShipTemplate().GetShip(shipId)
	if shipConfig == nil {
		return msg.ErrCode_SHIP_NOT_EXIST
	}

	for i := 0; i < len(shipConfig.ComposeItem); i++ {
		if !EnoughItem(p.GetUserId(),
			shipConfig.ComposeItem[i].ItemId, shipConfig.ComposeItem[i].ItemNum) {
			return msg.ErrCode_NO_ENOUGH_ITEM
		}
	}

	var notifyClientItems []uint32
	for i := 0; i < len(shipConfig.ComposeItem); i++ {
		CostItem(p.GetUserId(),
			shipConfig.ComposeItem[i].ItemId,
			shipConfig.ComposeItem[i].ItemNum,
			publicconst.ShipExchangeCostItem,
			false)
		notifyClientItems = append(notifyClientItems, shipConfig.ComposeItem[i].ItemId)
	}
	// 通知客户端
	//	event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(playerData, notifyClientItems, ListenNotifyClientItemEventEvent))
	// 通知客户端
	updateClientItemsChange(p.GetUserId(), notifyClientItems)
	event.EventMgr.PublishEvent(event.NewShipChangeEvent(p, shipId, 1,
		1, publicconst.Ship_Star_Level, ListenShipChangeEvent))

	AddShip(p, shipId, publicconst.ExchangeShip, true)
	UpdateTask(p, true, publicconst.TASK_COND_ACTIVE_SHIP, shipConfig.Rarity, 1) // 激活XX个XX品质库鲁
	return msg.ErrCode_SUCC
}

// addAllShipStarAttr 提升机甲等级属性
func addAllShipStarAttr(p *player.Player, ship *model.Ship, oldLevel int32) {
	//	s.addShipStarAttr(playerData, oldLevel, ship)
	calcShipAttr(ship)
	p.SaveShips()

	GlobalAttrChange(p, true)
}

// addAllShipLevelAttr 提升机甲等级属性
func addAllShipLevelAttr(p *player.Player, oldLevel, newLevel uint32) {
	// 通知客户端
	notifyMsg := &msg.NotifyShipsChange{}
	attrMap := template2.GetShipLevelTempalte().GetLevelRangeAttr(oldLevel, newLevel)
	for i := 0; i < len(p.UserData.Ships.Ships); i++ {
		addShipLevelAttr(p.UserData.Ships.Ships[i], attrMap)
		calcShipAttr(p.UserData.Ships.Ships[i])
		notifyMsg.Data = append(notifyMsg.Data, ToProtocolShip(p.UserData.Ships.Ships[i], p))
	}
	p.SaveShips()
	p.SendNotify(notifyMsg)
}

func getShip(p *player.Player, shipId uint32) *model.Ship {
	for i := 0; i < len(p.UserData.Ships.Ships); i++ {
		if p.UserData.Ships.Ships[i].Id == shipId {
			return p.UserData.Ships.Ships[i]
		}
	}
	return nil
}

// AddShip 添加机甲
func AddShip(p *player.Player, shipId uint32, source publicconst.ItemSource, notifyClient bool) *model.Ship {
	if getShip(p, shipId) != nil {
		log.Error("ship not nil", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("shipId", shipId))
		return nil
	}

	shipConfig := template2.GetShipTemplate().GetShip(shipId)
	if shipConfig == nil {
		log.Error("ship cfg nil", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("shipId", shipId))
		return nil
	}

	//attrMap := make(map[uint32]float32)
	//attrMap[1] = 1.0
	ship := model.NewShip(shipId)
	p.UserData.Ships.Ships = append(p.UserData.Ships.Ships, ship)

	// 初始化属性
	//if shipConfig != nil {
	for i := 0; i < len(shipConfig.InitAttr); i++ {
		attr := model.NewAttr(shipConfig.InitAttr[i].Id, shipConfig.InitAttr[i].Value)
		ship.Attrs[attr.Id] = attr
	}
	//}

	// 等级属性
	addShipLevelAttr(ship, template2.GetShipLevelTempalte().GetLevelAttr(p.UserData.Level))

	// 星级属性
	//s.addShipStarAttr(playerData, -1, ship)

	calcShipAttr(ship)

	//dao.ShipDao.AddShip(p.AccountInfo.AccountId, ship)

	// 计算终值
	var shipIds []uint32
	shipIds = append(shipIds, shipId)

	// 计算战力
	GlobalAttrChange(p, true)

	if source != publicconst.InitAddItem {
		if notifyClient {
			UpdateClientShipChange(p, shipIds)
		}
		event.EventMgr.PublishEvent(event.NewAddShipEvent(p, shipIds, notifyClient, ListenAddShipEvent))
		AddShipAppearances(p, shipId)
	} else {
		UpdateShipTreasure(p, shipId, false)
		UpdateShipPoker(p, shipId, false)
	}

	if shipStarConfig := template2.GetRoleStarTemplate().GetShipStarLevel(shipId, 0); shipStarConfig != nil {
		AddComboSkill(p, shipStarConfig.Data.PickSkillcombo)
	}

	ActivateAtlasByType(p, template2.HandBook_Type_Ship, []uint32{shipId})

	p.SaveShips()
	// tda
	tda.TdaKuluUnlock(p.ChannelId, p.TdaCommonAttr, shipId, shipConfig.Rarity)

	return ship
}

// addShipLevelAttr 添加机甲等级属性
func addShipLevelAttr(ship *model.Ship, attrMap map[uint32]*template2.AttrItem) []uint32 {
	var ret []uint32
	for id, attr := range attrMap {
		if data, ok := ship.Attrs[id]; ok {
			data.AddLevelValue(attr.Value)
		} else {
			data := model.NewAttr(id, 0)
			data.AddLevelValue(attr.Value)
		}
		ret = append(ret, id)
	}

	return ret
}

// calcShipAttr 计算机甲属性
func calcShipAttr(ship *model.Ship) {
	for id, data := range ship.Attrs {
		//	idConfig := template2.GetAttrListTemplate().GetAttr(id)
		starConfig := template2.GetRoleStarTemplate().GetShipStarLevel(ship.Id, ship.StarLevel)
		var finalValue float32
		if id == uint32(publicconst.Attack) || id == uint32(publicconst.Hp) ||
			id == uint32(publicconst.Defense) {
			finalValue = (data.InitValue + data.LevelValue + data.Add) * starConfig.Data.StarFactor
		} else {
			finalValue = data.InitValue + data.LevelValue + data.Add
		}
		data.SetFinalValue(finalValue)
	}
}

func getShipSupoortAttr(ship *model.Ship) []*model.Attr {
	shipConfig := template2.GetShipTemplate().GetShip(ship.Id)
	initAttr := template2.GetSupoortTemplate().GetSupportAttr(shipConfig.Rarity, ship.StarLevel)

	var ret []*model.Attr
	for i := 0; i < len(initAttr); i++ {
		if attrConfig := template2.GetAttrListTemplate().GetAttr(initAttr[i].Id); attrConfig != nil {
			data := model.NewAttr(initAttr[i].Id, initAttr[i].Value)
			data.CalcFinalValue()
			ret = append(ret, data)
		}
	}
	return ret
}

// UpdateClientShipChange 更新客户端机甲变化
func UpdateClientShipChange(p *player.Player, shipIds []uint32) {
	msg := &msg.NotifyShipsChange{}
	for k := 0; k < len(shipIds); k++ {
		if ship := getShip(p, shipIds[k]); ship != nil {
			msg.Data = append(msg.Data, ToProtocolShip(ship, p))
		}
	}
	p.SendNotify(msg)
}

// ToProtocolShip 转换成协议机甲
func ToProtocolShip(ship *model.Ship, p *player.Player) *msg.ShipInfo {
	ret := &msg.ShipInfo{ShipId: ship.Id, StarLevel: ship.StarLevel}
	var ids []uint32
	ids = append(ids, uint32(publicconst.Attack))
	ids = append(ids, uint32(publicconst.Hp))
	ids = append(ids, uint32(publicconst.Defense))
	if config := template2.GetRoleStarTemplate().GetShipStarLevel(ship.Id, ship.StarLevel); config != nil {
		for i := 0; i < len(config.StarAttr); i++ {
			ids = append(ids, config.StarAttr[i].Id)
		}
	}
	for id, data := range ship.Attrs {
		if tools.ListContain(ids, id) {
			ret.Attrs = append(ret.Attrs, &msg.Attr{
				Id:        data.Id,
				Value:     data.InitValue + data.LevelValue + data.Add,
				CalcValue: data.FinalValue,
			})
		}
	}
	ret.SupportAttr = ToProtocolAttrs(getShipSupoortAttr(ship))

	// 机甲皮肤 只返回已激活皮肤
	//coatMap := template2.GetCoatTemplate().GetModelCoatConfigs(int(ship.Id))
	//shipCfg := template2.GetShipTemplate().GetShip(ship.Id)
	//originalCoat := template2.GetCoatTemplate().GetCoatCfg(int(ship.Id), shipCfg.CoatModeOne)
	//for _, cfg := range coatMap {
	//	if item, ok := ship.CoatMap[cfg.CoatId]; !ok {
	//		if cfg.CoatId == originalCoat.CoatId { // 原皮
	//			item = model.NewCoatItem(cfg)
	//			ship.CoatMap[cfg.CoatId] = item
	//			p.SaveShips()
	//			ret.Coats = append(ret.Coats, item.ToProto())
	//		}
	//	} else {
	//		ret.Coats = append(ret.Coats, item.ToProto())
	//	}
	//}
	for _, item := range ship.CoatMap {
		ret.Coats = append(ret.Coats, item.ToProto())
	}
	return ret
}

// ToProtocolShips 转换成协议机甲
func ToProtocolShips(ships []*model.Ship, p *player.Player) []*msg.ShipInfo {
	var ret []*msg.ShipInfo
	for i := 0; i < len(ships); i++ {
		ret = append(ret, ToProtocolShip(ships[i], p))
	}
	return ret
}

// 激活皮肤
func ActiveCoat(p *player.Player, modelId, coatId int) (msg.ErrCode, *model.Ship) {
	cfg := template2.GetCoatTemplate().GetCoatCfg(modelId, coatId)
	if cfg == nil || len(cfg.CoatItem) == 0 { // 原皮不需要激活
		log.Error("active coat not exit", zap.Uint64("uid", p.GetUserId()), zap.Int("modelId", modelId), zap.Int("coatId", coatId))
		return msg.ERRCODE_COAT_NOT_EXIST, nil
	}

	//shipConfig := template2.GetShipTemplate().GetShip(uint32(modelId))
	//if shipConfig == nil {
	//	log.Error("active coat ship cfg not exit", zap.Uint64("uid", p.GetUserId()), zap.Int("modelId", modelId), zap.Int("coatId", coatId))
	//	return msg.ErrCode_SHIP_NOT_EXIST
	//}

	ship := getShip(p, uint32(modelId))
	if ship == nil {
		log.Error("active coat ship not exit", zap.Uint64("uid", p.GetUserId()), zap.Int("modelId", modelId), zap.Int("coatId", coatId))
		return msg.ErrCode_SHIP_NOT_EXIST, ship
	}

	if _, ok := ship.CoatMap[coatId]; ok {
		log.Error("active coat already active", zap.Uint64("uid", p.GetUserId()), zap.Int("modelId", modelId), zap.Int("coatId", coatId))
		return msg.ErrCode_COAT_ALREADY_ACTIVE, ship
	}

	for i := 0; i < len(cfg.CostItem); i++ {
		if !EnoughItem(p.GetUserId(),
			cfg.CostItem[i].ItemId, cfg.CostItem[i].ItemNum) {
			return msg.ErrCode_NO_ENOUGH_ITEM, ship
		}
	}

	var notifyClientItems []uint32
	for i := 0; i < len(cfg.CostItem); i++ {
		CostItem(p.GetUserId(),
			cfg.CostItem[i].ItemId,
			cfg.CostItem[i].ItemNum,
			publicconst.CoatActive,
			false)
		notifyClientItems = append(notifyClientItems, cfg.CostItem[i].ItemId)
	}

	ship.PutOnCoat(cfg, true)
	p.SaveShips()
	updateClientItemsChange(p.GetUserId(), notifyClientItems) // 通知客户端
	return msg.ErrCode_SUCC, ship
}

// 穿上皮肤
func PutOnCoat(p *player.Player, modelId, coatId int) (msg.ErrCode, *model.Ship) {
	cfg := template2.GetCoatTemplate().GetCoatCfg(modelId, coatId)
	if cfg == nil {
		log.Error("put on coat not exit", zap.Uint64("uid", p.GetUserId()), zap.Int("modelId", modelId), zap.Int("coatId", coatId))
		return msg.ERRCODE_COAT_NOT_EXIST, nil
	}

	ship := getShip(p, uint32(modelId))
	if ship == nil {
		log.Error("put on coat ship not exit", zap.Uint64("uid", p.GetUserId()), zap.Int("modelId", modelId), zap.Int("coatId", coatId))
		return msg.ErrCode_SHIP_NOT_EXIST, ship
	}

	if _, ok := ship.CoatMap[coatId]; !ok { // 穿上未激活皮肤
		log.Error("put on not active coat", zap.Uint64("uid", p.GetUserId()), zap.Int("modelId", modelId), zap.Int("coatId", coatId))
		return msg.ERRCODE_COAT_NOT_ACTIVE, ship
	}

	//for _, item := range ship.CoatMap {
	//	if item.Status == model.CoatOn {
	//		item.Status = model.CoatActive
	//	}
	//}
	//ship.CoatMap[coatId].Status = model.CoatOn
	ship.PutOnCoat(cfg, false)
	p.SaveShips()
	return msg.ErrCode_SUCC, ship
}
