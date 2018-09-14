package Common

import (
	"math/rand"
	"strings"
	"time"
)

//	Defined format layout time
const (
	JP_DATE_H_MONTH_DAY      = "1月2日"
	JP_DATE_H_YEAR_MONTH_DAY = "2016年1月2日"
	JP_DATE_F_MONTH_DAY      = "01月02日"
	JP_DATE_F_YEAR_MONTH_DAY = "2016年01月02日"

	/** 日付書式(yyyymmdd) */
	DATE_FORMAT_YMD = "20060102"

	/** 日付書式(yyyy/mm/dd) */
	DATE_FORMAT_YMD_SLASH = "2006/01/02"

	/** 日付書式(yyyy-mm-dd) */
	DATE_FORMAT_YMD_SUBSTRACT = "2006-01-02"

	/** 日付書式(yyyy年mm月dd日) */
	DATE_FORMAT_YMD_SLASH_JP = "2006年01月02日"

	/** 日付書式(yyyymmddHHMMSSS) */
	DATE_FORMAT_YMDHMS = "20060102150405"

	/** 日付書式(yyyy-mm-dd HH:MM:SS) */
	DATE_FORMAT_MYSQL_YMDHMS = "2006-01-02 15:04:05"

	DATE_FORMAT_MD = "1/2"

	DATE_FORMAT_ZERO_TIME = "20060102000000"

	/** 日付書式(yyyy) */
	DATE_FORMAT_Y = "2006"

	/** 日付書式(yyyymm) */
	DATE_FORMAT_YM = "200601"

	/** 日付書式(yyyy/mm) */
	DATE_FORMAT_YM_SLASH = "2006/01"

	/** 日付書式(mm/dd) */
	DATE_FORMAT_MD_SLASH = "01/02"

	/** 日付書式(mm) */
	DATE_FORMAT_M = "01"

	/** */
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//CurrentDate return current date in format yyyymmdd
func CurrentDate() string {
	now := time.Now()
	return now.Format(DATE_FORMAT_YMD)
}

//DayFromToday date of (today + n day) in format yyyymmdd
func DayFromToday(day int) string {
	now := time.Now().AddDate(0, 0, day)
	return now.Format(DATE_FORMAT_YMD)
}

//CurrentDateTime return current date time in format yyyymmddHHMMSS
func CurrentDateTime() string {
	date_format := DATE_FORMAT_YMDHMS
	now := time.Now()
	return now.Format(date_format)
}

//CurrentMySqlDateTime return current date time in format yyyy-mm-dd HH:MM:SS
func CurrentMySqlDateTime() string {
	return time.Now().Format(DATE_FORMAT_MYSQL_YMDHMS)
}

//DateAddSlash convert yyyymmdd to yyyy/mm/dd
func DateAddSlash(yyyymmdd string) string {
	if len(yyyymmdd) >= 8 {
		return_year := yyyymmdd[0:4]
		return_month := yyyymmdd[4:6]
		return_day := yyyymmdd[6:8]
		return_date := return_year + "/" + return_month + "/" + return_day
		return return_date
	}
	return yyyymmdd
}

//DateAddHyphen convert yyyymmdd to yyyy-mm-dd
func DateAddHyphen(yyyymmdd string) string {
	if len(yyyymmdd) >= 8 {
		return_year := yyyymmdd[0:4]
		return_month := yyyymmdd[4:6]
		return_day := yyyymmdd[6:8]
		return_date := return_year + "-" + return_month + "-" + return_day
		return return_date
	}
	return yyyymmdd
}

//DateRemoveSlash convert date_search (yyyy/mm/dd) to yyyymmdd
func DateRemoveSlash(date_search string) string {
	if len(date_search) == 10 {
		date_search = strings.Replace(date_search, "/", "", -1)
	}
	return date_search
}

//DateToJPDate convert yyyymmdd to yyyy年mm月dd日
func DateToJPDate(yyyymmdd string) string {
	if len(yyyymmdd) < 8 {
		return ""
	}
	return_year := yyyymmdd[0:4]
	return_month := yyyymmdd[4:6]
	return_day := yyyymmdd[6:8]
	return_date := return_year + "年" + return_month + "月" + return_day + "日"
	return return_date
}

//DateAddYMD add year, month, day to date_search, date format: yyyymmdd
func DateAddYMD(yyyymmdd string, year int, month int, day int) string {
	t, err := time.Parse(DATE_FORMAT_YMD, yyyymmdd)
	if err != nil {
		return ""
	}

	return t.AddDate(year, month, day).Format(DATE_FORMAT_YMD)
}

//DateAddDay add nDate to yyyymmdd
func DateAddDay(yyyymmdd string, nDay int) string {
	date_format := DATE_FORMAT_YMD
	t, err := time.Parse(date_format, yyyymmdd)
	if err != nil {
		return yyyymmdd
	}
	return t.AddDate(0, 0, nDay).Format(date_format)
}

//DateAddMonth add nMonth to yyyymmdd
func DateAddMonth(yyyymmdd string, nMonth int) string {
	date_format := DATE_FORMAT_YMD
	t, err := time.Parse(date_format, yyyymmdd)
	if err != nil {
		return yyyymmdd
	}
	return t.AddDate(0, nMonth, 0).Format(date_format)
}

//DateAddYear add nYear to yyyymmdd
func DateAddYear(yyyymmdd string, nYear int) string {
	date_format := DATE_FORMAT_YMD
	t, err := time.Parse(date_format, yyyymmdd)
	if err != nil {
		return yyyymmdd
	}
	return t.AddDate(nYear, 0, 0).Format(date_format)
}

//GetSevenDayInWeek return list dates (DD) from (date_search - 6) to date_search
//date_search format yyyy/mm/dd
func GetSevenDayInWeek(date_search string) ([]string, error) {
	var day_in_week []string
	layout := DATE_FORMAT_YMD_SLASH
	t, err := time.Parse(layout, date_search)
	if err != nil {
		return nil, err
	}
	if len(date_search) > 5 {
		for i := 0; i < 7; i++ {
			days := t.AddDate(0, 0, -i).Format(DATE_FORMAT_YMD_SLASH)
			day := days[5:]
			day_in_week = append(day_in_week, day)
		}
	}

	return day_in_week, err
}

func GetTwoMonthInYear(date_search string) ([]string, error) {
	var month_in_year []string
	layout := DATE_FORMAT_YMD_SLASH
	t, err := time.Parse(layout, date_search)
	if err != nil {
		return nil, err
	}
	if len(date_search) > 5 {
		for i := 1; i <= 2; i++ {
			months := t.AddDate(0, -i, 0).Format(DATE_FORMAT_YMD_SLASH)
			month := months[5:7]
			month_in_year = append(month_in_year, month)
		}
	}

	return month_in_year, err
}

//IsValidateDate return true if date_search in format yyyy/mm/dd
func IsValidateDate(date_search string) bool {
	_, err := time.Parse(DATE_FORMAT_YMD_SLASH, date_search)
	return (err == nil)
}

//CheckDateFormat return time.Time and error when parse date_search in format yyyy/mm/dd
func CheckDateFormat(strDate string) (date time.Time, err error) {
	if date, err = time.Parse(DATE_FORMAT_YMD_SLASH, strDate); err != nil {
		date, err = time.Parse(DATE_FORMAT_YMD, strDate)
	}
	return
}

//IsBeforeDate return true if date_from <= date_to, format yyyy/mm/dd
func IsBeforeDate(date_from string, date_to string) (bool, error) {
	date_format := DATE_FORMAT_YMD_SLASH

	t_from, err := time.Parse(date_format, date_from)

	if err != nil {
		return false, err
	}

	t_to, err := time.Parse(date_format, date_to)

	if err != nil {
		return false, err
	}

	if t_from.Before(t_to) || t_from.Equal(t_to) {
		return true, nil
	} else {
		return false, nil
	}
}

//IsBeforeDate return true if date_from + (years, months, days) <= date_to, format yyyy/mm/dd
func IsBeforeDateAdd(date_from string, date_to string, years int, months int, days int) (bool, error) {
	date_format := DATE_FORMAT_YMD

	t_from, err := time.Parse(date_format, date_from)
	t_from = t_from.AddDate(years, months, days)

	if err != nil {
		return false, err
	}

	t_to, err := time.Parse(date_format, date_to)

	if err != nil {
		return false, err
	}

	if t_from.Before(t_to) || t_from.Equal(t_to) {
		return true, nil
	} else {
		return false, nil
	}
}

// Get current datetime : format default yyyymmddhhmmss
func CurrentDatetime() string {
	now := time.Now()
	return now.Format(DATE_FORMAT_YMDHMS)
}

// Get last datetime : format default yyyymmddhhmmss
func DatetimePast(days int) string {
	return DatetimePastWithFormat(days, DATE_FORMAT_YMDHMS)
}

// Get last datetime with format input
func DatetimePastWithFormat(days int, format string) string {
	past := time.Now().AddDate(0, 0, -days)
	return past.Format(format)
}

// Get future datetime : format default yyyymmddhhmmss
func DatetimeFuture(days int) string {
	return DatetimeFutureWithFormat(days, DATE_FORMAT_YMDHMS)
}

// Get future datetime with format input
func DatetimeFutureWithFormat(days int, format string) string {
	future := time.Now().AddDate(0, 0, days)
	return future.Format(format)
}
