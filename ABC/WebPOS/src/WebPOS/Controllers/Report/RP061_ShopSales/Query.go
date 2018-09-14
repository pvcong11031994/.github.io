package RP061_ShopSales

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"fmt"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"strings"
	"time"
)

func Query(ctx *gf.Context) {

	var err error
	ctx.ViewBases = nil
	form := QueryForm{}
	ctx.Form.ReadStruct(&form)
	form.Limit = 1000

	// Check 店舗+++++++++++++++++++++++++++++++++++++++++++++
	if len(form.ShopCd) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_NO_SHOP
		ctx.View = "report/061_shop_sales/result_0.html"
		return
	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// Check date search by GroupType
	err, messageErr := ConvertDateSearchByGroupType(&form)
	Common.LogErr(err)
	if err != nil {
		ctx.ViewData["err_msg"] = messageErr
		ctx.View = "report/061_shop_sales/result_0.html"
		return
	}
	// Check 日付+++++++++++++++++++++++++++++++++++++++++++++
	selectedDateFrom := strings.TrimSpace(form.DateFrom)
	if selectedDateFrom != "" {
		if !Common.IsValidateDate(selectedDateFrom) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE
			ctx.View = "report/061_shop_sales/result_0.html"
			return
		}
	}
	selectedDateTo := strings.TrimSpace(form.DateTo)
	if selectedDateTo != "" {
		if !Common.IsValidateDate(selectedDateTo) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE
			ctx.View = "report/061_shop_sales/result_0.html"
			return
		}
	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//Check Input JANコード完全一致 and JANコード前方一致
	if len(strings.TrimSpace(form.JanArrays[0])) == 0 && form.JanSingle == "" {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_INPUT_JAN_ARRAY_AND_JAN_SINGLE
		ctx.View = "report/061_shop_sales/result_0.html"
		return
	}
	// check validate JANコード前方一致 (length >= 6)
	if len(strings.TrimSpace(form.JanSingle)) < 6 && len(strings.TrimSpace(form.JanArrays[0])) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_INPUT_JAN_SINGLE
		ctx.View = "report/061_shop_sales/result_0.html"
		return
	}
	// check JANコード完全一致
	// Cut JAN by every 13 character
	listJan := strings.Split(form.JanArrays[0], "\r\n")
	form.JanArrays = []string{}
	// Cut JAN from text area
	for _, v := range listJan {
		if len(v) > 13 {
			item := ""
			for i := 0; i < len(v); i++ {
				item += v[i : i+1]
				if (i+1)%13 == 0 || (i+1) == len(v) {
					if !isExistJan(form.JanArrays, item) {
						form.JanArrays = append(form.JanArrays, item)
					}
					item = ""
				}
			}
		} else if !isExistJan(form.JanArrays, v) && v != "" {
			form.JanArrays = append(form.JanArrays, v)
		}
	}
	//=======================================================
	if strings.TrimSpace(form.DateFrom) == "" {
		form.DateFrom = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	}
	if strings.TrimSpace(form.DateTo) == "" {
		form.DateTo = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	}
	//※Limit日付
	if form.DateFrom != "" && form.DateTo != "" {
		timeFrom, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, form.DateFrom)
		Common.LogErr(err)
		timeTo, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, form.DateTo)
		Common.LogErr(err)

		// Limit date search = 100
		if form.GroupType == GROUP_TYPE_DATE {
			timeFrom = timeFrom.AddDate(0, 0, RPComon.REPORT_LIMIT_DATE_SEARCH)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LIMIT_DATE
				ctx.View = "report/061_shop_sales/result_0.html"
				return
			}
			// Limit week search = 30
		} else if form.GroupType == GROUP_TYPE_WEEK {
			timeFrom = timeFrom.AddDate(0, 0, RPComon.REPORT_LIMIT_WEEK_SEARCH*7)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LIMIT_WEEK
				ctx.View = "report/061_shop_sales/result_0.html"
				return
			}
			// Limit month search = 13
		} else if form.GroupType == GROUP_TYPE_MONTH {
			timeFrom = timeFrom.AddDate(0, RPComon.REPORT_LIMIT_MONTH_SEARCH, 0)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LIMIT_MONTH
				ctx.View = "report/061_shop_sales/result_0.html"
				return
			}
		}
	}

	sql, _, _ := buildSql(form, ctx, false)

	// Load and save cache data
	randStringFromSQL := Common.GenerateString(sql)
	// Output query to log file
	Common.LogSQL(ctx, sql)
	data := &RpData{}
	// Check case download CSV
	cacheKey := ""
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		cacheKey = randStringFromSQL + "_KEY_FILE_CSV"
	} else {
		cacheKey = randStringFromSQL + Common.ConvertIntToString(form.Page) + Common.ConvertIntToString(form.Limit) + "_KEY_PAGING"
	}

	err = ctx.LoadCache(cacheKey, data)
	if err != nil {
		newData := &RpData{}
		if form.Page == 1 {
			newData, err = queryData(ctx, sql, form, randStringFromSQL)
		} else {
			newData, err = queryGetDataWithJobId(ctx, sql, form, randStringFromSQL)
		}
		// システムエラー
		if err != nil {
			ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
			ctx.View = RPComon.REPORT_ERROR_PATH_HTML
			Common.LogErr(err)
			return
		} else {
			data = newData
			ctx.SaveCache(cacheKey, data, 3600)
		}
	} else {
		// ========================================================================================
		// Output log search condition
		tag := "report=" + _REPORT_ID
		if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
			tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_DOWNLOAD_TEXT + `"`
		} else {
			tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
		}
		if len(form.JanArrays) > 0 {
			tag = tag + ",JAN IN (" + strings.Join(form.JanArrays, ",") + ")"
		}
		if len(form.JanSingle) > 0 {
			tag = tag + ",JAN LIKE " + `"` + form.JanSingle + `%"`
		}
		if form.GroupType == GROUP_TYPE_DATE {
			tag = tag + ",単位=" + `"` + GROUP_TYPE_DATE_TEXT + `"`
			tag = tag + ",期間=" + `"` + form.DateFrom + "～" + form.DateTo + `"`
		} else if form.GroupType == GROUP_TYPE_WEEK {
			tag = tag + ",単位=" + `"` + GROUP_TYPE_WEEK_TEXT + `"`
			tag = tag + ",期間=" + `"` + form.WeekFrom + "～" + form.WeekTo + `"`
		} else if form.GroupType == GROUP_TYPE_MONTH {
			tag = tag + ",単位=" + `"` + GROUP_TYPE_MONTH_TEXT + `"`
			tag = tag + ",期間=" + `"` + form.MonthFrom + "～" + form.MonthTo + `"`
		}
		tag = tag + ",店舗 IN (" + strings.Join(form.ShopCd, ",") + ")"
		queryLog := bq.QueryLog{
			Context:   ctx,
			Tag:       tag,
			StartAt:   time.Now(),
			QuerySize: 0,
			ExecTime:  0,
			State:     bq.QUERY_LOG_BEGIN,
		}
		RPComon.QueryLogHandle(&queryLog)
		queryLog = bq.QueryLog{
			Context:   ctx,
			Tag:       tag,
			StartAt:   time.Now(),
			QuerySize: 0,
			ExecTime:  0,
			State:     bq.QUERY_LOG_END,
		}
		RPComon.QueryLogHandle(&queryLog)
		// ========================================================================================
	}
	data.HeaderCols = []string{"略称", "順位", "JANコード", "品名", "著者", "出版社名", "発売日", "本体価格", "期間入荷累計", "期間売上累計", "在庫数", "初売上日", "期間売上合計"}

	//-------------------------------------------------------------
	/* 20170208 Common Download File */
	// Check write file
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		/* --------------------------------------------------------
		----------------------------------------------------------*/
		data.HeaderCols = []string{"共有書店コード", "店舗コード", "店舗名", "JANコード", "品名", "著者", "出版社名", "発売日", "本体価格", "期間入荷累計", "期間売上累計", "在庫数", "初売上日", "期間売上合計"}
		file := ""
		shortFile := ""
		if data.CountResultRows == 0 {
			ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
			ctx.View = "report/download/result_error.html"
			return
		} else {
			err, file, shortFile = WriteFile1(data.Rows, data.HeaderCols, ctx, form)
			// システムエラー
			if err != nil {
				ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
				ctx.View = RPComon.REPORT_ERROR_PATH_HTML
				return
			}
		}
		ctx.ViewData["download_file_name"] = file
		ctx.ViewData["download_short_file_name"] = shortFile
		ctx.ViewData["report_download"] = RPComon.PATH_REPORT_DOWN_LOAD_LINK
		ctx.View = "report/download/result_download.html"
		return
	}
	//-------------------------------------------------------------
	ctx.ViewData["data"] = data
	ctx.ViewData["form"] = form
	ctx.ViewData["ms_cumulative"] = fmt.Sprintf(MS_CUMULATIVE, form.DateTo)

	ctx.TemplateFunc["sum_format"] = Common.FormatNumber
	ctx.TemplateFunc["minus"] = Common.Minus
	ctx.TemplateFunc["checkSunday"] = Common.CheckSunday
	ctx.TemplateFunc["arr"] = Common.MakeArray

	if data.CountResultRows == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
		ctx.View = "report/061_shop_sales/result_0.html"
	} else {
		ctx.View = "report/061_shop_sales/result_4.html"
	}
}

