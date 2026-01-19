package event

import "gameserver/internal/game/player"

// PassMissionEvent 通关事件
type PassMissionEvent struct {
	IEvent
	PlayerData *player.Player
	MissionId  int
	f          EventFunc
}

func NewPassMissionEvent(player *player.Player, missionId int, cb EventFunc) *PassMissionEvent {
	return &PassMissionEvent{
		PlayerData: player,
		MissionId:  missionId,
		f:          cb,
	}
}

func (a *PassMissionEvent) RouteID() uint64 {
	return a.PlayerData.GetUserId()
}

func (a *PassMissionEvent) CallBack() EventFunc {
	return a.f
}

// FirstPassMissionEvent 首个通关的玩家
type FirstPassMissionEvent struct {
	IEvent
	PlayerData *player.Player
	Args       interface{}
	f          EventFunc
}

func NewFirstPassMissionEvent(player *player.Player, args interface{}, cb EventFunc) *FirstPassMissionEvent {
	return &FirstPassMissionEvent{
		PlayerData: player,
		Args:       args,
		f:          cb,
	}
}

func (a *FirstPassMissionEvent) RouteID() uint64 {
	return a.PlayerData.GetUserId()
}

func (a *FirstPassMissionEvent) CallBack() EventFunc {
	return a.f
}
