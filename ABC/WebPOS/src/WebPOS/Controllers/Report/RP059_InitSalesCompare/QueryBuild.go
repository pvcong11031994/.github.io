package RP059_InitSalesCompare

import (
	"github.com/goframework/gf"
	"strconv"
	"strings"
)

////Check condition and export query
//// tab JAN/ISBN
//func buildJan(form QueryForm, ctx *gf.Context) (sql string, listHeader []string) {
//
//	listHeader = append(
//		listHeader,
//		"",
//		"JAN",
//		"品名",
//		"著者名",
//		"出版社",
//		"本体価格",
//		"在庫数",
//		"最終売上日",
//		"売上累計",
//		"入荷累計",
//		"発売日",
//	)
//
//	parameter := map[string]interface{}{}
//	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)
//
//	//CONDITION
//	sqlCondition := ""
//	//=================================================================================================================
//	//SHOP_CD
//	sqlCondition += `
//		shop_code IN {{shop_cd}}
//	`
//	parameter["shop_cd"] = form.ShopCd
//	//=================================================================================================================
//	//JAN_CD
//	sqlCondition += `
//		AND jan_code IN {{jan_cd}}
//	`
//	parameter["jan_cd"] = form.JanArrays
//	//=================================================================================================================
//	// range data from 売上初速01 -> 売上初速40
//	selectRangeArray := []string{}
//	for index := 1; index <= form.SearchDateType; index++ {
//		listHeader = append(listHeader, strconv.Itoa(index)+"日目")
//		item := "SUM(sales_quantity_day"
//		if index < 10 {
//			item += "0" + strconv.Itoa(index)
//		} else {
//			item += strconv.Itoa(index)
//		}
//		item += ") sales_quantity_day" + strconv.Itoa(index)
//		selectRangeArray = append(selectRangeArray, item)
//	}
//	selectRangeStr := strings.Join(selectRangeArray, ",")
//	//=================================================================================================================
//	// Create SQL
//	sqlData := bq.NewCommand()
//	sqlData.Parameters = parameter
//	sqlData.CommandText = `
//#StandardSQL
//WITH
//	bqst AS (
//		SELECT
//			*
//		FROM {{@DATASET}}.bq_stock
//		WHERE
//		` + sqlCondition + `
//	),
//	bqst_max_date AS (
//		SELECT
//			jan_code,
//			shop_code,
//			MAX(create_datetime) create_datetime
//		FROM bqst
//		GROUP BY
//			jan_code,
//			shop_code
//	)
//
//	SELECT
//		bqst.jan_code,
//		MAX(product_name) product_name,
//		MAX(author_name) author_name,
//		IF(MAX(jan_maker_name) <> '' ,MAX(jan_maker_name),MAX(maker_name)) maker_name,
//		MAX(list_price) list_price,
//		SUM(stock_quantity) stock_quantity,
//		MAX(last_sales_date) last_sales_date,
//		SUM(cumulative_sales_quantity) cumulative_sales_quantity,
//		SUM(cumulative_receiving_quantity) cumulative_receiving_quantity,
//		IFNULL(MIN(IF(TRIM(selling_date) = "",NULL, selling_date)),"") selling_date,
//		` + selectRangeStr + `
//	FROM bqst
//	JOIN bqst_max_date
//		ON bqst.jan_code = bqst_max_date.jan_code
//		AND bqst.shop_code = bqst_max_date.shop_code
//		AND bqst.create_datetime = bqst_max_date.create_datetime
//	GROUP BY
//		bqst.jan_code
//	ORDER BY
//		selling_date DESC,
//		bqst.jan_code
//	LIMIT 15
//`
//
//	sql, err := sqlData.Build()
//	Common.LogErr(err)
//	return
//}
//
////Check condition and export query
//// Tab 雑誌コード
//func buildMagazine(form QueryForm, ctx *gf.Context) (sql string, listHeader []string) {
//
//	listHeader = append(
//		listHeader,
//		"",
//		"JAN",
//		"月号",
//		"本体価格",
//		"在庫数",
//		"最終売上日",
//		"売上累計",
//		"入荷累計",
//		"完売率",
//		"発売日",
//	)
//
//	parameter := map[string]interface{}{}
//	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)
//
//	//CONDITION
//	sqlCondition := ""
//	//=================================================================================================================
//	//SHOP_CD
//	sqlCondition += `
//		shop_code IN {{shop_cd}}
//		AND jan_grouping = '2'
//	`
//	parameter["shop_cd"] = form.ShopCd
//
//	//=================================================================================================================
//	// 雑誌コード
//	if form.ControlType == CONTROL_TYPE_MAGAZINE && form.MagazineCdSingle != "" {
//		sqlCondition += `
//			AND magazine_code = {{magazine_cd_single}}
//		`
//	}
//	parameter["magazine_cd_single"] = form.MagazineCdSingle
//	//=================================================================================================================
//	// range data from 売上初速01 -> 売上初速40
//	selectRangeArray := []string{}
//	for index := 1; index <= form.SearchDateType; index++ {
//		listHeader = append(listHeader, strconv.Itoa(index)+"日目")
//		item := "SUM(sales_quantity_day"
//		if index < 10 {
//			item += "0" + strconv.Itoa(index)
//		} else {
//			item += strconv.Itoa(index)
//		}
//		item += ") sales_quantity_day" + strconv.Itoa(index)
//		selectRangeArray = append(selectRangeArray, item)
//	}
//	selectRangeStr := strings.Join(selectRangeArray, ",")
//	//=================================================================================================================
//	// Create SQL
//	sqlData := bq.NewCommand()
//	sqlData.Parameters = parameter
//	sqlData.CommandText = `
//#StandardSQL
//WITH
//	bqst AS (
//		SELECT
//			*
//		FROM {{@DATASET}}.bq_stock
//		WHERE
//		` + sqlCondition + `
//	),
//	bqst_max_date AS (
//		SELECT
//			jan_code,
//			shop_code,
//			MAX(create_datetime) create_datetime
//		FROM bqst
//		GROUP BY
//			jan_code,
//			shop_code
//	)
//
//	SELECT
//		MAX(magazine_name) magazine_name,
//		MAX(author_name) author_name,
//		IF(MAX(jan_maker_name) <> '' ,MAX(jan_maker_name),MAX(maker_name)) maker_name,
//		bqst.jan_code,
//		MAX(jan_seq_number) jan_seq_number,
//		SUBSTR(bqst.jan_code,10,2) month_num,
//		MAX(list_price) list_price,
//		SUM(stock_quantity) stock_quantity,
//		MAX(last_sales_date) last_sales_date,
//		SUM(cumulative_sales_quantity) cumulative_sales_quantity,
//		SUM(cumulative_receiving_quantity) cumulative_receiving_quantity,
//		IFNULL(MIN(IF(TRIM(selling_date) = "",NULL, selling_date)),"") selling_date,
//		` + selectRangeStr + `
//	FROM bqst
//	JOIN bqst_max_date
//		ON bqst.jan_code = bqst_max_date.jan_code
//		AND bqst.shop_code = bqst_max_date.shop_code
//		AND bqst.create_datetime = bqst_max_date.create_datetime
//	GROUP BY
//		bqst.jan_code
//	ORDER BY
//		jan_seq_number DESC
//	`
//
//	sql, err := sqlData.Build()
//	Common.LogErr(err)
//	return
//}
//
////Check condition and export query
//// Data detail by JAN
//func buildDetail(form QueryForm, ctx *gf.Context) (sql string, listHeader []string) {
//
//	listHeader = append(
//		listHeader,
//		"店舗名",
//		"売上累計",
//		"入荷累計",
//		"在庫数",
//		"初売上日",
//	)
//
//	parameter := map[string]interface{}{}
//	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)
//
//	//CONDITION
//	sqlCondition := ""
//	//=================================================================================================================
//	//SHOP_CD
//	sqlCondition += `
//		shop_code IN {{shop_cd}}
//	`
//	parameter["shop_cd"] = form.ShopCd
//	//=================================================================================================================
//	// 雑誌コード
//	if form.ControlType == CONTROL_TYPE_MAGAZINE && form.MagazineCdSingle != "" {
//		sqlCondition += `
//			AND magazine_code = {{magazine_cd_single}}
//		`
//	}
//	parameter["magazine_cd_single"] = form.MagazineCdSingle
//	//=================================================================================================================
//	//JAN_CD
//	sqlCondition += `
//		AND jan_code = {{jan_cd}}
//	`
//	parameter["jan_cd"] = form.JanKey
//	//=================================================================================================================
//
//	// range data from 売上初速01 -> 売上初速40
//	selectRangeArray := []string{}
//	for index := 1; index <= form.SearchDateType; index++ {
//		listHeader = append(listHeader, strconv.Itoa(index)+"日目")
//		item := "SUM(sales_quantity_day"
//		if index < 10 {
//			item += "0" + strconv.Itoa(index)
//		} else {
//			item += strconv.Itoa(index)
//		}
//		item += ") sales_quantity_day" + strconv.Itoa(index)
//		selectRangeArray = append(selectRangeArray, item)
//	}
//	selectRangeStr := strings.Join(selectRangeArray, ",")
//	//=================================================================================================================
//	// Create SQL
//	sqlData := bq.NewCommand()
//	sqlData.Parameters = parameter
//	sqlData.CommandText = `
//#StandardSQL
//WITH
//	bqst AS (
//		SELECT
//			*
//		FROM {{@DATASET}}.bq_stock
//		WHERE
//		` + sqlCondition + `
//	),
//	bqst_max_date AS (
//		SELECT
//			jan_code,
//			shop_code,
//			MAX(create_datetime) create_datetime
//		FROM bqst
//		GROUP BY
//			jan_code,
//			shop_code
//	)
//
//	SELECT
//		MAX(shop_name) shop_name,
//		SUM(cumulative_sales_quantity) cumulative_sales_quantity,
//		SUM(cumulative_receiving_quantity) cumulative_receiving_quantity,
//		SUM(stock_quantity) stock_quantity,
//		MAX(shop_seq_number) shop_seq_number,
//		IFNULL(MIN(IF(TRIM(first_sales_date) = "",NULL, first_sales_date)),"") first_sales_date,
//		` + selectRangeStr + `
//	FROM bqst
//	JOIN bqst_max_date
//		ON bqst.jan_code = bqst_max_date.jan_code
//		AND bqst.shop_code = bqst_max_date.shop_code
//		AND bqst.create_datetime = bqst_max_date.create_datetime
//	GROUP BY
//		bqst.jan_code,
//		bqst.shop_code
//	ORDER BY
//		shop_seq_number
//	`
//
//	sql, err := sqlData.Build()
//	Common.LogErr(err)
//	return
//}
//
//// Check condition and export query
//// Data detail by JAN (export for CSV)
//func buildDetailForCSV(form QueryForm, ctx *gf.Context, listJanKey []string) (sql string, listHeader []string) {
//
//	listHeader = append(
//		listHeader,
//		"店舗名",
//		"売上累計",
//		"入荷累計",
//		"在庫数",
//		"初売上日",
//	)
//
//	parameter := map[string]interface{}{}
//	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)
//
//	//CONDITION
//	sqlCondition := ""
//	sortCondition := ""
//	sortSelected := ""
//	if form.ControlType == CONTROL_TYPE_JAN {
//		sortSelected += `
//			IFNULL(MIN(IF(TRIM(selling_date) = "",NULL, selling_date)),"") selling_date,
//		`
//		sortCondition += " selling_date DESC,"
//	} else {
//		sortSelected += " MAX(jan_seq_number) jan_seq_number,"
//		sortCondition += " jan_seq_number DESC,"
//	}
//	//=================================================================================================================
//	//SHOP_CD
//	sqlCondition += `
//		shop_code IN {{shop_cd}}
//	`
//	parameter["shop_cd"] = form.ShopCd
//	//=================================================================================================================
//	// 雑誌コード
//	if form.ControlType == CONTROL_TYPE_MAGAZINE && form.MagazineCdSingle != "" {
//		sqlCondition += `
//			AND magazine_code = {{magazine_cd_single}}
//		`
//	}
//	parameter["magazine_cd_single"] = form.MagazineCdSingle
//	//=================================================================================================================
//	//JAN_CD
//	sqlCondition += `
//		AND jan_code IN {{jan_cd}}
//	`
//	parameter["jan_cd"] = listJanKey
//	//=================================================================================================================
//	// range data from 売上初速01 -> 売上初速40
//	selectRangeArray := []string{}
//	for index := 1; index <= form.SearchDateType; index++ {
//		listHeader = append(listHeader, strconv.Itoa(index)+"日目")
//		item := "SUM(sales_quantity_day"
//		if index < 10 {
//			item += "0" + strconv.Itoa(index)
//		} else {
//			item += strconv.Itoa(index)
//		}
//		item += ") sales_quantity_day" + strconv.Itoa(index)
//		selectRangeArray = append(selectRangeArray, item)
//	}
//	selectRangeStr := strings.Join(selectRangeArray, ",")
//	//=================================================================================================================
//	// Create SQL
//	sqlData := bq.NewCommand()
//	sqlData.Parameters = parameter
//	sqlData.CommandText = `
//#StandardSQL
//WITH
//	bqst AS (
//		SELECT
//			*
//		FROM {{@DATASET}}.bq_stock
//		WHERE
//		` + sqlCondition + `
//	),
//	bqst_max_date AS (
//		SELECT
//			jan_code,
//			shop_code,
//			MAX(create_datetime) create_datetime
//		FROM bqst
//		GROUP BY
//			jan_code,
//			shop_code
//	)
//
//	SELECT
//		bqst.jan_code,
//		MAX(shop_name) shop_name,
//		SUM(cumulative_sales_quantity) cumulative_sales_quantity,
//		SUM(cumulative_receiving_quantity) cumulative_receiving_quantity,
//		SUM(stock_quantity) stock_quantity,
//		MAX(shop_seq_number) shop_seq_number,
//		IFNULL(MIN(IF(TRIM(first_sales_date) = "",NULL, first_sales_date)),"") first_sales_date,
//		` + sortSelected + selectRangeStr + `
//	FROM bqst
//	JOIN bqst_max_date
//		ON bqst.jan_code = bqst_max_date.jan_code
//		AND bqst.shop_code = bqst_max_date.shop_code
//		AND bqst.create_datetime = bqst_max_date.create_datetime
//	GROUP BY
//		bqst.jan_code,
//		bqst.shop_code
//	ORDER BY
//		` + sortCondition + `
//		bqst.jan_code,
//		shop_seq_number
//	`
//
//	sql, err := sqlData.Build()
//	Common.LogErr(err)
//	return
//}

