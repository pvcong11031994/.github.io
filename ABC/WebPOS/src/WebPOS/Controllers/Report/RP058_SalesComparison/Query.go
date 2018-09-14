package RP058_SalesComparison

import (
	"WebPOS/Common"
	favorite "WebPOS/Controllers/Report/RP060_FavoriteManagement"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"fmt"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"strconv"
	"strings"
	"time"
)

type ResultItem struct {
	Header string
	Value0, Value1, Value2, Value3, Value4, Value5, Value6, Value7,
	Value8, Value9, Value10, Value11, Value12, Value13, Value14 int64
	Tooltip0, Tooltip1, Tooltip2, Tooltip3, Tooltip4, Tooltip5, Tooltip6, Tooltip7,
	Tooltip8, Tooltip9, Tooltip10, Tooltip11, Tooltip12, Tooltip13, Tooltip14 string
}

func Query(ctx *gf.Context) {

	var err error
	var userJan favorite.UserJan
	ctx.ViewBases = nil
	user := WebApp.GetContextUser(ctx)
	form := QueryForm{}
	ctx.Form.ReadStruct(&form)

	// Check date search by GroupType
	err, messageErr := ConvertDateSearchByGroupType(&form)
	Common.LogErr(err)
	if err != nil || strings.Compare(messageErr, "") != 0 {
		ctx.ViewData["err_msg"] = messageErr
		ctx.View = "report/058_sales_comparison/result_0.html"
		return
	}

	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//Check Input JAN

	if len(strings.TrimSpace(form.JanArrays[0])) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_INPUT_JAN
		ctx.View = "report/058_sales_comparison/result_0.html"
		return
	}
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
		} else if !isExistJan(form.JanArrays, v) {
			form.JanArrays = append(form.JanArrays, v)
		}
	}

	//Cut JanArrays
	if len(form.JanArrays) > 15 {
		form.JanArrays = form.JanArrays[:15]

	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// Check 日付+++++++++++++++++++++++++++++++++++++++++++++
	selectedDateFrom := form.DateFrom
	if selectedDateFrom != "" {
		if !Common.IsValidateDate(selectedDateFrom) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE
			ctx.View = "report/058_sales_comparison/result_0.html"
			return
		}
	}
	selectedDateTo := form.DateTo
	if selectedDateTo != "" {
		if !Common.IsValidateDate(selectedDateTo) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE
			ctx.View = "report/058_sales_comparison/result_0.html"
			return
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
				ctx.View = "report/058_sales_comparison/result_0.html"
				return
			}
			// Limit week search = 30
		} else if form.GroupType == GROUP_TYPE_WEEK {
			timeFrom = timeFrom.AddDate(0, 0, RPComon.REPORT_LIMIT_WEEK_SEARCH*7)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LIMIT_WEEK
				ctx.View = "report/058_sales_comparison/result_0.html"
				return
			}
			// Limit month search = 13
		} else if form.GroupType == GROUP_TYPE_MONTH {
			timeFrom = timeFrom.AddDate(0, RPComon.REPORT_LIMIT_MONTH_SEARCH, 0)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LIMIT_MONTH
				ctx.View = "report/058_sales_comparison/result_0.html"
				return
			}
		}
	}

	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sm := Models.ShopMasterModel{ctx.DB}
	selectableShops, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)
	// Check 店舗+++++++++++++++++++++++++++++++++++++++++++++
	selectedShopCd := []string{}
	for _, shopCd := range form.ShopCd {
		for _, item := range selectableShops {
			if shopCd == item.ShopCD {
				selectedShopCd = append(selectedShopCd, shopCd)
				break
			}
		}
	}

	// Check Shop ++++++++++++++++++++++++++++++++
	if len(selectedShopCd) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_NO_SHOP
		ctx.View = "report/058_sales_comparison/result_0.html"
		return
	} else {
		form.ShopCd = selectedShopCd
	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	sql, listRange := buildSql(form, ctx)
	// Output query to log file
	Common.LogSQL(ctx, sql)

	// Load and save cache data
	randStringFromSQLItem := Common.GenerateString(sql + "Item")
	randStringFromSQLTotal := Common.GenerateString(sql + "Total")
	dataSingleItem := []SingleItem{}
	totalSingleItem := SingleItem{}

	errItem := ctx.LoadCache(randStringFromSQLItem, &dataSingleItem)
	errTotal := ctx.LoadCache(randStringFromSQLTotal, &totalSingleItem)
	if errItem != nil || errTotal != nil {
		dataSingleItemNew, totalSingleItemNew, err := queryData(ctx, sql, listRange, form)
		Common.LogErr(err)
		// システムエラー
		if err != nil {
			ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
			ctx.View = RPComon.REPORT_ERROR_PATH_HTML
			return
		}
		dataSingleItem = dataSingleItemNew
		totalSingleItem = totalSingleItemNew
		ctx.SaveCache(randStringFromSQLItem, dataSingleItem, 21600)
		ctx.SaveCache(randStringFromSQLTotal, totalSingleItem, 21600)
	} else {
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
		// set report name to import info log search charging
		ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END
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
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
		tag = tag + `,app_id="mBAWEB-v10a"`
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END
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

	dataSingleItemSort := []SingleItem{}
	for _, valueFormJan := range form.JanArrays {
		for _, valueSingleJan := range dataSingleItem {
			if valueFormJan == valueSingleJan.JanCd {
				dataSingleItemSort = append(dataSingleItemSort, valueSingleJan)
			}
		}
	}
	//-------------------------------------------------------------
	/* 20170208 Common Download File */
	// Check write file
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		/* --------------------------------------------------------
		----------------------------------------------------------*/
		file := ""
		shortFile := ""
		if len(dataSingleItemSort) == 0 {
			ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
			ctx.View = "report/download/result_error.html"
			return
		} else {
			err, file, shortFile = WriteFile4(dataSingleItemSort, form.SearchHandleType, form.GroupType, ctx, form)
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
	// save to favorite management
	if len(dataSingleItem) > 0 {
		userJan := Models.NewUserJanCodeModel(ctx.DB)
		err = userJan.BeginTrans()
		if err != nil {
			exterror.WrapExtError(err)
		}
		for i := 0; i < len(dataSingleItem); i++ {
			if len(dataSingleItem[i].JanCd) == 13 {
				args := []interface{}{}
				args = append(args, user.UserID, dataSingleItem[i].JanCd, dataSingleItem[i].GoodsName, dataSingleItem[i].PublisherName,
					dataSingleItem[i].AuthorName, dataSingleItem[i].SaleDate, dataSingleItem[i].Price, Common.INIT_MEMO, Common.INIT_PRIORITY_NUMBER)
				err = userJan.InsertUpdateUserJan(args)
				if err != nil {
					break
				}
				Common.LogOutput(fmt.Sprintf(Common.FAVORITE_MANAGEMENT_LOG, user.UserID, dataSingleItem[0].JanCd, Common.FAVORITE_MANAGEMENT_UPDATE_PROCESS))
			}
		}
		userJan.FinishTrans(&err)
		userJan.DeleteAutoUserJan(user.UserID)
	} else {
		for i := 0; i < len(form.JanArrays); i++ {
			userJan.JanCode = form.JanArrays[i]
			userJan.PriorityNumber = Common.INIT_PRIORITY_NUMBER
			err = favorite.UpdateJanFavorite(ctx, userJan)
		}
	}

	//-------------------------------------------------------------

	ctx.ViewData["data"] = dataSingleItemSort
	ctx.ViewData["total_count"] = len(dataSingleItemSort)
	ctx.ViewData["datatotal"] = totalSingleItem
	ctx.ViewData["rangeType"] = form.GroupType
	ctx.ViewData["listCol"] = listRange
	ctx.ViewData["form"] = form
	ctx.ViewData["ms_cumulative"] = fmt.Sprintf(MS_CUMULATIVE, form.DateTo)

	ctx.TemplateFunc["sum_format"] = Common.FormatNumber
	ctx.TemplateFunc["minus"] = Common.Minus
	ctx.TemplateFunc["checkSunday"] = Common.CheckSunday
	ctx.TemplateFunc["arr"] = Common.MakeArray

	if len(dataSingleItemSort) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
		ctx.View = "report/058_sales_comparison/result_0.html"
	} else {
		ctx.View = "report/058_sales_comparison/result_4.html"
		// create data for graph
		SetDataGraph(ctx, totalSingleItem, form.GroupType, listRange, dataSingleItemSort)
	}
}

// create data for graph
func SetDataGraph(ctx *gf.Context, totalSingleItem SingleItem, groupType string, listRange []ModelItems.MasterCalendarItem, dataSingleItem []SingleItem) {

	// List item show in graph
	listResultGraph := []ResultItem{}
	// Title show in vAxis google line chart
	vAxisTitle := ""
	// property show min value vAxis is 0
	viewWindowMode := true
	// option for hAxis.title (google line chart)
	showTextEvery := 5

	var dataMax int64 = 0
	for _, col := range listRange {
		headerItem := ""
		toolTipItem := ""
		switch groupType {
		case GROUP_TYPE_DATE:
			// display value YYYY/MM/DD
			toolTipItem += col.Mcyyyy + "年" + formatDateGraph(col.Mcmm) + "月" + formatDateGraph(col.Mcdd) + "日"
			headerItem = formatDateGraph(col.Mcmm) + "/" + formatDateGraph(col.Mcdd)
			vAxisTitle = "日別"
			showTextEvery = 10
		case GROUP_TYPE_WEEK:
			// display value YYYY/MM/DD～MM/DD
			toolTipItem = col.Mcyyyy + "/" + formatWeekDateGraph(col.Mcweekdate)
			headerItem = formatWeekDateGraph(col.Mcweekdate)
			vAxisTitle = "週別"
		case GROUP_TYPE_MONTH:
			// display value YYYY/MM
			toolTipItem += col.Mcyyyy + "年" + formatDateGraph(col.Mcmm) + "月"
			headerItem = formatDateGraph(col.Mcyyyy) + "/" + formatDateGraph(col.Mcmm)
			vAxisTitle = "月別"
		default:
			continue
		}
		getValueChart(totalSingleItem, dataSingleItem, col, &listResultGraph, headerItem, toolTipItem, &dataMax)
	}
	var Max int64
	if dataMax > 4 {
		if (dataMax % 4) != 0 {
			Max = (dataMax + (4 - (dataMax % 4))) + ((dataMax+(4-(dataMax%4)))*5)/100
		} else {
			Max = dataMax + dataMax*5/100
		}
	} else {
		Max = 4
	}

	ctx.ViewData["show_text_every"] = len(listResultGraph) / showTextEvery
	ctx.ViewData["len_data_single_item"] = len(dataSingleItem)
	ctx.ViewData["list_result_graph"] = listResultGraph
	ctx.ViewData["h_title"] = vAxisTitle
	ctx.ViewData["view_window_mode"] = viewWindowMode
	ctx.ViewData["max"] = Max
}

func isExistJan(listJan []string, janCheck string) bool {

	for _, v := range listJan {
		if v == janCheck {
			return true
		}
	}
	return false
}

func getValueMax(valueMax, valueItem int64) int64 {

	if valueMax < valueItem {
		valueMax = valueItem
	}
	return valueMax
}
func getValueChart(totalSingleItem SingleItem, dataSingleItem []SingleItem, col ModelItems.MasterCalendarItem,
	listResultGraph *[]ResultItem, headerItem, toolTipItem string, dataMax *int64) {

	// Get value
	valueItem0 := [15]int64{0}

	for k, valueSingle := range dataSingleItem {
		valueItem0[k] = totalSingleItem.SaleDay[valueSingle.JanCd][col.McKey]
		*dataMax = getValueMax(*dataMax, valueItem0[k])

	}
	*listResultGraph = append(*listResultGraph, ResultItem{
		Header:    headerItem,
		Value0:    valueItem0[0],
		Tooltip0:  toolTipItem + " A: " + strconv.FormatInt(valueItem0[0], 10) + "冊",
		Value1:    valueItem0[1],
		Tooltip1:  toolTipItem + " B: " + strconv.FormatInt(valueItem0[1], 10) + "冊",
		Value2:    valueItem0[2],
		Tooltip2:  toolTipItem + " C: " + strconv.FormatInt(valueItem0[2], 10) + "冊",
		Value3:    valueItem0[3],
		Tooltip3:  toolTipItem + " D: " + strconv.FormatInt(valueItem0[3], 10) + "冊",
		Value4:    valueItem0[4],
		Tooltip4:  toolTipItem + " E: " + strconv.FormatInt(valueItem0[4], 10) + "冊",
		Value5:    valueItem0[5],
		Tooltip5:  toolTipItem + " F: " + strconv.FormatInt(valueItem0[5], 10) + "冊",
		Value6:    valueItem0[6],
		Tooltip6:  toolTipItem + " G: " + strconv.FormatInt(valueItem0[6], 10) + "冊",
		Value7:    valueItem0[7],
		Tooltip7:  toolTipItem + " H: " + strconv.FormatInt(valueItem0[7], 10) + "冊",
		Value8:    valueItem0[8],
		Tooltip8:  toolTipItem + " I: " + strconv.FormatInt(valueItem0[8], 10) + "冊",
		Value9:    valueItem0[9],
		Tooltip9:  toolTipItem + " J: " + strconv.FormatInt(valueItem0[9], 10) + "冊",
		Value10:   valueItem0[10],
		Tooltip10: toolTipItem + " K: " + strconv.FormatInt(valueItem0[10], 10) + "冊",
		Value11:   valueItem0[11],
		Tooltip11: toolTipItem + " L: " + strconv.FormatInt(valueItem0[11], 10) + "冊",
		Value12:   valueItem0[12],
		Tooltip12: toolTipItem + " M: " + strconv.FormatInt(valueItem0[12], 10) + "冊",
		Value13:   valueItem0[13],
		Tooltip13: toolTipItem + " N: " + strconv.FormatInt(valueItem0[13], 10) + "冊",
		Value14:   valueItem0[14],
		Tooltip14: toolTipItem + " O: " + strconv.FormatInt(valueItem0[14], 10) + "冊",
	})
}

// Convert display string from "04" to "4"
func formatDateGraph(number string) string {

	if len(number) > 1 && number[0:1] == "0" {
		return number[1:]
	}
	return number
}

// Convert display string from "02/27～03/05" to "2/27～3/5"
func formatWeekDateGraph(str string) string {

	list := strings.Split(str, "～")
	if len(list) <= 1 {
		return str
	}
	//newStr := formatDateGraph(strings.Split(list[0], "/")[0]) + "/" + formatDateGraph(strings.Split(list[0], "/")[1])
	//newStr += "～"
	//newStr += formatDateGraph(strings.Split(list[1], "/")[0]) + "/" + formatDateGraph(strings.Split(list[1], "/")[1])
	arrFrom := strings.Split(list[0], "/")
	if len(arrFrom) <= 1 {
		return str
	}
	newStr := formatDateGraph(arrFrom[0]) + "/" + formatDateGraph(arrFrom[1])
	newStr += "～"
	arrTo := strings.Split(list[1], "/")
	if len(arrTo) <= 1 {
		return str
	}
	newStr += formatDateGraph(arrTo[0]) + "/" + formatDateGraph(arrTo[1])
	return newStr
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
				return err, RPComon.REPORT_ERROR_MONTH
			}
			form.DateFrom = monthFromTime.Format(Common.DATE_FORMAT_YMD_SLASH)
		}

		if strings.TrimSpace(form.MonthTo) == "" {
			form.DateTo = ""
		} else {
			monthToTime, err := time.Parse(Common.DATE_FORMAT_YM_SLASH, form.MonthTo)
			if err != nil {
				return err, RPComon.REPORT_ERROR_MONTH
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
				return err, RPComon.REPORT_ERROR_DATE
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
				return err, RPComon.REPORT_ERROR_DATE
			}
			weekToTimeFull, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, weekToTime.Format("2006")+"/"+arrWT[1])
			if err != nil {
				return err, RPComon.REPORT_ERROR_DATE
			}
			form.DateTo = weekToTimeFull.Format(Common.DATE_FORMAT_YMD_SLASH)
		}
	}
	return nil, ""
}

