package PublisherMakerBestSaleStockDownload_Maria

type SalesItem struct {
	SalesDateTime        string
	SharedBookStoreCode  string
	ShopCode             string
	ShopName             string
	JanCode              string
	ProductName          string
	MakerCode            string
	SalesTaxExcUnitPrice float64
	SalesBodyQuantity    int64
	StockQuantity        int64
}
