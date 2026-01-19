package main

import (
	"flag"
	"gameserver/internal/global"
	_ "net/http/pprof"
)

var (
	serverId   = flag.Uint("serverid", 0, "server id")
	serverName = flag.String("servername", "", "server name")
	localDb    = flag.String("localdb", "", "local db")
)

func main() {
	flag.Parse()
	global.Init(uint32(*serverId),*serverName,*localDb)
	global.Run()
}
