package tda

import (
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 常規活動開啟	開啟活動頁面時記錄
type EventOpen struct {
	*CommonAttr
	Activity_id      string `json:"activity_id"`      // 活動ID
	Activity_type_id string `json:"activity_type_id"` // 活動類型ID
}

func TdaEventOpen(channelId uint32, commonAttr *CommonAttr, activityId, activityTypeId uint32) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaEventOpen", func() {
		tdaData := &EventOpen{
			CommonAttr:       commonAttr,
			Activity_id:      strconv.Itoa(int(activityId)),
			Activity_type_id: strconv.Itoa(int(activityTypeId)),
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Event_Open

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaEventOpen tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaEventOpen track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 常規活動任務完成	完成活動時記錄，獲取最大獎、完成所有任務
type EventTaskComplete struct {
	*CommonAttr
	Activity_id      string `json:"activity_id"`      // 活動ID
	Activity_type_id string `json:"activity_type_id"` // 活動類型ID
	Task_id          string `json:"task_id"`          // 任务ID
	Task_type_id     string `json:"task_type_id"`     // 任务类型ID
}

func TdaEventTaskComplete(channelId uint32, commonAttr *CommonAttr, activityId, activityTypeId, taskId, taskTypeId uint32) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaEventTaskComplete", func() {
		tdaData := &EventTaskComplete{
			CommonAttr:       commonAttr,
			Activity_id:      strconv.Itoa(int(activityId)),
			Activity_type_id: strconv.Itoa(int(activityTypeId)),
			Task_id:          strconv.Itoa(int(taskId)),
			Task_type_id:     strconv.Itoa(int(taskTypeId)),
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Event_Task_Complete

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaEventTaskComplete tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaEventTaskComplete track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 常規活動領獎	點擊領取活動獎勵時記錄
type EventReward struct {
	*CommonAttr
	Activity_id      string  `json:"activity_id"`      // 活動ID
	Activity_type_id string  `json:"activity_type_id"` // 活動類型ID
	Task_id          string  `json:"task_id"`          // 任务ID
	Task_type_id     string  `json:"task_type_id"`     // 任务类型ID
	Item_gain        []*Item `json:"item_gain"`        // 获得道具列表
}

func TdaEventReward(channelId uint32, commonAttr *CommonAttr, activityId, activityTypeId, taskId, taskTypeId uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaEventReward", func() {
		tdaData := &EventReward{
			CommonAttr:       commonAttr,
			Activity_id:      strconv.Itoa(int(activityId)),
			Activity_type_id: strconv.Itoa(int(activityTypeId)),
			Task_id:          strconv.Itoa(int(taskId)),
			Task_type_id:     strconv.Itoa(int(taskTypeId)),
			Item_gain:        items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Event_Reward

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaEventReward tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaEventReward track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 沙漠大冒險抽獎	抽獎完成獲得獎勵時上報
type EventDessertLottery struct {
	*CommonAttr
	Lottery_ticket_used uint32  `json:"lottery_ticket_used"` // 抽獎券消耗
	Item_gain           []*Item `json:"item_gain"`           // 获得道具列表
}

func TdaEventDessertLottery(channelId uint32, commonAttr *CommonAttr, lotteryTicketUsed uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaEventDessertLottery", func() {
		tdaData := &EventDessertLottery{
			CommonAttr:          commonAttr,
			Lottery_ticket_used: lotteryTicketUsed,
			Item_gain:           items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Event_Dessert_Lottery

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaEventDessertLottery tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaEventDessertLottery track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 沙漠大冒險圖鑑解鎖	解鎖沙漠大冒險圖鑑時上報
type EventDessertCollection struct {
	*CommonAttr
	Collection_id string  `json:"collection_id"` // 圖鑑id
	Item_gain     []*Item `json:"item_gain"`     // 获得道具列表
}

func TdaEventDessertCollection(channelId uint32, commonAttr *CommonAttr, collectionId uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaEventDessertCollection", func() {
		tdaData := &EventDessertCollection{
			CommonAttr:    commonAttr,
			Collection_id: strconv.Itoa(int(collectionId)),
			Item_gain:     items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Event_Dessert_Collection

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaEventDessertCollection tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaEventDessertCollection track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 沙漠大冒險簽到	完成沙漠大冒險簽到時上報
type EventDessertSignin struct {
	*CommonAttr
	Item_gain                 []*Item `json:"item_gain"`                 // 获得道具列表
	Is_retro_signin           bool    `json:"is_retro_signin"`           // 是否為補簽
	Retro_signin_total        uint32  `json:"retro_signin_total"`        // 累計補簽次數
	Retro_signin_cost_diamond uint32  `json:"retro_signin_cost_diamond"` // 補簽消耗鑽石數
}

func TdaEventDessertSignin(channelId uint32, commonAttr *CommonAttr, isRetroSignin bool, retroSigninTotal, retroSigninCostDiamond uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaEventDessertSignin", func() {
		tdaData := &EventDessertSignin{
			CommonAttr:                commonAttr,
			Item_gain:                 items,
			Is_retro_signin:           isRetroSignin,
			Retro_signin_total:        retroSigninTotal,
			Retro_signin_cost_diamond: retroSigninCostDiamond,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Event_Dessert_Signin

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaEventDessertSignin tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaEventDessertSignin track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 主理人之路任務領獎	點擊領取主理人之路任務獎勵時記錄
type EventOpenServerReward struct {
	*CommonAttr
	Task_id      string  `json:"task_id"`      // 任务ID
	Task_type_id string  `json:"task_type_id"` // 任务类型ID
	Item_gain    []*Item `json:"item_gain"`    // 获得道具列表
}

func TdaEventOpenServerReward(channelId uint32, commonAttr *CommonAttr, taskId, taskTypeId uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaEventOpenServerReward", func() {
		tdaData := &EventOpenServerReward{
			CommonAttr:   commonAttr,
			Task_id:      strconv.Itoa(int(taskId)),
			Task_type_id: strconv.Itoa(int(taskTypeId)),
			Item_gain:    items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Event_Open_Server_Reward

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaEventOpenServerReward tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaEventOpenServerReward track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
