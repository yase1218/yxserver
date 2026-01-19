package model

import (
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"github.com/zy/game_data/template"
)

type UserData struct {
	UserId     uint64
	Version    uint64 // 数据版本号
	AccountId  string `bson:"accountid"` // 账号(平台openid)
	ServerId   uint32
	ServerName string
	ChannelId  uint32 // 没用？
	Nick       string // 昵称
	Level      uint32 // 等级
	IsNew      bool   // 是否是新手
	IsRegister bool   // 是否需要进入选角
	HeadImg    uint32 // 头像
	HeadFrame  uint32 // 头像框
	Title      uint32 // 称号

	//Roles           []*Role            // 驾驶员
	//StarPort        *AccountStarPort   // 星港
	//Explore         *Explore           // 世界探索
	BaseInfo        *UserBase
	StageInfo       *UserStage
	Items           *UserItems         // 道具
	Task            *AccountTask       // 任务
	Mission         *AccountMission    // 关卡
	Ships           *UserShips         // 机甲
	Team            *AccountTeam       // 编队
	Equip           *AccountEquip      // 装备
	Shop            *AccountShop       // 商店
	PlayMethod      *AccountPlayMethod // 玩法
	MailData        *UserMail          `bson:"mail"` // 邮件
	Weapon          *AccountWeapon     // 武器
	Treasure        *AccountTreasure   // 秘宝
	Poker           *AccountPoker      // 扑克
	AccountActivity *AccountActivity   // 身上所有活动
	CardPool        *AccountCardPool   // 卡池
	Appearance      *AccountAppearance // 外观
	PetData         *AccountPet        // 宠物系统
	//AccountFri      *AccountFriend     // 好友
	FriendData      *UserFriend      // 好友
	Fight           *Fight           // 战斗
	PeakFight       *PeakFight       // 巅峰战场
	Contract        *Contract        // 疯狂合约
	Desert          *DesertFight     // 沙漠大冒险
	Arena           *ArenaPlayerData // 竞技场
	Atlas           *Atlas           // 图鉴
	LuckSale        *LuckSale        // 幸运售货机
	FunctionPreview *FunctionPreview // 功能预览
	EquipStage      *UserEquipStage
	ResourcesPass   *UserResourcesPass   `bson:"resources_pass"` // 资源本
	Likes           *LikesInfo           // 点赞记录
	Ranks           *RankInfo            // 排行榜相关
	Personalized    *AccountPersonalized // 推送
	WeekPass        *WeekPass            // 周常本数据
}

type UserItems struct {
	Items []*Item
}

type UserShips struct {
	Ships []*Ship
}

type UserResourcesPass struct {
	PassList      []*ResourcesPass
	LastResetTime time.Time
}

type LikesInfo struct {
	LikesMap map[template.RankType]bool
}

type RankInfo struct {
	NormalPassRewardInfo map[uint32]bool // 普通关卡首通奖励领取记录
	ElitePassRewardInfo  map[uint32]bool // 精英关卡首通奖励领取记录
}

func (s *UserShips) GetShipCoatId(shipId uint32) int {
	var coatId int
	var ship *Ship
	for i := 0; i < len(s.Ships); i++ {
		if s.Ships[i].Id == shipId {
			ship = s.Ships[i]
		}
	}
	if ship == nil {
		log.Error("GetShipCoatId get nil", zap.Uint32("shipId", shipId))
		return 0
	}
	for _, item := range ship.CoatMap {
		if item.Status == CoatOn {
			coatId = item.Id
			break
		}
	}
	return coatId
}
