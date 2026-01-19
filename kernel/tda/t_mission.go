package tda

import (
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 主線關卡結束-后端
type MainBattleDoneServer struct {
	*CommonAttr
	Mainbattle_serial_num string `json:"mainbattle_serial_num"` // 戰鬥流水號 作為關聯同一場戰鬥資訊用
	Mainbattle_timespent  uint32 `json:"mainbattle_timespent"`  // 通關耗時 秒(s)
}

func TdaMainBattleDoneServer(channelId uint32, commonAttr *CommonAttr, mainbattle_serial_num string, Mainbattle_timespent uint32) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaLotteryResult", func() {
		tdaData := &MainBattleDoneServer{
			CommonAttr:            commonAttr,
			Mainbattle_serial_num: mainbattle_serial_num,
			Mainbattle_timespent:  Mainbattle_timespent,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Main_Battle_Done_Server

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaMainBattleDoneServer tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaMainBattleDoneServer track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
