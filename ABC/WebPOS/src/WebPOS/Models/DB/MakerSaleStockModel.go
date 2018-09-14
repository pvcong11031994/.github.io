package Models

import (
	"WebPOS/Common"
	"WebPOS/WebApp"
	"database/sql"
	"errors"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"strings"
)

type MakerSaleStockModel struct {
	DB *sql.DB
}

type MakerSSData struct {
	TransactionDate string `sql:"bqio_trn_date" header:"日付"`
	CommonShopCd    string `sql:"cm_shop_cd" header:"共有書店コード"`
	ShopCd          string `sql:"shm_shop_cd" header:"店舗コード"`
	ShopName        string `sql:"shm_shop_name" header:"店舗名"`
	ISBN            string `sql:"goods_cd" header:"ISBN"`
	GoodsName       string `sql:"goods_name" header:"書誌名"`
	PublisherName   string `sql:"publisher_name" header:"出版社"`
	Price           string `sql:"bqgm_price" header:"本体価格"`
	GoodsCount      string `sql:"bqio_goods_count"`
	StockCount      string `sql:"bqsc_stock_count"`
}

type MakerSSDataWithErr struct {
	Data *MakerSSData
	Err  error
}

const (
	_REPORT_NAME     = "メーカー売上・在庫ダウンロード_20170803"
	_REPORT_ID       = "maker_sale_stock_download"
	_REPORT_NAME_KEY = "REPORT_NAME_KEY"
	_BQ_DATA_LIMIT   = 10000000

	// 商品区分(和書: 1、雑誌（月刊誌）: 2、雑誌（週刊誌）: 3、雑誌（その他）: 4)
	_BQGM_GOODS_TYPE_BOOK         = 1
	_BQGM_GOODS_TYPE_ZASSHI_MONTH = 2
	_BQGM_GOODS_TYPE_ZASSHI_WEEK  = 3
	_BQGM_GOODS_TYPE_ZASSHI_OTHER = 4

	_BQGM_GOODS_TYPE_BOOK_TEXT         = "和書"
	_BQGM_GOODS_TYPE_ZASSHI_MONTH_TEXT = "雑誌（月刊誌）"
	_BQGM_GOODS_TYPE_ZASSHI_WEEK_TEXT  = "雑誌（週刊誌）"
	_BQGM_GOODS_TYPE_ZASSHI_OTHER_TEXT = "雑誌（その他）"
)

const (
	bqio_trn_date    = "bqio_trn_date"
	cm_shop_cd       = "cm_shop_cd"
	shm_shop_cd      = "shm_shop_cd"
	shm_shop_name    = "shm_shop_name"
	goods_cd         = "goods_cd"
	goods_name       = "goods_name"
	publisher_name   = "publisher_name"
	bqgm_price       = "bqgm_price"
	bqio_goods_count = "bqio_goods_count"
	bqsc_stock_count = "bqsc_stock_count"
)

