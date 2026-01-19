package service

import (
	"kernel/tools"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"

	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
)

type TaskCheckFunc func(*player.Player, model.ITask, uint32, []uint32, ...uint32) bool

type TaskInitFunc func(*player.Player, model.ITask) bool

var (
	taskCheckFuncMap map[publicconst.TaskCond]TaskCheckFunc
	taskInitMap      map[publicconst.TaskCond]TaskInitFunc
)

func InitTaskCheck() {
	taskCheckFuncMap = make(map[publicconst.TaskCond]TaskCheckFunc)
	taskCheckFuncMap[publicconst.TASK_COND_PASS_MISSION_ID] = compareArg0
	taskCheckFuncMap[publicconst.TASK_COND_SHIP_LEVEL] = setArg0
	taskCheckFuncMap[publicconst.TASK_COND_SHIP_RARITY_NUM] = setShipRarityNum
	taskCheckFuncMap[publicconst.TASK_COND_SHIP_STAR_NUM] = setShipStarNum
	taskCheckFuncMap[publicconst.TASK_COND_SPECIFY_SHIP_STAR] = compareArg0SetArg1
	taskCheckFuncMap[publicconst.TASK_COND_SHIP_SWITCH] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_GET_EQUIP_NUM] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_UPGRADE_EQUIP] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_ANY_EQUIP_LEVEL] = setAnyEquipLevel
	taskCheckFuncMap[publicconst.TASK_COND_ALL_EQUIP_POS_LEVEL] = setAllEquipPosLevel
	taskCheckFuncMap[publicconst.TASK_COND_COMPOSE_EQUIP] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_EQUIP_RARITY_NUM] = setEquipRarityNum
	taskCheckFuncMap[publicconst.TASK_COND_UPGRADE_WEAPON] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_ANY_WEAPON_LEVEL] = setAnyWeaponLevel
	taskCheckFuncMap[publicconst.TASK_COND_COMPLETE_DAILY_TASK_COUNT] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_COMPLETE_WEEKLY_TASK_COUNT] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_DAILY_ACTIVE_SCORE] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_WEEKLY_ACTIVE_SCORE] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_GET_ON_HOOK_REWARD] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_QUICK_ON_HOOK] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_BUY_AP] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_KILL_MONSTER_NUM] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_KILL_ELITE_MONSTER_BOSS_NUM] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_BATTLE_MISSIONTYPE] = compareArg0AddArg1
	taskCheckFuncMap[publicconst.TASK_COND_ADD_ITEM] = addItem
	taskCheckFuncMap[publicconst.TASK_COND_COST_ITEM] = compareArg0AddArg1
	taskCheckFuncMap[publicconst.TASK_COND_LOGIN] = compareLogin
	taskCheckFuncMap[publicconst.TASK_COND_SPECIFY_WEAPON_LEVEL] = compareArg0SetArg1
	taskCheckFuncMap[publicconst.TASK_COND_ALL_WEAPON_MIN_LEVEL] = setAllWeaponMinLevel
	taskCheckFuncMap[publicconst.TASK_COND_WEAPON_LIB_LEVEL] = setArg0
	taskCheckFuncMap[publicconst.TASK_COND_GET_WEAPON_NUM] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_PUT_ON_EQUIP] = setPutOnEquip
	taskCheckFuncMap[publicconst.TASK_COND_CLICK_UPDATE_NICK] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_ANY_BATTLE] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_GET_DAILY_BOX] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_LOTTERY] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_ANY_ONE_MISSION] = compareAnyMission
	taskCheckFuncMap[publicconst.TASK_COND_GET_CHIP_NUM] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_POKER_NUM] = comparePoker
	//taskCheckFuncMap[publicconst.TASK_COND_PET_ADVENTURE] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_GET_PET] = compareArg0AddArg1 // 激活XX个XX品质库鲁兽
	taskCheckFuncMap[publicconst.TASK_COND_PET_PART_CHARM] = setArg0
	taskCheckFuncMap[publicconst.TASK_COND_ALLIANCE_JOIN] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_ALLIANCE_BOSS_FIGHT] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_ALLIANCE_FINISH_TASK] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_PEAK_FIGHT_PK] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_PEAK_FIGHT_RANK] = compareArg0
	taskCheckFuncMap[publicconst.TASK_COND_EXPLORE_USE_CARD] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_EXPLORE_OCCUPY_BUILD] = compareArg0
	taskCheckFuncMap[publicconst.TASK_COND_EXPLORE_PROGRESS] = compareArg0
	taskCheckFuncMap[publicconst.TASK_COND_MISSION_BOX_REWARD] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_FIRST_PASS_CHALLENGE_MISSION] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_RESUME_DIAMOND_BUY_SHOP] = addArg0
	taskCheckFuncMap[publicconst.TASK_COND_PASS_TYPE_MISSION] = compareArg0AddArg1 // 完成XX关卡类型XX次
	taskCheckFuncMap[publicconst.TASK_COND_ACTIVE_SHIP] = compareArg0AddArg1       // 激活XX个XX品质库鲁 兑换
	taskCheckFuncMap[publicconst.TASK_COND_PUT_ON_DISK] = addArg0                  // 镶嵌XX个磁盘
	taskCheckFuncMap[publicconst.TASK_COND_ADD_DISK] = compareArg0AddArg1          // 获得XX个磁盘
	taskCheckFuncMap[publicconst.TASK_COND_ADD_RARITY_DISK] = compareArg01AddArg2  // 获得X个XX品质磁盘
	taskCheckFuncMap[publicconst.TASK_COND_PET_LEVEL] = setArg0                    // 库鲁兽等级达到XX级
	taskCheckFuncMap[publicconst.TASK_COND_RESOURCES_PASS] = addArg0               // 进行XX次素材本
	taskCheckFuncMap[publicconst.TASK_COND_SHIP_LOTTERY] = addArg0                 // 累计进行XX次库鲁招募
	taskCheckFuncMap[publicconst.TASK_COND_DISK_BIND_BOX_LOTTERY] = addArg0        // 累计进行XX次磁盘盲盒抽奖
	taskCheckFuncMap[publicconst.TASK_COND_SHIP_BEAST_EGG_LOTTERY] = addArg0       // 累计进行XX次库鲁兽扭蛋
	taskCheckFuncMap[publicconst.TASK_COND_REWARD_EQUIP] = addArg0                 // 装备本累计领取奖励1次

	taskInitMap = make(map[publicconst.TaskCond]TaskInitFunc)
	taskInitMap[publicconst.TASK_COND_SHIP_LEVEL] = initLevel
	taskInitMap[publicconst.TASK_COND_PASS_MISSION_ID] = initPassMission
	taskInitMap[publicconst.TASK_COND_SHIP_RARITY_NUM] = initShipRarityNum
	taskInitMap[publicconst.TASK_COND_SHIP_STAR_NUM] = initShipStarNum
	taskInitMap[publicconst.TASK_COND_SPECIFY_SHIP_STAR] = initPlayerShipStarLevel
	taskInitMap[publicconst.TASK_COND_ANY_EQUIP_LEVEL] = initAnyEquipLevel
	taskInitMap[publicconst.TASK_COND_ALL_EQUIP_POS_LEVEL] = initAllEquipPosLevel
	taskInitMap[publicconst.TASK_COND_EQUIP_RARITY_NUM] = initEquipRarityNum
	taskInitMap[publicconst.TASK_COND_ANY_WEAPON_LEVEL] = initAnyWeaponLevel
	taskInitMap[publicconst.TASK_COND_SPECIFY_WEAPON_LEVEL] = initSpecifyWeaponLevel
	taskInitMap[publicconst.TASK_COND_ALL_WEAPON_MIN_LEVEL] = initAllWeaponMinLevel
	taskInitMap[publicconst.TASK_COND_WEAPON_LIB_LEVEL] = initWeaponLibLevel
	taskInitMap[publicconst.TASK_COND_PUT_ON_EQUIP] = initPutOnEquip
	taskInitMap[publicconst.TASK_COND_GET_PET] = initPet
	taskInitMap[publicconst.TASK_COND_PET_PART_CHARM] = initPetPartCharm
	taskInitMap[publicconst.TASK_COND_ADD_ITEM] = initAddItem
	taskInitMap[publicconst.TASK_COND_PEAK_FIGHT_RANK] = initPeakFightRank
	taskInitMap[publicconst.TASK_COND_EXPLORE_OCCUPY_BUILD] = initExploreBuild
	taskInitMap[publicconst.TASK_COND_UPGRADE_WEAPON] = initUseHistoryData
	taskInitMap[publicconst.TASK_COND_MISSION_BOX_REWARD] = initUseMissionBoxReward // 领取X次通关宝箱
}

