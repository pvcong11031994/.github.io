package RP066_BestSales_Cloud

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"strconv"
	"strings"
	"time"
)

//Check condition and export query
func buildSql(form QueryForm, ctx *gf.Context) (sql string, colField [][]string, headerCol []string, exCols []string) {

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 日付 年月日
	dateSearchFrom := ""
	dateSearchTo := ""

	yearFrom := Common.ConvertIntToString(time.Now().Year())
	yearTo := Common.ConvertIntToString(time.Now().Year())
	monthFrom := time.Now().Format("01")
	monthTo := time.Now().Format("01")
	dayFrom := Common.ConvertIntToString(time.Now().Day())
	dayTo := Common.ConvertIntToString(time.Now().Day())

	if form.DateFrom != "" || form.DateTo != "" {
		if form.DateFrom != "" {
			dateSearchFrom = strings.Replace(form.DateFrom, "/", "", -1)
		}
		if form.DateTo != "" {
			dateSearchTo = strings.Replace(form.DateTo, "/", "", -1)
		}
	}

	if dateSearchFrom == "" {
		dateSearchFrom = yearFrom + monthFrom + dayFrom
	}
	if dateSearchTo == "" {
		dateSearchTo = yearTo + monthTo + dayTo
	}

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// get Date =================================
	colField = [][]string{}
	sqlRange := ""
	mcmd := Models.MasterCalendarModel{ctx.DB}
	listRange := []ModelItems.MasterCalendarItem{}
	stringsCol, _, _ := initDefaultLayout()
	var err error
	sqlWithMasterCalendar := ""
	sqlJoinMasterCalendar := ""
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION {
			switch form.GroupType {
			case GROUP_TYPE_DATE:
				listRange, err = mcmd.GetDay(dateSearchFrom, dateSearchTo)
				Common.LogErr(err)
				listDay := []string{}
				listMonth := []string{}
				listYear := []string{}
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
					listDay = append(listDay, item.Mcdd)
					listMonth = append(listMonth, item.Mcmm)
					listYear = append(listYear, item.Mcyyyy)
					exCols = append(exCols, item.Mcyyyy + item.Mcmm + item.Mcdd)
				}
				colField = append(colField, listYear, listMonth, listDay)
				headerCol = []string{stringsCol[mc_yyyy], stringsCol[mc_mm], stringsCol[mc_dd]}
			case GROUP_TYPE_WEEK:
				listRange, err = mcmd.GetWeek(dateSearchFrom, dateSearchTo)
				Common.LogErr(err)
				listWeek := []string{}
				listYear := []string{}
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
					listWeek = append(listWeek, item.Mcweekdate)
					listYear = append(listYear, item.Mcyyyy)
					exCols = append(exCols, item.Mcyyyy + item.Mcweeknum)
				}
				colField = append(colField, listYear, listWeek)
				headerCol = []string{stringsCol[mc_yyyy], stringsCol[mc_weekdate]}
			case GROUP_TYPE_MONTH:
				listRange, err = mcmd.GetMonth(dateSearchFrom, dateSearchTo)
				Common.LogErr(err)
				listMonth := []string{}
				listYear := []string{}
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
					listMonth = append(listMonth, item.Mcmm)
					listYear = append(listYear, item.Mcyyyy)
					exCols = append(exCols, item.Mcyyyy + item.Mcmm)
				}
				colField = append(colField, listYear, listMonth)
				headerCol = []string{stringsCol[mc_yyyy], stringsCol[mc_mm]}
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
	JOIN mc
		ON mc.mc_yyyymmdd = SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8)
		AND mc.mc_yyyy = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8),1,4)
		AND mc.mc_mm = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8),5,2)
		AND mc.mc_dd = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),' ',''),0,8),7,2)
	`		//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
		}

		// Add query when フォーマット choose 集計結果+店舗
		if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT_STORE {
			listShopName := []string{}
			for i, item := range form.ShopCd {
				sqlRange += `
			,SUM(
				CASE
				WHEN
					bqsl.shop_code = '` + item + `'
					AND SUBSTR(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),0,8) <= {{date_to}}
					AND bqsl.bookstore_biz_category = '40'
					THEN SAFE_CAST((IFNULL(bqsl.sales_body_quantity,0)) AS INT64)
				   ELSE 0
				END
			) as ` + `A`+item
				listShopName = append(listShopName, form.ShopName[i])
				exCols = append(exCols, item)
			}
			headerCol = listShopName
		}
		//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	}

	parameter := map[string]interface{}{}
	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)

	// condition control type
	conditionSelect := ""
	if form.ControlType == CONTROL_TYPE_BOOK {
		// ASO-5754 [BA]mBAWEB-v16a 売上ベスト：商品名などをCloudSQLから取得 - DEL START
		//conditionSelect = `MAX(author_name) author_name,`
		// ASO-5754 [BA]mBAWEB-v16a 売上ベスト：商品名などをCloudSQLから取得 - DEL END
	} else if form.ControlType == CONTROL_TYPE_MAGAZINE {
		conditionSelect = `MAX(SUBSTR(bqsl.jan_code, 5, 7)) magazine_code,`
	}

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 商品区分
	conditionBqslMagazine := ""
	if form.MagazineCodeWeek == BQSL_MAGAZINE_CODE_MONTH && form.ControlType == CONTROL_TYPE_MAGAZINE {
		if form.MagazineCodeMonth == BQSL_MAGAZINE_CODE_WEEK && form.MagazineCodeQuarter == BQSL_MAGAZINE_CODE_QUARTER {
			conditionBqslMagazine += `
			`
		} else if form.MagazineCodeMonth == BQSL_MAGAZINE_CODE_WEEK && form.MagazineCodeQuarter != BQSL_MAGAZINE_CODE_QUARTER {
			conditionBqslMagazine += `
		AND (magazine_code LIKE '0%' OR magazine_code LIKE '1%' OR magazine_code LIKE '2%' OR magazine_code LIKE '3%')
			`
		} else if form.MagazineCodeMonth != BQSL_MAGAZINE_CODE_WEEK && form.MagazineCodeQuarter == BQSL_MAGAZINE_CODE_QUARTER {
			conditionBqslMagazine += `
		AND (magazine_code LIKE '0%' OR magazine_code LIKE '1%')
		AND magazine_code NOT LIKE '2%'
		AND magazine_code NOT LIKE '3%'
			`
		} else {
			conditionBqslMagazine += `
		AND (magazine_code LIKE '0%' OR magazine_code LIKE '1%')
			`
		}
	} else if form.MagazineCodeMonth == BQSL_MAGAZINE_CODE_WEEK {
		if form.MagazineCodeQuarter == BQSL_MAGAZINE_CODE_QUARTER {
			conditionBqslMagazine += `
		AND (magazine_code LIKE '2%' OR magazine_code LIKE '3%')
		AND magazine_code NOT LIKE '0%'
		AND magazine_code NOT LIKE '1%'
			`
		} else {
			conditionBqslMagazine += `
		AND (magazine_code LIKE '2%' OR magazine_code LIKE '3%')
			`
		}
	} else if form.MagazineCodeQuarter == BQSL_MAGAZINE_CODE_QUARTER {
		conditionBqslMagazine += `
		AND (magazine_code NOT LIKE '0%'
		AND magazine_code NOT LIKE '1%'
		AND magazine_code NOT LIKE '2%'
		AND magazine_code NOT LIKE '3%')
		`
	}
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	conditionBqmg := ""
	// 特大分類
	// 大分類
	// 中分類
	// 小分類
	if len(form.MediaGroup1Cd) > 0 {
		conditionBqmg += `
    	AND bqsl.super_large_grouping_code IN {{bqsl_media_group1_cd}}
     `
		parameter["bqsl_media_group1_cd"] = form.MediaGroup1Cd
	}
	if len(form.MediaGroup2Cd) > 0 {
		conditionBqmg += `
    	AND bqsl.large_grouping_code  IN {{bqsl_media_group2_cd}}
    `
		parameter["bqsl_media_group2_cd"] = form.MediaGroup2Cd
	}
	if len(form.MediaGroup3Cd) > 0 {
		conditionBqmg += `
    	AND bqsl.middle_grouping_code IN {{bqsl_media_group3_cd}}
    `
		parameter["bqsl_media_group3_cd"] = form.MediaGroup3Cd

	}
	if len(form.MediaGroup4Cd) > 0 {
		conditionBqmg += `
    	AND bqsl.small_grouping_code IN {{bqsl_media_group4_cd}}
    `
		parameter["bqsl_media_group4_cd"] = form.MediaGroup4Cd
	}

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// メーカー (出版者記号)
	conditionJanMakerCode := ""
	if len(form.JanMakerCode) > 0 {
		janMakerCodeConditionArr := []string{}
		for i, v := range form.JanMakerCode {
			janMakerCode := "Jan_maker_code_" + strconv.Itoa(i)
			janMakerCodeConditionArr = append(janMakerCodeConditionArr, ` jan_maker_code = {{`+janMakerCode+`}} `)
			parameter[janMakerCode] = JAN_MAKER_CODE + v
		}
		conditionJanMakerCode += `
   	 	AND (` + strings.Join(janMakerCodeConditionArr, " OR ") + `)`
	}
	//=================================================================================================================

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 出版社コード
	conditionMakerCd := ""
	if len(form.MakerCd) > 0 {
		conditionMakerCd += " AND maker_code IN {{maker_list}}"
		parameter["maker_list"] = form.MakerCd
	}
	//=================================================================================================================

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 雑誌コード
	conditionMagazineCd := ""
	if len(form.MagazineCd) > 0 {
		magazineCdConditionArr := []string{}
		for i, v := range form.MagazineCd {
			magazineCd := "magazine_cd_" + strconv.Itoa(i)
			magazineCdConditionArr = append(magazineCdConditionArr, `magazine_code LIKE {{`+magazineCd+`}} `)
			if !strings.Contains(v, "%") {
				v = "%" + v + "%"
			}
			parameter[magazineCd] = v
		}
		conditionMagazineCd += `
    	AND (` + strings.Join(magazineCdConditionArr, " OR ") + `)`
	}
	//=================================================================================================================

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 件数
	conditionLimitCSV := ""
	if form.SearchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		conditionLimitCSV = `LIMIT ` + strconv.Itoa(form.Limit)
	}
	//=================================================================================================================

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// Tab
	conditionTab := ""
	if form.ControlType == CONTROL_TYPE_BOOK {
		conditionTab += `
		AND jan_grouping = '1'`
	} else if form.ControlType == CONTROL_TYPE_MAGAZINE {
		conditionTab += `
		AND jan_grouping = '2'`
	}

	//=================================================================================================================

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// condition shop
	parameter["shop_cd"] = form.ShopCd
	conditionDate := ""
	if dateSearchFrom != "" {
		conditionDate += `AND SUBSTR(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),0,8) >= {{date_from}}
		`
	}
	if dateSearchTo != "" {
		conditionDate += `AND SUBSTR(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),0,8) <= {{date_to}}
		`
	}
	parameter["date_from"] = dateSearchFrom
	parameter["date_to"] = dateSearchTo
	parameter["month_from"] = dateSearchFrom[:6]
	parameter["month_to"] = dateSearchTo[:6]
	//=================================================================================================================
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	sqlWithData := bq.NewCommand()
	sqlWithData.Parameters = parameter
	sqlWithData.CommandText = `
