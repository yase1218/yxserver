package service

import (
	"gameserver/internal/game/builder"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"msg"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

/*
 * LuckSaleTaskReward
 *  @Description: 任务完成
 *  @param player
 *  @param req
 *  @param ack
 *  @return error
 */
func LuckSaleTaskReward(p *player.Player, req *msg.LuckSaleTaskRewardReq, ack *msg.LuckSaleTaskRewardAck) {
	taskCfg := template.GetTaskTemplate().GetTask(req.GetTaskId())
	if taskCfg == nil {
		log.Error("task cfg nil", zap.Uint32("id", req.GetTaskId()))
		ack.Err = msg.ErrCode_CONFIG_NIL
		return
	}

	var task *model.LuckSaleTaskUnit
	for _, v := range p.UserData.LuckSale.Task {
		if v.TaskId == req.GetTaskId() {
			task = v
			break
		}
	}

	if task == nil {
		log.Error("task nil", zap.Uint32("taskId", req.GetTaskId()), zap.Uint64("accountId", p.GetUserId()))
		ack.Err = msg.ErrCode_SYSTEM_ERROR
		return
	}

	if task.State == uint32(publicconst.TASK_DONE) || task.Reward {
		log.Error("task has get reward", zap.Uint32("taskId", req.GetTaskId()), zap.Uint64("accountId", p.GetUserId()))
		ack.Err = msg.ErrCode_TASK_HAS_GET_REWARD
		return
	}
	if task.State == uint32(publicconst.TASK_ACCEPT) {
		log.Error("task not complete", zap.Uint32("taskId", req.GetTaskId()), zap.Uint64("accountId", p.GetUserId()))
		ack.Err = msg.ErrCode_TASK_NOT_COMPLETE
		return
	}

	notifyItems := make([]uint32, 0, 1)
	addItem := template.GetSystemItemTemplate().LuckSaleItem
	AddItem(p.GetUserId(), addItem, 1, publicconst.LuckSaleTaskReward, false)
	notifyItems = append(notifyItems, addItem)
	updateClientItemsChange(p.GetUserId(), notifyItems)

	task.State = uint32(publicconst.TASK_DONE)
	task.Reward = true

	if taskCfg.NextTaskId != 0 {
		newTaskCfg := template.GetTaskTemplate().GetTask(taskCfg.NextTaskId)
		if newTaskCfg == nil {
			log.Error("new task cfg nil", zap.Uint32("nextTaskId", taskCfg.NextTaskId))
			ack.Err = msg.ErrCode_CONFIG_NIL
			return
		}
		newDbTask := &model.LuckSaleTaskUnit{
			TaskId: newTaskCfg.Data.Id,
		}
		initTaskValue(p, newDbTask, publicconst.TaskCond(newTaskCfg.Data.TaskCondition))
		p.UserData.LuckSale.Task = append(p.UserData.LuckSale.Task, newDbTask)
		ack.NewTask = builder.BuildLuckSaleTask(newDbTask)
	}

	p.SaveLuckSale()

	ack.TaskId = req.GetTaskId()
}

/*
 * LuckSaleExtract
 *  @Description: 抽奖
 *  @param player
 *  @param req
 *  @param ack
 *  @return error
 */
func LuckSaleExtract(p *player.Player, req *msg.LuckSaleExtractReq, ack *msg.LuckSaleExtractAck) {
	costItemId := template.GetSystemItemTemplate().LuckSaleItem

	curJackpot := p.UserData.LuckSale.Jackpot
	luckSaleGroupId := template.GetSystemItemTemplate().LuckSaleGroupId
	curLuckSale := p.UserData.LuckSale.Data[curJackpot]
	randCfg := template.GetNoRepeatLootPoolTemplate().Rand(luckSaleGroupId, uint32(curJackpot), curLuckSale.Times, curLuckSale.Ids)
	if randCfg == nil {
		log.Error("rand err", zap.Uint64("accountId", p.GetUserId()))
		ack.Err = msg.ErrCode_SYSTEM_ERROR
		return
	}

	var notifyItems []uint32

	if res := CostItem(p.GetUserId(), costItemId, 1, publicconst.LuckSaleExtract, false); res != msg.ErrCode_SUCC {
		ack.Err = res
		return
	}
	notifyItems = append(notifyItems, costItemId)

	// add item
	if len(randCfg.Reward) > 0 {
		for _, item := range randCfg.Reward {
			addItems := AddItem(p.GetUserId(), item.ItemId, int32(item.ItemNum), publicconst.LuckSaleExtract, false)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		}
		ack.RewardItem = append(ack.RewardItem, TemplateItemToProtocolItems(randCfg.Reward)...)
	}
	updateClientItemsChange(p.GetUserId(), notifyItems)

	curLuckSale.Ids = append(curLuckSale.Ids, randCfg.Id)
	curLuckSale.Times++

	times := template.GetNoRepeatLootPoolTemplate().GetTimes(luckSaleGroupId, uint32(curJackpot))
	if curLuckSale.Times >= times {
		p.UserData.LuckSale.Jackpot++
		p.UserData.LuckSale.Data[p.UserData.LuckSale.Jackpot] = model.NewLuckSaleUnit()

		times := template.GetNoRepeatLootPoolTemplate().GetTimes(luckSaleGroupId, uint32(p.UserData.LuckSale.Jackpot))
		if times == 0 {
			p.UserData.LuckSale.Jackpot = -1
		}
	}

	p.SaveLuckSale()

	ack.Id = randCfg.Id
	ack.Jackpot = int32(p.UserData.LuckSale.Jackpot)
	if p.UserData.LuckSale.Jackpot == -1 {
		ack.Times = 0
	} else {
		ack.Times = p.UserData.LuckSale.Data[p.UserData.LuckSale.Jackpot].Times
	}

	//ServMgr.GetCommonService().AddStaticsData(p, publicconst.Statics_Add_Atlas, fmt.Sprintf("%d", randCfg.Id))
}

func UpdateLuckSaleTask(p *player.Player, cond uint32, args ...uint32) {
	pbTask := make([]*msg.Task, 0)
	updateTasks := make([]*model.LuckSaleTaskUnit, 0)

	for _, v := range p.UserData.LuckSale.Task {
		if v.State == uint32(publicconst.TASK_COMPLETE) ||
			v.State == uint32(publicconst.TASK_DONE) {
			continue
		}

		taskCfg := template.GetTaskTemplate().GetTask(v.TaskId)
		if taskCfg == nil {
			log.Error("task cfg nil", zap.Uint32("taskId", v.TaskId))
			continue
		}

		if taskCfg.Data.TaskCondition != cond {
			continue
		}
		if refreshTaskValue(p, v, publicconst.TaskCond(cond), taskCfg.Data.MaxValue, taskCfg.Data.Effect1, args...) {
			pbTask = append(pbTask, builder.BuildLuckSaleTask(v))
			updateTasks = append(updateTasks, v)
		}
	}

	if len(pbTask) > 0 {
		p.SaveLuckSale()

		pbMsg := &msg.LuckSaleTaskNtf{
			Tasks: pbTask,
		}
		p.SendNotify(pbMsg)
	}
}
