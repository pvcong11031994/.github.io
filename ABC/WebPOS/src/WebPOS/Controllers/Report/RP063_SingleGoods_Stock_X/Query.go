package RP063_SingleGoods_Stock_X

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"strconv"
	"strings"
	"time"
)

type SingleItemCumulativeResultItem struct {
	Header                   string
	ValueSales               int64
	TooltipSales             string
	ValueReceivingQuantity   int64
	TooltipReceivingQuantity string
	ValueSalesQuantityDay    int64
	TooltipSalesQuantityDay  string
}

func Query(ctx *gf.Context) {

	var err error
	ctx.ViewBases = nil
	user := WebApp.GetContextUser(ctx)
	form := QueryFormSingleGoods{}
	ctx.Form.ReadStruct(&form)

	// Check date search by GroupType
	err, messageErr := ConvertDateSearchByGroupType(&form)
	Common.LogErr(err)
	if err != nil || strings.Compare(messageErr, "") != 0 {
		ctx.ViewData["err_msg"] = messageErr
		ctx.View = "report/063_single_goods_stock_x/result_0.html"
		return
	}

	sm := Models.ShopMasterModel{ctx.DB}
	selectableShops, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)
	// システムエラー
	if err != nil {
		ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
		ctx.View = RPComon.REPORT_ERROR_PATH_HTML
		return
	}
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

	if len(selectedShopCd) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_NO_SHOP
		ctx.View = "report/063_single_goods_stock_x/result_0.html"
		return
	} else {
		form.ShopCd = selectedShopCd
	}

	// Check 日付+++++++++++++++++++++++++++++++++++++++++++++
	selectedDateFrom := form.DateFrom
	if selectedDateFrom != "" {
		if !Common.IsValidateDate(selectedDateFrom) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
			ctx.View = "report/063_single_goods_stock_x/result_0.html"
			return
		}
	}
	selectedDateTo := form.DateTo
	if selectedDateTo != "" {
		if !Common.IsValidateDate(selectedDateTo) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
			ctx.View = "report/063_single_goods_stock_x/result_0.html"
			return
		}
	}
	//=======================================================
	if strings.TrimSpace(form.DateFrom) == ""{
		form.DateFrom = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	}
	if strings.TrimSpace(form.DateTo) == ""{
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
				ctx.View = "report/063_single_goods_stock_x/result_0.html"
				return
			}
			// Limit week search = 30
		} else if form.GroupType == GROUP_TYPE_WEEK {
			timeFrom = timeFrom.AddDate(0, 0, RPComon.REPORT_LIMIT_WEEK_SEARCH*7)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LIMIT_WEEK
				ctx.View = "report/063_single_goods_stock_x/result_0.html"
				return
			}
			// Limit month search = 13
		} else if form.GroupType == GROUP_TYPE_MONTH {
			timeFrom = timeFrom.AddDate(0, RPComon.REPORT_LIMIT_MONTH_SEARCH, 0)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LIMIT_MONTH
				ctx.View = "report/063_single_goods_stock_x/result_0.html"
				return
			}
		}

	}

	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// Check JAN (単品推移画面)++++++++++++++++++++++++++++++++
	if form.FlagSingleItem != "" && len(form.JAN) < 6 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_INPUT_JAN
		ctx.View = "report/063_single_goods_stock_x/result_0.html"
		return
	}

	form.JAN = Common.GenerateJAN(form.JAN)
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// BUILD QUERY
	var sql string
	//var rows, cols, sums []string

	sql, listRange := buildSqlSingle(form, ctx)
	sqlStockAll := buildSqlStockAll(form, ctx)
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
		dataSingleItemNew, totalSingleItemNew, err := queryDataSingle(ctx, sql, listRange, form)
		dataSingleItemNew, totalSingleItemNew, err = queryDataStockAll(ctx, sqlStockAll, listRange, form,dataSingleItemNew, totalSingleItemNew)

		Common.LogErr(err)
		// システムエラー
		if err != nil {
			ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
			ctx.View = RPComon.REPORT_ERROR_PATH_HTML
			return
		} else {
			dataSingleItem = dataSingleItemNew
			totalSingleItem = totalSingleItemNew
			ctx.SaveCache(randStringFromSQLItem, dataSingleItem, 3600)
			ctx.SaveCache(randStringFromSQLTotal, totalSingleItem, 3600)
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
		if len(form.JAN) > 0 {
			tag = tag + ",JAN = " + `"` + form.JAN + `"`
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
	//-------------------------------------------------------------
	/* 20170208 Common Download File */
	// Check write file
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		/* --------------------------------------------------------
		----------------------------------------------------------*/
		file := ""
		shortFile := ""
		if len(dataSingleItem) == 0 {
			ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
			ctx.View = "report/download/result_error.html"
			return
		} else {
			err, file, shortFile = WriteFile4SingleNew(dataSingleItem, totalSingleItem, form.SearchHandleType, form.GroupType, listRange)
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
	//-------------------------------------------------------------
	ctx.ViewData["SUM_KEY_FIELD"] = RPComon.SUM_KEY_FIELD
	ctx.ViewData["data"] = dataSingleItem
	ctx.ViewData["datatotal"] = totalSingleItem
	ctx.ViewData["rangeType"] = form.GroupType
	ctx.ViewData["listCols"] = listRange

	ctx.TemplateFunc["sum_format"] = Common.FormatNumber
	ctx.TemplateFunc["minus"] = Common.Minus
	ctx.TemplateFunc["arr"] = Common.MakeArray
	ctx.TemplateFunc["checkSunday"] = Common.CheckSunday

	if len(dataSingleItem) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
		ctx.View = "report/063_single_goods_stock_x/result_0.html"
	} else {
		ctx.View = "report/063_single_goods_stock_x/single_item_transition_result.html"
		// create data for graph
		SetDataGraph(ctx, totalSingleItem, form.GroupType, listRange)
	}
}

// create data for graph
func SetDataGraph(ctx *gf.Context, dataSingleItem SingleItem, groupType string, listRange []ModelItems.MasterCalendarItem) {

	// List item show in graph
	listResultGraph := []SingleItemCumulativeResultItem{}
	// Title show in vAxis google line chart
	vAxisTitle := ""
	// property show min value vAxis is 0
	viewWindowMode := true
	// option for hAxis.title (google line chart)
	showTextEvery := 5

	var dataSalesMax int64 = 0
	var maxReceivingQuantity int64 = 0
	var maxSalesQuantity int64 = 0

	for _, col := range listRange {
		//colItem := data.Cols[col]
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
			//for keyItem, value := range colItem {
			//	toolTipItem += formatDateGraph(value) + data.HeaderCol[keyItem]
			//}
			toolTipItem += col.Mcyyyy + "年" + formatDateGraph(col.Mcmm) + "月"
			headerItem = formatDateGraph(col.Mcyyyy) + "/" + formatDateGraph(col.Mcmm)
			vAxisTitle = "月別"
		default:
			continue
		}

		// Get value 販売数
		valueSales := dataSingleItem.SaleDay[col.McKey]
		// Get value 期間入荷累計
		valueReceivingQuantity := dataSingleItem.ReceivingQuantityDay[col.McKey]
		// Get value 期間売上累計
		valueSalesQuantity := dataSingleItem.SalesQuantityDay[col.McKey]

		salesCount := valueSales
		if salesCount < 0 {
			viewWindowMode = false
		}

		// Get Max for value 期間入荷累計
		if valueReceivingQuantity > maxReceivingQuantity {
			maxReceivingQuantity = valueReceivingQuantity
		} else {
			valueReceivingQuantity = maxReceivingQuantity
		}

		// Get Max for value 期間売上累計
		if valueSalesQuantity > maxSalesQuantity {
			maxSalesQuantity = valueSalesQuantity
		} else {
			valueSalesQuantity = maxSalesQuantity
		}

		// Add list result
		listResultGraph = append(listResultGraph, SingleItemCumulativeResultItem{
			Header:                   headerItem,
			ValueSales:               valueSales,
			TooltipSales:             toolTipItem + " 販売数:" + strconv.FormatInt(valueSales, 10) + "冊",
			ValueReceivingQuantity:   valueReceivingQuantity,
			TooltipReceivingQuantity: toolTipItem + " 期間入荷累計:" + strconv.FormatInt(valueReceivingQuantity, 10) + "冊",
			ValueSalesQuantityDay:    valueSalesQuantity,
			TooltipSalesQuantityDay:  toolTipItem + " 期間売上累計:" + strconv.FormatInt(valueSalesQuantity, 10) + "冊",
		})

		// Get max for graph
		if dataSalesMax < valueSales {
			dataSalesMax = valueSales
		}
	}

	var Max int64
	if dataSalesMax > 4 {
		if (dataSalesMax % 4) != 0 {
			Max = dataSalesMax + (4 - (dataSalesMax % 4))
		} else {
			Max = dataSalesMax
		}
	} else {
		Max = 4
	}

	var MaxQuantity int64
	if maxSalesQuantity < maxReceivingQuantity {
		MaxQuantity = maxReceivingQuantity
	} else {
		MaxQuantity = maxSalesQuantity
	}
	if MaxQuantity > 4 {
		if (MaxQuantity % 4) != 0 {
			MaxQuantity = MaxQuantity + (4 - (MaxQuantity % 4))
		} else {
			MaxQuantity = MaxQuantity
		}
	} else {
		MaxQuantity = 4
	}

	ctx.ViewData["show_text_every"] = len(listResultGraph) / showTextEvery
	ctx.ViewData["list_result_graph"] = listResultGraph
	ctx.ViewData["h_title"] = vAxisTitle
	ctx.ViewData["view_window_mode"] = viewWindowMode
	ctx.ViewData["max"] = Max
	ctx.ViewData["max_quantity"] = MaxQuantity

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
func ConvertDateSearchByGroupType(form *QueryFormSingleGoods) (error, string) {

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
