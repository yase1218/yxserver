package model

import "time"

type AccountItem struct {
	AccountId int64
	Items     []*Item
}

type Item struct {
	Id         uint32
	Num        uint64
	LimitDate  uint32
	CreateTime uint32
	UpdateTime uint32
}

type SimpleItem struct {
	Id  uint32 `json:"id"`
	Num uint32 `json:"num"`
	Src uint32 // 来源 1 重复转换来的
}

type FloatItem struct {
	Id  uint32
	Num float64
}

func NewItem(id, limitDate uint32, num uint32) *Item {
	curTime := uint32(time.Now().Unix())
	return &Item{
		Id:         id,
		Num:        uint64(num),
		LimitDate:  limitDate,
		CreateTime: curTime,
	}
}

func NewAccountItem(accountId int64) *AccountItem {
	ret := &AccountItem{
		AccountId: accountId,
	}
	ret.Items = make([]*Item, 0, 0)
	return ret
}

func NewSimpleItem(id, num uint32) *SimpleItem {
	return &SimpleItem{
		Id:  id,
		Num: num,
	}
}
