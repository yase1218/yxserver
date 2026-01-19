package tda

import (
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 機甲獲取	解鎖新機甲時上報
type KuluUnlock struct {
	*CommonAttr
	Kulu_id    string `json:"kulu_id"`    // 機甲id
	Kulu_class string `json:"kulu_class"` // 機甲品階
}

func TdaKuluUnlock(channelId uint32, commonAttr *CommonAttr, kuluId, kuluClass uint32) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaKuluUnlock", func() {
		tdaData := &KuluUnlock{
			CommonAttr: commonAttr,
			Kulu_id:    strconv.Itoa(int(kuluId)),
			Kulu_class: strconv.Itoa(int(kuluClass)),
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Kulu_Unlock

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaKuluUnlock tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaKuluUnlock track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 機甲升級	升級機甲完成時上報
type KuluUpgrade struct {
	*CommonAttr
	Kulu_level uint32  `json:"kulu_level"` // 機甲等級
	Item_used  []*Item `json:"item_used"`  // 消耗道具列表
}

func TdaKuluUpgrade(channelId uint32, commonAttr *CommonAttr, kuluLevel uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaKuluUpgrade", func() {
		tdaData := &KuluUpgrade{
			CommonAttr: commonAttr,
			Kulu_level: kuluLevel,
			Item_used:  items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Kulu_Upgrade

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaKuluUpgrade tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaKuluUpgrade track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 機甲升星	升星機甲完成時上報
type KuluRankUp struct {
	*CommonAttr
	Kulu_id    string  `json:"kulu_id"`    // 機甲id
	Kulu_class string  `json:"kulu_class"` // 機甲品階
	Kulu_rank  string  `json:"kulu_rank"`  // 機甲星級
	Item_used  []*Item `json:"item_used"`  // 消耗道具列表
}

func TdaKuluRankUp(channelId uint32, commonAttr *CommonAttr, kuluId, kuluClass, kuluRank uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaKuluRankUp", func() {
		tdaData := &KuluRankUp{
			CommonAttr: commonAttr,
			Kulu_id:    strconv.Itoa(int(kuluId)),
			Kulu_class: strconv.Itoa(int(kuluClass)),
			Kulu_rank:  strconv.Itoa(int(kuluRank)),
			Item_used:  items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Kulu_Rank_Up

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaKuluRankUp tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaKuluRankUp track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
