package Account

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Dashboard"
	"WebPOS/ControllersApi/Utils"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"github.com/goframework/gf/html/template"
	"strings"
)

func SetContextMenu(ctx *gf.Context) {

	userInfo := WebApp.GetContextUser(ctx)
	um := Models.UserMasterModel{ctx.DB}

	mapEnableMenu, err := um.GetMenuPathByGroup()
	Common.LogErr(err)

	mcm := Models.MenuControlMasterModel{ctx.DB}
	listMenuLevel, _ := mcm.GetMenuLevel(userInfo.FlgMenuGroup)

	menu := WebApp.GetEnableMenu(mapEnableMenu, listMenuLevel, ctx)
	ctx.ViewData["MENU"] = template.HTML(menu.ToString())

	//Set custom report menu items
	srlm := Models.SettingReportLayoutModel{ctx.DB}
	reportMenuItems, err := srlm.GetUserReportMenu(userInfo.UserID)
	Common.LogErr(err)
	if reportMenuItems != nil {
		ctx.ViewData["user_report_menu"] = reportMenuItems
	}
}

func CheckMenuRole(ctx *gf.Context) {

	if ctx.UrlPath == PATH_LOGIN ||
		ctx.UrlPath == PATH_ROOT ||
		ctx.UrlPath == PATH_LOGOUT ||
		strings.HasPrefix(ctx.UrlPath, ApiUtils.ROUTE_API) ||
		ctx.UrlPath == Dashboard.PATH_TOP_PAGE {
		return
	}

	url := ctx.UrlPath
	if strings.HasSuffix(ctx.UrlPath, "_ajax") || strings.HasSuffix(ctx.UrlPath, "_download") {
		url = ctx.UrlPath[:strings.LastIndex(ctx.UrlPath, "/")]
		if ctx.UrlPath == PATH_MAKE_SALE_DOWNLOAD {
			url = ctx.UrlPath
		}
	}

	userInfo := WebApp.GetContextUser(ctx)

	mm := Models.MenuControlMasterModel{ctx.DB}
	if ctx.UrlPath != PATH_REPORT_DOWNLOAD {
		// ASO-5929 [BA]mBAWEB-v09f 店舗一覧ダウンロード-複数ファイル対応 - EDIT START
		//if !mm.CheckMenuByUrlByFlg(url[1:],userInfo.FlgMenuGroup) {
		//	ctx.Redirect("/")
		//}
		path := url
		if path == PATH_MAKE_SALE_DOWNLOAD {
			path = url[1:] + "?filename=" + ctx.Form.String("filename")
		} else {
			path = url[1:]
		}
		if !mm.CheckMenuByUrlByFlg(path,userInfo.FlgMenuGroup) {
			ctx.Redirect("/")
		}
		// ASO-5929 [BA]mBAWEB-v09f 店舗一覧ダウンロード-複数ファイル対応 - EDIT END
	}
}
