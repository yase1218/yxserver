package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

// RequestLoadTeamsHandle 加载编队
func RequestLoadTeamsHandle(packetId uint32, args interface{}, p *player.Player) {
	team := p.UserData.Team
	retMsg := &msg.ResponseLoadTeams{Result: msg.ErrCode_SUCC}
	if len(team.TeamData) > 0 {
		retMsg.Teams = append(retMsg.Teams, service.ToProtocolTeams(team.TeamData)...)
	}
	if len(team.BattleData) > 0 {
		retMsg.BattleTeams = append(retMsg.BattleTeams, service.ToProtocolBattleTeams(team.BattleData)...)
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestUpdateTeamHandle 更新编队
func RequestUpdateTeamHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestUpdateTeam)
	err := service.UpdateTeam(p, req.Data.TeamId,
		req.Data.ShipId, req.Data.RoleId, req.Data.SupportShipIds, req.Data.Equips)
	retMsg := &msg.ResponseUpdateTeam{Result: err, Data: req.Data}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestSetBattleTeamHandle 设置玩法编队
func RequestSetBattleTeamHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestSetBattleTeam)
	err := service.UpdateBattleTeam(p, req.BatType, req.TeamId)
	retMsg := &msg.ResponseSetBattleTeam{Result: err, BatType: req.BatType, TeamId: req.TeamId}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}
