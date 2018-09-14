package PublisherMakerBestSaleStockDownload_Maria

import "encoding/gob"

const (
	_MSS_DATA_DOWNLOAD_KEY_ = "_MSS_DATA_DOWNLOAD_KEY_"
	_REPORT_NAME_PUBLISHER  = "出版社別ダウンロード"
	_REPORT_ID_PUBLISHER    = "publisher_maker_sale_stock_download"

	TYPE_SEARCH_SALES           = "0"
	TYPE_SEARCH_STOCK           = "1"
	TYPE_SEARCH_SALES_RECEIVING = "2"
	TYPE_SEARCH_SALES_RETURN    = "3"

	TYPE_SEARCH_SALES_TEXT           = "売上＋在庫"
	TYPE_SEARCH_STOCK_TEXT           = "累計＋在庫"
	TYPE_SEARCH_SALES_RECEIVING_TEXT = "入荷＋売上"
	TYPE_SEARCH_SALES_RETURN_TEXT    = "入荷＋売上＋返品"
)

type Form struct {
	ShopCd    []string `form:"shop_cd"`
	DateFrom  string   `form:"date_from"`
	DateTo    string   `form:"date_to"`
	JAN       string   `form:"jan_cd"`
	DataMode  string   `form:"data_mode"`
}

type RpData struct {
	HeaderCols      []string
	Cols            [][]string
	Rows            [][]string
	CountResultRows int
}

func init() {

	gob.Register(Form{})
}

var (
	LIST_HEADER_SALES_AND_RECEIVING = []string{
		"日付",
		"共有書店コード",
		"店舗コード",
		"店舗名",
		"ISBN",
		"書誌名",
		"出版社",
		"本体価格",
		"入荷数",
		"売上数",
	}

	LIST_HEADER_SALES_AND_RETURN = []string{
		"日付",
		"共有書店コード",
		"店舗コード",
		"店舗名",
		"ISBN",
		"書誌名",
		"出版社",
		"本体価格",
		"入荷数",
		"売上数",
		"返品数",
	}

	LIST_HEADER_SALES = []string{
		"日付",
		"共有書店コード",
		"店舗コード",
		"店舗名",
		"ISBN",
		"書誌名",
		"出版社",
		"本体価格",
		"売上数",
		"在庫数",
	}

	LIST_HEADER_STOCK = []string{
		"共有書店コード",
		"店舗コード",
		"店舗名",
		"ISBN",
		"書誌名",
		"出版社",
		"本体価格",
		"入荷累計",
		"売上累計",
		"在庫数",
	}
)
