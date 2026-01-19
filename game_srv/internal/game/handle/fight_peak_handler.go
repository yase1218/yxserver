package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

// import (
// 	"github.com/v587-zyf/gc/log"
// 	"go.uber.org/zap"
// 	"kernel/tools"
// 	"msg"
// 	"server/internal/game/common"
// 	"server/internal/game/service"
// )

func PeakFight(pid uint32, args interface{}, p *player.Player) {
	//req := args.(*msg.PeakFightReq)
	service.LoadPeakFight(pid, p)
}

func PeakFightMatch(pid uint32, args interface{}, p *player.Player) {
	service.PeakFightMatch(pid, p)
}

func PeakFightCancelMatch(pid uint32, args interface{}, p *player.Player) {
	service.PeakFightCancelMatch(pid, p)
}

func PeakFightRank(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.PeakFightRankReq)
	service.PeakFightGetRank(pid, p, req)
}
