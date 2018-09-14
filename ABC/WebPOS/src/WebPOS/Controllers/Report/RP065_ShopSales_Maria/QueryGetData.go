package RP065_ShopSales_Maria

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"errors"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"strings"
	"WebPOS/Models/ModelItems"
	"github.com/goframework/gf/db"
	"strconv"
)

func queryData(ctx *gf.Context, sql string, form QueryForm, randStringFromSQL string, exCols []string) (*RpData, error) {

	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)

	data := RpData{
		HeaderCols: []string{},
		Rows:       []SingleItem{},
	}

	keyErr := errors.New("KEY_ERR")
	msgRetryTmp := strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	// ASO-5502 JOBのコンパイルが通らない（stg）
	retryMap := make(map[string]string)
	retryMap["CONFIG_RETRY_COUNT"] = WebApp.CONFIG_RETRY_COUNT
	retryMap["CONFIG_TIME_RETRY"] = WebApp.CONFIG_TIME_RETRY
	retryMap["CONFIG_LIST_CODE_HTTP"] = WebApp.CONFIG_LIST_CODE_HTTP
	conn, err := bq.NewConnection(keyFile, mailAccount, projectId, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return nil, exterror.WrapExtError(err)
	}

	// set report name to import info log search charging
	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
	//totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, _REPORT_ID)
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
	if len(form.JanSingle) > 0 {
		tag = tag + ",JAN LIKE " + `"` + form.JanSingle + `%"`
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
	if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT {
		tag = tag + ",フォーマット=" + `"` + DOWNLOAD_TYPE_TOTAL_RESULT_TEXT + `"`
	} else if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION {
		tag = tag + ",フォーマット=" + `"` + DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION_TEXT + `"`
	}
	// ========================================================================================
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, tag, ctx, msgRetryTmp, retryMap)

	intTotalRows := int(totalRows)
	ctx.Session.Values[RPComon.REPORT_QUERY_JOB_ID+randStringFromSQL] = jobId
	ctx.Session.Values[RPComon.REPORT_QUERY_JOB_ID_COUNT+randStringFromSQL] = intTotalRows
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}

	if totalRows > RPComon.BQ_DATA_LIMIT {
		return nil, exterror.WrapExtError(errors.New("Respone data too large"))
	}
	//=====================================================================================
	// set report name to import info log search charging
	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
	// Search count Data by shop
	// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
	//sqlCount, listRange, groupType := buildSql(form, ctx, true)
	sqlCount, _, _, exCols := buildSql(form, ctx, true)
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	_, jobIdCount, err := conn.QueryForResponseBySql(sqlCount, ctx, tag, ctx, msgRetryTmp, retryMap)

	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChanCount, err := conn.GetResponseData(jobIdCount, 0, RPComon.BQ_DATA_LIMIT, ctx, msgRetryTmp, retryMap)

	countShop := map[string]int64{}
	for {
		row := <-dataChanCount
		if row == nil {
			break
		}
		count := row.ValueMap["count_data"].Int()
		shop := row.ValueMap["shop_code"].String()
		ctx.Session.Values[_REPORT_ID+randStringFromSQL+shop] = count
		countShop[shop] = count
	}

	//=====================================================================================

	// set paging
	data.CountResultRows = intTotalRows
	data.PageCount = intTotalRows / form.Limit
	if intTotalRows%form.Limit > 0 {
		data.PageCount += 1
	}
	data.ThisPage = form.Page
	if data.ThisPage < 1 {
		data.ThisPage = 1
	}
	if data.ThisPage > data.PageCount {
		data.ThisPage = data.PageCount
	}

	data.ShowLineFrom = (data.ThisPage - 1) * form.Limit

	limitLengthData := form.Limit
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		data.ShowLineFrom = 0
		limitLengthData = RPComon.BQ_DATA_LIMIT
	}
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChan, err := conn.GetResponseData(jobId, 0, limitLengthData, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return nil, exterror.WrapExtError(err)

	}

	// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
	listJan := []string{}
	singleItem := []SingleItem{}

	shopCodeKey := ""
	item := SingleItem{}
	shm := Models.ShopMasterModel{ctx.DB}
	index := 0
	for {
		row := <-dataChan
		if row == nil {
			// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
			//if item.ShmShopCode != "" {
			//	item.DataCount = countShop[shopCodeKey]
			//	data.Rows = append(data.Rows, item)
			//}
			break
		}

		// Get list JAN from result with BigQuery
		if len(listJan) >= 1000 {
			break
		}
		listJan = append(listJan, row.ValueMap["jan_code"].String())

		shopCd := row.ValueMap["shop_code"].String()
		if index == 0 {
			shopCodeKey = shopCd
		}
		if shopCd != shopCodeKey {
			// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
			//item.DataCount = countShop[shopCodeKey]
			//data.Rows = append(data.Rows, item)
			item = SingleItem{}
			item.DataCount = countShop[shopCodeKey]
			shopCodeKey = shopCd
		}

		if item.ShmShopCode == "" {
			shopInfo, _ := shm.GetInfoShopByCD(shopCd)
			item.ShmSharedBookStoreCode = shopInfo.SharedBookStoreCode
			item.ShmShopCode = shopCd
			item.ShmShopName = shopInfo.ShopName
			item.ShmTelNo = shopInfo.TelNo
			item.ShmBizStartTime = shopInfo.BizStartTime
			item.ShmBizEndTime = shopInfo.BizEndTime
			item.ShmAddress = shopInfo.Address
			item.ShmShopNameShort = shopInfo.ShopNameShort
		}
		// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
		//singleRowData := []interface{}{}
		//singleRowData = append(
		//	singleRowData,
		//	item.ShmShopNameShort,
		//	row.ValueMap["rank"].String(),
		//	row.ValueMap["jan_code"].String(),
		//	row.ValueMap["product_name"].String(),
		//	row.ValueMap["author_name"].String(),
		//	row.ValueMap["maker_name"].String(),
		//	row.ValueMap["selling_date"].String(),
		//	row.ValueMap["sales_tax_exc_unit_price"].Float(),
		//	row.ValueMap["stok_cumulative_receiving_quantity"].Int(),
		//	row.ValueMap["stok_cumulative_sales_quantity"].Int(),
		//	row.ValueMap["stok_stock_quantity"].Int(),
		//	row.ValueMap["first_sales_date"].String(),
		//	row.ValueMap["sales_body_quantity"].Int(),
		//)
		//if len(listRange) > 0 {
		//	for _, item := range listRange {
		//		switch groupType {
		//		case GROUP_TYPE_DATE:
		//			singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcmm+item.Mcdd].Int())
		//		case GROUP_TYPE_WEEK:
		//			singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcweeknum].Int())
		//		case GROUP_TYPE_MONTH:
		//			singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcmm].Int())
		//		}
		//	}
		//}
		//item.Data = append(item.Data, singleRowData)
		//index++

		item.Rank = row.ValueMap["rank"].String()
		item.JanCd = row.ValueMap["jan_code"].String()
		item.ProductName = row.ValueMap["product_name"].String()
		item.AuthorName = row.ValueMap["author_name"].String()
		item.MakerName = row.ValueMap["maker_name"].String()
		item.SellingDate = row.ValueMap["selling_date"].String()
		item.SalesTaxExcUnitPrice = row.ValueMap["sales_tax_exc_unit_price"].Float()
		item.SalesBodyQuantity = row.ValueMap["sales_body_quantity"].Int()
		item.SalesBodyQuantityShop = []int64{}
		for _, v := range exCols {
			item.SalesBodyQuantityShop = append(item.SalesBodyQuantityShop, row.ValueMap["A" + v].Int())
		}
		// Add list result to an array struct
		singleItem = append(singleItem, item)
		index++
	}

	// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
	if len(listJan) > 0 {
		// Get 累計、在庫数、初売上日 from m_stock of MariaDB
		listMStock := []ModelItems.MStockItem{}
		sqlString := `
SELECT
	shop_code,
	jan_code,
	SUM(stock_quantity) stock_quantity,
	SUM(cumulative_receiving_quantity) cumulative_receiving_quantity,
	SUM(cumulative_sales_quantity) cumulative_sales_quantity,
	min(first_sales_date) first_sales_date
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
			return nil, exterror.WrapExtError(err)
		}
		defer rows.Close()
		for rows.Next() {
			newMStockItem := ModelItems.MStockItem{}
			err = db.SqlScanStruct(rows, &newMStockItem)
			if err != nil {
				return nil, exterror.WrapExtError(err)
			}
			listMStock = append(listMStock, newMStockItem)
		}

		itemTemp := SingleItem{}
		index = 0
		// Loop list data get from BigQuery
		for i, _ := range singleItem {
			// Loop list data get from m_stock in MariaDB
			for _, mStockItem := range listMStock {
				// Check exists JAN from result BigQuery in m_stock
				if strings.Compare(singleItem[i].JanCd, mStockItem.JanCode) == 0 &&
				strings.Compare(singleItem[i].ShmShopCode, mStockItem.ShopCode) == 0 {
					// In case of exists, update 累計、在庫数、初売上日 from m_stock to item of singleItem
					singleItem[i].StockQuantity = mStockItem.StockQuantity
					singleItem[i].CumulativeReceivingQuantity = mStockItem.CumulativeReceivingQuantity
					singleItem[i].CumulativeSalesQuantity = mStockItem.CumulativeSalesQuantity
					singleItem[i].FirstSalesDate = mStockItem.FirstSalesDate
					break
				}
			}

			// Check first shop in list
			if index == 0 {
				itemTemp = singleItem[i]
				shopCodeKey = singleItem[i].ShmShopCode
			}
			// Check change shop_code in list data
			if singleItem[i].ShmShopCode != shopCodeKey {
				// In case of change, add data from SingleItem to data.Rows
				itemTemp.DataCount = countShop[shopCodeKey]
				data.Rows = append(data.Rows, itemTemp)
				itemTemp = singleItem[i]
				shopCodeKey = singleItem[i].ShmShopCode
			}

			// Create detail data with JAN
			singleRowData := []interface{}{}
			singleRowData = append(
				singleRowData,
				singleItem[i].ShmShopNameShort,
				singleItem[i].Rank,
				singleItem[i].JanCd,
				singleItem[i].ProductName,
				singleItem[i].AuthorName,
				singleItem[i].MakerName,
				singleItem[i].SellingDate,
				singleItem[i].SalesTaxExcUnitPrice,
				singleItem[i].CumulativeReceivingQuantity,
				singleItem[i].CumulativeSalesQuantity,
				singleItem[i].StockQuantity,
				singleItem[i].FirstSalesDate,
				singleItem[i].SalesBodyQuantity,
			)
			if len(exCols) > 0 {
				for _, v := range singleItem[i].SalesBodyQuantityShop {
					singleRowData = append(singleRowData, strconv.Itoa(int(v)))
				}
			}
			itemTemp.Data = append(itemTemp.Data, singleRowData)
			index++
		}
		// Add last SingleItem to data.Rows
		if len(singleItem) > 0 {
			itemTemp.DataCount = countShop[shopCodeKey]
			data.Rows = append(data.Rows, itemTemp)
		}
	}

	// sort data by input
	if len(form.JanArrays) > 0 {
		sortDataByInput(&data, form.JanArrays)
	}
	//======================================================
	return &data, nil
}

func sortDataByInput(data *RpData, listJan []string) {
	for i, item := range data.Rows {
		temp := [][]interface{}{}
		for _, jan := range listJan {
			for _, row := range item.Data {
				if row[2] == jan {
					row[1] = len(temp) + 1
					temp = append(temp, row)
					break
				}
			}
		}
		data.Rows[i].Data = temp
	}
}

//Create data from queryBuild
func queryGetDataWithJobId(ctx *gf.Context, sql string, form QueryForm, randStringFromSQL string, exCols []string) (*RpData, error) {

	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)

	data := RpData{
		HeaderCols: []string{},
		Rows:       []SingleItem{},
	}

	keyErr := errors.New("KEY_ERR")
	msgRetryTmp := strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	// ASO-5502 JOBのコンパイルが通らない（stg）
	retryMap := make(map[string]string)
	retryMap["CONFIG_RETRY_COUNT"] = WebApp.CONFIG_RETRY_COUNT
	retryMap["CONFIG_TIME_RETRY"] = WebApp.CONFIG_TIME_RETRY
	retryMap["CONFIG_LIST_CODE_HTTP"] = WebApp.CONFIG_LIST_CODE_HTTP
	conn, err := bq.NewConnection(keyFile, mailAccount, projectId, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return &data, exterror.WrapExtError(err)
	}

	jobId, _ := ctx.Session.Values[RPComon.REPORT_QUERY_JOB_ID+randStringFromSQL]
	totalRows, _ := ctx.Session.Values[RPComon.REPORT_QUERY_JOB_ID_COUNT+randStringFromSQL]

	if jobId == nil || totalRows == nil {
		newData, err := queryData(ctx, sql, form, randStringFromSQL, exCols)
		return newData, err
	}

	// set paging
	intTotalRows := totalRows.(int)
	data.CountResultRows = intTotalRows
	data.PageCount = intTotalRows / form.Limit
	if intTotalRows%form.Limit > 0 {
		data.PageCount += 1
	}
	data.ThisPage = form.Page
	if data.ThisPage < 1 {
		data.ThisPage = 1
	}
	if data.ThisPage > data.PageCount {
		data.ThisPage = data.PageCount
	}

	data.ShowLineFrom = (data.ThisPage - 1) * form.Limit

	limitLengthData := form.Limit
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		data.ShowLineFrom = 0
		limitLengthData = RPComon.BQ_DATA_LIMIT
	}
	// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
	//_, listRange, groupType, _ := buildSql(form, ctx, false)
	_, _, _, exCols = buildSql(form, ctx, false)
	// Get data
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChan, err := conn.GetResponseData(jobId.(string), data.ShowLineFrom, limitLengthData, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return &data, exterror.WrapExtError(err)
	}

	// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
	listJan := []string{}
	singleItem := []SingleItem{}

	shopCodeKey := ""
	item := SingleItem{}
	shm := Models.ShopMasterModel{ctx.DB}
	index := 0
	for {
		row := <-dataChan
		if row == nil {
			// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
			//if item.ShmShopCode != "" {
			//	item.DataCount = countShop[shopCodeKey]
			//	data.Rows = append(data.Rows, item)
			//}
			break
		}

		// Get list JAN from result with BigQuery
		if len(listJan) >= 1000 {
			break
		}
		listJan = append(listJan, row.ValueMap["jan_code"].String())

		shopCd := row.ValueMap["shop_code"].String()
		if index == 0 {
			shopCodeKey = shopCd
		}
		if shopCd != shopCodeKey {
			// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
			//value := ctx.Session.Values[_REPORT_ID+randStringFromSQL+shopCodeKey]
			//item.DataCount = value.(int64)
			//data.Rows = append(data.Rows, item)
			item = SingleItem{}
			value := ctx.Session.Values[_REPORT_ID+randStringFromSQL+shopCodeKey]
			item.DataCount = value.(int64)
			shopCodeKey = shopCd
		}

		if item.ShmShopCode == "" {
			shopInfo, _ := shm.GetInfoShopByCD(shopCd)
			item.ShmSharedBookStoreCode = shopInfo.SharedBookStoreCode
			item.ShmShopCode = shopCd
			item.ShmShopName = shopInfo.ShopName
			item.ShmTelNo = shopInfo.TelNo
			item.ShmBizStartTime = shopInfo.BizStartTime
			item.ShmBizEndTime = shopInfo.BizEndTime
			item.ShmAddress = shopInfo.Address
			item.ShmShopNameShort = shopInfo.ShopNameShort
		}
		// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
		//singleRowData := []interface{}{}
		//singleRowData = append(
		//	singleRowData,
		//	item.ShmShopNameShort,
		//	row.ValueMap["rank"].String(),
		//	row.ValueMap["jan_code"].String(),
		//	row.ValueMap["product_name"].String(),
		//	row.ValueMap["author_name"].String(),
		//	row.ValueMap["maker_name"].String(),
		//	row.ValueMap["selling_date"].String(),
		//	row.ValueMap["sales_tax_exc_unit_price"].Float(),
		//	row.ValueMap["stok_cumulative_receiving_quantity"].Int(),
		//	row.ValueMap["stok_cumulative_sales_quantity"].Int(),
		//	row.ValueMap["stok_stock_quantity"].Int(),
		//	row.ValueMap["first_sales_date"].String(),
		//	row.ValueMap["sales_body_quantity"].Int(),
		//)
		//if len(listRange) > 0 {
		//	for _, item := range listRange {
		//		switch groupType {
		//		case GROUP_TYPE_DATE:
		//			singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcmm+item.Mcdd].Int())
		//		case GROUP_TYPE_WEEK:
		//			singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcweeknum].Int())
		//		case GROUP_TYPE_MONTH:
		//			singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcmm].Int())
		//		}
		//	}
		//}
		//item.Data = append(item.Data, singleRowData)
		//index++

		item.Rank = row.ValueMap["rank"].String()
		item.JanCd = row.ValueMap["jan_code"].String()
		item.ProductName = row.ValueMap["product_name"].String()
		item.AuthorName = row.ValueMap["author_name"].String()
		item.MakerName = row.ValueMap["maker_name"].String()
		item.SellingDate = row.ValueMap["selling_date"].String()
		item.SalesTaxExcUnitPrice = row.ValueMap["sales_tax_exc_unit_price"].Float()
		item.SalesBodyQuantity = row.ValueMap["sales_body_quantity"].Int()
		item.SalesBodyQuantityShop = []int64{}
		for _, v := range exCols {
			item.SalesBodyQuantityShop = append(item.SalesBodyQuantityShop, row.ValueMap["A" + v].Int())
		}
		// Add list result to an array struct
		singleItem = append(singleItem, item)
		index++
	}

	// ASO-5538 [BA]mBAWEB-v23a ：店舗別集計(Maria)　shop_sales_maria
	if len(listJan) > 0 {
		// Get 累計、在庫数、初売上日 from m_stock of MariaDB
		listMStock := []ModelItems.MStockItem{}
		sqlString := `
SELECT
	shop_code,
	jan_code,
	SUM(stock_quantity) stock_quantity,
	SUM(cumulative_receiving_quantity) cumulative_receiving_quantity,
	SUM(cumulative_sales_quantity) cumulative_sales_quantity,
	min(first_sales_date) first_sales_date
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
			return nil, exterror.WrapExtError(err)
		}
		defer rows.Close()
		for rows.Next() {
			newMStockItem := ModelItems.MStockItem{}
			err = db.SqlScanStruct(rows, &newMStockItem)
			if err != nil {
				return nil, exterror.WrapExtError(err)
			}
			listMStock = append(listMStock, newMStockItem)
		}

		itemTemp := SingleItem{}
		index = 0
		// Loop list data get from BigQuery
		for i, _ := range singleItem {
			// Loop list data get from m_stock in MariaDB
			for _, mStockItem := range listMStock {
				// Check exists JAN from result BigQuery in m_stock
				if strings.Compare(singleItem[i].JanCd, mStockItem.JanCode) == 0 &&
				strings.Compare(singleItem[i].ShmShopCode, mStockItem.ShopCode) == 0 {
					// In case of exists, update 累計、在庫数、初売上日 from m_stock to item of singleItem
					singleItem[i].StockQuantity = mStockItem.StockQuantity
					singleItem[i].CumulativeReceivingQuantity = mStockItem.CumulativeReceivingQuantity
					singleItem[i].CumulativeSalesQuantity = mStockItem.CumulativeSalesQuantity
					singleItem[i].FirstSalesDate = mStockItem.FirstSalesDate
					break
				}
			}

			// Check first shop in list
			if index == 0 {
				itemTemp = singleItem[i]
				shopCodeKey = singleItem[i].ShmShopCode
			}
			// Check change shop_code in list data
			if singleItem[i].ShmShopCode != shopCodeKey {
				// In case of change, add data from SingleItem to data.Rows
				value := ctx.Session.Values[_REPORT_ID + randStringFromSQL + shopCodeKey]
				itemTemp.DataCount = value.(int64)
				data.Rows = append(data.Rows, itemTemp)
				itemTemp = singleItem[i]
				shopCodeKey = singleItem[i].ShmShopCode
			}

			// Create detail data with JAN
			singleRowData := []interface{}{}
			singleRowData = append(
				singleRowData,
				singleItem[i].ShmShopNameShort,
				singleItem[i].Rank,
				singleItem[i].JanCd,
				singleItem[i].ProductName,
				singleItem[i].AuthorName,
				singleItem[i].MakerName,
				singleItem[i].SellingDate,
				singleItem[i].SalesTaxExcUnitPrice,
				singleItem[i].CumulativeReceivingQuantity,
				singleItem[i].CumulativeSalesQuantity,
				singleItem[i].StockQuantity,
				singleItem[i].FirstSalesDate,
				singleItem[i].SalesBodyQuantity,
			)
			if len(exCols) > 0 {
				for _, v := range singleItem[i].SalesBodyQuantityShop {
					singleRowData = append(singleRowData, strconv.Itoa(int(v)))
				}
			}
			itemTemp.Data = append(itemTemp.Data, singleRowData)
			index++
		}
		// Add last SingleItem to data.Rows
		if len(singleItem) > 0 {
			value := ctx.Session.Values[_REPORT_ID + randStringFromSQL + shopCodeKey]
			itemTemp.DataCount = value.(int64)
			data.Rows = append(data.Rows, itemTemp)
		}
	}

	// sort data by input
	if len(form.JanArrays) > 0 {
		sortDataByInput(&data, form.JanArrays)
	}
	//rankNo := data.ShowLineFrom
	//for {
	//	row := <-dataChan
	//	if row == nil {
	//		break
	//	}
	//
	//	rankNo++
	//	dataHaveRankNo := []interface{}{}
	//	dataHaveRankNo = append(dataHaveRankNo, strconv.Itoa(rankNo))
	//	for _, v := range row.ValueIF {
	//		dataHaveRankNo = append(dataHaveRankNo, v)
	//	}
	//	data.Rows = append(data.Rows, dataHaveRankNo)
	//}
	//======================================================
	return &data, nil
}
