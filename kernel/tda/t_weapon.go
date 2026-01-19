package tda

import (
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 武器艙系統解鎖	武器解鎖時上報
type WeaponSystemUnlock struct {
	*CommonAttr
	Subweapon_id string `json:"subweapon_id"` // 武器id
}

func TdaWeaponSystemUnlock(channelId uint32, commonAttr *CommonAttr, subweaponId uint32) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaWeaponSystemUnlock", func() {
		tdaData := &WeaponSystemUnlock{
			CommonAttr:   commonAttr,
			Subweapon_id: strconv.Itoa(int(subweaponId)),
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Weapon_System_Unlock

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaWeaponSystemUnlock tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(tdaData.CommonAttr.AccountId, tdaData.CommonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaWeaponSystemUnlock track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, tdaData.CommonAttr.AccountId, tdaData.CommonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 武器艙系統升級	武器升級時上報
type WeaponSystemUpgrade struct {
	*CommonAttr
	Subweapon_id    string  `json:"subweapon_id"`    // 武器id
	Subweapon_level uint32  `json:"subweapon_level"` // 武器等級
	Item_used       []*Item `json:"item_used"`       // 消耗道具列表
}

func TdaWeaponSystemUpgrade(channelId uint32, commonAttr *CommonAttr, subweaponId, subweaponLevel uint32, items []*Item) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaWeaponSystemUpgrade", func() {
		tdaData := &WeaponSystemUpgrade{
			CommonAttr:      commonAttr,
			Subweapon_id:    strconv.Itoa(int(subweaponId)),
			Subweapon_level: subweaponLevel,
			Item_used:       items,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Weapon_System_Upgrade

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaWeaponSystemUpgrade tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(tdaData.CommonAttr.AccountId, tdaData.CommonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaWeaponSystemUpgrade track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, tdaData.CommonAttr.AccountId, tdaData.CommonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
