package service

import (
	"errors"
	"fmt"
	"gameserver/internal/enum"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"gameserver/internal/rdb"
	"kernel/tools"
	"msg"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

type (
	PeakFightUser struct {
		id              uint64                  // uid || robotId
		isRobot         bool                    // 是否机器人
		matchId         uint32                  // battle_match_id
		stairs          enum.FightPeakStairType // 阶段
		stairs_data     float64                 // 阶段对应数值 1阶段->能量 2阶段->boss血量
		stairs_data_max int64                   // 阶段对应顶值 和stairs_data比对 1阶段存上限,2阶段存boss血量
		stairs_data_msg uint32                  // 阶段对应值(用于消息 机器人用)
		isSettlement    bool                    // 是否结算
		isDel           bool                    // 是否删除

		nick       string // 昵称
		head       uint32 // 头像
		headFrame  uint32 // 头像框
		ship       uint32 // 出战机甲
		shipClass  uint32 // 机甲品阶
		shipLv     uint32 // 机甲等级
		shipStarLv uint32 // 机甲星级
		power      uint32 // 战力

		broadcastMap     map[msg.PeakFightBroadcastType]uint32 // msg.PeakFightBroadcastType:【second|...】
		broadcastDataMap map[uint32]uint32                     // 【second|...】:bid
		broadcastNum     int

		ranking     int                  // 排名
		lastRanking int                  // 上次排名
		isCost      bool                 // 是否扣除体力
		status      enum.FightPeakStatus // 状态

		robotId       uint32    // 机器人id
		robotSpeed    float64   // 机器人基础速度
		robotHasten   float64   // 机器人进度系数
		robotSpeedAdd float64   // 机器人加速百分比(计算时要除100)
		robotAddRand  []float64 // 机器人加速触发概率区间
		robotPassTime []uint32  // 机器人通关时间 0:怪物阶段 1:boss阶段
		robotStartAt  time.Time // 机器人每阶段开始时间
	}
)

var (
	pf_roomId     uint64                    // 用于生成房间唯一id
	pf_readyRooms map[uint64]*PeakFightRoom // 正在匹配的房间 roomId->*PeakFightRoom
	pf_startRooms map[uint64]*PeakFightRoom // 开始的房间 roomId->*PeakFightRoom

	pf_rankCache   *collection.Cache
	pf_meRankCache *collection.Cache
)

func init() {
	pf_readyRooms = make(map[uint64]*PeakFightRoom)
	pf_startRooms = make(map[uint64]*PeakFightRoom)
}

func InitFightPeak() {
	var (
		err error
	)

	if pf_rankCache, err = collection.NewCache(rankCacheTime); err != nil {
		panic(fmt.Errorf("peak fight rank init rank cache err:%s", err.Error()))
	}

	if pf_meRankCache, err = collection.NewCache(rankCacheTime); err != nil {
		panic(fmt.Errorf("peak fight rank init me_rank cache err:%s", err.Error()))
	}
}

func LoadPeakFight(pid uint32, p *player.Player) {
	res := &msg.PeakFightAck{
		Result:        msg.ErrCode_SUCC,
		BattleMatchId: p.UserData.PeakFight.BattleMatchId,
		Cup:           p.UserData.PeakFight.Cup,
		Streak:        p.UserData.PeakFight.Streak,
		FreeTimes:     p.UserData.PeakFight.FreeTimes,
	}
	defer p.SendResponse(pid, res, res.Result)
}

func PeakFightMatch(pid uint32, p *player.Player) {
	res := &msg.PeakFightMatchAck{
		Result: msg.ErrCode_SUCC,
	}
	defer p.SendResponse(pid, res, res.Result)

	match_id := p.UserData.PeakFight.BattleMatchId

	if p.UserData.PeakFight.FreeTimes >= template.GetSystemItemTemplate().PeakFightResetCostNum {
		if !EnoughItem(p.GetUserId(), template.GetSystemItemTemplate().PeakFightCost[0], template.GetSystemItemTemplate().PeakFightCost[1]) {
			res.Result = msg.ErrCode_Fight_Peak_Tickets_Not_Enough
			return
		}
	}

	match_cfg := template.GetBattleMatchTemplate().GetCfg(match_id)
	if match_cfg == nil {
		log.Error("battle match cfg nil", zap.Uint32("battle match id", match_id), ZapUser(p))
		res.Result = msg.ErrCode_CONFIG_NIL
		return
	}
	var room *PeakFightRoom
	if p.UserData.PeakFight.Cup == 0 { // 海边派对新手第一把必定匹配指定机器人
		pf_roomId++
		room = CreatePeakFightRoom(pf_roomId, p.GetUserId(),
			match_cfg.MatchRangeStart, match_cfg.MatchRangeEnd, false)
		if room == nil {
			log.Error("create room err", zap.Uint32("battle match id", match_id))
			res.Result = msg.ErrCode_SYSTEM_ERROR
			return
		}
		SetPFReadyRoom(room.id, room)
	} else {
		for _, v := range pf_readyRooms {
			if match_id >= v.matchStart &&
				match_id <= v.matchEnd && v.canAdd {
				room = v
				break
			}
		}

		if room == nil {
			pf_roomId++
			room = CreatePeakFightRoom(pf_roomId, p.GetUserId(),
				match_cfg.MatchRangeStart, match_cfg.MatchRangeEnd, true)
			if room == nil {
				log.Error("create room err", zap.Uint32("battle match id", match_id))
				res.Result = msg.ErrCode_SYSTEM_ERROR
				return
			}
			SetPFReadyRoom(room.id, room)
		}
	}
	if room.maxMatchId < match_cfg.Id {
		room.maxMatchId = match_cfg.Id
	}
	log.Debug("PeakFightMatch", zap.Uint64("uid", p.GetUserId()), zap.Uint32("match_id", match_id), zap.Uint32("Cup", p.UserData.PeakFight.Cup), zap.Uint64("roomId", room.id))
	shipConfig := template.GetShipTemplate().GetShip(p.UserData.BaseInfo.ShipId)
	if shipConfig == nil {
		log.Error("ship cfg nil", zap.Uint64("uid", p.GetUserId()), zap.Uint32("shipId", p.UserData.BaseInfo.ShipId))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}
	var playerShipInfo *model.Ship
	for i := 0; i < len(p.UserData.Ships.Ships); i++ {
		if p.UserData.Ships.Ships[i].Id == p.UserData.BaseInfo.ShipId {
			playerShipInfo = p.UserData.Ships.Ships[i]
			break
		}
	}
	if playerShipInfo == nil {
		log.Error("player ship info nil", zap.Uint64("uid", p.GetUserId()), zap.Uint32("shipId", p.UserData.BaseInfo.ShipId))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	if !room.AddUser(p.GetUserId(), &PeakFightUser{
		id:              p.GetUserId(),
		isRobot:         false,
		matchId:         match_id,
		stairs:          enum.Fight_Peak_Stair_Pick_Item,
		stairs_data:     0,
		stairs_data_max: int64(match_cfg.Boss),
		nick:            p.UserData.Nick,
		head:            p.UserData.HeadImg,
		headFrame:       p.UserData.HeadFrame,
		ship:            p.UserData.BaseInfo.ShipId,

		shipClass:        shipConfig.Rarity,
		shipLv:           playerShipInfo.Level,
		shipStarLv:       playerShipInfo.StarLevel,
		power:            p.UserData.BaseInfo.Combat,
		broadcastMap:     make(map[msg.PeakFightBroadcastType]uint32),
		broadcastDataMap: make(map[uint32]uint32),
	}) {
		log.Error("room add user err", zap.Uint64("roomId", room.id), zap.Uint64("uid", p.GetUserId()))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}
	room.realUserNum++

	p.UserData.PeakFight.RoomId = room.id
}

func PeakFightCancelMatch(pid uint32, p *player.Player) {
	res := &msg.PeakFightCancelMatchAck{
		Result: msg.ErrCode_SUCC,
	}
	defer func() {
		if pid > 0 {
			p.SendResponse(pid, res, res.Result)
		}
	}()

	if p.UserData.PeakFight.RoomId == 0 {
		log.Error("room id 0", ZapUser(p))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	room := GetPFReadyRoom(p.UserData.PeakFight.RoomId)
	if room == nil {
		log.Error("room nil", zap.Uint64("roomId", p.UserData.PeakFight.RoomId))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	if !room.DelUser(p.GetUserId(), false, true) {
		log.Error("del user nil", zap.Uint64("roomId", p.UserData.PeakFight.RoomId), zap.Uint64("uid", p.GetUserId()))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}
	room.SendPeakFightMatchProgressNtf()

	p.UserData.PeakFight.RoomId = 0
}

func SendPeakFightNtf(p *player.Player) {
	p.SendNotify(&msg.PeakFightNtf{
		BattleMatchId: p.UserData.PeakFight.BattleMatchId,
		Cup:           p.UserData.PeakFight.Cup,
		Streak:        p.UserData.PeakFight.Streak,
	})
}

func RecordUserPeakFightInfo(p *player.Player, matchId uint32, roomId uint64, cup uint32) {
	p.UserData.PeakFight.BattleMatchId = matchId
	p.UserData.PeakFight.RoomId = roomId
	p.UserData.PeakFight.Cup = cup

	if cup != 0 {
		if err := rdb.AddPeakFightRank(p.GetUserId(), float64(cup)); err != nil {
			log.Error("add peak fight rank err", zap.Error(err))
		}
	}
}

func GetPFReadyRoom(roomId uint64) *PeakFightRoom {
	if room, ok := pf_readyRooms[roomId]; ok {
		return room
	}
	return nil
}

func SetPFReadyRoom(roomId uint64, room *PeakFightRoom) {
	pf_readyRooms[roomId] = room
}

func DelPFReadyRoom(roomId uint64) {
	delete(pf_readyRooms, roomId)
}

// =================================Start Room==========================================

func GetPFStartRoom(roomId uint64) *PeakFightRoom {
	if room, ok := pf_startRooms[roomId]; ok {
		return room
	}
	return nil
}

func SetPFStartRoom(roomId uint64, room *PeakFightRoom) {
	pf_startRooms[roomId] = room
}

func DelPFStartRoom(roomId uint64) {
	delete(pf_startRooms, roomId)
}

func PFUserExit(p *player.Player, exitType msg.EndFightType) error {
	if p == nil {
		log.Error("player nil")
		return errcode.ERR_PARAM
	}
	if p.UserData.PeakFight.RoomId == 0 {
		log.Error("room nil", zap.Uint64("uid", p.GetUserId()))
		return errcode.ERR_PARAM
	}
	if readyRoom := GetPFReadyRoom(p.UserData.PeakFight.RoomId); readyRoom != nil {
		PeakFightCancelMatch(0, p)
		return nil
	} else {
		startRoom := GetPFStartRoom(p.UserData.PeakFight.RoomId)
		if startRoom == nil {
			log.Error("room nil", zap.Uint64("roomId", p.UserData.PeakFight.RoomId))
			return errcode.ERR_PARAM
		}
		return startRoom.UserExit(p.GetUserId(), exitType)
	}
}

func OnPFLogOut(p *player.Player) {
	if p.UserData.PeakFight.RoomId != 0 {
		PFUserExit(p, msg.EndFightType_End_Fight_Disconnect)
	}

	RecordUserPeakFightInfo(p, p.UserData.PeakFight.BattleMatchId, 0, p.UserData.PeakFight.Cup)
}

func PFPickItem(p *player.Player, itemMap map[uint32]uint32) error {
	if p.UserData.PeakFight.RoomId == 0 {
		log.Error("room nil", zap.Uint64("uid", p.GetUserId()))
		return nil
	}
	room := GetPFStartRoom(p.UserData.PeakFight.RoomId)
	if room == nil {
		log.Error("room nil", zap.Uint64("roomId", p.UserData.PeakFight.RoomId))
		return errcode.ERR_PARAM
	}

	finalAddNum := uint32(0)
	for itemId, itemNum := range itemMap {

		// log.Debug("PickItem, ", zap.Any("itemId", itemId), zap.Any("num", itemNum))

		itemCfg := template.GetItemTemplate().GetItem(itemId)
		if itemCfg == nil {
			log.Error("item cfg nil", zap.Uint32("itemId", itemId))
			return errcode.ERR_CONFIG_NIL
		}
		if itemCfg.BigType != publicconst.Item_BigType_Peak_Fight_energy {
			log.Error("pick item not peak fight module", zap.Uint32("itemId", itemId))
			return errcode.ERR_PARAM
		}
		finalAddNum += itemCfg.EffectArgs[0] * itemNum
	}

	if finalAddNum > 0 {
		if err := room.PickItem(p.GetUserId(), finalAddNum); err != nil {
			log.Error("pick item err", zap.Uint64("roomId", room.id), zap.Uint64("uid", p.GetUserId()))
			return err
		}
	}

	return nil
}

func EnterPeakFight(p *player.Player) error {
	if p.UserData.PeakFight.RoomId == 0 {
		log.Error("room nil", zap.Uint64("uid", p.GetUserId()))
		return errcode.ERR_PARAM
	}
	room := GetPFStartRoom(p.UserData.PeakFight.RoomId)
	if room == nil {
		log.Error("room nil", zap.Uint64("roomId", p.UserData.PeakFight.RoomId))
		return errcode.ERR_PARAM
	}
	room.AddStartFight(p.GetUserId())

	processHistoryData(p, publicconst.TASK_COND_PEAK_FIGHT_PK, 0, 1)

	UpdateTask(
		p, true, publicconst.TASK_COND_PEAK_FIGHT_PK, 1)

	return nil
}

func PeakFightCheck(p *player.Player, stageId uint32) msg.ErrCode {
	// 校验功能模块是否开放
	if err := FunctionOpen(p, publicconst.PEAK_FIGHT); err != msg.ErrCode_SUCC {
		return err
	}

	// 校验当前赛季是否开启（当前时间和表里的周时间）
	if !CheckIsOpen() {
		return msg.ErrCode_FUNCTION_NOT_OPEN
	}

	return msg.ErrCode_SUCC
}

func CheckIsOpen() bool {
	curTime := tools.GetCurTime()
	battlePassCfg := template.GetBattlePassTemplate().GetCurSeason(curTime)
	if battlePassCfg == nil {
		log.Error("battle pass cfg nil", zap.Uint32("curTime", curTime))
		return false
	}

	return true
}
func PFCreateFightAck(uid uint64, ack *msg.FsCreateFightAck) {
	p := player.FindByUserId(uid)
	if p != nil {
		fightExtra := p.MakeProtocolBase()
		extraBytes, err := fightExtra.Marshal()
		if err != nil {
			log.Error("zombie extra marshal err", zap.Error(err),
				zap.Uint64("uid", p.GetUserId()), zap.Any("extra", fightExtra))
			return
		}
		SendFsEnterFight(p, ack.GetStageId(), ack.GetFightId(), extraBytes)
	}
}

// 广播
func PFBroadcast(p *player.Player, req *msg.FightBroadcastReq) error {
	if p.UserData.PeakFight.RoomId == 0 {
		log.Error("room nil", zap.Uint64("uid", p.GetUserId()))
		return errcode.ERR_PARAM
	}
	room := GetPFStartRoom(p.UserData.PeakFight.RoomId)
	if room == nil {
		log.Error("room nil", zap.Uint64("roomId", p.UserData.PeakFight.RoomId))
		return errcode.ERR_PARAM
	}

	ntf := &msg.FightBroadcastNtf{
		AccountId: p.GetUserId(),
		Id:        req.GetBroadcastId(),
		IsEmote:   false,
	}
	room.SendBroadcast(ntf)

	return nil
}

func GmSetPeakFightCup(p *player.Player, cup uint32) msg.ErrCode {
	i := uint32(1)
	for {
		cfg := template.GetBattleMatchTemplate().GetCfg(i)
		if cfg == nil || cup < cfg.CupEnd {
			break
		}
		i++
	}
	p.UserData.PeakFight.BattleMatchId = i
	p.UserData.PeakFight.Cup = cup

	if cup != 0 {
		if err := rdb.AddPeakFightRank(p.GetUserId(), float64(cup)); err != nil {
			log.Error("add peak fight rank err", zap.Error(err))
		}
	}
	SendPeakFightNtf(p)
	return msg.ErrCode_SUCC
}

func PeakFightGetRank(pid uint32, p *player.Player, req *msg.PeakFightRankReq) {
	res := &msg.PeakFightRankAck{
		Result:   msg.ErrCode_SUCC,
		RankType: req.RankType,
	}
	defer p.SendResponse(pid, res, res.Result)

	switch req.GetRankType() {
	case msg.PeakFightRankType_PeakFight_Rank:
		makePeakFightRank(p, req, res)
	case msg.PeakFightRankType_PeakFight_Rank_Friend:
		makePeakFightFriendRank(p, req, res)
	default:
		log.Error("invalid rank type", zap.Int32("type", int32(req.RankType)), ZapUser(p))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

}

func makePeakFightRank(p *player.Player, req *msg.PeakFightRankReq, ack *msg.PeakFightRankAck) error {
	rankData, err := pf_rankCache.Take(req.GetRankType().String(), func() (any, error) {
		rankData, err := rdb.GetPeakFightRank()
		if err != nil {
			log.Error("get peak fight rank err", zap.Error(err), zap.Any("rankType", req.GetRankType()))
			ack.Result = msg.ErrCode_SYSTEM_ERROR
			return nil, err
		}
		pbRankSlice := make([]*msg.PeakFightRankAck_PeakFightRankUnit, 0, len(rankData))
		for i, rankDataUnit := range rankData {
			accountId := utils.StrToInt64(rankDataUnit.Member.(string))

			_, accBasic := GetPlayerBasic(uint64(accountId))
			if accBasic == nil {
				log.Error("player basic nil", zap.Int64("accountId", accountId))
				continue
			}

			pbRankUnit := &msg.PeakFightRankAck_PeakFightRankUnit{
				Rank:      uint32(i + 1),
				Name:      accBasic.Nick,
				Head:      accBasic.HeadImg,
				Score:     utils.Float64ToUint64(rankDataUnit.Score),
				Title:     accBasic.Title,
				HeadFrame: accBasic.HeadFrame,
				AccountId: accountId,
				ShipId:    accBasic.ShipId,
			}
			pbRankSlice = append(pbRankSlice, pbRankUnit)
		}
		return pbRankSlice, nil
	})
	if err != nil {
		log.Error("get peak fight rank err", zap.Error(err), zap.Any("rankType", req.GetRankType()))
		ack.Result = msg.ErrCode_SYSTEM_ERROR
		return errcode.ERR_PARAM
	}
	ack.RankList = rankData.([]*msg.PeakFightRankAck_PeakFightRankUnit)

	meRankData, err := pf_meRankCache.Take(fmt.Sprintf("%s:%d", req.GetRankType().String(), p.GetUserId()), func() (any, error) {
		meRank, err := rdb.GetUserPeakFightRanking(p.GetUserId())
		if err != nil {
			log.Error("get me rank err", zap.Error(err),
				zap.Uint64("uid", p.GetUserId()))

			ack.Result = msg.ErrCode_SYSTEM_ERROR
			return nil, nil
		}
		meScore, err := rdb.GetUserPeakFightRankScore(p.GetUserId())
		if err != nil {
			log.Error("get me rank err", zap.Error(err),
				zap.Uint64("uid", p.GetUserId()))
			ack.Result = msg.ErrCode_SYSTEM_ERROR
			return nil, nil
		}

		return &msg.PeakFightRankAck_PeakFightRankUnit{
			Rank:      uint32(meRank + 1),
			Name:      p.UserData.Nick,
			Head:      p.UserData.HeadImg,
			Score:     utils.Float64ToUint64(meScore),
			Title:     p.UserData.Title,
			HeadFrame: p.UserData.HeadFrame,
			AccountId: int64(p.GetUserId()),
			ShipId:    p.UserData.BaseInfo.ShipId,
		}, nil
	})
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Error("get me rank err", zap.Error(err), zap.Uint64("uid", p.GetUserId()))
		ack.Result = msg.ErrCode_SYSTEM_ERROR
		return errcode.ERR_PARAM
	}
	if meRankData != nil {
		ack.Me = meRankData.(*msg.PeakFightRankAck_PeakFightRankUnit)
	}

	return nil
}

func makePeakFightFriendRank(p *player.Player, req *msg.PeakFightRankReq, ack *msg.PeakFightRankAck) error {
	//if len(playerData.AccountFri.Friend) == 0 {
	//	return nil
	//}

	rankData, err := pf_rankCache.Take(fmt.Sprintf("%s_friend:%d", req.GetRankType().String(), p.GetUserId()), func() (any, error) {
		pbRankSlice := make([]*msg.PeakFightRankAck_PeakFightRankUnit, 0, len(p.UserData.FriendData.Friend)+1)

		meScore, err := rdb.GetUserPeakFightRankScore(p.GetUserId())
		if err != nil {
			return nil, err
		}
		pbRankSlice = append(pbRankSlice, &msg.PeakFightRankAck_PeakFightRankUnit{
			Score:     utils.Float64ToUint64(meScore),
			AccountId: int64(p.GetUserId()),
			Name:      p.UserData.Nick,
			Head:      p.UserData.HeadImg,
			Title:     p.UserData.Title,
			HeadFrame: p.UserData.HeadFrame,
			ShipId:    p.UserData.BaseInfo.ShipId,
		})

		for friend_uid := range p.UserData.FriendData.Friend {
			friScore, err := rdb.GetUserPeakFightRankScore(friend_uid)
			if err != nil && !errors.Is(err, redis.Nil) {
				continue
			}
			unit := &msg.PeakFightRankAck_PeakFightRankUnit{
				Score:     utils.Float64ToUint64(friScore),
				AccountId: int64(friend_uid),
			}

			_, accBasic := GetPlayerBasic(friend_uid)
			if accBasic == nil {
				log.Error("player basic nil", zap.Uint64("friend uid", friend_uid))
				continue
			}

			unit.Name = accBasic.Nick
			unit.Head = accBasic.HeadImg
			unit.Title = accBasic.Title
			unit.HeadFrame = accBasic.HeadFrame
			unit.ShipId = accBasic.ShipId

			pbRankSlice = append(pbRankSlice, unit)
		}
		return pbRankSlice, nil
	})
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Error("get peak fight friend rank err", zap.Error(err), zap.Any("rankType", req.GetRankType()))
		return errcode.ERR_PARAM
	}
	if rankData != nil {
		ack.RankList = rankData.([]*msg.PeakFightRankAck_PeakFightRankUnit)
		// 重新排序
		sortRankList(ack.RankList)
		for i, unit := range ack.RankList {
			unit.Rank = uint32(i + 1)
			if unit.AccountId == int64(p.GetUserId()) {
				ack.Me = unit
			}
		}
	}

	return nil
}

func sortRankList(rankList []*msg.PeakFightRankAck_PeakFightRankUnit) {
	sort.Slice(rankList, func(i, j int) bool {
		return rankList[i].Score > rankList[j].Score
	})
}

func update_peak_fight(now time.Time) {
	for _, r := range pf_readyRooms {
		r.on_second_ticker(now)
	}

	for _, r := range pf_startRooms {
		r.on_second_ticker(now)
	}
}
