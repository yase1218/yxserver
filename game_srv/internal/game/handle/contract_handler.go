package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

func Contract(packetId uint32, args interface{}, p *player.Player) {
	//req := args.(*msg.ContractReq)
	res := &msg.ContractAck{}
	res.Result = service.Contract(p, res)
	p.SendResponse(packetId, res, res.Result)
}

func ContractSign(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.ContractSignReq)
	res := &msg.ContractSignAck{}
	res.Result = service.ContractSign(p, req, res)
	res.TaskId = req.GetTaskId()
	p.SendResponse(packetId, res, res.Result)
}

func ContractCancel(packetId uint32, args interface{}, p *player.Player) {
	//req := args.(*msg.ContractCancelReq)
	res := &msg.ContractCancelAck{}
	res.Result = service.ContractCancel(p, res)
	p.SendResponse(packetId, res, res.Result)
}

func ContractRand(packetId uint32, args interface{}, p *player.Player) {
	//req := args.(*msg.ContractRandReq)
	res := &msg.ContractRandAck{}
	res.Result = service.ContractRand(p, res)
	p.SendResponse(packetId, res, res.Result)
}

func ContractReward(packetId uint32, args interface{}, p *player.Player) {
	//req := args.(*msg.ContractRewardReq)
	res := &msg.ContractRewardAck{}
	res.Result = service.ContractReward(p, res)
	p.SendResponse(packetId, res, res.Result)
}
