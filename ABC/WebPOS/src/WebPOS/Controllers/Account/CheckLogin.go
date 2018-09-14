package Account

import (
	"WebPOS/Common"
	"WebPOS/ControllersApi/Utils"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"strings"
	"time"
)

// Check login for all request from client
func CheckLogin(ctx *gf.Context) {
	if ctx.UrlPath == PATH_LOGIN || strings.HasPrefix(ctx.UrlPath, ApiUtils.ROUTE_API) {
		return
	}

	userInfo := WebApp.GetSessionUser(ctx)

	// 再認証ユーザ
	if userInfo != nil {
		newUserInfo, err := (&Models.UserMasterModel{ctx.DB}).GetUserInfoById(userInfo.UserID)
		Common.LogErr(err)

		if err == nil &&
			newUserInfo.UserID == userInfo.UserID &&
			newUserInfo.UserPass == userInfo.UserPass {
			ulsm := Models.UserLoginStatusModel{ctx.DB}
			isRejected, err := ulsm.CheckRejected(userInfo.UserID, ctx.Session.ID)
			Common.LogErr(err)
			if isRejected {
				ctx.ClearSession()
				ctx.Session.AddFlash(LOGIN_REJECT_MESSAGE, SESSION_KEY_LOGIN_MSG)
				ctx.Session.AddFlash(userInfo.UserID, SESSION_KEY_LOGIN_ID)
				ctx.Redirect(PATH_LOGIN)
				return
			}

			rpm := Models.ReportDesignMaster{ctx.DB}
			if design, err := rpm.GetReportDesign(newUserInfo.ShopChainCd); err != nil {
				Common.LogErr(err)
			} else {
				newUserInfo.Design = design
			}

			userInfo = newUserInfo
		} else {
			ctx.ClearSession()
			userInfo = nil
		}
	}

	if userInfo != nil {

		WebApp.SetContextUser(ctx, userInfo)

		ctx.ViewData["UserName"] = userInfo.UserName

		ctx.ViewData["IsVJ"] = strings.Contains(userInfo.ShopChainCd, Common.CHAIN_VJ)

		ctx.ViewData["IsHeadquarter"] = (userInfo.FlgAuth == "1")

		// check login first time
		if strings.TrimSpace(userInfo.PwExpiringDate) == "" {
			ctx.ViewData["change_pass"] = true
			return
		}
		timeExpiringDate, err := time.Parse(Common.DATE_FORMAT_ZERO_TIME, userInfo.PwExpiringDate)
		Common.LogErr(err)

		toDayYMD := Common.CurrentDate()
		timeNow, _ := time.Parse(Common.DATE_FORMAT_YMD, toDayYMD)
		if err != nil || timeNow.After(timeExpiringDate) {
			ctx.ViewData["change_pass"] = true
			return
		}

		// Init design
		if userInfo.Design != nil {
			ctx.ViewData["HeaderColor"] = userInfo.Design.ColorServiceBar
			ctx.ViewData["MenuColor"] = userInfo.Design.ColorDashboardBar
			logo := userInfo.Design.LogoServiceBar
			if logo == "" {
				logo = WebApp.DEFAULT_LOGO_FILE
			}
			ctx.ViewData["Logo"] = logo
		} else {
			ctx.ViewData["HeaderColor"] = ""
			ctx.ViewData["MenuColor"] = ""
			ctx.ViewData["Logo"] = WebApp.DEFAULT_LOGO_FILE
		}

		SetContextMenu(ctx)

		// ASO-5929 [BA]mBAWEB-v09f 店舗一覧ダウンロード-複数ファイル対応 - ADD START
		path := strings.TrimLeft(ctx.UrlPath,"/")
		if path == PATH_MAKE_SALE_DOWNLOAD {
			path = path + "?filename=" + ctx.Form.String("filename")
		}
		// ASO-5929 [BA]mBAWEB-v09f 店舗一覧ダウンロード-複数ファイル対応 - ADD END
		//メニュー名
		mcm := Models.MenuControlMasterModel{ctx.DB}
		// ASO-5929 [BA]mBAWEB-v09f 店舗一覧ダウンロード-複数ファイル対応 - EDIT START
		//menu,err := mcm.GetMenuName(strings.TrimLeft(ctx.UrlPath,"/"),userInfo.FlgMenuGroup)
		menu,err := mcm.GetMenuName(path,userInfo.FlgMenuGroup)
		// ASO-5929 [BA]mBAWEB-v09f 店舗一覧ダウンロード-複数ファイル対応 - EDIT END
		Common.LogErr(err)
		ctx.ViewData["MenuName"] = menu.MenuName
		titlePage := ""
		if menu.MenuName != "" {
			titlePage = " - " + menu.MenuName
		}
		ctx.ViewData["TitlePageMenu"] = titlePage

		return
	}

	if ctx.UrlPath != PATH_LOGIN {
		Common.LogOutput("user session " + ctx.Session.ID + " timeout")
		ctx.Redirect(PATH_LOGIN)

		//MJ require: blank page
		//ctx.WriteHttpStatus(http.StatusForbidden)
		//ctx.WriteS("")
	}
}
