package model

import "time"

type ItemLog struct {
	AccountId  int64
	ItemId     uint32
	Delta      int32
	CurNum     int64
	Src        int32
	Remark     string
	UpdateTime uint32
}

func NewItemLog(accountId int64, itemId uint32, delta int32, curNum int64, src int32, remark string) *ItemLog {
	curTime := uint32(time.Now().Unix())
	return &ItemLog{
		AccountId:  accountId,
		ItemId:     itemId,
		Delta:      delta,
		CurNum:     curNum,
		Src:        src,
		Remark:     remark,
		UpdateTime: curTime,
	}
}
