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

func getTaskListByType(p *player.Player, taskType publicconst.TaskType) []*model.Task {
	var tasks []*model.Task
	switch taskType {
	case publicconst.DAILY_TASK:
		tasks = p.UserData.Task.DailyTasks
	case publicconst.WEEKLY_TASK:
		tasks = p.UserData.Task.WeeklyTasks
	case publicconst.MAIN_TASK:
		tasks = p.UserData.Task.MainTasks
	case publicconst.ACHIEVE_TASK:
		tasks = p.UserData.Task.AchieveTasks
	case publicconst.ALLIANCE_WEEKLY_TASK:
		tasks = p.UserData.Task.AllianceWeeklyTasks
	}
	return tasks
}

// GetTaskActiveReward 获取任务活跃度奖励
func GetTaskActiveReward(p *player.Player, taskType publicconst.TaskType, pos uint32) (msg.ErrCode, uint32) { // 根据pos获取活跃度奖励
	if taskType != publicconst.DAILY_TASK && taskType != publicconst.WEEKLY_TASK {
		return msg.ErrCode_INVALID_DATA, 0
	}

	var notifyItems []uint32
	var taskActiveReward uint32
	if taskType == publicconst.DAILY_TASK { // 日任务
		taskActives := template.GetTaskActiveTemplate().GetDailyTaskActiveValue(p.UserData.Task.DailyActiveValue)
		for i := 0; i < len(taskActives) && i <= int(pos); i++ {
			if p.UserData.Task.DailyActiveReward&(1<<i) == 0 {
				p.UserData.Task.DailyActiveReward |= 1 << i
				for k := 0; k < len(taskActives[i].RewardItems); k++ {
					addItems := AddItem(p.GetUserId(),
						taskActives[i].RewardItems[k].ItemId,
						int32(taskActives[i].RewardItems[k].ItemNum),
						publicconst.GetTaskActiveRewardAddItem, false)
					notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
				}
			}

		}
		if len(notifyItems) > 0 {
			p.SaveTask()
		}
		taskActiveReward = p.UserData.Task.DailyActiveReward

		// 行为统计
		UpdateTask(p, true, publicconst.TASK_COND_GET_DAILY_BOX, 1)
	} else { // 周任务
		taskActives := template.GetTaskActiveTemplate().GetWeeklyTaskActiveValue(p.UserData.Task.WeeklyActiveValue)
		for i := 0; i < len(taskActives) && i <= int(pos); i++ {
			if p.UserData.Task.WeeklyActiveReward&(1<<i) == 0 {
				p.UserData.Task.WeeklyActiveReward |= 1 << i
				for k := 0; k < len(taskActives[i].RewardItems); k++ {
					addItems := AddItem(p.GetUserId(),
						taskActives[i].RewardItems[k].ItemId,
						int32(taskActives[i].RewardItems[k].ItemNum),
						publicconst.GetTaskActiveRewardAddItem, false)
					notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
				}
			}
		}
		if len(notifyItems) > 0 {
			p.SaveTask()
		}
		taskActiveReward = p.UserData.Task.WeeklyActiveReward

	}

	if len(notifyItems) == 0 {
		return msg.ErrCode_NO_ENOUGH_TASK_ACTIVE, 0
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)
	return msg.ErrCode_SUCC, taskActiveReward
}

