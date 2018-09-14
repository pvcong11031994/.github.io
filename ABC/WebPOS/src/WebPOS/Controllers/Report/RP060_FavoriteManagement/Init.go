package RP060_FavoriteManagement

import "github.com/goframework/gf"

const (
	_REPORT_ID   = "favorite_management"
	_REPORT_NAME = "お気に入り管理"

	FAVORITE_MANAGEMENT                   = "/report/" + _REPORT_ID
	FAVORITE_MANAGEMENT_QUERY_LOAD_AJAX = "/report/" + _REPORT_ID + "/query_load_ajax"
	FAVORITE_MANAGEMENT_QUERY_DELETE_AJAX = "/report/" + _REPORT_ID + "/query_delete_ajax"
	FAVORITE_MANAGEMENT_QUERY_UPDATE_AJAX = "/report/" + _REPORT_ID + "/query_update_ajax"
	FAVORITE_MANAGEMENT_QUERY_INSERT_AJAX = "/report/" + _REPORT_ID + "/query_insert_ajax"

)

func Init() {

	// 初速比較
	gf.HandleGetPost(FAVORITE_MANAGEMENT, Search)
	gf.HandlePost(FAVORITE_MANAGEMENT_QUERY_LOAD_AJAX, Load)
	gf.HandlePost(FAVORITE_MANAGEMENT_QUERY_DELETE_AJAX, Delete)
	gf.HandlePost(FAVORITE_MANAGEMENT_QUERY_UPDATE_AJAX, Update)
	gf.HandlePost(FAVORITE_MANAGEMENT_QUERY_INSERT_AJAX, Insert)

}
