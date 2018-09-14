package RP067_ShopSales_Cloud

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"errors"
	"strconv"
	"strings"

	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
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
	// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
	tag = tag + `,app_id="mBAWEB-v25a"`
	// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END
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
	if totalRows == 0 {
		return &data, nil
	}
	//=====================================================================================
	// set report name to import info log search charging
	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
	// Search count Data by shop
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

	listShop := []string{}
	listJan := []string{}
	listedJanMap := map[string]bool{}
	singleItem := []SingleItem{}

	shopCodeKey := ""
	item := SingleItem{}
	shm := Models.ShopMasterModel{ctx.DB}
	index := 0
	for {
		row := <-dataChan
		if row == nil {
			break
		}

		janCd := row.ValueMap["jan_code"].String()
		if !listedJanMap[janCd] {
			listedJanMap[janCd] = true
			listJan = append(listJan, janCd)
		}

		shopCd := row.ValueMap["shop_code"].String()
		if index == 0 {
			shopCodeKey = shopCd
			listShop = append(listShop, shopCd)
		}
		if shopCd != shopCodeKey {
			item = SingleItem{}
			item.DataCount = countShop[shopCodeKey]
			shopCodeKey = shopCd
			listShop = append(listShop, shopCd)
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
			item.SalesBodyQuantityShop = append(item.SalesBodyQuantityShop, row.ValueMap["A"+v].Int())
		}

		singleItem = append(singleItem, item)
		index++
	}

	sqlString := `
SELECT
	shop_code,
	jan_code,
	stock_quantity,
	cumulative_receiving_quantity,
	cumulative_sales_quantity,
	first_sales_date
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
	//	return nil, exterror.WrapExtError(err)
	//}
	//
	//rows, err := dbCloud.Query(sqlString, args...)
	//if err != nil {
	//	return nil, exterror.WrapExtError(err)
	//}
	//
	//defer func() {
	//	rows.Close()
	//	dbCloud.Close()
	//}()
	rows, err := ctx.DB.Query(sqlString, args...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
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
			return nil, exterror.WrapExtError(err)
		}
		mStockDataMap[newMStockItem.ShopCode+newMStockItem.JanCode] = newMStockItem
	}

	itemTemp := SingleItem{}
	index = 0
	for i, _ := range singleItem {
		if mStockData, ok := mStockDataMap[singleItem[i].ShmShopCode+singleItem[i].JanCd]; ok {
			singleItem[i].StockQuantity = mStockData.StockQuantity
			singleItem[i].CumulativeReceivingQuantity = mStockData.CumulativeReceivingQuantity
			singleItem[i].CumulativeSalesQuantity = mStockData.CumulativeSalesQuantity
			singleItem[i].FirstSalesDate = mStockData.FirstSalesDate
		}

		if index == 0 {
			itemTemp = singleItem[i]
			shopCodeKey = singleItem[i].ShmShopCode
		}

		if singleItem[i].ShmShopCode != shopCodeKey {
			itemTemp.DataCount = countShop[shopCodeKey]
			data.Rows = append(data.Rows, itemTemp)
			itemTemp = singleItem[i]
			shopCodeKey = singleItem[i].ShmShopCode
		}

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

	if len(singleItem) > 0 {
		itemTemp.DataCount = countShop[shopCodeKey]
		data.Rows = append(data.Rows, itemTemp)
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
	_, _, _, exCols = buildSql(form, ctx, false)
	// Get data
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChan, err := conn.GetResponseData(jobId.(string), data.ShowLineFrom, limitLengthData, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return &data, exterror.WrapExtError(err)
	}

	listShop := []string{}
	listJan := []string{}
	listedJanMap := map[string]bool{}
	singleItem := []SingleItem{}

	shopCodeKey := ""
	item := SingleItem{}
	shm := Models.ShopMasterModel{ctx.DB}
	index := 0
	for {
		row := <-dataChan
		if row == nil {
			break
		}

		janCd := row.ValueMap["jan_code"].String()
		if !listedJanMap[janCd] {
			listedJanMap[janCd] = true
			listJan = append(listJan, janCd)
		}

		shopCd := row.ValueMap["shop_code"].String()
		if index == 0 {
			shopCodeKey = shopCd
			listShop = append(listShop, shopCd)
		}
		if shopCd != shopCodeKey {
			item = SingleItem{}
			value := ctx.Session.Values[_REPORT_ID+randStringFromSQL+shopCodeKey]
			item.DataCount = value.(int64)
			shopCodeKey = shopCd
			listShop = append(listShop, shopCd)
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
			item.SalesBodyQuantityShop = append(item.SalesBodyQuantityShop, row.ValueMap["A"+v].Int())
		}

		singleItem = append(singleItem, item)
		index++
	}

	sqlString := `
SELECT
	shop_code,
	jan_code,
	stock_quantity,
	cumulative_receiving_quantity,
	cumulative_sales_quantity,
	first_sales_date
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
	//	return nil, exterror.WrapExtError(err)
	//}
	//
	//rows, err := dbCloud.Query(sqlString, args...)
	//if err != nil {
	//	return nil, exterror.WrapExtError(err)
	//}
	//
	//defer func() {
	//	rows.Close()
	//	dbCloud.Close()
	//}()
	rows, err := ctx.DB.Query(sqlString, args...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
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
			return nil, exterror.WrapExtError(err)
		}
		mStockDataMap[newMStockItem.ShopCode+newMStockItem.JanCode] = newMStockItem
	}

	itemTemp := SingleItem{}
	index = 0
	for i, _ := range singleItem {
		if mStockData, ok := mStockDataMap[singleItem[i].ShmShopCode+singleItem[i].JanCd]; ok {
			singleItem[i].StockQuantity = mStockData.StockQuantity
			singleItem[i].CumulativeReceivingQuantity = mStockData.CumulativeReceivingQuantity
			singleItem[i].CumulativeSalesQuantity = mStockData.CumulativeSalesQuantity
			singleItem[i].FirstSalesDate = mStockData.FirstSalesDate
		}

		if index == 0 {
			itemTemp = singleItem[i]
			shopCodeKey = singleItem[i].ShmShopCode
		}

		if singleItem[i].ShmShopCode != shopCodeKey {
			value := ctx.Session.Values[_REPORT_ID+randStringFromSQL+shopCodeKey]
			itemTemp.DataCount = value.(int64)
			data.Rows = append(data.Rows, itemTemp)
			itemTemp = singleItem[i]
			shopCodeKey = singleItem[i].ShmShopCode
		}

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

	if len(singleItem) > 0 {
		value := ctx.Session.Values[_REPORT_ID+randStringFromSQL+shopCodeKey]
		itemTemp.DataCount = value.(int64)
		data.Rows = append(data.Rows, itemTemp)
	}

	// sort data by input
	if len(form.JanArrays) > 0 {
		sortDataByInput(&data, form.JanArrays)
	}

	return &data, nil
}