#StandardSQL
	WITH
	bq_sales_condition AS (
		SELECT
			jan_code,
			shop_code,
			sales_datetime,
			bookstore_biz_category,
			sales_body_quantity,
			super_large_grouping_code,
			large_grouping_code,
			middle_grouping_code,
			small_grouping_code
		FROM ` + "`{{@DATASET}}.bq_sales_*`" + ` bqsl
		WHERE
			SUBSTR(_TABLE_SUFFIX,-14,6) BETWEEN {{month_from}} AND {{month_to}}
			AND bqsl.shop_code IN {{shop_cd}}
			` + conditionBqslMagazine + conditionBqmg + conditionJanMakerCode + conditionMakerCd + conditionMagazineCd + conditionTab + `
			AND SUBSTR(REPLACE(REPLACE(sales_datetime, '-', ''),'/',''),0,8) >= {{date_from}}
			AND SUBSTR(REPLACE(REPLACE(sales_datetime, '-', ''),'/',''),0,8) <= {{date_to}}
	)
	--,
	--bq_sales_max_date AS (
	--	SELECT
	--		shop_code,
	--		jan_code,
	--		MAX(received_datetime) received_datetime
	--	FROM bq_sales_condition
	--	GROUP BY shop_code, jan_code
	--)
	` + sqlWithMasterCalendar + `

