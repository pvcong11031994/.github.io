package PublisherMakerBestSaleStockDownload_Cloud

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"time"
	"errors"
	"github.com/goframework/encode"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"WebPOS/WebApp"
	"compress/gzip"
	"WebPOS/Models/ModelItems"
	"github.com/goframework/gf/db"
	"fmt"
	"archive/zip"
)

// Check condition field and execute the query
func Query(ctx *gf.Context) {

	ctx.ViewBases = nil
	form := Form{}
	ctx.Form.ReadStruct(&form)

	// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - DEL START
	//if len(form.ShopCd) == 0 {
	//	ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_NO_SHOP
	//	ctx.View = "download/publishermakerbestsalestock_cloud/result_0.html"
	//	return
	//}
	////==========================================================
	//// 日付チェック+++++++++++++++++++++++++++++++++++++++++++++
	//if strings.TrimSpace(form.DateFrom) != "" && form.DataMode != TYPE_SEARCH_STOCK {
	//	if !Common.IsValidateDate(form.DateFrom) {
	//		//form.DateFrom = time.Now().Format(Common.DATE_FORMAT_YMD)
	//		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
	//		ctx.View = "download/publishermakerbestsalestock_cloud/result_0.html"
	//		return
	//	}
	//}
	//if strings.TrimSpace(form.DateTo) != "" && form.DataMode != TYPE_SEARCH_STOCK {
	//	if !Common.IsValidateDate(form.DateTo) {
	//		//form.DateTo = time.Now().Format(Common.DATE_FORMAT_YMD)
	//		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
	//		ctx.View = "download/publishermakerbestsalestock_cloud/result_0.html"
	//		return
	//	}
	//}
	////=======================================================
	//if strings.TrimSpace(form.DateFrom) == "" {
	//	form.DateFrom = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	//}
	//if strings.TrimSpace(form.DateTo) == "" {
	//	form.DateTo = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	//}
	////※Limit日付
	//if form.DataMode != TYPE_SEARCH_STOCK {
	//	timeFrom, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, form.DateFrom)
	//	Common.LogErr(err)
	//	timeTo, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, form.DateTo)
	//	Common.LogErr(err)
	//
	//	// Limit date queryAndGetData = 100
	//	timeFrom = timeFrom.AddDate(0, 0, RPComon.REPORT_LIMIT_DATE_SEARCH)
	//	if timeFrom.Before(timeTo) {
	//		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LIMIT_DATE
	//		ctx.View = "download/publishermakerbestsalestock_cloud/result_0.html"
	//		return
	//	}
	//}
	//
	////=======================================================
	////※Check JAN
	//if len(form.JAN) == 0 {
	//	ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_INPUT_JAN
	//	ctx.View = "download/publishermakerbestsalestock_cloud/result_0.html"
	//	return
	//} else if len(strings.TrimSpace(form.JAN)) < 6 {
	//	ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_INPUT_JAN_SINGLE
	//	ctx.View = "download/publishermakerbestsalestock_cloud/result_0.html"
	//	return
	//}
	// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - DEL END

	sql := ""
	listHeader := []string{}
	if form.DataMode == TYPE_SEARCH_STOCK {
		// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5599 - EDIT START
		//sql, listHeader = buildSqlStock(form, ctx)
		listHeader = LIST_HEADER_STOCK
		// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5599 - EDIT END
	} else if form.DataMode == TYPE_SEARCH_SALES {
		sql, listHeader = buildSqlSales(form, ctx)
	} else if form.DataMode == TYPE_SEARCH_SALES_RETURN || form.DataMode == TYPE_SEARCH_SALES_RECEIVING {
		sql, listHeader = buildSqlSalesReturnsReceive(form, ctx)
	}

	// Output query to log file
	Common.LogSQL(ctx, sql)
	data := &RpData{}

	filePathCSV := ""
	fileNameCSV := ""
	err := ctx.LoadCache(sql, data)
	if form.DataMode == TYPE_SEARCH_STOCK {
		err = errors.New(TYPE_SEARCH_STOCK_TEXT)
	}
	if err != nil {
		newData := &RpData{}
		newData, err, filePathCSV, fileNameCSV = queryData(ctx, sql, form, listHeader)
		// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - DEL START
		//// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5599 - ADD START
		//if form.DataMode == TYPE_SEARCH_STOCK {
		//	if newData == nil {
		//		return
		//	}
		//	data = newData
		//} else {
			// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5599 - ADD END
		// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - EDIT END
		// システムエラー
		if err != nil {
			ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
			ctx.View = RPComon.REPORT_ERROR_PATH_HTML
			Common.LogErr(err)
			return
		} else {
			data = newData
			ctx.SaveCache(sql, data, 3600)
		}
	} else {
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
		// set report name to import info log search charging
		ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME_PUBLISHER)
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END
		// ========================================================================================
		// Output log search condition
		tag := "report=" + _REPORT_ID_PUBLISHER
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_DOWNLOAD_TEXT + `"`
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END
		tag = tag + ",店舗 IN (" + strings.Join(form.ShopCd, ",") + ")"
		tag = tag + ",期間=" + `"` + form.DateFrom + "～" + form.DateTo + `"`
		if form.JAN != "" {
			tag = tag + ",JANコード LIKE " + `"` + form.JAN + `%"`
		}
		if form.DataMode == TYPE_SEARCH_SALES {
			tag = tag + ",フォーマット=" + `"` + TYPE_SEARCH_SALES_TEXT + `"`
		} else if form.DataMode == TYPE_SEARCH_STOCK {
			tag = tag + ",フォーマット=" + `"` + TYPE_SEARCH_STOCK_TEXT + `"`
		} else if form.DataMode == TYPE_SEARCH_SALES_RETURN {
			tag = tag + ",フォーマット=" + `"` + TYPE_SEARCH_SALES_RETURN_TEXT + `"`
		} else {
			tag = tag + ",フォーマット=" + `"` + TYPE_SEARCH_SALES_RECEIVING_TEXT + `"`
		}
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
		tag = tag + `,app_id="mBAWEB-v26a"`
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

	if data.CountResultRows == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
		ctx.View = "download/publishermakerbestsalestock_cloud/result_0.html"
		return
	}

	// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - EDIT START
	//// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5599 - ADD START
	//countResultRows := data.CountResultRows
	//recordLimitFrom := ctx.Config.Int(WebApp.CONFIG_DOWNLOAD_RECORD_LIMIT_FROM, 0)
	//recordLimitTo := ctx.Config.Int(WebApp.CONFIG_DOWNLOAD_RECORD_LIMIT_TO, 0)
	//
	//if form.DataMode == TYPE_SEARCH_STOCK && recordLimitFrom < countResultRows && countResultRows < recordLimitTo {
	//	getFileNameJAN := strings.TrimSuffix(form.JAN, "%")
	//	err, filePathLoadCache, fileNameLoadCache := writeOutputFileGZip(data, form.DataMode, getFileNameJAN)
	//
	//	if err != nil || filePathLoadCache == "" {
	//		Common.LogErr(err)
	//		ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
	//		ctx.View = RPComon.REPORT_ERROR_PATH_HTML
	//		return
	//	}
	//	filePathCSV = filePathLoadCache
	//	fileNameCSV = fileNameLoadCache
	//}
	//// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5599 - ADD END
	if form.DataMode == TYPE_SEARCH_STOCK {
		if form.ZipMode == ZIP_MODE {
			getFileNameJAN := strings.TrimSuffix(form.JAN, "%")
			err, filePathLoadCache, fileNameLoadCache := writeOutputFileZip(data, form.DataMode, getFileNameJAN)

			if err != nil || filePathLoadCache == "" {
				Common.LogErr(err)
				ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
				ctx.View = RPComon.REPORT_ERROR_PATH_HTML
				return
			}
			filePathCSV = filePathLoadCache
			fileNameCSV = fileNameLoadCache
		}
	}
	// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - EDIT END

	// Process load cache
	if filePathCSV == "" || fileNameCSV == "" {
		// CSVファイルをエクスポートする
		getFileNameJAN := strings.TrimSuffix(form.JAN, "%")
		err, filePathLoadCache, fileNameLoadCache := writeOutputFile(data, form.DataMode, getFileNameJAN)
		if err != nil || filePathLoadCache == "" {
			Common.LogErr(err)
			ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
			ctx.View = RPComon.REPORT_ERROR_PATH_HTML
			return
		}
		filePathCSV = filePathLoadCache
		fileNameCSV = fileNameLoadCache
	}

	ctx.ViewData["download_file_name"] = filePathCSV
	ctx.ViewData["download_short_file_name"] = fileNameCSV
	ctx.ViewData["download"] = RPComon.PATH_REPORT_DOWN_LOAD_LINK
	ctx.View = "download/publishermakerbestsalestock_cloud/result_download.html"
}

