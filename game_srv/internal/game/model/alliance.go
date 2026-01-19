package model

import "time"

// Alliance 联盟基础信息
type Alliance struct {
	ID             uint32    `bson:"_id"`         // 联盟ID
	Name           string    `bson:"name"`        // 联盟名称
	Banner         uint32    `bson:"banner"`      // 联盟旗帜
	Level          uint32    `bson:"level"`       // 联盟等级
	Declaration    string    `bson:"declaration"` // 联盟宣言
	LeaderID       int64     `bson:"leader_id"`   // 指挥官ID
	PowerRequired  uint32    `bson:"power_req"`   // 入盟战力要求
	MemberCount    uint32    `bson:"mem_count"`   // 当前成员数
	MaxMemberCount uint32    `bson:"max_count"`   // 最大成员数
	TotalPower     uint64    `bson:"total_power"` // 联盟总战力
	TotalTreasure  uint32    `bson:"treasure"`    // 联盟总宝箱数
	CreateTime     time.Time `bson:"create_time"` // 创建时间
	TreasureRule   uint8     `bson:"tres_rule"`   // 宝箱分配规则(1:平均分配 2:活跃度比例分配)
	Exp            uint32    `bson:"exp"`         // 联盟经验
	AutoJoin       bool      `bson:"auto_join"`   // 自动加入联盟
}

// AllianceMember 联盟成员信息
type AllianceMember struct {
	AllianceID        uint32    `bson:"alliance_id"`          // 联盟ID
	PlayerID          int64     `bson:"player_id"`            // 玩家ID
	Name              string    `bson:"name"`                 // 玩家名字
	HeadImg           uint32    `bson:"head_img"`             // 头像
	HeadFrame         uint32    `bson:"head_frame"`           // 头像框
	Position          uint8     `bson:"position"`             // 职位(1:指挥官 2:副指挥官 3:精英 4:普通成员)
	Power             uint32    `bson:"power"`                // 战力
	WeeklyActive      uint32    `bson:"weekly_act"`           // 本周活跃度
	LastOnline        time.Time `bson:"last_online"`          // 最后在线时间
	JoinTime          time.Time `bson:"join_time"`            // 加入时间
	IsMuted           bool      `bson:"is_muted"`             // 是否被禁言
	RedPacketIndexMax int64     `bson:"red_packet_index_max"` // 红包自增索引
}

// AllianceApplication 联盟申请记录
type AllianceApplication struct {
	ID         int64     `bson:"_id"`         // 申请ID
	AllianceID uint32    `bson:"alliance_id"` // 联盟ID
	PlayerID   int64     `bson:"player_id"`   // 申请者ID
	Name       string    `bson:"name"`        // 申请者名字
	HeadImg    uint32    `bson:"head_img"`    // 头像
	HeadFrame  uint32    `bson:"head_frame"`  // 头像框
	Power      uint32    `bson:"power"`       // 战力
	ApplyTime  time.Time `bson:"apply_time"`  // 申请时间
	Status     uint8     `bson:"status"`      // 状态(0:待处理 1:已同意 2:已拒绝)
}

// AllianceRedPacket 联盟红包
type AllianceRedPacket struct {
	ID          int64  `bson:"_id"`         // 红包ID
	RedPacketId uint32 `bson:"conf_id"`     // 红包配置ID
	AllianceID  uint32 `bson:"alliance_id"` // 联盟ID
	SenderID    int64  `bson:"sender_id"`   // 发送者ID
	SenderName  string `bson:"sender_name"` // 发送者名字
	HeadImg     uint32 `bson:"head_img"`    // 头像
	HeadFrame   uint32 `bson:"head_frame"`  // 头像框
	CreateTime  int64  `bson:"create_time"` // 创建时间
	ExpireTime  int64  `bson:"expire_time"` // 过期时间
}

// AllianceBossRecord 联盟BOSS战斗记录
type AllianceBossRecord struct {
	AllianceID uint32    `bson:"alliance_id"` // 联盟ID
	PlayerID   int64     `bson:"player_id"`   // 玩家ID
	BossID     uint32    `bson:"boss_id"`     // BOSS ID
	Damage     uint64    `bson:"damage"`      // 伤害值
	RecordTime time.Time `bson:"record_time"` // 记录时间
	WeekNumber uint32    `bson:"week_num"`    // 周数
}

// AllianceTask 联盟任务进度
type AllianceTask struct {
	AllianceID uint32    `bson:"alliance_id"` // 联盟ID
	PlayerID   int64     `bson:"player_id"`   // 玩家ID
	TaskID     uint32    `bson:"task_id"`     // 任务ID
	Progress   uint32    `bson:"progress"`    // 进度
	Status     uint8     `bson:"status"`      // 状态(0:进行中 1:已完成 2:已领取)
	UpdateTime time.Time `bson:"update_time"` // 更新时间
	WeekNumber uint32    `bson:"week_num"`    // 周数
}

// AllianceShopItem 联盟商店商品
type AllianceShopItem struct {
	ID           uint32    `bson:"_id"`          // 商品ID
	AllianceID   uint32    `bson:"alliance_id"`  // 联盟ID
	ItemID       uint32    `bson:"item_id"`      // 物品ID
	Price        uint32    `bson:"price"`        // 价格
	HasRedPacket bool      `bson:"has_redpack"`  // 是否带红包
	ShopType     uint32    `bson:"shop_type"`    // 商店类型(1:每日商店 2:每月商店)
	RefreshTime  time.Time `bson:"refresh_time"` // 刷新时间
}

// AllianceMemberBossDamage 联盟成员BOSS伤害排名
type AllianceMemberBossDamage struct {
	AllianceID uint32 `bson:"alliance_id"` // 联盟ID
	PlayerID   int64  `bson:"player_id"`   // 玩家ID
	BossID     uint32 `bson:"boss_id"`     // BOSS ID
	Damage     uint64 `bson:"damage"`      // 伤害值
	Name       string `bson:"name"`        // 玩家名字
	HeadImg    uint32 `bson:"head_img"`    // 头像
	HeadFrame  uint32 `bson:"head_frame"`  // 头像框
	Position   uint8  `bson:"position"`    // 职位
}

// AllianceBossDamage 联盟BOSS伤害排名
type AllianceBossDamage struct {
	AllianceID  uint32 `bson:"_id"`          // 联盟ID
	Name        string `bson:"name"`         // 联盟名称
	Flag        uint32 `bson:"flag"`         // 联盟旗帜
	Level       uint32 `bson:"level"`        // 联盟等级
	TotalDamage uint64 `bson:"total_damage"` // 总伤害
	MemberCount uint32 `bson:"mem_count"`    // 成员数
	LeaderName  string `bson:"leader_name"`  // 指挥官名字
}
