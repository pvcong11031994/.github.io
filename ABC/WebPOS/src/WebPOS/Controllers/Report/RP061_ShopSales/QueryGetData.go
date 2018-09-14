package RP061_ShopSales

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
)

func queryData(ctx *gf.Context, sql string, form QueryForm, randStringFromSQL string) (*RpData, error) {

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
	sqlCount, listRange, groupType := buildSql(form, ctx, true)
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

	shopCodeKey := ""
	item := SingleItem{}
	shm := Models.ShopMasterModel{ctx.DB}
	index := 0
	for {
		row := <-dataChan
		if row == nil {
			if item.ShmShopCode != "" {
				item.DataCount = countShop[shopCodeKey]
				data.Rows = append(data.Rows, item)
			}
			break
		}

		shopCd := row.ValueMap["shop_code"].String()
		if index == 0 {
			shopCodeKey = shopCd
		}
		if shopCd != shopCodeKey {
			item.DataCount = countShop[shopCodeKey]
			data.Rows = append(data.Rows, item)
			item = SingleItem{}
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
		singleRowData := []interface{}{}
		singleRowData = append(
			singleRowData,
			item.ShmShopNameShort,
			row.ValueMap["rank"].String(),
			row.ValueMap["jan_code"].String(),
			row.ValueMap["product_name"].String(),
			row.ValueMap["author_name"].String(),
			row.ValueMap["maker_name"].String(),
			row.ValueMap["selling_date"].String(),
			row.ValueMap["sales_tax_exc_unit_price"].Float(),
			row.ValueMap["stok_cumulative_receiving_quantity"].Int(),
			row.ValueMap["stok_cumulative_sales_quantity"].Int(),
			row.ValueMap["stok_stock_quantity"].Int(),
			row.ValueMap["first_sales_date"].String(),
			row.ValueMap["sales_body_quantity"].Int(),
		)
		if len(listRange) > 0 {
			for _, item := range listRange {
				switch groupType {
				case GROUP_TYPE_DATE:
					singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcmm+item.Mcdd].Int())
				case GROUP_TYPE_WEEK:
					singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcweeknum].Int())
				case GROUP_TYPE_MONTH:
					singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcmm].Int())
				}
			}
		}
		item.Data = append(item.Data, singleRowData)
		index++
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
func queryGetDataWithJobId(ctx *gf.Context, sql string, form QueryForm, randStringFromSQL string) (*RpData, error) {

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
		newData, err := queryData(ctx, sql, form, randStringFromSQL)
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
	_, listRange, groupType := buildSql(form, ctx, false)
	// Get data
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChan, err := conn.GetResponseData(jobId.(string), data.ShowLineFrom, limitLengthData, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return &data, exterror.WrapExtError(err)
	}

	shopCodeKey := ""
	item := SingleItem{}
	shm := Models.ShopMasterModel{ctx.DB}
	index := 0
	for {
		row := <-dataChan
		if row == nil {
			if item.ShmShopCode != "" {
				value := ctx.Session.Values[_REPORT_ID+randStringFromSQL+shopCodeKey]
				item.DataCount = value.(int64)
				data.Rows = append(data.Rows, item)
			}
			break
		}

		shopCd := row.ValueMap["shop_code"].String()
		if index == 0 {
			shopCodeKey = shopCd
		}
		if shopCd != shopCodeKey {
			value := ctx.Session.Values[_REPORT_ID+randStringFromSQL+shopCodeKey]
			item.DataCount = value.(int64)
			data.Rows = append(data.Rows, item)
			item = SingleItem{}
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
		singleRowData := []interface{}{}
		singleRowData = append(
			singleRowData,
			item.ShmShopNameShort,
			row.ValueMap["rank"].String(),
			row.ValueMap["jan_code"].String(),
			row.ValueMap["product_name"].String(),
			row.ValueMap["author_name"].String(),
			row.ValueMap["maker_name"].String(),
			row.ValueMap["selling_date"].String(),
			row.ValueMap["sales_tax_exc_unit_price"].Float(),
			row.ValueMap["stok_cumulative_receiving_quantity"].Int(),
			row.ValueMap["stok_cumulative_sales_quantity"].Int(),
			row.ValueMap["stok_stock_quantity"].Int(),
			row.ValueMap["first_sales_date"].String(),
			row.ValueMap["sales_body_quantity"].Int(),
		)
		if len(listRange) > 0 {
			for _, item := range listRange {
				switch groupType {
				case GROUP_TYPE_DATE:
					singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcmm+item.Mcdd].Int())
				case GROUP_TYPE_WEEK:
					singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcweeknum].Int())
				case GROUP_TYPE_MONTH:
					singleRowData = append(singleRowData, row.ValueMap["A"+item.Mcyyyy+item.Mcmm].Int())
				}
			}
		}
		item.Data = append(item.Data, singleRowData)
		index++
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
