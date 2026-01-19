package service

import (
	"context"
	"gameserver/internal/config"
	"gameserver/internal/game/builder"
	"gameserver/internal/game/charge"
	"gameserver/internal/game/player"
	"gameserver/internal/io_out"
	"kernel/kenum"
	"msg"
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	update_internal int64 = 10
	debug_ticker    int64 = 300
)

func Init(ctx context.Context) {
	InitFuncPreview()
	InitFightPeak()
	InitFiendOps(ctx)
	InitTaskCheck()
	InitGlobalMail(ctx)
	charge.InitOrderManager()
	InitBattlePassWeekField()
	log.Info("service init succ")
}

func OnSecondTicker(now time.Time) {
	if config.Conf.IsDebug() {
		if now.Unix()%debug_ticker == 0 {
			DebugLog("main logic running")
		}
	}
	update_equip_stage(now)
	update_players(now)
	update_peak_fight(now)
	update_zero_time_activate(now)
	UpdateChat(now)
}

func update_players(now time.Time) {
	for _, p := range player.AllPlayers() {
		update_player(p, now)
	}
}

func update_player(p *player.Player, now time.Time) {
	now_unix := now.Unix()
	if now_unix-p.LastTick < update_internal {
		return
	}

	if !utils.IsSameDay(p.LastTick, now_unix) {
		OnCrossDay(p)
	}

	if !utils.IsSameDay(now_unix, p.UserData.ResourcesPass.LastResetTime.Unix()) {
		handleResetPlayerResourcePass(p)
	}

	p.LastTick = now_unix

	ResetData(p, now)
	recoveryAp(p, now)
	// ClearPlayerRankInfo()
}

func PostPanic(info string) {
	if !config.Conf.Monitor {
		return
	}
	post_msg := &msg.MonitorMsg{
		Id:         config.Conf.ServerId,
		Type:       msg.Monitortype_MT_Panic,
		ServerName: "game",
		Data:       info,
	}

	sub := "game.monitor"
	//nats.Publish(sub, post_msg)
	io_out.Push(&io_out.OutMsg{
		Subject: sub,
		Msg:     post_msg,
	})
	time.Sleep(time.Second)
}

func DebugLog(info string, filelds ...zap.Field) {
	if config.Conf.IsDebug() {
		log.Debug("~~~~~~"+info, filelds...)
	}
}

func OnFsCreateFight(uid uint64, ack *msg.FsCreateFightAck) {
	stage_t := template.GetMissionTemplate().GetMission(int(ack.StageId))
	if stage_t == nil {
		log.Error("stage nil", zap.Uint32("stage id", ack.StageId))
		return
	}

	switch msg.BattleType(stage_t.Type) {
	case msg.BattleType_Battle_EquipStage:
		OnFsCreateEquipStageFight(uid, ack)
	// case msg.BattleType_Battle_Peak:
	// 	PFCreateFightAck(uid, ack)
	default:
		p := player.FindByUserId(uid)
		if p != nil {
			// LeaveFight(p)
			GlobalAttrChange(p, false)
			p.SetFsId(ack.GetFightSerId())
			p.UserData.Fight.FightStageId = int(ack.GetStageId())
			p.UserData.Fight.FightId = ack.GetFightId()
			p.IsSendEndFight = true
			p.SaveFight()

			fightExtra := p.MakeProtocolBase()
			extraBytes, err := fightExtra.Marshal()
			if err != nil {
				log.Error("zombie extra marshal err", zap.Error(err),
					zap.Uint64("accountId", p.GetUserId()), zap.Any("extra", fightExtra))
				return
			}

			bc := &BattleCache{
				FsId:     p.GetFsId(),
				BattleId: ack.FightId,
				Uid:      uid,
				DeadLine: time.Now().Add(time.Minute * time.Duration(kenum.FightDeadLine)),
				StageId:  ack.StageId,
			}

			AddBattleCache(bc)

			SendFsEnterFight(p, ack.GetStageId(), ack.GetFightId(), extraBytes)
		}
	}
}

/*
 * SendFsEnterFight
 *  @Description: 向战斗服发送进入场景 (服务器内部通信)
 *  @param p
 *  @param stageId
 *  @param fightId
 *  @param extraBytes
 */
func SendFsEnterFight(p *player.Player, stageId, battleId uint32, extraBytes []byte) {
	bc := GetBattleCacheByUid(p.GetUserId())
	if bc == nil {
		log.Error("fs enter fight bc is nil")
		return
	}

	var pbPet = GetFightPets(p)
	// if pet := p.GetPet(p.UserData.BaseInfo.UsePet); pet != nil {
	// 	pbPet = &msg.FsPet{
	// 		Pet:  builder.BuildPetUnit(pet),
	// 		Attr: builder.BuildAttrMap(pet.BaseAttr), // todo 职业属性?
	// 	}
	// }
	ships := CalcShipsAttr(p)
	user := builder.BuildFsUser(p, ships)
	if user == nil {
		log.Error("build fs user err", zap.Uint64("userID", p.GetUserId()))
		return
	}
	enterFightReq := &msg.FsEnterFightReq{
		BattleId: bc.BattleId,
		User:     user,
		Pet:      pbPet,
		Extra:    extraBytes,
	}
	SendToFight(p, 0, enterFightReq)
}

func update_zero_time_activate(now time.Time) {
	if now.Hour() == 0 && now.Minute() == 0 && now.Second() == 0 {
		HandleRefreshRankData()
	}
}

func SendToFight(p *player.Player, pid uint32, m proto.Message) {
	bc := GetBattleCacheByUid(p.GetUserId())

	if bc == nil {
		log.Error("bc is nil")
		return
	}

	if bc.FsId == 0 {
		log.Error("fsId is 0", zap.String("user", p.Info()))
		return
	}
	pbBytes, err := proto.Marshal(m)
	if err != nil {
		log.Error("marshal err", zap.Error(err), zap.Uint64("uid", p.GetUserId()))
		return
	}

	fight_msg := &msg.GameToFight{
		UserId:   p.GetUserId(),
		PacketId: pid,
		MsgId:    msg.PbProcessor.GetMsgId(m),
		BattleId: bc.BattleId,
		GameId:   config.Conf.ServerId,
		GateId:   p.GetGateId(),
		Data:     pbBytes,
	}

	if p.UserData.PeakFight.RoomId > 0 {
		room := GetPFReadyRoom(p.UserData.PeakFight.RoomId)
		if room != nil {
			// room.users
		}
	}

	sub := "game.fight." + strconv.Itoa(int(bc.FsId))
	//log.Debug("Fight <<<<============= Game", zap.Uint32("msg id", pid), zap.String("info", p.Info()), zap.Any("data", m))
	//nats.Publish(sub, fight_msg)
	io_out.Push(&io_out.OutMsg{
		Subject: sub,
		Msg:     fight_msg,
	})
}
