package service

import (
	"fmt"
	"msg"

	"github.com/zy/game_data/template"

	"gameserver/internal/game/event"
	"gameserver/internal/game/model"
	"gameserver/internal/publicconst"
)

// ListenAddItemEvent 添加道具事件
func ListenAddItemEvent(e event.IEvent) {
	addItemEvent := e.(*event.AddItemEvent)
	if addItemEvent == nil {
		return
	}
	//accountId := addItemEvent.PlayerData.GetAccountId()

	if addItemEvent.SendClient {
		updateClientItemChange(addItemEvent.PlayerData, addItemEvent.ItemId)
	}

	UpdateTask(addItemEvent.PlayerData, true, publicconst.TASK_COND_ADD_ITEM, addItemEvent.ItemId, addItemEvent.Num)
	itemConfig := template.GetItemTemplate().GetItem(addItemEvent.ItemId)
	if itemConfig != nil && itemConfig.BigType == 16 {
		UpdateTask(addItemEvent.PlayerData, true,
			publicconst.TASK_COND_GET_WEAPON_NUM, addItemEvent.Num)
		processHistoryData(addItemEvent.PlayerData, publicconst.TASK_COND_GET_EQUIP_NUM, 0, addItemEvent.Num)
	}
	processHistoryData(addItemEvent.PlayerData, publicconst.TASK_COND_ADD_ITEM, addItemEvent.ItemId, addItemEvent.Num)

	// if itemConfig != nil && itemConfig.BigType == 20 {
	// 	PutPetEgg(addItemEvent.PlayerData, addItemEvent.ItemId)
	// }

	// // 记录日志
	// itemlog := model.NewItemLog(accountId, addItemEvent.ItemId, int32(addItemEvent.Num), addItemEvent.CurNum, int32(addItemEvent.ItemSrc), "")
	// dao.ItemLogDao.AddItemLog(itemlog)
}

// ListenCostItemEvent 扣除道具事件
func ListenCostItemEvent(e event.IEvent) {
	costItemEvent := e.(*event.CostItemEvent)
	if costItemEvent == nil {
		return
	}

	// accountId := costItemEvent.PlayerData.GetAccountId()
	if costItemEvent.SendClient {
		updateClientItemChange(costItemEvent.PlayerData, costItemEvent.ItemId)
	}

	UpdateTask(costItemEvent.PlayerData, true, publicconst.TASK_COND_COST_ITEM, costItemEvent.ItemId, costItemEvent.Num)

	// // 记录日志
	// itemlog := model.NewItemLog(accountId, costItemEvent.ItemId, int32(-costItemEvent.Num), costItemEvent.CurNum, int32(costItemEvent.ItemSrc), "")
	// dao.ItemLogDao.AddItemLog(itemlog)
}

// ListenNotifyClientItemEventEvent
func ListenNotifyClientItemEventEvent(e event.IEvent) {
	event := e.(*event.NotifyClientItemEvent)
	if event == nil {
		return
	}
	if len(event.ItemIds) == 0 {
		return
	}
	updateClientItemsChange(event.PlayerData.GetUserId(), event.ItemIds)
}

// ListenGetTaskRewardEvent 领奖事件
func ListenGetTaskRewardEvent(e event.IEvent) {
	//getTaskRewardEvent := e.(*event.GetTaskRewardEvent)
	//if getTaskRewardEvent == nil {
	//	return
	//}
	//
	//jTask := template.GetTaskTemplate().GetTask(getTaskRewardEvent.TaskId)
	//if jTask == nil {
	//	return
	//}
	//
	//// 刷新活跃度
	////ServMgr.GetTaskService().RefreshTaskActive(getTaskRewardEvent.PlayerData, jTask)
	//
	//// 接下一级任务
	//if jTask.Data.NextId != 0 {
	//	// 主线任务完成需要删除主线任务
	//	if jTask.GetTaskType() == publicconst.MAIN_TASK {
	//		//	ServMgr.GetTaskService().DelTask(getTaskRewardEvent.PlayerData, getTaskRewardEvent.TaskId)
	//	}
	//
	//	if jNextTask := template.GetTaskTemplate().GetTask(jTask.Data.NextId); jNextTask != nil {
	//		if jNextTask.Data.UsePrevValue == 1 {
	//			//	ServMgr.GetTaskService().AddTask(getTaskRewardEvent.PlayerData, jNextTask.Data.Id, jTask.Data.TaskValue)
	//		} else {
	//			//	ServMgr.GetTaskService().AddTask(getTaskRewardEvent.PlayerData, jNextTask.Data.Id, 0)
	//		}
	//	}
	//}
}

