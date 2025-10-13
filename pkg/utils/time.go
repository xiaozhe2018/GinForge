package utils

import (
	"fmt"
	"time"
)

// TimeFormat 时间格式常量
const (
	TimeFormatDefault     = "2006-01-02 15:04:05"
	TimeFormatDate        = "2006-01-02"
	TimeFormatTime        = "15:04:05"
	TimeFormatDateTime    = "2006-01-02 15:04:05"
	TimeFormatISO8601     = "2006-01-02T15:04:05Z"
	TimeFormatRFC3339     = "2006-01-02T15:04:05Z07:00"
	TimeFormatTimestamp   = "2006-01-02 15:04:05.000"
	TimeFormatChinese     = "2006年01月02日 15:04:05"
	TimeFormatChineseDate = "2006年01月02日"
)

// Now 获取当前时间
func Now() time.Time {
	return time.Now()
}

// NowUnix 获取当前时间戳
func NowUnix() int64 {
	return time.Now().Unix()
}

// NowUnixNano 获取当前纳秒时间戳
func NowUnixNano() int64 {
	return time.Now().UnixNano()
}

// NowMilli 获取当前毫秒时间戳
func NowMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// FormatTime 格式化时间
func FormatTime(t time.Time, format string) string {
	if format == "" {
		format = TimeFormatDefault
	}
	return t.Format(format)
}

// FormatNow 格式化当前时间
func FormatNow(format string) string {
	return FormatTime(Now(), format)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr, format string) (time.Time, error) {
	if format == "" {
		format = TimeFormatDefault
	}
	return time.Parse(format, timeStr)
}

// ParseTimeInLocation 在指定时区解析时间
func ParseTimeInLocation(timeStr, format, location string) (time.Time, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}

	if format == "" {
		format = TimeFormatDefault
	}

	return time.ParseInLocation(format, timeStr, loc)
}

// UnixToTime 时间戳转时间
func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// UnixNanoToTime 纳秒时间戳转时间
func UnixNanoToTime(unixNano int64) time.Time {
	return time.Unix(0, unixNano)
}

// TimeToUnix 时间转时间戳
func TimeToUnix(t time.Time) int64 {
	return t.Unix()
}

// TimeToUnixNano 时间转纳秒时间戳
func TimeToUnixNano(t time.Time) int64 {
	return t.UnixNano()
}

// TimeToMilli 时间转毫秒时间戳
func TimeToMilli(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// AddDays 添加天数
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddMonths 添加月数
func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// AddYears 添加年数
func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

// AddHours 添加小时
func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// AddMinutes 添加分钟
func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// AddSeconds 添加秒数
func AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

// StartOfDay 获取一天的开始时间
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay 获取一天的结束时间
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// StartOfWeek 获取一周的开始时间
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // 周日为0，转换为7
	}
	return StartOfDay(t.AddDate(0, 0, -weekday+1))
}

// EndOfWeek 获取一周的结束时间
func EndOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // 周日为0，转换为7
	}
	return EndOfDay(t.AddDate(0, 0, 7-weekday))
}

// StartOfMonth 获取一月的开始时间
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 获取一月的结束时间
func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	nextMonth := month + 1
	if nextMonth > 12 {
		nextMonth = 1
		year++
	}
	return time.Date(year, nextMonth, 1, 0, 0, 0, 0, t.Location()).Add(-time.Nanosecond)
}

// StartOfYear 获取一年的开始时间
func StartOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 获取一年的结束时间
func EndOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year+1, 1, 1, 0, 0, 0, 0, t.Location()).Add(-time.Nanosecond)
}

// IsSameDay 检查是否为同一天
func IsSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsSameWeek 检查是否为同一周
func IsSameWeek(t1, t2 time.Time) bool {
	s1 := StartOfWeek(t1)
	s2 := StartOfWeek(t2)
	return s1.Equal(s2)
}

// IsSameMonth 检查是否为同一月
func IsSameMonth(t1, t2 time.Time) bool {
	y1, m1, _ := t1.Date()
	y2, m2, _ := t2.Date()
	return y1 == y2 && m1 == m2
}

// IsSameYear 检查是否为同一年
func IsSameYear(t1, t2 time.Time) bool {
	y1, _, _ := t1.Date()
	y2, _, _ := t2.Date()
	return y1 == y2
}

// DaysBetween 计算两个时间之间的天数
func DaysBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Hours() / 24)
}

// HoursBetween 计算两个时间之间的小时数
func HoursBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Hours())
}

// MinutesBetween 计算两个时间之间的分钟数
func MinutesBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Minutes())
}

// SecondsBetween 计算两个时间之间的秒数
func SecondsBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Seconds())
}

// IsLeapYear 检查是否为闰年
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DaysInMonth 获取某月的天数
func DaysInMonth(year, month int) int {
	if month < 1 || month > 12 {
		return 0
	}

	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	if month == 2 && IsLeapYear(year) {
		return 29
	}

	return daysInMonth[month-1]
}

// GetTimezone 获取时区
func GetTimezone(t time.Time) string {
	return t.Location().String()
}

// SetTimezone 设置时区
func SetTimezone(t time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

// IsWeekend 检查是否为周末
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekday 检查是否为工作日
func IsWeekday(t time.Time) bool {
	return !IsWeekend(t)
}

// GetWeekday 获取星期几
func GetWeekday(t time.Time) string {
	weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	return weekdays[t.Weekday()]
}

// GetWeekdayEnglish 获取英文星期几
func GetWeekdayEnglish(t time.Time) string {
	return t.Weekday().String()
}

// GetMonth 获取月份
func GetMonth(t time.Time) string {
	months := []string{"一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"}
	return months[t.Month()-1]
}

// GetMonthEnglish 获取英文月份
func GetMonthEnglish(t time.Time) string {
	return t.Month().String()
}

// HumanizeDuration 人性化显示时间间隔
func HumanizeDuration(duration time.Duration) string {
	if duration < time.Minute {
		return "刚刚"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d分钟前", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d小时前", hours)
	} else if duration < 30*24*time.Hour {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d天前", days)
	} else if duration < 365*24*time.Hour {
		months := int(duration.Hours() / (24 * 30))
		return fmt.Sprintf("%d个月前", months)
	} else {
		years := int(duration.Hours() / (24 * 365))
		return fmt.Sprintf("%d年前", years)
	}
}

// HumanizeTime 人性化显示时间
func HumanizeTime(t time.Time) string {
	now := Now()
	duration := now.Sub(t)
	return HumanizeDuration(duration)
}
