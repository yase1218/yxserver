package protocol

import (
	"encoding/json"
)

// 僵尸竞技场 被挑战的人信息
type ZombieColiseumExtra struct {
	Id      uint32         // accountId | robotId
	IsRobot bool           // 是否机器人
	List    map[int]uint32 // 波次:阵容 从1开始
	Combat  uint32         // 战力(防守方)
}

func (e *ZombieColiseumExtra) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *ZombieColiseumExtra) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}
