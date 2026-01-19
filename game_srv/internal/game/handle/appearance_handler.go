package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

// RequestLoadAppearanceHandle 加载外观
func RequestLoadAppearanceHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponseLoadAppearance{
		Result: msg.ErrCode_SUCC,
	}
	res.Data = service.ToProtocolAppearances(p.UserData.Appearance.Appearances)
	p.SendResponse(packetId, res, res.Result)
}

// RequestUseAppearanceHandle 使用外观
func RequestUseAppearanceHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestUseAppearance)
	res := &msg.ResponseUseAppearance{}
	res.Result = service.UseAppearance(p, req.Id)
	if res.Result == msg.ErrCode_SUCC {
		res.Id = req.Id
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestActiveAppearanceHandle 激活外观
func RequestActiveAppearanceHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestActiveAppearance)
	res := &msg.ResponseActiveAppearance{}
	res.Result = service.ActiveAppearance(p, req.Id)
	if res.Result == msg.ErrCode_SUCC {
		res.Id = req.Id
	}
	p.SendResponse(packetId, res, res.Result)
}
