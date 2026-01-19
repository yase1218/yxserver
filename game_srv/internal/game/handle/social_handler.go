package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

// RequestPlayerDetailHandle 玩家详细
func RequestPlayerDetailHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestPlayerDetailInfo)
	err, data := service.GetPlayerDetailInfo(uint64(req.AccountId))
	res := &msg.ResponsePlayerDetailInfo{
		Result: err,
	}
	if err == msg.ErrCode_SUCC {
		res.Info = service.ToProtocolPlayerDetail(data, p)
	}
	p.SendResponse(packetId, res, res.Result)
}
