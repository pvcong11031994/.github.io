package PublisherMakerBestSaleStockDownload_Maria

import (
	"github.com/goframework/gf"
	"os"
)

const (
	DOWNLOAD_TOKEN_REQUEST = "f"
)

func Download(ctx *gf.Context) {

	// ダウンロードトークン取得
	downloadToken := ctx.Form.String(DOWNLOAD_TOKEN_REQUEST)
	file, ok := ctx.Session.Values[_MSS_DATA_DOWNLOAD_KEY_+downloadToken]
	if ok && file != nil {
		filePath := file.(string)
		defer os.Remove(filePath)
		ctx.ServeStaticFile(filePath, true)
	} else {
		ctx.Redirect(PATH_DOWNLOAD_MAKER_DATA_SEARCH)
	}
}
