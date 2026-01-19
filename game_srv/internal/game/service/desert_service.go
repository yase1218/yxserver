package service

import (
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/errCode"
	"kernel/tools"
	"msg"

	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

func AddKillTimes(p *player.Player, killNum uint32) {
	p.UserData.Desert.KillTimes += killNum
	p.SendNotify(&msg.DesertKillTimesNtf{
		KillTimes: p.UserData.Desert.KillTimes,
	})
}

func DesertLampReward(p *player.Player, req *msg.DesertLampRewardReq, ack *msg.DesertLampRewardAck) error {
	rewardsCfg, ok := template.GetSystemItemTemplate().DesertLampReward[req.GetRewardTimes()]
	if !ok {
		log.Error("desert reward nil", zap.Uint32("times", req.GetRewardTimes()))
		return errcode.ERR_CONFIG_NIL
	}

	for _, times := range p.UserData.Desert.RewardTimes {
		if times == req.GetRewardTimes() {
			log.Error("desert repeated reward", zap.Uint64("uid", p.GetUserId()),
				zap.Uint32("times", req.GetRewardTimes()))
			return errCode.ERR_REPEATE_REWARD
		}
	}

	p.UserData.Desert.RewardTimes = append(p.UserData.Desert.RewardTimes, req.GetRewardTimes())

	var notifyItems []uint32
	var retItems []*msg.SimpleItem
	if len(rewardsCfg) > 0 {
		for _, item := range rewardsCfg {
			addItems := AddItem(p.GetUserId(), item.ItemId, int32(item.ItemNum), publicconst.PeakFight, false)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
			retItems = append(retItems, &msg.SimpleItem{ItemId: item.ItemId, ItemNum: item.ItemNum})
		}
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}
	if len(retItems) > 0 {
		ack.RewardItem = retItems
	}
	return nil
}

func DayResetDesert(p *player.Player, sendNtf bool) {
	resetTime := tools.GetDailyRefreshTime()
	if p.UserData.Desert.ResetDate != resetTime {
		p.UserData.Desert.KillTimes = 0
		p.UserData.Desert.RewardTimes = make([]uint32, 0)
		p.UserData.Desert.ResetDate = resetTime

		p.SaveDesert()
		if sendNtf {
			p.SendNotify(&msg.DesertResetNtf{
				RewardTimes: p.UserData.Desert.RewardTimes,
				KillTimes:   p.UserData.Desert.KillTimes,
			})
		}
	}
}
