package RP062_SearchGoods

import (
	"github.com/goframework/gf"
	"WebPOS/Controllers/Report/RPComon"
	"time"
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
)

//商品検索
func Search(ctx *gf.Context) {

	form := QueryForm{}
	ctx.Form.ReadStruct(&form)
	user := WebApp.GetContextUser(ctx)

	//Default date form/date to
	ctx.ViewData["default_date_from"] = time.Now().AddDate(0, -1, 0).Format(Common.DATE_FORMAT_YMD_SLASH)
	ctx.ViewData["default_date_to"] = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	//Get list shop
	sm := Models.ShopMasterModel{ctx.DB}
	listShop, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)
	// init shop
	ctx.ViewData["list_shop"] = listShop
	ctx.ViewData["form"] = form

	//Load form revert
	keyForm := ctx.Form.String("key_form")
	if keyForm != "" {
		formObj := RPComon.QueryFormSingleGoods{}
		ctx.LoadCache(keyForm, &formObj)
		ctx.ViewData["form"] = formObj
	}
	ctx.TemplateFunc["value"] = Common.GetValueArray
	ctx.ViewData["link_revert"] = PATH_REPORT_SEARCH_GOODS
	ctx.View = "report/062_search_goods/search.html"
}
