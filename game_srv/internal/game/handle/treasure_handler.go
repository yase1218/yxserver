package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

// RequestLoadRareTreasureHandle 加载秘宝
func RequestLoadRareTreasureHandle(packetId uint32, args interface{}, p *player.Player) {
	retMsg := &msg.ResponseLoadRareTreasure{Result: msg.ErrCode_SUCC}
	retMsg.CommonTreasure = p.UserData.Treasure.CommData
	retMsg.MissionTreasure = p.UserData.Treasure.MissData
	retMsg.Srt = service.ToProtocolShipTreasures(p.UserData.Treasure.ShipData)
	retMsg.Wrt = service.ToProtocolWeaponTreasures(p.UserData.Treasure.WeaponData)

	p.SendResponse(packetId, retMsg, retMsg.Result)
}
