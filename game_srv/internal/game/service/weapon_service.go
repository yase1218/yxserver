package service

import (
	"gameserver/internal/config"
	"gameserver/internal/game/event"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/report"
	"kernel/tools"
	"msg"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

// ActiveWeapon 激活武器
func ActiveWeapon(p *player.Player, weaponId uint32) msg.ErrCode {
	weapon := getWeapon(p, weaponId)
	if weapon != nil {
		return msg.ErrCode_WEAPON_HAS_ACTIVE
	}

	weaponConfig := template.GetWeaponTemplate().GetWeapon(weaponId)
	for i := 0; i < len(weaponConfig.CostItem); i++ {
		if !EnoughItem(p.GetUserId(), weaponConfig.CostItem[i].ItemId,
			weaponConfig.CostItem[i].ItemNum) {
			return msg.ErrCode_NO_ENOUGH_ITEM
		}
	}

	if len(weaponConfig.CostItem) == 0 {
		return msg.ErrCode_NO_ENOUGH_ITEM
	}

	// 扣除道具
	var notifyClientItems []uint32
	for i := 0; i < len(weaponConfig.CostItem); i++ {
		CostItem(p.GetUserId(),
			weaponConfig.CostItem[i].ItemId,
			weaponConfig.CostItem[i].ItemNum,
			publicconst.ActiveWeaponCostItem,
			false)
		notifyClientItems = append(notifyClientItems, weaponConfig.CostItem[i].ItemId)
	}

	// 通知客户端
	updateClientItemsChange(p.GetUserId(), notifyClientItems)

	AddWeapon(p, weaponId, publicconst.ActiveWeapon)
	return msg.ErrCode_SUCC
}

// SetSecondaryWeapon 设置副武器
func SetSecondaryWeapon(p *player.Player, pos, weaponId uint32) msg.ErrCode {
	weapon := getWeapon(p, weaponId)
	if weapon == nil {
		return msg.ErrCode_WEAPON_NOT_EXIST
	}

	var secondaryWeapon *model.SecondaryWeapon = nil
	for i := 0; i < len(p.UserData.Weapon.SecondaryWeapons); i++ {
		if p.UserData.Weapon.SecondaryWeapons[i].Pos == pos {
			secondaryWeapon = p.UserData.Weapon.SecondaryWeapons[i]
			break
		}
	}

	if secondaryWeapon == nil {
		if len(p.UserData.Weapon.SecondaryWeapons) >= int(template.GetSystemItemTemplate().SecondaryWeaponNum) {
			return msg.ErrCode_OVER_SECONDARY_NUM
		}

		secondaryWeapon = model.NewSecondaryWeapon(pos, weaponId)
		p.UserData.Weapon.SecondaryWeapons = append(p.UserData.Weapon.SecondaryWeapons, secondaryWeapon)
	} else {
		secondaryWeapon.WeaponId = weaponId
	}
	p.SaveWeapon()

	return msg.ErrCode_SUCC
}

// UpgradeWeapon 升级武器
func UpgradeWeapon(p *player.Player, weaponId uint32) (msg.ErrCode, uint32) {
	weapon := getWeapon(p, weaponId)
	var nextLevel uint32
	if weapon == nil {
		log.Error("weapon table nil", zap.Uint32("weaponId", weaponId), zap.Any("weapon", weapon))
		return msg.ErrCode_WEAPON_NOT_EXIST, 0
	}

	if weapon.Level >= p.UserData.Level {
		return msg.ErrCode_INVALID_DATA, 0
	}

	weaponLevel := template.GetWeaponLevelTemplate().GetWeaponLevel(weaponId, weapon.Level)
	if weaponLevel == nil {
		log.Error("weaponLevel table nil", zap.Uint32("weaponId", weaponId), zap.Uint32("lv", weapon.Level))
		return msg.ErrCode_WEAPON_NOT_EXIST, 0
	}
	nextLevel = weapon.Level + 1

	nextWeaponConfig := template.GetWeaponLevelTemplate().GetWeaponLevel(weaponId, nextLevel)
	if nextWeaponConfig == nil {
		return msg.ErrCode_WEAPON_LEVEL_FULL, 0
	}

	for i := 0; i < len(weaponLevel.CostItem); i++ {
		if !EnoughItem(p.GetUserId(),
			weaponLevel.CostItem[i].ItemId, weaponLevel.CostItem[i].ItemNum) {
			return msg.ErrCode_NO_ENOUGH_ITEM, 0
		}
	}

	// 扣除道具
	var notifyClientItems []uint32
	//tdaItems := make([]*tda.Item, 0, len(weaponLevel.CostItem))
	for i := 0; i < len(weaponLevel.CostItem); i++ {
		CostItem(p.GetUserId(),
			weaponLevel.CostItem[i].ItemId,
			weaponLevel.CostItem[i].ItemNum,
			publicconst.UpgradeWeaponCostItem,
			false)
		notifyClientItems = append(notifyClientItems, weaponLevel.CostItem[i].ItemId)
		//tdaItems = append(tdaItems, &tda.Item{ItemId: strconv.Itoa(int(weaponLevel.CostItem[i].ItemId)), ItemNum: weaponLevel.CostItem[i].ItemNum})
	}

	weapon.Level += 1
	weapon.UpdateTime = tools.GetCurTime()
	event.EventMgr.PublishEvent(event.NewWeaponUpgradeEvent(p, weaponId, weapon.Level, ListenWeaponUpgradeEvent))

	// 通知客户端
	updateClientItemsChange(p.GetUserId(), notifyClientItems)

	//if nextWeaponConfig != nil {
	p.UserData.Weapon.WeaponLibExp += nextWeaponConfig.GetWeaponLibExp
	upgradeWeaponLib(p)
	//}
	// 计算等级属性
	calcWeaponLevelAttr(p, weapon, nextWeaponConfig)
	notifyWeaponChange(p, []uint32{weaponId})

	p.SaveWeapon()

	AddComboSkill(p, nextWeaponConfig.PickSkillcombo)

	// tda
	// tda.TdaWeaponSystemUpgrade(p.ChannelId, p.TdaCommonAttr, weaponId, weapon.Level, tdaItems)

	// report.ReportWeaponLvup(p.ChannelId, p.GetUserId(), config.Conf.ServerId, weaponId, weapon.Level)
	return msg.ErrCode_SUCC, weapon.Level
}

func calcWeaponAttr(p *player.Player, weapon *model.Weapon) {
	for _, data := range weapon.Attrs {
		var finalValue float32
		finalValue = data.InitValue + data.LevelValue + data.Add
		data.SetFinalValue(finalValue)
	}
}

// calcWeaponLevelAttr 计算武器升级属性
func calcWeaponLevelAttr(p *player.Player, weapon *model.Weapon, config *template.JWeaponLevel) {
	if len(config.LevelAttr) == 0 {
		return
	}

	for i := 0; i < len(config.LevelAttr); i++ {
		if d, ok := weapon.Attrs[config.LevelAttr[i].Id]; ok {
			d.AddLevelValue(config.LevelAttr[i].Value)
		} else {
			d := model.NewAttr(config.LevelAttr[i].Id, 0)
			d.AddLevelValue(config.LevelAttr[i].Value)
			weapon.Attrs[config.LevelAttr[i].Id] = d
		}
	}

	// 计算终值
	calcWeaponAttr(p, weapon)
	p.SaveWeapon()

	// 计算全局属性
	GlobalAttrChange(p, true)
}

// GmUpgradeWeapon 通过gm命令设置武器等级
func GmUpgradeWeapon(p *player.Player, weaponId, lv uint32) msg.ErrCode {
	weapon := getWeapon(p, weaponId)
	if weapon == nil {
		log.Error("weapon table nil", zap.Uint32("weaponId", weaponId), zap.Any("weapon", weapon))
		return msg.ErrCode_WEAPON_NOT_EXIST
	}

	if weapon.Level >= p.UserData.Level {
		return msg.ErrCode_INVALID_DATA
	}

	weaponLevel := template.GetWeaponLevelTemplate().GetWeaponLevel(weaponId, weapon.Level)
	if weaponLevel == nil {
		log.Error("weaponLevel table nil", zap.Uint32("weaponId", weaponId), zap.Uint32("lv", weapon.Level))
		return msg.ErrCode_WEAPON_NOT_EXIST
	}

	nextWeaponConfig := template.GetWeaponLevelTemplate().GetWeaponLevel(weaponId, lv)
	if nextWeaponConfig == nil {
		return msg.ErrCode_WEAPON_LEVEL_FULL
	}

	weapon.Level = lv
	weapon.UpdateTime = tools.GetCurTime()
	event.EventMgr.PublishEvent(event.NewWeaponUpgradeEvent(p, weaponId, weapon.Level, ListenWeaponUpgradeEvent))

	p.UserData.Weapon.WeaponLibExp += nextWeaponConfig.GetWeaponLibExp
	upgradeWeaponLib(p)

	// 计算等级属性
	calcWeaponLevelAttr(p, weapon, nextWeaponConfig)
	notifyWeaponChange(p, []uint32{weaponId})
	p.SaveWeapon()

	AddComboSkill(p, nextWeaponConfig.PickSkillcombo)

	return msg.ErrCode_SUCC
}

// PassMissionAddWeapon 通关添加武器
func PassMissionAddWeapon(p *player.Player, missionId int) {
	missionConfig := template.GetMissionTemplate().GetMission(missionId)
	if missionConfig == nil {
		return
	}

	if len(missionConfig.WeaponID) == 0 {
		return
	}

	for i := 0; i < len(missionConfig.WeaponID); i++ {
		AddWeapon(p, uint32(missionConfig.WeaponID[i]), publicconst.PassMissionAddWeapon)
	}
}

// getWeapon 获取武器
func getWeapon(p *player.Player, weaponId uint32) *model.Weapon {
	for i := 0; i < len(p.UserData.Weapon.Weapons); i++ {
		if p.UserData.Weapon.Weapons[i].Id == weaponId {
			return p.UserData.Weapon.Weapons[i]
		}
	}
	return nil
}

// AddWeapon 添加武器
func AddWeapon(p *player.Player, weaponId uint32, source publicconst.ItemSource) (msg.ErrCode, *model.Weapon) {
	weaponConfig := template.GetWeaponTemplate().GetWeapon(weaponId)
	if weaponConfig == nil {
		return msg.ErrCode_INVALID_DATA, nil
	}

	if temp := getWeapon(p, weaponId); temp != nil {
		return msg.ErrCode_INVALID_DATA, nil
	}

	weapon := model.NewWeapon(weaponId, 1, weaponConfig.SkillId)
	p.UserData.Weapon.Weapons = append(p.UserData.Weapon.Weapons, weapon)

	weaponLevel := template.GetWeaponLevelTemplate().GetWeaponLevel(weaponId, 1)
	if weaponLevel == nil {
		log.Error("weaponLevel table nil", zap.Uint32("weaponId", weaponId))
		return msg.ErrCode_WEAPON_NOT_EXIST, nil
	}
	p.UserData.Weapon.WeaponLibExp += weaponLevel.GetWeaponLibExp
	upgradeWeaponLib(p)

	p.SaveWeapon()

	// 添加组合技能
	AddComboSkill(p, weaponLevel.PickSkillcombo)

	if source != publicconst.InitAddItem {
		notifyWeaponChange(p, []uint32{weaponId})
	}
	UpdateWeaponTreasure(p, weaponId)
	UpdateWeaponPoker(p, weaponId)
	UpdateTask(p, true, publicconst.TASK_COND_ANY_WEAPON_LEVEL)

	report.ReportWeaponAdd(p.ChannelId, p.GetUserId(), config.Conf.ServerId, weaponId, weapon.Level)
	return msg.ErrCode_SUCC, weapon
}

func getWeaponMaxLevel(p *player.Player) uint32 {
	var maxLevel uint32 = 0
	for i := 0; i < len(p.UserData.Weapon.Weapons); i++ {
		if p.UserData.Weapon.Weapons[i].Level > maxLevel {
			maxLevel = p.UserData.Weapon.Weapons[i].Level
		}
	}
	return maxLevel
}

// upgradeWeaponLib 升级武器库
func upgradeWeaponLib(p *player.Player) {
	level := template.GetWeaponLibTeamplate().GetLevelByExp(p.UserData.Weapon.WeaponLibExp)
	if level != nil && level.Level != p.UserData.Weapon.WeaponLibLevel {
		calcWeaponLibAttr(p, p.UserData.Weapon.WeaponLibLevel, level.Level)
		p.UserData.Weapon.WeaponLibLevel = level.Level
		event.EventMgr.PublishEvent(event.NewWeaponLibUpgradeEvent(p, p.UserData.Weapon.WeaponLibLevel, ListenWeaponLibUpgradeEvent))
	}
}

// calcWeaponWeaponAttr 计算武器库属性
func calcWeaponLibAttr(p *player.Player, start, end uint32) {
	attrs := template.GetWeaponLibTeamplate().GetLevelRangeAttr(start, end)
	if len(attrs) == 0 {
		return
	}

	for id, item := range attrs {
		if _, ok := p.UserData.Weapon.Attrs[id]; ok {
			p.UserData.Weapon.Attrs[id].AddLevelValue(item.Value)
		} else {
			d := model.NewAttr(id, 0)
			d.AddLevelValue(item.Value)
			p.UserData.Weapon.Attrs[id] = d
		}
	}

	for _, data := range p.UserData.Weapon.Attrs {
		var finalValue float32
		finalValue = data.InitValue + data.LevelValue + data.Add
		data.SetFinalValue(finalValue)
	}

	// 计算战力
	GlobalAttrChange(p, true)
}

// NotifyClientAddWeapon 通知客户端获得武器
func notifyWeaponChange(p *player.Player, weaponIds []uint32) {
	res := &msg.NotifyWeaponChange{}
	res.LibData = ToProtocolWeaponLib(p.UserData.Weapon)
	for k := 0; k < len(weaponIds); k++ {
		if weapon := getWeapon(p, weaponIds[k]); weapon != nil {
			res.Data = append(res.Data, ToProtocolWeapon(weapon))
		}
	}
	p.SendNotify(res)
}

func ToProtocolWeapon(data *model.Weapon) *msg.Weapon {
	ret := &msg.Weapon{
		Id:      data.Id,
		Level:   data.Level,
		SkillId: data.SkillId,
	}

	//for id, data := range data.Attrs {
	//	if id == uint32(publicconst.Attack) || id == uint32(publicconst.Hp) || id == uint32(publicconst.Defense) {
	//		ret.Attrs = append(ret.Attrs, &msg.Attr{
	//			Id:        data.Id,
	//			Value:     data.InitValue + data.LevelValue + data.Add,
	//			CalcValue: data.FinalValue,
	//		})
	//	}
	//}
	return ret
}

func ToProtocolWeaponLib(weapon *model.AccountWeapon) *msg.WeaponLib {
	ret := &msg.WeaponLib{}
	ret.Level = weapon.WeaponLibLevel
	ret.Exp = weapon.WeaponLibExp
	displayList := template.GetAttrListTemplate().GetDisplayAttr()
	for id, data := range weapon.Attrs {
		if id == 2 || id == 4 || id == 6 || id == 21 {
			for i := 0; i < len(displayList); i++ {
				if displayList[i].Id != id {
					continue
				}

				attrConfig := displayList[i]
				temp := &msg.Attr{
					Id:        attrConfig.Id,
					CalcValue: data.FinalValue,
					Value:     data.InitValue + data.LevelValue + data.Add,
				}

				if len(attrConfig.DisplayPara) > 0 {
					var finalValue float32
					for k := 0; k < len(attrConfig.DisplayPara); k++ {
						if data2, ok2 := weapon.Attrs[attrConfig.DisplayPara[k]]; ok2 {
							finalValue += data2.CalcFinalValue()
						}
					}
					temp.CalcValue = finalValue
				}
				ret.Attrs = append(ret.Attrs, temp)
				break
			}
		}
	}
	return ret
}

func ToProtocolWeapons(data []*model.Weapon) []*msg.Weapon {
	var ret []*msg.Weapon
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolWeapon(data[i]))
	}
	return ret
}

func ToProtocolSecondaryWeapon(data *model.SecondaryWeapon) *msg.WeaponPos {
	return &msg.WeaponPos{
		Pos: data.Pos,
		Id:  data.WeaponId,
	}
}

func ToProtocolSecondaryWeapons(data []*model.SecondaryWeapon) []*msg.WeaponPos {
	var ret []*msg.WeaponPos
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolSecondaryWeapon(data[i]))
	}
	return ret
}
