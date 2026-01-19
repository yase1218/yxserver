package model

import (
	"msg"
	"time"
)

type UserStage struct {
	MissionId       int    // 主线关卡id
	StageStar       uint32 `bson:"stage_star"`        // 关卡星数
	StageFirstEnter []int  `bson:"stage_first_enter"` // 关卡第一次进入标识
	StageFirstPass  []int  `bson:"stage_first_pass"`
}

type UserBase struct {
	ClientSettings map[uint32]string `bson:"client_settings"` // 客户端自定义内容

	IsEditName bool // 是否改过名字
	ApData     *ApInfo
	// MissionId                int              // 主线关卡id
	// StageStar                uint32           `bson:"stage_star"` // 关卡星数
	HookData                *OnHookData      // 挂机数据
	QuickOnHookData         *QuickOnHookData // 快速挂机数据
	NextDailyRefreshTime    uint32           // 日常刷新时间
	NextWeeklyRefreshTime   uint32           // 周常刷新时间
	NextRefreshActivityTime uint32           // 下次刷新活动时间
	Forbidden               bool             `bson:"forbidden"`      // 是否被封禁
	OnlineSeconds           uint32           `bson:"online_seconds"` // 累计在线时长(秒)
	LoginTime               uint32
	LogoutTime              uint32
	CreateTime              uint32
	UpdateTime              uint32
	ShipId                  uint32
	SupportId               []uint32
	GuideData               []*GuideInfo
	PopUps                  []*PopUpInfo
	LastActiveTime          uint32 // 上次活跃时间
	ActiveDay               uint32 // 活跃天数
	MaxActiveDay            uint32 // 最大连续活跃天数
	MissData                *MissionData
	BattleData              *msg.BattleData // 客户端上报的战斗数据 不知道还用不用
	PokerSlotCount          uint32
	//ChannelId                uint32
	ExtraInfo         string
	Ip                string
	VideoFlag         uint32
	QuestionIds       []string
	RewardQuestionIds []string
	Charge            []*ChargeInfo
	MonthCard         []*MonthcardInfo
	MainFund          []*MainFundInfo
	//InviteCode               string
	Ad                       []*AdInfo
	MonthCardDailyRewardTime uint32 // 月卡每日奖励时间
	Attrs                    map[uint32]*Attr
	RankLikesMap             map[uint32]uint32
	Combat                   uint32 // 战力
	DailyApData              []*DailyApInfo
	TalentData               *TalentInfo
	//PvPInfo                  *PvPDetailInfo
	ComboSkill          []uint32
	Adventures          map[uint32]*AdventureInfo
	RankMissionReward   []int
	UsePet              uint32
	ForbiddenChat       uint32 // 1 禁止聊天
	AllianceLeaveTime   uint32 `bson:"alliance_leave_time"`    // 最后申请时间
	AllianceBossHurtMax int64  `bson:"alliance_boss_hurt_max"` //联盟单次boss最大伤害

	LastLoginAt        time.Time                 `bson:"last_login_at"`        // 最后登录时间,隔天或大于24小时才更新,用于记录总登录天数
	LoginCnt           uint32                    `bson:"login_cnt"`            // 总登录天数
	FirstChargePackage []*FirstChargePackageData `bson:"first_charge_package"` // 首充礼包数据
	LastLogoutAt       time.Time                 `bson:"last_logout_at"`       // 最后登出时间,用于打点
}
