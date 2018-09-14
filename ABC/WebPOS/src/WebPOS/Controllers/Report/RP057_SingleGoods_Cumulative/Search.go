package RP057_SingleGoods_Cumulative

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"time"
)

// 単品推移（累計推移）
func Search(ctx *gf.Context) {

	form := QueryFormSingleGoods{}
	ctx.Form.ReadStruct(&form)
	user := WebApp.GetContextUser(ctx)

	sm := Models.ShopMasterModel{ctx.DB}
	listShop, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)

	// init shop
	ctx.ViewData["list_shop"] = listShop
	// init form
	if form.JAN == "" {
		form.DateFrom = time.Now().AddDate(0, -1, 0).Format(Common.DATE_FORMAT_YMD_SLASH)
		form.DateTo = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
		form.GroupType = "0"
	}

	//save cache form search
	keyForm := Common.GenerateString(_REPORT_ID + form.JAN + Common.CurrentDateTime())
	err = ctx.SaveCache(keyForm, form, 3600)
	Common.LogErr(err)

	// init form search
	ctx.ViewData["form"] = form
	ctx.ViewData["key_form"] = keyForm
	ctx.ViewData["link_revert"] = form.LinkRevert
	ctx.View = "report/057_single_goods_cumulative/search.html"
}
