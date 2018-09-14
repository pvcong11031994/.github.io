package Report

import (
	"WebPOS/Controllers/Report/Download"
	"WebPOS/Controllers/Report/RP052_ShopTotalSum"
	"WebPOS/Controllers/Report/RP053_BestSalesByStore"
	"WebPOS/Controllers/Report/RP055_BestSales"
	"WebPOS/Controllers/Report/RP057_SingleGoods_Cumulative"
	"WebPOS/Controllers/Report/RP058_SalesComparison"
	"WebPOS/Controllers/Report/RP059_InitSalesCompare"
	"WebPOS/Controllers/Report/RP060_FavoriteManagement"
	"WebPOS/Controllers/Report/RP061_ShopSales"
	"WebPOS/Controllers/Report/RP062_SearchGoods"
	"WebPOS/Controllers/Report/RP063_SingleGoods_Stock_X"
	"WebPOS/Controllers/Report/RP064_BestSales_Maria"
	"WebPOS/Controllers/Report/RP065_ShopSales_Maria"
	"WebPOS/Controllers/Report/RP066_BestSales_Cloud"
	"WebPOS/Controllers/Report/RP067_ShopSales_Cloud"
)

func Init() {

	//売上累計検索
	RP052_ShopTotalSum.Init()

	//店舗別売上ベスト
	RP053_BestSalesByStore.Init()

	//売上ベスト
	RP055_BestSales.Init()

	//単品推移（累計推移）
	RP057_SingleGoods_Cumulative.Init()

	//売上比較
	RP058_SalesComparison.Init()

	// 初速比較
	RP059_InitSalesCompare.Init()

	// お気に入り管理
	RP060_FavoriteManagement.Init()

	// 店舗別集計に変更
	RP061_ShopSales.Init()

	//商品検索
	RP062_SearchGoods.Init()

	// 単品推移(在庫参照)
	RP063_SingleGoods_Stock_X.Init()

	//売上ベストの累計と在庫
	RP064_BestSales_Maria.Init()

	// 店舗別集計に変更
	RP065_ShopSales_Maria.Init()

	//売上ベストCloudSQL版
	RP066_BestSales_Cloud.Init()

	// 店舗別集計(Cloud)
	RP067_ShopSales_Cloud.Init()

	// ダウンロード
	Download.Init()
}
