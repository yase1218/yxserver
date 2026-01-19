package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

func ClientSettings(packetId uint32, args interface{}, p *player.Player) {
	if p.UserData.BaseInfo.ClientSettings == nil {
		p.UserData.BaseInfo.ClientSettings = make(map[uint32]string)
		p.SaveBaseInfo()
	}
	res := &msg.ClientSettingsAck{
		Settings: p.UserData.BaseInfo.ClientSettings,
	}
	p.SendNotify(res)
}

func ClientSettingsUpdate(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.ClientSettingsUpdateReq)
	if len(req.GetSettings()) > 0 {
		if p.UserData.BaseInfo.ClientSettings == nil {
			p.UserData.BaseInfo.ClientSettings = make(map[uint32]string)
		}
		for _, kv := range req.GetSettings() {
			p.UserData.BaseInfo.ClientSettings[kv.Key] = kv.Value
		}
		p.SaveBaseInfo()
	}

	res := &msg.ClientSettingsUpdateAck{
		Settings: p.UserData.BaseInfo.ClientSettings,
	}
	p.SendNotify(res)
}

func ClientSettingsDelete(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.ClientSettingsDeleteReq)
	if len(req.GetKeys()) > 0 {
		if p.UserData.BaseInfo.ClientSettings == nil {
			p.UserData.BaseInfo.ClientSettings = make(map[uint32]string)
		}
		for _, v := range req.GetKeys() {
			delete(p.UserData.BaseInfo.ClientSettings, v)
		}
		p.SaveBaseInfo()
	}

	res := &msg.ClientSettingsDeleteAck{
		Keys: req.GetKeys(),
	}
	p.SendNotify(res)
}

func RejectReconnect(packetId uint32, args interface{}, p *player.Player) {
	service.LeaveFight(p)
}
