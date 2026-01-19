package model

import (
	"github.com/zy/game_data/template"
	"kernel/tools"
	"msg"
)

const (
	CoatUnActive = iota // 未激活
	CoatActive          // 已激活未穿上
	CoatOn              // 穿上
)

type CoatItem struct {
	Id         int   // 配置id
	Status     int32 // 状态
	CreateTime uint32
	UpdateTime uint32
}

func NewCoatItem(cfg *template.JCoatCfg) *CoatItem {
	now := tools.GetCurTime()
	item := &CoatItem{
		Id:         cfg.CoatId,
		CreateTime: now,
		UpdateTime: now,
	}
	//if len(cfg.CoatItem) == 0 {
	//	item.Status = CoatOn
	//}
	item.Status = CoatOn
	return item
}

func (item *CoatItem) ToProto() *msg.Coat { //
	return &msg.Coat{ItemId: int32(item.Id), Status: item.Status}
}

//func (a *AccountCoat) GetItem(id int) *CoatItem {
//	return a.ItemsMap[id]
//}
