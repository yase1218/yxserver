package model

import "kernel/tools"

type AccountMission struct {
	AccountId  uint64
	Missions   []*Mission
	Challenges []*Mission // 挑战关卡
	ExtraIds   []int
}

type Mission struct {
	MissionId      int
	CompleteTime   uint32 // 关卡完成时间
	KillMonsterNum uint32 // 杀怪数量
	BoxState       uint32 // 宝箱状态
	BoxRewardState uint32 // 宝箱奖励状态
	IsPass         bool   // 是否通关
	CreateTime     uint32
	UpdateTime     uint32
	BeforeStory    uint32
	AfterStory     uint32
}

// MissionReward 首次通关奖励
type MissionReward struct {
	ServerId  uint32
	Uid       uint64
	MissionId int
	Tp        int // 类型
	PassTime  uint32
}

func NewAccountMission(accountId uint64) *AccountMission {
	ret := &AccountMission{
		AccountId: accountId,
	}
	ret.Missions = make([]*Mission, 0, 0)
	ret.Challenges = make([]*Mission, 0, 0)
	return ret
}

func NewMission(missionId int, completeTime, killMonsterNum uint32, isPass bool) *Mission {
	curTime := tools.GetCurTime()
	newMission := &Mission{
		MissionId: missionId,
		//CompleteTime:   completeTime,
		KillMonsterNum: killMonsterNum,
		IsPass:         isPass,
		CreateTime:     curTime,
		UpdateTime:     0,
	}
	if isPass {
		newMission.CompleteTime = completeTime
	}

	return newMission
}

func NewMissionReward(accountId uint64, serverId, passTime uint32, missionId, tp int) *MissionReward {
	return &MissionReward{
		ServerId:  serverId,
		MissionId: missionId,
		Uid:       accountId,
		PassTime:  passTime,
		Tp:        tp,
	}
}
