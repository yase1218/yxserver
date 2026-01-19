package protocol

import (
	"encoding/json"
)

type PeakFightExtra struct {
	StageBuff []uint32 // 战场环境
}

func (e *PeakFightExtra) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *PeakFightExtra) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}
