package PublisherMakerBestSaleStockDownload_Maria

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/WebApp"
	"encoding/csv"
	"errors"
	"github.com/goframework/encode"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"WebPOS/Models/ModelItems"
	"github.com/goframework/gf/db"
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
	// ========================================================================================
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
	dataChan, err := conn.GetResponseData(jobId, 0, RPComon.BQ_DATA_LIMIT, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return &data, exterror.WrapExtError(err), filePath, fileName
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
	// ASO-5539 [BA]mBAWEB-v24a ：出版社ダウンロード(Maria)　bestsalestockdownload_maria
	listJan := make(map[string]string)
	salesItem := []SalesItem{}
	// Write all line data into CSV and get RpData for save cache
	if form.DataMode == TYPE_SEARCH_SALES {
		for {
			row := <-dataChan
			if row == nil {
				break
			}

			// ASO-5539 [BA]mBAWEB-v24a ：出版社ダウンロード(Maria)　bestsalestockdownload_maria
			//singleRowCSV := []string{}
			//singleRowCSV = append(singleRowCSV,
			//	row.ValueMap["sales_datetime"].String(),
			//	row.ValueMap["shared_book_store_code"].String(),
			//	row.ValueMap["shop_code"].String(),
			//	row.ValueMap["shop_name"].String(),
			//	row.ValueMap["jan_code"].String(),
			//	row.ValueMap["product_name"].String(),
			//	row.ValueMap["maker_code"].String(),
			//	convertFloat64ToString(row.ValueMap["sales_tax_exc_unit_price"].Float()),
			//	convertInt64ToString(row.ValueMap["sales_body_quantity"].Int()),
			//	convertInt64ToString(row.ValueMap["stock_quantity"].Int()),
			//)
			//csvWriter.Write(singleRowCSV)
			//data.Rows = append(data.Rows, singleRowCSV)
			//data.CountResultRows++
			listJan[row.ValueMap["jan_code"].String()] = row.ValueMap["jan_code"].String()
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
		// ASO-5539 [BA]mBAWEB-v24a ：出版社ダウンロード(Maria)　bestsalestockdownload_maria
		if len(listJan) > 0 {
			// Get 在庫数, 累計入荷数, 累計売上数, 初回売上日付 from m_stock of MariaDB
			listMStock := make(map[string]int64)
			sqlString := `
SELECT
	shop_code,
	jan_code,
	SUM(stock_quantity) stock_quantity
FROM
	m_stock
WHERE
	shop_code IN (?` + strings.Repeat(",?", len(form.ShopCd) - 1) + `)
	AND jan_code IN (?` + strings.Repeat(",?", len(listJan) - 1) + `)
	AND TRIM(first_sales_date) <> ''
GROUP BY
	shop_code,
	jan_code
`
			var args []interface{}
			for _, s := range form.ShopCd {
				args = append(args, s)
			}
			for _, s := range listJan {
				args = append(args, s)
			}
			rows, err := ctx.DB.Query(sqlString, args...)
			if err != nil {
				return nil, exterror.WrapExtError(err), filePath, fileName
			}
			defer rows.Close()
			for rows.Next() {
				newMStockItem := ModelItems.MStockItem{}
				err = db.SqlScanStruct(rows, &newMStockItem)
				if err != nil {
					return nil, exterror.WrapExtError(err), filePath, fileName
				}
				listMStock[newMStockItem.ShopCode + newMStockItem.JanCode] = newMStockItem.StockQuantity
			}
			for i, _ := range salesItem {
				if _, ok := listMStock[salesItem[i].ShopCode + salesItem[i].JanCode]; ok {
					salesItem[i].StockQuantity = listMStock[salesItem[i].ShopCode + salesItem[i].JanCode]
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
		for {
			row := <-dataChan
			if row == nil {
				break
			}

			singleRowCSV := []string{}
			singleRowCSV = append(singleRowCSV,
				row.ValueMap["shared_book_store_code"].String(),
				row.ValueMap["shop_code"].String(),
				row.ValueMap["shop_name"].String(),
				row.ValueMap["jan_code"].String(),
				row.ValueMap["product_name"].String(),
				row.ValueMap["maker_code"].String(),
				convertFloat64ToString(row.ValueMap["list_price"].Float()),
				convertInt64ToString(row.ValueMap["cumulative_receiving_quantity"].Int()),
				convertInt64ToString(row.ValueMap["cumulative_sales_quantity"].Int()),
				convertInt64ToString(row.ValueMap["stock_quantity"].Int()),
			)
			csvWriter.Write(singleRowCSV)
			data.Rows = append(data.Rows, singleRowCSV)
			data.CountResultRows++
		}
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