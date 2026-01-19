package service

import (
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tda"
	"msg"

	"github.com/zy/game_data/template"
)

// CalcFinalValue 计算终值
func CalcFinalValue(target map[uint32]*model.Attr) {
	// 先处理不需要公式计算的数值
	for id, data := range target {
		if config := template.GetAttrListTemplate().GetAttr(id); config != nil {
			// 固定数值
			if config.Total != 1 {
				data.CalcFinalValue()
			}
		}
	}

	for id, _ := range target {
		fieldConfig := template.GetAttrListTemplate().GetAttr(id)
		data, ok := target[id]
		if !ok {
			data = model.NewAttr(id, 0)
			target[id] = data
		}

		// 公式计算属性
		if fieldConfig.CalcType == 1 {
			a0 := getFieldCalcValue(fieldConfig.Paras[0], target)
			a1 := getFieldCalcValue(fieldConfig.Paras[1], target)
			a2 := getFieldCalcValue(fieldConfig.Paras[2], target)
			a3 := getFieldCalcValue(fieldConfig.Paras[3], target)
			a4 := getFieldCalcValue(fieldConfig.Paras[4], target)
			a5 := getFieldCalcValue(fieldConfig.Paras[5], target)
			a6 := getFieldCalcValue(fieldConfig.Paras[6], target)
			a7 := getFieldCalcValue(fieldConfig.Paras[7], target)
			a8 := getFieldCalcValue(fieldConfig.Paras[8], target)
			data.SetFinalValue((a0 + a1 + a2) * (1 + a3 + a4 + a5) * (1 + a6 + a7 + a8))
		} else if fieldConfig.CalcType == 2 {
			a0 := getFieldCalcValue(fieldConfig.Paras[0], target)
			a1 := getFieldCalcValue(fieldConfig.Paras[1], target)
			a2 := getFieldCalcValue(fieldConfig.Paras[2], target)
			data.SetFinalValue(a0 + a1 + a2)
		} else if fieldConfig.CalcType == 6 {
			a0 := getFieldCalcValue(fieldConfig.Paras[0], target)
			a1 := getFieldCalcValue(fieldConfig.Paras[1], target)
			a2 := getFieldCalcValue(fieldConfig.Paras[2], target)
			a3 := getFieldCalcValue(fieldConfig.Paras[3], target)
			a4 := getFieldCalcValue(fieldConfig.Paras[4], target)
			a5 := getFieldCalcValue(fieldConfig.Paras[5], target)
			data.SetFinalValue((a0 + a1 + a2) * (1 + a3 + a4 + a5))
		}
	}
}

func getFieldCalcValue(id uint32, target map[uint32]*model.Attr) float32 {
	config := template.GetAttrListTemplate().GetAttr(id)
	if config == nil {
		return 0
	}

	if data, ok := target[config.Id]; ok {
		return data.CalcFinalValue()
	}
	return 0
}

func CalcShipsAttr(p *player.Player) map[uint32]map[uint32]*model.Attr {
	res := make(map[uint32]map[uint32]*model.Attr)

	for _, ship := range p.UserData.Ships.Ships {
		if ship.Id == p.UserData.BaseInfo.ShipId {
			continue
		}
		res[ship.Id] = CalcShipFightAttr(p, ship)
	}
	return res
}

func CalcShipFightAttr(p *player.Player, ship *model.Ship) map[uint32]*model.Attr {
	// todo 计算出战属性
	// 支援库鲁的属性计算=(库鲁基础属性+当前装备属性+装备宝石属性)
	// *(1+局外升星加成)
	// *(1+武器库等级加成)
	// *(1+装备词条效果)
	// *(1+库鲁等级加成)

	temp := make(map[uint32]*model.Attr)

	for k, v := range ship.Attrs {
		temp[k] = model.DeepCopyAttr(v)
	}

	// 装备
	if p.UserData.Equip != nil {
		for i := 0; i < len(p.UserData.Equip.EquipPosData); i++ {
			mergeAttr(temp, p.UserData.Equip.EquipPosData[i].Attr)
			// 词条属性
			mergeAttr(temp, p.UserData.Equip.EquipPosData[i].AffixAttr)
		}
		mergeAttr(temp, GetGemAttrs(p))
	}

	// 武器库
	if p.UserData.Weapon != nil {
		mergeAttr(temp, p.UserData.Weapon.Attrs)
		// 武器
		for i := 0; i < len(p.UserData.Weapon.Weapons); i++ {
			mergeAttr(temp, p.UserData.Weapon.Weapons[i].Attrs)
		}
	}

	CalcFinalValue(temp)
	return temp
}

