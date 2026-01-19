package tools

import (
	"fmt"
	"time"

	"github.com/zy/game_data/template"
)

func GetCurTime() uint32 {
	return uint32(time.Now().Unix())
}

func GetCurMill() int64 {
	return time.Now().UnixMilli()
}

func GetMaxTime() uint32 {
	return uint32(time.Date(2037, 1, 1, 0, 0, 0, 0, &time.Location{}).Unix())
}

func GetDailyRefreshTime() uint32 {
	return GetHourRefreshTime(template.GetSystemItemTemplate().RefreshHour)
}

func GetDailyXRefreshTime(day uint32, hour uint32) uint32 {
	if day == 0 {
		return 0
	}
	if day > 1 {
		temp := getHourRefreshTime(hour)
		return uint32(temp.AddDate(0, 0, int(day-1)).Unix())
	}
	return GetHourRefreshTime(hour)
}

func GetHourRefreshTime(hour uint32) uint32 {
	return uint32(getHourRefreshTime(hour).Unix())
}

func GetTotalMin(hour, min uint32) uint32 {
	return hour*60 + min
}

func GetHourMinRefreshTime(hour, min uint32) uint32 {
	t := time.Now()
	var re time.Time
	curTotal := GetTotalMin(uint32(t.Hour()), uint32(t.Minute()))
	targetTotal := GetTotalMin(hour, min)
	if curTotal < targetTotal {
		re = time.Date(t.Year(), time.Month(t.Month()), t.Day(), int(hour), int(min), 0, 0, time.Local)
	} else {
		re = time.Date(t.Year(), time.Month(t.Month()), t.Day()+1, int(hour), int(min), 0, 0, time.Local)
	}
	return uint32(re.Unix())
}

func getHourRefreshTime(hour uint32) time.Time {
	t := time.Now()
	var re time.Time
	refreshHour := int(hour)
	if t.Hour() < refreshHour {
		re = time.Date(t.Year(), time.Month(t.Month()), t.Day(), refreshHour, 0, 0, 0, time.Local)
	} else {
		re = time.Date(t.Year(), time.Month(t.Month()), t.Day()+1, refreshHour, 0, 0, 0, time.Local)
	}
	return re
}

func GetWeeklyRefreshTime(hour uint32) uint32 {
	return uint32(getNextFirstDateOfWeek(hour).Unix())
}

func getFirstDateOfWeek(hour uint32) time.Time {
	t := time.Now()
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}

	refreshHour := int(hour)
	if offset == 0 && t.Hour() < refreshHour {
		offset = -7
	}
	return time.Date(t.Year(), t.Month(), t.Day(), refreshHour, 0, 0, 0, time.Local).
		AddDate(0, 0, offset)
}

// getNextFirstDateOfWeek 获取下周周一
func getNextFirstDateOfWeek(hour uint32) time.Time {
	return getFirstDateOfWeek(hour).
		AddDate(0, 0, 7)
}

func GetMonthRefreshTime() uint32 {
	t := time.Now()
	var re time.Time
	refreshHour := int(template.GetSystemItemTemplate().RefreshHour)
	if t.Day() == 1 {
		if t.Hour() < refreshHour {
			re = time.Date(t.Year(), time.Month(t.Month()), 1, refreshHour, 0, 0, 0, time.Local)
			return uint32(re.Unix())
		}
	}
	re = time.Date(t.Year(), time.Month(t.Month()+1), 1, refreshHour, 0, 0, 0, time.Local)
	return uint32(re.Unix())
}

func GetCurDateString() string {
	t := time.Now()
	return fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
}

func GetDateStart(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Local")                                                                                        //重要：获取时区
	theTime, _ := time.ParseInLocation(time.DateTime, fmt.Sprintf("%4d-%02d-%02d 00:00:00", t.Year(), t.Month(), t.Day()), loc) //使用模板在对应时区转化为time.time类型
	return theTime
}

func GetDateEnd(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Local")                                                                                        //重要：获取时区
	theTime, _ := time.ParseInLocation(time.DateTime, fmt.Sprintf("%4d-%02d-%02d 23:59:59", t.Year(), t.Month(), t.Day()), loc) //使用模板在对应时区转化为time.time类型
	return theTime
}

