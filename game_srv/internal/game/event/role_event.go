package event

import (
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
)

// AddRoleEvent 添加机甲
type AddRoleEvent struct {
	IEvent
	PlayerData   *player.Player
	RoleIds      []uint32
	NotifyClient bool
	f            EventFunc
}

func NewAddRoleEvent(player *player.Player, roleIds []uint32, notifyClient bool, cb EventFunc) *AddRoleEvent {
	return &AddRoleEvent{
		PlayerData:   player,
		RoleIds:      roleIds,
		NotifyClient: notifyClient,
		f:            cb,
	}
}

func (a *AddRoleEvent) RouteID() uint64 {
	return a.PlayerData.GetUserId()
}

func (a *AddRoleEvent) CallBack() EventFunc {
	return a.f
}

// RoleChangeEvent 驾驶员变化
type RoleChangeEvent struct {
	IEvent
	PlayerData *player.Player
	RoleId     uint32
	OldData    uint32
	NewData    uint32
	AttrType   publicconst.RoleAttrType
	f          EventFunc
}

func NewRoleChangeEvent(player *player.Player, roleId uint32, oldData, newData uint32, attrType publicconst.RoleAttrType, cb EventFunc) *RoleChangeEvent {
	return &RoleChangeEvent{
		PlayerData: player,
		RoleId:     roleId,
		OldData:    oldData,
		NewData:    newData,
		AttrType:   attrType,
		f:          cb,
	}
}

func (r *RoleChangeEvent) RouteID() uint64 {
	return r.PlayerData.GetUserId()
}

func (r *RoleChangeEvent) CallBack() EventFunc {
	return r.f
}

// NotifyClientRoleChangeEvent 驾驶员变化
type NotifyClientRoleChangeEvent struct {
	IEvent
	PlayerData *player.Player
	RoleIds    []uint32
	f          EventFunc
}

func NewNotifyClientRoleChangeEvent(player *player.Player, roleIds []uint32, cb EventFunc) *NotifyClientRoleChangeEvent {
	return &NotifyClientRoleChangeEvent{
		PlayerData: player,
		RoleIds:    roleIds,
		f:          cb,
	}
}

func (n *NotifyClientRoleChangeEvent) RouteID() uint64 {
	return n.PlayerData.GetUserId()
}

func (n *NotifyClientRoleChangeEvent) CallBack() EventFunc {
	return n.f
}
