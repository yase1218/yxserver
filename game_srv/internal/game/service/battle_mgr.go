package service

import (
	"kernel/kenum"
	"time"
)

type BattleCache struct {
	FsId     uint32 // 战斗服id
	BattleId uint32 // 战斗id
	Uid      uint64 // 角色id
	DeadLine time.Time
	StageId  uint32 // 关卡ID
}

type BattleCacheManager struct {
	BcByUser map[uint64]*BattleCache // 角色id:战斗缓存
}

var (
	bcManagerIns    *BattleCacheManager
	defaultDeadline = time.Now().Add(time.Minute * time.Duration(kenum.FightDeadLine))
)

func init() {
	bcManagerIns = &BattleCacheManager{
		BcByUser: make(map[uint64]*BattleCache),
	}
}

func AddBattleCache(bc *BattleCache) {
	bcManagerIns.BcByUser[bc.Uid] = bc
}

func RemoveBattleCacheByUid(uid uint64) {
	delete(bcManagerIns.BcByUser, uid)
}

func GetBattleCacheByUid(uid uint64) *BattleCache {
	if _, ok := bcManagerIns.BcByUser[uid]; !ok {
		newBc := &BattleCache{
			FsId:     0,
			BattleId: 0,
			Uid:      uid,
			StageId:  0,
			DeadLine: defaultDeadline,
		}

		bcManagerIns.BcByUser[uid] = newBc
	}

	return bcManagerIns.BcByUser[uid]
}

func SetBattleCacheFsId(uid uint64, fsId uint32) {
	bc := GetBattleCacheByUid(uid)
	bc.FsId = fsId
}
