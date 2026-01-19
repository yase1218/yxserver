package tapping

import (
	"kernel/tda"
	"time"
)

type Login struct {
	*tda.CommonAttr
	AccountId     string    `bson:"accountid"`                              // 账号(平台openid)
	RoleName      string    `json:"rolename" bson:"rolename"`               // 角色昵称
	AddDay        uint32    `json:"add_day" bson:"add_day"`                 // 累积登陆天数
	LastOutTime   time.Time `json:"last_out_time"  bson:"last_out_time"`    // 上次退出游戏时间
	Action        int       `json:"action" bson:"action"`                   // 登录类型 1正常登录；2顶号登录；3断线重连；4在线跨天
	IsFirst       int       `json:"is_first" bson:"is_first"`               // 是否账号首次登录, 1是；0否
	LastLoginTime time.Time `json:"last_login_time" bson:"last_login_time"` // 上次登录时间
}