//Check condition and export query
// tab JAN/ISBN
func buildJan(form QueryForm, ctx *gf.Context) (sql string, sqlCache string, listHeader []string) {

	listHeader = append(
		listHeader,
		"",
		"JAN",
		"品名",
		"著者名",
		"出版社",
		"本体価格",
		"在庫数",
		"最終売上日",
		"売上累計",
		"入荷累計",
		"発売日",
	)

	//=================================================================================================================
	// range data from 売上初速01 -> 売上初速40
	selectRangeArray := []string{}
	for index := 1; index <= form.SearchDateType; index++ {
		listHeader = append(listHeader, strconv.Itoa(index)+"日目")
		item := "SUM(d.sales_quantity_day"
		if index < 10 {
			item += "0" + strconv.Itoa(index)
		} else {
			item += strconv.Itoa(index)
		}
		item += ") sales_quantity_day" + strconv.Itoa(index)
		selectRangeArray = append(selectRangeArray, item)
	}
	selectRangeStr := strings.Join(selectRangeArray, ",")
	//=================================================================================================================
	sql = `
SELECT
	a.jan_code jan_code,
	MAX(b.product_name) product_name,
	MAX(b.author_name) author_name,
	MAX(b.maker_name) maker_name,
	MAX(b.list_price) list_price,
	SUM(a.stock_quantity) stock_quantity,
	MAX(a.last_sales_date) last_sales_date,
	SUM(a.cumulative_sales_quantity) cumulative_sales_quantity,
	SUM(a.cumulative_receiving_quantity) cumulative_receiving_quantity,
	IFNULL(MIN(IF(TRIM(b.selling_date) = "",NULL, b.selling_date)),"") selling_date,
	SUBSTR(a.jan_code,10,2) month_num,
	` + selectRangeStr + `
FROM
	m_stock a
LEFT OUTER JOIN m_jan b
	ON a.jan_code = b.jan_code
LEFT OUTER JOIN m_shop c
	ON a.shop_code = c.shop_code
LEFT OUTER JOIN m_initial_sales d
	ON a.shop_code = d.shop_code
	AND a.jan_code = d.jan_code
WHERE
	a.shop_code IN (?` + strings.Repeat(",?", len(form.ShopCd)-1) + `)
	AND a.jan_code IN (?` + strings.Repeat(",?", len(form.JanArrays)-1) + `)
GROUP BY
	a.jan_code
ORDER BY
	selling_date DESC,
	a.jan_code
LIMIT 15
`
	sqlCache = sql + sqlCache
	for _, s := range form.ShopCd {
		sqlCache = sqlCache + s
	}
	for _, s := range form.JanArrays {
		sqlCache = sqlCache + s
	}
	return
}

