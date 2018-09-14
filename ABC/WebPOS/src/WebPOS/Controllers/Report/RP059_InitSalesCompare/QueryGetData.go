package RP059_InitSalesCompare

import (
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/ModelItems"
	"encoding/csv"
	"github.com/goframework/gf"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
	"reflect"
	"strconv"
	"strings"
	"github.com/goframework/gcp/bq"
	"time"
)

////Create data from queryBuild
//func queryData(ctx *gf.Context, sql string, form QueryForm) (*RpData, *RpData, error) {
//
//	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
//	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
//	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)
//
//	data := RpData{
//		HeaderCols: []string{},
//		Rows:       [][]interface{}{},
//		GraphData:  map[int][]interface{}{},
//		MaxValue:   0,
//	}
//	dataDB := RpData{
//		HeaderCols: []string{},
//		Rows:       [][]interface{}{},
//		GraphData:  map[int][]interface{}{},
//		MaxValue:   0,
//	}
//
//	keyErr := errors.New("KEY_ERR")
//	msgRetryTmp := strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
//	// ASO-5502 JOBのコンパイルが通らない（stg）
//	retryMap := make(map[string]string)
//	retryMap["CONFIG_RETRY_COUNT"] = WebApp.CONFIG_RETRY_COUNT
//	retryMap["CONFIG_TIME_RETRY"] = WebApp.CONFIG_TIME_RETRY
//	retryMap["CONFIG_LIST_CODE_HTTP"] = WebApp.CONFIG_LIST_CODE_HTTP
//	conn, err := bq.NewConnection(keyFile, mailAccount, projectId, ctx, msgRetryTmp, retryMap)
//
//	if err != nil {
//		return &data, &dataDB, exterror.WrapExtError(err)
//	}
//
//	// set report name to import info log search charging
//	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
//	//totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, _REPORT_ID)
//	// ========================================================================================
//	// Output log search condition
//	tag := "report=" + _REPORT_ID
//	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
//		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_DOWNLOAD_TEXT + `"`
//	} else {
//		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
//	}
//	if form.SearchDateType == SEARCH_DATE_TYPE_40 {
//		tag = tag + ",期間=" + `"` + SEARCH_DATE_TYPE_40_TEXT + `"`
//	} else if form.SearchDateType == SEARCH_DATE_TYPE_14 {
//		tag = tag + ",期間=" + `"` + SEARCH_DATE_TYPE_14_TEXT + `"`
//	}
//	if form.ControlType == CONTROL_TYPE_JAN {
//		if len(form.JanArrays) > 0 {
//			tag = tag + ",JAN IN (" + strings.Join(form.JanArrays, ",") + ")"
//		}
//	}
//	if form.ControlType == CONTROL_TYPE_MAGAZINE {
//		if len(form.MagazineCdSingle) > 0 {
//			tag = tag + ",雑誌コード=" + `"` + form.MagazineCdSingle + `"`
//		}
//	}
//	tag = tag + ",店舗 IN (" + strings.Join(form.ShopCd, ",") + ")"
//	// ========================================================================================
//	keyErr = errors.New("KEY_ERR")
//	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
//	totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, tag, ctx, msgRetryTmp, retryMap)
//
//	if err != nil {
//		return &data, &dataDB, exterror.WrapExtError(err)
//	}
//	// set VJ_charging current search
//	if reportVJCharging, ok := ctx.Get(RPComon.REPORT_VJ_CHARGING_KEY); ok {
//		data.VJCharging = reportVJCharging.(int)
//	}
//	if totalRows > RPComon.BQ_DATA_LIMIT {
//		return &data, &dataDB, exterror.WrapExtError(errors.New("Respone data too large"))
//	} else if totalRows == 0 {
//		return &data, &dataDB, nil
//	}
//
//	// Get data
//	keyErr = errors.New("KEY_ERR")
//	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
//	dataChan, err := conn.GetResponseData(jobId, 0, RPComon.BQ_DATA_LIMIT, ctx, msgRetryTmp, retryMap)
//
//	if err != nil {
//		return &data, &dataDB, exterror.WrapExtError(err)
//	}
//	index := 0
//	if form.ControlType == CONTROL_TYPE_JAN {
//		rowsMapData := map[string][]interface{}{}
//		for {
//			row := <-dataChan
//			if row == nil {
//				break
//			}
//
//			// Add data for list rows
//			dataHaveRankNo := []interface{}{}
//			dataHaveRankNo = append(dataHaveRankNo,
//				row.ValueMap["jan_code"].String(),
//				row.ValueMap["product_name"].String(),
//				row.ValueMap["author_name"].String(),
//				row.ValueMap["maker_name"].String(),
//				row.ValueMap["list_price"].Float(),
//				row.ValueMap["stock_quantity"].Int(),
//				row.ValueMap["last_sales_date"].String(),
//				row.ValueMap["cumulative_sales_quantity"].Int(),
//				row.ValueMap["cumulative_receiving_quantity"].Int(),
//				row.ValueMap["selling_date"].String(),
//			)
//
//			// Get data 売上初速01 -> 売上初速40
//			for key := 1; key <= form.SearchDateType; key++ {
//				if data.GraphData[key] == nil {
//					data.GraphData[key] = []interface{}{}
//					data.GraphData[key] = append(data.GraphData[key], key)
//				}
//
//				value := row.ValueMap["sales_quantity_day"+strconv.Itoa(key)].Int()
//				data.GraphData[key] = append(
//					data.GraphData[key],
//					value,
//					strconv.Itoa(key)+"日目 "+ListRank[index]+": "+strconv.Itoa(int(value))+"冊",
//				)
//				dataHaveRankNo = append(dataHaveRankNo, value)
//				if value > data.MaxValue {
//					data.MaxValue = value
//				}
//			}
//			rowsMapData[row.ValueMap["jan_code"].String()] = dataHaveRankNo
//			// Set column name for graph
//			data.GraphSymbol = append(data.GraphSymbol, ListRank[index])
//			index++
//		}
//		for _, janCd := range form.JanArrays {
//			if rowsMapData[janCd] != nil {
//				data.Rows = append(data.Rows, rowsMapData[janCd])
//			}
//		}
//		dataDB = data
//	} else {
//		for {
//			row := <-dataChan
//			if row == nil {
//				break
//			}
//
//			// Get info magazine_name and maker_name (one time)
//			if index == 0 {
//				data.MagazineName = row.ValueMap["magazine_name"].String()
//				data.MakerName = row.ValueMap["maker_name"].String()
//			}
//
//			// Add data for list rows
//			dataHaveRankNo := []interface{}{}
//			// ASO-5486 初速比較 雑誌コードで検索した時、user_janへの登録がおかしくなる
//			// Init data for insert to database
//			dataHaveRankNoDB := []interface{}{}
//			sales := row.ValueMap["cumulative_sales_quantity"].Float()
//			received := row.ValueMap["cumulative_receiving_quantity"].Float()
//			rate := 0
//			if received > 0 {
//				rate = int(100 * (sales / received))
//			}
//			dataHaveRankNo = append(dataHaveRankNo,
//				row.ValueMap["jan_code"].String(),
//				row.ValueMap["month_num"].String(),
//				row.ValueMap["list_price"].Float(),
//				row.ValueMap["stock_quantity"].Int(),
//				row.ValueMap["last_sales_date"].String(),
//				sales,
//				received,
//				rate,
//				row.ValueMap["selling_date"].String(),
//			)
//			// ASO-5486 初速比較 雑誌コードで検索した時、user_janへの登録がおかしくなる
//			dataHaveRankNoDB = append(dataHaveRankNoDB,
//				row.ValueMap["jan_code"].String(),
//				row.ValueMap["magazine_name"].String(),
//				row.ValueMap["author_name"].String(),
//				row.ValueMap["maker_name"].String(),
//				row.ValueMap["list_price"].Float(),
//				row.ValueMap["stock_quantity"].Int(),
//				row.ValueMap["last_sales_date"].String(),
//				row.ValueMap["cumulative_sales_quantity"].Int(),
//				row.ValueMap["cumulative_receiving_quantity"].Int(),
//				row.ValueMap["selling_date"].String(),
//			)
//
//			// Get data 売上初速01 -> 売上初速40
//			for key := 1; key <= form.SearchDateType; key++ {
//				if data.GraphData[key] == nil {
//					data.GraphData[key] = []interface{}{}
//					data.GraphData[key] = append(data.GraphData[key], key)
//				}
//
//				value := row.ValueMap["sales_quantity_day"+strconv.Itoa(key)].Int()
//
//				// Graph only display 15 item
//				if index < 15 {
//					data.GraphData[key] = append(
//						data.GraphData[key],
//						value,
//						strconv.Itoa(key)+"日目 "+ListRank[index]+": "+strconv.Itoa(int(value))+"冊",
//					)
//					if value > data.MaxValue {
//						data.MaxValue = value
//					}
//				}
//				dataHaveRankNo = append(dataHaveRankNo, value)
//			}
//			data.Rows = append(data.Rows, dataHaveRankNo)
//			dataDB.Rows = append(dataDB.Rows, dataHaveRankNoDB)
//			// Set column name for graph
//			if index < 15 {
//				data.GraphSymbol = append(data.GraphSymbol, ListRank[index])
//			}
//			index++
//		}
//	}
//	data.MaxValue += (data.MaxValue * 5) / 100
//	data.CountResultRows = index
//	return &data, &dataDB, nil
//}
//
////Create data from queryBuild
//func queryDataDetail(ctx *gf.Context, sql string, form QueryForm) (*RpData, error) {
//
//	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
//	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
//	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)
//
//	data := RpData{
//		HeaderCols: []string{},
//		Rows:       [][]interface{}{},
//	}
//
//	keyErr := errors.New("KEY_ERR")
//	msgRetryTmp := strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
//	// ASO-5502 JOBのコンパイルが通らない（stg）
//	retryMap := make(map[string]string)
//	retryMap["CONFIG_RETRY_COUNT"] = WebApp.CONFIG_RETRY_COUNT
//	retryMap["CONFIG_TIME_RETRY"] = WebApp.CONFIG_TIME_RETRY
//	retryMap["CONFIG_LIST_CODE_HTTP"] = WebApp.CONFIG_LIST_CODE_HTTP
//	conn, err := bq.NewConnection(keyFile, mailAccount, projectId, ctx, msgRetryTmp, retryMap)
//
//	if err != nil {
//		return &data, exterror.WrapExtError(err)
//	}
//
//	// set report name to import info log search charging
//	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
//	//totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, _REPORT_ID)
//	// ========================================================================================
//	// Output log search condition
//	tag := "report=" + _REPORT_ID
//	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
//		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_DOWNLOAD_TEXT + `"`
//	} else {
//		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
//	}
//	if form.ControlType == CONTROL_TYPE_JAN {
//		if len(form.JanArrays) > 0 {
//			tag = tag + ",JAN IN (" + strings.Join(form.JanArrays, ",") + ")"
//		}
//	}
//	// ========================================================================================
//	keyErr = errors.New("KEY_ERR")
//	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
//	totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, tag, ctx, msgRetryTmp, retryMap)
//
//	if err != nil {
//		return &data, exterror.WrapExtError(err)
//	}
//	// set VJ_charging current search
//	if reportVJCharging, ok := ctx.Get(RPComon.REPORT_VJ_CHARGING_KEY); ok {
//		data.VJCharging = reportVJCharging.(int)
//	}
//	if totalRows > RPComon.BQ_DATA_LIMIT {
//		return &data, exterror.WrapExtError(errors.New("Respone data too large"))
//	} else if totalRows == 0 {
//		return &data, nil
//	}
//
//	// Get data
//	keyErr = errors.New("KEY_ERR")
//	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
//	dataChan, err := conn.GetResponseData(jobId, 0, RPComon.BQ_DATA_LIMIT, ctx, msgRetryTmp, retryMap)
//
//	if err != nil {
//		return &data, exterror.WrapExtError(err)
//	}
//
//	for {
//		row := <-dataChan
//		if row == nil {
//			break
//		}
//
//		// Add data for list rows
//		dataSingle := []interface{}{}
//		dataSingle = append(dataSingle,
//			row.ValueMap["shop_name"].String(),
//			row.ValueMap["cumulative_sales_quantity"].Int(),
//			row.ValueMap["cumulative_receiving_quantity"].Int(),
//			row.ValueMap["stock_quantity"].Int(),
//			row.ValueMap["first_sales_date"].String(),
//		)
//
//		// Get data 売上初速01 -> 売上初速40
//		for key := 1; key <= form.SearchDateType; key++ {
//			value := row.ValueMap["sales_quantity_day"+strconv.Itoa(key)].Int()
//			dataSingle = append(dataSingle, value)
//		}
//		data.Rows = append(data.Rows, dataSingle)
//		// Set column name for graph
//	}
//
//	return &data, nil
//}
//
////Create data from queryBuild
//func queryDataDetailForCSV(ctx *gf.Context, sql string, csvWriter *csv.Writer, listHeader []string, form QueryForm) error {
//
//	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
//	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
//	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)
//
//	keyErr := errors.New("KEY_ERR")
//	msgRetryTmp := strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
//	// ASO-5502 JOBのコンパイルが通らない（stg）
//	retryMap := make(map[string]string)
//	retryMap["CONFIG_RETRY_COUNT"] = WebApp.CONFIG_RETRY_COUNT
//	retryMap["CONFIG_TIME_RETRY"] = WebApp.CONFIG_TIME_RETRY
//	retryMap["CONFIG_LIST_CODE_HTTP"] = WebApp.CONFIG_LIST_CODE_HTTP
//	conn, err := bq.NewConnection(keyFile, mailAccount, projectId, ctx, msgRetryTmp, retryMap)
//
//	if err != nil {
//		return exterror.WrapExtError(err)
//	}
//
//	// set report name to import info log search charging
//	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
//	//totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, _REPORT_ID)
//	// ========================================================================================
//	// Output log search condition
//	tag := "report=" + _REPORT_ID
//	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
//		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_DOWNLOAD_TEXT + `"`
//	} else {
//		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
//	}
//	if form.ControlType == CONTROL_TYPE_JAN {
//		if len(form.JanArrays) > 0 {
//			tag = tag + ",JAN IN (" + strings.Join(form.JanArrays, ",") + ")"
//		}
//	}
//	// ========================================================================================
//	keyErr = errors.New("KEY_ERR")
//	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
//	totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, tag, ctx, msgRetryTmp, retryMap)
//
//	if err != nil {
//		return exterror.WrapExtError(err)
//	}
//
//	if totalRows > RPComon.BQ_DATA_LIMIT {
//		return exterror.WrapExtError(errors.New("Respone data too large"))
//	} else if totalRows == 0 {
//		return nil
//	}
//
//	// Get data
//	keyErr = errors.New("KEY_ERR")
//	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
//	dataChan, err := conn.GetResponseData(jobId, 0, RPComon.BQ_DATA_LIMIT, ctx, msgRetryTmp, retryMap)
//
//	if err != nil {
//		return exterror.WrapExtError(err)
//	}
//	if form.ControlType == CONTROL_TYPE_JAN {
//		rowsInterface := map[string][][]string{}
//		for {
//			row := <-dataChan
//			if row == nil {
//				break
//			}
//			janKeyRange := row.ValueMap["jan_code"].String()
//			if rowsInterface[janKeyRange] == nil {
//				rowsInterface[janKeyRange] = [][]string{}
//			}
//
//			// Add data for list rows
//			dataSingle := []string{}
//			dataSingle = append(dataSingle,
//				row.ValueMap["shop_name"].String(),
//				strconv.Itoa(int(row.ValueMap["cumulative_sales_quantity"].Int())),
//				strconv.Itoa(int(row.ValueMap["cumulative_receiving_quantity"].Int())),
//				strconv.Itoa(int(row.ValueMap["stock_quantity"].Int())),
//				row.ValueMap["first_sales_date"].String(),
//			)
//
//			// Get data 売上初速01 -> 売上初速40
//			for key := 1; key <= form.SearchDateType; key++ {
//				value := row.ValueMap["sales_quantity_day"+strconv.Itoa(key)].Int()
//				dataSingle = append(dataSingle, strconv.Itoa(int(value)))
//			}
//			rowsInterface[janKeyRange] = append(rowsInterface[janKeyRange], dataSingle)
//		}
//		for _, janCd := range form.JanArrays {
//			if rowsInterface[janCd] != nil {
//				csvWriter.Write([]string{})
//				csvWriter.Write([]string{"JAN", janCd})
//				csvWriter.Write(listHeader)
//				for _, row := range rowsInterface[janCd] {
//					csvWriter.Write(row)
//				}
//			}
//		}
//	} else {
//		index := 0
//		janKeyRange := ""
//		for {
//			row := <-dataChan
//			if row == nil {
//				break
//			}
//
//			if index == 0 {
//				csvWriter.Write([]string{})
//				janKeyRange = row.ValueMap["jan_code"].String()
//				csvWriter.Write([]string{"JAN", janKeyRange})
//				csvWriter.Write(listHeader)
//			}
//			janKeyItem := row.ValueMap["jan_code"].String()
//			if janKeyItem != janKeyRange {
//				janKeyRange = janKeyItem
//				csvWriter.Write([]string{})
//				csvWriter.Write([]string{"JAN", janKeyRange})
//				csvWriter.Write(listHeader)
//			}
//
//			// Add data for list rows
//			dataSingle := []string{}
//			dataSingle = append(dataSingle,
//				row.ValueMap["shop_name"].String(),
//				strconv.Itoa(int(row.ValueMap["cumulative_sales_quantity"].Int())),
//				strconv.Itoa(int(row.ValueMap["cumulative_receiving_quantity"].Int())),
//				strconv.Itoa(int(row.ValueMap["stock_quantity"].Int())),
//				row.ValueMap["first_sales_date"].String(),
//			)
//
//			// Get data 売上初速01 -> 売上初速40
//			for key := 1; key <= form.SearchDateType; key++ {
//				value := row.ValueMap["sales_quantity_day"+strconv.Itoa(key)].Int()
//				dataSingle = append(dataSingle, strconv.Itoa(int(value)))
//			}
//			csvWriter.Write(dataSingle)
//			index++
//		}
//	}
//
//	return nil
//}

