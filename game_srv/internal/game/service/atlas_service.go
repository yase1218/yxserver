package service

import (
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/errCode"
	"msg"

	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

/*
 * Reward
 *  @Description: 领取图鉴奖励
 *  @param playerData
 *  @param ack
 *  @return error
 */
func AtlasReward(p *player.Player, req *msg.AtlasRewardReq, ack *msg.AtlasRewardAck) error {
	atlasInfo, ok := p.UserData.Atlas.Data[req.GetId()]
	if !ok {
		return errcode.ERR_PARAM
	}

	if atlasInfo.Reward {
		return errCode.ERR_REPEATE_REWARD
	}

	handbookCfg := template.GetHandBookTemplate().GetCfg(req.GetId())
	if handbookCfg == nil {
		log.Error("handbook cfg nil", zap.Uint32("id", req.GetId()))
		return errcode.ERR_CONFIG_NIL
	}

	var notifyItems []uint32
	for _, item := range handbookCfg.Reward {
		addItems := AddItem(p.GetUserId(), item.ItemId, int32(item.ItemNum), publicconst.AtlasActivate, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	atlasInfo.Reward = true
	p.SaveAtlas()

	ack.Id = req.GetId()
	ack.RewardItem = TemplateItemToProtocolItems(handbookCfg.Reward)

	return nil
}

func ActivateAtlasByType(p *player.Player, cfgType uint32, params []uint32) {
	addMap := make(map[uint32]struct{})
	cfgMap := template.GetHandBookTemplate().GetCfgByType(cfgType)
	for _, v := range params {
		if cfg, ok := cfgMap[v]; ok {
			addMap[cfg.Id] = struct{}{}
		}
	}

	if len(addMap) > 0 {
		pbIdSlice := make([]uint32, 0, len(addMap))
		for k := range addMap {
			ActivateAtlas(p, k)
			pbIdSlice = append(pbIdSlice, k)
		}
		p.SendNotify(&msg.AtlasActivateNtf{Id: pbIdSlice})
	}
}

/*
 * ActiveAtlas
 *  @Description: 激活图鉴
 */
func ActivateAtlas(p *player.Player, id uint32) {
	if p.UserData.Atlas.Data[id] != nil {
		return
	}

	p.UserData.Atlas.Data[id] = &model.AtlasUnit{
		Id:     id,
		Reward: false,
	}
	p.SaveAtlas()

	// TODO: 统计
	// ServMgr.GetCommonService().AddStaticsData(p, publicconst.Statics_Add_Atlas, fmt.Sprintf("%id", id))

}
