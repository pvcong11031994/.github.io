package RP053_BestSalesByStore

import "github.com/goframework/gf"

const (
	_REPORT_ID                                 string = "best_sales_by_store"
	_REPORT_NAME                               string = "店舗別売上ベスト"
	PATH_REPORT_BEST_SALES_BY_STORE_SEARCH     string = "/report/" + _REPORT_ID
	PATH_REPORT_BEST_SALES_BY_STORE_QUERY_AJAX string = "/report/" + _REPORT_ID + "/query_ajax"
)

func Init() {

	// 売上ベスト帳票
	gf.HandleGet(PATH_REPORT_BEST_SALES_BY_STORE_SEARCH, Search)
	gf.HandlePost(PATH_REPORT_BEST_SALES_BY_STORE_QUERY_AJAX, Query)
}
