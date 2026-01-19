package service

import (
	"fmt"
	"kernel/tools"
	"msg"
	"time"

	"github.com/zy/game_data/template"

	"gameserver/internal/common"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
)

const (
	ActivePass_Info   uint32 = 1
	ActivePass_Reward uint32 = 2
)

func createActivityPass(player *player.Player, actId, actType, startTime, endTime uint32) *model.Activity {
	ret := model.NewActivity(actId, actType, startTime, endTime)

	ret.ActDatas = append(ret.ActDatas, model.NewActivityData(ActivePass_Info, 0, 0, 0, 0))
	ret.ActDatas = append(ret.ActDatas, model.NewActivityData(ActivePass_Reward, 0, -1, -1, 0))

	// tda
	//tda.TdaEventOpen(player.TdaCommonAttr, actId, actType)
	return ret
}

// addActivityExp 更新活跃战令经验值
func addActivityExp(p *player.Player, activity *model.Activity, exp uint32) {
	subAct := activity.ActDatas[0]
	subAct.Value1 += exp
	subAct.Value2 = int64(template.GetActivityPassTemplate().GetGradeByExp(activity.ActId, subAct.Value1))

	// 更新活动
	p.SaveAccountActivity()

	NtfActivityChange(p, activity, activity.ActDatas)
}

// getActivePassReward 获取活跃战令奖励
func getActivePassReward(activity *model.Activity, p *player.Player, subActId, num uint32) (msg.ErrCode, []*model.SimpleItem) {
	subAct1 := activity.ActDatas[0]
	subAct2 := activity.ActDatas[1]

	if subAct2.Value2 == subAct1.Value2 {
		return msg.ErrCode_ACTIVITY_HAS_GET_REWARD, nil
	}

	data := template.GetActivityPassTemplate().GetRangeGradeData(activity.ActId, subAct2.Value2, subAct1.Value2)

	// 发放奖励
	rewardItems := make(map[uint32]uint32)
	for i := 0; i < len(data); i++ {
		tools.MergeToMapItem(rewardItems, data[i].FreeReward)
	}

	if subAct1.Value3 == 1 {
		buyData := template.GetActivityPassTemplate().GetRangeGradeData(activity.ActId, subAct2.Value3, subAct1.Value2)
		for i := 0; i < len(buyData); i++ {
			tools.MergeToMapItem(rewardItems, buyData[i].BuyReward)
		}
	}

	var notifyItems []uint32
	for id, num := range rewardItems {
		addItems := AddItem(p.GetUserId(), id, int32(num), publicconst.ActivePassAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	// 更新活动
	subAct2.UpdateTime = tools.GetCurTime()
	subAct2.Value2 = subAct1.Value2
	if subAct1.Value3 == 1 {
		subAct2.Value3 = subAct1.Value2
	}
	p.SaveAccountActivity()

	NtfActivityChange(p, activity, []*model.ActivityData{subAct2})

	var retItems []*model.SimpleItem
	for id, num := range rewardItems {
		retItems = append(retItems, &model.SimpleItem{Id: id, Num: num})
	}

	//report.ReportPassReward(p.ChannelId, p.GetAccountId(), config.Conf.ServerId, activity.ActId, ProtocolSimpleItemsToItems(retItems))

	return msg.ErrCode_SUCC, retItems
}

func getActivePassRedPoint(activity *model.Activity) uint32 {
	subAct1 := activity.ActDatas[0]
	subAct2 := activity.ActDatas[1]

	if subAct2.Value2 < subAct1.Value2 {
		return 1
	}
	return 0
}

func buyPassGrade(p *player.Player, activity *model.Activity, grade uint32) msg.ErrCode {
	maxGrade := template.GetActivityPassTemplate().GetMaxGrade(activity.ActId)
	if uint32(activity.ActDatas[0].Value2) == grade || grade > uint32(maxGrade.Level) {
		return msg.ErrCode_INVALID_DATA
	}

	targetGrade := template.GetActivityPassTemplate().GetGrade(activity.ActId, grade)
	costNum := targetGrade.Exp - activity.ActDatas[0].Value1
	if !EnoughItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_DIAMOND), costNum) {
		return msg.ErrCode_NO_ENOUGH_ITEM
	}

	CostItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_DIAMOND), costNum, publicconst.BuyActivePassGradeCostItem, false)
	updateClientItemsChange(p.GetUserId(), []uint32{uint32(publicconst.ITEM_CODE_DIAMOND)})

	activity.ActDatas[0].Value1 += costNum
	activity.ActDatas[0].Value2 = int64(grade)
	activity.ActDatas[0].UpdateTime = tools.GetCurTime()
	p.SaveAccountActivity()

	NtfActivityChange(p, activity, activity.ActDatas)
	return msg.ErrCode_SUCC
}

