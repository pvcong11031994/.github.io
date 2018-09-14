package PublisherMakerBestSaleStockDownload_Maria

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"encoding/csv"
	"github.com/goframework/encode"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Check condition field and execute the query
func Query(ctx *gf.Context) {

	ctx.ViewBases = nil
	form := Form{}
	ctx.Form.ReadStruct(&form)

	if len(form.ShopCd) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_NO_SHOP
		ctx.View = "download/publishermakerbestsalestock_maria/result_0.html"
		return
	}
	//==========================================================
	// 日付チェック+++++++++++++++++++++++++++++++++++++++++++++
	if strings.TrimSpace(form.DateFrom) != "" && form.DataMode != TYPE_SEARCH_STOCK {
		if !Common.IsValidateDate(form.DateFrom) {
			//form.DateFrom = time.Now().Format(Common.DATE_FORMAT_YMD)
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
			ctx.View = "download/publishermakerbestsalestock_maria/result_0.html"
			return
		}
	}
	if strings.TrimSpace(form.DateTo) != "" && form.DataMode != TYPE_SEARCH_STOCK {
		if !Common.IsValidateDate(form.DateTo) {
			//form.DateTo = time.Now().Format(Common.DATE_FORMAT_YMD)
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_DATE_FORMAT
			ctx.View = "download/publishermakerbestsalestock_maria/result_0.html"
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
			ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LIMIT_DATE
			ctx.View = "download/publishermakerbestsalestock_maria/result_0.html"
			return
		}
	}

	//=======================================================
	//※Check JAN
	if len(form.JAN) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_INPUT_JAN
		ctx.View = "download/publishermakerbestsalestock_maria/result_0.html"
		return
	} else if len(strings.TrimSpace(form.JAN)) < 6 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_INPUT_JAN_SINGLE
		ctx.View = "download/publishermakerbestsalestock_maria/result_0.html"
		return
	}

	sql := ""
	listHeader := []string{}
	if form.DataMode == TYPE_SEARCH_STOCK {
		sql, listHeader = buildSqlStock(form, ctx)
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
	if err != nil {
		newData := &RpData{}
		newData, err, filePathCSV, fileNameCSV = queryData(ctx, sql, form, listHeader)
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
		// ========================================================================================
		// Output log search condition
		tag := "report=" + _REPORT_ID_PUBLISHER
		tag = tag + ",店舗 IN (" + strings.Join(form.ShopCd, ",") + ")"
		tag = tag + ",期間=" + `"` + form.DateFrom + "～" + form.DateTo + `"`
		if form.JAN != "" {
			tag = tag + ",JANコード LIKE " + `"` + form.JAN + `%"`
		}
		if form.DataMode == TYPE_SEARCH_SALES {
			tag = tag + ",フォーマット = " + TYPE_SEARCH_SALES_TEXT
		} else if form.DataMode == TYPE_SEARCH_STOCK {
			tag = tag + ",フォーマット = " + TYPE_SEARCH_STOCK_TEXT
		} else if form.DataMode == TYPE_SEARCH_SALES_RETURN {
			tag = tag + ",フォーマット = " + TYPE_SEARCH_SALES_RETURN_TEXT
		} else {
			tag = tag + ",フォーマット = " + TYPE_SEARCH_SALES_RECEIVING_TEXT
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

	if data.CountResultRows == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_SEARCH_RESULT_EMPTY
		ctx.View = "download/publishermakerbestsalestock_maria/result_0.html"
		return
	}

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
	ctx.View = "download/publishermakerbestsalestock_maria/result_download.html"
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
