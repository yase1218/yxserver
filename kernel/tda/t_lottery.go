package tda

import (
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 抽卡結果	抽卡動作完成時上報
type LotteryResult struct {
	*CommonAttr
	Reason     string  `json:"reason"`     // 操作行为
	Lottery_id string  `json:"lottery_id"` // 卡池
	Add_num    uint32  `json:"add_num"`    // 当前卡池抽奖次数
	Item_gain  []*Item `json:"item_gain"`  // 获得道具列表
	Item_use   []*Item `json:"item_used"`  // 消耗道具列表
}

func TdaLotteryResult(channelId uint32, commonAttr *CommonAttr, reason string, lotteryId, addNum uint32, gainItems, useItems []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaLotteryResult", func() {
		tdaData := &LotteryResult{
			CommonAttr: commonAttr,
			Reason:     reason,
			Lottery_id: strconv.Itoa(int(lotteryId)),
			Add_num:    addNum,
			Item_gain:  gainItems,
			Item_use:   useItems,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Lottery_Result

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaLotteryResult tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaLotteryResult track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
