package Account

import (
	"github.com/goframework/gf"
)

const (
	PATH_LOGIN_REQUIRE      string = "*"
	PATH_ROOT               string = "/"
	PATH_LOGIN              string = "/account/login"
	PATH_LOGOUT             string = "/account/logout"
	PATH_REPORT_DOWNLOAD    string = "/report/download"
	PATH_MAKE_SALE_DOWNLOAD string = "/download/makersalestock/shop_list_download"
)

func Init() {

	gf.Filter(PATH_LOGIN_REQUIRE, CheckLogin)
	gf.Filter(PATH_LOGIN_REQUIRE, CheckMenuRole)

	// ログイン
	gf.HandleGet(PATH_LOGIN, Login)
	gf.HandlePost(PATH_LOGIN, LoginPost)

	// ログアウト
	gf.HandleGet(PATH_LOGOUT, Logout)
}
