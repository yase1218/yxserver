package tda

import (
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 接取合約	接取合約任務
type ContractAccept struct {
	*CommonAttr
	Contract_id   string `json:"contract_id"`   // 合約ID
	Max_battle_id string `json:"max_battle_id"` // 當時的最高主線關卡進度
}

func TdaContractAccept(channelId uint32, commonAttr *CommonAttr, contractId, maxBattleId string) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaContractAccept", func() {
		tdaData := &ContractAccept{
			CommonAttr:    commonAttr,
			Contract_id:   contractId,
			Max_battle_id: maxBattleId,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Contract_Accept

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaContractAccept tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaContractAccept track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 領獎合約	合約任務領獎
type ContractReward struct {
	*CommonAttr
	Contract_id string  `json:"contract_id"` // 合約ID
	Item_gain   []*Item `json:"item_gain"`   // 获得道具列表
}

func TdaContractReward(channelId uint32, commonAttr *CommonAttr, contractId string, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaContractAccept", func() {
		tdaData := &ContractReward{
			CommonAttr:  commonAttr,
			Contract_id: contractId,
			Item_gain:   items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Contract_Reward

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaContractReward tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaContractReward track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
