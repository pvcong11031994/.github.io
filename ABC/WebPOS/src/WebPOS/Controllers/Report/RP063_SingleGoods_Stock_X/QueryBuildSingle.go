package RP063_SingleGoods_Stock_X

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"strings"
	"time"
)

func buildSqlSingle(form QueryFormSingleGoods, ctx *gf.Context) (sql string, listRange []ModelItems.MasterCalendarItem) {

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 日付 年月日
	dateSearchFrom := strings.Replace(form.DateFrom, "/", "", -1)
	dateSearchTo := strings.Replace(form.DateTo, "/", "", -1)

	mcmd := Models.MasterCalendarModel{ctx.DB}
	sqlRange := ""
	var err error
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
			sqlRange += `
				,SUM(
					CASE
					   WHEN
							mc.mc_yyyy = '` + item.Mcyyyy + `'
							AND mc.mc_mm = '` + item.Mcmm + `'
							AND mc.mc_dd = '` + item.Mcdd + `'
							AND bookstore_biz_category = '30'
						THEN SAFE_CAST((IFNULL(sales_body_quantity,0)) AS INT64)
					   ELSE 0
					END
				) as B` + item.Mcyyyy + item.Mcmm + item.Mcdd
			// 期間入荷累計
			sqlRange += `
				,MAX(
					CASE
					   WHEN
							mc.mc_yyyy = '` + item.Mcyyyy + `'
							AND mc.mc_mm = '` + item.Mcmm + `'
							AND mc.mc_dd = '` + item.Mcdd + `'
						THEN SAFE_CAST((IFNULL(cumulative_receiving_quantity,0)) AS INT64)
					   ELSE 0
					END
				) as C` + item.Mcyyyy + item.Mcmm + item.Mcdd
			// 期間売上累計
			sqlRange += `
				,MAX(
					CASE
					   WHEN
							mc.mc_yyyy = '` + item.Mcyyyy + `'
							AND mc.mc_mm = '` + item.Mcmm + `'
							AND mc.mc_dd = '` + item.Mcdd + `'
						THEN SAFE_CAST((IFNULL(cumulative_sales_quantity,0)) AS INT64)
					   ELSE 0
					END
				) as D` + item.Mcyyyy + item.Mcmm + item.Mcdd
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
			sqlRange += `
				,SUM(
					CASE
					   WHEN
							mc.mc_yyyy = '` + item.Mcyyyy + `'
							AND mc.mc_weekdate = '` + item.Mcweekdate + `'
							AND bookstore_biz_category = '30'
						THEN SAFE_CAST((IFNULL(sales_body_quantity,0)) AS INT64)
					   ELSE 0
					END
				) as B` + item.Mcyyyy + item.Mcweeknum

			// 期間入荷累計
			sqlRange += `
				,MAX(
					CASE
					   WHEN
							mc.mc_yyyy = '` + item.Mcyyyy + `'
							AND mc.mc_weekdate = '` + item.Mcweekdate + `'
						THEN SAFE_CAST((IFNULL(cumulative_receiving_quantity,0)) AS INT64)
					   ELSE 0
					END
				) as C` + item.Mcyyyy + item.Mcweeknum
			// 期間売上累計
			sqlRange += `
				,MAX(
					CASE
					   WHEN
							mc.mc_yyyy = '` + item.Mcyyyy + `'
							AND mc.mc_weekdate = '` + item.Mcweekdate + `'
						THEN SAFE_CAST((IFNULL(cumulative_sales_quantity,0)) AS INT64)
					   ELSE 0
					END
				) as D` + item.Mcyyyy + item.Mcweeknum
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
			sqlRange += `
				,SUM(
					CASE
					   WHEN
							mc.mc_yyyy = '` + item.Mcyyyy + `'
							AND mc.mc_mm = '` + item.Mcmm + `'
							AND bookstore_biz_category = '30'
						THEN SAFE_CAST((IFNULL(sales_body_quantity,0)) AS INT64)
					   ELSE 0
					END
				) as B` + item.Mcyyyy + item.Mcmm
			// 期間入荷累計
			sqlRange += `
				,MAX(
					CASE
					   WHEN
							mc.mc_yyyy = '` + item.Mcyyyy + `'
							AND mc.mc_mm = '` + item.Mcmm + `'
						THEN SAFE_CAST((IFNULL(cumulative_receiving_quantity,0)) AS INT64)
					   ELSE 0
					END
				) as C` + item.Mcyyyy + item.Mcmm
			// 期間売上累計
			sqlRange += `
				,MAX(
					CASE
					   WHEN
							mc.mc_yyyy = '` + item.Mcyyyy + `'
							AND mc.mc_mm = '` + item.Mcmm + `'
						THEN SAFE_CAST((IFNULL(cumulative_sales_quantity,0)) AS INT64)
					   ELSE 0
					END
				) as D` + item.Mcyyyy + item.Mcmm
		}

	}
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	parameter := map[string]interface{}{}
	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)

	//CONDITION
	//SHOP_CD
	parameter["shop_cd"] = form.ShopCd

	//JAN_CD
	if len(form.JAN) < 5 {
		return
	}
	parameter["jan_cd"] = form.JAN
	parameter["jan_index"] = form.JAN[4:5]

	//DATE
	parameter["date_from"] = dateSearchFrom
	parameter["date_to"] = dateSearchTo
	parameter["month_from"] = dateSearchFrom[:6]
	parameter["month_to"] = dateSearchTo[:6]
	dateCurrent := time.Now().Format(Common.DATE_FORMAT_YMD)
	parameter["month_current"] = dateCurrent[:6]

	sqlWithData := bq.NewCommand()
	sqlWithData.Parameters = parameter
	sqlWithData.CommandText = `
#StandardSQL
SELECT
	bqsl.shop_code AS shop_code,
	MAX(shop_name) AS shop_name,
	bqsl.jan_code AS jan_code,
	SUM(IF(bookstore_biz_category = '40', SAFE_CAST((IFNULL(sales_body_quantity,0)) AS INT64), 0)) sales_body_quantity
	` + sqlRange + `,
	MAX(bqsl.shop_seq_number) AS shop_seq_number
FROM ` + "`{{@DATASET}}.bq_sales_*`" + ` bqsl
JOIN {{@DATASET}}.master_calendar mc
	ON mc.mc_yyyymmdd = SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8)
	AND mc.mc_yyyy = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8),1,4)
	AND mc.mc_mm = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8),5,2)
	AND mc.mc_dd = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8),7,2)
	AND CONCAT(mc_yyyy , mc_mm , mc_dd) >= {{date_from}}  AND  CONCAT(mc_yyyy , mc_mm , mc_dd) <= {{date_to}}
WHERE
	SUBSTR(_TABLE_SUFFIX,-14,6) BETWEEN {{month_from}} AND {{month_to}}
	AND bqsl.shop_code IN {{shop_cd}}
	AND bqsl.jan_code = {{jan_cd}}
	AND SUBSTR(REPLACE(REPLACE(sales_datetime, '-', ''),'/',''),0,8) >= {{date_from}}
	AND SUBSTR(REPLACE(REPLACE(sales_datetime, '-', ''),'/',''),0,8) <= {{date_to}}
GROUP BY
	shop_code,
	jan_code
ORDER BY shop_seq_number
`

	sql, err = sqlWithData.Build()
	Common.LogErr(err)
	return

}

