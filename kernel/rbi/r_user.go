package rbi

import (
	"time"
)

// 玩家註冊資訊 玩家註冊時只會寫一筆（終生一筆）
type PlayerRegister struct {
	GameSvrId      string    // (必填)登录服务ID，分区服场景使用，合服后上报合服后服务器ID。
	DtEventTime    time.Time // (必填)游戏事件的时间, 格式 YYYY-MM-DD HH:MM:SS
	VGameAppid     string    // (必填)游戏APPID
	PlatID         int       // (必填)ios 0 /android 1/ PC 2 /harmonyos 3 /unknown -1
	IZoneAreaID    int       // (必填)注册服ID
	VOpenID        string    // (必填)用户OPENID号
	VRoleID        string    // (必填)玩家角色ID
	VRoleName      string    // (必填)玩家角色名
	VClientIP      string    // (必填)客户端IP(后台服务器记录与玩家通信时的IP地址)
	Region         string    // (必填)玩家註冊時登入的所屬地區，此欄位不會因為client ip而異動，註冊時為TW，登入登出時，也需打印TW
	Country        string    // (必填)ClientIP所在国家和地区
	GarenaOpenID   string    // (必填)Garena UID
	Timekey        int64     // (必填)当前时间的时间戳，即当前的unixtime
	ClientVersion  string    // 客户端版本
	SystemSoftware string    // 移动终端操作系统版本
	SystemHardware string    // 移动终端机型
	TelecomOper    string    // (必填)运营商
	Network        string    // 3G/WIFI/2G
	ScreenWidth    int       // 显示屏宽度
	ScreenHight    int       // 显示屏高度
	Density        float64   // 像素密度
	CpuHardware    string    // cpu类型|频率|核数
	Memory         int       // 内存信息(MB)
	GLRender       string    // opengl render信息
	GLVersion      string    // opengl版本信息
	DeviceId       string    // 设备ID,安卓上报IMEI(正确值为15-17位),IOS上报IDFA(正确值为36位)(报原始信息,不要加密)
	GenderType     int       // 性别 参考
}

func (data *PlayerRegister) Name() string {
	return "PlayerRegister"
}

// 玩家登入資訊 玩家有操作就會輸出一條
type PlayerLogin struct {
	GameSvrId      string    // (必填)登录服务ID，分区服场景使用，合服后上报合服后服务器ID。
	DtEventTime    time.Time // (必填)游戏事件的时间, 格式 YYYY-MM-DD HH:MM:SS
	VGameAppid     string    // (必填)游戏APPID
	PlatID         int       // (必填)ios 0 /android 1/ PC 2
	IZoneAreaID    int       // (必填)注册服ID
	VOpenID        string    // (必填)用户OPENID号
	VRoleID        string    // (必填)玩家角色ID
	VRoleName      string    // (必填)玩家角色名
	VClientIP      string    // (必填)客户端IP(后台服务器记录与玩家通信时的IP地址)
	Region         string    // (必填)玩家註冊時登入的所屬地區，此欄位不會因為client ip而異動，註冊時為TW，登入登出時，也需打印TW
	Country        string    // (必填)ClientIP所在国家和地区
	GarenaOpenID   string    // (必填)Garena UID
	Timekey        int64     // (必填)当前时间的时间戳，即当前的unixtime
	ClientVersion  string    // 客户端版本
	SystemSoftware string    // 移动终端操作系统版本
	SystemHardware string    // 移动终端机型
	TelecomOper    string    // (必填)运营商
	Network        string    // 3G/WIFI/2G
	ScreenWidth    int       // 显示屏宽度
	ScreenHight    int       // 显示屏高度
	Density        float64   // 像素密度
	CpuHardware    string    // cpu类型|频率|核数
	Memory         int       // 内存信息(MB)
	GLRender       string    // opengl render信息
	GLVersion      string    // opengl版本信息
	DeviceId       string    // 设备ID,安卓上报IMEI(正确值为15-17位),IOS上报IDFA(正确值为36位)(报原始信息,不要加密)
	GenderType     int       // 性别 参考
	ILevel         int       // (必填)等级
	RegisterTime   time.Time // 注册时间
	RoleTotalCash  int       // 玩家角色累计充值(元)
}

func (data *PlayerLogin) Name() string {
	return "PlayerLogin"
}

// 玩家登出資訊 玩家有操作就會輸出一條
type PlayerLogout struct {
	GameSvrId      string    // (必填)登录服务ID，分区服场景使用，合服后上报合服后服务器ID。
	DtEventTime    time.Time // (必填)游戏事件的时间, 格式 YYYY-MM-DD HH:MM:SS
	VGameAppid     string    // (必填)游戏APPID
	PlatID         int       // (必填)ios 0 /android 1/ PC 2
	IZoneAreaID    int       // (必填)注册服ID
	VOpenID        string    // (必填)用户OPENID号
	VRoleID        string    // (必填)玩家角色ID
	VRoleName      string    // (必填)玩家角色名
	VClientIP      string    // (必填)客户端IP(后台服务器记录与玩家通信时的IP地址)
	Region         string    // (必填)玩家註冊時登入的所屬地區，此欄位不會因為client ip而異動，註冊時為TW，登入登出時，也需打印TW
	Country        string    // (必填)ClientIP所在国家和地区
	GarenaOpenID   string    // (必填)Garena UID
	Timekey        int64     // (必填)当前时间的时间戳，即当前的unixtime
	ClientVersion  string    // 客户端版本
	SystemSoftware string    // 移动终端操作系统版本
	SystemHardware string    // 移动终端机型
	TelecomOper    string    // (必填)运营商
	Network        string    // 3G/WIFI/2G
	ScreenWidth    int       // 显示屏宽度
	ScreenHight    int       // 显示屏高度
	Density        float64   // 像素密度
	CpuHardware    string    // cpu类型|频率|核数
	Memory         int       // 内存信息(MB)
	GLRender       string    // opengl render信息
	GLVersion      string    // opengl版本信息
	DeviceId       string    // 设备ID,安卓上报IMEI(正确值为15-17位),IOS上报IDFA(正确值为36位)(报原始信息,不要加密)
	GenderType     int       // 性别 参考
	ILevel         int       // (必填)等级
	RegisterTime   time.Time // 注册时间
	RoleTotalCash  int       // 玩家角色累计充值(元)
}

func (data *PlayerLogout) Name() string {
	return "PlayerLogout"
}
