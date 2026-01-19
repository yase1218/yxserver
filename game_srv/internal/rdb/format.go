package rdb

import (
	"fmt"
	"gameserver/internal/config"
	"gameserver/internal/enum"
	"kernel/tools"
	"time"

	"github.com/zy/game_data/template"
)

// FormatAllianceRankSinglePower 成员战力榜（联盟内）
func FormatAllianceRankSinglePower(allianceId uint32) string {
	return fmt.Sprintf("%s:%d:%d", enum.Redis_Key_Alliance_Rank_Single_Power, config.Conf.ServerId, allianceId)
}

// FormatAllianceRankSingleActive 成员战力榜（活跃度）
func FormatAllianceRankSingleActive(allianceId uint32, d ...int) string {
	var date int
	if len(d) > 0 {
		date = d[0]
	} else {
		date = tools.GetYearWeekByOffset(time.Now(), int(template.GetSystemItemTemplate().RefreshHour))
	}

	return fmt.Sprintf("%s:%d:%d:%d",
		enum.Redis_Key_Alliance_Rank_Single_Active,
		config.Conf.ServerId,
		date,
		allianceId)
}

// FormatAllianceBossSingleRank 联盟BOSS伤害榜(个人)
func FormatAllianceBossSingleRank(allianceId uint32, d ...int) string {
	var date int
	if len(d) > 0 {
		date = d[0]
	} else {
		date = tools.GetYearWeekByOffset(time.Now(), int(template.GetSystemItemTemplate().RefreshHour))
	}

	return fmt.Sprintf("%s:%d:%d:%d",
		enum.Redis_Key_Alliance_Boss_Single_Rank,
		config.Conf.ServerId,
		date,
		allianceId)
}

// FormatAllianceBossRank 联盟BOSS伤害榜(联盟)
func FormatAllianceBossRank(d ...int) string {
	var date int
	if len(d) > 0 {
		date = d[0]
	} else {
		date = tools.GetYearWeekByOffset(time.Now(), int(template.GetSystemItemTemplate().RefreshHour))
	}

	return fmt.Sprintf("%s:%d:%d",
		enum.Redis_Key_Alliance_Boss_Rank,
		config.Conf.ServerId,
		date)
}

func FormatPeakFightRank(season uint32) string {
	return fmt.Sprintf("%s:%d:%d", enum.Redis_Key_PeakFight_Rank, config.Conf.ServerId, season)
}
