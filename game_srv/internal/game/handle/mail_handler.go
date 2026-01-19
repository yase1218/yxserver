package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

func RequestLoadMailHandle(pid uint32, args interface{}, p *player.Player) {
	service.LoadUserMail(pid, p)
}

func RequestReadMailHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestReadMail)
	service.ReadMail(pid, p, req)
}

func RequestGetMailRewardHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestGetMailReward)
	service.GetMailReward(pid, p, req)
}

func RequestBatchGetMailRewardHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestBatchGetMailReward)
	service.BatchGetMailReward(pid, p, req)
}

func RequestDelMailHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestDelMail)
	service.DelReadMail(pid, p, req)
}
