package handle

import (
	"gameserver/internal/game/builder"
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

/*
 * CreateFight
 *  @Description: 统一创建战斗
 *  @param packetId
 *  @param args
 *  @param p
 */
func CreateFight(packetId uint32, args interface{}, p *player.Player) {
	log.Debug("CreateFight msg", zap.Uint64("uid", p.GetUserId()), zap.Uint32("packetId", packetId))
	req := args.(*msg.CreateFightReq)
	res := &msg.CreateFightAck{
		Result:  msg.ErrCode_SUCC,
		StageId: req.GetStageId(),
		MapType: builder.BuildMapType(req.GetStageId()),
	}
	if err := service.CreateFight(p, req.StageId, packetId, nil); err != msg.ErrCode_SUCC {
		log.Error("create fight err", zap.String("err", err.String()), zap.Uint32("stageId", req.GetStageId()))
		res.Result = err
	}
	p.SendResponse(packetId, res, res.Result)
}

/*
 * StartFight
 *  @Description: 开始战斗 因为有开场动画,由客户端控制
 *  @param packetId
 *  @param args
 *  @param p
 */
func StartFight(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.StartFightReq)
	res := &msg.StartFightAck{
		Result: msg.ErrCode_SUCC,
	}
	if err := service.StartFight(p, req, packetId, nil); err != nil {
		log.Error("start fight err", zap.Error(err), zap.Uint64("userID", p.GetUserId()))
		res.Result = msg.ErrCode_SYSTEM_ERROR
	}
	p.SendResponse(packetId, res, res.Result)
}

/*
 * EndFight
 *  @Description: 客户端主动结束战斗
 *  @param packetId
 *  @param args
 *  @param p
 */
func EndFight(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.EndFightReq)
	res := &msg.EndFightAck{Result: msg.ErrCode_SUCC}
	if err := service.EndFight(p, req, packetId, nil); err != nil {
		log.Error("end fight err", zap.Error(err), zap.Uint64("userID", p.GetUserId()))
		res.Result = msg.ErrCode_SYSTEM_ERROR
	}
	p.SendResponse(packetId, res, res.Result)
}

/*
废弃 协议先保留
*/
func GetWeapon(packetId uint32, args interface{}, p *player.Player) {
	//req := args.(*msg.GetWeaponReq)
	res := &msg.GetWeaponAck{Result: msg.ErrCode_SUCC, WeaponIds: p.UserData.Fight.Weapons}
	p.SendResponse(packetId, res, res.Result)
}

/*
废弃 协议先保留
*/
func SelectWeapon(packetId uint32, args interface{}, p *player.Player) {
	var (
		//req = args.(*msg.SelectWeaponReq)
		res = &msg.SelectWeaponAck{Result: msg.ErrCode_SUCC}
	)

	// p.UserData.Fight.Weapons = req.GetWeaponIds()
	// p.SaveFight()
	p.SendResponse(packetId, res, res.Result)
}

/*
 * GetAccessories
 *  @Description: 获取流派
 *  @param packetId
 *  @param args
 *  @param p
 */
func GetAccessories(packetId uint32, args interface{}, p *player.Player) {
	//req := args.(*msg.GetAccessoryReq)
	res := &msg.GetAccessoryAck{Result: msg.ErrCode_SUCC, Faction: p.UserData.Fight.Faction}
	p.SendResponse(packetId, res, res.Result)
}

/*
 * SelectAccessories
 *  @Description: 选择流派
 *  @param packetId
 *  @param args
 *  @param p
 */
func SelectAccessories(packetId uint32, args interface{}, p *player.Player) {
	var (
		req = args.(*msg.SelectAccessoryReq)
		res = &msg.SelectAccessoryAck{Result: msg.ErrCode_SUCC}
	)

	p.UserData.Fight.Faction = req.GetFaction()
	p.SaveFight()
	p.SendResponse(packetId, res, res.Result)
}

/*
 * FightBroadcast
 *  @Description: 战斗广播
 *  @param packetId
 *  @param args
 *  @param p
 */
func FightBroadcast(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.FightBroadcastReq)
	res := new(msg.FightBroadcastAck)
	if err := service.Broadcast(p, req, packetId, nil); err != nil {
		log.Error("broadcast peak fight err", zap.Error(err))
	}
	p.SendResponse(packetId, res, msg.ErrCode_SUCC)
}

/*
 * FightMonster
 *  @Description: 战斗中有哪些怪物 用来激活图鉴
 *  @param packetId
 *  @param args
 *  @param p
 */
func FightMonster(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.FightMonsterReq)
	res := new(msg.FightMonsterAck)
	if err := service.FightMonster(p, req, packetId, nil); err != nil {
		log.Error("fight monster err", zap.Error(err))
	}
	p.SendResponse(packetId, res, msg.ErrCode_SUCC)
}

/*
 * FsCreateFightAck
 *  @Description: 战斗服创建战斗回调
 *  @param packetId
 *  @param args
 *  @param p
 */
