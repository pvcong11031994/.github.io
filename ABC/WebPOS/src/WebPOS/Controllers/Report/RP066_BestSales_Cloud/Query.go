package RP066_BestSales_Cloud

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"strings"
	"time"

	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
)

type RpData struct {
	HeaderCols []string
	Cols       [][]string
	Rows       [][]interface{}

	CountResultRows int

	ShowLineFrom int
	ShowLineTo   int

	PageCount  int
	ThisPage   int
	VJCharging int
}

//Check validate and show data to screen
func Query(ctx *gf.Context) {

	var err error
	ctx.ViewBases = nil
	user := WebApp.GetContextUser(ctx)
	form := QueryForm{}
	ctx.Form.ReadStruct(&form)

	// Check date search by GroupType
	err, messageErr := ConvertDateSearchByGroupType(&form)
	Common.LogErr(err)
	if err != nil || strings.Compare(messageErr, "") != 0 {
		ctx.ViewData["err_msg"] = messageErr
		ctx.View = "report/066_best_sales_cloud/result_0.html"
		return
	}
	if form.ControlType == CONTROL_TYPE_BOOK {
		form.LayoutRows = "rank_no, jan_code, product_name, author_name, publisher_name, selling_date, sales_tax_exc_unit_price," +
			" cumulative_receiving_quantity, cumulative_sales_quantity, stock_quantity, first_sales_date"
	} else if form.ControlType == CONTROL_TYPE_MAGAZINE {
		form.LayoutRows = "rank_no, jan_code, product_name, magazine_code, publisher_name, selling_date, sales_tax_exc_unit_price," +
			" cumulative_receiving_quantity, cumulative_sales_quantity, stock_quantity, first_sales_date"
	}
	form.LayoutColArr = strings.Split(form.LayoutCols, ",")
	form.LayoutRowArr = strings.Split(form.LayoutRows, ",")
	form.LayoutSumArr = strings.Split(form.LayoutSums, ",")

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
				form.ShopName = append(form.ShopName, item.ShopName)
				break
			}
		}
	}

	if len(selectedShopCd) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_NO_SHOP
		ctx.View = "report/066_best_sales_cloud/result_0.html"
		return
	} else {
		form.ShopCd = selectedShopCd
	}

	// Limit selected shop when select format download :集計結果＋店舗
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT_STORE {
			if len(selectedShopCd) > RPComon.REPORT_LIMIT_SHOP_SEARCH {
				// get 100 first shops
				form.ShopCd = selectedShopCd[0:100]
			}
		}
	}

	// Check 日付+++++++++++++++++++++++++++++++++++++++++++++
	selectedDateFrom := form.DateFrom
	if selectedDateFrom != "" {
		if !Common.IsValidateDate(selectedDateFrom) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
			ctx.View = "report/066_best_sales_cloud/result_0.html"
			return
		}
	}
	selectedDateTo := form.DateTo
	if selectedDateTo != "" {
		if !Common.IsValidateDate(selectedDateTo) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
			ctx.View = "report/066_best_sales_cloud/result_0.html"
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
				ctx.View = "report/066_best_sales_cloud/result_0.html"
				return
			}
			// Limit week search = 30
		} else if form.GroupType == GROUP_TYPE_WEEK {
			timeFrom = timeFrom.AddDate(0, 0, RPComon.REPORT_LIMIT_WEEK_SEARCH*7)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LIMIT_WEEK
				ctx.View = "report/066_best_sales_cloud/result_0.html"
				return
			}
			// Limit month search = 13
		} else if form.GroupType == GROUP_TYPE_MONTH {
			timeFrom = timeFrom.AddDate(0, RPComon.REPORT_LIMIT_MONTH_SEARCH, 0)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LIMIT_MONTH
				ctx.View = "report/066_best_sales_cloud/result_0.html"
				return
			}
		}

	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//Check len 出版社コード
	if len(form.MakerCd) > 0 {
		for _, value := range form.MakerCd {
			if len(value) != 4 {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LEN_MAKER_CODE
				ctx.View = "report/066_best_sales_cloud/result_0.html"
				return
			}
		}
	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sql, cols, headerCols, exCols := buildSql(form, ctx)

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
		cacheKey = randStringFromSQL + Common.ConvertIntToString(form.Page) + Common.ConvertIntToString(form.Limit)
	}
	err = ctx.LoadCache(cacheKey, data)
	if err != nil {
		newData := &RpData{}
		if form.Page == 1 {
			newData, err = queryData(ctx, sql, form, randStringFromSQL, exCols)
		} else {
			newData, err = queryGetDataWithJobId(ctx, sql, form, randStringFromSQL, exCols)
		}
		Common.LogErr(err)
		// システムエラー
		if err != nil {
			ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
			ctx.View = RPComon.REPORT_ERROR_PATH_HTML
			return
		} else {
			data = newData
			ctx.SaveCache(cacheKey, data, 3600)
		}
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
		if len(form.MediaGroup1Cd) > 0 {
			tag = tag + ",メディア大分類コード IN (" + strings.Join(form.MediaGroup1Cd, ",") + ")"
		}
		if len(form.MediaGroup2Cd) > 0 {
			tag = tag + ",メディア中分類コード IN (" + strings.Join(form.MediaGroup2Cd, ",") + ")"
		}
		if len(form.MediaGroup3Cd) > 0 {
			tag = tag + ",メディア中小分類コード IN (" + strings.Join(form.MediaGroup3Cd, ",") + ")"
		}
		if len(form.MediaGroup4Cd) > 0 {
			tag = tag + ",メディア小分類コード IN (" + strings.Join(form.MediaGroup4Cd, ",") + ")"
		}
		if len(form.MakerCd) > 0 {
			tag = tag + ",出版社 LIKE (" + Common.JoinArray(form.MakerCd, "%", "%", ",") + ")"
		}
		if form.ControlType == CONTROL_TYPE_BOOK {
			if len(form.JanMakerCode) > 0 {
				tag = tag + ",出版者記号 IN (" + Common.JoinArray(form.JanMakerCode, JAN_MAKER_CODE, "", ",") + ")"
			}
		}
		if form.ControlType == CONTROL_TYPE_MAGAZINE {
			if len(form.MagazineCd) > 0 {
				tag = tag + ",雑誌コード LIKE (" + Common.JoinArray(form.MagazineCd, "%", "%", ",") + ")"
			}
			if form.MagazineCodeWeek == BQSL_MAGAZINE_CODE_MONTH ||
				form.MagazineCodeMonth == BQSL_MAGAZINE_CODE_WEEK ||
				form.MagazineCodeQuarter == BQSL_MAGAZINE_CODE_QUARTER {
				arrGoodsType := []string{}
				if form.MagazineCodeWeek == BQSL_MAGAZINE_CODE_MONTH {
					arrGoodsType = append(arrGoodsType, BQSL_MAGAZINE_CODE_MONTH_TEXT)
				}
				if form.MagazineCodeMonth == BQSL_MAGAZINE_CODE_WEEK {
					arrGoodsType = append(arrGoodsType, BQSL_MAGAZINE_CODE_WEEK_TEXT)
				}
				if form.MagazineCodeQuarter == BQSL_MAGAZINE_CODE_QUARTER {
					arrGoodsType = append(arrGoodsType, BQSL_MAGAZINE_CODE_QUARTER_TEXT)
				}
				tag = tag + ",商品区分 IN (" + strings.Join(arrGoodsType, ",") + ")"
			}
		}
		if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT {
			tag = tag + ",フォーマット=" + `"` + DOWNLOAD_TYPE_TOTAL_RESULT_TEXT + `"`
		} else if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION {
			tag = tag + ",フォーマット=" + `"` + DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION_TEXT + `"`
		} else if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT_STORE {
			tag = tag + ",フォーマット=" + `"` + DOWNLOAD_TYPE_TOTAL_RESULT_STORE_TEXT + `"`
		}
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
		if form.ControlType == CONTROL_TYPE_BOOK {
			tag = tag + `,tab="書籍"`
		}
		if form.ControlType == CONTROL_TYPE_MAGAZINE {
			tag = tag + `,tab="雑誌"`
		}
		tag = tag + `,app_id="mBAWEB-v22a"`
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

	data.Cols = cols
	data.HeaderCols = headerCols
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
		if data.CountResultRows == 0 {
			ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
			ctx.View = "report/download/result_error.html"
			return
		} else {
			err, file, shortFile = WriteFile4(data, form.SearchHandleType, form.ControlType, form.DownloadType)
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
	//表示上限を絞ります。
	//選択可能件数は100件/500/1000
	//-------------------------------------------------------------
	//data.PageCount = data.CountResultRows / form.Limit
	//if data.CountResultRows%form.Limit > 0 {
	//	data.PageCount += 1
	//}
	//data.ThisPage = form.Page
	//if data.ThisPage < 1 {
	//	data.ThisPage = 1
	//}
	//if data.ThisPage > data.PageCount {
	//	data.ThisPage = data.PageCount
	//}
	//
	//data.ShowLineFrom = (data.ThisPage - 1) * form.Limit
	//data.ShowLineTo = data.ThisPage*form.Limit - 1
	//
	//rowsAll := data.Rows
	//data.Rows = [][]interface{}{}
	//for k, v := range rowsAll {
	//	if k > data.ShowLineTo {
	//		break
	//	}
	//	if k >= data.ShowLineFrom && k <= data.ShowLineTo {
	//		data.Rows = append(data.Rows, v)
	//	}
	//}
	//-------------------------------------------------------------
	ctx.ViewData["data"] = data
	ctx.ViewData["rand_string"] = randStringFromSQL
	ctx.ViewData["column_number"] = len(data.Cols) + 1
	ctx.ViewData["control_type"] = form.ControlType
	ctx.ViewData["total_count"] = data.CountResultRows
	ctx.ViewData["date_to"] = form.DateTo

	ctx.TemplateFunc["sum_format"] = Common.FormatNumber
	ctx.TemplateFunc["minus"] = Common.Minus
	ctx.TemplateFunc["arr"] = Common.MakeArray
	ctx.TemplateFunc["code_format"] = Common.FormatCodeAndMonth

	if data.CountResultRows == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
		ctx.View = "report/066_best_sales_cloud/result_0.html"
	} else if form.FlagSingleItem == "" {
		ctx.View = "report/066_best_sales_cloud/result_4.html"
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
