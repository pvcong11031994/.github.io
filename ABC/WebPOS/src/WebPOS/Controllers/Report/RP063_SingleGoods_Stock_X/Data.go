package RP063_SingleGoods_Stock_X

type SingleItem struct {
	JanCd                string
	GoodsName            string
	AuthorName           string
	PublisherName        string
	SaleDate             string
	Price                int64
	SaleTotalDate        int64

	ReturnTotal int64
	SaleTotal   int64
	StockTotal  int64

	FirstSaleDate        string
	ShopCd               string
	ShopName             string
	SaleDay              map[string]int64
	ReturnDay            map[string]int64
	ReceivingQuantityDay map[string]int64
	SalesQuantityDay     map[string]int64
	SaleTotalDay         int64
	StockCountByShop     int64

	StockCountByShopSearchDate     int64
}
