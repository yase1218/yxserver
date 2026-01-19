package service

import (
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/kenum"
	"kernel/protocol"
	"msg"
	"time"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

func LoadEquipStage(pid uint32, p *player.Player) {
	res := &msg.LoadEquipStageAck{
		Result: msg.ErrCode_SUCC,
		StageInfo: &msg.EquipStageInfo{
			RewardNum: p.UserData.EquipStage.RewardNum,
			RewardBuy: p.UserData.EquipStage.RewardBuy,
			Record:    p.UserData.EquipStage.RecordList(),
		},
	}
	defer p.SendResponse(pid, res, res.Result)
}

func LoadEquipStageDropRecord(pid uint32, p *player.Player) {
	res := &msg.EquipStageDropRecordAck{
		Result: msg.ErrCode_SUCC,
		Records: &msg.EquipStageDropRecords{
			Records:   make([]*msg.EquipStageDropRecord, 0, len(EquipStageRecords)),
			PlayerMap: make(map[uint64]*msg.PlayerSimpleInfo),
		},
	}
	defer p.SendResponse(pid, res, res.Result)

	for _, v := range EquipStageRecords {
		if _, ok := res.Records.PlayerMap[v.Uid]; ok {
			res.Records.Records = append(res.Records.Records, v)
			continue
		}

		err, p := GetPlayerSimpleInfo(v.Uid)
		if err != msg.ErrCode_SUCC {
			log.Error("GetPlayerSimpleInfo err", zap.Uint64("uid", v.Uid))
			continue
		}

		res.Records.Records = append(res.Records.Records, v)
		res.Records.PlayerMap[v.Uid] = ToPlayerSimpleInfo(p)
	}
}

func EquipStageTeamInvite(pid uint32, p *player.Player, req *msg.EquipStageTeamInviteReq) {
	res := &msg.EquipStageTeamInviteAck{
		Result:  msg.ErrCode_SUCC,
		StageId: req.StageId,
		Uid:     req.Uid,
	}
	defer p.SendResponse(pid, res, res.Result)

	DebugLog("EquipStageTeamInvite", zap.Any("req", req), ZapUser(p))

	stage_t := template.GetMissionTemplate().GetMission(int(req.StageId))
	if stage_t == nil {
		log.Error("stage cfg nil", zap.Uint32("stageId", req.StageId))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}
	// 解锁判定
	if err := CanLock(p, stage_t.UnlockCond); err != msg.ErrCode_SUCC {
		res.Result = msg.ErrCode_EquipStage_Lock
		return
	}

	//var team *EquipStageTeam
	//mid := GetEquipStageTeamMasterId(p.GetUserId())
	//if mid == 0 { // 没有队伍
	//	team = CreateEquipStageTeam(req.StageId, p.GetUserId())
	//} else { // 有队伍
	//	team = FindEquipStageTeam(mid)
	//	if team == nil {
	//		log.Error("team nil", ZapUser(p))
	//		res.Result = msg.ErrCode_SYSTEM_ERROR
	//		return
	//	}
	//	if team.StageId != req.StageId { // 副本id变化
	//		log.Error("stage err", zap.Uint32("req stage", req.StageId), zap.Uint32("team stage", team.StageId), ZapUser(p))
	//		res.Result = msg.ErrCode_SYSTEM_ERROR
	//		return
	//		// if p.GetUserId() == mid {// 是队长
	//
	//		// }
	//		// DissolveEquipStageTeam(p.GetUserId(), true)
	//		// CreateEquipStageTeam(req.StageId, p.GetUserId())
	//	}
	//	if len(team.Uids) >= EquipStageTeamCount { // 队伍人数满
	//		res.Result = msg.ErrCode_EquipStage_Team_Full
	//		return
	//	}
	//}
	mid := GetEquipStageTeamMasterId(p.GetUserId())
	team := FindEquipStageTeam(mid)
	if team == nil {
		log.Error("EquipStageTeamInvite team nil", zap.Any("req", req), ZapUser(p))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	if team.MatchState != EquipStage_UnMatch {
		log.Error("invite when match", zap.Uint64("uid", p.GetUserId()))
		res.Result = msg.ErrCode_EquipStage_Team_Mached
		return
	}

	target_player := player.FindByUserId(req.Uid)
	if target_player == nil {
		res.Result = msg.ErrCode_OTHER_NOT_ONLINE
		return
	}

	// 解锁判定
	if err := CanLock(target_player, stage_t.UnlockCond); err != msg.ErrCode_SUCC {
		res.Result = msg.ErrCode_EquipStage_Target_Lock
		return
	}

	if GetEquipStageTeamMasterId(req.Uid) > 0 { // 被邀请人已经在队伍
		res.Result = msg.ErrCode_EquipStage_Team_Already
		return
	}

	ntf := &msg.EquipStageTeamInviteNtf{
		StageId: req.StageId,
		TeamId:  req.TeamId,
		Player:  ToPlayerSimpleInfo(getSimpleInfoFromUser(p.UserData)),
	}
	target_player.SendNotify(ntf)
}

func EquipStageTeamAccept(pid uint32, p *player.Player, req *msg.EquipStageTeamAcceptReq) {
	res := &msg.EquipStageTeamAcceptAck{
		Result:  msg.ErrCode_SUCC,
		Accept:  req.Accept,
		Players: make([]*msg.PlayerSimpleInfo, 0),
		TeamId:  req.TeamId,
	}
	defer p.SendResponse(pid, res, res.Result)

	if GetEquipStageTeamMasterId(p.GetUserId()) != 0 {
		res.Result = msg.ErrCode_EquipStage_Team_Already
		return
	}

	inviter := player.FindByUserId(req.Uid) //邀请者
	if inviter == nil {
		log.Error("EquipStageTeamAccept master nil", zap.Uint64("master id", req.Uid), zap.Uint64("uid", p.GetUserId()))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	mid := GetEquipStageTeamMasterId(req.Uid)
	if mid == 0 {
		res.Result = msg.ErrCode_EquipStage_Team_Dissolve
		return
	}

	team := FindEquipStageTeam(mid)
	if team == nil {
		if req.Accept { // 接受
			res.Result = msg.ErrCode_EquipStage_Team_Dissolve
		}
		return
	}

	if team.StageId != req.StageId { // 副本id变化
		if req.Accept { // 接受
			res.Result = msg.ErrCode_EquipStage_Team_Dissolve
		}
		return
	}
	if team.TeamId != req.TeamId {
		log.Debug("EquipStageTeamAccept Dissolve", zap.String("TeamId", team.TeamId), zap.String("req teamID", req.TeamId))
		if req.Accept { // 接受
			res.Result = msg.ErrCode_EquipStage_Team_Dissolve
		}
		return
	}

	if req.Accept { // 接受
		if len(team.Uids) >= EquipStageTeamCount { // 队伍人数满
			res.Result = msg.ErrCode_EquipStage_Team_Full
			return
		}
		if team.MatchState != EquipStage_UnMatch { // 已匹配
			log.Error("accept when match", zap.Uint64("uid", p.GetUserId()))
			res.Result = msg.ErrCode_EquipStage_Team_Mached
			return
		}
		if GetEquipStageTeamMasterId(p.GetUserId()) > 0 { // 已经在队伍
			res.Result = msg.ErrCode_EquipStage_Team_Already
			return
		}

		stage := template.GetMissionTemplate().GetMission(int(req.StageId))
		if stage == nil {
			log.Error("EquipStageTeamAccept stage cfg nil", zap.Uint32("stageId", req.StageId))
			res.Result = msg.ErrCode_SYSTEM_ERROR
			return
		}
		// 解锁判定
		if err := CanLock(p, stage.UnlockCond); err != msg.ErrCode_SUCC {
			res.Result = msg.ErrCode_EquipStage_Lock
			return
		}

		team.AddMember(p.GetUserId())
		res.Players = equip_stage_team_proto_members(team, 0)

		team.BroadCast(&msg.EquipStageTeamJoinNtf{
			User: ToPlayerSimpleInfo(getSimpleInfoFromUser(p.UserData)),
		}, p.GetUserId())
	} else { // 拒绝
		// 目前不用通知邀请者
	}
}

func equip_stage_team_proto_members(team *EquipStageTeam, except_id uint64) []*msg.PlayerSimpleInfo {
	var members []*msg.PlayerSimpleInfo
	for _, mid := range team.Uids {
		if mid == except_id {
			continue
		}
		err, player := GetPlayerSimpleInfo(mid)
		if err != msg.ErrCode_SUCC {
			log.Error("GetPlayerSimpleInfo err", zap.Uint64("uid", mid), zap.Uint64("mid", mid))
			continue
		}
		members = append(members, ToPlayerSimpleInfo(player))
	}
	return members
}

func EquipStageTeamLeave(pid uint32, p *player.Player, offline bool) {
	res := &msg.EquipStageTeamLeaveAck{Result: msg.ErrCode_SUCC}
	defer func() {
		if offline {
			return
		}
		p.SendResponse(pid, res, res.Result)
	}()

	mid := GetEquipStageTeamMasterId(p.GetUserId())
	if mid == 0 { // 不是队长不能离队
		res.Result = msg.ErrCode_EquipStage_Team_Master_Leave
		if !offline {
			log.Error("EquipStageTeamLeave master nil")
		}
		return
	}

	// 队长不能离队
	if p.GetUserId() == mid {
		if offline {
			DissolveEquipStageTeam(p.GetUserId(), true)
			return
		}
		res.Result = msg.ErrCode_EquipStage_Team_Master_Leave
		return
	}

	team := FindEquipStageTeam(mid)
	if team == nil {
		log.Error("team is nil", zap.Uint64("uid", p.GetUserId()), zap.Uint64("mid", mid))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	if team.MatchState != EquipStage_UnMatch { // 已匹配
		log.Error("leave when match", zap.Uint64("uid", p.GetUserId()))
		res.Result = msg.ErrCode_EquipStage_Team_Mached
		return
	}

	team.RemoveMember(p.GetUserId())
}

func EquipStageRoomLeave(p *player.Player) {
	room := FindEquipStageRoom(p.GetUserId())
	if room == nil {
		//log.Error("room nil", zap.Uint64("uid", p.GetUserId()))
		return
	}
	if room.EndUser(p.GetUserId()) {
		DisEquipStageRoom(room.ID, false)
	}
}

func EquipStageTeamDissolve(pid uint32, p *player.Player) {
	res := &msg.EquipStageTeamDissolveAck{Result: msg.ErrCode_SUCC}
	defer p.SendResponse(pid, res, res.Result)

	mid := GetEquipStageTeamMasterId(p.GetUserId())
	if mid == 0 { // 没有自己的成员信息
		log.Error("EquipStageTeamDissolve master nil")
		return
	}

	if p.GetUserId() != mid { // 不是队长不能解散队伍
		res.Result = msg.ErrCode_EquipStage_Team_Master_Dissolve
		return
	}

	team := FindEquipStageTeam(mid)
	if team == nil {
		log.Error("team is nil", zap.Uint64("uid", p.GetUserId()), zap.Uint64("mid", mid))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	if team.MatchState != EquipStage_UnMatch { // 已匹配
		log.Error("dissolve when match", zap.Uint64("uid", p.GetUserId()))
		res.Result = msg.ErrCode_EquipStage_Team_Mached
		return
	}

	DissolveEquipStageTeam(p.GetUserId(), true)
}

func EquipStageMatch(pid uint32, p *player.Player, req *msg.EquipStageMatchReq) {
	res := &msg.EquipStageMatchAck{Result: msg.ErrCode_SUCC}
	defer p.SendResponse(pid, res, res.Result)

	DebugLog("EquipStageMatch", zap.Any("req", req), ZapUser(p))

	stage_t := template.GetMissionTemplate().GetMission(int(req.StageId))
	if stage_t == nil {
		log.Error("stage nil", zap.Uint32("stage id", req.StageId))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	mid := GetEquipStageTeamMasterId(p.GetUserId())
	var team *EquipStageTeam
	if mid == 0 { // 没队伍
		//team = CreateEquipStageTeam(req.StageId, p.GetUserId())
		res.Result = msg.ErrCode_EquipStage_Team_Dissolve
		return

	} else { // 有队伍
		if p.GetUserId() != mid { // 不是队长不能匹配
			res.Result = msg.ErrCode_EquipStage_Team_Master_Match
			return
		}

		team = FindEquipStageTeam(p.GetUserId())
		if team == nil {
			log.Error("team is nil", zap.Uint64("uid", p.GetUserId()))
			res.Result = msg.ErrCode_SYSTEM_ERROR
			return
		}

		for _, mid := range team.Uids {
			mbr := player.FindByUserId(mid)
			if mbr == nil {
				continue
			}

			if err := CanLock(mbr, stage_t.UnlockCond); err != msg.ErrCode_SUCC {
				log.Error("equip stage team member stage locked", ZapUser(mbr))
				res.Result = msg.ErrCode_EquipStage_Lock
				return
			}
		}

	}
	if team.MatchState != EquipStage_UnMatch { // 已匹配
		log.Error("match when match", zap.Uint64("uid", p.GetUserId()))
		res.Result = msg.ErrCode_EquipStage_Team_Mached
		return
	}
	team.StageId = req.StageId
	team.StartMatch()
}

func EquipStageMatchCacnel(pid uint32, p *player.Player, req *msg.EquipStageMatchCancelReq) {
	res := &msg.EquipStageMatchCancelAck{Result: msg.ErrCode_SUCC}
	defer p.SendResponse(pid, res, res.Result)

	mid := GetEquipStageTeamMasterId(p.GetUserId())
	if mid != p.GetUserId() { // 不是队长不能取消匹配
		res.Result = msg.ErrCode_EquipStage_Team_Master_Match
		return
	}
	team := FindEquipStageTeam(mid)
	if team == nil {
		log.Error("team is nil", zap.Uint64("uid", p.GetUserId()), zap.Uint64("mid", mid))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	if team.MatchState != EquipStage_Matching { // 不是匹配中的队伍
		log.Error("team not matching", zap.Uint64("uid", p.GetUserId()), zap.Uint64("mid", mid))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}
	cancel_match(team.StageId, p.GetUserId(), true)
}

func EquipStageAccept(pid uint32, p *player.Player, req *msg.EquipStageAcceptReq) {
	res := &msg.EquipStageAcceptAck{
		Result: msg.ErrCode_SUCC,
	}
	defer p.SendResponse(pid, res, res.Result)

	mid := GetEquipStageTeamMasterId(p.GetUserId())
	team := FindEquipStageTeam(mid)
	if team == nil {
		log.Error("team is nil", zap.Uint64("uid", p.GetUserId()), zap.Uint64("mid", mid))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	if team.MatchState != EquipStage_Matched { // 不是匹配完成的队伍
		log.Error("team not matched", zap.Uint64("uid", p.GetUserId()), zap.Uint64("mid", mid))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	room := FindEquipStageRoom(team.RoomId)
	if room == nil {
		log.Error("room nil",
			zap.Uint64("uid", p.GetUserId()),
			zap.Uint64("mid", mid),
			zap.Uint64("rid", team.RoomId),
		)
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	room.BroadCast(&msg.EquipStageAcceptNtf{
		Accept: req.Accept,
		Uid:    p.GetUserId(),
	}, 0)
	if req.Accept {
		room.Ready(p.GetUserId())
	} else {
		DisEquipStageRoom(team.RoomId, true)
	}

}

func EquipStageLoad(pid uint32, p *player.Player) {
	res := &msg.EquipStageLoadAck{Result: msg.ErrCode_SUCC}
	defer p.SendResponse(pid, res, res.Result)

	mid := GetEquipStageTeamMasterId(p.GetUserId())
	team := FindEquipStageTeam(mid)
	if team == nil {
		log.Error("team is nil", zap.Uint64("uid", p.GetUserId()), zap.Uint64("mid", mid))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}
	if team.MatchState != EquipStage_Matched { // 不是匹配完成的队伍
		log.Error("team not matched", zap.Uint64("uid", p.GetUserId()), zap.Uint64("mid", mid))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	room := FindEquipStageRoom(team.RoomId)
	if room == nil {
		log.Error("room nil",
			zap.Uint64("uid", p.GetUserId()),
			zap.Uint64("mid", mid),
			zap.Uint64("rid", team.RoomId),
		)
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	if !room.IsReady {
		log.Error("room not ready",
			zap.Uint64("uid", p.GetUserId()),
			zap.Uint64("mid", mid),
			zap.Uint64("rid", team.RoomId),
		)
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	room.BroadCast(&msg.EquipStageLoadNtf{
		Uid: p.GetUserId(),
	}, 0)

	if room.LoadFinish(p.GetUserId()) {
		//CreateStageFight(p, room.StatgeID, 0, tools.Uint64ToBytes(room.ID))
		start_msg := &msg.FsStartFight{
			FightId: room.FightId,
		}
		SendToFight(p, 0, start_msg)
	}
}

func OnFsCreateEquipStageFight(uid uint64, ack *msg.FsCreateFightAck) {
	room := GetEquipStageTeamRoom(uid)
	if room == nil {
		log.Error("room nil", zap.Uint64("uid", uid))
		return
	}
	stage_t := template.GetMissionTemplate().GetMission(int(ack.StageId))
	if stage_t == nil {
		log.Error("stage nil", zap.Uint32("stage id", ack.StageId))
		return
	}

	room.FightId = ack.GetFightId()

	now := time.Now()
	deadline := now.Add(time.Minute * time.Duration(kenum.FightDeadLine))
	for r_uid := range room.ReadyUids {
		m_p := player.FindByUserId(r_uid)
		if m_p == nil {
			log.Error("member nil", zap.Uint64("member id", r_uid))
			continue
		}
		m_p.FightType = msg.BattleType(stage_t.Type)
		m_p.LastFightTime = now
		m_p.SetFsId(ack.GetFightSerId())
		m_p.UserData.Fight.FightStageId = int(ack.GetStageId())
		m_p.UserData.Fight.FightId = ack.GetFightId()
		m_p.IsSendEndFight = true
		m_p.SaveFight()

		createFightNtf := &msg.CreateFightNtf{
			Result:  msg.ErrCode_SUCC,
			StageId: ack.GetStageId(),
			MapType: ack.GetMapType(),
		}
		m_p.SendNotify(createFightNtf)

		fightExtra := m_p.MakeProtocolBase()
		extraBytes, err := fightExtra.Marshal()
		if err != nil {
			log.Error("zombie extra marshal err", zap.Error(err),
				zap.Uint64("accountId", r_uid), zap.Any("extra", fightExtra))
			return
		}

		bc := &BattleCache{
			FsId:     m_p.GetFsId(),
			BattleId: ack.FightId,
			Uid:      r_uid,
			DeadLine: deadline,
			StageId:  ack.StageId,
		}

		AddBattleCache(bc)
		SendFsEnterFight(m_p, ack.GetStageId(), ack.GetFightId(), extraBytes)
	}
}

func OnFsEquipStageFightResult(uid uint64, ntf *msg.FsFightResultNtf) {
	extraInfo := new(protocol.EquipStageFightResultExtra)
	if err := extraInfo.Unmarshal(ntf.GetExtra()); err != nil {
		log.Error("extraInfo.Unmarshal err", zap.Error(err), zap.Any("ntf", ntf))
		return
	}
	room := GetEquipStageTeamRoom(uid)
	if room == nil {
		log.Error("room nil", zap.Uint64("user id", uid))
		return
	}

	p := player.FindByUserId(uid)
	if p != nil {

		u_ntf := &msg.EquipStageFigthResultNtf{
			StageId: ntf.StageId,
			Victory: ntf.Victory,
		}

		if ntf.Victory {
			p.EquipStageItems = extraInfo.Items
			u_ntf.PickItems = extraInfo.Items

			stage_t := template.GetMissionTemplate().GetMission(int(ntf.StageId))
			if stage_t != nil {
				for id, num := range stage_t.NormalReward {
					AddItem(p.GetUserId(), uint32(id), int32(num), publicconst.EquipStage, true)
				}
			}

			if rec, ok := p.UserData.EquipStage.Records[ntf.StageId]; ok {
				rec.Num += 1
			} else {
				p.UserData.EquipStage.Records[ntf.StageId] = &model.EquipStageRecord{
					Num: 1,
				}
			}
			p.SaveEquipStage()
		} else {
			if room.EndUser(p.GetUserId()) {
				DisEquipStageRoom(room.ID, false)
			}
		}

		p.SendNotify(u_ntf)
	} else {
		if room.EndUser(uid) {
			DisEquipStageRoom(room.ID, false)
		}
	}

	DebugLog("equip stage fight result", zap.Uint64("user id", uid))
}

func EquipStageReward(pid uint32, p *player.Player) {
	res := &msg.EquipStageRewardPickAck{Result: msg.ErrCode_SUCC}
	defer p.SendResponse(pid, res, res.Result)

	if p.UserData.EquipStage.RewardNum >= p.UserData.EquipStage.RewardNum+2 {
		res.Result = msg.ErrCode_EquipStage_Reward_Count
		return
	}

	p.UserData.EquipStage.RewardNum++

	ntf_items := make([]*msg.IDNum, 0, len(p.EquipStageItems))
	var resItems = make([]*model.SimpleItem, 0)
	for id, num := range p.EquipStageItems {
		if num > 0 {
			items := AddItem(p.GetUserId(), id, int32(num), publicconst.EquipStagePick, true)
			resItems = append(resItems, items...)
		}
	}

	for _, item := range resItems {
		ntf_items = append(ntf_items, &msg.IDNum{Id: item.Id, Num: item.Num})
	}

	var ntf_item *msg.IDNum
	if len(ntf_items) > 0 {
		ntf_item = ntf_items[0]
		for _, item := range ntf_items[1:] {
			item_t := template.GetItemTemplate().GetItem(item.Id)
			if item_t != nil && item_t.BigType == uint32(msg.ItemType_Item_Type_Equip) {
				ntf_item = item
				break
			}
		}
	}

	if ntf_item != nil {
		broadcast_msg := &msg.EquipStageRewardPickNtf{
			Uid:       p.GetUserId(),
			PickItems: make(map[uint32]uint32),
		}
		broadcast_msg.PickItems[ntf_item.Id] = ntf_item.Num
	}

	p.EquipStageItems = make(map[uint32]uint32)

	room := GetEquipStageTeamRoom(p.GetUserId())
	if room != nil {
		if ntf_item != nil {
			broadcast_msg := &msg.EquipStageRewardPickNtf{
				Uid:       p.GetUserId(),
				PickItems: make(map[uint32]uint32),
			}
			broadcast_msg.PickItems[ntf_item.Id] = ntf_item.Num
			room.BroadCast(broadcast_msg, 0)
		}

		if room.EndUser(p.GetUserId()) {
			DisEquipStageRoom(room.ID, false)
		}
	}
	p.SaveEquipStage()
	UpdateTask(p, true, publicconst.TASK_COND_REWARD_EQUIP, 1) // 装备本累计领取奖励1次
	processHistoryData(p, publicconst.TASK_COND_REWARD_EQUIP, 0, 1)
}

func EquipStageBuyReward(pid uint32, p *player.Player, req *msg.EquipStageBuyRewardReq) {
	res := &msg.EquipStageBuyRewardAck{Result: msg.ErrCode_SUCC}
	res.Num = req.Num
	defer p.SendResponse(pid, res, res.Result)

	if req.Num <= 0 {
		log.Error("req num invalid", zap.Uint32("num", req.Num), ZapUser(p))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	if req.Num+p.UserData.EquipStage.RewardBuy > template.GetSystemItemTemplate().EquipStageBuyCount {
		res.Result = msg.ErrCode_EquipStage_Reward_Buy_Max
		return
	}

	cost := make(map[uint32]uint32)
	for i := uint32(0); i < req.Num; i++ {
		price_idx := p.UserData.EquipStage.RewardBuy + i
		if price_idx >= uint32(len(template.GetSystemItemTemplate().EquipStagePrice)) {
			price_idx = uint32(len(template.GetSystemItemTemplate().EquipStagePrice)) - 1
		}

		cost_item := template.GetSystemItemTemplate().EquipStagePrice[price_idx]
		cost[cost_item.Id] += cost_item.Num
		if !EnoughItem(p.GetUserId(), cost_item.Id, cost_item.Num) {
			res.Result = msg.ErrCode_NO_ENOUGH_ITEM
			return
		}
	}

	for k, v := range cost {
		if !EnoughItem(p.GetUserId(), k, v) {
			res.Result = msg.ErrCode_NO_ENOUGH_ITEM
			return
		}
	}

	for k, v := range cost {
		CostItem(p.GetUserId(), k, v, publicconst.EquipStageBuyCount, true)
	}

	p.UserData.EquipStage.RewardBuy += req.Num
	p.SaveEquipStage()
}

func EquipStageLeave(p *player.Player) bool {
	room := GetEquipStageTeamRoom(p.GetUserId())
	if room != nil {
		if room.EndUser(p.GetUserId()) {
			DisEquipStageRoom(room.ID, false)
			return true
		}
	}
	return false
}

func RefreshEquipStage(p *player.Player) {
	p.UserData.EquipStage.RewardNum = 0
	p.UserData.EquipStage.RewardBuy = 0
	p.SaveEquipStage()
}

func FindEquipStage(p *player.Player, id uint32) *model.EquipStageRecord {
	if v, ok := p.UserData.EquipStage.Records[id]; ok {
		return v
	}
	return nil
}

func ReqCreateEquipStageTeam(pid uint32, p *player.Player, req *msg.RequestCreateEquipStageTeam) {
	res := &msg.ResponseCreateEquipStageTeam{
		Result: msg.ErrCode_SUCC,
	}
	defer p.SendResponse(pid, res, res.Result)

	DebugLog("ReqCreateEquipStageTeam", zap.Any("req", req), ZapUser(p))
	stage_t := template.GetMissionTemplate().GetMission(int(req.StageId))
	if stage_t == nil {
		log.Error("ReqCreateEquipStageTeam stage cfg nil", zap.Any("req", req))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}
	// 解锁判定
	if err := CanLock(p, stage_t.UnlockCond); err != msg.ErrCode_SUCC {
		res.Result = msg.ErrCode_EquipStage_Lock
		return
	}

	var team *EquipStageTeam
	mid := GetEquipStageTeamMasterId(p.GetUserId())
	if mid == 0 { // 没有队伍
		team = CreateEquipStageTeam(req.StageId, p.GetUserId())
	} else { // 有队伍
		team = FindEquipStageTeam(mid)
		if team == nil {
			log.Error("ReqCreateEquipStageTeam team nil", zap.Any("req", req), ZapUser(p))
			res.Result = msg.ErrCode_SYSTEM_ERROR
			return
		}
		log.Error("ReqCreateEquipStageTeam repeat creat team", zap.Any("req", req))
	}
	res.TeamId = team.TeamId
}
