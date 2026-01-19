package player

import (
	"errors"
	"fmt"
	"gameserver/internal/config"
	"gameserver/internal/io_out"
	"gameserver/internal/publicconst"
	"gameserver/internal/tapping"
	"kernel/protocol"
	"kernel/tda"
	"math/rand"
	"msg"
	"strconv"
	"time"

	"gameserver/internal/game/model"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// Player 玩家数据 池对象 注意重置脏数据
type Player struct {
	State         publicconst.PlayerState // 玩家状态
	UpdateTime    uint32                  // 更新时间
	ChannelId     uint32
	SdkChannelNo  string
	GateId        uint32
	SessionId     uint64          // gate上的session id
	TdaCommonAttr *tda.CommonAttr // 数数打点公共属性
	Os            int             // 1:android 2 ios 3 web 4 pc端

	UserData             *model.UserData     // 玩家数据
	MissionAdRewardItems []*model.SimpleItem // 关卡广告奖励 不入库
	CheckBattle          uint32

	// 聊天相关
	InWorldChannel    bool   // 是否在世界频道
	LastWorldChatTime uint32 // 上次世界聊天时间 做cd 用

	UpdateNickTime   uint32 // 修改昵称的时间
	ContractRandTime uint32 // 疯狂合约上次随机时间

	Rand           *rand.Rand       // 随机数
	FightType      msg.BattleType   // 正在战斗的类型
	LastFightTime  time.Time        // 上次开始战斗时间
	IsSendEndFight bool             // 是否发送退出战斗
	EndFightType   msg.EndFightType // 结束战斗类型
	LastTick       int64            // 上次update时间

	EquipStageItems map[uint32]uint32 // 物品id: 数量

	EquipUpgradeStageTime int64 // 装备升阶

	MailLoaded bool // 是否加载过邮件

	fsId uint32 // 当前fsid todo 临时记录 需要管理保存

	dirty map[string]struct{}
}

// 池对象 重置脏数据
func (p *Player) Init(now time.Time) {
	p.State = publicconst.Logining
	p.UpdateTime = uint32(now.Unix())
	p.ChannelId = 0
	p.GateId = 0
	p.SessionId = 0
	p.TdaCommonAttr = &tda.CommonAttr{}
	p.UserData = &model.UserData{}
	p.Rand = rand.New(rand.NewSource(time.Now().Unix()))
	p.EquipStageItems = make(map[uint32]uint32)
	p.EquipUpgradeStageTime = 0
	p.MailLoaded = false
	p.fsId = 0
	p.dirty = make(map[string]struct{})
}

func (p *Player) GetOpenId() string {
	if p.UserData == nil {
		return ""
	}
	return p.UserData.AccountId
}

// func (p *Player) GetAccountId() uint64 {
// 	if p.UserData == nil {
// 		return 0
// 	}
// 	return p.UserData.UserId
// }

func (p *Player) GetUserId() uint64 {
	if p.UserData == nil {
		return 0
	}
	return p.UserData.UserId
}
func (p *Player) GetNick() string {
	if p.UserData == nil {
		return ""
	}
	return p.UserData.Nick
}

func (p *Player) Info() string {
	if p.UserData == nil {
		return "nil player info!!!"
	}
	return fmt.Sprintf("a:%v i:%v n:%v", p.UserData.AccountId, p.UserData.UserId, p.UserData.Nick)
}

func (p *Player) BindGateId(gateId uint32) {
	p.GateId = gateId
}

func (p *Player) GetGateId() uint32 {
	return p.GateId
}

func (p *Player) BindSessionId(sessionId uint64) {
	p.SessionId = sessionId
}

func (p *Player) GetSessionId() uint64 {
	return p.SessionId
}

func (p *Player) ResetMissData() {
	if p.UserData.BaseInfo.MissData != nil {
		p.UserData.BaseInfo.MissData.MissionId = 0
		p.UserData.BaseInfo.MissData.StartTime = 0
		p.UserData.BaseInfo.MissData.Total = 0
		p.UserData.BaseInfo.MissData.Speed = 10
		p.UserData.BaseInfo.MissData.Ads = p.UserData.BaseInfo.MissData.Ads[0:0]
	}
}

//func (p *PlayerData) IsMsgToQuick() bool {
//	listBuckets := func() []uint32 {
//		var buckets []uint32
//		var count uint32 = 0
//		p.rw.Reduce(func(b *collection.Bucket[uint32]) {
//			count += b.Sum
//		})
//		buckets = append(buckets, count)
//		return buckets
//	}
//
//	return listBuckets()[0] > publicconst.MAX_MSG_COUN_MINUTE
//}
//
//func (p *PlayerData) AddMsgCount(count uint32) {
//	p.rw.Add(count)
//}

func (p *Player) IsOnline() bool {
	return p.State == publicconst.Online
}

func (p *Player) GetBattleId() uint32 {
	return p.UserData.Fight.FightId
}
func (p *Player) GetStageId() int {
	return p.UserData.Fight.FightStageId
}

// func (p *Player) GetAgent() gate.Agent {
// 	return p.PlayerAgent
// }

func (p *Player) GetPet(id uint32) *model.Pet {
	petConfig := template.GetPetTemplate().GetPet(id)
	if petConfig == nil {
		return nil
	}

	if p.UserData.PetData != nil && p.UserData.PetData.Pets != nil {
		for _, v := range p.UserData.PetData.Pets {
			if v.BaseId == id {
				return v
			}
		}
	}
	return nil
}

func (p *Player) MakeProtocolBase() *protocol.GameFightBaseExtra {
	return &protocol.GameFightBaseExtra{
		CommonAttr: p.TdaCommonAttr,
		ChannelId:  p.ChannelId,
		ServerId:   config.Conf.ServerId,
	}
}

// 额外技能
func (p *Player) GetExSkills() []uint32 {
	ret := make([]uint32, 0)
	for _, v := range p.UserData.Equip.EquipPosData {
		if len(v.AffixSkills) > 0 {
			ret = append(ret, v.AffixSkills...)
		}
	}
	return ret
}

func (p *Player) NtfPoker() {
	ntf := &msg.NotifyMissionPoker{
		Data: make([]uint32, 0, len(p.UserData.Poker.MissData)),
	}
	for _, v := range p.UserData.Poker.MissData {
		ntf.Data = append(ntf.Data, uint32(v))
	}
	p.SendNotify(ntf)
}

// func (p *Player) OnLogout(now time.Time) {
// 	log.Info("player logout", zap.Uint64("uid", p.GetUserId()))
// 	p.SaveAllSync(false)
// 	DelPlayer(p)
// }

func (p *Player) SendError(m proto.Message, err msg.ErrCode) {
	msg_id := msg.PbProcessor.GetMsgId(m)
	err_ntf := &msg.NotifyErrMsg{
		Id:     msg.MsgId(msg_id),
		Result: err,
	}
	p.SendNotify(err_ntf)
}

// 一般用于返回客户端请求 需要packet_id和错误码(客户端报错需要单独错误消息NotifyErrMsg)
func (p *Player) SendResponse(packet_id uint32, m proto.Message, ec msg.ErrCode) {
	traceFunc := func(mid, gid uint32, sid uint64, account, content string) {
		msgName := ""
		switch mid {
		case msg.MsgID_ResponseLoginId:
			msgName = "ResponseLogin"
		case msg.MsgID_RandomPlayerNameRespId:
			msgName = "RandomPlayerNameRes"
		case msg.MsgID_InitPlayerNameAndShipRspId:
			msgName = "InitPlayerNameAndShipRsp"
		default:
			return
		}
		log.Info("trace msg send to gate "+content,
			zap.Uint32("gate id", gid),
			zap.Uint32("msg", mid),
			zap.String("msgName", msgName),
			zap.Uint64("session id", sid),
			zap.String("account", account),
		)
	}

	msg_id := msg.PbProcessor.GetMsgId(m)
	gate_msg := &msg.GameToGate{
		GameId:    config.Conf.ServerId,
		Session:   p.GetSessionId(),
		AccountId: p.GetOpenId(),
		PacketId:  packet_id,
		MsgId:     msg_id,
	}
	sub := "game.gate." + strconv.Itoa(int(p.GetGateId()))
	if data, err := msg.PbProcessor.Marshal(gate_msg.PacketId, m); err == nil {
		gate_msg.Data = data

		if config.Conf.IsDebug() {
			//log.Debug("Gate <<<<============= Game", zap.Uint32("msg id", msg_id), zap.String("info", p.Info()), zap.Any("data", m))
		}
		//nats.Publish(sub, gate_msg)
		io_out.Push(&io_out.OutMsg{
			Subject: sub,
			Msg:     gate_msg,
		})
		traceFunc(msg_id, p.GetGateId(), p.GetSessionId(), p.GetOpenId(), "push to nats")
	} else {
		log.Error("marshal msg failed", zap.Error(err))
		traceFunc(msg_id, p.GetGateId(), p.GetSessionId(), p.GetOpenId(), "marshal failed")
	}

	if ec != msg.ErrCode_SUCC && ec != msg.ErrCode_ERR_NONE {
		err_msg := &msg.NotifyErrMsg{
			Id:     msg.MsgId(msg_id),
			Result: ec,
		}
		gate_err := &msg.GameToGate{
			GameId:    config.Conf.ServerId,
			Session:   p.GetSessionId(),
			AccountId: p.GetOpenId(),
			PacketId:  packet_id,
			MsgId:     uint32(msg.MsgId_ID_NotifyErrMsg),
		}
		if data, err := msg.PbProcessor.Marshal(gate_err.PacketId, err_msg); err == nil {
			gate_err.Data = data
			io_out.Push(&io_out.OutMsg{
				Subject: sub,
				Msg:     gate_err,
			})
		} else {
			log.Error("marshal err_msg failed", zap.Error(err))
		}
	}
}

// 一般用于主动推送 无packet_id和错误码
func (p *Player) SendNotify(m proto.Message) {
	msg_id := msg.PbProcessor.GetMsgId(m)
	gate_msg := &msg.GameToGate{
		GameId:    config.Conf.ServerId,
		Session:   p.GetSessionId(),
		AccountId: p.GetOpenId(),
		PacketId:  0,
		MsgId:     msg_id,
	}

	if data, err := msg.PbProcessor.Marshal(gate_msg.PacketId, m); err == nil {
		gate_msg.Data = data
		sub := "game.gate." + strconv.Itoa(int(p.GetGateId()))
		io_out.Push(&io_out.OutMsg{
			Subject: sub,
			Msg:     gate_msg,
		})
	} else {
		log.Error("marshal err", zap.Error(err))
	}
}

func (p *Player) SetFsId(fsId uint32) {
	p.fsId = fsId
}

func (p *Player) GetFsId() uint32 {
	return p.fsId
}

// 只有sdk打点
func (p *Player) TappingSdk(data map[string]any, eventName string) {
	if err := setCommonAttr(data, eventName); err != nil {
		log.Error("TappingSdk failed data not tda.ICommonAttr", zap.Error(err))
		return
	}

	if err := tapping.PushTapData(&tapping.TapData{
		AccountId: p.GetOpenId(),
		ChannelId: p.ChannelId,
		EventName: eventName,
		Switch:    tapping.TapSwitchShushu,
		Data:      data,
	}); err != nil {
		log.Error("TappingSdk tap data failed", zap.Uint64("uid", p.GetUserId()), zap.Error(err))
	}
}

// 只有本地打点
func (p *Player) TappingLocal(data map[string]any, eventName string) {
	if err := setCommonAttr(data, eventName); err != nil {
		log.Error("TappingLocal failed data not tda.ICommonAttr", zap.Error(err))
		return
	}

	if err := tapping.PushTapData(&tapping.TapData{
		AccountId: p.GetOpenId(),
		ChannelId: p.ChannelId,
		EventName: eventName,
		Switch:    tapping.TapSwitchLocal,
		Data:      data,
	}); err != nil {
		log.Error("TappingLocal tap data failed", zap.Uint64("uid", p.GetUserId()), zap.Error(err))
	}
}

// sdk和本地都打点
func (p *Player) TappingBoth(data any, eventName string) {
	if err := setCommonAttr(data, eventName); err != nil {
		log.Error("TappingBoth failed data not tda.ICommonAttr", zap.Error(err))
		return
	}

	if err := tapping.PushTapData(&tapping.TapData{
		AccountId: p.GetOpenId(),
		ChannelId: p.ChannelId,
		EventName: eventName,
		Switch:    tapping.TapSwitchAll,
		Data:      data,
	}); err != nil {
		log.Error("TappingBoth tap data failed", zap.Uint64("uid", p.GetUserId()), zap.Error(err))
	}
}

func setCommonAttr(data any, eventName string) error {
	if tapData, ok := data.(tda.ICommonAttr); !ok {
		err := errors.New(fmt.Sprintf("tap data not tda.ICommonAttr: %v", data))
		return err
	} else {
		tapData.SetEventName(eventName)
		tapData.SetEventTime(time.Now())
	}
	return nil
}
