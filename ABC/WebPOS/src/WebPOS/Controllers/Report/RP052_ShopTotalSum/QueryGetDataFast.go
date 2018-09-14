package RP052_ShopTotalSum

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/WebApp"
	"errors"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"strings"
)

func queryDataFastTable(ctx *gf.Context, sql string, data *DataSum, form QueryForm) error {

	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)

	keyErr := errors.New("KEY_ERR")
	msgRetryTmp := strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	retryMap := make(map[string]string)
	retryMap["CONFIG_RETRY_COUNT"] = WebApp.CONFIG_RETRY_COUNT
	retryMap["CONFIG_TIME_RETRY"] = WebApp.CONFIG_TIME_RETRY
	retryMap["CONFIG_LIST_CODE_HTTP"] = WebApp.CONFIG_LIST_CODE_HTTP
	conn, err := bq.NewConnection(keyFile, mailAccount, projectId, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return exterror.WrapExtError(err)
	}

	// set report name to import info log search charging
	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
	//totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, _REPORT_ID)
	tag := "report=" + _REPORT_ID
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV_TEXT + `"`
	} else {
		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
	}
	if len(form.MakerCd) > 0 {
		tag = tag + ",出版社 LIKE (" + Common.JoinArray(form.MakerCd, "%", "%", ",") + ")"
	}
	if len(form.MediaGroup1Cd) > 0 {
		tag = tag + ",メディア大分類コード IN (" + strings.Join(form.MediaGroup1Cd, ",") + ")"
	}
	if len(form.MediaGroup2Cd) > 0 {
		tag = tag + ",メディア中分類コード IN (" + strings.Join(form.MediaGroup2Cd, ",") + ")"
	}
	if len(form.MediaGroup3Cd) > 0 {
		tag = tag + ",メディア中小分類コード IN (" + strings.Join(form.MediaGroup3Cd, ",") + ")"
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
	if form.JAN != "" {
		tag = tag + ",JANコード LIKE " + `"` + form.JAN + `%"`
	}
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, tag, ctx, msgRetryTmp, retryMap)

	ctx.Session.Values[RPComon.REPORT_QUERY_JOB_ID] = jobId
	ctx.Session.Values[RPComon.REPORT_QUERY_JOB_ID_COUNT] = totalRows
	if err != nil {
		return exterror.WrapExtError(err)
	}
	data.CountResultRows = totalRows
	if totalRows == 0 {
		return nil
	}

	resultData := []ResultData{}

	// set VJ_charging current search
	if reportVJCharging, ok := ctx.Get(RPComon.REPORT_VJ_CHARGING_KEY); ok {
		data.VJCharging = reportVJCharging.(int)
	}

	if totalRows > RPComon.BQ_DATA_LIMIT {
		return exterror.WrapExtError(errors.New("Respone data too large"))
	}

	data.PageCount = totalRows / form.Limit
	if totalRows%form.Limit > 0 {
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
	data.ShowLineTo = data.ThisPage * form.Limit
	if data.ShowLineTo > totalRows {
		data.ShowLineTo = totalRows
	}

	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		data.ShowLineFrom = 0
		data.ShowLineTo = RPComon.BQ_DATA_LIMIT
	}
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChan, err := conn.GetResponseData(jobId, int(data.ShowLineFrom), int(data.ShowLineTo), ctx, msgRetryTmp, retryMap)

	if err != nil {
		return exterror.WrapExtError(err)
	}

	for {
		row := <-dataChan
		if row == nil {
			break
		}
		dataN := ResultData{}

		dataN.ShopCd = row.ValueMap["shm_shop_cd"].String()
		dataN.ShopName = row.ValueMap["shm_shop_name"].String()
		dataN.JanCd = row.ValueMap["bqio_jan_cd"].String()
		dataN.GoodsName = row.ValueMap["goods_name"].String()
		dataN.PublisherName = row.ValueMap["publisher_name"].String()
		dataN.Price = row.ValueMap["bqgm_price"].Int()
		dataN.StockCount = row.ValueMap["stock_count"].Int()
		dataN.SaleTotal = row.ValueMap["total_sales"].Int()
		dataN.SaleTotalDate = row.ValueMap["total_sales_date"].Int()

		resultData = append(resultData, dataN)
	}
	//======================================================
	data.ResultData = resultData

	return nil
}

func queryDataFastTableWithJobId(ctx *gf.Context, sql string, data *DataSum, form QueryForm) error {

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
		return exterror.WrapExtError(err)
	}

	// set report name to import info log search charging
	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
	jobIdT, _ := ctx.Session.Values[RPComon.REPORT_QUERY_JOB_ID]
	totalRowsT, _ := ctx.Session.Values[RPComon.REPORT_QUERY_JOB_ID_COUNT]
	jobId := jobIdT.(string)
	totalRows := totalRowsT.(int64)

	if err != nil {
		return exterror.WrapExtError(err)
	}
	data.CountResultRows = totalRows
	if totalRows == 0 {
		return nil
	}

	resultData := []ResultData{}

	// set VJ_charging current search
	if reportVJCharging, ok := ctx.Get(RPComon.REPORT_VJ_CHARGING_KEY); ok {
		data.VJCharging = reportVJCharging.(int)
	}

	if totalRows > RPComon.BQ_DATA_LIMIT {
		return exterror.WrapExtError(errors.New("Respone data too large"))
	}

	data.PageCount = totalRows / form.Limit
	if totalRows%form.Limit > 0 {
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
	data.ShowLineTo = data.ThisPage * form.Limit
	if data.ShowLineTo > totalRows {
		data.ShowLineTo = totalRows
	}

	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		data.ShowLineFrom = 0
		data.ShowLineTo = RPComon.BQ_DATA_LIMIT
	}
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChan, err := conn.GetResponseData(jobId, int(data.ShowLineFrom), int(data.ShowLineTo), ctx, msgRetryTmp, retryMap)

	if err != nil {
		return exterror.WrapExtError(err)
	}

	for {
		row := <-dataChan
		if row == nil {
			break
		}
		dataN := ResultData{}

		dataN.ShopCd = row.ValueMap["shm_shop_cd"].String()
		dataN.ShopName = row.ValueMap["shm_shop_name"].String()
		dataN.JanCd = row.ValueMap["bqio_jan_cd"].String()
		dataN.GoodsName = row.ValueMap["goods_name"].String()
		dataN.PublisherName = row.ValueMap["publisher_name"].String()
		dataN.Price = row.ValueMap["bqgm_price"].Int()
		dataN.StockCount = row.ValueMap["stock_count"].Int()
		dataN.SaleTotal = row.ValueMap["total_sales"].Int()
		dataN.SaleTotalDate = row.ValueMap["total_sales_date"].Int()

		resultData = append(resultData, dataN)
	}
	//======================================================
	data.ResultData = resultData

	return nil
}
