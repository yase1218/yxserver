package service

import (
	"gameserver/internal/game/common"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tools"
	"msg"
	"time"

	"github.com/zy/game_data/template"
	"google.golang.org/protobuf/proto"
)

// GetMonthCardMainScale 主线产出奖励 百分比
func GetMonthCardMainScale(p *player.Player) int {
	curTime := int(tools.GetCurTime())
	for i := 0; i < len(p.UserData.BaseInfo.MonthCard); i++ {
		config := template.GetMonthCardTemplate().GetMonthCard(p.UserData.BaseInfo.MonthCard[i].Id)
		if config.Data.Tpye == 2 && curTime <= p.UserData.BaseInfo.MonthCard[i].EndTime {
			return 100 + config.Data.MainlevelReward
		}
	}
	return 100
}

func CanLock(p *player.Player, conds []*template.UnlockCondition) msg.ErrCode {
	if !CanUnlock(p, conds) {
		return msg.ErrCode_FUNCTION_LOCK
	}
	return msg.ErrCode_SUCC
}

func GetMonthCardSweepTimes(p *player.Player) int {
	curTime := int(tools.GetCurTime())
	for i := 0; i < len(p.UserData.BaseInfo.MonthCard); i++ {
		config := template.GetMonthCardTemplate().GetMonthCard(p.UserData.BaseInfo.MonthCard[i].Id)
		if config.Data.Tpye == 2 && curTime <= p.UserData.BaseInfo.MonthCard[i].EndTime {
			return config.Data.QuestTimes
		}
	}
	return 0
}

func GetCreateRoleTimeRange(p *player.Player, startDay, endDay uint32) (uint32, uint32) {
	if startDay == 0 {
		startDay = 1
	}

	temp := time.Unix(int64(p.UserData.BaseInfo.CreateTime), 0)
	start := temp.AddDate(0, 0, int(startDay-1))
	end := temp.AddDate(0, 0, int(endDay-1))
	return uint32(tools.GetDateStart(start).Unix()), uint32(tools.GetDateEnd(end).Unix())
}

func GetOpenServerTimeRange(p *player.Player, startDay, endDay uint32) (uint32, uint32) {
	if startDay == 0 {
		startDay = 1
	}

	info := common.GetServerInfo()
	if info == nil {
		return 0, 0
	}

	temp := time.Unix(info.OpenTime, 0)
	start := temp.AddDate(0, 0, int(startDay-1))
	end := temp.AddDate(0, 0, int(endDay-1))
	return uint32(tools.GetDateStart(start).Unix()), uint32(tools.GetDateEnd(end).Unix())
}

func SendBannerMsg(bannerMsg *msg.InterNotifyBanner) {
	BoadCastMsg(bannerMsg)
}

func BoadCastMsg(m proto.Message) {
	users := player.AllPlayers()
	for _, v := range users {
		v.SendNotify(m)
	}
}

func GetMonthCardOnHookScale(p *player.Player) int {
	curTime := int(tools.GetCurTime())
	for i := 0; i < len(p.UserData.BaseInfo.MonthCard); i++ {
		config := template.GetMonthCardTemplate().GetMonthCard(p.UserData.BaseInfo.MonthCard[i].Id)
		if config.Data.Tpye == 2 && curTime <= p.UserData.BaseInfo.MonthCard[i].EndTime {
			return 100 + config.Data.OnHookReward
		}
	}
	return 100
}

func FunctionOpen(p *player.Player, id publicconst.FunctionId) msg.ErrCode {
	funcConfig := template.GetFunctionTemplate().GetFunction(uint32(id))
	if funcConfig == nil {
		return msg.ErrCode_FUNCTION_NOT_OPEN
	}

	if !CanUnlock(p, funcConfig.Conditions) {
		return msg.ErrCode_FUNCTION_LOCK
	}
	return msg.ErrCode_SUCC
}
