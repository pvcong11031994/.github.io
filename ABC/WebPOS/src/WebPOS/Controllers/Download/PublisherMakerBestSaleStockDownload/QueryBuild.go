package PublisherMakerBestSaleStockDownload

import (
	"WebPOS/Common"
	"WebPOS/WebApp"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"strings"
	"time"
)

// Get SQL BIGNET　売上+在庫 OR  売上+在庫+返品
func buildSqlSales(form Form, ctx *gf.Context) (sql string, listHeader []string) {

	listHeader = LIST_HEADER_SALES
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 日付 年月日
	parameter := map[string]interface{}{}
	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)

	conditionJAN := ""
	// JAN条件
	if form.JAN != "" {
		conditionJAN = `
			AND jan_code LIKE {{jan_cd}}`
		parameter["jan_cd"] = form.JAN + "%"
	}

	dateSearchFrom := strings.Replace(form.DateFrom, "/", "", -1)
	dateSearchTo := strings.Replace(form.DateTo, "/", "", -1)
	if strings.TrimSpace(dateSearchFrom) == "" {
		dateSearchFrom = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	}
	if strings.TrimSpace(dateSearchTo) == "" {
		dateSearchTo = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	}
	parameter["date_from"] = dateSearchFrom
	parameter["date_to"] = dateSearchTo
	parameter["month_from"] = dateSearchFrom[:6]
	parameter["month_to"] = dateSearchTo[:6]
	parameter["shop_cd"] = form.ShopCd

	sqlData := bq.NewCommand()
	sqlData.Parameters = parameter
	sqlData.CommandText = `
#StandardSQL
WITH
	stock as (
	SELECT
			bqst.jan_code stock_jan_code,
			bqst.shop_code stock_shop_code,
			SUM(stock_quantity) stock_quantity
		FROM {{@DATASET}}.bq_stock bqst
		JOIN (
			SELECT
				jan_code,
				shop_code,
				MAX(create_datetime) create_datetime
			FROM {{@DATASET}}.bq_stock
			WHERE
				shop_code IN {{shop_cd}}
				` + conditionJAN + `
			GROUP BY
				jan_code,
				shop_code
			) bqst_max_date
			ON bqst.jan_code = bqst_max_date.jan_code
			AND bqst.shop_code = bqst_max_date.shop_code
			AND bqst.create_datetime = bqst_max_date.create_datetime
		GROUP BY
			bqst.jan_code, bqst.shop_code
	)

SELECT
	SUBSTR(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''), 0, 8) AS sales_datetime,
	MAX(bqsl.shared_book_store_code) AS shared_book_store_code,
	bqsl.shop_code AS shop_code,
	MAX(bqsl.shop_name) AS shop_name,
	bqsl.jan_code AS  jan_code,
	MAX(bqsl.product_name) AS product_name,
	IF(MAX(bqsl.jan_maker_name) <> '' ,MAX(bqsl.jan_maker_name),MAX(IFNULL(bqsl.maker_name,''))) AS maker_code,
	MAX(SAFE_CAST(bqsl.sales_tax_exc_unit_price AS INT64)) AS sales_tax_exc_unit_price,
	MAX(stock.stock_quantity) AS stock_quantity,
	SUM(IF(bookstore_biz_category = '40', bqsl.sales_body_quantity, 0)) AS sales_body_quantity,
	SUM(IF(bookstore_biz_category = '20', bqsl.sales_body_quantity, 0)) AS return_body_quantity,
	MAX(shop_seq_number) shop_seq_number
FROM
	` + "`{{@DATASET}}.bq_sales_*`" + ` bqsl
LEFT JOIN stock
	ON bqsl.jan_code = stock_jan_code
	AND bqsl.shop_code = stock_shop_code
WHERE
	SUBSTR(_TABLE_SUFFIX,-14,6) BETWEEN {{month_from}} AND {{month_to}}
	AND SUBSTR(REPLACE(REPLACE(sales_datetime, '-', ''),'/',''), 0, 8) BETWEEN {{date_from}} AND {{date_to}}
	AND bqsl.shop_code IN {{shop_cd}}
	AND bookstore_biz_category = '40'
	` + conditionJAN + `
GROUP BY
	sales_datetime,
	shop_code,
	jan_code
ORDER BY
	sales_datetime, shop_seq_number, jan_code
`

	sql, err := sqlData.Build()
	Common.LogErr(err)
	return

}

