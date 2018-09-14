package UsersMaintenance

import (
	"github.com/goframework/gf"
)

const (
	PATH_MAINTENANCE_USER                     string = "/maintenance/user/*"
	PATH_MAINTENANCE_USER_SEARCH              string = "/maintenance/user/search"
	PATH_MAINTENANCE_USER_SEARCH_LIST         string = "/maintenance/user/list"
	PATH_MAINTENANCE_USER_NEW                 string = "/maintenance/user/new"
	PATH_MAINTENANCE_USER_NEW_CHECK           string = "/maintenance/user/new/check_ajax"
	PATH_MAINTENANCE_USER_CONFIRM             string = "/maintenance/user/confirm"
	PATH_MAINTENANCE_USER_CHANGE_PASS_VIEW    string = "/maintenance/user/pass"
	PATH_MAINTENANCE_USER_CHANGE_PASS_ACTION  string = "/maintenance/user/pass/change_ajax"
	PATH_MAINTENANCE_USER_REGISTER_CSV_VIEW   string = "/maintenance/user/regist"
	PATH_MAINTENANCE_USER_REGISTER_CSV_ACTION string = "/maintenance/user/regist/upload_ajax"
)

func Init() {

	// ユーザメンテ ナンス検索画面
	gf.HandleGet(PATH_MAINTENANCE_USER_SEARCH, Search)
	gf.HandlePost(PATH_MAINTENANCE_USER_SEARCH_LIST, List)

	// ユーザ新規追加画面
	gf.HandleGet(PATH_MAINTENANCE_USER_NEW, New)
	gf.HandlePost(PATH_MAINTENANCE_USER_NEW_CHECK, CheckAccount)

	gf.HandlePost(PATH_MAINTENANCE_USER_CONFIRM, Confirm)

	// パスワード変更
	gf.HandleGet(PATH_MAINTENANCE_USER_CHANGE_PASS_VIEW, ChangePassView)
	gf.HandlePost(PATH_MAINTENANCE_USER_CHANGE_PASS_ACTION, ChangePassAction)

	// CSVアップ ロード
	gf.HandleGet(PATH_MAINTENANCE_USER_REGISTER_CSV_VIEW, CsvUploadView)
	gf.HandlePost(PATH_MAINTENANCE_USER_REGISTER_CSV_ACTION, CsvUploadAction)
}