// configArgs:Effect1
func refreshTaskValue(p *player.Player, data model.ITask, cond publicconst.TaskCond, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	if f, ok := taskCheckFuncMap[cond]; ok {
		return f(p, data, maxValue, configArgs, args...)
	}
	return false
}

func initTaskValue(p *player.Player, data model.ITask, cond publicconst.TaskCond) bool {
	if f, ok := taskInitMap[cond]; ok {
		return f(p, data)
	}
	return false
}

func setArg0(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	if len(args) == 0 {
		return false
	}

	if data.GetTaskValue() < args[0] {
		data.SetTaskValue(args[0])
		if data.GetTaskValue() >= maxValue {
			data.SetTaskValue(maxValue)
			data.SetTaskState(publicconst.TASK_COMPLETE)
		}
		return true
	}
	return false
}

func compareArg0SetArg1(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	if len(args) != 2 {
		return false
	}

	if len(configArgs) == 1 {
		if (configArgs[0] == 0) || (args[0] == configArgs[0]) {
			if args[1] >= maxValue {
				data.SetTaskValue(maxValue)
				data.SetTaskState(publicconst.TASK_COMPLETE)
			} else {
				data.SetTaskValue(args[1])
			}
			return true
		}
	}
	return false
}

func compareArg0AddArg1(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	if len(args) != 2 {
		return false
	}

	if len(configArgs) == 1 {
		if (configArgs[0] == 0) || (args[0] == configArgs[0]) {
			data.AddTaskValue(args[1])
			if data.GetTaskValue() >= maxValue {
				data.SetTaskValue(maxValue)
				data.SetTaskState(publicconst.TASK_COMPLETE)
			}
			return true
		}
	}
	return false
}

