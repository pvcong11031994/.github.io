package RP053_BestSalesByStore

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"strings"
	"time"
)

//売上ベスト検索
func Search(ctx *gf.Context) {

	user := WebApp.GetContextUser(ctx)
	// 帳票項目
	srlm := Models.SettingReportLayoutModel{ctx.DB}
	var layout *ModelItems.SettingReportLayoutItem
	var err error
	menuId := ctx.Form.Int("menu")
	if menuId != 0 {
		layout, err = srlm.GetSettingReportLayoutByMenu(user.UserID, _REPORT_ID, menuId)
		ctx.ViewData["sub_report_name"] = layout.ReportName
		Common.LogErr(err)
	} else {
		layout, err = srlm.GetSettingReportLayout(user.UserID, _REPORT_ID)
		Common.LogErr(err)
	}

	if layout.SelectedCol[0] == "" &&
		layout.SelectedRow[0] == "" &&
		layout.SelectedSum[0] == "" {
		layout.SelectedCol[0] = "mc_yyyy,mc_mm,mc_weekdate,mc_dd"
		layout.SelectedRow[0] = "rank_no,bqio_jan_cd,goods_name,writer_name,publisher_name," +
			"sales_date,bqgm_price,total_sales,total_arrival,stock_count"
		layout.SelectedSum[0] = "bqio_goods_count"
	}

	stringsCol, stringsRow, stringsSum := initDefaultLayout()
	ctx.ViewData["layout_item_col"] = stringsCol
	ctx.ViewData["layout_item_row"] = stringsRow
	ctx.ViewData["layout_item_sum"] = stringsSum

	ctx.ViewData["layout_item_col_selected"] = strings.Join(layout.SelectedCol, ",")
	ctx.ViewData["layout_item_row_selected"] = strings.Join(layout.SelectedRow, ",")
	ctx.ViewData["layout_item_sum_selected"] = strings.Join(layout.SelectedSum, ",")

	sm := Models.ShopMasterModel{ctx.DB}
	listShop, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)

	// メディア分類
	bqctm := Models.BQCategoryModel{ctx.DB}
	listMedia, err := bqctm.GetMediaListByUser(user.UserID)
	Common.LogErr(err)

	ctx.ViewData["list_shop"] = listShop
	ctx.ViewData["list_media"] = listMedia

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

	ctx.View = "report/053_best_sales_by_store/search.html"
}