func ListenPassMissionEvent(e event.IEvent) {
	passMissionEvent := e.(*event.PassMissionEvent)
	if passMissionEvent == nil {
		return
	}

	// 通关添加任务
	tasks := template.GetTaskTemplate().GetMissionTasks(uint32(passMissionEvent.MissionId))
	var notifyTasks []*model.Task
	for i := 0; i < len(tasks); i++ {
		if CanAddTask(passMissionEvent.PlayerData, tasks[i]) {
			if task := AddTask(passMissionEvent.PlayerData, tasks[i], 0); task != nil {
				notifyTasks = append(notifyTasks, task)
			}
		}
	}
	NotifyClientTaskChange(passMissionEvent.PlayerData, notifyTasks)

	// 更新通关进度
	UpdateTask(passMissionEvent.PlayerData, true,
		publicconst.TASK_COND_PASS_MISSION_ID, uint32(passMissionEvent.MissionId))

	UpdateTask(passMissionEvent.PlayerData, true,
		publicconst.TASK_COND_ANY_ONE_MISSION, uint32(passMissionEvent.MissionId))

	// // 添加关卡流失统计 TODO
	// dao.UserStaticDao.AddLossMission(model.NewLossMission(tools.GetStaticTime(tools.GetCurTime()),
	// 	uint32(passMissionEvent.PlayerData.GetAccountId()), int(passMissionEvent.MissionId)))

	// 刷新活动
	RefreshActivity(passMissionEvent.PlayerData, true)

	UpdateFunctionPreview(passMissionEvent.PlayerData, msg.ConditionType_Condition_Pass_Mission)
}

func ListenFirstPassMissionEvent(e event.IEvent) {
	firstPassEvent := e.(*event.FirstPassMissionEvent)
	if firstPassEvent == nil {
		return
	}

	reward := firstPassEvent.Args.(*SimpleMissionReward)
	if _, info := GetPlayerSimpleInfo(reward.Uid); info != nil {
		reward.HeadFrame = info.HeadFrame
		reward.Nick = info.Name
		reward.Head = info.Head
	}

	notifyMsg := &msg.NotifyFirstPassMission{
		Data: ToProtocolRankMissionReward(reward),
	}
	BoadCastMsg(notifyMsg)

	bannerMsg := &msg.InterNotifyBanner{}
	bannerMsg.BannerType = uint32(msg.BannerType_Banner_Game)
	bannerMsg.Params = append(bannerMsg.Params, reward.Nick)
	bannerMsg.Params = append(bannerMsg.Params, fmt.Sprintf("%s_%d", template.T_Stage_FILE, reward.MissionId))

	if reward.Tp == int(msg.BattleType_Battle_Main) {
		bannerMsg.Content = fmt.Sprintf("%d", template.MainMissionFirstPassTickerID)
	} else if reward.Tp == int(msg.BattleType_Battle_Challenge) {
		bannerMsg.Content = fmt.Sprintf("%d", template.ChallengeMissionFirstPassTickerID)
	}

	BoadCastMsg(bannerMsg)
}

func ListenAddShipEvent(e event.IEvent) {
	event := e.(*event.AddShipEvent)
	if event == nil {
		return
	}

	for i := 0; i < len(event.ShipIds); i++ {
		UpdateShipTreasure(event.PlayerData, event.ShipIds[i], true)
		UpdateShipPoker(event.PlayerData, event.ShipIds[i], true)
	}
	UpdateTask(event.PlayerData, true, publicconst.TASK_COND_SHIP_RARITY_NUM)
}

func ListenShipChangeEvent(e event.IEvent) {
	event := e.(*event.ShipChangeEvent)
	if event == nil {
		return
	}
	switch event.AttrType {
	case publicconst.Ship_Level:
		UpdateTask(event.PlayerData, true,
			publicconst.TASK_COND_SHIP_LEVEL, event.ShipId, event.NewData)
	case publicconst.Ship_Star_Level:
		UpdateTask(event.PlayerData, true,
			publicconst.TASK_COND_SHIP_STAR_NUM)
		UpdateTask(event.PlayerData, true,
			publicconst.TASK_COND_SPECIFY_SHIP_STAR, event.ShipId, event.NewData)

		UpdateShipTreasure(event.PlayerData, event.ShipId, true)
		UpdateShipPoker(event.PlayerData, event.ShipId, true)

		UpdateActivity(event.PlayerData, publicconst.TASK_COND_SPECIFY_SHIP_STAR, event.ShipId, event.NewData)
	}
}

