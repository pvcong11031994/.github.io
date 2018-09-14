package MakerSaleStockDownload

import "github.com/goframework/gf"

const (
	PATH_DOWNLOAD_MAKER_DATA_SEARCH  string = "/download/makersalestock/search"
	PATH_DOWNLOAD_MAKER_DATA_AJAX    string = "/download/makersalestock/search/query_ajax"
	PATH_DOWNLOAD_MAKER_DATA_REQUEST string = "/download/makersalestock/search/query_download"
)

func Init() {

	gf.HandleGet(PATH_DOWNLOAD_MAKER_DATA_SEARCH, Search)
	gf.HandlePost(PATH_DOWNLOAD_MAKER_DATA_AJAX, Query)
	gf.HandleGet(PATH_DOWNLOAD_MAKER_DATA_REQUEST, Download)
}