func GetAllTaskReward(p *player.Player, taskType publicconst.TaskType) (msg.ErrCode, []*template.SimpleItem) {
	tasks := getTaskListByType(p, taskType)
	if len(tasks) == 0 {
		return msg.ErrCode_NO_COMPLETE_TASK, nil
	}

	var (
		addItemMap       = make(map[uint32]uint32) // 奖励
		notifyClientTask []*model.Task

		notifyItems []uint32
		finalItems  []*template.SimpleItem

		addDayNum     = uint32(0)
		addDayActive  = uint32(0) // 日活跃度
		addWeekNum    = uint32(0)
		addWeekActive = uint32(0) // 周活跃度
	)

	for {
		n := 0
		for _, v := range tasks {
			if v.TaskState != uint32(publicconst.TASK_COMPLETE) {
				continue
			}

			taskCfg := template.GetTaskTemplate().GetTask(v.TaskId)
			if taskCfg == nil {
				log.Error("task cfg nil", zap.Uint32("taskId", v.TaskId))
				continue
			}
			v.TaskState = uint32(publicconst.TASK_DONE)

			if taskCfg.Data.DayActive > 0 {
				addDayActive += taskCfg.Data.DayActive
			}
			if taskCfg.Data.WeekActive > 0 {
				addWeekActive += taskCfg.Data.WeekActive
			}

			//tdaItemSlice := make([]*tda.Item, 0, len(taskCfg.RewardItems))
			for _, vv := range taskCfg.RewardItems {
				addItemMap[vv.ItemId] += vv.ItemNum
				//tdaItemSlice = append(tdaItemSlice, &tda.Item{ItemId: strconv.Itoa(int(vv.ItemId)), ItemNum: vv.ItemNum})
			}

			taskType := publicconst.TaskType(taskCfg.Data.TaskType)
			if taskCfg.NextTaskId != 0 { // 添加后续任务
				if nextTasks := template.GetTaskTemplate().GetNextTasks(taskCfg.Data.Id); len(nextTasks) > 0 {
					for _, vv := range nextTasks {
						var initValue uint32 = 0
						if vv.Data.UsePreValue == 1 {
							initValue = v.TaskValue
						}
						if addTask := AddTask(p, vv, initValue); addTask != nil {
							notifyClientTask = append(notifyClientTask, addTask)
							tasks = append(tasks, addTask)
						}
					}
				}
				DeleteTask(p, taskType, v) // 删除已完成任务
			} else {
				notifyClientTask = append(notifyClientTask, v)
			}

			if taskType == publicconst.DAILY_TASK {
				addDayNum++
			} else if taskType == publicconst.WEEKLY_TASK {
				addWeekNum++
			}

			// tda
			// switch taskType {
			// case publicconst.DAILY_TASK:
			// 	tda.TdaTaskRoutineReward(p.ChannelId, p.TdaCommonAttr, v.TaskId, taskCfg.Data.TaskType, taskCfg.Data.DayActive, p.Task.DailyActiveValue, taskCfg.Data.DayActive)
			// case publicconst.WEEKLY_TASK:
			// 	tda.TdaTaskRoutineReward(p.ChannelId, p.TdaCommonAttr, v.TaskId, taskCfg.Data.TaskType, taskCfg.Data.WeekActive, p.Task.WeeklyActiveValue, taskCfg.Data.WeekActive)
			// case publicconst.MAIN_TASK:
			// 	tda.TdaTaskMainReward(p.ChannelId, p.TdaCommonAttr, v.TaskId, taskCfg.Data.TaskType, tdaItemSlice)
			// case publicconst.ACHIEVE_TASK:
			// 	tda.TdaTaskAchiReward(p.ChannelId, p.TdaCommonAttr, v.TaskId, taskCfg.Data.TaskType, tdaItemSlice)
			// case publicconst.ALLIANCE_WEEKLY_TASK:
			// 	tda.TdaTaskGuildReward(p.ChannelId, p.TdaCommonAttr, v.TaskId, taskCfg.Data.TaskType, tdaItemSlice)
			// }

			n++
		}

		if n == 0 {
			break
		}
	}

	if len(addItemMap) > 0 {
		for k, v := range addItemMap {
			addItems := AddItem(p.GetUserId(), k, int32(v), publicconst.GetTaskRewardAddItem, false)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
			finalItems = append(finalItems, ToTemplateItem(addItems)...)
		}
	}

	if addDayNum > 0 {
		UpdateTask(p, true, publicconst.TASK_COND_COMPLETE_DAILY_TASK_COUNT, addDayNum)
		processHistoryData(p, publicconst.TASK_COND_COMPLETE_DAILY_TASK_COUNT, 0, addDayNum)
	}
	if addWeekNum > 0 {
		UpdateTask(p, true, publicconst.TASK_COND_COMPLETE_WEEKLY_TASK_COUNT, addWeekNum)
		processHistoryData(p, publicconst.TASK_COND_COMPLETE_WEEKLY_TASK_COUNT, 0, addWeekNum)
	}
	save := false
	sendActive := false
	if addDayActive > 0 {
		p.UserData.Task.DailyActiveValue += addDayActive

		maxValue := template.GetTaskActiveTemplate().GetDailyTaskActiveMaxValue()
		if p.UserData.Task.DailyActiveValue > maxValue {
			p.UserData.Task.DailyActiveValue = maxValue
		}
		AddActivePassExp(p, addDayActive)
		UpdateTask(p, true, publicconst.TASK_COND_DAILY_ACTIVE_SCORE, addDayActive)
		processHistoryData(p, publicconst.TASK_COND_DAILY_ACTIVE_SCORE, 0, addDayActive)

		sendActive = true
	}
	if addWeekActive > 0 {
		p.UserData.Task.WeeklyActiveValue += addWeekActive

		maxValue := template.GetTaskActiveTemplate().GetWeeklyTaskActiveMaxValue()
		if p.UserData.Task.WeeklyActiveValue > maxValue {
			p.UserData.Task.WeeklyActiveValue = maxValue
		}
		UpdateTask(p, true, publicconst.TASK_COND_WEEKLY_ACTIVE_SCORE, addWeekActive)
		processHistoryData(p, publicconst.TASK_COND_WEEKLY_ACTIVE_SCORE, 0, addWeekActive)

		sendActive = true
	}

	if sendActive {
		save = true
		NotifyActiveChange(p)
	}

	if len(notifyClientTask) > 0 {
		save = true
		NotifyClientTaskChange(p, notifyClientTask)
	}

	if len(notifyItems) > 0 {
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	if save {
		p.SaveTask()
	}

	return msg.ErrCode_SUCC, finalItems
}

// GetTaskReward 获取任务奖励
func GetTaskReward(p *player.Player, taskId uint32) msg.ErrCode {
	jTask := template.GetTaskTemplate().GetTask(taskId)
	if jTask == nil {
		return msg.ErrCode_TASK_NOT_EXIST
	}

	task := getTaskByType(p, jTask.Data.Id, publicconst.TaskType(jTask.Data.TaskType))
	if task == nil {
		return msg.ErrCode_TASK_NOT_EXIST
	}

	if task.TaskState == uint32(publicconst.TASK_DONE) {
		return msg.ErrCode_TASK_HAS_GET_REWARD
	}

	if task.TaskState == uint32(publicconst.TASK_ACCEPT) {
		return msg.ErrCode_TASK_NOT_COMPLETE
	}

	task.TaskState = uint32(publicconst.TASK_DONE)

	// 增加活跃度
	var activeChange = false
	if jTask.Data.DayActive > 0 {
		p.UserData.Task.DailyActiveValue += jTask.Data.DayActive
		maxValue := template.GetTaskActiveTemplate().GetDailyTaskActiveMaxValue()
		if p.UserData.Task.DailyActiveValue > maxValue {
			p.UserData.Task.DailyActiveValue = maxValue
		}
		activeChange = true
		UpdateTask(p, true,
			publicconst.TASK_COND_DAILY_ACTIVE_SCORE, jTask.Data.DayActive)
		AddActivePassExp(p, jTask.Data.DayActive)
		processHistoryData(p, publicconst.TASK_COND_DAILY_ACTIVE_SCORE, 0, jTask.Data.DayActive)
	}
	if jTask.Data.WeekActive > 0 {
		p.UserData.Task.WeeklyActiveValue += jTask.Data.WeekActive
		maxValue := template.GetTaskActiveTemplate().GetWeeklyTaskActiveMaxValue()
		if p.UserData.Task.WeeklyActiveValue > maxValue {
			p.UserData.Task.WeeklyActiveValue = maxValue
		}
		activeChange = true
		UpdateTask(p, true,
			publicconst.TASK_COND_WEEKLY_ACTIVE_SCORE, jTask.Data.WeekActive)
		processHistoryData(p, publicconst.TASK_COND_WEEKLY_ACTIVE_SCORE, 0, jTask.Data.WeekActive)
	}
	if jTask.Data.TaskType == publicconst.ALLIANCE_WEEKLY_TASK {
		UpdateTask(p, true,
			publicconst.TASK_COND_ALLIANCE_FINISH_TASK, 1)
	}

	// 领取任务奖励
	var notifyItems []uint32
	//tdaReward := make([]*tda.Item, 0, len(jTask.RewardItems))
	for i := 0; i < len(jTask.RewardItems); i++ {
		addItems := AddItem(p.GetUserId(), jTask.RewardItems[i].ItemId,
			int32(jTask.RewardItems[i].ItemNum), publicconst.GetTaskRewardAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)

		//tdaReward = append(tdaReward, &tda.Item{ItemId: strconv.Itoa(int(jTask.RewardItems[i].ItemId)), ItemNum: jTask.RewardItems[i].ItemNum})
	}

	var notifyClientTask []*model.Task
	// 有后续任务
	taskType := publicconst.TaskType(jTask.Data.TaskType)
	if jTask.NextTaskId != 0 {
		if tasks := template.GetTaskTemplate().GetNextTasks(jTask.Data.Id); len(tasks) > 0 {
			for i := 0; i < len(tasks); i++ {
				var initValue uint32 = 0
				if tasks[i].Data.UsePreValue == 1 {
					initValue = task.TaskValue
				}
				if addTask := AddTask(p, tasks[i], initValue); addTask != nil {
					notifyClientTask = append(notifyClientTask, addTask)
				}
			}
		}
		DeleteTask(p, taskType, task)
	} else {
		notifyClientTask = append(notifyClientTask, task)
		p.SaveTask()
	}

	// 通知任务更新
	if len(notifyClientTask) > 0 {
		NotifyClientTaskChange(p, notifyClientTask)
	}

	// 通知活跃度更新
	if activeChange {
		p.SaveTask()
		NotifyActiveChange(p)
	}

	// 更新任务进度
	if taskType == publicconst.DAILY_TASK {
		UpdateTask(p, true,
			publicconst.TASK_COND_COMPLETE_DAILY_TASK_COUNT, 1)
		processHistoryData(p, publicconst.TASK_COND_COMPLETE_DAILY_TASK_COUNT, 0, 1)
	} else if taskType == publicconst.WEEKLY_TASK {
		UpdateTask(p, true,
			publicconst.TASK_COND_COMPLETE_WEEKLY_TASK_COUNT, 1)
		processHistoryData(p, publicconst.TASK_COND_COMPLETE_WEEKLY_TASK_COUNT, 0, 1)
	}

	// tda
	// switch taskType {
	// case publicconst.DAILY_TASK:
	// 	tda.TdaTaskRoutineReward(p.ChannelId, p.TdaCommonAttr, taskId, jTask.Data.TaskType, jTask.Data.DayActive, p.Task.DailyActiveValue, jTask.Data.DayActive)
	// case publicconst.WEEKLY_TASK:
	// 	tda.TdaTaskRoutineReward(p.ChannelId, p.TdaCommonAttr, taskId, jTask.Data.TaskType, jTask.Data.WeekActive, p.Task.WeeklyActiveValue, jTask.Data.WeekActive)
	// case publicconst.MAIN_TASK:
	// 	tda.TdaTaskMainReward(p.ChannelId, p.TdaCommonAttr, taskId, jTask.Data.TaskType, tdaReward)
	// case publicconst.ACHIEVE_TASK:
	// 	tda.TdaTaskAchiReward(p.ChannelId, p.TdaCommonAttr, taskId, jTask.Data.TaskType, tdaReward)
	// case publicconst.ALLIANCE_WEEKLY_TASK:
	// 	tda.TdaTaskGuildReward(p.ChannelId, p.TdaCommonAttr, taskId, jTask.Data.TaskType, tdaReward)
	// }

	// 通知道具更新
	//	event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(playerData, notifyItems, ListenNotifyClientItemEventEvent))

	updateClientItemsChange(p.GetUserId(), notifyItems)

	return msg.ErrCode_SUCC
}

// UpdateTask 更新任务进度  更新所有同cond的任务进度
func UpdateTask(p *player.Player, sendClient bool, cond publicconst.TaskCond, args ...uint32) {
	var tasks []*model.Task
	if temp := updateTaskByType(p, publicconst.DAILY_TASK, cond, args...); len(temp) > 0 {
		tasks = append(tasks, temp...)
	}
	if temp := updateTaskByType(p, publicconst.WEEKLY_TASK, cond, args...); len(temp) > 0 {
		tasks = append(tasks, temp...)
	}
	if temp := updateTaskByType(p, publicconst.MAIN_TASK, cond, args...); len(temp) > 0 {
		tasks = append(tasks, temp...)
	}
	if temp := updateTaskByType(p, publicconst.ACHIEVE_TASK, cond, args...); len(temp) > 0 {
		tasks = append(tasks, temp...)
	}
	if temp := updateTaskByType(p, publicconst.ALLIANCE_WEEKLY_TASK, cond, args...); len(temp) > 0 {
		tasks = append(tasks, temp...)
	}

	if len(tasks) > 0 && sendClient {
		NotifyClientTaskChange(p, tasks)
	}

	// 更新活动
	UpdateActivity(p, uint32(cond), args...)
	UpdateLuckSaleTask(p, uint32(cond), args...)
}

// updateTaskByType
func updateTaskByType(p *player.Player, taskType publicconst.TaskType, cond publicconst.TaskCond, args ...uint32) []*model.Task {
	if p.UserData.Task == nil {
		log.Error("task data nil", zap.Uint64("accountId", p.GetUserId()))
		return nil
	}
	var tasks []*model.Task
	switch taskType {
	case publicconst.DAILY_TASK:
		tasks = p.UserData.Task.DailyTasks
	case publicconst.WEEKLY_TASK:
		tasks = p.UserData.Task.WeeklyTasks
	case publicconst.MAIN_TASK:
		tasks = p.UserData.Task.MainTasks
	case publicconst.ACHIEVE_TASK:
		tasks = p.UserData.Task.AchieveTasks
	case publicconst.ALLIANCE_WEEKLY_TASK:
		tasks = p.UserData.Task.AllianceWeeklyTasks
	}
	ntf_map := make(map[uint32]*model.Task)
	var notifyTask []*model.Task
	for i := 0; i < len(tasks); i++ {
		jTask := template.GetTaskTemplate().GetTask(tasks[i].TaskId)
		if jTask == nil {
			log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", tasks[i].TaskId))
			continue
		}

		if jTask.Data.TaskCondition != uint32(cond) {
			continue
		}

		if tasks[i].TaskState == uint32(publicconst.TASK_COMPLETE) ||
			tasks[i].TaskState == uint32(publicconst.TASK_DONE) {
			continue
		}

		if refreshTaskValue(p, tasks[i], cond, jTask.Data.MaxValue, jTask.Data.Effect1, args...) {

			if _, ok := ntf_map[tasks[i].TaskId]; !ok {
				ntf_map[tasks[i].TaskId] = tasks[i]
			}
		}
	}
	for _, v := range ntf_map {
		notifyTask = append(notifyTask, v)
	}
	if len(tasks) > 0 {
		p.SaveTask()
	}
	return notifyTask
}

// CanAddTask 能否添加任务
func CanAddTask(p *player.Player, jTask *template.JTask) bool {
	if jTask.Data.PlayerLevel > 0 {
		if p.UserData.Level < jTask.Data.PlayerLevel {
			return false
		}
	}

	if jTask.Data.PreMission > 0 {
		if mission := findMission(p, int(jTask.Data.PreMission), true); mission == nil {
			return false
		}
	}

	if jTask.Data.PreTask > 0 {
		if !tools.ListContain(p.UserData.Task.FinshedTasks, jTask.Data.PreTask) {
			return false
		}
	}
	return true
}

// AddTask 添加任务
func AddTask(p *player.Player, jTask *template.JTask, initValue uint32) *model.Task {
	switch publicconst.TaskType(jTask.Data.TaskType) {
	case publicconst.ACHIEVE_TASK:
		if task := getTaskByType(p, jTask.Data.Id, publicconst.ACHIEVE_TASK); task == nil {
			task = model.NewTask(jTask.Data.Id, initValue)
			initTaskValue(p, task, publicconst.TaskCond(jTask.Data.TaskCondition))
			if jTask.Data.UseHisData == 1 {
				initHistoryDataTask(p, task)
			}
			p.UserData.Task.AchieveTasks = append(p.UserData.Task.AchieveTasks, task)
			p.SaveTask()
			return task
		}
	case publicconst.MAIN_TASK:
		if task := getTaskByType(p, jTask.Data.Id, publicconst.MAIN_TASK); task == nil {
			task = model.NewTask(jTask.Data.Id, initValue)
			initTaskValue(p, task, publicconst.TaskCond(jTask.Data.TaskCondition))
			if jTask.Data.UseHisData == 1 {
				initHistoryDataTask(p, task)
			}
			p.UserData.Task.MainTasks = append(p.UserData.Task.MainTasks, task)
			p.SaveTask()
			return task
		}
	}
	return nil
}

// getTaskByType 通过任务类型获得任务
func getTaskByType(p *player.Player, taskId uint32, taskType publicconst.TaskType) *model.Task {
	var tasks []*model.Task
	switch taskType {
	case publicconst.DAILY_TASK:
		tasks = p.UserData.Task.DailyTasks
	case publicconst.WEEKLY_TASK:
		tasks = p.UserData.Task.WeeklyTasks
	case publicconst.MAIN_TASK:
		tasks = p.UserData.Task.MainTasks
	case publicconst.ACHIEVE_TASK:
		tasks = p.UserData.Task.AchieveTasks
	case publicconst.ALLIANCE_WEEKLY_TASK:
		tasks = p.UserData.Task.AllianceWeeklyTasks
	}
	for i := 0; i < len(tasks); i++ {
		if tasks[i].TaskId == taskId {
			return tasks[i]
		}
	}
	return nil
}

// GetTasksByType 通过任务类型获得任务
func GetTasksByType(p *player.Player, taskType publicconst.TaskType) []*model.Task {
	if p.UserData.Task == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()))
		return nil
	}

	curTime := tools.GetCurTime()
	var ret []*model.Task
	switch taskType {
	case publicconst.DAILY_TASK:
		// 需要重置日常任务
		if curTime >= p.UserData.Task.NextDailyRefreshTime {
			p.UserData.Task.NextDailyRefreshTime = tools.GetDailyRefreshTime()
			p.UserData.Task.DailyTasks = p.UserData.Task.DailyTasks[0:0]
			p.UserData.Task.DailyActiveValue = 0
			p.UserData.Task.DailyActiveReward = 0
			tasks := template.GetTaskTemplate().GetTaskByType(uint32(publicconst.DAILY_TASK))
			for i := 0; i < len(tasks); i++ {
				p.UserData.Task.DailyTasks = append(p.UserData.Task.DailyTasks,
					model.NewTask(tasks[i].Data.Id, 0))
			}
			p.SaveTask()
		}
		ret = p.UserData.Task.DailyTasks
	case publicconst.WEEKLY_TASK:
		if curTime >= p.UserData.Task.NextWeeklyRefreshTime {
			p.UserData.Task.NextWeeklyRefreshTime = tools.GetWeeklyRefreshTime(template.GetSystemItemTemplate().RefreshHour)
			p.UserData.Task.WeeklyTasks = p.UserData.Task.WeeklyTasks[0:0]
			p.UserData.Task.WeeklyActiveValue = 0
			p.UserData.Task.WeeklyActiveReward = 0
			tasks := template.GetTaskTemplate().GetTaskByType(uint32(publicconst.WEEKLY_TASK))
			for i := 0; i < len(tasks); i++ {
				p.UserData.Task.WeeklyTasks = append(p.UserData.Task.WeeklyTasks,
					model.NewTask(tasks[i].Data.Id, 0))
			}

			p.SaveTask()
		}
		ret = p.UserData.Task.WeeklyTasks
	case publicconst.MAIN_TASK:
		if len(p.UserData.Task.MainTasks) == 0 {
			initTask := template.GetTaskTemplate().GetInitMainTask()
			for i := 0; i < len(initTask); i++ {
				AddTask(p, initTask[i], 0)
			}
		}
		ret = p.UserData.Task.MainTasks
	case publicconst.ACHIEVE_TASK:
		initTask := template.GetTaskTemplate().GetRootTask(2)
		for i := 0; i < len(initTask); i++ {
			exist := false
			for m := 0; m < len(p.UserData.Task.AchieveTasks); m++ {
				if taskConfig := template.GetTaskTemplate().GetTask(p.UserData.Task.AchieveTasks[m].TaskId); taskConfig != nil {
					if taskConfig.GetRooTTaskId() == initTask[i].Data.Id {
						exist = true
						break
					}
				}
			}
			if !exist {
				AddTask(p, initTask[i], 0)
			}
		}
		ret = p.UserData.Task.AchieveTasks
	case publicconst.ALLIANCE_WEEKLY_TASK:
		if curTime >= p.UserData.Task.NextAllianceWeeklyRefreshTime {
			p.UserData.Task.NextAllianceWeeklyRefreshTime = tools.GetWeeklyRefreshTime(template.GetSystemItemTemplate().RefreshHour)
			if p.UserData.Task.AllianceWeeklyTasks == nil {
				p.UserData.Task.AllianceWeeklyTasks = make([]*model.Task, 0)
			}
			p.UserData.Task.AllianceWeeklyTasks = p.UserData.Task.AllianceWeeklyTasks[0:0]
			p.UserData.Task.AllianceWeeklyActiveValue = 0
			p.UserData.Task.AllianceWeeklyActiveReward = 0
			tasks := template.GetTaskTemplate().GetTaskByType(uint32(publicconst.ALLIANCE_WEEKLY_TASK))
			for i := 0; i < len(tasks); i++ {
				p.UserData.Task.AllianceWeeklyTasks = append(p.UserData.Task.AllianceWeeklyTasks, model.NewTask(tasks[i].Data.Id, 0))
			}

			p.SaveTask()
		}
		ret = p.UserData.Task.AllianceWeeklyTasks
	}

	UpdateTask(p, false, publicconst.TASK_COND_LOGIN)
	return ret
}

