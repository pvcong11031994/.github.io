package RP059_InitSalesCompare

import (
	"WebPOS/Common"
	favorite "WebPOS/Controllers/Report/RP060_FavoriteManagement"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"fmt"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"strings"
	"time"
)

type RpData struct {
	HeaderCols      []string
	Rows            [][]interface{}
	GraphData       map[int][]interface{}
	GraphSymbol     []string
	CountResultRows int
	MagazineName    string
	MakerName       string
	MaxValue        int64
	VJCharging      int
}

//Check validate and show data to screen
func Query(ctx *gf.Context) {

	var err error
	var userJan favorite.UserJan
	ctx.ViewBases = nil
	form := QueryForm{}
	ctx.Form.ReadStruct(&form)
	user := WebApp.GetContextUser(ctx)

	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//Check Input JAN
	if form.ControlType == CONTROL_TYPE_JAN {
		if len(strings.TrimSpace(form.JanArrays[0])) == 0 {
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_INPUT_JAN
			ctx.View = "report/059_init_sales_compare/result_0.html"
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
	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//Check Input  雑誌コード
	if form.ControlType == CONTROL_TYPE_MAGAZINE && len(form.MagazineCdSingle) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_INPUT_MAGAZINE_CODE
		ctx.View = "report/059_init_sales_compare/result_0.html"
		return
	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//Check Input  店舗
	if len(form.ShopCd) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_NO_SHOP
		ctx.View = "report/059_init_sales_compare/result_0.html"
		return
	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sql := ""
	// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - ADD START
	sqlCache := ""
	// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - ADD END
	listHeader := []string{}
	if form.ControlType == CONTROL_TYPE_JAN {
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT START
		//sql, listHeader = buildJan(form, ctx)
		sql, sqlCache, listHeader = buildJan(form, ctx)
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT END
	} else {
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT START
		//sql, listHeader = buildMagazine(form, ctx)
		sql, sqlCache, listHeader = buildMagazine(form, ctx)
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT END
	}

	// Load and save cache data
	//// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT START
	//randStringFromSQL := Common.GenerateString(sql)
	//randStringFromSQLDB := Common.GenerateString(sql + "DB")
	//// Output query to log file
	//Common.LogSQL(ctx, sql)
	randStringFromSQL := Common.GenerateString(sqlCache)
	randStringFromSQLDB := Common.GenerateString(sqlCache + "DB")
	// Output query to log file
	Common.LogSQL(ctx, sqlCache)
	// ASO-5719 [BA]bq_log_search_chargingに項目を追加して検索条件の一部を追加 - ADD START
	// set report name to import info log search charging
	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
	// ASO-5719 [BA]bq_log_search_chargingに項目を追加して検索条件の一部を追加 - ADD END
	// ========================================================================================
	// Output log search condition
	tag := "report=" + _REPORT_ID
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_DOWNLOAD_TEXT + `"`
	} else {
		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
	}
	if form.SearchDateType == SEARCH_DATE_TYPE_40 {
		tag = tag + ",期間=" + `"` + SEARCH_DATE_TYPE_40_TEXT + `"`
	} else if form.SearchDateType == SEARCH_DATE_TYPE_14 {
		tag = tag + ",期間=" + `"` + SEARCH_DATE_TYPE_14_TEXT + `"`
	}
	if form.ControlType == CONTROL_TYPE_JAN {
		if len(form.JanArrays) > 0 {
			tag = tag + ",JAN IN (" + strings.Join(form.JanArrays, ",") + ")"
		}
	}
	if form.ControlType == CONTROL_TYPE_MAGAZINE {
		if len(form.MagazineCdSingle) > 0 {
			tag = tag + ",雑誌コード=" + `"` + form.MagazineCdSingle + `"`
		}
	}
	tag = tag + ",店舗 IN (" + strings.Join(form.ShopCd, ",") + ")"
	// ASO-5719 [BA]bq_log_search_chargingに項目を追加して検索条件の一部を追加 - ADD START
	if form.ControlType == CONTROL_TYPE_JAN {
		tag = tag + `,tab="JAN/ISBN"`
	}
	if form.ControlType == CONTROL_TYPE_MAGAZINE {
		tag = tag + `,tab="雑誌コード"`
	}
	tag = tag + `,app_id="mBAWEB-v11a"`
	// ASO-5719 [BA]bq_log_search_chargingに項目を追加して検索条件の一部を追加 - ADD END
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
	// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT END
	data := &RpData{}
	dataDB := &RpData{}
	err = ctx.LoadCache(randStringFromSQL, data)
	if err != nil {
		newData := &RpData{}
		newDataDB := &RpData{}
		newData, newDataDB, err = queryData(ctx, sql, form)
		// システムエラー
		if err != nil {
			ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
			ctx.View = RPComon.REPORT_ERROR_PATH_HTML
			Common.LogErr(err)
			return
		} else {
			data = newData
			dataDB = newDataDB
			ctx.SaveCache(randStringFromSQL, data, 3600)
			// ASO-5486 初速比較 雑誌コードで検索した時、user_janへの登録がおかしくなる
			// Save cache
			ctx.SaveCache(randStringFromSQLDB, dataDB, 3600)
		}
	} else {
		// ASO-5486 初速比較 雑誌コードで検索した時、user_janへの登録がおかしくなる
		// Load cache
		err = ctx.LoadCache(randStringFromSQLDB, dataDB)
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - DEL START
		//// ========================================================================================
		//// Output log search condition
		//tag := "report=" + _REPORT_ID
		//if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		//	tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_DOWNLOAD_TEXT + `"`
		//} else {
		//	tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
		//}
		//if form.SearchDateType == SEARCH_DATE_TYPE_40 {
		//	tag = tag + ",期間=" + `"` + SEARCH_DATE_TYPE_40_TEXT + `"`
		//} else if form.SearchDateType == SEARCH_DATE_TYPE_14 {
		//	tag = tag + ",期間=" + `"` + SEARCH_DATE_TYPE_14_TEXT + `"`
		//}
		//if form.ControlType == CONTROL_TYPE_JAN {
		//	if len(form.JanArrays) > 0 {
		//		tag = tag + ",JAN IN (" + strings.Join(form.JanArrays, ",") + ")"
		//	}
		//}
		//if form.ControlType == CONTROL_TYPE_MAGAZINE {
		//	if len(form.MagazineCdSingle) > 0 {
		//		tag = tag + ",雑誌コード=" + `"` + form.MagazineCdSingle + `"`
		//	}
		//}
		//tag = tag + ",店舗 IN (" + strings.Join(form.ShopCd, ",") + ")"
		//queryLog := bq.QueryLog{
		//	Context:   ctx,
		//	Tag:       tag,
		//	StartAt:   time.Now(),
		//	QuerySize: 0,
		//	ExecTime:  0,
		//	State:     bq.QUERY_LOG_BEGIN,
		//}
		//RPComon.QueryLogHandle(&queryLog)
		//queryLog = bq.QueryLog{
		//	Context:   ctx,
		//	Tag:       tag,
		//	StartAt:   time.Now(),
		//	QuerySize: 0,
		//	ExecTime:  0,
		//	State:     bq.QUERY_LOG_END,
		//}
		//RPComon.QueryLogHandle(&queryLog)
		//// ========================================================================================
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - DEL END
	}

	// init header data
	data.HeaderCols = listHeader
	//-------------------------------------------------------------
	/* Common Download File */
	// Check write file
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		file := ""
		shortFile := ""
		if data.CountResultRows == 0 {
			ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
			ctx.View = "report/download/result_error.html"
			return
		} else {
			err, file, shortFile = WriteFileCSV(data, form.ControlType, ctx, form)
			Common.LogErr(err)
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
	// save to favorite management
	if len(dataDB.Rows) > 0 {
		userJan := Models.NewUserJanCodeModel(ctx.DB)
		err = userJan.BeginTrans()
		if err != nil {
			fmt.Println(err)
		}
		for _, dataRow := range dataDB.Rows {
			if len(dataRow[0].(string)) == 13 {
				args := []interface{}{}
				args = append(args, user.UserID,
					dataRow[0],                  //JAN
					dataRow[1],                  //商品名
					dataRow[3],                  //出版社
					dataRow[2],                  //著者名
					dataRow[9],                  //発売日
					dataRow[4],                  //本体価格
					Common.INIT_MEMO,            //メモ
					Common.INIT_PRIORITY_NUMBER) //優先順位
				err = userJan.InsertUpdateUserJan(args)
				if err != nil {
					break
				}
				Common.LogOutput(fmt.Sprintf(Common.FAVORITE_MANAGEMENT_LOG, user.UserID, dataRow[0].(string), Common.FAVORITE_MANAGEMENT_UPDATE_PROCESS))
			}
		}
		userJan.FinishTrans(&err)
		userJan.DeleteAutoUserJan(user.UserID)
	} else if form.ControlType == CONTROL_TYPE_JAN {
		for i := 0; i < len(form.JanArrays); i++ {
			userJan.JanCode = form.JanArrays[i]
			userJan.PriorityNumber = Common.INIT_PRIORITY_NUMBER
			err = favorite.UpdateJanFavorite(ctx, userJan)
		}
	} else {
		userJan.JanCode = form.MagazineCdSingle
		userJan.PriorityNumber = Common.INIT_PRIORITY_NUMBER
		err = favorite.UpdateJanFavorite(ctx, userJan)
	}
	ctx.ViewData["data"] = data
	ctx.ViewData["list_rank"] = ListRank
	ctx.ViewData["form"] = form
	ctx.ViewData["control_type"] = form.ControlType
	ctx.TemplateFunc["code_format"] = Common.FormatCodeAndMonth

	ctx.TemplateFunc["sum_format"] = Common.FormatNumber
	ctx.TemplateFunc["loop"] = Common.LoopByLimitValue

	if data.CountResultRows == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
		ctx.View = "report/059_init_sales_compare/result_0.html"
	} else {
		ctx.View = "report/059_init_sales_compare/result_1.html"
	}
}

//Check validate and show data detail by JAN to screen
func QueryDetailByJan(ctx *gf.Context) {

	ctx.ViewBases = nil
	form := QueryForm{}
	ctx.Form.ReadStruct(&form)
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT START
	//sql, listHeader := buildDetail(form, ctx)
	//// Load and save cache data
	//randStringFromSQL := Common.GenerateString(sql)
	//// Output query to log file
	//Common.LogSQL(ctx, sql)
	sql, sqlCache, listHeader := buildDetail(form, ctx)
	// Load and save cache data
	randStringFromSQL := Common.GenerateString(sqlCache)
	// Output query to log file
	Common.LogSQL(ctx, sqlCache)
	// ASO-5719 [BA]bq_log_search_chargingに項目を追加して検索条件の一部を追加 - ADD START
	// set report name to import info log search charging
	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
	// ASO-5719 [BA]bq_log_search_chargingに項目を追加して検索条件の一部を追加 - ADD END
	// ========================================================================================
	// Output log search condition
	tag := "report=" + _REPORT_ID
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_DOWNLOAD_TEXT + `"`
	} else {
		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
	}
	if form.ControlType == CONTROL_TYPE_JAN {
		if len(form.JanArrays) > 0 {
			tag = tag + ",JAN IN (" + strings.Join(form.JanArrays, ",") + ")"
		}
	}
	// ASO-5719 [BA]bq_log_search_chargingに項目を追加して検索条件の一部を追加 - ADD START
	if form.ControlType == CONTROL_TYPE_JAN {
		tag = tag + `,tab="JAN/ISBN"`
	}
	if form.ControlType == CONTROL_TYPE_MAGAZINE {
		tag = tag + `,tab="雑誌コード"`
	}
	tag = tag + `,app_id="mBAWEB-v11a"`
	// ASO-5719 [BA]bq_log_search_chargingに項目を追加して検索条件の一部を追加 - ADD END
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
	// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT END
	data := &RpData{}
	err := ctx.LoadCache(randStringFromSQL, data)
	if err != nil {
		newData := &RpData{}
		newData, err = queryDataDetail(ctx, sql, form)
		// システムエラー
		if err != nil {
			ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
			ctx.View = RPComon.REPORT_ERROR_PATH_HTML
			Common.LogErr(err)
		} else {
			data = newData
			ctx.SaveCache(randStringFromSQL, data, 3600)
		}
	} else {
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - DEL START
		//// ========================================================================================
		//// Output log search condition
		//tag := "report=" + _REPORT_ID
		//if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		//	tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_DOWNLOAD_TEXT + `"`
		//} else {
		//	tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
		//}
		//if form.ControlType == CONTROL_TYPE_JAN {
		//	if len(form.JanArrays) > 0 {
		//		tag = tag + ",JAN IN (" + strings.Join(form.JanArrays, ",") + ")"
		//	}
		//}
		//queryLog := bq.QueryLog{
		//	Context:   ctx,
		//	Tag:       tag,
		//	StartAt:   time.Now(),
		//	QuerySize: 0,
		//	ExecTime:  0,
		//	State:     bq.QUERY_LOG_BEGIN,
		//}
		//RPComon.QueryLogHandle(&queryLog)
		//queryLog = bq.QueryLog{
		//	Context:   ctx,
		//	Tag:       tag,
		//	StartAt:   time.Now(),
		//	QuerySize: 0,
		//	ExecTime:  0,
		//	State:     bq.QUERY_LOG_END,
		//}
		//RPComon.QueryLogHandle(&queryLog)
		//// ========================================================================================
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - DEL END
	}

	// init header data
	data.HeaderCols = listHeader
	ctx.ViewData["data"] = data
	ctx.ViewData["form"] = form
	ctx.ViewData["control_type"] = form.ControlType

	ctx.TemplateFunc["sum_format"] = Common.FormatNumber
	ctx.TemplateFunc["loop"] = Common.LoopByLimitValue

	ctx.View = "report/059_init_sales_compare/result_detail.html"
}

// Check exist JAN from listJan
func isExistJan(listJan []string, janCheck string) bool {

	for _, v := range listJan {
		if v == janCheck {
			return true
		}
	}
	return false
}
