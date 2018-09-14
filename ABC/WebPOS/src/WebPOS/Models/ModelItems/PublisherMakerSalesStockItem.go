package ModelItems

type PublisherMakerSSData struct {
	TransactionDate string `sql:"bqio_trn_date"`
	CommonShopCd    string `sql:"cm_shop_cd" header:"共有書店コード"`
	ShopCd          string `sql:"shm_shop_cd" header:"店舗コード"`
	ShopName        string `sql:"shm_shop_name" header:"店舗名"`
	ISBN            string `sql:"goods_cd" header:"ISBN"`
	GoodsName       string `sql:"goods_name" header:"書誌名"`
	PublisherName   string `sql:"publisher_name" header:"出版社"`
	Price           string `sql:"bqgm_price" header:"本体価格"`
	GoodsCount_     string `sql:"bqio_goods_count"`
	StockCount_     string `sql:"bqsc_stock_count"`
}
