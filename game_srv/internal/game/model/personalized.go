package model

import (
	"github.com/zy/game_data/template"
	"kernel/tools"
	"msg"
)

type AccountPersonalized struct {
	AccountId uint64
	ItemsMap  map[int]*PersonalizedItem
}

type PersonalizedItem struct {
	Id              int    // 配置id
	PopNum          uint32 // 弹窗次数
	NextRefreshTime uint32 // 下次刷新时间
	OutTime         uint32 // 倒计时时间戳
	CreateTime      uint32
	UpdateTime      uint32
}

func NewPersonalizedItem(cfg *template.JPersonalizedCfg) *PersonalizedItem {
	now := tools.GetCurTime()
	item := &PersonalizedItem{
		Id:         cfg.TypeID,
		CreateTime: now,
		UpdateTime: now,
	}
	if cfg.OutTime > 0 {
		item.OutTime = tools.GetCurTime() + uint32(cfg.OutTime)
	}
	item.resetRefreshTime()
	return item
}

func (item *PersonalizedItem) resetRefreshTime() { //
	cfg := template.GetPersonalizedTemplate().GetPersonalizedCfg(item.Id)
	switch cfg.PersonalCycle {
	case 0: // 日刷新
		item.NextRefreshTime = tools.GetDailyRefreshTimeByHour(0) // 每天0点刷新
	case 1: // 周刷新
		item.NextRefreshTime = tools.GetWeeklyRefreshTime(0) // 每周一0点刷新
	case 2: // 月刷新
		item.NextRefreshTime = tools.GetMonthRefreshTimeByHour(0) // 每月1号0点刷新
	}
	curTime := tools.GetCurTime()
	item.PopNum = 1
	item.UpdateTime = curTime
	if cfg.OutTime > 0 {
		item.OutTime = curTime + uint32(cfg.OutTime)
	}
}

func (item *PersonalizedItem) ToProto() *msg.Personalized { //
	return &msg.Personalized{ItemId: uint32(item.Id), OutTime: item.OutTime}
}

func NewAccountPersonalized(uid uint64) *AccountPersonalized {
	ret := &AccountPersonalized{
		AccountId: uid,
		ItemsMap:  make(map[int]*PersonalizedItem),
	}
	return ret
}

//// 获取所有可推送数据
//func (a *AccountPersonalized) GetPersonalizedItems(personalizedType, stageId, itemId int) []*PersonalizedItem {
//	var items []*PersonalizedItem
//	configs := template.GetPersonalizedTemplate().GetPersonalizedConfigs(personalizedType, stageId, itemId)
//	for i := 0; i < len(configs); i++ {
//		if _, ok := a.ItemsMap[configs[i].TypeID]; !ok {
//			a.ItemsMap[configs[i].TypeID] = NewPersonalizedItem(configs[i].TypeID)
//		}
//		if a.CanPop(a.ItemsMap[configs[i].TypeID], configs[i]) {
//			items = append(items, a.ItemsMap[configs[i].TypeID])
//		}
//	}
//	return items
//}

// 是否可以推送
func (a *AccountPersonalized) CanPop(item *PersonalizedItem, cfg *template.JPersonalizedCfg) bool {
	curTime := tools.GetCurTime()
	//if cfg.OutTime > 0 && item.OutTime >= curTime { // 超过倒计时
	//	return true
	//}
	if item.NextRefreshTime > 0 { // 有刷新时间
		if item.NextRefreshTime > curTime { // 未到刷新时间
			if cfg.PersonalNum > 0 && item.PopNum >= uint32(cfg.PersonalNum) { // 推送次数达到上限
				if cfg.OutTime > 0 && item.OutTime >= curTime { // 在倒计时内
					return true
				}
				return false
			}
		} else {
			item.resetRefreshTime()
		}
	} else {
		if cfg.PersonalNum > 0 && item.PopNum >= uint32(cfg.PersonalNum) { // 推送次数达到上限
			return false
		}
	}

	if cfg.OutTime > 0 {
		if item.OutTime < curTime { // 已过时
			item.OutTime = curTime + uint32(cfg.OutTime)
			item.PopNum++
			item.UpdateTime = curTime
		}
	} else {
		item.PopNum++
		item.UpdateTime = curTime
	}
	return true
}

func (a *AccountPersonalized) GetItem(id int) *PersonalizedItem {
	return a.ItemsMap[id]
}
