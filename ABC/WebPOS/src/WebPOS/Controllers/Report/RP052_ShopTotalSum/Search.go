package RP052_ShopTotalSum

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"time"
)

//売上累計検索
func Search(ctx *gf.Context) {

	user := WebApp.GetContextUser(ctx)

	// メディア分類
	bqctm := Models.BQCategoryModel{ctx.DB}
	listMedia, err := bqctm.GetMediaListByUser(user.UserID)
	Common.LogErr(err)

	ctx.ViewData["list_media"] = listMedia
	ctx.ViewData["date_from"] = time.Now().AddDate(0, -1, 0).Format(Common.DATE_FORMAT_YMD_SLASH)
	ctx.ViewData["date_to"] = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)
	//【○月】
	ctx.ViewData["current_month"] = Common.ConvertStringToInt(time.Now().Format(Common.DATE_FORMAT_M))
	//【○-1月】
	ctx.ViewData["past_month"] = Common.ConvertStringToInt(time.Now().AddDate(0, -1, 0).Format(Common.DATE_FORMAT_M))
	//【○年】
	ctx.ViewData["current_year"] = time.Now().Format(Common.DATE_FORMAT_Y)
	//【○ｰ1年】
	ctx.ViewData["past_year"] = time.Now().AddDate(-1, 0, 0).Format(Common.DATE_FORMAT_Y)
	//【○ｰ2年】
	ctx.ViewData["two_past_year"] = time.Now().AddDate(-2, 0, 0).Format(Common.DATE_FORMAT_Y)

	ctx.View = "report/052_shop_total_sum/search.html"
}
