package model

import "kernel/tools"

type AccountWeapon struct {
	AccountId        int64
	WeaponLibExp     uint32
	WeaponLibLevel   uint32
	Attrs            map[uint32]*Attr
	Weapons          []*Weapon
	SecondaryWeapons []*SecondaryWeapon
}

type SecondaryWeapon struct {
	Pos        uint32
	WeaponId   uint32
	UpdateTime uint32
}

type Weapon struct {
	Id         uint32
	Level      uint32
	SkillId    uint32
	UpdateTime uint32
	Attrs      map[uint32]*Attr
}

func NewAccountWeapon(accountId int64) *AccountWeapon {
	return &AccountWeapon{
		AccountId:        accountId,
		WeaponLibLevel:   1,
		Weapons:          make([]*Weapon, 0, 0),
		SecondaryWeapons: make([]*SecondaryWeapon, 0, 0),
		Attrs:            make(map[uint32]*Attr),
	}
}

func NewWeapon(id, level, skill uint32) *Weapon {
	return &Weapon{
		Id:         id,
		Level:      level,
		SkillId:    skill,
		UpdateTime: tools.GetCurTime(),
		Attrs:      make(map[uint32]*Attr),
	}
}

func NewSecondaryWeapon(pos, weaponId uint32) *SecondaryWeapon {
	return &SecondaryWeapon{
		WeaponId:   weaponId,
		Pos:        pos,
		UpdateTime: tools.GetCurTime(),
	}
}
