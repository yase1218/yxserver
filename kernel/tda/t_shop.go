package tda

import (
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 世界探索抽卡前進	抽卡獲得獎勵時上報
type ShopBuy struct {
	*CommonAttr
	Store_type string `json:"store_type"` // 商店类型
	Goods_id   string `json:"goods_id"`   // 商品ID
	Goods_name string `json:"goods_name"` // 商品名称
	Goods_num  uint32 `json:"goods_num"`  // 数量
	Money_id   string `json:"money_id"`   // 貨幣 ID
	Cost_num   uint32 `json:"cost_num"`   // 消耗数量
}

func TdaShopBuy(channelId uint32, commonAttr *CommonAttr, storeType, goodsId, goodsNum, moneyId, costNum uint32) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaShopBuy", func() {
		tdaData := &ShopBuy{
			CommonAttr: commonAttr,
			Store_type: strconv.Itoa(int(storeType)),
			Goods_id:   strconv.Itoa(int(goodsId)),
			//Goods_name:
			Goods_num: goodsNum,
			Money_id:  strconv.Itoa(int(moneyId)),
			Cost_num:  costNum,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Shop_Buy

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaShopBuy tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaShopBuy track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
