package model

import (
	"kernel/tools"
	"msg"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GlobalAccount struct {
	UserId     string
	AccountId  int64
	ServerId   uint32
	InviteCode string
	Level      uint32
	Nick       string
}

type AccBasic struct {
	AccountId int64  // 账号
	Nick      string // 昵称
	HeadImg   uint32 // 头像
	HeadFrame uint32 // 头像框
	Title     uint32 // 称号
	ShipId    uint32 // 机甲id
}

type Account struct {
	ServerId   uint32
	ServerName string
	UserId     string
	AccountId  int64  // 账号
	Nick       string // 昵称
	Level      uint32 // 等级

	HeadImg   uint32 // 头像
	HeadFrame uint32 // 头像框
	Title     uint32 // 称号

	ClientSettings map[uint32]string `bson:"client_settings"` // 客户端自定义内容

	IsEditName               bool // 是否改过名字
	ApData                   *ApInfo
	MissionId                int              // 主线关卡id
	StageStar                uint32           `bson:"stage_star"` // 关卡星数
	HookData                 *OnHookData      // 挂机数据
	QuickOnHookData          *QuickOnHookData // 快速挂机数据
	NextDailyRefreshTime     uint32           // 日常刷新时间
	NextWeeklyRefreshTime    uint32           // 周常刷新时间
	IsNew                    bool             // 是否是新手
	Forbidden                bool             `bson:"forbidden"` // 是否被封禁
	GlobalMailId             int64            // 全局邮件id
	OnlineSeconds            uint32           `bson:"online_seconds"` // 累计在线时长(秒)
	LoginTime                uint32
	LogoutTime               uint32
	CreateTime               uint32
	UpdateTime               uint32
	ShipId                   uint32
	SupportId                []uint32
	GuideData                []*GuideInfo
	PopUps                   []*PopUpInfo
	LastActiveTime           uint32 // 上次活跃时间
	ActiveDay                uint32 // 活跃天数
	MaxActiveDay             uint32 // 最大连续活跃天数
	MissData                 *MissionData
	BattleData               *msg.BattleData // 战斗数据
	PokerSlotCount           uint32
	ChannelId                uint32
	ExtraInfo                string
	Ip                       string
	VideoFlag                uint32
	QuestionIds              []string
	RewardQuestionIds        []string
	Charge                   []*ChargeInfo
	MonthCard                []*MonthcardInfo
	MainFund                 []*MainFundInfo
	InviteCode               string
	Ad                       []*AdInfo
	MonthCardDailyRewardTime uint32 // 月卡每日奖励时间
	Attrs                    map[uint32]*Attr
	RankLikesMap             map[uint32]uint32
	Combat                   uint32 // 战力
	DailyApData              []*DailyApInfo
	TalentData               *TalentInfo
	PvPInfo                  *PvPDetailInfo
	ComboSkill               []uint32
	Adventures               map[uint32]*AdventureInfo
	RankMissionReward        []int
	UsePet                   uint32
	ForbiddenChat            uint32 // 1 禁止聊天
	AllianceLeaveTime        uint32 `bson:"alliance_leave_time"`    // 最后申请时间
	AllianceBossHurtMax      int64  `bson:"alliance_boss_hurt_max"` //联盟单次boss最大伤害
	StageFirstEnter          []int  `bson:"stage_first_enter"`      // 关卡第一次进入标识
	StageFirstPass           []int  `bson:"stage_first_pass"`

	LastLoginAt time.Time `bson:"last_login_at"` // 最后登录时间,隔天或大于24小时才更新,用于记录总登录天数
	LoginCnt    uint32    `bson:"login_cnt"`     // 总登录天数
}

type PvPAccount struct {
	AccountId uint32
	Nick      string // 昵称
	Level     uint32 // 等级
	HeadImg   uint32 // 头像
	HeadFrame uint32 // 头像框
	PvPInfo   *PvPDetailInfo
	ShipId    uint32
	Use       uint32 // 1 使用中
}

type PvPDetailInfo struct {
	SeasonId  uint32 // 赛季id
	Stage     uint32 // 阶
	Star      uint32 // 星
	StarScore uint32 // 保星积分
	Score     uint32 // 荣耀积分
	CWinTimes uint32 // 连胜次数
}

// ApInfo 体力信息
type ApInfo struct {
	BuyTimes         uint32 // 购买次数
	RecoverStartTime uint32 // 体力恢复开始时间
	NextBuyTime      uint32 // 下次购买时间
}

type GuideInfo struct {
	Id    uint32
	Value uint32
}

type ChargeInfo struct {
	Id        int
	Value     int
	ResetTime int
}

type MonthcardInfo struct {
	Id                int
	EndTime           int
	NextGetRewardTime int
}

type MainFundInfo struct {
	Id      int
	FreeId  int
	PayId   int
	BuyFlag int
}

type AdInfo struct {
	AdId          uint32
	Para          uint32
	Times         uint32
	Flag          bool // true 表示当前正在看的 false表示不是
	NextResetTime uint32
}

type PopUpInfo struct {
	Id        uint32
	PopUpType uint32
}

// OnHookData 挂机数据
type OnHookData struct {
	StartTime uint32       // 挂机开始时间(s)
	TotalTime uint32       // 挂机总时间(s)
	Items     []*FloatItem // 挂机道具
	IsNotice  bool
}

// QuickOnHookData 快速挂机数据
type QuickOnHookData struct {
	BuyTimes    uint32       // 已经购买的次数
	NextBuyTime uint32       // 下一次购买重置时间
	Items       []*FloatItem // 剩余道具
}

type MissionData struct {
	MissionId int
	StartTime uint32
	Total     uint32
	Speed     uint32
	Ads       []*AdInfo
}

type DailyApInfo struct {
	Id        uint32
	StartTime uint32
	EndTime   uint32
	State     uint32
}

type TalentInfo struct {
	NormalId uint32
	KeyId    uint32
	Attrs    map[uint32]*Attr
	Parts    []uint32
}

type AdventureInfo struct {
	Id         uint32
	State      uint32 // 1 完成 0 没有
	UpdateTime uint32
}

type FirstChargePackageData struct {
	Id         uint32
	State      []uint32 // 0可领取 1已领取 2明日可领
	LoginCount int32    //登录次数
}

func NewTalentInfo(normalId, keyId uint32) *TalentInfo {
	return &TalentInfo{
		NormalId: normalId,
		KeyId:    keyId,
		Attrs:    make(map[uint32]*Attr),
	}
}

func NewDailApInfo(id, start, end uint32) *DailyApInfo {
	return &DailyApInfo{
		Id:        id,
		StartTime: start,
		EndTime:   end,
	}
}

func NewMissionData(missionId int, startTime uint32) *MissionData {
	return &MissionData{
		MissionId: missionId,
		StartTime: startTime,
		Total:     0,
		Speed:     10,
	}
}

func NewChargeInfo(id, value int) *ChargeInfo {
	return &ChargeInfo{
		Id:    id,
		Value: value,
	}
}

func NewMonthcardInfo(id, endTime int) *MonthcardInfo {
	return &MonthcardInfo{
		Id:      id,
		EndTime: endTime,
	}
}

func NewMainFundInfo(id int) *MainFundInfo {
	return &MainFundInfo{
		Id: id,
	}
}

func NewAdInfo(adId, para, nextResetTime uint32) *AdInfo {
	return &AdInfo{
		AdId:          adId,
		Para:          para,
		Flag:          true,
		NextResetTime: nextResetTime,
	}
}

func NewPvpInfo(id uint32) *PvPDetailInfo {
	return &PvPDetailInfo{
		SeasonId: id,
		Stage:    1,
		Star:     1,
	}
}

func NewAdventureInfo(id, state uint32) *AdventureInfo {
	return &AdventureInfo{
		Id:         id,
		State:      state,
		UpdateTime: tools.GetCurTime(),
	}
}

type PlayerAllianceApplication struct {
	AllianceID primitive.ObjectID `bson:"alliance_id"` // 目标联盟ID
	ApplyTime  time.Time          `bson:"apply_time"`  // 申请时间
	Status     uint8              `bson:"status"`      // 申请状态
}

type PlayerAllianceInvite struct {
	AllianceID primitive.ObjectID `bson:"alliance_id"` // 联盟ID
	InviterID  int64              `bson:"inviter_id"`  // 邀请人accountID
	InviteTime time.Time          `bson:"invite_time"` // 邀请时间
	Status     uint8              `bson:"status"`      // 处理状态
}
