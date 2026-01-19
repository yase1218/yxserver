package model

import "kernel/tools"

type AccountShop struct {
	AccountId int64
	Items     []*ShopItem // 记录有刷新机制和购买的（无刷新机制）商品
}

type ShopItem struct {
	Id              uint32 // 配置id
	ItemId          uint32 // 道具id
	ItemNum         uint32
	BuyTimes        uint32 // 购买次数
	NextRefreshTime uint32 // 下次刷新时间
	Ver             uint32 // 控制手动刷新 TODO 没用上
	Lock            bool   // 是否锁
	Times           uint32 // 倍数
	ResetTimesTime  uint32 // 重置倍数时间 TODO 没用上
	CreateTime      uint32
	UpdateTime      uint32
	ConfigVer       uint32 // 配置版本提升则删除商品
}

type ShopItemVer struct {
	ShopItemId uint32
	Ver        uint32
	UpdateTime uint32
}

func NewShopItem(id uint32, lock bool, configVer uint32) *ShopItem {
	return &ShopItem{
		Id:         id,
		CreateTime: tools.GetCurTime(),
		Lock:       lock,
		Times:      1,
		ConfigVer:  configVer,
	}
}

func NewAccountShop(accountId int64) *AccountShop {
	ret := &AccountShop{
		AccountId: accountId,
	}
	ret.Items = make([]*ShopItem, 0, 0)
	return ret
}
