package tda

import (
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 日常任務領獎	點擊領取每日/每週任務獎勵時記錄
type TaskRoutineReward struct {
	*CommonAttr
	Task_id      string `json:"task_id"`      // 任务ID
	Task_type_id string `json:"task_type_id"` // 任务类型ID
	ExpNum       uint32 `json:"exp_num"`      // 获得日常/周常经验数量
	//Item_gain       []*Item `json:"item_gain"`       // 获得道具列表
	Activeness_add  uint32 `json:"activeness_add"`  // 獲取活躍度
	Activeness_accu uint32 `json:"activeness_accu"` // 該週期活躍度累計
}

func TdaTaskRoutineReward(channelId uint32, commonAttr *CommonAttr, taskId, taskTypeId, activenessAdd, activenessAccu, expNum uint32) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaTaskRoutineReward", func() {
		tdaData := &TaskRoutineReward{
			CommonAttr:   commonAttr,
			Task_id:      strconv.Itoa(int(taskId)),
			Task_type_id: strconv.Itoa(int(taskTypeId)),
			ExpNum:       expNum,
			//Item_gain:       items,
			Activeness_add:  activenessAdd,
			Activeness_accu: activenessAccu,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Task_Routine_Reward

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaTaskRoutineReward tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaTaskRoutineReward track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 主線任務領獎	點擊領取主線任務獎勵時記錄
type TaskMainReward struct {
	*CommonAttr
	Task_id        string  `json:"task_id"`        // 任务ID
	Task_type_id   string  `json:"task_type_id"`   // 任务类型ID
	Item_gain      []*Item `json:"item_gain"`      // 获得道具列表
	Main_work_time uint32  `json:"main_work_time"` // 距离上条完成主线的在线时长
}

func TdaTaskMainReward(channelId uint32, commonAttr *CommonAttr, taskId, taskTypeId uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaTaskMainReward", func() {
		tdaData := &TaskMainReward{
			CommonAttr:   commonAttr,
			Task_id:      strconv.Itoa(int(taskId)),
			Task_type_id: strconv.Itoa(int(taskTypeId)),
			Item_gain:    items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Task_Main_Reward

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaTaskMainReward tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaTaskMainReward track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 成就任務領獎	點擊領取成就任務獎勵時記錄
type TaskAchiReward struct {
	*CommonAttr
	Task_id      string  `json:"task_id"`      // 任务ID
	Task_type_id string  `json:"task_type_id"` // 任务类型ID
	Item_gain    []*Item `json:"item_gain"`    // 获得道具列表
}

func TdaTaskAchiReward(channelId uint32, commonAttr *CommonAttr, taskId, taskTypeId uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaTaskAchiReward", func() {
		tdaData := &TaskAchiReward{
			CommonAttr:   commonAttr,
			Task_id:      strconv.Itoa(int(taskId)),
			Task_type_id: strconv.Itoa(int(taskTypeId)),
			Item_gain:    items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Task_Achi_Reward

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaTaskAchiReward tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaTaskAchiReward track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 聯盟任務領獎	點擊領取聯盟任務獎勵時記錄
type TaskGuildReward struct {
	*CommonAttr
	Task_id      string  `json:"task_id"`      // 任务ID
	Task_type_id string  `json:"task_type_id"` // 任务类型ID
	Item_gain    []*Item `json:"item_gain"`    // 获得道具列表
}

func TdaTaskGuildReward(channelId uint32, commonAttr *CommonAttr, taskId, taskTypeId uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaTaskGuildReward", func() {
		tdaData := &TaskGuildReward{
			CommonAttr:   commonAttr,
			Task_id:      strconv.Itoa(int(taskId)),
			Task_type_id: strconv.Itoa(int(taskTypeId)),
			Item_gain:    items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Task_Guild_Reward

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaTaskGuildReward tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaTaskGuildReward track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 完成任務條件	完成任務條件時記錄，不分任務型態都記錄在完成任務的點位
type TaskAllComplete struct {
	*CommonAttr
	Task_id      string `json:"task_id"`      // 任务ID
	Task_type_id string `json:"task_type_id"` // 任务类型ID
}

func TdaTaskComplete(channelId uint32, commonAttr *CommonAttr, taskId, taskTypeId uint32) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaTaskComplete", func() {
		tdaData := &TaskAllComplete{
			CommonAttr:   commonAttr,
			Task_id:      strconv.Itoa(int(taskId)),
			Task_type_id: strconv.Itoa(int(taskTypeId)),
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Task_All_Complete

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaTaskComplete tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaTaskComplete track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
