package condition

import "gameserver/internal/game/player"

func checkArenaPkCnt(p *player.Player, args []uint32) ([]uint32, bool) {
	if p.UserData.Arena.TotalPkCnt >= args[0] {
		return []uint32{p.UserData.Arena.TotalPkCnt}, true
	}
	return nil, false
}
