package RP052_ShopTotalSum

import "github.com/goframework/gf"

const (
	_REPORT_ID                            string = "shop_total_sum"
	_REPORT_NAME                          string = "売上累計検索"
	PATH_REPORT_SHOP_TOTAL_SUM_SEARCH     string = "/report/" + _REPORT_ID
	PATH_REPORT_SHOP_TOTAL_SUM_QUERY_AJAX string = "/report/" + _REPORT_ID + "/query_ajax"
)

func Init() {

	// 売上累計検索
	gf.HandleGet(PATH_REPORT_SHOP_TOTAL_SUM_SEARCH, Search)
	gf.HandlePost(PATH_REPORT_SHOP_TOTAL_SUM_QUERY_AJAX, Query)

}