func QueryDetail(ctx *gf.Context) {

	ctx.ViewBases = nil
	form := QueryForm{}
	ctx.Form.ReadStruct(&form)

	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sql, listRange := buildDetailSql(form, ctx, form.JanKey)

	// Load and save cache data
	randStringFromSQLItem := Common.GenerateString(sql + "Item")
	randStringFromSQLTotal := Common.GenerateString(sql + "Total")
	// Output query to log file
	Common.LogSQL(ctx, sql)
	dataSingleItem := []SingleItem{}
	totalSingleItem := SingleItem{}
	errItem := ctx.LoadCache(randStringFromSQLItem, &dataSingleItem)
	errTotal := ctx.LoadCache(randStringFromSQLTotal, &totalSingleItem)
	if errItem != nil || errTotal != nil {
		dataSingleItemNew, totalSingleItemNew, err := queryDataDetail(ctx, sql, listRange, form, form.JanKey)
		Common.LogErr(err)
		// システムエラー
		if err != nil {
			ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
			ctx.View = RPComon.REPORT_ERROR_PATH_HTML
			return
		}
		dataSingleItem = dataSingleItemNew
		totalSingleItem = totalSingleItemNew
		ctx.SaveCache(randStringFromSQLItem, dataSingleItem, 21600)
		ctx.SaveCache(randStringFromSQLTotal, totalSingleItem, 21600)
	} else {
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
		// set report name to import info log search charging
		ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END
		// ========================================================================================
		// Output log search condition
		tag := "report=" + _REPORT_ID
		if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
			tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_DOWNLOAD_TEXT + `"`
		} else {
			tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
		}
		tag = tag + ",JAN = " + `"` + form.JanKey + `"`
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
		tag = tag + `,app_id="mBAWEB-v10a"`
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END
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
	}

	//-------------------------------------------------------------
	ctx.ViewData["data"] = dataSingleItem
	ctx.ViewData["datatotal"] = totalSingleItem
	ctx.ViewData["rangeType"] = form.GroupType
	ctx.ViewData["listCol"] = listRange

	ctx.TemplateFunc["sum_format"] = Common.FormatNumber
	ctx.TemplateFunc["minus"] = Common.Minus
	ctx.TemplateFunc["checkSunday"] = Common.CheckSunday
	ctx.TemplateFunc["arr"] = Common.MakeArray
	if len(dataSingleItem) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
		ctx.View = "report/058_sales_comparison/result_0.html"
	} else {
		ctx.View = "report/058_sales_comparison/result_detail.html"
	}
}
