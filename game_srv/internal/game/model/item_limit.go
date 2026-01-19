package model

import "time"

type (
	ItemLimitUnit struct {
		Num     int64
		ResetAt time.Time
	}
	AccountItemLimit struct {
		AccountId int64
		Items     map[uint32]*ItemLimitUnit
	}
)

func NewAccountItemLimit(accountId int64) *AccountItemLimit {
	return &AccountItemLimit{
		AccountId: accountId,
		Items:     make(map[uint32]*ItemLimitUnit),
	}
}

func (a *AccountItemLimit) Add(itemId uint32, num int64) {
	limit, ok := a.Items[itemId]
	if !ok {
		a.Items[itemId] = new(ItemLimitUnit)
		limit = a.Items[itemId]
	}
	limit.Num += num
}

func (a *AccountItemLimit) GetNum(itemId uint32) int64 {
	limit, ok := a.Items[itemId]
	if !ok {
		return 0
	}
	return limit.Num
}

func (a *AccountItemLimit) GetResetAt(itemId uint32) time.Time {
	limit, ok := a.Items[itemId]
	if !ok {
		return time.Time{}
	}
	return limit.ResetAt
}

func (a *AccountItemLimit) Reset(itemId uint32) {
	a.Items[itemId] = new(ItemLimitUnit)
}
