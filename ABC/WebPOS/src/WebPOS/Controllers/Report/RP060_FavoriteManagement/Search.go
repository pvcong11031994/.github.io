package RP060_FavoriteManagement

import (
	"WebPOS/Common"
	"github.com/goframework/gf"
	"time"
)

// 初速比較検索
func Search(ctx *gf.Context) {

	//init form search
	ctx.ViewData["default_date_from"] = time.Now().AddDate(0, -1, 0).Format(Common.DATE_FORMAT_YMD_SLASH)
	ctx.ViewData["default_date_to"] = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	//START --- Check Link Goods_Search
	form := QueryForm{}
	ctx.Form.ReadStruct(&form)
	ctx.ViewData["form"] = form
	if form.KeySearch != "" {
		InsertOrUpdate(ctx, form)
	}
	//END ---
	ctx.ViewData["link_revert"] = FAVORITE_MANAGEMENT
	ctx.View = "report/060_favorite_management/search.html"
}
