package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"gameserver/internal/publicconst"
	"msg"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

// RequestLoadItemsHandle 请求加载道具
func RequestLoadItemsHandle(packetId uint32, args interface{}, p *player.Player) {
	//err, items := service.LoadItems(playerData)
	log.Debug("load items msg", zap.Uint64("uid", p.GetUserId()), zap.Uint32("packetId", packetId))
	res := &msg.ResponseLoadItems{
		Result: msg.ErrCode_SUCC,
	}

	if len(p.UserData.Items.Items) > 0 {
		res.Data = append(res.Data, service.ToProtocolItems(p.UserData.Items.Items)...)
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestUseItemHandle 使用道具
func RequestUseItemHandle(packetId uint32, args interface{}, p *player.Player) {
	log.Debug("use item msg", zap.Uint64("uid", p.GetUserId()), zap.Uint32("packetId", packetId))
	req := args.(*msg.RequestUseItem)
	err, reward := service.UseItem(p, req.ItemId, req.ItemNum, publicconst.UseItem, req.SelectItems, true)
	res := &msg.ResponseUseItem{
		Result:     err,
		ItemId:     req.ItemId,
		RewardItem: service.ToProtocolSimpleItems(reward),
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestComposeItemHandle 道具合成
func RequestComposeItemHandle(packetId uint32, args interface{}, p *player.Player) {
	log.Debug("compose item msg", zap.Uint64("uid", p.GetUserId()), zap.Uint32("packetId", packetId))
	req := args.(*msg.RequestComposeItem)
	err, reward := service.ComposeItem(p, req.ItemId, req.ComposeItemId, publicconst.ComposeItem)
	res := &msg.ResponseComposeItem{
		Result:     err,
		RewardItem: service.ToProtocolSimpleItems(reward),
	}
	//log.Debug("RequestComposeItemHandle res", zap.Any("res", res))
	p.SendResponse(packetId, res, res.Result)
}