func buildSqlStockAll(form QueryFormSingleGoods, ctx *gf.Context) (sql string) {

	parameter := map[string]interface{}{}
	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)

	//CONDITION
	//SHOP_CD
	parameter["shop_cd"] = form.ShopCd

	//JAN_CD
	if len(form.JAN) < 5 {
		return
	}
	parameter["jan_cd"] = form.JAN
	parameter["jan_index"] = form.JAN[4:5]

	sqlWithData := bq.NewCommand()
	sqlWithData.Parameters = parameter
	sqlWithData.CommandText = `
#StandardSQL
SELECT
	bqst.jan_code jan_code,
	bqst.shop_code shop_code,
	MAX(product_name) product_name,
	MAX(author_name) author_name,
	IF(MAX(jan_maker_name) <> '' ,MAX(jan_maker_name),MAX(maker_name)) maker_name,
	MAX(IFNULL(SUBSTR(REPLACE(selling_date, '-', ''),0,8),'')) selling_date,
	MAX(list_price) list_price,
	IFNULL(MIN(IF(TRIM(first_sales_date) = "",NULL, SUBSTR(REPLACE(first_sales_date, '-', ''),0,8))),"") first_sales_date,
	SUM(stock_quantity) stock_quantity,
	SUM(cumulative_sales_quantity) cumulative_sales_quantity,
	SUM(cumulative_receiving_quantity) cumulative_receiving_quantity
FROM ` + "`{{@DATASET}}.bq_stock_*`" + ` bqst
JOIN (
	SELECT
		jan_code,
		shop_code,
		MAX(create_datetime) create_datetime
	FROM ` + "`{{@DATASET}}.bq_stock_*`" + `
	WHERE
		shop_code IN {{shop_cd}}
		AND jan_code = {{jan_cd}}
		AND _TABLE_SUFFIX = {{jan_index}}
	GROUP BY
		jan_code,
		shop_code
	) bqst_max_date
	ON bqst.jan_code = bqst_max_date.jan_code
	AND bqst.shop_code = bqst_max_date.shop_code
	AND bqst.create_datetime = bqst_max_date.create_datetime
WHERE
	_TABLE_SUFFIX = {{jan_index}}
	AND bqst.shop_code IN {{shop_cd}}
	AND bqst.jan_code = {{jan_cd}}
GROUP BY
	bqst.jan_code,
	bqst.shop_code
`

	sql, err := sqlWithData.Build()
	Common.LogErr(err)
	return

}
