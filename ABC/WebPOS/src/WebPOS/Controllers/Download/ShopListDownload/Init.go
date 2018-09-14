package ShopListDownload

import "github.com/goframework/gf"

const (
	PATH_DOWNLOAD_SHOP_LIST string = "/download/makersalestock/shop_list_download"
)

func Init() {

	gf.HandleGet(PATH_DOWNLOAD_SHOP_LIST, Download)
}