func FsCreateFightAck(packetId uint32, m *msg.FightToGame) {
	_, _, ack_pb, er := msg.PbProcessor.Unmarshal(m.Data)
	if er != nil {
		log.Error("fs create fight ack unmarshal err", zap.Error(er))
		return
	}
	ack := ack_pb.(*msg.FsCreateFightAck)

	if _, ok := msg.OutMap_name[int32(ack.GetStageId())]; ok {
		p := player.FindByUserId(m.UserId)
		if p != nil {
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
			service.SendFsEnterFight(p, ack.GetStageId(), ack.GetFightId(), extraBytes)
		}

	} else {
		stage_t := template.GetMissionTemplate().GetMission(int(ack.GetStageId()))
		if stage_t == nil {
			log.Error("stage nil", zap.Uint32("stage id", ack.GetStageId()))
			return
		}

		switch msg.BattleType(stage_t.Type) {
		case msg.BattleType_Battle_Peak:
			fallthrough
		case msg.BattleType_Battle_Zombie:
			fallthrough
		case msg.BattleType_Battle_Challenge:
			fallthrough
		case msg.BattleType_Battle_Desert:
			fallthrough
		case msg.BattleType_Battle_Union:
			fallthrough
		case msg.BattleType_Battle_EquipStage:
			fallthrough
		case msg.BattleType_Battle_Contract:
			fallthrough
		case msg.BattleType_Battle_Secret:
			fallthrough
		case msg.BattleType_Battle_Main:
			service.OnFsCreateFight(m.UserId, ack)
		default:
			log.Error("unknown fight type", zap.Uint64("accountId", m.UserId),
				zap.String("fightType", msg.BattleType(stage_t.Type).String()))
		}
	}

}

/*
 * FsFightResultNtf
 *  @Description: 战斗服战斗结果通知
 *  @param packetId
 *  @param args
 *  @param p
 */
func FsFightResultNtf(packetId uint32, m *msg.FightToGame) {
	//var (
	//	ntf = args.(*msg.FsFightResultNtf)
	//)
	// ntf := new(msg.FsFightResultNtf)
	// if err := proto.Unmarshal(args.([]byte), ntf); err != nil {
	// 	log.Error("fight result ntf unmarshal err", zap.Error(err))
	// 	return
	// }

	_, _, ntf_pb, er := msg.PbProcessor.Unmarshal(m.Data)
	if er != nil {
		log.Error("fs create fight ack unmarshal err", zap.Error(er))
		return
	}
	ntf := ntf_pb.(*msg.FsFightResultNtf)

	log.Info("fight result ntf", zap.Any("ntf", ntf))

	// out map no reward
	_, ok := msg.OutMap_name[int32(ntf.StageId)]
	if !ok {
		service.RemoveBattleCacheByUid(m.UserId)

		p := player.FindByUserId(m.UserId) // p可能为nil 已下线
		if p != nil {
			if !ntf.GetVictory() && ntf.GetReason() != msg.EndFightType_End_Fight_None {
				p.IsSendEndFight = true
			} else {
				p.IsSendEndFight = false
			}
		}

		stageCfg := template.GetMissionTemplate().GetMission(int(ntf.StageId))
		if stageCfg == nil {
			log.Error("stage cfg not exist", zap.Uint32("stageId", ntf.StageId))
			return
		}

		switch msg.BattleType(stageCfg.Type) {
		case msg.BattleType_Battle_Main:
			if p != nil {
				service.MainFightAfter(p, ntf)
			}
		case msg.BattleType_Battle_Challenge:
			if p != nil {
				service.ChallengeFightAfter(p, ntf)
			}
		case msg.BattleType_Battle_Desert, msg.BattleType_Battle_Union:
			if p != nil {
				service.PlayerMethodFightAfter(p, ntf)
			}
		case msg.BattleType_Battle_Peak:
			// todo 巅峰战场 结算
			if p != nil {
				if !ntf.GetVictory() {
					if err := service.PFUserExit(p, ntf.GetReason()); err != nil {
						log.Error("peak fight user exit err", zap.Error(err))
						return
					}
				}
			}
		case msg.BattleType_Battle_Zombie:
			// todo 竞技场 结算
			//service.ServMgr.GetArenaService().FightResult(playerData, ntf.GetVictory(), ntf.KillCnt, ntf.FightTime)
		case msg.BattleType_Battle_EquipStage:
			service.OnFsEquipStageFightResult(m.UserId, ntf)
			// todo 装备本结算
		case msg.BattleType_Battle_Contract:
			if p != nil {
				service.ContractFightAfter(p, ntf)
			}
		case msg.BattleType_Battle_Secret:
			if p != nil {
				service.SecretFightAfter(p, ntf)
			}
		}
		if p != nil {
			service.LeaveFight(p)
		}
	}
}

/*
 * FsPickItemNtf
 *  @Description: 战斗服拾取物品通知
 *  @param packetId
 *  @param args
 *  @param p
 */
func FsPickItemNtf(packetId uint32, m *msg.FightToGame) {
	_, _, ntf_pb, er := msg.PbProcessor.Unmarshal(m.Data)
	if er != nil {
		log.Error("fs create fight ack unmarshal err", zap.Error(er))
		return
	}
	ntf := ntf_pb.(*msg.FsPickItemNtf)

	p := player.FindByUserId(m.UserId)
	if p == nil {
		return
	}

	fightStageId := p.UserData.Fight.FightStageId
	if _, ok := msg.OutMap_name[int32(fightStageId)]; !ok {
		stageCfg := template.GetMissionTemplate().GetMission(fightStageId)
		if stageCfg == nil {
			log.Error("stage cfg nil", zap.Int("stageId", fightStageId))
		} else {
			switch msg.BattleType(stageCfg.Type) {
			case msg.BattleType_Battle_Peak:
				if err := service.PFPickItem(p, ntf.GetItems()); err != nil {
					log.Error("pick item err", zap.Error(err), zap.Uint64("uid", m.UserId), zap.Any("ntf", ntf))
				}
			}
		}
	}
}

func BattlePassWeekFieldsHandle(packetId uint32, args interface{}, p *player.Player) {
	log.Debug("BattlePassWeekFields msg", zap.Uint64("uid", p.GetUserId()), zap.Uint32("packetId", packetId))
	res := &msg.ResponseBattlePassWeekFields{
		Result: msg.ErrCode_SUCC,
	}
	res.Fields, res.OutTime = service.GetBattlePassWeekFields()
	p.SendResponse(packetId, res, res.Result)
}
