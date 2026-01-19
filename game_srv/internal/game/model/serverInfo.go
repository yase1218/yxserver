package model

import (
	"time"
)

type ServerInfo struct {
	InitServerId uint32
	ServerId     uint32
	ServerName   string
	RegistNum    int64
	UpdateTime   uint32
	CreateTime   uint32
	RegistLimit  uint32
	OnlineLimit  uint32
	OnlineFull   uint32
	OnlineNum    uint32
	DisplayTime  uint32
	Opentime     uint32
	ServerAddr   string
}

type ServerConfig struct {
	Id            uint32
	DefaultServer uint32
}

func NewServerInfo(serverId uint32) *ServerInfo {
	curTime := uint32(time.Now().Unix())
	return &ServerInfo{
		InitServerId: serverId,
		ServerId:     serverId,
		UpdateTime:   curTime,
		CreateTime:   curTime,
		RegistLimit:  5000,
		OnlineLimit:  5000,
		OnlineFull:   3000,
	}
}
