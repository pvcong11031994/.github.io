package RP055_BestSales

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/WebApp"
	"errors"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"strconv"
	"strings"
)

//Create data from queryBuild
func queryData(ctx *gf.Context, sql string, form QueryForm, randStringFromSQL string) (*RpData, error) {

	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)

	data := RpData{
		HeaderCols: []string{},
		Cols:       [][]string{},
		Rows:       [][]interface{}{},
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
	if len(form.MediaGroup1Cd) > 0 {
		tag = tag + ",メディア大分類コード IN (" + strings.Join(form.MediaGroup1Cd, ",") + ")"
	}
	if len(form.MediaGroup2Cd) > 0 {
		tag = tag + ",メディア中分類コード IN (" + strings.Join(form.MediaGroup2Cd, ",") + ")"
	}
	if len(form.MediaGroup3Cd) > 0 {
		tag = tag + ",メディア中小分類コード IN (" + strings.Join(form.MediaGroup3Cd, ",") + ")"
	}
	if len(form.MediaGroup4Cd) > 0 {
		tag = tag + ",メディア小分類コード IN (" + strings.Join(form.MediaGroup4Cd, ",") + ")"
	}
	if len(form.MakerCd) > 0 {
		tag = tag + ",出版社コード IN (" + Common.JoinArray(form.MakerCd, "'", "'", ",") + ")"
	}
	if form.ControlType == CONTROL_TYPE_BOOK {
		if len(form.JanMakerCode) > 0 {
			tag = tag + ",出版者記号 IN (" + Common.JoinArray(form.JanMakerCode, JAN_MAKER_CODE, "", ",") + ")"
		}
	}
	if form.ControlType == CONTROL_TYPE_MAGAZINE {
		if len(form.MagazineCd) > 0 {
			tag = tag + ",雑誌コード LIKE (" + Common.JoinArray(form.MagazineCd, "%", "%", ",") + ")"
		}
		if form.MagazineCodeWeek == BQSL_MAGAZINE_CODE_MONTH ||
		form.MagazineCodeMonth == BQSL_MAGAZINE_CODE_WEEK ||
		form.MagazineCodeQuarter == BQSL_MAGAZINE_CODE_QUARTER {
			arrGoodsType := []string{}
			if form.MagazineCodeWeek == BQSL_MAGAZINE_CODE_MONTH {
				arrGoodsType = append(arrGoodsType, BQSL_MAGAZINE_CODE_MONTH_TEXT)
			}
			if form.MagazineCodeMonth == BQSL_MAGAZINE_CODE_WEEK {
				arrGoodsType = append(arrGoodsType, BQSL_MAGAZINE_CODE_WEEK_TEXT)
			}
			if form.MagazineCodeQuarter == BQSL_MAGAZINE_CODE_QUARTER {
				arrGoodsType = append(arrGoodsType, BQSL_MAGAZINE_CODE_QUARTER_TEXT)
			}
			tag = tag + ",商品区分 IN (" + strings.Join(arrGoodsType, ",") + ")"
		}
	}
	if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT {
		tag = tag + ",フォーマット=" + `"` + DOWNLOAD_TYPE_TOTAL_RESULT_TEXT + `"`
	} else if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION {
		tag = tag + ",フォーマット=" + `"` + DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION_TEXT + `"`
	} else if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT_STORE {
		tag = tag + ",フォーマット=" + `"` + DOWNLOAD_TYPE_TOTAL_RESULT_STORE_TEXT + `"`
	}
	// ========================================================================================
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, tag, ctx, msgRetryTmp, retryMap)

	intTotalRows := int(totalRows)
	ctx.Session.Values[RPComon.REPORT_QUERY_JOB_ID+randStringFromSQL] = jobId
	ctx.Session.Values[RPComon.REPORT_QUERY_JOB_ID_COUNT+randStringFromSQL] = intTotalRows
	if err != nil {
		return &data, exterror.WrapExtError(err)
	}
	// set VJ_charging current search
	if reportVJCharging, ok := ctx.Get(RPComon.REPORT_VJ_CHARGING_KEY); ok {
		data.VJCharging = reportVJCharging.(int)
	}
	if totalRows > RPComon.BQ_DATA_LIMIT {
		return &data, exterror.WrapExtError(errors.New("Respone data too large"))
	} else if totalRows == 0 {
		return &data, nil
	}

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

	// Get data
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChan, err := conn.GetResponseDataNew(jobId, data.ShowLineFrom, limitLengthData, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return &data, exterror.WrapExtError(err)
	}

	rankNo := data.ShowLineFrom
	for {
		row := <-dataChan
		if row == nil {
			break
		}

		rankNo++
		dataHaveRankNo := []interface{}{}
		dataHaveRankNo = append(dataHaveRankNo, strconv.Itoa(rankNo))
		for _, v := range row.ValueIF {
			dataHaveRankNo = append(dataHaveRankNo, v)
		}
		data.Rows = append(data.Rows, dataHaveRankNo)
	}
	return &data, nil
}

//Create data from queryBuild
func queryGetDataWithJobId(ctx *gf.Context, sql string, form QueryForm, randStringFromSQL string) (*RpData, error) {

	keyFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_KEY_FILE)
	mailAccount := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_MAIL_ACCOUNT)
	projectId := ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_GCP_PROJECT_ID)

	data := RpData{
		HeaderCols: []string{},
		Cols:       [][]string{},
		Rows:       [][]interface{}{},
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
		return queryData(ctx, sql, form, randStringFromSQL)
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

	// Get data
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChan, err := conn.GetResponseDataNew(jobId.(string), data.ShowLineFrom, limitLengthData, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return &data, exterror.WrapExtError(err)
	}

	rankNo := data.ShowLineFrom
	for {
		row := <-dataChan
		if row == nil {
			break
		}

		rankNo++
		dataHaveRankNo := []interface{}{}
		dataHaveRankNo = append(dataHaveRankNo, strconv.Itoa(rankNo))
		for _, v := range row.ValueIF {
			dataHaveRankNo = append(dataHaveRankNo, v)
		}
		data.Rows = append(data.Rows, dataHaveRankNo)
	}
	return &data, nil
}
