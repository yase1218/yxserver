package service

import (
	"gameserver/internal/game/player"
	"msg"

	"github.com/zy/game_data/template"
)

func EnterBlackBoss(p *player.Player, btType uint32) (msg.ErrCode, uint32, *msg.WeekPassInfo) {
	playMethodInfo := p.UserData.PlayMethod
	var data = &msg.WeekPassInfo{}
	for _, v := range playMethodInfo.Data {
		if v.BtType == int(btType) {
			data.Type = uint32(v.BtType)
			data.MaxDamage = uint32(v.MaxDamage)
			data.TotalTimes = uint32(v.TotalTimes)
			break
		}
	}
	rank := getDesertRank(p)
	return msg.ErrCode_SUCC, rank, data
}

func getDesertRank(p *player.Player) uint32 {
	var rank uint32 = 0
	rankInfo, _, _ := GetRankData(p, template.WeekPassBlackBoosDamage)
	for k, v := range rankInfo {
		for _, j := range v.PlayerInfo {
			if j.AccountId == p.UserData.UserId {
				rank = uint32(k) + 1
				break
			}
		}

		if rank != 0 {
			break
		}
	}

	return rank
}

func OnCrossDayFreshWeekPass(p *player.Player) {
	p.UserData.WeekPass.ContractInfo = make(map[uint32]bool)
	p.UserData.WeekPass.SecretBoxState = 0
	p.UserData.WeekPass.SecretCount = 0
	p.SaveWeekPass()
}
