package event

import (
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
)

// AddWeaponEvent 添加武器
type AddWeaponEvent struct {
	IEvent
	PlayerData *player.Player
	WeaponIds  []uint32
	Source     publicconst.ItemSource
	f          EventFunc
}

func NewAddWeaponEvent(player *player.Player, weaponIds []uint32, source publicconst.ItemSource, cb EventFunc) *AddWeaponEvent {
	return &AddWeaponEvent{
		PlayerData: player,
		WeaponIds:  weaponIds,
		f:          cb,
		Source:     source,
	}
}

func (a *AddWeaponEvent) RouteID() uint64 {
	return a.PlayerData.GetUserId()
}

func (a *AddWeaponEvent) CallBack() EventFunc {
	return a.f
}

// WeaponUpgradeEvent 武器升级
type WeaponUpgradeEvent struct {
	IEvent
	PlayerData *player.Player
	WeaponId   uint32
	NewLevel   uint32
	f          EventFunc
}

func NewWeaponUpgradeEvent(player *player.Player, weaponId, level uint32, cb EventFunc) *WeaponUpgradeEvent {
	return &WeaponUpgradeEvent{
		PlayerData: player,
		WeaponId:   weaponId,
		f:          cb,
		NewLevel:   level,
	}
}

func (w *WeaponUpgradeEvent) RouteID() uint64 {
	return w.PlayerData.GetUserId()
}

func (w *WeaponUpgradeEvent) CallBack() EventFunc {
	return w.f
}

// WeaponLibUpgradeEvent
type WeaponLibUpgradeEvent struct {
	IEvent
	PlayerData *player.Player
	NewLevel   uint32
	f          EventFunc
}

func NewWeaponLibUpgradeEvent(player *player.Player, level uint32, cb EventFunc) *WeaponLibUpgradeEvent {
	return &WeaponLibUpgradeEvent{
		PlayerData: player,
		f:          cb,
		NewLevel:   level,
	}
}

func (w *WeaponLibUpgradeEvent) RouteID() uint64 {
	return w.PlayerData.GetUserId()
}

func (w *WeaponLibUpgradeEvent) CallBack() EventFunc {
	return w.f
}
