package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

func LoadEquipStageReqHandle(pid uint32, args interface{}, p *player.Player) {
	service.LoadEquipStage(pid, p)
}

func EquipStageDropRecordReqHandle(pid uint32, args interface{}, p *player.Player) {
	service.LoadEquipStageDropRecord(pid, p)
}

func EquipStageTeamInviteReqHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.EquipStageTeamInviteReq)
	service.EquipStageTeamInvite(pid, p, req)
}

func EquipStageTeamAcceptReqHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.EquipStageTeamAcceptReq)
	service.EquipStageTeamAccept(pid, p, req)
}

func EquipStageTeamLeaveReq(pid uint32, args interface{}, p *player.Player) {
	service.EquipStageTeamLeave(pid, p, false)
}

func EquipStageTeamDissolveReqHandle(pid uint32, args interface{}, p *player.Player) {
	service.EquipStageTeamDissolve(pid, p)
}

func EquipStageMatchReqHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.EquipStageMatchReq)
	service.EquipStageMatch(pid, p, req)
}

func EquipStageAcceptReqHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.EquipStageAcceptReq)
	service.EquipStageAccept(pid, p, req)
}

func EquipStageMatchCancelReqHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.EquipStageMatchCancelReq)
	service.EquipStageMatchCacnel(pid, p, req)
}

func EquipStageLoadReqHandle(pid uint32, args interface{}, p *player.Player) {
	service.EquipStageLoad(pid, p)
}

func EquipStageRewardPickReqHandle(pid uint32, args interface{}, p *player.Player) {
	service.EquipStageReward(pid, p)
}

func EquipStageBuyRewardReqHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.EquipStageBuyRewardReq)
	service.EquipStageBuyReward(pid, p, req)
}

func EquipStageLeaveReqHandle(pid uint32, args interface{}, p *player.Player) {
	service.EquipStageLeave(p)
}

func CreateEquipStageTeamHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestCreateEquipStageTeam)
	service.ReqCreateEquipStageTeam(pid, p, req)
}
