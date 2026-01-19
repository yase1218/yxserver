package service

import (
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"math"
	"msg"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

type unlockIdx int

var (
	LockUnlockIdxCfgId1 unlockIdx = 38001
	LockUnlockIdxCfgId2 unlockIdx = 38002
	unLockMap           map[int]unlockIdx
)

func init() {
	unLockMap = make(map[int]unlockIdx)
	unLockMap[1] = LockUnlockIdxCfgId1
	unLockMap[2] = LockUnlockIdxCfgId1
}

func GetPlayerPets(p *player.Player) (int32, []*msg.FightInfo, []*msg.ActPet) {
	petsInfo := p.UserData.PetData
	lv := petsInfo.Lv
	actPetList := make([]*msg.ActPet, 0)
	var list = make([]*msg.FightInfo, 0)

	for _, pet := range petsInfo.Pets {
		data := &msg.ActPet{
			PetId:  int32(pet.BaseId),
			StarLv: int32(pet.StarLv),
		}
		actPetList = append(actPetList, data)

		if pet.Loc != model.Rest {
			fightInfo := &msg.FightInfo{}
			fightInfo.Id = pet.BaseId
			fightInfo.LocIdx = uint32(pet.LocIdx)
			list = append(list, fightInfo)
		}
	}

	return lv, list, actPetList
}

// 激活宠物
func ActTargetPet(p *player.Player, petId uint32) (msg.ErrCode, []uint32) {
	petsInfo := p.UserData.PetData

	if _, ok := petsInfo.Pets[petId]; ok {
		return msg.ErrCode_SUCC, nil
	}

	cfg := template.GetPetTemplate().GetPet(petId)
	if cfg == nil {
		return msg.ErrCode_CONFIG_NIL, nil
	}

	cost := cfg.ActivateItem
	res := CostItem(p.GetUserId(), uint32(cost[0]), uint32(cost[1]), publicconst.PetActivate, true)
	if res != msg.ErrCode_SUCC {
		return res, nil
	}

	activatePet := model.NewPet(petId)

	petsInfo.Pets[petId] = activatePet
	p.UserData.PetData = petsInfo
	p.SavePetData()

	actPetsList := actPetsToProtocol(petsInfo.Pets)
	UpdateTask(p, true, publicconst.TASK_COND_GET_PET, cfg.Rarity, 1) // 激活XX个XX品质库鲁兽
	return msg.ErrCode_SUCC, actPetsList
}

func ChangePetState(p *player.Player, pets []*msg.FightInfo) (msg.ErrCode, []*msg.FightInfo) {
	petsInfo := p.UserData.PetData.Pets

	if len(pets) <= 0 {
		for _, v := range petsInfo {
			if v.Loc != model.Rest {
				v.LocIdx = -1
				v.Loc = model.Rest
			}
		}
	} else {
		for i := range pets {
			if pets[i].Id == 0 {
				break
			}

			if pets[i].LocIdx > 0 {
				locIdx := pets[i].LocIdx
				if _, ok := unLockMap[int(locIdx)]; ok {
					cfg := template.GetFunctionTemplate().GetFunction(uint32(unLockMap[int(locIdx)]))
					if cfg == nil {
						return msg.ErrCode_CONFIG_NIL, nil
					}

					isUnlock := CanUnlock(p, cfg.Conditions)
					if !isUnlock {
						return msg.ErrCode_Pet_Loc_Is_Unlock, nil
					}
				} else {
					return msg.ErrCode_Pet_Loc_Is_Not_Exist, nil
				}

			}

			if _, ok := petsInfo[uint32(pets[i].Id)]; !ok {
				return msg.ErrCode_Pet_Is_Not_Exist, nil
			}

			for _, v := range petsInfo {
				if v.BaseId != uint32(pets[i].Id) && v.Loc == model.Fights && v.LocIdx == int16(i) {
					v.Loc = model.Rest
					v.LocIdx = -1
					break
				}
			}

			pet := petsInfo[uint32(pets[i].Id)]
			if pets[i].LocIdx == 0 {
				pet.Loc = model.Fights
			} else {
				pet.Loc = model.HelpWar
			}

			pet.LocIdx = int16(pets[i].LocIdx)
			petsInfo[uint32(pets[i].Id)] = pet
		}
	}

	p.UserData.PetData.Pets = petsInfo
	p.SavePetData()
	return msg.ErrCode_SUCC, pets

}

func UpdatePetStarLv(p *player.Player, petId uint32) (msg.ErrCode, uint32) {
	petsInfo := p.UserData.PetData
	pet, ok := petsInfo.Pets[petId]
	if !ok {
		return msg.ErrCode_Pet_Is_Not_Exist, 0
	}

	petRankCfgMap := template.GetPetRankTemplate().GetPet(petId)
	currentLv := pet.StarLv
	userId := p.GetUserId()

	totalCost := make(map[uint32]uint32, 5)
	maxLv := currentLv

	for nextLv := currentLv + 1; ; nextLv++ {
		rankCfg := petRankCfgMap[uint32(nextLv)]
		if rankCfg == nil {
			break
		}

		for i := 0; i < len(rankCfg.Items); i++ {
			itemId := uint32(rankCfg.Items[i][0])
			itemCount := uint32(rankCfg.Items[i][1])
			totalCost[itemId] += itemCount
		}

		enough := true
		for itemId, required := range totalCost {
			if GetItemNum(userId, itemId) < uint64(required) {
				enough = false
				for i := range rankCfg.Items {
					totalCost[uint32(rankCfg.Items[i][0])] -= uint32(rankCfg.Items[i][1])
				}
				break
			}
		}
		if !enough {
			break
		}

		maxLv = nextLv
	}

	if maxLv == currentLv {
		return msg.ErrCode_SUCC, uint32(currentLv)
	}

	for itemId, count := range totalCost {
		if code := CostItem(userId, itemId, count, publicconst.PetLvUp, true); code != msg.ErrCode_SUCC {
			return code, 0
		}
	}

	petsInfo.Pets[petId].StarLv = maxLv
	p.UserData.PetData = petsInfo
	p.SavePetData()

	return msg.ErrCode_SUCC, uint32(maxLv)
}

func UpdatePetLv(p *player.Player) (msg.ErrCode, uint32) {
	nextLv := p.UserData.PetData.Lv + 1
	cfg := template.GetPetLevelTemplate().GetPetNextLvCfg(nextLv)
	if cfg == nil {
		return msg.ErrCode_CONFIG_NIL, 0
	}

	items := cfg.Items
	userId := p.GetUserId()
	for v := range items {
		res := CostItem(userId, uint32(items[v][0]), uint32(items[v][1]), publicconst.PetLvUp, true)
		if res != msg.ErrCode_SUCC {
			return res, 0
		}
	}

	p.UserData.PetData.Lv = nextLv
	p.SavePetData()
	//tempLv := p.UserData.PetData.Lv
	UpdateTask(p, true, publicconst.TASK_COND_PET_LEVEL, uint32(nextLv)) // 库鲁兽等级达到XX级
	return msg.ErrCode_SUCC, uint32(nextLv)
}

func calcPetsPower(p *player.Player) int32 {
	petsInfo := p.UserData.PetData
	var tempAttrs map[uint32]*model.Attr
	lvLimit := p.UserData.PetData.Lv - 1
	for _, v := range petsInfo.Pets {
		petsCfg := template.GetPetTemplate().GetPet(v.BaseId)
		isAddBase := false
		for i := 0; i < len(petsCfg.BaseAttrs); i++ {
			if i == int(lvLimit) {
				InitAttr(tempAttrs, petsCfg.BaseAttrs[i])
				isAddBase = true
			}

			if isAddBase {
				break
			}
		}

		isAddAct := false
		for i := 0; i < len(petsCfg.ActivateAttrs); i++ {
			if i == int(lvLimit) {
				InitAttr(tempAttrs, petsCfg.ActivateAttrs[i])
				isAddAct = true
			}

			if isAddAct {
				break
			}
		}
	}

	var totolPower float32 = 0
	for id, item := range tempAttrs {
		if idConfig := template.GetAttrListTemplate().GetAttr(id); idConfig != nil {
			if item.Add > idConfig.CombatInit {
				totolPower += (item.Add - idConfig.CombatInit) * idConfig.CombatFactor
			}
		}
	}

	var resPower int32 = 0
	res := math.Floor(float64(totolPower))
	resPower = int32(res)

	return resPower
}

func actPetsToProtocol(pets map[uint32]*model.Pet) []uint32 {
	actPets := make([]uint32, 0)
	for _, v := range pets {
		actPets = append(actPets, v.BaseId)
	}

	return actPets
}

func AddPetByItem(p *player.Player, petId uint32, ntf bool) {
	petsInfo := p.UserData.PetData.Pets

	if v, ok := petsInfo[petId]; ok {
		cfg := template.GetPetTemplate().GetPet(v.BaseId)
		if cfg == nil {
			log.Error("pet config is nil", zap.String("accountId", p.UserData.AccountId), zap.Uint32("petId", petId))
			return
		}

		repeateItem := cfg.ItemRepeat
		AddItem(p.GetUserId(), uint32(repeateItem[0]), repeateItem[1], publicconst.PetActivate, ntf)
	} else {
		pet := model.NewPet(petId)
		petsInfo[petId] = pet
		p.UserData.PetData.Pets = petsInfo
		p.SavePetData()

		lv, list, actPetList := GetPlayerPets(p)

		msg := &msg.GetPlayerPetsInfoResp{}
		msg.ActPets = actPetList
		msg.Lv = lv
		msg.FightList = list
		p.SendNotify(msg)
	}

}

func GetFightPets(p *player.Player) *msg.FsPet {
	pets := p.UserData.PetData.Pets

	if len(pets) <= 0 {
		return nil
	}

	var fightPet = &msg.FsPet{}
	var helpwardSkilss []uint32 = make([]uint32, 0)
	for _, pet := range pets {
		if pet.Loc == model.Fights {
			calcAttr := calcPetAttr(p, pet.BaseId, uint32(pet.StarLv))
			skills := calcPetSkills(pet.BaseId, uint32(pet.StarLv), true)
			protoAttr := mergetAttrToProtocolAttr(calcAttr)
			fightPet.Pet = &msg.PetUnit{
				BaseId:      pet.BaseId,
				CareerAttrs: protoAttr,
				SkillIds:    skills,
			}
		}

		if pet.Loc == model.HelpWar {
			skills := calcPetSkills(pet.BaseId, uint32(pet.StarLv), false)
			helpwardSkilss = append(helpwardSkilss, skills...)
		}
	}

	skillsMap := make(map[uint32]bool, 0)
	finalSkillsList := make([]uint32, 0)
	for i := range helpwardSkilss {
		skillId := helpwardSkilss[i]
		if _, ok := skillsMap[skillId]; !ok {
			skillsMap[skillId] = true
			finalSkillsList = append(finalSkillsList, skillId)
		}
	}
	if fightPet.Pet != nil && len(finalSkillsList) > 0 {
		fightPet.Pet.SkillIds = append(fightPet.Pet.SkillIds, finalSkillsList...)
	}

	log.Debug("fight pet data", zap.Any("pet", fightPet))
	return fightPet
}

func calcPetAttr(p *player.Player, targetPetId uint32, starLv uint32) map[uint32]float64 {
	petData := p.UserData.PetData
	lv := petData.Lv

	// 升级属性
	lvCfg := template.GetPetLevelTemplate().GetPetLvCfg()
	lvAttrMap := make(map[uint32]float64)
	for _, cfg := range lvCfg {
		if lv == cfg.RoleLevel && len(cfg.LevelAttr) != 0 {
			attrId := cfg.LevelAttr[0]
			attrVal := cfg.LevelAttr[1]
			if _, ok := lvAttrMap[uint32(attrId)]; !ok {
				lvAttrMap[uint32(attrId)] = float64(attrVal)
			} else {
				lvAttrMap[uint32(attrId)] += float64(attrVal)
			}
		}
	}

	petCfg := template.GetPetTemplate().GetPet(targetPetId)
	if petCfg == nil {
		log.Error("pet config is nil", zap.Uint32("petId", targetPetId))
		return nil
	}

	// 激活属性
	petActAttr := petCfg.ActivateAttrs
	actAttrMap := make(map[uint32]float64)
	for _, v := range petActAttr {
		for i := range v {
			attr := v[i]
			if _, ok := actAttrMap[attr.Id]; !ok {
				actAttrMap[attr.Id] = float64(attr.Value)
			} else {
				actAttrMap[attr.Id] += float64(attr.Value)
			}
		}
	}

	// 基础属性
	basetAttr := make(map[uint32]float64, 0)
	baseAttrCfg := petCfg.BaseAttrs
	for _, v := range baseAttrCfg {
		for i := range v {
			attr := v[i]
			if _, ok := basetAttr[attr.Id]; !ok {
				basetAttr[attr.Id] = float64(attr.Value)
			} else {
				basetAttr[attr.Id] += float64(attr.Value)
			}
		}
	}

	// 升星属性
	petStarCfg := template.GetPetRankTemplate().GetPet(targetPetId)
	if petStarCfg == nil {
		log.Error("pet rank config is nil", zap.Uint32("petId", targetPetId))
		return nil
	}

	starAttrMap := make(map[uint32]float64)
	for lv, cfg := range petStarCfg {
		if lv <= starLv {
			for i := range cfg.RankAttr {
				attrs := cfg.RankAttr[i]
				for j := range attrs {
					attr := attrs[j]
					if _, ok := starAttrMap[attr.Id]; !ok {
						starAttrMap[attr.Id] = float64(attr.Value)
					} else {
						starAttrMap[attr.Id] += float64(attr.Value)
					}
				}
			}
		}
	}

	// 合并
	finalAttrMap := make(map[uint32]float64, 0)
	mergeAttrMap(lvAttrMap, finalAttrMap)
	mergeAttrMap(actAttrMap, finalAttrMap)
	mergeAttrMap(starAttrMap, finalAttrMap)
	mergeAttrMap(basetAttr, finalAttrMap)
	log.Debug("final attr maps", zap.Any("attrs", finalAttrMap))
	return finalAttrMap
}

func calcPetSkills(targetPetId uint32, starLv uint32, isFight bool) []uint32 {
	petCfg := template.GetPetTemplate().GetPet(targetPetId)
	if petCfg == nil {
		log.Error("pet config is nil", zap.Uint32("petId", targetPetId))
		return nil
	}

	baseSkills := petCfg.FightSkills

	petStarCfg := template.GetPetRankTemplate().GetPet(targetPetId)
	if petStarCfg == nil {
		log.Error("pet rank config is nil", zap.Uint32("petId", targetPetId))
		return nil
	}

	var starSkills []uint32 = make([]uint32, 0)
	for lv, cfg := range petStarCfg {
		if lv <= starLv {
			for i := range cfg.PickSkill {
				starSkills = append(starSkills, uint32(cfg.PickSkill[i]))
			}

			if !isFight {
				for i := range cfg.HelpWarSkill {
					starSkills = append(starSkills, uint32(cfg.HelpWarSkill[i]))
				}
			}
		}
	}

	finalSkillsList := make([]uint32, 0)
	if isFight {
		finalSkillsList = append(finalSkillsList, baseSkills...)
	}

	finalSkillsList = append(finalSkillsList, starSkills...)
	return finalSkillsList
}

func mergeAttrMap(source map[uint32]float64, target map[uint32]float64) map[uint32]float64 {
	for id, val := range source {
		if _, ok := target[id]; !ok {
			target[id] = val
		} else {
			target[id] += val
		}
	}

	return target
}

func mergetAttrToProtocolAttr(attrMap map[uint32]float64) []*msg.Attr {
	attrList := make([]*msg.Attr, 0)
	for attrId, attrVal := range attrMap {
		attr := &msg.Attr{
			Id:        attrId,
			CalcValue: float32(attrVal),
		}

		attrList = append(attrList, attr)
	}

	return attrList
}
