package event

import (
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
)

// AddShipEvent 添加机甲
type AddShipEvent struct {
	IEvent
	PlayerData   *player.Player
	ShipIds      []uint32
	NotifyClient bool
	f            EventFunc
}

func NewAddShipEvent(player *player.Player, shipIds []uint32, notifyClient bool, cb EventFunc) *AddShipEvent {
	return &AddShipEvent{
		PlayerData:   player,
		ShipIds:      shipIds,
		NotifyClient: notifyClient,
		f:            cb,
	}
}

func (a *AddShipEvent) RouteID() uint64 {
	return a.PlayerData.GetUserId()
}

func (a *AddShipEvent) CallBack() EventFunc {
	return a.f
}

// ShipChangeEvent 机甲变化
type ShipChangeEvent struct {
	IEvent
	PlayerData *player.Player
	ShipId     uint32
	OldData    uint32
	NewData    uint32
	AttrType   publicconst.ShipAttrType
	f          EventFunc
}

func NewShipChangeEvent(player *player.Player, shipId uint32, oldData, newData uint32, attrType publicconst.ShipAttrType, cb EventFunc) *ShipChangeEvent {
	return &ShipChangeEvent{
		PlayerData: player,
		ShipId:     shipId,
		OldData:    oldData,
		NewData:    newData,
		AttrType:   attrType,
		f:          cb,
	}
}

func (s *ShipChangeEvent) RouteID() uint64 {
	return s.PlayerData.GetUserId()
}

func (s *ShipChangeEvent) CallBack() EventFunc {
	return s.f
}

// NotifyClientShipChangeEvent 机甲变化
type NotifyClientShipChangeEvent struct {
	IEvent
	PlayerData *player.Player
	ShipIds    []uint32
	f          EventFunc
}

func NewNotifyClientShipChangeEvent(player *player.Player, shipIds []uint32, cb EventFunc) *NotifyClientShipChangeEvent {
	return &NotifyClientShipChangeEvent{
		PlayerData: player,
		ShipIds:    shipIds,
		f:          cb,
	}
}

func (n *NotifyClientShipChangeEvent) RouteID() uint64 {
	return n.PlayerData.GetUserId()
}

func (n *NotifyClientShipChangeEvent) CallBack() EventFunc {
	return n.f
}
