package service

import (
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"kernel/tools"
	"msg"

	template2 "github.com/zy/game_data/template"
)

// UpdateTeam 更新编队
func UpdateTeam(p *player.Player, teamId, shipId, roleId uint32, supportShips, equip []uint32) msg.ErrCode {
	if shipId == 0 {
		return msg.ErrCode_TEAM_SHIP_NOT_EMPTY
	}
	if roleId == 0 {
		return msg.ErrCode_TEAM_ROLE_NOT_EMPTY
	}
	if uint32(len(supportShips)) > template2.GetSystemItemTemplate().TeamSupportShipNum {
		return msg.ErrCode_TEAM_SUPPORT_SHIP_OVER_LIMIT
	}
	if uint32(len(equip)) > template2.GetSystemItemTemplate().TeamEquipNum {
		return msg.ErrCode_TEAM_EQUIP_OVER_LIMIT
	}

	if ship := getShip(p, shipId); ship == nil {
		return msg.ErrCode_SHIP_NOT_EXIST
	}
	// if role := getRole(p, roleId); role == nil {
	// 	return msg.ErrCode_ROLE_NOT_EXIST
	// }
	for i := 0; i < len(supportShips); i++ {
		if supportShips[i] == 0 {
			continue
		}
		if supportShips[i] == shipId {
			return msg.ErrCode_INVALID_DATA
		}

		if ship := getShip(p, supportShips[i]); ship == nil {
			return msg.ErrCode_SHIP_NOT_EXIST
		}
	}

	for i := 0; i < len(equip); i++ {
		equipId := equip[i]
		if equipId > 0 {
			if equip := getEquip(p, equipId); equip == nil {
				return msg.ErrCode_EQUIP_NOT_EXIST
			}
			jEquip := template2.GetEquipTemplate().GetEquip(equipId)
			if jEquip.Data.Pos != uint32(i+1) {
				return msg.ErrCode_INVALID_DATA
			}
		}
	}

	team := getTeam(p, teamId)
	if team == nil {
		return msg.ErrCode_TEAM_NOT_EXIST
	} else {
		bUpdate := true
		if team.ShipId == shipId &&
			team.RoleId == roleId &&
			tools.ListUint32Equal(team.SupportShip, supportShips) &&
			tools.ListUint32Equal(team.Equips, equip) {
			bUpdate = false
		}
		if bUpdate {
			team.ShipId = shipId
			team.RoleId = roleId
			team.SupportShip = supportShips
			team.Equips = equip
			p.SaveTeam()
		}
	}
	return msg.ErrCode_SUCC
}

// equipChangeUpdateTeam 装备升阶变化影响的队伍变化
func equipChangeUpdateTeam(p *player.Player, equipId uint32) []uint32 {
	var ret []uint32
	//for i := 0; i < len(playerData.Team.TeamData); i++ {
	//	team := playerData.Team.TeamData[i]
	//	for k := 0; k < len(team.Equips); k++ {
	//		if team.Equips[k] == equipId {
	//			ret = append(ret, team.TeamId)
	//		}
	//	}
	//}
	return ret
}

func updateTeamEquip(p *player.Player, teamIds []uint32, equipId uint32) {
	if len(teamIds) == 0 {
		return
	}

	jEquip := template2.GetEquipTemplate().GetEquip(equipId)
	if jEquip == nil {
		return
	}

	res := &msg.NotifyTeamChange{}
	for i := 0; i < len(teamIds); i++ {
		if team := getTeam(p, teamIds[i]); team != nil {
			pos := jEquip.Data.Pos
			team.Equips[pos-1] = jEquip.Data.Id
			res.Data = append(res.Data, ToProtocolTeam(team))
		}
	}

	if len(res.Data) > 0 {
		p.SaveTeam()
		p.SendNotify(res)
	}
}

// UpdateBattleTeam 更新战斗编队
func UpdateBattleTeam(p *player.Player, battleType msg.BattleType, teamId uint32) msg.ErrCode {
	if battleType != msg.BattleType_Battle_Main && battleType != msg.BattleType_Battle_Challenge {
		return msg.ErrCode_TEAM_BATTLE_TYPE_NOT_EXIST
	}

	if team := getTeam(p, teamId); team == nil {
		return msg.ErrCode_TEAM_NOT_EXIST
	}

	battleTeam := getBattleTeam(p, battleType)
	if battleTeam == nil {
		battleTeam = model.NewBattleTeam(uint32(battleType), teamId)
		p.UserData.Team.BattleData = append(p.UserData.Team.BattleData, battleTeam)
	} else {
		battleTeam.TeamId = teamId
		p.UserData.Team.BattleData = append(p.UserData.Team.BattleData, battleTeam)
	}
	p.SaveTeam()
	return msg.ErrCode_SUCC
}

// getBattleTeam 获得编队
func getBattleTeam(p *player.Player, battleType msg.BattleType) *model.BattleTeam {
	for i := 0; i < len(p.UserData.Team.BattleData); i++ {
		if p.UserData.Team.BattleData[i].BattleType == uint32(battleType) {
			return p.UserData.Team.BattleData[i]
		}
	}
	return nil
}

// getTeam 获得编队
func getTeam(p *player.Player, teamId uint32) *model.Team {
	for i := 0; i < len(p.UserData.Team.TeamData); i++ {
		if p.UserData.Team.TeamData[i].TeamId == teamId {
			return p.UserData.Team.TeamData[i]
		}
	}
	return nil
}

func ToProtocolTeam(team *model.Team) *msg.Team {
	return &msg.Team{
		TeamId:         team.TeamId,
		ShipId:         team.ShipId,
		RoleId:         team.RoleId,
		SupportShipIds: team.SupportShip,
		Equips:         team.Equips,
	}
}

func ToProtocolTeams(team []*model.Team) []*msg.Team {
	var ret []*msg.Team
	for i := 0; i < len(team); i++ {
		ret = append(ret, ToProtocolTeam(team[i]))
	}
	return ret
}

func ToProtocolBattleTeam(team *model.BattleTeam) *msg.BattleTeam {
	return &msg.BattleTeam{
		BatType: msg.BattleType(team.BattleType),
		TeamId:  team.TeamId,
	}
}

func ToProtocolBattleTeams(team []*model.BattleTeam) []*msg.BattleTeam {
	var ret []*msg.BattleTeam
	for i := 0; i < len(team); i++ {
		ret = append(ret, ToProtocolBattleTeam(team[i]))
	}
	return ret
}
