package rbi

import "time"

type OnlineCnt struct {
	DtEventTime      time.Time // (必填) 格式 YYYY-MM-DD HH:MM:SS
	Login_channel    string    // (必填)登陆渠道
	Timekey          int64     // (必填)当前时间的时间戳，即当前的unixtime
	Country          string    // (必填)国家名称，使用中文
	CountryID        int       // (必填)国家ID(遵循iso 3166-1 2位标准)
	Zoneareaid       int       // (必填)分区id
	Onlinecntios     int       // (必填)ios在线人数
	Onlinecntandroid int       // (必填)android在线人数
	Onlinecntpc      int       // (必填)PC在线人数
	Onlinecntall     int       // (必填)全平台在线人数
	Onlinecnt        int       // (必填)在线人数，以上4个类型之和
	Registernum      int       // (必填)總註冊數
}

func (data *OnlineCnt) Name() string {
	return "OnlineCnt"
}