// Write data to file csv
func writeOutputFile(data *RpData, dataMode, janCode string) (error, string, string) {

	var csvWriter *csv.Writer = nil
	var csvFile *os.File = nil
	var err error

	tmpPath, _ := filepath.Abs("./tmp")
	fileDir := filepath.Join(tmpPath, Common.CurrentDate(), Common.RandString(8))
	os.MkdirAll(fileDir, os.ModePerm)
	filePath := ""
	fileName := ""
	if dataMode == TYPE_SEARCH_SALES {
		fileName = "sales." + janCode + ".csv"
	} else if dataMode == TYPE_SEARCH_STOCK {
		fileName = "stock." + janCode + ".csv"
	} else if dataMode == TYPE_SEARCH_SALES_RETURN {
		fileName = "sales_and_return." + janCode + ".csv"
	} else {
		fileName = "sales_and_receiving." + janCode + ".csv"
	}
	filePath = filepath.Join(fileDir, fileName)

	csvFile, err = os.Create(filePath)
	if err != nil {
		return exterror.WrapExtError(err), "", ""
	}
	csvWriter = csv.NewWriter(encode.NewEncoder(encode.ENCODER_SHIFT_JIS).NewWriter(csvFile))
	csvWriter.UseCRLF = true

	csvWriter.Write(data.HeaderCols)
	for _, item := range data.Rows {
		csvWriter.Write(item)
	}

	if csvWriter != nil {
		csvWriter.Flush()
		csvFile.Close()
	}
	return nil, filePath, fileName
}