func compareArg01AddArg2(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	if len(args) != 3 || len(configArgs) != 2 {
		return false
	}
	if args[0] == configArgs[0] && configArgs[1] <= args[1] {
		data.AddTaskValue(args[2])
		if data.GetTaskValue() >= maxValue {
			data.SetTaskValue(maxValue)
			data.SetTaskState(publicconst.TASK_COMPLETE)
		}
		return true
	}
	return false
}

func addItem(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	if len(args) != 2 {
		return false
	}

	if len(configArgs) == 1 {
		if (configArgs[0] == 0) || (args[0] == configArgs[0]) {
			data.AddTaskValue(args[1])
			if data.GetTaskValue() >= maxValue {
				data.SetTaskValue(maxValue)
				data.SetTaskState(publicconst.TASK_COMPLETE)
			}
			if jTask := template.GetTaskTemplate().GetTask(data.GetTaskId()); jTask != nil && jTask.Data.UseHisData == 1 {
				if historyData := GetHistoryData(p, publicconst.TASK_COND_ADD_ITEM, args[0]); historyData == nil {
					historyData := model.NewTaskHistroyData(uint32(publicconst.TASK_COND_ADD_ITEM),
						0, args[0])
					addHistoryData(p, historyData)
				}
			}
			return true
		}
	}
	return false
}

