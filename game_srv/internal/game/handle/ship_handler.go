package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

// RequestLoadShipsHandle 加载机甲
func RequestLoadShipsHandle(packetId uint32, args interface{}, p *player.Player) {
	//err, ships := service.ServMgr.GetShipService().LoadShips(p)
	res := &msg.ResponseLoadShips{
		Result: msg.ErrCode_SUCC,
	}
	ships := p.UserData.Ships.Ships

	if len(ships) > 0 {
		res.Data = append(res.Data, service.ToProtocolShips(ships, p)...)
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestShipUpgradeHandle 升级机甲
func RequestShipUpgradeHandle(packetId uint32, args interface{}, p *player.Player) {
	//req := args.(*msg.RequestShipUpgrade)
	//err, level := service.ServMgr.GetShipService().UpgradeShip(playerData, req.ShipId)
	//retMsg := &msg.ResponseShipUpgrade{
	//	Result: err,
	//	ShipId: req.ShipId,
	//	Level:  level,
	//}
	//tools.SendMsg(playerData.PlayerAgent, retMsg, msg.MsgId_ID_ResponseShipUpgrade, packetId, retMsg.Result)
}

// RequestShipStarUpgradeHandle 升星机甲
func RequestShipStarUpgradeHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestShipStarUpgrade)
	err, oldLevel, newLevel := service.UpgradeStarShip(p, req.ShipId, req.UpgradeMax)
	res := &msg.ResponseShipStarUpgrade{
		Result:       err,
		ShipId:       req.ShipId,
		StarLevel:    newLevel,
		OldStarLevel: oldLevel,
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestExchangeShipHandle 请求机甲兑换
func RequestExchangeShipHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestExchangeShip)
	err := service.ExchangeShip(p, req.ShipId)
	res := &msg.ResponseExchangeShip{
		Result: err,
		Ship: &msg.ShipInfo{
			ShipId: req.ShipId,
		},
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestShipPreviewHandle 请求机甲预览
func RequestShipPreviewHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestShipPreview)
	res := &msg.ResponseShipPreview{}
	result, data := service.GetShipPreview(p, req.ShipId)
	res.Result = result
	if result == msg.ErrCode_SUCC {
		res.Data = service.ToProtocolShip(data, p)
	}
	p.SendResponse(packetId, res, res.Result)
}

// 皮肤激活请求
func RequestActiveCoatHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestActiveCoat)
	err, ship := service.ActiveCoat(p, int(req.ModelId), int(req.CoatId))
	if err == msg.ErrCode_SUCC {
		notify := &msg.NotifyShipsChange{}
		notify.Data = append(notify.Data, service.ToProtocolShip(ship, p))
		p.SendNotify(notify)
	}
	res := &msg.ResponseActiveCoat{
		Result:  err,
		ModelId: req.ModelId,
		CoatId:  req.CoatId,
	}
	p.SendResponse(packetId, res, res.Result)
}

// 皮肤穿上请求
func RequestPutOnCoatHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestPutOnCoat)
	err, ship := service.PutOnCoat(p, int(req.ModelId), int(req.CoatId))
	if err == msg.ErrCode_SUCC {
		notify := &msg.NotifyShipsChange{}
		notify.Data = append(notify.Data, service.ToProtocolShip(ship, p))
		p.SendNotify(notify)
	}
	res := &msg.ResponsePutOnCoat{
		Result:  err,
		ModelId: req.ModelId,
		CoatId:  req.CoatId,
	}
	p.SendResponse(packetId, res, res.Result)
}