func GetCurDate() time.Time {
	curTime := time.Now()
	loc, _ := time.LoadLocation("Local")                                                                                                          //重要：获取时区
	theTime, _ := time.ParseInLocation(time.DateTime, fmt.Sprintf("%4d-%02d-%02d 00:00:00", curTime.Year(), curTime.Month(), curTime.Day()), loc) //使用模板在对应时区转化为time.time类型
	return theTime
}

func GetDateFromStr(d string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local") //重要：获取时区
	return time.ParseInLocation(time.DateOnly, d, loc)
}

func GetCurDateStart() uint32 {
	//转化为时间戳 类型是int64
	return uint32(GetDateStart(time.Now()).Unix())
}

func GetStrDateTime(curTime uint32) string {
	if curTime == 0 {
		return ""
	}
	timestamp := int64(curTime)
	loc, _ := time.LoadLocation("Local")
	t := time.Unix(timestamp, 0).In(loc) // 转换为北京时间
	return t.Format(time.DateTime)       // 输出格式化后的时间，例如：2021-04-12 17:41:03                             // 输出格式化后的时间，例如：2021-04-12 17:41:03
}

func GetCurDateEnd() uint32 {
	curTime := time.Now()
	loc, _ := time.LoadLocation("Local")                                                                                                          //重要：获取时区
	theTime, _ := time.ParseInLocation(time.DateTime, fmt.Sprintf("%4d-%02d-%02d 23:59:59", curTime.Year(), curTime.Month(), curTime.Day()), loc) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()                                                                                                                          //转化为时间戳 类型是int64
	return uint32(sr)
}

func GetDiffDay(start, end time.Time) uint32 {
	d1 := GetDateStart(start)
	d2 := GetDateStart(end)
	return uint32(d2.Sub(d1).Hours() / 24)
}

// GetDailyStartTime 获得每天开始时间
func GetDailyStartTime(starttime uint32, refreshHour uint32) uint32 {
	t := time.Unix(int64(starttime), 0)
	var re time.Time
	re = time.Date(t.Year(), time.Month(t.Month()), t.Day(), int(refreshHour), 0, 0, 0, time.Local)
	return uint32(re.Unix())
}

type DailyOnline struct {
	Day     uint32
	Seconds uint32
}

func GetDate(data uint32) time.Time {
	t := time.Unix(int64(data), 0)
	var re time.Time
	re = time.Date(t.Year(), time.Month(t.Month()), t.Day(), 0, 0, 0, 0, time.Local)
	return re
}

func GetStrDate(curTime uint32) string {
	if curTime == 0 {
		return ""
	}
	timestamp := int64(curTime)
	loc, _ := time.LoadLocation("Local")
	t := time.Unix(timestamp, 0).In(loc) // 转换为北京时间
	return t.Format(time.DateOnly)       // 输出格式化后的时间，例如：2021-04-12 17:41:03                             // 输出格式化后的时间，例如：2021-04-12 17:41:03
}

// GetStaticTime 获得统计时间
func GetStaticTime(data uint32) uint32 {
	temp := GetDailyStartTime(data, 0)
	return uint32(GetDate(temp).Unix())
}

func CalcOnlineTime(start, end uint32) []*DailyOnline {
	var ret []*DailyOnline
	t1 := GetDailyStartTime(start, 0)

	for {
		// 下一天开始时间
		t1 += 24 * 3600
		if t1 > end {
			break
		}
		ret = append(ret, &DailyOnline{
			Day:     uint32(GetDate(start).Unix()),
			Seconds: t1 - start,
		})
		start = t1
	}

	// 剩下的start end 在同一天时间
	if start < end {
		ret = append(ret, &DailyOnline{
			Day:     uint32(GetDate(start).Unix()),
			Seconds: end - start,
		})
	}
	return ret
}

// GetDateFromTimestamp 从时间戳获取日期字符串，格式为YYYYMMDD
func GetDateFromTimestamp(timestamp uint32) string {
	t := time.Unix(int64(timestamp), 0)
	return fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
}

// 获取今天早上零点
func GetTodayZeroTime() uint32 {
	now := time.Now()
	// 获取当天零点时间
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 转换为时间戳
	timestamp := midnight.Unix()
	return uint32(timestamp)
}

// 20250328	YMD格式
func GetDateIntByOffset(t time.Time, offset int) int {
	y, m, d := t.Add(-time.Hour * time.Duration(offset)).Date()
	date := y*10000 + int(m)*100 + d
	return date
}

