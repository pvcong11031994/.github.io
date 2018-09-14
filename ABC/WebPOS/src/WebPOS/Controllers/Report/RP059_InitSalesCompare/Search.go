package RP059_InitSalesCompare

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"time"
)

// 初速比較検索
func Search(ctx *gf.Context) {

	form := QueryForm{
		ControlType: "1",
	}
	ctx.Form.ReadStruct(&form)

	// Get info user
	user := WebApp.GetContextUser(ctx)
	// Get List shop by user role
	sm := Models.ShopMasterModel{ctx.DB}
	listShop, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)

	ctx.ViewData["form"] = form
	ctx.ViewData["list_shop"] = listShop

	//Load form revert
	keyForm := ctx.Form.String("key_form")
	if keyForm != "" {
		formObj := RPComon.QueryFormSingleGoods{}
		ctx.LoadCache(keyForm, &formObj)
		ctx.ViewData["form"] = formObj
	}

	ctx.TemplateFunc["plus"] = Common.Plus
	ctx.TemplateFunc["checkLen"] = Common.CheckLenArray
	ctx.TemplateFunc["cvString"] = Common.ConvertArrayToString
	// init form search
	ctx.ViewData["default_date_from"] = time.Now().AddDate(0, -1, 0).Format(Common.DATE_FORMAT_YMD_SLASH)
	ctx.ViewData["default_date_to"] = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	ctx.ViewData["link_revert"] = PATH_REPORT_INIT_SALES_COMPARE_SEARCH

	ctx.View = "report/059_init_sales_compare/search.html"
}
