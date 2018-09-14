package RP055_BestSales

const (
	shop_code   = "shop_code"
	shop_name   = "shop_name"
	mc_yyyy     = "mc_yyyy"
	mc_mm       = "mc_mm"
	mc_weekdate = "mc_weekdate"
	mc_dd       = "mc_dd"

	bccm_media_group1_cd = "bqmg_media_group1_cd"
	bccm_media_group2_cd = "bqmg_media_group2_cd"
	bccm_media_group3_cd = "bqmg_media_group3_cd"
	bccm_media_group4_cd = "bqmg_media_group4_cd"

	rank_no                  = "rank_no"
	jan_code                 = "jan_code"
	product_name             = "product_name"
	author_name              = "author_name"
	publisher_name           = "publisher_name"
	selling_date             = "selling_date"
	sales_tax_exc_unit_price = "sales_tax_exc_unit_price"

	sales_body_quantity = "sales_body_quantity"
	magazine_code       = "magazine_code"

	// add new (comment https://dev-backlog.rsp.honto.jp/backlog/view/FUJI-4654#comment-3209911)
	// 在庫数
	stock_quantity = "stock_quantity"
	// 期間入荷累計
	cumulative_receiving_quantity = "cumulative_receiving_quantity"
	// 期間売上累計
	cumulative_sales_quantity = "cumulative_sales_quantity"
	// 初売上日
	first_sales_date = "first_sales_date"

	// Const Group type search ("0": Date, "1": Week, "2": Month)
	GROUP_TYPE_DATE  = "0"
	GROUP_TYPE_WEEK  = "1"
	GROUP_TYPE_MONTH = "2"

	GROUP_TYPE_DATE_TEXT  = "日別"
	GROUP_TYPE_WEEK_TEXT  = "週別"
	GROUP_TYPE_MONTH_TEXT = "月別"

	// 商品区分(雑誌（月刊誌）: 1、雑誌（週刊誌）: 2、雑誌（季刊誌））: 3)
	BQSL_MAGAZINE_CODE_MONTH   = "1"
	BQSL_MAGAZINE_CODE_WEEK    = "2"
	BQSL_MAGAZINE_CODE_QUARTER = "3"

	BQSL_MAGAZINE_CODE_MONTH_TEXT   = "雑誌（月刊誌）"
	BQSL_MAGAZINE_CODE_WEEK_TEXT    = "雑誌（週刊誌）"
	BQSL_MAGAZINE_CODE_QUARTER_TEXT = "雑誌（季刊誌）"

	// 書籍 : 1, 雑誌: 2
	CONTROL_TYPE_BOOK     = "1"
	CONTROL_TYPE_MAGAZINE = "2"

	//Download type
	DOWNLOAD_TYPE_TOTAL_RESULT            = "0"
	DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION = "1"
	DOWNLOAD_TYPE_TOTAL_RESULT_STORE      = "2"

	DOWNLOAD_TYPE_TOTAL_RESULT_TEXT            = "集計結果"
	DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION_TEXT = "集計結果+推移"
	DOWNLOAD_TYPE_TOTAL_RESULT_STORE_TEXT      = "集計結果＋店舗"

	//Jan maker code default
	JAN_MAKER_CODE = "9784"
)

func initDefaultLayout() (cols map[string]string, rows map[string]string, sums map[string]string) {

	cols = map[string]string{
		shop_code:   "店舗コード",
		shop_name:   "店舗名称",
		mc_yyyy:     "年",
		mc_mm:       "月",
		mc_weekdate: "週",
		mc_dd:       "日",
	}

	rows = map[string]string{
		rank_no: "順位",

		bccm_media_group1_cd: "特大分類",
		bccm_media_group2_cd: "大分類",
		bccm_media_group3_cd: "中分類",
		bccm_media_group4_cd: "小分類",

		jan_code:                 "ＪＡＮ",
		product_name:             "品名",
		author_name:              "著者",
		publisher_name:           "出版社名",
		selling_date:             "発売日",
		sales_tax_exc_unit_price: "本体価格",

		stock_quantity:                "在庫数",
		cumulative_receiving_quantity: "期間入荷累計",
		cumulative_sales_quantity:     "期間売上累計",
		first_sales_date:              "初売上日",
		magazine_code:                 "雑誌コード+月号",
	}

	sums = map[string]string{
		sales_body_quantity: "期間売上合計",
	}

	return
}
