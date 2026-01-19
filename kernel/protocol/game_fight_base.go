package protocol

import (
	"encoding/json"
	"kernel/tda"
)

type GameFightBaseExtra struct {
	*tda.CommonAttr
	ChannelId uint32
	ServerId  uint32
}

func NewGameFightBaseExtra(commonAttr *tda.CommonAttr, channelId, serverId uint32) *GameFightBaseExtra {
	return &GameFightBaseExtra{
		CommonAttr: commonAttr,
		ChannelId:  channelId,
		ServerId:   serverId,
	}
}

func (e *GameFightBaseExtra) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *GameFightBaseExtra) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}