func compareLogin(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	curTime := tools.GetCurTime()
	if curTime >= data.GetExtraPara() {
		data.SetExtraPara(tools.GetDailyRefreshTime())
		data.AddTaskValue(1)
		if data.GetTaskValue() >= maxValue {
			data.SetTaskValue(maxValue)
			data.SetTaskState(publicconst.TASK_COMPLETE)
		}
		processHistoryData(p, publicconst.TASK_COND_LOGIN, 0, 1)
		return true
	}
	return false
}

// compareArg0 比较参数0 相等则任务完成
func compareArg0(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	if len(args) != 1 {
		return false
	}
	if len(configArgs) == 0 { // 不需要比较第一个参数
		data.AddTaskValue(1)
		data.SetTaskState(publicconst.TASK_COMPLETE)
		return true
	} else {
		if (configArgs[0] == 0) || (args[0] == configArgs[0]) {
			data.AddTaskValue(1)
			data.SetTaskState(publicconst.TASK_COMPLETE)
			return true
		}
	}
	return false
}

func addArg0(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	if len(args) == 0 {
		return false
	}

	data.AddTaskValue(args[0])
	if data.GetTaskValue() >= maxValue {
		data.SetTaskValue(maxValue)
		data.SetTaskState(publicconst.TASK_COMPLETE)
	}
	return true
}

func setShipRarityNum(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	return initShipRarityNum(p, data)
}

func setShipStarNum(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	return initShipStarNum(p, data)
}

func setAnyEquipLevel(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	return initAnyEquipLevel(p, data)
}

func setAllEquipPosLevel(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	return initAllEquipPosLevel(p, data)
}

func setEquipRarityNum(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	return initEquipRarityNum(p, data)
}

func setAnyWeaponLevel(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	return initAnyWeaponLevel(p, data)
}

func setAllWeaponMinLevel(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	return initAllWeaponMinLevel(p, data)
}

func setPutOnEquip(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	return initPutOnEquip(p, data)
}

func compareAnyMission(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	if len(args) != 1 || len(configArgs) == 0 {
		return false
	}

	if tools.ListContain(configArgs, args[0]) {
		data.SetTaskValue(1)
		data.SetTaskState(publicconst.TASK_COMPLETE)
		return true
	}
	return false
}

func comparePoker(p *player.Player, data model.ITask, maxValue uint32, configArgs []uint32, args ...uint32) bool {
	if len(configArgs) == 0 || len(args) == 0 {
		return false
	}
	for i := 0; i < len(args); i += 2 {
		if args[i] == configArgs[0] {
			data.AddTaskValue(args[i+1])
			if data.GetTaskValue() >= maxValue {
				data.SetTaskValue(maxValue)
				data.SetTaskState(publicconst.TASK_COMPLETE)
			}
			return true
		}
	}
	return false
}

func initTaskLevel(p *player.Player, cond publicconst.TaskCond, data model.ITask) bool {
	data.SetTaskValue(p.UserData.Level)
	return true
}

// initPassMission
func initPassMission(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		return false
	}

	if len(jTask.Data.Effect1) > 0 {
		stageId := jTask.Data.Effect1[0]
		missionConfig := template.GetMissionTemplate().GetMission(int(stageId))
		if missionConfig == nil {
			return false
		}

		isMain := true
		if missionConfig.Type != 1 {
			isMain = false
		}

		if mission := findMission(p, int(stageId), isMain); mission != nil && mission.IsPass {
			data.SetTaskValue(1)
			data.SetTaskState(publicconst.TASK_COMPLETE)
			return true
		}
	}
	return false
}

