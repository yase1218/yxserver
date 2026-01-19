package service

import (
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"time"

	"github.com/v587-zyf/gc/log"
)

func LoginCorrect(p *player.Player, now time.Time) bool {
	return correct_test(p, now) ||
		correct_ship(p, now)
}

func correct_test(p *player.Player, now time.Time) bool {
	return false
}

func correct_ship(p *player.Player, now time.Time) bool {
	if !p.UserData.IsRegister { // 已经完成创建角色改名和机甲
		if p.UserData.BaseInfo.ShipId == 0 { // 出站机甲未入库容错
			if len(p.UserData.Ships.Ships) == 0 { // 机甲背包空
				log.Error("!!!player ships empty after register", ZapUser(p))
				p.UserData.IsRegister = true // 强制玩家再次完成创建角色改名和机甲
				p.SaveBaseInfo()
			} else { // 机甲背包有数据
				p.UserData.BaseInfo.ShipId = p.UserData.Ships.Ships[0].Id
				p.SaveBaseInfo()
				return true
			}
		} else { // 机甲背包未入库容错
			if len(p.UserData.Ships.Ships) == 0 {
				AddShip(p, p.UserData.BaseInfo.ShipId, publicconst.InitAddItem, false)
			}
		}
	}
	return false
}
