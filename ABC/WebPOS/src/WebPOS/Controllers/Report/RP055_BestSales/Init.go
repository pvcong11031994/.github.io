package RP055_BestSales

import "github.com/goframework/gf"

const (
	_REPORT_ID   string = "best_sales"
	_REPORT_NAME string = "売上ベスト"

	PATH_REPORT_BEST_SALES_SEARCH     string = "/report/" + _REPORT_ID
	PATH_REPORT_BEST_SALES_QUERY_AJAX string = "/report/" + _REPORT_ID + "/query_ajax"
)

func Init() {

	// 売上ベスト
	gf.HandleGetPost(PATH_REPORT_BEST_SALES_SEARCH, Search)
	gf.HandlePost(PATH_REPORT_BEST_SALES_QUERY_AJAX, Query)
}
