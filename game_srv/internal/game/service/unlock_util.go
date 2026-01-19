package service

import (
	"gameserver/internal/game/common"
	"kernel/tools"
	"time"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"

	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
)

type UnlockFunc func(*player.Player, []uint32) bool

var (
	unlockFuncMap map[publicconst.UnlockType]UnlockFunc
)

func init() {
	unlockFuncMap = make(map[publicconst.UnlockType]UnlockFunc)
	unlockFuncMap[publicconst.Unlock_Mission] = unlockMission
	unlockFuncMap[publicconst.Unlock_Player_Level] = unlockPlayerLevel
	unlockFuncMap[publicconst.Unlock_Time] = unlockTime
	unlockFuncMap[publicconst.Unlock_Create_Role_Time] = unlockCreateRoleTime
	unlockFuncMap[publicconst.Unlock_Shop_Item] = isUnlockShopItem
	//unlockFuncMap[publicconst.Unlock_Challenge_Mission] = unlockChallengeMission
	unlockFuncMap[publicconst.Unlock_Challenge_Mission] = unlockChallengeMission
	unlockFuncMap[publicconst.Unlock_Open_Server_Time] = unlockOpenServerTime
	unlockFuncMap[publicconst.Unlock_Pet_Id] = unlockPetFunction
}

func CanUnlockOneCond(p *player.Player, cond *template.UnlockCondition) bool {
	if f, ok := unlockFuncMap[publicconst.UnlockType(cond.Id)]; ok {
		return f(p, cond.Values)
	}
	return true
}

func CanUnlock(p *player.Player, conds []*template.UnlockCondition) bool {
	for i := 0; i < len(conds); i++ {
		if f, ok := unlockFuncMap[publicconst.UnlockType(conds[i].Id)]; ok {
			if !f(p, conds[i].Values) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func unlockMission(p *player.Player, values []uint32) bool {
	if len(values) == 0 {
		return false
	}
	if p.UserData.Mission == nil {
		log.Error("mission data nil", zap.Uint64("accountId", p.GetUserId()))
		return false
	}
	if mission := findMission(p, int(values[0]), true); mission != nil {
		if mission.IsPass {
			return true
		}
	}

	if eq_stage := FindEquipStage(p, values[0]); eq_stage != nil {
		return true
	}

	return false
}

func unlockPlayerLevel(p *player.Player, values []uint32) bool {
	if len(values) == 0 {
		return false
	}
	if p.UserData.Level >= values[0] {
		return true
	}
	return false
}

func unlockTime(p *player.Player, values []uint32) bool {
	if len(values) == 0 {
		return false
	}
	curTime := tools.GetCurTime()
	if curTime >= values[0] && curTime <= values[1] {
		return true
	}
	return false
}

func unlockCreateRoleTime(p *player.Player, values []uint32) bool {
	days := tools.GetDiffDay(time.Unix(int64(p.UserData.BaseInfo.CreateTime), 0), time.Now()) + 1
	if len(values) == 1 {
		if days >= values[0] {
			return true
		}
	} else if len(values) == 2 {
		if days >= values[0] && days <= values[1] {
			return true
		}
	}
	return false
}

func unlockOpenServerTime(p *player.Player, values []uint32) bool {
	info := common.GetServerInfo()
	if info == nil || info.OpenTime == 0 {
		return false
	}

	days := tools.GetDiffDay(time.Unix(info.OpenTime, 0), time.Now()) + 1
	if len(values) == 1 {
		if days >= values[0] {
			return true
		}
	} else if len(values) == 2 {
		if days >= values[0] && days <= values[1] {
			return true
		}
	}
	return false
}

func isUnlockShopItem(p *player.Player, values []uint32) bool {
	if len(values) == 0 {
		return false
	}
	if item := getShopItem(p, values[0]); item != nil {
		itemConfig := template.GetShopTemplate().GetShopItem(values[0])
		if itemConfig != nil && item.BuyTimes >= itemConfig.LimitNum {
			return true
		}
	}
	return false
}

func unlockChallengeMission(p *player.Player, values []uint32) bool {
	if len(values) == 0 {
		return false
	}
	if p.UserData.PlayMethod == nil {
		log.Error("play method data nil", zap.Uint64("accountId", p.GetUserId()))
		return false
	}
	//if mission := getMission(p, int(msg.BattleType_Battle_Challenge), int(values[0])); mission != nil && mission.IsPass {
	//	return true
	//}
	if mission := findMission(p, int(values[0]), false); mission != nil {
		if mission.IsPass {
			return true
		}
	}
	return false
}

func unlockPetFunction(p *player.Player, values []uint32) bool {
	if len(values) == 0 {
		return false
	}
	if p.UserData.PetData == nil {
		log.Error("pet data nil", zap.Uint64("accountId", p.GetUserId()))
		return false
	}
	return p.UserData.BaseInfo.UsePet >= values[0]
}
