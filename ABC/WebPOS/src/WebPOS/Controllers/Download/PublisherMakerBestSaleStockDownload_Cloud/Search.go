package PublisherMakerBestSaleStockDownload_Cloud

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"time"

	"github.com/goframework/gf"
)

func Search(ctx *gf.Context) {

	user := WebApp.GetContextUser(ctx)
	// 店舗リストを取得する
	sm := Models.ShopMasterModel{ctx.DB}
	listShop, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)

	ctx.ViewData["list_shop"] = listShop
	ctx.ViewData["option_flag_return"] = user.OptionFlgReturn
	ctx.ViewData["date"] = time.Now().AddDate(0, 0, -1).Format(Common.DATE_FORMAT_YMD_SLASH)

	ctx.ViewData["link"] = PATH_DOWNLOAD_MAKER_DATA_AJAX
	ctx.ViewData["link_download"] = PATH_DOWNLOAD_MAKER_DATA_REQUEST
	ctx.ViewData["type_search_sales"] = TYPE_SEARCH_SALES_TEXT
	ctx.ViewData["type_search_stock"] = TYPE_SEARCH_STOCK_TEXT
	ctx.ViewData["type_search_sales_and_return"] = TYPE_SEARCH_SALES_RETURN_TEXT
	ctx.ViewData["type_search_sales_and_receiving"] = TYPE_SEARCH_SALES_RECEIVING_TEXT

	ctx.ViewData["type_search_sales_value"] = TYPE_SEARCH_SALES
	ctx.ViewData["type_search_stock_value"] = TYPE_SEARCH_STOCK
	ctx.ViewData["type_search_sales_and_return_value"] = TYPE_SEARCH_SALES_RETURN
	ctx.ViewData["type_search_sales_and_receiving_value"] = TYPE_SEARCH_SALES_RECEIVING

	ctx.View = "download/publishermakerbestsalestock_cloud/search.html"
}
