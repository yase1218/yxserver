package model

import (
	"github.com/zy/game_data/template"
	"kernel/tools"
)

type AccountPlayMethod struct {
	AccountId       int64
	NextRefreshTime uint32
	Data            []*PlayMethodData
}

type PlayMethodData struct {
	BtType      int
	TotalTimes  int
	WeaponIds   []uint32
	MissionData []*Mission
	MaxDamage   int
}

func NewPlayMethodInfo(btType, totalTimes int) *PlayMethodData {
	return &PlayMethodData{
		BtType:      btType,
		TotalTimes:  totalTimes,
		WeaponIds:   make([]uint32, 0),
		MissionData: make([]*Mission, 0),
		MaxDamage:   0,
	}
}

type SimpleMission struct {
	MissionId  int
	IsPass     bool // 是否通关
	UpdateTime uint32
}

func NewAccountPlayMethod(accountId int64, lst []*template.JPlayMethod) *AccountPlayMethod {
	ret := &AccountPlayMethod{
		AccountId: accountId,
	}
	for i := 0; i < len(lst); i++ {
		ret.Data = append(ret.Data, NewPlayMethodInfo(lst[i].Type, lst[i].Limit))
	}
	return ret
}

func NewSimpleMission(missionId int, isPass bool) *SimpleMission {
	return &SimpleMission{
		MissionId:  missionId,
		IsPass:     isPass,
		UpdateTime: tools.GetCurTime(),
	}
}
