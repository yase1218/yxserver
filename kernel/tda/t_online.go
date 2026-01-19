package tda

import (
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

type OnlineCnt struct {
	*CommonAttr
	Onlinecnt_ios int    `json:"onlinecnt_ios"` // ios在線人數
	Onlinecnt_aos int    `json:"onlinecnt_aos"` // aos在線人數
	Onlinecnt_pc  int    `json:"onlinecnt_pc"`  // pc在線人數
	Onlinecnt     int    `json:"onlinecnt"`     // 總在線人數
	Country       string `json:"country"`       // 歸屬國家
	Region        string `json:"region"`        // 歸屬區域
	Server_id     string `json:"server_id"`     // 登入服务器id
	Loginchannel  string `json:"loginchannel"`  // 登入渠道
	Platform      string `json:"platform"`      // 平台标识
}

func TdaOnlineCnt(channelId uint32, tdaData *OnlineCnt) {
	if !Send() {
		return
	}
	// tda event onlinecnt
	go tools.GoSafe("TdaOnlineCnt", func() {
		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaOnlineCnt tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track("0", "0", tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaOnlineCnt track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, "0", "0", tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
