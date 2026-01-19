package handle

import (
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"msg"

	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
)

// RequestLoadShopHandle 加载商城
func RequestLoadShopHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestLoadShop)
	log.Debug("load shop msg", zap.Uint64("uid", p.GetUserId()), zap.Uint32("packetId", packetId), zap.Any("req", req))
	//err, items := service.GetShopItem(p, req.ShopType)
	err, items := service.GetShopItems(p, req.ShopType)
	res := &msg.ResponseLoadShop{
		Result: err,
	}
	if err == msg.ErrCode_SUCC {
		res.Items = service.ToProtocolShopItems(items)
		res.ShopType = req.ShopType
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestGetShopItemHandle 获取指定的商品（同类型所有商品）
func RequestGetShopItemHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestGetShopItem)
	log.Debug("RequestGetShopItemHandle msg", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Uint32("packetId", packetId))
	err, items := service.GetShopItemByIds(p, req.Ids)
	res := &msg.ResponseGetShopItem{
		Result: err,
	}
	if err == msg.ErrCode_SUCC {
		res.Items = service.ToProtocolShopItems(items)
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestBuyShopItemHandle 购买商品
func RequestBuyShopItemHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestBuyShopItem)
	//err, rewardItems, nextItems, buyTimes := service.BuyShopItem(p, req.ShopItemId, req.BuyNum)
	err, rewardItems, nextItems, buyTimes := service.BuyShopItemWithCost(p, req.ShopItemId, req.BuyNum, req.CostItemId) // 购买指定商品
	res := &msg.ResponseBuyShopItem{
		Result:     err,
		ShopItemId: req.ShopItemId,
		BuyNum:     req.BuyNum,
	}
	if err == msg.ErrCode_SUCC {
		res.BuyTimes = buyTimes
		res.GetItems = service.ToProtocolSimpleItems(rewardItems)
		res.NextShopItem = service.ToProtocolShopItems(nextItems)
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestRefreshShopHandle 刷新商品
func RequestRefreshShopHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestRefreshShop)
	log.Debug("refresh shop msg", zap.Uint64("uid", p.GetUserId()), zap.Uint32("packetId", packetId), zap.Any("req", req))
	err, items := service.RefreshShop(p, req.ShopItemId)
	res := &msg.ResponseRefreshShop{
		Result: err,
		Data:   service.ToProtocolShopItems(items),
	}
	p.SendResponse(packetId, res, res.Result)
}

// 获取首充礼包数据
func RequestFirstChargePackageHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.FirstChargePackageResp{}
	firstChargePackageData := p.UserData.BaseInfo.FirstChargePackage
	if len(firstChargePackageData) > 0 {
		for _, v := range firstChargePackageData {
			data := &msg.FirstChargePackageData{}
			data.PackageId = v.Id
			data.State = v.State
			res.PackageData = append(res.PackageData, data)
		}
	}
	p.SendResponse(packetId, res, msg.ErrCode_SUCC)
}

// 领取首充礼包天数奖励
func RequestFirstChargePackageGetHandler(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.FirstChargePackageGetReq)
	resp := &msg.FirstChargePackageGetResp{}
	code, itemList := service.GetFirstChargePackageReward(p, req.ConfigId, req.Day)
	resp.Result = code
	resp.ItemList = itemList

	firstChargePackageData := p.UserData.BaseInfo.FirstChargePackage
	if len(firstChargePackageData) > 0 {
		for _, v := range firstChargePackageData {
			if v.Id == req.ConfigId {
				data := &msg.FirstChargePackageData{}
				data.PackageId = v.Id
				data.State = v.State
				resp.PackageData = data
				break
			}
		}
	}

	p.SendResponse(packetId, resp, resp.Result)
}