// GlobalAttrChange 计算终值
func GlobalAttrChange(p *player.Player, notifyClient bool) {
	// 出战  支援 装备 武器 武器库
	p.UserData.BaseInfo.Attrs = make(map[uint32]*model.Attr)

	addShipStarAttr(p)
	addSuitAttr(p)
	if p.UserData.Appearance != nil {
		mergeAttr(p.UserData.BaseInfo.Attrs, p.UserData.Appearance.Attrs)
	}

	// 出战
	if ship := getShip(p, p.UserData.BaseInfo.ShipId); ship != nil {
		mergeAttr(p.UserData.BaseInfo.Attrs, ship.Attrs)
	}

	// 支援
	for i := 0; i < len(p.UserData.BaseInfo.SupportId); i++ {
		if ship := getShip(p, p.UserData.BaseInfo.SupportId[i]); ship != nil {
			mergeAttr(p.UserData.BaseInfo.Attrs, ListToMapAttr(getShipSupoortAttr(ship)))
		}
	}

	// 装备
	if p.UserData.Equip != nil {
		for i := 0; i < len(p.UserData.Equip.EquipPosData); i++ {
			mergeAttr(p.UserData.BaseInfo.Attrs, p.UserData.Equip.EquipPosData[i].Attr)
			// 词条属性
			mergeAttr(p.UserData.BaseInfo.Attrs, p.UserData.Equip.EquipPosData[i].AffixAttr)
		}
		mergeAttr(p.UserData.BaseInfo.Attrs, GetGemAttrs(p))
	}

	// 武器库
	if p.UserData.Weapon != nil {
		mergeAttr(p.UserData.BaseInfo.Attrs, p.UserData.Weapon.Attrs)
		// 武器
		for i := 0; i < len(p.UserData.Weapon.Weapons); i++ {
			mergeAttr(p.UserData.BaseInfo.Attrs, p.UserData.Weapon.Weapons[i].Attrs)
		}
	}

	// 宠物
	if p.UserData.PetData != nil {
		// TODO 宠物暂留
		// for _, v := range p.UserData.PetData.Pets {
		// mergeAttrAll(p.UserData.BaseInfo.Attrs, v.BaseAttr)
		// mergeAttrAll(p.UserData.BaseInfo.Attrs, v.CareerAttr)
		// }
	}

	// 默认加暴击伤害150
	_, ok := p.UserData.BaseInfo.Attrs[uint32(msg.Attribute_Attr_Crit_Damage_Rate)]
	if !ok {
		p.UserData.BaseInfo.Attrs[uint32(msg.Attribute_Attr_Crit_Damage_Rate)] = model.NewAttr(uint32(msg.Attribute_Attr_Crit_Damage_Rate), 150)
	} else {
		p.UserData.BaseInfo.Attrs[uint32(msg.Attribute_Attr_Crit_Damage_Rate)].InitValue += 150
	}

	CalcFinalValue(p.UserData.BaseInfo.Attrs)
	//dao.AccountDao.UpdateGlobalAttrs(p.GetAccountId(), p.UserData.BaseInfo.Attrs)
	p.SaveBaseInfo()

	calcCombat(p, true)

	if notifyClient {
		notifyMsg := &msg.NotifyCombatChange{}
		notifyMsg.Combat = p.UserData.BaseInfo.Combat
		notifyMsg.Attrs = ToProtocolGlobalAttr(p.UserData.BaseInfo.Attrs)
		p.SendNotify(notifyMsg)
	}
}

