package service

import (
	"fmt"
	"gameserver/internal/config"
	"gameserver/internal/game/charge"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tools"
	"math/rand"
	"msg"
	"time"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func ChargeCallBack(playerData *player.Player, chargeId, num int) {
	chargeConfig := template.GetChargeTemplate().GetCharge(chargeId)
	if chargeConfig == nil && (chargeId != 2001 && chargeId != 2002 && chargeId != 2003) {
		log.Error("ChargeCallBack err", zap.Uint64("AccountId", playerData.GetUserId()), zap.Int("chargeId", chargeId), zap.Int("num", num))
		return
	}

	if chargeConfig == nil {
		// 首充礼包回调 占位
		cfg := template.GetOneChargeTemplate().GetTargetCfg(uint32(chargeId))
		if cfg != nil {
			firstChargePackageCB(playerData, cfg)
		}
		return
	}

	switch publicconst.ChargeType(chargeConfig.Type) {
	case publicconst.Charge_Diamond, publicconst.Charge_First:
		commonCharge(playerData, chargeConfig)
	case publicconst.Charge_MonthCard:
		monthcardCharge(playerData, chargeConfig)
	case publicconst.Charge_MainFund:
		mainFundCharge(playerData, chargeConfig)
	case publicconst.Charge_ActivePass:
		activePassCharge(playerData, chargeConfig)
	case publicconst.Charge_TaskPass:
		taskPassCharge(playerData, chargeConfig)
	case publicconst.Charge_MissionGift, publicconst.Charge_PerGift:
		itemConfig := template.GetShopTemplate().GetChargeShop(uint32(chargeConfig.Id))
		log.Debug("ChargeService ChargeCallBack, ", zap.Int("ChargeType", chargeConfig.Type),
			zap.Uint64("accountId", playerData.GetUserId()), zap.Int("chargeConfig id", chargeConfig.Id))
		BuyShopItemCallBack(playerData, itemConfig)
	case publicconst.Charge_ActivitySign:
		signCharge(playerData, chargeConfig)
	}
}

func commonCharge(p *player.Player, config *template.JCharge) {
	exist := false
	isFirst := false
	for i := 0; i < len(p.UserData.BaseInfo.Charge); i++ {
		if p.UserData.BaseInfo.Charge[i].Id == config.Id {
			if p.UserData.BaseInfo.Charge[i].Value == 0 {
				p.UserData.BaseInfo.Charge[i].Value = 1
				isFirst = true
			}
			exist = true
			break
		}
	}
	if !exist {
		isFirst = true
		p.UserData.BaseInfo.Charge = append(p.UserData.BaseInfo.Charge, model.NewChargeInfo(config.Id, 1))
	}

	totalItems := make(map[uint32]uint32)
	for i := 0; i < len(config.RewardItems); i++ {
		if _, ok := totalItems[config.RewardItems[i].ItemId]; ok {
			totalItems[config.RewardItems[i].ItemId] += config.RewardItems[i].ItemNum
		} else {
			totalItems[config.RewardItems[i].ItemId] = config.RewardItems[i].ItemNum
		}
	}

	extraItems := config.ExtraItems
	if isFirst {
		extraItems = config.FirstChargeItems
	}

	for i := 0; i < len(extraItems); i++ {
		if _, ok := totalItems[extraItems[i].ItemId]; ok {
			totalItems[extraItems[i].ItemId] += extraItems[i].ItemNum
		} else {
			totalItems[extraItems[i].ItemId] = extraItems[i].ItemNum
		}
	}

	// 领取任务奖励
	var notifyItems []uint32
	var finalItems []*template.SimpleItem
	for id, num := range totalItems {
		addItems := AddItem(p.GetUserId(),
			id,
			int32(num),
			publicconst.ChargeAddItem,
			false)
		finalItems = append(finalItems, ToTemplateItem(addItems)...)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}

	p.SaveBaseInfo()
	updateClientItemsChange(p.GetUserId(), notifyItems)

	res := &msg.NotifyRecharge{}
	res.RechargeId = uint32(config.Id)

	res.GetItems = TemplateItemToProtocolItems(finalItems)
	p.SendNotify(res)
}

// monthcardCharge 月卡充值
func monthcardCharge(p *player.Player, config *template.JCharge) {
	pos := -1
	for i := 0; i < len(p.UserData.BaseInfo.MonthCard); i++ {
		if p.UserData.BaseInfo.MonthCard[i].Id == config.Id {
			pos = i
			break
		}
	}

	monthCardConfig := template.GetMonthCardTemplate().GetMonthCard(config.Id)
	if monthCardConfig.Data.Tpye == 2 {
		AddPlayMethodTimes(p, int(msg.BattleType_Battle_Coin), int(monthCardConfig.Data.QuestTimes))
		AddPlayMethodTimes(p, int(msg.BattleType_Battle_Equip), int(monthCardConfig.Data.QuestTimes))
		AddPlayMethodTimes(p, int(msg.BattleType_Battle_Weapon), int(monthCardConfig.Data.QuestTimes))
	}

	var info *model.MonthcardInfo
	// 月卡续费
	if pos == -1 {
		endTime := int(tools.GetDateStart(time.Now()).AddDate(0, 0, int(monthCardConfig.Data.Time)).Unix())
		info = model.NewMonthcardInfo(config.Id, endTime)
		p.UserData.BaseInfo.MonthCard = append(p.UserData.BaseInfo.MonthCard, info)
	} else {
		info = p.UserData.BaseInfo.MonthCard[pos]
		curTime := int(tools.GetCurTime())
		if curTime >= p.UserData.BaseInfo.MonthCard[pos].EndTime {
			p.UserData.BaseInfo.MonthCard[pos].EndTime = int(tools.GetDateStart(time.Now()).AddDate(0, 0, int(monthCardConfig.Data.Time)).Unix())
			p.UserData.BaseInfo.MonthCard[pos].NextGetRewardTime = 0
		} else {
			p.UserData.BaseInfo.MonthCard[pos].EndTime += monthCardConfig.Data.Time * publicconst.DAY_SECONDS
		}
	}
	p.SaveBaseInfo()

	// 发放月卡奖励
	var notifyItems []uint32
	var finalItems []*template.SimpleItem
	for i := 0; i < len(config.RewardItems); i++ {
		addItems := AddItem(p.GetUserId(),
			config.RewardItems[i].ItemId,
			int32(config.RewardItems[i].ItemNum),
			publicconst.MonthCardAddItem,
			false)
		finalItems = append(finalItems, ToTemplateItem(addItems)...)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	res := &msg.NotifyMonthCardCharge{}
	res.CardInfo = &msg.MonthcardInfo{
		MonthCardId:       uint32(info.Id),
		EndTime:           uint32(info.EndTime),
		NextGetRewardTime: uint32(info.NextGetRewardTime),
	}

	res.GetItems = TemplateItemToProtocolItems(finalItems)
	p.SendNotify(res)
}

// mainFundCharge 主线基金充值
func mainFundCharge(p *player.Player, config *template.JCharge) {
	var info *model.MainFundInfo
	for i := 0; i < len(p.UserData.BaseInfo.MainFund); i++ {
		if p.UserData.BaseInfo.MainFund[i].Id == config.Id {
			info = p.UserData.BaseInfo.MainFund[i]
			break
		}
	}
	if info == nil {
		info = model.NewMainFundInfo(config.Id)
		p.UserData.BaseInfo.MainFund = append(p.UserData.BaseInfo.MainFund, info)
	}
	info.BuyFlag = 1

	itemMap, _ := getMainFundReward(p, config.Id)

	// 领取任务奖励
	var notifyItems []uint32
	var finalItems []*template.SimpleItem
	for id, num := range itemMap {
		addItems := AddItem(p.GetUserId(), uint32(id), int32(num), publicconst.BuyMainFundAddItem, false)
		finalItems = append(finalItems, ToTemplateItem(addItems)...)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}

	p.SaveBaseInfo()

	if len(notifyItems) > 0 {
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	res := &msg.NotifyChargeFund{}
	res.GetItems = TemplateItemToProtocolItems(finalItems)
	res.FundInfo = &msg.MainFundInfo{
		FundId:      uint32(info.Id),
		RewardMaxId: uint32(info.FreeId),
		BuyFlag:     uint32(info.BuyFlag),
	}
	p.SendNotify(res)
}

// activePassCharge 活跃战令充值
func activePassCharge(playerData *player.Player, config *template.JCharge) {
	buyPass(playerData, config.Id)
}

func signCharge(playerData *player.Player, config *template.JCharge) {
	buySign(playerData, config.Id)
}

// taskPassCharge 购买任务战令
func taskPassCharge(playerData *player.Player, config *template.JCharge) {
	buyTaskPass(playerData, config.Id)
}

func firstChargePackageCB(p *player.Player, cfg *template.JOneCharge) {
	chargeInfo := p.UserData.BaseInfo.FirstChargePackage
	isFind := false
	for i := 0; i < len(chargeInfo); i++ {
		if chargeInfo[i].Id == uint32(cfg.Id) {
			isFind = true
			break
		}
	}

	if isFind {
		log.Info("player has been first charge package", zap.String("account", p.UserData.AccountId), zap.Int("charge id", cfg.Id))
		return
	}

	stateList := make([]uint32, len(cfg.ParsedItemRewards))
	for i := 0; i < len(cfg.ParsedItemRewards); i++ {
		if i == 0 {
			stateList[i] = Obtained
		} else {
			stateList[i] = Forbidden
		}
	}

	chargeData := &model.FirstChargePackageData{
		Id:         uint32(cfg.Id),
		LoginCount: 1,
		State:      stateList,
	}

	chargeInfo = append(chargeInfo, chargeData)
	p.UserData.BaseInfo.FirstChargePackage = chargeInfo
	p.SaveBaseInfo()
}

// 创建订单
func CreateOrder(p *player.Player, req *msg.RequestCreateOrder) (msg.ErrCode, *model.OrderInfo) {
	//gameOrderNo := uuid.New().String()
	gameOrderNo := fmt.Sprintf("%d_%s_%s_%d_%d", config.Conf.ServerId, p.UserData.AccountId, time.Now().Format("20060102150405.000000"), req.ChargeId, rand.Int31n(999999))
	cfg := template.GetChargeTemplate().GetCharge(int(req.ChargeId))
	if cfg == nil {
		log.Error("CreateOrder GetCharge err", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req))
		return msg.ERRCODE_SYSTEM_ERROR, nil
	}
	now := time.Now()
	order := &model.OrderInfo{
		Id:        primitive.NewObjectID(),
		ChargeID:  int(req.ChargeId),
		OrderId:   gameOrderNo,
		UserId:    p.GetUserId(),
		ThirdNo:   "",
		ChannelNo: p.SdkChannelNo,
		Currency:  "CNY",                  // 第一阶段写死“人民币”
		Money:     int(cfg.CostRMB * 100), // 配置表里单位是元
		ProductId: model.GetProductId(cfg),
		AccountId: p.UserData.AccountId,
		CreateAt:  now,
		UpdateAt:  now,
	}
	//b, err := json.Marshal(order)
	//if err != nil {
	//	log.Error("CreateOrder Marshal order err", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Error(err))
	//	return msg.ERRCODE_SYSTEM_ERROR, nil
	//}
	//key := charge.GetOrderRedisKey(p.UserData.AccountId, config.Conf.ServerId)
	//rc := rdb_single.Get()
	//err = rc.HSet(rdb_single.GetCtx(), key, gameOrderNo, string(b)).Err()
	//if err != nil {
	//	log.Error("CreateOrder redis set order err", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Error(err))
	//	return msg.ERRCODE_SYSTEM_ERROR, nil
	//}
	//err = rc.HExpire(rdb_single.GetCtx(), key, model.RedisOrderTime, gameOrderNo).Err()
	//if err != nil {
	//	log.Error("CreateOrder redis expire order err", zap.Uint64("uid", p.GetUserId()), zap.String("key", key), zap.String("gameOrderNo", gameOrderNo), zap.String("b", string(b)), zap.Error(err))
	//	return msg.ERRCODE_SYSTEM_ERROR, nil
	//}
	//err = rc.SetEx(rdb_single.GetCtx(), key, string(b), model.RedisOrderTime).Err()

	if err := charge.CreateOrder(order); err != nil {
		log.Error("CreateOrder mongo err", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Error(err))
		return msg.ERRCODE_SYSTEM_ERROR, nil
	}
	charge.GetOrderManager().AddOrder(order)
	log.Debug("CreateOrder success", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req))
	return msg.ErrCode_SUCC, order
}

// 登录发货已验证订单
func UserOrdersShipment(p *player.Player) {
	orders := charge.GetOrderManager().GetOrderByUid(p.GetUserId())
	for i := 0; i < len(orders); i++ {
		if orders[i].Status == model.OrderVerified {
			ChargeCallBack(p, orders[i].ChargeID, 1)
			orders[i].Status = model.OrderShipment
			orders[i].UpdateAt = time.Now()
			orders[i].SaveStatus()
			ClearOutPut(p, orders[i].ChargeID, false) // 推送过期时间清零
		}
	}
}
