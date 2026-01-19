package rdb

import (
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"time"
)

const (
	// 一分钟的秒数
	TD_OneMinuteSecond int64 = 60
	// 一小时的秒数
	TD_OneHourSecond int64 = 60 * 60
	// 一天的秒数
	TD_OneDaySecond int64 = 60 * 60 * 24
	// 一周的秒数
	TD_OneWeekSecond int64 = 60 * 60 * 24 * 7
)

// Second int64 转成秒数
func Second(t int64) time.Duration {
	return time.Second * time.Duration(t)
}

// GetSecond 根据传入的小时 分钟 秒数 计算成秒数
func GetSecond(hour, minute, second int64) int64 {
	return (hour * TD_OneHourSecond) + (minute * TD_OneMinuteSecond) + second
}

// ParseTime 解析时间字符串为时间戳
func ParseTime(src string) int64 {
	TimeLoc := time.Local
	sTime, err := time.ParseInLocation("2006-01-02-15-04-05", src, TimeLoc)
	if err != nil {
		log.Error("ParseTime error: %v", zap.Error(err))
		return 0
	}

	return sTime.Unix()
}
