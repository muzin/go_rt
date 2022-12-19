package datetime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 时间解析字符串
const (
	yearTimeParseConstant        = "2006"
	monthTimeParseConstant       = "01"
	dayTimeParseConstant         = "02"
	hourTimeParseConstant        = "15"
	minuteTimeParseConstant      = "04"
	secondTimeParseConstant      = "05"
	millisecondTimeParseConstant = "000"
	microsecondTimeParseConstant = "000000"
	nanosecondTimeParseConstant  = "000000000"
)

/**
 * 将日期格式化成指定模板的字符串
 * @param date 日期对象
 * @param format 日期格式化的模板
 *   format 支持一下格式:
 *     yyyy - 年
 *     MM - 月
 *     dd - 日
 *     hh - 时
 *     mm - 分
 *     ss - 秒
 *     ms - 毫秒
 * @example
 *   times := "2020-09-18 15:04:05"
 *	 s, _ := time.Parse("2006-01-02 15:04:05", times)
 */
func Format(date time.Time, format string) string {

	year := date.Year()
	month := date.Month()
	day := date.Day()
	hour := date.Hour()
	minute := date.Minute()
	second := date.Second()
	millisecond := date.Nanosecond() / 1e6
	nanosecond := date.Nanosecond()

	var newStr = format

	yearStr := strconv.Itoa(year)
	monthStr := strconv.Itoa(int(month))
	dayStr := strconv.Itoa(day)
	hourStr := strconv.Itoa(hour)
	minuteStr := strconv.Itoa(minute)
	secondStr := strconv.Itoa(second)
	millisecondStr := strconv.Itoa(millisecond)
	nanosecondStr := strconv.Itoa(nanosecond)

	if len(monthStr) == 1 {
		monthStr = "0" + monthStr
	}
	if len(dayStr) == 1 {
		dayStr = "0" + dayStr
	}
	if len(hourStr) == 1 {
		hourStr = "0" + hourStr
	}
	if len(minuteStr) == 1 {
		minuteStr = "0" + minuteStr
	}
	if len(secondStr) == 1 {
		secondStr = "0" + secondStr
	}
	//if len(millisecondStr) == 1 	{ millisecondStr = "00" + millisecondStr
	//}else if len(millisecondStr) == 2 	{ millisecondStr = "0" + millisecondStr }

	millisecondStr = fmt.Sprintf("%03v", millisecond)
	nanosecondStr = fmt.Sprintf("%09v", nanosecond)

	newStr = strings.ReplaceAll(newStr, "yyyy", yearStr)
	newStr = strings.ReplaceAll(newStr, "MM", monthStr)
	newStr = strings.ReplaceAll(newStr, "dd", dayStr)
	newStr = strings.ReplaceAll(newStr, "hh", hourStr)
	newStr = strings.ReplaceAll(newStr, "mm", minuteStr)
	newStr = strings.ReplaceAll(newStr, "ss", secondStr)
	newStr = strings.ReplaceAll(newStr, "ms", millisecondStr)
	newStr = strings.ReplaceAll(newStr, "ns", nanosecondStr)

	return newStr
}

func Parse(datestr string, format string) (time.Time, error) {

	var layout = format

	layout = strings.ReplaceAll(layout, "yyyy", yearTimeParseConstant)
	layout = strings.ReplaceAll(layout, "MM", monthTimeParseConstant)
	layout = strings.ReplaceAll(layout, "dd", dayTimeParseConstant)
	layout = strings.ReplaceAll(layout, "hh", hourTimeParseConstant)
	layout = strings.ReplaceAll(layout, "mm", minuteTimeParseConstant)
	layout = strings.ReplaceAll(layout, "ss", secondTimeParseConstant)
	layout = strings.ReplaceAll(layout, "ms", millisecondTimeParseConstant)

	datetime, err := time.Parse(layout, datestr)

	return datetime, err
}
