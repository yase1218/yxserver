package service

import (
	"gameserver/internal/enum"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tools"
	"math/rand"
	"msg"
	"sort"
	"sync/atomic"
	"time"

	"kernel/iface"

	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

type PeakFightRoom struct {
	id         uint64 // 房间id
	matchUid   uint64 // 根据哪个玩家设置段位匹配区间
	maxMatchId uint32 // 最大的段位id
	matchStart uint32 // 段位开始 battle_match_id
	matchEnd   uint32 // 段位结束 battle_match_id
	stageId    uint32 // 分配的关卡 按段位最高的来 房间所有人都用这一个

	startAt       time.Time                 // 开始时间
	hasRobot      bool                      // 是否有机器人
	robotId       uint64                    // 用于生成机器人唯一id
	users         map[uint64]*PeakFightUser // accountId||robotId->*PeakFightUser
	userNum       int32                     // 当前玩家数量(含机器人)
	realUserNum   int32                     // 真实玩家数量
	startFightIds []uint64                  // 已进入战斗人数
	rankInc       int32                     // 当前排名 从1开始,正着数
	rankDec       int32                     // 当前排名 从4开始,倒着数
	refMonsterId  int                       // 刷新bossId
	monsterId     int                       // bossId
	hasName       map[string]struct{}

	//tdaPvpOps []*tda.PvpOpUnit // 对手信息 (打点需要)

	broadcastMap       map[msg.PeakFightBroadcastType]uint32 // msg.PeakFightBroadcastType:【second|...】
	broadcastSecondMap map[uint32]uint32                     // 【second|...】:bid
	broadcastNum       int

	robotStartFight bool // 机器人开始计算进度标识
	done            chan struct{}

	end_matched          bool
	end_match_tm         uint64
	destroyed            bool
	destroy_tm           uint64
	robot_start_fight_tm uint64
	robot_calc           bool
	robot_calc_tm        uint64
	start_fight          bool
	start_fight_tm       uint64
	canAdd               bool
}

func CreatePeakFightRoom(id uint64, uid uint64, ms, me uint32, canAdd bool) *PeakFightRoom {
	now := time.Now()
	refMonsterId := template.GetBattlePassTemplate().GetCurSeason(tools.GetCurTime()).BOSSRefresh
	monsterRefCfg := template.GetMonsterRefreshTemplate().GetMonsterRefresh(refMonsterId)
	if monsterRefCfg == nil {
		log.Error("monster refresh cfg nil", zap.Int("refreshId", refMonsterId))
		return nil
	}
	return &PeakFightRoom{
		id:            id,
		matchUid:      uid,
		matchStart:    ms,
		matchEnd:      me,
		robotId:       0,
		users:         make(map[uint64]*PeakFightUser),
		startFightIds: make([]uint64, 0, 4),
		rankInc:       1,
		rankDec:       4,
		refMonsterId:  refMonsterId,
		monsterId:     int(monsterRefCfg.Monsters[0].Id),

		hasName:            make(map[string]struct{}),
		broadcastMap:       make(map[msg.PeakFightBroadcastType]uint32),
		broadcastSecondMap: make(map[uint32]uint32),
		done:               make(chan struct{}, 1),

		end_match_tm:         uint64(now.Add(time.Second * enum.Fight_Peak_Room_Timeout_Seconds).Unix()),
		destroy_tm:           uint64(now.Add(time.Second * enum.Fight_Peak_Room_Destroy_Seconds).Unix()),
		robot_start_fight_tm: uint64(now.Add(time.Second * 2).Unix()),
		robot_calc_tm:        uint64(now.Add(time.Second * enum.Fight_Peak_Room_Robot_Check_Seconds).Unix()),
		canAdd:               canAdd,
	}
}

func (r *PeakFightRoom) Stop(exitType msg.EndFightType) {
	for _, u := range r.users {
		r.UserExit(u.id, exitType)
	}
}

func (r *PeakFightRoom) on_second_ticker(now time.Time) {
	cur := uint64(now.Unix())
	if r.destroyed {
		return
	}
	if cur >= r.destroy_tm {
		r.Stop(msg.EndFightType_End_Fight_Time)
		r.destroyed = true
		return
	}

	if r.end_match_tm > 0 && cur >= r.end_match_tm && !r.end_matched {
		r.AddRobot()
		r.ChangeReadyToStart()
		r.end_matched = true
	}

	if r.robot_start_fight_tm > 0 && cur >= r.robot_start_fight_tm && !r.startAt.IsZero() {
		r.RobotStartFight()
	}

	if r.robot_calc_tm > 0 && cur >= r.robot_calc_tm && r.hasRobot && r.robotStartFight {
		r.CalcRobotProgress()
	}

	if r.broadcastNum > 0 {
		seconds := uint32(time.Now().Sub(r.startAt).Seconds()) * 100

		if _, ok := r.broadcastMap[msg.PeakFightBroadcastType_Broadcast_Time]; ok {
			if bid, sec_ok := r.broadcastSecondMap[seconds]; sec_ok && bid != 0 {
				for _, u := range r.users {
					if u.isSettlement || u.isRobot {
						continue
					}

					p := player.FindByUserId(u.id)
					if p == nil {
						log.Error("player nil", zap.Uint64("uid", u.id))
						continue
					}

					p.SendNotify(&msg.FightBroadcastNtf{
						AccountId: u.id,
						IsEmote:   false,
						Id:        bid,
					})
				}
				r.broadcastSecondMap[seconds] = 0
			}
		}
	}

	if r.start_fight_tm > 0 && cur >= r.start_fight_tm && !r.start_fight {
		r.StartFight()
	}
}

func (r *PeakFightRoom) CalcRobotProgress() {
	checkHasRobot := false
	delRobotIds := make([]uint64, 0, 3)
	for _, u := range r.users {
		if u.isSettlement || u.isDel || !u.isRobot {
			continue
		}

		if u.robotStartAt.IsZero() || time.Now().Before(u.robotStartAt) {
			continue
		}
		speed := u.robotSpeed
		if u.ranking > u.lastRanking {
			// 检查参数合法性
			if u.robotAddRand[0] > u.robotAddRand[1] {
				log.Error("robotAddRand range invalid",
					zap.Float64("min", u.robotAddRand[0]),
					zap.Float64("max", u.robotAddRand[1]))
				continue
			}

			randomFloat := u.robotAddRand[0] +
				rand.Float64()*(u.robotAddRand[1]-u.robotAddRand[0])
			threshold := 1.0 + rand.Float64()*99.0
			if randomFloat >= threshold {
				speed += u.robotSpeedAdd / 100
			}
		}
		battleMatchCfg := template.GetBattleMatchTemplate().GetCfg(u.matchId)
		if battleMatchCfg == nil {
			log.Error("battleMatchCfg nil", zap.Uint32("match id", u.matchId))
			continue
		}
		switch u.stairs {
		case enum.Fight_Peak_Stair_Pick_Item:
			ret := (speed * u.robotHasten / float64(u.robotPassTime[0])) * enum.Fight_Peak_Robot_Speed * 100
			u.stairs_data += ret
			u.stairs_data_msg = uint32(u.stairs_data * float64(battleMatchCfg.Boss) / 100)
			if u.stairs_data_msg >= uint32(u.stairs_data_max) {
				u.stairs_data_msg = uint32(u.stairs_data_max)
				if battleMatchCfg.ClearanceEnergy == 0 {
					delRobotIds = append(delRobotIds, u.id)
					checkHasRobot = true
					r.rankInc++
					u.status = enum.Fight_Peak_Status_Finish
				} else {
					u.stairs = enum.Fight_Peak_Stair_Attack_Boss
					u.stairs_data_max = int64(battleMatchCfg.ClearanceEnergy)
					u.robotStartAt = time.Now()
				}
			}
		case enum.Fight_Peak_Stair_Attack_Boss:
			ret := (speed * u.robotHasten / float64(u.robotPassTime[0]+u.robotPassTime[1])) * enum.Fight_Peak_Robot_Speed * 100

			u.stairs_data += ret
			u.stairs_data_msg = uint32(u.stairs_data * float64(battleMatchCfg.ClearanceEnergy) / 200)

			if u.stairs_data_msg >= uint32(u.stairs_data_max) {
				u.stairs_data_msg = uint32(u.stairs_data_max)
				delRobotIds = append(delRobotIds, u.id)
				checkHasRobot = true
			}
		}
	}
	r.SortRanking()
	r.SendPeakFightProgressNtf()

	for _, id := range delRobotIds {
		u := r.GetUser(id)
		if u != nil {
			u.isSettlement = true
		}

		r.DelUser(id, true, false)
	}

	if checkHasRobot {
		hasRobot := false
		for _, u := range r.users {
			if u.isRobot && !u.isDel {
				hasRobot = true
				break
			}
		}
		r.hasRobot = hasRobot
	}
}

func (r *PeakFightRoom) RobotStartFight() {
	if !r.robotStartFight {
		for _, roomUser := range r.users {
			if !roomUser.isRobot {
				continue
			}
			delay := time.Duration(rand.Int63n(5000)+1000) * time.Millisecond
			//log.Debug("robot start", zap.Int64("accountId", ru.id))
			roomUser.robotStartAt = time.Now().Add(delay)
		}

		r.robotStartFight = true
	}
}

func (r *PeakFightRoom) destroy(exitType msg.EndFightType) {
	for _, u := range r.users {
		if !u.isRobot &&
			!u.isSettlement &&
			!u.isDel {
			r.UserExit(u.id, exitType)
		}
	}
}

func (r *PeakFightRoom) UserExit(uid uint64, exitType msg.EndFightType) error {
	roomUser := r.GetUser(uid)
	if roomUser == nil {
		log.Error("room user not found", zap.Uint64("roomId", r.id), zap.Uint64("uid", uid))
		return errcode.ERR_PARAM
	}

	//if !r.startAt.IsZero() && !roomUser.isDel {
	if !roomUser.isDel {
		switch exitType {
		case msg.EndFightType_End_Fight_Dead:
			fallthrough
		case msg.EndFightType_End_Fight_Exit:
			fallthrough
		case msg.EndFightType_End_Fight_Time:
			roomUser.status = enum.Fight_Peak_Status_Dead

			rank := int(atomic.LoadInt32(&r.rankDec))
			roomUser.ranking = rank
			roomUser.isSettlement = true
			atomic.AddInt32(&r.rankDec, -1)
		}

		if err := r.Settlement(roomUser, roomUser.ranking); err != nil {
			return err
		}
		log.Info("UserExit", zap.Uint64("accountId", uid), zap.Int("rank", roomUser.ranking),
			zap.Int32("rankDec", r.rankDec), zap.Int("exitType", int(exitType)))
	}

	return nil
}

func (r *PeakFightRoom) GetUser(id uint64) *PeakFightUser {
	if user, ok := r.users[id]; ok {
		return user
	}
	return nil
}

func (r *PeakFightRoom) Settlement(roomUser *PeakFightUser, rank int) error {
	if roomUser.isDel {
		return nil
	}

	if rank == 0 {
		rank = 4
	}

	log.Info("peak fight settlement", zap.Uint64("accountId", roomUser.id), zap.Int("rank", rank),
		zap.Uint64("roomId", r.id), zap.Time("startAt", r.startAt), zap.Int32("rankDec", r.rankDec))

	battleMatchCfg := template.GetBattleMatchTemplate().GetCfg(roomUser.matchId)
	if battleMatchCfg == nil {
		log.Error("battle match cfg nil", zap.Uint32("battle match id", roomUser.matchId))
		return errcode.ERR_CONFIG_NIL
	}

	p := player.FindByUserId(roomUser.id)
	if p == nil {
		log.Error("player nil", zap.Uint64("uid", roomUser.id))
		return errcode.ERR_USER_DATA_NOT_FOUND
	}
	var (
		cups   int
		getExp uint32
		oldCup = uint32(p.UserData.PeakFight.Cup)
	)

	//winOrLose := "1"
	if rank != 1 {
		//winOrLose = "2"
		p.UserData.PeakFight.Streak = 0
	} else {
		p.UserData.PeakFight.Streak++
	}

	// 添加奖杯
	if battleMatchCfg.WinCup[rank-1] != 0 {
		winCupNum := battleMatchCfg.WinCup[rank-1]
		cups = int(oldCup) + winCupNum

		if winCupNum > 0 {
			for uint32(cups) >= template.GetBattleMatchTemplate().GetCfg(roomUser.matchId).CupEnd {
				if template.GetBattleMatchTemplate().GetCfg(roomUser.matchId+1) == nil {
					break
				}
				roomUser.matchId++
			}
		} else {
			for uint32(cups) < template.GetBattleMatchTemplate().GetCfg(roomUser.matchId).CupStart {
				if template.GetBattleMatchTemplate().GetCfg(roomUser.matchId-1) == nil {
					break
				}
				roomUser.matchId--
			}
		}
	}

	SendPeakFightNtf(p)
	UpdateCommonRankInfo(p, uint32(cups), template.SessionRank)

	activity := getActivityByType(p, uint32(publicconst.TaskPass))
	if activity == nil {
		log.Error("activity task pass nil", zap.Uint64("uid", roomUser.id))
		return errcode.ERR_STANDARD_ERR
	}
	passTaskOldLv := activity.ActDatas[0].Value2
	passTaskOldExp := activity.ActDatas[0].Value1
	// 添加通行证经验
	if battleMatchCfg.WinExp[rank-1] > 0 {
		maxCfg := template.GetTaskPassTemplate().GetMaxGrade(activity.ActId)
		if activity.ActDatas[0].Value2 != int64(maxCfg.Level) {
			activity.ActDatas[0].Value1 += battleMatchCfg.WinExp[rank-1]
			activity.ActDatas[0].Value2 = int64(template.GetTaskPassTemplate().GetGradeByExp(activity.ActId, activity.ActDatas[0].Value1))
			activity.ActDatas[0].UpdateTime = tools.GetCurTime()
			getExp = battleMatchCfg.WinExp[rank-1]
			p.SaveAccountActivity()
		}
	}

	// 添加胜利奖励
	if battleMatchCfg.WinRewardMap[rank] != nil {
		var notifyItems []uint32
		if len(battleMatchCfg.WinRewardMap[rank]) > 0 {
			for _, item := range battleMatchCfg.WinRewardMap[rank] {
				addItems := AddItem(p.GetUserId(), item.ItemId, int32(item.ItemNum), publicconst.PeakFight, false)
				notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
			}
			updateClientItemsChange(p.GetUserId(), notifyItems)
		}
	}

	if cfg := template.GetBattleMatchTemplate().GetCfg(roomUser.matchId); cfg != nil {
		UpdateTask(p, true, publicconst.TASK_COND_PEAK_FIGHT_RANK, cfg.Id)
	}

	// 推送
	p.SendNotify(&msg.PeakFightEndNtf{
		Result:         msg.ErrCode_SUCC,
		StageId:        r.stageId,
		Rank:           uint32(rank),
		BattleMatchId:  roomUser.matchId,
		Cup:            int32(oldCup),
		CupGet:         int32(battleMatchCfg.WinCup[rank-1]),
		PassTaskLv:     uint32(passTaskOldLv),
		PassTaskExp:    uint32(passTaskOldExp),
		PassTaskLvNew:  uint32(activity.ActDatas[0].Value2),
		PassTaskExpNew: activity.ActDatas[0].Value1,
		PassTaskExpGet: getExp,
		Streak:         p.UserData.PeakFight.Streak,
	})

	RecordUserPeakFightInfo(p, roomUser.matchId, 0, uint32(cups))
	LeaveFight(p)

	// tda
	// tda.TdaPvpEnd(p.ChannelId, p.TdaCommonAttr, strconv.Itoa(rank), "正常结束",
	// 	uint32(battleMatchCfg.WinCup[rank-1]), oldCup+uint32(battleMatchCfg.WinCup[rank-1]),
	// 	uint32(time.Now().Sub(r.startAt).Seconds()), r.tdaPvpOps)

	r.DelUser(roomUser.id, false, false)

	return nil
}

func (r *PeakFightRoom) DelUser(id uint64, isRobot, isDel bool) bool {
	//log.Debug("del user", zap.Int64("id", id), zap.Bool("is robot", isRobot), zap.Bool("is del", isDel))
	if isDel {
		delete(r.users, id)
		r.userNum--
		r.realUserNum--

		log.Info("del user", zap.Int32("user num", r.userNum), zap.Int32("real user num", r.realUserNum),
			zap.Uint64("room id", r.id), zap.Uint64("uid", id))
		// 没人了,删除房间
		if r.realUserNum == 0 ||
			(r.userNum == 1 && r.realUserNum == 1) {
			r.destroy(msg.EndFightType_End_Fight_Time)

			DelPFReadyRoom(r.id)
			DelPFStartRoom(r.id)
		} else {
			if r.matchUid == id {
				// 重新设置匹配区间
				for _, v := range r.users {
					if v.isRobot {
						continue
					}
					battle_cfg := template.GetBattleMatchTemplate().GetCfg(v.matchId)
					if battle_cfg == nil {
						continue
					}
					r.matchStart = battle_cfg.MatchRangeStart
					r.matchEnd = battle_cfg.MatchRangeEnd
					break
				}
			}
		}
		return true
	} else {
		roomUser, ok := r.users[id]
		//_, ok := r.users.Load(id)
		if ok {
			if roomUser.isDel {
				return true
			}

			//roomUser.isSettlement = true
			roomUser.isDel = true
			r.userNum--
			if !isRobot {
				r.realUserNum--
			}
			//fmt.Println(r.userNum)

			log.Info("del user", zap.Int32("user num", r.userNum), zap.Int32("real user num", r.realUserNum),
				zap.Uint64("room id", r.id), zap.Uint64("accountId", id))
			// 没人了,删除房间
			if r.realUserNum == 0 ||
				(r.userNum == 1 && r.realUserNum == 1) {
				r.destroy(msg.EndFightType_End_Fight_Time)

				DelPFReadyRoom(r.id)
				DelPFStartRoom(r.id)
			}
		}
	}

	return false
}

func (r *PeakFightRoom) AddRobot() {
	if r.userNum < enum.Fight_Peak_Room_User_Num_Max {
		// 人数不足,添加机器人
		battleMatchCfg := template.GetBattleMatchTemplate().GetCfg(r.maxMatchId)
		if battleMatchCfg == nil {
			log.Error("battle match cfg nil", zap.Uint32("id", r.matchStart))
			return
		}

		loop_cnt := 0
		safe_cnt := 1000
		needNum := enum.Fight_Peak_Room_User_Num_Max - int(r.userNum)
		for i := 0; i < needNum; i++ {
			loop_cnt++
			if loop_cnt > safe_cnt {
				log.Error("safe cnt", zap.Int("safe_cnt", safe_cnt))
				break
			}

			randRobotId := battleMatchCfg.RandRobot()
			robotCfg := template.GetRobotTemplate().GetPeakFightCfg(randRobotId)
			if robotCfg == nil {
				log.Error("robot cfg nil", zap.Uint32("id", randRobotId))
				i--
				continue
			}

			r.robotId++

			nickId := template.GetRandomNameTemplate().RandOne()
			if _, ok := r.hasName[nickId]; ok {
				nickId = template.GetRandomNameTemplate().RandOne()
			}
			// r.hasName[nickId] = struct{}{}

			// nickId := template.GetRandomNameTemplate().RandOne()
			// name := template.GetLanguageTemplate().GetContent(nickId)

			newRobot := &PeakFightUser{
				id:      r.robotId,
				isRobot: true,
				//stageId:         battleMatchCfg.RandStage(),
				matchId:         uint32(rand.Intn(int(r.matchEnd-r.matchStart+1))) + r.matchStart,
				stairs:          enum.Fight_Peak_Stair_Pick_Item,
				stairs_data:     0,
				stairs_data_max: int64(battleMatchCfg.Boss),
				stairs_data_msg: 0,

				nick:      nickId,
				head:      robotCfg.RandHead(),
				headFrame: robotCfg.RandHeadFrame(),
				ship:      robotCfg.RandShip(),

				robotId:       randRobotId,
				robotSpeed:    robotCfg.RobotExtra[0][0],
				robotHasten:   robotCfg.RobotExtra[1][0] + rand.Float64()*(robotCfg.RobotExtra[1][1]-robotCfg.RobotExtra[1][0]),
				robotSpeedAdd: robotCfg.RobotExtra[2][0],
				robotAddRand:  robotCfg.RobotExtra[3],
				robotPassTime: []uint32{
					robotCfg.RandPassTime(1),
					robotCfg.RandPassTime(2),
				},
			}
			if !r.AddUser(r.robotId, newRobot) {
				//log.Error("add robot fail")
				i--
				continue
			}
			r.startFightIds = append(r.startFightIds, r.robotId)
		}

		r.hasRobot = true
	}
}

func (r *PeakFightRoom) AddUser(id uint64, newUser *PeakFightUser) bool {
	battleMatchCfg := template.GetBattleMatchTemplate().GetCfg(newUser.matchId)
	if battleMatchCfg == nil {
		log.Error("battle match cfg nil", zap.Uint32("battle match id", newUser.matchId))
		return false
	}

	r.users[id] = newUser

	r.userNum++
	r.SendPeakFightMatchProgressNtf()
	//fmt.Println(r.userNum)

	// 延迟10毫秒，不然会少一条匹配数量消息 todo 消息底层需要优化，不然解决不了
	//time.Sleep(10 * time.Millisecond)

	// 匹配成功,房间全是玩家
	if r.userNum == enum.Fight_Peak_Room_User_Num_Max {
		r.ChangeReadyToStart()
		r.SendPeakFightMatchNtf()
	}

	return true
}

func (r *PeakFightRoom) ChangeReadyToStart() {
	for _, u := range r.users {
		if !u.isRobot && !u.isCost {
			p := player.FindByUserId(u.id)
			if p == nil {
				log.Error("player nil", zap.Uint64("uid", u.id))
				continue
			}
			if p.UserData.PeakFight.FreeTimes < template.GetSystemItemTemplate().PeakFightResetCostNum {
				p.UserData.PeakFight.FreeTimes++
				p.SavePeakFight()
			} else {
				itemId := template.GetSystemItemTemplate().PeakFightCost[0]
				num := template.GetSystemItemTemplate().PeakFightCost[1]
				CostItem(u.id, itemId,
					num, publicconst.PeakFightMatch, true)
			}
			u.isCost = true
		}
	}

	SetPFStartRoom(r.id, r)
	DelPFReadyRoom(r.id)
}

func (r *PeakFightRoom) SendPeakFightMatchNtf() {
	battleMatchCfg := template.GetBattleMatchTemplate().GetCfg(r.maxMatchId)
	if battleMatchCfg == nil {
		log.Error("battle match cfg nil", zap.Uint32("id", r.matchStart))
		return
	}

	r.stageId = battleMatchCfg.RandStage()

	// 先组织房间玩家信息
	peakFightMatchUnitSlice := make([]*msg.PeakFightMatchNtf_PeakFightMatchUnit, 0, 4)

	for _, u := range r.users {
		peakFightMatchUnitSlice = append(peakFightMatchUnitSlice, &msg.PeakFightMatchNtf_PeakFightMatchUnit{
			AccountId:     int64(u.id),
			Nick:          u.nick,
			Head:          u.head,
			HeadFrame:     u.headFrame,
			ShipId:        u.ship,
			BattleMatchId: u.matchId,
			IsRobot:       u.isRobot,
		})
	}

	for _, u := range r.users {
		if u.isRobot {
			continue
		}
		p := player.FindByUserId(u.id)
		if p == nil {
			log.Error("player nil", zap.Uint64("uid", u.id))
			continue
		}
		p.SendNotify(&msg.PeakFightMatchNtf{
			StageId:   r.stageId,
			Users:     peakFightMatchUnitSlice,
			FreeTimes: p.UserData.PeakFight.FreeTimes,
		})
	}
}

func (r *PeakFightRoom) SendPeakFightMatchProgressNtf() {
	r.SendRoomUserMsg(&msg.PeakFightMatchProgressNtf{
		Num: uint32(r.userNum),
	})
}

func (r *PeakFightRoom) SendRoomUserMsg(msg iface.IProtoMessage) {
	for _, u := range r.users {
		if u.isSettlement || u.isDel || u.isRobot {
			continue
		}

		p := player.FindByUserId(u.id)
		if p == nil {
			continue
		}
		p.SendNotify(msg)
	}
}

func (r *PeakFightRoom) SortRanking() {
	rankingSlice := make([]*PeakFightUser, 0, 4)

	for _, u := range r.users {
		rankingSlice = append(rankingSlice, u)
	}

	sort.Slice(rankingSlice, func(i, j int) bool {
		return rankingSlice[i].stairs > rankingSlice[j].stairs ||
			(rankingSlice[i].stairs == rankingSlice[j].stairs && rankingSlice[i].stairs_data > rankingSlice[j].stairs_data)
	})

	// 结算的排名固定
	for k, v := range rankingSlice {
		if v.ranking == 0 {
			continue
		}
		vRank := v.ranking - 1
		if v.isSettlement && vRank != k {
			rankingSlice[k], rankingSlice[vRank] = rankingSlice[vRank], rankingSlice[k]
		}
	}
	for _, u := range r.users {
		for rsr, ranking := range rankingSlice {
			if u.id == ranking.id {
				if u.ranking != rsr+1 {
					u.lastRanking = u.ranking
					u.ranking = rsr + 1
				}
			}
		}
	}
}

func (r *PeakFightRoom) SendPeakFightProgressNtf() {
	// 先组织房间玩家信息
	peakFightProgressUnitSlice := make([]*msg.PeakFightProgressNtf_PeakFightProgressUnit, 0, 4)
	for _, u := range r.users {

		unit := &msg.PeakFightProgressNtf_PeakFightProgressUnit{
			AccountId:   int64(u.id),
			Stairs:      uint32(u.stairs),
			Progress:    uint32(u.stairs_data),
			ProgressMax: u.stairs_data_max,
			Ranking:     uint32(u.ranking),
			Status:      uint32(u.status),
		}
		if u.isRobot {
			unit.Progress = u.stairs_data_msg
		}
		peakFightProgressUnitSlice = append(peakFightProgressUnitSlice, unit)
	}

	ntf := &msg.PeakFightProgressNtf{
		Progress: peakFightProgressUnitSlice,
	}
	r.SendRoomUserMsg(ntf)
}

func (r *PeakFightRoom) PickItem(uid uint64, num uint32) error {
	p := player.FindByUserId(uid)
	if p == nil {
		return errcode.ERR_USER_DATA_NOT_FOUND
	}
	u := r.GetUser(uid)
	if u == nil {
		log.Error("room user not found", zap.Uint64("roomId", r.id), zap.Uint64("uid", uid))
		return errcode.ERR_PARAM
	}

	if u.isDel {
		return nil
	}
	r.RobotStartFight()

	isEnd := false

	if u.broadcastNum > 0 {
		_, ok := u.broadcastMap[msg.PeakFightBroadcastType_Broadcast_Energy]
		if ok {
			for i := u.stairs_data; i <= u.stairs_data+float64(num); i++ {
				idx := uint32(i)
				if bid, ok := u.broadcastDataMap[idx]; ok && bid != 0 {
					ntf := &msg.FightBroadcastNtf{
						AccountId: uid,
						IsEmote:   false,
						Id:        bid,
					}
					r.SendBroadcast(ntf)
					u.broadcastDataMap[idx] = 0
				}
			}
		}
	}

	u.stairs_data += float64(num)

	switch u.stairs {
	case enum.Fight_Peak_Stair_Pick_Item:
		// 如果到达max 刷boss 更新阶段
		if u.stairs_data >= float64(u.stairs_data_max) {
			u.stairs_data = float64(u.stairs_data_max)

			battleMatchCfg := template.GetBattleMatchTemplate().GetCfg(u.matchId)
			if battleMatchCfg == nil {
				log.Error("battle match cfg nil", zap.Uint32("battle match id", u.matchId), zap.Uint64("uid", uid))
				return errcode.ERR_CONFIG_NIL
			}
			if battleMatchCfg.ClearanceEnergy == 0 {
				isEnd = true
			} else {
				//roomUser.stairs_data = 0
				u.stairs = enum.Fight_Peak_Stair_Attack_Boss

				monsterCfg := template.GetMonsterTemplate().GetMonster(r.monsterId)
				if monsterCfg == nil {
					log.Error("monster cfg nil", zap.Int("monsterId", r.monsterId))
					return errcode.ERR_CONFIG_NIL
				}
				p := player.FindByUserId(uid)
				if p == nil {
					log.Error("player nil", zap.Uint64("uid", uid))
					return errcode.ERR_USER_DATA_NOT_FOUND
				}
				stageCfg := template.GetMissionTemplate().GetMission(p.UserData.Fight.FightStageId)
				if stageCfg == nil {
					log.Error("stage cfg nil", zap.Int("stageId", p.UserData.Fight.FightStageId))
					return errcode.ERR_CONFIG_NIL
				}

				// 20250507 boss血量改为拾取数量
				//roomUser.stairs_data_max = int64(float32(monsterCfg.GetInitHp()) * (stageCfg.HpRate / 100))

				u.stairs_data_max = int64(battleMatchCfg.ClearanceEnergy)

				if r.refMonsterId != 0 {
					SendToFight(p, 0, &msg.FsRefreshMonsterNtf{
						FightId:           p.UserData.Fight.FightId,
						StageMonsterRefId: []uint32{uint32(r.refMonsterId)},
						NewAt:             true,
					})
				}
			}
		}
	case enum.Fight_Peak_Stair_Attack_Boss:
		if u.stairs_data >= float64(u.stairs_data_max) {
			u.stairs_data = float64(u.stairs_data_max)

			isEnd = true
		}
	}

	r.SortRanking()
	r.SendPeakFightProgressNtf()

	if isEnd {
		rank := int(atomic.LoadInt32(&r.rankInc))
		u.ranking = rank
		u.isSettlement = true
		u.status = enum.Fight_Peak_Status_Finish
		atomic.AddInt32(&r.rankInc, 1)
		if err := r.Settlement(u, u.ranking); err != nil {
			return err
		}
	}

	return nil
}

func (r *PeakFightRoom) SendBroadcast(ntf *msg.FightBroadcastNtf) {
	r.SendRoomUserMsg(ntf)
}

func (r *PeakFightRoom) SendPeakFightEnterNtf() {
	r.SendRoomUserMsg(&msg.PeakFightEnterNtf{AccountId: r.startFightIds})
}

func (r *PeakFightRoom) AddStartFight(id uint64) {

	r.startFightIds = append(r.startFightIds, id)
	r.SendPeakFightEnterNtf()

	if len(r.startFightIds) == enum.Fight_Peak_Room_User_Num_Max {
		r.start_fight_tm = uint64(time.Now().Add(time.Second * 5).Unix())
		// 延迟5秒进入战斗
		// time.AfterFunc(5*time.Second, func() {
		// 	battlePassCfg := template.GetBattlePassTemplate().GetCurSeason(tools.GetCurTime())
		// 	if battlePassCfg == nil {
		// 		log.Error("battle pass cfg nil", zap.Uint32("curTime", tools.GetCurTime()))
		// 		return
		// 	}

		// 	tmpUserBroadcastMap := make(map[msg.PeakFightBroadcastType]map[uint32]uint32)
		// 	if len(battlePassCfg.Broadcast) > 0 {
		// 		for _, bid := range battlePassCfg.Broadcast {

		// 			broadcastCfg := template.GetBroadcastTemplate().GetCfg(bid)
		// 			if broadcastCfg == nil {
		// 				log.Error("broadcast cfg nil", zap.Uint32("bid", bid))
		// 				continue
		// 			}
		// 			if tmpUserBroadcastMap[broadcastCfg.MsgType] == nil {
		// 				tmpUserBroadcastMap[broadcastCfg.MsgType] = make(map[uint32]uint32)
		// 			}
		// 			switch broadcastCfg.MsgType {
		// 			case msg.PeakFightBroadcastType_Broadcast_None:
		// 			case msg.PeakFightBroadcastType_Broadcast_Kill_Monster:
		// 				fallthrough
		// 			case msg.PeakFightBroadcastType_Broadcast_Refresh_Monster:
		// 				for _, v := range broadcastCfg.ScriptKey {
		// 					tmpUserBroadcastMap[broadcastCfg.MsgType][v] = bid
		// 				}
		// 			case msg.PeakFightBroadcastType_Broadcast_Time:
		// 				// 时间放房间里 广播
		// 				seconds := broadcastCfg.ScriptKey[0] * 100
		// 				r.broadcastMap[broadcastCfg.MsgType] = seconds

		// 				r.broadcastSecondMap[seconds] = bid

		// 				r.broadcastNum++
		// 				continue
		// 			case msg.PeakFightBroadcastType_Broadcast_Energy:
		// 				tmpUserBroadcastMap[broadcastCfg.MsgType][broadcastCfg.ScriptKey[0]] = bid
		// 			default:
		// 				log.Error("peak fight broadcast type not set", zap.Int("type", int(broadcastCfg.MsgType)), zap.Uint32("bid", bid))
		// 				continue
		// 			}
		// 		}
		// 	}

		// 	//r.tdaPvpOps = make([]*tda.PvpOpUnit, 0, r.realUserNum)

		// 	for _, u := range r.users {
		// 		if u.status == 0 {
		// 			u.status = enum.Fight_Peak_Status_Fighting

		// 			if !u.isRobot {
		// 				p := player.FindByUserId(u.id)
		// 				if p == nil {
		// 					log.Error("player nil", zap.Uint64("uid", u.id))
		// 					continue
		// 				}
		// 				start_msg := &msg.FsStartFight{
		// 					FightId: p.GetFightId(),
		// 				}
		// 				p.Send2Fight(0, start_msg)

		// 				for k, v := range tmpUserBroadcastMap {
		// 					for kk, vv := range v {
		// 						u.broadcastMap[k] = kk
		// 						u.broadcastDataMap[kk] = vv
		// 					}

		// 					u.broadcastNum += len(v)
		// 				}
		// 			}
		// 		}
		// 		// if !userInfo.isRobot {
		// 		// 	r.tdaPvpOps = append(r.tdaPvpOps, &tda.PvpOpUnit{
		// 		// 		AccountId:  strconv.FormatInt(userInfo.id, 10),
		// 		// 		Robot:      0,
		// 		// 		Kulu_id:    strconv.Itoa(int(userInfo.ship)),
		// 		// 		Kulu_class: strconv.Itoa(int(userInfo.shipClass)),
		// 		// 		Kulu_rank:  strconv.Itoa(int(userInfo.shipStarLv)),
		// 		// 		Kulu_level: userInfo.shipLv,
		// 		// 		Power:      userInfo.power,
		// 		// 	})
		// 		// } else {
		// 		// 	r.tdaPvpOps = append(r.tdaPvpOps, &tda.PvpOpUnit{
		// 		// 		AccountId: strconv.Itoa(int(userInfo.robotId)),
		// 		// 		Robot:     1,
		// 		// 	})
		// 		// }
		// 	}

		// 	// tda
		// 	// r.users.Range(func(key, value any) bool {
		// 	// 	userInfo := value.(*PeakFightUser)
		// 	// 	if !userInfo.isRobot {
		// 	// 		playerData := common.PlayerMgr.FindPlayerData(userInfo.id)
		// 	// 		if playerData != nil {
		// 	// 			tda.TdaPvpMatch(playerData.ChannelId, playerData.TdaCommonAttr, strconv.Itoa(userInfo.ranking), "", r.tdaPvpOps)
		// 	// 		}
		// 	// 	}
		// 	// 	return true
		// 	// })

		// 	r.startAt = time.Now()
		// })
	}
}

func (r *PeakFightRoom) StartFight() {
	r.start_fight = true
	r.robot_calc = true
	battlePassCfg := template.GetBattlePassTemplate().GetCurSeason(tools.GetCurTime())
	if battlePassCfg == nil {
		log.Error("battle pass cfg nil", zap.Uint32("curTime", tools.GetCurTime()))
		return
	}

	tmpUserBroadcastMap := make(map[msg.PeakFightBroadcastType]map[uint32]uint32)
	if len(battlePassCfg.Broadcast) > 0 {
		for _, bid := range battlePassCfg.Broadcast {

			broadcastCfg := template.GetBroadcastTemplate().GetCfg(bid)
			if broadcastCfg == nil {
				log.Error("broadcast cfg nil", zap.Uint32("bid", bid))
				continue
			}
			if tmpUserBroadcastMap[broadcastCfg.MsgType] == nil {
				tmpUserBroadcastMap[broadcastCfg.MsgType] = make(map[uint32]uint32)
			}
			switch broadcastCfg.MsgType {
			case msg.PeakFightBroadcastType_Broadcast_None:
			case msg.PeakFightBroadcastType_Broadcast_Kill_Monster:
				fallthrough
			case msg.PeakFightBroadcastType_Broadcast_Refresh_Monster:
				for _, v := range broadcastCfg.ScriptKey {
					tmpUserBroadcastMap[broadcastCfg.MsgType][v] = bid
				}
			case msg.PeakFightBroadcastType_Broadcast_Time:
				// 时间放房间里 广播
				seconds := broadcastCfg.ScriptKey[0] * 100
				r.broadcastMap[broadcastCfg.MsgType] = seconds

				r.broadcastSecondMap[seconds] = bid

				r.broadcastNum++
				continue
			case msg.PeakFightBroadcastType_Broadcast_Energy:
				tmpUserBroadcastMap[broadcastCfg.MsgType][broadcastCfg.ScriptKey[0]] = bid
			default:
				log.Error("peak fight broadcast type not set", zap.Int("type", int(broadcastCfg.MsgType)), zap.Uint32("bid", bid))
				continue
			}
		}
	}

	//r.tdaPvpOps = make([]*tda.PvpOpUnit, 0, r.realUserNum)

	for _, u := range r.users {
		if u.status == 0 {
			u.status = enum.Fight_Peak_Status_Fighting

			if !u.isRobot {
				p := player.FindByUserId(u.id)
				if p == nil {
					log.Error("player nil", zap.Uint64("uid", u.id))
					continue
				}
				start_msg := &msg.FsStartFight{
					FightId: p.GetBattleId(),
				}
				SendToFight(p, 0, start_msg)

				for k, v := range tmpUserBroadcastMap {
					for kk, vv := range v {
						u.broadcastMap[k] = kk
						u.broadcastDataMap[kk] = vv
					}

					u.broadcastNum += len(v)
				}
			}
		}
		// if !userInfo.isRobot {
		// 	r.tdaPvpOps = append(r.tdaPvpOps, &tda.PvpOpUnit{
		// 		AccountId:  strconv.FormatInt(userInfo.id, 10),
		// 		Robot:      0,
		// 		Kulu_id:    strconv.Itoa(int(userInfo.ship)),
		// 		Kulu_class: strconv.Itoa(int(userInfo.shipClass)),
		// 		Kulu_rank:  strconv.Itoa(int(userInfo.shipStarLv)),
		// 		Kulu_level: userInfo.shipLv,
		// 		Power:      userInfo.power,
		// 	})
		// } else {
		// 	r.tdaPvpOps = append(r.tdaPvpOps, &tda.PvpOpUnit{
		// 		AccountId: strconv.Itoa(int(userInfo.robotId)),
		// 		Robot:     1,
		// 	})
		// }
	}

	// tda
	// r.users.Range(func(key, value any) bool {
	// 	userInfo := value.(*PeakFightUser)
	// 	if !userInfo.isRobot {
	// 		playerData := common.PlayerMgr.FindPlayerData(userInfo.id)
	// 		if playerData != nil {
	// 			tda.TdaPvpMatch(playerData.ChannelId, playerData.TdaCommonAttr, strconv.Itoa(userInfo.ranking), "", r.tdaPvpOps)
	// 		}
	// 	}
	// 	return true
	// })

	r.startAt = time.Now()
}
