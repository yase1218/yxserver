package service

import (
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"kernel/tools"
	"msg"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

// UpdateMissionPoker 更新关卡扑克
func UpdateMissionPoker(p *player.Player, missionId int) {
	if p.UserData.Poker == nil {
		log.Error("poker nil", zap.Uint64("accountId", p.GetUserId()))
		return
	}
	if missionId == 0 {
		return
	}

	oldCount := len(p.UserData.Poker.MissData)
	missionConfig := template.GetMissionTemplate().GetMission(missionId)
	if missionConfig == nil {
		log.Error("mission cfg nil", zap.Int("missionId", missionId))
		return
	}
	for i := 0; i < len(missionConfig.UnlockPoker); i++ {
		if !tools.ListIntContain(p.UserData.Poker.MissData, missionConfig.UnlockPoker[i]) {
			p.UserData.Poker.MissData = append(p.UserData.Poker.MissData, missionConfig.UnlockPoker[i])
		}
	}

	// 有变化
	if oldCount != len(p.UserData.Poker.MissData) {
		p.SavePoker()
	}

	SendPokerNtf(p)
}

func SendPokerNtf(p *player.Player) {
	res := &msg.NotifyMissionPoker{
		Data: make([]uint32, 0, len(p.UserData.Poker.MissData)),
	}
	for _, v := range p.UserData.Poker.MissData {
		res.Data = append(res.Data, uint32(v))
	}
	p.SendNotify(res)
}

// UpdateShipPoker 更新机甲扑克
func UpdateShipPoker(p *player.Player, shipId uint32, sendClient bool) {
	if p.UserData.Poker == nil {
		log.Error("poker nil", zap.Uint64("accountId", p.GetUserId()))
	}

	data := getShip(p, shipId)
	if data == nil {
		log.Error("ship not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("shipId", shipId))
		return
	}

	shipConfig := template.GetRoleStarTemplate().GetShipStarLevel(data.Id, data.StarLevel)
	if shipConfig.Data.Poker == 0 {
		return
	}

	pos := -1
	for i := 0; i < len(p.UserData.Poker.ShipData); i++ {
		if p.UserData.Poker.ShipData[i].ShipId == data.Id {
			pos = i
			break
		}
	}

	var shipPoker *model.ShipPoker
	if pos == -1 {
		shipPoker = model.NewShipPoker(data.Id)
		updateShipPoker(shipPoker, shipConfig)
		p.UserData.Poker.ShipData = append(p.UserData.Poker.ShipData, shipPoker)
	} else {
		shipPoker = p.UserData.Poker.ShipData[pos]
		updateShipPoker(shipPoker, shipConfig)
	}

	p.SavePoker()

	if sendClient {
		p.SendNotify(&msg.NotifyShipPoker{
			Data: ToProtocolShipPoker(shipPoker),
		})
	}
}

// updateShipPoker 更新机甲扑克
func updateShipPoker(shipPoker *model.ShipPoker, shipConfig *template.JRoleStarLevel) {
	if !tools.ListContain(shipPoker.Poker, shipConfig.Data.Poker) {
		shipPoker.Poker = append(shipPoker.Poker, shipConfig.Data.Poker)
	}
}

// UpdateWeaponPoker 更新武器扑克
func UpdateWeaponPoker(p *player.Player, weaponId uint32) {
	if p.UserData.Poker == nil {
		log.Error("poker nil", zap.Uint64("accountId", p.GetUserId()))
		return
	}

	data := getWeapon(p, weaponId)
	if data == nil {
		log.Error("weapon not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("weaponId", weaponId))
		return
	}

	pos := -1
	for i := 0; i < len(p.UserData.Poker.WeaponData); i++ {
		if p.UserData.Poker.WeaponData[i].WeaponId == data.Id {
			pos = i
			break
		}
	}

	weaponLevelConfig := template.GetWeaponLevelTemplate().GetWeaponLevel(data.Id, data.Level)
	var poker *model.WeaponPoker
	if pos == -1 {
		poker = model.NewWeaponPoker(data.Id, weaponLevelConfig.Poker)
		p.UserData.Poker.WeaponData = append(p.UserData.Poker.WeaponData, poker)
	} else {
		poker = p.UserData.Poker.WeaponData[pos]
		if !tools.ListContain(poker.Poker, weaponLevelConfig.Poker) {
			poker.Poker = append(poker.Poker, weaponLevelConfig.Poker)
		}
	}

	p.SavePoker()
	res := &msg.NotifyWeaponPoker{}
	res.Data = ToProtocolWeaponPoker(poker)
	p.SendNotify(res)
}

func ToProtocolShipPoker(data *model.ShipPoker) *msg.ShipPoker {
	return &msg.ShipPoker{
		ShipId: data.ShipId,
		Poker:  data.Poker,
	}
}

func ToProtocolShipPokers(data []*model.ShipPoker) []*msg.ShipPoker {
	var ret []*msg.ShipPoker
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolShipPoker(data[i]))
	}
	return ret
}

func ToProtocolWeaponPoker(data *model.WeaponPoker) *msg.WeaponPoker {
	return &msg.WeaponPoker{WeaponId: data.WeaponId, Poker: data.Poker}
}

func ToProtocolWeaponPokers(data []*model.WeaponPoker) []*msg.WeaponPoker {
	var ret []*msg.WeaponPoker
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolWeaponPoker(data[i]))
	}
	return ret
}
