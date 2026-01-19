package service

import (
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tools"
	"msg"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

type AdFunc func(*player.Player, *template.JAd, uint32)

var (
	adFuncMap map[publicconst.AdType]AdFunc
)

func init() {
	adFuncMap = make(map[publicconst.AdType]AdFunc)
	adFuncMap[publicconst.AD_Battle_Life] = battleCommonAd
	adFuncMap[publicconst.Ad_Battle_Select_Weapon] = battleCommonAd
	adFuncMap[publicconst.Ad_Battle_Settle_Chip] = battleCommonAd
	adFuncMap[publicconst.Ad_Battle_Part_Shop] = battleCommonAd
	adFuncMap[publicconst.Ad_Battle_Construction_Lottery] = battleCommonAd
	adFuncMap[publicconst.Ad_Add_Ap] = buyApAd
	adFuncMap[publicconst.Ad_Add_OnHook] = getQuickHookRewardAd
	adFuncMap[publicconst.Ad_Coin_Sweep_Timess] = addSweepTimesAd
	adFuncMap[publicconst.Ad_Equip_Sweep_Times] = addSweepTimesAd
	adFuncMap[publicconst.Ad_Weapon_Sweep_times] = addSweepTimesAd
	adFuncMap[publicconst.Ad_Add_Free_Lottery] = addFreeLotteryAd
	adFuncMap[publicconst.Ad_Shop_Buy] = buyShopItemAd
	adFuncMap[publicconst.Ad_Battle_Main_Settle] = mainBattleSettleAd
}

// mainBattleSettle 主线关卡结算
func mainBattleSettleAd(p *player.Player, adConfig *template.JAd, para uint32) {
	adInfo := getCurAd(p)
	if adInfo == nil {
		log.Error("adinfo nil", zap.Uint64("accountid", p.GetUserId()))
		return
	}
	adInfo.Times += 1
	adInfo.Flag = false
	p.SaveBaseInfo()

	var notifyClientItem []uint32
	var finalItems []*model.SimpleItem
	for _, v := range p.MissionAdRewardItems {
		temp := AddItem(p.GetUserId(),
			v.Id,
			int32(v.Num),
			publicconst.AdMissionAddItem, false)
		notifyClientItem = append(notifyClientItem, GetSimpleItemIds(temp)...)
		finalItems = append(finalItems, &model.SimpleItem{
			Id:  v.Id,
			Num: uint32(v.Num),
		})
	}
	updateClientItemsChange(p.GetUserId(), notifyClientItem)

	p.MissionAdRewardItems = nil

	p.SendNotify(&msg.NotifyAdInfo{
		Data: &msg.AdInfo{
			AdId:  adInfo.AdId,
			Times: adInfo.Times,
		},
	})
}

// buyShopItem 购买商城商品
func buyShopItemAd(p *player.Player, adConfig *template.JAd, para uint32) {
	adInfo := getCurAd(p)
	if adInfo == nil {
		return
	}
	adInfo.Times += 1
	adInfo.Flag = false
	p.SaveBaseInfo()

	itemConfig := template.GetShopTemplate().GetShopItem(adInfo.Para)
	if itemConfig.CostAd != adInfo.AdId {
		return
	}

	log.Debug("AdService buyShopItem, ",
		zap.Any("accountId", p.GetUserId()), zap.Any("adid", adInfo.AdId), zap.Any("Para", adInfo.Para))
	BuyShopItemCallBack(p, itemConfig)

	p.SendNotify(&msg.NotifyAdInfo{
		Data: &msg.AdInfo{
			AdId:  adInfo.AdId,
			Times: adInfo.Times,
		},
	})
}

// addFreeLottery 添加免费抽奖
func addFreeLotteryAd(p *player.Player, adConfig *template.JAd, para uint32) {
	adInfo := getCurAd(p)
	if adInfo == nil {
		return
	}
	adInfo.Times += 1
	adInfo.Flag = false
	p.SaveBaseInfo()

	if cardInfo := getCardPool(p, adInfo.Para); cardInfo != nil {
		cardInfo.FreeTimes += 1
		cardInfo.LotteryTotalTimes += 1
		p.SaveCardPool()
	}

	// 通知广告变化
	p.SendNotify(&msg.NotifyAdInfo{
		Data: &msg.AdInfo{
			AdId:  adInfo.AdId,
			Times: adInfo.Times,
		},
	})
}

// addSweepTimes 增加扫荡次数
func addSweepTimesAd(p *player.Player, adConfig *template.JAd, para uint32) {
	adInfo := getCurAd(p)
	if adInfo == nil {
		log.Error("adinfo nil", zap.Uint64("accountid", p.GetUserId()))
		return
	}
	adInfo.Times += 1
	adInfo.Flag = false
	p.SaveBaseInfo()

	// 增加玩法扫荡次数
	AddPlayMethodTimes(p, int(adInfo.Para), 1)

	// 通知广告变化
	p.SendNotify(&msg.NotifyAdInfo{
		Data: &msg.AdInfo{
			AdId:  adInfo.AdId,
			Times: adInfo.Times,
		},
	})
}

// getQuickHookReward 获得快速挂机奖励
func getQuickHookRewardAd(p *player.Player, adConfig *template.JAd, para uint32) {
	adInfo := getCurAd(p)
	if adInfo == nil {
		return
	}
	adInfo.Times += 1
	adInfo.Flag = false
	p.SaveBaseInfo()

	addItems := calcQuickOnHookReward(p, false)
	var notifyClientItem []uint32
	var finalItems []*model.SimpleItem
	for i := 0; i < len(addItems); i++ {
		temp := AddItem(p.GetUserId(), addItems[i].Id, int32(addItems[i].Num), publicconst.QuickOnHookAddItem, false)
		notifyClientItem = append(notifyClientItem, GetSimpleItemIds(temp)...)
		//	finalItems = append(finalItems, temp...)

		finalItems = append(finalItems, &model.SimpleItem{
			Id:  addItems[i].Id,
			Num: uint32(addItems[i].Num),
		})
	}
	updateClientItemsChange(p.GetUserId(), notifyClientItem)

	// 增加购买次数
	p.UserData.BaseInfo.QuickOnHookData.BuyTimes += 1
	p.SaveBaseInfo()

	// 通知奖励获得
	p.SendNotify(&msg.NotifyRewardItem{
		GetItems: ToProtocolSimpleItems(finalItems),
	})

	// 通知广告变化
	p.SendNotify(&msg.NotifyAdInfo{
		Data: &msg.AdInfo{
			AdId:  adInfo.AdId,
			Times: adInfo.Times,
		},
	})
}

func battleCommonAd(p *player.Player, adConfig *template.JAd, para uint32) {
	if p.UserData.BaseInfo.MissData != nil {
		var adInfo *model.AdInfo
		for i := 0; i < len(p.UserData.BaseInfo.MissData.Ads); i++ {
			if p.UserData.BaseInfo.MissData.Ads[i].Flag {
				adInfo = p.UserData.BaseInfo.MissData.Ads[i]
				break
			}
		}
		adInfo.Times += 1
		adInfo.Flag = false
		p.SaveBaseInfo()

		p.SendNotify(&msg.NotifyAdInfo{
			Data: &msg.AdInfo{
				AdId:  adInfo.AdId,
				Times: adInfo.Times,
			},
		})
	}
}

func buyApAd(p *player.Player, adConfig *template.JAd, para uint32) {
	adInfo := getCurAd(p)
	if adInfo == nil {
		log.Error("adinfo nil", zap.Uint64("accountid", p.GetUserId()))
		return
	}
	adInfo.Times += 1
	adInfo.Flag = false
	p.SaveBaseInfo()

	// 通知客户端
	var itemIds []uint32
	var finalItems []*model.SimpleItem
	for i := 0; i < len(adConfig.RewardItems); i++ {
		addItems := AddItem(p.GetUserId(),
			adConfig.RewardItems[i].ItemId,
			int32(adConfig.RewardItems[i].ItemNum),
			publicconst.AdAddItem, false)
		itemIds = append(itemIds, adConfig.RewardItems[i].ItemId)
		finalItems = append(finalItems, addItems...)
	}

	// 通知客户端
	updateClientItemsChange(p.GetUserId(), itemIds)

	// 通知奖励获得
	p.SendNotify(&msg.NotifyRewardItem{
		GetItems: ToProtocolSimpleItems(finalItems),
	})

	// 通知广告变化
	p.SendNotify(&msg.NotifyAdInfo{
		Data: &msg.AdInfo{
			AdId:  adInfo.AdId,
			Times: adInfo.Times,
		}})
}

// getCurAd 获得当前广告
func getCurAd(p *player.Player) *model.AdInfo {
	for i := 0; i < len(p.UserData.BaseInfo.Ad); i++ {
		if p.UserData.BaseInfo.Ad[i].Flag {
			return p.UserData.BaseInfo.Ad[i]
		}
	}
	return nil
}

// StartAd 开始广告
func StartAd(p *player.Player, adId, para uint32) msg.ErrCode {
	adConfig := template.GetAdTemplate().GetAd(adId)
	if adConfig == nil {
		return msg.ErrCode_AD_NOT_EXIST
	}

	adType := publicconst.AdType(adConfig.AdType)
	// 战内
	if adType == publicconst.AD_Battle_Life || adType == publicconst.Ad_Battle_Select_Weapon ||
		adType == publicconst.Ad_Battle_Settle_Chip || adType == publicconst.Ad_Battle_Part_Shop ||
		adType == publicconst.Ad_Battle_Construction_Lottery {
		if p.UserData.BaseInfo.MissData == nil || p.UserData.BaseInfo.MissData.MissionId == 0 {
			return msg.ErrCode_AD_NOT_IN_BATTLE
		}

		var adInfo *model.AdInfo
		for i := 0; i < len(p.UserData.BaseInfo.MissData.Ads); i++ {
			if p.UserData.BaseInfo.MissData.Ads[i].AdId == adId {
				adInfo = p.UserData.BaseInfo.MissData.Ads[i]
				break
			}
		}
		if adInfo == nil {
			adInfo = model.NewAdInfo(adId, para, 0)
			p.UserData.BaseInfo.MissData.Ads = append(p.UserData.BaseInfo.MissData.Ads, adInfo)
		}

		if adInfo.Times >= adConfig.MaxTimes {
			log.Error("ad time too many", zap.Uint32("adId", adId), zap.Uint32("times", adInfo.Times),
				zap.Uint32("cfgTimes", adConfig.MaxTimes))
			return msg.ErrCode_AD_TIMES_FULL
		}
		adInfo.Para = para
		adInfo.Flag = true
	} else {
		var adInfo *model.AdInfo
		for i := 0; i < len(p.UserData.BaseInfo.Ad); i++ {
			if p.UserData.BaseInfo.Ad[i].AdId == adId {
				adInfo = p.UserData.BaseInfo.Ad[i]
				break
			}
		}

		if adInfo != nil && adInfo.Times >= adConfig.MaxTimes {
			log.Error("ad time too many", zap.Uint32("adId", adId), zap.Uint32("times", adInfo.Times),
				zap.Uint32("cfgTimes", adConfig.MaxTimes))
			return msg.ErrCode_AD_TIMES_FULL
		}
		if adInfo == nil {
			adInfo = model.NewAdInfo(adId, para, tools.GetDailyRefreshTime())
			p.UserData.BaseInfo.Ad = append(p.UserData.BaseInfo.Ad, adInfo)
		} else {
			adInfo.Flag = true
			adInfo.Para = para
		}
		p.SaveBaseInfo()
	}
	return msg.ErrCode_SUCC
}

// ResetAdData 重置广告数据
func ResetAdData(p *player.Player) {
	curTime := tools.GetCurTime()
	update := false
	for i := 0; i < len(p.UserData.BaseInfo.Ad); i++ {
		if curTime >= p.UserData.BaseInfo.Ad[i].NextResetTime {
			p.UserData.BaseInfo.Ad[i].Times = 0
			p.UserData.BaseInfo.Ad[i].Para = 0
			p.UserData.BaseInfo.Ad[i].Flag = false
			p.UserData.BaseInfo.Ad[i].NextResetTime = tools.GetDailyRefreshTime()
			update = true
		}
	}
	if update {
		p.SaveBaseInfo()
	}
}

// AdCallBack 广告回调
func AdCallBack(p *player.Player) {
	var adId uint32 = 0
	var para uint32 = 0
	if p.UserData.BaseInfo.MissData != nil {
		for i := 0; i < len(p.UserData.BaseInfo.MissData.Ads); i++ {
			if p.UserData.BaseInfo.MissData.Ads[i].Flag {
				adId = p.UserData.BaseInfo.MissData.Ads[i].AdId
				para = p.UserData.BaseInfo.MissData.Ads[i].Para
				break
			}
		}
	}
	if adId == 0 {
		for i := 0; i < len(p.UserData.BaseInfo.Ad); i++ {
			if p.UserData.BaseInfo.Ad[i].Flag {
				adId = p.UserData.BaseInfo.Ad[i].AdId
				para = p.UserData.BaseInfo.Ad[i].Para
				break
			}
		}
	}

	adConfig := template.GetAdTemplate().GetAd(adId)
	if adConfig == nil {
		return
	}
	if f, ok := adFuncMap[publicconst.AdType(adConfig.AdType)]; ok {
		f(p, adConfig, para)
	}
}
