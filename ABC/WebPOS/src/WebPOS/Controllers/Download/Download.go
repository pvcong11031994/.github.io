package Download

import (
	"WebPOS/Controllers/Download/MakerSaleStockDownload"
	"WebPOS/Controllers/Download/PublisherMakerBestSaleStockDownload"
	"WebPOS/Controllers/Download/PublisherMakerBestSaleStockDownload_Cloud"
	"WebPOS/Controllers/Download/PublisherMakerBestSaleStockDownload_Maria"
	"WebPOS/Controllers/Download/ShopListDownload"
)

func Init() {

	// メーカー売上・在庫ダウンロード_20170803
	MakerSaleStockDownload.Init()
	//店舗一覧ダウンロード
	ShopListDownload.Init()
	//出版社別ダウンロード
	PublisherMakerBestSaleStockDownload.Init()
	//出版社ダウンロード(Maria)
	PublisherMakerBestSaleStockDownload_Maria.Init()
	//出版社ダウンロード(Cloud)
	PublisherMakerBestSaleStockDownload_Cloud.Init()
}
