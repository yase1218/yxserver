package service

import (
	"fmt"
	"gameserver/internal/game/player"
	"msg"
	"sort"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

const (
	EquipStageRecordMax = 50
	EquipStageTeamCount = 3
	EquipStageMatchMax  = 10
)

var (
	EquipStageRecords []*msg.EquipStageDropRecord

	// 原生队伍
	EquipStageTeamMap map[uint64]*EquipStageTeam // 队长id:队伍

	// 所有队员
	EquipStageTeamMemberMap map[uint64]uint64 // 队员id:队长id

	// 匹配队伍
	EquipStageMatchTeamMap map[uint32][]*EquipStageTeam // 副本id:{队伍人数:{队长id:队伍}}
)

type EquipStageMatchResult struct {
	StageId uint32
	Teams   []*EquipStageTeam
}

func init() {
	EquipStageRecords = make([]*msg.EquipStageDropRecord, 0, EquipStageRecordMax)
	EquipStageTeamMap = make(map[uint64]*EquipStageTeam)
	EquipStageTeamMemberMap = make(map[uint64]uint64)
	EquipStageMatchTeamMap = make(map[uint32][]*EquipStageTeam)
}

func AddEquipStageRecord(record *msg.EquipStageDropRecord) {
	if len(EquipStageRecords) >= EquipStageRecordMax {
		EquipStageRecords = EquipStageRecords[1:]
	}
	EquipStageRecords = append(EquipStageRecords, record)
}

func CreateEquipStageTeam(stageId uint32, master_id uint64) *EquipStageTeam {
	DebugLog("create equip stage team", zap.Uint32("stage id", stageId), zap.Uint64("master id", master_id))
	team := &EquipStageTeam{
		StageId: stageId,
		Uids:    []uint64{master_id},
		TeamId:  fmt.Sprintf("%d_%d", master_id, time.Now().UnixMilli()),
	}

	add_team(team)

	return team
}

func FindEquipStageTeam(master_id uint64) *EquipStageTeam {
	return EquipStageTeamMap[master_id]
}

func IsInEquipStageTeam(uid uint64) bool {
	if _, ok := EquipStageTeamMemberMap[uid]; ok {
		return true
	}
	return false
}

func GetEquipStageTeamMasterId(uid uint64) uint64 {
	if mid, ok := EquipStageTeamMemberMap[uid]; ok {
		return mid
	}
	return 0
}

func GetEquipStageTeam(uid uint64) *EquipStageTeam {
	if team, ok := EquipStageTeamMap[GetEquipStageTeamMasterId(uid)]; ok {
		return team
	}
	return nil
}

func GetEquipStageTeamRoom(uid uint64) *EquipStageRoom {
	team := GetEquipStageTeam(uid)
	if team == nil {
		return nil
	}
	return FindEquipStageRoom(team.RoomId)
}

func EquipStageTeamMatch(team *EquipStageTeam) {
	if team == nil {
		return
	}

	if _, ok := EquipStageMatchTeamMap[team.StageId]; !ok {
		EquipStageMatchTeamMap[team.StageId] = make([]*EquipStageTeam, 0)
	}
	EquipStageMatchTeamMap[team.StageId] = append(EquipStageMatchTeamMap[team.StageId], team)
}

func add_team(team *EquipStageTeam) {
	if team == nil {
		return
	}

	if len(team.Uids) <= 0 {
		return
	}
	master_id := team.Uids[0]
	EquipStageTeamMap[master_id] = team
	add_member(master_id, master_id)
}

func DissolveEquipStageTeam(master_id uint64, ntf bool) {
	remove_team(master_id, ntf)
}

func remove_team(master_id uint64, ntf bool) {
	team, ok := EquipStageTeamMap[master_id]
	if !ok {
		log.Error("team not exist", zap.Uint64("uid", master_id))
		return
	}

	if ntf {
		team.BroadCast(&msg.EquipStageTeamDissolveNtf{}, 0)
	}

	stage_id := team.StageId

	// 删除原生队伍
	delete(EquipStageTeamMap, master_id)
	DebugLog("remove team", zap.Uint64("uid", master_id))

	// 删除所有队员
	for _, mid := range team.Uids {
		remove_member(mid)
		DebugLog("remove team member", zap.Uint64("master id", master_id), zap.Uint64("member id", mid))

		memeber_player := player.FindByUserId(mid)
		if memeber_player != nil {
			memeber_player.SendNotify(&msg.EquipStageTeamDissolveNtf{})
		}
	}

	// 删除匹配队伍
	cancel_match(stage_id, master_id, false)
}

func cancel_match(stage_id uint32, master_id uint64, ntf bool) {
	if _, ok := EquipStageMatchTeamMap[stage_id]; !ok {
		EquipStageMatchTeamMap[stage_id] = make([]*EquipStageTeam, 0)
	}
	for i, team := range EquipStageMatchTeamMap[stage_id] {
		if team.Uids[0] == master_id {
			if ntf {
				team.BroadCast(&msg.EquipStageMatchCancelNtf{}, 0)
			}
			EquipStageMatchTeamMap[stage_id] = append(EquipStageMatchTeamMap[stage_id][:i], EquipStageMatchTeamMap[stage_id][i+1:]...)
			team.MatchState = EquipStage_UnMatch
			break
		}
	}
}

func add_member(uid, mid uint64) {
	EquipStageTeamMemberMap[uid] = mid
}

func remove_member(uid uint64) {
	delete(EquipStageTeamMemberMap, uid)
}

func update_equip_stage_match(now time.Time) {
	rooms := make([]*EquipStageMatchResult, 0)

	for stage_id, wait_teams := range EquipStageMatchTeamMap {
		if len(wait_teams) <= 0 {
			continue
		}
		// 按队伍人数降序排列，优先匹配人数多的队伍
		sort.Slice(wait_teams, func(i, j int) bool {
			return wait_teams[i].MemberCount() > wait_teams[j].MemberCount()
		})

		// 使用动态规划找到所有可能的完美匹配
		used := make([]bool, len(wait_teams))

		for i := 0; i < len(wait_teams); i++ {
			if used[i] {
				continue
			}
			// 超时跳过匹配
			if wait_teams[i].MatchTime.Add(time.Second * EquipStageMatchMax).Before(now) {
				rooms = append(rooms, &EquipStageMatchResult{
					StageId: stage_id,
					Teams:   []*EquipStageTeam{wait_teams[i]},
				})
				used[i] = true
				continue
			}

			// 单个队伍正好满足人数要求
			if wait_teams[i].MemberCount() == EquipStageTeamCount {
				rooms = append(rooms, &EquipStageMatchResult{
					StageId: stage_id,
					Teams:   []*EquipStageTeam{wait_teams[i]},
				})
				used[i] = true
				continue
			}

			// 寻找组合匹配
			combination := findCombination(wait_teams, i, used, wait_teams[i].MemberCount(), []int{i})
			if combination != nil {
				// 创建匹配结果
				matchedTeams := make([]*EquipStageTeam, 0, len(combination))
				for _, idx := range combination {
					matchedTeams = append(matchedTeams, wait_teams[idx])
					used[idx] = true
				}

				rooms = append(rooms, &EquipStageMatchResult{
					StageId: stage_id,
					Teams:   matchedTeams,
				})
			}
		}
		new_teams := make([]*EquipStageTeam, 0)
		for i, team := range wait_teams {
			if !used[i] {
				new_teams = append(new_teams, team)
			}
		}
		EquipStageMatchTeamMap[stage_id] = new_teams
	}

	for _, room := range rooms {
		ntf := &msg.EquipStageMatchResultNtf{
			Players: make([]*msg.PlayerSimpleInfo, 0),
		}

		for _, team := range room.Teams {
			for _, uid := range team.Uids {
				err, player := GetPlayerSimpleInfo(uid)
				if err != msg.ErrCode_SUCC {
					log.Error("GetPlayerSimpleInfo err", zap.Uint64("uid", uid))
					continue
				}
				ntf.Players = append(ntf.Players, ToPlayerSimpleInfo(player))
			}
		}

		CreateEquipStageRoom(room.StageId, room.Teams)

		for _, team := range room.Teams {
			team.BroadCast(ntf, 0)
		}
	}
}

// 寻找满足人数要求的组合
func findCombination(wait_teams []*EquipStageTeam, start_i int, used []bool, current_sum uint32, current_combination []int) []int {
	if current_sum == EquipStageTeamCount {
		return current_combination
	}

	if current_sum > EquipStageTeamCount {
		return nil
	}

	for i := start_i + 1; i < len(wait_teams); i++ {
		if used[i] {
			continue
		}

		newSum := current_sum + wait_teams[i].MemberCount()
		if newSum <= EquipStageTeamCount {
			new_combination := append([]int{}, current_combination...)
			new_combination = append(new_combination, i)

			result := findCombination(wait_teams, i, used, newSum, new_combination)
			if result != nil {
				return result
			}
		}
	}

	return nil
}

func update_equip_stage_room(now time.Time) {
	cancel_list := make([]uint64, 0)
	for _, room := range EquipStageRoomMap {
		// if room.RemoveDeadLine > 0 && now.Unix() >= int64(room.RemoveDeadLine) {
		// 	cancel_list = append(cancel_list, room.ID)
		// }

		if now.Sub(room.CreateAt) >= time.Minute*10 {
			cancel_list = append(cancel_list, room.ID)
			DebugLog("EquipStageRoom deadline", zap.Uint64("id", room.ID), zap.Time("create at", room.CreateAt))
		}

		if now.Unix() > int64(room.ReadyDeadLine) && !room.IsReady {
			cancel_list = append(cancel_list, room.ID)
			DebugLog("EquipStageRoom ready deadline", zap.Uint64("id", room.ID), zap.Time("create at", room.CreateAt))
		}
	}
	for _, id := range cancel_list {
		DisEquipStageRoom(id, true)
	}
}

func DisEquipStageRoom(id uint64, ntf bool) {
	DebugLog("dis equip stage room", zap.Uint64("id", id))
	room, ok := EquipStageRoomMap[id]
	if !ok {
		return
	}
	room.OnDis(ntf)
	delete(EquipStageRoomMap, id)
}

func update_equip_stage(now time.Time) {
	update_equip_stage_match(now)
	update_equip_stage_room(now)
}
