package enum

import "time"

const (
	Redis_Alliance_Rank_Expire  = time.Hour * 24 * 30 * 3
	Redis_PeakFight_Rank_Expire = time.Hour * 24 * 30 * 3
)

const (
	Redis_Key_Alliance_Rank_Single_Power  = "alliance_rank_single_power"  // 成员战力榜（联盟内）
	Redis_Key_Alliance_Rank_Single_Active = "alliance_rank_single_active" // 成员活跃榜（联盟内）
	Redis_Key_Alliance_Boss_Single_Rank   = "alliance_boss_single_rank"   // 联盟BOSS伤害榜(个人)
	Redis_Key_Alliance_Boss_Rank          = "alliance_boss_rank"          // 联盟BOSS伤害榜(联盟)

	Redis_Key_PeakFight_Rank = "peak_fight_rank" // 巅峰战场排行榜
)