//Create data from queryBuild
func queryData(ctx *gf.Context, sql string, form QueryForm) (*RpData, *RpData, error) {

	data := RpData{
		HeaderCols: []string{},
		Rows:       [][]interface{}{},
		GraphData:  map[int][]interface{}{},
		MaxValue:   0,
	}
	dataDB := RpData{
		HeaderCols: []string{},
		Rows:       [][]interface{}{},
		GraphData:  map[int][]interface{}{},
		MaxValue:   0,
	}
	// ASO-5719 [BA]bq_log_search_chargingに項目を追加して検索条件の一部を追加 - ADD START
	// set report name to import info log search charging
	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
	// ASO-5719 [BA]bq_log_search_chargingに項目を追加して検索条件の一部を追加 - ADD END

	index := 0
	rowsMapData := map[string][]interface{}{}
	if form.ControlType == CONTROL_TYPE_JAN {
		// Get data from CLOUD
		var args []interface{}
		for _, s := range form.ShopCd {
			args = append(args, s)
		}
		for _, s := range form.JanArrays {
			args = append(args, s)
		}
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT START
		//dbCloud, err := CloudSQL.Connect()
		//if err != nil {
		//	return &data, &dataDB, exterror.WrapExtError(err)
		//}
		//
		//rows, err := dbCloud.Query(sql, args...)
		//if err != nil {
		//	return &data, &dataDB, exterror.WrapExtError(err)
		//}
		//
		//defer func() {
		//	rows.Close()
		//	dbCloud.Close()
		//}()
		rows, err := ctx.DB.Query(sql, args...)
		if err != nil {
			return &data, &dataDB, exterror.WrapExtError(err)
		}

		defer func() {
			rows.Close()
		}()
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END

		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - EDIT START
		listDataGraphMap := make(map[string]map[int]int64)
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - EDIT END
		for rows.Next() {
			newMStockItem := ModelItems.MStockItem{}
			err = db.SqlScanStruct(rows, &newMStockItem)
			if err != nil {
				return &data, &dataDB, exterror.WrapExtError(err)
			}

			// Add data for list rows
			dataHaveRankNo := []interface{}{}
			dataHaveRankNo = append(dataHaveRankNo,
				newMStockItem.JanCode,
				newMStockItem.ProductName,
				newMStockItem.AuthorName,
				newMStockItem.MakerName,
				newMStockItem.ListPrice,
				newMStockItem.StockQuantity,
				newMStockItem.LastSalesDate,
				newMStockItem.CumulativeSalesQuantity,
				newMStockItem.CumulativeReceivingQuantity,
				newMStockItem.SellingDate,
			)

			// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - EDIT START
			//// Get data 売上初速01 -> 売上初速40
			//for key := 1; key <= form.SearchDateType; key++ {
			//    if data.GraphData[key] == nil {
			//        data.GraphData[key] = []interface{}{}
			//        data.GraphData[key] = append(data.GraphData[key], key)
			//    }
			//
			//    r := reflect.ValueOf(newMStockItem)
			//    value := int64(reflect.Indirect(r).FieldByName("SalesQuantityDay" + strconv.Itoa(key)).Int())
			//    data.GraphData[key] = append(
			//        data.GraphData[key],
			//        value,
			//        strconv.Itoa(key)+"日目 "+ListRank[index]+": "+strconv.Itoa(int(value))+"冊",
			//    )
			//    dataHaveRankNo = append(dataHaveRankNo, value)
			//    if value > data.MaxValue {
			//        data.MaxValue = value
			//    }
			//}
			//rowsMapData[newMStockItem.JanCode] = dataHaveRankNo
			//// Set column name for graph
			//data.GraphSymbol = append(data.GraphSymbol, ListRank[index])
			dataGraphMap := make(map[int]int64)
			for key := 1; key <= form.SearchDateType; key++ {
				r := reflect.ValueOf(newMStockItem)
				value := int64(reflect.Indirect(r).FieldByName("SalesQuantityDay" + strconv.Itoa(key)).Int())
				dataGraphMap[key] = value
				dataHaveRankNo = append(dataHaveRankNo, value)
				if value > data.MaxValue {
					data.MaxValue = value
				}
			}
			listDataGraphMap[newMStockItem.JanCode] = dataGraphMap
			rowsMapData[newMStockItem.JanCode] = dataHaveRankNo
			// Set column name for graph
			data.GraphSymbol = append(data.GraphSymbol, ListRank[index])
			// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - EDIT END
			index++
		}
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
		// Sort data from map mapping with list JAN input
		for i, janCd := range form.JanArrays {
			if listDataGraphMap[janCd] != nil {
				for key := 1; key <= form.SearchDateType; key++ {
					if data.GraphData[key] == nil {
						data.GraphData[key] = []interface{}{}
						data.GraphData[key] = append(data.GraphData[key], key)
					}
					value := int64(listDataGraphMap[janCd][key])
					data.GraphData[key] = append(
						data.GraphData[key],
						value,
						strconv.Itoa(key)+"日目 "+ListRank[i]+": "+strconv.Itoa(int(value))+"冊",
					)
				}
			}
		}
		// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - START END
		for _, janCd := range form.JanArrays {
			if rowsMapData[janCd] != nil {
				data.Rows = append(data.Rows, rowsMapData[janCd])
			}
		}
		dataDB = data
	} else {
		// Get data from CLOUD
		var args []interface{}
		for _, s := range form.ShopCd {
			args = append(args, s)
		}
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT START
		//dbCloud, err := CloudSQL.Connect()
		//if err != nil {
		//	return &data, &dataDB, exterror.WrapExtError(err)
		//}
		//
		//rows, err := dbCloud.Query(sql, args...)
		//if err != nil {
		//	return &data, &dataDB, exterror.WrapExtError(err)
		//}
		//
		//defer func() {
		//	rows.Close()
		//	dbCloud.Close()
		//}()
		rows, err := ctx.DB.Query(sql, args...)
		if err != nil {
			return &data, &dataDB, exterror.WrapExtError(err)
		}

		defer func() {
			rows.Close()
		}()
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END

		for rows.Next() {
			newMStockItem := ModelItems.MStockItem{}
			err = db.SqlScanStruct(rows, &newMStockItem)
			if err != nil {
				return &data, &dataDB, exterror.WrapExtError(err)
			}

			// Get info magazine_name and maker_name (one time)
			if index == 0 {
				data.MagazineName = newMStockItem.MagazineName
				data.MakerName = newMStockItem.MakerName
			}

			// Add data for list rows
			dataHaveRankNo := []interface{}{}
			// ASO-5486 初速比較 雑誌コードで検索した時、user_janへの登録がおかしくなる
			// Init data for insert to database
			dataHaveRankNoDB := []interface{}{}
			sales := float64(newMStockItem.CumulativeSalesQuantity)
			received := float64(newMStockItem.CumulativeReceivingQuantity)
			rate := 0
			if received > 0 {
				rate = int(100 * (sales / received))
			}
			dataHaveRankNo = append(dataHaveRankNo,
				newMStockItem.JanCode,
				newMStockItem.ProductName,
				newMStockItem.MagazineName,
				newMStockItem.MakerName,
				//newMStockItem.MonthNum,
				newMStockItem.ListPrice,
				newMStockItem.StockQuantity,
				newMStockItem.LastSalesDate,
				sales,
				received,
				rate,
				newMStockItem.SellingDate,
			)
			// ASO-5486 初速比較 雑誌コードで検索した時、user_janへの登録がおかしくなる
			// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT START
			//dataHaveRankNoDB = append(dataHaveRankNoDB,
			//	newMStockItem.JanCode,
			//	newMStockItem.MagazineName,
			//	newMStockItem.AuthorName,
			//	newMStockItem.MakerName,
			//	newMStockItem.ListPrice,
			//	newMStockItem.StockQuantity,
			//	newMStockItem.LastSalesDate,
			//	newMStockItem.CumulativeSalesQuantity,
			//	newMStockItem.CumulativeReceivingQuantity,
			//	newMStockItem.SellingDate,
			//)
			dataHaveRankNoDB = append(dataHaveRankNoDB,
				newMStockItem.JanCode,
				newMStockItem.ProductName,
				newMStockItem.AuthorName,
				newMStockItem.MakerName,
				newMStockItem.ListPrice,
				newMStockItem.StockQuantity,
				newMStockItem.LastSalesDate,
				newMStockItem.CumulativeSalesQuantity,
				newMStockItem.CumulativeReceivingQuantity,
				newMStockItem.SellingDate,
			)
			// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT END

			// Get data 売上初速01 -> 売上初速40
			for key := 1; key <= form.SearchDateType; key++ {
				if data.GraphData[key] == nil {
					data.GraphData[key] = []interface{}{}
					data.GraphData[key] = append(data.GraphData[key], key)
				}

				r := reflect.ValueOf(newMStockItem)
				value := int64(reflect.Indirect(r).FieldByName("SalesQuantityDay" + strconv.Itoa(key)).Int())

				// Graph only display 15 item
				if index < 15 {
					data.GraphData[key] = append(
						data.GraphData[key],
						value,
						strconv.Itoa(key)+"日目 "+ListRank[index]+": "+strconv.Itoa(int(value))+"冊",
					)
					if value > data.MaxValue {
						data.MaxValue = value
					}
				}
				dataHaveRankNo = append(dataHaveRankNo, value)
			}
			data.Rows = append(data.Rows, dataHaveRankNo)
			dataDB.Rows = append(dataDB.Rows, dataHaveRankNoDB)
			// Set column name for graph
			if index < 15 {
				data.GraphSymbol = append(data.GraphSymbol, ListRank[index])
			}
			index++
		}
		for _, janCd := range form.JanArrays {
			if rowsMapData[janCd] != nil {
				data.Rows = append(data.Rows, rowsMapData[janCd])
			}
		}
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - DEL START
		//dataDB = data
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - DEL END
	}
	data.MaxValue += (data.MaxValue * 5) / 100
	data.CountResultRows = index

	return &data, &dataDB, nil
}

