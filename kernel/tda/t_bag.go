package tda

import (
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

type Item struct {
	ItemId  string `json:"item_id"`  // 道具id
	ItemNum uint32 `json:"item_num"` // 道具数量
}

// 道具变动	道具发生变化时记录
type ItemChange struct {
	*CommonAttr
	Change_type      string `json:"change_type"`      // 變動類型 增加 add = 0 / 减少 reduce = 1
	Item_id          string `json:"item_id"`          // 物品 ID
	Item_name        string `json:"item_name"`        // 物品名稱
	Change_reason    string `json:"change_reason"`    // 變動原因
	Change_subreason string `json:"change_subreason"` // 變動次級原因
	Change_num       uint32 `json:"change_num"`       // 变动数量
	Change_after     uint32 `json:"change_after"`     // 变化后数量
}

// 貨幣变动	貨幣发生变化时记录
type MoneyChange struct {
	*CommonAttr
	Change_type      string `json:"change_type"`      // 變動類型 增加 add = 0 / 减少 reduce = 1
	Money_id         string `json:"money_id"`         // 物品 ID
	Money_name       string `json:"money_name"`       // 物品名稱
	Change_reason    string `json:"change_reason"`    // 變動原因
	Change_subreason string `json:"change_subreason"` // 變動次級原因
	Change_num       uint32 `json:"change_num"`       // 变动数量
	Change_after     uint32 `json:"change_after"`     // 变化后数量
}

func TdaItemChange(channelId uint32, commonAttr *CommonAttr, changeType string, itemId, num, source, after uint32) {
	if !Send() {
		return
	}
	// tda event logout
	go tools.GoSafe("TdaItemChange", func() {
		tdaData := &ItemChange{
			CommonAttr:  commonAttr,
			Change_type: changeType,
			Item_id:     strconv.Itoa(int(itemId)),
			//Money_name:       itemConfig.,
			Change_reason: strconv.Itoa(int(source)),
			//Change_subreason: "",
			Change_num:   num,
			Change_after: after,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Item_Change

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaItemChange tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaItemChange track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
