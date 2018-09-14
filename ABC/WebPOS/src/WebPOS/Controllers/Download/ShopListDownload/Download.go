package ShopListDownload

import (
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"os"
)

func Download(ctx *gf.Context) {

	// ダウンロードトークン取得
	// ASO-5929 [BA]mBAWEB-v09f 店舗一覧ダウンロード-複数ファイル対応 - EDIT START
	//pathFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_PATH_FILE_SHOP_LIST_DOWNLOAD)
	fileName := ctx.Form.String("filename")
	pathFile := ctx.Config.StrOrEmpty(WebApp.CONFIG_PATH_FILE_SHOP_LIST_DOWNLOAD) + fileName
	// ASO-5929 [BA]mBAWEB-v09f 店舗一覧ダウンロード-複数ファイル対応 - EDIT START
	if _, err := os.Stat(pathFile); os.IsNotExist(err) {
		ctx.View = "download/shoplistdownload/result_error.html"
	 } else {
		ctx.ServeStaticFile(pathFile, true)
	}
}
