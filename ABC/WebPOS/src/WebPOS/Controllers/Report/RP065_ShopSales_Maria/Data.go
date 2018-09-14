package RP065_ShopSales_Maria

type SingleItem struct {
	ShmSharedBookStoreCode      string
	ShmShopCode                 string
	ShmShopName                 string
	ShmTelNo                    string
	ShmBizStartTime             string
	ShmBizEndTime               string
	ShmAddress                  string
	ShmShopNameShort            string
	Rank                        string
	JanCd                       string
	ProductName                 string
	AuthorName                  string
	MakerName                   string
	SellingDate                 string
	SalesTaxExcUnitPrice        float64
	CumulativeReceivingQuantity int64
	CumulativeSalesQuantity     int64
	StockQuantity               int64
	FirstSalesDate              string
	SalesBodyQuantity           int64
	SalesBodyQuantityShop       []int64

	DataCount int64
	Data      [][]interface{}
}
type RpData struct {
	HeaderCols      []string
	Rows            []SingleItem
	CountResultRows int
	ShowLineFrom    int
	ShowLineTo      int
	PageCount       int
	ThisPage        int
	VJCharging      int
}
