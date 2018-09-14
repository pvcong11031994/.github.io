package RP062_SearchGoods

import "github.com/goframework/gf"

const (
	_REPORT_ID   string = "search_goods"
	_REPORT_NAME string = "商品検索"

	PATH_REPORT_SEARCH_GOODS            string = "/report/" + _REPORT_ID
	PATH_REPORT_SEARCH_GOODS_QUERY_AJAX string = "/report/" + _REPORT_ID + "/query_ajax"
)

func Init() {

	// 商品検索
	gf.HandleGetPost(PATH_REPORT_SEARCH_GOODS, Search)
	gf.HandlePost(PATH_REPORT_SEARCH_GOODS_QUERY_AJAX, Query)
}
