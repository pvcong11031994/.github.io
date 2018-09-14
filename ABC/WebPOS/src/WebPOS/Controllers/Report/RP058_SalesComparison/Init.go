package RP058_SalesComparison

import "github.com/goframework/gf"

const (
	_REPORT_ID                               string = "sales_comparison"
	_REPORT_NAME                             string = "売上比較"
	PATH_REPORT_SALES_COMPARISON_SEARCH      string = "/report/" + _REPORT_ID
	PATH_REPORT_SALES_COMPARISON_QUERY_AJAX  string = "/report/" + _REPORT_ID + "/query_ajax"
	PATH_REPORT_SALES_COMPARISON_DETAIL_AJAX string = "/report/" + _REPORT_ID + "/query_detail_ajax"
)

func Init() {

	// 売上比較 画面
	gf.HandleGetPost(PATH_REPORT_SALES_COMPARISON_SEARCH, Search)
	gf.HandlePost(PATH_REPORT_SALES_COMPARISON_QUERY_AJAX, Query)
	gf.HandlePost(PATH_REPORT_SALES_COMPARISON_DETAIL_AJAX, QueryDetail)
}