SELECT *
	FROM (
	SELECT
		bqsl.jan_code AS jan_code,
		-- MAX(product_name) product_name,
		` + conditionSelect + `
		-- IF(MAX(jan_maker_name) <> '' ,MAX(jan_maker_name),MAX(maker_name)) maker_name,
		-- MAX(IFNULL(SUBSTR(REPLACE(bqsl.selling_date, '-', ''),0,8),'')) selling_date,
		-- MAX(SAFE_CAST(sales_tax_exc_unit_price AS INT64)) sales_tax_exc_unit_price,
		--SUM(SAFE_CAST(IF(bqsl.received_datetime < bqsl_stock.received_datetime,0, bqsl.cumulative_receiving_quantity) AS INT64)) AS stok_cumulative_receiving_quantity,
		--SUM(SAFE_CAST(IF(bqsl.received_datetime < bqsl_stock.received_datetime,0, bqsl.cumulative_sales_quantity) AS INT64)) AS stok_cumulative_sales_quantity,
		--SUM(SAFE_CAST(IF(bqsl.received_datetime < bqsl_stock.received_datetime,0, bqsl.stock_quantity) AS INT64)) AS stok_stock_quantity,
		--IFNULL(MIN(IF(TRIM(first_sales_date) = "",NULL, first_sales_date)),"") first_sales_date,
		--SUM(IF(SUBSTR(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),0,8) <= {{date_to}}, SAFE_CAST((IFNULL(sales_body_quantity,0)) AS INT64), 0)) sales_body_quantity
		--SUM(IF(SUBSTR(REPLACE(REPLACE(bqsl.sales_datetime, '-', ''),'/',''),0,8) <= {{date_to}} AND bqsl.bookstore_biz_category = '40', SAFE_CAST((IFNULL(sales_body_quantity,0)) AS INT64), 0)) sales_body_quantity
		SUM(IF(bqsl.bookstore_biz_category = '40', SAFE_CAST((IFNULL(sales_body_quantity,0)) AS INT64), 0)) sales_body_quantity
		` + sqlRange + `
	FROM bq_sales_condition bqsl
	--JOIN bq_sales_max_date bqsl_stock
	--	ON bqsl.shop_code = bqsl_stock.shop_code
	--	AND bqsl.jan_code = bqsl_stock.jan_code
	` + sqlJoinMasterCalendar + `
	GROUP BY
		jan_code)
ORDER BY
	sales_body_quantity DESC
` + conditionLimitCSV
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sql, err = sqlWithData.Build()
	Common.LogErr(err)
	return
}