// 202515	第几周
func GetYearWeekByOffset(t time.Time, offset int) int {
	year, week := t.Add(-time.Hour * time.Duration(offset)).ISOWeek()
	date := year*100 + week
	return date
}

// 获取utc时间偏移小时
func GetTimeOffsetHour() int {
	_, offset := time.Now().Zone()
	return offset / 3600
}

func GetCurrentTimezoneOffset() int {
	_, offset := time.Now().Zone()
	offsetHours := offset / 3600

	if offsetHours >= 0 {
		return offsetHours
	}
	return -offsetHours
}

func GetWeekCount(timestamp int64) int {
	var startTs int64 = timestamp

	startTime := time.Unix(startTs, 0)
	currentTime := time.Now()

	if startTime.After(currentTime) {
		return 0
	}

	var startDaysToMonday int
	switch startTime.Weekday() {
	case time.Sunday:
		startDaysToMonday = -6
	case time.Monday:
		startDaysToMonday = 0
	case time.Tuesday:
		startDaysToMonday = -1
	case time.Wednesday:
		startDaysToMonday = -2
	case time.Thursday:
		startDaysToMonday = -3
	case time.Friday:
		startDaysToMonday = -4
	case time.Saturday:
		startDaysToMonday = -5
	}

	startMonday := startTime.AddDate(0, 0, startDaysToMonday)
	startMonday = time.Date(startMonday.Year(), startMonday.Month(), startMonday.Day(), 0, 0, 0, 0, startMonday.Location())

	startSunday := startMonday.AddDate(0, 0, 6)

	var currentDaysToMonday int
	switch currentTime.Weekday() {
	case time.Sunday:
		currentDaysToMonday = -6
	case time.Monday:
		currentDaysToMonday = 0
	case time.Tuesday:
		currentDaysToMonday = -1
	case time.Wednesday:
		currentDaysToMonday = -2
	case time.Thursday:
		currentDaysToMonday = -3
	case time.Friday:
		currentDaysToMonday = -4
	case time.Saturday:
		currentDaysToMonday = -5
	}

	currentMonday := currentTime.AddDate(0, 0, currentDaysToMonday)
	currentMonday = time.Date(currentMonday.Year(), currentMonday.Month(), currentMonday.Day(), 0, 0, 0, 0, currentMonday.Location())

	daysBetweenMondays := currentMonday.Sub(startMonday).Hours() / 24
	completeWeeks := int(daysBetweenMondays / 7)

	totalWeeks := completeWeeks + 1

	if currentTime.Before(startSunday) && currentTime.After(startTime) {
		totalWeeks = 1
	}

	return totalWeeks
}

// ===================== stopwatch =========================//
type Stopwatch struct {
	startTime   time.Time
	accumulated time.Duration
	isRunning   bool
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{
		accumulated: 0,
		isRunning:   false,
	}
}

func (s *Stopwatch) Start() {
	if !s.isRunning {
		s.startTime = time.Now()
		s.isRunning = true
	}
}

func (s *Stopwatch) Pause() {
	if s.isRunning {
		s.accumulated += time.Since(s.startTime)
		s.isRunning = false
	}
}

func (s *Stopwatch) Reset() {
	s.accumulated = 0
	s.isRunning = false
}

func (s *Stopwatch) Elapsed() time.Duration {
	if s.isRunning {
		return s.accumulated + time.Since(s.startTime)
	}
	return s.accumulated
}

func GetDailyRefreshTimeByHour(hour uint32) uint32 {
	return GetHourRefreshTime(hour)
}

func GetMonthRefreshTimeByHour(refreshHour int) uint32 {
	t := time.Now()
	var re time.Time
	if t.Day() == 1 {
		if t.Hour() < refreshHour {
			re = time.Date(t.Year(), time.Month(t.Month()), 1, refreshHour, 0, 0, 0, time.Local)
			return uint32(re.Unix())
		}
	}
	re = time.Date(t.Year(), time.Month(t.Month()+1), 1, refreshHour, 0, 0, 0, time.Local)
	return uint32(re.Unix())
}

// 获取间隔天数时间
func GetOffsetDays(offset int) time.Time {
	return time.Now().Add(time.Hour * time.Duration(offset*24))
}
