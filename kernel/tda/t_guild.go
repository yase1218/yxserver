package tda

import (
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 加入/退出公会	加入or 退出公会时触发
type GuildJoinOrLeave struct {
	*CommonAttr
	Guild_level uint32 `json:"guild_level"` // 公会等级
	Guild_id    string `json:"guild_id"`    // 公会ID
	Guild_name  string `json:"guild_name"`  // 公会名称
	Reason      string `json:"reason"`      // 加入/退出
}

func TdaGuildJoinOrLeave(channelId uint32, commonAttr *CommonAttr, guildLv, guildId uint32, guildName, resson string) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaGuildJoinOrLeave", func() {
		tdaData := &GuildJoinOrLeave{
			CommonAttr:  commonAttr,
			Guild_level: guildLv,
			Guild_id:    strconv.Itoa(int(guildId)),
			Guild_name:  guildName,
			Reason:      resson,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Guild_JoinOrLeave

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaGuildJoinOrLeave tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaGuildJoinOrLeave track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 创建公会	自己创建公会时触发
type GuildCreate struct {
	*CommonAttr
	Guild_level uint32 `json:"guild_level"` // 公会等级
	Guild_id    string `json:"guild_id"`    // 公会ID
	Guild_name  string `json:"guild_name"`  // 公会名称
}

func TdaGuildCreate(channelId uint32, commonAttr *CommonAttr, guildLv, guildId uint32, guildName string) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaGuildCreate", func() {
		tdaData := &GuildCreate{
			CommonAttr:  commonAttr,
			Guild_level: guildLv,
			Guild_id:    strconv.Itoa(int(guildId)),
			Guild_name:  guildName,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Guild_Create

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaGuildCreate tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaGuildCreate track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 公會活躍度	公會活躍度變動時上報
type GuildActiveNess struct {
	*CommonAttr
	Guild_level       uint32 `json:"guild_level"`       // 数值
	Guild_id          string `json:"guild_id"`          // 字符串
	Guild_name        string `json:"guild_name"`        // 字符串
	Guild_role_num    uint32 `json:"guild_role_num"`    // 数值
	Activeness_change uint32 `json:"activeness_change"` // 数值
}

func TdaGuildActiveNess(channelId uint32, commonAttr *CommonAttr, guildLv, guildId uint32, guildName string, guildRoleNum, activeChange uint32) {
	if !Send() {
		return
	}
	go tools.GoSafe("TdaGuildActiveNess", func() {
		tdaData := &GuildActiveNess{
			CommonAttr:        commonAttr,
			Guild_level:       guildLv,
			Guild_id:          strconv.Itoa(int(guildId)),
			Guild_name:        guildName,
			Guild_role_num:    guildRoleNum,
			Activeness_change: activeChange,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Guild_Active_Ness

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("tda TdaGuildActiveNess tdaStructToMap err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaGuildActiveNess track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}
