package service

import (
	"fmt"
	"gameserver/internal/common"
	"gameserver/internal/config"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/report"
	"kernel/tools"
	"msg"
	"time"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

const (
	Desert_Sign   uint32 = 1 // 签到 value1 是否购买 value2 签到天数  value3 下一次签到时间 value4 补签天数
	Desert_Battle uint32 = 2 //
	Desert_Info   uint32 = 3 // 沙漠活动信息 value1 神偷次数
)

const (
	TaskPass_Info   uint32 = 1
	TaskPass_Reward uint32 = 2
)

const (
	OpenServer_Info uint32 = 1
)

type ActivityFunc func(player *player.Player, configData []*template.JActivityConfig, notfiyClient bool)
type updateActivityFunc func(playerData *player.Player, activity *model.Activity, subActId uint32, args ...uint32)
type getActivityRewardFunc func(activity *model.Activity, playerData *player.Player, subActId, num uint32) (msg.ErrCode, []*model.SimpleItem)
type getActivityRedPointFunc func(activity *model.Activity) uint32
type refreshActivityDataFunc func(playerData *player.Player, activity *model.Activity)
type createActivityFunc func(player *player.Player, actId, actType, startTime, endTime uint32) *model.Activity
type deleteActivityFunc func(player *player.Player, activity *model.Activity)

var (
	activityFuncMap            map[publicconst.ActivityType]ActivityFunc
	updateActivityMap          map[publicconst.ActivityType]updateActivityFunc
	getActivityRewardMap       map[publicconst.ActivityType]getActivityRewardFunc
	getActivityRedPointFuncMap map[publicconst.ActivityType]getActivityRedPointFunc
	refreshActivityDataFuncMap map[publicconst.ActivityType]refreshActivityDataFunc
	createActivityMap          map[publicconst.ActivityType]createActivityFunc
	deleteActivityMap          map[publicconst.ActivityType]deleteActivityFunc
)

func init() {
	activityFuncMap = make(map[publicconst.ActivityType]ActivityFunc)
	activityFuncMap[publicconst.CardPoolActivity] = processCardPoolActivity
	activityFuncMap[publicconst.LoginActivity] = processCommonActivity
	activityFuncMap[publicconst.ActivePass] = processCommonActivity
	activityFuncMap[publicconst.TaskPass] = processCommonActivity
	activityFuncMap[publicconst.OpenServer] = processCommonActivity
	activityFuncMap[publicconst.Desert] = processCommonActivity

	// 更新活动进度
	updateActivityMap = make(map[publicconst.ActivityType]updateActivityFunc)
	updateActivityMap[publicconst.TaskPass] = updateActDataFunc
	updateActivityMap[publicconst.OpenServer] = updateActDataFunc
	updateActivityMap[publicconst.Desert] = updateActDataFunc

	// 获得活动奖励
	getActivityRewardMap = make(map[publicconst.ActivityType]getActivityRewardFunc)
	getActivityRewardMap[publicconst.LoginActivity] = getLoginReward
	getActivityRewardMap[publicconst.ActivePass] = getActivePassReward
	getActivityRewardMap[publicconst.TaskPass] = getTaskPassReward
	getActivityRewardMap[publicconst.OpenServer] = getOpenServerReward
	getActivityRewardMap[publicconst.Desert] = getDesertReward

	getActivityRedPointFuncMap = make(map[publicconst.ActivityType]getActivityRedPointFunc)
	getActivityRedPointFuncMap[publicconst.LoginActivity] = getLoginRedPoint
	getActivityRedPointFuncMap[publicconst.ActivePass] = getActivePassRedPoint
	getActivityRedPointFuncMap[publicconst.TaskPass] = getTaskPassRedPoint
	getActivityRedPointFuncMap[publicconst.OpenServer] = getOpenServerRedPoint
	getActivityRedPointFuncMap[publicconst.Desert] = getDesertRedPoint

	refreshActivityDataFuncMap = make(map[publicconst.ActivityType]refreshActivityDataFunc)
	refreshActivityDataFuncMap[publicconst.LoginActivity] = refreshLoginData
	refreshActivityDataFuncMap[publicconst.TaskPass] = refreshTaskData
	refreshActivityDataFuncMap[publicconst.Desert] = refreshDesertData

	createActivityMap = make(map[publicconst.ActivityType]createActivityFunc)
	createActivityMap[publicconst.LoginActivity] = createLoginActivity
	createActivityMap[publicconst.ActivePass] = createActivityPass
	createActivityMap[publicconst.TaskPass] = createTaskPass
	createActivityMap[publicconst.OpenServer] = createOpenServer
	createActivityMap[publicconst.Desert] = createDesertActivity

	// 删除活动处理
	deleteActivityMap = make(map[publicconst.ActivityType]deleteActivityFunc)
	deleteActivityMap[publicconst.ActivePass] = deleteActivePass
	deleteActivityMap[publicconst.TaskPass] = deleteTaskPass
	deleteActivityMap[publicconst.OpenServer] = deleteOpenServer
}

// deleteOpenServer
func deleteOpenServer(p *player.Player, activity *model.Activity) {
	subAct1 := activity.ActDatas[0]
	rewardItems := make(map[uint32]uint32)
	// 任务奖励
	for i := 0; i < len(activity.ActDatas); i++ {
		if activity.ActDatas[i].State == uint32(publicconst.TASK_COMPLETE) {
			if jTask := template.GetTaskTemplate().GetTask(activity.ActDatas[i].SubActId); jTask != nil {
				// 更新经验数值
				subAct1.Value1 += jTask.Data.PassExp
				tools.MergeToMapItem(rewardItems, jTask.RewardItems)
			}
		}
	}

	// 等级奖励
	maxLevel := template.GetOpenServerTemplate().GetLevelByExp(subAct1.Value1)
	if uint32(subAct1.Value2) < maxLevel {
		for i := uint32(subAct1.Value2 + 1); i <= maxLevel; i++ {
			if data := template.GetOpenServerTemplate().GetByLevel(i); data != nil {
				tools.MergeToMapItem(rewardItems, data.Reward)
			}
		}
	}

	var items []*model.SimpleItem
	for id, num := range rewardItems {
		items = append(items, &model.SimpleItem{Id: id, Num: num})
	}

	// 补发邮件
	mailConfig := template.GetMailTemplate().GetMail(template.GetSystemItemTemplate().OpenServerMailId)

	endTime := time.Now().AddDate(0, 0, 60)
	mail := model.NewMail(common.GenSnowFlake(), fmt.Sprintf("%v", mailConfig.Title),
		fmt.Sprintf("%v", mailConfig.Content), items, uint32(endTime.Unix()))
	mail.MailType = 1
	AddSystemMail(p, mail)
}

// deleteTaskPass 删除任务战令
func deleteTaskPass(p *player.Player, activity *model.Activity) {
	subAct1 := activity.ActDatas[0]
	subAct2 := activity.ActDatas[1]
	rewardItems := make(map[uint32]uint32)

	// 任务奖励
	for i := 0; i < len(activity.ActDatas); i++ {
		if activity.ActDatas[i].State == uint32(publicconst.TASK_COMPLETE) {
			if jTask := template.GetTaskTemplate().GetTask(activity.ActDatas[i].SubActId); jTask != nil {
				// 更新经验数值
				subAct1.Value1 += jTask.Data.PassExp
				subAct1.Value2 = int64(template.GetTaskPassTemplate().GetGradeByExp(activity.ActId, subAct1.Value1))
				subAct1.UpdateTime = tools.GetCurTime()
				tools.MergeToMapItem(rewardItems, jTask.RewardItems)
			}
		}
	}

	// 等级奖励
	data := template.GetTaskPassTemplate().GetRangeGradeData(activity.ActId, subAct2.Value2, subAct1.Value2)
	for i := 0; i < len(data); i++ {
		tools.MergeToMapItem(rewardItems, data[i].FreeReward)
		if subAct1.Value3 == 1 {
			tools.MergeToMapItem(rewardItems, data[i].BuyReward)
		}
	}

	// 额外奖励
	maxGrade := template.GetTaskPassTemplate().GetMaxGrade(activity.ActId)
	if subAct1.Value1 > maxGrade.Exp {
		scale := (subAct1.Value1 - maxGrade.Exp) / 100
		items := template.GetSystemItemTemplate().GetTaskPassBoxReward(activity.ActId, scale)
		tools.MergeToMapItem(rewardItems, items)
	}

	var items []*model.SimpleItem
	for id, num := range rewardItems {
		items = append(items, &model.SimpleItem{Id: id, Num: num})
	}

	// 补发邮件
	mailConfig := template.GetMailTemplate().GetMail(template.GetSystemItemTemplate().TaskPassMailId)

	endTime := time.Now().AddDate(0, 0, 60)
	mail := model.NewMail(common.GenSnowFlake(), fmt.Sprintf("%v", mailConfig.Title),
		fmt.Sprintf("%v", mailConfig.Content), items, uint32(endTime.Unix()))
	mail.MailType = 1
	AddSystemMail(p, mail)
}

// createDesertActivity 创建沙漠活动
func createDesertActivity(p *player.Player, actId, actType, startTime, endTime uint32) *model.Activity {
	actConfig := template.GetActivityConfigTemplate().GetActivityConfig(actId)
	ret := model.NewActivity(actId, actType, startTime, endTime)
	// 签到活动
	ret.ActDatas = append(ret.ActDatas, model.NewActivityData(Desert_Sign, 0, 0, 0, 0))
	ret.ActDatas = append(ret.ActDatas, model.NewActivityData(Desert_Battle, 0, int64(tools.GetWeeklyRefreshTime(0)), 0, 0))
	ret.ActDatas = append(ret.ActDatas, model.NewActivityData(Desert_Info, 0, 0, 0, 0))

	// 创建卡池
	temp := tools.GetDateStart(time.Unix(int64(startTime), 0))
	lotteryEndTime := temp.AddDate(0, 0, actConfig.Para[0][1]).Unix()

	var addCardPool []*model.CardPool
	for i := 0; i < len(actConfig.Para[1]); i++ {
		if ret := RefreshCardPoolActivity(p,
			uint32(actConfig.Para[1][i]), startTime, endTime, uint32(lotteryEndTime)); ret != nil {
			addCardPool = append(addCardPool, ret)
		}
	}

	// 通知奖池变化
	if len(addCardPool) > 0 {
		notifyMsg := &msg.NotifyCardPoolInfoChange{}
		for i := 0; i < len(addCardPool); i++ {
			notifyMsg.Data = append(notifyMsg.Data, ToProtocolCardPool(addCardPool[i]))
		}
		p.SendNotify(notifyMsg)
	}

	// 初始化任务
	var tasks []*template.JTask
	tasks = append(tasks, template.GetTaskTemplate().GetTaskByType(9)...)
	for i := 0; i < len(tasks); i++ {
		data := model.NewActivityData(tasks[i].Data.Id, 0, 0, 0, 0)
		initTaskValue(p, data, publicconst.TaskCond(tasks[i].Data.TaskCondition))
		ret.ActDatas = append(ret.ActDatas, data)
	}

	// tda
	//tda.TdaEventOpen(player.TdaCommonAttr, actId, actType)
	return ret
}

func createOpenServer(p *player.Player, actId, actType, startTime, endTime uint32) *model.Activity {
	ret := model.NewActivity(actId, actType, startTime, endTime)
	ret.ActDatas = append(ret.ActDatas, model.NewActivityData(OpenServer_Info, 0, 0, 0, 0))

	// 初始化任务
	var tasks []*template.JTask
	tasks = append(tasks, template.GetTaskTemplate().GetTaskByType(8)...)
	for i := 0; i < len(tasks); i++ {
		data := model.NewActivityData(tasks[i].Data.Id, 0, 0, 0, 0)
		initTaskValue(p, data, publicconst.TaskCond(tasks[i].Data.TaskCondition))
		ret.ActDatas = append(ret.ActDatas, data)
	}

	// tda
	//tda.TdaEventOpen(player.TdaCommonAttr, actId, actType)
	return ret
}

func createTaskPass(p *player.Player, actId, actType, startTime, endTime uint32) *model.Activity {
	ret := model.NewActivity(actId, actType, startTime, endTime)

	ret.ActDatas = append(ret.ActDatas, model.NewActivityData(TaskPass_Info, 0, 0, 0, 0))
	ret.ActDatas = append(ret.ActDatas, model.NewActivityData(TaskPass_Reward, 0, -1, -1, 0))

	// 初始化任务
	var tasks []*template.JTask
	tasks = append(tasks, template.GetTaskTemplate().GetTaskByType(5)...)
	tasks = append(tasks, template.GetTaskTemplate().GetTaskByType(6)...)
	tasks = append(tasks, template.GetTaskTemplate().GetTaskByType(7)...)
	for i := 0; i < len(tasks); i++ {
		var nextResetTime uint32 = 0
		if tasks[i].Data.TaskType == 5 {
			nextResetTime = tools.GetHourRefreshTime(template.GetSystemItemTemplate().RefreshHour)
		} else if tasks[i].Data.TaskType == 6 {
			nextResetTime = tools.GetWeeklyRefreshTime(template.GetSystemItemTemplate().RefreshHour)
		}
		data := model.NewActivityData(tasks[i].Data.Id, 0, 0, int64(tasks[i].Data.TaskType), int64(nextResetTime))
		initTaskValue(p, data, publicconst.TaskCond(tasks[i].Data.TaskCondition))
		ret.ActDatas = append(ret.ActDatas, data)
	}

	// tda
	//tda.TdaEventOpen(player.TdaCommonAttr, actId, actType)
	return ret
}

func createLoginActivity(p *player.Player, actId, actType, startTime, endTime uint32) *model.Activity {
	ret := model.NewActivity(actId, actType, startTime, endTime)
	// maxDays := template.GetLoginActivityTemplate().GetLoginMaxDays(actId)
	// 计算创角多少天了
	// days := tools.GetDiffDay(time.Unix(int64(player.AccountInfo.CreateTime), 0), time.Now()) + 1
	// if days > maxDays {
	// 	days = maxDays
	// }

	//新活动开启从第一天开始计算
	days := 1
	ret.ActDatas = append(ret.ActDatas, model.NewActivityData(1, 0, int64(tools.GetHourRefreshTime(0)), int64(days), 0))

	// tda
	//tda.TdaEventOpen(player.TdaCommonAttr, actId, actType)
	return ret
}
func refreshDesertData(p *player.Player, activity *model.Activity) {
	activity.ActDatas[2].Value1 = 0
	p.SaveAccountActivity()

	NtfActivityChange(p, activity, []*model.ActivityData{activity.ActDatas[2]})
}

// buySign 购买签到
func buySign(p *player.Player, passId int) {
	activity := getActivityByType(p, publicconst.Desert)
	if activity == nil {
		return
	}

	// 发放奖励
	rewardItems := make(map[uint32]uint32)
	if activity.ActDatas[0].Value1 == 0 && activity.ActDatas[0].Value2 > 0 {
		var notifyItems []uint32
		for i := 1; i <= int(activity.ActDatas[0].Value2); i++ {
			temp := template.GetSignActivityTemplate().GetSign(activity.ActId, uint32(i))
			for k := 0; k < len(temp.BuyReward); k++ {
				if _, ok := rewardItems[temp.BuyReward[k].ItemId]; ok {
					rewardItems[temp.BuyReward[k].ItemId] += temp.BuyReward[k].ItemNum
				} else {
					rewardItems[temp.BuyReward[k].ItemId] = temp.BuyReward[k].ItemNum
				}

				addItems := AddItem(p.GetUserId(),
					temp.BuyReward[k].ItemId, int32(temp.BuyReward[k].ItemNum),
					publicconst.DesertSignAddItem, false)
				notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
			}
		}
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	activity.ActDatas[0].Value1 = 1
	activity.ActDatas[0].UpdateTime = tools.GetCurTime()
	p.SaveAccountActivity()

	NtfActivityChange(p, activity, activity.ActDatas)

	// 通知奖励获得
	if len(rewardItems) > 0 {
		var retItems []*model.SimpleItem
		for id, num := range rewardItems {
			retItems = append(retItems, &model.SimpleItem{Id: id, Num: num})
		}
		res2 := &msg.NotifyRewardItem{}
		res2.GetItems = ToProtocolSimpleItems(retItems)
		p.SendNotify(res2)
	}
}

func refreshTaskData(p *player.Player, activity *model.Activity) {
	curTime := tools.GetCurTime()
	for i := 0; i < len(activity.ActDatas); i++ {
		if activity.ActDatas[i].Value4 > 0 && curTime >= uint32(activity.ActDatas[i].Value4) {
			activity.ActDatas[i].Value1 = 0
			activity.ActDatas[i].State = 0
			activity.ActDatas[i].UpdateTime = curTime
			if activity.ActDatas[i].Value3 == 5 {
				activity.ActDatas[i].Value4 = int64(tools.GetHourRefreshTime(template.GetSystemItemTemplate().RefreshHour))
			} else if activity.ActDatas[i].Value3 == 6 {
				activity.ActDatas[i].Value4 = int64(tools.GetWeeklyRefreshTime(template.GetSystemItemTemplate().RefreshHour))
			}
			p.SaveAccountActivity()
		}
	}
}
func refreshLoginData(p *player.Player, activity *model.Activity) {
	curTime := tools.GetCurTime()
	if curTime < uint32(activity.ActDatas[0].Value2) {
		return
	}

	maxDays := template.GetLoginActivityTemplate().GetLoginMaxDays(activity.ActId)
	if activity.ActDatas[0].Value1 >= maxDays || activity.ActDatas[0].Value3 >= int64(maxDays) {
		return
	}
	activity.ActDatas[0].Value2 = int64(tools.GetHourRefreshTime(0))
	activity.ActDatas[0].Value3 += 1

	// 更新活动
	p.SaveAccountActivity()
}

func getDesertRedPoint(activity *model.Activity) uint32 {
	if !isValidActivity(activity) {
		return 0
	}

	maxDays := template.GetSignActivityTemplate().GetSignMaxDays(activity.ActId)
	curTime := tools.GetCurTime()
	if uint32(activity.ActDatas[0].Value2) < maxDays && curTime >= uint32(activity.ActDatas[0].Value3) {
		return 1
	}

	for i := 0; i < len(activity.ActDatas); i++ {
		if activity.ActDatas[i].State == uint32(publicconst.TASK_COMPLETE) {
			return 1
		}
	}
	return 0
}

func getOpenServerRedPoint(activity *model.Activity) uint32 {
	subAct1 := activity.ActDatas[0]
	days := tools.GetDiffDay(time.Unix(int64(activity.StartTime), 0), time.Now()) + 1
	for i := 0; i < len(activity.ActDatas); i++ {
		if jTask := template.GetTaskTemplate().GetTask(activity.ActDatas[i].SubActId); jTask != nil {
			if days < jTask.Data.Param {
				continue
			}
		}
		if activity.ActDatas[i].State == uint32(publicconst.TASK_COMPLETE) {
			return 1
		}
	}

	maxLevel := template.GetOpenServerTemplate().GetLevelByExp(subAct1.Value1)
	if uint32(subAct1.Value2) < maxLevel {
		return 1
	}
	return 0
}

func getTaskPassRedPoint(activity *model.Activity) uint32 {
	subAct1 := activity.ActDatas[0]
	subAct2 := activity.ActDatas[1]

	for i := 0; i < len(activity.ActDatas); i++ {
		if activity.ActDatas[i].State == uint32(publicconst.TASK_COMPLETE) {
			return 1
		}
	}

	if subAct2.Value2 < subAct1.Value2 {
		return 1
	}

	maxGrade := template.GetTaskPassTemplate().GetMaxGrade(activity.ActId)
	if subAct1.Value1 > maxGrade.Exp && ((subAct1.Value1-maxGrade.Exp)/100) > 1 {
		return 1
	}
	return 0
}

func getLoginRedPoint(activity *model.Activity) uint32 {
	if !isValidActivity(activity) {
		return 0
	}

	maxDays := template.GetLoginActivityTemplate().GetLoginMaxDays(activity.ActId)
	if activity.ActDatas[0].Value1 >= maxDays {
		return 0
	}

	curTime := tools.GetCurTime()
	if curTime >= uint32(activity.ActDatas[0].Value2) || activity.ActDatas[0].Value1 < uint32(activity.ActDatas[0].Value3) {
		return 1
	}

	return 0
}

// getDesertTaskReward 获取沙漠活动任务奖励
func getDesertTaskReward(activity *model.Activity, p *player.Player, subActId uint32) (msg.ErrCode, []*model.SimpleItem) {
	pos := -1
	for i := 0; i < len(activity.ActDatas); i++ {
		if activity.ActDatas[i].SubActId == subActId {
			pos = i
			break
		}
	}

	if pos == -1 {
		return msg.ErrCode_TASK_NOT_EXIST, nil
	}

	if activity.ActDatas[pos].State == uint32(publicconst.TASK_ACCEPT) {
		return msg.ErrCode_NO_COMPLETE_TASK, nil
	}

	if activity.ActDatas[pos].State == uint32(publicconst.TASK_DONE) {
		return msg.ErrCode_TASK_HAS_GET_REWARD, nil
	}

	var changeAct []*model.ActivityData
	rewardItems := make(map[uint32]uint32)

	if jTask := template.GetTaskTemplate().GetTask(activity.ActDatas[pos].SubActId); jTask != nil {
		tools.MergeToMapItem(rewardItems, jTask.RewardItems)
		activity.ActDatas[pos].State = uint32(publicconst.TASK_DONE)
		p.SaveAccountActivity()
		changeAct = append(changeAct, activity.ActDatas[pos])
	}

	var notifyItems []uint32
	for id, num := range rewardItems {
		addItems := AddItem(p.GetUserId(), id, int32(num), publicconst.OpenServerAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	NtfActivityChange(p, activity, changeAct)

	var retItems []*model.SimpleItem
	for id, num := range rewardItems {
		retItems = append(retItems, &model.SimpleItem{Id: id, Num: num})
	}

	// tda

	return msg.ErrCode_SUCC, retItems
}

// getDesertSignReward 获取沙漠活动签到奖励
func getDesertSignReward(activity *model.Activity, p *player.Player) (msg.ErrCode, []*model.SimpleItem) {
	curTime := tools.GetCurTime()
	maxDays := template.GetSignActivityTemplate().GetSignMaxDays(activity.ActId)
	if uint32(activity.ActDatas[0].Value2) >= maxDays {
		return msg.ErrCode_INVALID_DATA, nil
	}

	startTime := tools.GetDateStart(time.Unix(int64(activity.StartTime), 0))
	endTime := startTime.AddDate(0, 0, int(maxDays))
	if curTime >= uint32(endTime.Unix()) {
		return msg.ErrCode_ACTIVITY_SIGN_OVER_TIME, nil
	}

	if curTime < uint32(activity.ActDatas[0].Value3) {
		return msg.ErrCode_HAS_GET_ACTIVITY_SIGN_REWARD, nil
	}

	// 发放奖励
	rewardItems := make(map[uint32]uint32)
	activity.ActDatas[0].Value2 += 1
	activity.ActDatas[0].Value3 = int64(tools.GetHourRefreshTime(0))
	// 更新活动
	p.SaveAccountActivity()

	signConfig := template.GetSignActivityTemplate().GetSign(activity.ActId, uint32(activity.ActDatas[0].Value2))
	temp := signConfig.GetReward(activity.ActDatas[0].Value1)
	for id, num := range temp {
		if _, ok := rewardItems[id]; ok {
			rewardItems[id] += num
		} else {
			rewardItems[id] = num
		}
	}

	var notifyItems []uint32
	for id, num := range rewardItems {
		addItems := AddItem(p.GetUserId(), id, int32(num), publicconst.DesertSignAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	NtfActivityChange(p, activity, activity.ActDatas)

	var retItems []*model.SimpleItem
	for id, num := range rewardItems {
		retItems = append(retItems, &model.SimpleItem{Id: id, Num: num})
	}

	return msg.ErrCode_SUCC, retItems
}

// getDesertReward 领取沙漠活动奖励
func getDesertReward(activity *model.Activity, p *player.Player, subActId, num uint32) (msg.ErrCode, []*model.SimpleItem) {
	if subActId == Desert_Sign {
		return getDesertSignReward(activity, p)
	}
	return getDesertTaskReward(activity, p, subActId)
}

// getOpenServerActiveReward 获取活跃度奖励
func getOpenServerActiveReward(activity *model.Activity, p *player.Player) (msg.ErrCode, []*model.SimpleItem) {
	subAct1 := activity.ActDatas[0]
	var changeAct []*model.ActivityData
	rewardItems := make(map[uint32]uint32)

	maxLevel := template.GetOpenServerTemplate().GetLevelByExp(subAct1.Value1)
	if uint32(subAct1.Value2) == maxLevel {
		return msg.ErrCode_INVALID_DATA, nil
	}

	for i := uint32(subAct1.Value2 + 1); i <= maxLevel; i++ {
		if data := template.GetOpenServerTemplate().GetByLevel(i); data != nil {
			tools.MergeToMapItem(rewardItems, data.Reward)
		}
	}

	subAct1.Value2 = int64(maxLevel)
	p.SaveAccountActivity()
	changeAct = append(changeAct, subAct1)

	var notifyItems []uint32
	for id, num := range rewardItems {
		addItems := AddItem(p.GetUserId(), id, int32(num), publicconst.OpenServerActiveAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	NtfActivityChange(p, activity, changeAct)

	var retItems []*model.SimpleItem
	for id, num := range rewardItems {
		retItems = append(retItems, &model.SimpleItem{Id: id, Num: num})
	}
	return msg.ErrCode_SUCC, retItems
}
func getOpenServerTaskReward(activity *model.Activity, p *player.Player, subActId uint32) (msg.ErrCode, []*model.SimpleItem) {
	subAct1 := activity.ActDatas[0]
	pos := -1
	for i := 0; i < len(activity.ActDatas); i++ {
		if activity.ActDatas[i].SubActId == subActId {
			pos = i
			break
		}
	}

	if pos == -1 {
		return msg.ErrCode_TASK_NOT_EXIST, nil
	}

	if activity.ActDatas[pos].State == uint32(publicconst.TASK_ACCEPT) {
		return msg.ErrCode_NO_COMPLETE_TASK, nil
	}

	if activity.ActDatas[pos].State == uint32(publicconst.TASK_DONE) {
		return msg.ErrCode_TASK_HAS_GET_REWARD, nil
	}

	var changeAct []*model.ActivityData
	rewardItems := make(map[uint32]uint32)
	jTask := template.GetTaskTemplate().GetTask(activity.ActDatas[pos].SubActId)
	if jTask != nil {
		// 更新经验数值
		subAct1.Value1 += jTask.Data.PassExp
		subAct1.UpdateTime = tools.GetCurTime()

		tools.MergeToMapItem(rewardItems, jTask.RewardItems)

		activity.ActDatas[pos].State = uint32(publicconst.TASK_DONE)
		p.SaveAccountActivity()
		changeAct = append(changeAct, activity.ActDatas[pos])
	}
	changeAct = append(changeAct, subAct1)
	p.SaveAccountActivity()

	var notifyItems []uint32
	for id, num := range rewardItems {
		addItems := AddItem(p.GetUserId(), id, int32(num), publicconst.OpenServerAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	NtfActivityChange(p, activity, changeAct)

	var retItems []*model.SimpleItem
	for id, num := range rewardItems {
		retItems = append(retItems, &model.SimpleItem{Id: id, Num: num})
	}

	return msg.ErrCode_SUCC, retItems
}

// getOpenServerReward 获取开服奖励
func getOpenServerReward(activity *model.Activity, p *player.Player, subActId, num uint32) (msg.ErrCode, []*model.SimpleItem) {
	if subActId == 1 {
		return getOpenServerActiveReward(activity, p)
	}
	return getOpenServerTaskReward(activity, p, subActId)
}

// getGradeReward 领取等级奖励
func getGradeReward(activity *model.Activity, p *player.Player) (msg.ErrCode, []*model.SimpleItem) {
	subAct1 := activity.ActDatas[0]
	subAct2 := activity.ActDatas[1]

	// 发放奖励
	rewardItems := make(map[uint32]uint32)
	maxGrade := template.GetTaskPassTemplate().GetMaxGrade(activity.ActId)

	curGrade := template.GetTaskPassTemplate().GetGrade(activity.ActId, uint32(subAct1.Value2))
	freeData := template.GetTaskPassTemplate().GetRangeGradeData(activity.ActId, subAct2.Value2, subAct1.Value2)
	for i := 0; i < len(freeData); i++ {
		tools.MergeToMapItem(rewardItems, freeData[i].FreeReward)
	}

	if subAct1.Value3 == 1 {
		buyData := template.GetTaskPassTemplate().GetRangeGradeData(activity.ActId, subAct2.Value3, subAct1.Value2)
		for i := 0; i < len(buyData); i++ {
			tools.MergeToMapItem(rewardItems, buyData[i].BuyReward)
		}
	}

	if subAct1.Value2 <= int64(maxGrade.Level) {
		subAct2.UpdateTime = tools.GetCurTime()
		subAct2.Value2 = subAct1.Value2
		if subAct1.Value3 == 1 {
			subAct2.Value3 = subAct1.Value2
		}
		subAct2.Value4 = int64(curGrade.Exp)
		p.SaveAccountActivity()
	}

	if subAct1.Value1 > maxGrade.Exp {
		if scale := (subAct1.Value1 - uint32(subAct2.Value4)) / 100; scale > 0 {
			subAct2.UpdateTime = tools.GetCurTime()
			subAct2.Value4 += int64(100 * scale)
			p.SaveAccountActivity()
			items := template.GetSystemItemTemplate().GetTaskPassBoxReward(activity.ActId, scale)
			tools.MergeToMapItem(rewardItems, items)
		}
	}

	if len(rewardItems) == 0 {
		return msg.ErrCode_ACTIVITY_HAS_GET_REWARD, nil
	}

	var notifyItems []uint32
	for id, num := range rewardItems {
		addItems := AddItem(p.GetUserId(), id, int32(num), publicconst.TaskPassAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	NtfActivityChange(p, activity, []*model.ActivityData{activity.ActDatas[0], activity.ActDatas[1]})

	var retItems []*model.SimpleItem
	for id, num := range rewardItems {
		retItems = append(retItems, &model.SimpleItem{Id: id, Num: num})
	}
	return msg.ErrCode_SUCC, retItems
}

func getTaskReward(activity *model.Activity, p *player.Player) (msg.ErrCode, []*model.SimpleItem) {
	subAct1 := activity.ActDatas[0]
	var changeAct []*model.ActivityData
	rewardItems := make(map[uint32]uint32)
	for i := 0; i < len(activity.ActDatas); i++ {
		if activity.ActDatas[i].State == uint32(publicconst.TASK_COMPLETE) {
			if jTask := template.GetTaskTemplate().GetTask(activity.ActDatas[i].SubActId); jTask != nil {
				// 更新经验数值
				subAct1.Value1 += jTask.Data.PassExp
				subAct1.Value2 = int64(template.GetTaskPassTemplate().GetGradeByExp(activity.ActId, subAct1.Value1))
				subAct1.UpdateTime = tools.GetCurTime()

				tools.MergeToMapItem(rewardItems, jTask.RewardItems)

				activity.ActDatas[i].State = uint32(publicconst.TASK_DONE)
				p.SaveAccountActivity()
				changeAct = append(changeAct, activity.ActDatas[i])
			}
		}
	}
	if len(changeAct) == 0 {
		return msg.ErrCode_ACTIVITY_HAS_GET_REWARD, nil
	}
	changeAct = append(changeAct, subAct1)
	p.SaveAccountActivity()

	var notifyItems []uint32
	for id, num := range rewardItems {
		addItems := AddItem(p.GetUserId(), id, int32(num), publicconst.TaskPassAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	NtfActivityChange(p, activity, changeAct)

	var retItems []*model.SimpleItem
	for id, num := range rewardItems {
		retItems = append(retItems, &model.SimpleItem{Id: id, Num: num})
	}
	return msg.ErrCode_SUCC, retItems
}

// getTaskPassReward 获取任务奖励
func getTaskPassReward(activity *model.Activity, p *player.Player, subActId, num uint32) (msg.ErrCode, []*model.SimpleItem) {
	if subActId == 1 {
		return getGradeReward(activity, p)
	}
	return getTaskReward(activity, p)
}

// getLoginReward 获取登录奖励
func getLoginReward(activity *model.Activity, p *player.Player, subActId, num uint32) (msg.ErrCode, []*model.SimpleItem) {
	if activity.ActDatas[0].Value1 == uint32(activity.ActDatas[0].Value3) {
		return msg.ErrCode_ACTIVITY_HAS_GET_REWARD, nil
	}

	maxDays := template.GetLoginActivityTemplate().GetLoginMaxDays(activity.ActId)
	if activity.ActDatas[0].Value1 >= maxDays {
		return msg.ErrCode_INVALID_DATA, nil
	}

	startDay := activity.ActDatas[0].Value1
	activity.ActDatas[0].Value1 = uint32(activity.ActDatas[0].Value3)
	activity.ActDatas[0].Value2 = int64(tools.GetHourRefreshTime(0))

	// 发放奖励
	rewardItems := make(map[uint32]uint32)
	for i := startDay + 1; i <= activity.ActDatas[0].Value1; i++ {
		temp := template.GetLoginActivityTemplate().GetLoginReward(activity.ActId, i)
		if temp == nil {
			continue
		}

		for k := 0; k < len(temp.Reward); k++ {
			if count, ok := rewardItems[temp.Reward[k].ItemId]; ok {
				rewardItems[temp.Reward[k].ItemId] = count + temp.Reward[k].ItemNum
			} else {
				rewardItems[temp.Reward[k].ItemId] = temp.Reward[k].ItemNum
			}
		}
	}

	var notifyItems []uint32
	for id, num := range rewardItems {
		addItems := AddItem(p.GetUserId(), id, int32(num), publicconst.LoginActivityAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}

	updateClientItemsChange(p.GetUserId(), notifyItems)

	// 更新活动
	p.SaveAccountActivity()

	NtfActivityChange(p, activity, activity.ActDatas)

	// 关闭活动
	if activity.ActDatas[0].Value1 >= maxDays {

		p.UserData.AccountActivity.HisData = append(p.UserData.AccountActivity.HisData, activity.ActId)
		p.SaveAccountActivity()

		if checkAllLoginActivityEnd(p) {
			delAllLoginActivity(p)
		}
	}

	var retItems []*model.SimpleItem
	for id, num := range rewardItems {
		retItems = append(retItems, &model.SimpleItem{Id: id, Num: num})
	}

	report.ReportLoginSignReward(p.ChannelId, p.GetUserId(), config.Conf.ServerId, activity.ActId, ProtocolSimpleItemsToItems(retItems))

	return msg.ErrCode_SUCC, retItems
}

// updateActDataFunc 更新活动进度
func updateActDataFunc(p *player.Player, activity *model.Activity, cond uint32, args ...uint32) {
	var changeAct []*model.ActivityData
	for i := 0; i < len(activity.ActDatas); i++ {
		jTask := template.GetTaskTemplate().GetTask(activity.ActDatas[i].GetTaskId())
		if jTask == nil {
			continue
		}

		if jTask.Data.TaskCondition != uint32(cond) {
			continue
		}

		if activity.ActDatas[i].State == uint32(publicconst.TASK_COMPLETE) ||
			activity.ActDatas[i].State == uint32(publicconst.TASK_DONE) {
			continue
		}

		if refreshTaskValue(p, activity.ActDatas[i], publicconst.TaskCond(cond), jTask.Data.MaxValue, jTask.Data.Effect1, args...) {
			changeAct = append(changeAct, activity.ActDatas[i])
			p.SaveAccountActivity()
		}
	}

	if len(changeAct) == 0 {
		return
	}

	NtfActivityChange(p, activity, changeAct)
}

// processCommonActivity 通用活动
func processCommonActivity(p *player.Player, configData []*template.JActivityConfig, notfiyClient bool) {
	curTime := tools.GetCurTime()

	var addActivity []*model.Activity
	for i := 0; i < len(configData); i++ {
		if hasHisAct(p, configData[i].Id) {
			continue
		}
		activity := getActivity(p, configData[i].Id)
		if activity != nil || !CanUnlock(p, configData[i].Cond) {
			continue
		}

		timeCond := configData[i].GetTimeCond()
		startTime := curTime
		endTime := tools.GetMaxTime()

		if timeCond != nil && timeCond.Id == 3 {
			startTime = timeCond.Values[0]
			endTime = timeCond.Values[1]
		} else if timeCond != nil && (timeCond.Id == 4 || timeCond.Id == 7) {
			temp1 := timeCond.Values[0]
			var temp2 uint32 = 1000
			if len(timeCond.Values) > 1 {
				temp2 = timeCond.Values[1]
			}
			if timeCond.Id == 4 {
				startTime, endTime = GetCreateRoleTimeRange(p, temp1, temp2)
			} else {
				startTime, endTime = GetOpenServerTimeRange(p, temp1, temp2)
			}
		} else if timeCond == nil {
			// 沙漠活动
			if configData[i].ActivityType == uint32(publicconst.Desert) {
				temp := tools.GetDateStart(time.Unix(int64(startTime), 0))
				start := temp.AddDate(0, 0, 0)
				end := temp.AddDate(0, 0, configData[i].Para[0][0])

				startTime = uint32(start.Unix())
				endTime = uint32(end.Unix())
			}
		}

		// 删除上一期过期活动
		deleteOvertimeAct(p, configData[i].ActivityType)
		if f, ok := createActivityMap[publicconst.ActivityType(configData[i].ActivityType)]; ok {
			activity = f(p, configData[i].Id, configData[i].ActivityType, startTime, endTime)
		}
		addActivity = append(addActivity, activity)
		p.UserData.AccountActivity.Activities = append(p.UserData.AccountActivity.Activities, activity)
	}
	if len(addActivity) > 0 {
		p.SaveAccountActivity()
		if notfiyClient {
			for i := 0; i < len(addActivity); i++ {
				NtfActivityChange(p, addActivity[i], addActivity[i].ActDatas)
			}
		}
	}
}

// processCardPoolActivity 抽卡活动
func processCardPoolActivity(p *player.Player, configData []*template.JActivityConfig, notfiyClient bool) {
	curTime := tools.GetCurTime()

	var addCardPool []*model.CardPool
	for i := 0; i < len(configData); i++ {
		// 没有解锁
		if !CanUnlock(p, configData[i].Cond) {
			continue
		}

		timeCond := configData[i].GetTimeCond()
		startTime := curTime
		endTime := tools.GetMaxTime()

		if timeCond != nil && timeCond.Id == 3 {
			startTime = timeCond.Values[0]
			endTime = timeCond.Values[1]
		} else if timeCond != nil && timeCond.Id == 4 {
			temp1 := timeCond.Values[0]
			var temp2 uint32 = 10000
			if len(timeCond.Values) > 1 {
				temp2 = timeCond.Values[1]
			}
			startTime, endTime = GetCreateRoleTimeRange(p, temp1, temp2)
		}

		// 在活动范围内
		if ret := RefreshCardPoolActivity(p, uint32(configData[i].Para[0][0]), startTime, endTime, 0); ret != nil {
			addCardPool = append(addCardPool, ret)
		}
	}

	// 通知奖池变化
	if len(addCardPool) > 0 {
		notifyMsg := &msg.NotifyCardPoolInfoChange{}
		for i := 0; i < len(addCardPool); i++ {
			notifyMsg.Data = append(notifyMsg.Data, ToProtocolCardPool(addCardPool[i]))
		}
		p.SendNotify(notifyMsg)
	}
}

// UpdateActivity 更新活动
func UpdateActivity(p *player.Player, cond uint32, args ...uint32) {
	if p.UserData.AccountActivity == nil {
		log.Error("activity data nil", zap.Uint64("accountId", p.GetUserId()))
		return
	}
	for i := 0; i < len(p.UserData.AccountActivity.Activities); i++ {
		if f, ok := updateActivityMap[publicconst.ActivityType(p.UserData.AccountActivity.Activities[i].ActType)]; ok {
			f(p, p.UserData.AccountActivity.Activities[i], cond, args...)
		}
	}
}

func RefreshActivity(p *player.Player, notfiyClient bool) {
	mapData := template.GetActivityConfigTemplate().GetAllTypeActivityConfig()
	for actType, configData := range mapData {
		if f, ok := activityFuncMap[publicconst.ActivityType(actType)]; ok {
			f(p, configData, notfiyClient)
		}
	}
}

func updateActivitiesData(p *player.Player) {
	for i := 0; i < len(p.UserData.AccountActivity.Activities); i++ {
		activity := p.UserData.AccountActivity.Activities[i]
		updateActivityData(p, activity)
	}
}

func updateActivityData(p *player.Player, activity *model.Activity) {
	if f, ok := refreshActivityDataFuncMap[publicconst.ActivityType(activity.ActType)]; ok {
		if isValidActivity(activity) {
			f(p, activity)
		}
	}
}

// isValidActivity 是否是有效活动
func isValidActivity(data *model.Activity) bool {
	curTime := tools.GetCurTime()
	if curTime >= data.StartTime && curTime <= data.EndTime {
		return true
	}
	return false
}

// getActivity
func getActivity(p *player.Player, actId uint32) *model.Activity {
	for i := 0; i < len(p.UserData.AccountActivity.Activities); i++ {
		ret := p.UserData.AccountActivity.Activities[i]
		if ret.ActId == actId {
			if !isValidActivity(ret) {
				return nil
			}
			return p.UserData.AccountActivity.Activities[i]
		}
	}
	return nil
}

// hasHisAct 是否有历史活动
func hasHisAct(p *player.Player, actId uint32) bool {
	if p.UserData.AccountActivity == nil {
		log.Error("activity data nil", zap.Uint64("accountId", p.GetUserId()))
		return false
	}
	return tools.ListContain(p.UserData.AccountActivity.HisData, actId)
}

// deleteOvertimeAct 删除过期活动
func deleteOvertimeAct(p *player.Player, actType uint32) {
	pos := -1
	for i := 0; i < len(p.UserData.AccountActivity.Activities); i++ {
		ret := p.UserData.AccountActivity.Activities[i]
		if ret.ActType == actType && ret.EndTime < uint32(time.Now().Unix()) {
			pos = i
		}
	}

	if pos >= 0 {
		act := p.UserData.AccountActivity.Activities[pos]
		log.Info("deleteOvertimeAct", zap.Uint64("accountId", p.GetUserId()),
			zap.Uint32("actType", actType), zap.Reflect("activityInfo", act))
		if f, ok := deleteActivityMap[publicconst.ActivityType(act.ActType)]; ok {
			f(p, act)
		}
		p.UserData.AccountActivity.Activities = append(p.UserData.AccountActivity.Activities[0:pos], p.UserData.AccountActivity.Activities[pos+1:]...)
		p.SaveAccountActivity()
	}
}

func NtfActivityChange(p *player.Player, activity *model.Activity, changeAct []*model.ActivityData) {
	if p == nil || activity == nil || len(changeAct) == 0 {
		return
	}

	cur := time.Now().Unix()
	if !(int64(activity.StartTime) <= cur && cur <= int64(activity.EndTime)) {
		return
	}

	res := &msg.NotifyActDataChange{}
	cache := &msg.ActDataChange{}
	cache.ActConfig = &msg.ActConfig{ActId: activity.ActId, StartTime: activity.StartTime, EndTime: activity.EndTime}
	for _, subAct := range changeAct {
		cache.Data = append(cache.Data, ToProtocolActData(subAct))
	}
	res.Acts = append(res.Acts, cache)

	p.SendNotify(res)
}

func ToProtocolActData(data *model.ActivityData) *msg.ActData {
	return &msg.ActData{
		SubActId: data.SubActId,
		Value1:   data.Value1,
		Value2:   data.Value2,
		Value3:   data.Value3,
		Value4:   data.Value4,
		State:    msg.ActState(data.State),
	}
}

func checkAllLoginActivityEnd(p *player.Player) bool {
	allFinish := true

	for _, v := range p.UserData.AccountActivity.Activities {

		if v.ActType != uint32(publicconst.LoginActivity) {
			continue
		}

		maxDays := template.GetLoginActivityTemplate().GetLoginMaxDays(v.ActId)
		if v.ActDatas[0].Value1 < maxDays {
			allFinish = false
			break
		}
	}

	return allFinish
}

func delAllLoginActivity(player *player.Player) bool {
	allFinish := true

	var ids []uint32
	for _, v := range player.UserData.AccountActivity.Activities {
		if v.ActType != uint32(publicconst.LoginActivity) {
			continue
		}

		v.EndTime = v.StartTime
		ids = append(ids, v.ActId)
	}

	for _, id := range ids {
		DeleteActivity(player, id)
	}

	return allFinish
}

// DeleteActivity 删除活动
func DeleteActivity(p *player.Player, actId uint32) {
	data := p.UserData.AccountActivity.Activities
	for i := 0; i < len(data); i++ {
		if data[i].ActId == actId {
			log.Info("deleteActivity", zap.Uint64("accountId", p.GetUserId()),
				zap.Uint32("actId", actId), zap.Reflect("activityInfo", data[i]))
			data = append(data[:i], data[i+1:]...)
			p.UserData.AccountActivity.Activities = data
			p.SaveAccountActivity()
			break
		}
	}
}

// GetActivityRedPoint 活动红点
func GetActivityRedPoint(p *player.Player) *msg.RedPointInfo {
	if p.UserData.AccountActivity == nil {
		log.Error("activity nil", zap.Uint64("accountid", p.GetUserId()))
		return nil
	}

	var ret = &msg.RedPointInfo{
		RdType: msg.RedPointType_Act_Point,
	}
	for i := 0; i < len(p.UserData.AccountActivity.Activities); i++ {
		activity := p.UserData.AccountActivity.Activities[i]
		if f, ok := getActivityRedPointFuncMap[publicconst.ActivityType(activity.ActType)]; ok {
			if f(activity) > 0 {
				ret.RdData = append(ret.RdData, activity.ActId)
			}
		}
	}
	return ret
}

// getActivityByType
func getActivityByType(p *player.Player, actType uint32) *model.Activity {
	for i := 0; i < len(p.UserData.AccountActivity.Activities); i++ {
		ret := p.UserData.AccountActivity.Activities[i]
		if ret.ActType == actType {
			return ret
		}
	}
	return nil
}

func GmFinishTodayOpenTask(p *player.Player) msg.ErrCode {
	openActivity := uint32(10401)
	activity := getActivity(p, openActivity)
	if activity == nil {
		return msg.ErrCode_ACTIVITY_NOT_EXIST
	}

	for i := 0; i < len(activity.ActDatas); i++ {
		subTask := activity.ActDatas[i]
		if subTask.State >= uint32(publicconst.TASK_COMPLETE) {
			continue
		}
		subTask.State = uint32(publicconst.TASK_COMPLETE)
		subTask.UpdateTime = tools.GetCurTime()
	}
	return msg.ErrCode_SUCC
}

func GmDayResetActivity(p *player.Player) {
	if p.UserData.AccountActivity == nil {
		return
	}

	for i := 0; i < len(p.UserData.AccountActivity.Activities); i++ {
		activity := p.UserData.AccountActivity.Activities[i]

		for i := 0; i < len(activity.ActDatas); i++ {
			needSave := false
			if activity.ActDatas[i].Value3 == 5 {
				needSave = true
				activity.ActDatas[i].Value4 = int64(tools.GetHourRefreshTime(template.GetSystemItemTemplate().RefreshHour))
			} else if activity.ActDatas[i].Value3 == 6 {
				needSave = true
				activity.ActDatas[i].Value4 = int64(tools.GetWeeklyRefreshTime(template.GetSystemItemTemplate().RefreshHour))
			}

			if needSave {
				p.SaveAccountActivity()
			}
		}
	}
}

// buyTaskPass 购买任务战令
func buyTaskPass(p *player.Player, passId int) {
	activity := getActivityByType(p, uint32(publicconst.TaskPass))
	if activity == nil {
		return
	}
	activity.ActDatas[0].Value3 = 1
	activity.ActDatas[0].UpdateTime = tools.GetCurTime()
	p.SaveAccountActivity()

	_, rewardItems := getGradeReward(activity, p)

	var retItems []*model.SimpleItem
	for i := 0; i < len(rewardItems); i++ {
		retItems = append(retItems, &model.SimpleItem{Id: rewardItems[i].Id, Num: rewardItems[i].Num})
	}
	res2 := &msg.NotifyRewardItem{}
	res2.GetItems = ToProtocolSimpleItems(retItems)
	p.SendNotify(res2)
}

// AddActivePassExp 添加活跃战令经验
func AddActivePassExp(p *player.Player, exp uint32) {
	for i := 0; i < len(p.UserData.AccountActivity.Activities); i++ {
		if p.UserData.AccountActivity.Activities[i].ActType == uint32(publicconst.ActivePass) {
			addActivityExp(p, p.UserData.AccountActivity.Activities[i], exp)
		}
	}
}

func ToProtocolActConfig(data *model.Activity) *msg.ActConfig {
	return &msg.ActConfig{
		ActId:     data.ActId,
		StartTime: data.StartTime,
		EndTime:   data.EndTime,
	}
}

func ToProtocolActConfigs(data []*model.Activity) []*msg.ActConfig {
	var ret []*msg.ActConfig
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolActConfig(data[i]))
	}
	return ret
}

func ToProtocolActDatas(data []*model.ActivityData) []*msg.ActData {
	var ret []*msg.ActData
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolActData(data[i]))
	}
	return ret
}

// GetActivityData 获得活动数据
func GetActivityData(p *player.Player, actId uint32) (msg.ErrCode, []*model.ActivityData) {
	activity := getActivity(p, actId)
	if activity == nil {
		return msg.ErrCode_ACTIVITY_NOT_EXIST, nil
	}

	updateActivityData(p, activity)
	return msg.ErrCode_SUCC, activity.ActDatas
}

// GetActivityReward 获取活动奖励
func GetActivityReward(p *player.Player, actId, subActId, num uint32) (msg.ErrCode, []*model.SimpleItem) {
	activity := getActivity(p, actId)
	if activity == nil {
		return msg.ErrCode_INVALID_DATA, nil
	}

	if f, ok := getActivityRewardMap[publicconst.ActivityType(activity.ActType)]; ok {
		activity := getActivity(p, actId)
		if activity == nil || !isValidActivity(activity) {
			return msg.ErrCode_ACTIVITY_NOT_EXIST, nil
		}
		return f(activity, p, subActId, num)
	}
	return msg.ErrCode_INVALID_DATA, nil
}

// BuyPassGrade 购买战令等级
func BuyPassGrade(p *player.Player, actId, grade uint32) msg.ErrCode {
	activity := getActivity(p, actId)
	if activity == nil || !isValidActivity(activity) {
		return msg.ErrCode_ACTIVITY_NOT_EXIST
	}
	if activity.ActType == uint32(publicconst.TaskPass) {
		return buyTaskPassGrade(p, activity, grade)
	} else if activity.ActType == uint32(publicconst.ActivePass) {
		return buyPassGrade(p, activity, grade)
	}
	return msg.ErrCode_INVALID_DATA
}

func buyTaskPassGrade(p *player.Player, activity *model.Activity, grade uint32) msg.ErrCode {
	temp := activity.ActDatas[0].Value1 % 100
	costNum := 100 - temp
	if !EnoughItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_DIAMOND), costNum) {
		return msg.ErrCode_NO_ENOUGH_ITEM
	}

	CostItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_DIAMOND), costNum, publicconst.BuyTaskPassGradeCostItem, false)
	updateClientItemsChange(p.GetUserId(), []uint32{uint32(publicconst.ITEM_CODE_DIAMOND)})

	activity.ActDatas[0].Value1 += costNum
	activity.ActDatas[0].Value2 = int64(template.GetTaskPassTemplate().GetGradeByExp(activity.ActId, activity.ActDatas[0].Value1))
	activity.ActDatas[0].UpdateTime = tools.GetCurTime()
	p.SaveAccountActivity()

	NtfActivityChange(p, activity, activity.ActDatas)
	return msg.ErrCode_SUCC
}

// SuppSign 请求补签
func SuppSign(p *player.Player, actId uint32) (msg.ErrCode, []*model.SimpleItem) {
	activity := getActivity(p, actId)
	if activity == nil || !isValidActivity(activity) {
		return msg.ErrCode_ACTIVITY_NOT_EXIST, nil
	}
	if activity.ActType == uint32(publicconst.Desert) {
		return suppSign(p, activity)
	}
	return msg.ErrCode_INVALID_DATA, nil
}

// suppSign 补签
func suppSign(p *player.Player, activity *model.Activity) (msg.ErrCode, []*model.SimpleItem) {
	maxDays := template.GetSignActivityTemplate().GetSignMaxDays(activity.ActId)
	if uint32(activity.ActDatas[0].Value2) == maxDays {
		return msg.ErrCode_NOT_SUPP_SIGN, nil
	}

	curTime := tools.GetCurTime()
	startTime := tools.GetDateStart(time.Unix(int64(activity.StartTime), 0))
	endTime := startTime.AddDate(0, 0, int(maxDays))
	if curTime < uint32(endTime.Unix()) {
		return msg.ErrCode_SUPP_SIGN_TIME_NOT_OPEN, nil
	}

	pos := int(activity.ActDatas[0].Value4)
	arr := template.GetSystemItemTemplate().SuppSignDiamondArr
	num := arr[len(arr)-1]
	if pos < len(arr) {
		num = arr[pos]
	}

	if !EnoughItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_DIAMOND), num) {
		return msg.ErrCode_NO_ENOUGH_ITEM, nil
	}

	// 发放奖励
	rewardItems := make(map[uint32]uint32)
	activity.ActDatas[0].Value2 += 1
	activity.ActDatas[0].Value4 += 1
	// 更新活动
	p.SaveAccountActivity()

	signConfig := template.GetSignActivityTemplate().GetSign(activity.ActId, uint32(activity.ActDatas[0].Value2))
	temp := signConfig.GetReward(activity.ActDatas[0].Value1)
	for id, num := range temp {
		if _, ok := rewardItems[id]; ok {
			rewardItems[id] += num
		} else {
			rewardItems[id] = num
		}
	}

	var notifyItems []uint32
	CostItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_DIAMOND),
		num, publicconst.DesrtSuppSignCostItem, false)
	notifyItems = append(notifyItems, uint32(publicconst.ITEM_CODE_DIAMOND))
	for id, num := range rewardItems {
		addItems := AddItem(p.GetUserId(), id, int32(num), publicconst.DesertSuppSignAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	NtfActivityChange(p, activity, activity.ActDatas)

	var retItems []*model.SimpleItem
	//tdaItems := make([]*tda.Item, 0, len(rewardItems))
	for id, num := range rewardItems {
		retItems = append(retItems, &model.SimpleItem{Id: id, Num: num})
		//tdaItems = append(tdaItems, &tda.Item{ItemId: strconv.Itoa(int(id)), ItemNum: num})
	}

	// write desert retroSigninNum
	p.UserData.Desert.RetroSignin++
	p.SaveDesert()

	// tda
	//tda.TdaEventDessertSignin(p.ChannelId, p.TdaCommonAttr, true, p.Desert.RetroSignin, num, tdaItems)

	return msg.ErrCode_SUCC, retItems
}

// EnterDesert 进入沙漠活动
// func EnterDesert(p *player.Player) (msg.ErrCode, uint32, []*model.ActivityData) {
// 	activity := getActivityByType(p, uint32(publicconst.Desert))
// 	if activity == nil || !isValidActivity(activity) {
// 		return msg.ErrCode_ACTIVITY_NOT_EXIST, 0, nil
// 	}
// 	if activity.ActType != uint32(publicconst.Desert) {
// 		return msg.ErrCode_ACTIVITY_NOT_EXIST, 0, nil
// 	}
// 	return enterDesert(p, activity)
// }

// // enterDesert 进入沙漠
// func enterDesert(p *player.Player, activity *model.Activity) (msg.ErrCode, uint32, []*model.ActivityData) {
// 	//rankData := ServMgr.GetRankService().GetDesertRank().GetRankData()

// 	// TODO rank
// 	var rank uint32 = 0
// 	// for i := 0; i < len(rankData); i++ {
// 	// 	if rankData[i].AccountId == uint32(p.GetUserId()) {
// 	// 		rank = uint32(i + 1)
// 	// 		break
// 	// 	}
// 	// }

// 	return msg.ErrCode_SUCC, rank, activity.ActDatas
// }

// GetActPreviewReward 请求获得活动预览奖励
func GetActPreviewReward(p *player.Player, tp uint32) (msg.ErrCode, []*model.SimpleItem) {
	if tools.ListContain(p.UserData.AccountActivity.PreRewardTps, tp) {
		return msg.ErrCode_ACTIVITY_HAS_GET_REWARD, nil
	}
	if tp == uint32(publicconst.Desert) {
		return getDesertActPreviewReward(p)
	}
	return msg.ErrCode_INVALID_DATA, nil
}

// getDesertActPreviewReward 沙漠活动预览奖励
func getDesertActPreviewReward(p *player.Player) (msg.ErrCode, []*model.SimpleItem) {
	rewardItems := template.GetSystemItemTemplate().DesertInfoReward
	var ret []*model.SimpleItem
	var notifyItems []uint32
	for i := 0; i < len(rewardItems); i++ {
		addItems := AddItem(p.GetUserId(),
			rewardItems[i].ItemId, int32(rewardItems[i].ItemNum),
			publicconst.DesertActPreviewAddItem, false)
		ret = append(ret, addItems...)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	p.UserData.AccountActivity.PreRewardTps = append(p.UserData.AccountActivity.PreRewardTps,
		uint32(publicconst.Desert))
	p.SaveAccountActivity()
	return msg.ErrCode_SUCC, ret
}