//Create data from queryBuild
func queryDataDetail(ctx *gf.Context, sql string, form QueryForm) (*RpData, error) {

	data := RpData{
		HeaderCols: []string{},
		Rows:       [][]interface{}{},
	}

	// Get data from CLOUD
	var args []interface{}
	for _, s := range form.ShopCd {
		args = append(args, s)
	}
	for _, s := range form.JanArrays {
		args = append(args, s)
	}
	// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT START
	//dbCloud, err := CloudSQL.Connect()
	//if err != nil {
	//	return &data, exterror.WrapExtError(err)
	//}
	//
	//rows, err := dbCloud.Query(sql, args...)
	//if err != nil {
	//	return &data, exterror.WrapExtError(err)
	//}
	//
	//defer func() {
	//	rows.Close()
	//	dbCloud.Close()
	//}()
	rows, err := ctx.DB.Query(sql, args...)
	if err != nil {
		return &data, exterror.WrapExtError(err)
	}

	defer func() {
		rows.Close()
	}()
	// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END

	for rows.Next() {
		newMStockItem := ModelItems.MStockItem{}
		err = db.SqlScanStruct(rows, &newMStockItem)
		if err != nil {
			return &data, exterror.WrapExtError(err)
		}

		// Add data for list rows
		dataSingle := []interface{}{}
		dataSingle = append(dataSingle,
			newMStockItem.ShopName,
			newMStockItem.CumulativeSalesQuantity,
			newMStockItem.CumulativeReceivingQuantity,
			newMStockItem.StockQuantity,
			newMStockItem.FirstSalesDate,
		)

		// Get data 売上初速01 -> 売上初速40
		for key := 1; key <= form.SearchDateType; key++ {
			r := reflect.ValueOf(newMStockItem)
			value := int64(reflect.Indirect(r).FieldByName("SalesQuantityDay" + strconv.Itoa(key)).Int())
			dataSingle = append(dataSingle, value)
		}
		data.Rows = append(data.Rows, dataSingle)
		// Set column name for graph
	}

	return &data, nil
}

