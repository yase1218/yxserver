package player

import (
	"sync"
	"time"
)

var (
	mgr *PlayerMgr
)

type PlayerMgr struct {
	//sync.RWMutex
	Account map[string]*Player
	UserId  map[uint64]*Player
	Nick    map[string]*Player
	Pool    sync.Pool
}

func init() {
	mgr = &PlayerMgr{
		Account: make(map[string]*Player),
		UserId:  make(map[uint64]*Player),
		Nick:    make(map[string]*Player),
	}
	mgr.Pool.New = func() interface{} {
		return &Player{}
	}
}

func CreatePlayer(now time.Time) *Player {
	p := mgr.Pool.Get().(*Player)
	p.Init(now)
	return p
}

func AddPlayer(p *Player) {
	// mgr.Lock()
	// defer mgr.Unlock()
	mgr.Account[p.GetOpenId()] = p
	mgr.UserId[p.GetUserId()] = p
	mgr.Nick[p.UserData.Nick] = p
}

func FindByAccount(accountId string) *Player {
	// mgr.RLock()
	// defer mgr.RUnlock()
	return mgr.Account[accountId]
}

func FindByUserId(userId uint64) *Player {
	// mgr.RLock()
	// defer mgr.RUnlock()
	return mgr.UserId[userId]
}

func FindByNick(nick string) *Player {
	// mgr.RLock()
	// defer mgr.RUnlock()
	return mgr.Nick[nick]
}

func UpdateNick(oldNick, newNick string) {
	// mgr.Lock()
	// defer mgr.Unlock()
	p := FindByNick(oldNick)
	if p == nil {
		return
	}
	delete(mgr.Nick, oldNick)
	mgr.Nick[newNick] = p
}

func DelPlayer(p *Player) {
	if p == nil {
		return
	}
	// mgr.Lock()
	// defer mgr.Unlock()

	if p.GetOpenId() != "" {
		delete(mgr.Account, p.GetOpenId())
	}
	if p.GetUserId() != 0 {
		delete(mgr.UserId, p.GetUserId())
	}
	mgr.Pool.Put(p)
}

func AllPlayers() map[uint64]*Player {
	// mgr.RLock()
	// defer mgr.RUnlock()
	return mgr.UserId
}

func OnlineNum() uint32 {
	// mgr.RLock()
	// defer mgr.RUnlock()
	return uint32(len(mgr.UserId))
}

func Stop() {
	// mgr.RLock()
	// defer mgr.RUnlock()
	for _, p := range AllPlayers() {
		p.SaveDirtySync()
	}
}
