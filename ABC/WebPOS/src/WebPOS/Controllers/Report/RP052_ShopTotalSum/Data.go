package RP052_ShopTotalSum

type DataSum struct {
	ResultData      []ResultData
	CountResultRows int64
	ShowLineFrom    int64
	ShowLineTo      int64
	PageCount       int64
	ThisPage        int64
	VJCharging      int
}

type ResultData struct {
	ShopCd        string
	ShopName      string
	JanCd         string
	GoodsName     string
	PublisherName string
	Price         int64
	SaleTotal     int64
	SaleTotalDate int64
	StockCount    int64
}
