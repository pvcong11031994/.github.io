package RP058_SalesComparison

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"errors"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
	"strings"
)

func queryData(ctx *gf.Context, sql string, listRange []ModelItems.MasterCalendarItem, form QueryForm) ([]SingleItem, SingleItem, error) {

	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)

	singleItem := []SingleItem{}
	totalSingleItem := SingleItem{
		SaleDay: make(map[string]map[string]int64),
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
		return nil, totalSingleItem, exterror.WrapExtError(err)
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
	// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
	tag = tag + `,app_id="mBAWEB-v10a"`
	// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END
	// ========================================================================================
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, tag, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return nil, totalSingleItem, exterror.WrapExtError(err)
	}

	if totalRows > RPComon.BQ_DATA_LIMIT {
		return nil, totalSingleItem, exterror.WrapExtError(errors.New("Respone data too large"))
	}
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChan, err := conn.GetResponseData(jobId, 0, RPComon.BQ_DATA_LIMIT, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return nil, totalSingleItem, exterror.WrapExtError(err)
	}

	// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
	listJan := make(map[string]string)
	// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5560#comment-3260201
	//listShop := make(map[string]string)
	strJanCode := ""
	rs := SingleItem{
		SaleDay: make(map[string]map[string]int64),
	}
	listShopOfJAN := make(map[string]string)
	count := 0
	for {
		row := <-dataChan
		if row == nil {
			break
		}
		// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
		// Get list JAN and shop from result query in bq_sales
		listJan[row.ValueMap["jan_code"].String()] = row.ValueMap["jan_code"].String()
		// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5560#comment-3260201
		//listShop[row.ValueMap["shop_code"].String()] = row.ValueMap["shop_code"].String()

		// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
		if strJanCode == "" {
			strJanCode = row.ValueMap["jan_code"].String()
		}
		if strings.Compare(strJanCode, row.ValueMap["jan_code"].String()) != 0 {
			rs.ListShopCd = listShopOfJAN
			singleItem = append(singleItem, rs)
			rs = SingleItem{
				SaleDay: make(map[string]map[string]int64),
			}
			strJanCode = row.ValueMap["jan_code"].String()
			listShopOfJAN = make(map[string]string)
			// ASO-5829 [BA]mBAWEB-v10a 売上比較：バグ修正＞最小のJANについてグラフに1店舗分の売上数しか表示されない - EDIT START
			//count++
			count = 0
			// ASO-5829 [BA]mBAWEB-v10a 売上比較：バグ修正＞最小のJANについてグラフに1店舗分の売上数しか表示されない - EDIT END
		}
		//rs := SingleItem{
		//	SaleDay: make(map[string]map[string]int64),
		//}
		listShopOfJAN[row.ValueMap["shop_code"].String()] = row.ValueMap["shop_code"].String()
		rs.JanCd = row.ValueMap["jan_code"].String()
		// ASO-5836 [BA]mBAWEB-v10a 売上比較：商品名などをCloudSQLから取得 - DEL START
		//rs.GoodsName = row.ValueMap["product_name"].String()
		//rs.AuthorName = row.ValueMap["author_name"].String()
		//rs.PublisherName = row.ValueMap["maker_name"].String()
		//rs.SaleDate = row.ValueMap["selling_date"].String()
		//rs.Price = row.ValueMap["sales_tax_exc_unit_price"].Int()
		// ASO-5836 [BA]mBAWEB-v10a 売上比較：商品名などをCloudSQLから取得 - DEL START
		// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
		//rs.FirstSaleDate = row.ValueMap["first_sales_date"].String()
		//rs.ReturnTotal = row.ValueMap["stok_cumulative_receiving_quantity"].Int()
		//rs.SaleTotal = row.ValueMap["stok_cumulative_sales_quantity"].Int()
		//rs.SaleTotalDate = row.ValueMap["sales_body_quantity"].Int()
		rs.SaleTotalDate += row.ValueMap["sales_body_quantity"].Int()
		if count == 0 {
			totalSingleItem.SaleTotalDate = row.ValueMap["sales_body_quantity"].Int()
			// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
			//totalSingleItem.StockCurCount = row.ValueMap["stok_stock_quantity"].Int()
			//totalSingleItem.SaleTotal = row.ValueMap["stok_cumulative_sales_quantity"].Int()
			//totalSingleItem.ReturnTotal = row.ValueMap["stok_cumulative_receiving_quantity"].Int()
		} else {
			totalSingleItem.SaleTotalDate = totalSingleItem.SaleTotalDate + row.ValueMap["sales_body_quantity"].Int()
			// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
			//totalSingleItem.StockCurCount = totalSingleItem.StockCurCount + row.ValueMap["stok_stock_quantity"].Int()
			//totalSingleItem.SaleTotal = totalSingleItem.SaleTotal + row.ValueMap["stok_cumulative_sales_quantity"].Int()
			//totalSingleItem.ReturnTotal = totalSingleItem.ReturnTotal + row.ValueMap["stok_cumulative_receiving_quantity"].Int()
		}
		// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
		//rs.StockCurCount = row.ValueMap["stok_stock_quantity"].Int()
		if totalSingleItem.SaleDay[rs.JanCd] == nil {
			totalSingleItem.SaleDay[rs.JanCd] = make(map[string]int64)
		}
		if rs.SaleDay[rs.JanCd] == nil {
			rs.SaleDay[rs.JanCd] = make(map[string]int64)
		}
		for _, item := range listRange {

			keySales := item.McKey
			rs.SaleDay[rs.JanCd][keySales] = row.ValueMap["A"+keySales].Int()
			if count == 0 {
				totalSingleItem.SaleDay[rs.JanCd][keySales] = row.ValueMap["A"+keySales].Int()
			} else {
				totalSingleItem.SaleDay[rs.JanCd][keySales] = totalSingleItem.SaleDay[rs.JanCd][keySales] + row.ValueMap["A"+keySales].Int()
			}

		}
		// ASO-5829 [BA]mBAWEB-v10a 売上比較：バグ修正＞最小のJANについてグラフに1店舗分の売上数しか表示されない - ADD START
		count++
		// ASO-5829 [BA]mBAWEB-v10a 売上比較：バグ修正＞最小のJANについてグラフに1店舗分の売上数しか表示されない - ADD END
		// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
		//singleItem = append(singleItem, rs)
		//count++
	}
	// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
	rs.ListShopCd = listShopOfJAN
	singleItem = append(singleItem, rs)
	if len(listJan) > 0 {
		// Get 在庫数, 累計入荷数, 累計売上数, 初回売上日付 from m_stock of MariaDB
		listMStock := make(map[string]ModelItems.MStockItem)
		sqlString := `
SELECT
	jan_code,
	SUM(stock_quantity) stock_quantity,
	SUM(cumulative_receiving_quantity) cumulative_receiving_quantity,
	SUM(cumulative_sales_quantity) cumulative_sales_quantity,
	MIN(CASE WHEN TRIM(first_sales_date) <> '' THEN first_sales_date END) first_sales_date
FROM
	m_stock
WHERE
	shop_code IN (?` + strings.Repeat(",?", len(form.ShopCd)-1) + `)
	AND jan_code IN (?` + strings.Repeat(",?", len(listJan)-1) + `)
GROUP BY
	jan_code
`
		var args []interface{}
		for _, s := range form.ShopCd {
			args = append(args, s)
		}
		for _, s := range listJan {
			args = append(args, s)
		}
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT START
		//// ASO-5560 [BA]mBAWEB-v10a ：売上比較　sales_comparison（CLOUD対応）
		//dbCloud, err := CloudSQL.Connect()
		//if err != nil {
		//	return nil, totalSingleItem, exterror.WrapExtError(err)
		//}
		//rows, err := dbCloud.Query(sqlString, args...)
		//if err != nil {
		//	return nil, totalSingleItem, exterror.WrapExtError(err)
		//}
		//defer func() {
		//	rows.Close()
		//	dbCloud.Close()
		//}()
		rows, err := ctx.DB.Query(sqlString, args...)
		if err != nil {
			return nil, totalSingleItem, exterror.WrapExtError(err)
		}
		defer func() {
			rows.Close()
		}()
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END
		for rows.Next() {
			newMStockItem := ModelItems.MStockItem{}
			err = db.SqlScanStruct(rows, &newMStockItem)
			if err != nil {
				return nil, totalSingleItem, exterror.WrapExtError(err)
			}
			// Save result query from CloudSQL to listMStock with key is JAN
			// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5560#comment-3260201
			//listMStock[newMStockItem.ShopCode+newMStockItem.JanCode] = newMStockItem
			listMStock[newMStockItem.JanCode] = newMStockItem
		}

		// ASO-5836 [BA]mBAWEB-v10a 売上比較：商品名などをCloudSQLから取得 - ADD START
		listMJan := make(map[string]ModelItems.MJanItem)
		sqlStringMJan := `
SELECT
	jan_code,
	product_name,
	author_name,
	maker_name,
	selling_date,
	list_price
FROM
	m_jan
WHERE
	jan_code IN (?` + strings.Repeat(",?", len(listJan)-1) + `)
`
		var argsMJan []interface{}
		for _, s := range listJan {
			argsMJan = append(argsMJan, s)
		}
		rows, err = ctx.DB.Query(sqlStringMJan, argsMJan...)
		if err != nil {
			return nil, totalSingleItem, exterror.WrapExtError(err)
		}
		for rows.Next() {
			newMJanItem := ModelItems.MJanItem{}
			err = db.SqlScanStruct(rows, &newMJanItem)
			if err != nil {
				return nil, totalSingleItem, exterror.WrapExtError(err)
			}
			listMJan[newMJanItem.JanCode] = newMJanItem
		}
		// ASO-5836 [BA]mBAWEB-v10a 売上比較：商品名などをCloudSQLから取得 - ADD END

		// Update info get from CloudSQL to result from bq_sales
		for i, _ := range singleItem {
			// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5560#comment-3260201
			//for _, shop := range singleItem[i].ListShopCd {
			//	if mStockData, ok := listMStock[shop+singleItem[i].JanCd]; ok {
			if mStockData, ok := listMStock[singleItem[i].JanCd]; ok {
				// 在庫数
				singleItem[i].StockCurCount += mStockData.StockQuantity
				// 累計売上数
				singleItem[i].SaleTotal += mStockData.CumulativeSalesQuantity
				// 累計入荷数
				singleItem[i].ReturnTotal += mStockData.CumulativeReceivingQuantity
				// 初回売上日付
				singleItem[i].FirstSaleDate = getMinItemStr(singleItem[i].FirstSaleDate, mStockData.FirstSalesDate)

				// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5560#comment-3260201
				//// 在庫数
				//totalSingleItem.StockCurCount += mStockData.StockQuantity
				//// 累計売上数
				//totalSingleItem.SaleTotal += mStockData.CumulativeSalesQuantity
				//// 累計入荷数
				//totalSingleItem.ReturnTotal += mStockData.CumulativeReceivingQuantity
			}
			// ASO-5836 [BA]mBAWEB-v10a 売上比較：商品名などをCloudSQLから取得 - ADD START
			if mJanData, ok := listMJan[singleItem[i].JanCd]; ok {
				singleItem[i].GoodsName = mJanData.ProductName
				singleItem[i].AuthorName = mJanData.AuthorName
				singleItem[i].PublisherName = mJanData.MakerName
				singleItem[i].SaleDate = mJanData.SellingDate
				singleItem[i].Price = int64(mJanData.ListPrice)
			}
			// ASO-5836 [BA]mBAWEB-v10a 売上比較：商品名などをCloudSQLから取得 - ADD END
			//}
		}
	}
	//======================================================
	return singleItem, totalSingleItem, nil
}

