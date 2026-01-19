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

// UpdateMissionTreasure 更新关卡秘宝
func UpdateMissionTreasure(p *player.Player, missionId uint32) {
	if p.UserData.Treasure == nil {
		log.Error("treasure nil", zap.Uint64("accountid", p.GetUserId()))
		return
	}

	oldCount := len(p.UserData.Treasure.MissData)
	missionConfig := template.GetMissionTemplate().GetMission(int(missionId))
	for i := 0; i < len(missionConfig.UnlockSkillpro); i++ {
		if !tools.ListContain(p.UserData.Treasure.MissData, uint32(missionConfig.UnlockSkillpro[i])) {
			p.UserData.Treasure.MissData = append(p.UserData.Treasure.MissData, uint32(missionConfig.UnlockSkillpro[i]))
		}
	}

	// 有变化
	if oldCount != len(p.UserData.Treasure.MissData) {
		p.SaveTreasure()
		res := &msg.NotifyMissionRareTreasure{}
		res.Treasure = p.UserData.Treasure.MissData
		p.SendNotify(res)
	}
}

func GetShipTreasure(p *player.Player) *model.ShipTreasure {
	res := &model.ShipTreasure{
		WarTreasure:     make([]uint32, 0),
		SupportTreasure: make([]uint32, 0),
	}

	for _, v := range p.UserData.Treasure.ShipData {
		if v.ShipId == p.UserData.BaseInfo.ShipId {
			return v
		}
	}

	return res
}

// UpdateShipTreasure 更新机甲秘宝
func UpdateShipTreasure(p *player.Player, shipId uint32, sendClient bool) {
	if p.UserData.Treasure == nil {
		//LoadTreasure(p)
		log.Error("treasure nil", zap.Uint64("accountId", p.GetUserId()))
		return
	}

	data := getShip(p, shipId)
	if data == nil {
		log.Error("ship not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("shipId", shipId))
		return
	}

	shipConfig := template.GetRoleStarTemplate().GetShipStarLevel(data.Id, data.StarLevel)
	if len(shipConfig.Data.PickSkillpro) == 0 && len(shipConfig.Data.DefaultSkillpro) == 0 &&
		len(shipConfig.Data.BanSkillpro) == 0 {
		return
	}

	pos := -1
	for i := 0; i < len(p.UserData.Treasure.ShipData); i++ {
		if p.UserData.Treasure.ShipData[i].ShipId == data.Id {
			pos = i
			break
		}
	}

	var shipTreasure *model.ShipTreasure
	if pos == -1 {
		shipTreasure = model.NewShipTreasure(data.Id)
		updateShipTreasure(shipTreasure, shipConfig)
		p.UserData.Treasure.ShipData = append(p.UserData.Treasure.ShipData, shipTreasure)
		//dao.TreasureDao.AddShipTreasure(p.GetUserId(), shipTreasure)
	} else {
		shipTreasure = p.UserData.Treasure.ShipData[pos]
		updateShipTreasure(shipTreasure, shipConfig)
		//dao.TreasureDao.UpdateShipTreasure(p.GetUserId(), shipTreasure)
	}
	p.SaveTreasure()

	if sendClient {
		p.SendNotify(&msg.NotifyShipRareTreasure{
			Data: ToProtocolShipTreasure(shipTreasure),
		})
	}
}

func updateShipTreasure(treasure *model.ShipTreasure, shipConfig *template.JRoleStarLevel) {
	for _, v := range shipConfig.Data.PickSkillpro {
		if !tools.ListContain(treasure.WarTreasure, uint32(v)) {
			treasure.WarTreasure = append(treasure.WarTreasure, uint32(v))
		}
	}

	for _, v := range shipConfig.Data.DefaultSkillpro {
		if !tools.ListContain(treasure.SupportTreasure, uint32(v)) {
			treasure.SupportTreasure = append(treasure.SupportTreasure, uint32(v))
		}
	}

	for _, v := range shipConfig.Data.BanSkillpro {
		treasure.WarTreasure = tools.ListRemove(treasure.WarTreasure, uint32(v))
		treasure.SupportTreasure = tools.ListRemove(treasure.WarTreasure, uint32(v))
	}
}

func removeBanSkill(treasure *model.WeaponTreasure, weaponConfig *template.JWeaponLevel) {
	for i := 0; i < len(weaponConfig.BanSkillpro); i++ {
		treasure.Treasure = tools.ListRemove(treasure.Treasure, weaponConfig.BanSkillpro[i])
	}
}

// UpdateWeaponTreasure 更新武器秘宝
func UpdateWeaponTreasure(p *player.Player, weaponId uint32) {
	if p.UserData.Treasure == nil {
		log.Error("Treasure nil", zap.Uint64("accountId", p.GetUserId()))
		return
	}

	data := getWeapon(p, weaponId)
	if data == nil {
		return
	}

	pos := -1
	for i := 0; i < len(p.UserData.Treasure.WeaponData); i++ {
		if p.UserData.Treasure.WeaponData[i].WeaponId == data.Id {
			pos = i
			break
		}
	}

	weaponConfig := template.GetWeaponLevelTemplate().GetWeaponLevel(data.Id, data.Level)
	if len(weaponConfig.PickSkillpro) == 0 {
		return
	}
	var treasure *model.WeaponTreasure
	if pos == -1 {
		treasure = model.NewWeaponTreasure(data.Id, weaponConfig.PickSkillpro)
		removeBanSkill(treasure, weaponConfig)
		p.UserData.Treasure.WeaponData = append(p.UserData.Treasure.WeaponData, treasure)
		p.SaveTreasure()
	} else {
		treasure = p.UserData.Treasure.WeaponData[pos]
		for i := 0; i < len(weaponConfig.PickSkillpro); i++ {
			if !tools.ListContain(treasure.Treasure, weaponConfig.PickSkillpro[i]) {
				treasure.Treasure = append(treasure.Treasure, weaponConfig.PickSkillpro[i])
			}
		}
		removeBanSkill(treasure, weaponConfig)
		p.SaveTreasure()
	}

	res := &msg.NotifyWeaponRareTreasure{}
	res.Data = ToProtoclWeaponTreasure(treasure)
	p.SendNotify(res)
}

func ToProtocolShipTreasure(data *model.ShipTreasure) *msg.ShipRareTreasure {
	return &msg.ShipRareTreasure{
		ShipId:          data.ShipId,
		WarTreasure:     data.WarTreasure,
		SupportTreasure: data.SupportTreasure,
	}
}

func ToProtocolShipTreasures(data []*model.ShipTreasure) []*msg.ShipRareTreasure {
	var ret []*msg.ShipRareTreasure
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolShipTreasure(data[i]))
	}
	return ret
}

func ToProtoclWeaponTreasure(data *model.WeaponTreasure) *msg.WeaponRareTreasure {
	return &msg.WeaponRareTreasure{WeaponId: data.WeaponId, Treasure: data.Treasure}
}

func ToProtocolWeaponTreasures(data []*model.WeaponTreasure) []*msg.WeaponRareTreasure {
	var ret []*msg.WeaponRareTreasure
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtoclWeaponTreasure(data[i]))
	}
	return ret
}
