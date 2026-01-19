package tda

import (
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

type PvpOpUnit struct {
	AccountId  string `json:"role_id"`    // 對手角色id
	Robot      uint32 `json:"robot"`      // 對手角色id
	Kulu_id    string `json:"kulu_id"`    // 對手機甲id
	Kulu_class string `json:"kulu_class"` // 對手機甲品階
	Kulu_rank  string `json:"kulu_rank"`  // 對手機甲星級
	Kulu_level uint32 `json:"kulu_level"` // 對手機甲等級
	Power      uint32 `json:"power"`      // 對手角色戰力
}

// pvp 開始匹配	pvp 匹配開始時上報
type PvpMatch struct {
	*CommonAttr
	Success_or_fail string       `json:"success_or_fail"` // 匹配成功/失敗 "1=成功, ０=失敗"
	Fail_reason     string       `json:"fail_reason"`     // 匹配成功/失敗原因 如: 無人匹配、自行退出、斷線
	Match_time      uint32       `json:"match_time"`      // 匹配花費時間 秒(s)
	Ops             []*PvpOpUnit `json:"ops"`             // 对手(不算机器人)
}

func TdaPvpMatch(channelId uint32, commonAttr *CommonAttr, successOrFail, reason string, ops []*PvpOpUnit) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaPvpMatch", func() {
		tdaData := &PvpMatch{
			CommonAttr:      commonAttr,
			Success_or_fail: successOrFail,
			Fail_reason:     reason,
			//Match_time:      0,
			Ops: ops,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Pvp_Match

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaPvpMatch tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaPvpMatch track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// pvp 結算	pvp 結算時上報
type PvpEnd struct {
	*CommonAttr
	PvpMatchRank   string       `json:"pvp_match_rank"` // 匹配成功/失敗 "1=win, lose=2"
	Pvp_grade      uint32       `json:"pvp_grade"`      // pvp 積分變動
	Pvp_rank       uint32       `json:"pvp_rank"`       // pvp 當前積分
	Pvp_End_Reason string       `json:"pvp_end_reason"` // 结束原因
	Pvp_Duration   uint32       `json:"pvp_duration"`   // 当局时长(秒)
	Ops            []*PvpOpUnit `json:"ops"`            // 对手(不算机器人)
}

func TdaPvpEnd(channelId uint32, commonAttr *CommonAttr, pvpMatchRank, Pvp_End_Reason string, pvpGrade, pvpRank, pvp_duration uint32, ops []*PvpOpUnit) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaPvpEnd", func() {
		tdaData := &PvpEnd{
			CommonAttr:     commonAttr,
			PvpMatchRank:   pvpMatchRank,
			Pvp_grade:      pvpGrade,
			Pvp_rank:       pvpRank,
			Pvp_End_Reason: Pvp_End_Reason,
			Pvp_Duration:   pvp_duration,
			Ops:            ops,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Pvp_End

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaPvpEnd tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaPvpEnd track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
