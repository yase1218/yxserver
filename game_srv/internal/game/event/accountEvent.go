package event

import "gameserver/internal/game/player"

type LevelChangeEvent struct {
	PlayerData *player.Player
	OldLevel   uint32
	NewLevel   uint32
	f          EventFunc
}

func NewLevelChangeEvent(playerData *player.Player, oldLevel, newLevel uint32, cb EventFunc) *LevelChangeEvent {
	return &LevelChangeEvent{
		PlayerData: playerData,
		OldLevel:   oldLevel,
		NewLevel:   newLevel,
		f:          cb,
	}
}

func (a *LevelChangeEvent) RouteID() uint64 {
	return a.PlayerData.GetUserId()
}

func (a *LevelChangeEvent) CallBack() EventFunc {
	return a.f
}