// Write data to file gzip
func writeOutputFileGZip(data *RpData, dataMode, janCode string) (error, string, string) {

	var csvWriter *csv.Writer
	var gzWriter *gzip.Writer
	var csvFile *os.File = nil
	var err error

	tmpPath, _ := filepath.Abs("./tmp")
	fileDir := filepath.Join(tmpPath, Common.CurrentDate(), Common.RandString(8))
	os.MkdirAll(fileDir, os.ModePerm)
	filePath := ""
	fileName := "stock." + janCode + ".csv.gzip"
	filePath = filepath.Join(fileDir, fileName)

	csvFile, err = os.Create(filePath)
	if err != nil {
		return exterror.WrapExtError(err), "", ""
	}
	gzWriter = gzip.NewWriter(csvFile)
	csvWriter = csv.NewWriter(encode.NewEncoder(encode.ENCODER_SHIFT_JIS).NewWriter(gzWriter))
	csvWriter.UseCRLF = true

	csvWriter.Write(data.HeaderCols)
	for _, item := range data.Rows {
		csvWriter.Write(item)
	}

	csvWriter.Flush()
	if gzWriter != nil {
		gzWriter.Flush()
		gzWriter.Close()
	}
	return nil, filePath, fileName
}

// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - EDIT START
/*
	Write data to file zip
 */
func writeOutputFileZip(data *RpData, dataMode, janCode string) (error, string, string) {

	var csvWriter *csv.Writer
	var zWriter *zip.Writer
	var csvFile *os.File = nil
	var err error

	tmpPath, _ := filepath.Abs("./tmp")
	fileDir := filepath.Join(tmpPath, Common.CurrentDate(), Common.RandString(8))
	os.MkdirAll(fileDir, os.ModePerm)
	filePath := ""
	fileNameCSV := "stock." + janCode + ".csv"
	fileNameZIP := fileNameCSV + ".zip"
	filePath = filepath.Join(fileDir, fileNameZIP)

	csvFile, err = os.Create(filePath)
	if err != nil {
		return exterror.WrapExtError(err), "", ""
	}
	zWriter = zip.NewWriter(csvFile)
	csvOut, err := zWriter.Create(fileNameCSV)
	if err != nil {
		return exterror.WrapExtError(err), "", ""
	}
	csvWriter = csv.NewWriter(encode.NewEncoder(encode.ENCODER_SHIFT_JIS).NewWriter(csvOut))
	csvWriter.UseCRLF = true

	csvWriter.Write(data.HeaderCols)
	for _, item := range data.Rows {
		csvWriter.Write(item)
	}

	csvWriter.Flush()
	if zWriter != nil {
		zWriter.Flush()
		zWriter.Close()
	}
	return nil, filePath, fileNameZIP
}
// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - EDIT END

// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - ADD START
/*
	In case of search stock, check number of records
 */
