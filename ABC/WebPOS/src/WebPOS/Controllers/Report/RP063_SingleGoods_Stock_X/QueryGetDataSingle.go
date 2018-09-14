package RP063_SingleGoods_Stock_X

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"errors"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"strings"
	"time"
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

	count := 0
	for {
		row := <-dataChan
		if row == nil {
			break
		}

		rs := SingleItem{
			SaleDay:              make(map[string]int64),
			ReturnDay:            make(map[string]int64),
			ReceivingQuantityDay: make(map[string]int64),
			SalesQuantityDay:     make(map[string]int64),
		}
		rs.ShopCd = row.ValueMap["shop_code"].String()
		rs.ShopName = row.ValueMap["shop_name"].String()
		rs.JanCd = row.ValueMap["jan_code"].String()
		rs.SaleTotalDate = row.ValueMap["sales_body_quantity"].Int()
		if count == 0 {
			totalSingleItem.SaleTotalDate = row.ValueMap["sales_body_quantity"].Int()
		} else {
			totalSingleItem.SaleTotalDate = totalSingleItem.SaleTotalDate + row.ValueMap["sales_body_quantity"].Int()
		}
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
	//======================================================

	return singleItem, totalSingleItem, nil
}

func queryDataStockAll(ctx *gf.Context, sql string, listRange []ModelItems.MasterCalendarItem, form QueryFormSingleGoods, listSingleItem []SingleItem, totalSingleItem SingleItem) ([]SingleItem, SingleItem, error) {

	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)
	keyErr := errors.New("KEY_ERR")
	msgRetryTmp := strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	// ASO-5502 JOBのコンパイルが通らない（stg）
	retryMap := make(map[string]string)
	retryMap["CONFIG_RETRY_COUNT"] = WebApp.CONFIG_RETRY_COUNT
	retryMap["CONFIG_TIME_RETRY"] = WebApp.CONFIG_TIME_RETRY
	retryMap["CONFIG_LIST_CODE_HTTP"] = WebApp.CONFIG_LIST_CODE_HTTP
	conn, err := bq.NewConnection(keyFile, mailAccount, projectId, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return listSingleItem, totalSingleItem, exterror.WrapExtError(err)
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
		return listSingleItem, totalSingleItem, exterror.WrapExtError(err)
	}

	listStockByShop := map[string]int64{}
	minFirstSalesDate := time.Now().AddDate(0, 0, 1).Format(Common.DATE_FORMAT_YMD)

	index := 0
	for {
		row := <-dataChan
		if row == nil {
			break
		}

		if index == 0 {
			totalSingleItem.ShopCd = row.ValueMap["shop_code"].String()
			totalSingleItem.JanCd = row.ValueMap["jan_code"].String()
			totalSingleItem.GoodsName = row.ValueMap["product_name"].String()
			totalSingleItem.AuthorName = row.ValueMap["author_name"].String()
			totalSingleItem.PublisherName = row.ValueMap["maker_name"].String()
			totalSingleItem.SaleDate = row.ValueMap["selling_date"].String()
			totalSingleItem.Price = int64(row.ValueMap["list_price"].Float())
			totalSingleItem.FirstSaleDate = getMinItemStr(minFirstSalesDate, row.ValueMap["first_sales_date"].String())
		} else {
			totalSingleItem.ShopCd = getMaxItemStr(totalSingleItem.ShopCd, row.ValueMap["shop_code"].String())
			totalSingleItem.JanCd = getMaxItemStr(totalSingleItem.JanCd, row.ValueMap["jan_code"].String())
			totalSingleItem.GoodsName = getMaxItemStr(totalSingleItem.GoodsName, row.ValueMap["product_name"].String())
			totalSingleItem.AuthorName = getMaxItemStr(totalSingleItem.AuthorName, row.ValueMap["author_name"].String())
			totalSingleItem.PublisherName = getMaxItemStr(totalSingleItem.PublisherName, row.ValueMap["maker_name"].String())
			totalSingleItem.SaleDate = getMaxItemStr(totalSingleItem.SaleDate, row.ValueMap["selling_date"].String())
			totalSingleItem.Price = getMaxItemInt64(totalSingleItem.Price, int64(row.ValueMap["list_price"].Int()))
			totalSingleItem.FirstSaleDate = getMinItemStr(totalSingleItem.FirstSaleDate, row.ValueMap["first_sales_date"].String())
		}

		shopCd := row.ValueMap["shop_code"].String()
		totalSingleItem.StockTotal += row.ValueMap["stock_quantity"].Int()
		totalSingleItem.SaleTotal += row.ValueMap["cumulative_sales_quantity"].Int()
		totalSingleItem.ReturnTotal += row.ValueMap["cumulative_receiving_quantity"].Int()
		listStockByShop[shopCd] = row.ValueMap["stock_quantity"].Int()
		index++
	}
	//======================================================

	if totalSingleItem.FirstSaleDate == time.Now().AddDate(0, 0, 1).Format(Common.DATE_FORMAT_YMD) {
		totalSingleItem.FirstSaleDate = ""
	}
	for key, value := range listSingleItem {
		listSingleItem[key].StockCountByShop = listStockByShop[value.ShopCd]
		totalSingleItem.StockCountByShopSearchDate += listStockByShop[value.ShopCd]
	}
	return listSingleItem, totalSingleItem, nil
}

func getMaxItemStr(itemMax, itemCheck string) string {
	if itemCheck > itemMax {
		return itemCheck
	}
	return itemMax
}
func getMaxItemInt64(itemMax, itemCheck int64) int64 {
	if itemCheck > itemMax {
		return itemCheck
	}
	return itemMax
}
func getMinItemStr(itemMin, itemCheck string) string {
	if itemCheck == "" {
		return itemMin
	}
	if itemCheck < itemMin {
		return itemCheck
	}
	return itemMin
}
