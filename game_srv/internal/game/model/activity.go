package model

import (
	"kernel/tools"

	"gameserver/internal/publicconst"
)

type AccountActivity struct {
	AccountId    int64
	HisData      []uint32
	Activities   []*Activity
	PreRewardTps []uint32
}

func NewAccountActivity(accountId int64) *AccountActivity {
	return &AccountActivity{
		AccountId:  accountId,
		Activities: make([]*Activity, 0, 0),
	}
}

type Activity struct {
	ActId     uint32
	ParentId  uint32
	ActType   uint32
	StartTime uint32
	EndTime   uint32
	ActDatas  []*ActivityData
}

type ActivityData struct {
	SubActId   uint32
	Value1     uint32
	Value2     int64
	Value3     int64
	Value4     int64
	State      uint32
	UpdateTime uint32
}

func (a *ActivityData) GetTaskId() uint32 {
	return a.SubActId
}

func (a *ActivityData) GetTaskValue() uint32 {
	return a.Value1
}

func (a *ActivityData) SetTaskValue(value uint32) {
	a.Value1 = value
}

func (a *ActivityData) AddTaskValue(add uint32) {
	a.Value1 += add
}

func (a *ActivityData) SetTaskCompleteTime(time uint32) {
	a.UpdateTime = time
}

func (a *ActivityData) GetTaskCompleteTime() uint32 {
	return a.UpdateTime
}

func (a *ActivityData) GetTaskState() publicconst.TaskState {
	return publicconst.TaskState(a.State)
}

func (a *ActivityData) SetTaskState(state publicconst.TaskState) {
	a.State = uint32(state)
}

func (a *ActivityData) SetExtraPara(value uint32) {
	a.Value2 = int64(value)
}

func (a *ActivityData) GetExtraPara() uint32 {
	return uint32(a.Value2)
}

func NewActivity(actId, actType, startTime, endTime uint32) *Activity {
	return &Activity{
		ActId:     actId,
		ActType:   actType,
		StartTime: startTime,
		EndTime:   endTime,
		ActDatas:  make([]*ActivityData, 0, 0),
	}
}

func NewActivityData(subActId uint32, value1 uint32, value2, value3, value4 int64) *ActivityData {
	return &ActivityData{
		SubActId:   subActId,
		Value1:     value1,
		Value2:     value2,
		Value3:     value3,
		Value4:     value4,
		UpdateTime: tools.GetCurTime(),
	}
}

// 活动是否进行中
func (a *AccountActivity) IsActivityInProgress(activityId uint32) bool {
	for i := 0; i < len(a.Activities); i++ {
		if a.Activities[i].ActId == activityId {
			for j := 0; j < len(a.Activities[i].ActDatas); j++ {
				if a.Activities[i].ActDatas[j].State != uint32(publicconst.TASK_DONE) {
					return true
				}
			}
		}
	}
	return false
}
