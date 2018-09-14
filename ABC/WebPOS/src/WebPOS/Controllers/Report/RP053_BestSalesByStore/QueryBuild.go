package RP053_BestSalesByStore

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/WebApp"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"strconv"
	"strings"
	"time"
)

func buildSql(form QueryForm, ctx *gf.Context) (sql string, rowField []string, colField []string, sumField []string) {

	switch form.GroupType {
	case GROUP_TYPE_DATE:
		form.LayoutColArr = []string{mc_yyyy, mc_mm, mc_dd}
	case GROUP_TYPE_WEEK:
		form.LayoutColArr = []string{mc_yyyy, mc_weekdate}
	case GROUP_TYPE_MONTH:
		form.LayoutColArr = []string{mc_yyyy, mc_mm}
	}

	sumField = []string{}
	sumFieldSelect := []string{}
	for _, s := range form.LayoutSumArr {
		switch s {
		case bqio_goods_count:
			sumField = append(sumField, s)
			sumFieldSelect = append(sumFieldSelect, "SUM(SAFE_CAST((IFNULL("+s+",0)) AS INT64)) "+s)
		}
	}

	colField = []string{}
	colFieldSelect := []string{}
	colFieldBlank := []string{}
	for _, s := range form.LayoutColArr {
		switch s {
		case mc_yyyy,
			mc_mm,
			mc_weekdate,
			mc_dd:
			colField = append(colField, s)
			colFieldSelect = append(colFieldSelect, "MAX(IFNULL("+s+",'')) "+s)
			colFieldBlank = append(colFieldBlank, `'' `+s)
		}
	}

	orderRowFields := []string{}
	for i, v := range sumField {
		if i > 3 {
			break
		}
		orderRowFields = append(orderRowFields, v+" DESC")
	}

	rowField = []string{}
	rowFieldSelect := []string{}
	rowFieldBlank := []string{}
	for _, s := range form.LayoutRowArr {
		switch s {
		case bqio_jan_cd,
			goods_name,
			writer_name,
			publisher_name,
			bqgm_price,
			rank_no,
			stock_count,
			total_arrival,
			total_sales:
			rowField = append(rowField, s)
			if s == bqgm_price ||
				s == rank_no ||
				s == total_arrival ||
				s == total_sales ||
				s == stock_count {
				rowFieldBlank = append(rowFieldBlank, `SAFE_CAST(0 AS INT64) `+s)
			} else {
				rowFieldBlank = append(rowFieldBlank, `'' `+s)
			}

			if s == rank_no {
				rowFieldSelect = append(rowFieldSelect, "ROW_NUMBER() OVER (ORDER BY SUM(SAFE_CAST((IFNULL(bqio_goods_count,0)) AS INT64)) DESC) "+s)
				orderRowFields = []string{rank_no}
			} else if s == stock_count {
				rowFieldSelect = append(rowFieldSelect, "MAX(IFNULL(bqsc_stock_count,SAFE_CAST(0 AS INT64))) "+s)
			} else if s == total_arrival {
				rowFieldSelect = append(rowFieldSelect, "MAX(IFNULL(total_arrival,SAFE_CAST(0 AS INT64))) "+s)
			} else if s == total_sales {
				rowFieldSelect = append(rowFieldSelect, "MAX(IFNULL(total_sales,SAFE_CAST(0 AS INT64))) "+s)
			} else {
				rowFieldSelect = append(rowFieldSelect, "MAX("+s+") "+s)
			}
		}
	}

	//================================================================================================================================================
	//================================================================================================================================================
	joinWithBqCategory := ""
	if len(form.MediaGroup1Cd) > 0 ||
		len(form.MediaGroup2Cd) > 0 ||
		len(form.MediaGroup3Cd) > 0 {
		joinWithBqCategory += `
LEFT JOIN {{@DATASET}}.bq_category_ms bccm
	ON bccm.bccm_media_cd = bqio.bqio_media_cd
`
	}
	//================================================================================================================================================
	//================================================================================================================================================
	colKey := ""
	groupCol := ""
	orderCol := ""
	arrOrderCol := []string{}
	compareColNotGet := []string{}
	for _, s := range colField {
		switch s {
		case mc_yyyy,
			mc_mm,
			mc_weekdate,
			mc_dd:
			arrOrderCol = append(arrOrderCol, "IFNULL("+s+", '')")
			compareColNotGet = append(compareColNotGet, "")
		}
		orderCol += s + ","
	}
	if len(arrOrderCol) == 0 {
		colKey = RPComon.NO_KEY
	} else {
		colKey = `CAST(FARM_FINGERPRINT(CONCAT(` + strings.Join(arrOrderCol, ",'|',") + `)) AS STRING)`
	}

	groupCol = "col_key "

	rowKey := bqio_jan_cd
	orderRow := ""
	orderRow = strings.Join(orderRowFields, ", ") + ","

	parameter := map[string]interface{}{}
	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sqlHeaderCol := ""
	isCalendar := false
	selectCalendar := ""
	whereCalendar := " ON mc.mc_yyyymmdd = SUBSTR(REPLACE(REPLACE(REPLACE(bqio.bqio_trn_date, '-', ''),'/',''),' ',''),0,8) "
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	for _, s := range form.LayoutColArr {
		switch s {
		case mc_yyyy:
			isCalendar = true
			whereCalendar += " AND mc.mc_yyyy = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqio.bqio_trn_date, '-', ''),'/',''),' ',''),0,8),1,4) "
			selectCalendar = selectCalendar + s + ","
		case mc_mm:
			isCalendar = true
			whereCalendar += " AND mc.mc_mm = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqio.bqio_trn_date, '-', ''),'/',''),' ',''),0,8),5,2) "
			selectCalendar = selectCalendar + s + ","
		case mc_weekdate:
			isCalendar = true
			whereCalendar += " AND mc.mc_mm = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqio.bqio_trn_date, '-', ''),'/',''),' ',''),0,8),5,2) "
			whereCalendar += " AND mc.mc_dd = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqio.bqio_trn_date, '-', ''),'/',''),' ',''),0,8),7,2) "
			selectCalendar = selectCalendar + "mc_dd,mc_mm,mc_weekdate,"
		case mc_dd:
			isCalendar = true
			whereCalendar += " AND mc.mc_dd = SUBSTR(SUBSTR(REPLACE(REPLACE(REPLACE(bqio.bqio_trn_date, '-', ''),'/',''),' ',''),0,8),7,2) "
			selectCalendar = selectCalendar + s + ","
		default:
			// not thing
			continue
		}
	}

	//================================================================================================================================================
	//================================================================================================================================================
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// condition
	// 店舗
	conditionBqshm := ""
	if len(form.ShopCd) > 0 {
		conditionBqshm += `
	AND CONCAT(bqio.bqio_servername , '|' , bqio.bqio_shop_cd) = {{shop_cd}}    `
		parameter["shop_cd"] = form.ShopCd
	}

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

	conditionBqcl := ""
	conditionBqio := ""
	if dateSearchFrom != "" {
		if isCalendar {
			conditionBqcl += " AND  CONCAT(mc.mc_yyyy , mc.mc_mm , mc.mc_dd) >= {{calendar_from}} "
		}
		conditionBqio += " AND SUBSTR(REPLACE(REPLACE(bqio.bqio_trn_date, '-', ''),'/',''),0,8) >= {{trn_date_from}} "
		parameter["calendar_from"] = dateSearchFrom
		parameter["trn_date_from"] = dateSearchFrom
	}
	if dateSearchTo != "" {
		if isCalendar {
			conditionBqcl += " AND  CONCAT(mc.mc_yyyy , mc.mc_mm , mc.mc_dd) <= {{calendar_to}} "
		}
		conditionBqio += " AND SUBSTR(REPLACE(REPLACE(bqio.bqio_trn_date, '-', ''),'/',''),0,8) <= {{trn_date_to}} "
		parameter["calendar_to"] = dateSearchTo
		parameter["trn_date_to"] = dateSearchTo
	}
	conditionBqio += `
		AND bqio_jan_cd NOT LIKE '99%'
`
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	conditionBqmg := ""
	// メディア大分類
	// メディア中分類
	// メディア中小分類
	if len(form.MediaGroup3Cd) > 0 {
		conditionBqmg += `
	AND SUBSTR(bqio.bqio_media_cd,0, 6) IN {{bqio_media_group3_cd}}
`
		parameter["bqio_media_group3_cd"] = form.MediaGroup3Cd
	} else if len(form.MediaGroup2Cd) > 0 {
		conditionBqmg += `
	AND SUBSTR(bqio.bqio_media_cd,0, 4) IN {{bqio_media_group2_cd}}
`
		parameter["bqio_media_group2_cd"] = form.MediaGroup2Cd
	} else if len(form.MediaGroup1Cd) > 0 {
		conditionBqmg += `
	AND SUBSTR(bqio.bqio_media_cd,0, 2) IN {{bqio_media_group1_cd}}
`
		parameter["bqio_media_group1_cd"] = form.MediaGroup1Cd
	}
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	conditionBqmk := ""
	// メーカー
	if len(form.MakerCd) > 0 {
		makerCdConditionArr := []string{}
		for i, v := range form.MakerCd {
			makerCd := "maker_cd_" + strconv.Itoa(i)
			makerCdConditionArr = append(makerCdConditionArr, ` publisher_name LIKE {{`+makerCd+`}} `)
			if !strings.Contains(v, "%") {
				v = "%" + v + "%"
			}
			parameter[makerCd] = v
		}
		conditionBqmk += `
	AND (` + strings.Join(makerCdConditionArr, " OR ") + `)`
	}
	//================================================================================================================================================
	//================================================================================================================================================
	//サーバー名
	conditionServer := `
	AND (_TABLE_SUFFIX LIKE {{server_name}})`

	whereSC := `
	AND SUBSTR(_TABLE_SUFFIX,-14,6) BETWEEN {{month_from}} AND {{month_to}}
	AND SUBSTR(_TABLE_SUFFIX,-8) = '00000000'
	AND _TABLE_SUFFIX LIKE '%_1_%'
	`
	sqlSC := " `{{@DATASET}}.bq_inout_*`" + ` bqio`
	if isCalendar {
		sqlSC = `
	{{@DATASET}}.master_calendar mc
LEFT JOIN ` + "`{{@DATASET}}.bq_inout_*`" + ` bqio
	` + whereCalendar + whereSC + conditionServer + conditionBqio + conditionBqshm
		conditionServer = ""
		whereSC = ""
		conditionBqio = ""
		conditionBqshm = ""
	}
	dateSearchFrom = strings.Replace(strings.Replace(dateSearchFrom, "-", "", -1), "/", "", -1)
	dateSearchTo = strings.Replace(strings.Replace(dateSearchTo, "-", "", -1), "/", "", -1)

	//parameter["server_name"] = "%" + strings.Split(form.ShopCd, "|")[0] + "%"
	arrShopCd := strings.Split(form.ShopCd, "|")
	serverName := ""
	if len(arrShopCd) > 0 {
		serverName = arrShopCd[0]
	}
	parameter["server_name"] = "%" + serverName + "%"
	parameter["date_from"] = dateSearchFrom
	parameter["date_to"] = dateSearchTo
	parameter["month_from"] = dateSearchFrom[:6]
	parameter["month_to"] = dateSearchTo[:6]

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sqlInOutGoods := bq.NewCommand()
	sqlInOutGoods.Parameters = parameter
	sqlInOutGoods.CommandText = `
SELECT
` + selectCalendar + `
	bqio.bqio_shop_cd AS bqio_shop_cd,
	bqio.bqio_servername AS bqio_servername,
	shm.shm_server_name AS shm_server_name,
	shm.shm_shop_cd AS shm_shop_cd,
	shm.shm_shop_name AS shm_shop_name,
	shm.shm_product_cd AS shm_product_cd,
	bqio.bqio_jan_cd AS bqio_jan_cd,
	bqio.bqio_media_cd AS bqio_media_cd,
	SUBSTR(bqio.bqio_media_cd,0,2) AS bqio_media_group1_cd,
	SUBSTR(bqio.bqio_media_cd,0,2) AS bqio_media_group2_cd,
	SUBSTR(bqio.bqio_media_cd,0,4) AS bqio_media_group3_cd,
	bqio.bqio_media_cd AS bqio_media_group4_cd,
	bqgm.bqgm_goods_type AS  goods_type,
	SUBSTR(bqgm.bqgm_media_cd,0, 2) AS bqgm_media_group1_cd,
	SUBSTR(bqgm.bqgm_media_cd,0, 4) AS bqgm_media_group2_cd,
	SUBSTR(bqgm.bqgm_media_cd,0, 6) AS bqgm_media_group3_cd,
	bqgm.bqgm_SYOZAICODE2 AS bqgm_SYOZAICODE2,
	bqgm.bqgm_media_cd AS bqgm_media_group4_cd,
	bqgm.bqgm_category_cd AS category_cd,
	bqgm.bqgm_standard_number AS standard_number,
	bqgm.bqgm_goods_name AS goods_name,
	bqstc.bqsc_cum_receiving_quantity total_arrival,
	bqstc.bqsc_cum_sales_quantity total_sales,
	IF(bqgm.bqgm_goods_type IN('20','21'),
		bqgm.bqgm_writer_name,
		bqgm.bqgm_artist_name) AS writer_name,
	IF(bqgm.bqgm_goods_type IN('20','21'),
		bqgm.bqgm_publisher_name,
		IFNULL(bqmk.bqmk_maker_name,'')) AS publisher_name,
	bqgm.bqgm_sales_date AS sales_date,
	SAFE_CAST(bqgm.bqgm_price AS INT64) AS bqgm_price ,
	bqio.bqio_goods_count AS bqio_goods_count,
	bqstc.bqsc_stock_count AS bqsc_stock_count,
	bqgm.bqgm_isbn AS bqgm_isbn,
	bqgm.bqgm_goods_type AS bqgm_goods_type,
	bqgm.bqgm_kihon_magazinecode AS bqgm_kihon_magazinecode
FROM ` + sqlSC + `
LEFT JOIN ` + "`{{@DATASET}}.bq_stock_*`" + ` bqstc
	ON bqio.bqio_shop_cd = bqstc.bqsc_shop_cd
	AND bqio.bqio_servername = bqstc.bqsc_servername
	AND bqio.bqio_jan_cd = bqstc.bqsc_jan_cd
	AND (bqstc._TABLE_SUFFIX LIKE {{server_name}})
------------------------------------------------------------------
LEFT JOIN {{@DATASET}}.shop_master shm
	ON shm.shm_server_name = bqio.bqio_servername
	AND shm.shm_shop_cd = bqio.bqio_shop_cd
------------------------------------------------------------------
LEFT JOIN {{@DATASET}}.bq_goods_master bqgm
	ON bqio.bqio_jan_cd = bqgm.bqgm_jan_cd
------------------------------------------------------------------
LEFT JOIN {{@DATASET}}.bq_maker bqmk
	ON bqgm.bqgm_label_cd = bqmk.bqmk_label_cd
	AND bqio.bqio_servername = bqmk.bqmk_servername
WHERE
	TRUE
	` + conditionBqcl + conditionBqshm + conditionServer + conditionBqio + whereSC

	parameter["#sql_io_goods_media"] = sqlInOutGoods
	//parameter["server_name"] = "%" + strings.Split(form.ShopCd, "|")[0] + "%"
	parameter["server_name"] = "%" + serverName + "%"

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	if len(colField) > 0 {
		cols := ""
		if len(colFieldSelect) > 0 {
			cols = strings.Join(colFieldSelect, ",\r\n\t") + `,`
		}
		rows := ""
		if len(rowFieldBlank) > 0 {
			rows = strings.Join(rowFieldBlank, ",\r\n\t") + `,`
		}
		sqlHeaderColTemplate := `
SELECT
	1 group_code,
	` + colKey + ` col_key,
	` + RPComon.SUM_KEY + ` row_key,
	` + cols + `
	` + rows + `
	` + strings.Join(sumFieldSelect, ",\r\n\t") + `
FROM ` + ` ({{#sql_io_goods_media}}) bqio ` + `
` + joinWithBqCategory + `
WHERE TRUE
GROUP BY
    ` + groupCol + `
`
		cmdHeaderCol := bq.NewCommand()
		cmdHeaderCol.CommandText = sqlHeaderColTemplate
		cmdHeaderCol.Parameters = parameter
		var err error
		sqlHeaderCol, err = cmdHeaderCol.Build()
		Common.LogErr(exterror.WrapExtError(err))
	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sqlHeaderRow := ``
	if len(rowField) > 0 {
		cols := ""
		if len(colFieldBlank) > 0 {
			cols = strings.Join(colFieldBlank, ",\r\n\t") + `,`
		}
		rows := ""
		if len(rowFieldSelect) > 0 {
			rows = strings.Join(rowFieldSelect, ",\r\n\t") + `,`
		}
		sqlHeaderRowTemplate := `
SELECT
	2 group_code,
	` + RPComon.SUM_KEY + ` col_key,
	` + rowKey + ` row_key,
	` + cols + `
	` + rows + `
	` + strings.Join(sumFieldSelect, ",\r\n\t") + `
FROM ` + ` ({{#sql_io_goods_media}}) bqio ` + `
` + joinWithBqCategory + `
WHERE TRUE ` + conditionBqmg + `
	AND bqio_jan_cd <> ''
GROUP BY row_key
HAVING TRUE ` + conditionBqmk
		cmdHeaderRow := bq.NewCommand()
		cmdHeaderRow.CommandText = sqlHeaderRowTemplate
		cmdHeaderRow.Parameters = parameter
		var err error
		sqlHeaderRow, err = cmdHeaderRow.Build()
		Common.LogErr(exterror.WrapExtError(err))
	}
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sqlData := ""

	if len(colField) > 0 && len(rowField) > 0 {
		cols := ""
		if len(colFieldBlank) > 0 {
			cols = strings.Join(colFieldBlank, ",\r\n\t") + `,`
		}
		rows := ""
		if len(rowFieldBlank) > 0 {
			rows = strings.Join(rowFieldBlank, ",\r\n\t") + `,`
		}

		sqlDataTemplate := `
SELECT
	3 group_code,
	` + colKey + ` col_key,
	` + rowKey + ` row_key,
	` + cols + `
	` + rows + `
	` + strings.Join(sumFieldSelect, ",\r\n\t") + `
FROM ` + ` ({{#sql_io_goods_media}}) bqio ` + `
` + joinWithBqCategory + `
WHERE TRUE ` + `
GROUP BY col_key, row_key
HAVING
	row_key IS NOT NULL
	AND col_key IS NOT NULL
`
		cmdData := bq.NewCommand()
		cmdData.CommandText = sqlDataTemplate
		cmdData.Parameters = parameter
		var err error
		sqlData, err = cmdData.Build()
		Common.LogErr(exterror.WrapExtError(err))
	}

	if sqlHeaderCol != "" {
		sqlHeaderCol = "(" + sqlHeaderCol + ") UNION ALL"
	}
	if sqlHeaderRow != "" {
		sqlHeaderRow = "(" + sqlHeaderRow + ") UNION ALL"
	}
	if sqlData != "" {
		sqlData = "(" + sqlData + ")"
	}

	sql = `
	#StandardSQL
	SELECT *
	FROM ` + sqlHeaderCol + sqlHeaderRow + sqlData + `
	ORDER BY  group_code , ` + orderCol + strings.TrimRight(orderRow, ",")
	return
}
