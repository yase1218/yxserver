package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

// RequestLoadPokerHandle 加载扑克
func RequestLoadPokerHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponseLoadPoker{Result: msg.ErrCode_SUCC}
	res.CommonPokerData = p.UserData.Poker.CommData
	for _, v := range p.UserData.Poker.MissData {
		res.MissionPokerData = append(res.MissionPokerData, uint32(v))
	}
	res.ShipPokerData = service.ToProtocolShipPokers(p.UserData.Poker.ShipData)
	res.WeaponPokerData = service.ToProtocolWeaponPokers(p.UserData.Poker.WeaponData)
	p.SendResponse(packetId, res, res.Result)
}