//Create data from queryBuild
func queryDataDetailForCSV(ctx *gf.Context, sql string, csvWriter *csv.Writer, listJanKey []string, listHeader []string, form QueryForm) error {

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
	//// ========================================================================================

	// Get data from CLOUD
	if form.ControlType == CONTROL_TYPE_JAN {
		var args []interface{}
		for _, s := range form.ShopCd {
			args = append(args, s)
		}
		for _, s := range listJanKey {
			args = append(args, s)
		}
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT START
		//dbCloud, err := CloudSQL.Connect()
		//if err != nil {
		//	return exterror.WrapExtError(err)
		//}
		//
		//rows, err := dbCloud.Query(sql, args...)
		//if err != nil {
		//	return exterror.WrapExtError(err)
		//}
		//
		//defer func() {
		//	rows.Close()
		//	dbCloud.Close()
		//}()
		rows, err := ctx.DB.Query(sql, args...)
		if err != nil {
			return exterror.WrapExtError(err)
		}

		defer func() {
			rows.Close()
		}()
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END

		rowsInterface := map[string][][]string{}
		for rows.Next() {
			newMStockItem := ModelItems.MStockItem{}
			err = db.SqlScanStruct(rows, &newMStockItem)
			if err != nil {
				return exterror.WrapExtError(err)
			}

			janKeyRange := newMStockItem.JanCode
			if rowsInterface[janKeyRange] == nil {
				rowsInterface[janKeyRange] = [][]string{}
			}

			// Add data for list rows
			dataSingle := []string{}
			dataSingle = append(dataSingle,
				newMStockItem.ShopName,
				strconv.Itoa(int(newMStockItem.CumulativeSalesQuantity)),
				strconv.Itoa(int(newMStockItem.CumulativeReceivingQuantity)),
				strconv.Itoa(int(newMStockItem.StockQuantity)),
				newMStockItem.FirstSalesDate,
			)

			// Get data 売上初速01 -> 売上初速40
			for key := 1; key <= form.SearchDateType; key++ {
				r := reflect.ValueOf(newMStockItem)
				value := int64(reflect.Indirect(r).FieldByName("SalesQuantityDay" + strconv.Itoa(key)).Int())
				dataSingle = append(dataSingle, strconv.Itoa(int(value)))
			}
			rowsInterface[janKeyRange] = append(rowsInterface[janKeyRange], dataSingle)
		}
		for _, janCd := range form.JanArrays {
			if rowsInterface[janCd] != nil {
				csvWriter.Write([]string{})
				csvWriter.Write([]string{"JAN", janCd})
				csvWriter.Write(listHeader)
				for _, row := range rowsInterface[janCd] {
					csvWriter.Write(row)
				}
			}
		}
	} else {
		index := 0
		janKeyRange := ""
		var args []interface{}
		for _, s := range form.ShopCd {
			args = append(args, s)
		}
		for _, s := range listJanKey {
			args = append(args, s)
		}
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT START
		//dbCloud, err := CloudSQL.Connect()
		//if err != nil {
		//	return exterror.WrapExtError(err)
		//}
		//
		//rows, err := dbCloud.Query(sql, args...)
		//if err != nil {
		//	return exterror.WrapExtError(err)
		//}
		//
		//defer func() {
		//	rows.Close()
		//	dbCloud.Close()
		//}()
		rows, err := ctx.DB.Query(sql, args...)
		if err != nil {
			return exterror.WrapExtError(err)
		}

		defer func() {
			rows.Close()
		}()
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END

		for rows.Next() {
			newMStockItem := ModelItems.MStockItem{}
			err = db.SqlScanStruct(rows, &newMStockItem)
			if err != nil {
				return exterror.WrapExtError(err)
			}

			if index == 0 {
				csvWriter.Write([]string{})
				janKeyRange = newMStockItem.JanCode
				csvWriter.Write([]string{"JAN", janKeyRange})
				csvWriter.Write(listHeader)
			}
			janKeyItem := newMStockItem.JanCode
			if janKeyItem != janKeyRange {
				janKeyRange = janKeyItem
				csvWriter.Write([]string{})
				csvWriter.Write([]string{"JAN", janKeyRange})
				csvWriter.Write(listHeader)
			}

			// Add data for list rows
			dataSingle := []string{}
			dataSingle = append(dataSingle,
				newMStockItem.ShopName,
				strconv.Itoa(int(newMStockItem.CumulativeSalesQuantity)),
				strconv.Itoa(int(newMStockItem.CumulativeReceivingQuantity)),
				strconv.Itoa(int(newMStockItem.StockQuantity)),
				newMStockItem.FirstSalesDate,
			)

			// Get data 売上初速01 -> 売上初速40
			for key := 1; key <= form.SearchDateType; key++ {
				r := reflect.ValueOf(newMStockItem)
				value := int64(reflect.Indirect(r).FieldByName("SalesQuantityDay" + strconv.Itoa(key)).Int())
				dataSingle = append(dataSingle, strconv.Itoa(int(value)))
			}
			csvWriter.Write(dataSingle)
			index++
		}
	}

	return nil
}