// Get SQL BIGNET　在庫のみ
func buildSqlStock(form Form, ctx *gf.Context) (sql string, listHeader []string) {

	listHeader = LIST_HEADER_STOCK
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 日付 年月日
	parameter := map[string]interface{}{}
	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)

	conditionJAN := ""
	// JAN条件
	if form.JAN != "" {
		conditionJAN = `
			AND jan_code LIKE {{jan_cd}}`
		parameter["jan_cd"] = form.JAN + "%"
	}

	parameter["shop_cd"] = form.ShopCd
	sqlData := bq.NewCommand()
	sqlData.Parameters = parameter
	sqlData.CommandText = `
#StandardSQL
SELECT
	bqst.jan_code jan_code,
	bqst.shop_code shop_code,
	MAX(shop_name) shop_name,
	MAX(shared_book_store_code) shared_book_store_code,
	MAX(product_name) product_name,
	IF(MAX(jan_maker_name) <> '' ,MAX(jan_maker_name),MAX(maker_name)) maker_code,
	MAX(list_price) list_price,
	SUM(stock_quantity) stock_quantity,
	MAX(shop_seq_number) shop_seq_number,
	MAX(jan_seq_number) jan_seq_number
FROM {{@DATASET}}.bq_stock bqst
JOIN (
	SELECT
		jan_code,
		shop_code,
		MAX(create_datetime) create_datetime
	FROM {{@DATASET}}.bq_stock
	WHERE
		shop_code IN {{shop_cd}}
		` + conditionJAN + `
	GROUP BY
		jan_code,
		shop_code
	) bqst_max_date
	ON bqst.jan_code = bqst_max_date.jan_code
	AND bqst.shop_code = bqst_max_date.shop_code
	AND bqst.create_datetime = bqst_max_date.create_datetime
GROUP BY
	bqst.jan_code, bqst.shop_code
ORDER BY
	shop_seq_number, jan_seq_number
`

	sql, err := sqlData.Build()
	Common.LogErr(err)
	return
}

// Get SQL 入荷+売上+返品 (not stock)
func buildSqlSalesReturnsReceive(form Form, ctx *gf.Context) (sql string, listHeader []string) {

	listHeader = LIST_HEADER_SALES_AND_RECEIVING
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 日付 年月日
	parameter := map[string]interface{}{}
	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)

	conditionJAN := ""
	// JAN条件
	if form.JAN != "" {
		conditionJAN = `
			AND jan_code LIKE {{jan_cd}}`
		parameter["jan_cd"] = form.JAN + "%"
	}
	dateSearchFrom := strings.Replace(form.DateFrom, "/", "", -1)
	dateSearchTo := strings.Replace(form.DateTo, "/", "", -1)
	if strings.TrimSpace(dateSearchFrom) == "" {
		dateSearchFrom = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	}
	if strings.TrimSpace(dateSearchTo) == "" {
		dateSearchTo = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	}
	parameter["date_from"] = dateSearchFrom
	parameter["date_to"] = dateSearchTo
	parameter["month_from"] = dateSearchFrom[:6]
	parameter["month_to"] = dateSearchTo[:6]
	parameter["shop_cd"] = form.ShopCd

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// process 返品
	selectReturn := ""
	conditionReturn := ""
	if form.DataMode == TYPE_SEARCH_SALES_RETURN {
		selectReturn += `
			SUM(IF(bookstore_biz_category = '20', bqsl.sales_body_quantity, 0)) AS return_body_quantity,
		`
		conditionReturn += `
			OR bookstore_biz_category = '20'
		`
		listHeader = LIST_HEADER_SALES_AND_RETURN
	}
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sqlData := bq.NewCommand()
	sqlData.Parameters = parameter
	sqlData.CommandText = `
#StandardSQL
SELECT
	SUBSTR(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''), 0, 8) AS sales_datetime,
	MAX(bqsl.shared_book_store_code) AS shared_book_store_code,
	bqsl.shop_code AS shop_code,
	MAX(bqsl.shop_name) AS shop_name,
	bqsl.jan_code AS  jan_code,
	MAX(bqsl.product_name) AS product_name,
	IF(MAX(bqsl.jan_maker_name) <> '' ,MAX(bqsl.jan_maker_name),MAX(IFNULL(bqsl.maker_name,''))) AS maker_code,
	MAX(SAFE_CAST(bqsl.sales_tax_exc_unit_price AS INT64)) AS sales_tax_exc_unit_price,
	SUM(IF(bookstore_biz_category = '30', bqsl.sales_body_quantity, 0)) AS receiving_body_quantity,
	SUM(IF(bookstore_biz_category = '40', bqsl.sales_body_quantity, 0)) AS sales_body_quantity,
	` + selectReturn + `
	MAX(shop_seq_number) shop_seq_number
FROM
	` + "`{{@DATASET}}.bq_sales_*`" + ` bqsl
WHERE
	SUBSTR(_TABLE_SUFFIX,-14,6) BETWEEN {{month_from}} AND {{month_to}}
	AND SUBSTR(REPLACE(REPLACE(sales_datetime, '-', ''),'/',''), 0, 8) BETWEEN {{date_from}} AND {{date_to}}
	AND bqsl.shop_code IN {{shop_cd}}
	AND (bookstore_biz_category = '40' OR bookstore_biz_category = '30' ` + conditionReturn + ` )
	` + conditionJAN + `
GROUP BY
	sales_datetime,
	shop_code,
	jan_code
ORDER BY
	sales_datetime, shop_seq_number, jan_code
`

	sql, err := sqlData.Build()
	Common.LogErr(err)
	return

}
