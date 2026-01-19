package event

import (
	"gameserver/internal/game/player"
)

// NotifyEquipChangeEvent 装备变化
type NotifyEquipChangeEvent struct {
	IEvent
	PlayerData *player.Player
	EquipIds   []uint32
	f          EventFunc
}

func NewNotifyEquipChangeEvent(p *player.Player, equipIds []uint32, cb EventFunc) *NotifyEquipChangeEvent {
	return &NotifyEquipChangeEvent{
		PlayerData: p,
		EquipIds:   equipIds,
		f:          cb,
	}
}

func (n *NotifyEquipChangeEvent) RouteID() uint64 {
	return n.PlayerData.GetUserId()
}

func (n *NotifyEquipChangeEvent) CallBack() EventFunc {
	return n.f
}
