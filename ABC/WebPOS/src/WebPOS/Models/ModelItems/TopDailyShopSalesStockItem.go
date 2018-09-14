package ModelItems

type TopDailyShopSalesStockItem struct {
	CreateDate  string `sql:"tdsss_create_date"`
	UpdateDate  string `sql:"tdsss_update_date"`
	Servername  string `sql:"tdsss_servername"`
	Dbname      string `sql:"tdsss_dbname"`
	FranchiseCd string `sql:"tdsss_franchise_cd"`
	ChainCd     string `sql:"tdsss_chain_cd"`
	ShopCd      string `sql:"tdsss_shop_cd"`
	Date        string `sql:"tdsss_date"`
	Year        string `sql:"tdsss_year"`
	Month       string `sql:"tdsss_month"`
	Day         string `sql:"tdsss_day"`
	GoodsCount  int    `sql:"tdsss_goods_count"`
	StockCount  int    `sql:"tdsss_stock_count"`
}