// DeleteTask 删除任务
func DeleteTask(p *player.Player, taskType publicconst.TaskType, task *model.Task) {
	tasks := GetTasksByType(p, taskType)
	for i := 0; i < len(tasks); i++ {
		if tasks[i].TaskId == task.TaskId {
			tasks = append(tasks[0:i], tasks[i+1:]...)
			break
		}
	}
	switch taskType {
	case publicconst.MAIN_TASK:
		p.UserData.Task.MainTasks = tasks
	case publicconst.ACHIEVE_TASK:
		p.UserData.Task.AchieveTasks = tasks
	}
	p.UserData.Task.FinshedTasks = append(p.UserData.Task.FinshedTasks, task.TaskId)
	p.SaveTask()
}

// NotifyActiveChange 通知活跃度变化
func NotifyActiveChange(p *player.Player) {
	p.SendNotify(&msg.NotifyTaskActiveChange{
		DailyActive:          p.UserData.Task.DailyActiveValue,
		WeeklyActive:         p.UserData.Task.WeeklyActiveValue,
		AllianceWeeklyActive: p.UserData.Task.AllianceWeeklyActiveValue,
	})
}

// NotifyClientTaskChange 通知客户端任务变化
func NotifyClientTaskChange(p *player.Player, tasks []*model.Task) {
	if p == nil {
		return
	}

	if len(tasks) == 0 {
		return
	}

	p.SendNotify(&msg.NotifyTaskChange{
		Data: ToProtocolTasks(tasks),
	})
}