func (this *MakerSaleStockModel) Search(makerCd, dateFrom, dateTo, JAN, dataMode string, goodsType int, shopCd []string, ctx *gf.Context) chan MakerSSDataWithErr {

	outputChan := make(chan MakerSSDataWithErr)

	go func() {

		query := bq.NewCommand()

		parameter := map[string]interface{}{}
		parameter["@dataset"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)

		// 店舗条件
		conditionBqshm := ""
		if len(shopCd) > 0 {
			conditionBqshm += `
	AND (shm.shm_server_name + '|' + shm.shm_shop_cd) IN {{shop_cd}}`
			parameter["shop_cd"] = shopCd
		}

		// 出力欄条件
		conditionDataMode := ""
		switch dataMode {
		case "0":
			conditionDataMode += `
	SUM(INTEGER(IFNULL(bqio.bqio_goods_count,0))) bqio_goods_count,
	'' bqsc_stock_count,`
		case "1":
			conditionDataMode += `
	'' bqio_goods_count,
	FIRST(IFNULL(bqstc.bqsc_stock_count,INTEGER(0))) bqsc_stock_count,
	`
		}

		condition := ""
		// JAN条件
		if JAN != "" {
			condition += `
		AND bqio_jan_cd = {{jan_cd}}`
		}

		// 商品区分条件
		switch goodsType {
		case _BQGM_GOODS_TYPE_BOOK:
			condition += `
		AND bqgm.bqgm_goods_type = '20'
		`
		case _BQGM_GOODS_TYPE_ZASSHI_MONTH:
			condition += `
		AND bqgm.bqgm_goods_type = '21'
		AND (bqgm.bqgm_kihon_magazinecode LIKE '0%' OR bqgm.bqgm_kihon_magazinecode LIKE '1%')
		`
		case _BQGM_GOODS_TYPE_ZASSHI_WEEK:
			condition += `
		AND bqgm.bqgm_goods_type = '21'
		AND (bqgm.bqgm_kihon_magazinecode LIKE '2%' OR bqgm.bqgm_kihon_magazinecode LIKE '3%')
		`
		case _BQGM_GOODS_TYPE_ZASSHI_OTHER:
			condition += `
		AND bqgm.bqgm_goods_type = '21'
		AND bqgm.bqgm_kihon_magazinecode NOT LIKE '0%'
		AND bqgm.bqgm_kihon_magazinecode NOT LIKE '1%'
		AND bqgm.bqgm_kihon_magazinecode NOT LIKE '2%'
		AND bqgm.bqgm_kihon_magazinecode NOT LIKE '3%'
		`
		}

		dateSearchFrom := strings.Replace(dateFrom, "/", "", -1)
		dateSearchTo := strings.Replace(dateTo, "/", "", -1)

		parameter["date_from"] = dateSearchFrom
		parameter["date_to"] = dateSearchTo
		parameter["month_from"] = dateSearchFrom[:6]
		parameter["month_to"] = dateSearchTo[:6]
		parameter["maker_cd"] = "%" + makerCd + "%"
		parameter["jan_cd"] = JAN
		parameter["server_name"] = parseServerList(shopCd)
		parameter["#sqlStock"] = createStockQuery(dateSearchTo, shopCd, ctx)
		query.Parameters = parameter

		query.CommandText = `
SELECT * FROM (
	SELECT
		LEFT(REPLACE(bqio_trn_date,'-','/'), 10) bqio_trn_date,
		FIRST(bqsc.shm_shared_book_store_code) cm_shop_cd,
		bqsc.shm_shop_cd shm_shop_cd,
		FIRST(bqsc.shm_shop_name) shm_shop_name,
		bqio_jan_cd  goods_cd,
		FIRST(bqgm.bqgm_goods_name) goods_name,
		FIRST(IF(bqgm.bqgm_goods_type IN('20','21'),
				bqgm.bqgm_publisher_name,
				IFNULL(bqmk.bqmk_maker_name,''))) publisher_name,
		FIRST(bqio_price) bqgm_price,
		` + conditionDataMode + `
	FROM
	-----店舗-----
	(
		SELECT
			shm.shm_server_name shm_server_name,
			shm.shm_shop_cd shm_shop_cd,
			shm.shm_shop_name shm_shop_name,
			shm.shm_shared_book_store_code shm_shared_book_store_code,
		FROM
			[{{@dataset}}.shop_master] shm
		WHERE TRUE ` + conditionBqshm + `
		GROUP BY
			shm_server_name,
			shm_shop_cd,
			shm_shop_name,
			shm_shared_book_store_code,
	) bqsc
	-----店舗-----

	-----取引-----
	LEFT JOIN EACH
	(
		SELECT
			bqio_servername,
			bqio_shop_cd,
			bqio_jan_cd,
			bqio_bq_work_type,
			bqio_trn_date,
			bqio_goods_count,
			bqio_goods_name,
			bqio_price,
			bqio_dbname,
		FROM
			TABLE_QUERY({{@dataset}},"
		REGEXP_MATCH(TABLE_ID, '^bq_inout_[A-Z][A-Z0-9]+_\\d{1,2}_\\d{14}$')
		AND LEFT(RIGHT(TABLE_ID, 14), 6) BETWEEN {{month_from}} AND {{month_to}}
		AND REGEXP_EXTRACT(TABLE_ID,'bq_inout_([A-Z][A-Z0-9]+)') IN {{server_name}}
		AND REGEXP_MATCH(TABLE_ID, '_(1)|(6)|(13)_')
		"), (
		SELECT * FROM (
			SELECT
				'' bqio_servername,
				'' bqio_shop_cd,
				'' bqio_jan_cd,
				'' bqio_bq_work_type,
				'' bqio_trn_date,
				0 bqio_goods_count,
				'' bqio_goods_name,
				0 bqio_price,
				'' bqio_dbname,
			)
		WHERE FALSE
		)
		WHERE
			LEFT(REPLACE(REPLACE(bqio_trn_date, '-', ''),'/',''), 8) BETWEEN {{date_from}} AND {{date_to}}
	) bqio
	ON
		bqio.bqio_shop_cd = bqsc.shm_shop_cd
		AND bqio.bqio_servername = bqsc.shm_server_name
	-----取引-----

	-----在庫-----
	LEFT JOIN EACH
	({{#sqlStock}}) bqstc
	ON
		bqio.bqio_shop_cd = bqstc.bqsc_shop_cd
		AND bqio.bqio_servername = bqstc.bqsc_servername
		AND bqio.bqio_jan_cd = bqstc.bqsc_jan_cd
	-----在庫-----

	-----商品-----
	LEFT JOIN EACH [{{@dataset}}.bq_goods_master] bqgm
	ON
		bqio.bqio_jan_cd = bqgm.bqgm_jan_cd
	-----商品-----

	-----出版社-----
	LEFT JOIN EACH [{{@dataset}}.bq_maker] bqmk
	ON
		bqgm.bqgm_label_cd = bqmk.bqmk_label_cd
		AND bqio.bqio_servername = bqmk.bqmk_servername
	-----出版社-----
	WHERE TRUE
` + condition + `
	GROUP BY
		bqio_trn_date,
		shm_shop_cd,
		goods_cd
) bq
WHERE
	bq.publisher_name LIKE {{maker_cd}}
`

		sqlSum, err := query.Build()
		Common.LogErr(exterror.WrapExtError(err))
		// Output query to log file
		Common.LogSQL(ctx, sqlSum)

		// BQサービス接続
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
			outputChan <- MakerSSDataWithErr{nil, exterror.WrapExtError(err)}
			close(outputChan)
			return
		}

		// 検索ログ設定
		ctx.SetSessionFlash(_REPORT_NAME_KEY, _REPORT_NAME)
		//totalRows, jobId, err := conn.QueryForResponseBySql(sqlSum, ctx, _REPORT_ID)
		tag := "report=" + _REPORT_ID
		if makerCd != "" {
			tag = tag + ",出版社=" + `"` + makerCd + `"`
		}
		tag = tag + ",店舗 IN (" + strings.Join(shopCd, ",") + ")"
		tag = tag + ",期間=" + `"` + dateFrom + "～" + dateTo + `"`
		if JAN != "" {
			tag = tag + ",JAN=" + `"` + JAN + `"`
		}
		if goodsType == _BQGM_GOODS_TYPE_BOOK {
			tag = tag + ",商品区分=" + `"` + _BQGM_GOODS_TYPE_BOOK_TEXT + `"`
		} else if goodsType == _BQGM_GOODS_TYPE_ZASSHI_MONTH {
			tag = tag + ",商品区分=" + `"` + _BQGM_GOODS_TYPE_ZASSHI_MONTH_TEXT + `"`
		} else if goodsType == _BQGM_GOODS_TYPE_ZASSHI_WEEK {
			tag = tag + ",商品区分=" + `"` + _BQGM_GOODS_TYPE_ZASSHI_WEEK_TEXT + `"`
		} else if goodsType == _BQGM_GOODS_TYPE_ZASSHI_OTHER {
			tag = tag + ",商品区分=" + `"` + _BQGM_GOODS_TYPE_ZASSHI_OTHER_TEXT + `"`
		}
		if dataMode == "0" {
			tag = tag + ",データモード=" + `"売上"`
		} else if dataMode == "1" {
			tag = tag + ",データモード=" + `"在庫"`
		}
		keyErr = errors.New("KEY_ERR")
		msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
		totalRows, jobId, err := conn.QueryForResponseBySql(sqlSum, ctx, tag, ctx, msgRetryTmp, retryMap)

		if err != nil {
			outputChan <- MakerSSDataWithErr{nil, exterror.WrapExtError(err)}
			close(outputChan)
			return
		}

		// 最大限処理
		if totalRows > _BQ_DATA_LIMIT {
			outputChan <- MakerSSDataWithErr{nil, exterror.WrapExtError(errors.New("Respone data too large"))}
			close(outputChan)
			return
		}
		keyErr = errors.New("KEY_ERR")
		msgRetryTmp = strings.Replace(exterror.WrapExtError(keyErr).Error(), Common.REPORT_MSG_ERROR, Common.REPORT_MSG_RETRY, -1)
		dataChan, err := conn.GetResponseData(jobId, 0, _BQ_DATA_LIMIT, ctx, msgRetryTmp, retryMap)

		if err != nil {
			outputChan <- MakerSSDataWithErr{nil, exterror.WrapExtError(err)}
			close(outputChan)
			return
		}

		// データ取得
		for {
			row, ok := <-dataChan
			if !ok {
				break
			}

			newMSSData := MakerSSData{}
			newMSSData.TransactionDate = row.ValueMap[bqio_trn_date].String()
			newMSSData.CommonShopCd = row.ValueMap[cm_shop_cd].String()
			newMSSData.ShopCd = row.ValueMap[shm_shop_cd].String()
			newMSSData.ShopName = row.ValueMap[shm_shop_name].String()
			newMSSData.ISBN = row.ValueMap[goods_cd].String()
			newMSSData.GoodsName = row.ValueMap[goods_name].String()
			newMSSData.PublisherName = row.ValueMap[publisher_name].String()
			newMSSData.Price = row.ValueMap[bqgm_price].String()
			switch dataMode {
			case "0":
				newMSSData.GoodsCount = row.ValueMap[bqio_goods_count].String()
			case "1":
				newMSSData.StockCount = row.ValueMap[bqsc_stock_count].String()
			}

			outputChan <- MakerSSDataWithErr{&newMSSData, nil}
		}

		close(outputChan)
	}()

	return outputChan
}

