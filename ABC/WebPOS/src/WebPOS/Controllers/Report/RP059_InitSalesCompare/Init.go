package RP059_InitSalesCompare

import "github.com/goframework/gf"

const (
	_REPORT_ID   = "init_sales_compare"
	_REPORT_NAME = "初速比較"

	PATH_REPORT_INIT_SALES_COMPARE_SEARCH            = "/report/" + _REPORT_ID
	PATH_REPORT_INIT_SALES_COMPARE_QUERY_AJAX        = "/report/" + _REPORT_ID + "/query_ajax"
	PATH_REPORT_INIT_SALES_COMPARE_QUERY_DETAIL_AJAX = "/report/" + _REPORT_ID + "/query_detail_ajax"

	// init layout
	// JAN/ISBN : 1, 雑誌: 2
	CONTROL_TYPE_JAN      = "1"
	CONTROL_TYPE_MAGAZINE = "2"

	SEARCH_DATE_TYPE_40 = 40
	SEARCH_DATE_TYPE_14 = 14

	SEARCH_DATE_TYPE_40_TEXT = "40日"
	SEARCH_DATE_TYPE_14_TEXT = "14日"
)

func Init() {

	// 初速比較
	gf.HandleGetPost(PATH_REPORT_INIT_SALES_COMPARE_SEARCH, Search)
	gf.HandlePost(PATH_REPORT_INIT_SALES_COMPARE_QUERY_AJAX, Query)
	gf.HandlePost(PATH_REPORT_INIT_SALES_COMPARE_QUERY_DETAIL_AJAX, QueryDetailByJan)
}
