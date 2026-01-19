package model

import "kernel/tools"

type GlobalInfo struct {
	ServerId             uint32
	NextDesertSettleTime uint32
}

func NewGlobalInfo(serverId uint32) *GlobalInfo {
	return &GlobalInfo{
		ServerId:             serverId,
		NextDesertSettleTime: tools.GetWeeklyRefreshTime(0),
	}
}
