package event

import "gameserver/internal/game/player"

type MailEvent struct {
	PlayerData *player.Player
	MailId     int64
	IsDelete   bool
	f          EventFunc
}

func NewMailEvent(playerData *player.Player, mailId int64, isDelete bool, cb EventFunc) *MailEvent {
	return &MailEvent{
		PlayerData: playerData,
		MailId:     mailId,
		IsDelete:   isDelete,
		f:          cb,
	}
}

func (a *MailEvent) RouteID() uint64 {
	return a.PlayerData.GetUserId()
}

func (a *MailEvent) CallBack() EventFunc {
	return a.f
}
