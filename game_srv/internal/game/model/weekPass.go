package model

type WeekPass struct {
	ContractInfo   map[uint32]bool // 通关状态
	SecretCount    int32           // 秘境挑战次数
	SecretBoxState uint32          // 秘境任务状况
}
