package handle

import (
	"msg"

	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
)

// RequestGMCommandHandle 处理GM命令
func RequestGMCommandHandle(packetId uint32, args interface{}, p *player.Player) {
	request := args.(*msg.RequestGMCommand)
	err := service.ProcessCommand(request.CommandId, request.Content, p)
	res := &msg.ResponseGMCommand{
		Result: err,
	}
	p.SendResponse(packetId, res, res.Result)
}
