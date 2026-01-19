package tda

import (
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"kernel/tools"
)

type CommonUser struct {
	Reg_country          string    `json:"reg_country"`          // 註冊國家 (注册时记录)
	Reg_region           string    `json:"reg_region"`           // 註冊區域 (注册时记录)
	Reg_severid          string    `json:"reg_severid"`          // 註冊服务器 (注册时记录)
	Reg_loginchannel     string    `json:"reg_loginchannel"`     // 註冊渠道 (注册时记录，google，apple，garena，fb，line，guest)
	U_country            string    `json:"country"`              // 歸屬國家 (註冊時所在國家的国家地区代码(ISO 3166-1 alpha-2，即两位大写英文字母)，基本與 reg_country 相同，但可能透過客服換區)
	U_region             string    `json:"region"`               // 歸屬區域 (註冊時角色所在區服，基本與 reg_region 相同，但可能透過客服換區)
	Server_id            string    `json:"server_id"`            // 服务器 (更新时覆盖，合服后上报合服后服务器ID)
	Gop_open_id          string    `json:"gop_open_id"`          // gopenid (注册时记录)
	Role_id              string    `json:"role_id"`              // 角色id (注册时记录)
	Rolename             string    `json:"rolename"`             // 角色昵称 (更新时覆盖)
	App_version          string    `json:"app_version"`          // 客户端版本号 (更新时覆盖)
	Platform             string    `json:"platform"`             // 平台标识 (更新时覆盖，平台标识：android，IOS，Simulator、PC)
	Current_level        string    `json:"current_level"`        // 当前等级 (更新时覆盖)
	Highest_power        uint32    `json:"highest_power"`        // 战斗力 (歷史最高戰力)
	Current_vip_level    uint32    `json:"current_vip_level"`    // 当前VIP等级 (每次vip等级变动时设置)
	Current_monthpass    string    `json:"current_monthpass"`    // 当前月卡状态 (每次状态变更时设置)
	Main_task            string    `json:"main_task"`            // 当前主线任务栏编号 (每次主线任务进度变更时设置)
	Max_battle_id        string    `json:"max_battle_id"`        // 当前最大关卡 (每次主线任务进度变更时设置)
	Register_time        time.Time `json:"register_time"`        // 角色注册时间 (注册时记录)
	Server_time          time.Time `json:"server_time"`          // 开服时间 (注册时记录，所在服务器的开服时间)
	First_login_time     time.Time `json:"first_login_time"`     // 首次登录时间 (首次登陆时记录)
	Last_login_time      time.Time `json:"last_login_time"`      // 最后登录时间 (每次充值时覆盖原来记录)
	First_pay_time       time.Time `json:"first_pay_time"`       // 首次充值时间 (首次充值时记录)
	Last_pay_time        time.Time `json:"last_pay_time"`        // 最后充值时间 (每次充值时覆盖原来记录)
	Total_pay_amount     uint32    `json:"total_pay_amount"`     // 累计付费金额 (每次付费完成时更新覆盖)
	Total_pay_num        uint32    `json:"total_pay_num"`        // 累计付费次数 (每次付费完成时更新覆盖)
	Total_ad_num         uint32    `json:"total_ad_num"`         // 累计观看广告次数 (广告观看后更新覆盖，如果没有广告就不用这个字段)
	Continuous_login_day uint32    `json:"continuous_login_day"` // 连续登录天数 (每次登录时设置)
	Total_login_day      uint32    `json:"total_login_day"`      // 累计登录天数 (每次登录时更新覆盖)
	Tmp_diamond          int64     `json:"tmp_diamond"`          // 当前钻石数量 (每次登出时设置)
	Tmp_gold             int64     `json:"tmp_gold"`             // 当前金币 (每次登出时设置)
	Tmp_stamina          int64     `json:"tmp_stamina"`          // 当前体力 (每次登出时设置)
}

func TdaUpdateCommonUser(accountId, distinctId string, data *CommonUser) {
	if !Send() {
		return
	}
	// tda update common user
	go tools.GoSafe("TdaCommonUser", func() {
		if tdaDataMap, err := FlattenStructToMap(data); err != nil {
			log.Error("tdaStructToMap err", zap.Error(err), zap.Reflect("data", data))
		} else {
			if err = GetTa().UserSet(accountId, distinctId, tdaDataMap); err != nil {
				log.Error("tda common user err", zap.Error(err))
			}
		}
	})
}