// ListenNotifyClientShipChangeEvent 通知客户端更新机甲
func ListenNotifyClientShipChangeEvent(e event.IEvent) {
	event := e.(*event.NotifyClientShipChangeEvent)
	if event == nil {
		return
	}
	UpdateClientShipChange(event.PlayerData, event.ShipIds)
}

// ListenAddRoleEvent 添加驾驶员
func ListenAddRoleEvent(e event.IEvent) {

}

func ListenRoleChangeEvent(e event.IEvent) {
	event := e.(*event.RoleChangeEvent)
	if event == nil {
		return
	}
	switch event.AttrType {
	case publicconst.Role_Star_Level:
	case publicconst.Role_Favor_Level:
	}
}

func ListenClientRoleChangeEvent(e event.IEvent) {
	// event := e.(*event.NotifyClientRoleChangeEvent)
	// if event == nil {
	// 	return
	// }
	// ServMgr.GetRoleService().UpdateClientRoleChange(event.PlayerData, event.RoleIds)
}

func ListenLevelChangeEvent(e event.IEvent) {
	event := e.(*event.LevelChangeEvent)
	if event == nil {
		return
	}

	// 通知客户端账号变化
	NotifyAccountChange(event.PlayerData)

	// 获取等级任务
	var notifyTasks []*model.Task
	for i := event.OldLevel + 1; i <= event.NewLevel; i++ {
		tasks := template.GetTaskTemplate().GetLevelTasks(i)
		for k := 0; k < len(tasks); k++ {
			if CanAddTask(event.PlayerData, tasks[k]) {
				if task := AddTask(event.PlayerData, tasks[k], 0); task != nil {
					notifyTasks = append(notifyTasks, task)
				}
			}
		}
	}
	NotifyClientTaskChange(event.PlayerData, notifyTasks)

	// 更新玩家等级任务进度
	UpdateTask(event.PlayerData, true,
		publicconst.TASK_COND_SHIP_LEVEL, event.PlayerData.UserData.Level)

	// 添加等级流失数据 TODO
	// lossLevel := model.NewLossLevel(tools.GetStaticTime(tools.GetCurTime()),
	// 	uint32(event.PlayerData.GetAccountId()), event.PlayerData.UserData.Level, event.PlayerData.UserData.ChannelId)
	// dao.UserStaticDao.AddLossLevel(lossLevel)

	// 刷新活动
	RefreshActivity(event.PlayerData, true)
}

func ListenNotifyEquipChangeEvent(e event.IEvent) {
	event := e.(*event.NotifyEquipChangeEvent)
	if event == nil {
		return
	}
	updateClientEquipChange(event.PlayerData, event.EquipIds)
}

// func ListenMailEvent(e event.IEvent) {
// 	mailEvent := e.(*event.MailEvent)
// 	if mailEvent == nil {
// 		return
// 	}
// 	notifyClientMail(mailEvent.PlayerData, mailEvent.MailId, mailEvent.IsDelete)
// }

func ListenWeaponUpgradeEvent(e event.IEvent) {
	upgradeEvent := e.(*event.WeaponUpgradeEvent)
	if upgradeEvent == nil {
		return
	}

	var ids []uint32
	ids = append(ids, upgradeEvent.WeaponId)
	UpdateWeaponTreasure(upgradeEvent.PlayerData, upgradeEvent.WeaponId)
	UpdateWeaponPoker(upgradeEvent.PlayerData, upgradeEvent.WeaponId)
	UpdateTask(upgradeEvent.PlayerData, true, publicconst.TASK_COND_ANY_WEAPON_LEVEL)
	UpdateTask(upgradeEvent.PlayerData, true, publicconst.TASK_COND_SPECIFY_WEAPON_LEVEL, upgradeEvent.WeaponId, upgradeEvent.NewLevel)
	UpdateTask(upgradeEvent.PlayerData, true, publicconst.TASK_COND_ALL_WEAPON_MIN_LEVEL)

	UpdateTask(upgradeEvent.PlayerData, true, publicconst.TASK_COND_UPGRADE_WEAPON, 1)
	processHistoryData(upgradeEvent.PlayerData, publicconst.TASK_COND_UPGRADE_WEAPON, 0, 1)
}

func ListenWeaponLibUpgradeEvent(e event.IEvent) {
	libUpgradeEvent := e.(*event.WeaponLibUpgradeEvent)
	if libUpgradeEvent == nil {
		return
	}
	UpdateTask(libUpgradeEvent.PlayerData, true, publicconst.TASK_COND_WEAPON_LIB_LEVEL, libUpgradeEvent.NewLevel)
}
