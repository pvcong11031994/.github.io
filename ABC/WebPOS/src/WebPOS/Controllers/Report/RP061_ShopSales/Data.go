package RP061_ShopSales

type SingleItem struct {
	ShmSharedBookStoreCode string
	ShmShopCode            string
	ShmShopName            string
	ShmTelNo               string
	ShmBizStartTime        string
	ShmBizEndTime          string
	ShmAddress             string
	ShmShopNameShort       string

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
