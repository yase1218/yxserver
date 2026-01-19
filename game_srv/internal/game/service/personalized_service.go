package service

import (
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"github.com/zy/game_data/template"
	"kernel/tools"
	"msg"
	"sort"
)

// 获取可用推送
func GetPersonalizedItems(p *player.Player, personalizedType, itemId int) (msg.ErrCode, []*msg.Personalized) {
	var msgItems []*msg.Personalized
	var items []*model.PersonalizedItem
	configs := template.GetPersonalizedTemplate().GetPersonalizedConfigs(personalizedType, itemId)
	var dbSave bool
	for i := 0; i < len(configs); i++ {
		if _, ok := p.UserData.Personalized.ItemsMap[configs[i].TypeID]; !ok {
			p.UserData.Personalized.ItemsMap[configs[i].TypeID] = model.NewPersonalizedItem(configs[i])
			dbSave = true
		}
		item := p.UserData.Personalized.ItemsMap[configs[i].TypeID]
		if p.UserData.Personalized.CanPop(item, configs[i]) { // 是否可推送
			if configs[i].PersonalDay > 0 { // 活动相关，活动必须在进行中
				if p.UserData.AccountActivity.IsActivityInProgress(uint32(configs[i].PersonalDay)) {
					items = append(items, item)
					dbSave = true
				}
			} else {
				items = append(items, item)
				dbSave = true
			}
		}
	}
	if dbSave {
		p.SavePersonalized()
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreateTime < items[j].CreateTime
	})
	for i := 0; i < len(items); i++ {
		msgItems = append(msgItems, items[i].ToProto())
	}
	return msg.ErrCode_SUCC, msgItems
}

func ClearOutPut(p *player.Player, chargeId int, notify bool) {
	if shopCfg := template.GetShopTemplate().GetChargeShop(uint32(chargeId)); shopCfg != nil {
		if personalizeCfg := template.GetPersonalizedTemplate().GetPersonalizedCfgByShopID(int(shopCfg.Id)); personalizeCfg != nil {
			if personalizeItem := p.UserData.Personalized.GetItem(personalizeCfg.TypeID); personalizeItem != nil {
				personalizeItem.OutTime = 0
				p.SavePersonalized()
				if notify {
					p.SendNotify(&msg.NotifyClosePersonalized{PersonalizedId: int32(personalizeItem.Id)})
				}
			}
		}
	}
}

func GetUnOutTimePersonalizedItems(p *player.Player) []*msg.Personalized {
	var msgItems []*msg.Personalized
	var items []*model.PersonalizedItem
	curTime := tools.GetCurTime()
	for _, item := range p.UserData.Personalized.ItemsMap {
		if item.OutTime > curTime {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreateTime < items[j].CreateTime
	})
	for i := 0; i < len(items); i++ {
		msgItems = append(msgItems, items[i].ToProto())
	}
	return msgItems
}
