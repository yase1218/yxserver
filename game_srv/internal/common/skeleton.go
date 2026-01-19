package common

import (
	"gameserver/internal/config"
	"leaf/chanrpc"
	"leaf/module"
)

func NewSkeleton() *module.Skeleton {
	skeleton := &module.Skeleton{
		GoLen:              config.GoLen,
		TimerDispatcherLen: config.TimerDispatcherLen,
		AsynCallLen:        config.AsynCallLen,
		ChanRPCServer:      chanrpc.NewServer(config.ChanRPCLen),
	}
	skeleton.Init()
	return skeleton
}
