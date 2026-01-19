package service

import (
	common2 "gameserver/internal/common"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tools"
	"msg"
	"sort"
	"time"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

type EquipList struct {
	Equips []*model.Equip
	Num    uint32
}

func manageEquip(p *player.Player, used_equip map[uint32]struct{}) map[uint32]map[uint32]map[uint32]*EquipList {
	equip_manager := make(map[uint32]map[uint32]map[uint32]*EquipList)

	for _, v := range p.UserData.Equip.EquipData {
		j_equip := template.GetEquipTemplate().GetEquip(v.Id)
		if j_equip == nil {
			log.Error("equip cfg nil", zap.Uint32("equipId", v.Id))
			continue
		}

		if equip_manager[j_equip.Data.Pos] == nil {
			equip_manager[j_equip.Data.Pos] = make(map[uint32]map[uint32]*EquipList)
		}

		if equip_manager[j_equip.Data.Pos][j_equip.Data.Rarity] == nil {
			equip_manager[j_equip.Data.Pos][j_equip.Data.Rarity] = make(map[uint32]*EquipList)
		}

		if equip_manager[j_equip.Data.Pos][j_equip.Data.Rarity][j_equip.Data.SmallStage] == nil {
			equip_manager[j_equip.Data.Pos][j_equip.Data.Rarity][j_equip.Data.SmallStage] = &EquipList{
				Equips: make([]*model.Equip, 0),
				Num:    0,
			}
		}
		equip_manager[j_equip.Data.Pos][j_equip.Data.Rarity][j_equip.Data.SmallStage].Equips =
			append(equip_manager[j_equip.Data.Pos][j_equip.Data.Rarity][j_equip.Data.SmallStage].Equips, v)

		if used_equip == nil {
			equip_manager[j_equip.Data.Pos][j_equip.Data.Rarity][j_equip.Data.SmallStage].Num += v.Num
		} else {

			if _, ok := used_equip[v.Id]; ok { // 装备中
				if v.Num >= 1 {
					equip_manager[j_equip.Data.Pos][j_equip.Data.Rarity][j_equip.Data.SmallStage].Num += v.Num - 1
				}
			} else {
				equip_manager[j_equip.Data.Pos][j_equip.Data.Rarity][j_equip.Data.SmallStage].Num += v.Num
			}
		}
	}

	for _, pos := range equip_manager {
		for _, rarity := range pos {
			for _, smallStage := range rarity {
				sort.Slice(smallStage.Equips, func(i, j int) bool {
					return smallStage.Equips[i].Num < smallStage.Equips[j].Num
				})
			}
		}
	}
	return equip_manager
}

// updateEquip 更新装备
func updateEquip(p *player.Player, equip *model.Equip, save bool) {
	if equip.Num <= 0 {
		for i := 0; i < len(p.UserData.Equip.EquipData); i++ {
			if p.UserData.Equip.EquipData[i].Id == equip.Id {
				p.UserData.Equip.EquipData = append(p.UserData.Equip.EquipData[0:i], p.UserData.Equip.EquipData[i+1:]...)
				break
			}
		}
	}
	if save {
		p.SaveEquip()
	}
}

func oneKeyUpgradeEquipPos(p *player.Player, add_equip_map, reduce_map map[uint32]*model.Equip) []*model.EquipPos {
	update_pos := make([]*model.EquipPos, 0)
	equip_manager := manageEquip(p, nil)
	cur_time := tools.GetCurTime()

	for _, pos_data := range p.UserData.Equip.EquipPosData {
		DebugLog("oneKeyUpgradeEquipPos 1", ZapUser(p))
		j_e := template.GetEquipTemplate().GetEquip(pos_data.EquipId)
		if j_e == nil {
			continue
		}
		equip_data := getEquip(p, pos_data.EquipId)
		if equip_data == nil {
			continue
		}

		if _, ok := equip_manager[j_e.Data.Pos]; !ok {
			continue
		}

		if _, ok := equip_manager[j_e.Data.Pos][j_e.Data.Rarity]; !ok {
			continue
		}

		elist, ok := equip_manager[j_e.Data.Pos][j_e.Data.Rarity][j_e.Data.SmallStage]
		if !ok {
			continue
		}
		lst := elist.Equips

		totoal_num := elist.Num

		// 自身不能算入材料 总数需要减1
		// if pos_data.EquipId == j_e.Data.Id && totoal_num > 0 {
		// 	totoal_num -= 1
		// }

		need_num := template.GetEquipUpConditionTemplate().GetNeedNum(j_e.Data.Pos, j_e.Data.Rarity, j_e.Data.SmallStage)
		// enough
		if totoal_num < need_num {
			continue
		}

		costed_num := uint32(0) // 已经消耗
	LOOP:
		for i := uint32(0); i < need_num; i++ {
			DebugLog("oneKeyUpgradeEquipPos 3", ZapUser(p))
			for _, material := range lst {
				DebugLog("oneKeyUpgradeEquipPos 4", ZapUser(p))
				if material.Num <= 0 { // correct
					updateEquip(p, material, false)
					if _, ok := reduce_map[material.Id]; !ok {
						reduce_map[material.Id] = material
					}
					continue
				}

				// pos_self := uint32(0)
				// 如果是自己
				// if material.Id == pos_data.EquipId &&
				// 	material.Num == 1 {
				// 	continue
				// }

				// 轮训每次消耗1个
				if costed_num+1 >= need_num { // 够
					costed_num += 1
					material.Num -= 1
					updateEquip(p, material, false)
					if _, ok := reduce_map[material.Id]; !ok {
						reduce_map[material.Id] = material
					}
					break LOOP
				} else { // 不够
					costed_num += 1
					material.Num -= 1
					updateEquip(p, material, false)
					if _, ok := reduce_map[material.Id]; !ok {
						reduce_map[material.Id] = material
					}
				}
			}
		}

		// 删除原装备
		if equip_data.Num > 0 {
			updateEquip(p, equip_data, false)
			if _, ok := reduce_map[equip_data.Id]; !ok {
				reduce_map[equip_data.Id] = equip_data
			}
		} else {
			log.Error("pos equip num <= 0 before remove self",
				zap.Uint32("pos", pos_data.EquipId),
				zap.Uint64("player_id", p.GetUserId()),
			)
		}

		// 添加新装备并穿戴
		new_e := AddEquip(p, j_e.Data.Rank, 1, false)
		if _, ok := add_equip_map[new_e.Id]; !ok {
			add_equip_map[new_e.Id] = new_e
		}

		pos_data.EquipId = j_e.Data.Rank
		pos_data.UpdateTime = cur_time
		updateEquipPos(p, pos_data, j_e, true)
		update_pos = append(update_pos, pos_data)
	}

	return update_pos
}

func oneKeyUpgradeEquipBag(p *player.Player, add_equip_map, reduce_map map[uint32]*model.Equip) {
	used_equip := make(map[uint32]struct{})

	for _, pos := range p.UserData.Equip.EquipPosData {
		DebugLog("oneKeyUpgradeEquipBag 1", ZapUser(p))
		if _, ok := used_equip[pos.EquipId]; ok {
			continue
		}

		used_equip[pos.EquipId] = struct{}{}
	}
	equip_manager := manageEquip(p, used_equip)
	safe_cnt := 10000
	for pos, pos_map := range equip_manager {

		DebugLog("oneKeyUpgradeEquipBag 2", ZapUser(p))
		for rarity, rarity_map := range pos_map {
			DebugLog("oneKeyUpgradeEquipBag 3", ZapUser(p))
			for stage, elst := range rarity_map {

				loop_cnt := 0
				DebugLog("oneKeyUpgradeEquipBag 4", ZapUser(p))

				// 合成一个所需数量
				need_num := template.GetEquipUpConditionTemplate().GetNeedNum(pos, rarity, stage)

				if elst.Num < need_num {
					continue
				}

				total_cnt := uint32(elst.Num / need_num) // 可合成数 todo 有装备再减一
				com_cnt := uint32(0)                     // 当前合成数
				cost_cnt := uint32(0)                    // 当前合成消耗

				reduce_e := make(map[uint32]*model.Equip)

				for com_cnt < total_cnt { // 未完成
					loop_cnt++
					if loop_cnt > safe_cnt {
						log.Error("may be dead loop")
						break
					}
					DebugLog("oneKeyUpgradeEquipBag 5", ZapUser(p))
					for _, e := range elst.Equips {

						DebugLog("oneKeyUpgradeEquipBag 6", ZapUser(p))
						if e.Num == 0 {
							DebugLog("oneKeyUpgradeEquipBag 7", ZapUser(p))
							continue
						}

						// if _, ok := used_equip[e.Id]; ok {
						// 	if e.Num == 1 {
						// 		DebugLog("oneKeyUpgradeEquipBag 8", ZapUser(p))
						// 		continue
						// 	}
						// }

						e.Num--
						cost_cnt++
						if _, ok := reduce_e[e.Id]; !ok {
							reduce_e[e.Id] = e
						}
						// 合成一个
						if cost_cnt >= need_num {
							com_cnt++
							cost_cnt = 0
							break
						}
					}
				}

				for i := uint32(0); i < total_cnt; i++ {

					DebugLog("oneKeyUpgradeEquipBag 7", ZapUser(p))
					rand_j := template.GetEquipTemplate().RandEquip(p.Rand, pos, rarity, stage)
					if rand_j == nil {
						log.Error("RandEquip error", zap.Uint64("AccountId", p.GetUserId()), zap.Int("pos", int(pos)), zap.Int("rarity", int(rarity)), zap.Int("stage", int(stage)))
						continue
					}

					new_e := AddEquip(p, rand_j.Data.Id, 1, false)

					if _, ok := add_equip_map[new_e.Id]; !ok {
						add_equip_map[new_e.Id] = new_e
					}
				}

				for _, v := range reduce_e {

					DebugLog("oneKeyUpgradeEquipBag 7", ZapUser(p))
					updateEquip(p, v, false)
					if _, ok := reduce_map[v.Id]; !ok {
						reduce_map[v.Id] = v
					}
				}
			}
		}
	}
}

// getSuitPosMaxquality 获取套装最高品质装备
func getSuitPosMaxquality(p *player.Player, suitId, pos uint32) *template.JEquip {
	var ret *template.JEquip = nil
	for i := 0; i < len(p.UserData.Equip.EquipData); i++ {
		equip := p.UserData.Equip.EquipData[i]
		if configData := template.GetEquipTemplate().GetEquip(equip.Id); configData != nil {
			if configData.Data.OutfitsType == suitId && configData.Data.Pos == pos {
				if ret == nil {
					ret = configData
				} else {
					if configData.HasHighQuality(ret) {
						ret = configData
					}
				}
			}
		}
	}
	return ret
}

func getSuit(p *player.Player, suitId uint32) *model.SuitInfo {
	for i := 0; i < len(p.UserData.Equip.SuitReward); i++ {
		if p.UserData.Equip.SuitReward[i].SuitId == suitId {
			return p.UserData.Equip.SuitReward[i]
		}
	}
	return nil
}

func addSuit(p *player.Player, data *model.SuitInfo) {
	p.UserData.Equip.SuitReward = append(p.UserData.Equip.SuitReward, data)
	p.SaveEquip()
}

// GetSuitPosReward 获得套装部位奖励
func GetSuitPosReward(p *player.Player, suit, pos uint32) (msg.ErrCode, []*template.SimpleItem, *template.JEquip) {
	equips := template.GetEquipTemplate().GetSuitPosEquip(suit, pos)
	if len(equips) == 0 {
		return msg.ErrCode_SUIT_NOT_EXIST, nil, nil
	}

	maxQuality := getSuitPosMaxquality(p, suit, pos)
	if maxQuality == nil {
		return msg.ErrCode_NO_SUIT_REWARD, nil, nil
	}

	suitInfo := getSuit(p, suit)

	var minRarity uint32 = 0
	var minSmallStage uint32 = 0
	if suitInfo != nil {
		posIndex := -1
		for i := 0; i < len(suitInfo.PosData); i++ {
			if suitInfo.PosData[i].Pos == pos {
				if posEquip := template.GetEquipTemplate().GetEquip(suitInfo.PosData[i].EquipId); posEquip != nil {
					minRarity = posEquip.Data.Rarity
					minSmallStage = posEquip.Data.SmallStage
				}

				if !maxQuality.HasHighQuality2(minRarity, minSmallStage) {
					return msg.ErrCode_NO_SUIT_REWARD, nil, nil
				}

				posIndex = i
				break
			}
		}
		if posIndex == -1 {
			suitInfo.PosData = append(suitInfo.PosData, model.NewSuitPosInfo(pos, maxQuality.Data.Id))
		} else {
			suitInfo.PosData[posIndex].EquipId = maxQuality.Data.Id
		}
		p.SaveEquip()
	} else {
		suitInfo = model.NewSuitInfo(suit)
		suitInfo.PosData = append(suitInfo.PosData, model.NewSuitPosInfo(pos, maxQuality.Data.Id))
		addSuit(p, suitInfo)
	}
	temp := template.GetEquipTemplate().GetSuitRangeEquip(suit, pos, minRarity, minSmallStage, maxQuality.Data.Rarity, maxQuality.Data.SmallStage)
	rewardItems := template.GetEquipTemplate().GetEquipRewardItems(temp)

	var notifyItems []uint32
	for i := 0; i < len(rewardItems); i++ {
		addItems := AddItem(p.GetUserId(),
			rewardItems[i].ItemId,
			int32(rewardItems[i].ItemNum),
			publicconst.SuitRewardAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	//event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(playerData, notifyItems, ListenNotifyClientItemEventEvent))

	updateClientItemsChange(p.GetUserId(), notifyItems)

	return msg.ErrCode_SUCC, rewardItems, maxQuality
}

// PutInSuit 放入套装
func PutInSuit(p *player.Player, suitId uint32) msg.ErrCode {
	var equipIds []uint32
	var skillId uint32 = 0
	count := 0
	for i := 1; i <= 6; i++ {
		if maxQuality := getSuitPosMaxquality(p, suitId, uint32(i)); maxQuality != nil {
			skillId = maxQuality.Data.SkillId
			equipIds = append(equipIds, maxQuality.Data.Id)
			count += 1
		} else {
			equipIds = append(equipIds, 0)
		}
	}

	if count != 6 {
		skillId = 0
	}

	var data *model.EquipSuit
	data = getEquipSuit(p, suitId)
	if data != nil {
		if !tools.ListUint32Equal(equipIds, data.EquipIds) {
			data.EquipIds = equipIds
			data.SkillId = skillId
		}
	} else {
		data = model.NewEquipSuit(suitId, equipIds, skillId)
		p.UserData.Equip.EquipSuits = append(p.UserData.Equip.EquipSuits, data)
	}
	p.SaveEquip()

	notifySuitChange(p, data)
	return msg.ErrCode_SUCC
}

// refreshSuit 刷新套装
func refreshSuit(p *player.Player, suitId uint32) {
	suit := getEquipSuit(p, suitId)
	if suit == nil {
		return
	}

	var equipIds []uint32
	var skillId uint32 = 0
	count := 0
	for i := 1; i <= 6; i++ {
		if maxQuality := getSuitPosMaxquality(p, suitId, uint32(i)); maxQuality != nil {
			skillId = maxQuality.Data.SkillId
			equipIds = append(equipIds, maxQuality.Data.Id)
			count += 1
		} else {
			equipIds = append(equipIds, 0)
		}
	}

	if count != 6 {
		skillId = 0
	}

	// 没有变化
	if tools.ListUint32Equal(equipIds, suit.EquipIds) {
		return
	}

	suit.EquipIds = equipIds
	suit.SkillId = skillId
	p.SaveEquip()

	notifySuitChange(p, suit)
}

// UseSuit 请求使用套装
func UseSuit(p *player.Player, suitId uint32) msg.ErrCode {
	data := getEquipSuit(p, suitId)
	if data == nil {
		return msg.ErrCode_SUIT_NOT_EXIST
	}
	refreshSuit(p, suitId)

	ids := getEquipPosIds(p)
	if tools.ListUint32Equal(ids, data.EquipIds) {
		return msg.ErrCode_SUCC
	}

	p.UserData.Equip.UseEquipSuit = suitId
	p.UserData.Equip.SyncEquipPos = true

	for i := 1; i <= 6; i++ {
		equipId := data.EquipIds[i-1]
		if equipId == 0 {
			continue
		}
		jEquip := template.GetEquipTemplate().GetEquip(equipId)
		if jEquip == nil {
			continue
		}
		data := getEquipPos(p, uint32(i))
		if data == nil {
			data = model.NewEquipPos(uint32(i), equipId)
			p.UserData.Equip.EquipPosData = append(p.UserData.Equip.EquipPosData, data)
		} else {
			data.EquipId = equipId
		}
		data.Attr = make(map[uint32]*model.Attr)
		InitAttr(data.Attr, jEquip.InitAttr)
		SetAttrLevelValue(data.Attr, jEquip.GetLevelAttr(data.Level))
		calcEquipAffix(p, data)
		calcEquipPosAttr(p, data)
	}

	p.SaveEquip()
	// 通知装备槽位变化
	notifyMsg := &msg.NotifyEquipSlotChange{}
	notifyMsg.Data = ToProtocolEquipPosList(p.UserData.Equip.EquipPosData)
	p.SendNotify(notifyMsg)

	// 计算战力
	GlobalAttrChange(p, true)
	return msg.ErrCode_SUCC
}

/*
 * OneKeyUpLvEquip
 *  @Description: 一键升级装备 槽位
 *  @param playerData
 */
func OneKeyUpLvEquip(packetId uint32, p *player.Player) {
	res := &msg.ResponseAllEquipUpgrade{
		Result: msg.ErrCode_SUCC,
	}
	defer p.SendResponse(packetId, res, res.Result)
	var (
		upNum        uint32
		ntfItems     []uint32
		costMap      = make(map[uint32]uint32)
		filterMap    = make(map[int]bool, 6)
		changeEquips = make(map[uint32]*model.EquipPos)
		notifyMsg    = new(msg.NotifyEquipSlotChange)
	)

LOOP:
	for {
		continueNum := 0
		for i := 1; i <= 6; i++ {
			flag, ok := filterMap[i]
			if !ok {
				filterMap[i] = false
			} else {
				if flag {
					continueNum++
					continue
				}
			}

			// get equip
			equipPos := getEquipPos(p, uint32(i))
			if equipPos == nil {
				filterMap[i] = true
				continue
			}

			equipCfg := template.GetEquipTemplate().GetEquip(equipPos.EquipId)
			if equipCfg == nil {
				log.Error("equip cfg nil", zap.Uint32("equipId", equipPos.EquipId))
				filterMap[i] = true
				continue
			}

			if equipPos.Level >= equipCfg.Data.LevelMax {
				filterMap[i] = true
				continue
			}

			equipLvCfg := template.GetEquipLevelTemplate().GetEquipLevel(uint32(i), equipPos.Level)
			if equipLvCfg == nil {
				log.Error("equip lv cfg nil", zap.Uint32("equipId", equipPos.EquipId), zap.Uint32("level", equipPos.Level))
				filterMap[i] = true
				continue
			}

			enough := true
			for _, v := range equipLvCfg.CostItem {
				if !EnoughItem(p.GetUserId(), v.ItemId, costMap[v.ItemId]+v.ItemNum) {
					enough = false
					break
				}
			}
			if !enough {
				filterMap[i] = true
				continue
			}

			for _, v := range equipLvCfg.CostItem {
				costMap[v.ItemId] += v.ItemNum
			}

			upNum++
			equipPos.Level++
			changeEquips[uint32(i)] = equipPos
		}

		if continueNum == 6 {
			break LOOP
		}
	}

	if len(costMap) > 0 {
		for k, v := range costMap {
			CostItem(p.GetUserId(), k, v, publicconst.EquipAutoUpgradeCostItem, false)
			ntfItems = append(ntfItems, k)
		}

		if len(changeEquips) > 0 {
			for _, v := range changeEquips {
				notifyMsg.Data = append(notifyMsg.Data, ToProtocolEquipPos(v))

				equipCfg := template.GetEquipTemplate().GetEquip(v.EquipId)
				if equipCfg == nil {
					log.Error("equip cfg nil", zap.Uint32("equipId", v.EquipId))
					continue
				}

				setEquipPosLevelAttr(p, v, equipCfg, false)

				// //tdaItems := make([]*tda.Item, 0, len(costMap))
				// for id, nm := range costMap {
				// 	tdaItems = append(tdaItems, &tda.Item{ItemId: strconv.Itoa(int(id)), ItemNum: nm})
				// }
				// tda.TdaEquipUpgrade(p.ChannelId, p.TdaCommonAttr, k, v.Level, equipCfg.GenRarity(), tdaItems)
			}
		}
		GlobalAttrChange(p, true)
		p.SaveEquip()

		updateClientItemsChange(p.GetUserId(), ntfItems)

		if len(changeEquips) > 0 {
			p.SendNotify(notifyMsg)
		}

		UpdateTask(p, true, publicconst.TASK_COND_UPGRADE_EQUIP, upNum)
		processHistoryData(p, publicconst.TASK_COND_UPGRADE_EQUIP, 0, upNum)

		UpdateTask(p, true, publicconst.TASK_COND_ANY_EQUIP_LEVEL)
		UpdateTask(p, true, publicconst.TASK_COND_ALL_EQUIP_POS_LEVEL)
	}

	return
}

func GetGemAffixRealTypeAndAffix(p *player.Player, uuid uint64) *common2.UintPair {
	gem_id := GetGemId(uuid)
	j_gem := template.GetGemTemplate().GetGem(gem_id)
	if j_gem == nil {
		log.Error("get gem err",
			zap.Uint64("account_id", p.GetUserId()),
			zap.Uint64("uuid", uuid))
		return &common2.UintPair{First: 0, Second: 0}
	}

	affix_type := GetGemAffixType(uuid)
	if affix_type == 0 {
		// 没有洗练 取配置词条
		if len(j_gem.Attrs) > 0 {
			affix_type = uint32(msg.GemAffixType_GemAffix_Attr)
			attr_id := j_gem.Attrs[0].Id
			return &common2.UintPair{First: affix_type, Second: attr_id}
		} else if j_gem.Skills > 0 {
			affix_type = uint32(msg.GemAffixType_GemAffix_Skill)
			skill_id := j_gem.Skills
			j_skill := template.GetSkillTemplate().GetSkill(int(skill_id))
			if j_skill == nil {
				log.Error("get skill err when GetGemAffixRealType",
					zap.Uint64("account_id", p.GetUserId()),
					zap.Uint64("uuid", uuid))
				return &common2.UintPair{First: 0, Second: 0}
			}
			return &common2.UintPair{First: affix_type, Second: uint32(j_skill.BaseId)}
		} else if j_gem.Buffs > 0 {
			affix_type = uint32(msg.GemAffixType_GemAffix_Buf)
			buff_id := j_gem.Buffs
			j_buff := template.GetBuffTemplate().GetCfg(int(buff_id))
			if j_buff == nil {
				log.Error("get buff err when GetGemAffixRealType",
					zap.Uint64("account_id", p.GetUserId()),
					zap.Uint64("uuid", uuid))
				return &common2.UintPair{First: 0, Second: 0}
			}

			return &common2.UintPair{First: affix_type, Second: uint32(j_buff.Group)}
		} else {
			// 有洗练
			affix_idx := GetGemAffixIdx(uuid)
			j_refresh := template.GetGemRefreshTemplate().GetGemRefresh(gem_id)
			if j_refresh == nil {
				log.Error("get gem refresh err when GetGemAffixRealType",
					zap.Uint64("account_id", p.GetUserId()),
					zap.Uint64("uuid", uuid))
				return &common2.UintPair{First: 0, Second: 0}
			}
			switch affix_type {
			case uint32(msg.GemAffixType_GemAffix_Attr):
				if affix_idx >= uint32(len(j_refresh.Attrs)) {
					log.Error("affix idx err when GetGemAffixRealType",
						zap.Uint64("account_id", p.GetUserId()),
						zap.Uint64("uuid", uuid),
						zap.Uint32("affix_type", affix_type),
						zap.Uint32("affix_idx", affix_idx),
					)
					return &common2.UintPair{First: 0, Second: 0}
				}
				return &common2.UintPair{First: affix_type, Second: j_refresh.Attrs[affix_idx].Id}
			case uint32(msg.GemAffixType_GemAffix_Skill):
				if affix_idx >= uint32(len(j_refresh.Skills)) {
					log.Error("affix idx err when GetGemAffixRealType",
						zap.Uint64("account_id", p.GetUserId()),
						zap.Uint64("uuid", uuid),
						zap.Uint32("affix_type", affix_type),
						zap.Uint32("affix_idx", affix_idx),
					)
					return &common2.UintPair{First: 0, Second: 0}
				}

				skill_id := j_refresh.Skills[affix_idx].Id
				j_skill := template.GetSkillTemplate().GetSkill(int(skill_id))
				if j_skill == nil {
					log.Error("get skill err when GetGemAffixRealType",
						zap.Uint64("account_id", p.GetUserId()),
						zap.Uint64("uuid", uuid))
					return &common2.UintPair{First: 0, Second: 0}
				}
				return &common2.UintPair{First: affix_type, Second: uint32(j_skill.BaseId)}
			case uint32(msg.GemAffixType_GemAffix_Buf):
				if affix_idx >= uint32(len(j_refresh.Buffs)) {
					log.Error("affix idx err when GetGemAffixRealType",
						zap.Uint64("account_id", p.GetUserId()),
						zap.Uint64("uuid", uuid),
						zap.Uint32("affix_type", affix_type),
						zap.Uint32("affix_idx", affix_idx),
					)
					return &common2.UintPair{First: 0, Second: 0}
				}

				buff_id := j_refresh.Buffs[affix_idx].Id
				j_buff := template.GetBuffTemplate().GetCfg(int(buff_id))
				if j_buff == nil {
					log.Error("get buff err when GetGemAffixRealType",
						zap.Uint64("account_id", p.GetUserId()),
						zap.Uint64("uuid", uuid))
					return &common2.UintPair{First: 0, Second: 0}
				}
				return &common2.UintPair{First: affix_type, Second: uint32(j_buff.Group)}
			default:
				log.Error("get affix type err when GetGemAffixRealType",
					zap.Uint64("account_id", p.GetUserId()),
					zap.Uint64("uuid", uuid),
					zap.Uint32("affix_type", affix_type),
				)
				return &common2.UintPair{First: 0, Second: 0}
			}
		}
	}
	return &common2.UintPair{First: 0, Second: 0}
}

func SocketGem(p *player.Player, pos uint32, uuid uint64, slot int32) (msg.ErrCode, uint32) {
	if pos >= uint32(msg.EquipPos_EquipPos_Max) ||
		pos <= uint32(msg.EquipPos_EquipPos_None) {
		log.Error("pos error when SocketGem", zap.Uint64("AccountId", p.GetUserId()), zap.Uint32("pos", pos))
		return msg.ErrCode_Gem_Pos_Error, 0
	}
	gem_new := getGem(p, uuid)
	if gem_new == nil {
		return msg.ErrCode_Gem_Nil, 0
	}
	gem_id_new := GetGemId(uuid)
	j_gem_new := template.GetGemTemplate().GetGem(gem_id_new)
	if j_gem_new == nil {
		log.Error("get gem cfg err when SocketGem",
			zap.Uint64("account_id", p.GetUserId()),
			zap.Uint64("uuid", uuid),
		)
		return msg.ErrCode_CONFIG_NIL, 0
	}

	if j_gem_new.Pos != pos {
		log.Error("pos not match when SocketGem",
			zap.Uint64("AccountId", p.GetUserId()), zap.Uint32("pos", pos), zap.Uint64("gem uuid", uuid))
		return msg.ErrCode_Gem_Pos_Error, 0
	}

	new_affix := GetGemAffixRealTypeAndAffix(p, uuid)
	pos_gems := p.UserData.Equip.GemPos[pos-1]
	if slot >= int32(len(pos_gems)) {
		log.Error("slot error when SocketGem",
			zap.Uint64("AccountId", p.GetUserId()), zap.Uint32("pos", pos), zap.Int32("slot", slot))
		return msg.ErrCode_Gem_Slot_Error, 0
	}

	first_nil := int32(-1)   // 第一个空位
	socket_slot := int32(-1) // 最终镶嵌位置
	if slot > 0 {            // 主动替换
		for i, old_gem := range pos_gems {
			if old_gem == 0 {
				continue
			}
			old_affix := GetGemAffixRealTypeAndAffix(p, old_gem)
			if i == int(slot) { // 选中slot当前宝石是否低级
				if !new_affix.EqualWith(old_affix) {
					continue
				}
				j_gem_old := template.GetGemTemplate().GetGem(GetGemId(old_gem))
				if j_gem_old == nil {
					log.Error("get gem cfg err when SocketGem",
						zap.Uint64("account_id", p.GetUserId()),
						zap.Uint64("uuid", old_gem),
					)
					return msg.ErrCode_CONFIG_NIL, 0
				}

				if j_gem_new.Quality <= j_gem_old.Quality {
					return msg.ErrCode_Gem_Affix_Low, 0
				}
			} else { // 其他位置是否有同词条
				if new_affix != old_affix {
					continue
				}

				if new_affix.EqualWith(old_affix) {
					return msg.ErrCode_Gem_Affix_Repeat, 0
				}
			}
		}
		socket_slot = slot // 可以替换
	} else { // 自动查找镶嵌位
		// 是否已有同词条
		for i, old_gem := range pos_gems {
			if old_gem == 0 {
				if first_nil < 0 {
					first_nil = int32(i)
				}
				continue
			}
			old_affix := GetGemAffixRealTypeAndAffix(p, old_gem)
			if !new_affix.EqualWith(old_affix) {
				continue
			}
			// 词条相同 判断宝石品质 该判断机制存疑!
			// 如果新宝石品质高于旧宝石品质,则替换
			gem_id_old := GetGemId(old_gem)
			j_gem_old := template.GetGemTemplate().GetGem(gem_id_old)
			if j_gem_old == nil {
				log.Error("get gem err",
					zap.Uint64("account_id", p.GetUserId()),
					zap.Uint64("uuid", old_gem))
				continue
			}

			if j_gem_new.Quality <= j_gem_old.Quality {
				return msg.ErrCode_Gem_Affix_Low, 0
			}

			// 找到替换位
			socket_slot = int32(i)
			break
		}
	}
	if socket_slot < 0 {
		// 没有替换位
		if first_nil < 0 {
			return msg.ErrCode_Gem_Socket_Error, 0
		}
		socket_slot = first_nil
	}

	pos_gems[socket_slot] = uuid //设置镶嵌宝石
	p.SaveEquip()
	GlobalAttrChange(p, true)
	UpdateTask(p, true, publicconst.TASK_COND_PUT_ON_DISK, 1) // 镶嵌XX个磁盘
	return msg.ErrCode_SUCC, uint32(socket_slot)
}

func UnSocketGem(p *player.Player, pos, slot uint32) msg.ErrCode {
	if pos >= uint32(msg.EquipPos_EquipPos_Max) ||
		pos <= uint32(msg.EquipPos_EquipPos_None) {
		log.Error("pos error when SocketGem", zap.Uint64("AccountId", p.GetUserId()), zap.Uint32("pos", pos))
		return msg.ErrCode_Gem_Pos_Error
	}
	if slot >= uint32(len(p.UserData.Equip.GemPos[pos-1])) {
		log.Error("slot error when SocketGem", zap.Uint64("AccountId", p.GetUserId()), zap.Uint32("pos", pos), zap.Uint32("slot", slot))
		return msg.ErrCode_Gem_Slot_Error
	}

	if p.UserData.Equip.GemPos[pos-1][slot] == 0 {
		log.Error("gem not socket when UnSocketGem", zap.Uint64("AccountId", p.GetUserId()), zap.Uint32("pos", pos), zap.Uint32("slot", slot))
		return msg.ErrCode_SUCC
	}
	p.UserData.Equip.GemPos[pos-1][slot] = 0
	p.SaveEquip()
	GlobalAttrChange(p, true)

	return msg.ErrCode_SUCC
}

func LockGem(p *player.Player, uuid uint64, lock bool) msg.ErrCode {
	gem := getGem(p, uuid)
	if gem == nil {
		return msg.ErrCode_Gem_Nil
	}

	j_gem := template.GetGemTemplate().GetGem(GetGemId(uuid))
	if j_gem == nil {
		log.Error("get gem err when LockGem",
			zap.Uint64("account_id", p.GetUserId()),
			zap.Uint64("uuid", uuid))
		return msg.ErrCode_CONFIG_NIL
	}

	pos := j_gem.Pos
	if pos >= uint32(len(p.UserData.Equip.GemPos)) {
		log.Error("pos error when LockGem",
			zap.Uint64("account_id", p.GetUserId()),
			zap.Uint64("uuid", uuid),
			zap.Uint32("pos", pos))
		return msg.ErrCode_Gem_Pos_Error
	}

	for _, v := range p.UserData.Equip.GemPos[pos] {
		if v == uuid {
			return msg.ErrCode_Gem_In_Use
		}
	}

	gem.Lock = lock
	p.SaveEquip()

	return msg.ErrCode_SUCC
}

type NumGems struct {
	TotalNum uint32
	Gems     []*model.GemBagSlot
}

func ComposeGem(p *player.Player) (msg.ErrCode, []*msg.Gem) {
	change_gems := make(map[uint64]*model.GemBagSlot, 0)
	get_gems := make([]*msg.Gem, 0)
	gem_by_pos_quality := make(map[uint32]map[uint32]*NumGems)

	rm_gems := make(map[uint64]struct{})
	used_gems := make(map[uint64]struct{})
	new_gems := make(map[uint64]uint32)

	for _, pos := range p.UserData.Equip.GemPos {
		for _, gem := range pos {
			if gem == 0 {
				continue
			}
			used_gems[gem] = struct{}{}
		}
	}

	for _, v := range p.UserData.Equip.GemBag {
		gem_id := GetGemId(v.Uuid)
		j_gem := template.GetGemTemplate().GetGem(gem_id)
		if j_gem == nil {
			log.Error("get gem cfg err when ComposeGem",
				zap.Uint64("account_id", p.GetUserId()),
				zap.Uint32("gem_id", gem_id))
			continue
		}

		if v.Lock {
			continue
		}

		if j_gem.Quality >= uint32(msg.GemQuality_GemQuality_Max-1) {
			continue
		}

		if _, ok := gem_by_pos_quality[j_gem.Pos]; !ok {
			gem_by_pos_quality[j_gem.Pos] = make(map[uint32]*NumGems)
		}
		if _, ok := gem_by_pos_quality[j_gem.Pos][j_gem.Quality]; !ok {
			gem_by_pos_quality[j_gem.Pos][j_gem.Quality] = &NumGems{
				TotalNum: 0,
				Gems:     make([]*model.GemBagSlot, 0),
			}
		}
		ng := gem_by_pos_quality[j_gem.Pos][j_gem.Quality]
		ng.Gems = append(gem_by_pos_quality[j_gem.Pos][j_gem.Quality].Gems, v)
		if v.Lock { // 锁定不计数
			if v.Num > 0 {
				ng.TotalNum += v.Num - 1
			}
		} else {
			if _, ok := used_gems[v.Uuid]; ok { // 没锁定 穿戴不计数
				if v.Num > 0 {
					ng.TotalNum += v.Num - 1
				}
			} else {
				ng.TotalNum += v.Num
			}

		}
	}

	if len(gem_by_pos_quality) == 0 {
		log.Error("no gem composed when ComposeGem",
			zap.Uint64("account_id", p.GetUserId()))
		return msg.ErrCode_SUCC, get_gems
	}

	// sort
	for pos, quality_map := range gem_by_pos_quality {
		for quality, ng := range quality_map {
			sort.Slice(ng.Gems, func(i, j int) bool {
				return ng.Gems[i].Num < ng.Gems[j].Num
			})
			total_cnt := ng.TotalNum / template.GetSystemItemTemplate().GemComposeCost // 可合成数
			com_cnt := uint32(0)                                                       // 当前合成数
			cost_cnt := uint32(0)                                                      // 当前合成消耗

			for com_cnt < total_cnt { // 未完成
				//rm_idx := 0 不支持锁定和穿戴保留逻辑
				for _, gem := range ng.Gems {
					if gem.Num == 0 {
						//rm_idx++
						if _, ok := rm_gems[gem.Uuid]; !ok {
							rm_gems[gem.Uuid] = struct{}{}
						}
						continue
					}

					gem.Num--
					cost_cnt++
					if gem.Num == 0 {
						//rm_idx++
						if _, ok := rm_gems[gem.Uuid]; !ok {
							rm_gems[gem.Uuid] = struct{}{}
						}
					}

					if _, ok := change_gems[gem.Uuid]; !ok {
						change_gems[gem.Uuid] = gem
					}

					// 合成一个
					if cost_cnt >= template.GetSystemItemTemplate().GemComposeCost {
						com_cnt++
						cost_cnt = 0
						break
					}
				}
			}

			// 生成新宝石
			for total_cnt > 0 {
				total_cnt--

				quality_add := quality + 1
				j_rand_gem := template.GetGemTemplate().RandGem(p.Rand, pos, quality_add)
				if j_rand_gem == nil {
					log.Error("rand gem err when ComposeGem",
						zap.Uint64("account_id", p.GetUserId()),
						zap.Uint32("pos", pos),
						zap.Uint32("quality_add", quality_add))
					continue
				}
				new_gem := AddGem(p, j_rand_gem.Id, 1, 0, 0, false)
				if _, ok := new_gems[new_gem.Uuid]; !ok {
					new_gems[new_gem.Uuid] = 1
				} else {
					new_gems[new_gem.Uuid] += 1
				}
				if _, ok := change_gems[new_gem.Uuid]; !ok {
					change_gems[new_gem.Uuid] = new_gem
				}
			}
		}
	}
	for k := range rm_gems {
		delete(p.UserData.Equip.GemBag, k)
	}
	if len(change_gems) > 0 {
		p.SaveEquip()
		NotifyGemmapChange(p, change_gems)
	}

	for k, v := range new_gems {
		get_gems = append(get_gems, &msg.Gem{
			Uuid: k,
			Num:  v,
		})
	}
	return msg.ErrCode_SUCC, get_gems
}

func RefreshGem(p *player.Player, uuid uint64) (msg.ErrCode, uint64) {
	gem := getGem(p, uuid)
	if gem == nil {
		return msg.ErrCode_Gem_Nil, 0
	}

	gem_id := GetGemId(uuid)

	j_gem := template.GetGemTemplate().GetGem(gem_id)
	if j_gem == nil {
		log.Error("get gem err when RefreshGem",
			zap.Uint64("account_id", p.GetUserId()),
			zap.Uint64("uuid", uuid))
		return msg.ErrCode_CONFIG_NIL, 0
	}

	if j_gem.Pos >= uint32(len(p.UserData.Equip.GemPos)) {
		log.Error("pos error when RefreshGem",
			zap.Uint64("account_id", p.GetUserId()),
			zap.Uint64("uuid", uuid),
			zap.Uint32("pos", j_gem.Pos))
		return msg.ErrCode_Gem_Pos_Error, 0
	}

	j_refresh := template.GetGemRefreshTemplate().GetGemRefresh(GetGemId(uuid))
	if j_refresh == nil {
		log.Error("get gem refresh err when RefreshGem",
			zap.Uint64("account_id", p.GetUserId()),
			zap.Uint64("uuid", uuid))
		return msg.ErrCode_CONFIG_NIL, 0
	}

	// 是否已装备
	pos_gems := p.UserData.Equip.GemPos[j_gem.Pos]
	for _, v := range pos_gems {
		if v == uuid {
			if gem.Num <= 1 {
				return msg.ErrCode_Gem_In_Use, 0
			}
			break
		}
	}

	af, ai := template.GetGemRefreshTemplate().RandAffix(p.Rand, gem_id)
	if af == 0 {
		log.Error("get gem refresh err when RandAffix",
			zap.Uint64("account_id", p.GetUserId()),
			zap.Uint64("uuid", uuid),
		)
		return msg.ErrCode_SYSTEM_ERROR, 0
	}

	new_gem := AddGem(p, gem_id, 1, ai, af, true)
	DelGem(p, uuid, 1)
	return msg.ErrCode_SUCC, new_gem.Uuid
}

func DelGem(p *player.Player, uuid uint64, num uint32) *model.GemBagSlot {
	gem := getGem(p, uuid)
	if gem == nil {
		return nil
	}

	if gem.Num > num {
		gem.Num -= num
	} else {
		delete(p.UserData.Equip.GemBag, uuid)
	}
	p.SaveEquip()
	return gem
}

func getEquipSuit(p *player.Player, suitId uint32) *model.EquipSuit {
	for i := 0; i < len(p.UserData.Equip.EquipSuits); i++ {
		if p.UserData.Equip.EquipSuits[i].SuitId == suitId {
			return p.UserData.Equip.EquipSuits[i]
		}
	}
	return nil
}

func getEquipPosIds(p *player.Player) []uint32 {
	var ret []uint32
	for i := 1; i <= 6; i++ {
		if data := getEquipPos(p, uint32(i)); data != nil {
			ret = append(ret, data.EquipId)
		} else {
			ret = append(ret, 0)
		}
	}
	return ret
}

func notifySuitChange(p *player.Player, suit *model.EquipSuit) {
	// 通知客户端刷新
	res := &msg.NotifyEquipSuitChange{}
	res.Data = ToProtocolEquipSuit(suit)
	p.SendNotify(res)
}

// syncSuit 同步套装
func SyncSuit(p *player.Player) {
	if !p.UserData.Equip.SyncEquipPos {
		return
	}

	suit := getEquipSuit(p, p.UserData.Equip.UseEquipSuit)
	if suit == nil {
		return
	}

	ids := getEquipPosIds(p)
	if tools.ListUint32Equal(ids, suit.EquipIds) {
		return
	}

	suit.EquipIds = ids
	p.SaveEquip()
	notifySuitChange(p, suit)
}

/*
 * OneKeyUpgradeEquip
 *  @Description: 一键升级装备 装备品质 消耗装备获得品质高的装备
 *  @param playerData
 */
func OneKeyUpgradeEquip(p *player.Player) {
	var (
		ack                = &msg.ResponseEquipUpgradeStage{Result: msg.ErrCode_SUCC, IsAuto: 1, GetEquip: make([]*msg.Equip, 0)}
		equipChangeNtf     = &msg.NotifyEquipChange{Data: make([]*msg.Equip, 0)}
		equipSlotChangeNtf = &msg.NotifyEquipSlotChange{Data: make([]*msg.UseEquip, 0)}
	)

	curTime := time.Now().UnixMilli()
	if curTime-p.EquipUpgradeStageTime < 500 {
		ack.Result = msg.ErrCode_EQUIP_BATCH_UPGRADE_STAGE_IN_CD
		p.SendNotify(ack)
		return
	}

	add_map := make(map[uint32]*model.Equip)
	reduce_map := make(map[uint32]*model.Equip)
	pos_pos := oneKeyUpgradeEquipPos(p, add_map, reduce_map)
	oneKeyUpgradeEquipBag(p, add_map, reduce_map)

	msg_add := ToProtocolEquipMap(add_map)

	msg_reduece := ToProtocolEquipMap(reduce_map)

	ack.GetEquip = append(ack.GetEquip, msg_add...)

	equipChangeNtf.Data = append(equipChangeNtf.Data, msg_add...)
	equipChangeNtf.Data = append(equipChangeNtf.Data, msg_reduece...)

	equipSlotChangeNtf.Data = append(equipSlotChangeNtf.Data, ToProtocolEquipPosList(pos_pos)...)

	UpdateTask(p, true, publicconst.TASK_COND_COMPOSE_EQUIP, uint32(len(msg_add)))
	UpdateTask(p, true, publicconst.TASK_COND_EQUIP_RARITY_NUM)

	// TODO tda

	SyncSuit(p)

	p.EquipUpgradeStageTime = curTime

	p.SaveEquip()
	GlobalAttrChange(p, true)
	p.SendNotify(equipChangeNtf)
	p.SendNotify(equipSlotChangeNtf)
	p.SendNotify(ack)
}

// UpgradeStageEquip 升阶装备
func UpgradeStageEquip(p *player.Player, equipId uint32, cost_equips []uint32) (msg.ErrCode, []*msg.UseEquip, []*model.Equip, []*msg.Equip) {
	change_equip_ids := make(map[uint32]struct{})
	equip := getEquip(p, equipId)
	if equip == nil {
		return msg.ErrCode_EQUIP_NOT_EXIST, nil, nil, nil
	}

	jEquip := template.GetEquipTemplate().GetEquip(equipId)
	if jEquip == nil {
		return msg.ErrCode_EQUIP_NOT_EXIST, nil, nil, nil
	}

	if jEquip.Data.OpenRank == 2 {
		return msg.ErrCode_EQUIP_NOT_UPGRADE_STAGE, nil, nil, nil
	}

	num := uint32(0)
	ce_id_counts := make(map[uint32]uint32)
	for _, ce_id := range cost_equips { //TODO self
		if ce_id == 0 {
			continue
		}

		num++
		ce_id_counts[ce_id]++
	}

	if num+1 != template.GetEquipUpConditionTemplate().GetNeedNum(jEquip.Data.Pos, jEquip.Data.Rarity, jEquip.Data.SmallStage) {
		return msg.ErrCode_EQUIP_NO_ENOUGH_EQUIP, nil, nil, nil
	}

	// 装备中
	posEquip := getEquipPos(p, jEquip.Data.Pos)
	for ce_id, ce_cnt := range ce_id_counts {
		ce_data := getEquip(p, ce_id)
		if ce_data == nil {
			return msg.ErrCode_EQUIP_NOT_EXIST, nil, nil, nil
		}
		j_ce := template.GetEquipTemplate().GetEquip(ce_id)
		if j_ce == nil {
			return msg.ErrCode_EQUIP_NOT_EXIST, nil, nil, nil
		}
		if j_ce.Data.Pos != jEquip.Data.Pos ||
			j_ce.Data.Rarity != jEquip.Data.Rarity ||
			j_ce.Data.SmallStage != jEquip.Data.SmallStage {
			return msg.ErrCode_EQUIP_NOT_UPGRADE_STAGE, nil, nil, nil
		}

		// 如果装备中 需保留一个
		if posEquip != nil {
			if posEquip.EquipId == ce_id {
				if ce_data.Num <= ce_cnt {
					return msg.ErrCode_EQUIP_NO_ENOUGH_EQUIP, nil, nil, nil
				}
			}
		}
	}
	for ce_id, ce_cnt := range ce_id_counts {
		ce_data := getEquip(p, ce_id)
		if ce_data == nil {
			continue
		}
		ce_data.Num -= ce_cnt
		updateEquip(p, ce_data, false)
		change_equip_ids[ce_id] = struct{}{}
	}

	get_e_map := make(map[uint32]*msg.Equip)
	var getEquips []*msg.Equip
	var ret []*msg.UseEquip

	if posEquip != nil {
		jNewEquip := template.GetEquipTemplate().GetEquip(jEquip.Data.Rank)
		jPosEquip := template.GetEquipTemplate().GetEquip(posEquip.EquipId)
		if posEquip.EquipId == equipId && jNewEquip.HasHighQuality(jPosEquip) {
			posEquip.EquipId = jEquip.Data.Rank
			posEquip.UpdateTime = tools.GetCurTime()
			if posEquip.EquipId == 0 {
				log.Error("UpgradeStageEquip err", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("equipId", posEquip.EquipId))
				//log.Errorf("UpgradeStageEquip accountid:%v equipId:%v", playerData.GetUserId(), posEquip.EquipId)
			}
			updateEquipPos(p, posEquip, jNewEquip, true)
			ret = append(ret, ToProtocolEquipPos(posEquip))
		}
	}

	// 原装备要删除
	if equip.Num > 0 {
		equip.Num--
		updateEquip(p, equip, false)
	}

	// if _, ok := get_e_map[equipId]; !ok {
	// 	get_e_map[equipId] = &msg.Equip{
	// 		Id:  equipId,
	// 		Num: 1,
	// 	}
	// } else {
	// 	get_e_map[equipId].Num++
	// }

	AddEquip(p, jEquip.Data.Rank, 1, true)
	get_e_map[jEquip.Data.Rank] = &msg.Equip{
		Id:  jEquip.Data.Rank,
		Num: 1,
	}
	//if len(teamIds) > 0 {
	//	ServMgr.GetTeamService().updateTeamEquip(playerData, teamIds, jEquip.Data.Rank)
	//}

	UpdateTask(p, true, publicconst.TASK_COND_COMPOSE_EQUIP, 1)
	UpdateTask(p, true, publicconst.TASK_COND_EQUIP_RARITY_NUM)

	change_equip_ids[jEquip.Data.Rank] = struct{}{}

	var change_equips []*model.Equip

	for change_id := range change_equip_ids {
		if equip := getEquip(p, change_id); equip != nil {
			change_equips = append(change_equips, equip)
		} else {
			change_equips = append(change_equips, model.NewEquip(change_id, 0))
		}
	}
	for _, ge := range get_e_map {
		getEquips = append(getEquips, ge)
	}
	return msg.ErrCode_SUCC, ret, change_equips, getEquips
}

// getEquip 获得装备
func getEquip(p *player.Player, equipId uint32) *model.Equip {
	for i := 0; i < len(p.UserData.Equip.EquipData); i++ {
		if p.UserData.Equip.EquipData[i].Id == equipId {
			return p.UserData.Equip.EquipData[i]
		}
	}
	return nil
}

// UpgradeEquipPos 升级装备部位
func UpgradeEquipPos(p *player.Player, pos uint32) (msg.ErrCode, *model.EquipPos) {
	equipPos := getEquipPos(p, pos)
	if equipPos == nil {
		return msg.ErrCode_EQUIP_POS_NO_EQUIP, nil
	}

	jEquip := template.GetEquipTemplate().GetEquip(equipPos.EquipId)
	if jEquip == nil {
		return msg.ErrCode_EQUIP_NOT_EXIST, nil
	}
	if equipPos.Level >= jEquip.Data.LevelMax {
		return msg.ErrCode_EQUIP_LEVEL_FULL, nil
	}

	equipLevel := template.GetEquipLevelTemplate().GetEquipLevel(pos, equipPos.Level)
	if equipLevel == nil || len(equipLevel.CostItem) == 0 {
		return msg.ErrCode_EQUIP_LEVEL_FULL, nil
	}

	//tdaItems := make([]*tda.Item, 0, len(equipLevel.CostItem))
	for i := 0; i < len(equipLevel.CostItem); i++ {
		if !EnoughItem(p.GetUserId(), equipLevel.CostItem[i].ItemId,
			equipLevel.CostItem[i].ItemNum) {
			return msg.ErrCode_NO_ENOUGH_ITEM, nil
		}
		// tdaItems = append(tdaItems, &tda.Item{ItemId: strconv.Itoa(int(equipLevel.CostItem[i].ItemId)), ItemNum: equipLevel.CostItem[i].ItemNum})
	}

	var notifyItems []uint32
	for i := 0; i < len(equipLevel.CostItem); i++ {
		CostItem(p.GetUserId(), equipLevel.CostItem[i].ItemId,
			equipLevel.CostItem[i].ItemNum, publicconst.EquipUpgradeCostItem, false)
		notifyItems = append(notifyItems, equipLevel.CostItem[i].ItemId)
	}

	equipPos.Level += 1
	setEquipPosLevelAttr(p, equipPos, jEquip, true)
	p.SaveEquip()

	// 通知客户端
	updateClientItemsChange(p.GetUserId(), notifyItems)

	//event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(playerData, notifyItems, ListenNotifyClientItemEventEvent))

	//
	UpdateTask(p, true, publicconst.TASK_COND_UPGRADE_EQUIP, 1)
	UpdateTask(p, true, publicconst.TASK_COND_ANY_EQUIP_LEVEL)
	UpdateTask(p, true, publicconst.TASK_COND_ALL_EQUIP_POS_LEVEL)

	processHistoryData(p, publicconst.TASK_COND_UPGRADE_EQUIP, 0, 1)

	// tda
	//tda.TdaEquipUpgrade(p.ChannelId, p.TdaCommonAttr, pos, equipPos.Level, jEquip.GenRarity(), tdaItems)

	return msg.ErrCode_SUCC, equipPos
}

// AutoUpgradeEquipPos 自动升级装备部位
func AutoUpgradeEquipPos(p *player.Player, pos uint32) (msg.ErrCode, *model.EquipPos) {
	equipPos := getEquipPos(p, pos)
	if equipPos == nil {
		return msg.ErrCode_EQUIP_POS_NO_EQUIP, nil
	}

	jEquip := template.GetEquipTemplate().GetEquip(equipPos.EquipId)
	if jEquip == nil {
		return msg.ErrCode_EQUIP_NOT_EXIST, nil
	}

	if equipPos.Level >= jEquip.Data.LevelMax {
		return msg.ErrCode_EQUIP_LEVEL_FULL, nil
	}

	var levels []*template.JEquipLevel
	for i := equipPos.Level; i < jEquip.Data.LevelMax; i++ {
		equipLevel := template.GetEquipLevelTemplate().GetEquipLevel(pos, i)
		if equipLevel != nil && len(equipLevel.CostItem) > 0 {
			hasNext := true
			for k := 0; k < len(equipLevel.CostItem); k++ {
				if !EnoughItem(p.GetUserId(),
					equipLevel.CostItem[k].ItemId, equipLevel.CostItem[k].ItemNum) {
					hasNext = false
					break
				}
			}
			if !hasNext {
				break
			} else {
				levels = append(levels, equipLevel)
			}
		}
	}

	if len(levels) == 0 {
		return msg.ErrCode_NO_ENOUGH_ITEM, nil
	}

	var notifyItems []uint32
	//tdaItems := make([]*tda.Item, 0, len(levels))
	for k := 0; k < len(levels); k++ {
		for i := 0; i < len(levels[k].CostItem); i++ {
			CostItem(p.GetUserId(), levels[k].CostItem[i].ItemId,
				levels[k].CostItem[i].ItemNum, publicconst.EquipAutoUpgradeCostItem, false)
			notifyItems = append(notifyItems, levels[k].CostItem[i].ItemId)
			//tdaItems = append(tdaItems, &tda.Item{ItemId: strconv.Itoa(int(levels[k].CostItem[i].ItemId)), ItemNum: levels[k].CostItem[i].ItemNum})
		}
	}

	equipPos.Level = levels[len(levels)-1].Data.Level + 1
	setEquipPosLevelAttr(p, equipPos, jEquip, true)
	p.SaveEquip()

	//event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(playerData, notifyItems, ListenNotifyClientItemEventEvent))

	// 通知客户端
	updateClientItemsChange(p.GetUserId(), notifyItems)

	UpdateTask(p, true, publicconst.TASK_COND_UPGRADE_EQUIP, 1)
	UpdateTask(p, true, publicconst.TASK_COND_ANY_EQUIP_LEVEL)
	UpdateTask(p, true, publicconst.TASK_COND_ALL_EQUIP_POS_LEVEL)

	processHistoryData(p, publicconst.TASK_COND_UPGRADE_EQUIP, 0, 1)

	// tda
	//tda.TdaEquipUpgrade(p.ChannelId, p.TdaCommonAttr, pos, equipPos.Level, jEquip.GenRarity(), tdaItems)

	return msg.ErrCode_SUCC, equipPos
}

// EquipPos 部位安装装备
func EquipPos(p *player.Player, pos, equipId uint32) (msg.ErrCode, uint32, *model.EquipPos) {
	if pos > 10 {
		return msg.ErrCode_INVALID_DATA, 0, nil
	}

	jEquip := template.GetEquipTemplate().GetEquip(equipId)
	if jEquip == nil {
		return msg.ErrCode_EQUIP_NOT_EXIST, 0, nil
	}

	if jEquip.Data.Pos != pos {
		return msg.ErrCode_EQUIP_POS_INVALID, 0, nil
	}

	equip := getEquip(p, equipId)
	if equip == nil {
		return msg.ErrCode_EQUIP_NOT_EXIST, 0, nil
	}

	equipPos := getEquipPos(p, pos)
	if equipPos == nil {
		equipPos = model.NewEquipPos(pos, equipId)
		addEquipPos(p, equipPos, jEquip)
	} else {
		equipPos.EquipId = equipId
		equipPos.UpdateTime = tools.GetCurTime()

		if equipPos.EquipId == 0 {
			log.Error("EquipPos is zero", zap.Uint64("AccountId", p.GetUserId()), zap.Uint32("pos", pos))
			//log.Errorf("EquipPos accountid:%v equipId:%v", playerData.GetUserId(), equipPos.EquipId)
		}
		updateEquipPos(p, equipPos, jEquip, true)
	}

	equipLevel := equipPos.Level
	if equipPos.Level > jEquip.Data.LevelMax {
		equipLevel = jEquip.Data.LevelMax
	}

	UpdateTask(p, true,
		publicconst.TASK_COND_PUT_ON_EQUIP)

	// 手动操作之后取消同步
	cancelSuitSync(p)
	return msg.ErrCode_SUCC, equipLevel, equipPos
}

// cancelSuitSync 取消套装同步
func cancelSuitSync(p *player.Player) {
	if !p.UserData.Equip.SyncEquipPos {
		return
	}
	p.UserData.Equip.SyncEquipPos = false
	p.SaveEquip()
}

// updateEquipPos 更新装备位置
func updateEquipPos(p *player.Player, data *model.EquipPos, jEquip *template.JEquip, notifyClient bool) {
	data.Attr = make(map[uint32]*model.Attr)
	InitAttr(data.Attr, jEquip.InitAttr)
	SetAttrLevelValue(data.Attr, jEquip.GetLevelAttr(data.Level))
	calcEquipAffix(p, data)
	calcEquipPosAttr(p, data)
	p.SaveEquip()

	// 计算战力
	GlobalAttrChange(p, notifyClient)
}

// GemUuidParam 100
func GemUuidParam() uint64 {
	return uint64(msg.GemUuidParam_GemUuidParam_Rate)
}

func GetGemId(uuid uint64) uint32 {
	return uint32(uuid / GemUuidParam() / GemUuidParam())
}

// GetGemAffixType 洗练词条类型 0表示没有洗练 词条属性为原始值
func GetGemAffixType(uuid uint64) uint32 {
	return uint32(uuid % GemUuidParam())
}

// GetGemAffixIdx 洗练词条池索引
func GetGemAffixIdx(uuid uint64) uint32 {
	return uint32((uuid / GemUuidParam()) % GemUuidParam())
}
func GetGemAttrs(p *player.Player) map[uint32]*model.Attr {
	ret := make(map[uint32]*model.Attr)
	for _, pos := range p.UserData.Equip.GemPos {
		for _, uuid := range pos {
			if uuid == 0 {
				continue
			}
			affix_type := GetGemAffixType(uuid)

			if affix_type == 0 {
				j_gem := template.GetGemTemplate().GetGem(GetGemId(uuid))
				if j_gem == nil {
					log.Error("get gem err when GetGemAttrs",
						zap.Uint64("account_id", p.GetUserId()),
						zap.Uint64("uuid", uuid))
					continue
				}
				for _, attr := range j_gem.Attrs {
					if v, ok := ret[attr.Id]; ok {
						v.InitValue += attr.Value
						v.FinalValue += attr.Value
					} else {
						ret[attr.Id] = &model.Attr{
							Id:         attr.Id,
							InitValue:  attr.Value,
							FinalValue: attr.Value,
						}
					}

				}
			} else {
				if affix_type != uint32(msg.GemAffixType_GemAffix_Attr) {
					continue
				}

				j_refresh := template.GetGemRefreshTemplate().GetGemRefresh(GetGemId(uuid))
				if j_refresh == nil {
					log.Error("get gem refresh err when GetGemAttrs",
						zap.Uint64("account_id", p.GetUserId()),
						zap.Uint64("uuid", uuid))
					continue
				}

				affix_idx := GetGemAffixIdx(uuid)
				if affix_idx >= uint32(len(j_refresh.Attrs)) {
					log.Error("affix idx err when GetGemAttrs",
						zap.Uint64("account_id", p.GetUserId()),
						zap.Uint64("uuid", uuid),
						zap.Uint32("affix_type", affix_type),
						zap.Uint32("affix_idx", affix_idx),
					)
					continue
				}

				attr := j_refresh.Attrs[affix_idx]
				if v, ok := ret[attr.Id]; ok {
					v.InitValue += attr.Value
					v.FinalValue += attr.Value
				} else {
					ret[attr.Id] = &model.Attr{
						Id:         attr.Id,
						InitValue:  attr.Value,
						FinalValue: attr.Value,
					}
				}
			}
		}
	}

	return ret
}

// addEquip 添加装备
func addEquip(p *player.Player, equipId, num uint32, save bool) *model.Equip {
	equip := model.NewEquip(equipId, num)
	p.UserData.Equip.EquipData = append(p.UserData.Equip.EquipData, equip)
	if save {
		p.SaveEquip()
	}
	return equip
}
func AddEquip(p *player.Player, equipId, num uint32, save bool) *model.Equip {
	equip := getEquip(p, equipId)
	if equip != nil {
		equip.Num += num
		if save {
			p.SaveEquip()
		}
	} else {
		equip = addEquip(p, equipId, num, save)
	}
	UpdateTask(p, true, publicconst.TASK_COND_GET_EQUIP_NUM, num)
	UpdateTask(p, true, publicconst.TASK_COND_EQUIP_RARITY_NUM)
	processHistoryData(p, publicconst.TASK_COND_GET_EQUIP_NUM, 0, num)
	return equip
}
func GenGemUuid(id, affix_idx, affix_type uint32) uint64 {
	// uuid : 配置id * 100 * 100 + 洗练词条索引 * 100 + 洗练词条类型 (用于同配置id且同词条可堆叠)
	// 洗练词条类型为0(kenum.GemAffixType_None) 表示没有洗练(使用默认词条)
	// e.g. 配置id:123456 洗练词条索引:1 洗练词条类型:2  uuid:1234560102
	return uint64(id)*GemUuidParam()*GemUuidParam() +
		uint64(affix_idx)*GemUuidParam() +
		uint64(affix_type)
}

func AddGem(p *player.Player, gemId, num, affix_idx, affix_type uint32, save bool) *model.GemBagSlot {
	uuid := GenGemUuid(gemId, affix_idx, affix_type)

	gem := getGem(p, uuid)
	if gem != nil {
		gem.Num += num
		if save {
			p.SaveEquip()
		}
		return gem
	} else {
		gem = addGem(p, uuid, num, save)

	}
	return gem
}

func addGem(p *player.Player, uuid uint64, num uint32, save bool) *model.GemBagSlot {
	gem := model.NewGem(uuid, num)
	if save {
		p.SaveEquip()
	}
	p.UserData.Equip.GemBag[uuid] = gem
	return gem
}

func getGem(p *player.Player, uuid uint64) *model.GemBagSlot {
	if v, ok := p.UserData.Equip.GemBag[uuid]; ok {
		return v
	}
	return nil
}

func isAllPosOverLevel(p *player.Player, level uint32) bool {
	for pos := 1; pos <= 6; pos++ {
		if posData := getEquipPos(p, uint32(pos)); posData != nil {
			if posData.Level < level {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// addEquipPos 添加装备位置
func addEquipPos(p *player.Player, data *model.EquipPos, jEquip *template.JEquip) {
	data.Attr = make(map[uint32]*model.Attr)
	InitAttr(data.Attr, jEquip.InitAttr)
	SetAttrLevelValue(data.Attr, jEquip.GetLevelAttr(data.Level))
	calcEquipAffix(p, data)
	calcEquipPosAttr(p, data)

	p.UserData.Equip.EquipPosData = append(p.UserData.Equip.EquipPosData, data)
	p.SaveEquip()
	// 计算战力
	GlobalAttrChange(p, true)
}

func calcEquipAffix(p *player.Player, data *model.EquipPos) {
	j_equip := template.GetEquipTemplate().GetEquip(data.EquipId)
	if j_equip == nil {
		log.Error("j_equip is nil",
			zap.Uint64("AccountId", p.GetUserId()),
			zap.Uint32("equipId", data.EquipId))
		return
	}

	affix_attrs := template.GetEquipAffixTemplate().GetAffixAttrs(j_equip.Data.TypeId, j_equip.Data.Rarity)
	if affix_attrs == nil {
		log.Error("attr_affix is nil",
			zap.Uint64("AccountId", p.GetUserId()),
			zap.Uint32("equipId", data.EquipId))
		return
	}

	affix_skill := template.GetEquipAffixTemplate().GetAffixSkills(j_equip.Data.TypeId, j_equip.Data.Rarity)
	if affix_skill == nil {
		log.Error("skill_affix is nil",
			zap.Uint64("AccountId", p.GetUserId()),
			zap.Uint32("equipId", data.EquipId))
		return
	}

	data.AffixAttr = make(map[uint32]*model.Attr)
	for _, a := range affix_attrs {
		if _, ok := data.AffixAttr[a.Id]; ok {
			data.AffixAttr[a.Id].InitValue += a.Value
			data.AffixAttr[a.Id].FinalValue += a.Value
			continue
		} else {
			data.AffixAttr[a.Id] = &model.Attr{
				Id:         a.Id,
				InitValue:  a.Value,
				FinalValue: a.Value,
			}
		}
	}
	data.AffixSkills = append(data.AffixSkills, affix_skill...)
}

// getEquipPos
func getEquipPos(p *player.Player, pos uint32) *model.EquipPos {
	for i := 0; i < len(p.UserData.Equip.EquipPosData); i++ {
		if p.UserData.Equip.EquipPosData[i].Pos == pos {
			return p.UserData.Equip.EquipPosData[i]
		}
	}
	return nil
}

func getEquipRarityNum(p *player.Player, big, samll uint32) uint32 {
	var ret uint32 = 0
	if p.UserData.Equip == nil {
		return ret
	}
	for i := 0; i < len(p.UserData.Equip.EquipData); i++ {
		equipId := p.UserData.Equip.EquipData[i].Id
		if equipConfig := template.GetEquipTemplate().GetEquip(equipId); equipConfig != nil {
			if equipConfig.HasGtOrEqQuality(big, samll) {
				ret++
			}
		}
	}
	return ret
}

func calcEquipPosAttr(p *player.Player, data *model.EquipPos) {
	for _, a := range data.Attr {
		finalValue := a.InitValue + a.LevelValue + a.Add
		a.SetFinalValue(finalValue)
	}
	for _, a := range data.AffixAttr {
		a.SetFinalValue(a.InitValue)
	}
}

func setEquipPosLevelAttr(p *player.Player, data *model.EquipPos, jEquip *template.JEquip, notifyClient bool) {
	SetAttrLevelValue(data.Attr, jEquip.GetLevelAttr(data.Level))
	calcEquipPosAttr(p, data)

	// 计算战力
	GlobalAttrChange(p, notifyClient)
}

// GmUpgradeEquipPos gm命令升级装备部位
func GmUpgradeEquipPos(p *player.Player, pos, lv uint32) msg.ErrCode {
	equipPos := getEquipPos(p, pos)
	if equipPos == nil {
		return msg.ErrCode_EQUIP_POS_NO_EQUIP
	}

	jEquip := template.GetEquipTemplate().GetEquip(equipPos.EquipId)
	if jEquip == nil {
		return msg.ErrCode_EQUIP_NOT_EXIST
	}

	if lv >= jEquip.Data.LevelMax {
		lv = jEquip.Data.LevelMax
	} else if equipPos.Level <= 0 {
		lv = 1
	}

	equipPos.Level = lv
	setEquipPosLevelAttr(p, equipPos, jEquip, true)
	p.SaveEquip()

	UpdateTask(p, true, publicconst.TASK_COND_UPGRADE_EQUIP, 1)
	UpdateTask(p, true, publicconst.TASK_COND_ANY_EQUIP_LEVEL)
	UpdateTask(p, true, publicconst.TASK_COND_ALL_EQUIP_POS_LEVEL)

	processHistoryData(p, publicconst.TASK_COND_UPGRADE_EQUIP, 0, 1)

	return msg.ErrCode_SUCC
}

// updateClientEquipChange 通知客户端装备变化
func updateClientEquipChange(p *player.Player, equipIds []uint32) {
	res := &msg.NotifyEquipChange{}
	for k := 0; k < len(equipIds); k++ {
		if equip := getEquip(p, equipIds[k]); equip != nil {
			res.Data = append(res.Data, ToProtocolEquip(equip))
		} else {
			res.Data = append(res.Data, &msg.Equip{
				Id:  equipIds[k],
				Num: 0,
			})
		}
	}
	p.SendNotify(res)
}

// getEquipPosMaxLevel 装备部位最高等级
func getEquipPosMaxLevel(p *player.Player) uint32 {
	var maxLevel uint32 = 0
	for i := 0; i < len(p.UserData.Equip.EquipPosData); i++ {
		if p.UserData.Equip.EquipPosData[i].Level > maxLevel {
			maxLevel = p.UserData.Equip.EquipPosData[i].Level
		}
	}
	return maxLevel
}

func NotifyGemChange(p *player.Player, gem *model.GemBagSlot) {
	p.SendNotify(&msg.GemNtf{
		Gems: ToProtocolGems([]*model.GemBagSlot{gem}),
	})
}

func NotifyGemmapChange(p *player.Player, gems map[uint64]*model.GemBagSlot) {
	p.SendNotify(&msg.GemNtf{
		Gems: ToProtocolGemmap(gems),
	})
}
