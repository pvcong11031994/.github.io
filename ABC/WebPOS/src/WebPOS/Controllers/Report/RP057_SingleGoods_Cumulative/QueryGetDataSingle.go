package RP057_SingleGoods_Cumulative

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

func queryDataSingle(ctx *gf.Context, sql string, listRange []ModelItems.MasterCalendarItem, form QueryFormSingleGoods) ([]SingleItem, SingleItem, error) {

	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)

	singleItem := []SingleItem{}
	totalSingleItem := SingleItem{
		SaleDay:              make(map[string]int64),
		ReturnDay:            make(map[string]int64),
		ReceivingQuantityDay: make(map[string]int64),
		SalesQuantityDay:     make(map[string]int64),
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
	if len(form.JAN) > 0 {
		tag = tag + ",JAN = " + `"` + form.JAN + `"`
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
	tag = tag + `,app_id="mBAWEB-v02a"`
	// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END
	// ========================================================================================
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, tag, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return nil, totalSingleItem, exterror.WrapExtError(err)
	}

	data := RPComon.ReportData{
		ListColKey: []string{},
		ListRowKey: []string{},

		Cols: map[string][]string{},
		Rows: map[string][]string{},
		Data: map[string]map[string][]interface{}{},
	}
	//dataShopName := ResultData{}
	// set VJ_charging current search
	if reportVJCharging, ok := ctx.Get(RPComon.REPORT_VJ_CHARGING_KEY); ok {
		data.VJCharging = reportVJCharging.(int)
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

	// ASO-5540 [BA]mBAWEB-v02a ：単品推移　single_goods
	listJan := make(map[string]string)
	listShop := make(map[string]string)
	count := 0
	for {
		row := <-dataChan
		if row == nil {
			break
		}

		// ASO-5540 [BA]mBAWEB-v02a ：単品推移　single_goods
		// Get list JAN and shop from result query in bq_sales
		listJan[row.ValueMap["jan_code"].String()] = row.ValueMap["jan_code"].String()
		listShop[row.ValueMap["shop_code"].String()] = row.ValueMap["shop_code"].String()

		rs := SingleItem{
			SaleDay:              make(map[string]int64),
			ReturnDay:            make(map[string]int64),
			ReceivingQuantityDay: make(map[string]int64),
			SalesQuantityDay:     make(map[string]int64),
		}
		rs.ShopCd = row.ValueMap["shop_code"].String()
		rs.ShopName = row.ValueMap["shop_name"].String()
		rs.JanCd = row.ValueMap["jan_code"].String()
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		//rs.GoodsName = row.ValueMap["product_name"].String()
		//rs.AuthorName = row.ValueMap["author_name"].String()
		//rs.PublisherName = row.ValueMap["maker_name"].String()
		//rs.SaleDate = row.ValueMap["selling_date"].String()
		//rs.Price = row.ValueMap["sales_tax_exc_unit_price"].Int()
		// ASO-5540 [BA]mBAWEB-v02a ：単品推移　single_goods()
		//rs.FirstSaleDate = row.ValueMap["first_sales_date"].String()
		//rs.ReturnTotal = row.ValueMap["stok_cumulative_receiving_quantity"].Int()
		//rs.SaleTotal = row.ValueMap["stok_cumulative_sales_quantity"].Int()
		// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5597#comment-3267291
		//rs.ShopSeqNumber = row.ValueMap["shop_seq_number"].String()
		rs.SharedBookStoreCode = row.ValueMap["shared_book_store_code"].String()

		rs.SaleTotalDate = row.ValueMap["sales_body_quantity"].Int()
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		rs.ReturnTotalDate = row.ValueMap["receiving_body_quantity"].Int()
		if count == 0 {
			totalSingleItem.SaleTotalDate = row.ValueMap["sales_body_quantity"].Int()
			// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
			totalSingleItem.ReturnTotalDate = row.ValueMap["receiving_body_quantity"].Int()
			// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
			// ASO-5540 [BA]mBAWEB-v02a ：単品推移　single_goods
			//totalSingleItem.StockCurCount = row.ValueMap["stok_stock_quantity"].Int()
			//totalSingleItem.SaleTotal = row.ValueMap["stok_cumulative_sales_quantity"].Int()
			//totalSingleItem.ReturnTotal = row.ValueMap["stok_cumulative_receiving_quantity"].Int()
		} else {
			totalSingleItem.SaleTotalDate = totalSingleItem.SaleTotalDate + row.ValueMap["sales_body_quantity"].Int()
			// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
			totalSingleItem.ReturnTotalDate = totalSingleItem.ReturnTotalDate + row.ValueMap["receiving_body_quantity"].Int()
			// ASO-5540 [BA]mBAWEB-v02a ：単品推移　single_goods
			//totalSingleItem.StockCurCount = totalSingleItem.StockCurCount + row.ValueMap["stok_stock_quantity"].Int()
			//totalSingleItem.SaleTotal = totalSingleItem.SaleTotal + row.ValueMap["stok_cumulative_sales_quantity"].Int()
			//totalSingleItem.ReturnTotal = totalSingleItem.ReturnTotal + row.ValueMap["stok_cumulative_receiving_quantity"].Int()
		}
		//rs.StockCurCount = row.ValueMap["stok_stock_quantity"].Int()
		var maxReceive int64 = 0
		var maxSales int64 = 0
		for _, item := range listRange {
			key := item.McKey

			rs.SaleDay[key] = row.ValueMap["A"+key].Int()
			rs.ReturnDay[key] = row.ValueMap["B"+key].Int()
			rs.ReceivingQuantityDay[key] = row.ValueMap["C"+key].Int()
			rs.SalesQuantityDay[key] = row.ValueMap["D"+key].Int()

			valueReceiving := row.ValueMap["C"+key].Int()
			if valueReceiving > maxReceive {
				maxReceive = valueReceiving
			} else {
				valueReceiving = maxReceive
			}

			valueSales := row.ValueMap["D"+key].Int()
			if valueSales > maxSales {
				maxSales = valueSales
			} else {
				valueSales = maxSales
			}

			if count == 0 {
				totalSingleItem.SaleDay[key] = row.ValueMap["A"+key].Int()
				totalSingleItem.ReturnDay[key] = row.ValueMap["B"+key].Int()
				totalSingleItem.ReceivingQuantityDay[key] = valueReceiving
				totalSingleItem.SalesQuantityDay[key] = valueSales
			} else {
				totalSingleItem.SaleDay[key] += row.ValueMap["A"+key].Int()
				totalSingleItem.ReturnDay[key] += row.ValueMap["B"+key].Int()

				totalSingleItem.ReceivingQuantityDay[key] += valueReceiving
				totalSingleItem.SalesQuantityDay[key] += valueSales
			}
		}

		singleItem = append(singleItem, rs)
		count++

	}
	// ASO-5540 [BA]mBAWEB-v02a ：単品推移　single_goods
	if len(listJan) > 0 {
		// ASO-5559 [BA]mBAWEB-v02a ：単品推移　single_goods（CLOUD対応）
		// Get 在庫数, 累計入荷数, 累計売上数, 初回売上日付 from m_stock of CloudSQL
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
		// ASO-5559 [BA]mBAWEB-v02a ：単品推移　single_goods（CLOUD対応）
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

		// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5559#comment-3260016
		listMStockAll := make(map[string]ModelItems.MStockItem)
		sqlString = `
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
		args = nil
		for _, s := range form.ShopCd {
			args = append(args, s)
		}
		for _, s := range listJan {
			args = append(args, s)
		}
		// ASO-5559 [BA]mBAWEB-v02a ：単品推移　single_goods（CLOUD対応）
		// ASO-5618 [BA]単品推移：セッションが残り続ける
		//dbCloud, err = CloudSQL.Connect()
		//if err != nil {
		//	return nil, totalSingleItem, exterror.WrapExtError(err)
		//}
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT START
		//rows, err = dbCloud.Query(sqlString, args...)
		rows, err = ctx.DB.Query(sqlString, args...)
		// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END
		if err != nil {
			return nil, totalSingleItem, exterror.WrapExtError(err)
		}
		// ASO-5618 [BA]単品推移：セッションが残り続ける
		//defer func() {
		//	rows.Close()
		//	dbCloud.Close()
		//}()
		for rows.Next() {
			newMStockItem := ModelItems.MStockItem{}
			err = db.SqlScanStruct(rows, &newMStockItem)
			if err != nil {
				return nil, totalSingleItem, exterror.WrapExtError(err)
			}
			// Save result query from CloudSQL to listMStock with key is JAN
			listMStockAll[newMStockItem.JanCode] = newMStockItem
		}

		// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5559#comment-3260016
		//minFirstSalesDate := ""
		// Update info get from CloudSQL to result from bq_sales
		for i, _ := range singleItem {
			if mStockData, ok := listMStock[singleItem[i].ShopCd+singleItem[i].JanCd]; ok {
				// 在庫数
				singleItem[i].StockCurCount = mStockData.StockQuantity
				// 累計売上数
				singleItem[i].SaleTotal = mStockData.CumulativeSalesQuantity
				// 累計入荷数
				singleItem[i].ReturnTotal = mStockData.CumulativeReceivingQuantity
				// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5559#comment-3260016
				//// 初回売上日付
				//minFirstSalesDate = getMinItemStr(minFirstSalesDate, mStockData.FirstSalesDate)
				//
				// 在庫数
				totalSingleItem.StockCurCount += mStockData.StockQuantity
				// 累計売上数
				totalSingleItem.SaleTotal += mStockData.CumulativeSalesQuantity
				// 累計入荷数
				totalSingleItem.ReturnTotal += mStockData.CumulativeReceivingQuantity
			}
		}
		if len(singleItem) > 0 {
			// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5559#comment-3260016
			//singleItem[0].FirstSaleDate = minFirstSalesDate
			rs := SingleItem{}

			// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
			sqlString = `
SELECT
	jan_code,
	MAX(product_name) product_name,
	MAX(author_name) author_name,
	MAX(maker_name) maker_name,
	MAX(selling_date) selling_date,
	MAX(list_price) list_price,
	CONCAT(MAX(stock_inf_category), ' ', MAX(stock_inf), ' (', MAX(stock_inf_update_datetime), '時点)') stock_inf
FROM
	m_jan
WHERE
	jan_code IN (?` + strings.Repeat(",?", len(listJan)-1) + `)
GROUP BY
	jan_code
`
			args = nil
			for _, s := range listJan {
				args = append(args, s)
			}
			// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT START
			//rows, err = dbCloud.Query(sqlString, args...)
			rows, err = ctx.DB.Query(sqlString, args...)
			// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - EDIT END
			if err != nil {
				return nil, totalSingleItem, exterror.WrapExtError(err)
			}
			for rows.Next() {
				newMJanItem := ModelItems.MJanItem{}
				err = db.SqlScanStruct(rows, &newMJanItem)
				if err != nil {
					return nil, totalSingleItem, exterror.WrapExtError(err)
				}
				rs.GoodsName = newMJanItem.ProductName
				rs.AuthorName = newMJanItem.AuthorName
				rs.PublisherName = newMJanItem.MakerName
				rs.SaleDate = newMJanItem.SellingDate
				rs.Price = newMJanItem.ListPrice
				// ASO-5651 [BA]mBAWEB-v02a 単品推移：商品情報に出版社在庫の表示を追加
				rs.StockInf = newMJanItem.StockInf
			}
			rs.JanCd = singleItem[0].JanCd
			//rs.GoodsName = singleItem[0].GoodsName
			//rs.AuthorName = singleItem[0].AuthorName
			//rs.PublisherName = singleItem[0].PublisherName
			//rs.SaleDate = singleItem[0].SaleDate
			//rs.Price = singleItem[0].Price
			if mStockData, ok := listMStockAll[rs.JanCd]; ok {
				// 累計入荷数
				rs.ReturnTotal = mStockData.CumulativeReceivingQuantity
				// 累計売上数
				rs.SaleTotal = mStockData.CumulativeSalesQuantity
				// 在庫数
				rs.StockCurCount = mStockData.StockQuantity
				// 初回売上日付
				rs.FirstSaleDate = mStockData.FirstSalesDate
			}
			singleItem = append(singleItem, rs)
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