// buyPass 购买活跃战令
func buyPass(p *player.Player, passId int) {
	activity := getActivityByType(p, uint32(publicconst.ActivePass))
	activity.ActDatas[0].Value3 = 1
	activity.ActDatas[0].UpdateTime = tools.GetCurTime()

	subAct1 := activity.ActDatas[0]
	subAct2 := activity.ActDatas[1]

	freeData := template.GetActivityPassTemplate().GetRangeGradeData(activity.ActId, subAct2.Value2, subAct1.Value2)
	buyData := template.GetActivityPassTemplate().GetRangeGradeData(activity.ActId, subAct2.Value3, subAct1.Value2)

	// 发放奖励
	rewardItems := make(map[uint32]uint32)
	for i := 0; i < len(freeData); i++ {
		tools.MergeToMapItem(rewardItems, freeData[i].FreeReward)
	}
	for i := 0; i < len(buyData); i++ {
		tools.MergeToMapItem(rewardItems, buyData[i].BuyReward)
	}

	var notifyItems []uint32
	for id, num := range rewardItems {
		addItems := AddItem(p.GetUserId(), id, int32(num), publicconst.ActivePassBuyAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	// 更新已经领取奖励的等级
	subAct2.Value2 = subAct1.Value2
	subAct2.Value3 = subAct1.Value2
	p.SaveAccountActivity()

	NtfActivityChange(p, activity, activity.ActDatas)

	// 通知奖励获得
	var retItems []*model.SimpleItem
	for id, num := range rewardItems {
		retItems = append(retItems, &model.SimpleItem{Id: id, Num: num})
	}
	res2 := &msg.NotifyRewardItem{}
	res2.GetItems = ToProtocolSimpleItems(retItems)
	p.SendNotify(res2)
}

// deleteActivePass 删除过期的活动
func deleteActivePass(p *player.Player, activity *model.Activity) {
	subAct1 := activity.ActDatas[0]
	subAct2 := activity.ActDatas[1]
	data := template.GetActivityPassTemplate().GetRangeGradeData(activity.ActId, subAct2.Value2, subAct1.Value2)
	// 发放奖励
	rewardItems := make(map[uint32]uint32)
	for i := 0; i < len(data); i++ {
		tools.MergeToMapItem(rewardItems, data[i].FreeReward)
		if subAct1.Value3 == 1 {
			tools.MergeToMapItem(rewardItems, data[i].BuyReward)
		}
	}
	if len(rewardItems) == 0 {
		return
	}

	var items []*model.SimpleItem
	for id, num := range rewardItems {
		items = append(items, &model.SimpleItem{Id: id, Num: num})
	}

	// 补发邮件
	mailConfig := template.GetMailTemplate().GetMail(template.GetSystemItemTemplate().ActivePassMailId)

	endTime := time.Now().AddDate(0, 0, 60)
	mail := model.NewMail(common.GenSnowFlake(), fmt.Sprintf("%v", mailConfig.Title),
		fmt.Sprintf("%v", mailConfig.Content), items, uint32(endTime.Unix()))
	mail.MailType = 1
	AddSystemMail(p, mail)
}
