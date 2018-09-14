package RP057_SingleGoods_Cumulative

type SingleItem struct {
	JanCd                string
	GoodsName            string
	AuthorName           string
	PublisherName        string
	SaleDate             string
	Price                int64
	StockInf             string
	SaleTotalDate        int64
	ReturnTotalDate      int64
	ReturnTotal          int64
	SaleTotal            int64
	StockCurCount        int64
	FirstSaleDate        string
	ShopCd               string
	ShopName             string
	SharedBookStoreCode  string
	SaleDay              map[string]int64
	ReturnDay            map[string]int64
	ReceivingQuantityDay map[string]int64
	SalesQuantityDay     map[string]int64
	SaleTotalDay         int64
	StockCount           int64
}