// initPassMissionNum 通关主线关卡数量
func initPassMissionNum(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		return false
	}

	if len(jTask.Data.Effect1) > 0 {
		num := passMissionNum(p, true)
		data.SetTaskValue(num)
		if num >= jTask.Data.Effect1[0] {
			data.SetTaskState(publicconst.TASK_COMPLETE)
		}
	}
	return false
}

func initPlayerLevel(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		return false
	}

	data.SetTaskValue(p.UserData.Level)
	if len(jTask.Data.Effect1) > 0 {
		if data.GetTaskValue() >= jTask.Data.Effect1[0] {
			data.SetTaskValue(jTask.Data.Effect1[0])
			data.SetTaskState(publicconst.TASK_COMPLETE)
		}
		return true
	}
	return false
}

func initPlayerShipLevel(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		//log.Errorf("account:%v taskId:%v not exist", playData.GetUserId(), data.GetTaskId())
		return false
	}

	if len(jTask.Data.Effect1) == 2 {
		if ship := getShip(p, jTask.Data.Effect1[0]); ship != nil {
			if ship.Level >= jTask.Data.Effect1[1] {
				data.SetTaskValue(jTask.Data.Effect1[1])
				data.SetTaskState(publicconst.TASK_COMPLETE)
			} else {
				data.SetTaskValue(ship.Level)
			}
			return true
		}
	}
	return false
}

func initShipLevelNum(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		return false
	}

	if len(jTask.Data.Effect1) == 2 {
		ships := p.UserData.Ships.Ships
		var num uint32 = 0
		for i := 0; i < len(ships); i++ {
			if ships[i].Level >= jTask.Data.Effect1[1] {
				num += 1
			}
		}
		oldValue := data.GetTaskValue()
		if num > oldValue {
			if num >= jTask.Data.Effect1[0] {
				data.SetTaskValue(jTask.Data.Effect1[0])
				data.SetTaskState(publicconst.TASK_COMPLETE)
			} else {
				data.SetTaskValue(num)
			}
			return true
		}
	}
	return false
}

func initShipRarityNum(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		return false
	}

	if len(jTask.Data.Effect1) == 1 {
		ships := p.UserData.Ships.Ships
		var num uint32 = 0
		for i := 0; i < len(ships); i++ {
			if shipConfig := template.GetShipTemplate().GetShip(ships[i].Id); shipConfig != nil {
				if shipConfig.Rarity >= jTask.Data.Effect1[0] {
					num += 1
				}
			}
		}
		oldValue := data.GetTaskValue()
		if num > oldValue {
			if num >= jTask.Data.MaxValue {
				data.SetTaskValue(jTask.Data.MaxValue)
				data.SetTaskState(publicconst.TASK_COMPLETE)
			} else {
				data.SetTaskValue(num)
			}
			return true
		}
	}
	return false
}

func initPlayerShipStarLevel(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		return false
	}

	if len(jTask.Data.Effect1) == 1 {
		if ship := getShip(p, jTask.Data.Effect1[0]); ship != nil {
			if ship.StarLevel >= jTask.Data.MaxValue {
				data.SetTaskValue(jTask.Data.MaxValue)
				data.SetTaskState(publicconst.TASK_COMPLETE)
			} else {
				data.SetTaskValue(ship.StarLevel)
			}
			return true
		}
	}
	return false
}

func initAnyEquipLevel(playData *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", playData.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		return false
	}

	maxLevel := getEquipPosMaxLevel(playData)
	if maxLevel > data.GetTaskValue() {
		data.SetTaskValue(maxLevel)
		if data.GetTaskValue() >= jTask.Data.MaxValue {
			data.SetTaskValue(jTask.Data.MaxValue)
			data.SetTaskState(publicconst.TASK_COMPLETE)
		}
		return true
	}
	return false
}

func initAllEquipPosLevel(playData *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", playData.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		return false
	}

	result := isAllPosOverLevel(playData, jTask.Data.Effect1[0])
	if result {
		data.SetTaskValue(jTask.Data.MaxValue)
		data.SetTaskState(publicconst.TASK_COMPLETE)
		return true
	}
	return false
}