//Check condition and export query
// Tab 雑誌コード
func buildMagazine(form QueryForm, ctx *gf.Context) (sql string, sqlCache string, listHeader []string) {

	listHeader = append(
		listHeader,
		"",
		"JAN",
		"品名",
		"雑誌コード+月号",
		"出版社名",
		//"月号",
		"本体価格",
		"在庫数",
		"最終売上日",
		"売上累計",
		"入荷累計",
		"完売率",
		"発売日",
	)

	//=================================================================================================================
	// range data from 売上初速01 -> 売上初速40
	selectRangeArray := []string{}
	for index := 1; index <= form.SearchDateType; index++ {
		listHeader = append(listHeader, strconv.Itoa(index)+"日目")
		item := "SUM(d.sales_quantity_day"
		if index < 10 {
			item += "0" + strconv.Itoa(index)
		} else {
			item += strconv.Itoa(index)
		}
		item += ") sales_quantity_day" + strconv.Itoa(index)
		selectRangeArray = append(selectRangeArray, item)
	}
	selectRangeStr := strings.Join(selectRangeArray, ",")
	//=================================================================================================================
	sql = `
SELECT
	MAX(SUBSTR(a.jan_code, 5, 7)) magazine_name,
	MAX(b.author_name) author_name,
	MAX(b.maker_name) maker_name,
	a.jan_code jan_code,
	MAX(b.product_name) product_name,
	-- MAX(b.jan_seq_number) jan_seq_number,
	MAX(b.magazine_seq_number) magazine_seq_number,
	-- SUBSTR(a.jan_code,10,2) month_num,
	MAX(b.list_price) list_price,
	SUM(a.stock_quantity) stock_quantity,
	MAX(a.last_sales_date) last_sales_date,
	SUM(a.cumulative_sales_quantity) cumulative_sales_quantity,
	SUM(a.cumulative_receiving_quantity) cumulative_receiving_quantity,
	IFNULL(MIN(IF(TRIM(b.selling_date) = "",NULL, b.selling_date)),"") selling_date,
	` + selectRangeStr + `
FROM
	m_stock a
LEFT OUTER JOIN m_jan b
	ON a.jan_code = b.jan_code
LEFT OUTER JOIN m_shop c
	ON a.shop_code = c.shop_code
LEFT OUTER JOIN m_initial_sales d
	ON a.shop_code = d.shop_code
	AND a.jan_code = d.jan_code
WHERE
	a.shop_code IN (?` + strings.Repeat(",?", len(form.ShopCd)-1) + `)
	AND a.jan_code LIKE '4910` + form.MagazineCdSingle + `%'
GROUP BY
	a.jan_code
ORDER BY
	-- jan_seq_number DESC
	magazine_seq_number DESC
`
	sqlCache = sql + sqlCache
	for _, s := range form.ShopCd {
		sqlCache = sqlCache + s
	}
	sqlCache = sqlCache + form.MagazineCdSingle
	return
}

