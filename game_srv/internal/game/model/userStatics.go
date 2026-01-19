package model

import "kernel/tools"

type UserStatics struct {
	AccountId   int64
	OnlineTime  uint32
	StaticsTime uint32
	ChannelId   uint32
}

func NewUserStatics(accountId int64, staticsTime, onlineTime, channelId uint32) *UserStatics {
	return &UserStatics{
		AccountId:   accountId,
		StaticsTime: staticsTime,
		OnlineTime:  onlineTime,
		ChannelId:   channelId,
	}
}

type Action struct {
	Id         uint32
	Params     string
	UpdateTime uint32
}

func NewAction(id uint32, params string) *Action {
	return &Action{
		Id:         id,
		Params:     params,
		UpdateTime: tools.GetCurTime(),
	}
}

type UserDailyAction struct {
	AccountId   int64
	StaticsTime uint32
	Actions     []*Action
}

func NewUserDailyAction(accountId int64, staticTime uint32) *UserDailyAction {
	return &UserDailyAction{
		AccountId:   accountId,
		StaticsTime: staticTime,
	}
}

type DNUStatics struct {
	Date       uint32
	ChannelId  uint32
	AccountId  uint32
	CreateTime uint32
	Ip         string
	ExtraInfo  string
	UserId     string
}

func NewDNUStatics(d, channelId, accountId uint32, createTime uint32, ip, extra, userId string) *DNUStatics {
	return &DNUStatics{
		Date:       d,
		ChannelId:  channelId,
		AccountId:  accountId,
		CreateTime: createTime,
		Ip:         ip,
		ExtraInfo:  extra,
		UserId:     userId,
	}
}

type OnlineCount struct {
	Date     string
	ServerId uint32
	Minute   uint32
	Data     []*ChannelOnlineCount
}

func NewOnlineCount(curDate string, serverId, minute uint32, data []*ChannelOnlineCount) *OnlineCount {
	return &OnlineCount{
		Date:     curDate,
		ServerId: serverId,
		Minute:   minute,
		Data:     data,
	}
}

type ChannelOnlineCount struct {
	ChannelId uint32
	Count     uint32
}

func NewChannelOnlineCount(chId, count uint32) *ChannelOnlineCount {
	return &ChannelOnlineCount{
		ChannelId: chId,
		Count:     count,
	}
}

type GuideStep struct {
	Date      uint32
	AccountId uint32
	ChannelId uint32
	Step      uint32
	SubStep   uint32
}

func NewGuideStep(d, accountId, chId, step, subStep uint32) *GuideStep {
	return &GuideStep{
		d,
		accountId,
		chId,
		step,
		subStep,
	}
}

type LossLevel struct {
	Date      uint32
	AccountId uint32
	Level     uint32
	ChannelId uint32
}

func NewLossLevel(d, accountId, level, chId uint32) *LossLevel {
	return &LossLevel{
		d,
		accountId,
		level,
		chId,
	}
}

type LossMission struct {
	Date      uint32
	AccountId uint32
	MissionId int
}

func NewLossMission(d, accountId uint32, missionId int) *LossMission {
	return &LossMission{
		d,
		accountId,
		missionId,
	}
}

type CdkRecord struct {
	Cdk        string
	AccountId  int64
	Nick       string
	UpdateTime uint32
}

func NewCdkRecord(cdk string, accountId int64, nick string) *CdkRecord {
	return &CdkRecord{
		Cdk:        cdk,
		AccountId:  accountId,
		Nick:       nick,
		UpdateTime: tools.GetCurTime(),
	}
}
