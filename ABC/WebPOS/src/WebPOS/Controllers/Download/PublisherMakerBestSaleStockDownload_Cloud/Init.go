package PublisherMakerBestSaleStockDownload_Cloud

import "github.com/goframework/gf"

const (
	PATH_DOWNLOAD_MAKER_DATA_SEARCH     string = "/download/makersalestock/bestsalestockdownload_cloud"
	PATH_DOWNLOAD_MAKER_DATA_AJAX       string = "/download/makersalestock/bestsalestockdownload_cloud/query_ajax"
	PATH_DOWNLOAD_MAKER_DATA_CHECK_AJAX string = "/download/makersalestock/bestsalestockdownload_cloud/check_query_ajax"
	PATH_DOWNLOAD_MAKER_DATA_REQUEST    string = "/download/makersalestock/bestsalestockdownload_cloud/query_download"
)

func Init() {

	gf.HandleGet(PATH_DOWNLOAD_MAKER_DATA_SEARCH, Search)
	gf.HandlePost(PATH_DOWNLOAD_MAKER_DATA_AJAX, Query)
	gf.HandlePost(PATH_DOWNLOAD_MAKER_DATA_CHECK_AJAX, Query_Check)
	gf.HandleGet(PATH_DOWNLOAD_MAKER_DATA_REQUEST, Download)
}
