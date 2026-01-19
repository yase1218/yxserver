package tda

import (
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 世界探索拆卡包	拆卡包獲得獎勵時上報
type AdventureGetCard struct {
	*CommonAttr
	Card_rank string  `json:"card_rank"` // 卡包品階
	Item_gain []*Item `json:"item_gain"` // 获得道具列表
}

func TdaAdventureGetCard(channelId uint32, commonAttr *CommonAttr, cardRank uint32, items []*Item, post func(string)) {
	if !Send() {
		return
	}
	go tools.GoSafePost("TdaAdventureGetCard", func() {
		tdaData := &AdventureGetCard{
			CommonAttr: commonAttr,
			Card_rank:  strconv.Itoa(int(cardRank)),
			Item_gain:  items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Adventrue_Get_Card

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaAdventureGetCard tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaAdventureGetCard track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	}, post)
}

// 世界探索抽卡前進	抽卡獲得獎勵時上報
type AdventureUseCard struct {
	*CommonAttr
	Adventure_id string  `json:"adventure_id"` // 事件id
	Item_gain    []*Item `json:"item_gain"`    // 获得道具列表
}

func TdaAdventureUseCard(channelId uint32, commonAttr *CommonAttr, adventureId uint32, items []*Item, post func(string)) {
	if !Send() {
		return
	}
	go tools.GoSafePost("TdaAdventureUseCard", func() {
		tdaData := &AdventureUseCard{
			CommonAttr:   commonAttr,
			Adventure_id: strconv.Itoa(int(adventureId)),
			Item_gain:    items,
		}
		tdaData.CommonAttr.EventName = Tda_Adventrue_Use_Card

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaAdventureUseCard tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaAdventureUseCard track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	}, post)
}
