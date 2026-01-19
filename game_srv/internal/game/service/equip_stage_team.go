package service

import (
	"gameserver/internal/game/player"
	"msg"
	"time"

	"google.golang.org/protobuf/proto"
)

type EquipStageMatchState uint32

const (
	EquipStage_UnMatch EquipStageMatchState = iota
	EquipStage_Matching
	EquipStage_Matched
)

type EquipStageTeam struct {
	StageId uint32   // 副本id
	Uids    []uint64 // 成员ids 首位是队长

	MatchState EquipStageMatchState
	MatchTime  time.Time
	RoomId     uint64
	TeamId     string
}

func (m *EquipStageTeam) MasterId() uint64 {
	if len(m.Uids) == 0 {
		return 0
	}
	return m.Uids[0]
}

func (m *EquipStageTeam) MemberCount() uint32 {
	return uint32(len(m.Uids))
}

func (m *EquipStageTeam) AddMember(uid uint64) {
	if len(m.Uids) >= EquipStageTeamCount {
		return
	}
	m.Uids = append(m.Uids, uid)
	add_member(uid, m.Uids[0])
}

func (m *EquipStageTeam) RemoveMember(uid uint64) {
	for i := 0; i < len(m.Uids); i++ {
		if m.Uids[i] == uid {
			m.Uids = append(m.Uids[:i], m.Uids[i+1:]...)
			break
		}
	}
	remove_member(uid)

	ntf := &msg.EquipStageTeamLeaveNtf{
		Uid: uid,
	}
	m.BroadCast(ntf, 0)
}

func (m *EquipStageTeam) BroadCast(message proto.Message, except_id uint64) {
	for _, mid := range m.Uids {
		if mid == except_id {
			continue
		}
		player := player.FindByUserId(mid)
		if player != nil {
			player.SendNotify(message)
		}
	}
}

func (m *EquipStageTeam) StartMatch() {
	m.MatchState = EquipStage_Matching
	m.MatchTime = time.Now()
	EquipStageTeamMatch(m)
	m.BroadCast(&msg.EquipStageMatchNtf{}, 0)
}