func initEquipRarityNum(playData *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", playData.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		return false
	}

	if len(jTask.Data.Effect1) == 2 {
		num := getEquipRarityNum(playData, jTask.Data.Effect1[0], jTask.Data.Effect1[1])
		if num > data.GetTaskValue() {
			data.SetTaskValue(num)
			if data.GetTaskValue() >= jTask.Data.MaxValue {
				data.SetTaskValue(jTask.Data.MaxValue)
				data.SetTaskState(publicconst.TASK_COMPLETE)
			}
			return true
		}
	}
	return false
}

func initPeakFightRank(playData *player.Player, data model.ITask) bool {
	// jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	// if jTask == nil {
	// 	log.Error("initPeakFightRank task not found", zap.Uint64("accountId", playData.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
	// 	return false
	// }

	// ServMgr.GetPeakFightService().LoadPeakFight(playData)

	// if playData.PeakFight == nil {
	// 	return false
	// }

	// matchId := playData.PeakFight.BattleMatchId

	// if cfg := template.GetBattleMatchTemplate().GetCfg(matchId); cfg != nil {
	// 	ServMgr.GetTaskService().UpdateTask(
	// 		playData, true, publicconst.TASK_COND_PEAK_FIGHT_RANK, cfg.BattleLevel)
	// }
	return true
}

func initExploreBuild(playData *player.Player, data model.ITask) bool {
	// jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	// if jTask == nil {
	// 	log.Error("initExploreBuild task not found", zap.Uint64("accountId", playData.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
	// 	return false
	// }

	// ServMgr.GetExploreService().LoadExplore(playData)
	// if playData.Explore != nil && playData.Explore.Stages != nil {
	// 	for _, stage := range playData.Explore.Stages {

	// 		for _, build := range stage.Buildings {

	// 			if build.BuildingId == jTask.Data.Effect1[0] {
	// 				data.SetTaskValue(1)
	// 				data.SetTaskState(publicconst.TASK_COMPLETE)
	// 				return true
	// 			}

	// 		}
	// 	}
	// }

	return false
}

func initPet(playData *player.Player, data model.ITask) bool {
	// jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	// if jTask == nil {
	// 	log.Error("task not found", zap.Uint64("accountId", playData.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
	// 	//log.Errorf("account:%v taskId:%v not exist", playData.GetUserId(), data.GetTaskId())
	// 	return false
	// }
	// ServMgr.GetPetService().LoadPet(playData)

	// for _, petData := range playData.AccountPet.Pets {
	// 	if petData.Id == jTask.Data.Effect1[0] {
	// 		data.SetTaskValue(1)
	// 		data.SetTaskState(publicconst.TASK_COMPLETE)
	// 		return true
	// 	}
	// }
	return false
}

func initPetPartCharm(playData *player.Player, data model.ITask) bool {
	// jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	// if jTask == nil {
	// 	log.Error("task not found", zap.Int64("accountId", playData.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
	// 	//log.Errorf("account:%v taskId:%v not exist", playData.GetUserId(), data.GetTaskId())
	// 	return false
	// }
	// ServMgr.GetPetService().LoadPet(playData)
	// if playData.AccountPet.Charm == 0 {
	// 	return false
	// }

	// data.SetTaskValue(playData.AccountPet.Charm)
	// if data.GetTaskValue() >= jTask.Data.MaxValue {
	// 	data.SetTaskValue(jTask.Data.MaxValue)
	// 	data.SetTaskState(publicconst.TASK_COMPLETE)
	// }
	return true
}