func isExistJan(listJan []string, janCheck string) bool {

	for _, v := range listJan {
		if v == janCheck {
			return true
		}
	}
	return false
}

// Check date search by GroupType
func ConvertDateSearchByGroupType(form *QueryForm) (error, string) {

	switch form.GroupType {
	case GROUP_TYPE_MONTH:
		// convert YYYY/MM to YYYY/MM/DD
		if strings.TrimSpace(form.MonthFrom) == "" {
			form.DateFrom = ""
		} else {
			monthFromTime, err := time.Parse(Common.DATE_FORMAT_YM_SLASH, form.MonthFrom)
			if err != nil {
				return err, RPComon.REPORT_ERROR_MONTH_FORMAT
			}
			form.DateFrom = monthFromTime.Format(Common.DATE_FORMAT_YMD_SLASH)
		}

		if strings.TrimSpace(form.MonthTo) == "" {
			form.DateTo = ""
		} else {
			monthToTime, err := time.Parse(Common.DATE_FORMAT_YM_SLASH, form.MonthTo)
			if err != nil {
				return err, RPComon.REPORT_ERROR_MONTH_FORMAT
			}
			form.DateTo = monthToTime.AddDate(0, 1, -1).Format(Common.DATE_FORMAT_YMD_SLASH)
		}
	case GROUP_TYPE_WEEK:
		// convert YYYY/MM/DD～MM/DD to YYYY/MM/DD
		arrWF := strings.Split(form.WeekFrom, "～")
		if strings.TrimSpace(form.WeekFrom) == "" || len(arrWF) != 2 {
			form.DateFrom = ""
		} else {
			if len(arrWF) <= 0 {
				return nil, RPComon.REPORT_ERROR_WEEK_FORMAT
			}
			weekFromTime, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, arrWF[0])
			if err != nil {
				return err, RPComon.REPORT_ERROR_WEEK_FORMAT
			}
			form.DateFrom = weekFromTime.Format(Common.DATE_FORMAT_YMD_SLASH)
		}

		arrWT := strings.Split(form.WeekTo, "～")
		if strings.TrimSpace(form.WeekTo) == "" || len(arrWT) != 2 {
			form.DateTo = ""
		} else {
			if len(arrWT) <= 0 {
				return nil, RPComon.REPORT_ERROR_WEEK_FORMAT
			}
			weekToTime, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, arrWT[0])
			if err != nil {
				return err, RPComon.REPORT_ERROR_WEEK_FORMAT
			}
			weekToTimeFull, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, weekToTime.Format("2006")+"/"+arrWT[1])
			if err != nil {
				return err, RPComon.REPORT_ERROR_WEEK_FORMAT
			}
			form.DateTo = weekToTimeFull.Format(Common.DATE_FORMAT_YMD_SLASH)
		}
	}
	return nil, ""
}