func addShipStarAttr(p *player.Player) {
	for i := 0; i < len(p.UserData.Ships.Ships); i++ {
		attrs := template.GetRoleStarTemplate().GetShipRangeStarAttr(p.UserData.Ships.Ships[i].Id,
			-1, int32(p.UserData.Ships.Ships[i].StarLevel))
		for id, data := range attrs {
			if _, ok := p.UserData.BaseInfo.Attrs[id]; ok {
				p.UserData.BaseInfo.Attrs[id].InitValue += data.Value
			} else {
				attr := model.NewAttr(id, data.Value)
				p.UserData.BaseInfo.Attrs[id] = attr
			}
		}
	}
}

func addSuitAttr(p *player.Player) {
	// 套装属性
	var ids []uint32
	if p.UserData.Equip != nil {
		for i := 0; i < len(p.UserData.Equip.EquipPosData); i++ {
			ids = append(ids, p.UserData.Equip.EquipPosData[i].EquipId)
		}
	}
	suitAttr := template.GetEquipTemplate().GetSuitAttr(ids)
	for i := 0; i < len(suitAttr); i++ {
		id := suitAttr[i].Id
		value := suitAttr[i].Value
		if _, ok := p.UserData.BaseInfo.Attrs[id]; ok {
			p.UserData.BaseInfo.Attrs[id].InitValue += value
		} else {
			attr := model.NewAttr(id, value)
			p.UserData.BaseInfo.Attrs[id] = attr
		}
	}
}

func ListToMapAttr(lst []*model.Attr) map[uint32]*model.Attr {
	ret := make(map[uint32]*model.Attr)
	for i := 0; i < len(lst); i++ {
		ret[lst[i].Id] = lst[i]
	}
	return ret
}

// calcCombat 计算战力
func calcCombat(p *player.Player, notifyClient bool) {
	totalAttr := make(map[uint32]*model.Attr)
	// 全局属性
	mergeCombatMap(totalAttr, p.UserData.BaseInfo.Attrs)

	var total float32 = 0
	for id, item := range totalAttr {
		if idConfig := template.GetAttrListTemplate().GetAttr(id); idConfig != nil {
			if item.Add > idConfig.CombatInit {
				total += (item.Add - idConfig.CombatInit) * idConfig.CombatFactor
			}
		}
	}

	if p.UserData.BaseInfo.Combat < uint32(total) {
		// tda update power
		tdaData := &tda.CommonUser{
			Highest_power: p.UserData.BaseInfo.Combat,
		}
		tda.TdaUpdateCommonUser(p.TdaCommonAttr.AccountId, p.TdaCommonAttr.DistinctId, tdaData)
	}

	if p.UserData.BaseInfo.Combat != uint32(total) {
		p.UserData.BaseInfo.Combat = uint32(total)
		//dao.AccountDao.UpdateCombat(playerData.GetAccountId(), playerData.AccountInfo.Combat)
		p.SaveBaseInfo()
	}
}

func mergeCombatMap(targetAttr map[uint32]*model.Attr, source map[uint32]*model.Attr) {
	for id, item := range source {
		if d, ok := targetAttr[id]; ok {
			d.Add += item.GetRawValue()
		} else {
			d := model.NewAttr(id, 0)
			d.Add = item.GetRawValue()
			targetAttr[id] = d
		}
	}
}

func mergeAttr(targetAttr map[uint32]*model.Attr, source map[uint32]*model.Attr) {
	for id, item := range source {
		idConfig := template.GetAttrListTemplate().GetAttr(id)
		var attrId uint32
		switch idConfig.EffectType {
		case template.Attr_Effect_Type_All:
			attrId = id - 200
		default:
			attrId = id
		}
		if d, ok := targetAttr[attrId]; ok {
			if idConfig.ValueType == 2 {
				d.Add += item.GetRawValue()
			} else {
				d.Add += item.FinalValue
			}
		} else {
			d := model.NewAttr(attrId, 0)
			if idConfig.ValueType == 2 {
				d.Add = item.GetRawValue()
			} else {
				d.Add = item.FinalValue
			}
			targetAttr[attrId] = d
		}
	}
}

