package model

import (
	template2 "github.com/zy/game_data/template"
	"kernel/tools"
)

type AccountShip struct {
	AccountId int64
	Ships     []*Ship
}

// Ship 机甲
type Ship struct {
	Id          uint32
	Level       uint32
	StarLevel   uint32
	CreateTime  uint32
	UpdateTime  uint32
	Attrs       map[uint32]*Attr
	SupportAttr []*Attr
	CoatMap     map[int]*CoatItem // 皮肤
}

func NewShip(id uint32) *Ship {
	coatCfg := template2.GetShipTemplate().GetShip(id).GetOriginalCoat()
	return &Ship{
		Id:         id,
		Level:      1,
		StarLevel:  0,
		CreateTime: tools.GetCurTime(),
		Attrs:      make(map[uint32]*Attr),
		CoatMap:    map[int]*CoatItem{coatCfg.CoatId: NewCoatItem(coatCfg)},
	}
}

func (s *Ship) GetCoatAttr() []*template2.AttrItem {
	for _, item := range s.CoatMap {
		if item.Status == CoatOn {
			coatCfg := template2.GetCoatTemplate().GetCoatCfg(int(s.Id), item.Id)
			return coatCfg.InitAttr
		}
	}
	return nil
}

func (s *Ship) PutOnCoat(coatCfg *template2.JCoatCfg, isActive bool) {
	for _, item := range s.CoatMap {
		if item.Status == CoatOn {
			item.Status = CoatActive
			itemCfg := template2.GetCoatTemplate().GetCoatCfg(int(s.Id), item.Id)
			for i := 0; i < len(itemCfg.InitAttr); i++ {
				if attr, ok := s.Attrs[itemCfg.InitAttr[i].Id]; ok {
					attr.DelValue(itemCfg.InitAttr[i].Value)
				}
			}
		}
	}
	for i := 0; i < len(coatCfg.InitAttr); i++ {
		if attr, ok := s.Attrs[coatCfg.InitAttr[i].Id]; ok {
			attr.AddValue(coatCfg.InitAttr[i].Value)
		} else {
			s.Attrs[coatCfg.InitAttr[i].Id] = &Attr{
				Id:  coatCfg.InitAttr[i].Id,
				Add: coatCfg.InitAttr[i].Value,
			}
		}
	}
	if isActive {
		s.CoatMap[coatCfg.CoatId] = NewCoatItem(coatCfg)
	} else {
		s.CoatMap[coatCfg.CoatId].Status = CoatOn
	}
}
