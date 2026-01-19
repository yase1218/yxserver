package condition

import (
	"gameserver/internal/game/common"
	"gameserver/internal/game/player"
	"kernel/tools"
	"time"
)

func initOpenServerDays(p *player.Player, args []uint32) ([]uint32, bool) {
	if days := common.GetOpenServerDays(); days >= args[0] {
		return []uint32{common.GetOpenServerDays()}, true
	}

	return nil, false
}

func initAccountDays(p *player.Player, args []uint32) ([]uint32, bool) {
	days := tools.GetDiffDay(time.Unix(int64(p.UserData.BaseInfo.CreateTime), 0), time.Now()) + 1
	if days >= args[0] {
		return []uint32{days}, true
	}

	return nil, false
}
