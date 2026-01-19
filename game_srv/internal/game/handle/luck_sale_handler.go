package handle

import (
	"gameserver/internal/game/builder"
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

func LuckSale(packetId uint32, args interface{}, p *player.Player) {
	res := builder.BuildLuckSaleAck(p.UserData)
	p.SendResponse(packetId, res, msg.ErrCode_SUCC)
}

func LuckSaleExtract(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.LuckSaleExtractReq)
	res := &msg.LuckSaleExtractAck{
		Err: msg.ErrCode_SUCC,
	}
	service.LuckSaleExtract(p, req, res)
	p.SendResponse(packetId, res, res.Err)
}

func LuckSaleTaskReward(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.LuckSaleTaskRewardReq)
	res := &msg.LuckSaleTaskRewardAck{
		Err: msg.ErrCode_SUCC,
	}
	service.LuckSaleTaskReward(p, req, res)
	p.SendResponse(packetId, res, res.Err)
}
