package condition

import "gameserver/internal/game/player"

func initPet(p *player.Player, args []uint32) ([]uint32, bool) {
	if p.GetPet(args[0]) != nil {
		return args, true
	}

	return nil, false
}
