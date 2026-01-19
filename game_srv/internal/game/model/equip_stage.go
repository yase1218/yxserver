package model

type UserEquipStage struct {
	RewardNum uint32 `bson:"reward_num"` // 当天已领奖次数
	RewardBuy uint32 `bson:"reward_buy"` // 购买次数

	Records map[uint32]*EquipStageRecord `bson:"records"` // key: 副本id, value: 挑战记录
}

type EquipStageRecord struct {
	Num uint32 `bson:"num"` // 挑战次数
}

func (u *UserEquipStage) RecordList() []uint32 {
	var list []uint32
	for k := range u.Records {
		list = append(list, k)
	}
	return list
}
