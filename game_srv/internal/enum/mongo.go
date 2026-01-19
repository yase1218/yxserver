package enum

const (
	DB_SPACE_WAR_GLOBAL = "space_war_global%d"
	DB_SPACE_WAR_LOCAL  = "space_war_local%d"
)

const (
	DB_LOCAL_COL_SYSTEM = "system"
)

type SystemKey string

const (
	DB_SYS_KEY_ALLIANCE_RANK_BOSS_SINGLE_DAILY SystemKey = "alliance_rank_boss_single_daily" // 联盟boss个人伤害榜排行榜日时间
	DB_SYS_KEY_ALLIANCE_RANK_BOSS_DAILY        SystemKey = "alliance_rank_boss_daily"        // 联盟boss伤害榜排行榜日时间
	DB_SYS_KEY_PEAK_FIGHT_RANK_DAY             SystemKey = "peak_fight_rank_boss_day"        // 巅峰战场排行榜日
	DB_SYS_KEY_PEAK_FIGHT_RANK_WEEK            SystemKey = "peak_fight_rank_boss_week"       // 巅峰战场排行榜周
	DB_SYS_KEY_PEAK_FIGHT_RANK_SEASON          SystemKey = "peak_fight_rank_boss_season"     // 巅峰战场排行榜赛季
)
