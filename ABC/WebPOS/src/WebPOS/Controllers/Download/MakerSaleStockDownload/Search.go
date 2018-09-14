package MakerSaleStockDownload

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"time"
)

const (
	_LIMIT_CONCURRENT_DOWNLOAD_REQUEST = 10
)

var mPoolChan = make(chan bool, _LIMIT_CONCURRENT_DOWNLOAD_REQUEST)

func Search(ctx *gf.Context) {

	user := WebApp.GetContextUser(ctx)

	// 店舗リストを取得する
	sm := Models.ShopMasterModel{ctx.DB}
	listShop, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)

	// 出版社名を取得する
	mmm := Models.MakerMasterModel{ctx.DB}
	maker, err := mmm.GetPublisherInfoByUser(user.UserID)
	Common.LogErr(err)

	ctx.ViewData["list_shop"] = listShop
	ctx.ViewData["maker"] = maker.MakerName
	ctx.ViewData["date"] = time.Now().Format(Common.DATE_FORMAT_YMD_SLASH)

	ctx.ViewData["link"] = PATH_DOWNLOAD_MAKER_DATA_AJAX
	ctx.ViewData["link_download"] = PATH_DOWNLOAD_MAKER_DATA_REQUEST

	ctx.View = "download/makersalestock/search.html"
}
