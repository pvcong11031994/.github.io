package RP053_BestSalesByStore

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

func queryData(ctx *gf.Context, sql string, rowField []string, colField []string, sumField []string, form QueryForm) (*RPComon.ReportData, error) {

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
		return nil, exterror.WrapExtError(err)
	}

	// set report name to import info log search charging
	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
	//totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, _REPORT_ID)
	// ========================================================================================
	// Output log search condition
	tag := "report=" + _REPORT_ID
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV_TEXT + `"`
	} else {
		tag = tag + ",handle=" + `"` + RPComon.REPORT_SEARCH_TYPE_HANDLE_TEXT + `"`
	}
	if form.GroupType == GROUP_TYPE_DATE {
		tag = tag + ",表示=" + `"` + GROUP_TYPE_DATE_TEXT + `"`
		tag = tag + ",期間=" + `"` + form.DateFrom + "～" + form.DateTo + `"`
	} else if form.GroupType == GROUP_TYPE_WEEK {
		tag = tag + ",表示=" + `"` + GROUP_TYPE_WEEK_TEXT + `"`
		tag = tag + ",期間=" + `"` + form.WeekFrom + "～" + form.WeekTo + `"`
	} else if form.GroupType == GROUP_TYPE_MONTH {
		tag = tag + ",表示=" + `"` + GROUP_TYPE_MONTH_TEXT + `"`
		tag = tag + ",期間=" + `"` + form.MonthFrom + "～" + form.MonthTo + `"`
	}
	tag = tag + ",店舗=" + `"` + form.ShopCd + `"`
	if len(form.MediaGroup1Cd) > 0 {
		tag = tag + ",メディア大分類コード IN (" + strings.Join(form.MediaGroup1Cd, ",") + ")"
	}
	if len(form.MediaGroup2Cd) > 0 {
		tag = tag + ",メディア中分類コード IN (" + strings.Join(form.MediaGroup2Cd, ",") + ")"
	}
	if len(form.MediaGroup3Cd) > 0 {
		tag = tag + ",メディア中小分類コード IN (" + strings.Join(form.MediaGroup3Cd, ",") + ")"
	}
	if len(form.MakerCd) > 0 {
		tag = tag + ",出版社 LIKE (" + Common.JoinArray(form.MakerCd, "%", "%", ",") + ")"
	}
	// ========================================================================================
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	totalRows, jobId, err := conn.QueryForResponseBySql(sql, ctx, tag, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return nil, exterror.WrapExtError(err)
	}

	data := RPComon.ReportData{
		ListColKey: []string{},
		ListRowKey: []string{},

		Cols: map[string][]string{},
		Rows: map[string][]string{},
		Data: map[string]map[string][]interface{}{},
	}
	// set VJ_charging current search
	if reportVJCharging, ok := ctx.Get(RPComon.REPORT_VJ_CHARGING_KEY); ok {
		data.VJCharging = reportVJCharging.(int)
	}
	if totalRows > RPComon.BQ_DATA_LIMIT {
		return nil, exterror.WrapExtError(errors.New("Respone data too large"))
	}
	keyErr = errors.New("KEY_ERR")
	msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
	dataChan, err := conn.GetResponseData(jobId, 0, RPComon.BQ_DATA_LIMIT, ctx, msgRetryTmp, retryMap)

	if err != nil {
		return nil, exterror.WrapExtError(err)
	}

	data.Data[RPComon.SUM_KEY_FIELD] = map[string][]interface{}{}
	data.Data[RPComon.NO_KEY_FIELD] = map[string][]interface{}{}

	rowCount := 0

	for {
		row := <-dataChan
		if row == nil {
			break
		}

		data.CountResultRows++

		sumData := []interface{}{}

		for _, k := range sumField {
			if (*row.TypeMap)[k] == "INTEGER" {
				sumData = append(sumData, row.ValueMap[k].Int())
			} else if (*row.TypeMap)[k] == "FLOAT" {
				sumData = append(sumData, row.ValueMap[k].Float())
			} else {
				sumData = append(sumData, row.ValueMap[k].String())
			}
		}
		colKey := row.ValueMap["col_key"].String()
		rowKey := row.ValueMap["row_key"].String()

		switch row.ValueMap["group_code"].String() {
		case "1":
			data.ListColKey = append(data.ListColKey, colKey)
			colHeader := []string{}
			for _, k := range colField {
				colHeader = append(colHeader, row.ValueMap[k].String())
			}
			data.Cols[colKey] = colHeader
		case "2":
			rowCount++
			data.ListRowKey = append(data.ListRowKey, rowKey)
			rowHeader := []string{}
			for _, k := range rowField {
				rowHeader = append(rowHeader, row.ValueMap[k].String())
			}
			data.Rows[rowKey] = rowHeader
			data.Data[rowKey] = map[string][]interface{}{}
		}

		if data.Data[rowKey] != nil {
			data.Data[rowKey][colKey] = sumData
		}
	}
	//======================================================

	mapColName, mapRowName, mapSumName := initDefaultLayout()
	formatSum := []string{}
	headerSum := []string{}
	for _, v := range sumField {
		headerSum = append(headerSum, mapSumName[v])
		formatSum = append(formatSum, "number")
	}

	headerRow := []string{}
	for _, v := range rowField {
		headerRow = append(headerRow, mapRowName[v])
	}

	headerCol := []string{}
	for _, v := range colField {
		headerCol = append(headerCol, mapColName[v])
	}

	data.FormatSum = formatSum
	data.HeaderCol = headerCol
	data.HeaderRow = headerRow
	data.HeaderSum = headerSum

	return &data, nil
}
