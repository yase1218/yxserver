package tda

import (
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

type RankUnit struct {
	AccountId int64  `json:"role_id"` // 账号
	Ranking   uint32 `json:"ranking"` // 排名
	Score     uint64 `json:"score"`   // 分数
}

// 主線關卡排行榜	主線關卡排行榜變動時上報
type RankChangeMainBattle struct {
	List []*RankUnit `json:"rank_list"` // 排行列表
}

func TdaRankChangeMainBattle(list []*RankUnit) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaRankChangeMainBattle", func() {
		tdaData := &RankChangeMainBattle{
			List: list,
		}

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaRankChangeMainBattle tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track("0", "0", Tda_Rank_Change_Main_Battle, tdaDataMap); err != nil {
				log.Error("tda TdaRankChangeMainBattle track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(0, "0", "0", Tda_Rank_Change_Main_Battle, tdaDataMap)
		}
	})
}

// 主線菁英關卡排行榜	主線菁英關卡排行榜變動時上報
type RankChangeEliteMainBattle struct {
	List []*RankUnit `json:"rank_list"` // 排行列表
}

func TdaRankChangeEliteMainBattle(list []*RankUnit) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaRankChangeEliteMainBattle", func() {
		tdaData := &RankChangeMainBattle{
			List: list,
		}

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaRankChangeEliteMainBattle tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track("0", "0", Tda_Rank_Change_Elite_Main_Battle, tdaDataMap); err != nil {
				log.Error("tda TdaRankChangeEliteMainBattle track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(0, "0", "0", Tda_Rank_Change_Elite_Main_Battle, tdaDataMap)
		}
	})
}

// PVP 排行榜	PVP 排行榜
type RankChangePvp struct {
	List []*RankUnit `json:"rank_list"` // 排行列表
}

func TdaRankChangePvp(list []*RankUnit) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaRankChangePvp", func() {
		tdaData := &RankChangePvp{
			List: list,
		}

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaRankChangePvp tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track("0", "0", Tda_Rank_Change_Pvp, tdaDataMap); err != nil {
				log.Error("tda TdaRankChangePvp track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(0, "0", "0", Tda_Rank_Change_Pvp, tdaDataMap)
		}
	})
}

// 沙漠密藏排行榜	沙漠密藏排行榜
type RankChangeEventDesert struct {
	List []*RankUnit `json:"rank_list"` // 排行列表
}

func TdaRankChangeEventDesert(list []*RankUnit) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaRankChangeEventDesert", func() {
		tdaData := &RankChangeEventDesert{
			List: list,
		}

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaRankChangeEventDesert tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track("0", "0", Tda_Rank_Change_Event_Desert, tdaDataMap); err != nil {
				log.Error("tda TdaRankChangeEventDesert track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(0, "0", "0", Tda_Rank_Change_Event_Desert, tdaDataMap)
		}
	})
}

type GuildRankUnit struct {
	*RankUnit
	GuildId uint32 `json:"guild_id"` // 联盟id
}

// 聯盟boss_聯盟傷害	聯盟總傷害排行
type RankChangeGuildBossGuild struct {
	List []*GuildRankUnit `json:"rank_list"` // 排行列表
}

func TdaRankChangeGuildBossGuild(list []*GuildRankUnit) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaRankChangeGuildBossGuild", func() {
		tdaData := &RankChangeGuildBossGuild{
			List: list,
		}

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaRankChangeGuildBossGuild tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {

			if err := GetTa().Track("0", "0", Tda_Rank_Change_Guild_Boss_Guild, tdaDataMap); err != nil {
				log.Error("tda TdaRankChangeGuildBossGuild track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(0, "0", "0", Tda_Rank_Change_Guild_Boss_Guild, tdaDataMap)
		}
	})
}

// 聯盟boss_個人傷害	聯盟個人傷害排行
type RankChangeGuildBossSelf struct {
	List []*GuildRankUnit `json:"rank_list"` // 排行列表
}

func TdaRankChangeGuildBossSelf(list []*GuildRankUnit) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaRankChangeGuildBossSelf", func() {
		tdaData := &RankChangeGuildBossSelf{
			List: list,
		}

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaRankChangeGuildBossSelf tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track("0", "0", Tda_Rank_Change_Guild_Boss_Self, tdaDataMap); err != nil {
				log.Error("tda TdaRankChangeGuildBossSelf track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(0, "0", "0", Tda_Rank_Change_Guild_Boss_Self, tdaDataMap)
		}
	})
}

// 聯盟成員戰力	聯盟內成員戰力排行
type RankChangeGuildPower struct {
	List []*GuildRankUnit `json:"rank_list"` // 排行列表
}

func TdaRankChangeGuildPower(list []*GuildRankUnit) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaRankChangeGuildBossSelf", func() {
		tdaData := &RankChangeGuildPower{
			List: list,
		}

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaRankChangeGuildPower tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track("0", "0", Tda_Rank_Change_Guild_Power, tdaDataMap); err != nil {
				log.Error("tda TdaRankChangeGuildPower track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(0, "0", "0", Tda_Rank_Change_Guild_Power, tdaDataMap)
		}
	})
}

// 聯盟成員活躍度	聯盟內成員活躍排行
type RankChangeGuildActive struct {
	List []*GuildRankUnit `json:"rank_list"` // 排行列表
}

func TdaRankChangeGuildActive(list []*GuildRankUnit) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaRankChangeGuildBossSelf", func() {
		tdaData := &RankChangeGuildActive{
			List: list,
		}

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaRankChangeGuildActive tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track("0", "0", Tda_Rank_Change_Guild_Active, tdaDataMap); err != nil {
				log.Error("tda v track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(0, "0", "0", Tda_Rank_Change_Guild_Active, tdaDataMap)
		}
	})
}
