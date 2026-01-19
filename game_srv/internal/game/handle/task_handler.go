package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"gameserver/internal/publicconst"
	"msg"
)

// RequestLoadTaskByTypeHandle 通过类型加载任务
func RequestLoadTaskByTypeHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestLoadTaskByType)
	tasks := service.GetTasksByType(p, publicconst.TaskType(req.ReqType))

	retMsg := &msg.ResponseLoadTaskByType{
		Result:  msg.ErrCode_SUCC,
		ReqType: req.ReqType,
		Tasks:   service.ToProtocolTasks(tasks),
	}
	switch req.ReqType {
	case msg.TaskType_Daily_Task:
		retMsg.ActInfo = &msg.ActiveInfo{Value: p.UserData.Task.DailyActiveValue, State: p.UserData.Task.DailyActiveReward}
	case msg.TaskType_Weekly_Task:
		retMsg.ActInfo = &msg.ActiveInfo{Value: p.UserData.Task.WeeklyActiveValue, State: p.UserData.Task.WeeklyActiveReward}
	case msg.TaskType_Alliance_Weekly_Task:
		retMsg.ActInfo = &msg.ActiveInfo{Value: p.UserData.Task.AllianceWeeklyActiveValue, State: p.UserData.Task.AllianceWeeklyActiveReward}
	default:

	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestGetTaskRewardHandle 获取任务奖励
func RequestGetTaskRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestGetTaskReward)
	err := service.GetTaskReward(p, req.TaskId)
	retMsg := &msg.ResponseGetTaskReward{
		Result: err,
		TaskId: req.TaskId,
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestBatchGetTaskRewardHandle 批量获取任务奖励
func RequestBatchGetTaskRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestBatchGetTaskReward)
	//err, items := service.ServMgr.GetTaskService().BatchGetTaskReward(playerData, publicconst.TaskType(req.ReqType))
	err, items := service.GetAllTaskReward(p, publicconst.TaskType(req.ReqType))
	retMsg := &msg.ResponseBatchGetTaskReward{
		Result:   err,
		ReqType:  req.ReqType,
		GetItems: service.TemplateItemToProtocolItems(items),
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestGetActiveRewardHandle 获得活跃度奖励
func RequestGetActiveRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestGetActiveReward)
	err, activeReward := service.GetTaskActiveReward(p, publicconst.TaskType(req.ReqType), req.Pos)
	retMsg := &msg.ResponseGetActiveReward{
		Result:  err,
		ReqType: req.ReqType,
		State:   activeReward,
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}
