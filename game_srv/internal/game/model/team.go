package model

import "kernel/tools"

type AccountTeam struct {
	AccountId  int64
	TeamData   []*Team
	BattleData []*BattleTeam
}

type Team struct {
	TeamId      uint32
	ShipId      uint32   // 机甲
	RoleId      uint32   // 驾驶员
	SupportShip []uint32 // 支援机甲
	Equips      []uint32 // 装备
	CreateTime  uint32
	UpdateTime  uint32
}

// BattleTeam 玩法编队
type BattleTeam struct {
	BattleType uint32
	TeamId     uint32
	CreateTime uint32
	UpdateTime uint32
}

func NewAccountTeam(accountId int64) *AccountTeam {
	ret := &AccountTeam{
		AccountId: accountId,
	}
	ret.TeamData = make([]*Team, 0, 0)
	ret.BattleData = make([]*BattleTeam, 0, 0)
	return ret
}

func NewTeam(teamId, shipId, roleId uint32, supportShip, equips []uint32) *Team {
	curTime := tools.GetCurTime()
	return &Team{
		TeamId:      teamId,
		ShipId:      shipId,
		RoleId:      roleId,
		SupportShip: supportShip,
		Equips:      equips,
		CreateTime:  curTime,
		UpdateTime:  0,
	}
}

func NewBattleTeam(battleType uint32, teamId uint32) *BattleTeam {
	return &BattleTeam{
		BattleType: battleType,
		TeamId:     teamId,
		CreateTime: tools.GetCurTime(),
		UpdateTime: 0,
	}
}
