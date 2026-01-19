package condition

import (
	"gameserver/internal/game/player"
)

func initContractRand(p *player.Player, args []uint32) ([]uint32, bool) {
	if p.UserData.Contract.StageEventId == args[0] &&
		p.UserData.Contract.FinishCount >= args[1] {
		return []uint32{p.UserData.Contract.FinishCount}, true
	}
	return nil, false
}

func initContractKillMonster(p *player.Player, args []uint32) ([]uint32, bool) {
	if p.UserData.Contract.FinishCount >= args[0] {
		return []uint32{p.UserData.Contract.FinishCount}, true
	}
	return nil, false
}
