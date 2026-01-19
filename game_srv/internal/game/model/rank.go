package model

import "kernel/tools"

// MissionRank 关卡排行榜
type MissionRank struct {
	AccountId   uint32
	ServerId    uint32
	MissionId   int    // 关卡id
	MissionTime uint32 // 关卡时间
	Likes       uint32
	PassTime    uint32
	Tp          int // 关卡类型
}

type PetRank struct {
	AccountId  uint32
	ServerId   uint32
	PetId      uint32 // 宠物id
	UpdateTime uint32
	Likes      uint32
}

func NewPetRank(accountId, serverId, petId, updateTime uint32) *PetRank {
	return &PetRank{
		AccountId:  accountId,
		ServerId:   serverId,
		PetId:      petId,
		UpdateTime: updateTime,
	}
}

type SpecialMissionRank struct {
	ServerId   uint32
	MissionId  int
	AccountId  uint32
	PassTime   uint32
	UpdateTime uint32
}

func NewSpecialMissionRank(serverId, accountId, passTime uint32, missionId int) *SpecialMissionRank {
	return &SpecialMissionRank{
		ServerId:   serverId,
		MissionId:  missionId,
		AccountId:  accountId,
		PassTime:   passTime,
		UpdateTime: tools.GetCurTime(),
	}
}

type DesertHurtRank struct {
	ServerId   uint32
	AccountId  uint32
	Hurt       uint64
	UpdateTime uint32
}

func NewDesertHurtRank(serverId, accountId uint32, hurt uint64) *DesertHurtRank {
	return &DesertHurtRank{
		ServerId:   serverId,
		AccountId:  accountId,
		Hurt:       hurt,
		UpdateTime: tools.GetCurTime(),
	}
}