// GetHistoryData 获得历史数量
func GetHistoryData(p *player.Player, cond publicconst.TaskCond, condPara uint32) *model.TaskHistroyData {
	if p.UserData.Task != nil && p.UserData.Task.HistoryData != nil {
		for i := 0; i < len(p.UserData.Task.HistoryData); i++ {
			if p.UserData.Task.HistoryData[i].TaskCond == uint32(cond) {
				if condPara == 0 {
					return p.UserData.Task.HistoryData[i]
				} else if condPara == p.UserData.Task.HistoryData[i].CondPara {
					return p.UserData.Task.HistoryData[i]
				}
			}
		}
	}
	return nil
}

// addHistoryData 添加历史数据
func addHistoryData(p *player.Player, data *model.TaskHistroyData) {
	p.UserData.Task.HistoryData = append(p.UserData.Task.HistoryData, data)
	p.SaveTask()
}

// updateHistoryData 更新历史数据
func updateHistoryData(p *player.Player, data *model.TaskHistroyData) {
	p.SaveTask()
	//p.SaveTask()
}

// processHistoryData 通过条件类型存数据
func processHistoryData(p *player.Player, cond publicconst.TaskCond, condPara uint32, args ...uint32) {
	switch cond {
	case publicconst.TASK_COND_COMPLETE_DAILY_TASK_COUNT,
		publicconst.TASK_COND_COMPLETE_WEEKLY_TASK_COUNT,
		publicconst.TASK_COND_DAILY_ACTIVE_SCORE,
		publicconst.TASK_COND_WEEKLY_ACTIVE_SCORE,
		publicconst.TASK_COND_GET_EQUIP_NUM,
		publicconst.TASK_COND_GET_WEAPON_NUM,
		publicconst.TASK_COND_KILL_MONSTER_NUM,
		publicconst.TASK_COND_KILL_ELITE_MONSTER_BOSS_NUM,
		publicconst.TASK_COND_GET_ON_HOOK_REWARD,
		publicconst.TASK_COND_PEAK_FIGHT_PK,
		publicconst.TASK_COND_UPGRADE_EQUIP,
		publicconst.TASK_COND_FIRST_PASS_CHALLENGE_MISSION,
		publicconst.TASK_COND_RESUME_DIAMOND_BUY_SHOP,
		publicconst.TASK_COND_UPGRADE_WEAPON:
		if historyData := GetHistoryData(p, cond, condPara); historyData != nil {
			historyData.TaskValue += args[0]
			updateHistoryData(p, historyData)
		} else {
			historyData := model.NewTaskHistroyData(uint32(cond), args[0], condPara)
			addHistoryData(p, historyData)
		}
	case publicconst.TASK_COND_ADD_ITEM:
		if historyData := GetHistoryData(p, cond, condPara); historyData != nil {
			historyData.TaskValue += args[0]
			updateHistoryData(p, historyData)
		}
	case publicconst.TASK_COND_LOGIN:
		if historyData := GetHistoryData(p, cond, 0); historyData != nil {
			curTime := tools.GetCurTime()
			if curTime >= historyData.ExtraPara {
				historyData.TaskValue += args[0]
				historyData.ExtraPara = tools.GetDailyRefreshTime()
				updateHistoryData(p, historyData)
			}
		} else {
			historyData := model.NewTaskHistroyData(uint32(cond), args[0], 0)
			historyData.ExtraPara = tools.GetDailyRefreshTime()
			addHistoryData(p, historyData)
		}
	case publicconst.TASK_COND_LOTTERY, publicconst.TASK_COND_DISK_BIND_BOX_LOTTERY:
		if historyData := GetHistoryData(p, cond, condPara); historyData != nil {
			historyData.TaskValue += args[0]
			updateHistoryData(p, historyData)
		} else {
			historyData := model.NewTaskHistroyData(uint32(cond), args[0], condPara)
			addHistoryData(p, historyData)
		}
	}
}

func ToProtocolTasks(data []*model.Task) []*msg.Task {
	var ret []*msg.Task
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolTask(data[i]))
	}
	return ret
}

func ToProtocolTask(data *model.Task) *msg.Task {
	return &msg.Task{
		TaskId:       data.TaskId,
		TaskValue:    data.TaskValue,
		State:        msg.TaskState(data.TaskState),
		CompleteTime: data.CompleteTime,
	}
}