func initPutOnEquip(playData *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", playData.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		//log.Errorf("account:%v taskId:%v not exist", playData.GetUserId(), data.GetTaskId())
		return false
	}

	var num uint32 = 0
	for i := 1; i <= 6; i++ {
		if posData := getEquipPos(playData, uint32(i)); posData != nil {
			if posData.EquipId > 0 {
				num += 1
			}
		}
	}
	if num > data.GetTaskValue() {
		data.SetTaskValue(num)
		if data.GetTaskValue() >= jTask.Data.MaxValue {
			data.SetTaskValue(jTask.Data.MaxValue)
			data.SetTaskState(publicconst.TASK_COMPLETE)
		}
		return true
	}
	return false
}

func initAnyWeaponLevel(playData *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", playData.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		//log.Errorf("account:%v taskId:%v not exist", playData.GetUserId(), data.GetTaskId())
		return false
	}

	maxLevel := getWeaponMaxLevel(playData)
	if maxLevel > data.GetTaskValue() {
		data.SetTaskValue(maxLevel)
		if data.GetTaskValue() >= jTask.Data.MaxValue {
			data.SetTaskValue(jTask.Data.MaxValue)
			data.SetTaskState(publicconst.TASK_COMPLETE)
		}
		return true
	}
	return false
}

func initSpecifyWeaponLevel(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		//log.Errorf("account:%v taskId:%v not exist", playData.GetUserId(), data.GetTaskId())
		return false
	}
	if len(jTask.Data.Effect1) < 1 {
		return false
	}
	if p.UserData.Weapon == nil {
		log.Error("weapon not found", zap.Uint64("accountId", p.GetUserId()))
		return false
	}
	for i := 0; i < len(p.UserData.Weapon.Weapons); i++ {
		if p.UserData.Weapon.Weapons[i].Id == jTask.Data.Effect1[0] {
			data.SetTaskValue(p.UserData.Weapon.Weapons[i].Level)
			if data.GetTaskValue() >= jTask.Data.MaxValue {
				data.SetTaskValue(jTask.Data.MaxValue)
				data.SetTaskState(publicconst.TASK_COMPLETE)
			}
			return true
		}
	}

	return false
}

func initAllWeaponMinLevel(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		//log.Errorf("account:%v taskId:%v not exist", playData.GetUserId(), data.GetTaskId())
		return false
	}
	if len(jTask.Data.Effect1) < 1 {
		return false
	}

	if p.UserData.Weapon == nil {
		log.Error("weapon not found", zap.Uint64("accountId", p.GetUserId()))
		return false
	}
	var minLevel uint32 = 1000
	allWeapon := template.GetWeaponTemplate().GetAllWeapon()
	for i := 0; i < len(allWeapon); i++ {
		if weapon := getWeapon(p, allWeapon[i].Id); weapon != nil {
			if weapon.Level < minLevel {
				minLevel = weapon.Level
			}
		} else {
			minLevel = 0
			break
		}
	}

	if minLevel > data.GetTaskValue() {
		data.SetTaskValue(minLevel)
		if data.GetTaskValue() >= jTask.Data.MaxValue {
			data.SetTaskValue(jTask.Data.MaxValue)
			data.SetTaskState(publicconst.TASK_COMPLETE)
		}
		return true
	}
	return false
}

func initWeaponLibLevel(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		//log.Errorf("account:%v taskId:%v not exist", playData.GetUserId(), data.GetTaskId())
		return false
	}
	if len(jTask.Data.Effect1) < 1 {
		return false
	}

	if p.UserData.Weapon == nil {
		log.Error("weapon not found", zap.Uint64("accountId", p.GetUserId()))
		return false
	}
	if p.UserData.Weapon.WeaponLibLevel > data.GetTaskValue() {
		data.SetTaskValue(p.UserData.Weapon.WeaponLibLevel)
		if data.GetTaskValue() >= jTask.Data.MaxValue {
			data.SetTaskValue(jTask.Data.MaxValue)
			data.SetTaskState(publicconst.TASK_COMPLETE)
		}
		return true
	}
	return false
}

