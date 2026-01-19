package event

import (
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
)

// AddItemEvent 添加道具
type AddItemEvent struct {
	IEvent
	PlayerData *player.Player
	ItemId     uint32
	Num        uint32
	CurNum     int64
	Remark     string
	ItemSrc    publicconst.ItemSource
	SendClient bool
	f          EventFunc
}

func NewAddItemEvent(p *player.Player, id, num uint32, curNum int64, remark string, src publicconst.ItemSource, sendClient bool, cb EventFunc) *AddItemEvent {
	return &AddItemEvent{
		PlayerData: p,
		ItemId:     id,
		Num:        num,
		CurNum:     curNum,
		ItemSrc:    src,
		SendClient: sendClient,
		f:          cb,
	}
}

func (a *AddItemEvent) RouteID() uint64 {
	return a.PlayerData.GetUserId()
}

func (a *AddItemEvent) CallBack() EventFunc {
	return a.f
}

// CostItemEvent 消耗道具事件
type CostItemEvent struct {
	IEvent
	PlayerData *player.Player
	ItemId     uint32
	Num        uint32
	CurNum     int64
	Remark     string
	ItemSrc    publicconst.ItemSource
	SendClient bool
	f          EventFunc
}

func NewCostItemEvent(player *player.Player, id, num uint32, curNum int64, remark string, src publicconst.ItemSource, notifyClient bool, cb EventFunc) *CostItemEvent {
	return &CostItemEvent{
		PlayerData: player,
		ItemId:     id,
		Num:        num,
		CurNum:     curNum,
		ItemSrc:    src,
		f:          cb,
		SendClient: notifyClient,
	}
}

func (c *CostItemEvent) RouteID() uint64 {
	return c.PlayerData.GetUserId()
}

func (c *CostItemEvent) CallBack() EventFunc {
	return c.f
}

// NotifyClientItemEvent 通知客户端道具事件
type NotifyClientItemEvent struct {
	IEvent
	PlayerData *player.Player
	ItemIds    []uint32
	f          EventFunc
}

func NewNotifyClientItemEvent(player *player.Player, itemIds []uint32, cb EventFunc) *NotifyClientItemEvent {
	return &NotifyClientItemEvent{
		PlayerData: player,
		ItemIds:    itemIds,
		f:          cb,
	}
}

func (c *NotifyClientItemEvent) RouteID() uint64 {
	return c.PlayerData.GetUserId()
}

func (c *NotifyClientItemEvent) CallBack() EventFunc {
	return c.f
}
