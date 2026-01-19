package handle

import (
	"gameserver/internal/game/builder"
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func Atlas(packetId uint32, args interface{}, p *player.Player) {
	res := new(msg.AtlasAck)
	res.Atlases = builder.BuildAtlas(p.UserData.Atlas)
	p.SendNotify(res)
}

func AtlasReward(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.AtlasRewardReq)
	res := &msg.AtlasRewardAck{
		Id: req.Id,
	}
	if err := service.AtlasReward(p, req, res); err != nil {
		log.Error("atlas err", zap.Uint64("uid", p.GetUserId()), zap.Uint32("id", req.Id), zap.Error(err))
		return
	}
	p.SendNotify(res)
}
