package tda

import (
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

// 用户/账号注册 (账号) 注册完成时上报
type Register struct {
	*CommonAttr
	First_register string `json:"first_register" canEmpty:"true"` // 是否首次註冊, 跨服建立新角色後
}

func TdaRegister(channelId uint32, accountId, distinctId string, commonAttr *CommonAttr) {
	if !Send() {
		return
	}
	// tda event register
	go tools.GoSafe("tda event register", func() {
		tdaData := Register{
			CommonAttr:     commonAttr,
			First_register: "1",
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Register

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("TdaRegister err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaRegister track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 用户/账号登录	(账号) 登录成功后上报
type Login struct {
	*CommonAttr
	Rolename      string    `json:"rolename"`      // 角色昵称
	Add_day       uint32    `json:"add_day"`       // 累积登陆天数
	Last_out_time time.Time `json:"last_out_time"` // 上次退出游戏时间
}

func TdaLogin(channelId uint32, roleName string, addDay uint32, lastOutTime time.Time, commonAttr *CommonAttr) {
	if !Send() {
		return
	}
	// tda event login
	go tools.GoSafe("tda event login", func() {
		tdaData := &Login{
			CommonAttr:    commonAttr,
			Rolename:      roleName,
			Add_day:       addDay,
			Last_out_time: lastOutTime,
		}
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Login

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("TdaLogin err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {
			if err := GetTa().Track(commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaLogin track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, commonAttr.AccountId, commonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 用户/账号登出	(账号) 登出成功后上报
type EquipUnit struct {
	Pos    uint32 `json:"pos"`     // 位置
	ItemId uint32 `json:"item_id"` // 物品id
	Lv     uint32 `json:"lv"`      // 等级
	Rarity uint32 `json:"rarity"`  // 品阶(稀有度)
}
type Logout struct {
	*CommonAttr
	Rolename               string       `json:"rolename"`               // 角色昵称
	Online_time            uint32       `json:"online_time"`            // 本次在线时长
	Add_time               uint32       `json:"add_time"`               // 角色累计在线时长
	Online_time_mainbattle uint32       `json:"online_time_mainbattle"` // 參與主線關卡在線時長
	Online_time_explore    uint32       `json:"online_time_explore"`    // 參與世界探索在線時長
	Online_time_pvp        uint32       `json:"online_time_pvp"`        // 參與pvp在線時長
	Online_time_event      uint32       `json:"online_time_event"`      // 參與活動在線時長
	Tmp_kulu               string       `json:"tmp_kulu"`               // 角色
	Tmp_level              uint32       `json:"tmp_level"`              // 角色等级
	Tmp_power              uint32       `json:"tmp_power"`              // 角色戰力
	Tmp_equip              []*EquipUnit `json:"tmp_equip"`              // 裝備
	Tmp_vip_level          uint32       `json:"tmp_vip_level"`          // 当前vip等级
	Tmp_monthpass          uint32       `json:"tmp_monthpass"`          // 当前月卡状态
	Tmp_main_id            uint32       `json:"tmp_main_id"`            // 当前主线任务进度
	Max_battle_id          int          `json:"max_battle_id"`          // 当前最大关卡
	Tmp_diamond            int64        `json:"tmp_diamond"`            // 当前钻石数量
	Tmp_gold               int64        `json:"tmp_gold"`               // 当前金币
	Tmp_stamina            int64        `json:"tmp_stamina"`            // 当前体力
	Total_pay_amount       uint32       `json:"total_pay_amount"`       // 当前付费金额
	Total_pay_num          uint32       `json:"total_pay_num"`          // 当前付费次数
	Total_ad_num           uint32       `json:"total_ad_num"`           // 累计观看广告次数
}

func TdaLogout(channelId uint32, tdaData *Logout) {
	if !Send() {
		return
	}
	// tda event logout
	go tools.GoSafe("tda event logout", func() {
		tdaData.CommonAttr.EventTime = time.Now()
		tdaData.CommonAttr.EventName = Tda_Logout

		if tdaDataMap, err := FlattenStructToMap(tdaData); err != nil {
			log.Error("TdaLogout err", zap.Error(err), zap.Reflect("data", tdaData))
		} else {

			if err := GetTa().Track(tdaData.CommonAttr.AccountId, tdaData.CommonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap); err != nil {
				log.Error("tda TdaLogout track err", zap.Error(err), zap.Reflect("data", tdaData))
			}

			UploadLog(channelId, tdaData.CommonAttr.AccountId, tdaData.CommonAttr.DistinctId, tdaData.CommonAttr.EventName, tdaDataMap)
		}
	})
}

// 創建角色	用户創建角色，產生新的roleid时上报
type CreateRole struct {
	*CommonAttr
}

// 綁定帳號	用戶完成帳號綁定時上報
type BindAccount struct {
	*CommonAttr
	Bind_platform      string `json:"bind_platform"`      // 綁定平台渠道
	Main_bind_platform string `json:"main_bind_platform"` // 綁定主平台渠道
}
