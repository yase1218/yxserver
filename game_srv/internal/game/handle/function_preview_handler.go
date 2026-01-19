package handle

import (
	"gameserver/internal/game/builder"
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

func FunctionPreview(packetId uint32, args interface{}, p *player.Player) {
	res := new(msg.FunctionPreviewAck)
	builder.BuildFunctionPreviewAck(p.UserData, res)
	p.SendResponse(packetId, res, msg.ErrCode_SUCC)
}

func FunctionPreviewReward(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.FunctionPreviewRewardReq)
	res := &msg.FunctionPreviewRewardAck{
		Err: msg.ErrCode_SUCC,
	}
	service.FunctionPreviewReward(p, req, res)
	p.SendResponse(packetId, res, res.Err)
}

func HandleGetResourcesPassBaseData(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.GetResourcesPassBaseDataReq)
	res := service.GetResourcesPassBaseData(p, req.PassType)
	if res != nil {
		p.SendResponse(packetId, res, msg.ErrCode_SUCC)
	}
}

func HandleBuyResoucePassItem(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.BuyResoucesPassItemReq)
	code, count := service.BuyResoucesPassItem(p, req.PassType, req.Num)
	res := &msg.BuyResoucesPassItemResp{
		Code:     code,
		BuyCount: count,
		PassType: req.PassType,
	}

	p.SendResponse(packetId, res, res.Code)
}

func HandleResourcesAttack(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.ResoucePassAttackReq)
	code, num := service.ResourcePassAttack(p, req.PassType)
	res := &msg.ResoucePassAttackResp{}
	res.Code = code
	res.ItemNum = num
	res.Data = service.GetResourcesPassBaseData(p, req.PassType)
	p.SendResponse(packetId, res, code)
}

func HandleResourceStateChange(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.ResourcePassStateChangeReq)
	code := service.ResourcePassStateUpdate(p, req.PassType, req.State)
	resp := &msg.ResourcePassStateChangeResp{}
	resp.Code = code
	resp.PassType = req.PassType
	p.SendResponse(packetId, resp, resp.Code)
}

func HandleGetResourcesPassRankList(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.ResourcePassRankListReq)
	code, list := service.ResourcesPassRankList(p, req.PassType)
	resp := &msg.ResourcePassRankListResp{}
	resp.List = list
	p.SendResponse(packetId, resp, code)
}
