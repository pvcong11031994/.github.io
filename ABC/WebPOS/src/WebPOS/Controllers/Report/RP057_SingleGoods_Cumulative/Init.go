package RP057_SingleGoods_Cumulative

import "github.com/goframework/gf"

const (
	_REPORT_ID                                     string = "single_goods_cumulative"
	_REPORT_NAME                                   string = "単品推移"
	PATH_REPORT_SINGLE_GOODS_CUMULATIVE_SEARCH     string = "/report/" + _REPORT_ID
	PATH_REPORT_SINGLE_GOODS_CUMULATIVE_QUERY_AJAX string = "/report/" + _REPORT_ID + "/query_ajax"

	// Const Group type search ("0": Date, "1": Week, "2": Month, "3": Shop, "4": Super)
	GROUP_TYPE_DATE  = "0"
	GROUP_TYPE_WEEK  = "1"
	GROUP_TYPE_MONTH = "2"

	GROUP_TYPE_DATE_TEXT  = "日別"
	GROUP_TYPE_WEEK_TEXT  = "週別"
	GROUP_TYPE_MONTH_TEXT = "月別"
)

func Init() {

	// 単品推移（累計推移）画面
	gf.HandleGetPost(PATH_REPORT_SINGLE_GOODS_CUMULATIVE_SEARCH, Search)
	gf.HandlePost(PATH_REPORT_SINGLE_GOODS_CUMULATIVE_QUERY_AJAX, Query)
}
