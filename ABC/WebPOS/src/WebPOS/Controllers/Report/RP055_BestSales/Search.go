package RP055_BestSales

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

	ctx.ViewData["form"] = RPComon.QueryFormSingleGoods{
		// init form
		DateFrom:         time.Now().AddDate(0, -1, 0).Format(Common.DATE_FORMAT_YMD_SLASH),
		DateTo:           time.Now().Format(Common.DATE_FORMAT_YMD_SLASH),
		GroupType:        "0",
		Limit:            100,
		DownloadType:     "0",
		ControlType:      "1",
		Page:             1,
		MagazineCodeWeek: "1",
	}

	// load form revert
	keyForm := ctx.Form.String("key_form")
	if keyForm != "" {
		formObj := RPComon.QueryFormSingleGoods{}
		ctx.LoadCache(keyForm, &formObj)
		ctx.ViewData["form"] = formObj
	}

	ctx.ViewData["default_from"] = time.Now().AddDate(0, -1, 0).Format(Common.DATE_FORMAT_YMD_SLASH)
	ctx.ViewData["default_to"] = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)

	user := WebApp.GetContextUser(ctx)
	// 帳票項目
	stringsCol, stringsRow, stringsSum := initDefaultLayout()
	ctx.ViewData["layout_item_col"] = stringsCol
	ctx.ViewData["layout_item_row"] = stringsRow
	ctx.ViewData["layout_item_sum"] = stringsSum

	sm := Models.ShopMasterModel{ctx.DB}
	listShop, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)

	// メディア分類
	bqctm := Models.BQCategoryModel{ctx.DB}
	listMedia2, err := bqctm.ListFullGenreByJanGroup([]string{"1", "2"})
	Common.LogErr(err)

	ctx.ViewData["list_shop"] = listShop
	ctx.ViewData["list_media2"] = listMedia2
	ctx.TemplateFunc["plus"] = Common.Plus

	// init date
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

	ctx.ViewData["link_revert"] = PATH_REPORT_BEST_SALES_SEARCH
	ctx.View = "report/055_best_sales/search.html"
}