// サーバ名と店舗名を分ける
func parseServerList(serverShops []string) []string {

	listServerName := []string{}
	mapServerName := map[string]bool{}
	for _, sh := range serverShops {
		//serverName := strings.Split(sh, "|")[0]
		serverName := ""
		arrsh := strings.Split(sh, "|")
		if len(arrsh) > 0 {
			serverName = arrsh[0]
		}
		if !mapServerName[serverName] && strings.Compare(serverName, "") != 0 {
			mapServerName[serverName] = true
			listServerName = append(listServerName, serverName)
		}
	}
	return listServerName
}

func createStockQuery(date string, shopCd []string, ctx *gf.Context) bq.Command {
	var query string

	if strings.HasPrefix(Common.CurrentDate(), date[:6]) {
		query = `
SELECT
	bqsc.bqsc_servername bqsc_servername,
	bqsc.bqsc_shop_cd bqsc_shop_cd,
	bqsc.bqsc_jan_cd bqsc_jan_cd,
	bqsc.bqsc_stock_count + IFNULL(bqio.bqio_goods_count,0) bqsc_stock_count,
FROM (
	SELECT
		bqsc_servername,
		bqsc_shop_cd,
		bqsc_jan_cd,
		FIRST(bqsc_stock_count) bqsc_stock_count,
	FROM
		TABLE_QUERY({{@dataset}}, "
	REGEXP_MATCH(TABLE_ID,'^bq_stok_cur_[A-Z][A-Z0-9]')
	"), (
		SELECT * FROM (
		SELECT
			'' bqsc_servername,
			'' bqsc_shop_cd,
			'' bqsc_jan_cd,
			0 bqsc_stock_count,
		)
		WHERE FALSE
	)
	GROUP BY
		bqsc_servername,
		bqsc_shop_cd,
		bqsc_jan_cd,
) bqsc
LEFT JOIN EACH (
	SELECT
		bqio_shop_cd,
		bqio_servername,
		bqio_jan_cd,
		bqio_trn_date,
		SUM(IF(bqio_bq_work_type IN ('00000001', '20000008', '20000009', '20000012', '21000005'), - IFNULL(bqio_goods_count, 0), IFNULL(bqio_goods_count,0))) AS bqio_goods_count
	FROM
		TABLE_QUERY({{@dataset}}, "
	REGEXP_MATCH(TABLE_ID, '^bq_inout_[A-Z][A-Z0-9]+_\\d{1,2}_\\d{14}$')
	AND LEFT(RIGHT(TABLE_ID, 14), 6) = {{month}}
	AND REGEXP_EXTRACT(TABLE_ID,'bq_inout_([A-Z][A-Z0-9]+)') IN {{server_name}}
	AND REGEXP_MATCH(TABLE_ID, '_(1)|(6)|(13)_')
	"), (
		SELECT
			'' bqio_servername,
			'' bqio_shop_cd,
			'' bqio_jan_cd,
			'' bqio_bq_work_type,
			0 bqio_goods_count,
			'' bqio_trn_date,
	)
	WHERE
		bqio_jan_cd NOT LIKE '99%'
		AND LEFT(REPLACE(REPLACE(bqio_trn_date, '-', ''),'/',''), 8) > {{date}}
	GROUP BY
		bqio_shop_cd,
		bqio_servername,
		bqio_jan_cd,
		bqio_trn_date,
) bqio
ON
	bqio.bqio_shop_cd = bqsc.bqsc_shop_cd
	AND bqio.bqio_servername = bqsc.bqsc_servername
	AND bqio.bqio_jan_cd = bqsc.bqsc_jan_cd
`
	} else {
		query = `
SELECT
	bqsc.bqsc_servername bqsc_servername,
	bqsc.bqsc_shop_cd bqsc_shop_cd,
	bqsc.bqsc_jan_cd bqsc_jan_cd,
	bqsc.bqsc_stock_count + IFNULL(bqio.bqio_goods_count,0) bqsc_stock_count,
FROM (
	SELECT
		bqsc_servername,
		bqsc_shop_cd,
		bqsc_jan_cd,
		FIRST(bqsc_stock_count) bqsc_stock_count,
	FROM
		TABLE_QUERY({{@dataset}}, "
	REGEXP_MATCH(TABLE_ID,'^bq_stok_cur_[A-Z][A-Z0-9]')
	")
	GROUP BY
		bqsc_servername,
		bqsc_shop_cd,
		bqsc_jan_cd,
	) bqsc
LEFT JOIN EACH (
	SELECT
		bqio_shop_cd,
		bqio_servername,
		bqio_jan_cd,
		SUM(IF(bqio_bq_work_type IN ('00000001', '20000008', '20000009', '20000012', '21000005'), - IFNULL(bqio_goods_count, 0), IFNULL(bqio_goods_count,0))) AS bqio_goods_count
	FROM
		TABLE_QUERY({{@dataset}}, "
	REGEXP_MATCH(TABLE_ID, '^bq_inout_[A-Z][A-Z0-9]+_\\d{1,2}_\\d{14}$')
	AND LEFT(RIGHT(TABLE_ID, 14), 6) = {{month}}
	AND REGEXP_EXTRACT(TABLE_ID,'bq_inout_([A-Z][A-Z0-9]+)') IN {{server_name}}
	AND REGEXP_MATCH(TABLE_ID, '_(1)|(6)|(13)_')
	"), (
		SELECT
			'' bqio_servername,
			'' bqio_shop_cd,
			'' bqio_jan_cd,
			'' bqio_bq_work_type,
			'' bqio_trn_date,
			0 bqio_goods_count,
			'' bqio_goods_name,
			0 bqio_price,
	)
	WHERE
		bqio_jan_cd NOT LIKE '99%'
		AND LEFT(REPLACE(REPLACE(bqio_trn_date, '-', ''),'/',''), 8) > {{date}}
	GROUP BY
		bqio_shop_cd,
		bqio_servername,
		bqio_jan_cd,
) bqio
ON
	bqio.bqio_shop_cd = bqsc.bqsc_shop_cd
	AND bqio.bqio_servername = bqsc.bqsc_servername
	AND bqio.bqio_jan_cd = bqsc.bqsc_jan_cd
`
	}

	cmd := bq.NewCommand()
	cmd.CommandText = query
	date = strings.Replace(strings.Replace(date, "-", "", -1), "/", "", -1)

	cmd.Parameters["date"] = date
	cmd.Parameters["@dataset"] = ctx.Config.Str(WebApp.CONFIG_KEY_BQ_DATASET, WebApp.DEFAULT_BQ_DATASET)
	cmd.Parameters["server_name"] = parseServerList(shopCd)
	cmd.Parameters["month"] = date[:6]

	return cmd
}
