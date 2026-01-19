package model

type LocType int

const (
	Rest    LocType = iota // 休息
	Fights                 // 出战
	HelpWar                // 助战
)

type AccountPet struct {
	Pets map[uint32]*Pet
	Lv   int32 // 等级
}

func NewAccountPet() *AccountPet {
	ret := &AccountPet{
		Pets: make(map[uint32]*Pet),
		Lv:   1,
	}
	return ret
}

type Pet struct {
	BaseId uint32
	Loc    LocType
	LocIdx int16
	StarLv int16
}

func NewPet(baseId uint32) *Pet {
	return &Pet{
		BaseId: baseId,
		Loc:    Rest,
		LocIdx: -1,
		StarLv: 0,
	}
}
