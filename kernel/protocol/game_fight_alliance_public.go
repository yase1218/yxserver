package protocol

import "encoding/json"

type AlliancePublicExtra struct {
	AllianceId uint32
}

func (e *AlliancePublicExtra) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *AlliancePublicExtra) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}
