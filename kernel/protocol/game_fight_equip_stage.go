package protocol

import "encoding/json"

type EquipStageFightExtra struct{}

func (e *EquipStageFightExtra) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *EquipStageFightExtra) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

type EquipStageFightResultExtra struct {
	Items map[uint32]uint32 // 物品id: 数量
}

func (e *EquipStageFightResultExtra) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *EquipStageFightResultExtra) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}
