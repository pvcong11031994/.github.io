package RP052_ShopTotalSum

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

func buildSqlFast(form QueryForm, ctx *gf.Context) (sql string) {

	parameter := map[string]interface{}{}
	parameter["@DATASET"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// condition
	// 店舗
	conditionBqshm := ""
	if len(form.ShopCd) > 0 {
		conditionBqshm += `
    AND CONCAT(shm.shm_server_name ,'|', shm.shm_shop_cd) IN {{shop_cd}}    `
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

	conditionBqio := `
		AND bqio.bqio_jan_cd NOT LIKE '99%'
	`
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	conditionBqmg := ""
	// メディア大分類
	// メディア中分類
	// メディア小分類
	if len(form.MediaGroup3Cd) > 0 {
		conditionBqmg += `
    AND (
        SUBSTR(bqio.bqio_media_cd,1, 6) IN {{gm_media_group3_cd}}
        )
    `
		parameter["gm_media_group3_cd"] = form.MediaGroup3Cd
	} else if len(form.MediaGroup2Cd) > 0 {
		conditionBqmg += `
    AND (
        SUBSTR(bqio.bqio_media_cd,1, 4) IN {{gm_media_group2_cd}}
        )
    `
		parameter["gm_media_group2_cd"] = form.MediaGroup2Cd
	} else if len(form.MediaGroup1Cd) > 0 {
		conditionBqmg += `
    AND (
        SUBSTR(bqio.bqio_media_cd,1, 2) IN {{gm_media_group1_cd}}
        )
     `
		parameter["gm_media_group1_cd"] = form.MediaGroup1Cd
	}

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// JAN
	if form.JAN != "" {
		conditionBqio += `
		AND (bqio.bqio_jan_cd LIKE {{jan_cd}} OR bqgm.bqgm_isbn LIKE {{jan_cd}})
    `
		parameter["jan_cd"] = form.JAN + "%"
	}

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	conditionBqmk := ""
	// メーカー
	if len(form.MakerCd) > 0 {
		makerCdConditionArr := []string{}
		for i, v := range form.MakerCd {
			makerCd := "maker_cd_" + strconv.Itoa(i)
			makerCdConditionArr = append(makerCdConditionArr, ` bqgm.bqgm_publisher_name LIKE {{`+makerCd+`}} OR bqmk.bqmk_maker_name LIKE {{`+makerCd+`}}`)
			if !strings.Contains(v, "%") {
				v = "%" + v + "%"
			}
			parameter[makerCd] = v
		}
		conditionBqmk += `
    AND (` + strings.Join(makerCdConditionArr, " OR ") + `)`
	}

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	conditionServerBqio := ""
	conditionServerBqsc := ""
	//サーバー名
	listServer := RPComon.ParseServerList(form.ShopCd)
	serverBqioConditionArr := []string{}
	serverBqscConditionArr := []string{}
	for i, v := range listServer {
		serverN := "server_" + strconv.Itoa(i)
		serverBqioConditionArr = append(serverBqioConditionArr, ` bqio._TABLE_SUFFIX LIKE {{`+serverN+`}} `)
		serverBqscConditionArr = append(serverBqioConditionArr, ` bqsc._TABLE_SUFFIX = {{`+serverN+`}} `)
		if !strings.Contains(v, "%") {
			v = "%" + v + "%"
		}
		parameter[serverN] = v
	}
	conditionServerBqio += `
    AND (` + strings.Join(serverBqioConditionArr, " OR ") + `)`
	conditionServerBqsc += `
    AND (` + strings.Join(serverBqscConditionArr, " OR ") + `)`
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	sqlSC := `
#StandardSQL
SELECT
	shm.shm_shop_name AS shm_shop_name,
	shm.shm_shop_cd AS shm_shop_cd,
	bqio.bqio_jan_cd AS bqio_jan_cd,
	bqgm.bqgm_goods_name AS goods_name,
	IF(bqgm.bqgm_goods_type IN('20','21'),
           bqgm.bqgm_publisher_name,
		    IFNULL(bqmk.bqmk_maker_name,'')) AS publisher_name,
	bqgm.bqgm_price AS bqgm_price,
	bqsc.bqsc_stock_count AS stock_count,
	bqsc.bqsc_cum_sales_quantity AS total_sales,
	SUM(IFNULL(bqio.bqio_goods_count,0)) AS total_sales_date
FROM ` + "`{{@DATASET}}.bq_inout_*`" + `  bqio
LEFT JOIN ` + "`{{@DATASET}}.bq_stock_*`" + `  bqsc
ON bqsc.bqsc_servername = bqio.bqio_servername
AND bqsc.bqsc_shop_cd = bqio.bqio_shop_cd
AND bqsc.bqsc_jan_cd = bqio.bqio_jan_cd
` + conditionServerBqsc + `
LEFT JOIN {{@DATASET}}.shop_master  shm
	ON shm.shm_server_name = bqio.bqio_servername
	AND shm.shm_shop_cd = bqio.bqio_shop_cd
LEFT JOIN {{@DATASET}}.bq_goods_master bqgm
 	ON bqio.bqio_jan_cd = bqgm.bqgm_jan_cd
LEFT JOIN {{@DATASET}}.bq_maker bqmk
 ON bqgm.bqgm_label_cd = bqmk.bqmk_label_cd
 AND bqio.bqio_servername = bqmk.bqmk_servername
WHERE bqio.bqio_trn_date BETWEEN {{calendar_from}} AND {{calendar_to}}
 	AND SUBSTR(bqio._TABLE_SUFFIX,-14,6) BETWEEN {{month_from}} AND {{month_to}}
	AND SUBSTR(bqio._TABLE_SUFFIX,-8) = '00000000'
	AND bqio._TABLE_SUFFIX LIKE '%_1_%'
` + conditionServerBqio + `
` + conditionBqmk + `
` + conditionBqio + `
` + conditionBqshm + `
` + conditionBqmg + `
GROUP BY
	shm.shm_shop_name,
	shm.shm_shop_cd,
	bqio.bqio_jan_cd,
	bqgm.bqgm_goods_name,
	bqgm.bqgm_goods_type,
	bqgm.bqgm_publisher_name,
	bqmk.bqmk_maker_name,
	bqgm.bqgm_price,
	bqsc.bqsc_stock_count,
	bqsc.bqsc_cum_sales_quantity
	`
	parameter["calendar_from"] = dateSearchFrom
	parameter["calendar_to"] = dateSearchTo
	parameter["month_from"] = dateSearchFrom[:6]
	parameter["month_to"] = dateSearchTo[:6]
	parameter["server_name"] = RPComon.ParseServerList(form.ShopCd)

	sqlS := bq.NewCommand()
	sqlS.Parameters = parameter
	sqlS.CommandText = sqlSC

	var err error
	sql, err = sqlS.Build()
	Common.LogErr(exterror.WrapExtError(err))

	return
}
