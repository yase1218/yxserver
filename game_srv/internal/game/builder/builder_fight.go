package builder

import (
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"msg"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func BuildMapType(stageId uint32) msg.MapType {
	_, ok := msg.OutMap_name[int32(stageId)]
	if ok {
		return msg.MapType_Map_Type_Out
	}

	return msg.MapType_Map_Type_Fight
}

func BuildFsUser(p *player.Player, ships map[uint32]map[uint32]*model.Attr) *msg.FsUser {
	var (
		fightType       = p.FightType
		pbSkillProSlice = make([]uint32, 0)
		pbWeaponMap     = make(map[uint32]uint32, len(p.UserData.Weapon.Weapons))
		coatId          = p.UserData.Ships.GetShipCoatId(p.UserData.BaseInfo.ShipId) // 机甲皮肤

	)
	pbShips := make(map[uint32]*msg.FsShips)
	for id, shipAttr := range ships {
		pbShips[id] = &msg.FsShips{
			SkillPros: make([]uint32, 0),
			Attrs:     BuildFightAttrMap(shipAttr),
			CoatId:    int32(p.UserData.Ships.GetShipCoatId(id)),
		}
	}
	var playerShip *model.Ship
	for _, ship := range p.UserData.Ships.Ships {
		if ship.Id == p.UserData.BaseInfo.ShipId {
			playerShip = ship
			break
		}
	}
	if playerShip == nil {
		log.Error("player ship nil", zap.Uint64("accountId", p.GetUserId()),
			zap.Uint32("shipId", p.UserData.BaseInfo.ShipId))
		return nil
	}

	for _, v := range p.UserData.Weapon.Weapons {
		pbWeaponMap[v.Id] = v.Level
	}

	for _, v := range p.UserData.Treasure.ShipData {
		if v.ShipId == p.UserData.BaseInfo.ShipId {
			pbSkillProSlice = append(pbSkillProSlice, v.WarTreasure...)
		} else {
			pbShips[v.ShipId].SkillPros = append(pbShips[v.ShipId].SkillPros, v.WarTreasure...)
		}
	}

	for _, v := range p.UserData.Treasure.WeaponData {
		pbSkillProSlice = append(pbSkillProSlice, v.Treasure...)
	}

	pbSkillProSlice = append(pbSkillProSlice, p.UserData.Treasure.MissData...)
	pbSkillProSlice = append(pbSkillProSlice, p.UserData.Treasure.CommData...)

	pbPokers := make([]uint32, 0, len(p.UserData.Poker.MissData))
	for _, v := range p.UserData.Poker.MissData {
		pbPokers = append(pbPokers, uint32(v))
	}

	return &msg.FsUser{
		Player:    BuildPlayerUnit(p),
		SkillPro:  pbSkillProSlice,
		Combat:    p.UserData.BaseInfo.Combat,
		FightType: fightType,
		ShipLv:    playerShip.StarLevel,
		Pokers:    pbPokers,
		Attr:      BuildFightAttrMap(p.UserData.BaseInfo.Attrs),
		Weapons:   pbWeaponMap,
		Ships:     pbShips,
		CoatId:    int32(coatId),
	}
}

func get_player_method(p *player.Player, btType int) *model.PlayMethodData {
	for i := 0; i < len(p.UserData.PlayMethod.Data); i++ {
		if p.UserData.PlayMethod.Data[i].BtType == btType {
			return p.UserData.PlayMethod.Data[i]
		}
	}
	return nil
}

func get_fight_weapons(p *player.Player) []uint32 {
	res := make([]uint32, 0)
	playMethod := get_player_method(p, int(p.FightType))
	if playMethod == nil {
		return res
	}
	return playMethod.WeaponIds
}

func BuildPlayerUnit(p *player.Player) *msg.PlayerUnit {
	pbData := &msg.PlayerUnit{
		AccountId:   p.GetUserId(),
		Name:        p.UserData.Nick,
		ShipId:      p.UserData.BaseInfo.ShipId,
		SupportId:   p.UserData.BaseInfo.SupportId,
		ComboSkills: p.UserData.BaseInfo.ComboSkill,
		WeaponIds:   make(map[uint32]uint32, 0),
		SkillIds:    make([]uint32, 0),
		Fresh:       len(p.UserData.Mission.Missions) == 0,
		CoatId:      uint32(p.UserData.Ships.GetShipCoatId(p.UserData.BaseInfo.ShipId)),
	}
	weapons := get_fight_weapons(p)

	for _, v := range weapons {
		for _, vv := range p.UserData.Weapon.Weapons {
			if vv.Id == v {
				pbData.WeaponIds[v] = vv.Level
				break
			}
		}
	}
	pbData.SkillIds = append(pbData.SkillIds, p.GetExSkills()...)
	return pbData
}

// todo 重构宠物后看怎么改
func BuildPetUnit(pet *model.Pet) *msg.PetUnit {
	if pet == nil {
		return nil
	}

	ret := &msg.PetUnit{
		BaseId: pet.BaseId,
		// Id:         pet.Id,
		// StartTime:  pet.StartTime,
		// EndTime:    pet.EndTime,
		// PetDressUp: &msg.PetDressUpInfo{
		// SuitId: pet.PetDressUp.SuitId,
		// BodyId: pet.PetDressUp.BodyId,
		// Up:     pet.PetDressUp.Up,
		// Down:   pet.PetDressUp.Down,
		// Head:   pet.PetDressUp.Head,
		// Hand:   pet.PetDressUp.Hand,
		// Foot:   pet.PetDressUp.Foot,
		// },
		// Name:     pet.Name,
		// SkillIds: pet.Skill,
	}

	// petConfig := template.GetPetTemplate().GetPet(pet.Id)
	//for i := 0; i < len(petConfig.BaseAttrs); i++ {
	//	if data, ok := pet.BaseAttr[petConfig.BaseAttrs[i].Id]; ok {
	//		ret.Attr = append(ret.BaseAttrs, BuildAttr(data))
	//	}
	//}

	// for i := 0; i < len(petConfig.CareerAttrs); i++ {
	// 	if data, ok := pet.CareerAttr[petConfig.CareerAttrs[i].Id]; ok {
	// 		ret.CareerAttrs = append(ret.CareerAttrs, BuildAttr(data))
	// 	}
	// }

	return ret
}

func BuildFightAttrMap(data map[uint32]*model.Attr) map[uint32]float64 {
	retMap := make(map[uint32]float64, len(data))
	for _, v := range data {
		ret := BuildAttr(v)
		retMap[ret.Id] = float64(ret.CalcValue)
	}
	return retMap
}
