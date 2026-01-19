package handle

import (
	"gameserver/internal/game/player"
	"msg"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

type GateHandlerFunc func(uint32, interface{}, *player.Player)
type FightHandlerFunc func(uint32, *msg.FightToGame)

var gate_handler_map map[msg.MsgId]GateHandlerFunc
var fight_handler_map map[msg.MsgId]FightHandlerFunc

func get_gate_hanler(id uint32) GateHandlerFunc {
	return gate_handler_map[msg.MsgId(id)]
}

func get_fight_hanler(id uint32) FightHandlerFunc {
	return fight_handler_map[msg.MsgId(id)]
}

func reg_gate(id msg.MsgId, handler GateHandlerFunc) {
	gate_handler_map[msg.MsgId(id)] = handler
}

func reg_fight(id msg.MsgId, handler FightHandlerFunc) {
	fight_handler_map[msg.MsgId(id)] = handler
}

func init() {
	gate_handler_map = make(map[msg.MsgId]GateHandlerFunc)
	fight_handler_map = make(map[msg.MsgId]FightHandlerFunc)
}

func HandleGateMsg(msg_id, packet_id uint32, args interface{}, p *player.Player) {
	f := get_gate_hanler(msg_id)
	if f == nil {
		log.Error("gate msg handler not found", zap.Uint32("msgId", msg_id), zap.Uint32("packetId", packet_id))
		return
	} else {
		f(packet_id, args, p)
	}
}

func HandleFightMsg(msg_id uint32, m *msg.FightToGame) {
	f := get_fight_hanler(msg_id)
	if f == nil {
		log.Error("fight msg handler not found", zap.Uint32("msgId", msg_id))
		return
	} else {
		f(msg_id, m)
	}
}