//Check condition and export query
// Data detail by JAN
func buildDetail(form QueryForm, ctx *gf.Context) (sql string, sqlCache string, listHeader []string) {

	listHeader = append(
		listHeader,
		"店舗名",
		"売上累計",
		"入荷累計",
		"在庫数",
		"初売上日",
	)

	//=================================================================================================================
	// range data from 売上初速01 -> 売上初速40
	selectRangeArray := []string{}
	for index := 1; index <= form.SearchDateType; index++ {
		listHeader = append(listHeader, strconv.Itoa(index)+"日目")
		item := "MAX(d.sales_quantity_day"
		if index < 10 {
			item += "0" + strconv.Itoa(index)
		} else {
			item += strconv.Itoa(index)
		}
		item += ") sales_quantity_day" + strconv.Itoa(index)
		selectRangeArray = append(selectRangeArray, item)
	}
	selectRangeStr := strings.Join(selectRangeArray, ",")
	//=================================================================================================================
	sql = `
SELECT
	MAX(a.shop_code) shop_code,
	MAX(c.shop_name) shop_name,
	SUM(a.cumulative_sales_quantity) cumulative_sales_quantity,
	SUM(a.cumulative_receiving_quantity) cumulative_receiving_quantity,
	SUM(a.stock_quantity) stock_quantity,
	MAX(c.shop_seq_number) shop_seq_number,
	MAX(d.first_sales_date) first_sales_date,
	` + selectRangeStr + `
FROM
	m_stock a
LEFT OUTER JOIN m_shop c
	ON a.shop_code = c.shop_code
LEFT OUTER JOIN m_initial_sales d
	ON a.shop_code = d.shop_code
	AND a.jan_code = d.jan_code
WHERE
	a.shop_code IN (?` + strings.Repeat(",?", len(form.ShopCd)-1) + `)
	AND a.jan_code = '` + form.JanKey + `'
GROUP BY
	a.jan_code,
	a.shop_code
ORDER BY
	shop_seq_number
`
	sqlCache = sql + sqlCache
	for _, s := range form.ShopCd {
		sqlCache = sqlCache + s
	}
	sqlCache = sqlCache + form.JanKey
	return
}

