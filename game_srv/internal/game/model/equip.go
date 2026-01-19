package model

import (
	"kernel/tools"
	"msg"

	"github.com/zy/game_data/template"
)

type AccountEquip struct {
	AccountId    int64
	EquipData    []*Equip
	EquipPosData []*EquipPos
	SuitReward   []*SuitInfo
	UseEquipSuit uint32 // 使用的装备套装
	SyncEquipPos bool   // 是否需要同步装备部位
	EquipSuits   []*EquipSuit
	GemBag       map[uint64]*GemBagSlot // key: uuid
	GemPos       [][]uint64             // key: pos value: uuid
}

type SuitPosInfo struct {
	Pos     uint32
	EquipId uint32
}

type SuitInfo struct {
	SuitId  uint32
	PosData []*SuitPosInfo
}

// Equip 装备
type Equip struct {
	Id         uint32
	Num        uint32
	CreateTime uint32
	UpdateTime uint32
}

// EquipPos 装备位置信息
type EquipPos struct {
	Pos         uint32
	Level       uint32
	EquipId     uint32
	UpdateTime  uint32
	Attr        map[uint32]*Attr
	AffixAttr   map[uint32]*Attr
	AffixSkills []uint32
}

// EquipSuit 装备套装
type EquipSuit struct {
	SuitId   uint32
	EquipIds []uint32 // 固定6个 没有的部位填0
	SkillId  uint32   // 技能id
}

// GemBagSlot 宝石背包格子数据
type GemBagSlot struct {
	Uuid       uint64 // 配置表id*GemUuidParam*GemUuidParam + AffixIdx*GemUuidParam + AffixType
	Num        uint32 // 数量
	Lock       bool   // 是否锁定
	CreateTime uint32
	UpdateTime uint32
}

func NewGem(uuid uint64, num uint32) *GemBagSlot {
	return &GemBagSlot{
		Uuid:       uuid,
		Num:        num,
		CreateTime: tools.GetCurTime(),
	}
}

func NewEquip(id, num uint32) *Equip {
	return &Equip{
		Id:         id,
		Num:        num,
		CreateTime: tools.GetCurTime(),
	}
}

func NewEquipPos(pos, equipId uint32) *EquipPos {
	return &EquipPos{
		Pos:         pos,
		EquipId:     equipId,
		Level:       1,
		UpdateTime:  tools.GetCurTime(),
		Attr:        make(map[uint32]*Attr),
		AffixAttr:   make(map[uint32]*Attr),
		AffixSkills: make([]uint32, 0),
	}
}

func NewAccountEquip(accountId int64) *AccountEquip {
	ret := &AccountEquip{
		AccountId: accountId,
	}
	ret.EquipData = make([]*Equip, 0, 0)
	ret.EquipPosData = make([]*EquipPos, 0, 0)
	ret.SuitReward = make([]*SuitInfo, 0, 0)
	ret.EquipSuits = make([]*EquipSuit, 0, 0)
	ret.GemBag = make(map[uint64]*GemBagSlot)
	ret.GemPos = make([][]uint64, msg.EquipPos_EquipPos_Max-1)

	for i := 0; i < len(ret.GemPos); i++ {
		ret.GemPos[i] = make([]uint64, template.GetSystemItemTemplate().GemSlotMax)
		for j := 0; j < len(ret.GemPos[i]); j++ {
			ret.GemPos[i][j] = 0
		}
	}

	return ret
}

func NewSuitInfo(suitId uint32) *SuitInfo {
	return &SuitInfo{
		SuitId:  suitId,
		PosData: make([]*SuitPosInfo, 0, 0),
	}
}

func NewSuitPosInfo(pos, equipId uint32) *SuitPosInfo {
	return &SuitPosInfo{
		Pos:     pos,
		EquipId: equipId,
	}
}

func NewEquipSuit(id uint32, equipIds []uint32, skillId uint32) *EquipSuit {
	return &EquipSuit{
		SuitId:   id,
		EquipIds: equipIds,
		SkillId:  skillId,
	}
}
