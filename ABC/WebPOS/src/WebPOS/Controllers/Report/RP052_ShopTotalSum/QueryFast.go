package RP052_ShopTotalSum

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"strings"
	"time"
)

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
		ctx.View = "report/052_shop_total_sum/result_0.html"
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

	for _, item := range selectableShops {
		shopCd := item.ServerName + "|" + item.ShopCD
		selectedShopCd = append(selectedShopCd, shopCd)
	}

	if len(selectedShopCd) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_NO_SHOP
		ctx.View = "report/052_shop_total_sum/result_0.html"
		return
	} else {
		form.ShopCd = selectedShopCd
	}

	// Check 日付+++++++++++++++++++++++++++++++++++++++++++++
	selectedDateFrom := form.DateFrom
	if selectedDateFrom != "" {
		if !Common.IsValidateDate(selectedDateFrom) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
			ctx.View = "report/052_shop_total_sum/result_0.html"
			return
		}
	}
	selectedDateTo := form.DateTo
	if selectedDateTo != "" {
		if !Common.IsValidateDate(selectedDateTo) {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
			ctx.View = "report/052_shop_total_sum/result_0.html"
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
				ctx.View = "report/052_shop_total_sum/result_0.html"
				return
			}
		} else if form.GroupType == GROUP_TYPE_WEEK {
			timeFrom = timeFrom.AddDate(0, 0, 700)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_100_WEEK
				ctx.View = "report/052_shop_total_sum/result_0.html"
				return
			}
		} else if form.GroupType == GROUP_TYPE_MONTH {
			timeFrom = timeFrom.AddDate(0, 100, 0)
			if timeFrom.Before(timeTo) {
				ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_100_MONTH
				ctx.View = "report/052_shop_total_sum/result_0.html"
				return
			}
		}
	}
	// Check JAN ++++++++++++++++++++++++++++++++
	//if len(form.JAN) == 0 {
	//	ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_INPUT_JAN
	//	ctx.View = "report/051_best_sales_goods_transition/result_0.html"
	//	return
	//}
	//form.JAN = Common.GenerateJAN(form.JAN)
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// ctx.ViewData["form"] = form
	sql := buildSqlFast(form, ctx)
	// Output query to log file
	Common.LogSQL(ctx, sql)

	data := DataSum{}
	if form.Page == 1 || form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		err = queryDataFastTable(ctx, sql, &data, form)
	} else {
		err = queryDataFastTableWithJobId(ctx, sql, &data, form)
	}
	// システムエラー
	if err != nil {
		ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
		ctx.View = RPComon.REPORT_ERROR_PATH_HTML
		return
	}
	//-------------------------------------------------------------
	/* 20170208 Common Download File */
	// Check write file
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		file := ""
		shortFile := ""
		if data.CountResultRows == 0 {
			ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
			ctx.View = "report/download/result_error.html"
			return
		} else {
			err, file, shortFile = WriteFileFast(&data)
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

	ctx.ViewData["data"] = data
	ctx.TemplateFunc["sum_format"] = Common.FormatNumber
	ctx.TemplateFunc["minus"] = Common.Minus
	ctx.TemplateFunc["arr"] = Common.MakeArray

	if data.CountResultRows == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
		ctx.View = "report/052_shop_total_sum/result_0.html"
	} else {
		ctx.View = "report/052_shop_total_sum/result_fast.html"
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
