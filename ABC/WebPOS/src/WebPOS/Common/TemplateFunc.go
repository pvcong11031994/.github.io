package Common

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func formatNumber3(str string) string {
	strFormatted := ""
	if str[0] == '-' {
		strFormatted = "-"
		str = str[1:]
	}
	if len(str) > 3 {
		firstIndex := len(str) % 3
		if firstIndex > 0 {
			strFormatted += str[:firstIndex] + ","
			str = str[firstIndex:]
		}
		for len(str) > 3 {
			strFormatted += str[:3] + ","
			str = str[3:]
		}
	}
	strFormatted += str
	return strFormatted
}

func FormatNumber(obj interface{}) string {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.String,
		reflect.Float32,
		reflect.Float64:
		str := fmt.Sprintf("%v", obj)
		if _, err := strconv.ParseInt(str, 10, 64); err == nil {
			return formatNumber3(str)
		}
	}

	if obj != nil {
		return fmt.Sprintf("%v", obj)
	}

	return ""
}

func Minus(a interface{}, b interface{}) interface{} {
	if reflect.TypeOf(a).Kind() == reflect.Float32 ||
		reflect.TypeOf(b).Kind() == reflect.Float32 ||
		reflect.TypeOf(a).Kind() == reflect.Float64 ||
		reflect.TypeOf(b).Kind() == reflect.Float64 {

		fA := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
		fB := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
		return fA - fB
	}

	iA := reflect.ValueOf(a).Convert(reflect.TypeOf(int64(0))).Int()
	iB := reflect.ValueOf(b).Convert(reflect.TypeOf(int64(0))).Int()
	return iA - iB
}

func Plus(a interface{}, b interface{}) interface{} {
	if reflect.TypeOf(a).Kind() == reflect.Float32 ||
		reflect.TypeOf(b).Kind() == reflect.Float32 ||
		reflect.TypeOf(a).Kind() == reflect.Float64 ||
		reflect.TypeOf(b).Kind() == reflect.Float64 {

		fA := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
		fB := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
		return fA + fB
	}

	iA := reflect.ValueOf(a).Convert(reflect.TypeOf(int64(0))).Int()
	iB := reflect.ValueOf(b).Convert(reflect.TypeOf(int64(0))).Int()
	return iA + iB
}

func Multiply(a interface{}, b interface{}) interface{} {
	if reflect.TypeOf(a).Kind() == reflect.Float32 ||
		reflect.TypeOf(b).Kind() == reflect.Float32 ||
		reflect.TypeOf(a).Kind() == reflect.Float64 ||
		reflect.TypeOf(b).Kind() == reflect.Float64 {

		fA := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
		fB := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
		return fA * fB
	}

	iA := reflect.ValueOf(a).Convert(reflect.TypeOf(int64(0))).Int()
	iB := reflect.ValueOf(b).Convert(reflect.TypeOf(int64(0))).Int()
	return iA * iB
}

func Divide(a interface{}, b interface{}) interface{} {
	if reflect.TypeOf(a).Kind() == reflect.Float32 ||
		reflect.TypeOf(b).Kind() == reflect.Float32 ||
		reflect.TypeOf(a).Kind() == reflect.Float64 ||
		reflect.TypeOf(b).Kind() == reflect.Float64 {

		fA := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
		fB := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
		return fA / fB
	}

	iA := reflect.ValueOf(a).Convert(reflect.TypeOf(int64(0))).Int()
	iB := reflect.ValueOf(b).Convert(reflect.TypeOf(int64(0))).Int()
	return iA / iB
}

func MakeArray(from interface{}, to interface{}) []int64 {
	iFrom := reflect.ValueOf(from).Convert(reflect.TypeOf(int64(0))).Int()
	iTo := reflect.ValueOf(to).Convert(reflect.TypeOf(int64(0))).Int()
	arr := make([]int64, iTo-iFrom+1)
	for i := range arr {
		arr[i] = iFrom + int64(i)
	}

	return arr
}

func CheckNumber(f interface{}) bool {
	var isNum = false
	switch reflect.ValueOf(f).Kind() {
	case reflect.String:
		if _, err := strconv.Atoi(f.(string)); err == nil {
			isNum = true
		} else {
			isNum = false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		isNum = true
	default:
		isNum = false
	}
	return isNum
}

//Check Sunday. Return 1: Sunday, 0: not Sunday
func CheckSunday(dateStr string) int {

	var isSunday = 0
	if len(dateStr) == 8 {
		dateStamp, err := time.Parse(DATE_FORMAT_YMD, dateStr)
		if err != nil {
			LogErr(err)
		} else {
			if dateStamp.Weekday() == time.Sunday {
				isSunday = 1
			} else {
				isSunday = 0
			}
		}
	}
	return isSunday
}

//Check Sunday len = 3 and dateStr array
func CheckSundayArray(dateStr []string) int {

	var isSunday = 0
	if len(dateStr) == 3 {
		resultDate := strings.Join(dateStr, "")
		dateStamp, err := time.Parse(DATE_FORMAT_YMD, resultDate)
		if err != nil {
			LogErr(err)
		} else {
			if dateStamp.Weekday() == time.Sunday {
				isSunday = 1
			} else {
				isSunday = 0
			}
		}
	}
	return isSunday
}

func LoopByLimitValue(n, startIndex int) []int {
	var arr = make([]int, n)
	for i := startIndex; i <= n; i++ {
		arr[i-startIndex] = i
	}
	return arr
}

func FormatDateTime(date string) string {

	if len(date) > 12 {
		return date[12:]
	}
	return date
}

func CheckLenArray(arrayJan []string) int {

	if len(arrayJan) > 0 {
		return 1
	}
	return 0
}

func ConvertArrayToString(arrayJan []string) string {

	var stringJan = ""
	if len(arrayJan) > 0 {
		stringJan = arrayJan[0]
	}
	return stringJan
}

func GetValueArray(arrayGroupType []string) string {

	if len(arrayGroupType) == 1 {
		return arrayGroupType[0]
	}
	return "0"
}

func FormatCodeAndMonth(strValue string) string {

	if len(strValue) > 6 {
		// ASO-5442 売上ベスト ※雑誌の週刊誌で月号の表示がおかしい
		//if strValue[:1] == "2" || strValue[:1] == "3" {
		//	strCode := strValue[:4]
		//	strWeek := strValue[4:5]
		//	strMonth := strValue[5:]
		//	return strCode + "+" + strMonth + "+" + strWeek
		//} else {
		//	strCode := strValue[:5]
		//	strMonth := strValue[5:]
		//	return strCode + "+" + strMonth
		//}
		strCode := strValue[:5]
		strMonth := strValue[5:]
		return strCode + "-" + strMonth
	}
	return strValue
}
