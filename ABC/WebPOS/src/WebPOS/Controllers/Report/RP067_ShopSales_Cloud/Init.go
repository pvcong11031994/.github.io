package RP067_ShopSales_Cloud

import "github.com/goframework/gf"

const (
	_REPORT_ID          string = "shop_sales_cloud"
	_REPORT_ID_DOWNLOAD string = "shop_sales"
	_REPORT_NAME        string = "店舗別集計(Cloud)"

	PATH_REPORT_SHOP_SALES_SEARCH     string = "/report/" + _REPORT_ID
	PATH_REPORT_SHOP_SALES_QUERY_AJAX string = "/report/" + _REPORT_ID + "/query_ajax"
)

func Init() {

	// 店舗別集計(Cloud)
	gf.HandleGetPost(PATH_REPORT_SHOP_SALES_SEARCH, Search)
	gf.HandlePost(PATH_REPORT_SHOP_SALES_QUERY_AJAX, Query)
}
