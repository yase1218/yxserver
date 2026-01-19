package model

type AccountCardPool struct {
	AccountId int64
	CardPools []*CardPool
}

type CardPool struct {
	CardId            uint32
	FreeTimes         uint32
	NextFreeTime      uint32
	TenQuaranteeTimes uint32 // 10连抽次数
	BigQuaranteeTimes uint32 // 大保底次数
	TotalTimes        uint32 // 总次数
	NextResetTime     uint32 // 下一次重置次数时间
	StartTime         uint32 // 开始时间
	EndTime           uint32 // 结束时间
	LotteryEndTime    uint32
	Status            uint32 // 是否 有效
	FirstLottery      uint32
	FirstLotteryTen   uint32
	LotteryTotalTimes uint32 // 抽奖总次数
}

func NewAccountCardPool(accountId int64) *AccountCardPool {
	ret := &AccountCardPool{
		AccountId: accountId,
	}
	ret.CardPools = make([]*CardPool, 0, 0)
	return ret
}

func NewCardPool(id uint32, freeTimes, nextFreeTime, nextResetTime, startTime, endTime, lotteryEndTime uint32) *CardPool {
	return &CardPool{
		CardId:         id,
		FreeTimes:      freeTimes,
		NextFreeTime:   nextFreeTime,
		NextResetTime:  nextResetTime,
		StartTime:      startTime,
		EndTime:        endTime,
		Status:         1,
		LotteryEndTime: lotteryEndTime,
	}
}
