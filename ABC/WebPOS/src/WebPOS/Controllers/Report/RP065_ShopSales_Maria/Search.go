package RP065_ShopSales_Maria

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"time"
)

//売上ベスト
func Search(ctx *gf.Context) {

	form := QueryForm{}
	ctx.Form.ReadStruct(&form)
	user := WebApp.GetContextUser(ctx)

	//Default date form/date to
	ctx.ViewData["default_date_from"] = time.Now().AddDate(0, -1, 0).Format(Common.DATE_FORMAT_YMD_SLASH)
	ctx.ViewData["default_date_to"] = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)

	sm := Models.ShopMasterModel{ctx.DB}
	listShop, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)

	ctx.ViewData["form"] = form
	// init shop
	ctx.ViewData["list_shop"] = listShop
	//Load form revert
	keyForm := ctx.Form.String("key_form")
	if keyForm != "" {
		formObj := RPComon.QueryFormSingleGoods{}
		ctx.LoadCache(keyForm, &formObj)
		ctx.ViewData["form"] = formObj
	}

	ctx.TemplateFunc["plus"] = Common.Plus
	// init form search
	ctx.ViewData["link_revert"] = PATH_REPORT_SHOP_SALES_SEARCH
	ctx.View = "report/065_shop_sales_maria/search.html"
}
