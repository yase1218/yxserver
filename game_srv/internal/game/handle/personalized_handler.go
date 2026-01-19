package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"msg"
)

// 推送列表
func RequestPersonalizedActHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestPersonalized)
	log.Debug("RequestPersonalizedActHandle", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Uint32("packetId", packetId))
	err, data := service.GetPersonalizedItems(p, int(req.Type), int(req.ItemId))
	res := &msg.ResponsePersonalized{
		Result: err,
		Items:  data,
	}
	p.SendResponse(packetId, res, res.Result)
}

// 未过时推送
func RequestUnOutTimePersonalizedHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestUnOutTimePersonalized)
	log.Debug("RequestUnOutTimePersonalizedHandle", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Uint32("packetId", packetId))
	data := service.GetUnOutTimePersonalizedItems(p)
	res := &msg.ResponseUnOutTimePersonalized{
		Result: msg.ErrCode_SUCC,
		Items:  data,
	}
	p.SendResponse(packetId, res, res.Result)
}
