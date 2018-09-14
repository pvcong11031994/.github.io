package RP053_BestSalesByStore

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"github.com/goframework/gf"
	"strings"
	"time"
	"github.com/goframework/gcp/bq"
)

func Query(ctx *gf.Context) {

	var err error
	ctx.ViewBases = nil
	form := QueryForm{}
	ctx.Form.ReadStruct(&form)

	// Check date search by GroupType
	err, messageErr := ConvertDateSearchByGroupType(&form)
	Common.LogErr(err)
	if err != nil || strings.Compare(messageErr, "") != 0 {
		ctx.ViewData["err_msg"] = messageErr
		ctx.View = "report/053_best_sales_by_store/result_0.html"
		return
	}

	form.LayoutColArr = strings.Split(form.LayoutCols, ",")
	form.LayoutRowArr = strings.Split(form.LayoutRows, ",")
	form.LayoutSumArr = strings.Split(form.LayoutSums, ",")

	arrShopCd := strings.Split(form.ShopCd, "|")
	strShopNm := ""
	if len(arrShopCd) > 2 {
		strShopNm = arrShopCd[2]
		form.ShopCd = arrShopCd[0] + "|" + arrShopCd[1]
	}

	Common.LogErr(err)
	// Check 店舗+++++++++++++++++++++++++++++++++++++++++++++
	if strShopNm == "" {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_NO_SHOP
		ctx.View = "report/053_best_sales_by_store/result_0.html"
		return
	}
	// Check 日付+++++++++++++++++++++++++++++++++++++++++++++
	selectedDateFrom := form.DateFrom
	if selectedDateFrom != "" {
		if !Common.IsValidateDate(selectedDateFrom) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
			ctx.View = "report/053_best_sales_by_store/result_0.html"
			return
		}
	}
	selectedDateTo := form.DateTo
	if selectedDateTo != "" {
		if !Common.IsValidateDate(selectedDateTo) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
			ctx.View = "report/053_best_sales_by_store/result_0.html"
			return
		}
	}
	//=======================================================
	//※Limit日付
	if form.DateFrom != "" && form.DateTo != "" {
		timeFrom, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, form.DateFrom)
		Common.LogErr(err)
		timeTo, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, form.DateTo)
		Common.LogErr(err)

		if form.GroupType == GROUP_TYPE_DATE {
			timeFrom = timeFrom.AddDate(0, 0, 100)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_100_DATE
				ctx.View = "report/053_best_sales_by_store/result_0.html"
				return
			}
		} else if form.GroupType == GROUP_TYPE_WEEK {
			timeFrom = timeFrom.AddDate(0, 0, 700)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_100_WEEK
				ctx.View = "report/053_best_sales_by_store/result_0.html"
				return
			}
		} else if form.GroupType == GROUP_TYPE_MONTH {
			timeFrom = timeFrom.AddDate(0, 100, 0)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_100_MONTH
				ctx.View = "report/053_best_sales_by_store/result_0.html"
				return
			}
		}

	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sql, rows, cols, sums := buildSql(form, ctx)
	// Output query to log file
	Common.LogSQL(ctx, sql)

	// Load and save cache data
	randStringFromSQL := Common.GenerateString(sql)
	data := &RPComon.ReportData{}
	err = ctx.LoadCache(randStringFromSQL, data)
	if err != nil {
		newData, err := queryData(ctx, sql, rows, cols, sums, form)
		Common.LogErr(err)
		// システムエラー
		if err != nil {
			ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
			ctx.View = RPComon.REPORT_ERROR_PATH_HTML
			return
		} else {
			data = newData
			ctx.SaveCache(randStringFromSQL, data, 21600)
		}
	} else {
		// ========================================================================================
		// Output log search condition
		tag := "report=" + _REPORT_ID
		if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
			tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV_TEXT + `"`
		} else {
			tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
		}
		if form.GroupType == GROUP_TYPE_DATE {
			tag = tag + ",表示=" + `"` + GROUP_TYPE_DATE_TEXT + `"`
			tag = tag + ",期間=" + `"` + form.DateFrom + "～" + form.DateTo + `"`
		} else if form.GroupType == GROUP_TYPE_WEEK {
			tag = tag + ",表示=" + `"` + GROUP_TYPE_WEEK_TEXT + `"`
			tag = tag + ",期間=" + `"` + form.WeekFrom + "～" + form.WeekTo + `"`
		} else if form.GroupType == GROUP_TYPE_MONTH {
			tag = tag + ",表示=" + `"` + GROUP_TYPE_MONTH_TEXT + `"`
			tag = tag + ",期間=" + `"` + form.MonthFrom + "～" + form.MonthTo + `"`
		}
		tag = tag + ",店舗=" + `"` + form.ShopCd + `"`
		if len(form.MediaGroup1Cd) > 0 {
			tag = tag + ",メディア大分類コード IN (" + strings.Join(form.MediaGroup1Cd, ",") + ")"
		}
		if len(form.MediaGroup2Cd) > 0 {
			tag = tag + ",メディア中分類コード IN (" + strings.Join(form.MediaGroup2Cd, ",") + ")"
		}
		if len(form.MediaGroup3Cd) > 0 {
			tag = tag + ",メディア中小分類コード IN (" + strings.Join(form.MediaGroup3Cd, ",") + ")"
		}
		if len(form.MakerCd) > 0 {
			tag = tag + ",出版社 LIKE (" + Common.JoinArray(form.MakerCd, "%", "%", ",") + ")"
		}
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
	//-------------------------------------------------------------
	/* 20170208 Common Download File */
	// Check write file
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		/* --------------------------------------------------------
		0. If count result of data is 0
		Output error empty data report file
		1. If length of rows and cols is 0
		Write file not ROW | COL
		2. If length of cols is 0
		Write file not COL
		3. If length of rows is 0
		Write file not ROW
		4. Write file ROW | COL | SUM
		----------------------------------------------------------*/
		file := ""
		shortFile := ""
		if len(data.ListRowKey) == 0 {
			ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
			ctx.View = "report/download/result_error.html"
			return
		} else {
			err, file, shortFile = WriteFile4(data, form.SearchHandleType, form.GroupType)
			Common.LogErr(err)
			// システムエラー
			if err != nil {
				ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
				ctx.View = RPComon.REPORT_ERROR_PATH_HTML
				return
			}
		}
		ctx.ViewData["search_handle_type"] = form.SearchHandleType
		ctx.ViewData["download_file_name"] = file
		ctx.ViewData["download_short_file_name"] = shortFile
		ctx.ViewData["report_download"] = RPComon.PATH_REPORT_DOWN_LOAD_LINK
		ctx.View = "report/download/result_download.html"
		return
	}
	data.PageCount = len(data.ListRowKey)/RPComon.LINE_PER_PAGE + 1
	if data.ThisPage < 1 {
		data.ThisPage = 1
	}
	data.ThisPage = form.Page
	if data.ThisPage > data.PageCount {
		data.ThisPage = data.PageCount
	}

	data.ShowLineFrom = (data.ThisPage - 1) * RPComon.LINE_PER_PAGE
	data.ShowLineTo = data.ThisPage*RPComon.LINE_PER_PAGE - 1

	listRow := data.ListRowKey
	data.ListRowKey = []string{}
	for key, value := range listRow {
		if key > data.ShowLineTo {
			break
		}
		if key >= data.ShowLineFrom {
			data.ListRowKey = append(data.ListRowKey, value)
		}
	}
	//-------------------------------------------------------------
	ctx.ViewData["NO_KEY_FIELD"] = RPComon.NO_KEY_FIELD
	ctx.ViewData["SUM_KEY_FIELD"] = RPComon.SUM_KEY_FIELD
	ctx.ViewData["data"] = data
	ctx.ViewData["rand_string"] = randStringFromSQL
	ctx.ViewData["column_number"] = len(data.HeaderCol) + 1
	ctx.ViewData["shop_name"] = strShopNm
	ctx.ViewData["total_count"] = len(listRow)
	ctx.ViewData["groupType"] = form.GroupType

	ctx.TemplateFunc["sum_format"] = Common.FormatNumber
	ctx.TemplateFunc["minus"] = Common.Minus
	ctx.TemplateFunc["arr"] = Common.MakeArray
	ctx.TemplateFunc["CheckSundayArray"] = Common.CheckSundayArray

	if len(data.ListRowKey) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
		ctx.View = "report/053_best_sales_by_store/result_0.html"
	} else {
		ctx.View = "report/053_best_sales_by_store/result_4.html"
	}
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
