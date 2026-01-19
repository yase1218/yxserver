package model

import "gameserver/internal/publicconst"

type ITask interface {
	GetTaskId() uint32
	GetTaskValue() uint32
	SetTaskValue(value uint32)
	AddTaskValue(add uint32)
	SetTaskCompleteTime(time uint32)
	GetTaskCompleteTime() uint32
	GetTaskState() publicconst.TaskState
	SetTaskState(state publicconst.TaskState)
	SetExtraPara(value uint32)
	GetExtraPara() uint32
}
