package model

import "kernel/tools"

type AccountAppearance struct {
	AccountId   int64
	Appearances []*Appearance
	Attrs       map[uint32]*Attr
}

type Appearance struct {
	Id         uint32
	LimitDate  uint32
	UpdateTime uint32
}

func NewAccountAppearance(accountId int64, defaultIds map[uint32]uint32) *AccountAppearance {
	ret := &AccountAppearance{
		AccountId: accountId,
		Attrs:     make(map[uint32]*Attr),
	}
	for _, id := range defaultIds {
		ret.Appearances = append(ret.Appearances, NewAppearance(id))
	}
	return ret
}

func NewAppearance(id uint32) *Appearance {
	return &Appearance{
		Id:         id,
		UpdateTime: tools.GetCurTime(),
	}
}
