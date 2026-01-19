package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"kernel/tools"
	"msg"
)

// RequestPlayMethodDataHandle 请求玩法数据
func RequestPlayMethodDataHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponsePlayMethodData{
		Data:   service.ToProtocolPlayMethodDatas(p.UserData.PlayMethod.Data, p.UserData.Fight.Faction),
		Result: msg.ErrCode_SUCC,
	}

	if tools.GetCurTime() >= p.UserData.PlayMethod.NextRefreshTime {
		res.FirstLogin = 1
	}
	res.Data = service.ToProtocolPlayMethodDatas(p.UserData.PlayMethod.Data, p.UserData.Fight.Faction)

	p.SendResponse(packetId, res, res.Result)
}

// RequestPlayMethodStartBattleHandle 开始战斗
func RequestPlayMethodStartBattleHandle(packetId uint32, args interface{}, playerData *player.Player) {
	// req := args.(*msg.RequestPlayMethodStartBattle)
	// retMsg := &msg.ResponsePlayMethodStartBattle{
	// 	Result:    msg.ErrCode_SUCC,
	// 	MissionId: req.MissionId,
	// }
	// if err := service.ServMgr.GetCommonService().FunctionOpen(playerData, publicconst.PlayMethod); err != msg.ErrCode_SUCC {
	// 	retMsg.Result = err
	// 	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
	// 	return
	// }

	// retMsg.Result = service.ServMgr.GetPlayMethodService().StartBattle(playerData, int(req.MissionId))
	// if retMsg.Result == msg.ErrCode_SUCC {
	// 	retMsg.Attrs = service.ToProtocolAttrs2(playerData.AccountInfo.Attrs)
	// 	retMsg.ShipId = playerData.AccountInfo.ShipId
	// 	retMsg.SupportId = playerData.AccountInfo.SupportId
	// 	retMsg.ComboSkills = playerData.AccountInfo.ComboSkill
	// }
	// tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
}

// RequestPlayMethodEndBattleHandle 结束战斗
func RequestPlayMethodEndBattleHandle(packetId uint32, args interface{}, playerData *player.Player) {
	// req := args.(*msg.RequestPlayMethodEndBattle)
	// retMsg := &msg.ResponsePlayMethodEndBattle{
	// 	Result: msg.ErrCode_SUCC,
	// }
	// if playerData.AccountInfo.MissData != nil {
	// 	retMsg.MissionId = uint32(playerData.AccountInfo.MissData.MissionId)
	// }

	// if err := service.ServMgr.GetCommonService().FunctionOpen(playerData, publicconst.PlayMethod); err != msg.ErrCode_SUCC {
	// 	retMsg.Result = err
	// 	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
	// 	return
	// }

	// err, isPass, items := service.ServMgr.GetPlayMethodService().EndBattle(playerData, req.Result)
	// retMsg.Result = err
	// if isPass {
	// 	retMsg.IsPass = 1
	// }
	// if len(items) > 0 {
	// 	retMsg.GetItems = append(retMsg.GetItems, service.ToProtocolSimpleItems(items)...)
	// }

	// tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
}

// RequestPlayMethodSwapHandle 扫荡
func RequestPlayMethodSwapHandle(packetId uint32, args interface{}, playerData *player.Player) {
	// req := args.(*msg.RequestPlayMethodSwap)
	// retMsg := &msg.ResponsePlayMethodSwap{
	// 	Result: msg.ErrCode_SUCC,
	// }
	// if err := service.ServMgr.GetCommonService().FunctionOpen(playerData, publicconst.PlayMethod); err != msg.ErrCode_SUCC {
	// 	retMsg.Result = err
	// 	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
	// 	return
	// }

	// result, items := service.ServMgr.GetPlayMethodService().SwapBattle(playerData, int(req.BtType))
	// retMsg.Result = result
	// if result == msg.ErrCode_SUCC {
	// 	retMsg.GetItems = service.ToProtocolSimpleItems(items)
	// }
	// tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
}

// RequestPlayMethodUpdateWeaponHandle 保存武器
func RequestPlayMethodUpdateWeaponHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestPlayMethodUpdateWeapon)
	res := &msg.ResponsePlayMethodUpdateWeapon{
		Result: msg.ErrCode_SUCC,
	}

	res.Result = service.UpdateWeapon(p, int(req.BtType), req.WeaponIds)
	res.BtType = req.BtType
	res.WeaponIds = req.WeaponIds
	p.SendResponse(packetId, res, res.Result)
}
