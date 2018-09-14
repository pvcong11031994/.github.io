package ApiUtils

import (
	"bytes"
	"fmt"
	"github.com/goframework/gf"
	"github.com/goframework/gf/cfg"
	"reflect"
	"strconv"
)

// Verify api key input with key config
func VerifyApiKey(apiKey string) bool {
	//
	cfgInfo := GetConfig()
	cfgApiKey := cfgInfo.StrOrEmpty(_API_KEY)
	if cfgApiKey == apiKey {
		return true
	}
	return false
}

// Create offset,limit from page,display value input
func ParseOffset(pageCountInput, displayCountInput interface{}) (int, int) {
	//
	pageCount, displayCount := interfaceToInt(pageCountInput), interfaceToInt(displayCountInput)
	offsetCount, limitCount := 0, 0
	//
	limitCount = displayCount
	if displayCount > DB_INFO_MAX_DISPLAY_COUNT {
		limitCount = DB_INFO_MAX_DISPLAY_COUNT
	}
	if displayCount < 0 {
		limitCount = 0
	}
	//
	if pageCount > DB_INFO_MAX_PAGE_NUM {
		pageCount = DB_INFO_MAX_PAGE_NUM
	}
	if pageCount < 0 {
		pageCount = 0
	}
	offsetCount = (pageCount - 1) * limitCount
	if offsetCount < 0 {
		offsetCount = 0
	}
	return offsetCount, limitCount
}

// Convert interface type string to int
func interfaceToInt(itemConvert interface{}) int {
	st := reflect.ValueOf(itemConvert)
	switch st.Kind() {
	case reflect.String:
		num, err := strconv.Atoi(st.String())
		if err != nil {
			break
		}
		return num
	default:
		return 0
	}
	return 0
}

// Get config file
func GetConfig() *cfg.Cfg {
	var cfgInfo cfg.Cfg = cfg.Cfg{}
	cfgInfo.Load(gf.GetConfigPath())
	return &cfgInfo
}

//
func ResponseWithCallbackEmbed(embed string, data []byte) []byte {
	callbackContent := bytes.NewBufferString(fmt.Sprintf("%v(", embed))
	callbackContent.Write(data)
	callbackContent.WriteString(");")
	return callbackContent.Bytes()
}
