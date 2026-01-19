package tda

import (
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 裝備升級	裝備升級完成時上報
type EquipUpgrade struct {
	*CommonAttr
	Equip_part_id string  `json:"equip_part_id"` // 裝備部位id
	Equip_level   uint32  `json:"equip_level"`   // 裝備等級
	Equip_rank    string  `json:"equip_rank"`    // 裝備品階
	Item_used     []*Item `json:"item_used"`     // 消耗道具列表
}

func TdaEquipUpgrade(channelId uint32, commonAttr *CommonAttr, equipPartId, equipLevel, equipRank uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaEquipUpgrade", func() {
		tdaData := &EquipUpgrade{
			CommonAttr:    commonAttr,
			Equip_part_id: strconv.Itoa(int(equipPartId)),
			Equip_level:   equipLevel,
			Equip_rank:    strconv.Itoa(int(equipRank)),
			Item_used:     items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Equip_Upgrade

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaEquipUpgrade tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaEquipUpgrade track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 裝備進階	裝備進階完成時上報
type EquipRankUp struct {
	*CommonAttr
	Equip_part_id string  `json:"equip_part_id"` // 裝備部位id
	Equip_level   uint32  `json:"equip_level"`   // 裝備等級
	Equip_rank    string  `json:"equip_rank"`    // 裝備品階
	Item_used     []*Item `json:"item_used"`     // 消耗道具列表
}

func TdaEquipRankUp(channelId uint32, commonAttr *CommonAttr, equipPartId, equipLevel, equipRank uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaEquipRankUp", func() {
		tdaData := &EquipRankUp{
			CommonAttr:    commonAttr,
			Equip_part_id: strconv.Itoa(int(equipPartId)),
			Equip_level:   equipLevel,
			Equip_rank:    strconv.Itoa(int(equipRank)),
			Item_used:     items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Equip_Rank_Up

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaEquipRankUp tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaEquipRankUp track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
