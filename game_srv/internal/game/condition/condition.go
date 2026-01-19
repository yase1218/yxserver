package condition

import (
	"gameserver/internal/game/player"
	"msg"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

var (
	conditionMgr *Condition
)

type checkFn func(*player.Player, []uint32) ([]uint32, bool)

type Condition struct {
	checkers map[uint32]checkFn
}

func init() {
	conditionMgr = NewCondition()
	conditionMgr.Init()
}

func NewCondition() *Condition {
	return &Condition{
		checkers: make(map[uint32]checkFn),
	}
}

func GetCondition() *Condition {
	return conditionMgr
}

func (c *Condition) Init() {
	c.checkers[uint32(msg.ConditionType_Condition_Pass_Mission)] = checkMissionPass
	c.checkers[uint32(msg.ConditionType_Condition_Account_Days)] = initAccountDays
	c.checkers[uint32(msg.ConditionType_Condition_Open_Server_Days)] = initOpenServerDays
	c.checkers[uint32(msg.ConditionType_Condition_Pet)] = initPet
	c.checkers[uint32(msg.ConditionType_Condition_Contract_Rand)] = initContractRand
	c.checkers[uint32(msg.ConditionType_Condition_Contract_Kill_Monster)] = initContractKillMonster
	c.checkers[uint32(msg.ConditionType_Condition_Arena_Pk_Cnt)] = checkArenaPkCnt
}

func (c *Condition) Check(p *player.Player, args []uint32) ([]uint32, bool) {
	if len(args) == 0 {
		return []uint32{0}, true
	}

	if fn, ok := c.checkers[args[0]]; ok {
		return fn(p, args[1:])
	} else {
		log.Warn("condition not set", zap.Uint32s("args", args))
	}

	return []uint32{0}, false
}
