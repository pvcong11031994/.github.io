package RP061_ShopSales

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"strings"
	"time"
)

func buildSql(form QueryForm, ctx *gf.Context, isCountData bool) (sql string, listRange []ModelItems.MasterCalendarItem, groupType string) {

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 日付 年月日
	dateSearchFrom := time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	dateSearchTo := time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)

	if form.DateFrom != "" {
		dateSearchFrom = strings.Replace(form.DateFrom, "/", "", -1)
	}
	if form.DateTo != "" {
		dateSearchTo = strings.Replace(form.DateTo, "/", "", -1)
	}

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// get Date =================================
	sqlRange := ""
	mcmd := Models.MasterCalendarModel{ctx.DB}
	listRange = []ModelItems.MasterCalendarItem{}
	sqlWithMasterCalendar := ""
	sqlJoinMasterCalendar := ""
	groupType = form.GroupType
	var err error
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV &&
		form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION {
		switch form.GroupType {
		case GROUP_TYPE_DATE:
			listRange, err = mcmd.GetDay(dateSearchFrom, dateSearchTo)
			Common.LogErr(err)
			for _, item := range listRange {
				sqlRange += `
			,SUM(
				CASE
				   WHEN
						mc.mc_yyyy = '` + item.Mcyyyy + `'
						AND mc.mc_mm = '` + item.Mcmm + `'
						AND mc.mc_dd = '` + item.Mcdd + `'
						AND bookstore_biz_category = '40'
					THEN SAFE_CAST((IFNULL(sales_body_quantity,0)) AS INT64)
				   ELSE 0
				END
			) as A` + item.Mcyyyy + item.Mcmm + item.Mcdd
			}
		case GROUP_TYPE_WEEK:
			listRange, err = mcmd.GetWeek(dateSearchFrom, dateSearchTo)
			Common.LogErr(err)
			for _, item := range listRange {
				sqlRange += `
			,SUM(
				CASE
				   WHEN
						mc.mc_yyyy = '` + item.Mcyyyy + `'
						AND mc.mc_weekdate = '` + item.Mcweekdate + `'
						AND bookstore_biz_category = '40'
					THEN SAFE_CAST((IFNULL(sales_body_quantity,0)) AS INT64)
				   ELSE 0
				END
			) as A` + item.Mcyyyy + item.Mcweeknum
			}
		case GROUP_TYPE_MONTH:
			listRange, err = mcmd.GetMonth(dateSearchFrom, dateSearchTo)
			Common.LogErr(err)
			for _, item := range listRange {
				sqlRange += `
			,SUM(
				CASE
				   WHEN
						mc.mc_yyyy = '` + item.Mcyyyy + `'
						AND mc.mc_mm = '` + item.Mcmm + `'
						AND bookstore_biz_category = '40'
					THEN SAFE_CAST((IFNULL(sales_body_quantity,0)) AS INT64)
				   ELSE 0
				END
			) as A` + item.Mcyyyy + item.Mcmm
			}
		}

		//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
		sqlWithMasterCalendar = `
	,mc AS (
		SELECT
			mc_yyyymmdd
			,mc_yyyy
			,mc_mm
			,mc_dd
			,mc_weekdate
			,mc_weeknum
		FROM
			{{@DATASET}}.master_calendar mc
		WHERE
			CONCAT(mc_yyyy , mc_mm , mc_dd) >= {{date_from}}  AND  CONCAT(mc_yyyy , mc_mm , mc_dd) <= {{date_to}}
	)
	`
		sqlJoinMasterCalendar = `
	LEFT JOIN mc
		ON mc.mc_yyyymmdd = SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8)
		AND mc.mc_yyyy = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8),1,4)
		AND mc.mc_mm = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8),5,2)
		AND mc.mc_dd = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8),7,2)
	`
		//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	}
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	parameter := map[string]interface{}{}
	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)

	//CONDITION
	//SHOP_CD
	parameter["shop_cd"] = form.ShopCd

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//JAN_CD_Array
	conditionJanCd := ""
	//conditionOrderBy := ""
	if len(form.JanArrays) > 0 {
		conditionJanCd = "AND jan_code IN {{jan_cd}}"
		parameter["jan_cd"] = form.JanArrays
	}
	if form.JanSingle != "" {
		conditionJanCd = "AND jan_code LIKE {{jan_cd}}"
		form.JanSingle += "%"
		parameter["jan_cd"] = form.JanSingle
	}

	//=================================================================================================================

	//DATE
	parameter["date_from"] = dateSearchFrom
	parameter["date_to"] = dateSearchTo
	parameter["month_from"] = dateSearchFrom[:6]
	parameter["month_to"] = dateSearchTo[:6]

	sqlWithData := bq.NewCommand()
	sqlWithData.Parameters = parameter
	if !isCountData {
		sqlWithData.CommandText = `
#StandardSQL
	WITH
	bqsl_condition AS (
		SELECT *
		FROM ` + "`{{@DATASET}}.bq_sales_*`" + ` bqsl
		WHERE
			SUBSTR(_TABLE_SUFFIX,-14,6) BETWEEN {{month_from}} AND {{month_to}}
			AND bqsl.shop_code IN {{shop_cd}}
			` + conditionJanCd + `
			AND SUBSTR(REPLACE(REPLACE(sales_datetime, '-', ''),'/',''),0,8) >= {{date_from}}
			AND SUBSTR(REPLACE(REPLACE(sales_datetime, '-', ''),'/',''),0,8) <= {{date_to}}
	),
	bq_sales_max_date AS (
		SELECT
			shop_code,
			jan_code,
			MAX(received_datetime) received_datetime
		FROM bqsl_condition
		GROUP BY shop_code, jan_code
	)
	` + sqlWithMasterCalendar + `

SELECT
	*,
	RANK() OVER (
		PARTITION BY shop_code
		ORDER BY
			shop_seq_number,
			shop_code,
			selling_date DESC,
			jan_code DESC
		) rank
FROM (
	SELECT
		bqsl.shop_code AS shop_code,
		bqsl.jan_code AS jan_code,
		MAX(product_name) product_name,
		MAX(author_name) author_name,
		IF(MAX(jan_maker_name) <> '' ,MAX(jan_maker_name),MAX(maker_name)) maker_name,
		MAX(IFNULL(SUBSTR(REPLACE(bqsl.selling_date, '-', ''),0,8),'')) selling_date,
		MAX(SAFE_CAST(sales_tax_exc_unit_price AS INT64)) sales_tax_exc_unit_price,
		SUM(SAFE_CAST(IF(bqsl.received_datetime < bqsl_md.received_datetime,0, bqsl.cumulative_receiving_quantity) AS INT64)) AS stok_cumulative_receiving_quantity,
		SUM(SAFE_CAST(IF(bqsl.received_datetime < bqsl_md.received_datetime,0, bqsl.cumulative_sales_quantity) AS INT64)) AS stok_cumulative_sales_quantity,
		SUM(SAFE_CAST(IF(bqsl.received_datetime < bqsl_md.received_datetime,0, bqsl.stock_quantity) AS INT64)) AS stok_stock_quantity,
		IFNULL(MIN(IF(TRIM(first_sales_date) = "",NULL, first_sales_date)),"") first_sales_date,
		SUM(IF(bookstore_biz_category = '40', SAFE_CAST((IFNULL(sales_body_quantity,0)) AS INT64), 0)) sales_body_quantity
		` + sqlRange + `,
		IFNULL(MIN(IF(TRIM(bqsl.sales_datetime) = "", NULL, bqsl.sales_datetime)), "") AS sales_datetime,
		MAX(shop_seq_number) shop_seq_number
	FROM bqsl_condition bqsl
	JOIN bq_sales_max_date bqsl_md
		ON bqsl.shop_code = bqsl_md.shop_code
		AND bqsl.jan_code = bqsl_md.jan_code
		` + sqlJoinMasterCalendar + `
	GROUP BY
		jan_code,
		shop_code
	)
ORDER BY
	shop_seq_number,
	shop_code,
	selling_date DESC,
	jan_code DESC
`
	} else {
		sqlWithData.CommandText = `
#StandardSQL
SELECT
	COUNT (*) count_data,
	shop_code
FROM (
	SELECT
		bqsl.shop_code AS shop_code,
		bqsl.jan_code AS jan_code
	FROM ` + "`{{@DATASET}}.bq_sales_*`" + ` bqsl
	WHERE
		SUBSTR(_TABLE_SUFFIX,-14,6) BETWEEN {{month_from}} AND {{month_to}}
		AND bqsl.shop_code IN {{shop_cd}}
		` + conditionJanCd + `
		AND SUBSTR(REPLACE(REPLACE(sales_datetime, '-', ''),'/',''),0,8) >= {{date_from}}
		AND SUBSTR(REPLACE(REPLACE(sales_datetime, '-', ''),'/',''),0,8) <= {{date_to}}
	GROUP BY
		jan_code,
		shop_code
	)
GROUP BY
	shop_code
`
	}
	sql, err = sqlWithData.Build()
	Common.LogErr(err)
	return
}
