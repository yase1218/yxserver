package model

import (
	"fmt"
	"msg"
)

// Explore 世界探索数据
type Explore struct {
	AccountId        string                   `bson:"account_id"`          // 账号ID
	CurrStageId      uint32                   `bson:"curr_stage_id"`       // 当前关卡ID
	Stages           map[uint32]*ExploreStage `bson:"stages"`              // 已解锁关卡 key:stageId
	RevivePoint      uint32                   `bson:"revive_point"`        // 复活点
	FreeCardCount    uint32                   `bson:"free_card_count"`     // 免费卡包剩余数量
	FreeCardLastTime uint32                   `bson:"free_card_last_time"` // 最后一次获得免费卡包的时间戳
	UseCardCount     uint32                   `bson:"use_card_count"`      // 获得免费卡的总次数
}

// ExploreStage 世界探索关卡
type ExploreStage struct {
	StageId             uint32               `bson:"stage_id"`               // 关卡ID
	UnlockTime          uint32               `bson:"unlock_time"`            // 解锁时间
	FirstEnter          bool                 `bson:"first_enter"`            // 是否第一次进入
	PassTime            uint32               `bson:"pass_time"`              // 通过时间
	CurrPoint           uint32               `bson:"curr_point"`             // 当前玩家所在的路点, 范围 [1, LastPoint]
	CurrPointTime       uint32               `bson:"curr_point_time"`        // 到达当前点的时间
	LastPoint           uint32               `bson:"last_point"`             // 最后一个占领路点
	LastPointUnlockTime uint32               `bson:"last_point_unlock_time"` // 最后一个占领路点解锁时间
	Buildings           map[string]*Building `bson:"buildings"`              // 建筑 key:buildingId
	Buffs               []uint32             `bson:"buffs"`                  // 关卡BUFF列表
}

// Building 建筑
type Building struct {
	WaypointId      uint32 `bson:"waypoint_id"`       // 路点ID
	BuildingId      uint32 `bson:"building_id"`       // 建筑ID
	Level           uint32 `bson:"level"`             // 当前等级
	CreateTime      uint32 `bson:"create_time"`       // 创建时间(第一次占领)
	UnlockTime      uint32 `bson:"unlock_time"`       // 解锁时间(消耗道具解锁)
	LastCollectTime uint32 `bson:"last_collect_time"` // 最后收集时间
}

// NewExplore 创建一个新的世界探索数据
func NewExplore(accountId string) *Explore {
	return &Explore{
		AccountId:        accountId,
		CurrStageId:      0,
		Stages:           make(map[uint32]*ExploreStage),
		RevivePoint:      0,
		FreeCardCount:    5, // 初始每天5次免费卡包
		FreeCardLastTime: 0, // 默认为0表示从未使用过
		UseCardCount:     0,
	}
}

// NewExploreStage 创建一个新的世界探索关卡
func NewExploreStage(stageId, unlockTime uint32) *ExploreStage {
	return &ExploreStage{
		StageId:             stageId,
		UnlockTime:          unlockTime,
		FirstEnter:          false,
		PassTime:            0,
		CurrPoint:           1,
		CurrPointTime:       0,
		LastPoint:           1,
		LastPointUnlockTime: 0,
		Buildings:           make(map[string]*Building),
	}
}

// NewBuilding 创建一个新的建筑
func NewBuilding(waypointId, buildingId, createTime uint32) *Building {
	return &Building{
		WaypointId:      waypointId,
		BuildingId:      buildingId,
		Level:           1,
		CreateTime:      createTime,
		UnlockTime:      0,
		LastCollectTime: createTime,
	}
}

// GetBuildingKey 获取建筑的唯一键
func GetBuildingKey(stageId, waypointId, buildingId uint32) string {
	return fmt.Sprintf("%d_%d_%d", stageId, waypointId, buildingId)
}

// ToSimple 转换为简单数据结构
func (s *ExploreStage) ToSimple() *msg.ExploreStageSimple {
	return &msg.ExploreStageSimple{
		StageId:             s.StageId,
		UnlockTime:          s.UnlockTime,
		FirstEnter:          s.FirstEnter,
		PassTime:            s.PassTime,
		CurrPoint:           s.CurrPoint,
		CurrPointTime:       s.CurrPointTime,
		LastPoint:           s.LastPoint,
		LastPointUnlockTime: s.LastPointUnlockTime,
	}
}

// ConvertToExploreBuilding 转换为传输结构
func (b *Building) ConvertToExploreBuilding(stageId uint32) *msg.ExploreBuilding {
	return &msg.ExploreBuilding{
		StageId:         stageId,
		WaypointId:      b.WaypointId,
		BuildingId:      b.BuildingId,
		Level:           b.Level,
		CreateTime:      b.CreateTime,
		UnlockTime:      b.UnlockTime,
		LastCollectTime: b.LastCollectTime,
	}
}

// UpdateCurrPointAndCheckLastPoint 更新当前点位，并检查是否需要更新最后占领点位
// 如果当前点位大于最后占领点位，则同时更新最后占领点位
// 返回是否更新了最后占领点位
func (s *ExploreStage) UpdateCurrPointAndCheckLastPoint(currPoint, currTime uint32) bool {
	s.CurrPoint = currPoint
	s.CurrPointTime = currTime

	// 如果当前点位大于最后占领点位，同时更新最后占领点位
	if currPoint >= s.LastPoint {
		//log.Debugf("UpdateCurrPointAndCheckLastPoint First currPoint=%d , LastPoint=%d", currPoint, s.LastPoint)
		s.LastPoint = currPoint
		s.LastPointUnlockTime = currTime
		return true
	}

	return false
}
