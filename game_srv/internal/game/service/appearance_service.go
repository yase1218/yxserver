package service

import (
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"msg"

	"github.com/zy/game_data/template"
)

func getAppearance(p *player.Player, id uint32) *model.Appearance {
	for i := 0; i < len(p.UserData.Appearance.Appearances); i++ {
		if p.UserData.Appearance.Appearances[i].Id == id {
			return p.UserData.Appearance.Appearances[i]
		}
	}
	return nil
}

// ActiveAppearance 激活外观
func ActiveAppearance(p *player.Player, id uint32) msg.ErrCode {
	config := template.GetAppearanceTemplate().GetAppearance(id)
	if config == nil {
		return msg.ErrCode_APPEARANCE_NOT_EXIST
	}

	if appearance := getAppearance(p, id); appearance != nil {
		return msg.ErrCode_APPEARANCE_EXIST
	}

	for i := 0; i < len(config.CostItems); i++ {
		if !EnoughItem(p.GetUserId(),
			config.CostItems[i].ItemId, config.CostItems[i].ItemNum) {
			return msg.ErrCode_NO_ENOUGH_ITEM
		}
	}

	var notifyClientItems []uint32
	for i := 0; i < len(config.CostItems); i++ {
		CostItem(p.GetUserId(),
			config.CostItems[i].ItemId,
			config.CostItems[i].ItemNum,
			publicconst.ActiveAppearance,
			false)
		notifyClientItems = append(notifyClientItems, config.CostItems[i].ItemId)
	}
	updateClientItemsChange(p.GetUserId(), notifyClientItems)
	AddAppearances(p, []*model.Appearance{model.NewAppearance(id)}, true)
	return msg.ErrCode_SUCC
}

// UseAppearance 使用外观
func UseAppearance(p *player.Player, id uint32) msg.ErrCode {
	appearance := getAppearance(p, id)
	if appearance == nil {
		return msg.ErrCode_APPEARANCE_NOT_EXIST
	}

	config := template.GetAppearanceTemplate().GetAppearance(id)
	if config == nil {
		return msg.ErrCode_APPEARANCE_NOT_EXIST
	}

	if config.Type == 1 {
		if p.UserData.HeadImg == id {
			return msg.ErrCode_INVALID_DATA
		}
		p.UserData.HeadImg = id
		p.SaveHeadImg()
		NotifyAccountChange(p)

		//common.PlayerMgr.UpdatePlayerBasic(&model.AccBasic{AccountId: p.GetUserId(), HeadImg: id})
	} else if config.Type == 2 {
		if p.UserData.HeadFrame == id {
			return msg.ErrCode_INVALID_DATA
		}
		p.UserData.HeadFrame = id
		p.SaveHeadFrame()
		NotifyAccountChange(p)

		//common.PlayerMgr.UpdatePlayerBasic(&model.AccBasic{AccountId: p.GetUserId(), HeadFrame: id})
	} else if config.Type == 3 {
		if p.UserData.Title == id {
			return msg.ErrCode_INVALID_DATA
		}
		p.UserData.Title = id
		p.SaveTitle()
		NotifyAccountChange(p)

		//common.PlayerMgr.UpdatePlayerBasic(&model.AccBasic{AccountId: p.GetUserId(), Title: id})
	}

	attrs := make(map[uint32]*model.Attr)
	var ids []uint32
	ids = append(ids, p.UserData.HeadImg)
	ids = append(ids, p.UserData.HeadFrame)
	ids = append(ids, p.UserData.Title)

	for i := 0; i < len(ids); i++ {
		if config := template.GetAppearanceTemplate().GetAppearance(ids[i]); config != nil {
			for m := 0; m < len(config.Attr); m++ {
				if tmp, ok := attrs[config.Attr[m].Id]; ok {
					tmp.AddValue(config.Attr[m].Value)
				} else {
					d := model.NewAttr(config.Attr[m].Id, 0)
					d.Add = config.Attr[m].Value
					attrs[config.Attr[m].Id] = d
				}
			}
		}
	}

	update := false
	for id, data := range attrs {
		if tmp, ok := p.UserData.Appearance.Attrs[id]; ok {
			if tmp.Add != data.Add {
				update = true
				break
			}
		} else {
			update = true
			break
		}
	}

	if update {
		p.UserData.Appearance.Attrs = attrs
		calcAppearanceAttr(p)
		p.SaveAppearance()
		// 计算全局属性
		GlobalAttrChange(p, true)
	}

	return msg.ErrCode_SUCC
}

func calcAppearanceAttr(p *player.Player) {
	for _, data := range p.UserData.Appearance.Attrs {
		var finalValue float32
		finalValue = data.InitValue + data.LevelValue + data.Add
		data.SetFinalValue(finalValue)
	}
}

// AddAppearances 添加外观
func AddAppearances(p *player.Player, data []*model.Appearance, notifyClient bool) {
	p.UserData.Appearance.Appearances = append(p.UserData.Appearance.Appearances, data...)
	//dao.AppearanceDao.AddAppearances(p.GetUserId(), data)
	p.SaveAppearance()

	if notifyClient {
		p.SendNotify(&msg.NotifyAppearanceChange{
			Data: ToProtocolAppearances(data),
		})
	}
}

func AddShipAppearances(p *player.Player, shipId uint32) {
	appearance := template.GetAppearanceTemplate().GetShipAppearance(shipId)
	if len(appearance) == 0 {
		return
	}

	var data []*model.Appearance
	for i := 0; i < len(appearance); i++ {
		data = append(data, model.NewAppearance(appearance[i].Id))
	}
	AddAppearances(p, data, true)
}

func ToProtocolAppearance(data *model.Appearance) *msg.AppearanceData {
	return &msg.AppearanceData{
		Id: data.Id,
	}
}

func ToProtocolAppearances(data []*model.Appearance) []*msg.AppearanceData {
	var ret []*msg.AppearanceData
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolAppearance(data[i]))
	}
	return ret
}
