package RP058_SalesComparison

type SingleItem struct {
	JanCd         string
	GoodsName     string
	AuthorName    string
	PublisherName string
	SaleDate      string
	Price         int64
	SaleTotalDate int64
	ReturnTotal   int64
	SaleTotal     int64
	StockCurCount int64
	FirstSaleDate string
	ShopCd        string
	ListShopCd    map[string]string
	ShopName      string
	SaleDay       map[string]map[string]int64
}
