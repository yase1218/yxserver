package model

import (
	"gameserver/internal/publicconst"
	"kernel/tools"
)

type AccountTask struct {
	AccountId                 int64
	DailyActiveValue          uint32 // 日常活跃度
	WeeklyActiveValue         uint32 // 周长活跃度
	AllianceWeeklyActiveValue uint32 // 联盟周任务活跃度(获取联盟勋章数量)

	NextDailyRefreshTime          uint32 // 下一次日常刷新时间
	NextWeeklyRefreshTime         uint32 // 下一次周常刷新时间
	NextAllianceWeeklyRefreshTime uint32 // 下一次联盟周任务刷新时间

	DailyActiveReward          uint32 // 日常活跃度奖励标志
	WeeklyActiveReward         uint32 // 周常活跃度奖励标志
	AllianceWeeklyActiveReward uint32 // 联盟周任务活跃度奖励标志

	DailyTasks          []*Task  // 日常任务
	WeeklyTasks         []*Task  // 周长任务
	AchieveTasks        []*Task  // 成就任务
	MainTasks           []*Task  // 主线任务
	AllianceWeeklyTasks []*Task  // 联盟每周任务
	FinshedTasks        []uint32 // 已经完成的任务 不包括日常周常

	HistoryData []*TaskHistroyData // 历史数据
}

type Task struct {
	TaskId       uint32
	TaskValue    uint32 // 进度
	TaskState    uint32
	ExtraPara    uint32 // 额外参数
	CompleteTime uint32 // 任务完成时间
	CreateTime   uint32
	UpdateTime   uint32
}

// TaskHistroyData 任务历史数据
type TaskHistroyData struct {
	TaskCond  uint32 // 任务条件
	TaskValue uint32 // 任务数值
	CondPara  uint32 // 条件参数
	ExtraPara uint32
}

func NewTask(taskId, taskValue uint32) *Task {
	curTime := tools.GetCurTime()
	return &Task{
		TaskId:     taskId,
		TaskValue:  taskValue,
		CreateTime: curTime,
		UpdateTime: curTime,
	}
}

func NewTaskHistroyData(taskCond, taskValue, condPara uint32) *TaskHistroyData {
	return &TaskHistroyData{
		TaskCond:  taskCond,
		TaskValue: taskValue,
		CondPara:  condPara,
	}
}

func NewAccountTask(accountId int64) *AccountTask {
	ret := &AccountTask{
		AccountId: accountId,
	}
	ret.DailyTasks = make([]*Task, 0, 0)
	ret.WeeklyTasks = make([]*Task, 0, 0)
	ret.AchieveTasks = make([]*Task, 0, 0)
	ret.MainTasks = make([]*Task, 0, 0)
	ret.FinshedTasks = make([]uint32, 0, 0)
	ret.HistoryData = make([]*TaskHistroyData, 0)
	return ret
}

func (t *Task) GetTaskId() uint32 {
	return t.TaskId
}

func (t *Task) GetTaskValue() uint32 {
	return t.TaskValue
}
func (t *Task) SetTaskValue(value uint32) {
	t.TaskValue = value
}
func (t *Task) AddTaskValue(add uint32) {
	t.TaskValue += add
}
func (t *Task) GetTaskState() publicconst.TaskState {
	return publicconst.TaskState(t.TaskState)
}

func (t *Task) SetTaskState(state publicconst.TaskState) {
	t.TaskState = uint32(state)
}

func (t *Task) SetTaskCompleteTime(time uint32) {
	t.CompleteTime = time
}

func (t *Task) GetTaskCompleteTime() uint32 {
	return t.CompleteTime
}

func (t *Task) SetExtraPara(value uint32) {
	t.ExtraPara = value
}

func (t *Task) GetExtraPara() uint32 {
	return t.ExtraPara
}