func initLevel(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		//log.Errorf("account:%v taskId:%v not exist", playData.GetUserId(), data.GetTaskId())
		return false
	}

	data.SetTaskValue(p.UserData.Level)
	if p.UserData.Level >= jTask.Data.MaxValue {
		//data.SetTaskValue(jTask.Data.MaxValue)
		data.SetTaskState(publicconst.TASK_COMPLETE)
	}
	return true
}

func initShipStarNum(p *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		log.Error("task not found", zap.Uint64("accountId", p.GetUserId()), zap.Uint32("taskId", data.GetTaskId()))
		//log.Errorf("account:%v taskId:%v not exist", playData.GetUserId(), data.GetTaskId())
		return false
	}

	if len(jTask.Data.Effect1) == 1 {
		ships := p.UserData.Ships.Ships
		var num uint32 = 0
		for i := 0; i < len(ships); i++ {
			if ships[i].StarLevel >= jTask.Data.Effect1[0] {
				num += 1
			}
		}
		oldValue := data.GetTaskValue()
		if num > oldValue {
			if num >= jTask.Data.MaxValue {
				data.SetTaskValue(jTask.Data.MaxValue)
				data.SetTaskState(publicconst.TASK_COMPLETE)
			} else {
				data.SetTaskValue(num)
			}
			return true
		}
	}
	return false
}

func initHistoryDataTask(playData *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		return false
	}
	var condExtra uint32 = 0
	if len(jTask.Data.Effect1) > 0 {
		condExtra = jTask.Data.Effect1[0]
	}

	historyData := GetHistoryData(playData, publicconst.TaskCond(jTask.Data.TaskCondition), condExtra)
	if historyData != nil {
		if historyData.TaskValue >= jTask.Data.MaxValue {
			data.SetTaskValue(jTask.Data.MaxValue)
			data.SetTaskState(publicconst.TASK_COMPLETE)
		} else {
			// 登陆需要特殊处理
			if jTask.Data.TaskCondition == publicconst.TASK_COND_LOGIN {
				dailyRefreshTime := tools.GetDailyRefreshTime()
				if dailyRefreshTime == historyData.ExtraPara && historyData.TaskValue > 0 {
					data.SetTaskValue(historyData.TaskValue - 1)
				} else {
					data.SetTaskValue(historyData.TaskValue)
				}
			} else {
				data.SetTaskValue(historyData.TaskValue)
			}
		}
		return true
	}
	return false
}

func initAddItem(playData *player.Player, data model.ITask) bool {
	if jTask := template.GetTaskTemplate().GetTask(data.GetTaskId()); jTask != nil && jTask.Data.UseHisData == 1 {
		return initHistoryDataTask(playData, data)
	}
	return false
}

func initUseHistoryData(playData *player.Player, data model.ITask) bool {
	if jTask := template.GetTaskTemplate().GetTask(data.GetTaskId()); jTask != nil && jTask.Data.UseHisData == 1 {
		return initHistoryDataTask(playData, data)
	}
	return false
}

func initUseMissionBoxReward(playData *player.Player, data model.ITask) bool {
	jTask := template.GetTaskTemplate().GetTask(data.GetTaskId())
	if jTask == nil {
		return false
	}
	var boxCount uint32 = 0
	for i := 0; i < len(playData.UserData.Mission.Missions); i++ {
		if playData.UserData.Mission.Missions[i].BoxRewardState > 0 {
			boxCount++
		}
	}
	for i := 0; i < len(playData.UserData.Mission.Challenges); i++ {
		if playData.UserData.Mission.Challenges[i].BoxRewardState > 0 {
			boxCount++
		}
	}

	if boxCount >= jTask.Data.MaxValue {
		data.SetTaskValue(jTask.Data.MaxValue)
		data.SetTaskState(publicconst.TASK_COMPLETE)
	} else {
		data.SetTaskValue(boxCount)
	}
	return true
}