func Query_Check(ctx *gf.Context) {

	ctx.ViewBases = nil
	form := Form{}
	ctx.Form.ReadStruct(&form)

	if len(form.ShopCd) == 0 {
		ctx.JsonResponse = map[string]interface{}{
			"Errors": true,
			"Msg": RPComon.REPORT_ERROR_NO_SHOP,
		}
		return
	}
	//==========================================================
	// 日付チェック+++++++++++++++++++++++++++++++++++++++++++++
	if strings.TrimSpace(form.DateFrom) != "" && form.DataMode != TYPE_SEARCH_STOCK {
		if !Common.IsValidateDate(form.DateFrom) {
			ctx.JsonResponse = map[string]interface{}{
				"Errors": true,
				"Msg": RPComon.REPORT_ERROR_DATE_FORMAT,
			}
			return
		}
	}
	if strings.TrimSpace(form.DateTo) != "" && form.DataMode != TYPE_SEARCH_STOCK {
		if !Common.IsValidateDate(form.DateTo) {
			ctx.JsonResponse = map[string]interface{}{
				"Errors": true,
				"Msg": RPComon.REPORT_ERROR_DATE_FORMAT,
			}
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
	if form.DataMode != TYPE_SEARCH_STOCK {
		timeFrom, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, form.DateFrom)
		Common.LogErr(err)
		timeTo, err := time.Parse(Common.DATE_FORMAT_YMD_SLASH, form.DateTo)
		Common.LogErr(err)

		// Limit date queryAndGetData = 100
		timeFrom = timeFrom.AddDate(0, 0, RPComon.REPORT_LIMIT_DATE_SEARCH)
		if timeFrom.Before(timeTo) {
			ctx.JsonResponse = map[string]interface{}{
				"Errors": true,
				"Msg": RPComon.REPORT_ERROR_LIMIT_DATE,
			}
			return
		}
	}

	//=======================================================
	//※Check JAN
	if len(form.JAN) == 0 {
		ctx.JsonResponse = map[string]interface{}{
			"Errors": true,
			"Msg": RPComon.REPORT_ERROR_INPUT_JAN,
		}
		return
	} else if len(strings.TrimSpace(form.JAN)) < 6 {
		ctx.JsonResponse = map[string]interface{}{
			"Errors": true,
			"Msg": RPComon.REPORT_ERROR_INPUT_JAN_SINGLE,
		}
		return
	}

	if form.DataMode == TYPE_SEARCH_STOCK {
		listShop := form.ShopCd
		sqlString := `
SELECT
	MAX(c.shared_book_store_code) shared_book_store_code,
	a.shop_code shop_code,
	MAX(c.shop_name) shop_name,
	a.jan_code jan_code,
	MAX(b.product_name) product_name,
	MAX(b.maker_name) maker_name,
	MAX(b.list_price) list_price,
	SUM(a.cumulative_receiving_quantity) cumulative_receiving_quantity,
	SUM(a.cumulative_sales_quantity) cumulative_sales_quantity,
	SUM(a.stock_quantity) stock_quantity
FROM
	m_stock a
LEFT OUTER JOIN m_jan b
	ON a.jan_code = b.jan_code
LEFT OUTER JOIN m_shop c
	ON a.shop_code = c.shop_code
WHERE
	a.shop_code IN (?` + strings.Repeat(",?", len(listShop) - 1) + `)
	AND a.jan_code LIKE '` + form.JAN + `%'
GROUP BY
	a.jan_code,
	a.shop_code
`
		var args []interface{}
		for _, s := range listShop {
			args = append(args, s)
		}

		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - DELETE START
		//dbCloud, err := CloudSQL.Connect()
		//if err != nil {
		//	return
		//}
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - DELETE START

		// Get total records
		sqlStringTotalRecords := `
SELECT
	COUNT(*) total_records
FROM
( ` + sqlString + `
) A
		`
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT START
		//rows, err := dbCloud.Query(sqlStringTotalRecords, args...)
		//if err != nil {
		//	return
		//}
		//
		//defer func() {
		//	rows.Close()
		//	dbCloud.Close()
		//}()
		rows, err := ctx.DB.Query(sqlStringTotalRecords, args...)
		if err != nil {
			return
		}

		defer func() {
			rows.Close()
		}()
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END

		intTotalRecords := 0
		for rows.Next() {
			newMStockItem := ModelItems.MStockItem{}
			err = db.SqlScanStruct(rows, &newMStockItem)
			if err != nil {
				return
			}
			intTotalRecords = int(newMStockItem.TotalRecords)
		}
		sizeLimitFrom := ctx.Config.Int(WebApp.CONFIG_DOWNLOAD_RECORD_LIMIT_FROM, 0)
		sizeLimitTo := ctx.Config.Int(WebApp.CONFIG_DOWNLOAD_RECORD_LIMIT_TO, 0)
		if intTotalRecords >= sizeLimitTo {
			ctx.JsonResponse = map[string]interface{}{
				"Errors": true,
				"Msg": fmt.Sprintf(RPComon.REPORT_DOWNLOAD_FILE_SIZE_OVER, Common.FormatNumber(intTotalRecords)),
			}
		} else if intTotalRecords >= sizeLimitFrom {
			ctx.JsonResponse = map[string]interface{}{
				"Success": false,
				"Records": Common.FormatNumber(intTotalRecords),
			}
		} else if intTotalRecords == 0 {
			ctx.JsonResponse = map[string]interface{}{
				"Errors": true,
				"Msg": RPComon.REPORT_SEARCH_RESULT_EMPTY,
			}
		} else {
			// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - DEL START
			//Query(ctx)
			// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - DEL END
		}
	} else {
		// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - DEL START
		//Query(ctx)
		// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - DEL END
	}
}
// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - ADD END