package tda

import "time"

type ICommonAttr interface {
	SetEventName(eventName string)
	SetEventTime(now time.Time)
}

type CommonAttr struct {
	//Device_id       string    `json:"#device_id"`       // 设备id
	//Os              string    `json:"#os"`              // 操作系统
	//Os_version      string    `json:"#os_version"`      // 操作系统版本
	//App_version     string    `json:"#app_version"`     // app版本号
	//Manufacturer    string    `json:"#manufacturer"`    // 设备制造商
	//Device_model    string    `json:"#device_model"`    // 手机型号
	//Screen_height   string    `json:"#screen_height"`   // 屏幕高度
	//Screen_width    string    `json:"#screen_width"`    // 屏幕宽度
	//Ram             string    `json:"#ram"`             // 设备运行内存状态
	//Disk            string    `json:"#disk"`            // 设备存储空间状态
	//Network_type    string    `json:"#network_type"`    // 网络状态
	//Carrier         string    `json:"#carrier"`         // 网络运营商
	//Country         string    `json:"#country"`         // 国家地区
	//Country_code    string    `json:"#country_code"`    // 国家地区代码
	//System_language string    `json:"#system_language"` // 系统语言
	UserId     string    `json:"user_id"`     // 用户唯一 ID  用户在系统中的唯一用户标识 openId
	AccountId  string    `json:"account_id"`  // 账户 ID  相當於 roleID uid
	DistinctId string    `json:"distinct_id"` // 访客 ID
	EventName  string    `json:"event_name"`  // 事件名称
	EventTime  time.Time `json:"event_time"`  // 事件时间
	//ZoneOffset      string    `json:"#zone_offset"`     // 时区偏移
	//GopOpenId       string    `json:"gop_open_id"`      // gopenid  garena帳號 ID
	//UCountry        string    `json:"country"`          // 歸屬國家  註冊時所在國家的国家地区代码(ISO 3166-1 alpha-2，即两位大写英文字母)
	//URegion         string    `json:"region"`           // 歸屬區域  註冊時角色所在區服
	//SeverId         string    `json:"server_id"`        // 登入服务器id  登录服务ID，合服后上报合服后服务器ID
	//Platform        string    `json:"platform"`         // 平台标识  平台标识：android，IOS，Simulator、PC
	//LoginChannel    string    `json:"loginchannel"`     // 登入渠道  google，apple，garena，fb，line，guest
	//
	//Account        string    `json:"account"`          // 账号ID 雷霆SDK登录回调接口下返回的userId，雷霆官网下的数据形式为【8位数字+字母】的随机组合。
	//Channel        string    `json:"channel"`          // 渠道ID 事件发生当时的渠道ID，可通过雷霆SDK获取。
	//Ip             string    `json:"#ip"`              // 客户端ip地址,调用雷霆SDK【getReportBaselnfo("2")】方法，内含ipv4及ipv6的信息，按需选用。
	//BaseInfo       any       `json:"base_info"`        // 由雷霆SDK整合的从客户端获取的设备及包体信息，从SDK通过【getReportBaselnfo("2")】方法获取并将整个JSON拼接到日志中即可。
	//RoleCreateTime time.Time `json:"role_create_time"` // 角色成功创建的时间。
	//RoleLevel      int       `json:"rolelevel"`        // 角色当前等级。
	//RoleName       string    `json:"rolename"`         // 角色当前昵称。
	//RolePaid       float64   `json:"rolepaid"`         // 角色创建至今累计付费金额。
	//Media          string   `json:"media"`            // 媒体编号
}

func (c *CommonAttr) SetEventName(eventName string) {
	c.EventName = eventName
}

func (c *CommonAttr) SetEventTime(now time.Time) {
	c.EventTime = now
}