func mergeAttrAll(targetAttr map[uint32]*model.Attr, source map[uint32]*model.Attr) {
	for id, item := range source {
		if id < 200 {
			continue
		}

		idConfig := template.GetAttrListTemplate().GetAttr(id)
		var attrId uint32
		switch idConfig.EffectType {
		case template.Attr_Effect_Type_All:
			attrId = id - 200
		default:
			attrId = id
		}
		if d, ok := targetAttr[attrId]; ok {
			if idConfig.ValueType == 2 {
				d.Add += item.GetRawValue()
			} else {
				d.Add += item.FinalValue
			}
		} else {
			d := model.NewAttr(attrId, 0)
			if idConfig.ValueType == 2 {
				d.Add = item.GetRawValue()
			} else {
				d.Add = item.FinalValue
			}
			targetAttr[attrId] = d
		}
	}
}

func getAttrsFromMap(attrs map[uint32]*model.Attr, attrIds []uint32) []*model.Attr {
	var ret []*model.Attr
	for m := 0; m < len(attrIds); m++ {
		if data, ok := attrs[attrIds[m]]; ok {
			ret = append(ret, data)
		}
	}
	return ret
}

// InitAttr 初始化属性
func InitAttr(data map[uint32]*model.Attr, initAttr []*template.AttrItem) []uint32 {
	var ret []uint32
	for i := 0; i < len(initAttr); i++ {
		if d, ok := data[initAttr[i].Id]; ok {
			d.AddInitValue(initAttr[i].Value)
		} else {
			data[initAttr[i].Id] = model.NewAttr(initAttr[i].Id, initAttr[i].Value)
		}
		ret = append(ret, initAttr[i].Id)
	}
	return ret
}

//// AddAttrLevelValue 增加等级属性
//func AddAttrLevelValue(data map[uint32]*model.Attr, levelAttr []*template.AttrItem) []uint32 {
//	var ret []uint32
//	for i := 0; i < len(levelAttr); i++ {
//		if d, ok := data[levelAttr[i].Id]; ok {
//			d.AddLevelValue(levelAttr[i].Value)
//		} else {
//			d := model.NewAttr(levelAttr[i].Id, 0)
//			d.SetLevelValue(levelAttr[i].Value)
//			data[levelAttr[i].Id] = d
//		}
//		ret = append(ret, levelAttr[i].Id)
//	}
//	return ret
//}

// SetAttrLevelValue 增加等级属性
func SetAttrLevelValue(data map[uint32]*model.Attr, levelAttr []*template.AttrItem) []uint32 {
	var ret []uint32
	for i := 0; i < len(levelAttr); i++ {
		if d, ok := data[levelAttr[i].Id]; ok {
			d.SetLevelValue(levelAttr[i].Value)
		} else {
			d := model.NewAttr(levelAttr[i].Id, 0)
			d.SetLevelValue(levelAttr[i].Value)
			data[levelAttr[i].Id] = d
		}
		ret = append(ret, levelAttr[i].Id)
	}

	return ret
}

func ToProtocolGlobalAttr(attrs map[uint32]*model.Attr) []*msg.Attr {
	var ret []*msg.Attr
	for id, attr := range attrs {
		if id == uint32(publicconst.Attack) ||
			id == uint32(publicconst.Hp) ||
			id == uint32(publicconst.Defense) {
			ret = append(ret, ToProtocolAttr(attr))
		}
	}
	return ret
}

func ToProtocolAttr(data *model.Attr) *msg.Attr {
	return &msg.Attr{
		Id:        data.Id,
		Value:     data.InitValue + data.LevelValue + data.Add,
		CalcValue: data.FinalValue,
	}
}

func ToProtocolAttrs(data []*model.Attr) []*msg.Attr {
	var ret []*msg.Attr
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolAttr(data[i]))
	}
	return ret
}

func ToProtocolAttrs2(data map[uint32]*model.Attr) []*msg.Attr {
	var ret []*msg.Attr
	for _, d := range data {
		ret = append(ret, ToProtocolAttr(d))
	}
	return ret
}
