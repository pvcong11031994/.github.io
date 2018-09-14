package RP064_BestSales_Maria

type SalesItem struct {
	Rank                        string
	JanCd                       string
	ProductName                 string
	AuthorName                  string
	MagazineCode                string
	MakerName                   string
	SellingDate                 string
	SalesTaxExcUnitPrice        float64
	StockQuantity               int64
	CumulativeReceivingQuantity int64
	CumulativeSalesQuantity     int64
	FirstSalesDate              string
	SalesBodyQuantity           int64
	SalesBodyQuantityShop       []int64
}