func queryDataDetail(ctx *gf.Context, sql string, listRange []ModelItems.MasterCalendarItem, form QueryForm, jancd string) ([]SingleItem, SingleItem, error) {

	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)

	singleItem := []SingleItem{}
	totalSingleItem := SingleItem{
		SaleDay: make(map[string]map[string]int64),
	}
	totalSingleItem.JanCd = jancd

	keyErr := errors.New("KEY_ERR")
	msgRetryTmp := strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	// ASO-5502 JOBのコンパイルが通らない（stg）
	retryMap := make(map[string]string)
	retryMap["CONFIG_RETRY_COUNT"] = WebApp.CONFIG_RETRY_COUNT
	retryMap["CONFIG_TIME_RETRY"] = WebApp.CONFIG_TIME_RETRY
	retryMap["CONFIG_LIST_CODE_HTTP"] = WebApp.CONFIG_LIST_CODE_HTTP
	conn, err := bq.NewConnection(keyFile, mailAccount, projectId, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return nil, totalSingleItem, exterror.WrapExtError(err)
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
	tag = tag + ",JAN = " + `"` + jancd + `"`
	// ========================================================================================
	keyErr = errors.New("KEY_ERR")
	// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
	tag = tag + `,app_id="mBAWEB-v10a"`
	// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, tag, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return nil, totalSingleItem, exterror.WrapExtError(err)
	}

	if totalRows > RPComon.BQ_DATA_LIMIT {
		return nil, totalSingleItem, exterror.WrapExtError(errors.New("Respone data too large"))
	}
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChan, err := conn.GetResponseData(jobId, 0, RPComon.BQ_DATA_LIMIT, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return nil, totalSingleItem, exterror.WrapExtError(err)
	}

	// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
	listJan := make(map[string]string)
	listShop := make(map[string]string)
	count := 0
	for {
		row := <-dataChan
		if row == nil {
			break
		}
		// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
		// Get list JAN and shop from result query in bq_sales
		listJan[row.ValueMap["jan_code"].String()] = row.ValueMap["jan_code"].String()
		listShop[row.ValueMap["shop_code"].String()] = row.ValueMap["shop_code"].String()

		rs := SingleItem{
			SaleDay: make(map[string]map[string]int64),
		}
		rs.JanCd = row.ValueMap["jan_code"].String()
		rs.ShopCd = row.ValueMap["shop_code"].String()
		rs.ShopName = row.ValueMap["shop_name"].String()
		// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
		//rs.ReturnTotal = row.ValueMap["stok_cumulative_receiving_quantity"].Int()
		//rs.SaleTotal = row.ValueMap["stok_cumulative_sales_quantity"].Int()
		rs.SaleTotalDate = row.ValueMap["sales_body_quantity"].Int()
		if count == 0 {
			totalSingleItem.SaleTotalDate = row.ValueMap["sales_body_quantity"].Int()
			// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
			//totalSingleItem.StockCurCount = row.ValueMap["stok_stock_quantity"].Int()
			//totalSingleItem.SaleTotal = row.ValueMap["stok_cumulative_sales_quantity"].Int()
			//totalSingleItem.ReturnTotal = row.ValueMap["stok_cumulative_receiving_quantity"].Int()
		} else {
			totalSingleItem.SaleTotalDate = totalSingleItem.SaleTotalDate + row.ValueMap["sales_body_quantity"].Int()
			// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
			//totalSingleItem.StockCurCount = totalSingleItem.StockCurCount + row.ValueMap["stok_stock_quantity"].Int()
			//totalSingleItem.SaleTotal = totalSingleItem.SaleTotal + row.ValueMap["stok_cumulative_sales_quantity"].Int()
			//totalSingleItem.ReturnTotal = totalSingleItem.ReturnTotal + row.ValueMap["stok_cumulative_receiving_quantity"].Int()
		}
		// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
		//rs.StockCurCount = row.ValueMap["stok_stock_quantity"].Int()
		if totalSingleItem.SaleDay[rs.JanCd] == nil {
			totalSingleItem.SaleDay[rs.JanCd] = make(map[string]int64)
		}
		if rs.SaleDay[rs.JanCd] == nil {
			rs.SaleDay[rs.JanCd] = make(map[string]int64)
		}
		for _, item := range listRange {
			keySales := item.McKey
			rs.SaleDay[rs.JanCd][keySales] = row.ValueMap["A"+keySales].Int()
			if count == 0 {
				totalSingleItem.SaleDay[rs.JanCd][keySales] = row.ValueMap["A"+keySales].Int()
			} else {
				totalSingleItem.SaleDay[rs.JanCd][keySales] = totalSingleItem.SaleDay[rs.JanCd][keySales] + row.ValueMap["A"+keySales].Int()
			}

		}
		singleItem = append(singleItem, rs)
		count++
	}
	// ASO-5541 [BA]mBAWEB-v10a ：売上比較　sales_comparison
	if len(listJan) > 0 {
		// Get 在庫数, 累計入荷数, 累計売上数, 初回売上日付 from m_stock of MariaDB
		listMStock := make(map[string]ModelItems.MStockItem)
		sqlString := `
SELECT
	shop_code,
	jan_code,
	SUM(stock_quantity) stock_quantity,
	SUM(cumulative_receiving_quantity) cumulative_receiving_quantity,
	SUM(cumulative_sales_quantity) cumulative_sales_quantity,
	MIN(CASE WHEN TRIM(first_sales_date) <> '' THEN first_sales_date END) first_sales_date
FROM
	m_stock
WHERE
	shop_code IN (?` + strings.Repeat(",?", len(listShop)-1) + `)
	AND jan_code IN (?` + strings.Repeat(",?", len(listJan)-1) + `)
GROUP BY
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
		//// ASO-5560 [BA]mBAWEB-v10a ：売上比較　sales_comparison（CLOUD対応）
		//dbCloud, err := CloudSQL.Connect()
		//if err != nil {
		//	return nil, totalSingleItem, exterror.WrapExtError(err)
		//}
		//rows, err := dbCloud.Query(sqlString, args...)
		//if err != nil {
		//	return nil, totalSingleItem, exterror.WrapExtError(err)
		//}
		//defer func() {
		//	rows.Close()
		//	dbCloud.Close()
		//}()
		rows, err := ctx.DB.Query(sqlString, args...)
		if err != nil {
			return nil, totalSingleItem, exterror.WrapExtError(err)
		}
		defer func() {
			rows.Close()
		}()
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END
		for rows.Next() {
			newMStockItem := ModelItems.MStockItem{}
			err = db.SqlScanStruct(rows, &newMStockItem)
			if err != nil {
				return nil, totalSingleItem, exterror.WrapExtError(err)
			}
			// Save result query from CloudSQL to listMStock with key is Shop and JAN
			listMStock[newMStockItem.ShopCode+newMStockItem.JanCode] = newMStockItem
		}
		minFirstSalesDate := ""
		// Update info get from CloudSQL to result from bq_sales
		for i, _ := range singleItem {
			if mStockData, ok := listMStock[singleItem[i].ShopCd+singleItem[i].JanCd]; ok {
				// 在庫数
				singleItem[i].StockCurCount = mStockData.StockQuantity
				// 累計売上数
				singleItem[i].SaleTotal = mStockData.CumulativeSalesQuantity
				// 累計入荷数
				singleItem[i].ReturnTotal = mStockData.CumulativeReceivingQuantity
				// 初回売上日付
				minFirstSalesDate = getMinItemStr(minFirstSalesDate, mStockData.FirstSalesDate)

				// 在庫数
				totalSingleItem.StockCurCount += mStockData.StockQuantity
				// 累計売上数
				totalSingleItem.SaleTotal += mStockData.CumulativeSalesQuantity
				// 累計入荷数
				totalSingleItem.ReturnTotal += mStockData.CumulativeReceivingQuantity
			}
		}
		// 初回売上日付
		if len(singleItem) > 0 {
			singleItem[0].FirstSaleDate = minFirstSalesDate
		}
	}
	//======================================================
	return singleItem, totalSingleItem, nil
}

func getMinItemStr(itemMin, itemCheck string) string {
	if itemMin == "" {
		return itemCheck
	}
	if itemCheck == "" {
		return itemMin
	}
	if itemCheck < itemMin {
		return itemCheck
	}
	return itemMin
}
