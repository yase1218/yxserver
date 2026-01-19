package handle

import "gameserver/internal/game/player"

// RequestEnterPvPHandle 进入pvp
func RequestEnterPvPHandle(packetId uint32, args interface{}, p *player.Player) {
	//retMsg := &msg.ResponseEnterPvP{
	//	Result: service.ServMgr.GetPvPService().EnterPvP(playerData),
	//}
	//if retMsg.Result == msg.ErrCode_SUCC {
	//	retMsg.PvpInfo = service.ToProtocolPvPInfo(playerData.AccountInfo.PvPInfo)
	//}
	//util.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
}

// RequestStartMatchHandle 开始匹配
func RequestStartMatchHandle(packetId uint32, args interface{}, p *player.Player) {
	//retMsg := &msg.ResponseStartMatch{
	//	Result: service.ServMgr.GetPvPService().StartMatch(playerData),
	//}
	//util.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
}

// RequestLoadingHandle loading开始
func RequestLoadingHandle(packetId uint32, args interface{}, p *player.Player) {
	//retMsg := &msg.ResponseLoading{
	//	Result: service.ServMgr.GetPvPService().BattleLoading(playerData),
	//}
	//util.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
}

// RequestPvPBattleEndHandle
func RequestPvPBattleEndHandle(packetId uint32, args interface{}, p *player.Player) {
	//req := args.(*msg.RequestPvPBattleEnd)
	//retMsg := msg.ResponsePvPBattleEnd{
	//	Result: service.ServMgr.GetPvPService().BattleEnd(playerData, req.Data),
	//}
	//util.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
}
