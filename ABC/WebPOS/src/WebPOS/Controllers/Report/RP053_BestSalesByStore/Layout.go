package RP053_BestSalesByStore

const (
	mc_yyyy       = "mc_yyyy"
	mc_mm         = "mc_mm"
	mc_weekdate   = "mc_weekdate"
	mc_dd         = "mc_dd"

	rank_no         = "rank_no"
	bqio_jan_cd     = "bqio_jan_cd"
	goods_name      = "goods_name"
	writer_name     = "writer_name"
	publisher_name  = "publisher_name"
	bqgm_price      = "bqgm_price"

	bqio_goods_count = "bqio_goods_count"

	// add new (comment https://dev-backlog.rsp.honto.jp/backlog/view/FUJI-4654#comment-3209911)
	// 在庫
	stock_count = "stock_count"
	// 入荷累計
	total_arrival = "total_arrival"
	// 売上累計
	total_sales = "total_sales"

	// Const Group type search ("0": Date, "1": Week, "2": Month)
	GROUP_TYPE_DATE  = "0"
	GROUP_TYPE_WEEK  = "1"
	GROUP_TYPE_MONTH = "2"

	GROUP_TYPE_DATE_TEXT  = "日別"
	GROUP_TYPE_WEEK_TEXT  = "週別"
	GROUP_TYPE_MONTH_TEXT = "月別"
)

func initDefaultLayout() (cols map[string]string, rows map[string]string, sums map[string]string) {

	cols = map[string]string{
		mc_yyyy:       "年",
		mc_mm:         "月",
		mc_weekdate:   "週",
		mc_dd:         "日",
	}

	rows = map[string]string{
		rank_no: "順位",

		bqio_jan_cd:     "ＪＡＮ",
		goods_name:      "品名",
		writer_name:     "著者",
		publisher_name:  "出版社",
		bqgm_price:      "本体価格",

		stock_count:   "在庫",
		total_arrival: "入荷累計",
		total_sales:   "売上累計",
	}

	sums = map[string]string{
		bqio_goods_count: "販売数",
	}

	return
}
