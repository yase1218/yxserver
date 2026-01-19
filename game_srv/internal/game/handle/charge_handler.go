package handle

import (
	"gameserver/internal/config"
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"msg"
)

// 创建订单
func RequestCreateOrderHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestCreateOrder)
	log.Debug("RequestCreateOrderHandle", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Uint32("packetId", packetId))
	errCode, order := service.CreateOrder(p, req)
	res := &msg.ResponseCreateOrder{
		Result: errCode,
	}
	if order != nil {
		res.ChargeId = req.ChargeId
		res.OrderId = order.OrderId
		res.Money = int32(order.Money)
		res.ProductId = order.ProductId
		res.ZoneId = int32(config.Conf.ServerId)
		res.Currency = order.Currency
		res.NotifyUri = config.Conf.Leiting.NotifyUrl
	}
	log.Debug("RequestCreateOrderHandle res", zap.Uint64("uid", p.GetUserId()), zap.Any("res", res))
	p.SendResponse(packetId, res, res.Result)
}
