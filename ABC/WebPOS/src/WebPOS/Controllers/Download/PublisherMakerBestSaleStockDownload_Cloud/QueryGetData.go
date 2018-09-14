package PublisherMakerBestSaleStockDownload_Cloud

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"github.com/goframework/encode"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
	"time"
)

//Create data from queryBuild
func queryData(ctx *gf.Context, sql string, form Form, listHeader []string) (*RpData, error, string, string) {

	filePath := ""
	fileName := ""
	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)

	data := RpData{
		HeaderCols: []string{},
		Cols:       [][]string{},
		Rows:       [][]string{},
	}

	keyErr := errors.New("KEY_ERR")
	msgRetryTmp := strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	retryMap := make(map[string]string)
	retryMap["CONFIG_RETRY_COUNT"] = WebApp.CONFIG_RETRY_COUNT
	retryMap["CONFIG_TIME_RETRY"] = WebApp.CONFIG_TIME_RETRY
	retryMap["CONFIG_LIST_CODE_HTTP"] = WebApp.CONFIG_LIST_CODE_HTTP
	conn, err := bq.NewConnection(keyFile, mailAccount, projectId, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return &data, exterror.WrapExtError(err), filePath, fileName
	}

	// set report name to import info log search charging
	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME_PUBLISHER)
	//totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, _REPORT_ID)
	// ========================================================================================
	// Output log search condition
	// Output log queryAndGetData condition
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
	// ========================================================================================
	dataChan := make(chan *bq.SingleRow)
	if form.DataMode != TYPE_SEARCH_STOCK {
		keyErr = errors.New("KEY_ERR")
		msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
		totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, tag, ctx, msgRetryTmp, retryMap)

		if err != nil {
			return &data, exterror.WrapExtError(err), filePath, fileName
		}
		if totalRows > RPComon.BQ_DATA_LIMIT {
			return &data, exterror.WrapExtError(errors.New("Respone data too large")), filePath, fileName
		} else if totalRows == 0 {
			return &data, nil, filePath, fileName
		}

		// Get data
		keyErr = errors.New("KEY_ERR")
		msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
		dataChan, err = conn.GetResponseData(jobId, 0, RPComon.BQ_DATA_LIMIT, ctx, msgRetryTmp, retryMap)

		if err != nil {
			return &data, exterror.WrapExtError(err), filePath, fileName
		}
	} else {
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
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
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END
	}

	// Create file CSV
	//=================================================================================================================
	var csvWriter *csv.Writer = nil
	var csvFile *os.File = nil

	janCode := strings.TrimSuffix(form.JAN, "%")
	tmpPath, _ := filepath.Abs("./tmp")
	fileDir := filepath.Join(tmpPath, Common.CurrentDate(), Common.RandString(8))
	os.MkdirAll(fileDir, os.ModePerm)
	if form.DataMode == TYPE_SEARCH_SALES {
		fileName = "sales." + janCode + ".csv"
	} else if form.DataMode == TYPE_SEARCH_STOCK {
		fileName = "stock." + janCode + ".csv"
	} else if form.DataMode == TYPE_SEARCH_SALES_RETURN {
		fileName = "sales_and_return." + janCode + ".csv"
	} else if form.DataMode == TYPE_SEARCH_SALES_RECEIVING {
		fileName = "sales_and_receiving." + janCode + ".csv"
	}
	filePath = filepath.Join(fileDir, fileName)

	csvFile, err = os.Create(filePath)
	if err != nil {
		return &data, exterror.WrapExtError(err), filePath, fileName
	}
	csvWriter = csv.NewWriter(encode.NewEncoder(encode.ENCODER_SHIFT_JIS).NewWriter(csvFile))
	csvWriter.UseCRLF = true
	csvWriter.Write(listHeader)
	//==================================================================================================================
	listShop := []string{}
	listedShopMap := map[string]bool{}
	listJan := []string{}
	listedJanMap := map[string]bool{}
	salesItem := []SalesItem{}
	// Write all line data into CSV and get RpData for save cache
	if form.DataMode == TYPE_SEARCH_SALES {
		for {
			row := <-dataChan
			if row == nil {
				break
			}

			shopCd := row.ValueMap["shop_code"].String()
			if !listedShopMap[shopCd] {
				listedShopMap[shopCd] = true
				listShop = append(listShop, shopCd)
			}
			janCd := row.ValueMap["jan_code"].String()
			if !listedJanMap[janCd] {
				listedJanMap[janCd] = true
				listJan = append(listJan, janCd)
			}

			rs := SalesItem{}
			rs.SalesDateTime = row.ValueMap["sales_datetime"].String()
			rs.SharedBookStoreCode = row.ValueMap["shared_book_store_code"].String()
			rs.ShopCode = row.ValueMap["shop_code"].String()
			rs.ShopName = row.ValueMap["shop_name"].String()
			rs.JanCode = row.ValueMap["jan_code"].String()
			rs.ProductName = row.ValueMap["product_name"].String()
			rs.MakerCode = row.ValueMap["maker_code"].String()
			rs.SalesTaxExcUnitPrice = row.ValueMap["sales_tax_exc_unit_price"].Float()
			rs.SalesBodyQuantity = row.ValueMap["sales_body_quantity"].Int()
			salesItem = append(salesItem, rs)
		}

		if len(listJan) > 0 {
			sqlString := `
SELECT
	shop_code,
	jan_code,
	stock_quantity
FROM
	m_stock
WHERE
	    shop_code IN (?` + strings.Repeat(",?", len(listShop)-1) + `)
	AND jan_code IN (?` + strings.Repeat(",?", len(listJan)-1) + `)
ORDER BY
	shop_code,
	jan_code
`
			var args []interface{}
			for _, s := range listShop {
				args = append(args, s)
			}

			for _, s := range listJan {
				args = append(args, s)
			}

			// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT START
			//dbCloud, err := CloudSQL.Connect()
			//if err != nil {
			//	return nil, exterror.WrapExtError(err), filePath, fileName
			//}
			//
			//rows, err := dbCloud.Query(sqlString, args...)
			//if err != nil {
			//	return nil, exterror.WrapExtError(err), filePath, fileName
			//}
			//
			//defer func() {
			//	rows.Close()
			//	dbCloud.Close()
			//}()
			rows, err := ctx.DB.Query(sqlString, args...)
			if err != nil {
				return nil, exterror.WrapExtError(err), filePath, fileName
			}

			defer func() {
				rows.Close()
			}()
			// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END

			mStockDataMap := map[string]ModelItems.MStockItem{}
			for rows.Next() {
				newMStockItem := ModelItems.MStockItem{}
				err = db.SqlScanStruct(rows, &newMStockItem)
				if err != nil {
					return nil, exterror.WrapExtError(err), filePath, fileName
				}
				mStockDataMap[newMStockItem.ShopCode+newMStockItem.JanCode] = newMStockItem
			}

			for i, _ := range salesItem {
				if mStockData, ok := mStockDataMap[salesItem[i].ShopCode+salesItem[i].JanCode]; ok {
					salesItem[i].StockQuantity = mStockData.StockQuantity
				}
				singleRowCSV := []string{}
				singleRowCSV = append(singleRowCSV,
					salesItem[i].SalesDateTime,
					salesItem[i].SharedBookStoreCode,
					salesItem[i].ShopCode,
					salesItem[i].ShopName,
					salesItem[i].JanCode,
					salesItem[i].ProductName,
					salesItem[i].MakerCode,
					convertFloat64ToString(salesItem[i].SalesTaxExcUnitPrice),
					convertInt64ToString(salesItem[i].SalesBodyQuantity),
					convertInt64ToString(salesItem[i].StockQuantity),
				)
				csvWriter.Write(singleRowCSV)
				data.Rows = append(data.Rows, singleRowCSV)
				data.CountResultRows++
			}
		}
	} else if form.DataMode == TYPE_SEARCH_SALES_RETURN {
		for {
			row := <-dataChan
			if row == nil {
				break
			}

			singleRowCSV := []string{}
			singleRowCSV = append(singleRowCSV,
				row.ValueMap["sales_datetime"].String(),
				row.ValueMap["shared_book_store_code"].String(),
				row.ValueMap["shop_code"].String(),
				row.ValueMap["shop_name"].String(),
				row.ValueMap["jan_code"].String(),
				row.ValueMap["product_name"].String(),
				row.ValueMap["maker_code"].String(),
				convertFloat64ToString(row.ValueMap["sales_tax_exc_unit_price"].Float()),
				convertInt64ToString(row.ValueMap["receiving_body_quantity"].Int()),
				convertInt64ToString(row.ValueMap["sales_body_quantity"].Int()),
				convertInt64ToString(row.ValueMap["return_body_quantity"].Int()),
			)
			csvWriter.Write(singleRowCSV)
			data.Rows = append(data.Rows, singleRowCSV)
			data.CountResultRows++
		}
	} else if form.DataMode == TYPE_SEARCH_STOCK {
		// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5599 - EDIT START
		//for {
		//	row := <-dataChan
		//	if row == nil {
		//		break
		//	}
		//
		//	singleRowCSV := []string{}
		//	singleRowCSV = append(singleRowCSV,
		//		row.ValueMap["shared_book_store_code"].String(),
		//		row.ValueMap["shop_code"].String(),
		//		row.ValueMap["shop_name"].String(),
		//		row.ValueMap["jan_code"].String(),
		//		row.ValueMap["product_name"].String(),
		//		row.ValueMap["maker_code"].String(),
		//		convertFloat64ToString(row.ValueMap["list_price"].Float()),
		//		convertInt64ToString(row.ValueMap["cumulative_receiving_quantity"].Int()),
		//		convertInt64ToString(row.ValueMap["cumulative_sales_quantity"].Int()),
		//		convertInt64ToString(row.ValueMap["stock_quantity"].Int()),
		//	)
		//	csvWriter.Write(singleRowCSV)
		//	data.Rows = append(data.Rows, singleRowCSV)
		//	data.CountResultRows++
		//}

		listShop = form.ShopCd
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
	SUM(a.stock_quantity) stock_quantity,
	MAX(c.shop_seq_number) shop_seq_number,
	MAX(b.magazine_seq_number) magazine_seq_number
FROM
	m_stock a
LEFT OUTER JOIN m_jan b
	ON a.jan_code = b.jan_code
LEFT OUTER JOIN m_shop c
	ON a.shop_code = c.shop_code
WHERE
	a.shop_code IN (?` + strings.Repeat(",?", len(listShop)-1) + `)
	AND a.jan_code LIKE '` + form.JAN + `%'
GROUP BY
	a.jan_code,
	a.shop_code
ORDER BY
	c.shop_seq_number, b.magazine_seq_number, a.jan_code
`
		var args []interface{}
		for _, s := range listShop {
			args = append(args, s)
		}
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - DEL START
		//dbCloud, err := CloudSQL.Connect()
		//if err != nil {
		//	return nil, exterror.WrapExtError(err), filePath, fileName
		//}
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - DEL END

		// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - DEL START
//		// Get total records
//		sqlStringTotalRecords := `
//SELECT
//	COUNT(*) total_records
//FROM
//( ` + sqlString + `
//) A
//		`
//		rows, err := dbCloud.Query(sqlStringTotalRecords, args...)
//		if err != nil {
//			return nil, exterror.WrapExtError(err), filePath, fileName
//		}
//
//		defer func() {
//			rows.Close()
//			dbCloud.Close()
//		}()
//
//		intTotalRecords := 0
//		for rows.Next() {
//			newMStockItem := ModelItems.MStockItem{}
//			err = db.SqlScanStruct(rows, &newMStockItem)
//			if err != nil {
//				return nil, exterror.WrapExtError(err), filePath, fileName
//			}
//			intTotalRecords = int(newMStockItem.TotalRecords)
//		}
//		sizeLimitTo := ctx.Config.Int(WebApp.CONFIG_DOWNLOAD_RECORD_LIMIT_TO, 0)
//		if intTotalRecords >= sizeLimitTo {
//			ctx.ViewData["err_msg"] = fmt.Sprintf(RPComon.REPORT_DOWNLOAD_FILE_SIZE_OVER, strconv.Itoa(intTotalRecords))
//			ctx.View = "download/publishermakerbestsalestock_cloud/result_0.html"
//			return nil, exterror.WrapExtError(errors.New(TYPE_SEARCH_STOCK_TEXT)), filePath, fileName
//		}
		// ASO-5599 [BA]mBAWEB-v03a 出版社ダウンロード：累計+在庫のCLOUD化 - DEL END

		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT START
		//rows, err := dbCloud.Query(sqlString, args...)
		//if err != nil {
		//	return nil, exterror.WrapExtError(err), filePath, fileName
		//}
		//
		//// ASO-5748 [BA]mBAWEB-v03a 出版社ダウンロード：1回のダウンロード実行で2回ログが吐き出される - ADD START
		//defer func() {
		//	rows.Close()
		//	dbCloud.Close()
		//}()
		//// ASO-5748 [BA]mBAWEB-v03a 出版社ダウンロード：1回のダウンロード実行で2回ログが吐き出される - ADD END
		rows, err := ctx.DB.Query(sqlString, args...)
		if err != nil {
			return nil, exterror.WrapExtError(err), filePath, fileName
		}

		defer func() {
			rows.Close()
		}()
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END

		for rows.Next() {
			newMStockItem := ModelItems.MStockItem{}
			err = db.SqlScanStruct(rows, &newMStockItem)
			if err != nil {
				return nil, exterror.WrapExtError(err), filePath, fileName
			}
			singleRowCSV := []string{}
			singleRowCSV = append(singleRowCSV,
				newMStockItem.SharedBookStoreCode,
				newMStockItem.ShopCode,
				newMStockItem.ShopName,
				newMStockItem.JanCode,
				newMStockItem.ProductName,
				newMStockItem.MakerName,
				convertInt64ToString(newMStockItem.ListPrice),
				convertInt64ToString(newMStockItem.CumulativeReceivingQuantity),
				convertInt64ToString(newMStockItem.CumulativeSalesQuantity),
				convertInt64ToString(newMStockItem.StockQuantity),
			)
			csvWriter.Write(singleRowCSV)
			data.Rows = append(data.Rows, singleRowCSV)
			data.CountResultRows++
		}
		// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5599 - EDIT END
	} else if form.DataMode == TYPE_SEARCH_SALES_RECEIVING {
		for {
			row := <-dataChan
			if row == nil {
				break
			}

			singleRowCSV := []string{}
			singleRowCSV = append(singleRowCSV,
				row.ValueMap["sales_datetime"].String(),
				row.ValueMap["shared_book_store_code"].String(),
				row.ValueMap["shop_code"].String(),
				row.ValueMap["shop_name"].String(),
				row.ValueMap["jan_code"].String(),
				row.ValueMap["product_name"].String(),
				row.ValueMap["maker_code"].String(),
				convertFloat64ToString(row.ValueMap["sales_tax_exc_unit_price"].Float()),
				convertInt64ToString(row.ValueMap["receiving_body_quantity"].Int()),
				convertInt64ToString(row.ValueMap["sales_body_quantity"].Int()),
			)
			csvWriter.Write(singleRowCSV)
			data.Rows = append(data.Rows, singleRowCSV)
			data.CountResultRows++
		}
	}

	data.HeaderCols = listHeader

	if csvWriter != nil {
		csvWriter.Flush()
		csvFile.Close()
	}

	return &data, nil, filePath, fileName
}

func convertFloat64ToString(value float64) string {
	return strconv.Itoa(int(value))
}

func convertInt64ToString(value int64) string {
	return strconv.Itoa(int(value))
}
