package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func Desert(packetId uint32, args interface{}, p *player.Player) {
	ack := new(msg.DesertAck)
	ack.KillTimes = p.UserData.Desert.KillTimes
	ack.RewardTimes = p.UserData.Desert.RewardTimes
	p.SendResponse(packetId, ack, msg.ErrCode_SUCC)
}

func DesertLampReward(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.DesertLampRewardReq)
	ack := new(msg.DesertLampRewardAck)
	if err := service.DesertLampReward(p, req, ack); err != nil {
		log.Error("desert lamp reward err", zap.String("err", err.Error()))
	}
	ack.RewardTimes = req.GetRewardTimes()
	p.SendResponse(packetId, ack, msg.ErrCode_SUCC)
}
