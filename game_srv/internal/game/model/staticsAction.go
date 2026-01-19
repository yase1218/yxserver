package model

import (
	"fmt"
	"kernel/tools"
)

type StaticsAction struct {
	ActionId   uint32 `json:"actionId"`
	AccountId  string `json:"accountId"` // 设备号
	ChannelId  uint32 `json:"channelId"` // 渠道id
	Paras      string `json:"paras"`     // 参数
	Ip         string `json:"ip"`
	UpdateTime uint32 `json:"updateTime"`
}

func NewStaticsAction(actionId uint32, accountId, channelId uint32, para, ip string) *StaticsAction {
	return &StaticsAction{
		ActionId:   actionId,
		AccountId:  fmt.Sprintf("%v", accountId),
		ChannelId:  channelId,
		Paras:      para,
		Ip:         ip,
		UpdateTime: tools.GetCurTime(),
	}
}