// Check condition and export query
// Data detail by JAN (export for CSV)
func buildDetailForCSV(form QueryForm, ctx *gf.Context, listJanKey []string) (sql string, sqlCache string, listHeader []string) {

	listHeader = append(
		listHeader,
		"店舗名",
		"売上累計",
		"入荷累計",
		"在庫数",
		"初売上日",
	)
	//CONDITION
	sortCondition := ""
	sortSelected := ""
	if form.ControlType == CONTROL_TYPE_JAN {
		sortSelected += `
			IFNULL(MIN(IF(TRIM(b.selling_date) = "",NULL, b.selling_date)),"") selling_date,
		`
		sortCondition += " selling_date DESC,"
	} else {
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT START
		//sortSelected += " MAX(c.shop_seq_number) shop_seq_number,"
		//sortCondition += " jan_seq_number DESC,"
		sortSelected += " MAX(b.magazine_seq_number) magazine_seq_number,"
		sortCondition += " magazine_seq_number DESC,"
		// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT END
	}

	//=================================================================================================================
	// range data from 売上初速01 -> 売上初速40
	selectRangeArray := []string{}
	for index := 1; index <= form.SearchDateType; index++ {
		listHeader = append(listHeader, strconv.Itoa(index)+"日目")
		item := "MAX(d.sales_quantity_day"
		if index < 10 {
			item += "0" + strconv.Itoa(index)
		} else {
			item += strconv.Itoa(index)
		}
		item += ") sales_quantity_day" + strconv.Itoa(index)
		selectRangeArray = append(selectRangeArray, item)
	}
	selectRangeStr := strings.Join(selectRangeArray, ",")
	//=================================================================================================================
	sql = `
SELECT
	a.jan_code jan_code,
	MAX(a.shop_code) shop_code,
	MAX(c.shop_name) shop_name,
	SUM(a.cumulative_sales_quantity) cumulative_sales_quantity,
	SUM(a.cumulative_receiving_quantity) cumulative_receiving_quantity,
	SUM(a.stock_quantity) stock_quantity,
	MAX(c.shop_seq_number) shop_seq_number,
	MAX(d.first_sales_date) first_sales_date,
	` + sortSelected + selectRangeStr + `
FROM
	m_stock a
LEFT OUTER JOIN m_jan b
	ON a.jan_code = b.jan_code
LEFT OUTER JOIN m_shop c
	ON a.shop_code = c.shop_code
LEFT OUTER JOIN m_initial_sales d
	ON a.shop_code = d.shop_code
	AND a.jan_code = d.jan_code
WHERE
	a.shop_code IN (?` + strings.Repeat(",?", len(form.ShopCd)-1) + `)
	AND a.jan_code IN (?` + strings.Repeat(",?", len(listJanKey)-1) + `)
GROUP BY
	a.jan_code,
	a.shop_code
ORDER BY
	` + sortCondition + `
	a.jan_code,
	shop_seq_number
`
	sqlCache = sql + sqlCache
	for _, s := range form.ShopCd {
		sqlCache = sqlCache + s
	}
	for _, s := range listJanKey {
		sqlCache = sqlCache + s
	}
	return
}
