package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
	"time"

	"github.com/v587-zyf/gc/utils"
)

// RequestGetWeekPassActHandle 请求获取周常活动数据
func RequestGetWeekPassActHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestGetWeekPassAct)
	err, rank, data := service.EnterBlackBoss(p, req.Type)
	res := &msg.ResponseGetWeekPassAct{
		Result: err,
	}
	if err == msg.ErrCode_SUCC {
		res.Rank = rank
		res.Data = data
		res.SettleTime = uint32(utils.GetWeekdayZeroTime(time.Friday, false).Unix())
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestSetFactionAndWeaponHandle 请求设置流派和副武器
func RequestSetFactionAndWeaponHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.SetWeekPassDeputyWeaponAndFactionReq)
	res := &msg.SetWeekPassDeputyWeaponAndFactionResp{}
	code := service.UpdateWeapon(p, int(req.BtType), req.WeaponIds)
	res.BtType = req.BtType
	res.Code = code
	res.Type = req.Type
	if code != msg.ErrCode_SUCC {
		p.SendResponse(packetId, res, res.Code)
		return
	}
	res.WeaponIds = req.WeaponIds

	p.UserData.Fight.Faction = req.Type
	p.SaveFight()
	p.SendResponse(packetId, res, res.Code)
}

// RequestGetContractInfo 获取合约数据
func RequestGetContractInfo(packetId uint32, args interface{}, p *player.Player) {
	resp := &msg.GetContractInfoResp{
		PassInfo: p.UserData.WeekPass.ContractInfo,
	}
	p.SendResponse(packetId, resp, msg.ErrCode_SUCC)
}

// RequestSecretInfo 获取秘境数据
func RequestSecretInfo(packetId uint32, args interface{}, p *player.Player) {
	resp := &msg.GetSecretInfoResp{
		Count: uint32(p.UserData.WeekPass.SecretCount),
		State: p.UserData.WeekPass.SecretBoxState,
	}
	p.SendResponse(packetId, resp, msg.ErrCode_SUCC)
}
