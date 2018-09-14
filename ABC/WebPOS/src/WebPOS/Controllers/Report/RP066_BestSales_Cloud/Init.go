package RP066_BestSales_Cloud

import "github.com/goframework/gf"

const (
	_REPORT_ID          string = "best_sales_cloud"
	_REPORT_ID_DOWNLOAD string = "best_sales"
	_REPORT_NAME        string = "売上ベスト(Cloud)"

	PATH_REPORT_BEST_SALES_SEARCH     string = "/report/" + _REPORT_ID
	PATH_REPORT_BEST_SALES_QUERY_AJAX string = "/report/" + _REPORT_ID + "/query_ajax"
)

func Init() {

	// 売上ベスト
	gf.HandleGetPost(PATH_REPORT_BEST_SALES_SEARCH, Search)
	gf.HandlePost(PATH_REPORT_BEST_SALES_QUERY_AJAX, Query)
}
