package service

import (
	"gameserver/internal/game/player"
	"kernel/tools"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	EquipStageRoomID  uint64
	EquipStageRoomMap map[uint64]*EquipStageRoom
)

func init() {
	EquipStageRoomID = 0
	EquipStageRoomMap = make(map[uint64]*EquipStageRoom)
}

type EquipStageRoom struct {
	ID             uint64
	StatgeID       uint32
	Teams          []*EquipStageTeam
	ReadyDeadLine  uint64
	ReadyUids      map[uint64]struct{}
	IsReady        bool
	LoadFinishUids map[uint64]struct{}
	EndUids        map[uint64]struct{}
	StartDeadLine  uint64
	FightId        uint32

	//RemoveDeadLine uint64

	CreateAt time.Time
}

func GenEquipStageRoomID() uint64 {
	EquipStageRoomID++
	return EquipStageRoomID
}

func CreateEquipStageRoom(stage_id uint32, teams []*EquipStageTeam) *EquipStageRoom {
	id := GenEquipStageRoomID()
	for _, team := range teams {
		team.MatchState = EquipStage_Matched
		team.RoomId = id
		DebugLog("matched team", zap.Any("team", team))
	}
	room := &EquipStageRoom{
		ID:             id,
		StatgeID:       stage_id,
		Teams:          teams,
		ReadyDeadLine:  uint64(time.Now().Add(time.Second * 20).Unix()),
		ReadyUids:      make(map[uint64]struct{}),
		LoadFinishUids: make(map[uint64]struct{}),
		EndUids:        make(map[uint64]struct{}),
		//RemoveDeadLine: uint64(time.Now().Add(time.Minute * 10).Unix()),
		CreateAt: time.Now(),
	}
	EquipStageRoomMap[id] = room
	DebugLog("create equip stage room", zap.Uint64("room id", id))
	return room
}

func FindEquipStageRoom(id uint64) *EquipStageRoom {
	return EquipStageRoomMap[id]
}

func (r *EquipStageRoom) Ready(uid uint64) {
	r.ReadyUids[uid] = struct{}{}

	if len(r.ReadyUids) == r.MemberCount() {
		r.IsReady = true
		p := r.Player()
		if p != nil {
			CreateFight(p, r.StatgeID, 0, tools.Uint64ToBytes(r.ID))
		}
	}
}

func (r *EquipStageRoom) MemberCount() int {
	count := uint32(0)
	for _, team := range r.Teams {
		count += team.MemberCount()
	}
	return int(count)
}

func (r *EquipStageRoom) Player() *player.Player {
	// todo 当前战斗创建以来玩家对象发送消息 这里返回第一个玩家
	for _, team := range r.Teams {
		for _, uid := range team.Uids {
			p := player.FindByUserId(uid)
			if p != nil {
				return p
			}
		}
	}
	return nil
}

func (r *EquipStageRoom) OnDis(ntf bool) {
	for _, team := range r.Teams {
		for _, uid := range team.Uids {
			p := player.FindByUserId(uid)
			if p == nil {
				continue
			}
		}
		DissolveEquipStageTeam(team.MasterId(), ntf)
	}
}

func (r *EquipStageRoom) LoadFinish(uid uint64) bool {
	r.LoadFinishUids[uid] = struct{}{}
	return len(r.LoadFinishUids) == r.MemberCount()
}

func (r *EquipStageRoom) EndUser(uid uint64) bool {
	r.EndUids[uid] = struct{}{}
	return len(r.EndUids) == r.MemberCount()
}

func (r *EquipStageRoom) BroadCast(m proto.Message, except_id uint64) {
	for _, team := range r.Teams {
		team.BroadCast(m, except_id)
	}
}
